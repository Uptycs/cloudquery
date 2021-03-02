/**
 * Copyright (c) 2020-present, The cloudquery authors
 *
 * This source code is licensed as defined by the LICENSE file found in the
 * root directory of this source tree.
 *
 * SPDX-License-Identifier: (Apache-2.0 OR GPL-2.0-only)
 */

package elb

import (
	"context"
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/Uptycs/cloudquery/utilities"

	"github.com/Uptycs/basequery-go/plugin/table"
	extaws "github.com/Uptycs/cloudquery/extension/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing"
)

// DescribeLoadBalancersColumns returns the list of columns in the table
func DescribeLoadBalancersColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("account_id"),
		table.TextColumn("region_code"),
		table.TextColumn("availability_zones"),
		table.TextColumn("backend_server_descriptions"),
		//table.BigIntColumn("backend_server_descriptions_instance_port"),
		//table.TextColumn("backend_server_descriptions_policy_names"),
		table.TextColumn("canonical_hosted_zone_name"),
		table.TextColumn("canonical_hosted_zone_name_id"),
		table.TextColumn("created_time"),
		//table.BigIntColumn("created_time_ext"),
		//table.TextColumn("created_time_loc"),
		//table.BigIntColumn("created_time_loc_cache_end"),
		//table.BigIntColumn("created_time_loc_cache_start"),
		//table.TextColumn("created_time_loc_cache_zone"),
		//table.TextColumn("created_time_loc_cache_zone_is_dst"),
		//table.TextColumn("created_time_loc_cache_zone_name"),
		//table.IntegerColumn("created_time_loc_cache_zone_offset"),
		//table.TextColumn("created_time_loc_name"),
		//table.TextColumn("created_time_loc_tx"),
		//table.IntegerColumn("created_time_loc_tx_index"),
		//table.TextColumn("created_time_loc_tx_isstd"),
		//table.TextColumn("created_time_loc_tx_isutc"),
		//table.BigIntColumn("created_time_loc_tx_when"),
		//table.TextColumn("created_time_loc_zone"),
		//table.TextColumn("created_time_loc_zone_is_dst"),
		//table.TextColumn("created_time_loc_zone_name"),
		//table.IntegerColumn("created_time_loc_zone_offset"),
		//table.BigIntColumn("created_time_wall"),
		table.TextColumn("dns_name"),
		table.TextColumn("health_check"),
		//table.BigIntColumn("health_check_healthy_threshold"),
		//table.BigIntColumn("health_check_interval"),
		//table.TextColumn("health_check_target"),
		//table.BigIntColumn("health_check_timeout"),
		//table.BigIntColumn("health_check_unhealthy_threshold"),
		table.TextColumn("instances"),
		//table.TextColumn("instances_instance_id"),
		table.TextColumn("listener_descriptions"),
		//table.TextColumn("listener_descriptions_listener"),
		//table.BigIntColumn("listener_descriptions_listener_instance_port"),
		//table.TextColumn("listener_descriptions_listener_instance_protocol"),
		//table.BigIntColumn("listener_descriptions_listener_load_balancer_port"),
		//table.TextColumn("listener_descriptions_listener_protocol"),
		//table.TextColumn("listener_descriptions_listener_ssl_certificate_id"),
		//table.TextColumn("listener_descriptions_policy_names"),
		table.TextColumn("load_balancer_name"),
		table.TextColumn("policies"),
		//table.TextColumn("policies_app_cookie_stickiness_policies"),
		//table.TextColumn("policies_app_cookie_stickiness_policies_cookie_name"),
		//table.TextColumn("policies_app_cookie_stickiness_policies_policy_name"),
		//table.TextColumn("policies_lb_cookie_stickiness_policies"),
		//table.BigIntColumn("policies_lb_cookie_stickiness_policies_cookie_expiration_period"),
		//table.TextColumn("policies_lb_cookie_stickiness_policies_policy_name"),
		//table.TextColumn("policies_other_policies"),
		table.TextColumn("scheme"),
		table.TextColumn("security_groups"),
		//table.TextColumn("source_security_group"),
		//table.TextColumn("source_security_group_group_name"),
		//table.TextColumn("source_security_group_owner_alias"),
		table.TextColumn("subnets"),
		table.TextColumn("vpc_id"),
	}
}

// DescribeLoadBalancersGenerate returns the rows in the table for all configured accounts
func DescribeLoadBalancersGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	if len(utilities.ExtConfiguration.ExtConfAws.Accounts) == 0 && extaws.ShouldProcessAccount("aws_elb_loadbalancer", utilities.AwsAccountID) {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_elb_loadbalancer",
			"account":   "default",
		}).Info("processing account")
		results, err := processAccountDescribeLoadBalancers(osqCtx, queryContext, nil)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, results...)
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfAws.Accounts {
			if !extaws.ShouldProcessAccount("aws_elb_loadbalancer", account.ID) {
				continue
			}
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_elb_loadbalancer",
				"account":   account.ID,
			}).Info("processing account")
			results, err := processAccountDescribeLoadBalancers(osqCtx, queryContext, &account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}

	return resultMap, nil
}

func processRegionDescribeLoadBalancers(osqCtx context.Context, queryContext table.QueryContext, tableConfig *utilities.TableConfig, account *utilities.ExtensionConfigurationAwsAccount, region types.Region) ([]map[string]string, error) {
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
		"tableName": "aws_elb_loadbalancer",
		"account":   accountId,
		"region":    *region.RegionName,
	}).Debug("processing region")

	svc := elasticloadbalancing.NewFromConfig(*sess)
	params := &elasticloadbalancing.DescribeLoadBalancersInput{}

	paginator := elasticloadbalancing.NewDescribeLoadBalancersPaginator(svc, params)

	for {
		page, err := paginator.NextPage(osqCtx)
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_elb_loadbalancer",
				"account":   accountId,
				"region":    *region.RegionName,
				"task":      "DescribeLoadBalancers",
				"errString": err.Error(),
			}).Error("failed to process region")
			return resultMap, err
		}
		byteArr, err := json.Marshal(page)
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_elb_loadbalancer",
				"account":   accountId,
				"region":    *region.RegionName,
				"task":      "DescribeLoadBalancers",
				"errString": err.Error(),
			}).Error("failed to marshal response")
			return nil, err
		}
		table := utilities.NewTable(byteArr, tableConfig)
		for _, row := range table.Rows {
			if !extaws.ShouldProcessRow(osqCtx, queryContext, "aws_elb_loadbalancer", accountId, *region.RegionName, row) {
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

func processAccountDescribeLoadBalancers(osqCtx context.Context, queryContext table.QueryContext, account *utilities.ExtensionConfigurationAwsAccount) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	awsSession, err := extaws.GetAwsConfig(account, "us-east-1")
	if err != nil {
		return resultMap, err
	}
	regions, err := extaws.FetchRegions(osqCtx, awsSession)
	if err != nil {
		return resultMap, err
	}
	tableConfig, ok := utilities.TableConfigurationMap["aws_elb_loadbalancer"]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_elb_loadbalancer",
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found")
	}
	for _, region := range regions {
		accountId := utilities.AwsAccountID
		if account != nil {
			accountId = account.ID
		}
		if !extaws.ShouldProcessRegion("aws_elb_loadbalancer", accountId, *region.RegionName) {
			continue
		}
		result, err := processRegionDescribeLoadBalancers(osqCtx, queryContext, tableConfig, account, region)
		if err != nil {
			continue
		}
		resultMap = append(resultMap, result...)
	}
	return resultMap, nil
}
