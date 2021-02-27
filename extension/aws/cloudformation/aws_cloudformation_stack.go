/**
 * Copyright (c) 2020-present, The cloudquery authors
 *
 * This source code is licensed as defined by the LICENSE file found in the
 * root directory of this source tree.
 *
 * SPDX-License-Identifier: (Apache-2.0 OR GPL-2.0-only)
 */

package cloudformation

import (
	"context"
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/Uptycs/cloudquery/utilities"

	"github.com/Uptycs/basequery-go/plugin/table"
	extaws "github.com/Uptycs/cloudquery/extension/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

// DescribeStacksColumns returns the list of columns in the table
func DescribeStacksColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("account_id"),
		table.TextColumn("region_code"),
		//table.TextColumn("values"),
		table.TextColumn("capabilities"),
		table.TextColumn("change_set_id"),
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
		table.TextColumn("deletion_time"),
		//table.BigIntColumn("deletion_time_ext"),
		//table.TextColumn("deletion_time_loc"),
		//table.BigIntColumn("deletion_time_loc_cache_end"),
		//table.BigIntColumn("deletion_time_loc_cache_start"),
		//table.TextColumn("deletion_time_loc_cache_zone"),
		//table.TextColumn("deletion_time_loc_cache_zone_is_dst"),
		//table.TextColumn("deletion_time_loc_cache_zone_name"),
		//table.IntegerColumn("deletion_time_loc_cache_zone_offset"),
		//table.TextColumn("deletion_time_loc_name"),
		//table.TextColumn("deletion_time_loc_tx"),
		//table.IntegerColumn("deletion_time_loc_tx_index"),
		//table.TextColumn("deletion_time_loc_tx_isstd"),
		//table.TextColumn("deletion_time_loc_tx_isutc"),
		//table.BigIntColumn("deletion_time_loc_tx_when"),
		//table.TextColumn("deletion_time_loc_zone"),
		//table.TextColumn("deletion_time_loc_zone_is_dst"),
		//table.TextColumn("deletion_time_loc_zone_name"),
		//table.IntegerColumn("deletion_time_loc_zone_offset"),
		//table.BigIntColumn("deletion_time_wall"),
		table.TextColumn("description"),
		table.TextColumn("disable_rollback"),
		table.TextColumn("drift_information"),
		//table.TextColumn("drift_information_last_check_timestamp"),
		//table.BigIntColumn("drift_information_last_check_timestamp_ext"),
		//table.TextColumn("drift_information_last_check_timestamp_loc"),
		//table.BigIntColumn("drift_information_last_check_timestamp_loc_cache_end"),
		//table.BigIntColumn("drift_information_last_check_timestamp_loc_cache_start"),
		//table.TextColumn("drift_information_last_check_timestamp_loc_cache_zone"),
		//table.TextColumn("drift_information_last_check_timestamp_loc_cache_zone_is_dst"),
		//table.TextColumn("drift_information_last_check_timestamp_loc_cache_zone_name"),
		//table.IntegerColumn("drift_information_last_check_timestamp_loc_cache_zone_offset"),
		//table.TextColumn("drift_information_last_check_timestamp_loc_name"),
		//table.TextColumn("drift_information_last_check_timestamp_loc_tx"),
		//table.IntegerColumn("drift_information_last_check_timestamp_loc_tx_index"),
		//table.TextColumn("drift_information_last_check_timestamp_loc_tx_isstd"),
		//table.TextColumn("drift_information_last_check_timestamp_loc_tx_isutc"),
		//table.BigIntColumn("drift_information_last_check_timestamp_loc_tx_when"),
		//table.TextColumn("drift_information_last_check_timestamp_loc_zone"),
		//table.TextColumn("drift_information_last_check_timestamp_loc_zone_is_dst"),
		//table.TextColumn("drift_information_last_check_timestamp_loc_zone_name"),
		//table.IntegerColumn("drift_information_last_check_timestamp_loc_zone_offset"),
		//table.BigIntColumn("drift_information_last_check_timestamp_wall"),
		//table.TextColumn("drift_information_stack_drift_status"),
		table.TextColumn("enable_termination_protection"),
		table.TextColumn("last_updated_time"),
		//table.BigIntColumn("last_updated_time_ext"),
		//table.TextColumn("last_updated_time_loc"),
		//table.BigIntColumn("last_updated_time_loc_cache_end"),
		//table.BigIntColumn("last_updated_time_loc_cache_start"),
		//table.TextColumn("last_updated_time_loc_cache_zone"),
		//table.TextColumn("last_updated_time_loc_cache_zone_is_dst"),
		//table.TextColumn("last_updated_time_loc_cache_zone_name"),
		//table.IntegerColumn("last_updated_time_loc_cache_zone_offset"),
		//table.TextColumn("last_updated_time_loc_name"),
		//table.TextColumn("last_updated_time_loc_tx"),
		//table.IntegerColumn("last_updated_time_loc_tx_index"),
		//table.TextColumn("last_updated_time_loc_tx_isstd"),
		//table.TextColumn("last_updated_time_loc_tx_isutc"),
		//table.BigIntColumn("last_updated_time_loc_tx_when"),
		//table.TextColumn("last_updated_time_loc_zone"),
		//table.TextColumn("last_updated_time_loc_zone_is_dst"),
		//table.TextColumn("last_updated_time_loc_zone_name"),
		//table.IntegerColumn("last_updated_time_loc_zone_offset"),
		//table.BigIntColumn("last_updated_time_wall"),
		table.TextColumn("notification_arns"),
		table.TextColumn("outputs"),
		//table.TextColumn("outputs_description"),
		//table.TextColumn("outputs_export_name"),
		//table.TextColumn("outputs_output_key"),
		//table.TextColumn("outputs_output_value"),
		table.TextColumn("parameters"),
		//table.TextColumn("parameters_parameter_key"),
		//table.TextColumn("parameters_parameter_value"),
		//table.TextColumn("parameters_resolved_value"),
		//table.TextColumn("parameters_use_previous_value"),
		table.TextColumn("parent_id"),
		table.TextColumn("role_arn"),
		table.TextColumn("rollback_configuration"),
		//table.IntegerColumn("rollback_configuration_monitoring_time_in_minutes"),
		//table.TextColumn("rollback_configuration_rollback_triggers"),
		//table.TextColumn("rollback_configuration_rollback_triggers_arn"),
		//table.TextColumn("rollback_configuration_rollback_triggers_type"),
		table.TextColumn("root_id"),
		table.TextColumn("stack_id"),
		table.TextColumn("stack_name"),
		table.TextColumn("stack_status"),
		//table.TextColumn("stack_status_reason"),
		table.TextColumn("tags"),
		//table.TextColumn("tags_key"),
		//table.TextColumn("tags_value"),
		table.IntegerColumn("timeout_in_minutes"),
	}
}

// DescribeStacksGenerate returns the rows in the table for all configured accounts
func DescribeStacksGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	if len(utilities.ExtConfiguration.ExtConfAws.Accounts) == 0 && extaws.ShouldProcessAccount("aws_cloudformation_stack", utilities.AwsAccountID) {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_cloudformation_stack",
			"account":   "default",
		}).Info("processing account")
		results, err := processAccountDescribeStacks(osqCtx, queryContext, nil)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, results...)
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfAws.Accounts {
			if !extaws.ShouldProcessAccount("aws_cloudformation_stack", account.ID) {
				continue
			}
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_cloudformation_stack",
				"account":   account.ID,
			}).Info("processing account")
			results, err := processAccountDescribeStacks(osqCtx, queryContext, &account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}

	return resultMap, nil
}

func processRegionDescribeStacks(osqCtx context.Context, queryContext table.QueryContext, tableConfig *utilities.TableConfig, account *utilities.ExtensionConfigurationAwsAccount, region types.Region) ([]map[string]string, error) {
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
		"tableName": "aws_cloudformation_stack",
		"account":   accountId,
		"region":    *region.RegionName,
	}).Debug("processing region")

	svc := cloudformation.NewFromConfig(*sess)
	params := &cloudformation.DescribeStacksInput{}

	paginator := cloudformation.NewDescribeStacksPaginator(svc, params)

	for {
		page, err := paginator.NextPage(osqCtx)
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_cloudformation_stack",
				"account":   accountId,
				"region":    *region.RegionName,
				"task":      "DescribeStacks",
				"errString": err.Error(),
			}).Error("failed to process region")
			return resultMap, err
		}
		byteArr, err := json.Marshal(page)
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_cloudformation_stack",
				"account":   accountId,
				"region":    *region.RegionName,
				"task":      "DescribeStacks",
				"errString": err.Error(),
			}).Error("failed to marshal response")
			return nil, err
		}
		table := utilities.NewTable(byteArr, tableConfig)
		for _, row := range table.Rows {
			if !extaws.ShouldProcessRow(osqCtx, queryContext, "aws_cloudformation_stack", accountId, *region.RegionName, row) {
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

func processAccountDescribeStacks(osqCtx context.Context, queryContext table.QueryContext, account *utilities.ExtensionConfigurationAwsAccount) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	awsSession, err := extaws.GetAwsConfig(account, "us-east-1")
	if err != nil {
		return resultMap, err
	}
	regions, err := extaws.FetchRegions(osqCtx, awsSession)
	if err != nil {
		return resultMap, err
	}
	tableConfig, ok := utilities.TableConfigurationMap["aws_cloudformation_stack"]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_cloudformation_stack",
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found")
	}
	for _, region := range regions {
		accountId := utilities.AwsAccountID
		if account != nil {
			accountId = account.ID
		}
		if !extaws.ShouldProcessRegion("aws_cloudformation_stack", accountId, *region.RegionName) {
			continue
		}
		result, err := processRegionDescribeStacks(osqCtx, queryContext, tableConfig, account, region)
		if err != nil {
			continue
		}
		resultMap = append(resultMap, result...)
	}
	return resultMap, nil
}
