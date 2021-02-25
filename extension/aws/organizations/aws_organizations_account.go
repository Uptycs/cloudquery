/**
 * Copyright (c) 2020-present, The cloudquery authors
 *
 * This source code is licensed as defined by the LICENSE file found in the
 * root directory of this source tree.
 *
 * SPDX-License-Identifier: (Apache-2.0 OR GPL-2.0-only)
 */

package organizations

import (
	"context"
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/Uptycs/cloudquery/utilities"

	"github.com/Uptycs/basequery-go/plugin/table"
	extaws "github.com/Uptycs/cloudquery/extension/aws"
	"github.com/aws/aws-sdk-go-v2/service/organizations"
)

// ListAccountsColumns returns the list of columns in the table
func ListAccountsColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("account_id"),
		table.TextColumn("arn"),
		table.TextColumn("email"),
		table.TextColumn("id"),
		table.TextColumn("joined_method"),
		table.TextColumn("joined_timestamp"),
		//table.BigIntColumn("joined_timestamp_ext"),
		//table.TextColumn("joined_timestamp_loc"),
		//table.BigIntColumn("joined_timestamp_loc_cache_end"),
		//table.BigIntColumn("joined_timestamp_loc_cache_start"),
		//table.TextColumn("joined_timestamp_loc_cache_zone"),
		//table.TextColumn("joined_timestamp_loc_cache_zone_is_dst"),
		//table.TextColumn("joined_timestamp_loc_cache_zone_name"),
		//table.IntegerColumn("joined_timestamp_loc_cache_zone_offset"),
		//table.TextColumn("joined_timestamp_loc_name"),
		//table.TextColumn("joined_timestamp_loc_tx"),
		//table.IntegerColumn("joined_timestamp_loc_tx_index"),
		//table.TextColumn("joined_timestamp_loc_tx_isstd"),
		//table.TextColumn("joined_timestamp_loc_tx_isutc"),
		//table.BigIntColumn("joined_timestamp_loc_tx_when"),
		//table.TextColumn("joined_timestamp_loc_zone"),
		//table.TextColumn("joined_timestamp_loc_zone_is_dst"),
		//table.TextColumn("joined_timestamp_loc_zone_name"),
		//table.IntegerColumn("joined_timestamp_loc_zone_offset"),
		//table.BigIntColumn("joined_timestamp_wall"),
		table.TextColumn("name"),
		table.TextColumn("status"),
		//table.TextColumn("values"),

	}
}

// ListAccountsGenerate returns the rows in the table for all configured accounts
func ListAccountsGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	if len(utilities.ExtConfiguration.ExtConfAws.Accounts) == 0 && extaws.ShouldProcessAccount("aws_organizations_account", utilities.AwsAccountID) {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_organizations_account",
			"account":   "default",
		}).Info("processing account")
		results, err := processAccountListAccounts(osqCtx, queryContext, nil)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, results...)
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfAws.Accounts {
			if !extaws.ShouldProcessAccount("aws_organizations_account", account.ID) {
				continue
			}
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_organizations_account",
				"account":   account.ID,
			}).Info("processing account")
			results, err := processAccountListAccounts(osqCtx, queryContext, &account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}

	return resultMap, nil
}

func processGlobalListAccounts(osqCtx context.Context, queryContext table.QueryContext, tableConfig *utilities.TableConfig, account *utilities.ExtensionConfigurationAwsAccount) ([]map[string]string, error) {
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
		"tableName": "aws_organizations_account",
		"account":   accountId,
		"region":    "aws-global",
	}).Debug("processing region")

	svc := organizations.NewFromConfig(*sess)
	params := &organizations.ListAccountsInput{}

	paginator := organizations.NewListAccountsPaginator(svc, params)

	for {
		page, err := paginator.NextPage(osqCtx)
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_organizations_account",
				"account":   accountId,
				"region":    "aws-global",
				"task":      "ListAccounts",
				"errString": err.Error(),
			}).Error("failed to process region")
			return resultMap, err
		}
		byteArr, err := json.Marshal(page)
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_organizations_account",
				"account":   accountId,
				"region":    "aws-global",
				"task":      "ListAccounts",
				"errString": err.Error(),
			}).Error("failed to marshal response")
			return nil, err
		}
		table := utilities.NewTable(byteArr, tableConfig)
		for _, row := range table.Rows {
			if !extaws.ShouldProcessRow(osqCtx, queryContext, "aws_organizations_account", accountId, "aws-global", row) {
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

func processAccountListAccounts(osqCtx context.Context, queryContext table.QueryContext, account *utilities.ExtensionConfigurationAwsAccount) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	tableConfig, ok := utilities.TableConfigurationMap["aws_organizations_account"]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_organizations_account",
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found")
	}
	result, err := processGlobalListAccounts(osqCtx, queryContext, tableConfig, account)
	if err != nil {
		return resultMap, err
	}
	resultMap = append(resultMap, result...)
	return resultMap, nil
}
