package compute

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	extgcp "github.com/Uptycs/cloudquery/extension/gcp"
	"github.com/Uptycs/cloudquery/utilities"
	"github.com/kolide/osquery-go/plugin/table"

	"google.golang.org/api/option"

	compute "google.golang.org/api/compute/v1"
)

type myGcpComputeNetworkItemsContainer struct {
	Items []*compute.Network `json:"items"`
}

func GcpComputeNetworkColumns() []table.ColumnDefinition {
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
		//table.TextColumn("routing_config"),
		table.TextColumn("routing_config_routing_mode"),
		//table.TextColumn("self_link"),
		table.TextColumn("subnetworks"),
	}
}

func GcpComputeNetworkGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	var _ = queryContext
	ctx, cancel := context.WithCancel(osqCtx)
	defer cancel()

	resultMap := make([]map[string]string, 0)

	if len(utilities.ExtConfiguration.ExtConfGcp.Accounts) == 0 {
		results, err := processAccountGcpComputeNetwork(ctx, nil)
		if err == nil {
			resultMap = append(resultMap, results...)
		}
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfGcp.Accounts {
			results, err := processAccountGcpComputeNetwork(ctx, &account)
			if err != nil {
				// TODO: Continue to next account or return error ?
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}
	return resultMap, nil
}

func getGcpComputeNetworkNewServiceForAccount(ctx context.Context, account *utilities.ExtensionConfigurationGcpAccount) (*compute.Service, string) {
	var projectID = ""
	var service *compute.Service
	var err error
	if account != nil {
		projectID = account.ProjectId
		service, err = compute.NewService(ctx, option.WithCredentialsFile(account.KeyFile))
	} else {
		projectID = utilities.DefaultGcpProjectID
		service, err = compute.NewService(ctx)
	}
	if err != nil {
		fmt.Println("NewService() error: ", err)
		return nil, ""
	}
	return service, projectID
}

func processAccountGcpComputeNetwork(ctx context.Context,
	account *utilities.ExtensionConfigurationGcpAccount) ([]map[string]string, error) {

	resultMap := make([]map[string]string, 0)

	service, projectID := getGcpComputeNetworkNewServiceForAccount(ctx, account)
	if service == nil {
		return resultMap, fmt.Errorf("failed to initialize compute.Service")
	}
	myApiService := compute.NewNetworksService(service)
	if myApiService == nil {
		fmt.Println("compute.NewNetworksService() returned nil")
		return resultMap, fmt.Errorf("compute.NewNetworksService() returned nil")
	}

	aggListCall := myApiService.List(projectID)
	if aggListCall == nil {
		fmt.Println("aggListCall is nil")
		return resultMap, nil
	}
	itemsContainer := myGcpComputeNetworkItemsContainer{Items: make([]*compute.Network, 0)}
	if err := aggListCall.Pages(ctx, func(page *compute.NetworkList) error {

		itemsContainer.Items = append(itemsContainer.Items, page.Items...)

		return nil
	}); err != nil {
		fmt.Println("aggListCal.Page: ", err)
		return resultMap, nil
	}

	byteArr, err := json.Marshal(itemsContainer)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
	//fmt.Printf("%+v\n", string(byteArr))
	tableConfig, ok := utilities.TableConfigurationMap["gcp_compute_network"]
	if !ok {
		fmt.Println("getTableConfig: ", err)
		return resultMap, fmt.Errorf("table configuration not found")
	}
	jsonTable := utilities.Table{}
	jsonTable.Init(byteArr, tableConfig.MaxLevel, tableConfig.GetParsedAttributeConfigMap())
	for _, row := range jsonTable.Rows {
		result := extgcp.RowToMap(row, projectID, "", tableConfig)
		resultMap = append(resultMap, result)
	}

	return resultMap, nil
}
