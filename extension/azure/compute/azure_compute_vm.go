package compute

import (
	"context"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"sync"

	"github.com/Uptycs/cloudquery/extension/azure"
	extazure "github.com/Uptycs/cloudquery/extension/azure"

	"github.com/Uptycs/cloudquery/utilities"
	"github.com/kolide/osquery-go/plugin/table"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2018-06-01/compute"
)

// VirtualMachinesColumns returns the list of columns in the table
func VirtualMachinesColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("subscription_id"),
		table.TextColumn("id"),
		table.TextColumn("identity"),
		//table.TextColumn("identity_principal_id"),
		//table.TextColumn("identity_tenant_id"),
		//table.TextColumn("identity_type"),
		//table.TextColumn("identity_user_assigned_identities"),
		table.TextColumn("location"),
		table.TextColumn("name"),
		table.TextColumn("plan"),
		//table.TextColumn("plan_name"),
		//table.TextColumn("plan_product"),
		//table.TextColumn("plan_promotion_code"),
		//table.TextColumn("plan_publisher"),
		//table.TextColumn("properties"),
		table.TextColumn("additional_capabilities"),
		//table.TextColumn("additional_capabilities_ultra_ssd_enabled"),
		table.TextColumn("availability_set"),
		//table.TextColumn("availability_set_id"),
		table.TextColumn("diagnostics_profile"),
		//table.TextColumn("diagnostics_profile_boot_diagnostics"),
		//table.TextColumn("diagnostics_profile_boot_diagnostics_enabled"),
		//table.TextColumn("diagnostics_profile_boot_diagnostics_storage_uri"),
		table.TextColumn("hardware_profile"),
		//table.TextColumn("hardware_profile_vm_size"),
		table.TextColumn("instance_view"),
		//table.TextColumn("instance_view_boot_diagnostics"),
		//table.TextColumn("instance_view_boot_diagnostics_console_screenshot_blob_uri"),
		//table.TextColumn("instance_view_boot_diagnostics_serial_console_log_blob_uri"),
		//table.TextColumn("instance_view_boot_diagnostics_status"),
		//table.TextColumn("instance_view_boot_diagnostics_status_code"),
		//table.TextColumn("instance_view_boot_diagnostics_status_display_status"),
		//table.TextColumn("instance_view_boot_diagnostics_status_level"),
		//table.TextColumn("instance_view_boot_diagnostics_status_message"),
		//table.TextColumn("instance_view_boot_diagnostics_status_time"),
		//table.TextColumn("instance_view_computer_name"),
		//table.TextColumn("instance_view_disks"),
		//table.TextColumn("instance_view_disks_encryption_settings"),
		//table.TextColumn("instance_view_disks_encryption_settings_disk_encryption_key"),
		//table.TextColumn("instance_view_disks_encryption_settings_disk_encryption_key_secret_url"),
		//table.TextColumn("instance_view_disks_encryption_settings_disk_encryption_key_source_vault"),
		//table.TextColumn("instance_view_disks_encryption_settings_disk_encryption_key_source_vault_id"),
		//table.TextColumn("instance_view_disks_encryption_settings_enabled"),
		//table.TextColumn("instance_view_disks_encryption_settings_key_encryption_key"),
		//table.TextColumn("instance_view_disks_encryption_settings_key_encryption_key_key_url"),
		//table.TextColumn("instance_view_disks_encryption_settings_key_encryption_key_source_vault"),
		//table.TextColumn("instance_view_disks_encryption_settings_key_encryption_key_source_vault_id"),
		//table.TextColumn("instance_view_disks_name"),
		//table.TextColumn("instance_view_disks_statuses"),
		//table.TextColumn("instance_view_disks_statuses_code"),
		//table.TextColumn("instance_view_disks_statuses_display_status"),
		//table.TextColumn("instance_view_disks_statuses_level"),
		//table.TextColumn("instance_view_disks_statuses_message"),
		//table.TextColumn("instance_view_disks_statuses_time"),
		//table.TextColumn("instance_view_extensions"),
		//table.TextColumn("instance_view_extensions_name"),
		//table.TextColumn("instance_view_extensions_statuses"),
		//table.TextColumn("instance_view_extensions_statuses_code"),
		//table.TextColumn("instance_view_extensions_statuses_display_status"),
		//table.TextColumn("instance_view_extensions_statuses_level"),
		//table.TextColumn("instance_view_extensions_statuses_message"),
		//table.TextColumn("instance_view_extensions_statuses_time"),
		//table.TextColumn("instance_view_extensions_substatuses"),
		//table.TextColumn("instance_view_extensions_substatuses_code"),
		//table.TextColumn("instance_view_extensions_substatuses_display_status"),
		//table.TextColumn("instance_view_extensions_substatuses_level"),
		//table.TextColumn("instance_view_extensions_substatuses_message"),
		//table.TextColumn("instance_view_extensions_substatuses_time"),
		//table.TextColumn("instance_view_extensions_type"),
		//table.TextColumn("instance_view_extensions_type_handler_version"),
		//table.TextColumn("instance_view_maintenance_redeploy_status"),
		//table.TextColumn("instance_view_maintenance_redeploy_status_is_customer_initiated_maintenance_allowed"),
		//table.TextColumn("instance_view_maintenance_redeploy_status_last_operation_message"),
		//table.TextColumn("instance_view_maintenance_redeploy_status_last_operation_result_code"),
		//table.TextColumn("instance_view_maintenance_redeploy_status_maintenance_window_end_time"),
		//table.TextColumn("instance_view_maintenance_redeploy_status_maintenance_window_start_time"),
		//table.TextColumn("instance_view_maintenance_redeploy_status_pre_maintenance_window_end_time"),
		//table.TextColumn("instance_view_maintenance_redeploy_status_pre_maintenance_window_start_time"),
		//table.TextColumn("instance_view_os_name"),
		//table.TextColumn("instance_view_os_version"),
		//table.IntegerColumn("instance_view_platform_fault_domain"),
		//table.IntegerColumn("instance_view_platform_update_domain"),
		//table.TextColumn("instance_view_rdp_thumb_print"),
		//table.TextColumn("instance_view_statuses"),
		//table.TextColumn("instance_view_statuses_code"),
		//table.TextColumn("instance_view_statuses_display_status"),
		//table.TextColumn("instance_view_statuses_level"),
		//table.TextColumn("instance_view_statuses_message"),
		//table.TextColumn("instance_view_statuses_time"),
		//table.TextColumn("instance_view_vm_agent"),
		//table.TextColumn("instance_view_vm_agent_extension_handlers"),
		//table.TextColumn("instance_view_vm_agent_extension_handlers_status"),
		//table.TextColumn("instance_view_vm_agent_extension_handlers_status_code"),
		//table.TextColumn("instance_view_vm_agent_extension_handlers_status_display_status"),
		//table.TextColumn("instance_view_vm_agent_extension_handlers_status_level"),
		//table.TextColumn("instance_view_vm_agent_extension_handlers_status_message"),
		//table.TextColumn("instance_view_vm_agent_extension_handlers_status_time"),
		//table.TextColumn("instance_view_vm_agent_extension_handlers_type"),
		//table.TextColumn("instance_view_vm_agent_extension_handlers_type_handler_version"),
		//table.TextColumn("instance_view_vm_agent_statuses"),
		//table.TextColumn("instance_view_vm_agent_statuses_code"),
		//table.TextColumn("instance_view_vm_agent_statuses_display_status"),
		//table.TextColumn("instance_view_vm_agent_statuses_level"),
		//table.TextColumn("instance_view_vm_agent_statuses_message"),
		//table.TextColumn("instance_view_vm_agent_statuses_time"),
		//table.TextColumn("instance_view_vm_agent_vm_agent_version"),
		table.TextColumn("license_type"),
		table.TextColumn("network_profile"),
		//table.TextColumn("network_profile_network_interfaces"),
		//table.TextColumn("network_profile_network_interfaces_id"),
		table.TextColumn("os_profile"),
		//table.TextColumn("os_profile_admin_password"),
		//table.TextColumn("os_profile_admin_username"),
		//table.TextColumn("os_profile_allow_extension_operations"),
		//table.TextColumn("os_profile_computer_name"),
		//table.TextColumn("os_profile_custom_data"),
		//table.TextColumn("os_profile_linux_configuration"),
		//table.TextColumn("os_profile_linux_configuration_disable_password_authentication"),
		//table.TextColumn("os_profile_linux_configuration_provision_vm_agent"),
		//table.TextColumn("os_profile_linux_configuration_ssh"),
		//table.TextColumn("os_profile_linux_configuration_ssh_public_keys"),
		//table.TextColumn("os_profile_linux_configuration_ssh_public_keys_key_data"),
		//table.TextColumn("os_profile_linux_configuration_ssh_public_keys_path"),
		//table.TextColumn("os_profile_secrets"),
		//table.TextColumn("os_profile_secrets_source_vault"),
		//table.TextColumn("os_profile_secrets_source_vault_id"),
		//table.TextColumn("os_profile_secrets_vault_certificates"),
		//table.TextColumn("os_profile_secrets_vault_certificates_certificate_store"),
		//table.TextColumn("os_profile_secrets_vault_certificates_certificate_url"),
		//table.TextColumn("os_profile_windows_configuration"),
		//table.TextColumn("os_profile_windows_configuration_additional_unattend_content"),
		//table.TextColumn("os_profile_windows_configuration_additional_unattend_content_component_name"),
		//table.TextColumn("os_profile_windows_configuration_additional_unattend_content_content"),
		//table.TextColumn("os_profile_windows_configuration_additional_unattend_content_pass_name"),
		//table.TextColumn("os_profile_windows_configuration_additional_unattend_content_setting_name"),
		//table.TextColumn("os_profile_windows_configuration_enable_automatic_updates"),
		//table.TextColumn("os_profile_windows_configuration_provision_vm_agent"),
		//table.TextColumn("os_profile_windows_configuration_time_zone"),
		//table.TextColumn("os_profile_windows_configuration_win_rm"),
		//table.TextColumn("os_profile_windows_configuration_win_rm_listeners"),
		//table.TextColumn("os_profile_windows_configuration_win_rm_listeners_certificate_url"),
		//table.TextColumn("os_profile_windows_configuration_win_rm_listeners_protocol"),
		table.TextColumn("provisioning_state"),
		table.TextColumn("proximity_placement_group"),
		//table.TextColumn("proximity_placement_group_id"),
		table.TextColumn("storage_profile"),
		//table.TextColumn("storage_profile_data_disks"),
		//table.TextColumn("storage_profile_data_disks_caching"),
		//table.TextColumn("storage_profile_data_disks_create_option"),
		//table.IntegerColumn("storage_profile_data_disks_disk_size_gb"),
		//table.TextColumn("storage_profile_data_disks_image"),
		//table.TextColumn("storage_profile_data_disks_image_uri"),
		//table.IntegerColumn("storage_profile_data_disks_lun"),
		//table.TextColumn("storage_profile_data_disks_managed_disk"),
		//table.TextColumn("storage_profile_data_disks_managed_disk_id"),
		//table.TextColumn("storage_profile_data_disks_managed_disk_storage_account_type"),
		//table.TextColumn("storage_profile_data_disks_name"),
		//table.TextColumn("storage_profile_data_disks_vhd"),
		//table.TextColumn("storage_profile_data_disks_vhd_uri"),
		//table.TextColumn("storage_profile_data_disks_write_accelerator_enabled"),
		//table.TextColumn("storage_profile_image_reference"),
		//table.TextColumn("storage_profile_image_reference_id"),
		//table.TextColumn("storage_profile_image_reference_offer"),
		//table.TextColumn("storage_profile_image_reference_publisher"),
		//table.TextColumn("storage_profile_image_reference_sku"),
		//table.TextColumn("storage_profile_image_reference_version"),
		//table.TextColumn("storage_profile_os_disk"),
		//table.TextColumn("storage_profile_os_disk_caching"),
		//table.TextColumn("storage_profile_os_disk_create_option"),
		//table.TextColumn("storage_profile_os_disk_diff_disk_settings"),
		//table.TextColumn("storage_profile_os_disk_diff_disk_settings_option"),
		//table.IntegerColumn("storage_profile_os_disk_disk_size_gb"),
		//table.TextColumn("storage_profile_os_disk_encryption_settings"),
		//table.TextColumn("storage_profile_os_disk_encryption_settings_disk_encryption_key"),
		//table.TextColumn("storage_profile_os_disk_encryption_settings_disk_encryption_key_secret_url"),
		//table.TextColumn("storage_profile_os_disk_encryption_settings_disk_encryption_key_source_vault"),
		//table.TextColumn("storage_profile_os_disk_encryption_settings_disk_encryption_key_source_vault_id"),
		//table.TextColumn("storage_profile_os_disk_encryption_settings_enabled"),
		//table.TextColumn("storage_profile_os_disk_encryption_settings_key_encryption_key"),
		//table.TextColumn("storage_profile_os_disk_encryption_settings_key_encryption_key_key_url"),
		//table.TextColumn("storage_profile_os_disk_encryption_settings_key_encryption_key_source_vault"),
		//table.TextColumn("storage_profile_os_disk_encryption_settings_key_encryption_key_source_vault_id"),
		//table.TextColumn("storage_profile_os_disk_image"),
		//table.TextColumn("storage_profile_os_disk_image_uri"),
		//table.TextColumn("storage_profile_os_disk_managed_disk"),
		//table.TextColumn("storage_profile_os_disk_managed_disk_id"),
		//table.TextColumn("storage_profile_os_disk_managed_disk_storage_account_type"),
		//table.TextColumn("storage_profile_os_disk_name"),
		//table.TextColumn("storage_profile_os_disk_os_type"),
		//table.TextColumn("storage_profile_os_disk_vhd"),
		//table.TextColumn("storage_profile_os_disk_vhd_uri"),
		//table.TextColumn("storage_profile_os_disk_write_accelerator_enabled"),
		table.TextColumn("vm_id"),
		table.TextColumn("resources"),
		//table.TextColumn("resources_id"),
		//table.TextColumn("resources_location"),
		//table.TextColumn("resources_name"),
		//table.TextColumn("resources_tags"),
		//table.TextColumn("resources_type"),
		table.TextColumn("tags"),
		table.TextColumn("type"),
		table.TextColumn("zones"),
	}
}

// VirtualMachinesGenerate returns the rows in the table for all configured accounts
func VirtualMachinesGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	if len(utilities.ExtConfiguration.ExtConfAzure.Accounts) == 0 {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "azure_compute_vm",
			"account":   "default",
		}).Info("processing account")
		results, err := processAccountVirtualMachines(nil)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, results...)
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfAzure.Accounts {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "azure_compute_vm",
				"account":   account.SubscriptionID,
			}).Info("processing account")
			results, err := processAccountVirtualMachines(&account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}

	return resultMap, nil
}

func processAccountVirtualMachines(account *utilities.ExtensionConfigurationAzureAccount) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	var wg sync.WaitGroup
	session, err := azure.GetAuthSession(account)
	if err != nil {
		return resultMap, err
	}
	groups, err := azure.GetGroups(session)

	if err != nil {
		return resultMap, err
	}

	wg.Add(len(groups))

	tableConfig, ok := utilities.TableConfigurationMap["azure_compute_vm"]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "azure_compute_vm",
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found")
	}

	for _, group := range groups {
		go getVirtualMachines(session, group, &wg, &resultMap, tableConfig)
	}
	wg.Wait()
	return resultMap, nil
}

func getVirtualMachines(session *azure.AzureSession, rg string, wg *sync.WaitGroup, resultMap *[]map[string]string, tableConfig *utilities.TableConfig) {
	defer wg.Done()

	svcClient := compute.NewVirtualMachinesClient(session.SubscriptionId)
	svcClient.Authorizer = session.Authorizer

	for resourceItr, err := svcClient.ListComplete(context.Background(), rg); resourceItr.NotDone(); err = resourceItr.Next() {
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName":     "azure_compute_vm",
				"resourceGroup": rg,
				"errString":     err.Error(),
			}).Error("failed to get resource list")
			continue
		}

		resource := resourceItr.Value()
		byteArr, err := json.Marshal(resource)
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName":     "azure_compute_vm",
				"resourceGroup": rg,
				"errString":     err.Error(),
			}).Error("failed to marshal response")
			continue
		}
		table := utilities.NewTable(byteArr, tableConfig)
		for _, row := range table.Rows {
			result := extazure.RowToMap(row, session.SubscriptionId, "", rg, tableConfig)
			*resultMap = append(*resultMap, result)
		}
	}
}
