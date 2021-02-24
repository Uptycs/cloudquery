/**
 * Copyright (c) 2020-present, The cloudquery authors
 *
 * This source code is licensed as defined by the LICENSE file found in the
 * root directory of this source tree.
 *
 * SPDX-License-Identifier: (Apache-2.0 OR GPL-2.0-only)
 */

package dns

import (
	"context"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"

	"github.com/Uptycs/basequery-go/plugin/table"
	extgcp "github.com/Uptycs/cloudquery/extension/gcp"
	"github.com/Uptycs/cloudquery/utilities"

	"google.golang.org/api/option"

	gcpdns "google.golang.org/api/dns/v1beta2"
)

type myGcpDNSManagedZonesItemsContainer struct {
	Items []*gcpdns.ManagedZone `json:"items"`
}

// GcpDNSManagedZonesColumns returns the list of columns for gcp_dns_managed_zone
func GcpDNSManagedZonesColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("project_id"),
		table.TextColumn("creation_time"),
		table.TextColumn("description"),
		table.TextColumn("dns_name"),
		table.TextColumn("dnssec_config"),
		//table.TextColumn("dnssec_config_default_key_specs"),
		//table.TextColumn("dnssec_config_default_key_specs_algorithm"),
		//table.BigIntColumn("dnssec_config_default_key_specs_key_length"),
		//table.TextColumn("dnssec_config_default_key_specs_key_type"),
		//table.TextColumn("dnssec_config_default_key_specs_kind"),
		//table.TextColumn("dnssec_config_kind"),
		//table.TextColumn("dnssec_config_non_existence"),
		//table.TextColumn("dnssec_config_state"),
		table.TextColumn("forwarding_config"),
		//table.TextColumn("forwarding_config_kind"),
		//table.TextColumn("forwarding_config_target_name_servers"),
		//table.TextColumn("forwarding_config_target_name_servers_forwarding_path"),
		//table.TextColumn("forwarding_config_target_name_servers_ipv4_address"),
		//table.TextColumn("forwarding_config_target_name_servers_kind"),
		table.BigIntColumn("id"),
		table.TextColumn("kind"),
		table.TextColumn("labels"),
		table.TextColumn("name"),
		table.TextColumn("name_server_set"),
		table.TextColumn("name_servers"),
		table.TextColumn("peering_config"),
		//table.TextColumn("peering_config_kind"),
		//table.TextColumn("peering_config_target_network"),
		//table.TextColumn("peering_config_target_network_deactivate_time"),
		//table.TextColumn("peering_config_target_network_kind"),
		//table.TextColumn("peering_config_target_network_network_url"),
		table.TextColumn("private_visibility_config"),
		//table.TextColumn("private_visibility_config_kind"),
		//table.TextColumn("private_visibility_config_networks"),
		//table.TextColumn("private_visibility_config_networks_kind"),
		//table.TextColumn("private_visibility_config_networks_network_url"),
		table.TextColumn("reverse_lookup_config"),
		//table.TextColumn("reverse_lookup_config_kind"),
		table.TextColumn("service_directory_config"),
		//table.TextColumn("service_directory_config_kind"),
		//table.TextColumn("service_directory_config_namespace"),
		//table.TextColumn("service_directory_config_namespace_deletion_time"),
		//table.TextColumn("service_directory_config_namespace_kind"),
		//table.TextColumn("service_directory_config_namespace_namespace_url"),
		table.TextColumn("visibility"),
	}
}

// GcpDNSManagedZonesGenerate returns the rows in the table for all configured accounts
func GcpDNSManagedZonesGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	ctx, cancel := context.WithCancel(osqCtx)
	defer cancel()

	resultMap := make([]map[string]string, 0)

	if len(utilities.ExtConfiguration.ExtConfGcp.Accounts) == 0 && extgcp.ShouldProcessProject("gcp_dns_managed_zone", utilities.DefaultGcpProjectID) {
		results, err := processAccountGcpDNSManagedZones(ctx, queryContext, nil)
		if err == nil {
			resultMap = append(resultMap, results...)
		}
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfGcp.Accounts {
			if !extgcp.ShouldProcessProject("gcp_dns_managed_zone", account.ProjectID) {
				continue
			}
			results, err := processAccountGcpDNSManagedZones(ctx, queryContext, &account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}
	return resultMap, nil
}

func getGcpDNSManagedZonesNewServiceForAccount(ctx context.Context, account *utilities.ExtensionConfigurationGcpAccount) (*gcpdns.Service, string) {
	var projectID string
	var service *gcpdns.Service
	var err error
	if account != nil && account.KeyFile != "" {
		projectID = account.ProjectID
		service, err = gcpdns.NewService(ctx, option.WithCredentialsFile(account.KeyFile))
	} else if account != nil && account.ProjectID != "" {
		projectID = account.ProjectID
		service, err = gcpdns.NewService(ctx)
	} else {
		projectID = utilities.DefaultGcpProjectID
		service, err = gcpdns.NewService(ctx)
	}
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_dns_managed_zone",
			"projectId": projectID,
			"errString": err.Error(),
		}).Error("failed to create service")
		return nil, ""
	}
	return service, projectID
}

func processAccountGcpDNSManagedZones(ctx context.Context, queryContext table.QueryContext,
	account *utilities.ExtensionConfigurationGcpAccount) ([]map[string]string, error) {

	resultMap := make([]map[string]string, 0)

	service, projectID := getGcpDNSManagedZonesNewServiceForAccount(ctx, account)
	if service == nil {
		return resultMap, fmt.Errorf("failed to initialize gcpdns.Service")
	}

	listCall := service.ManagedZones.List(projectID)
	if listCall == nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_dns_managed_zone",
			"projectId": projectID,
		}).Debug("list call is nil")
		return resultMap, nil
	}
	itemsContainer := myGcpDNSManagedZonesItemsContainer{Items: make([]*gcpdns.ManagedZone, 0)}
	if err := listCall.Pages(ctx, func(page *gcpdns.ManagedZonesListResponse) error {

		itemsContainer.Items = append(itemsContainer.Items, page.ManagedZones...)

		return nil
	}); err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_dns_managed_zone",
			"projectId": projectID,
			"errString": err.Error(),
		}).Error("failed to get aggregate list page")
		return resultMap, nil
	}

	byteArr, err := json.Marshal(itemsContainer)
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_dns_managed_zone",
			"errString": err.Error(),
		}).Error("failed to marshal response")
		return resultMap, err
	}
	tableConfig, ok := utilities.TableConfigurationMap["gcp_dns_managed_zone"]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_dns_managed_zone",
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found for \"gcp_dns_managed_zone\"")
	}
	jsonTable := utilities.NewTable(byteArr, tableConfig)
	for _, row := range jsonTable.Rows {
		if !extgcp.ShouldProcessRow(ctx, queryContext, "gcp_dns_managed_zone", projectID, "", row) {
			continue
		}
		result := extgcp.RowToMap(row, projectID, "", tableConfig)
		resultMap = append(resultMap, result)
	}

	return resultMap, nil
}
