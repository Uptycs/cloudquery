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

type myGcpFileInstancesItemsContainer struct {
	Items []*gcpfile.Instance `json:"items"`
}

// GcpFileInstancesColumns returns the list of columns for gcp_file_instance
func GcpFileInstancesColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("project_id"),
		table.TextColumn("create_time"),
		table.TextColumn("description"),
		table.TextColumn("etag"),
		table.TextColumn("file_shares"),
		//table.BigIntColumn("file_shares_capacity_gb"),
		//table.TextColumn("file_shares_name"),
		//table.TextColumn("file_shares_nfs_export_options"),
		//table.TextColumn("file_shares_nfs_export_options_access_mode"),
		//table.BigIntColumn("file_shares_nfs_export_options_anon_gid"),
		//table.BigIntColumn("file_shares_nfs_export_options_anon_uid"),
		//table.TextColumn("file_shares_nfs_export_options_ip_ranges"),
		//table.TextColumn("file_shares_nfs_export_options_squash_mode"),
		//table.TextColumn("file_shares_source_backup"),
		table.TextColumn("labels"),
		table.TextColumn("name"),
		table.TextColumn("networks"),
		//table.TextColumn("networks_ip_addresses"),
		//table.TextColumn("networks_modes"),
		//table.TextColumn("networks_network"),
		//table.TextColumn("networks_reserved_ip_range"),
		table.TextColumn("state"),
		table.TextColumn("status_message"),
		table.TextColumn("tier"),
	}
}

// GcpFileInstancesGenerate returns the rows in the table for all configured accounts
func GcpFileInstancesGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	ctx, cancel := context.WithCancel(osqCtx)
	defer cancel()

	resultMap := make([]map[string]string, 0)

	if len(utilities.ExtConfiguration.ExtConfGcp.Accounts) == 0 && extgcp.ShouldProcessProject("gcp_file_instance", utilities.DefaultGcpProjectID) {
		results, err := processAccountGcpFileInstances(ctx, queryContext, nil)
		if err == nil {
			resultMap = append(resultMap, results...)
		}
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfGcp.Accounts {
			if !extgcp.ShouldProcessProject("gcp_file_instance", account.ProjectID) {
				continue
			}
			results, err := processAccountGcpFileInstances(ctx, queryContext, &account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}
	return resultMap, nil
}

func getGcpFileInstancesNewServiceForAccount(ctx context.Context, account *utilities.ExtensionConfigurationGcpAccount) (*gcpfile.Service, string) {
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
			"tableName": "gcp_file_instance",
			"projectId": projectID,
			"errString": err.Error(),
		}).Error("failed to create service")
		return nil, ""
	}
	return service, projectID
}

func processAccountGcpFileInstances(ctx context.Context, queryContext table.QueryContext,
	account *utilities.ExtensionConfigurationGcpAccount) ([]map[string]string, error) {

	resultMap := make([]map[string]string, 0)

	service, projectID := getGcpFileInstancesNewServiceForAccount(ctx, account)
	if service == nil {
		return resultMap, fmt.Errorf("failed to initialize gcpfile.Service")
	}

	listCall := service.Projects.Locations.Instances.List("projects/" + projectID + "/locations/-")
	if listCall == nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_file_instance",
			"projectId": projectID,
		}).Debug("list call is nil")
		return resultMap, nil
	}
	itemsContainer := myGcpFileInstancesItemsContainer{Items: make([]*gcpfile.Instance, 0)}
	if err := listCall.Pages(ctx, func(page *gcpfile.ListInstancesResponse) error {

		itemsContainer.Items = append(itemsContainer.Items, page.Instances...)

		return nil
	}); err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_file_instance",
			"projectId": projectID,
			"errString": err.Error(),
		}).Error("failed to get aggregate list page")
		return resultMap, nil
	}

	byteArr, err := json.Marshal(itemsContainer)
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_file_instance",
			"errString": err.Error(),
		}).Error("failed to marshal response")
		return resultMap, err
	}
	tableConfig, ok := utilities.TableConfigurationMap["gcp_file_instance"]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_file_instance",
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found for \"gcp_file_instance\"")
	}
	jsonTable := utilities.NewTable(byteArr, tableConfig)
	for _, row := range jsonTable.Rows {
		if !extgcp.ShouldProcessRow(ctx, queryContext, "gcp_file_instance", projectID, "", row) {
			continue
		}
		result := extgcp.RowToMap(row, projectID, "", tableConfig)
		resultMap = append(resultMap, result)
	}

	return resultMap, nil
}
