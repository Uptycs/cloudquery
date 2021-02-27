/**
 * Copyright (c) 2020-present, The cloudquery authors
 *
 * This source code is licensed as defined by the LICENSE file found in the
 * root directory of this source tree.
 *
 * SPDX-License-Identifier: (Apache-2.0 OR GPL-2.0-only)
 */

package sql

import (
	"context"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"

	"github.com/Uptycs/basequery-go/plugin/table"
	extgcp "github.com/Uptycs/cloudquery/extension/gcp"
	"github.com/Uptycs/cloudquery/utilities"

	"google.golang.org/api/option"

	gcpsql "google.golang.org/api/sqladmin/v1beta4"
)

type myGcpSQLInstancesItemsContainer struct {
	Items []*gcpsql.DatabaseInstance `json:"items"`
}

// GcpSQLInstancesColumns returns the list of columns for gcp_sql_instance
func GcpSQLInstancesColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("project_id"),
		table.TextColumn("backend_type"),
		table.TextColumn("connection_name"),
		table.BigIntColumn("current_disk_size"),
		table.TextColumn("database_version"),
		table.TextColumn("disk_encryption_configuration"),
		//table.TextColumn("disk_encryption_configuration_kind"),
		//table.TextColumn("disk_encryption_configuration_kms_key_name"),
		table.TextColumn("disk_encryption_status"),
		//table.TextColumn("disk_encryption_status_kind"),
		//table.TextColumn("disk_encryption_status_kms_key_version_name"),
		table.TextColumn("etag"),
		table.TextColumn("failover_replica"),
		//table.TextColumn("failover_replica_available"),
		//table.TextColumn("failover_replica_name"),
		table.TextColumn("gce_zone"),
		table.TextColumn("instance_type"),
		table.TextColumn("ip_addresses"),
		//table.TextColumn("ip_addresses_ip_address"),
		//table.TextColumn("ip_addresses_time_to_retire"),
		//table.TextColumn("ip_addresses_type"),
		table.TextColumn("ipv6_address"),
		table.TextColumn("kind"),
		table.TextColumn("master_instance_name"),
		table.BigIntColumn("max_disk_size"),
		table.TextColumn("name"),
		table.TextColumn("on_premises_configuration"),
		//table.TextColumn("on_premises_configuration_ca_certificate"),
		//table.TextColumn("on_premises_configuration_client_certificate"),
		//table.TextColumn("on_premises_configuration_client_key"),
		//table.TextColumn("on_premises_configuration_dump_file_path"),
		//table.TextColumn("on_premises_configuration_host_port"),
		//table.TextColumn("on_premises_configuration_kind"),
		//table.TextColumn("on_premises_configuration_password"),
		//table.TextColumn("on_premises_configuration_username"),
		table.TextColumn("project"),
		table.TextColumn("region"),
		table.TextColumn("replica_configuration"),
		//table.TextColumn("replica_configuration_failover_target"),
		//table.TextColumn("replica_configuration_kind"),
		//table.TextColumn("replica_configuration_mysql_replica_configuration"),
		//table.TextColumn("replica_configuration_mysql_replica_configuration_ca_certificate"),
		//table.TextColumn("replica_configuration_mysql_replica_configuration_client_certificate"),
		//table.TextColumn("replica_configuration_mysql_replica_configuration_client_key"),
		//table.BigIntColumn("replica_configuration_mysql_replica_configuration_connect_retry_interval"),
		//table.TextColumn("replica_configuration_mysql_replica_configuration_dump_file_path"),
		//table.TextColumn("replica_configuration_mysql_replica_configuration_kind"),
		//table.BigIntColumn("replica_configuration_mysql_replica_configuration_master_heartbeat_period"),
		//table.TextColumn("replica_configuration_mysql_replica_configuration_password"),
		//table.TextColumn("replica_configuration_mysql_replica_configuration_ssl_cipher"),
		//table.TextColumn("replica_configuration_mysql_replica_configuration_username"),
		//table.TextColumn("replica_configuration_mysql_replica_configuration_verify_server_certificate"),
		table.TextColumn("replica_names"),
		//table.TextColumn("root_password"),
		table.TextColumn("satisfies_pzs"),
		table.TextColumn("scheduled_maintenance"),
		//table.TextColumn("scheduled_maintenance_can_defer"),
		//table.TextColumn("scheduled_maintenance_can_reschedule"),
		//table.TextColumn("scheduled_maintenance_start_time"),
		table.TextColumn("secondary_gce_zone"),
		//table.TextColumn("self_link"),
		//table.TextColumn("server_ca_cert"),
		//table.TextColumn("server_ca_cert_cert"),
		//table.TextColumn("server_ca_cert_cert_serial_number"),
		//table.TextColumn("server_ca_cert_common_name"),
		//table.TextColumn("server_ca_cert_create_time"),
		//table.TextColumn("server_ca_cert_expiration_time"),
		//table.TextColumn("server_ca_cert_instance"),
		//table.TextColumn("server_ca_cert_kind"),
		//table.TextColumn("server_ca_cert_self_link"),
		//table.TextColumn("server_ca_cert_sha1_fingerprint"),
		//table.TextColumn("service_account_email_address"),
		table.TextColumn("settings"),
		//table.TextColumn("settings_activation_policy"),
		//table.TextColumn("settings_active_directory_config"),
		//table.TextColumn("settings_active_directory_config_domain"),
		//table.TextColumn("settings_active_directory_config_kind"),
		//table.TextColumn("settings_authorized_gae_applications"),
		//table.TextColumn("settings_availability_type"),
		//table.TextColumn("settings_backup_configuration"),
		//table.TextColumn("settings_backup_configuration_backup_retention_settings"),
		//table.BigIntColumn("settings_backup_configuration_backup_retention_settings_retained_backups"),
		//table.TextColumn("settings_backup_configuration_backup_retention_settings_retention_unit"),
		//table.TextColumn("settings_backup_configuration_binary_log_enabled"),
		//table.TextColumn("settings_backup_configuration_enabled"),
		//table.TextColumn("settings_backup_configuration_kind"),
		//table.TextColumn("settings_backup_configuration_location"),
		//table.TextColumn("settings_backup_configuration_point_in_time_recovery_enabled"),
		//table.TextColumn("settings_backup_configuration_replication_log_archiving_enabled"),
		//table.TextColumn("settings_backup_configuration_start_time"),
		//table.BigIntColumn("settings_backup_configuration_transaction_log_retention_days"),
		//table.TextColumn("settings_collation"),
		//table.TextColumn("settings_crash_safe_replication_enabled"),
		//table.BigIntColumn("settings_data_disk_size_gb"),
		//table.TextColumn("settings_data_disk_type"),
		//table.TextColumn("settings_database_flags"),
		//table.TextColumn("settings_database_flags_name"),
		//table.TextColumn("settings_database_flags_value"),
		//table.TextColumn("settings_database_replication_enabled"),
		//table.TextColumn("settings_deny_maintenance_periods"),
		//table.TextColumn("settings_deny_maintenance_periods_end_date"),
		//table.TextColumn("settings_deny_maintenance_periods_start_date"),
		//table.TextColumn("settings_deny_maintenance_periods_time"),
		//table.TextColumn("settings_insights_config"),
		//table.TextColumn("settings_insights_config_query_insights_enabled"),
		//table.BigIntColumn("settings_insights_config_query_string_length"),
		//table.TextColumn("settings_insights_config_record_application_tags"),
		//table.TextColumn("settings_insights_config_record_client_address"),
		//table.TextColumn("settings_ip_configuration"),
		//table.TextColumn("settings_ip_configuration_authorized_networks"),
		//table.TextColumn("settings_ip_configuration_authorized_networks_expiration_time"),
		//table.TextColumn("settings_ip_configuration_authorized_networks_kind"),
		//table.TextColumn("settings_ip_configuration_authorized_networks_name"),
		//table.TextColumn("settings_ip_configuration_authorized_networks_value"),
		//table.TextColumn("settings_ip_configuration_ipv4_enabled"),
		//table.TextColumn("settings_ip_configuration_private_network"),
		//table.TextColumn("settings_ip_configuration_require_ssl"),
		//table.TextColumn("settings_kind"),
		//table.TextColumn("settings_location_preference"),
		//table.TextColumn("settings_location_preference_follow_gae_application"),
		//table.TextColumn("settings_location_preference_kind"),
		//table.TextColumn("settings_location_preference_secondary_zone"),
		//table.TextColumn("settings_location_preference_zone"),
		//table.TextColumn("settings_maintenance_window"),
		//table.BigIntColumn("settings_maintenance_window_day"),
		//table.BigIntColumn("settings_maintenance_window_hour"),
		//table.TextColumn("settings_maintenance_window_kind"),
		//table.TextColumn("settings_maintenance_window_update_track"),
		//table.TextColumn("settings_pricing_plan"),
		//table.TextColumn("settings_replication_type"),
		//table.BigIntColumn("settings_settings_version"),
		//table.TextColumn("settings_storage_auto_resize"),
		//table.BigIntColumn("settings_storage_auto_resize_limit"),
		//table.TextColumn("settings_tier"),
		//table.TextColumn("settings_user_labels"),
		table.TextColumn("state"),
		table.TextColumn("suspension_reason"),
	}
}

// GcpSQLInstancesGenerate returns the rows in the table for all configured accounts
func GcpSQLInstancesGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	ctx, cancel := context.WithCancel(osqCtx)
	defer cancel()

	resultMap := make([]map[string]string, 0)

	if len(utilities.ExtConfiguration.ExtConfGcp.Accounts) == 0 && extgcp.ShouldProcessProject("gcp_sql_instance", utilities.DefaultGcpProjectID) {
		results, err := processAccountGcpSQLInstances(ctx, queryContext, nil)
		if err == nil {
			resultMap = append(resultMap, results...)
		}
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfGcp.Accounts {
			if !extgcp.ShouldProcessProject("gcp_sql_instance", account.ProjectID) {
				continue
			}
			results, err := processAccountGcpSQLInstances(ctx, queryContext, &account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}
	return resultMap, nil
}

func getGcpSQLInstancesNewServiceForAccount(ctx context.Context, account *utilities.ExtensionConfigurationGcpAccount) (*gcpsql.Service, string) {
	var projectID string
	var service *gcpsql.Service
	var err error
	if account != nil && account.KeyFile != "" {
		projectID = account.ProjectID
		service, err = gcpsql.NewService(ctx, option.WithCredentialsFile(account.KeyFile))
	} else if account != nil && account.ProjectID != "" {
		projectID = account.ProjectID
		service, err = gcpsql.NewService(ctx)
	} else {
		projectID = utilities.DefaultGcpProjectID
		service, err = gcpsql.NewService(ctx)
	}
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_sql_instance",
			"projectId": projectID,
			"errString": err.Error(),
		}).Error("failed to create service")
		return nil, ""
	}
	return service, projectID
}

func processAccountGcpSQLInstances(ctx context.Context, queryContext table.QueryContext,
	account *utilities.ExtensionConfigurationGcpAccount) ([]map[string]string, error) {

	resultMap := make([]map[string]string, 0)

	service, projectID := getGcpSQLInstancesNewServiceForAccount(ctx, account)
	if service == nil {
		return resultMap, fmt.Errorf("failed to initialize gcpsql.Service")
	}

	listCall := service.Instances.List(projectID)
	if listCall == nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_sql_instance",
			"projectId": projectID,
		}).Debug("list call is nil")
		return resultMap, nil
	}
	itemsContainer := myGcpSQLInstancesItemsContainer{Items: make([]*gcpsql.DatabaseInstance, 0)}
	if err := listCall.Pages(ctx, func(page *gcpsql.InstancesListResponse) error {

		itemsContainer.Items = append(itemsContainer.Items, page.Items...)

		return nil
	}); err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_sql_instance",
			"projectId": projectID,
			"errString": err.Error(),
		}).Error("failed to get aggregate list page")
		return resultMap, nil
	}

	byteArr, err := json.Marshal(itemsContainer)
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_sql_instance",
			"errString": err.Error(),
		}).Error("failed to marshal response")
		return resultMap, err
	}
	tableConfig, ok := utilities.TableConfigurationMap["gcp_sql_instance"]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_sql_instance",
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found for \"gcp_sql_instance\"")
	}
	jsonTable := utilities.NewTable(byteArr, tableConfig)
	for _, row := range jsonTable.Rows {
		if !extgcp.ShouldProcessRow(ctx, queryContext, "gcp_sql_instance", projectID, "", row) {
			continue
		}
		result := extgcp.RowToMap(row, projectID, "", tableConfig)
		resultMap = append(resultMap, result)
	}

	return resultMap, nil
}
