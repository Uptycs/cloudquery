/**
 * Copyright (c) 2020-present, The cloudquery authors
 *
 * This source code is licensed as defined by the LICENSE file found in the
 * root directory of this source tree.
 *
 * SPDX-License-Identifier: (Apache-2.0 OR GPL-2.0-only)
 */

package run

import (
	"context"
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/Uptycs/basequery-go/plugin/table"
	extgcp "github.com/Uptycs/cloudquery/extension/gcp"
	"github.com/Uptycs/cloudquery/utilities"

	"google.golang.org/api/option"

	gcprun "google.golang.org/api/run/v1"
)

type myGcpCloudRunRevisionsItemsContainer struct {
	Items []*gcprun.Revision `json:"items"`
}

// GcpCloudRunRevisionsColumns returns the list of columns for gcp_cloud_run_revision
func GcpCloudRunRevisionsColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("project_id"),
		table.TextColumn("api_version"),
		table.TextColumn("kind"),
		table.TextColumn("metadata"),
		//table.TextColumn("metadata_annotations"),
		//table.TextColumn("metadata_cluster_name"),
		//table.TextColumn("metadata_creation_timestamp"),
		//table.BigIntColumn("metadata_deletion_grace_period_seconds"),
		//table.TextColumn("metadata_deletion_timestamp"),
		//table.TextColumn("metadata_finalizers"),
		//table.TextColumn("metadata_generate_name"),
		//table.BigIntColumn("metadata_generation"),
		//table.TextColumn("metadata_labels"),
		//table.TextColumn("metadata_name"),
		//table.TextColumn("metadata_namespace"),
		//table.TextColumn("metadata_owner_references"),
		//table.TextColumn("metadata_owner_references_api_version"),
		//table.TextColumn("metadata_owner_references_block_owner_deletion"),
		//table.TextColumn("metadata_owner_references_controller"),
		//table.TextColumn("metadata_owner_references_kind"),
		//table.TextColumn("metadata_owner_references_name"),
		//table.TextColumn("metadata_owner_references_uid"),
		//table.TextColumn("metadata_resource_version"),
		//table.TextColumn("metadata_self_link"),
		//table.TextColumn("metadata_uid"),
		table.TextColumn("spec"),
		//table.BigIntColumn("spec_container_concurrency"),
		//table.TextColumn("spec_containers"),
		//table.TextColumn("spec_containers_args"),
		//table.TextColumn("spec_containers_command"),
		//table.TextColumn("spec_containers_env"),
		//table.TextColumn("spec_containers_env_from"),
		//table.TextColumn("spec_containers_env_from_config_map_ref"),
		//table.TextColumn("spec_containers_env_from_config_map_ref_local_object_reference"),
		//table.TextColumn("spec_containers_env_from_config_map_ref_local_object_reference_name"),
		//table.TextColumn("spec_containers_env_from_config_map_ref_name"),
		//table.TextColumn("spec_containers_env_from_config_map_ref_optional"),
		//table.TextColumn("spec_containers_env_from_prefix"),
		//table.TextColumn("spec_containers_env_from_secret_ref"),
		//table.TextColumn("spec_containers_env_from_secret_ref_local_object_reference"),
		//table.TextColumn("spec_containers_env_from_secret_ref_local_object_reference_name"),
		//table.TextColumn("spec_containers_env_from_secret_ref_name"),
		//table.TextColumn("spec_containers_env_from_secret_ref_optional"),
		//table.TextColumn("spec_containers_env_name"),
		//table.TextColumn("spec_containers_env_value"),
		//table.TextColumn("spec_containers_env_value_from"),
		//table.TextColumn("spec_containers_env_value_from_config_map_key_ref"),
		//table.TextColumn("spec_containers_env_value_from_config_map_key_ref_key"),
		//table.TextColumn("spec_containers_env_value_from_config_map_key_ref_local_object_reference"),
		//table.TextColumn("spec_containers_env_value_from_config_map_key_ref_local_object_reference_name"),
		//table.TextColumn("spec_containers_env_value_from_config_map_key_ref_name"),
		//table.TextColumn("spec_containers_env_value_from_config_map_key_ref_optional"),
		//table.TextColumn("spec_containers_env_value_from_secret_key_ref"),
		//table.TextColumn("spec_containers_env_value_from_secret_key_ref_key"),
		//table.TextColumn("spec_containers_env_value_from_secret_key_ref_local_object_reference"),
		//table.TextColumn("spec_containers_env_value_from_secret_key_ref_local_object_reference_name"),
		//table.TextColumn("spec_containers_env_value_from_secret_key_ref_name"),
		//table.TextColumn("spec_containers_env_value_from_secret_key_ref_optional"),
		//table.TextColumn("spec_containers_image"),
		//table.TextColumn("spec_containers_image_pull_policy"),
		//table.TextColumn("spec_containers_liveness_probe"),
		//table.TextColumn("spec_containers_liveness_probe_exec"),
		//table.TextColumn("spec_containers_liveness_probe_exec_command"),
		//table.BigIntColumn("spec_containers_liveness_probe_failure_threshold"),
		//table.TextColumn("spec_containers_liveness_probe_http_get"),
		//table.TextColumn("spec_containers_liveness_probe_http_get_host"),
		//table.TextColumn("spec_containers_liveness_probe_http_get_http_headers"),
		//table.TextColumn("spec_containers_liveness_probe_http_get_http_headers_name"),
		//table.TextColumn("spec_containers_liveness_probe_http_get_http_headers_value"),
		//table.TextColumn("spec_containers_liveness_probe_http_get_path"),
		//table.TextColumn("spec_containers_liveness_probe_http_get_scheme"),
		//table.BigIntColumn("spec_containers_liveness_probe_initial_delay_seconds"),
		//table.BigIntColumn("spec_containers_liveness_probe_period_seconds"),
		//table.BigIntColumn("spec_containers_liveness_probe_success_threshold"),
		//table.TextColumn("spec_containers_liveness_probe_tcp_socket"),
		//table.TextColumn("spec_containers_liveness_probe_tcp_socket_host"),
		//table.BigIntColumn("spec_containers_liveness_probe_tcp_socket_port"),
		//table.BigIntColumn("spec_containers_liveness_probe_timeout_seconds"),
		//table.TextColumn("spec_containers_name"),
		//table.TextColumn("spec_containers_ports"),
		//table.BigIntColumn("spec_containers_ports_container_port"),
		//table.TextColumn("spec_containers_ports_name"),
		//table.TextColumn("spec_containers_ports_protocol"),
		//table.TextColumn("spec_containers_readiness_probe"),
		//table.TextColumn("spec_containers_readiness_probe_exec"),
		//table.TextColumn("spec_containers_readiness_probe_exec_command"),
		//table.BigIntColumn("spec_containers_readiness_probe_failure_threshold"),
		//table.TextColumn("spec_containers_readiness_probe_http_get"),
		//table.TextColumn("spec_containers_readiness_probe_http_get_host"),
		//table.TextColumn("spec_containers_readiness_probe_http_get_http_headers"),
		//table.TextColumn("spec_containers_readiness_probe_http_get_http_headers_name"),
		//table.TextColumn("spec_containers_readiness_probe_http_get_http_headers_value"),
		//table.TextColumn("spec_containers_readiness_probe_http_get_path"),
		//table.TextColumn("spec_containers_readiness_probe_http_get_scheme"),
		//table.BigIntColumn("spec_containers_readiness_probe_initial_delay_seconds"),
		//table.BigIntColumn("spec_containers_readiness_probe_period_seconds"),
		//table.BigIntColumn("spec_containers_readiness_probe_success_threshold"),
		//table.TextColumn("spec_containers_readiness_probe_tcp_socket"),
		//table.TextColumn("spec_containers_readiness_probe_tcp_socket_host"),
		//table.BigIntColumn("spec_containers_readiness_probe_tcp_socket_port"),
		//table.BigIntColumn("spec_containers_readiness_probe_timeout_seconds"),
		//table.TextColumn("spec_containers_resources"),
		//table.TextColumn("spec_containers_resources_limits"),
		//table.TextColumn("spec_containers_resources_requests"),
		//table.TextColumn("spec_containers_security_context"),
		//table.BigIntColumn("spec_containers_security_context_run_as_user"),
		//table.TextColumn("spec_containers_termination_message_path"),
		//table.TextColumn("spec_containers_termination_message_policy"),
		//table.TextColumn("spec_containers_volume_mounts"),
		//table.TextColumn("spec_containers_volume_mounts_mount_path"),
		//table.TextColumn("spec_containers_volume_mounts_name"),
		//table.TextColumn("spec_containers_volume_mounts_read_only"),
		//table.TextColumn("spec_containers_volume_mounts_sub_path"),
		//table.TextColumn("spec_containers_working_dir"),
		//table.TextColumn("spec_service_account_name"),
		//table.BigIntColumn("spec_timeout_seconds"),
		//table.TextColumn("spec_volumes"),
		//table.TextColumn("spec_volumes_config_map"),
		//table.BigIntColumn("spec_volumes_config_map_default_mode"),
		//table.TextColumn("spec_volumes_config_map_items"),
		//table.TextColumn("spec_volumes_config_map_items_key"),
		//table.BigIntColumn("spec_volumes_config_map_items_mode"),
		//table.TextColumn("spec_volumes_config_map_items_path"),
		//table.TextColumn("spec_volumes_config_map_name"),
		//table.TextColumn("spec_volumes_config_map_optional"),
		//table.TextColumn("spec_volumes_name"),
		//table.TextColumn("spec_volumes_secret"),
		//table.BigIntColumn("spec_volumes_secret_default_mode"),
		//table.TextColumn("spec_volumes_secret_items"),
		//table.TextColumn("spec_volumes_secret_items_key"),
		//table.BigIntColumn("spec_volumes_secret_items_mode"),
		//table.TextColumn("spec_volumes_secret_items_path"),
		//table.TextColumn("spec_volumes_secret_optional"),
		//table.TextColumn("spec_volumes_secret_secret_name"),
		table.TextColumn("status"),
		//table.TextColumn("status_conditions"),
		//table.TextColumn("status_conditions_last_transition_time"),
		//table.TextColumn("status_conditions_message"),
		//table.TextColumn("status_conditions_reason"),
		//table.TextColumn("status_conditions_severity"),
		//table.TextColumn("status_conditions_status"),
		//table.TextColumn("status_conditions_type"),
		//table.TextColumn("status_image_digest"),
		//table.TextColumn("status_log_url"),
		//table.BigIntColumn("status_observed_generation"),
		//table.TextColumn("status_service_name"),

	}
}

// GcpCloudRunRevisionsGenerate returns the rows in the table for all configured accounts
func GcpCloudRunRevisionsGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	var _ = queryContext
	ctx, cancel := context.WithCancel(osqCtx)
	defer cancel()

	resultMap := make([]map[string]string, 0)

	if len(utilities.ExtConfiguration.ExtConfGcp.Accounts) == 0 {
		results, err := processAccountGcpCloudRunRevisions(ctx, nil)
		if err == nil {
			resultMap = append(resultMap, results...)
		}
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfGcp.Accounts {
			results, err := processAccountGcpCloudRunRevisions(ctx, &account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}
	return resultMap, nil
}

func getGcpCloudRunRevisionsNewServiceForAccount(ctx context.Context, account *utilities.ExtensionConfigurationGcpAccount) (*gcprun.APIService, string) {
	var projectID string
	var service *gcprun.APIService
	var err error
	if account != nil && account.KeyFile != "" {
		projectID = account.ProjectID
		service, err = gcprun.NewService(ctx, option.WithCredentialsFile(account.KeyFile))
	} else if account != nil && account.ProjectID != "" {
		projectID = account.ProjectID
		service, err = gcprun.NewService(ctx)
	} else {
		projectID = utilities.DefaultGcpProjectID
		service, err = gcprun.NewService(ctx)
	}
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_cloud_run_revision",
			"projectId": projectID,
			"errString": err.Error(),
		}).Error("failed to create service")
		return nil, ""
	}
	return service, projectID
}

func processAccountGcpCloudRunRevisions(ctx context.Context,
	account *utilities.ExtensionConfigurationGcpAccount) ([]map[string]string, error) {

	resultMap := make([]map[string]string, 0)

	service, projectID := getGcpCloudRunRevisionsNewServiceForAccount(ctx, account)
	if service == nil {
		return resultMap, fmt.Errorf("failed to initialize gcprun.APIService")
	}

	listCall := service.Projects.Locations.Revisions.List("projects/" + projectID + "/locations/-")
	if listCall == nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_cloud_run_revision",
			"projectId": projectID,
		}).Debug("list call is nil")
		return resultMap, nil
	}
	itemsContainer := myGcpCloudRunRevisionsItemsContainer{Items: make([]*gcprun.Revision, 0)}
	rsp, err := listCall.Do()
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_cloud_run_revision",
			"projectId": projectID,
			"errString": err.Error(),
		}).Error("failed listCall.Do()")
		return resultMap, nil
	}

	itemsContainer.Items = rsp.Items

	byteArr, err := json.Marshal(itemsContainer)
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_cloud_run_revision",
			"errString": err.Error(),
		}).Error("failed to marshal response")
		return resultMap, err
	}
	tableConfig, ok := utilities.TableConfigurationMap["gcp_cloud_run_revision"]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_cloud_run_revision",
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found for \"gcp_cloud_run_revision\"")
	}
	jsonTable := utilities.NewTable(byteArr, tableConfig)
	for _, row := range jsonTable.Rows {
		result := extgcp.RowToMap(row, projectID, "", tableConfig)
		resultMap = append(resultMap, result)
	}

	return resultMap, nil
}
