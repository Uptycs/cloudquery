/**
 * Copyright (c) 2020-present, The cloudquery authors
 *
 * This source code is licensed as defined by the LICENSE file found in the
 * root directory of this source tree.
 *
 * SPDX-License-Identifier: (Apache-2.0 OR GPL-2.0-only)
 */

package ec2

import (
	"context"
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/Uptycs/cloudquery/utilities"

	"github.com/Uptycs/basequery-go/plugin/table"
	extaws "github.com/Uptycs/cloudquery/extension/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

// DescribeSnapshotsColumns returns the list of columns in the table
func DescribeSnapshotsColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("account_id"),
		table.TextColumn("region_code"),
		table.TextColumn("data_encryption_key_id"),
		table.TextColumn("description"),
		table.TextColumn("encrypted"),
		table.TextColumn("kms_key_id"),
		table.TextColumn("owner_alias"),
		table.TextColumn("owner_id"),
		table.TextColumn("progress"),
		table.TextColumn("snapshot_id"),
		table.TextColumn("start_time"),
		//table.BigIntColumn("start_time_ext"),
		//table.TextColumn("start_time_loc"),
		//table.BigIntColumn("start_time_loc_cache_end"),
		//table.BigIntColumn("start_time_loc_cache_start"),
		//table.TextColumn("start_time_loc_cache_zone"),
		//table.TextColumn("start_time_loc_cache_zone_is_dst"),
		//table.TextColumn("start_time_loc_cache_zone_name"),
		//table.IntegerColumn("start_time_loc_cache_zone_offset"),
		//table.TextColumn("start_time_loc_name"),
		//table.TextColumn("start_time_loc_tx"),
		//table.IntegerColumn("start_time_loc_tx_index"),
		//table.TextColumn("start_time_loc_tx_isstd"),
		//table.TextColumn("start_time_loc_tx_isutc"),
		//table.BigIntColumn("start_time_loc_tx_when"),
		//table.TextColumn("start_time_loc_zone"),
		//table.TextColumn("start_time_loc_zone_is_dst"),
		//table.TextColumn("start_time_loc_zone_name"),
		//table.IntegerColumn("start_time_loc_zone_offset"),
		//table.BigIntColumn("start_time_wall"),
		table.TextColumn("state"),
		table.TextColumn("state_message"),
		table.TextColumn("tags"),
		//table.TextColumn("tags_key"),
		//table.TextColumn("tags_value"),
		table.TextColumn("volume_id"),
		table.BigIntColumn("volume_size"),
	}
}

// DescribeSnapshotsGenerate returns the rows in the table for all configured accounts
func DescribeSnapshotsGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	if len(utilities.ExtConfiguration.ExtConfAws.Accounts) == 0 && extaws.ShouldProcessAccount("aws_ec2_snapshot", utilities.AwsAccountID) {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_ec2_snapshot",
			"account":   "default",
		}).Info("processing account")
		results, err := processAccountDescribeSnapshots(osqCtx, queryContext, nil)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, results...)
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfAws.Accounts {
			if !extaws.ShouldProcessAccount("aws_ec2_snapshot", account.ID) {
				continue
			}
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_ec2_snapshot",
				"account":   account.ID,
			}).Info("processing account")
			results, err := processAccountDescribeSnapshots(osqCtx, queryContext, &account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}

	return resultMap, nil
}

func processRegionDescribeSnapshots(osqCtx context.Context, queryContext table.QueryContext, tableConfig *utilities.TableConfig, account *utilities.ExtensionConfigurationAwsAccount, region types.Region) ([]map[string]string, error) {
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
		"tableName": "aws_ec2_snapshot",
		"account":   accountId,
		"region":    *region.RegionName,
	}).Debug("processing region")

	svc := ec2.NewFromConfig(*sess)
	params := &ec2.DescribeSnapshotsInput{}

	paginator := ec2.NewDescribeSnapshotsPaginator(svc, params)

	for {
		page, err := paginator.NextPage(osqCtx)
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_ec2_snapshot",
				"account":   accountId,
				"region":    *region.RegionName,
				"task":      "DescribeSnapshots",
				"errString": err.Error(),
			}).Error("failed to process region")
			return resultMap, err
		}
		byteArr, err := json.Marshal(page)
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_ec2_snapshot",
				"account":   accountId,
				"region":    *region.RegionName,
				"task":      "DescribeSnapshots",
				"errString": err.Error(),
			}).Error("failed to marshal response")
			return nil, err
		}
		table := utilities.NewTable(byteArr, tableConfig)
		for _, row := range table.Rows {
			if !extaws.ShouldProcessRow(osqCtx, queryContext, "aws_ec2_snapshot", accountId, *region.RegionName, row) {
				continue
			}
			result := extaws.RowToMap(row, accountId, *region.RegionName, tableConfig)
			resultMap = append(resultMap, result)
		}
		if !paginator.HasMorePages() {
			break
		}
	}
	return resultMap, nil
}

func processAccountDescribeSnapshots(osqCtx context.Context, queryContext table.QueryContext, account *utilities.ExtensionConfigurationAwsAccount) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	awsSession, err := extaws.GetAwsConfig(account, "us-east-1")
	if err != nil {
		return resultMap, err
	}
	regions, err := extaws.FetchRegions(osqCtx, awsSession)
	if err != nil {
		return resultMap, err
	}
	tableConfig, ok := utilities.TableConfigurationMap["aws_ec2_snapshot"]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_ec2_snapshot",
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found")
	}
	for _, region := range regions {
		accountId := utilities.AwsAccountID
		if account != nil {
			accountId = account.ID
		}
		if !extaws.ShouldProcessRegion("aws_ec2_snapshot", accountId, *region.RegionName) {
			continue
		}
		result, err := processRegionDescribeSnapshots(osqCtx, queryContext, tableConfig, account, region)
		if err != nil {
			continue
		}
		resultMap = append(resultMap, result...)
	}
	return resultMap, nil
}
