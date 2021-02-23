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

type myGcpComputeInstancesItemsContainer struct {
	Items []*compute.Instance `json:"items"`
}

// GcpComputeInstancesColumns returns the list of columns for gcp_compute_instance
func (handler *GcpComputeHandler) GcpComputeInstancesColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("project_id"),
		table.TextColumn("can_ip_forward"),
		table.TextColumn("confidential_instance_config"),
		//table.TextColumn("confidential_instance_config_enable_confidential_compute"),
		table.TextColumn("cpu_platform"),
		table.TextColumn("creation_timestamp"),
		table.TextColumn("deletion_protection"),
		table.TextColumn("description"),
		table.TextColumn("disks"),
		//table.TextColumn("disks_auto_delete"),
		//table.TextColumn("disks_boot"),
		//table.TextColumn("disks_device_name"),
		//table.TextColumn("disks_disk_encryption_key"),
		//table.TextColumn("disks_disk_encryption_key_kms_key_name"),
		//table.TextColumn("disks_disk_encryption_key_kms_key_service_account"),
		//table.TextColumn("disks_disk_encryption_key_raw_key"),
		//table.TextColumn("disks_disk_encryption_key_sha256"),
		//table.BigIntColumn("disks_disk_size_gb"),
		//table.TextColumn("disks_guest_os_features"),
		//table.TextColumn("disks_guest_os_features_type"),
		//table.BigIntColumn("disks_index"),
		//table.TextColumn("disks_initialize_params"),
		//table.TextColumn("disks_initialize_params_description"),
		//table.TextColumn("disks_initialize_params_disk_name"),
		//table.BigIntColumn("disks_initialize_params_disk_size_gb"),
		//table.TextColumn("disks_initialize_params_disk_type"),
		//table.TextColumn("disks_initialize_params_labels"),
		//table.TextColumn("disks_initialize_params_on_update_action"),
		//table.TextColumn("disks_initialize_params_resource_policies"),
		//table.TextColumn("disks_initialize_params_source_image"),
		//table.TextColumn("disks_initialize_params_source_image_encryption_key"),
		//table.TextColumn("disks_initialize_params_source_image_encryption_key_kms_key_name"),
		//table.TextColumn("disks_initialize_params_source_image_encryption_key_kms_key_service_account"),
		//table.TextColumn("disks_initialize_params_source_image_encryption_key_raw_key"),
		//table.TextColumn("disks_initialize_params_source_image_encryption_key_sha256"),
		//table.TextColumn("disks_initialize_params_source_snapshot"),
		//table.TextColumn("disks_initialize_params_source_snapshot_encryption_key"),
		//table.TextColumn("disks_initialize_params_source_snapshot_encryption_key_kms_key_name"),
		//table.TextColumn("disks_initialize_params_source_snapshot_encryption_key_kms_key_service_account"),
		//table.TextColumn("disks_initialize_params_source_snapshot_encryption_key_raw_key"),
		//table.TextColumn("disks_initialize_params_source_snapshot_encryption_key_sha256"),
		//table.TextColumn("disks_interface"),
		//table.TextColumn("disks_kind"),
		//table.TextColumn("disks_licenses"),
		//table.TextColumn("disks_mode"),
		//table.TextColumn("disks_shielded_instance_initial_state"),
		//table.TextColumn("disks_shielded_instance_initial_state_dbs"),
		//table.TextColumn("disks_shielded_instance_initial_state_dbs_content"),
		//table.TextColumn("disks_shielded_instance_initial_state_dbs_file_type"),
		//table.TextColumn("disks_shielded_instance_initial_state_dbxs"),
		//table.TextColumn("disks_shielded_instance_initial_state_dbxs_content"),
		//table.TextColumn("disks_shielded_instance_initial_state_dbxs_file_type"),
		//table.TextColumn("disks_shielded_instance_initial_state_keks"),
		//table.TextColumn("disks_shielded_instance_initial_state_keks_content"),
		//table.TextColumn("disks_shielded_instance_initial_state_keks_file_type"),
		//table.TextColumn("disks_shielded_instance_initial_state_pk"),
		//table.TextColumn("disks_shielded_instance_initial_state_pk_content"),
		//table.TextColumn("disks_shielded_instance_initial_state_pk_file_type"),
		//table.TextColumn("disks_source"),
		//table.TextColumn("disks_type"),
		table.TextColumn("display_device"),
		//table.TextColumn("display_device_enable_display"),
		table.TextColumn("fingerprint"),
		table.TextColumn("guest_accelerators"),
		//table.BigIntColumn("guest_accelerators_accelerator_count"),
		//table.TextColumn("guest_accelerators_accelerator_type"),
		table.TextColumn("hostname"),
		table.BigIntColumn("id"),
		table.TextColumn("kind"),
		table.TextColumn("label_fingerprint"),
		table.TextColumn("labels"),
		table.TextColumn("last_start_timestamp"),
		table.TextColumn("last_stop_timestamp"),
		table.TextColumn("last_suspended_timestamp"),
		table.TextColumn("machine_type"),
		table.TextColumn("metadata"),
		//table.TextColumn("metadata_fingerprint"),
		//table.TextColumn("metadata_items"),
		//table.TextColumn("metadata_items_key"),
		//table.TextColumn("metadata_items_value"),
		//table.TextColumn("metadata_kind"),
		table.TextColumn("min_cpu_platform"),
		table.TextColumn("name"),
		table.TextColumn("network_interfaces"),
		//table.TextColumn("network_interfaces_access_configs"),
		//table.TextColumn("network_interfaces_access_configs_kind"),
		//table.TextColumn("network_interfaces_access_configs_name"),
		//table.TextColumn("network_interfaces_access_configs_nat_ip"),
		//table.TextColumn("network_interfaces_access_configs_network_tier"),
		//table.TextColumn("network_interfaces_access_configs_public_ptr_domain_name"),
		//table.TextColumn("network_interfaces_access_configs_set_public_ptr"),
		//table.TextColumn("network_interfaces_access_configs_type"),
		//table.TextColumn("network_interfaces_alias_ip_ranges"),
		//table.TextColumn("network_interfaces_alias_ip_ranges_ip_cidr_range"),
		//table.TextColumn("network_interfaces_alias_ip_ranges_subnetwork_range_name"),
		//table.TextColumn("network_interfaces_fingerprint"),
		//table.TextColumn("network_interfaces_ipv6_address"),
		//table.TextColumn("network_interfaces_kind"),
		//table.TextColumn("network_interfaces_name"),
		//table.TextColumn("network_interfaces_network"),
		//table.TextColumn("network_interfaces_network_ip"),
		//table.TextColumn("network_interfaces_subnetwork"),
		table.TextColumn("private_ipv6_google_access"),
		table.TextColumn("reservation_affinity"),
		//table.TextColumn("reservation_affinity_consume_reservation_type"),
		//table.TextColumn("reservation_affinity_key"),
		//table.TextColumn("reservation_affinity_values"),
		table.TextColumn("resource_policies"),
		table.TextColumn("scheduling"),
		//table.TextColumn("scheduling_automatic_restart"),
		//table.BigIntColumn("scheduling_min_node_cpus"),
		//table.TextColumn("scheduling_node_affinities"),
		//table.TextColumn("scheduling_node_affinities_key"),
		//table.TextColumn("scheduling_node_affinities_operator"),
		//table.TextColumn("scheduling_node_affinities_values"),
		//table.TextColumn("scheduling_on_host_maintenance"),
		//table.TextColumn("scheduling_preemptible"),
		//table.TextColumn("self_link"),
		table.TextColumn("service_accounts"),
		//table.TextColumn("service_accounts_email"),
		//table.TextColumn("service_accounts_scopes"),
		//table.TextColumn("shielded_instance_config"),
		//table.TextColumn("shielded_instance_config_enable_integrity_monitoring"),
		//table.TextColumn("shielded_instance_config_enable_secure_boot"),
		//table.TextColumn("shielded_instance_config_enable_vtpm"),
		//table.TextColumn("shielded_instance_integrity_policy"),
		//table.TextColumn("shielded_instance_integrity_policy_update_auto_learn_policy"),
		table.TextColumn("start_restricted"),
		table.TextColumn("status"),
		table.TextColumn("status_message"),
		table.TextColumn("tags"),
		//table.TextColumn("tags_fingerprint"),
		//table.TextColumn("tags_items"),
		table.TextColumn("zone"),
	}
}

// GcpComputeInstancesGenerate returns the rows in the table for all configured accounts
func (handler *GcpComputeHandler) GcpComputeInstancesGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	var _ = queryContext
	ctx, cancel := context.WithCancel(osqCtx)
	defer cancel()

	resultMap := make([]map[string]string, 0)

	if len(utilities.ExtConfiguration.ExtConfGcp.Accounts) == 0 {
		results, err := handler.processAccountGcpComputeInstances(ctx, nil)
		if err == nil {
			resultMap = append(resultMap, results...)
		}
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfGcp.Accounts {
			results, err := handler.processAccountGcpComputeInstances(ctx, &account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}
	return resultMap, nil
}

func (handler *GcpComputeHandler) getGcpComputeInstancesNewServiceForAccount(ctx context.Context, account *utilities.ExtensionConfigurationGcpAccount) (*compute.Service, string) {
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
			"tableName": "gcp_compute_instance",
			"projectId": projectID,
			"errString": err.Error(),
		}).Error("failed to create service")
		return nil, ""
	}
	return service, projectID
}

func (handler *GcpComputeHandler) processAccountGcpComputeInstances(ctx context.Context,
	account *utilities.ExtensionConfigurationGcpAccount) ([]map[string]string, error) {

	resultMap := make([]map[string]string, 0)

	service, projectID := handler.getGcpComputeInstancesNewServiceForAccount(ctx, account)
	if service == nil {
		return resultMap, fmt.Errorf("failed to initialize compute.Service")
	}
	myAPIService := handler.svcInterface.NewInstancesService(service)
	if myAPIService == nil {
		return resultMap, fmt.Errorf("NewInstancesService() returned nil")
	}

	aggListCall := handler.svcInterface.InstancesAggregatedList(myAPIService, projectID)
	if aggListCall == nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_compute_instance",
			"projectId": projectID,
		}).Debug("aggregate list call is nil")
		return resultMap, nil
	}
	itemsContainer := myGcpComputeInstancesItemsContainer{Items: make([]*compute.Instance, 0)}
	if err := handler.svcInterface.InstancesPages(ctx, aggListCall, func(page *compute.InstanceAggregatedList) error {

		for _, item := range page.Items {
			for _, inst := range item.Instances {
				zonePathSplit := strings.Split(inst.Zone, "/")
				inst.Zone = zonePathSplit[len(zonePathSplit)-1]
			}
			itemsContainer.Items = append(itemsContainer.Items, item.Instances...)
		}

		return nil
	}); err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_compute_instance",
			"projectId": projectID,
			"errString": err.Error(),
		}).Error("failed to get aggregate list page")
		return resultMap, nil
	}

	byteArr, err := json.Marshal(itemsContainer)
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_compute_instance",
			"errString": err.Error(),
		}).Error("failed to marshal response")
		return resultMap, err
	}
	tableConfig, ok := utilities.TableConfigurationMap["gcp_compute_instance"]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_compute_instance",
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found for \"gcp_compute_instance\"")
	}
	jsonTable := utilities.NewTable(byteArr, tableConfig)
	for _, row := range jsonTable.Rows {
		result := extgcp.RowToMap(row, projectID, "", tableConfig)
		resultMap = append(resultMap, result)
	}

	return resultMap, nil
}
