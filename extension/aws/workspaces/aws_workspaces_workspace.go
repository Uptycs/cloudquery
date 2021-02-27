/**
 * Copyright (c) 2020-present, The cloudquery authors
 *
 * This source code is licensed as defined by the LICENSE file found in the
 * root directory of this source tree.
 *
 * SPDX-License-Identifier: (Apache-2.0 OR GPL-2.0-only)
 */

package workspaces

import (
	"context"
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/Uptycs/cloudquery/utilities"

	"github.com/Uptycs/basequery-go/plugin/table"
	extaws "github.com/Uptycs/cloudquery/extension/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/workspaces"
)

// DescribeWorkspacesColumns returns the list of columns in the table
func DescribeWorkspacesColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("account_id"),
		table.TextColumn("region_code"),
		//table.TextColumn("values"),
		table.TextColumn("bundle_id"),
		table.TextColumn("computer_name"),
		table.TextColumn("directory_id"),
		table.TextColumn("error_code"),
		table.TextColumn("error_message"),
		table.TextColumn("ip_address"),
		table.TextColumn("modification_states"),
		//table.TextColumn("modification_states_resource"),
		//table.TextColumn("modification_states_state"),
		table.TextColumn("root_volume_encryption_enabled"),
		table.TextColumn("state"),
		table.TextColumn("subnet_id"),
		table.TextColumn("user_name"),
		table.TextColumn("user_volume_encryption_enabled"),
		table.TextColumn("volume_encryption_key"),
		table.TextColumn("workspace_id"),
		table.TextColumn("workspace_properties"),
		//table.TextColumn("workspace_properties_compute_type_name"),
		//table.IntegerColumn("workspace_properties_root_volume_size_gib"),
		//table.TextColumn("workspace_properties_running_mode"),
		//table.IntegerColumn("workspace_properties_running_mode_auto_stop_timeout_in_minutes"),
		//table.IntegerColumn("workspace_properties_user_volume_size_gib"),

	}
}

// DescribeWorkspacesGenerate returns the rows in the table for all configured accounts
func DescribeWorkspacesGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	if len(utilities.ExtConfiguration.ExtConfAws.Accounts) == 0 {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_workspaces_workspace",
			"account":   "default",
		}).Info("processing account")
		results, err := processAccountDescribeWorkspaces(nil)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, results...)
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfAws.Accounts {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_workspaces_workspace",
				"account":   account.ID,
			}).Info("processing account")
			results, err := processAccountDescribeWorkspaces(&account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}

	return resultMap, nil
}

func processRegionDescribeWorkspaces(tableConfig *utilities.TableConfig, account *utilities.ExtensionConfigurationAwsAccount, region types.Region) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	sess, err := extaws.GetAwsConfig(account, *region.RegionName)
	if err != nil {
		return resultMap, err
	}

	accountId := utilities.AwsAccountID
	if account != nil {
		accountId = account.ID
	}

	utilities.GetLogger().WithFields(log.Fields{
		"tableName": "aws_workspaces_workspace",
		"account":   accountId,
		"region":    *region.RegionName,
	}).Debug("processing region")

	svc := workspaces.NewFromConfig(*sess)
	params := &workspaces.DescribeWorkspacesInput{}

	paginator := workspaces.NewDescribeWorkspacesPaginator(svc, params)

	for {
		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_workspaces_workspace",
				"account":   accountId,
				"region":    *region.RegionName,
				"task":      "DescribeWorkspaces",
				"errString": err.Error(),
			}).Error("failed to process region")
			return resultMap, err
		}
		byteArr, err := json.Marshal(page)
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_workspaces_workspace",
				"account":   accountId,
				"region":    *region.RegionName,
				"task":      "DescribeWorkspaces",
				"errString": err.Error(),
			}).Error("failed to marshal response")
			return nil, err
		}
		table := utilities.NewTable(byteArr, tableConfig)
		for _, row := range table.Rows {
			result := extaws.RowToMap(row, accountId, *region.RegionName, tableConfig)
			resultMap = append(resultMap, result)
		}
		if !paginator.HasMorePages() {
			break
		}
	}
	return resultMap, nil
}

func processAccountDescribeWorkspaces(account *utilities.ExtensionConfigurationAwsAccount) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	awsSession, err := extaws.GetAwsConfig(account, "us-east-1")
	if err != nil {
		return resultMap, err
	}
	regions, err := extaws.FetchRegions(context.TODO(), awsSession)
	if err != nil {
		return resultMap, err
	}
	tableConfig, ok := utilities.TableConfigurationMap["aws_workspaces_workspace"]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_workspaces_workspace",
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found")
	}
	for _, region := range regions {
		result, err := processRegionDescribeWorkspaces(tableConfig, account, region)
		if err != nil {
			continue
		}
		resultMap = append(resultMap, result...)
	}
	return resultMap, nil
}
