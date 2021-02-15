/**
 * Copyright (c) 2020-present, The cloudquery authors
 *
 * This source code is licensed as defined by the LICENSE file found in the
 * root directory of this source tree.
 *
 * SPDX-License-Identifier: (Apache-2.0 OR GPL-2.0-only)
 */

package config

import (
	"context"
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/Uptycs/cloudquery/utilities"

	extaws "github.com/Uptycs/cloudquery/extension/aws"
	"github.com/aws/aws-sdk-go-v2/service/configservice"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/Uptycs/basequery-go/plugin/table"
)

// DescribeDeliveryChannelsColumns returns the list of columns in the table
func DescribeDeliveryChannelsColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("account_id"),
		table.TextColumn("region_code"),
		table.TextColumn("config_snapshot_delivery_properties"),
		//table.TextColumn("config_snapshot_delivery_properties_delivery_frequency"),
		table.TextColumn("name"),
		table.TextColumn("s3_bucket_name"),
		table.TextColumn("s3_key_prefix"),
		table.TextColumn("sns_topic_arn"),
	}
}

// DescribeDeliveryChannelsGenerate returns the rows in the table for all configured accounts
func DescribeDeliveryChannelsGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	if len(utilities.ExtConfiguration.ExtConfAws.Accounts) == 0 {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_config_delivery_channel",
			"account":   "default",
		}).Info("processing account")
		results, err := processAccountDescribeDeliveryChannels(nil)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, results...)
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfAws.Accounts {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_config_delivery_channel",
				"account":   account.ID,
			}).Info("processing account")
			results, err := processAccountDescribeDeliveryChannels(&account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}

	return resultMap, nil
}

func processRegionDescribeDeliveryChannels(tableConfig *utilities.TableConfig, account *utilities.ExtensionConfigurationAwsAccount, region types.Region) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	sess, err := extaws.GetAwsConfig(account, *region.RegionName)
	if err != nil {
		return resultMap, err
	}

	accountId := utilities.AwsAccountID
	if account != nil {
		accountId = account.ID
	}

	utilities.GetLogger().WithFields(log.Fields{
		"tableName": "aws_config_delivery_channel",
		"account":   accountId,
		"region":    *region.RegionName,
	}).Debug("processing region")

	svc := configservice.NewFromConfig(*sess)
	params := &configservice.DescribeDeliveryChannelsInput{}

	result, err := svc.DescribeDeliveryChannels(context.TODO(), params)
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_config_delivery_channel",
			"account":   accountId,
			"region":    *region.RegionName,
			"task":      "DescribeDeliveryChannels",
			"errString": err.Error(),
		}).Error("failed to process region")
		return resultMap, err
	}

	byteArr, err := json.Marshal(result)
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_config_delivery_channel",
			"account":   accountId,
			"region":    *region.RegionName,
			"errString": err.Error(),
		}).Error("failed to marshal response")
		return resultMap, err
	}
	table := utilities.NewTable(byteArr, tableConfig)
	for _, row := range table.Rows {
		result := extaws.RowToMap(row, accountId, *region.RegionName, tableConfig)
		resultMap = append(resultMap, result)
	}
	return resultMap, nil
}

func processAccountDescribeDeliveryChannels(account *utilities.ExtensionConfigurationAwsAccount) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	awsSession, err := extaws.GetAwsConfig(account, "us-east-1")
	if err != nil {
		return resultMap, err
	}
	regions, err := extaws.FetchRegions(context.TODO(), awsSession)
	if err != nil {
		return resultMap, err
	}
	tableConfig, ok := utilities.TableConfigurationMap["aws_config_delivery_channel"]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_config_delivery_channel",
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found")
	}
	for _, region := range regions {
		result, err := processRegionDescribeDeliveryChannels(tableConfig, account, region)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, result...)
	}
	return resultMap, nil
}
