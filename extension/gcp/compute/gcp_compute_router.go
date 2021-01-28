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

type myGcpComputeRoutersItemsContainer struct {
	Items []*compute.Router `json:"items"`
}

// GcpComputeRoutersColumns returns the list of columns for gcp_compute_router
func (handler *GcpComputeHandler) GcpComputeRoutersColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("project_id"),
		table.TextColumn("bgp"),
		table.TextColumn("bgp_peers"),
		//table.TextColumn("bgp_peers_advertise_mode"),
		//table.TextColumn("bgp_peers_advertised_groups"),
		//table.TextColumn("bgp_peers_advertised_ip_ranges"),
		//table.TextColumn("bgp_peers_advertised_ip_ranges_description"),
		//table.TextColumn("bgp_peers_advertised_ip_ranges_range"),
		//table.BigIntColumn("bgp_peers_advertised_route_priority"),
		//table.TextColumn("bgp_peers_interface_name"),
		//table.TextColumn("bgp_peers_ip_address"),
		//table.TextColumn("bgp_peers_management_type"),
		//table.TextColumn("bgp_peers_name"),
		//table.BigIntColumn("bgp_peers_peer_asn"),
		//table.TextColumn("bgp_peers_peer_ip_address"),
		//table.TextColumn("bgp_advertise_mode"),
		//table.TextColumn("bgp_advertised_groups"),
		//table.TextColumn("bgp_advertised_ip_ranges"),
		//table.TextColumn("bgp_advertised_ip_ranges_description"),
		//table.TextColumn("bgp_advertised_ip_ranges_range"),
		//table.BigIntColumn("bgp_asn"),
		table.TextColumn("creation_timestamp"),
		table.TextColumn("description"),
		table.BigIntColumn("id"),
		table.TextColumn("interfaces"),
		//table.TextColumn("interfaces_ip_range"),
		//table.TextColumn("interfaces_linked_interconnect_attachment"),
		//table.TextColumn("interfaces_linked_vpn_tunnel"),
		//table.TextColumn("interfaces_management_type"),
		//table.TextColumn("interfaces_name"),
		table.TextColumn("kind"),
		table.TextColumn("name"),
		table.TextColumn("nats"),
		//table.TextColumn("nats_drain_nat_ips"),
		//table.TextColumn("nats_enable_endpoint_independent_mapping"),
		//table.BigIntColumn("nats_icmp_idle_timeout_sec"),
		//table.TextColumn("nats_log_config"),
		//table.TextColumn("nats_log_config_enable"),
		//table.TextColumn("nats_log_config_filter"),
		//table.BigIntColumn("nats_min_ports_per_vm"),
		//table.TextColumn("nats_name"),
		//table.TextColumn("nats_nat_ip_allocate_option"),
		//table.TextColumn("nats_nat_ips"),
		//table.TextColumn("nats_source_subnetwork_ip_ranges_to_nat"),
		//table.TextColumn("nats_subnetworks"),
		//table.TextColumn("nats_subnetworks_name"),
		//table.TextColumn("nats_subnetworks_secondary_ip_range_names"),
		//table.TextColumn("nats_subnetworks_source_ip_ranges_to_nat"),
		//table.BigIntColumn("nats_tcp_established_idle_timeout_sec"),
		//table.BigIntColumn("nats_tcp_transitory_idle_timeout_sec"),
		//table.BigIntColumn("nats_udp_idle_timeout_sec"),
		table.TextColumn("network"),
		table.TextColumn("region"),
		//table.TextColumn("self_link"),

	}
}

// GcpComputeRoutersGenerate returns the rows in the table for all configured accounts
func (handler *GcpComputeHandler) GcpComputeRoutersGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	var _ = queryContext
	ctx, cancel := context.WithCancel(osqCtx)
	defer cancel()

	resultMap := make([]map[string]string, 0)

	if len(utilities.ExtConfiguration.ExtConfGcp.Accounts) == 0 {
		results, err := handler.processAccountGcpComputeRouters(ctx, nil)
		if err == nil {
			resultMap = append(resultMap, results...)
		}
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfGcp.Accounts {
			results, err := handler.processAccountGcpComputeRouters(ctx, &account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}
	return resultMap, nil
}

func (handler *GcpComputeHandler) getGcpComputeRoutersNewServiceForAccount(ctx context.Context, account *utilities.ExtensionConfigurationGcpAccount) (*compute.Service, string) {
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
			"tableName": "gcp_compute_router",
			"projectId": projectID,
			"errString": err.Error(),
		}).Error("failed to create service")
		return nil, ""
	}
	return service, projectID
}

func (handler *GcpComputeHandler) processAccountGcpComputeRouters(ctx context.Context,
	account *utilities.ExtensionConfigurationGcpAccount) ([]map[string]string, error) {

	resultMap := make([]map[string]string, 0)

	service, projectID := handler.getGcpComputeRoutersNewServiceForAccount(ctx, account)
	if service == nil {
		return resultMap, fmt.Errorf("failed to initialize compute.Service")
	}
	myAPIService := handler.svcInterface.NewRoutersService(service)
	if myAPIService == nil {
		return resultMap, fmt.Errorf("NewRoutersService() returned nil")
	}

	aggListCall := handler.svcInterface.RoutersAggregatedList(myAPIService, projectID)
	if aggListCall == nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_compute_router",
			"projectId": projectID,
		}).Debug("aggregate list call is nil")
		return resultMap, nil
	}
	itemsContainer := myGcpComputeRoutersItemsContainer{Items: make([]*compute.Router, 0)}
	if err := handler.svcInterface.RoutersPages(ctx, aggListCall, func(page *compute.RouterAggregatedList) error {

		for _, item := range page.Items {

			itemsContainer.Items = append(itemsContainer.Items, item.Routers...)
		}

		return nil
	}); err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_compute_router",
			"projectId": projectID,
			"errString": err.Error(),
		}).Error("failed to get aggregate list page")
		return resultMap, nil
	}

	byteArr, err := json.Marshal(itemsContainer)
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_compute_router",
			"errString": err.Error(),
		}).Error("failed to marshal response")
		return resultMap, err
	}
	tableConfig, ok := utilities.TableConfigurationMap["gcp_compute_router"]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_compute_router",
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found for \"gcp_compute_router\"")
	}
	jsonTable := utilities.NewTable(byteArr, tableConfig)
	for _, row := range jsonTable.Rows {
		result := extgcp.RowToMap(row, projectID, "", tableConfig)
		resultMap = append(resultMap, result)
	}

	return resultMap, nil
}
