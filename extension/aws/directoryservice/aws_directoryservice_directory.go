/**
 * Copyright (c) 2020-present, The cloudquery authors
 *
 * This source code is licensed as defined by the LICENSE file found in the
 * root directory of this source tree.
 *
 * SPDX-License-Identifier: (Apache-2.0 OR GPL-2.0-only)
 */

package directoryservice

import (
	"context"
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/Uptycs/cloudquery/utilities"

	"github.com/Uptycs/basequery-go/plugin/table"
	extaws "github.com/Uptycs/cloudquery/extension/aws"
	"github.com/aws/aws-sdk-go-v2/service/directoryservice"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

// DescribeDirectoriesColumns returns the list of columns in the table
func DescribeDirectoriesColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("account_id"),
		table.TextColumn("region_code"),
		table.TextColumn("access_url"),
		table.TextColumn("alias"),
		table.TextColumn("connect_settings"),
		//table.TextColumn("connect_settings_availability_zones"),
		//table.TextColumn("connect_settings_connect_ips"),
		//table.TextColumn("connect_settings_customer_user_name"),
		//table.TextColumn("connect_settings_security_group_id"),
		//table.TextColumn("connect_settings_subnet_ids"),
		//table.TextColumn("connect_settings_vpc_id"),
		table.TextColumn("description"),
		table.IntegerColumn("desired_number_of_domain_controllers"),
		table.TextColumn("directory_id"),
		table.TextColumn("dns_ip_addrs"),
		table.TextColumn("edition"),
		table.TextColumn("launch_time"),
		//table.BigIntColumn("launch_time_ext"),
		//table.TextColumn("launch_time_loc"),
		//table.BigIntColumn("launch_time_loc_cache_end"),
		//table.BigIntColumn("launch_time_loc_cache_start"),
		//table.TextColumn("launch_time_loc_cache_zone"),
		//table.TextColumn("launch_time_loc_cache_zone_is_dst"),
		//table.TextColumn("launch_time_loc_cache_zone_name"),
		//table.IntegerColumn("launch_time_loc_cache_zone_offset"),
		//table.TextColumn("launch_time_loc_extend"),
		//table.TextColumn("launch_time_loc_name"),
		//table.TextColumn("launch_time_loc_tx"),
		//table.IntegerColumn("launch_time_loc_tx_index"),
		//table.TextColumn("launch_time_loc_tx_isstd"),
		//table.TextColumn("launch_time_loc_tx_isutc"),
		//table.BigIntColumn("launch_time_loc_tx_when"),
		//table.TextColumn("launch_time_loc_zone"),
		//table.TextColumn("launch_time_loc_zone_is_dst"),
		//table.TextColumn("launch_time_loc_zone_name"),
		//table.IntegerColumn("launch_time_loc_zone_offset"),
		//table.BigIntColumn("launch_time_wall"),
		table.TextColumn("name"),
		table.TextColumn("owner_directory_description"),
		//table.TextColumn("owner_directory_description_account_id"),
		//table.TextColumn("owner_directory_description_directory_id"),
		//table.TextColumn("owner_directory_description_dns_ip_addrs"),
		//table.TextColumn("owner_directory_description_radius_settings"),
		//table.TextColumn("owner_directory_description_radius_settings_authentication_protocol"),
		//table.TextColumn("owner_directory_description_radius_settings_display_label"),
		//table.IntegerColumn("owner_directory_description_radius_settings_radius_port"),
		//table.IntegerColumn("owner_directory_description_radius_settings_radius_retries"),
		//table.TextColumn("owner_directory_description_radius_settings_radius_servers"),
		//table.IntegerColumn("owner_directory_description_radius_settings_radius_timeout"),
		//table.TextColumn("owner_directory_description_radius_settings_shared_secret"),
		//table.TextColumn("owner_directory_description_radius_settings_use_same_username"),
		//table.TextColumn("owner_directory_description_radius_status"),
		//table.TextColumn("owner_directory_description_vpc_settings"),
		//table.TextColumn("owner_directory_description_vpc_settings_availability_zones"),
		//table.TextColumn("owner_directory_description_vpc_settings_security_group_id"),
		//table.TextColumn("owner_directory_description_vpc_settings_subnet_ids"),
		//table.TextColumn("owner_directory_description_vpc_settings_vpc_id"),
		table.TextColumn("radius_settings"),
		//table.TextColumn("radius_settings_authentication_protocol"),
		//table.TextColumn("radius_settings_display_label"),
		//table.IntegerColumn("radius_settings_radius_port"),
		//table.IntegerColumn("radius_settings_radius_retries"),
		//table.TextColumn("radius_settings_radius_servers"),
		//table.IntegerColumn("radius_settings_radius_timeout"),
		//table.TextColumn("radius_settings_shared_secret"),
		//table.TextColumn("radius_settings_use_same_username"),
		table.TextColumn("radius_status"),
		table.TextColumn("regions_info"),
		//table.TextColumn("regions_info_additional_regions"),
		//table.TextColumn("regions_info_primary_region"),
		table.TextColumn("share_method"),
		table.TextColumn("share_notes"),
		table.TextColumn("share_status"),
		table.TextColumn("short_name"),
		table.TextColumn("size"),
		table.TextColumn("sso_enabled"),
		table.TextColumn("stage"),
		//table.TextColumn("stage_last_updated_date_time"),
		//table.BigIntColumn("stage_last_updated_date_time_ext"),
		//table.TextColumn("stage_last_updated_date_time_loc"),
		//table.BigIntColumn("stage_last_updated_date_time_loc_cache_end"),
		//table.BigIntColumn("stage_last_updated_date_time_loc_cache_start"),
		//table.TextColumn("stage_last_updated_date_time_loc_cache_zone"),
		//table.TextColumn("stage_last_updated_date_time_loc_cache_zone_is_dst"),
		//table.TextColumn("stage_last_updated_date_time_loc_cache_zone_name"),
		//table.IntegerColumn("stage_last_updated_date_time_loc_cache_zone_offset"),
		//table.TextColumn("stage_last_updated_date_time_loc_extend"),
		//table.TextColumn("stage_last_updated_date_time_loc_name"),
		//table.TextColumn("stage_last_updated_date_time_loc_tx"),
		//table.IntegerColumn("stage_last_updated_date_time_loc_tx_index"),
		//table.TextColumn("stage_last_updated_date_time_loc_tx_isstd"),
		//table.TextColumn("stage_last_updated_date_time_loc_tx_isutc"),
		//table.BigIntColumn("stage_last_updated_date_time_loc_tx_when"),
		//table.TextColumn("stage_last_updated_date_time_loc_zone"),
		//table.TextColumn("stage_last_updated_date_time_loc_zone_is_dst"),
		//table.TextColumn("stage_last_updated_date_time_loc_zone_name"),
		//table.IntegerColumn("stage_last_updated_date_time_loc_zone_offset"),
		//table.BigIntColumn("stage_last_updated_date_time_wall"),
		//table.TextColumn("stage_reason"),
		table.TextColumn("type"),
		table.TextColumn("vpc_settings"),
		//table.TextColumn("vpc_settings_availability_zones"),
		//table.TextColumn("vpc_settings_security_group_id"),
		//table.TextColumn("vpc_settings_subnet_ids"),
		//table.TextColumn("vpc_settings_vpc_id"),
		//table.TextColumn("values"),

	}
}

// DescribeDirectoriesGenerate returns the rows in the table for all configured accounts
func DescribeDirectoriesGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	if len(utilities.ExtConfiguration.ExtConfAws.Accounts) == 0 {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_directoryservice_directory",
			"account":   "default",
		}).Info("processing account")
		results, err := processAccountDescribeDirectories(nil)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, results...)
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfAws.Accounts {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_directoryservice_directory",
				"account":   account.ID,
			}).Info("processing account")
			results, err := processAccountDescribeDirectories(&account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}

	return resultMap, nil
}

func processRegionDescribeDirectories(tableConfig *utilities.TableConfig, account *utilities.ExtensionConfigurationAwsAccount, region types.Region) ([]map[string]string, error) {
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
		"tableName": "aws_directoryservice_directory",
		"account":   accountId,
		"region":    *region.RegionName,
	}).Debug("processing region")

	svc := directoryservice.NewFromConfig(*sess)
	params := &directoryservice.DescribeDirectoriesInput{}

	result, err := svc.DescribeDirectories(context.TODO(), params)
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_directoryservice_directory",
			"account":   accountId,
			"region":    *region.RegionName,
			"task":      "DescribeDirectories",
			"errString": err.Error(),
		}).Error("failed to process region")
		return resultMap, err
	}

	byteArr, err := json.Marshal(result)
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_directoryservice_directory",
			"account":   accountId,
			"region":    *region.RegionName,
			"errString": err.Error(),
		}).Error("failed to marshal response")
		return resultMap, err
	}
	table := utilities.NewTable(byteArr, tableConfig)
	for _, row := range table.Rows {
		result := extaws.RowToMap(row, accountId, *region.RegionName, tableConfig)
		resultMap = append(resultMap, result)
	}
	return resultMap, nil
}

func processAccountDescribeDirectories(account *utilities.ExtensionConfigurationAwsAccount) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	awsSession, err := extaws.GetAwsConfig(account, "us-east-1")
	if err != nil {
		return resultMap, err
	}
	regions, err := extaws.FetchRegions(context.TODO(), awsSession)
	if err != nil {
		return resultMap, err
	}
	tableConfig, ok := utilities.TableConfigurationMap["aws_directoryservice_directory"]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_directoryservice_directory",
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found")
	}
	for _, region := range regions {
		result, err := processRegionDescribeDirectories(tableConfig, account, region)
		if err != nil {
			continue
		}
		resultMap = append(resultMap, result...)
	}
	return resultMap, nil
}
