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

// ListRootsColumns returns the list of columns in the table
func ListRootsColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("account_id"),
		//table.TextColumn("values"),
		table.TextColumn("arn"),
		table.TextColumn("id"),
		table.TextColumn("name"),
		table.TextColumn("policy_types"),
		table.TextColumn("policy_types_status"),
		//table.TextColumn("policy_types_type"),

	}
}

// ListRootsGenerate returns the rows in the table for all configured accounts
func ListRootsGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	if len(utilities.ExtConfiguration.ExtConfAws.Accounts) == 0 {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_organizations_list_roots",
			"account":   "default",
		}).Info("processing account")
		results, err := processAccountListRoots(nil)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, results...)
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfAws.Accounts {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_organizations_list_roots",
				"account":   account.ID,
			}).Info("processing account")
			results, err := processAccountListRoots(&account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}

	return resultMap, nil
}

func processGlobalListRoots(tableConfig *utilities.TableConfig, account *utilities.ExtensionConfigurationAwsAccount) ([]map[string]string, error) {
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
		"tableName": "aws_organizations_list_roots",
		"account":   accountId,
		"region":    "aws-global",
	}).Debug("processing region")

	svc := organizations.NewFromConfig(*sess)
	params := &organizations.ListRootsInput{}

	paginator := organizations.NewListRootsPaginator(svc, params)

	for {
		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_organizations_list_roots",
				"account":   accountId,
				"region":    "aws-global",
				"task":      "ListRoots",
				"errString": err.Error(),
			}).Error("failed to process region")
			return resultMap, err
		}
		byteArr, err := json.Marshal(page)
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_organizations_list_roots",
				"account":   accountId,
				"region":    "aws-global",
				"task":      "ListRoots",
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

func processAccountListRoots(account *utilities.ExtensionConfigurationAwsAccount) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	tableConfig, ok := utilities.TableConfigurationMap["aws_organizations_list_roots"]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_organizations_list_roots",
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found")
	}
	result, err := processGlobalListRoots(tableConfig, account)
	if err != nil {
		return resultMap, err
	}
	resultMap = append(resultMap, result...)
	return resultMap, nil
}
