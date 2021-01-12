package compute

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/Uptycs/cloudquery/extension/azure"
	extazure "github.com/Uptycs/cloudquery/extension/azure"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2018-06-01/compute"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-02-01/resources"
	"github.com/Uptycs/cloudquery/utilities"
	"github.com/kolide/osquery-go/plugin/table"

	"github.com/pkg/errors"
)

func ComputeVmColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
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
		table.TextColumn("properties_additional_capabilities"),
		//table.TextColumn("properties_additional_capabilities_ultra_ssd_enabled"),
		table.TextColumn("properties_availability_set"),
		//table.TextColumn("properties_availability_set_id"),
		table.TextColumn("properties_billing_profile"),
		//table.DoubleColumn("properties_billing_profile_max_price"),
		table.TextColumn("properties_diagnostics_profile"),
		//table.TextColumn("properties_diagnostics_profile_boot_diagnostics"),
		//table.TextColumn("properties_diagnostics_profile_boot_diagnostics_enabled"),
		//table.TextColumn("properties_diagnostics_profile_boot_diagnostics_storage_uri"),
		table.TextColumn("properties_eviction_policy"),
		table.TextColumn("properties_extensions_time_budget"),
		table.TextColumn("properties_hardware_profile"),
		//table.TextColumn("properties_hardware_profile_vm_size"),
		table.TextColumn("properties_host"),
		table.TextColumn("properties_host_group"),
		//table.TextColumn("properties_host_group_id"),
		//table.TextColumn("properties_host_id"),
		table.TextColumn("properties_instance_view"),
		//table.TextColumn("properties_instance_view_assigned_host"),
		//table.TextColumn("properties_instance_view_boot_diagnostics"),
		//table.TextColumn("properties_instance_view_boot_diagnostics_console_screenshot_blob_uri"),
		//table.TextColumn("properties_instance_view_boot_diagnostics_serial_console_log_blob_uri"),
		//table.TextColumn("properties_instance_view_boot_diagnostics_status"),
		//table.TextColumn("properties_instance_view_boot_diagnostics_status_code"),
		//table.TextColumn("properties_instance_view_boot_diagnostics_status_display_status"),
		//table.TextColumn("properties_instance_view_boot_diagnostics_status_level"),
		//table.TextColumn("properties_instance_view_boot_diagnostics_status_message"),
		//table.TextColumn("properties_instance_view_boot_diagnostics_status_time"),
		//table.TextColumn("properties_instance_view_computer_name"),
		//table.TextColumn("properties_instance_view_disks"),
		//table.TextColumn("properties_instance_view_disks_encryption_settings"),
		//table.TextColumn("properties_instance_view_disks_encryption_settings_disk_encryption_key"),
		//table.TextColumn("properties_instance_view_disks_encryption_settings_disk_encryption_key_secret_url"),
		//table.TextColumn("properties_instance_view_disks_encryption_settings_disk_encryption_key_source_vault"),
		//table.TextColumn("properties_instance_view_disks_encryption_settings_disk_encryption_key_source_vault_id"),
		//table.TextColumn("properties_instance_view_disks_encryption_settings_enabled"),
		//table.TextColumn("properties_instance_view_disks_encryption_settings_key_encryption_key"),
		//table.TextColumn("properties_instance_view_disks_encryption_settings_key_encryption_key_key_url"),
		//table.TextColumn("properties_instance_view_disks_encryption_settings_key_encryption_key_source_vault"),
		//table.TextColumn("properties_instance_view_disks_encryption_settings_key_encryption_key_source_vault_id"),
		//table.TextColumn("properties_instance_view_disks_name"),
		//table.TextColumn("properties_instance_view_disks_statuses"),
		//table.TextColumn("properties_instance_view_disks_statuses_code"),
		//table.TextColumn("properties_instance_view_disks_statuses_display_status"),
		//table.TextColumn("properties_instance_view_disks_statuses_level"),
		//table.TextColumn("properties_instance_view_disks_statuses_message"),
		//table.TextColumn("properties_instance_view_disks_statuses_time"),
		//table.TextColumn("properties_instance_view_extensions"),
		//table.TextColumn("properties_instance_view_extensions_name"),
		//table.TextColumn("properties_instance_view_extensions_statuses"),
		//table.TextColumn("properties_instance_view_extensions_statuses_code"),
		//table.TextColumn("properties_instance_view_extensions_statuses_display_status"),
		//table.TextColumn("properties_instance_view_extensions_statuses_level"),
		//table.TextColumn("properties_instance_view_extensions_statuses_message"),
		//table.TextColumn("properties_instance_view_extensions_statuses_time"),
		//table.TextColumn("properties_instance_view_extensions_substatuses"),
		//table.TextColumn("properties_instance_view_extensions_substatuses_code"),
		//table.TextColumn("properties_instance_view_extensions_substatuses_display_status"),
		//table.TextColumn("properties_instance_view_extensions_substatuses_level"),
		//table.TextColumn("properties_instance_view_extensions_substatuses_message"),
		//table.TextColumn("properties_instance_view_extensions_substatuses_time"),
		//table.TextColumn("properties_instance_view_extensions_type"),
		//table.TextColumn("properties_instance_view_extensions_type_handler_version"),
		//table.TextColumn("properties_instance_view_hyper_v_generation"),
		//table.TextColumn("properties_instance_view_maintenance_redeploy_status"),
		//table.TextColumn("properties_instance_view_maintenance_redeploy_status_is_customer_initiated_maintenance_allowed"),
		//table.TextColumn("properties_instance_view_maintenance_redeploy_status_last_operation_message"),
		//table.TextColumn("properties_instance_view_maintenance_redeploy_status_last_operation_result_code"),
		//table.TextColumn("properties_instance_view_maintenance_redeploy_status_maintenance_window_end_time"),
		//table.TextColumn("properties_instance_view_maintenance_redeploy_status_maintenance_window_start_time"),
		//table.TextColumn("properties_instance_view_maintenance_redeploy_status_pre_maintenance_window_end_time"),
		//table.TextColumn("properties_instance_view_maintenance_redeploy_status_pre_maintenance_window_start_time"),
		//table.TextColumn("properties_instance_view_os_name"),
		//table.TextColumn("properties_instance_view_os_version"),
		//table.TextColumn("properties_instance_view_patch_status"),
		//table.TextColumn("properties_instance_view_patch_status_available_patch_summary"),
		//table.TextColumn("properties_instance_view_patch_status_available_patch_summary_assessment_activity_id"),
		//table.IntegerColumn("properties_instance_view_patch_status_available_patch_summary_critical_and_security_patch_count"),
		//table.TextColumn("properties_instance_view_patch_status_available_patch_summary_error"),
		//table.TextColumn("properties_instance_view_patch_status_available_patch_summary_error_code"),
		//table.TextColumn("properties_instance_view_patch_status_available_patch_summary_error_details"),
		//table.TextColumn("properties_instance_view_patch_status_available_patch_summary_error_details_code"),
		//table.TextColumn("properties_instance_view_patch_status_available_patch_summary_error_details_message"),
		//table.TextColumn("properties_instance_view_patch_status_available_patch_summary_error_details_target"),
		//table.TextColumn("properties_instance_view_patch_status_available_patch_summary_error_innererror"),
		//table.TextColumn("properties_instance_view_patch_status_available_patch_summary_error_innererror_errordetail"),
		//table.TextColumn("properties_instance_view_patch_status_available_patch_summary_error_innererror_exceptiontype"),
		//table.TextColumn("properties_instance_view_patch_status_available_patch_summary_error_message"),
		//table.TextColumn("properties_instance_view_patch_status_available_patch_summary_error_target"),
		//table.TextColumn("properties_instance_view_patch_status_available_patch_summary_last_modified_time"),
		//table.IntegerColumn("properties_instance_view_patch_status_available_patch_summary_other_patch_count"),
		//table.TextColumn("properties_instance_view_patch_status_available_patch_summary_reboot_pending"),
		//table.TextColumn("properties_instance_view_patch_status_available_patch_summary_start_time"),
		//table.TextColumn("properties_instance_view_patch_status_available_patch_summary_status"),
		//table.TextColumn("properties_instance_view_patch_status_last_patch_installation_summary"),
		//table.TextColumn("properties_instance_view_patch_status_last_patch_installation_summary_error"),
		//table.TextColumn("properties_instance_view_patch_status_last_patch_installation_summary_error_code"),
		//table.TextColumn("properties_instance_view_patch_status_last_patch_installation_summary_error_details"),
		//table.TextColumn("properties_instance_view_patch_status_last_patch_installation_summary_error_details_code"),
		//table.TextColumn("properties_instance_view_patch_status_last_patch_installation_summary_error_details_message"),
		//table.TextColumn("properties_instance_view_patch_status_last_patch_installation_summary_error_details_target"),
		//table.TextColumn("properties_instance_view_patch_status_last_patch_installation_summary_error_innererror"),
		//table.TextColumn("properties_instance_view_patch_status_last_patch_installation_summary_error_innererror_errordetail"),
		//table.TextColumn("properties_instance_view_patch_status_last_patch_installation_summary_error_innererror_exceptiontype"),
		//table.TextColumn("properties_instance_view_patch_status_last_patch_installation_summary_error_message"),
		//table.TextColumn("properties_instance_view_patch_status_last_patch_installation_summary_error_target"),
		//table.IntegerColumn("properties_instance_view_patch_status_last_patch_installation_summary_excluded_patch_count"),
		//table.IntegerColumn("properties_instance_view_patch_status_last_patch_installation_summary_failed_patch_count"),
		//table.TextColumn("properties_instance_view_patch_status_last_patch_installation_summary_installation_activity_id"),
		//table.IntegerColumn("properties_instance_view_patch_status_last_patch_installation_summary_installed_patch_count"),
		//table.TextColumn("properties_instance_view_patch_status_last_patch_installation_summary_last_modified_time"),
		//table.TextColumn("properties_instance_view_patch_status_last_patch_installation_summary_maintenance_window_exceeded"),
		//table.IntegerColumn("properties_instance_view_patch_status_last_patch_installation_summary_not_selected_patch_count"),
		//table.IntegerColumn("properties_instance_view_patch_status_last_patch_installation_summary_pending_patch_count"),
		//table.TextColumn("properties_instance_view_patch_status_last_patch_installation_summary_reboot_status"),
		//table.TextColumn("properties_instance_view_patch_status_last_patch_installation_summary_start_time"),
		//table.TextColumn("properties_instance_view_patch_status_last_patch_installation_summary_started_by"),
		//table.TextColumn("properties_instance_view_patch_status_last_patch_installation_summary_status"),
		//table.IntegerColumn("properties_instance_view_platform_fault_domain"),
		//table.IntegerColumn("properties_instance_view_platform_update_domain"),
		//table.TextColumn("properties_instance_view_rdp_thumb_print"),
		//table.TextColumn("properties_instance_view_statuses"),
		//table.TextColumn("properties_instance_view_statuses_code"),
		//table.TextColumn("properties_instance_view_statuses_display_status"),
		//table.TextColumn("properties_instance_view_statuses_level"),
		//table.TextColumn("properties_instance_view_statuses_message"),
		//table.TextColumn("properties_instance_view_statuses_time"),
		//table.TextColumn("properties_instance_view_vm_agent"),
		//table.TextColumn("properties_instance_view_vm_agent_extension_handlers"),
		//table.TextColumn("properties_instance_view_vm_agent_extension_handlers_status"),
		//table.TextColumn("properties_instance_view_vm_agent_extension_handlers_status_code"),
		//table.TextColumn("properties_instance_view_vm_agent_extension_handlers_status_display_status"),
		//table.TextColumn("properties_instance_view_vm_agent_extension_handlers_status_level"),
		//table.TextColumn("properties_instance_view_vm_agent_extension_handlers_status_message"),
		//table.TextColumn("properties_instance_view_vm_agent_extension_handlers_status_time"),
		//table.TextColumn("properties_instance_view_vm_agent_extension_handlers_type"),
		//table.TextColumn("properties_instance_view_vm_agent_extension_handlers_type_handler_version"),
		//table.TextColumn("properties_instance_view_vm_agent_statuses"),
		//table.TextColumn("properties_instance_view_vm_agent_statuses_code"),
		//table.TextColumn("properties_instance_view_vm_agent_statuses_display_status"),
		//table.TextColumn("properties_instance_view_vm_agent_statuses_level"),
		//table.TextColumn("properties_instance_view_vm_agent_statuses_message"),
		//table.TextColumn("properties_instance_view_vm_agent_statuses_time"),
		//table.TextColumn("properties_instance_view_vm_agent_vm_agent_version"),
		//table.TextColumn("properties_instance_view_vm_health"),
		//table.TextColumn("properties_instance_view_vm_health_status"),
		//table.TextColumn("properties_instance_view_vm_health_status_code"),
		//table.TextColumn("properties_instance_view_vm_health_status_display_status"),
		//table.TextColumn("properties_instance_view_vm_health_status_level"),
		//table.TextColumn("properties_instance_view_vm_health_status_message"),
		//table.TextColumn("properties_instance_view_vm_health_status_time"),
		table.TextColumn("properties_license_type"),
		table.TextColumn("properties_network_profile"),
		//table.TextColumn("properties_network_profile_network_interfaces"),
		//table.TextColumn("properties_network_profile_network_interfaces_id"),
		//table.TextColumn("properties_network_profile_network_interfaces_properties"),
		//table.TextColumn("properties_network_profile_network_interfaces_properties_primary"),
		table.TextColumn("properties_os_profile"),
		//table.TextColumn("properties_os_profile_admin_password"),
		//table.TextColumn("properties_os_profile_admin_username"),
		//table.TextColumn("properties_os_profile_allow_extension_operations"),
		//table.TextColumn("properties_os_profile_computer_name"),
		//table.TextColumn("properties_os_profile_custom_data"),
		//table.TextColumn("properties_os_profile_linux_configuration"),
		//table.TextColumn("properties_os_profile_linux_configuration_disable_password_authentication"),
		//table.TextColumn("properties_os_profile_linux_configuration_provision_vm_agent"),
		//table.TextColumn("properties_os_profile_linux_configuration_ssh"),
		//table.TextColumn("properties_os_profile_linux_configuration_ssh_public_keys"),
		//table.TextColumn("properties_os_profile_linux_configuration_ssh_public_keys_key_data"),
		//table.TextColumn("properties_os_profile_linux_configuration_ssh_public_keys_path"),
		//table.TextColumn("properties_os_profile_require_guest_provision_signal"),
		//table.TextColumn("properties_os_profile_secrets"),
		//table.TextColumn("properties_os_profile_secrets_source_vault"),
		//table.TextColumn("properties_os_profile_secrets_source_vault_id"),
		//table.TextColumn("properties_os_profile_secrets_vault_certificates"),
		//table.TextColumn("properties_os_profile_secrets_vault_certificates_certificate_store"),
		//table.TextColumn("properties_os_profile_secrets_vault_certificates_certificate_url"),
		//table.TextColumn("properties_os_profile_windows_configuration"),
		//table.TextColumn("properties_os_profile_windows_configuration_additional_unattend_content"),
		//table.TextColumn("properties_os_profile_windows_configuration_additional_unattend_content_component_name"),
		//table.TextColumn("properties_os_profile_windows_configuration_additional_unattend_content_content"),
		//table.TextColumn("properties_os_profile_windows_configuration_additional_unattend_content_pass_name"),
		//table.TextColumn("properties_os_profile_windows_configuration_additional_unattend_content_setting_name"),
		//table.TextColumn("properties_os_profile_windows_configuration_enable_automatic_updates"),
		//table.TextColumn("properties_os_profile_windows_configuration_patch_settings"),
		//table.TextColumn("properties_os_profile_windows_configuration_patch_settings_patch_mode"),
		//table.TextColumn("properties_os_profile_windows_configuration_provision_vm_agent"),
		//table.TextColumn("properties_os_profile_windows_configuration_time_zone"),
		//table.TextColumn("properties_os_profile_windows_configuration_win_rm"),
		//table.TextColumn("properties_os_profile_windows_configuration_win_rm_listeners"),
		//table.TextColumn("properties_os_profile_windows_configuration_win_rm_listeners_certificate_url"),
		//table.TextColumn("properties_os_profile_windows_configuration_win_rm_listeners_protocol"),
		table.TextColumn("properties_priority"),
		table.TextColumn("properties_provisioning_state"),
		table.TextColumn("properties_proximity_placement_group"),
		//table.TextColumn("properties_proximity_placement_group_id"),
		table.TextColumn("properties_security_profile"),
		//table.TextColumn("properties_security_profile_encryption_at_host"),
		table.TextColumn("properties_storage_profile"),
		//table.TextColumn("properties_storage_profile_data_disks"),
		//table.TextColumn("properties_storage_profile_data_disks_caching"),
		//table.TextColumn("properties_storage_profile_data_disks_create_option"),
		//table.IntegerColumn("properties_storage_profile_data_disks_disk_iops_read_write"),
		//table.IntegerColumn("properties_storage_profile_data_disks_disk_m_bps_read_write"),
		//table.IntegerColumn("properties_storage_profile_data_disks_disk_size_gb"),
		//table.TextColumn("properties_storage_profile_data_disks_image"),
		//table.TextColumn("properties_storage_profile_data_disks_image_uri"),
		//table.IntegerColumn("properties_storage_profile_data_disks_lun"),
		//table.TextColumn("properties_storage_profile_data_disks_managed_disk"),
		//table.TextColumn("properties_storage_profile_data_disks_managed_disk_disk_encryption_set"),
		//table.TextColumn("properties_storage_profile_data_disks_managed_disk_disk_encryption_set_id"),
		//table.TextColumn("properties_storage_profile_data_disks_managed_disk_id"),
		//table.TextColumn("properties_storage_profile_data_disks_managed_disk_storage_account_type"),
		//table.TextColumn("properties_storage_profile_data_disks_name"),
		//table.TextColumn("properties_storage_profile_data_disks_to_be_detached"),
		//table.TextColumn("properties_storage_profile_data_disks_vhd"),
		//table.TextColumn("properties_storage_profile_data_disks_vhd_uri"),
		//table.TextColumn("properties_storage_profile_data_disks_write_accelerator_enabled"),
		//table.TextColumn("properties_storage_profile_image_reference"),
		//table.TextColumn("properties_storage_profile_image_reference_exact_version"),
		//table.TextColumn("properties_storage_profile_image_reference_id"),
		//table.TextColumn("properties_storage_profile_image_reference_offer"),
		//table.TextColumn("properties_storage_profile_image_reference_publisher"),
		//table.TextColumn("properties_storage_profile_image_reference_sku"),
		//table.TextColumn("properties_storage_profile_image_reference_version"),
		//table.TextColumn("properties_storage_profile_os_disk"),
		//table.TextColumn("properties_storage_profile_os_disk_caching"),
		//table.TextColumn("properties_storage_profile_os_disk_create_option"),
		//table.TextColumn("properties_storage_profile_os_disk_diff_disk_settings"),
		//table.TextColumn("properties_storage_profile_os_disk_diff_disk_settings_option"),
		//table.TextColumn("properties_storage_profile_os_disk_diff_disk_settings_placement"),
		//table.IntegerColumn("properties_storage_profile_os_disk_disk_size_gb"),
		//table.TextColumn("properties_storage_profile_os_disk_encryption_settings"),
		//table.TextColumn("properties_storage_profile_os_disk_encryption_settings_disk_encryption_key"),
		//table.TextColumn("properties_storage_profile_os_disk_encryption_settings_disk_encryption_key_secret_url"),
		//table.TextColumn("properties_storage_profile_os_disk_encryption_settings_disk_encryption_key_source_vault"),
		//table.TextColumn("properties_storage_profile_os_disk_encryption_settings_disk_encryption_key_source_vault_id"),
		//table.TextColumn("properties_storage_profile_os_disk_encryption_settings_enabled"),
		//table.TextColumn("properties_storage_profile_os_disk_encryption_settings_key_encryption_key"),
		//table.TextColumn("properties_storage_profile_os_disk_encryption_settings_key_encryption_key_key_url"),
		//table.TextColumn("properties_storage_profile_os_disk_encryption_settings_key_encryption_key_source_vault"),
		//table.TextColumn("properties_storage_profile_os_disk_encryption_settings_key_encryption_key_source_vault_id"),
		//table.TextColumn("properties_storage_profile_os_disk_image"),
		//table.TextColumn("properties_storage_profile_os_disk_image_uri"),
		//table.TextColumn("properties_storage_profile_os_disk_managed_disk"),
		//table.TextColumn("properties_storage_profile_os_disk_managed_disk_disk_encryption_set"),
		//table.TextColumn("properties_storage_profile_os_disk_managed_disk_disk_encryption_set_id"),
		//table.TextColumn("properties_storage_profile_os_disk_managed_disk_id"),
		//table.TextColumn("properties_storage_profile_os_disk_managed_disk_storage_account_type"),
		//table.TextColumn("properties_storage_profile_os_disk_name"),
		//table.TextColumn("properties_storage_profile_os_disk_os_type"),
		//table.TextColumn("properties_storage_profile_os_disk_vhd"),
		//table.TextColumn("properties_storage_profile_os_disk_vhd_uri"),
		//table.TextColumn("properties_storage_profile_os_disk_write_accelerator_enabled"),
		table.TextColumn("properties_virtual_machine_scale_set"),
		//table.TextColumn("properties_virtual_machine_scale_set_id"),
		table.TextColumn("properties_vm_id"),
		table.TextColumn("resources"),
		//table.TextColumn("resources_id"),
		//table.TextColumn("resources_location"),
		//table.TextColumn("resources_name"),
		//table.TextColumn("resources_properties"),
		//table.TextColumn("resources_properties_auto_upgrade_minor_version"),
		//table.TextColumn("resources_properties_enable_automatic_upgrade"),
		//table.TextColumn("resources_properties_force_update_tag"),
		//table.TextColumn("resources_properties_instance_view"),
		//table.TextColumn("resources_properties_instance_view_name"),
		//table.TextColumn("resources_properties_instance_view_statuses"),
		//table.TextColumn("resources_properties_instance_view_statuses_code"),
		//table.TextColumn("resources_properties_instance_view_statuses_display_status"),
		//table.TextColumn("resources_properties_instance_view_statuses_level"),
		//table.TextColumn("resources_properties_instance_view_statuses_message"),
		//table.TextColumn("resources_properties_instance_view_statuses_time"),
		//table.TextColumn("resources_properties_instance_view_substatuses"),
		//table.TextColumn("resources_properties_instance_view_substatuses_code"),
		//table.TextColumn("resources_properties_instance_view_substatuses_display_status"),
		//table.TextColumn("resources_properties_instance_view_substatuses_level"),
		//table.TextColumn("resources_properties_instance_view_substatuses_message"),
		//table.TextColumn("resources_properties_instance_view_substatuses_time"),
		//table.TextColumn("resources_properties_instance_view_type"),
		//table.TextColumn("resources_properties_instance_view_type_handler_version"),
		//table.TextColumn("resources_properties_protected_settings"),
		//table.TextColumn("resources_properties_provisioning_state"),
		//table.TextColumn("resources_properties_publisher"),
		//table.TextColumn("resources_properties_settings"),
		//table.TextColumn("resources_properties_type"),
		//table.TextColumn("resources_properties_type_handler_version"),
		//table.TextColumn("resources_tags"),
		//table.TextColumn("resources_type"),
		table.TextColumn("tags"),
		table.TextColumn("type"),
		table.TextColumn("zones"),
	}
}

func ComputeVmGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	if len(utilities.ExtConfiguration.ExtConfAzure.Accounts) == 0 {
		fmt.Println("Processing default account")
		results, err := processAccountComputeVms(nil)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, results...)
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfAzure.Accounts {
			fmt.Println("Processing account:" + account.SubscriptionId)
			results, err := processAccountComputeVms(&account)
			if err != nil {
				// TODO: Continue to next account or return error ?
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}

	return resultMap, nil
}

func processAccountComputeVms(account *utilities.ExtensionConfigurationAzureAccount) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	var wg sync.WaitGroup
	session, err := azure.GetAuthSession(account)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	groups, err := getGroups(session)

	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	wg.Add(len(groups))

	tableConfig, ok := utilities.TableConfigurationMap["azure_compute_vm"]
	if !ok {
		fmt.Println("getTableConfig: ", err)
		log.Fatal(err)
		return resultMap, fmt.Errorf("table configuration not found")
	}

	for _, group := range groups {
		go getVM(session, group, &wg, &resultMap, tableConfig)
	}
	wg.Wait()
	return resultMap, nil
}

func getGroups(session *azure.AzureSession) ([]string, error) {
	tab := make([]string, 0)
	var err error

	grClient := resources.NewGroupsClient(session.SubscriptionId)
	grClient.Authorizer = session.Authorizer

	for list, err := grClient.ListComplete(context.Background(), "", nil); list.NotDone(); err = list.Next() {
		if err != nil {
			return nil, errors.Wrap(err, "error traverising RG list")
		}
		rgName := *list.Value().Name
		tab = append(tab, rgName)
	}
	return tab, err
}

func getVM(session *azure.AzureSession, rg string, wg *sync.WaitGroup, resultMap *[]map[string]string, tableConfig *utilities.TableConfig) {
	defer wg.Done()

	vmClient := compute.NewVirtualMachinesClient(session.SubscriptionId)
	vmClient.Authorizer = session.Authorizer

	for vm, err := vmClient.ListComplete(context.Background(), rg); vm.NotDone(); err = vm.Next() {
		if err != nil {
			log.Print("got error while traverising RG list: ", err)
		}

		i := vm.Value()
		byteArr, err := json.Marshal(i)
		if err != nil {
			fmt.Println("getImages marshal: ", err)
			log.Fatal(err)
			continue
		}
		fmt.Println("Data:" + string(byteArr))
		table := utilities.Table{}
		table.Init(byteArr, tableConfig.MaxLevel, tableConfig.GetParsedAttributeConfigMap())
		for _, row := range table.Rows {
			result := extazure.RowToMap(row, session.SubscriptionId, "", rg, tableConfig)
			*resultMap = append(*resultMap, result)
		}
	}
}
