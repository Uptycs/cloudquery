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
	if len(utilities.ExtConfiguration.ExtConfAws.Accounts) == 0 && extaws.ShouldProcessAccount("aws_s3_bucket", utilities.AwsAccountID) {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_s3_bucket",
			"account":   "default",
		}).Info("processing account")
		results, err := processAccountListBuckets(osqCtx, queryContext, nil)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, results...)
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfAws.Accounts {
			if !extaws.ShouldProcessAccount("aws_s3_bucket", account.ID) {
				continue
			}
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_s3_bucket",
				"account":   account.ID,
			}).Info("processing account")
			results, err := processAccountListBuckets(osqCtx, queryContext, &account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}

	return resultMap, nil
}

func getBucketLocation(osqCtx context.Context, queryContext table.QueryContext, svc *s3.Client, bucketName *string) (string, error) {
	bucketLocationInput := s3.GetBucketLocationInput{Bucket: bucketName}
	getBucketLocationOutput, err := svc.GetBucketLocation(osqCtx, &bucketLocationInput)
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

func addBucketToRegionBucketList(osqCtx context.Context, queryContext table.QueryContext, svc *s3.Client, bucket types.Bucket) error {
	bucketRegion, err := getBucketLocation(osqCtx, queryContext, svc, bucket.Name)
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

func (bucket *s3BucketInfo) getBucketEncryption(osqCtx context.Context, queryContext table.QueryContext, svc *s3.Client) {
	input := s3.GetBucketEncryptionInput{Bucket: &bucket.Name}
	output, err := svc.GetBucketEncryption(osqCtx, &input)
	if err != nil {
		return
	}
	bucket.ServerSideEncryptionConfiguration = output.ServerSideEncryptionConfiguration
}

func (bucket *s3BucketInfo) getBucketVersioning(osqCtx context.Context, queryContext table.QueryContext, svc *s3.Client) {
	input := s3.GetBucketVersioningInput{Bucket: &bucket.Name}
	output, err := svc.GetBucketVersioning(osqCtx, &input)
	if err != nil {
		return
	}
	bucket.MfaDelete = string(output.MFADelete)
	bucket.VersioningStatus = string(output.Status)
}

func (bucket *s3BucketInfo) getBucketAcl(osqCtx context.Context, queryContext table.QueryContext, svc *s3.Client) {
	input := s3.GetBucketAclInput{Bucket: &bucket.Name}
	output, err := svc.GetBucketAcl(osqCtx, &input)
	if err != nil {
		return
	}
	bucket.AclOwner = output.Owner
	bucket.AclGrants = output.Grants
}

func (bucket *s3BucketInfo) getBucketWebsite(osqCtx context.Context, queryContext table.QueryContext, svc *s3.Client) {
	bucket.WebsiteEnabled = false
	input := s3.GetBucketWebsiteInput{Bucket: &bucket.Name}
	output, err := svc.GetBucketWebsite(osqCtx, &input)
	if err != nil {
		return
	}
	if output != nil {
		bucket.WebsiteEnabled = true
	}
	bucket.WebsiteRedirection = output.RedirectAllRequestsTo
}

func (bucket *s3BucketInfo) getBucketPublicAccessBlock(osqCtx context.Context, queryContext table.QueryContext, svc *s3.Client) {
	input := s3.GetPublicAccessBlockInput{Bucket: &bucket.Name}
	output, err := svc.GetPublicAccessBlock(osqCtx, &input)
	if err != nil {
		return
	}
	bucket.PublicAccessBlockConfig = output.PublicAccessBlockConfiguration
}

func (bucket *s3BucketInfo) getBucketPolicyStatus(osqCtx context.Context, queryContext table.QueryContext, svc *s3.Client) {
	input := s3.GetBucketPolicyStatusInput{Bucket: &bucket.Name}
	output, err := svc.GetBucketPolicyStatus(osqCtx, &input)
	if err != nil {
		return
	}
	bucket.PolicyStatus = output.PolicyStatus
}

func (bucket *s3BucketInfo) getBucketAccelerateConfiguration(osqCtx context.Context, queryContext table.QueryContext, svc *s3.Client) {
	input := s3.GetBucketAccelerateConfigurationInput{Bucket: &bucket.Name}
	output, err := svc.GetBucketAccelerateConfiguration(osqCtx, &input)
	if err != nil {
		return
	}
	bucket.AccelerateConfigurationStatus = string(output.Status)
}

func (bucket *s3BucketInfo) getObjectLockConfiguration(osqCtx context.Context, queryContext table.QueryContext, svc *s3.Client) {
	bucket.ObjectLockConfigurationEnabled = false
	input := s3.GetObjectLockConfigurationInput{Bucket: &bucket.Name}
	output, err := svc.GetObjectLockConfiguration(osqCtx, &input)
	if err != nil {
		return
	}
	if output != nil && output.ObjectLockConfiguration != nil {
		bucket.ObjectLockConfigurationEnabled = true
	}
}

func (bucket *s3BucketInfo) getBucketLifecycleConfiguration(osqCtx context.Context, queryContext table.QueryContext, svc *s3.Client) {
	bucket.LifecycleConfigurationEnabled = false
	input := s3.GetBucketLifecycleConfigurationInput{Bucket: &bucket.Name}
	output, err := svc.GetBucketLifecycleConfiguration(osqCtx, &input)
	if err != nil {
		return
	}
	if output != nil && len(output.Rules) > 0 {
		bucket.LifecycleConfigurationEnabled = true
	}
}

func (bucket *s3BucketInfo) getBucketTags(osqCtx context.Context, queryContext table.QueryContext, svc *s3.Client) {
	input := s3.GetBucketTaggingInput{Bucket: &bucket.Name}
	output, err := svc.GetBucketTagging(osqCtx, &input)
	if err != nil {
		return
	}
	bucket.Tags = output.TagSet
}

func (bucket *s3BucketInfo) getBucketNotificationConfiguration(osqCtx context.Context, queryContext table.QueryContext, svc *s3.Client) {
	bucket.NotificationEnabled = false
	input := s3.GetBucketNotificationConfigurationInput{Bucket: &bucket.Name}
	output, err := svc.GetBucketNotificationConfiguration(osqCtx, &input)
	if err != nil {
		return
	}
	if output != nil {
		bucket.NotificationEnabled = true
	}
}

func (bucket *s3BucketInfo) getBucketCorsConfiguration(osqCtx context.Context, queryContext table.QueryContext, svc *s3.Client) {
	bucket.CorsEnabled = false
	input := s3.GetBucketCorsInput{Bucket: &bucket.Name}
	output, err := svc.GetBucketCors(osqCtx, &input)
	if err != nil {
		return
	}
	if output != nil && len(output.CORSRules) > 0 {
		bucket.CorsEnabled = true
	}
}

func (bucket *s3BucketInfo) getBucketPolicy(osqCtx context.Context, queryContext table.QueryContext, svc *s3.Client) {
	input := s3.GetBucketPolicyInput{Bucket: &bucket.Name}
	output, err := svc.GetBucketPolicy(osqCtx, &input)
	if err != nil {
		return
	}
	bucket.Policy = output.Policy
}

func processBucket(osqCtx context.Context, queryContext table.QueryContext, tableConfig *utilities.TableConfig, account *utilities.ExtensionConfigurationAwsAccount, region string, bucket *s3BucketInfo) ([]map[string]string, error) {
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
	bucket.getBucketAccelerateConfiguration(osqCtx, queryContext, svc)
	bucket.getBucketAcl(osqCtx, queryContext, svc)
	bucket.getBucketCorsConfiguration(osqCtx, queryContext, svc)
	bucket.getBucketEncryption(osqCtx, queryContext, svc)
	bucket.getBucketLifecycleConfiguration(osqCtx, queryContext, svc)
	bucket.getBucketNotificationConfiguration(osqCtx, queryContext, svc)
	bucket.getBucketPolicy(osqCtx, queryContext, svc)
	bucket.getBucketPolicyStatus(osqCtx, queryContext, svc)
	bucket.getBucketPublicAccessBlock(osqCtx, queryContext, svc)
	bucket.getBucketTags(osqCtx, queryContext, svc)
	bucket.getBucketVersioning(osqCtx, queryContext, svc)
	bucket.getBucketWebsite(osqCtx, queryContext, svc)
	bucket.getObjectLockConfiguration(osqCtx, queryContext, svc)
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
		if !extaws.ShouldProcessRow(osqCtx, queryContext, "aws_s3_bucket", accountId, region, row) {
			continue
		}
		result := extaws.RowToMap(row, accountId, region, tableConfig)
		resultMap = append(resultMap, result)
	}
	return resultMap, nil
}

func processListBuckets(osqCtx context.Context, queryContext table.QueryContext, tableConfig *utilities.TableConfig, account *utilities.ExtensionConfigurationAwsAccount) ([]map[string]string, error) {
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
	output, err := svc.ListBuckets(osqCtx, params)
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
		addBucketToRegionBucketList(osqCtx, queryContext, svc, bucket)
	}
	// Process all buckets
	for region, regionBucketList := range regionBuckets {
		accountId := utilities.AwsAccountID
		if account != nil {
			accountId = account.ID
		}
		if !extaws.ShouldProcessRegion("aws_s3_bucket", accountId, region) {
			continue
		}
		for _, regionBucket := range regionBucketList.buckets {
			result, err := processBucket(osqCtx, queryContext, tableConfig, account, region, &regionBucket)
			if err == nil {
				resultMap = append(resultMap, result...)
			}
		}
	}
	return resultMap, nil
}

func processAccountListBuckets(osqCtx context.Context, queryContext table.QueryContext, account *utilities.ExtensionConfigurationAwsAccount) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	tableConfig, ok := utilities.TableConfigurationMap["aws_s3_bucket"]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_s3_bucket",
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found")
	}
	result, err := processListBuckets(osqCtx, queryContext, tableConfig, account)
	if err != nil {
		return resultMap, err
	}
	resultMap = append(resultMap, result...)
	return resultMap, nil
}
