package compute

import (
	"context"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"

	extgcp "github.com/Uptycs/cloudquery/extension/gcp"
	"github.com/Uptycs/cloudquery/utilities"
	"github.com/kolide/osquery-go/plugin/table"

	"google.golang.org/api/option"

	compute "google.golang.org/api/compute/v1"
)

type myGcpComputeNetworksItemsContainer struct {
	Items []*compute.Network `json:"items"`
}

func (handler *GcpComputeHandler) GcpComputeNetworksColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("project_id"),
		table.TextColumn("ipv4_range"),
		table.TextColumn("auto_create_subnetworks"),
		table.TextColumn("creation_timestamp"),
		table.TextColumn("description"),
		table.TextColumn("gateway_ipv4"),
		table.BigIntColumn("id"),
		table.TextColumn("kind"),
		table.BigIntColumn("mtu"),
		table.TextColumn("name"),
		table.TextColumn("peerings"),
		//table.TextColumn("peerings_auto_create_routes"),
		//table.TextColumn("peerings_exchange_subnet_routes"),
		//table.TextColumn("peerings_export_custom_routes"),
		//table.TextColumn("peerings_export_subnet_routes_with_public_ip"),
		//table.TextColumn("peerings_import_custom_routes"),
		//table.TextColumn("peerings_import_subnet_routes_with_public_ip"),
		//table.TextColumn("peerings_name"),
		//table.TextColumn("peerings_network"),
		//table.BigIntColumn("peerings_peer_mtu"),
		//table.TextColumn("peerings_state"),
		//table.TextColumn("peerings_state_details"),
		table.TextColumn("routing_config"),
		//table.TextColumn("routing_config_routing_mode"),
		//table.TextColumn("self_link"),
		table.TextColumn("subnetworks"),
	}
}

func (handler *GcpComputeHandler) GcpComputeNetworksGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	var _ = queryContext
	ctx, cancel := context.WithCancel(osqCtx)
	defer cancel()

	resultMap := make([]map[string]string, 0)

	if len(utilities.ExtConfiguration.ExtConfGcp.Accounts) == 0 {
		results, err := handler.processAccountGcpComputeNetworks(ctx, nil)
		if err == nil {
			resultMap = append(resultMap, results...)
		}
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfGcp.Accounts {
			results, err := handler.processAccountGcpComputeNetworks(ctx, &account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}
	return resultMap, nil
}

func (handler *GcpComputeHandler) getGcpComputeNetworksNewServiceForAccount(ctx context.Context, account *utilities.ExtensionConfigurationGcpAccount) (*compute.Service, string) {
	var projectID = ""
	var service *compute.Service
	var err error
	if account != nil {
		projectID = account.ProjectId
		service, err = handler.svcInterface.NewService(ctx, option.WithCredentialsFile(account.KeyFile))
	} else {
		projectID = utilities.DefaultGcpProjectID
		service, err = handler.svcInterface.NewService(ctx)
	}
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_compute_network",
			"projectId": projectID,
			"errString": err.Error(),
		}).Error("failed to create service")
		return nil, ""
	}
	return service, projectID
}

func (handler *GcpComputeHandler) processAccountGcpComputeNetworks(ctx context.Context,
	account *utilities.ExtensionConfigurationGcpAccount) ([]map[string]string, error) {

	resultMap := make([]map[string]string, 0)

	service, projectID := handler.getGcpComputeNetworksNewServiceForAccount(ctx, account)
	if service == nil {
		return resultMap, fmt.Errorf("failed to initialize compute.Service")
	}
	myApiService := handler.svcInterface.NewNetworksService(service)
	if myApiService == nil {
		return resultMap, fmt.Errorf("NewNetworksService() returned nil")
	}

	aggListCall := handler.svcInterface.NetworksList(myApiService, projectID)
	if aggListCall == nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_compute_network",
			"projectId": projectID,
		}).Debug("aggregate list call is nil")
		return resultMap, nil
	}
	itemsContainer := myGcpComputeNetworksItemsContainer{Items: make([]*compute.Network, 0)}
	if err := handler.svcInterface.NetworksPages(aggListCall, ctx, func(page *compute.NetworkList) error {

		itemsContainer.Items = append(itemsContainer.Items, page.Items...)

		return nil
	}); err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_compute_network",
			"projectId": projectID,
			"errString": err.Error(),
		}).Error("failed to get aggregate list page")
		return resultMap, nil
	}

	byteArr, err := json.Marshal(itemsContainer)
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_compute_network",
			"errString": err.Error(),
		}).Error("failed to marshal response")
		return resultMap, err
	}
	tableConfig, ok := utilities.TableConfigurationMap["gcp_compute_network"]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_compute_network",
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found for \"gcp_compute_network\"")
	}
	jsonTable := utilities.NewTable(byteArr, tableConfig)
	for _, row := range jsonTable.Rows {
		result := extgcp.RowToMap(row, projectID, "", tableConfig)
		resultMap = append(resultMap, result)
	}

	return resultMap, nil
}
