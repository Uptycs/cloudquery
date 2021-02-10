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

// DescribeRouteTablesColumns returns the list of columns in the table
func DescribeRouteTablesColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("account_id"),
		table.TextColumn("region_code"),
		table.TextColumn("associations"),
		//table.TextColumn("associations_association_state"),
		//table.TextColumn("associations_association_state_state"),
		//table.TextColumn("associations_association_state_status_message"),
		//table.TextColumn("associations_gateway_id"),
		//table.TextColumn("associations_main"),
		//table.TextColumn("associations_route_table_association_id"),
		//table.TextColumn("associations_route_table_id"),
		//table.TextColumn("associations_subnet_id"),
		table.TextColumn("owner_id"),
		table.TextColumn("propagating_vgws"),
		//table.TextColumn("propagating_vgws_gateway_id"),
		table.TextColumn("route_table_id"),
		table.TextColumn("routes"),
		//table.TextColumn("routes_carrier_gateway_id"),
		//table.TextColumn("routes_destination_cidr_block"),
		//table.TextColumn("routes_destination_ipv6_cidr_block"),
		//table.TextColumn("routes_destination_prefix_list_id"),
		//table.TextColumn("routes_egress_only_internet_gateway_id"),
		//table.TextColumn("routes_gateway_id"),
		//table.TextColumn("routes_instance_id"),
		//table.TextColumn("routes_instance_owner_id"),
		//table.TextColumn("routes_local_gateway_id"),
		//table.TextColumn("routes_nat_gateway_id"),
		//table.TextColumn("routes_network_interface_id"),
		//table.TextColumn("routes_origin"),
		//table.TextColumn("routes_state"),
		//table.TextColumn("routes_transit_gateway_id"),
		//table.TextColumn("routes_vpc_peering_connection_id"),
		table.TextColumn("tags"),
		//table.TextColumn("tags_key"),
		//table.TextColumn("tags_value"),
		table.TextColumn("vpc_id"),
	}
}

// DescribeRouteTablesGenerate returns the rows in the table for all configured accounts
func DescribeRouteTablesGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	if len(utilities.ExtConfiguration.ExtConfAws.Accounts) == 0 {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_ec2_route_table",
			"account":   "default",
		}).Info("processing account")
		results, err := processAccountDescribeRouteTables(nil)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, results...)
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfAws.Accounts {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_ec2_route_table",
				"account":   account.ID,
			}).Info("processing account")
			results, err := processAccountDescribeRouteTables(&account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}

	return resultMap, nil
}

func processRegionDescribeRouteTables(tableConfig *utilities.TableConfig, account *utilities.ExtensionConfigurationAwsAccount, region types.Region) ([]map[string]string, error) {
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
		"tableName": "aws_ec2_route_table",
		"account":   accountId,
		"region":    *region.RegionName,
	}).Debug("processing region")

	svc := ec2.NewFromConfig(*sess)
	params := &ec2.DescribeRouteTablesInput{}

	paginator := ec2.NewDescribeRouteTablesPaginator(svc, params)

	for {
		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_ec2_route_table",
				"account":   accountId,
				"region":    *region.RegionName,
				"task":      "DescribeRouteTables",
				"errString": err.Error(),
			}).Error("failed to process region")
			return resultMap, err
		}
		byteArr, err := json.Marshal(page)
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_ec2_route_table",
				"account":   accountId,
				"region":    *region.RegionName,
				"task":      "DescribeRouteTables",
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

func processAccountDescribeRouteTables(account *utilities.ExtensionConfigurationAwsAccount) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	awsSession, err := extaws.GetAwsConfig(account, "us-east-1")
	if err != nil {
		return resultMap, err
	}
	regions, err := extaws.FetchRegions(context.TODO(), awsSession)
	if err != nil {
		return resultMap, err
	}
	tableConfig, ok := utilities.TableConfigurationMap["aws_ec2_route_table"]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_ec2_route_table",
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found")
	}
	for _, region := range regions {
		result, err := processRegionDescribeRouteTables(tableConfig, account, region)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, result...)
	}
	return resultMap, nil
}
