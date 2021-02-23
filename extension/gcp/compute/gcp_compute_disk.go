/**
 * Copyright (c) 2020-present, The cloudquery authors
 *
 * This source code is licensed as defined by the LICENSE file found in the
 * root directory of this source tree.
 *
 * SPDX-License-Identifier: (Apache-2.0 OR GPL-2.0-only)
 */

package compute

import (
	"context"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"strings"

	"github.com/Uptycs/basequery-go/plugin/table"
	extgcp "github.com/Uptycs/cloudquery/extension/gcp"
	"github.com/Uptycs/cloudquery/utilities"

	"google.golang.org/api/option"

	compute "google.golang.org/api/compute/v1"
)

type myGcpComputeDisksItemsContainer struct {
	Items []*compute.Disk `json:"items"`
}

// GcpComputeDisksColumns returns the list of columns for gcp_compute_disk
func (handler *GcpComputeHandler) GcpComputeDisksColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("project_id"),
		table.TextColumn("creation_timestamp"),
		table.TextColumn("description"),
		table.TextColumn("disk_encryption_key"),
		//table.TextColumn("disk_encryption_key_kms_key_name"),
		//table.TextColumn("disk_encryption_key_kms_key_service_account"),
		//table.TextColumn("disk_encryption_key_raw_key"),
		//table.TextColumn("disk_encryption_key_sha256"),
		table.TextColumn("guest_os_features"),
		//table.TextColumn("guest_os_features_type"),
		table.BigIntColumn("id"),
		table.TextColumn("kind"),
		table.TextColumn("label_fingerprint"),
		table.TextColumn("labels"),
		table.TextColumn("last_attach_timestamp"),
		table.TextColumn("last_detach_timestamp"),
		table.TextColumn("license_codes"),
		table.TextColumn("licenses"),
		table.TextColumn("name"),
		table.TextColumn("options"),
		table.BigIntColumn("physical_block_size_bytes"),
		table.TextColumn("region"),
		table.TextColumn("replica_zones"),
		table.TextColumn("resource_policies"),
		//table.TextColumn("self_link"),
		table.BigIntColumn("size_gb"),
		table.TextColumn("source_disk"),
		table.TextColumn("source_disk_id"),
		table.TextColumn("source_image"),
		table.TextColumn("source_image_encryption_key"),
		//table.TextColumn("source_image_encryption_key_kms_key_name"),
		//table.TextColumn("source_image_encryption_key_kms_key_service_account"),
		//table.TextColumn("source_image_encryption_key_raw_key"),
		//table.TextColumn("source_image_encryption_key_sha256"),
		table.TextColumn("source_image_id"),
		table.TextColumn("source_snapshot"),
		//table.TextColumn("source_snapshot_encryption_key"),
		//table.TextColumn("source_snapshot_encryption_key_kms_key_name"),
		//table.TextColumn("source_snapshot_encryption_key_kms_key_service_account"),
		//table.TextColumn("source_snapshot_encryption_key_raw_key"),
		//table.TextColumn("source_snapshot_encryption_key_sha256"),
		table.TextColumn("source_snapshot_id"),
		table.TextColumn("status"),
		table.TextColumn("type"),
		table.TextColumn("users"),
		table.TextColumn("zone"),
	}
}

// GcpComputeDisksGenerate returns the rows in the table for all configured accounts
func (handler *GcpComputeHandler) GcpComputeDisksGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	ctx, cancel := context.WithCancel(osqCtx)
	defer cancel()

	resultMap := make([]map[string]string, 0)

	if len(utilities.ExtConfiguration.ExtConfGcp.Accounts) == 0 && extgcp.ShouldProcessProject("gcp_compute_disk", utilities.DefaultGcpProjectID) {
		results, err := handler.processAccountGcpComputeDisks(ctx, queryContext, nil)
		if err == nil {
			resultMap = append(resultMap, results...)
		}
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfGcp.Accounts {
			if !extgcp.ShouldProcessProject("gcp_compute_disk", account.ProjectID) {
				continue
			}
			results, err := handler.processAccountGcpComputeDisks(ctx, queryContext, &account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}
	return resultMap, nil
}

func (handler *GcpComputeHandler) getGcpComputeDisksNewServiceForAccount(ctx context.Context, account *utilities.ExtensionConfigurationGcpAccount) (*compute.Service, string) {
	var projectID string
	var service *compute.Service
	var err error
	if account != nil {
		projectID = account.ProjectID
		service, err = handler.svcInterface.NewService(ctx, option.WithCredentialsFile(account.KeyFile))
	} else {
		projectID = utilities.DefaultGcpProjectID
		service, err = handler.svcInterface.NewService(ctx)
	}
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_compute_disk",
			"projectId": projectID,
			"errString": err.Error(),
		}).Error("failed to create service")
		return nil, ""
	}
	return service, projectID
}

func (handler *GcpComputeHandler) processAccountGcpComputeDisks(ctx context.Context, queryContext table.QueryContext,
	account *utilities.ExtensionConfigurationGcpAccount) ([]map[string]string, error) {

	resultMap := make([]map[string]string, 0)

	service, projectID := handler.getGcpComputeDisksNewServiceForAccount(ctx, account)
	if service == nil {
		return resultMap, fmt.Errorf("failed to initialize compute.Service")
	}
	myAPIService := handler.svcInterface.NewDisksService(service)
	if myAPIService == nil {
		return resultMap, fmt.Errorf("NewDisksService() returned nil")
	}

	aggListCall := handler.svcInterface.DisksAggregatedList(myAPIService, projectID)
	if aggListCall == nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_compute_disk",
			"projectId": projectID,
		}).Debug("aggregate list call is nil")
		return resultMap, nil
	}
	itemsContainer := myGcpComputeDisksItemsContainer{Items: make([]*compute.Disk, 0)}
	if err := handler.svcInterface.DisksPages(ctx, aggListCall, func(page *compute.DiskAggregatedList) error {

		for _, item := range page.Items {
			for _, inst := range item.Disks {
				zonePathSplit := strings.Split(inst.Zone, "/")
				inst.Zone = zonePathSplit[len(zonePathSplit)-1]
			}
			itemsContainer.Items = append(itemsContainer.Items, item.Disks...)
		}

		return nil
	}); err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_compute_disk",
			"projectId": projectID,
			"errString": err.Error(),
		}).Error("failed to get aggregate list page")
		return resultMap, nil
	}

	byteArr, err := json.Marshal(itemsContainer)
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_compute_disk",
			"errString": err.Error(),
		}).Error("failed to marshal response")
		return resultMap, err
	}
	tableConfig, ok := utilities.TableConfigurationMap["gcp_compute_disk"]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_compute_disk",
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found for \"gcp_compute_disk\"")
	}
	jsonTable := utilities.NewTable(byteArr, tableConfig)
	for _, row := range jsonTable.Rows {
		if !extgcp.ShouldProcessRow(ctx, queryContext, "gcp_compute_disk", projectID, "", row) {
			continue
		}
		result := extgcp.RowToMap(row, projectID, "", tableConfig)
		resultMap = append(resultMap, result)
	}

	return resultMap, nil
}
