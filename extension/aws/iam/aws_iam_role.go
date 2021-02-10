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

// ListRolesColumns returns the list of columns in the table
func ListRolesColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("account_id"),
		table.TextColumn("arn"),
		table.TextColumn("assume_role_policy_document"),
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
		table.TextColumn("description"),
		table.BigIntColumn("max_session_duration"),
		table.TextColumn("path"),
		table.TextColumn("permissions_boundary"),
		//table.TextColumn("permissions_boundary_permissions_boundary_arn"),
		//table.TextColumn("permissions_boundary_permissions_boundary_type"),
		table.TextColumn("role_id"),
		table.TextColumn("role_last_used"),
		//table.TextColumn("role_last_used_last_used_date"),
		//table.BigIntColumn("role_last_used_last_used_date_ext"),
		//table.TextColumn("role_last_used_last_used_date_loc"),
		//table.BigIntColumn("role_last_used_last_used_date_loc_cache_end"),
		//table.BigIntColumn("role_last_used_last_used_date_loc_cache_start"),
		//table.TextColumn("role_last_used_last_used_date_loc_cache_zone"),
		//table.TextColumn("role_last_used_last_used_date_loc_cache_zone_is_dst"),
		//table.TextColumn("role_last_used_last_used_date_loc_cache_zone_name"),
		//table.IntegerColumn("role_last_used_last_used_date_loc_cache_zone_offset"),
		//table.TextColumn("role_last_used_last_used_date_loc_name"),
		//table.TextColumn("role_last_used_last_used_date_loc_tx"),
		//table.IntegerColumn("role_last_used_last_used_date_loc_tx_index"),
		//table.TextColumn("role_last_used_last_used_date_loc_tx_isstd"),
		//table.TextColumn("role_last_used_last_used_date_loc_tx_isutc"),
		//table.BigIntColumn("role_last_used_last_used_date_loc_tx_when"),
		//table.TextColumn("role_last_used_last_used_date_loc_zone"),
		//table.TextColumn("role_last_used_last_used_date_loc_zone_is_dst"),
		//table.TextColumn("role_last_used_last_used_date_loc_zone_name"),
		//table.IntegerColumn("role_last_used_last_used_date_loc_zone_offset"),
		//table.BigIntColumn("role_last_used_last_used_date_wall"),
		//table.TextColumn("role_last_used_region"),
		table.TextColumn("role_name"),
		table.TextColumn("tags"),
		//table.TextColumn("tags_key"),
		//table.TextColumn("tags_value"),

	}
}

// ListRolesGenerate returns the rows in the table for all configured accounts
func ListRolesGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	if len(utilities.ExtConfiguration.ExtConfAws.Accounts) == 0 {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_iam_role",
			"account":   "default",
		}).Info("processing account")
		results, err := processAccountListRoles(nil)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, results...)
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfAws.Accounts {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_iam_role",
				"account":   account.ID,
			}).Info("processing account")
			results, err := processAccountListRoles(&account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}

	return resultMap, nil
}

func processGlobalListRoles(tableConfig *utilities.TableConfig, account *utilities.ExtensionConfigurationAwsAccount) ([]map[string]string, error) {
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
		"tableName": "aws_iam_role",
		"account":   accountId,
		"region":    "aws-global",
	}).Debug("processing region")

	svc := iam.NewFromConfig(*sess)
	params := &iam.ListRolesInput{}

	paginator := iam.NewListRolesPaginator(svc, params)

	for {
		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_iam_role",
				"account":   accountId,
				"region":    "aws-global",
				"task":      "ListRoles",
				"errString": err.Error(),
			}).Error("failed to process region")
			return resultMap, err
		}
		byteArr, err := json.Marshal(page)
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_iam_role",
				"account":   accountId,
				"region":    "aws-global",
				"task":      "ListRoles",
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

func processAccountListRoles(account *utilities.ExtensionConfigurationAwsAccount) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	tableConfig, ok := utilities.TableConfigurationMap["aws_iam_role"]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_iam_role",
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found")
	}
	result, err := processGlobalListRoles(tableConfig, account)
	if err != nil {
		return resultMap, err
	}
	resultMap = append(resultMap, result...)
	return resultMap, nil
}
