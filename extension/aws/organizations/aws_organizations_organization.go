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

// DescribeOrganizationColumns returns the list of columns in the table
func DescribeOrganizationColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("account_id"),
		table.TextColumn("arn"),
		table.TextColumn("available_policy_types"),
		//table.TextColumn("available_policy_types_status"),
		//table.TextColumn("available_policy_types_type"),
		table.TextColumn("feature_set"),
		table.TextColumn("id"),
		table.TextColumn("master_account_arn"),
		table.TextColumn("master_account_email"),
		table.TextColumn("master_account_id"),
		//table.TextColumn("values"),

	}
}

// DescribeOrganizationGenerate returns the rows in the table for all configured accounts
func DescribeOrganizationGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	if len(utilities.ExtConfiguration.ExtConfAws.Accounts) == 0 && extaws.ShouldProcessAccount("aws_organizations_organization", utilities.AwsAccountID) {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_organizations_organization",
			"account":   "default",
		}).Info("processing account")
		results, err := processAccountDescribeOrganization(osqCtx, queryContext, nil)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, results...)
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfAws.Accounts {
			if !extaws.ShouldProcessAccount("aws_organizations_organization", account.ID) {
				continue
			}
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_organizations_organization",
				"account":   account.ID,
			}).Info("processing account")
			results, err := processAccountDescribeOrganization(osqCtx, queryContext, &account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}

	return resultMap, nil
}

func processGlobalDescribeOrganization(osqCtx context.Context, queryContext table.QueryContext, tableConfig *utilities.TableConfig, account *utilities.ExtensionConfigurationAwsAccount) ([]map[string]string, error) {
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
		"tableName": "aws_organizations_organization",
		"account":   accountId,
		"region":    "aws-global",
	}).Debug("processing region")

	svc := organizations.NewFromConfig(*sess)
	params := &organizations.DescribeOrganizationInput{}

	result, err := svc.DescribeOrganization(osqCtx, params)
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_organizations_organization",
			"account":   accountId,
			"region":    "aws-global",
			"task":      "DescribeOrganization",
			"errString": err.Error(),
		}).Error("failed to process region")
		return resultMap, err
	}

	byteArr, err := json.Marshal(result)
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_organizations_organization",
			"account":   accountId,
			"region":    "aws-global",
			"errString": err.Error(),
		}).Error("failed to marshal response")
		return resultMap, err
	}
	table := utilities.NewTable(byteArr, tableConfig)
	for _, row := range table.Rows {
		if !extaws.ShouldProcessRow(osqCtx, queryContext, "aws_organizations_organization", accountId, "aws-global", row) {
			continue
		}
		result := extaws.RowToMap(row, accountId, "aws-global", tableConfig)
		resultMap = append(resultMap, result)
	}
	return resultMap, nil
}

func processAccountDescribeOrganization(osqCtx context.Context, queryContext table.QueryContext, account *utilities.ExtensionConfigurationAwsAccount) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	tableConfig, ok := utilities.TableConfigurationMap["aws_organizations_organization"]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_organizations_organization",
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found")
	}
	result, err := processGlobalDescribeOrganization(osqCtx, queryContext, tableConfig, account)
	if err != nil {
		return resultMap, err
	}
	resultMap = append(resultMap, result...)
	return resultMap, nil
}
