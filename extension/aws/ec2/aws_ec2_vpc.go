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

// DescribeVpcsColumns returns the list of columns in the table
func DescribeVpcsColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("account_id"),
		table.TextColumn("region_code"),
		table.TextColumn("cidr_block"),
		table.TextColumn("cidr_block_association_set"),
		//table.TextColumn("cidr_block_association_set_association_id"),
		//table.TextColumn("cidr_block_association_set_cidr_block"),
		//table.TextColumn("cidr_block_association_set_cidr_block_state"),
		//table.TextColumn("cidr_block_association_set_cidr_block_state_state"),
		//table.TextColumn("cidr_block_association_set_cidr_block_state_status_message"),
		table.TextColumn("dhcp_options_id"),
		table.TextColumn("instance_tenancy"),
		table.TextColumn("ipv6_cidr_block_association_set"),
		//table.TextColumn("ipv6_cidr_block_association_set_association_id"),
		//table.TextColumn("ipv6_cidr_block_association_set_ipv6_cidr_block"),
		//table.TextColumn("ipv6_cidr_block_association_set_ipv6_cidr_block_state"),
		//table.TextColumn("ipv6_cidr_block_association_set_ipv6_cidr_block_state_state"),
		//table.TextColumn("ipv6_cidr_block_association_set_ipv6_cidr_block_state_status_message"),
		//table.TextColumn("ipv6_cidr_block_association_set_ipv6_pool"),
		//table.TextColumn("ipv6_cidr_block_association_set_network_border_group"),
		table.TextColumn("is_default"),
		table.TextColumn("owner_id"),
		table.TextColumn("state"),
		table.TextColumn("tags"),
		table.TextColumn("tags_key"),
		table.TextColumn("tags_value"),
		table.TextColumn("vpc_id"),
	}
}

// DescribeVpcsGenerate returns the rows in the table for all configured accounts
func DescribeVpcsGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	if len(utilities.ExtConfiguration.ExtConfAws.Accounts) == 0 && extaws.ShouldProcessAccount("aws_ec2_vpc", utilities.AwsAccountID) {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_ec2_vpc",
			"account":   "default",
		}).Info("processing account")
		results, err := processAccountDescribeVpcs(osqCtx, queryContext, nil)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, results...)
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfAws.Accounts {
			if !extaws.ShouldProcessAccount("aws_ec2_vpc", account.ID) {
				continue
			}
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_ec2_vpc",
				"account":   account.ID,
			}).Info("processing account")
			results, err := processAccountDescribeVpcs(osqCtx, queryContext, &account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}

	return resultMap, nil
}

func processRegionDescribeVpcs(osqCtx context.Context, queryContext table.QueryContext, tableConfig *utilities.TableConfig, account *utilities.ExtensionConfigurationAwsAccount, region types.Region) ([]map[string]string, error) {
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
		"tableName": "aws_ec2_vpc",
		"account":   accountId,
		"region":    *region.RegionName,
	}).Debug("processing region")

	svc := ec2.NewFromConfig(*sess)
	params := &ec2.DescribeVpcsInput{}

	paginator := ec2.NewDescribeVpcsPaginator(svc, params)

	for {
		page, err := paginator.NextPage(osqCtx)
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_ec2_vpc",
				"account":   accountId,
				"region":    *region.RegionName,
				"task":      "DescribeVpcs",
				"errString": err.Error(),
			}).Error("failed to process region")
			return resultMap, err
		}
		byteArr, err := json.Marshal(page)
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_ec2_vpc",
				"account":   accountId,
				"region":    *region.RegionName,
				"task":      "DescribeVpcs",
				"errString": err.Error(),
			}).Error("failed to marshal response")
			return nil, err
		}
		table := utilities.NewTable(byteArr, tableConfig)
		for _, row := range table.Rows {
			if !extaws.ShouldProcessRow(osqCtx, queryContext, "aws_ec2_vpc", accountId, *region.RegionName, row) {
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

func processAccountDescribeVpcs(osqCtx context.Context, queryContext table.QueryContext, account *utilities.ExtensionConfigurationAwsAccount) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	awsSession, err := extaws.GetAwsConfig(account, "us-east-1")
	if err != nil {
		return resultMap, err
	}
	regions, err := extaws.FetchRegions(osqCtx, awsSession)
	if err != nil {
		return resultMap, err
	}
	tableConfig, ok := utilities.TableConfigurationMap["aws_ec2_vpc"]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_ec2_vpc",
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found")
	}
	for _, region := range regions {
		accountId := utilities.AwsAccountID
		if account != nil {
			accountId = account.ID
		}
		if !extaws.ShouldProcessRegion("aws_ec2_vpc", accountId, *region.RegionName) {
			continue
		}
		result, err := processRegionDescribeVpcs(osqCtx, queryContext, tableConfig, account, region)
		if err != nil {
			continue
		}
		resultMap = append(resultMap, result...)
	}
	return resultMap, nil
}
