/**
 * Copyright (c) 2020-present, The cloudquery authors
 *
 * This source code is licensed as defined by the LICENSE file found in the
 * root directory of this source tree.
 *
 * SPDX-License-Identifier: (Apache-2.0 OR GPL-2.0-only)
 */

package container

import (
	"context"
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/Uptycs/basequery-go/plugin/table"
	extgcp "github.com/Uptycs/cloudquery/extension/gcp"
	"github.com/Uptycs/cloudquery/utilities"

	"google.golang.org/api/option"

	gcpcontainer "google.golang.org/api/container/v1beta1"
)

type myGcpContainerClustersItemsContainer struct {
	Items []*gcpcontainer.Cluster `json:"items"`
}

// GcpContainerClustersColumns returns the list of columns for gcp_container_cluster
func GcpContainerClustersColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("project_id"),
		table.TextColumn("addons_config"),
		//table.TextColumn("addons_config_cloud_run_config"),
		//table.TextColumn("addons_config_cloud_run_config_disabled"),
		//table.TextColumn("addons_config_cloud_run_config_load_balancer_type"),
		//table.TextColumn("addons_config_config_connector_config"),
		//table.TextColumn("addons_config_config_connector_config_enabled"),
		//table.TextColumn("addons_config_dns_cache_config"),
		//table.TextColumn("addons_config_dns_cache_config_enabled"),
		//table.TextColumn("addons_config_gce_persistent_disk_csi_driver_config"),
		//table.TextColumn("addons_config_gce_persistent_disk_csi_driver_config_enabled"),
		//table.TextColumn("addons_config_horizontal_pod_autoscaling"),
		//table.TextColumn("addons_config_horizontal_pod_autoscaling_disabled"),
		//table.TextColumn("addons_config_http_load_balancing"),
		//table.TextColumn("addons_config_http_load_balancing_disabled"),
		//table.TextColumn("addons_config_istio_config"),
		//table.TextColumn("addons_config_istio_config_auth"),
		//table.TextColumn("addons_config_istio_config_disabled"),
		//table.TextColumn("addons_config_kalm_config"),
		//table.TextColumn("addons_config_kalm_config_enabled"),
		//table.TextColumn("addons_config_kubernetes_dashboard"),
		//table.TextColumn("addons_config_kubernetes_dashboard_disabled"),
		//table.TextColumn("addons_config_network_policy_config"),
		//table.TextColumn("addons_config_network_policy_config_disabled"),
		table.TextColumn("authenticator_groups_config"),
		//table.TextColumn("authenticator_groups_config_enabled"),
		//table.TextColumn("authenticator_groups_config_security_group"),
		table.TextColumn("autoscaling"),
		//table.TextColumn("autoscaling_autoprovisioning_locations"),
		//table.TextColumn("autoscaling_autoprovisioning_node_pool_defaults"),
		//table.TextColumn("autoscaling_autoprovisioning_node_pool_defaults_boot_disk_kms_key"),
		//table.BigIntColumn("autoscaling_autoprovisioning_node_pool_defaults_disk_size_gb"),
		//table.TextColumn("autoscaling_autoprovisioning_node_pool_defaults_disk_type"),
		//table.TextColumn("autoscaling_autoprovisioning_node_pool_defaults_management"),
		//table.TextColumn("autoscaling_autoprovisioning_node_pool_defaults_management_auto_repair"),
		//table.TextColumn("autoscaling_autoprovisioning_node_pool_defaults_management_auto_upgrade"),
		//table.TextColumn("autoscaling_autoprovisioning_node_pool_defaults_management_upgrade_options"),
		//table.TextColumn("autoscaling_autoprovisioning_node_pool_defaults_management_upgrade_options_auto_upgrade_start_time"),
		//table.TextColumn("autoscaling_autoprovisioning_node_pool_defaults_management_upgrade_options_description"),
		//table.TextColumn("autoscaling_autoprovisioning_node_pool_defaults_min_cpu_platform"),
		//table.TextColumn("autoscaling_autoprovisioning_node_pool_defaults_oauth_scopes"),
		//table.TextColumn("autoscaling_autoprovisioning_node_pool_defaults_service_account"),
		//table.TextColumn("autoscaling_autoprovisioning_node_pool_defaults_shielded_instance_config"),
		//table.TextColumn("autoscaling_autoprovisioning_node_pool_defaults_shielded_instance_config_enable_integrity_monitoring"),
		//table.TextColumn("autoscaling_autoprovisioning_node_pool_defaults_shielded_instance_config_enable_secure_boot"),
		//table.TextColumn("autoscaling_autoprovisioning_node_pool_defaults_upgrade_settings"),
		//table.BigIntColumn("autoscaling_autoprovisioning_node_pool_defaults_upgrade_settings_max_surge"),
		//table.BigIntColumn("autoscaling_autoprovisioning_node_pool_defaults_upgrade_settings_max_unavailable"),
		//table.TextColumn("autoscaling_autoscaling_profile"),
		//table.TextColumn("autoscaling_enable_node_autoprovisioning"),
		//table.TextColumn("autoscaling_resource_limits"),
		//table.BigIntColumn("autoscaling_resource_limits_maximum"),
		//table.BigIntColumn("autoscaling_resource_limits_minimum"),
		//table.TextColumn("autoscaling_resource_limits_resource_type"),
		table.TextColumn("binary_authorization"),
		//table.TextColumn("binary_authorization_enabled"),
		table.TextColumn("cluster_ipv4_cidr"),
		table.TextColumn("cluster_telemetry"),
		//table.TextColumn("cluster_telemetry_type"),
		table.TextColumn("conditions"),
		//table.TextColumn("conditions_canonical_code"),
		//table.TextColumn("conditions_code"),
		//table.TextColumn("conditions_message"),
		//table.TextColumn("confidential_nodes"),
		//table.TextColumn("confidential_nodes_enabled"),
		table.TextColumn("create_time"),
		table.TextColumn("current_master_version"),
		table.BigIntColumn("current_node_count"),
		table.TextColumn("current_node_version"),
		table.TextColumn("database_encryption"),
		//table.TextColumn("database_encryption_key_name"),
		//table.TextColumn("database_encryption_state"),
		table.TextColumn("default_max_pods_constraint"),
		//table.BigIntColumn("default_max_pods_constraint_max_pods_per_node"),
		table.TextColumn("description"),
		table.TextColumn("enable_kubernetes_alpha"),
		table.TextColumn("enable_tpu"),
		table.TextColumn("endpoint"),
		table.TextColumn("expire_time"),
		table.TextColumn("initial_cluster_version"),
		table.BigIntColumn("initial_node_count"),
		table.TextColumn("instance_group_urls"),
		table.TextColumn("ip_allocation_policy"),
		//table.TextColumn("ip_allocation_policy_allow_route_overlap"),
		//table.TextColumn("ip_allocation_policy_cluster_ipv4_cidr"),
		//table.TextColumn("ip_allocation_policy_cluster_ipv4_cidr_block"),
		//table.TextColumn("ip_allocation_policy_cluster_secondary_range_name"),
		//table.TextColumn("ip_allocation_policy_create_subnetwork"),
		//table.TextColumn("ip_allocation_policy_node_ipv4_cidr"),
		//table.TextColumn("ip_allocation_policy_node_ipv4_cidr_block"),
		//table.TextColumn("ip_allocation_policy_services_ipv4_cidr"),
		//table.TextColumn("ip_allocation_policy_services_ipv4_cidr_block"),
		//table.TextColumn("ip_allocation_policy_services_secondary_range_name"),
		//table.TextColumn("ip_allocation_policy_subnetwork_name"),
		//table.TextColumn("ip_allocation_policy_tpu_ipv4_cidr_block"),
		//table.TextColumn("ip_allocation_policy_use_ip_aliases"),
		//table.TextColumn("ip_allocation_policy_use_routes"),
		table.TextColumn("label_fingerprint"),
		table.TextColumn("legacy_abac"),
		//table.TextColumn("legacy_abac_enabled"),
		table.TextColumn("location"),
		table.TextColumn("locations"),
		table.TextColumn("logging_service"),
		table.TextColumn("maintenance_policy"),
		//table.TextColumn("maintenance_policy_resource_version"),
		//table.TextColumn("maintenance_policy_window"),
		//table.TextColumn("maintenance_policy_window_daily_maintenance_window"),
		//table.TextColumn("maintenance_policy_window_daily_maintenance_window_duration"),
		//table.TextColumn("maintenance_policy_window_daily_maintenance_window_start_time"),
		//table.TextColumn("maintenance_policy_window_maintenance_exclusions"),
		//table.TextColumn("maintenance_policy_window_recurring_window"),
		//table.TextColumn("maintenance_policy_window_recurring_window_recurrence"),
		//table.TextColumn("maintenance_policy_window_recurring_window_window"),
		//table.TextColumn("maintenance_policy_window_recurring_window_window_end_time"),
		//table.TextColumn("maintenance_policy_window_recurring_window_window_start_time"),
		table.TextColumn("master"),
		table.TextColumn("master_auth"),
		//table.TextColumn("master_auth_client_certificate"),
		//table.TextColumn("master_auth_client_certificate_config"),
		//table.TextColumn("master_auth_client_certificate_config_issue_client_certificate"),
		//table.TextColumn("master_auth_client_key"),
		//table.TextColumn("master_auth_cluster_ca_certificate"),
		//table.TextColumn("master_auth_password"),
		//table.TextColumn("master_auth_username"),
		//table.TextColumn("master_authorized_networks_config"),
		//table.TextColumn("master_authorized_networks_config_cidr_blocks"),
		//table.TextColumn("master_authorized_networks_config_cidr_blocks_cidr_block"),
		//table.TextColumn("master_authorized_networks_config_cidr_blocks_display_name"),
		//table.TextColumn("master_authorized_networks_config_enabled"),
		table.TextColumn("master_ipv4_cidr_block"),
		table.TextColumn("monitoring_service"),
		table.TextColumn("name"),
		table.TextColumn("network"),
		table.TextColumn("network_config"),
		//table.TextColumn("network_config_datapath_provider"),
		//table.TextColumn("network_config_default_snat_status"),
		//table.TextColumn("network_config_default_snat_status_disabled"),
		//table.TextColumn("network_config_enable_intra_node_visibility"),
		//table.TextColumn("network_config_network"),
		//table.TextColumn("network_config_subnetwork"),
		table.TextColumn("network_policy"),
		//table.TextColumn("network_policy_enabled"),
		//table.TextColumn("network_policy_provider"),
		table.TextColumn("node_config"),
		//table.TextColumn("node_config_accelerators"),
		//table.BigIntColumn("node_config_accelerators_accelerator_count"),
		//table.TextColumn("node_config_accelerators_accelerator_type"),
		//table.TextColumn("node_config_boot_disk_kms_key"),
		//table.BigIntColumn("node_config_disk_size_gb"),
		//table.TextColumn("node_config_disk_type"),
		//table.TextColumn("node_config_ephemeral_storage_config"),
		//table.BigIntColumn("node_config_ephemeral_storage_config_local_ssd_count"),
		//table.TextColumn("node_config_image_type"),
		//table.TextColumn("node_config_kubelet_config"),
		//table.TextColumn("node_config_kubelet_config_cpu_cfs_quota"),
		//table.TextColumn("node_config_kubelet_config_cpu_cfs_quota_period"),
		//table.TextColumn("node_config_kubelet_config_cpu_manager_policy"),
		//table.TextColumn("node_config_labels"),
		//table.TextColumn("node_config_linux_node_config"),
		//table.TextColumn("node_config_linux_node_config_sysctls"),
		//table.BigIntColumn("node_config_local_ssd_count"),
		//table.TextColumn("node_config_machine_type"),
		//table.TextColumn("node_config_metadata"),
		//table.TextColumn("node_config_min_cpu_platform"),
		//table.TextColumn("node_config_node_group"),
		//table.TextColumn("node_config_oauth_scopes"),
		//table.TextColumn("node_config_preemptible"),
		//table.TextColumn("node_config_reservation_affinity"),
		//table.TextColumn("node_config_reservation_affinity_consume_reservation_type"),
		//table.TextColumn("node_config_reservation_affinity_key"),
		//table.TextColumn("node_config_reservation_affinity_values"),
		//table.TextColumn("node_config_sandbox_config"),
		//table.TextColumn("node_config_sandbox_config_sandbox_type"),
		//table.TextColumn("node_config_sandbox_config_type"),
		//table.TextColumn("node_config_service_account"),
		//table.TextColumn("node_config_shielded_instance_config"),
		//table.TextColumn("node_config_shielded_instance_config_enable_integrity_monitoring"),
		//table.TextColumn("node_config_shielded_instance_config_enable_secure_boot"),
		//table.TextColumn("node_config_tags"),
		//table.TextColumn("node_config_taints"),
		//table.TextColumn("node_config_taints_effect"),
		//table.TextColumn("node_config_taints_key"),
		//table.TextColumn("node_config_taints_value"),
		//table.TextColumn("node_config_workload_metadata_config"),
		//table.TextColumn("node_config_workload_metadata_config_mode"),
		//table.TextColumn("node_config_workload_metadata_config_node_metadata"),
		table.BigIntColumn("node_ipv4_cidr_size"),
		table.TextColumn("node_pools"),
		//table.TextColumn("node_pools_autoscaling"),
		//table.TextColumn("node_pools_autoscaling_autoprovisioned"),
		//table.TextColumn("node_pools_autoscaling_enabled"),
		//table.BigIntColumn("node_pools_autoscaling_max_node_count"),
		//table.BigIntColumn("node_pools_autoscaling_min_node_count"),
		//table.TextColumn("node_pools_conditions"),
		//table.TextColumn("node_pools_conditions_canonical_code"),
		//table.TextColumn("node_pools_conditions_code"),
		//table.TextColumn("node_pools_conditions_message"),
		//table.TextColumn("node_pools_config"),
		//table.TextColumn("node_pools_config_accelerators"),
		//table.BigIntColumn("node_pools_config_accelerators_accelerator_count"),
		//table.TextColumn("node_pools_config_accelerators_accelerator_type"),
		//table.TextColumn("node_pools_config_boot_disk_kms_key"),
		//table.BigIntColumn("node_pools_config_disk_size_gb"),
		//table.TextColumn("node_pools_config_disk_type"),
		//table.TextColumn("node_pools_config_ephemeral_storage_config"),
		//table.BigIntColumn("node_pools_config_ephemeral_storage_config_local_ssd_count"),
		//table.TextColumn("node_pools_config_image_type"),
		//table.TextColumn("node_pools_config_kubelet_config"),
		//table.TextColumn("node_pools_config_kubelet_config_cpu_cfs_quota"),
		//table.TextColumn("node_pools_config_kubelet_config_cpu_cfs_quota_period"),
		//table.TextColumn("node_pools_config_kubelet_config_cpu_manager_policy"),
		//table.TextColumn("node_pools_config_labels"),
		//table.TextColumn("node_pools_config_linux_node_config"),
		//table.TextColumn("node_pools_config_linux_node_config_sysctls"),
		//table.BigIntColumn("node_pools_config_local_ssd_count"),
		//table.TextColumn("node_pools_config_machine_type"),
		//table.TextColumn("node_pools_config_metadata"),
		//table.TextColumn("node_pools_config_min_cpu_platform"),
		//table.TextColumn("node_pools_config_node_group"),
		//table.TextColumn("node_pools_config_oauth_scopes"),
		//table.TextColumn("node_pools_config_preemptible"),
		//table.TextColumn("node_pools_config_reservation_affinity"),
		//table.TextColumn("node_pools_config_reservation_affinity_consume_reservation_type"),
		//table.TextColumn("node_pools_config_reservation_affinity_key"),
		//table.TextColumn("node_pools_config_reservation_affinity_values"),
		//table.TextColumn("node_pools_config_sandbox_config"),
		//table.TextColumn("node_pools_config_sandbox_config_sandbox_type"),
		//table.TextColumn("node_pools_config_sandbox_config_type"),
		//table.TextColumn("node_pools_config_service_account"),
		//table.TextColumn("node_pools_config_shielded_instance_config"),
		//table.TextColumn("node_pools_config_shielded_instance_config_enable_integrity_monitoring"),
		//table.TextColumn("node_pools_config_shielded_instance_config_enable_secure_boot"),
		//table.TextColumn("node_pools_config_tags"),
		//table.TextColumn("node_pools_config_taints"),
		//table.TextColumn("node_pools_config_taints_effect"),
		//table.TextColumn("node_pools_config_taints_key"),
		//table.TextColumn("node_pools_config_taints_value"),
		//table.TextColumn("node_pools_config_workload_metadata_config"),
		//table.TextColumn("node_pools_config_workload_metadata_config_mode"),
		//table.TextColumn("node_pools_config_workload_metadata_config_node_metadata"),
		//table.BigIntColumn("node_pools_initial_node_count"),
		//table.TextColumn("node_pools_instance_group_urls"),
		//table.TextColumn("node_pools_locations"),
		//table.TextColumn("node_pools_management"),
		//table.TextColumn("node_pools_management_auto_repair"),
		//table.TextColumn("node_pools_management_auto_upgrade"),
		//table.TextColumn("node_pools_management_upgrade_options"),
		//table.TextColumn("node_pools_management_upgrade_options_auto_upgrade_start_time"),
		//table.TextColumn("node_pools_management_upgrade_options_description"),
		//table.TextColumn("node_pools_max_pods_constraint"),
		//table.BigIntColumn("node_pools_max_pods_constraint_max_pods_per_node"),
		//table.TextColumn("node_pools_name"),
		//table.BigIntColumn("node_pools_pod_ipv4_cidr_size"),
		//table.TextColumn("node_pools_self_link"),
		//table.TextColumn("node_pools_status"),
		//table.TextColumn("node_pools_status_message"),
		//table.TextColumn("node_pools_upgrade_settings"),
		//table.BigIntColumn("node_pools_upgrade_settings_max_surge"),
		//table.BigIntColumn("node_pools_upgrade_settings_max_unavailable"),
		//table.TextColumn("node_pools_version"),
		table.TextColumn("notification_config"),
		//table.TextColumn("notification_config_pubsub"),
		//table.TextColumn("notification_config_pubsub_enabled"),
		//table.TextColumn("notification_config_pubsub_topic"),
		table.TextColumn("pod_security_policy_config"),
		//table.TextColumn("pod_security_policy_config_enabled"),
		table.TextColumn("private_cluster"),
		table.TextColumn("private_cluster_config"),
		//table.TextColumn("private_cluster_config_enable_private_endpoint"),
		//table.TextColumn("private_cluster_config_enable_private_nodes"),
		//table.TextColumn("private_cluster_config_master_global_access_config"),
		//table.TextColumn("private_cluster_config_master_global_access_config_enabled"),
		//table.TextColumn("private_cluster_config_master_ipv4_cidr_block"),
		//table.TextColumn("private_cluster_config_peering_name"),
		//table.TextColumn("private_cluster_config_private_endpoint"),
		//table.TextColumn("private_cluster_config_public_endpoint"),
		table.TextColumn("release_channel"),
		//table.TextColumn("release_channel_channel"),
		table.TextColumn("resource_labels"),
		table.TextColumn("resource_usage_export_config"),
		//table.TextColumn("resource_usage_export_config_bigquery_destination"),
		//table.TextColumn("resource_usage_export_config_bigquery_destination_dataset_id"),
		//table.TextColumn("resource_usage_export_config_consumption_metering_config"),
		//table.TextColumn("resource_usage_export_config_consumption_metering_config_enabled"),
		//table.TextColumn("resource_usage_export_config_enable_network_egress_metering"),
		//table.TextColumn("self_link"),
		table.TextColumn("services_ipv4_cidr"),
		table.TextColumn("shielded_nodes"),
		//table.TextColumn("shielded_nodes_enabled"),
		table.TextColumn("status"),
		table.TextColumn("status_message"),
		table.TextColumn("subnetwork"),
		table.TextColumn("tpu_config"),
		//table.TextColumn("tpu_config_enabled"),
		//table.TextColumn("tpu_config_ipv4_cidr_block"),
		//table.TextColumn("tpu_config_use_service_networking"),
		table.TextColumn("tpu_ipv4_cidr_block"),
		table.TextColumn("vertical_pod_autoscaling"),
		//table.TextColumn("vertical_pod_autoscaling_enabled"),
		table.TextColumn("workload_identity_config"),
		//table.TextColumn("workload_identity_config_identity_namespace"),
		//table.TextColumn("workload_identity_config_identity_provider"),
		//table.TextColumn("workload_identity_config_workload_pool"),
		table.TextColumn("zone"),
	}
}

// GcpContainerClustersGenerate returns the rows in the table for all configured accounts
func GcpContainerClustersGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	var _ = queryContext
	ctx, cancel := context.WithCancel(osqCtx)
	defer cancel()

	resultMap := make([]map[string]string, 0)

	if len(utilities.ExtConfiguration.ExtConfGcp.Accounts) == 0 {
		results, err := processAccountGcpContainerClusters(ctx, nil)
		if err == nil {
			resultMap = append(resultMap, results...)
		}
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfGcp.Accounts {
			results, err := processAccountGcpContainerClusters(ctx, &account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}
	return resultMap, nil
}

func getGcpContainerClustersNewServiceForAccount(ctx context.Context, account *utilities.ExtensionConfigurationGcpAccount) (*gcpcontainer.Service, string) {
	var projectID string
	var service *gcpcontainer.Service
	var err error
	if account != nil && account.KeyFile != "" {
		projectID = account.ProjectID
		service, err = gcpcontainer.NewService(ctx, option.WithCredentialsFile(account.KeyFile))
	} else if account != nil && account.ProjectID != "" {
		projectID = account.ProjectID
		service, err = gcpcontainer.NewService(ctx)
	} else {
		projectID = utilities.DefaultGcpProjectID
		service, err = gcpcontainer.NewService(ctx)
	}
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_container_cluster",
			"projectId": projectID,
			"errString": err.Error(),
		}).Error("failed to create service")
		return nil, ""
	}
	return service, projectID
}

func processAccountGcpContainerClusters(ctx context.Context,
	account *utilities.ExtensionConfigurationGcpAccount) ([]map[string]string, error) {

	resultMap := make([]map[string]string, 0)

	service, projectID := getGcpContainerClustersNewServiceForAccount(ctx, account)
	if service == nil {
		return resultMap, fmt.Errorf("failed to initialize gcpcontainer.Service")
	}

	listCall := service.Projects.Locations.Clusters.List("projects/" + projectID + "/locations/-")
	if listCall == nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_container_cluster",
			"projectId": projectID,
		}).Debug("list call is nil")
		return resultMap, nil
	}
	itemsContainer := myGcpContainerClustersItemsContainer{Items: make([]*gcpcontainer.Cluster, 0)}
	rsp, err := listCall.Do()
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_container_cluster",
			"projectId": projectID,
			"errString": err.Error(),
		}).Error("failed listCall.Do()")
		return resultMap, nil
	}

	itemsContainer.Items = rsp.Clusters

	byteArr, err := json.Marshal(itemsContainer)
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_container_cluster",
			"errString": err.Error(),
		}).Error("failed to marshal response")
		return resultMap, err
	}
	tableConfig, ok := utilities.TableConfigurationMap["gcp_container_cluster"]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_container_cluster",
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found for \"gcp_container_cluster\"")
	}
	jsonTable := utilities.NewTable(byteArr, tableConfig)
	for _, row := range jsonTable.Rows {
		result := extgcp.RowToMap(row, projectID, "", tableConfig)
		resultMap = append(resultMap, result)
	}

	return resultMap, nil
}
