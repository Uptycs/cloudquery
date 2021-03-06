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

// ListUsersColumns returns the list of columns in the table
func ListUsersColumns() []table.ColumnDefinition {
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
		table.TextColumn("password_last_used"),
		//table.BigIntColumn("password_last_used_ext"),
		//table.TextColumn("password_last_used_loc"),
		//table.BigIntColumn("password_last_used_loc_cache_end"),
		//table.BigIntColumn("password_last_used_loc_cache_start"),
		//table.TextColumn("password_last_used_loc_cache_zone"),
		//table.TextColumn("password_last_used_loc_cache_zone_is_dst"),
		//table.TextColumn("password_last_used_loc_cache_zone_name"),
		//table.IntegerColumn("password_last_used_loc_cache_zone_offset"),
		//table.TextColumn("password_last_used_loc_name"),
		//table.TextColumn("password_last_used_loc_tx"),
		//table.IntegerColumn("password_last_used_loc_tx_index"),
		//table.TextColumn("password_last_used_loc_tx_isstd"),
		//table.TextColumn("password_last_used_loc_tx_isutc"),
		//table.BigIntColumn("password_last_used_loc_tx_when"),
		//table.TextColumn("password_last_used_loc_zone"),
		//table.TextColumn("password_last_used_loc_zone_is_dst"),
		//table.TextColumn("password_last_used_loc_zone_name"),
		//table.IntegerColumn("password_last_used_loc_zone_offset"),
		//table.BigIntColumn("password_last_used_wall"),
		table.TextColumn("path"),
		table.TextColumn("permissions_boundary"),
		//table.TextColumn("permissions_boundary_permissions_boundary_arn"),
		//table.TextColumn("permissions_boundary_permissions_boundary_type"),
		table.TextColumn("tags"),
		//table.TextColumn("tags_key"),
		//table.TextColumn("tags_value"),
		table.TextColumn("user_id"),
		table.TextColumn("user_name"),
	}
}

// ListUsersGenerate returns the rows in the table for all configured accounts
func ListUsersGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	if len(utilities.ExtConfiguration.ExtConfAws.Accounts) == 0 && extaws.ShouldProcessAccount("aws_iam_user", utilities.AwsAccountID) {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_iam_user",
			"account":   "default",
		}).Info("processing account")
		results, err := processAccountListUsers(osqCtx, queryContext, nil)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, results...)
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfAws.Accounts {
			if !extaws.ShouldProcessAccount("aws_iam_user", account.ID) {
				continue
			}
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_iam_user",
				"account":   account.ID,
			}).Info("processing account")
			results, err := processAccountListUsers(osqCtx, queryContext, &account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}

	return resultMap, nil
}

func processGlobalListUsers(osqCtx context.Context, queryContext table.QueryContext, tableConfig *utilities.TableConfig, account *utilities.ExtensionConfigurationAwsAccount) ([]map[string]string, error) {
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
		"tableName": "aws_iam_user",
		"account":   accountId,
		"region":    "aws-global",
	}).Debug("processing region")

	svc := iam.NewFromConfig(*sess)
	params := &iam.ListUsersInput{}

	paginator := iam.NewListUsersPaginator(svc, params)

	for {
		page, err := paginator.NextPage(osqCtx)
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_iam_user",
				"account":   accountId,
				"region":    "aws-global",
				"task":      "ListUsers",
				"errString": err.Error(),
			}).Error("failed to process region")
			return resultMap, err
		}
		byteArr, err := json.Marshal(page)
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_iam_user",
				"account":   accountId,
				"region":    "aws-global",
				"task":      "ListUsers",
				"errString": err.Error(),
			}).Error("failed to marshal response")
			return nil, err
		}
		table := utilities.NewTable(byteArr, tableConfig)
		for _, row := range table.Rows {
			if !extaws.ShouldProcessRow(osqCtx, queryContext, "aws_iam_user", accountId, "aws-global", row) {
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

func processAccountListUsers(osqCtx context.Context, queryContext table.QueryContext, account *utilities.ExtensionConfigurationAwsAccount) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	tableConfig, ok := utilities.TableConfigurationMap["aws_iam_user"]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_iam_user",
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found")
	}
	result, err := processGlobalListUsers(osqCtx, queryContext, tableConfig, account)
	if err != nil {
		return resultMap, err
	}
	resultMap = append(resultMap, result...)
	return resultMap, nil
}
