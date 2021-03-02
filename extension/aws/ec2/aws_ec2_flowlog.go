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

// DescribeFlowLogsColumns returns the list of columns in the table
func DescribeFlowLogsColumns() []table.ColumnDefinition {
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
		table.TextColumn("deliver_logs_error_message"),
		table.TextColumn("deliver_logs_permission_arn"),
		table.TextColumn("deliver_logs_status"),
		table.TextColumn("flow_log_id"),
		table.TextColumn("flow_log_status"),
		table.TextColumn("log_destination"),
		table.TextColumn("log_destination_type"),
		table.TextColumn("log_format"),
		table.TextColumn("log_group_name"),
		table.BigIntColumn("max_aggregation_interval"),
		table.TextColumn("resource_id"),
		table.TextColumn("tags"),
		//table.TextColumn("tags_key"),
		//table.TextColumn("tags_value"),
		table.TextColumn("traffic_type"),
	}
}

// DescribeFlowLogsGenerate returns the rows in the table for all configured accounts
func DescribeFlowLogsGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	if len(utilities.ExtConfiguration.ExtConfAws.Accounts) == 0 && extaws.ShouldProcessAccount("aws_ec2_flowlog", utilities.AwsAccountID) {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_ec2_flowlog",
			"account":   "default",
		}).Info("processing account")
		results, err := processAccountDescribeFlowLogs(osqCtx, queryContext, nil)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, results...)
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfAws.Accounts {
			if !extaws.ShouldProcessAccount("aws_ec2_flowlog", account.ID) {
				continue
			}
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_ec2_flowlog",
				"account":   account.ID,
			}).Info("processing account")
			results, err := processAccountDescribeFlowLogs(osqCtx, queryContext, &account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}

	return resultMap, nil
}

func processRegionDescribeFlowLogs(osqCtx context.Context, queryContext table.QueryContext, tableConfig *utilities.TableConfig, account *utilities.ExtensionConfigurationAwsAccount, region types.Region) ([]map[string]string, error) {
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
		"tableName": "aws_ec2_flowlog",
		"account":   accountId,
		"region":    *region.RegionName,
	}).Debug("processing region")

	svc := ec2.NewFromConfig(*sess)
	params := &ec2.DescribeFlowLogsInput{}

	paginator := ec2.NewDescribeFlowLogsPaginator(svc, params)

	for {
		page, err := paginator.NextPage(osqCtx)
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_ec2_flowlog",
				"account":   accountId,
				"region":    *region.RegionName,
				"task":      "DescribeFlowLogs",
				"errString": err.Error(),
			}).Error("failed to process region")
			return resultMap, err
		}
		byteArr, err := json.Marshal(page)
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_ec2_flowlog",
				"account":   accountId,
				"region":    *region.RegionName,
				"task":      "DescribeFlowLogs",
				"errString": err.Error(),
			}).Error("failed to marshal response")
			return nil, err
		}
		table := utilities.NewTable(byteArr, tableConfig)
		for _, row := range table.Rows {
			if !extaws.ShouldProcessRow(osqCtx, queryContext, "aws_ec2_flowlog", accountId, *region.RegionName, row) {
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

func processAccountDescribeFlowLogs(osqCtx context.Context, queryContext table.QueryContext, account *utilities.ExtensionConfigurationAwsAccount) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	awsSession, err := extaws.GetAwsConfig(account, "us-east-1")
	if err != nil {
		return resultMap, err
	}
	regions, err := extaws.FetchRegions(osqCtx, awsSession)
	if err != nil {
		return resultMap, err
	}
	tableConfig, ok := utilities.TableConfigurationMap["aws_ec2_flowlog"]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_ec2_flowlog",
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found")
	}
	for _, region := range regions {
		accountId := utilities.AwsAccountID
		if account != nil {
			accountId = account.ID
		}
		if !extaws.ShouldProcessRegion("aws_ec2_flowlog", accountId, *region.RegionName) {
			continue
		}
		result, err := processRegionDescribeFlowLogs(osqCtx, queryContext, tableConfig, account, region)
		if err != nil {
			continue
		}
		resultMap = append(resultMap, result...)
	}
	return resultMap, nil
}
