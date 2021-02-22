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

type myGcpComputeVpnTunnelsItemsContainer struct {
	Items []*compute.VpnTunnel `json:"items"`
}

// GcpComputeVpnTunnelsColumns returns the list of columns for gcp_compute_vpn_tunnel
func (handler *GcpComputeHandler) GcpComputeVpnTunnelsColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("project_id"),
		table.TextColumn("creation_timestamp"),
		table.TextColumn("description"),
		table.TextColumn("detailed_status"),
		table.BigIntColumn("id"),
		table.BigIntColumn("ike_version"),
		table.TextColumn("kind"),
		table.TextColumn("local_traffic_selector"),
		table.TextColumn("name"),
		table.TextColumn("peer_external_gateway"),
		table.BigIntColumn("peer_external_gateway_interface"),
		table.TextColumn("peer_gcp_gateway"),
		table.TextColumn("peer_ip"),
		table.TextColumn("region"),
		table.TextColumn("remote_traffic_selector"),
		table.TextColumn("router"),
		//table.TextColumn("self_link"),
		table.TextColumn("shared_secret"),
		table.TextColumn("shared_secret_hash"),
		table.TextColumn("status"),
		table.TextColumn("target_vpn_gateway"),
		table.TextColumn("vpn_gateway"),
		table.BigIntColumn("vpn_gateway_interface"),
	}
}

// GcpComputeVpnTunnelsGenerate returns the rows in the table for all configured accounts
func (handler *GcpComputeHandler) GcpComputeVpnTunnelsGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	var _ = queryContext
	ctx, cancel := context.WithCancel(osqCtx)
	defer cancel()

	resultMap := make([]map[string]string, 0)

	if len(utilities.ExtConfiguration.ExtConfGcp.Accounts) == 0 {
		results, err := handler.processAccountGcpComputeVpnTunnels(ctx, nil)
		if err == nil {
			resultMap = append(resultMap, results...)
		}
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfGcp.Accounts {
			results, err := handler.processAccountGcpComputeVpnTunnels(ctx, &account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}
	return resultMap, nil
}

func (handler *GcpComputeHandler) getGcpComputeVpnTunnelsNewServiceForAccount(ctx context.Context, account *utilities.ExtensionConfigurationGcpAccount) (*compute.Service, string) {
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
			"tableName": "gcp_compute_vpn_tunnel",
			"projectId": projectID,
			"errString": err.Error(),
		}).Error("failed to create service")
		return nil, ""
	}
	return service, projectID
}

func (handler *GcpComputeHandler) processAccountGcpComputeVpnTunnels(ctx context.Context,
	account *utilities.ExtensionConfigurationGcpAccount) ([]map[string]string, error) {

	resultMap := make([]map[string]string, 0)

	service, projectID := handler.getGcpComputeVpnTunnelsNewServiceForAccount(ctx, account)
	if service == nil {
		return resultMap, fmt.Errorf("failed to initialize compute.Service")
	}
	myAPIService := handler.svcInterface.NewVpnTunnelsService(service)
	if myAPIService == nil {
		return resultMap, fmt.Errorf("NewVpnTunnelsService() returned nil")
	}

	aggListCall := handler.svcInterface.VpnTunnelsAggregatedList(myAPIService, projectID)
	if aggListCall == nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_compute_vpn_tunnel",
			"projectId": projectID,
		}).Debug("aggregate list call is nil")
		return resultMap, nil
	}
	itemsContainer := myGcpComputeVpnTunnelsItemsContainer{Items: make([]*compute.VpnTunnel, 0)}
	if err := handler.svcInterface.VpnTunnelsPages(ctx, aggListCall, func(page *compute.VpnTunnelAggregatedList) error {

		for _, item := range page.Items {

			itemsContainer.Items = append(itemsContainer.Items, item.VpnTunnels...)
		}

		return nil
	}); err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_compute_vpn_tunnel",
			"projectId": projectID,
			"errString": err.Error(),
		}).Error("failed to get aggregate list page")
		return resultMap, nil
	}

	byteArr, err := json.Marshal(itemsContainer)
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_compute_vpn_tunnel",
			"errString": err.Error(),
		}).Error("failed to marshal response")
		return resultMap, err
	}
	tableConfig, ok := utilities.TableConfigurationMap["gcp_compute_vpn_tunnel"]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_compute_vpn_tunnel",
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found for \"gcp_compute_vpn_tunnel\"")
	}
	jsonTable := utilities.NewTable(byteArr, tableConfig)
	for _, row := range jsonTable.Rows {
		result := extgcp.RowToMap(row, projectID, "", tableConfig)
		resultMap = append(resultMap, result)
	}

	return resultMap, nil
}
