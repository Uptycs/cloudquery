/**
 * Copyright (c) 2020-present, The cloudquery authors
 *
 * This source code is licensed as defined by the LICENSE file found in the
 * root directory of this source tree.
 *
 * SPDX-License-Identifier: (Apache-2.0 OR GPL-2.0-only)
 */

package efs

import (
	"context"
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/Uptycs/cloudquery/utilities"

	"github.com/Uptycs/basequery-go/plugin/table"
	extaws "github.com/Uptycs/cloudquery/extension/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/efs"
)

// DescribeFileSystemsColumns returns the list of columns in the table
func DescribeFileSystemsColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("account_id"),
		table.TextColumn("region_code"),
		table.TextColumn("creation_time"),
		//table.BigIntColumn("creation_time_ext"),
		//table.TextColumn("creation_time_loc"),
		//table.BigIntColumn("creation_time_loc_cache_end"),
		//table.BigIntColumn("creation_time_loc_cache_start"),
		//table.TextColumn("creation_time_loc_cache_zone"),
		//table.TextColumn("creation_time_loc_cache_zone_is_dst"),
		//table.TextColumn("creation_time_loc_cache_zone_name"),
		//table.IntegerColumn("creation_time_loc_cache_zone_offset"),
		//table.TextColumn("creation_time_loc_name"),
		//table.TextColumn("creation_time_loc_tx"),
		//table.IntegerColumn("creation_time_loc_tx_index"),
		//table.TextColumn("creation_time_loc_tx_isstd"),
		//table.TextColumn("creation_time_loc_tx_isutc"),
		//table.BigIntColumn("creation_time_loc_tx_when"),
		//table.TextColumn("creation_time_loc_zone"),
		//table.TextColumn("creation_time_loc_zone_is_dst"),
		//table.TextColumn("creation_time_loc_zone_name"),
		//table.IntegerColumn("creation_time_loc_zone_offset"),
		//table.BigIntColumn("creation_time_wall"),
		table.TextColumn("creation_token"),
		table.TextColumn("encrypted"),
		table.TextColumn("file_system_arn"),
		table.TextColumn("file_system_id"),
		table.TextColumn("kms_key_id"),
		table.TextColumn("life_cycle_state"),
		table.TextColumn("name"),
		table.BigIntColumn("number_of_mount_targets"),
		table.TextColumn("owner_id"),
		table.TextColumn("performance_mode"),
		table.DoubleColumn("provisioned_throughput_in_mibps"),
		table.TextColumn("size_in_bytes"),
		//table.TextColumn("size_in_bytes_timestamp"),
		//table.BigIntColumn("size_in_bytes_timestamp_ext"),
		//table.TextColumn("size_in_bytes_timestamp_loc"),
		//table.BigIntColumn("size_in_bytes_timestamp_loc_cache_end"),
		//table.BigIntColumn("size_in_bytes_timestamp_loc_cache_start"),
		//table.TextColumn("size_in_bytes_timestamp_loc_cache_zone"),
		//table.TextColumn("size_in_bytes_timestamp_loc_cache_zone_is_dst"),
		//table.TextColumn("size_in_bytes_timestamp_loc_cache_zone_name"),
		//table.IntegerColumn("size_in_bytes_timestamp_loc_cache_zone_offset"),
		//table.TextColumn("size_in_bytes_timestamp_loc_name"),
		//table.TextColumn("size_in_bytes_timestamp_loc_tx"),
		//table.IntegerColumn("size_in_bytes_timestamp_loc_tx_index"),
		//table.TextColumn("size_in_bytes_timestamp_loc_tx_isstd"),
		//table.TextColumn("size_in_bytes_timestamp_loc_tx_isutc"),
		//table.BigIntColumn("size_in_bytes_timestamp_loc_tx_when"),
		//table.TextColumn("size_in_bytes_timestamp_loc_zone"),
		//table.TextColumn("size_in_bytes_timestamp_loc_zone_is_dst"),
		//table.TextColumn("size_in_bytes_timestamp_loc_zone_name"),
		//table.IntegerColumn("size_in_bytes_timestamp_loc_zone_offset"),
		//table.BigIntColumn("size_in_bytes_timestamp_wall"),
		//table.BigIntColumn("size_in_bytes_value"),
		//table.BigIntColumn("size_in_bytes_value_in_ia"),
		//table.BigIntColumn("size_in_bytes_value_in_standard"),
		table.TextColumn("tags"),
		//table.TextColumn("tags_key"),
		//table.TextColumn("tags_value"),
		table.TextColumn("throughput_mode"),
	}
}

// DescribeFileSystemsGenerate returns the rows in the table for all configured accounts
func DescribeFileSystemsGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	if len(utilities.ExtConfiguration.ExtConfAws.Accounts) == 0 {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_efs_file_system",
			"account":   "default",
		}).Info("processing account")
		results, err := processAccountDescribeFileSystems(nil)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, results...)
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfAws.Accounts {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_efs_file_system",
				"account":   account.ID,
			}).Info("processing account")
			results, err := processAccountDescribeFileSystems(&account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}

	return resultMap, nil
}

func processRegionDescribeFileSystems(tableConfig *utilities.TableConfig, account *utilities.ExtensionConfigurationAwsAccount, region types.Region) ([]map[string]string, error) {
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
		"tableName": "aws_efs_file_system",
		"account":   accountId,
		"region":    *region.RegionName,
	}).Debug("processing region")

	svc := efs.NewFromConfig(*sess)
	params := &efs.DescribeFileSystemsInput{}

	paginator := efs.NewDescribeFileSystemsPaginator(svc, params)

	for {
		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_efs_file_system",
				"account":   accountId,
				"region":    *region.RegionName,
				"task":      "DescribeFileSystems",
				"errString": err.Error(),
			}).Error("failed to process region")
			return resultMap, err
		}
		byteArr, err := json.Marshal(page)
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_efs_file_system",
				"account":   accountId,
				"region":    *region.RegionName,
				"task":      "DescribeFileSystems",
				"errString": err.Error(),
			}).Error("failed to marshal response")
			return nil, err
		}
		table := utilities.NewTable(byteArr, tableConfig)
		for _, row := range table.Rows {
			result := extaws.RowToMap(row, accountId, *region.RegionName, tableConfig)
			resultMap = append(resultMap, result)
		}
		if !paginator.HasMorePages() {
			break
		}
	}
	return resultMap, nil
}

func processAccountDescribeFileSystems(account *utilities.ExtensionConfigurationAwsAccount) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	awsSession, err := extaws.GetAwsConfig(account, "us-east-1")
	if err != nil {
		return resultMap, err
	}
	regions, err := extaws.FetchRegions(context.TODO(), awsSession)
	if err != nil {
		return resultMap, err
	}
	tableConfig, ok := utilities.TableConfigurationMap["aws_efs_file_system"]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_efs_file_system",
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found")
	}
	for _, region := range regions {
		result, err := processRegionDescribeFileSystems(tableConfig, account, region)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, result...)
	}
	return resultMap, nil
}
