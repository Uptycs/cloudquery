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

// ListGroupsColumns returns the list of columns in the table
func ListGroupsColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("account_id"),
		table.TextColumn("arn"),
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
		table.TextColumn("group_id"),
		table.TextColumn("group_name"),
		table.TextColumn("path"),
	}
}

// ListGroupsGenerate returns the rows in the table for all configured accounts
func ListGroupsGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	if len(utilities.ExtConfiguration.ExtConfAws.Accounts) == 0 && extaws.ShouldProcessAccount("aws_iam_group", utilities.AwsAccountID) {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_iam_group",
			"account":   "default",
		}).Info("processing account")
		results, err := processAccountListGroups(osqCtx, queryContext, nil)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, results...)
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfAws.Accounts {
			if !extaws.ShouldProcessAccount("aws_iam_group", account.ID) {
				continue
			}
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_iam_group",
				"account":   account.ID,
			}).Info("processing account")
			results, err := processAccountListGroups(osqCtx, queryContext, &account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}

	return resultMap, nil
}

func processGlobalListGroups(osqCtx context.Context, queryContext table.QueryContext, tableConfig *utilities.TableConfig, account *utilities.ExtensionConfigurationAwsAccount) ([]map[string]string, error) {
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
		"tableName": "aws_iam_group",
		"account":   accountId,
		"region":    "aws-global",
	}).Debug("processing region")

	svc := iam.NewFromConfig(*sess)
	params := &iam.ListGroupsInput{}

	paginator := iam.NewListGroupsPaginator(svc, params)

	for {
		page, err := paginator.NextPage(osqCtx)
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_iam_group",
				"account":   accountId,
				"region":    "aws-global",
				"task":      "ListGroups",
				"errString": err.Error(),
			}).Error("failed to process region")
			return resultMap, err
		}
		byteArr, err := json.Marshal(page)
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_iam_group",
				"account":   accountId,
				"region":    "aws-global",
				"task":      "ListGroups",
				"errString": err.Error(),
			}).Error("failed to marshal response")
			return nil, err
		}
		table := utilities.NewTable(byteArr, tableConfig)
		for _, row := range table.Rows {
			if !extaws.ShouldProcessRow(osqCtx, queryContext, "aws_iam_group", accountId, "aws-global", row) {
				continue
			}
			result := extaws.RowToMap(row, accountId, "aws-global", tableConfig)
			resultMap = append(resultMap, result)
		}
		if !paginator.HasMorePages() {
			break
		}
	}
	return resultMap, nil
}

func processAccountListGroups(osqCtx context.Context, queryContext table.QueryContext, account *utilities.ExtensionConfigurationAwsAccount) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	tableConfig, ok := utilities.TableConfigurationMap["aws_iam_group"]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_iam_group",
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found")
	}
	result, err := processGlobalListGroups(osqCtx, queryContext, tableConfig, account)
	if err != nil {
		return resultMap, err
	}
	resultMap = append(resultMap, result...)
	return resultMap, nil
}
