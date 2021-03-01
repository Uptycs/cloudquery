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

// DescribeImagesColumns returns the list of columns in the table
func DescribeImagesColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("account_id"),
		table.TextColumn("region_code"),
		table.TextColumn("architecture"),
		table.TextColumn("owner_id"),
		table.TextColumn("platform"),
		table.TextColumn("usage_operation"),
		table.TextColumn("block_device_mappings"),
		//table.TextColumn("block_device_mappings_ebs"),
		//table.IntegerColumn("block_device_mappings_ebs_volume_size"),
		//table.TextColumn("block_device_mappings_ebs_volume_type"),
		//table.TextColumn("block_device_mappings_ebs_kms_key_id"),
		//table.IntegerColumn("block_device_mappings_ebs_throughput"),
		//table.TextColumn("block_device_mappings_ebs_encrypted"),
		//table.TextColumn("block_device_mappings_ebs_delete_on_termination"),
		//table.IntegerColumn("block_device_mappings_ebs_iops"),
		//table.TextColumn("block_device_mappings_ebs_snapshot_id"),
		//table.TextColumn("block_device_mappings_no_device"),
		//table.TextColumn("block_device_mappings_device_name"),
		//table.TextColumn("block_device_mappings_virtual_name"),
		table.TextColumn("root_device_type"),
		table.TextColumn("tags"),
		//table.TextColumn("tags_key"),
		//table.TextColumn("tags_value"),
		table.TextColumn("image_id"),
		table.TextColumn("image_type"),
		table.TextColumn("product_codes"),
		//table.TextColumn("product_codes_product_code_id"),
		//table.TextColumn("product_codes_product_code_type"),
		table.TextColumn("ramdisk_id"),
		table.TextColumn("root_device_name"),
		table.TextColumn("creation_date"),
		table.TextColumn("image_location"),
		table.TextColumn("kernel_id"),
		table.TextColumn("state"),
		table.TextColumn("hypervisor"),
		table.TextColumn("image_owner_alias"),
		table.TextColumn("name"),
		table.TextColumn("sriov_net_support"),
		table.TextColumn("state_reason"),
		//table.TextColumn("state_reason_code"),
		//table.TextColumn("state_reason_message"),
		table.TextColumn("public"),
		table.TextColumn("platform_details"),
		table.TextColumn("description"),
		table.TextColumn("ena_support"),
		table.TextColumn("virtualization_type"),
	}
}

// DescribeImagesGenerate returns the rows in the table for all configured accounts
func DescribeImagesGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	if len(utilities.ExtConfiguration.ExtConfAws.Accounts) == 0 && extaws.ShouldProcessAccount("aws_ec2_image", utilities.AwsAccountID) {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_ec2_image",
			"account":   "default",
		}).Info("processing account")
		results, err := processAccountDescribeImages(osqCtx, queryContext, nil)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, results...)
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfAws.Accounts {
			if !extaws.ShouldProcessAccount("aws_ec2_image", account.ID) {
				continue
			}
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_ec2_image",
				"account":   account.ID,
			}).Info("processing account")
			results, err := processAccountDescribeImages(osqCtx, queryContext, &account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}

	return resultMap, nil
}

func updateFilters(page *ec2.DescribeInstancesOutput, filters map[*string]bool) {
	for _, reservation := range page.Reservations {
		for _, instance := range reservation.Instances {
			filters[instance.ImageId] = true
		}
	}
}

func processDescribeImages(osqCtx context.Context, queryContext table.QueryContext, tableConfig *utilities.TableConfig, accountId string, svc *ec2.Client, region *types.Region, params *ec2.DescribeImagesInput) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	output, err := svc.DescribeImages(osqCtx, params)
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_ec2_image",
			"account":   accountId,
			"region":    *region.RegionName,
			"errString": err.Error(),
		}).Error("failed to get images")
		return resultMap, err
	}
	byteArr, err := json.Marshal(output)
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_ec2_image",
			"account":   accountId,
			"region":    *region.RegionName,
			"errString": err.Error(),
		}).Error("failed to marshal response")
		return resultMap, err
	}
	table := utilities.NewTable(byteArr, tableConfig)
	for _, row := range table.Rows {
		if !extaws.ShouldProcessRow(osqCtx, queryContext, "aws_ec2_image", accountId, *region.RegionName, row) {
			continue
		}
		result := extaws.RowToMap(row, accountId, *region.RegionName, tableConfig)
		resultMap = append(resultMap, result)
	}
	return resultMap, nil
}

func getImages(osqCtx context.Context, queryContext table.QueryContext, tableConfig *utilities.TableConfig, accountId string, svc *ec2.Client, region *types.Region, filters map[*string]bool) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	params := &ec2.DescribeImagesInput{}
	for key := range filters {
		params.ImageIds = append(params.ImageIds, *key)
		if len(params.ImageIds) >= 50 {
			result, err := processDescribeImages(osqCtx, queryContext, tableConfig, accountId, svc, region, params)
			if err != nil {
				return resultMap, err
			}
			resultMap = append(resultMap, result...)
			// reset params
			params = &ec2.DescribeImagesInput{}
		}
	}
	if len(params.ImageIds) > 0 {
		result, err := processDescribeImages(osqCtx, queryContext, tableConfig, accountId, svc, region, params)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, result...)
	}
	return resultMap, nil
}

func processRegionDescribeImages(osqCtx context.Context, queryContext table.QueryContext, tableConfig *utilities.TableConfig, account *utilities.ExtensionConfigurationAwsAccount, region types.Region) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	sess, err := extaws.GetAwsConfig(account, *region.RegionName)
	if err != nil {
		return resultMap, err
	}

	accountId := utilities.AwsAccountID
	if account != nil {
		accountId = account.ID
	}
	svc := ec2.NewFromConfig(*sess)
	params := &ec2.DescribeInstancesInput{}

	filters := make(map[*string]bool)
	paginator := ec2.NewDescribeInstancesPaginator(svc, params)

	for {
		page, err := paginator.NextPage(osqCtx)
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_ec2_image",
				"account":   accountId,
				"region":    *region.RegionName,
				"task":      "DescribeInstances",
				"errString": err.Error(),
			}).Error("failed to process region")
			return resultMap, err
		}
		updateFilters(page, filters)
		if !paginator.HasMorePages() {
			break
		}
	}

	resultMap, err = getImages(osqCtx, queryContext, tableConfig, accountId, svc, &region, filters)
	return resultMap, err
}

func processAccountDescribeImages(osqCtx context.Context, queryContext table.QueryContext, account *utilities.ExtensionConfigurationAwsAccount) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	awsSession, err := extaws.GetAwsConfig(account, "us-east-1")
	if err != nil {
		return resultMap, err
	}
	regions, err := extaws.FetchRegions(osqCtx, awsSession)
	if err != nil {
		return resultMap, err
	}
	tableConfig, ok := utilities.TableConfigurationMap["aws_ec2_image"]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_ec2_image",
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found")
	}
	for _, region := range regions {
		accountId := utilities.AwsAccountID
		if account != nil {
			accountId = account.ID
		}
		if !extaws.ShouldProcessRegion("aws_ec2_image", accountId, *region.RegionName) {
			continue
		}
		result, err := processRegionDescribeImages(osqCtx, queryContext, tableConfig, account, region)
		if err != nil {
			continue
		}
		resultMap = append(resultMap, result...)
	}
	return resultMap, nil
}
