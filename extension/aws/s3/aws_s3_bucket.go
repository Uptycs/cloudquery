package s3

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Uptycs/cloudquery/utilities"
	"log"
	"os"
	"strings"

	extaws "github.com/Uptycs/cloudquery/extension/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/kolide/osquery-go/plugin/table"
)

type s3BucketInfo struct {
	Name                              string
	CreationTime                      string
	ServerSideEncryptionConfiguration *s3.ServerSideEncryptionConfiguration
	MfaDelete                         *string
	VersioningStatus                  *string
	AclOwner                          *s3.Owner
	AclGrants                         []*s3.Grant
	WebsiteEnabled                    bool
	WebsiteRedirection                *s3.RedirectAllRequestsTo
	PublicAccessBlockConfig           *s3.PublicAccessBlockConfiguration
	PolicyStatus                      *s3.PolicyStatus
	AccelerateConfigurationStatus     *string
	ObjectLockConfigurationEnabled    bool
	LifecycleConfigurationEnabled     bool
	NotificationEnabled               bool
	CorsEnabled                       bool
	Policy                            *string
	Tags                              []*s3.Tag
}

type s3BucketInfoList struct {
	buckets []s3BucketInfo
}

var (
	// Map of region to buckets
	regionBuckets map[string]s3BucketInfoList
)

func ListBucketsColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("account_id"),
		table.TextColumn("region_code"),
		table.TextColumn("region"),
		table.TextColumn("name"),
		table.TextColumn("creation_time"),
		table.TextColumn("server_side_encryption_configuration"),
		table.TextColumn("mfa_delete"),
		table.TextColumn("versioning_status"),
		table.TextColumn("acl_owner"),
		table.TextColumn("acl_grants"),
		table.TextColumn("website_enabled"),
		table.TextColumn("website_redirection"),
		table.TextColumn("public_access_block_config"),
		table.TextColumn("policy_status"),
		table.TextColumn("accelerate_configuration_status"),
		table.TextColumn("object_lock_configuration_enabled"),
		table.TextColumn("lifecycle_configuration_enabled"),
		table.TextColumn("notification_enabled"),
		table.TextColumn("cors_enabled"),
		table.TextColumn("policy"),
		table.TextColumn("tags"),
	}
}

func ListBucketsGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	if len(utilities.ExtConfiguration.ExtConfAws.Accounts) == 0 {
		fmt.Println("Processing default account")
		results, err := processAccountListBuckets(nil)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, results...)
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfAws.Accounts {
			fmt.Println("Processing account:" + account.ID)
			results, err := processAccountListBuckets(&account)
			if err != nil {
				// TODO: Continue to next account or return error ?
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}

	return resultMap, nil
}

func getBucketLocation(svc *s3.S3, bucketName *string) (string, error) {
	bucketLocationInput := s3.GetBucketLocationInput{Bucket: bucketName}
	getBucketLocationOutput, err := svc.GetBucketLocation(&bucketLocationInput)
	if err != nil {
		fmt.Println("Failed to get bucket location ", err)
		log.Fatal(err)
		return "", err
	}
	if getBucketLocationOutput.LocationConstraint == nil || len(*getBucketLocationOutput.LocationConstraint) == 0 {
		// Default is us-east-1
		return "us-east-1", nil
	} else if strings.EqualFold(*getBucketLocationOutput.LocationConstraint, "EU") {
		return "us-west-1", nil
	} else {
		return *getBucketLocationOutput.LocationConstraint, nil
	}
}

func addBucketToRegionBucketList(svc *s3.S3, bucket *s3.Bucket) error {
	bucketRegion, err := getBucketLocation(svc, bucket.Name)
	if err != nil {
		return err
	}
	bucketInfo := s3BucketInfo{
		Name:         *bucket.Name,
		CreationTime: bucket.CreationDate.String(),
	}
	bucketList, ok := regionBuckets[bucketRegion]
	if !ok {
		bucketList = s3BucketInfoList{buckets: make([]s3BucketInfo, 0)}
		regionBuckets[bucketRegion] = bucketList
	}
	bucketList.buckets = append(bucketList.buckets, bucketInfo)
	regionBuckets[bucketRegion] = bucketList
	return nil
}

func (bucket *s3BucketInfo) getBucketEncryption(svc *s3.S3) {
	input := s3.GetBucketEncryptionInput{Bucket: &bucket.Name}
	output, err := svc.GetBucketEncryption(&input)
	if err != nil {
		// fmt.Println("Failed to get bucket encryption ", err)
		return
	}
	bucket.ServerSideEncryptionConfiguration = output.ServerSideEncryptionConfiguration
}

func (bucket *s3BucketInfo) getBucketVersioning(svc *s3.S3) {
	input := s3.GetBucketVersioningInput{Bucket: &bucket.Name}
	output, err := svc.GetBucketVersioning(&input)
	if err != nil {
		// fmt.Println("Failed to get bucket versioning ", err)
		return
	}
	bucket.MfaDelete = output.MFADelete
	bucket.VersioningStatus = output.Status
}

func (bucket *s3BucketInfo) getBucketAcl(svc *s3.S3) {
	input := s3.GetBucketAclInput{Bucket: &bucket.Name}
	output, err := svc.GetBucketAcl(&input)
	if err != nil {
		// fmt.Println("Failed to get bucket ACL ", err)
		return
	}
	bucket.AclOwner = output.Owner
	bucket.AclGrants = output.Grants
}

func (bucket *s3BucketInfo) getBucketWebsite(svc *s3.S3) {
	bucket.WebsiteEnabled = false
	input := s3.GetBucketWebsiteInput{Bucket: &bucket.Name}
	output, err := svc.GetBucketWebsite(&input)
	if err != nil {
		// fmt.Println("Failed to get bucket website ", err)
		return
	}
	if output != nil {
		bucket.WebsiteEnabled = true
	}
	bucket.WebsiteRedirection = output.RedirectAllRequestsTo
}

func (bucket *s3BucketInfo) getBucketPublicAccessBlock(svc *s3.S3) {
	input := s3.GetPublicAccessBlockInput{Bucket: &bucket.Name}
	output, err := svc.GetPublicAccessBlock(&input)
	if err != nil {
		// fmt.Println("Failed to get bucket public access block ", err)
		return
	}
	bucket.PublicAccessBlockConfig = output.PublicAccessBlockConfiguration
}

func (bucket *s3BucketInfo) getBucketPolicyStatus(svc *s3.S3) {
	input := s3.GetBucketPolicyStatusInput{Bucket: &bucket.Name}
	output, err := svc.GetBucketPolicyStatus(&input)
	if err != nil {
		// fmt.Println("Failed to get bucket policy status ", err)
		return
	}
	bucket.PolicyStatus = output.PolicyStatus
}

func (bucket *s3BucketInfo) getBucketAccelerateConfiguration(svc *s3.S3) {
	input := s3.GetBucketAccelerateConfigurationInput{Bucket: &bucket.Name}
	output, err := svc.GetBucketAccelerateConfiguration(&input)
	if err != nil {
		// fmt.Println("Failed to get bucket accelerate configuration ", err)
		return
	}
	bucket.AccelerateConfigurationStatus = output.Status
}

func (bucket *s3BucketInfo) getObjectLockConfiguration(svc *s3.S3) {
	bucket.ObjectLockConfigurationEnabled = false
	input := s3.GetObjectLockConfigurationInput{Bucket: &bucket.Name}
	output, err := svc.GetObjectLockConfiguration(&input)
	if err != nil {
		// fmt.Println("Failed to get object lock configuration ", err)
		return
	}
	if output != nil && output.ObjectLockConfiguration != nil {
		bucket.ObjectLockConfigurationEnabled = true
	}
}

func (bucket *s3BucketInfo) getBucketLifecycleConfiguration(svc *s3.S3) {
	bucket.LifecycleConfigurationEnabled = false
	input := s3.GetBucketLifecycleConfigurationInput{Bucket: &bucket.Name}
	output, err := svc.GetBucketLifecycleConfiguration(&input)
	if err != nil {
		// fmt.Println("Failed to get bucket lifecycle configuration ", err)
		return
	}
	if output != nil && len(output.Rules) > 0 {
		bucket.LifecycleConfigurationEnabled = true
	}
}

func (bucket *s3BucketInfo) getBucketTags(svc *s3.S3) {
	input := s3.GetBucketTaggingInput{Bucket: &bucket.Name}
	output, err := svc.GetBucketTagging(&input)
	if err != nil {
		// fmt.Println("Failed to get bucket tags ", err)
		return
	}
	bucket.Tags = output.TagSet
}

func (bucket *s3BucketInfo) getBucketNotificationConfiguration(svc *s3.S3) {
	bucket.NotificationEnabled = false
	input := s3.GetBucketNotificationConfigurationRequest{Bucket: &bucket.Name}
	output, err := svc.GetBucketNotificationConfiguration(&input)
	if err != nil {
		// fmt.Println("Failed to get bucket notification configuration ", err)
		return
	}
	if output != nil {
		bucket.NotificationEnabled = true
	}
}

func (bucket *s3BucketInfo) getBucketCorsConfiguration(svc *s3.S3) {
	bucket.CorsEnabled = false
	input := s3.GetBucketCorsInput{Bucket: &bucket.Name}
	output, err := svc.GetBucketCors(&input)
	if err != nil {
		// fmt.Println("Failed to get bucket cors configuration ", err)
		return
	}
	if output != nil && len(output.CORSRules) > 0 {
		bucket.CorsEnabled = true
	}
}

func (bucket *s3BucketInfo) getBucketPolicy(svc *s3.S3) {
	input := s3.GetBucketPolicyInput{Bucket: &bucket.Name}
	output, err := svc.GetBucketPolicy(&input)
	if err != nil {
		// fmt.Println("Failed to get bucket policy ", err)
		return
	}
	bucket.Policy = output.Policy
}

func processBucket(tableConfig *utilities.TableConfig, account *utilities.ExtensionConfigurationAwsAccount, region string, bucket *s3BucketInfo) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	sess, err := extaws.GetAwsSession(account, region)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	accountId := utilities.AwsAccountId
	if account != nil {
		accountId = account.ID
	}
	svc := s3.New(sess)
	fmt.Println("Processing bucket:" + bucket.Name)
	bucket.getBucketAccelerateConfiguration(svc)
	bucket.getBucketAcl(svc)
	bucket.getBucketCorsConfiguration(svc)
	bucket.getBucketEncryption(svc)
	bucket.getBucketLifecycleConfiguration(svc)
	bucket.getBucketNotificationConfiguration(svc)
	bucket.getBucketPolicy(svc)
	bucket.getBucketPolicyStatus(svc)
	bucket.getBucketPublicAccessBlock(svc)
	bucket.getBucketTags(svc)
	bucket.getBucketVersioning(svc)
	bucket.getBucketWebsite(svc)
	bucket.getObjectLockConfiguration(svc)
	//byteArr, _ := json.MarshalIndent(bucket, "", "  ")
	//fmt.Printf("%s\n", string(byteArr))
	byteArr, err := json.Marshal(bucket)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return resultMap, err
	}
	table := utilities.Table{}
	table.Init(byteArr, tableConfig.MaxLevel, tableConfig.GetParsedAttributeConfigMap())
	for _, row := range table.Rows {
		result := extaws.RowToMap(row, accountId, region, tableConfig)
		resultMap = append(resultMap, result)
	}
	fmt.Printf("Map size:%d\n", len(resultMap))
	return resultMap, nil
}

func processListBuckets(tableConfig *utilities.TableConfig, account *utilities.ExtensionConfigurationAwsAccount) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	sess, err := extaws.GetAwsSession(account, "us-west-1")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	svc := s3.New(sess)
	params := &s3.ListBucketsInput{}

	// Get list of buckets
	output, err := svc.ListBuckets(params)
	if err != nil {
		fmt.Println("ListBuckets.Page: ", err)
		log.Fatal(err)
		return resultMap, err
	}
	regionBuckets = make(map[string]s3BucketInfoList)
	// Get bucket region and put that bucket in that bucketList
	for _, bucket := range output.Buckets {
		addBucketToRegionBucketList(svc, bucket)
	}
	// Process all buckets
	for region, regionBucketList := range regionBuckets {
		for _, regionBucket := range regionBucketList.buckets {
			result, err := processBucket(tableConfig, account, region, &regionBucket)
			if err == nil {
				resultMap = append(resultMap, result...)
			}
		}
	}
	return resultMap, nil
}

func processAccountListBuckets(account *utilities.ExtensionConfigurationAwsAccount) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	tableConfig, ok := utilities.TableConfigurationMap["aws_s3_bucket"]
	if !ok {
		fmt.Println("failed to get TableConfig")
		return resultMap, fmt.Errorf("table configuration not found")
	}
	result, err := processListBuckets(tableConfig, account)
	if err != nil {
		log.Fatal(err)
		return resultMap, err
	}
	resultMap = append(resultMap, result...)
	return resultMap, nil
}
