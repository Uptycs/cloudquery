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

func GcpComputeInstanceColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("project_id"),
		table.BigIntColumn("id"),
		table.TextColumn("hostname"),
		table.TextColumn("name"),
		table.TextColumn("status"),
		table.TextColumn("kind"),
		table.TextColumn("tags"),
		table.TextColumn("creation_timestamp"),
		table.TextColumn("description"),
		table.TextColumn("machine_type"),
		table.TextColumn("status_message"),
		table.TextColumn("zone"),
		table.TextColumn("can_ip_forward"),
		table.TextColumn("cpu_platform"),
		table.TextColumn("label_fingerprint"),
		table.TextColumn("min_cpu_platform"),
		table.TextColumn("start_restricted"),
		table.TextColumn("deletion_protection"),
		table.TextColumn("fingerprint"),
		table.TextColumn("private_ipv6Google_access"),
		table.TextColumn("last_start_timestamp"),
		table.TextColumn("last_stop_timestamp"),
		table.TextColumn("last_suspended_timestamp"),
		table.TextColumn("display_device_enable_display"),
		table.TextColumn("reservation_affinity_consume_reservation_type"),
		table.TextColumn("reservation_affinity_key"),
		table.TextColumn("scheduling_on_host_maintenance"),
		table.TextColumn("scheduling_automatic_restart"),
		table.TextColumn("scheduling_preemptible"),
		table.TextColumn("scheduling_min_node_cpus"),
		table.TextColumn("shielded_instance_config_enable_secure_boot"),
		table.TextColumn("shielded_instance_config_enable_vtpm"),
		table.TextColumn("shielded_instance_config_enable_integrity_monitoring"),
		table.TextColumn("shielded_instance_integrity_policy_update_auto_learn_policy"),
		table.TextColumn("confidential_instance_config_enable_confidential_compute"),
	}
}

func GcpComputeInstanceGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	var _ = queryContext
	ctx, cancel := context.WithCancel(osqCtx)
	defer cancel()

	resultMap := make([]map[string]string, 0)

	for _, account := range utilities.ExtConfiguration.ExtConfGcp.Accounts {
		results, err := processAccountGcpComputeInstance(ctx, &account)
		if err != nil {
			// TODO: Continue to next account or return error ?
			continue
		}
		resultMap = append(resultMap, results...)
	}
	return resultMap, nil
}

func processAccountGcpComputeInstance(ctx context.Context,
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

	zoneList := extgcp.GetZones(ctx, myApiService, account.ProjectId)
	for _, zone := range zoneList {
		listCall := myApiService.List(account.ProjectId, zone)
		if listCall == nil {
			fmt.Println("listCall is nil")
			return resultMap, nil
		}
		if err := listCall.Pages(ctx, func(page *compute.InstanceList) error {
			byteArr, err := json.Marshal(page)
			if err != nil {
				fmt.Printf("error: %v\n", err)
				os.Exit(1)
			}
			//fmt.Printf("%+v\n", string(byteArr))
			jsonTable := utilities.Table{}
			jsonTable.Init(byteArr, tableConfig.MaxLevel, tableConfig.GetParsedAttributeConfigMap())
			for _, row := range jsonTable.Rows {
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
