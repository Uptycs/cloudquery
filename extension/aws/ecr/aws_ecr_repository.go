/**
 * Copyright (c) 2020-present, The cloudquery authors
 *
 * This source code is licensed as defined by the LICENSE file found in the
 * root directory of this source tree.
 *
 * SPDX-License-Identifier: (Apache-2.0 OR GPL-2.0-only)
 */

package ecr

import (
	"context"
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/Uptycs/cloudquery/utilities"

	"github.com/Uptycs/basequery-go/plugin/table"
	extaws "github.com/Uptycs/cloudquery/extension/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
)

// DescribeRepositoriesColumns returns the list of columns in the table
func DescribeRepositoriesColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("account_id"),
		table.TextColumn("region_code"),
		table.TextColumn("created_at"),
		//table.BigIntColumn("created_at_ext"),
		//table.TextColumn("created_at_loc"),
		//table.BigIntColumn("created_at_loc_cache_end"),
		//table.BigIntColumn("created_at_loc_cache_start"),
		//table.TextColumn("created_at_loc_cache_zone"),
		//table.TextColumn("created_at_loc_cache_zone_is_dst"),
		//table.TextColumn("created_at_loc_cache_zone_name"),
		//table.IntegerColumn("created_at_loc_cache_zone_offset"),
		//table.TextColumn("created_at_loc_extend"),
		//table.TextColumn("created_at_loc_name"),
		//table.TextColumn("created_at_loc_tx"),
		//table.IntegerColumn("created_at_loc_tx_index"),
		//table.TextColumn("created_at_loc_tx_isstd"),
		//table.TextColumn("created_at_loc_tx_isutc"),
		//table.BigIntColumn("created_at_loc_tx_when"),
		//table.TextColumn("created_at_loc_zone"),
		//table.TextColumn("created_at_loc_zone_is_dst"),
		//table.TextColumn("created_at_loc_zone_name"),
		//table.IntegerColumn("created_at_loc_zone_offset"),
		//table.BigIntColumn("created_at_wall"),
		table.TextColumn("encryption_configuration"),
		//table.TextColumn("encryption_configuration_encryption_type"),
		//table.TextColumn("encryption_configuration_kms_key"),
		table.TextColumn("image_scanning_configuration"),
		//table.TextColumn("image_scanning_configuration_scan_on_push"),
		table.TextColumn("image_tag_mutability"),
		table.TextColumn("registry_id"),
		table.TextColumn("repository_arn"),
		table.TextColumn("repository_name"),
		table.TextColumn("repository_uri"),
		//table.TextColumn("values"),

	}
}

// DescribeRepositoriesGenerate returns the rows in the table for all configured accounts
func DescribeRepositoriesGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	if len(utilities.ExtConfiguration.ExtConfAws.Accounts) == 0 && extaws.ShouldProcessAccount("aws_ecr_repository", utilities.AwsAccountID) {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_ecr_repository",
			"account":   "default",
		}).Info("processing account")
		results, err := processAccountDescribeRepositories(osqCtx, queryContext, nil)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, results...)
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfAws.Accounts {
			if !extaws.ShouldProcessAccount("aws_ecr_repository", account.ID) {
				continue
			}
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_ecr_repository",
				"account":   account.ID,
			}).Info("processing account")
			results, err := processAccountDescribeRepositories(osqCtx, queryContext, &account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}

	return resultMap, nil
}

func processRegionDescribeRepositories(osqCtx context.Context, queryContext table.QueryContext, tableConfig *utilities.TableConfig, account *utilities.ExtensionConfigurationAwsAccount, region types.Region) ([]map[string]string, error) {
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
		"tableName": "aws_ecr_repository",
		"account":   accountId,
		"region":    *region.RegionName,
	}).Debug("processing region")

	svc := ecr.NewFromConfig(*sess)
	params := &ecr.DescribeRepositoriesInput{}

	paginator := ecr.NewDescribeRepositoriesPaginator(svc, params)

	for {
		page, err := paginator.NextPage(osqCtx)
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_ecr_repository",
				"account":   accountId,
				"region":    *region.RegionName,
				"task":      "DescribeRepositories",
				"errString": err.Error(),
			}).Error("failed to process region")
			return resultMap, err
		}
		byteArr, err := json.Marshal(page)
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_ecr_repository",
				"account":   accountId,
				"region":    *region.RegionName,
				"task":      "DescribeRepositories",
				"errString": err.Error(),
			}).Error("failed to marshal response")
			return nil, err
		}
		table := utilities.NewTable(byteArr, tableConfig)
		for _, row := range table.Rows {
			if !extaws.ShouldProcessRow(osqCtx, queryContext, "aws_ecr_repository", accountId, *region.RegionName, row) {
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

func processAccountDescribeRepositories(osqCtx context.Context, queryContext table.QueryContext, account *utilities.ExtensionConfigurationAwsAccount) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	awsSession, err := extaws.GetAwsConfig(account, "us-east-1")
	if err != nil {
		return resultMap, err
	}
	regions, err := extaws.FetchRegions(osqCtx, awsSession)
	if err != nil {
		return resultMap, err
	}
	tableConfig, ok := utilities.TableConfigurationMap["aws_ecr_repository"]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_ecr_repository",
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found")
	}
	for _, region := range regions {
		accountId := utilities.AwsAccountID
		if account != nil {
			accountId = account.ID
		}
		if !extaws.ShouldProcessRegion("aws_ecr_repository", accountId, *region.RegionName) {
			continue
		}
		result, err := processRegionDescribeRepositories(osqCtx, queryContext, tableConfig, account, region)
		if err != nil {
			continue
		}
		resultMap = append(resultMap, result...)
	}
	return resultMap, nil
}
