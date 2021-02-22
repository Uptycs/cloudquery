/**
 * Copyright (c) 2020-present, The cloudquery authors
 *
 * This source code is licensed as defined by the LICENSE file found in the
 * root directory of this source tree.
 *
 * SPDX-License-Identifier: (Apache-2.0 OR GPL-2.0-only)
 */

package s3

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Uptycs/cloudquery/utilities"

	log "github.com/sirupsen/logrus"

	"github.com/Uptycs/basequery-go/plugin/table"
	extaws "github.com/Uptycs/cloudquery/extension/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type s3BucketInfo struct {
	Name                              string
	CreationTime                      string
	ServerSideEncryptionConfiguration *types.ServerSideEncryptionConfiguration
	MfaDelete                         string
	VersioningStatus                  string
	AclOwner                          *types.Owner
	AclGrants                         []types.Grant
	WebsiteEnabled                    bool
	WebsiteRedirection                *types.RedirectAllRequestsTo
	PublicAccessBlockConfig           *types.PublicAccessBlockConfiguration
	PolicyStatus                      *types.PolicyStatus
	AccelerateConfigurationStatus     string
	ObjectLockConfigurationEnabled    bool
	LifecycleConfigurationEnabled     bool
	NotificationEnabled               bool
	CorsEnabled                       bool
	Policy                            *string
	Tags                              []types.Tag
}

type s3BucketInfoList struct {
	buckets []s3BucketInfo
}

var (
	// Map of region to buckets
	regionBuckets map[string]s3BucketInfoList
)

// ListBucketsColumns returns the list of columns in the table
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

// ListBucketsGenerate returns the rows in the table for all configured accounts
func ListBucketsGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	if len(utilities.ExtConfiguration.ExtConfAws.Accounts) == 0 {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_s3_bucket",
			"account":   "default",
		}).Info("processing account")
		results, err := processAccountListBuckets(nil)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, results...)
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfAws.Accounts {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_s3_bucket",
				"account":   account.ID,
			}).Info("processing account")
			results, err := processAccountListBuckets(&account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}

	return resultMap, nil
}

func getBucketLocation(svc *s3.Client, bucketName *string) (string, error) {
	bucketLocationInput := s3.GetBucketLocationInput{Bucket: bucketName}
	getBucketLocationOutput, err := svc.GetBucketLocation(context.TODO(), &bucketLocationInput)
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_s3_bucket",
			"bucket":    *bucketName,
			"errString": err.Error(),
		}).Error("failed to get bucket location")
		return "", err
	}
	if len(getBucketLocationOutput.LocationConstraint) == 0 {
		// Default is us-east-1
		return "us-east-1", nil
	} else if getBucketLocationOutput.LocationConstraint == types.BucketLocationConstraintEu {
		return "us-west-1", nil
	} else {
		return string(getBucketLocationOutput.LocationConstraint), nil
	}
}

func addBucketToRegionBucketList(svc *s3.Client, bucket types.Bucket) error {
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

func (bucket *s3BucketInfo) getBucketEncryption(svc *s3.Client) {
	input := s3.GetBucketEncryptionInput{Bucket: &bucket.Name}
	output, err := svc.GetBucketEncryption(context.TODO(), &input)
	if err != nil {
		return
	}
	bucket.ServerSideEncryptionConfiguration = output.ServerSideEncryptionConfiguration
}

func (bucket *s3BucketInfo) getBucketVersioning(svc *s3.Client) {
	input := s3.GetBucketVersioningInput{Bucket: &bucket.Name}
	output, err := svc.GetBucketVersioning(context.TODO(), &input)
	if err != nil {
		return
	}
	bucket.MfaDelete = string(output.MFADelete)
	bucket.VersioningStatus = string(output.Status)
}

func (bucket *s3BucketInfo) getBucketAcl(svc *s3.Client) {
	input := s3.GetBucketAclInput{Bucket: &bucket.Name}
	output, err := svc.GetBucketAcl(context.TODO(), &input)
	if err != nil {
		return
	}
	bucket.AclOwner = output.Owner
	bucket.AclGrants = output.Grants
}

func (bucket *s3BucketInfo) getBucketWebsite(svc *s3.Client) {
	bucket.WebsiteEnabled = false
	input := s3.GetBucketWebsiteInput{Bucket: &bucket.Name}
	output, err := svc.GetBucketWebsite(context.TODO(), &input)
	if err != nil {
		return
	}
	if output != nil {
		bucket.WebsiteEnabled = true
	}
	bucket.WebsiteRedirection = output.RedirectAllRequestsTo
}

func (bucket *s3BucketInfo) getBucketPublicAccessBlock(svc *s3.Client) {
	input := s3.GetPublicAccessBlockInput{Bucket: &bucket.Name}
	output, err := svc.GetPublicAccessBlock(context.TODO(), &input)
	if err != nil {
		return
	}
	bucket.PublicAccessBlockConfig = output.PublicAccessBlockConfiguration
}

func (bucket *s3BucketInfo) getBucketPolicyStatus(svc *s3.Client) {
	input := s3.GetBucketPolicyStatusInput{Bucket: &bucket.Name}
	output, err := svc.GetBucketPolicyStatus(context.TODO(), &input)
	if err != nil {
		return
	}
	bucket.PolicyStatus = output.PolicyStatus
}

func (bucket *s3BucketInfo) getBucketAccelerateConfiguration(svc *s3.Client) {
	input := s3.GetBucketAccelerateConfigurationInput{Bucket: &bucket.Name}
	output, err := svc.GetBucketAccelerateConfiguration(context.TODO(), &input)
	if err != nil {
		return
	}
	bucket.AccelerateConfigurationStatus = string(output.Status)
}

func (bucket *s3BucketInfo) getObjectLockConfiguration(svc *s3.Client) {
	bucket.ObjectLockConfigurationEnabled = false
	input := s3.GetObjectLockConfigurationInput{Bucket: &bucket.Name}
	output, err := svc.GetObjectLockConfiguration(context.TODO(), &input)
	if err != nil {
		return
	}
	if output != nil && output.ObjectLockConfiguration != nil {
		bucket.ObjectLockConfigurationEnabled = true
	}
}

func (bucket *s3BucketInfo) getBucketLifecycleConfiguration(svc *s3.Client) {
	bucket.LifecycleConfigurationEnabled = false
	input := s3.GetBucketLifecycleConfigurationInput{Bucket: &bucket.Name}
	output, err := svc.GetBucketLifecycleConfiguration(context.TODO(), &input)
	if err != nil {
		return
	}
	if output != nil && len(output.Rules) > 0 {
		bucket.LifecycleConfigurationEnabled = true
	}
}

func (bucket *s3BucketInfo) getBucketTags(svc *s3.Client) {
	input := s3.GetBucketTaggingInput{Bucket: &bucket.Name}
	output, err := svc.GetBucketTagging(context.TODO(), &input)
	if err != nil {
		return
	}
	bucket.Tags = output.TagSet
}

func (bucket *s3BucketInfo) getBucketNotificationConfiguration(svc *s3.Client) {
	bucket.NotificationEnabled = false
	input := s3.GetBucketNotificationConfigurationInput{Bucket: &bucket.Name}
	output, err := svc.GetBucketNotificationConfiguration(context.TODO(), &input)
	if err != nil {
		return
	}
	if output != nil {
		bucket.NotificationEnabled = true
	}
}

func (bucket *s3BucketInfo) getBucketCorsConfiguration(svc *s3.Client) {
	bucket.CorsEnabled = false
	input := s3.GetBucketCorsInput{Bucket: &bucket.Name}
	output, err := svc.GetBucketCors(context.TODO(), &input)
	if err != nil {
		return
	}
	if output != nil && len(output.CORSRules) > 0 {
		bucket.CorsEnabled = true
	}
}

func (bucket *s3BucketInfo) getBucketPolicy(svc *s3.Client) {
	input := s3.GetBucketPolicyInput{Bucket: &bucket.Name}
	output, err := svc.GetBucketPolicy(context.TODO(), &input)
	if err != nil {
		return
	}
	bucket.Policy = output.Policy
}

func processBucket(tableConfig *utilities.TableConfig, account *utilities.ExtensionConfigurationAwsAccount, region string, bucket *s3BucketInfo) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	sess, err := extaws.GetAwsConfig(account, region)
	if err != nil {
		return resultMap, err
	}
	accountId := utilities.AwsAccountID
	if account != nil {
		accountId = account.ID
	}
	svc := s3.NewFromConfig(*sess)
	utilities.GetLogger().WithFields(log.Fields{
		"tableName": "aws_s3_bucket",
		"bucket":    bucket.Name,
	}).Debug("processing bucket")
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
	byteArr, err := json.Marshal(bucket)
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_s3_bucket",
			"account":   accountId,
			"region":    region,
			"bucket":    bucket.Name,
			"errString": err.Error(),
		}).Error("failed to marshal response")
		return resultMap, err
	}
	table := utilities.NewTable(byteArr, tableConfig)
	for _, row := range table.Rows {
		result := extaws.RowToMap(row, accountId, region, tableConfig)
		resultMap = append(resultMap, result)
	}
	return resultMap, nil
}

func processListBuckets(tableConfig *utilities.TableConfig, account *utilities.ExtensionConfigurationAwsAccount) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	sess, err := extaws.GetAwsConfig(account, "us-west-1")
	if err != nil {
		return resultMap, err
	}

	svc := s3.NewFromConfig(*sess)
	params := &s3.ListBucketsInput{}

	accountId := utilities.AwsAccountID
	if account != nil {
		accountId = account.ID
	}

	// Get list of buckets
	output, err := svc.ListBuckets(context.TODO(), params)
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_s3_bucket",
			"account":   accountId,
			"errString": err.Error(),
		}).Error("failed to get bucket list")
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
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_s3_bucket",
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found")
	}
	result, err := processListBuckets(tableConfig, account)
	if err != nil {
		return resultMap, err
	}
	resultMap = append(resultMap, result...)
	return resultMap, nil
}
