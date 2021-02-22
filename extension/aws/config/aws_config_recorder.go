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

	"github.com/Uptycs/basequery-go/plugin/table"
	extaws "github.com/Uptycs/cloudquery/extension/aws"
	"github.com/aws/aws-sdk-go-v2/service/configservice"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

// DescribeConfigurationRecordersColumns returns the list of columns in the table
func DescribeConfigurationRecordersColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("account_id"),
		table.TextColumn("region_code"),
		table.TextColumn("name"),
		table.TextColumn("recording_group"),
		table.TextColumn("recording_group_all_supported"),
		table.TextColumn("recording_group_include_global_resource_types"),
		table.TextColumn("recording_group_resource_types"),
		table.TextColumn("role_arn"),
	}
}

// DescribeConfigurationRecordersGenerate returns the rows in the table for all configured accounts
func DescribeConfigurationRecordersGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	if len(utilities.ExtConfiguration.ExtConfAws.Accounts) == 0 {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_config_recorder",
			"account":   "default",
		}).Info("processing account")
		results, err := processAccountDescribeConfigurationRecorders(nil)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, results...)
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfAws.Accounts {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_config_recorder",
				"account":   account.ID,
			}).Info("processing account")
			results, err := processAccountDescribeConfigurationRecorders(&account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}

	return resultMap, nil
}

func processRegionDescribeConfigurationRecorders(tableConfig *utilities.TableConfig, account *utilities.ExtensionConfigurationAwsAccount, region types.Region) ([]map[string]string, error) {
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
		"tableName": "aws_config_recorder",
		"account":   accountId,
		"region":    *region.RegionName,
	}).Debug("processing region")

	svc := configservice.NewFromConfig(*sess)
	params := &configservice.DescribeConfigurationRecordersInput{}

	result, err := svc.DescribeConfigurationRecorders(context.TODO(), params)
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_config_recorder",
			"account":   accountId,
			"region":    *region.RegionName,
			"task":      "DescribeConfigurationRecorders",
			"errString": err.Error(),
		}).Error("failed to process region")
		return resultMap, err
	}

	byteArr, err := json.Marshal(result)
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_config_recorder",
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

func processAccountDescribeConfigurationRecorders(account *utilities.ExtensionConfigurationAwsAccount) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	awsSession, err := extaws.GetAwsConfig(account, "us-east-1")
	if err != nil {
		return resultMap, err
	}
	regions, err := extaws.FetchRegions(context.TODO(), awsSession)
	if err != nil {
		return resultMap, err
	}
	tableConfig, ok := utilities.TableConfigurationMap["aws_config_recorder"]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_config_recorder",
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found")
	}
	for _, region := range regions {
		result, err := processRegionDescribeConfigurationRecorders(tableConfig, account, region)
		if err != nil {
			continue
		}
		resultMap = append(resultMap, result...)
	}
	return resultMap, nil
}
