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

type myGcpComputeInterconnectsItemsContainer struct {
	Items []*compute.Interconnect `json:"items"`
}

func (handler *GcpComputeHandler) GcpComputeInterconnectsColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("project_id"),
		table.TextColumn("admin_enabled"),
		table.TextColumn("circuit_infos"),
		//table.TextColumn("circuit_infos_customer_demarc_id"),
		//table.TextColumn("circuit_infos_google_circuit_id"),
		//table.TextColumn("circuit_infos_google_demarc_id"),
		table.TextColumn("creation_timestamp"),
		table.TextColumn("customer_name"),
		table.TextColumn("description"),
		table.TextColumn("expected_outages"),
		//table.TextColumn("expected_outages_affected_circuits"),
		//table.TextColumn("expected_outages_description"),
		//table.BigIntColumn("expected_outages_end_time"),
		//table.TextColumn("expected_outages_issue_type"),
		//table.TextColumn("expected_outages_name"),
		//table.TextColumn("expected_outages_source"),
		//table.BigIntColumn("expected_outages_start_time"),
		//table.TextColumn("expected_outages_state"),
		table.TextColumn("google_ip_address"),
		table.TextColumn("google_reference_id"),
		table.BigIntColumn("id"),
		table.TextColumn("interconnect_attachments"),
		table.TextColumn("interconnect_type"),
		table.TextColumn("kind"),
		table.TextColumn("link_type"),
		table.TextColumn("location"),
		table.TextColumn("name"),
		table.TextColumn("noc_contact_email"),
		table.TextColumn("operational_status"),
		table.TextColumn("peer_ip_address"),
		table.BigIntColumn("provisioned_link_count"),
		table.BigIntColumn("requested_link_count"),
		//table.TextColumn("self_link"),
		table.TextColumn("state"),
	}
}

func (handler *GcpComputeHandler) GcpComputeInterconnectsGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	var _ = queryContext
	ctx, cancel := context.WithCancel(osqCtx)
	defer cancel()

	resultMap := make([]map[string]string, 0)

	if len(utilities.ExtConfiguration.ExtConfGcp.Accounts) == 0 {
		results, err := handler.processAccountGcpComputeInterconnects(ctx, nil)
		if err == nil {
			resultMap = append(resultMap, results...)
		}
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfGcp.Accounts {
			results, err := handler.processAccountGcpComputeInterconnects(ctx, &account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}
	return resultMap, nil
}

func (handler *GcpComputeHandler) getGcpComputeInterconnectsNewServiceForAccount(ctx context.Context, account *utilities.ExtensionConfigurationGcpAccount) (*compute.Service, string) {
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
			"tableName": "gcp_compute_interconnect",
			"projectId": projectID,
			"errString": err.Error(),
		}).Error("failed to create service")
		return nil, ""
	}
	return service, projectID
}

func (handler *GcpComputeHandler) processAccountGcpComputeInterconnects(ctx context.Context,
	account *utilities.ExtensionConfigurationGcpAccount) ([]map[string]string, error) {

	resultMap := make([]map[string]string, 0)

	service, projectID := handler.getGcpComputeInterconnectsNewServiceForAccount(ctx, account)
	if service == nil {
		return resultMap, fmt.Errorf("failed to initialize compute.Service")
	}
	myApiService := handler.svcInterface.NewInterconnectsService(service)
	if myApiService == nil {
		return resultMap, fmt.Errorf("NewInterconnectsService() returned nil")
	}

	aggListCall := handler.svcInterface.InterconnectsList(myApiService, projectID)
	if aggListCall == nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_compute_interconnect",
			"projectId": projectID,
		}).Debug("aggregate list call is nil")
		return resultMap, nil
	}
	itemsContainer := myGcpComputeInterconnectsItemsContainer{Items: make([]*compute.Interconnect, 0)}
	if err := handler.svcInterface.InterconnectsPages(aggListCall, ctx, func(page *compute.InterconnectList) error {

		itemsContainer.Items = append(itemsContainer.Items, page.Items...)

		return nil
	}); err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_compute_interconnect",
			"projectId": projectID,
			"errString": err.Error(),
		}).Error("failed to get aggregate list page")
		return resultMap, nil
	}

	byteArr, err := json.Marshal(itemsContainer)
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_compute_interconnect",
			"errString": err.Error(),
		}).Error("failed to marshal response")
		return resultMap, err
	}
	tableConfig, ok := utilities.TableConfigurationMap["gcp_compute_interconnect"]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_compute_interconnect",
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found for \"gcp_compute_interconnect\"")
	}
	jsonTable := utilities.NewTable(byteArr, tableConfig)
	for _, row := range jsonTable.Rows {
		result := extgcp.RowToMap(row, projectID, "", tableConfig)
		resultMap = append(resultMap, result)
	}

	return resultMap, nil
}
