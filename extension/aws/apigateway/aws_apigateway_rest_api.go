/**
 * Copyright (c) 2020-present, The cloudquery authors
 *
 * This source code is licensed as defined by the LICENSE file found in the
 * root directory of this source tree.
 *
 * SPDX-License-Identifier: (Apache-2.0 OR GPL-2.0-only)
 */

package apigateway

import (
	"context"
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/Uptycs/cloudquery/utilities"

	"github.com/Uptycs/basequery-go/plugin/table"
	extaws "github.com/Uptycs/cloudquery/extension/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

// GetRestApisColumns returns the list of columns in the table
func GetRestApisColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("account_id"),
		table.TextColumn("region_code"),
		table.TextColumn("api_key_source"),
		table.TextColumn("binary_media_types"),
		table.TextColumn("created_date"),
		//table.BigIntColumn("created_date_ext"),
		//table.TextColumn("created_date_loc"),
		//table.BigIntColumn("created_date_loc_cache_end"),
		//table.BigIntColumn("created_date_loc_cache_start"),
		//table.TextColumn("created_date_loc_cache_zone"),
		//table.TextColumn("created_date_loc_cache_zone_is_dst"),
		//table.TextColumn("created_date_loc_cache_zone_name"),
		//table.IntegerColumn("created_date_loc_cache_zone_offset"),
		//table.TextColumn("created_date_loc_extend"),
		//table.TextColumn("created_date_loc_name"),
		//table.TextColumn("created_date_loc_tx"),
		//table.IntegerColumn("created_date_loc_tx_index"),
		//table.TextColumn("created_date_loc_tx_isstd"),
		//table.TextColumn("created_date_loc_tx_isutc"),
		//table.BigIntColumn("created_date_loc_tx_when"),
		//table.TextColumn("created_date_loc_zone"),
		//table.TextColumn("created_date_loc_zone_is_dst"),
		//table.TextColumn("created_date_loc_zone_name"),
		//table.IntegerColumn("created_date_loc_zone_offset"),
		//table.BigIntColumn("created_date_wall"),
		table.TextColumn("description"),
		table.TextColumn("disable_execute_api_endpoint"),
		table.TextColumn("endpoint_configuration"),
		//table.TextColumn("endpoint_configuration_types"),
		//table.TextColumn("endpoint_configuration_vpc_endpoint_ids"),
		table.TextColumn("id"),
		table.IntegerColumn("minimum_compression_size"),
		table.TextColumn("name"),
		table.TextColumn("policy"),
		table.TextColumn("tags"),
		table.TextColumn("version"),
		table.TextColumn("warnings"),
		//table.TextColumn("values"),

	}
}

// GetRestApisGenerate returns the rows in the table for all configured accounts
func GetRestApisGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	if len(utilities.ExtConfiguration.ExtConfAws.Accounts) == 0 {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_apigateway_rest_api",
			"account":   "default",
		}).Info("processing account")
		results, err := processAccountGetRestApis(nil)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, results...)
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfAws.Accounts {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_apigateway_rest_api",
				"account":   account.ID,
			}).Info("processing account")
			results, err := processAccountGetRestApis(&account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}

	return resultMap, nil
}

func processRegionGetRestApis(tableConfig *utilities.TableConfig, account *utilities.ExtensionConfigurationAwsAccount, region types.Region) ([]map[string]string, error) {
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
		"tableName": "aws_apigateway_rest_api",
		"account":   accountId,
		"region":    *region.RegionName,
	}).Debug("processing region")

	svc := apigateway.NewFromConfig(*sess)
	params := &apigateway.GetRestApisInput{}

	paginator := apigateway.NewGetRestApisPaginator(svc, params)

	for {
		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_apigateway_rest_api",
				"account":   accountId,
				"region":    *region.RegionName,
				"task":      "GetRestApis",
				"errString": err.Error(),
			}).Error("failed to process region")
			return resultMap, err
		}
		byteArr, err := json.Marshal(page)
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_apigateway_rest_api",
				"account":   accountId,
				"region":    *region.RegionName,
				"task":      "GetRestApis",
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

func processAccountGetRestApis(account *utilities.ExtensionConfigurationAwsAccount) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	awsSession, err := extaws.GetAwsConfig(account, "us-east-1")
	if err != nil {
		return resultMap, err
	}
	regions, err := extaws.FetchRegions(context.TODO(), awsSession)
	if err != nil {
		return resultMap, err
	}
	tableConfig, ok := utilities.TableConfigurationMap["aws_apigateway_rest_api"]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_apigateway_rest_api",
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found")
	}
	for _, region := range regions {
		result, err := processRegionGetRestApis(tableConfig, account, region)
		if err != nil {
			continue
		}
		resultMap = append(resultMap, result...)
	}
	return resultMap, nil
}
