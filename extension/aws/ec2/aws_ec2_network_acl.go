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

// DescribeNetworkAclsColumns returns the list of columns in the table
func DescribeNetworkAclsColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("account_id"),
		table.TextColumn("region_code"),
		table.TextColumn("associations"),
		//table.TextColumn("associations_network_acl_association_id"),
		//table.TextColumn("associations_network_acl_id"),
		//table.TextColumn("associations_subnet_id"),
		table.TextColumn("entries"),
		//table.TextColumn("entries_cidr_block"),
		//table.TextColumn("entries_egress"),
		//table.TextColumn("entries_icmp_type_code"),
		//table.BigIntColumn("entries_icmp_type_code_code"),
		//table.BigIntColumn("entries_icmp_type_code_type"),
		//table.TextColumn("entries_ipv6_cidr_block"),
		//table.TextColumn("entries_port_range"),
		//table.BigIntColumn("entries_port_range_from"),
		//table.BigIntColumn("entries_port_range_to"),
		//table.TextColumn("entries_protocol"),
		//table.TextColumn("entries_rule_action"),
		//table.BigIntColumn("entries_rule_number"),
		table.TextColumn("is_default"),
		table.TextColumn("network_acl_id"),
		table.TextColumn("owner_id"),
		table.TextColumn("tags"),
		//table.TextColumn("tags_key"),
		//table.TextColumn("tags_value"),
		table.TextColumn("vpc_id"),
	}
}

// DescribeNetworkAclsGenerate returns the rows in the table for all configured accounts
func DescribeNetworkAclsGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	if len(utilities.ExtConfiguration.ExtConfAws.Accounts) == 0 {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_ec2_network_acl",
			"account":   "default",
		}).Info("processing account")
		results, err := processAccountDescribeNetworkAcls(nil)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, results...)
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfAws.Accounts {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_ec2_network_acl",
				"account":   account.ID,
			}).Info("processing account")
			results, err := processAccountDescribeNetworkAcls(&account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}

	return resultMap, nil
}

func processRegionDescribeNetworkAcls(tableConfig *utilities.TableConfig, account *utilities.ExtensionConfigurationAwsAccount, region types.Region) ([]map[string]string, error) {
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
		"tableName": "aws_ec2_network_acl",
		"account":   accountId,
		"region":    *region.RegionName,
	}).Debug("processing region")

	svc := ec2.NewFromConfig(*sess)
	params := &ec2.DescribeNetworkAclsInput{}

	paginator := ec2.NewDescribeNetworkAclsPaginator(svc, params)

	for {
		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_ec2_network_acl",
				"account":   accountId,
				"region":    *region.RegionName,
				"task":      "DescribeNetworkAcls",
				"errString": err.Error(),
			}).Error("failed to process region")
			return resultMap, err
		}
		byteArr, err := json.Marshal(page)
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_ec2_network_acl",
				"account":   accountId,
				"region":    *region.RegionName,
				"task":      "DescribeNetworkAcls",
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

func processAccountDescribeNetworkAcls(account *utilities.ExtensionConfigurationAwsAccount) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	awsSession, err := extaws.GetAwsConfig(account, "us-east-1")
	if err != nil {
		return resultMap, err
	}
	regions, err := extaws.FetchRegions(context.TODO(), awsSession)
	if err != nil {
		return resultMap, err
	}
	tableConfig, ok := utilities.TableConfigurationMap["aws_ec2_network_acl"]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_ec2_network_acl",
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found")
	}
	for _, region := range regions {
		result, err := processRegionDescribeNetworkAcls(tableConfig, account, region)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, result...)
	}
	return resultMap, nil
}
