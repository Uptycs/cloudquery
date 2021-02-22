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

	"github.com/Uptycs/basequery-go/plugin/table"
	extgcp "github.com/Uptycs/cloudquery/extension/gcp"
	"github.com/Uptycs/cloudquery/utilities"

	"google.golang.org/api/option"

	compute "google.golang.org/api/compute/v1"
)

type myGcpComputeImagesItemsContainer struct {
	Items []*compute.Image `json:"items"`
}

// GcpComputeImagesColumns returns the list of columns for gcp_compute_image
func (handler *GcpComputeHandler) GcpComputeImagesColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("project_id"),
		table.BigIntColumn("archive_size_bytes"),
		table.TextColumn("creation_timestamp"),
		table.TextColumn("deprecated"),
		//table.TextColumn("deprecated_deleted"),
		//table.TextColumn("deprecated_deprecated"),
		//table.TextColumn("deprecated_obsolete"),
		//table.TextColumn("deprecated_replacement"),
		//table.TextColumn("deprecated_state"),
		table.TextColumn("description"),
		table.BigIntColumn("disk_size_gb"),
		table.TextColumn("family"),
		table.TextColumn("guest_os_features"),
		//table.TextColumn("guest_os_features_type"),
		table.BigIntColumn("id"),
		table.TextColumn("image_encryption_key"),
		//table.TextColumn("image_encryption_key_kms_key_name"),
		//table.TextColumn("image_encryption_key_kms_key_service_account"),
		//table.TextColumn("image_encryption_key_raw_key"),
		//table.TextColumn("image_encryption_key_sha256"),
		table.TextColumn("kind"),
		table.TextColumn("label_fingerprint"),
		table.TextColumn("labels"),
		table.TextColumn("license_codes"),
		table.TextColumn("licenses"),
		table.TextColumn("name"),
		table.TextColumn("raw_disk"),
		//table.TextColumn("raw_disk_container_type"),
		//table.TextColumn("raw_disk_sha1_checksum"),
		//table.TextColumn("raw_disk_source"),
		table.TextColumn("self_link"),
		table.TextColumn("shielded_instance_initial_state"),
		//table.TextColumn("shielded_instance_initial_state_dbs"),
		//table.TextColumn("shielded_instance_initial_state_dbs_content"),
		//table.TextColumn("shielded_instance_initial_state_dbs_file_type"),
		//table.TextColumn("shielded_instance_initial_state_dbxs"),
		//table.TextColumn("shielded_instance_initial_state_dbxs_content"),
		//table.TextColumn("shielded_instance_initial_state_dbxs_file_type"),
		//table.TextColumn("shielded_instance_initial_state_keks"),
		//table.TextColumn("shielded_instance_initial_state_keks_content"),
		//table.TextColumn("shielded_instance_initial_state_keks_file_type"),
		//table.TextColumn("shielded_instance_initial_state_pk"),
		//table.TextColumn("shielded_instance_initial_state_pk_content"),
		//table.TextColumn("shielded_instance_initial_state_pk_file_type"),
		table.TextColumn("source_disk"),
		table.TextColumn("source_disk_encryption_key"),
		//table.TextColumn("source_disk_encryption_key_kms_key_name"),
		//table.TextColumn("source_disk_encryption_key_kms_key_service_account"),
		//table.TextColumn("source_disk_encryption_key_raw_key"),
		//table.TextColumn("source_disk_encryption_key_sha256"),
		table.TextColumn("source_disk_id"),
		table.TextColumn("source_image"),
		table.TextColumn("source_image_encryption_key"),
		//table.TextColumn("source_image_encryption_key_kms_key_name"),
		//table.TextColumn("source_image_encryption_key_kms_key_service_account"),
		//table.TextColumn("source_image_encryption_key_raw_key"),
		//table.TextColumn("source_image_encryption_key_sha256"),
		table.TextColumn("source_image_id"),
		table.TextColumn("source_snapshot"),
		table.TextColumn("source_snapshot_encryption_key"),
		//table.TextColumn("source_snapshot_encryption_key_kms_key_name"),
		//table.TextColumn("source_snapshot_encryption_key_kms_key_service_account"),
		//table.TextColumn("source_snapshot_encryption_key_raw_key"),
		//table.TextColumn("source_snapshot_encryption_key_sha256"),
		table.TextColumn("source_snapshot_id"),
		table.TextColumn("source_type"),
		table.TextColumn("status"),
		table.TextColumn("storage_locations"),
	}
}

// GcpComputeImagesGenerate returns the rows in the table for all configured accounts
func (handler *GcpComputeHandler) GcpComputeImagesGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	var _ = queryContext
	ctx, cancel := context.WithCancel(osqCtx)
	defer cancel()

	resultMap := make([]map[string]string, 0)

	if len(utilities.ExtConfiguration.ExtConfGcp.Accounts) == 0 {
		results, err := handler.processAccountGcpComputeImages(ctx, nil)
		if err == nil {
			resultMap = append(resultMap, results...)
		}
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfGcp.Accounts {
			results, err := handler.processAccountGcpComputeImages(ctx, &account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}
	return resultMap, nil
}

func (handler *GcpComputeHandler) getGcpComputeImagesNewServiceForAccount(ctx context.Context, account *utilities.ExtensionConfigurationGcpAccount) (*compute.Service, string) {
	var projectID string
	var service *compute.Service
	var err error
	if account != nil && account.KeyFile != "" {
		projectID = account.ProjectID
		service, err = handler.svcInterface.NewService(ctx, option.WithCredentialsFile(account.KeyFile))
	} else if account != nil && account.ProjectID != "" {
		projectID = account.ProjectID
		service, err = handler.svcInterface.NewService(ctx)
	} else {
		projectID = utilities.DefaultGcpProjectID
		service, err = handler.svcInterface.NewService(ctx)
	}
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_compute_image",
			"projectId": projectID,
			"errString": err.Error(),
		}).Error("failed to create service")
		return nil, ""
	}
	return service, projectID
}

func (handler *GcpComputeHandler) processAccountGcpComputeImages(ctx context.Context,
	account *utilities.ExtensionConfigurationGcpAccount) ([]map[string]string, error) {

	resultMap := make([]map[string]string, 0)

	service, projectID := handler.getGcpComputeImagesNewServiceForAccount(ctx, account)
	if service == nil {
		return resultMap, fmt.Errorf("failed to initialize compute.Service")
	}
	myAPIService := handler.svcInterface.NewImagesService(service)
	if myAPIService == nil {
		return resultMap, fmt.Errorf("NewImagesService() returned nil")
	}

	aggListCall := handler.svcInterface.ImagesList(myAPIService, projectID)
	if aggListCall == nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_compute_image",
			"projectId": projectID,
		}).Debug("aggregate list call is nil")
		return resultMap, nil
	}
	itemsContainer := myGcpComputeImagesItemsContainer{Items: make([]*compute.Image, 0)}
	if err := handler.svcInterface.ImagesPages(ctx, aggListCall, func(page *compute.ImageList) error {

		itemsContainer.Items = append(itemsContainer.Items, page.Items...)

		return nil
	}); err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_compute_image",
			"projectId": projectID,
			"errString": err.Error(),
		}).Error("failed to get aggregate list page")
		return resultMap, nil
	}

	byteArr, err := json.Marshal(itemsContainer)
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_compute_image",
			"errString": err.Error(),
		}).Error("failed to marshal response")
		return resultMap, err
	}
	tableConfig, ok := utilities.TableConfigurationMap["gcp_compute_image"]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_compute_image",
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found for \"gcp_compute_image\"")
	}
	jsonTable := utilities.NewTable(byteArr, tableConfig)
	for _, row := range jsonTable.Rows {
		result := extgcp.RowToMap(row, projectID, "", tableConfig)
		resultMap = append(resultMap, result)
	}

	return resultMap, nil
}
