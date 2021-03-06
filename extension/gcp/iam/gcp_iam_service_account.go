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

	"github.com/Uptycs/basequery-go/plugin/table"
	extgcp "github.com/Uptycs/cloudquery/extension/gcp"
	"github.com/Uptycs/cloudquery/utilities"

	"google.golang.org/api/option"

	gcpiam "google.golang.org/api/iam/v1"
)

type myGcpIamServiceAccountsItemsContainer struct {
	Items []*gcpiam.ServiceAccount `json:"items"`
}

// GcpIamServiceAccountsColumns returns the list of columns for gcp_iam_service_account
func GcpIamServiceAccountsColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("description"),
		table.TextColumn("disabled"),
		table.TextColumn("display_name"),
		table.TextColumn("email"),
		table.TextColumn("etag"),
		table.TextColumn("name"),
		table.TextColumn("oauth2_client_id"),
		table.TextColumn("project_id"),
		table.TextColumn("unique_id"),
	}
}

// GcpIamServiceAccountsGenerate returns the rows in the table for all configured accounts
func GcpIamServiceAccountsGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	ctx, cancel := context.WithCancel(osqCtx)
	defer cancel()

	resultMap := make([]map[string]string, 0)

	if len(utilities.ExtConfiguration.ExtConfGcp.Accounts) == 0 && extgcp.ShouldProcessProject("gcp_iam_service_account", utilities.DefaultGcpProjectID) {
		results, err := processAccountGcpIamServiceAccounts(ctx, queryContext, nil)
		if err == nil {
			resultMap = append(resultMap, results...)
		}
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfGcp.Accounts {
			if !extgcp.ShouldProcessProject("gcp_iam_service_account", account.ProjectID) {
				continue
			}
			results, err := processAccountGcpIamServiceAccounts(ctx, queryContext, &account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}
	return resultMap, nil
}

func getGcpIamServiceAccountsNewServiceForAccount(ctx context.Context, account *utilities.ExtensionConfigurationGcpAccount) (*gcpiam.Service, string) {
	var projectID string
	var service *gcpiam.Service
	var err error
	if account != nil && account.KeyFile != "" {
		projectID = account.ProjectID
		service, err = gcpiam.NewService(ctx, option.WithCredentialsFile(account.KeyFile))
	} else if account != nil && account.ProjectID != "" {
		projectID = account.ProjectID
		service, err = gcpiam.NewService(ctx)
	} else {
		projectID = utilities.DefaultGcpProjectID
		service, err = gcpiam.NewService(ctx)
	}
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_iam_service_account",
			"projectId": projectID,
			"errString": err.Error(),
		}).Error("failed to create service")
		return nil, ""
	}
	return service, projectID
}

func processAccountGcpIamServiceAccounts(ctx context.Context, queryContext table.QueryContext,
	account *utilities.ExtensionConfigurationGcpAccount) ([]map[string]string, error) {

	resultMap := make([]map[string]string, 0)

	service, projectID := getGcpIamServiceAccountsNewServiceForAccount(ctx, account)
	if service == nil {
		return resultMap, fmt.Errorf("failed to initialize gcpiam.Service")
	}

	listCall := service.Projects.ServiceAccounts.List("projects/" + projectID)
	if listCall == nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_iam_service_account",
			"projectId": projectID,
		}).Debug("list call is nil")
		return resultMap, nil
	}
	itemsContainer := myGcpIamServiceAccountsItemsContainer{Items: make([]*gcpiam.ServiceAccount, 0)}
	if err := listCall.Pages(ctx, func(page *gcpiam.ListServiceAccountsResponse) error {

		itemsContainer.Items = append(itemsContainer.Items, page.Accounts...)

		return nil
	}); err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_iam_service_account",
			"projectId": projectID,
			"errString": err.Error(),
		}).Error("failed to get aggregate list page")
		return resultMap, nil
	}

	byteArr, err := json.Marshal(itemsContainer)
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_iam_service_account",
			"errString": err.Error(),
		}).Error("failed to marshal response")
		return resultMap, err
	}
	tableConfig, ok := utilities.TableConfigurationMap["gcp_iam_service_account"]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_iam_service_account",
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found for \"gcp_iam_service_account\"")
	}
	jsonTable := utilities.NewTable(byteArr, tableConfig)
	for _, row := range jsonTable.Rows {
		if !extgcp.ShouldProcessRow(ctx, queryContext, "gcp_iam_service_account", projectID, "", row) {
			continue
		}
		result := extgcp.RowToMap(row, projectID, "", tableConfig)
		resultMap = append(resultMap, result)
	}

	return resultMap, nil
}
