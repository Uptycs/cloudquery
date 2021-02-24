/**
 * Copyright (c) 2020-present, The cloudquery authors
 *
 * This source code is licensed as defined by the LICENSE file found in the
 * root directory of this source tree.
 *
 * SPDX-License-Identifier: (Apache-2.0 OR GPL-2.0-only)
 */

package file

import (
	"context"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"

	"github.com/Uptycs/basequery-go/plugin/table"
	extgcp "github.com/Uptycs/cloudquery/extension/gcp"
	"github.com/Uptycs/cloudquery/utilities"

	"google.golang.org/api/option"

	gcpfile "google.golang.org/api/file/v1beta1"
)

type myGcpFileBackupsItemsContainer struct {
	Items []*gcpfile.Backup `json:"items"`
}

// GcpFileBackupsColumns returns the list of columns for gcp_file_backup
func GcpFileBackupsColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("project_id"),
		table.BigIntColumn("capacity_gb"),
		table.TextColumn("create_time"),
		table.TextColumn("description"),
		table.BigIntColumn("download_bytes"),
		table.TextColumn("labels"),
		table.TextColumn("name"),
		table.TextColumn("source_file_share"),
		table.TextColumn("source_instance"),
		//table.TextColumn("source_instance_tier"),
		table.TextColumn("state"),
		table.BigIntColumn("storage_bytes"),
	}
}

// GcpFileBackupsGenerate returns the rows in the table for all configured accounts
func GcpFileBackupsGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	ctx, cancel := context.WithCancel(osqCtx)
	defer cancel()

	resultMap := make([]map[string]string, 0)

	if len(utilities.ExtConfiguration.ExtConfGcp.Accounts) == 0 && extgcp.ShouldProcessProject("gcp_file_backup", utilities.DefaultGcpProjectID) {
		results, err := processAccountGcpFileBackups(ctx, queryContext, nil)
		if err == nil {
			resultMap = append(resultMap, results...)
		}
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfGcp.Accounts {
			if !extgcp.ShouldProcessProject("gcp_file_backup", account.ProjectID) {
				continue
			}
			results, err := processAccountGcpFileBackups(ctx, queryContext, &account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}
	return resultMap, nil
}

func getGcpFileBackupsNewServiceForAccount(ctx context.Context, account *utilities.ExtensionConfigurationGcpAccount) (*gcpfile.Service, string) {
	var projectID string
	var service *gcpfile.Service
	var err error
	if account != nil && account.KeyFile != "" {
		projectID = account.ProjectID
		service, err = gcpfile.NewService(ctx, option.WithCredentialsFile(account.KeyFile))
	} else if account != nil && account.ProjectID != "" {
		projectID = account.ProjectID
		service, err = gcpfile.NewService(ctx)
	} else {
		projectID = utilities.DefaultGcpProjectID
		service, err = gcpfile.NewService(ctx)
	}
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_file_backup",
			"projectId": projectID,
			"errString": err.Error(),
		}).Error("failed to create service")
		return nil, ""
	}
	return service, projectID
}

func processAccountGcpFileBackups(ctx context.Context, queryContext table.QueryContext,
	account *utilities.ExtensionConfigurationGcpAccount) ([]map[string]string, error) {

	resultMap := make([]map[string]string, 0)

	service, projectID := getGcpFileBackupsNewServiceForAccount(ctx, account)
	if service == nil {
		return resultMap, fmt.Errorf("failed to initialize gcpfile.Service")
	}

	listCall := service.Projects.Locations.Backups.List("projects/" + projectID + "/locations/-")
	if listCall == nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_file_backup",
			"projectId": projectID,
		}).Debug("list call is nil")
		return resultMap, nil
	}
	itemsContainer := myGcpFileBackupsItemsContainer{Items: make([]*gcpfile.Backup, 0)}
	if err := listCall.Pages(ctx, func(page *gcpfile.ListBackupsResponse) error {

		itemsContainer.Items = append(itemsContainer.Items, page.Backups...)

		return nil
	}); err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_file_backup",
			"projectId": projectID,
			"errString": err.Error(),
		}).Error("failed to get aggregate list page")
		return resultMap, nil
	}

	byteArr, err := json.Marshal(itemsContainer)
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_file_backup",
			"errString": err.Error(),
		}).Error("failed to marshal response")
		return resultMap, err
	}
	tableConfig, ok := utilities.TableConfigurationMap["gcp_file_backup"]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_file_backup",
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found for \"gcp_file_backup\"")
	}
	jsonTable := utilities.NewTable(byteArr, tableConfig)
	for _, row := range jsonTable.Rows {
		if !extgcp.ShouldProcessRow(ctx, queryContext, "gcp_file_backup", projectID, "", row) {
			continue
		}
		result := extgcp.RowToMap(row, projectID, "", tableConfig)
		resultMap = append(resultMap, result)
	}

	return resultMap, nil
}
