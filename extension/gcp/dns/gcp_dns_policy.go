package dns

import (
	"context"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"

	extgcp "github.com/Uptycs/cloudquery/extension/gcp"
	"github.com/Uptycs/cloudquery/utilities"
	"github.com/kolide/osquery-go/plugin/table"

	"google.golang.org/api/option"

	gcpdns "google.golang.org/api/dns/v1beta2"
)

type myGcpDnsPoliciesItemsContainer struct {
	Items []*gcpdns.Policy `json:"items"`
}

func GcpDnsPoliciesColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("project_id"),
		table.TextColumn("alternative_name_server_config"),
		//table.TextColumn("alternative_name_server_config_kind"),
		//table.TextColumn("alternative_name_server_config_target_name_servers"),
		//table.TextColumn("alternative_name_server_config_target_name_servers_forwarding_path"),
		//table.TextColumn("alternative_name_server_config_target_name_servers_ipv4_address"),
		//table.TextColumn("alternative_name_server_config_target_name_servers_kind"),
		table.TextColumn("description"),
		table.TextColumn("enable_inbound_forwarding"),
		table.TextColumn("enable_logging"),
		table.BigIntColumn("id"),
		table.TextColumn("kind"),
		table.TextColumn("name"),
		table.TextColumn("networks"),
		//table.TextColumn("networks_kind"),
		//table.TextColumn("networks_network_url"),

	}
}

func GcpDnsPoliciesGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	var _ = queryContext
	ctx, cancel := context.WithCancel(osqCtx)
	defer cancel()

	resultMap := make([]map[string]string, 0)

	if len(utilities.ExtConfiguration.ExtConfGcp.Accounts) == 0 {
		results, err := processAccountGcpDnsPolicies(ctx, nil)
		if err == nil {
			resultMap = append(resultMap, results...)
		}
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfGcp.Accounts {
			results, err := processAccountGcpDnsPolicies(ctx, &account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}
	return resultMap, nil
}

func getGcpDnsPoliciesNewServiceForAccount(ctx context.Context, account *utilities.ExtensionConfigurationGcpAccount) (*gcpdns.Service, string) {
	var projectID = ""
	var service *gcpdns.Service
	var err error
	if account != nil {
		projectID = account.ProjectId
		service, err = gcpdns.NewService(ctx, option.WithCredentialsFile(account.KeyFile))
	} else {
		projectID = utilities.DefaultGcpProjectID
		service, err = gcpdns.NewService(ctx)
	}
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_dns_policy",
			"projectId": projectID,
			"errString": err.Error(),
		}).Error("failed to create service")
		return nil, ""
	}
	return service, projectID
}

func processAccountGcpDnsPolicies(ctx context.Context,
	account *utilities.ExtensionConfigurationGcpAccount) ([]map[string]string, error) {

	resultMap := make([]map[string]string, 0)

	service, projectID := getGcpDnsPoliciesNewServiceForAccount(ctx, account)
	if service == nil {
		return resultMap, fmt.Errorf("failed to initialize gcpdns.Service")
	}

	listCall := service.Policies.List(projectID)
	if listCall == nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_dns_policy",
			"projectId": projectID,
		}).Debug("list call is nil")
		return resultMap, nil
	}
	itemsContainer := myGcpDnsPoliciesItemsContainer{Items: make([]*gcpdns.Policy, 0)}
	if err := listCall.Pages(ctx, func(page *gcpdns.PoliciesListResponse) error {

		itemsContainer.Items = append(itemsContainer.Items, page.Policies...)

		return nil
	}); err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_dns_policy",
			"projectId": projectID,
			"errString": err.Error(),
		}).Error("failed to get aggregate list page")
		return resultMap, nil
	}

	byteArr, err := json.Marshal(itemsContainer)
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_dns_policy",
			"errString": err.Error(),
		}).Error("failed to marshal response")
		return resultMap, err
	}
	tableConfig, ok := utilities.TableConfigurationMap["gcp_dns_policy"]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_dns_policy",
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found for \"gcp_dns_policy\"")
	}
	jsonTable := utilities.NewTable(byteArr, tableConfig)
	for _, row := range jsonTable.Rows {
		result := extgcp.RowToMap(row, projectID, "", tableConfig)
		resultMap = append(resultMap, result)
	}

	return resultMap, nil
}
