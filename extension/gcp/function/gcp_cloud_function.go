/**
 * Copyright (c) 2020-present, The cloudquery authors
 *
 * This source code is licensed as defined by the LICENSE file found in the
 * root directory of this source tree.
 *
 * SPDX-License-Identifier: (Apache-2.0 OR GPL-2.0-only)
 */

package function

import (
	"context"
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/Uptycs/basequery-go/plugin/table"
	extgcp "github.com/Uptycs/cloudquery/extension/gcp"
	"github.com/Uptycs/cloudquery/utilities"

	"google.golang.org/api/option"

	gcpfunction "google.golang.org/api/cloudfunctions/v1beta2"
)

type myGcpCloudFunctionsItemsContainer struct {
	Items []*gcpfunction.CloudFunction `json:"items"`
}

// GcpCloudFunctionsColumns returns the list of columns for gcp_cloud_function
func GcpCloudFunctionsColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("project_id"),
		table.BigIntColumn("available_memory_mb"),
		table.TextColumn("entry_point"),
		table.TextColumn("environment_variables"),
		table.TextColumn("event_trigger"),
		//table.TextColumn("event_trigger_event_type"),
		//table.TextColumn("event_trigger_failure_policy"),
		//table.TextColumn("event_trigger_failure_policy_retry"),
		//table.TextColumn("event_trigger_resource"),
		//table.TextColumn("event_trigger_service"),
		table.TextColumn("https_trigger"),
		//table.TextColumn("https_trigger_url"),
		table.TextColumn("labels"),
		table.TextColumn("latest_operation"),
		table.BigIntColumn("max_instances"),
		table.TextColumn("name"),
		table.TextColumn("network"),
		table.TextColumn("runtime"),
		table.TextColumn("service_account"),
		table.TextColumn("source_archive_url"),
		table.TextColumn("source_repository"),
		table.TextColumn("source_repository_url"),
		//table.TextColumn("source_repository_branch"),
		//table.TextColumn("source_repository_deployed_revision"),
		//table.TextColumn("source_repository_repository_url"),
		//table.TextColumn("source_repository_revision"),
		//table.TextColumn("source_repository_source_path"),
		//table.TextColumn("source_repository_tag"),
		table.TextColumn("source_upload_url"),
		table.TextColumn("status"),
		table.TextColumn("timeout"),
		table.TextColumn("update_time"),
		table.BigIntColumn("version_id"),
		table.TextColumn("vpc_connector"),
	}
}

// GcpCloudFunctionsGenerate returns the rows in the table for all configured accounts
func GcpCloudFunctionsGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	var _ = queryContext
	ctx, cancel := context.WithCancel(osqCtx)
	defer cancel()

	resultMap := make([]map[string]string, 0)

	if len(utilities.ExtConfiguration.ExtConfGcp.Accounts) == 0 {
		results, err := processAccountGcpCloudFunctions(ctx, nil)
		if err == nil {
			resultMap = append(resultMap, results...)
		}
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfGcp.Accounts {
			results, err := processAccountGcpCloudFunctions(ctx, &account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}
	return resultMap, nil
}

func getGcpCloudFunctionsNewServiceForAccount(ctx context.Context, account *utilities.ExtensionConfigurationGcpAccount) (*gcpfunction.Service, string) {
	var projectID string
	var service *gcpfunction.Service
	var err error
	if account != nil {
		projectID = account.ProjectID
		service, err = gcpfunction.NewService(ctx, option.WithCredentialsFile(account.KeyFile))
	} else {
		projectID = utilities.DefaultGcpProjectID
		service, err = gcpfunction.NewService(ctx)
	}
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_cloud_function",
			"projectId": projectID,
			"errString": err.Error(),
		}).Error("failed to create service")
		return nil, ""
	}
	return service, projectID
}

func processAccountGcpCloudFunctions(ctx context.Context,
	account *utilities.ExtensionConfigurationGcpAccount) ([]map[string]string, error) {

	resultMap := make([]map[string]string, 0)

	service, projectID := getGcpCloudFunctionsNewServiceForAccount(ctx, account)
	if service == nil {
		return resultMap, fmt.Errorf("failed to initialize gcpfunction.Service")
	}

	listCall := service.Projects.Locations.Functions.List("projects/" + projectID + "/locations/-")
	if listCall == nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_cloud_function",
			"projectId": projectID,
		}).Debug("list call is nil")
		return resultMap, nil
	}
	itemsContainer := myGcpCloudFunctionsItemsContainer{Items: make([]*gcpfunction.CloudFunction, 0)}
	rsp, err := listCall.Do()
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_cloud_function",
			"projectId": projectID,
			"errString": err.Error(),
		}).Error("failed listCall.Do()")
		return resultMap, nil
	}

	itemsContainer.Items = rsp.Functions

	byteArr, err := json.Marshal(itemsContainer)
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_cloud_function",
			"errString": err.Error(),
		}).Error("failed to marshal response")
		return resultMap, err
	}
	tableConfig, ok := utilities.TableConfigurationMap["gcp_cloud_function"]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_cloud_function",
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found for \"gcp_cloud_function\"")
	}
	jsonTable := utilities.NewTable(byteArr, tableConfig)
	for _, row := range jsonTable.Rows {
		result := extgcp.RowToMap(row, projectID, "", tableConfig)
		resultMap = append(resultMap, result)
	}

	return resultMap, nil
}
