package ec2

import (
	"context"
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/Uptycs/cloudquery/utilities"

	extaws "github.com/Uptycs/cloudquery/extension/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/kolide/osquery-go/plugin/table"
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
	if len(utilities.ExtConfiguration.ExtConfAws.Accounts) == 0 {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_ec2_flowlog",
			"account":   "default",
		}).Info("processing account")
		results, err := processAccountDescribeFlowLogs(nil)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, results...)
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfAws.Accounts {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_ec2_flowlog",
				"account":   account.ID,
			}).Info("processing account")
			results, err := processAccountDescribeFlowLogs(&account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}

	return resultMap, nil
}

func processRegionDescribeFlowLogs(tableConfig *utilities.TableConfig, account *utilities.ExtensionConfigurationAwsAccount, region *ec2.Region) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	sess, err := extaws.GetAwsSession(account, *region.RegionName)
	if err != nil {
		return resultMap, err
	}

	accountId := utilities.AwsAccountId
	if account != nil {
		accountId = account.ID
	}

	utilities.GetLogger().WithFields(log.Fields{
		"tableName": "aws_ec2_flowlog",
		"account":   accountId,
		"region":    *region.RegionName,
	}).Debug("processing region")

	svc := ec2.New(sess)
	params := &ec2.DescribeFlowLogsInput{}

	err = svc.DescribeFlowLogsPages(params,
		func(page *ec2.DescribeFlowLogsOutput, lastPage bool) bool {
			byteArr, err := json.Marshal(page)
			if err != nil {
				utilities.GetLogger().WithFields(log.Fields{
					"tableName": "aws_ec2_flowlog",
					"account":   accountId,
					"region":    *region.RegionName,
					"errString": err.Error(),
				}).Error("failed to marshal response")
				return lastPage
			}
			table := utilities.NewTable(byteArr, tableConfig)
			for _, row := range table.Rows {
				result := extaws.RowToMap(row, accountId, *region.RegionName, tableConfig)
				resultMap = append(resultMap, result)
			}
			return lastPage
		})
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
	return resultMap, nil
}

func processAccountDescribeFlowLogs(account *utilities.ExtensionConfigurationAwsAccount) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	awsSession, err := extaws.GetAwsSession(account, "us-east-1")
	if err != nil {
		return resultMap, err
	}
	regions, err := extaws.FetchRegions(awsSession)
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
		result, err := processRegionDescribeFlowLogs(tableConfig, account, region)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, result...)
	}
	return resultMap, nil
}
