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
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/kolide/osquery-go/plugin/table"
)

// DescribeSecurityGroupsColumns returns the list of columns in the table
func DescribeSecurityGroupsColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("account_id"),
		table.TextColumn("region_code"),
		table.TextColumn("description"),
		table.TextColumn("group_id"),
		table.TextColumn("group_name"),
		table.TextColumn("ip_permissions"),
		//table.TextColumn("ip_permissions_egress"),
		//table.BigIntColumn("ip_permissions_egress_from_port"),
		//table.TextColumn("ip_permissions_egress_ip_protocol"),
		//table.TextColumn("ip_permissions_egress_ip_ranges"),
		//table.TextColumn("ip_permissions_egress_ip_ranges_cidr_ip"),
		//table.TextColumn("ip_permissions_egress_ip_ranges_description"),
		//table.TextColumn("ip_permissions_egress_ipv6_ranges"),
		//table.TextColumn("ip_permissions_egress_ipv6_ranges_cidr_ipv6"),
		//table.TextColumn("ip_permissions_egress_ipv6_ranges_description"),
		//table.TextColumn("ip_permissions_egress_prefix_list_ids"),
		//table.TextColumn("ip_permissions_egress_prefix_list_ids_description"),
		//table.TextColumn("ip_permissions_egress_prefix_list_ids_prefix_list_id"),
		//table.BigIntColumn("ip_permissions_egress_to_port"),
		//table.TextColumn("ip_permissions_egress_user_id_group_pairs"),
		//table.TextColumn("ip_permissions_egress_user_id_group_pairs_description"),
		//table.TextColumn("ip_permissions_egress_user_id_group_pairs_group_id"),
		//table.TextColumn("ip_permissions_egress_user_id_group_pairs_group_name"),
		//table.TextColumn("ip_permissions_egress_user_id_group_pairs_peering_status"),
		//table.TextColumn("ip_permissions_egress_user_id_group_pairs_user_id"),
		//table.TextColumn("ip_permissions_egress_user_id_group_pairs_vpc_id"),
		//table.TextColumn("ip_permissions_egress_user_id_group_pairs_vpc_peering_connection_id"),
		//table.BigIntColumn("ip_permissions_from_port"),
		//table.TextColumn("ip_permissions_ip_protocol"),
		//table.TextColumn("ip_permissions_ip_ranges"),
		//table.TextColumn("ip_permissions_ip_ranges_cidr_ip"),
		//table.TextColumn("ip_permissions_ip_ranges_description"),
		//table.TextColumn("ip_permissions_ipv6_ranges"),
		//table.TextColumn("ip_permissions_ipv6_ranges_cidr_ipv6"),
		//table.TextColumn("ip_permissions_ipv6_ranges_description"),
		//table.TextColumn("ip_permissions_prefix_list_ids"),
		//table.TextColumn("ip_permissions_prefix_list_ids_description"),
		//table.TextColumn("ip_permissions_prefix_list_ids_prefix_list_id"),
		//table.BigIntColumn("ip_permissions_to_port"),
		//table.TextColumn("ip_permissions_user_id_group_pairs"),
		//table.TextColumn("ip_permissions_user_id_group_pairs_description"),
		//table.TextColumn("ip_permissions_user_id_group_pairs_group_id"),
		//table.TextColumn("ip_permissions_user_id_group_pairs_group_name"),
		//table.TextColumn("ip_permissions_user_id_group_pairs_peering_status"),
		//table.TextColumn("ip_permissions_user_id_group_pairs_user_id"),
		//table.TextColumn("ip_permissions_user_id_group_pairs_vpc_id"),
		//table.TextColumn("ip_permissions_user_id_group_pairs_vpc_peering_connection_id"),
		table.TextColumn("owner_id"),
		table.TextColumn("tags"),
		//table.TextColumn("tags_key"),
		//table.TextColumn("tags_value"),
		table.TextColumn("vpc_id"),
	}
}

// DescribeSecurityGroupsGenerate returns the rows in the table for all configured accounts
func DescribeSecurityGroupsGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	if len(utilities.ExtConfiguration.ExtConfAws.Accounts) == 0 {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_ec2_security_group",
			"account":   "default",
		}).Info("processing account")
		results, err := processAccountDescribeSecurityGroups(nil)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, results...)
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfAws.Accounts {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_ec2_security_group",
				"account":   account.ID,
			}).Info("processing account")
			results, err := processAccountDescribeSecurityGroups(&account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}

	return resultMap, nil
}

func processRegionDescribeSecurityGroups(tableConfig *utilities.TableConfig, account *utilities.ExtensionConfigurationAwsAccount, region types.Region) ([]map[string]string, error) {
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
		"tableName": "aws_ec2_security_group",
		"account":   accountId,
		"region":    *region.RegionName,
	}).Debug("processing region")

	svc := ec2.NewFromConfig(*sess)
	params := &ec2.DescribeSecurityGroupsInput{}

	paginator := ec2.NewDescribeSecurityGroupsPaginator(svc, params)

	for {
		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_ec2_security_group",
				"account":   accountId,
				"region":    *region.RegionName,
				"task":      "DescribeSecurityGroups",
				"errString": err.Error(),
			}).Error("failed to process region")
			return resultMap, err
		}
		byteArr, err := json.Marshal(page)
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_ec2_security_group",
				"account":   accountId,
				"region":    *region.RegionName,
				"task":      "DescribeSecurityGroups",
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

func processAccountDescribeSecurityGroups(account *utilities.ExtensionConfigurationAwsAccount) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	awsSession, err := extaws.GetAwsConfig(account, "us-east-1")
	if err != nil {
		return resultMap, err
	}
	regions, err := extaws.FetchRegions(context.TODO(), awsSession)
	if err != nil {
		return resultMap, err
	}
	tableConfig, ok := utilities.TableConfigurationMap["aws_ec2_security_group"]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_ec2_security_group",
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found")
	}
	for _, region := range regions {
		result, err := processRegionDescribeSecurityGroups(tableConfig, account, region)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, result...)
	}
	return resultMap, nil
}
