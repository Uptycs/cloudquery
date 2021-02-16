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

// ListDelegatedAdministratorsColumns returns the list of columns in the table
func ListDelegatedAdministratorsColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("account_id"),
		table.TextColumn("arn"),
		table.TextColumn("delegation_enabled_date"),
		//table.BigIntColumn("delegation_enabled_date_ext"),
		//table.TextColumn("delegation_enabled_date_loc"),
		//table.BigIntColumn("delegation_enabled_date_loc_cache_end"),
		//table.BigIntColumn("delegation_enabled_date_loc_cache_start"),
		//table.TextColumn("delegation_enabled_date_loc_cache_zone"),
		//table.TextColumn("delegation_enabled_date_loc_cache_zone_is_dst"),
		//table.TextColumn("delegation_enabled_date_loc_cache_zone_name"),
		//table.IntegerColumn("delegation_enabled_date_loc_cache_zone_offset"),
		//table.TextColumn("delegation_enabled_date_loc_name"),
		//table.TextColumn("delegation_enabled_date_loc_tx"),
		//table.IntegerColumn("delegation_enabled_date_loc_tx_index"),
		//table.TextColumn("delegation_enabled_date_loc_tx_isstd"),
		//table.TextColumn("delegation_enabled_date_loc_tx_isutc"),
		//table.BigIntColumn("delegation_enabled_date_loc_tx_when"),
		//table.TextColumn("delegation_enabled_date_loc_zone"),
		//table.TextColumn("delegation_enabled_date_loc_zone_is_dst"),
		//table.TextColumn("delegation_enabled_date_loc_zone_name"),
		//table.IntegerColumn("delegation_enabled_date_loc_zone_offset"),
		//table.BigIntColumn("delegation_enabled_date_wall"),
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

// ListDelegatedAdministratorsGenerate returns the rows in the table for all configured accounts
func ListDelegatedAdministratorsGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	if len(utilities.ExtConfiguration.ExtConfAws.Accounts) == 0 {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_organizations_delegated_administrator",
			"account":   "default",
		}).Info("processing account")
		results, err := processAccountListDelegatedAdministrators(nil)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, results...)
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfAws.Accounts {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_organizations_delegated_administrator",
				"account":   account.ID,
			}).Info("processing account")
			results, err := processAccountListDelegatedAdministrators(&account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}

	return resultMap, nil
}

func processGlobalListDelegatedAdministrators(tableConfig *utilities.TableConfig, account *utilities.ExtensionConfigurationAwsAccount) ([]map[string]string, error) {
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
		"tableName": "aws_organizations_delegated_administrator",
		"account":   accountId,
		"region":    "aws-global",
	}).Debug("processing region")

	svc := organizations.NewFromConfig(*sess)
	params := &organizations.ListDelegatedAdministratorsInput{}

	paginator := organizations.NewListDelegatedAdministratorsPaginator(svc, params)

	for {
		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_organizations_delegated_administrator",
				"account":   accountId,
				"region":    "aws-global",
				"task":      "ListDelegatedAdministrators",
				"errString": err.Error(),
			}).Error("failed to process region")
			return resultMap, err
		}
		byteArr, err := json.Marshal(page)
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_organizations_delegated_administrator",
				"account":   accountId,
				"region":    "aws-global",
				"task":      "ListDelegatedAdministrators",
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

func processAccountListDelegatedAdministrators(account *utilities.ExtensionConfigurationAwsAccount) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	tableConfig, ok := utilities.TableConfigurationMap["aws_organizations_delegated_administrator"]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_organizations_delegated_administrator",
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found")
	}
	result, err := processGlobalListDelegatedAdministrators(tableConfig, account)
	if err != nil {
		return resultMap, err
	}
	resultMap = append(resultMap, result...)
	return resultMap, nil
}
