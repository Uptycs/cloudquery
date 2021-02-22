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

// DescribeVolumesColumns returns the list of columns in the table
func DescribeVolumesColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("account_id"),
		table.TextColumn("region_code"),
		table.TextColumn("attachments"),
		//table.TextColumn("attachments_attach_time"),
		//table.BigIntColumn("attachments_attach_time_ext"),
		//table.TextColumn("attachments_attach_time_loc"),
		//table.BigIntColumn("attachments_attach_time_loc_cache_end"),
		//table.BigIntColumn("attachments_attach_time_loc_cache_start"),
		//table.TextColumn("attachments_attach_time_loc_cache_zone"),
		//table.TextColumn("attachments_attach_time_loc_cache_zone_is_dst"),
		//table.TextColumn("attachments_attach_time_loc_cache_zone_name"),
		//table.IntegerColumn("attachments_attach_time_loc_cache_zone_offset"),
		//table.TextColumn("attachments_attach_time_loc_name"),
		//table.TextColumn("attachments_attach_time_loc_tx"),
		//table.IntegerColumn("attachments_attach_time_loc_tx_index"),
		//table.TextColumn("attachments_attach_time_loc_tx_isstd"),
		//table.TextColumn("attachments_attach_time_loc_tx_isutc"),
		//table.BigIntColumn("attachments_attach_time_loc_tx_when"),
		//table.TextColumn("attachments_attach_time_loc_zone"),
		//table.TextColumn("attachments_attach_time_loc_zone_is_dst"),
		//table.TextColumn("attachments_attach_time_loc_zone_name"),
		//table.IntegerColumn("attachments_attach_time_loc_zone_offset"),
		//table.BigIntColumn("attachments_attach_time_wall"),
		//table.TextColumn("attachments_delete_on_termination"),
		//table.TextColumn("attachments_device"),
		//table.TextColumn("attachments_instance_id"),
		//table.TextColumn("attachments_state"),
		//table.TextColumn("attachments_volume_id"),
		table.TextColumn("availability_zone"),
		table.TextColumn("create_time"),
		//table.BigIntColumn("create_time_ext"),
		//table.TextColumn("create_time_loc"),
		//table.BigIntColumn("create_time_loc_cache_end"),
		//table.BigIntColumn("create_time_loc_cache_start"),
		//table.TextColumn("create_time_loc_cache_zone"),
		//table.TextColumn("create_time_loc_cache_zone_is_dst"),
		//table.TextColumn("create_time_loc_cache_zone_name"),
		//table.IntegerColumn("create_time_loc_cache_zone_offset"),
		//table.TextColumn("create_time_loc_name"),
		//table.TextColumn("create_time_loc_tx"),
		//table.IntegerColumn("create_time_loc_tx_index"),
		//table.TextColumn("create_time_loc_tx_isstd"),
		//table.TextColumn("create_time_loc_tx_isutc"),
		//table.BigIntColumn("create_time_loc_tx_when"),
		//table.TextColumn("create_time_loc_zone"),
		//table.TextColumn("create_time_loc_zone_is_dst"),
		//table.TextColumn("create_time_loc_zone_name"),
		//table.IntegerColumn("create_time_loc_zone_offset"),
		//table.BigIntColumn("create_time_wall"),
		table.TextColumn("encrypted"),
		table.TextColumn("fast_restored"),
		table.BigIntColumn("iops"),
		table.TextColumn("kms_key_id"),
		table.TextColumn("multi_attach_enabled"),
		table.TextColumn("outpost_arn"),
		table.BigIntColumn("size"),
		table.TextColumn("snapshot_id"),
		table.TextColumn("state"),
		table.TextColumn("tags"),
		//table.TextColumn("tags_key"),
		//table.TextColumn("tags_value"),
		table.BigIntColumn("throughput"),
		table.TextColumn("volume_id"),
		table.TextColumn("volume_type"),
	}
}

// DescribeVolumesGenerate returns the rows in the table for all configured accounts
func DescribeVolumesGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	if len(utilities.ExtConfiguration.ExtConfAws.Accounts) == 0 {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_ec2_volume",
			"account":   "default",
		}).Info("processing account")
		results, err := processAccountDescribeVolumes(nil)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, results...)
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfAws.Accounts {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_ec2_volume",
				"account":   account.ID,
			}).Info("processing account")
			results, err := processAccountDescribeVolumes(&account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}

	return resultMap, nil
}

func processRegionDescribeVolumes(tableConfig *utilities.TableConfig, account *utilities.ExtensionConfigurationAwsAccount, region types.Region) ([]map[string]string, error) {
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
		"tableName": "aws_ec2_volume",
		"account":   accountId,
		"region":    *region.RegionName,
	}).Debug("processing region")

	svc := ec2.NewFromConfig(*sess)
	params := &ec2.DescribeVolumesInput{}

	paginator := ec2.NewDescribeVolumesPaginator(svc, params)

	for {
		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_ec2_volume",
				"account":   accountId,
				"region":    *region.RegionName,
				"task":      "DescribeVolumes",
				"errString": err.Error(),
			}).Error("failed to process region")
			return resultMap, err
		}
		byteArr, err := json.Marshal(page)
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_ec2_volume",
				"account":   accountId,
				"region":    *region.RegionName,
				"task":      "DescribeVolumes",
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

func processAccountDescribeVolumes(account *utilities.ExtensionConfigurationAwsAccount) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	awsSession, err := extaws.GetAwsConfig(account, "us-east-1")
	if err != nil {
		return resultMap, err
	}
	regions, err := extaws.FetchRegions(context.TODO(), awsSession)
	if err != nil {
		return resultMap, err
	}
	tableConfig, ok := utilities.TableConfigurationMap["aws_ec2_volume"]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_ec2_volume",
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found")
	}
	for _, region := range regions {
		result, err := processRegionDescribeVolumes(tableConfig, account, region)
		if err != nil {
			continue
		}
		resultMap = append(resultMap, result...)
	}
	return resultMap, nil
}
