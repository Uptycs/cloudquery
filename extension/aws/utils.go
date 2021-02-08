/**
 * Copyright (c) 2020-present, The cloudquery authors
 *
 * This source code is licensed as defined by the LICENSE file found in the
 * root directory of this source tree.
 *
 * SPDX-License-Identifier: (Apache-2.0 OR GPL-2.0-only)
 */

package aws

import (
	"context"
	"time"

	"github.com/Uptycs/cloudquery/utilities"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	log "github.com/sirupsen/logrus"
)

// GetAwsConfig creates an AWS Config for given account.
// If account is nil, it creates a default config
func GetAwsConfig(account *utilities.ExtensionConfigurationAwsAccount, regionCode string) (*aws.Config, error) {
	if account == nil {
		utilities.GetLogger().Debug("creating default session")
		return getDefaultAwsConfig(regionCode)
	}

	if len(account.ProfileName) != 0 && len(account.RoleArn) == 0 {
		utilities.GetLogger().Debug("creating session using profile")
		return getAwsConfigForProfile(account, regionCode)
	} else if len(account.RoleArn) != 0 {
		utilities.GetLogger().Debug("creating session using roleArn")
		return getAwsConfigForRole(account, regionCode)
	} else {
		utilities.GetLogger().Debug("creating default session")
		return getDefaultAwsConfig(regionCode)
	}
}

func getAwsConfigForProfile(account *utilities.ExtensionConfigurationAwsAccount, regionCode string) (*aws.Config, error) {
	utilities.GetLogger().WithFields(log.Fields{
		"account": account.ID,
		"region":  regionCode,
		"profile": account.ProfileName,
	}).Debug("creating config")
	credentialFiles := make([]string, 0)
	credentialFiles = append(credentialFiles, account.CredentialFile)
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(regionCode),
		config.WithSharedCredentialsFiles(credentialFiles),
		config.WithSharedConfigProfile(account.ProfileName),
	)
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"account":   account.ID,
			"profile":   account.ProfileName,
			"errString": err.Error(),
		}).Error("failed to create config")
		return nil, err
	}
	return &cfg, nil
}

func getAwsConfigForRole(account *utilities.ExtensionConfigurationAwsAccount, regionCode string) (*aws.Config, error) {
	utilities.GetLogger().WithFields(log.Fields{
		"account": account.ID,
		"region":  regionCode,
		"profile": account.ProfileName,
		"role":    account.RoleArn,
	}).Debug("creating config")
	credentialFiles := make([]string, 0)
	credentialFiles = append(credentialFiles, account.CredentialFile)
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(regionCode),
		config.WithSharedCredentialsFiles(credentialFiles),
		config.WithSharedConfigProfile(account.ProfileName),
	)
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"account":   account.ID,
			"role":      account.RoleArn,
			"errString": err.Error(),
		}).Error("failed to create config")
		return nil, err
	}
	// Create the credentials from AssumeRoleProvider to assume the role
	// referenced by the role ARN.
	stsSvc := sts.NewFromConfig(cfg)
	creds := stscreds.NewAssumeRoleProvider(stsSvc, account.RoleArn, func(options *stscreds.AssumeRoleOptions) {
		options.Duration = time.Duration(60) * time.Minute
		options.ExternalID = &account.ExternalID
	})
	cfg.Credentials = aws.NewCredentialsCache(creds)
	return &cfg, nil
}

func getDefaultAwsConfig(regionCode string) (*aws.Config, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(regionCode),
	)
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"account":   "default",
			"region":    regionCode,
			"errString": err.Error(),
		}).Error("failed to create config")
		return nil, err
	}
	return &cfg, nil
}

// FetchRegions returns the list of regions for given AWS config
func FetchRegions(ctx context.Context, awsConfig *aws.Config) ([]types.Region, error) {
	svc := ec2.NewFromConfig(*awsConfig)
	awsRegions, err := svc.DescribeRegions(ctx, &ec2.DescribeRegionsInput{})
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"errString": err.Error(),
		}).Error("failed to get regions")
		return nil, err
	}
	return awsRegions.Regions, nil
}

// RowToMap converts JSON row into osquery row.
// If configured it will copy some metadata vaues into appropriate columns
func RowToMap(row map[string]interface{}, accountId string, region string, tableConfig *utilities.TableConfig) map[string]string {
	result := make(map[string]string)

	if len(tableConfig.Aws.AccountIDAttribute) != 0 {
		result[tableConfig.Aws.AccountIDAttribute] = accountId
	}
	if len(tableConfig.Aws.RegionCodeAttribute) != 0 {
		result[tableConfig.Aws.RegionCodeAttribute] = region
	}
	if len(tableConfig.Aws.RegionAttribute) != 0 {
		result[tableConfig.Aws.RegionAttribute] = region // TODO: Fix it
	}

	result = utilities.RowToMap(result, row, tableConfig)
	return result
}
