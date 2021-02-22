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

// DescribeNatGatewaysColumns returns the list of columns in the table
func DescribeNatGatewaysColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("account_id"),
		table.TextColumn("region_code"),
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
		table.TextColumn("delete_time"),
		//table.BigIntColumn("delete_time_ext"),
		//table.TextColumn("delete_time_loc"),
		//table.BigIntColumn("delete_time_loc_cache_end"),
		//table.BigIntColumn("delete_time_loc_cache_start"),
		//table.TextColumn("delete_time_loc_cache_zone"),
		//table.TextColumn("delete_time_loc_cache_zone_is_dst"),
		//table.TextColumn("delete_time_loc_cache_zone_name"),
		//table.IntegerColumn("delete_time_loc_cache_zone_offset"),
		//table.TextColumn("delete_time_loc_name"),
		//table.TextColumn("delete_time_loc_tx"),
		//table.IntegerColumn("delete_time_loc_tx_index"),
		//table.TextColumn("delete_time_loc_tx_isstd"),
		//table.TextColumn("delete_time_loc_tx_isutc"),
		//table.BigIntColumn("delete_time_loc_tx_when"),
		//table.TextColumn("delete_time_loc_zone"),
		//table.TextColumn("delete_time_loc_zone_is_dst"),
		//table.TextColumn("delete_time_loc_zone_name"),
		//table.IntegerColumn("delete_time_loc_zone_offset"),
		//table.BigIntColumn("delete_time_wall"),
		table.TextColumn("failure_code"),
		table.TextColumn("failure_message"),
		table.TextColumn("nat_gateway_addresses"),
		//table.TextColumn("nat_gateway_addresses_allocation_id"),
		//table.TextColumn("nat_gateway_addresses_network_interface_id"),
		//table.TextColumn("nat_gateway_addresses_private_ip"),
		//table.TextColumn("nat_gateway_addresses_public_ip"),
		table.TextColumn("nat_gateway_id"),
		table.TextColumn("provisioned_bandwidth"),
		//table.TextColumn("provisioned_bandwidth_provision_time"),
		//table.BigIntColumn("provisioned_bandwidth_provision_time_ext"),
		//table.TextColumn("provisioned_bandwidth_provision_time_loc"),
		//table.BigIntColumn("provisioned_bandwidth_provision_time_loc_cache_end"),
		//table.BigIntColumn("provisioned_bandwidth_provision_time_loc_cache_start"),
		//table.TextColumn("provisioned_bandwidth_provision_time_loc_cache_zone"),
		//table.TextColumn("provisioned_bandwidth_provision_time_loc_cache_zone_is_dst"),
		//table.TextColumn("provisioned_bandwidth_provision_time_loc_cache_zone_name"),
		//table.IntegerColumn("provisioned_bandwidth_provision_time_loc_cache_zone_offset"),
		//table.TextColumn("provisioned_bandwidth_provision_time_loc_name"),
		//table.TextColumn("provisioned_bandwidth_provision_time_loc_tx"),
		//table.IntegerColumn("provisioned_bandwidth_provision_time_loc_tx_index"),
		//table.TextColumn("provisioned_bandwidth_provision_time_loc_tx_isstd"),
		//table.TextColumn("provisioned_bandwidth_provision_time_loc_tx_isutc"),
		//table.BigIntColumn("provisioned_bandwidth_provision_time_loc_tx_when"),
		//table.TextColumn("provisioned_bandwidth_provision_time_loc_zone"),
		//table.TextColumn("provisioned_bandwidth_provision_time_loc_zone_is_dst"),
		//table.TextColumn("provisioned_bandwidth_provision_time_loc_zone_name"),
		//table.IntegerColumn("provisioned_bandwidth_provision_time_loc_zone_offset"),
		//table.BigIntColumn("provisioned_bandwidth_provision_time_wall"),
		//table.TextColumn("provisioned_bandwidth_provisioned"),
		//table.TextColumn("provisioned_bandwidth_request_time"),
		//table.BigIntColumn("provisioned_bandwidth_request_time_ext"),
		//table.TextColumn("provisioned_bandwidth_request_time_loc"),
		//table.BigIntColumn("provisioned_bandwidth_request_time_loc_cache_end"),
		//table.BigIntColumn("provisioned_bandwidth_request_time_loc_cache_start"),
		//table.TextColumn("provisioned_bandwidth_request_time_loc_cache_zone"),
		//table.TextColumn("provisioned_bandwidth_request_time_loc_cache_zone_is_dst"),
		//table.TextColumn("provisioned_bandwidth_request_time_loc_cache_zone_name"),
		//table.IntegerColumn("provisioned_bandwidth_request_time_loc_cache_zone_offset"),
		//table.TextColumn("provisioned_bandwidth_request_time_loc_name"),
		//table.TextColumn("provisioned_bandwidth_request_time_loc_tx"),
		//table.IntegerColumn("provisioned_bandwidth_request_time_loc_tx_index"),
		//table.TextColumn("provisioned_bandwidth_request_time_loc_tx_isstd"),
		//table.TextColumn("provisioned_bandwidth_request_time_loc_tx_isutc"),
		//table.BigIntColumn("provisioned_bandwidth_request_time_loc_tx_when"),
		//table.TextColumn("provisioned_bandwidth_request_time_loc_zone"),
		//table.TextColumn("provisioned_bandwidth_request_time_loc_zone_is_dst"),
		//table.TextColumn("provisioned_bandwidth_request_time_loc_zone_name"),
		//table.IntegerColumn("provisioned_bandwidth_request_time_loc_zone_offset"),
		//table.BigIntColumn("provisioned_bandwidth_request_time_wall"),
		//table.TextColumn("provisioned_bandwidth_requested"),
		//table.TextColumn("provisioned_bandwidth_status"),
		table.TextColumn("state"),
		table.TextColumn("subnet_id"),
		table.TextColumn("tags"),
		//table.TextColumn("tags_key"),
		//table.TextColumn("tags_value"),
		table.TextColumn("vpc_id"),
	}
}

// DescribeNatGatewaysGenerate returns the rows in the table for all configured accounts
func DescribeNatGatewaysGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	if len(utilities.ExtConfiguration.ExtConfAws.Accounts) == 0 {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_ec2_nat_gateway",
			"account":   "default",
		}).Info("processing account")
		results, err := processAccountDescribeNatGateways(nil)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, results...)
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfAws.Accounts {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_ec2_nat_gateway",
				"account":   account.ID,
			}).Info("processing account")
			results, err := processAccountDescribeNatGateways(&account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}

	return resultMap, nil
}

func processRegionDescribeNatGateways(tableConfig *utilities.TableConfig, account *utilities.ExtensionConfigurationAwsAccount, region types.Region) ([]map[string]string, error) {
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
		"tableName": "aws_ec2_nat_gateway",
		"account":   accountId,
		"region":    *region.RegionName,
	}).Debug("processing region")

	svc := ec2.NewFromConfig(*sess)
	params := &ec2.DescribeNatGatewaysInput{}

	paginator := ec2.NewDescribeNatGatewaysPaginator(svc, params)

	for {
		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_ec2_nat_gateway",
				"account":   accountId,
				"region":    *region.RegionName,
				"task":      "DescribeNatGateways",
				"errString": err.Error(),
			}).Error("failed to process region")
			return resultMap, err
		}
		byteArr, err := json.Marshal(page)
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_ec2_nat_gateway",
				"account":   accountId,
				"region":    *region.RegionName,
				"task":      "DescribeNatGateways",
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

func processAccountDescribeNatGateways(account *utilities.ExtensionConfigurationAwsAccount) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	awsSession, err := extaws.GetAwsConfig(account, "us-east-1")
	if err != nil {
		return resultMap, err
	}
	regions, err := extaws.FetchRegions(context.TODO(), awsSession)
	if err != nil {
		return resultMap, err
	}
	tableConfig, ok := utilities.TableConfigurationMap["aws_ec2_nat_gateway"]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_ec2_nat_gateway",
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found")
	}
	for _, region := range regions {
		result, err := processRegionDescribeNatGateways(tableConfig, account, region)
		if err != nil {
			continue
		}
		resultMap = append(resultMap, result...)
	}
	return resultMap, nil
}
