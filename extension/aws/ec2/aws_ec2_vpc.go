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

	extaws "github.com/Uptycs/cloudquery/extension/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/kolide/osquery-go/plugin/table"
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
	if len(utilities.ExtConfiguration.ExtConfAws.Accounts) == 0 {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_ec2_vpc",
			"account":   "default",
		}).Info("processing account")
		results, err := processAccountDescribeVpcs(nil)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, results...)
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfAws.Accounts {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_ec2_vpc",
				"account":   account.ID,
			}).Info("processing account")
			results, err := processAccountDescribeVpcs(&account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}

	return resultMap, nil
}

func processRegionDescribeVpcs(tableConfig *utilities.TableConfig, account *utilities.ExtensionConfigurationAwsAccount, region *ec2.Region) ([]map[string]string, error) {
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
		"tableName": "aws_ec2_vpc",
		"account":   accountId,
		"region":    *region.RegionName,
	}).Debug("processing region")

	svc := ec2.New(sess)
	params := &ec2.DescribeVpcsInput{}

	err = svc.DescribeVpcsPages(params,
		func(page *ec2.DescribeVpcsOutput, lastPage bool) bool {
			byteArr, err := json.Marshal(page)
			if err != nil {
				utilities.GetLogger().WithFields(log.Fields{
					"tableName": "aws_ec2_vpc",
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
			"tableName": "aws_ec2_vpc",
			"account":   accountId,
			"region":    *region.RegionName,
			"task":      "DescribeVpcs",
			"errString": err.Error(),
		}).Error("failed to process region")
		return resultMap, err
	}
	return resultMap, nil
}

func processAccountDescribeVpcs(account *utilities.ExtensionConfigurationAwsAccount) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	awsSession, err := extaws.GetAwsSession(account, "us-east-1")
	if err != nil {
		return resultMap, err
	}
	regions, err := extaws.FetchRegions(awsSession)
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
		result, err := processRegionDescribeVpcs(tableConfig, account, region)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, result...)
	}
	return resultMap, nil
}
