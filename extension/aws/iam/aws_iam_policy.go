/**
 * Copyright (c) 2020-present, The cloudquery authors
 *
 * This source code is licensed as defined by the LICENSE file found in the
 * root directory of this source tree.
 *
 * SPDX-License-Identifier: (Apache-2.0 OR GPL-2.0-only)
 */

package iam

import (
	"context"
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/Uptycs/cloudquery/utilities"

	"github.com/Uptycs/basequery-go/plugin/table"
	extaws "github.com/Uptycs/cloudquery/extension/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
)

// ListPoliciesColumns returns the list of columns in the table
func ListPoliciesColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("account_id"),
		table.TextColumn("arn"),
		table.BigIntColumn("attachment_count"),
		table.TextColumn("create_date"),
		//table.BigIntColumn("create_date_ext"),
		//table.TextColumn("create_date_loc"),
		//table.BigIntColumn("create_date_loc_cache_end"),
		//table.BigIntColumn("create_date_loc_cache_start"),
		//table.TextColumn("create_date_loc_cache_zone"),
		//table.TextColumn("create_date_loc_cache_zone_is_dst"),
		//table.TextColumn("create_date_loc_cache_zone_name"),
		//table.IntegerColumn("create_date_loc_cache_zone_offset"),
		//table.TextColumn("create_date_loc_name"),
		//table.TextColumn("create_date_loc_tx"),
		//table.IntegerColumn("create_date_loc_tx_index"),
		//table.TextColumn("create_date_loc_tx_isstd"),
		//table.TextColumn("create_date_loc_tx_isutc"),
		//table.BigIntColumn("create_date_loc_tx_when"),
		//table.TextColumn("create_date_loc_zone"),
		//table.TextColumn("create_date_loc_zone_is_dst"),
		//table.TextColumn("create_date_loc_zone_name"),
		//table.IntegerColumn("create_date_loc_zone_offset"),
		//table.BigIntColumn("create_date_wall"),
		table.TextColumn("default_version_id"),
		table.TextColumn("description"),
		table.TextColumn("is_attachable"),
		table.TextColumn("path"),
		table.BigIntColumn("permissions_boundary_usage_count"),
		table.TextColumn("policy_id"),
		table.TextColumn("policy_name"),
		table.TextColumn("update_date"),
		//table.BigIntColumn("update_date_ext"),
		//table.TextColumn("update_date_loc"),
		//table.BigIntColumn("update_date_loc_cache_end"),
		//table.BigIntColumn("update_date_loc_cache_start"),
		//table.TextColumn("update_date_loc_cache_zone"),
		//table.TextColumn("update_date_loc_cache_zone_is_dst"),
		//table.TextColumn("update_date_loc_cache_zone_name"),
		//table.IntegerColumn("update_date_loc_cache_zone_offset"),
		//table.TextColumn("update_date_loc_name"),
		//table.TextColumn("update_date_loc_tx"),
		//table.IntegerColumn("update_date_loc_tx_index"),
		//table.TextColumn("update_date_loc_tx_isstd"),
		//table.TextColumn("update_date_loc_tx_isutc"),
		//table.BigIntColumn("update_date_loc_tx_when"),
		//table.TextColumn("update_date_loc_zone"),
		//table.TextColumn("update_date_loc_zone_is_dst"),
		//table.TextColumn("update_date_loc_zone_name"),
		//table.IntegerColumn("update_date_loc_zone_offset"),
		//table.BigIntColumn("update_date_wall"),

	}
}

// ListPoliciesGenerate returns the rows in the table for all configured accounts
func ListPoliciesGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	if len(utilities.ExtConfiguration.ExtConfAws.Accounts) == 0 {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_iam_policy",
			"account":   "default",
		}).Info("processing account")
		results, err := processAccountListPolicies(nil)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, results...)
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfAws.Accounts {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_iam_policy",
				"account":   account.ID,
			}).Info("processing account")
			results, err := processAccountListPolicies(&account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}

	return resultMap, nil
}

func processGlobalListPolicies(tableConfig *utilities.TableConfig, account *utilities.ExtensionConfigurationAwsAccount) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	sess, err := extaws.GetAwsConfig(account, "aws-global")
	if err != nil {
		return resultMap, err
	}

	accountId := utilities.AwsAccountID
	if account != nil {
		accountId = account.ID
	}

	utilities.GetLogger().WithFields(log.Fields{
		"tableName": "aws_iam_policy",
		"account":   accountId,
		"region":    "aws-global",
	}).Debug("processing region")

	svc := iam.NewFromConfig(*sess)
	params := &iam.ListPoliciesInput{}

	paginator := iam.NewListPoliciesPaginator(svc, params)

	for {
		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_iam_policy",
				"account":   accountId,
				"region":    "aws-global",
				"task":      "ListPolicies",
				"errString": err.Error(),
			}).Error("failed to process region")
			return resultMap, err
		}
		byteArr, err := json.Marshal(page)
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_iam_policy",
				"account":   accountId,
				"region":    "aws-global",
				"task":      "ListPolicies",
				"errString": err.Error(),
			}).Error("failed to marshal response")
			return nil, err
		}
		table := utilities.NewTable(byteArr, tableConfig)
		for _, row := range table.Rows {
			result := extaws.RowToMap(row, accountId, "aws-global", tableConfig)
			resultMap = append(resultMap, result)
		}
		if !paginator.HasMorePages() {
			break
		}
	}
	return resultMap, nil
}

func processAccountListPolicies(account *utilities.ExtensionConfigurationAwsAccount) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	tableConfig, ok := utilities.TableConfigurationMap["aws_iam_policy"]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_iam_policy",
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found")
	}
	result, err := processGlobalListPolicies(tableConfig, account)
	if err != nil {
		return resultMap, err
	}
	resultMap = append(resultMap, result...)
	return resultMap, nil
}
