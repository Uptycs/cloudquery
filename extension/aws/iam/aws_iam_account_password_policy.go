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

	extaws "github.com/Uptycs/cloudquery/extension/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/kolide/osquery-go/plugin/table"
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
	if len(utilities.ExtConfiguration.ExtConfAws.Accounts) == 0 {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_iam_account_password_policy",
			"account":   "default",
		}).Info("processing account")
		results, err := processAccountGetAccountPasswordPolicy(nil)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, results...)
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfAws.Accounts {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_iam_account_password_policy",
				"account":   account.ID,
			}).Info("processing account")
			results, err := processAccountGetAccountPasswordPolicy(&account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}

	return resultMap, nil
}

func processGlobalGetAccountPasswordPolicy(tableConfig *utilities.TableConfig, account *utilities.ExtensionConfigurationAwsAccount) ([]map[string]string, error) {
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

	result, err := svc.GetAccountPasswordPolicy(context.TODO(), params)
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
		result := extaws.RowToMap(row, accountId, "aws-global", tableConfig)
		resultMap = append(resultMap, result)
	}
	return resultMap, nil
}

func processAccountGetAccountPasswordPolicy(account *utilities.ExtensionConfigurationAwsAccount) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	tableConfig, ok := utilities.TableConfigurationMap["aws_iam_account_password_policy"]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_iam_account_password_policy",
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found")
	}
	result, err := processGlobalGetAccountPasswordPolicy(tableConfig, account)
	if err != nil {
		return resultMap, err
	}
	resultMap = append(resultMap, result...)
	return resultMap, nil
}
