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
	if len(utilities.ExtConfiguration.ExtConfAws.Accounts) == 0 {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_ec2_snapshot",
			"account":   "default",
		}).Info("processing account")
		results, err := processAccountDescribeSnapshots(nil)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, results...)
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfAws.Accounts {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_ec2_snapshot",
				"account":   account.ID,
			}).Info("processing account")
			results, err := processAccountDescribeSnapshots(&account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}

	return resultMap, nil
}

func processRegionDescribeSnapshots(tableConfig *utilities.TableConfig, account *utilities.ExtensionConfigurationAwsAccount, region *ec2.Region) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	sess, err := extaws.GetAwsSession(account, *region.RegionName)
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

	svc := ec2.New(sess)
	params := &ec2.DescribeSnapshotsInput{}

	err = svc.DescribeSnapshotsPages(params,
		func(page *ec2.DescribeSnapshotsOutput, lastPage bool) bool {
			byteArr, err := json.Marshal(page)
			if err != nil {
				utilities.GetLogger().WithFields(log.Fields{
					"tableName": "aws_ec2_snapshot",
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
			"tableName": "aws_ec2_snapshot",
			"account":   accountId,
			"region":    *region.RegionName,
			"task":      "DescribeSnapshots",
			"errString": err.Error(),
		}).Error("failed to process region")
		return resultMap, err
	}
	return resultMap, nil
}

func processAccountDescribeSnapshots(account *utilities.ExtensionConfigurationAwsAccount) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	awsSession, err := extaws.GetAwsSession(account, "us-east-1")
	if err != nil {
		return resultMap, err
	}
	regions, err := extaws.FetchRegions(awsSession)
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
		result, err := processRegionDescribeSnapshots(tableConfig, account, region)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, result...)
	}
	return resultMap, nil
}
