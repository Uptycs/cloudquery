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

// GetAccountPasswordPolicyColumns returns the list of columns in the table
func GetAccountPasswordPolicyColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("account_id"),
		table.TextColumn("allow_users_to_change_password"),
		table.TextColumn("expire_passwords"),
		table.TextColumn("hard_expiry"),
		table.BigIntColumn("max_password_age"),
		table.BigIntColumn("minimum_password_length"),
		table.BigIntColumn("password_reuse_prevention"),
		table.TextColumn("require_lowercase_characters"),
		table.TextColumn("require_numbers"),
		table.TextColumn("require_symbols"),
		table.TextColumn("require_uppercase_characters"),
	}
}

// GetAccountPasswordPolicyGenerate returns the rows in the table for all configured accounts
func GetAccountPasswordPolicyGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	if len(utilities.ExtConfiguration.ExtConfAws.Accounts) == 0 && extaws.ShouldProcessAccount("aws_iam_account_password_policy", utilities.AwsAccountID) {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_iam_account_password_policy",
			"account":   "default",
		}).Info("processing account")
		results, err := processAccountGetAccountPasswordPolicy(osqCtx, queryContext, nil)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, results...)
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfAws.Accounts {
			if !extaws.ShouldProcessAccount("aws_iam_account_password_policy", account.ID) {
				continue
			}
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_iam_account_password_policy",
				"account":   account.ID,
			}).Info("processing account")
			results, err := processAccountGetAccountPasswordPolicy(osqCtx, queryContext, &account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}

	return resultMap, nil
}

func processGlobalGetAccountPasswordPolicy(osqCtx context.Context, queryContext table.QueryContext, tableConfig *utilities.TableConfig, account *utilities.ExtensionConfigurationAwsAccount) ([]map[string]string, error) {
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
		"tableName": "aws_iam_account_password_policy",
		"account":   accountId,
		"region":    "aws-global",
	}).Debug("processing region")

	svc := iam.NewFromConfig(*sess)
	params := &iam.GetAccountPasswordPolicyInput{}

	result, err := svc.GetAccountPasswordPolicy(osqCtx, params)
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_iam_account_password_policy",
			"account":   accountId,
			"region":    "aws-global",
			"task":      "GetAccountPasswordPolicy",
			"errString": err.Error(),
		}).Error("failed to process region")
		return resultMap, err
	}

	byteArr, err := json.Marshal(result)
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_iam_account_password_policy",
			"account":   accountId,
			"region":    "aws-global",
			"errString": err.Error(),
		}).Error("failed to marshal response")
		return resultMap, err
	}
	table := utilities.NewTable(byteArr, tableConfig)
	for _, row := range table.Rows {
		if !extaws.ShouldProcessRow(osqCtx, queryContext, "aws_iam_account_password_policy", accountId, "aws-global", row) {
			continue
		}
		result := extaws.RowToMap(row, accountId, "aws-global", tableConfig)
		resultMap = append(resultMap, result)
	}
	return resultMap, nil
}

func processAccountGetAccountPasswordPolicy(osqCtx context.Context, queryContext table.QueryContext, account *utilities.ExtensionConfigurationAwsAccount) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	tableConfig, ok := utilities.TableConfigurationMap["aws_iam_account_password_policy"]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_iam_account_password_policy",
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found")
	}
	result, err := processGlobalGetAccountPasswordPolicy(osqCtx, queryContext, tableConfig, account)
	if err != nil {
		return resultMap, err
	}
	resultMap = append(resultMap, result...)
	return resultMap, nil
}
