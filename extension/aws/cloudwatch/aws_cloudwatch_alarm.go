/**
 * Copyright (c) 2020-present, The cloudquery authors
 *
 * This source code is licensed as defined by the LICENSE file found in the
 * root directory of this source tree.
 *
 * SPDX-License-Identifier: (Apache-2.0 OR GPL-2.0-only)
 */

package cloudwatch

import (
	"context"
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/Uptycs/cloudquery/utilities"

	"github.com/Uptycs/basequery-go/plugin/table"
	extaws "github.com/Uptycs/cloudquery/extension/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

// DescribeAlarmsColumns returns the list of columns in the table
func DescribeAlarmsColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("account_id"),
		table.TextColumn("region_code"),
		table.TextColumn("actions_enabled"),
		table.TextColumn("alarm_actions"),
		table.TextColumn("alarm_arn"),
		table.TextColumn("alarm_configuration_updated_timestamp"),
		//table.BigIntColumn("alarm_configuration_updated_timestamp_ext"),
		//table.TextColumn("alarm_configuration_updated_timestamp_loc"),
		//table.BigIntColumn("alarm_configuration_updated_timestamp_loc_cache_end"),
		//table.BigIntColumn("alarm_configuration_updated_timestamp_loc_cache_start"),
		//table.TextColumn("alarm_configuration_updated_timestamp_loc_cache_zone"),
		//table.TextColumn("alarm_configuration_updated_timestamp_loc_cache_zone_is_dst"),
		//table.TextColumn("alarm_configuration_updated_timestamp_loc_cache_zone_name"),
		//table.IntegerColumn("alarm_configuration_updated_timestamp_loc_cache_zone_offset"),
		//table.TextColumn("alarm_configuration_updated_timestamp_loc_name"),
		//table.TextColumn("alarm_configuration_updated_timestamp_loc_tx"),
		//table.IntegerColumn("alarm_configuration_updated_timestamp_loc_tx_index"),
		//table.TextColumn("alarm_configuration_updated_timestamp_loc_tx_isstd"),
		//table.TextColumn("alarm_configuration_updated_timestamp_loc_tx_isutc"),
		//table.BigIntColumn("alarm_configuration_updated_timestamp_loc_tx_when"),
		//table.TextColumn("alarm_configuration_updated_timestamp_loc_zone"),
		//table.TextColumn("alarm_configuration_updated_timestamp_loc_zone_is_dst"),
		//table.TextColumn("alarm_configuration_updated_timestamp_loc_zone_name"),
		//table.IntegerColumn("alarm_configuration_updated_timestamp_loc_zone_offset"),
		//table.BigIntColumn("alarm_configuration_updated_timestamp_wall"),
		table.TextColumn("alarm_description"),
		table.TextColumn("alarm_name"),
		table.TextColumn("comparison_operator"),
		table.BigIntColumn("datapoints_to_alarm"),
		table.TextColumn("dimensions"),
		table.TextColumn("dimensions_name"),
		table.TextColumn("dimensions_value"),
		table.TextColumn("evaluate_low_sample_count_percentile"),
		table.BigIntColumn("evaluation_periods"),
		table.TextColumn("extended_statistic"),
		table.TextColumn("insufficient_data_actions"),
		table.TextColumn("metric_name"),
		table.TextColumn("metrics"),
		//table.TextColumn("metrics_expression"),
		//table.TextColumn("metrics_id"),
		table.TextColumn("metrics_label"),
		table.TextColumn("metrics_metric_stat"),
		//table.TextColumn("metrics_metric_stat_metric"),
		//table.TextColumn("metrics_metric_stat_metric_dimensions"),
		//table.TextColumn("metrics_metric_stat_metric_dimensions_name"),
		//table.TextColumn("metrics_metric_stat_metric_dimensions_value"),
		//table.TextColumn("metrics_metric_stat_metric_metric_name"),
		//table.TextColumn("metrics_metric_stat_metric_namespace"),
		table.BigIntColumn("metrics_metric_stat_period"),
		table.TextColumn("metrics_metric_stat_stat"),
		table.TextColumn("metrics_metric_stat_unit"),
		table.BigIntColumn("metrics_period"),
		table.TextColumn("metrics_return_data"),
		table.TextColumn("namespace"),
		table.TextColumn("ok_actions"),
		table.BigIntColumn("period"),
		table.TextColumn("state_reason"),
		table.TextColumn("state_reason_data"),
		table.TextColumn("state_updated_timestamp"),
		//table.BigIntColumn("state_updated_timestamp_ext"),
		//table.TextColumn("state_updated_timestamp_loc"),
		//table.BigIntColumn("state_updated_timestamp_loc_cache_end"),
		//table.BigIntColumn("state_updated_timestamp_loc_cache_start"),
		//table.TextColumn("state_updated_timestamp_loc_cache_zone"),
		//table.TextColumn("state_updated_timestamp_loc_cache_zone_is_dst"),
		//table.TextColumn("state_updated_timestamp_loc_cache_zone_name"),
		//table.IntegerColumn("state_updated_timestamp_loc_cache_zone_offset"),
		//table.TextColumn("state_updated_timestamp_loc_name"),
		//table.TextColumn("state_updated_timestamp_loc_tx"),
		//table.IntegerColumn("state_updated_timestamp_loc_tx_index"),
		//table.TextColumn("state_updated_timestamp_loc_tx_isstd"),
		//table.TextColumn("state_updated_timestamp_loc_tx_isutc"),
		//table.BigIntColumn("state_updated_timestamp_loc_tx_when"),
		//table.TextColumn("state_updated_timestamp_loc_zone"),
		//table.TextColumn("state_updated_timestamp_loc_zone_is_dst"),
		//table.TextColumn("state_updated_timestamp_loc_zone_name"),
		//table.IntegerColumn("state_updated_timestamp_loc_zone_offset"),
		table.BigIntColumn("state_updated_timestamp_wall"),
		table.TextColumn("state_value"),
		table.TextColumn("statistic"),
		table.DoubleColumn("threshold"),
		table.TextColumn("threshold_metric_id"),
		table.TextColumn("treat_missing_data"),
		table.TextColumn("unit"),
	}
}

// DescribeAlarmsGenerate returns the rows in the table for all configured accounts
func DescribeAlarmsGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	if len(utilities.ExtConfiguration.ExtConfAws.Accounts) == 0 {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_cloudwatch_alarm",
			"account":   "default",
		}).Info("processing account")
		results, err := processAccountDescribeAlarms(nil)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, results...)
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfAws.Accounts {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_cloudwatch_alarm",
				"account":   account.ID,
			}).Info("processing account")
			results, err := processAccountDescribeAlarms(&account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}

	return resultMap, nil
}

func processRegionDescribeAlarms(tableConfig *utilities.TableConfig, account *utilities.ExtensionConfigurationAwsAccount, region types.Region) ([]map[string]string, error) {
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
		"tableName": "aws_cloudwatch_alarm",
		"account":   accountId,
		"region":    *region.RegionName,
	}).Debug("processing region")

	svc := cloudwatch.NewFromConfig(*sess)
	params := &cloudwatch.DescribeAlarmsInput{}

	paginator := cloudwatch.NewDescribeAlarmsPaginator(svc, params)

	for {
		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_cloudwatch_alarm",
				"account":   accountId,
				"region":    *region.RegionName,
				"task":      "DescribeAlarms",
				"errString": err.Error(),
			}).Error("failed to process region")
			return resultMap, err
		}
		byteArr, err := json.Marshal(page)
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_cloudwatch_alarm",
				"account":   accountId,
				"region":    *region.RegionName,
				"task":      "DescribeAlarms",
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

func processAccountDescribeAlarms(account *utilities.ExtensionConfigurationAwsAccount) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	awsSession, err := extaws.GetAwsConfig(account, "us-east-1")
	if err != nil {
		return resultMap, err
	}
	regions, err := extaws.FetchRegions(context.TODO(), awsSession)
	if err != nil {
		return resultMap, err
	}
	tableConfig, ok := utilities.TableConfigurationMap["aws_cloudwatch_alarm"]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_cloudwatch_alarm",
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found")
	}
	for _, region := range regions {
		result, err := processRegionDescribeAlarms(tableConfig, account, region)
		if err != nil {
			continue
		}
		resultMap = append(resultMap, result...)
	}
	return resultMap, nil
}
