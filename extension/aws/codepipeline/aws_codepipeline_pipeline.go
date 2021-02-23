/**
 * Copyright (c) 2020-present, The cloudquery authors
 *
 * This source code is licensed as defined by the LICENSE file found in the
 * root directory of this source tree.
 *
 * SPDX-License-Identifier: (Apache-2.0 OR GPL-2.0-only)
 */

package codepipeline

import (
	"context"
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/Uptycs/cloudquery/utilities"

	"github.com/Uptycs/basequery-go/plugin/table"
	extaws "github.com/Uptycs/cloudquery/extension/aws"
	"github.com/aws/aws-sdk-go-v2/service/codepipeline"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

// ListPipelinesColumns returns the list of columns in the table
func ListPipelinesColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("account_id"),
		table.TextColumn("region_code"),
		table.TextColumn("created"),
		//table.BigIntColumn("created_ext"),
		//table.TextColumn("created_loc"),
		//table.BigIntColumn("created_loc_cache_end"),
		//table.BigIntColumn("created_loc_cache_start"),
		//table.TextColumn("created_loc_cache_zone"),
		//table.TextColumn("created_loc_cache_zone_is_dst"),
		//table.TextColumn("created_loc_cache_zone_name"),
		//table.IntegerColumn("created_loc_cache_zone_offset"),
		//table.TextColumn("created_loc_extend"),
		//table.TextColumn("created_loc_name"),
		//table.TextColumn("created_loc_tx"),
		//table.IntegerColumn("created_loc_tx_index"),
		//table.TextColumn("created_loc_tx_isstd"),
		//table.TextColumn("created_loc_tx_isutc"),
		//table.BigIntColumn("created_loc_tx_when"),
		//table.TextColumn("created_loc_zone"),
		//table.TextColumn("created_loc_zone_is_dst"),
		//table.TextColumn("created_loc_zone_name"),
		//table.IntegerColumn("created_loc_zone_offset"),
		//table.BigIntColumn("created_wall"),
		table.TextColumn("name"),
		table.TextColumn("updated"),
		//table.BigIntColumn("updated_ext"),
		//table.TextColumn("updated_loc"),
		//table.BigIntColumn("updated_loc_cache_end"),
		//table.BigIntColumn("updated_loc_cache_start"),
		//table.TextColumn("updated_loc_cache_zone"),
		//table.TextColumn("updated_loc_cache_zone_is_dst"),
		//table.TextColumn("updated_loc_cache_zone_name"),
		//table.IntegerColumn("updated_loc_cache_zone_offset"),
		//table.TextColumn("updated_loc_extend"),
		//table.TextColumn("updated_loc_name"),
		//table.TextColumn("updated_loc_tx"),
		//table.IntegerColumn("updated_loc_tx_index"),
		//table.TextColumn("updated_loc_tx_isstd"),
		//table.TextColumn("updated_loc_tx_isutc"),
		//table.BigIntColumn("updated_loc_tx_when"),
		//table.TextColumn("updated_loc_zone"),
		//table.TextColumn("updated_loc_zone_is_dst"),
		//table.TextColumn("updated_loc_zone_name"),
		//table.IntegerColumn("updated_loc_zone_offset"),
		//table.BigIntColumn("updated_wall"),
		table.IntegerColumn("version"),
		//table.TextColumn("values"),

	}
}

// ListPipelinesGenerate returns the rows in the table for all configured accounts
func ListPipelinesGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	if len(utilities.ExtConfiguration.ExtConfAws.Accounts) == 0 {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_codepipeline_pipeline",
			"account":   "default",
		}).Info("processing account")
		results, err := processAccountListPipelines(nil)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, results...)
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfAws.Accounts {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_codepipeline_pipeline",
				"account":   account.ID,
			}).Info("processing account")
			results, err := processAccountListPipelines(&account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}

	return resultMap, nil
}

func processRegionListPipelines(tableConfig *utilities.TableConfig, account *utilities.ExtensionConfigurationAwsAccount, region types.Region) ([]map[string]string, error) {
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
		"tableName": "aws_codepipeline_pipeline",
		"account":   accountId,
		"region":    *region.RegionName,
	}).Debug("processing region")

	svc := codepipeline.NewFromConfig(*sess)
	params := &codepipeline.ListPipelinesInput{}

	paginator := codepipeline.NewListPipelinesPaginator(svc, params)

	for {
		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_codepipeline_pipeline",
				"account":   accountId,
				"region":    *region.RegionName,
				"task":      "ListPipelines",
				"errString": err.Error(),
			}).Error("failed to process region")
			return resultMap, err
		}
		byteArr, err := json.Marshal(page)
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_codepipeline_pipeline",
				"account":   accountId,
				"region":    *region.RegionName,
				"task":      "ListPipelines",
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

func processAccountListPipelines(account *utilities.ExtensionConfigurationAwsAccount) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	awsSession, err := extaws.GetAwsConfig(account, "us-east-1")
	if err != nil {
		return resultMap, err
	}
	regions, err := extaws.FetchRegions(context.TODO(), awsSession)
	if err != nil {
		return resultMap, err
	}
	tableConfig, ok := utilities.TableConfigurationMap["aws_codepipeline_pipeline"]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_codepipeline_pipeline",
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found")
	}
	for _, region := range regions {
		result, err := processRegionListPipelines(tableConfig, account, region)
		if err != nil {
			continue
		}
		resultMap = append(resultMap, result...)
	}
	return resultMap, nil
}
