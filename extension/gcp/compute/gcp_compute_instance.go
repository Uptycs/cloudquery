package compute

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	extgcp "github.com/Uptycs/cloudquery/extension/gcp"
	"github.com/Uptycs/cloudquery/utilities"
	"github.com/kolide/osquery-go/plugin/table"

	"google.golang.org/api/option"

	compute "google.golang.org/api/compute/v1"
)

func GcpComputeInstanceColumns() []table.ColumnDefinition {
	var _, _ = strconv.Atoi("123") // Disables warning when strcov is not used
	return []table.ColumnDefinition{
		table.TextColumn("project_id"),
		table.TextColumn("zone"),
		table.TextColumn("id"),
		table.TextColumn("hostname"),
		table.TextColumn("name"),
		table.TextColumn("status"),
		table.TextColumn("kind"),
		table.TextColumn("tags"),
	}
}

func GcpComputeInstanceGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	var _ = queryContext
	ctx, cancel := context.WithCancel(osqCtx)
	defer cancel()

	resultMap := make([]map[string]string, 0)

	for _, account := range utilities.ExtConfiguration.ExtConfGcp.Accounts {
		results, err := GcpComputeInstanceProcessAccount(ctx, &account)
		if err != nil {
			// TODO: Continue to next account or return error ?
			continue
		}
		resultMap = append(resultMap, results...)
	}
	return resultMap, nil
}

func GcpComputeInstanceProcessAccount(ctx context.Context,
	account *utilities.ExtensionConfigurationGcpAccount) ([]map[string]string, error) {

	resultMap := make([]map[string]string, 0)

	service, err := compute.NewService(ctx, option.WithCredentialsFile(account.KeyFile))
	if err != nil {
		fmt.Println("NewService() error: ", err)
		return resultMap, err
	}
	myApiService := compute.NewInstancesService(service)
	if myApiService == nil {
		fmt.Println("compute.NewInstancesService() returned nil")
		return resultMap, fmt.Errorf("compute.NewInstancesService() returned nil")
	}

	tableConfig, ok := utilities.TableConfigurationMap["gcp_compute_instance"]
	if !ok {
		fmt.Println("getTableConfig: ", err)
		return resultMap, fmt.Errorf("table configuration not found")
	}

	for _, zone := range tableConfig.Gcp.Zones {
		listCall := myApiService.List(account.ProjectId, zone)
		if listCall == nil {
			fmt.Println("listCall is nil")
			return resultMap, nil
		}
		if err := listCall.Pages(ctx, func(page *compute.InstanceList) error {
			byteArr, err := json.Marshal(page)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				os.Exit(1)
			}
			table := utilities.Table{}
			table.Init(byteArr, tableConfig.MaxLevel, tableConfig.GetParsedAttributeConfigMap())
			for _, row := range table.Rows {
				result := extgcp.RowToMap(row, account.ProjectId, zone, tableConfig)
				resultMap = append(resultMap, result)
			}
			return nil
		}); err != nil {
			fmt.Println("listCall.Page: ", err)
			//log.Fatal(err)
			return resultMap, err
		}
	}

	return resultMap, nil
}
