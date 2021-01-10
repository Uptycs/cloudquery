package compute

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	extgcp "github.com/Uptycs/cloudquery/extension/gcp"
	"github.com/Uptycs/cloudquery/utilities"
	"github.com/kolide/osquery-go/plugin/table"

	"google.golang.org/api/option"

	compute "google.golang.org/api/compute/v1"
)

type myGcpComputeInstanceItemsContainer struct {
	Items []*compute.Instance `json:"items"`
}

func GcpComputeInstanceColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("project_id"),
		//table.TextColumn("guest_accelerators"),
		table.TextColumn("min_cpu_platform"),
		table.TextColumn("label_fingerprint"),
		table.TextColumn("network_interfaces"),
		table.TextColumn("service_accounts"),
		//table.TextColumn("shielded_instance_integrity_policy"),
		table.TextColumn("shielded_instance_integrity_policy_update_auto_learn_policy"),
		table.TextColumn("start_restricted"),
		table.TextColumn("tags"),
		//table.TextColumn("tags_fingerprint"),
		//table.TextColumn("tags_items"),
		//table.TextColumn("confidential_instance_config"),
		table.TextColumn("confidential_instance_config_enable_confidential_compute"),
		table.TextColumn("id"),
		table.TextColumn("name"),
		table.TextColumn("private_ipv6_google_access"),
		table.TextColumn("creation_timestamp"),
		table.TextColumn("description"),
		table.TextColumn("reservation_affinity"),
		table.TextColumn("reservation_affinity_consume_reservation_type"),
		table.TextColumn("reservation_affinity_key"),
		table.TextColumn("reservation_affinity_values"),
		table.TextColumn("kind"),
		table.TextColumn("last_start_timestamp"),
		table.TextColumn("resource_policies"),
		table.TextColumn("cpu_platform"),
		table.TextColumn("disks"),
		//table.TextColumn("display_device"),
		table.TextColumn("display_device_enable_display"),
		table.TextColumn("hostname"),
		table.TextColumn("status_message"),
		table.TextColumn("zone"),
		table.TextColumn("deletion_protection"),
		table.TextColumn("last_stop_timestamp"),
		table.TextColumn("last_suspended_timestamp"),
		table.TextColumn("scheduling"),
		//table.TextColumn("scheduling_preemptible"),
		//table.TextColumn("scheduling_automatic_restart"),
		//table.TextColumn("scheduling_location_hint"),
		//table.IntegerColumn("scheduling_min_node_cpus"),
		//table.TextColumn("scheduling_node_affinities"),
		//table.TextColumn("scheduling_on_host_maintenance"),
		table.TextColumn("metadata"),
		//table.TextColumn("metadata_kind"),
		//table.TextColumn("metadata_fingerprint"),
		//table.TextColumn("metadata_items"),
		//table.TextColumn("self_link"),
		table.TextColumn("shielded_instance_config"),
		//table.TextColumn("shielded_instance_config_enable_secure_boot"),
		//table.TextColumn("shielded_instance_config_enable_vtpm"),
		//table.TextColumn("shielded_instance_config_enable_integrity_monitoring"),
		table.TextColumn("status"),
		table.TextColumn("can_ip_forward"),
		table.TextColumn("fingerprint"),
		table.TextColumn("machine_type"),
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

	aggListCall := myApiService.AggregatedList(account.ProjectId)
	if aggListCall == nil {
		fmt.Println("aggListCall is nil")
		return resultMap, nil
	}
	itemsContainer := myGcpComputeInstanceItemsContainer{Items: make([]*compute.Instance, 0)}
	if err := aggListCall.Pages(ctx, func(page *compute.InstanceAggregatedList) error {

		for _, item := range page.Items {
			for _, inst := range item.Instances {
				zonePathSplit := strings.Split(inst.Zone, "/")
				inst.Zone = zonePathSplit[len(zonePathSplit)-1]
			}
			itemsContainer.Items = append(itemsContainer.Items, item.Instances...)
		}

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
	tableConfig, ok := utilities.TableConfigurationMap["gcp_compute_instance"]
	if !ok {
		fmt.Println("getTableConfig: ", err)
		return resultMap, fmt.Errorf("table configuration not found")
	}
	jsonTable := utilities.Table{}
	jsonTable.Init(byteArr, tableConfig.MaxLevel, tableConfig.GetParsedAttributeConfigMap())
	for _, row := range jsonTable.Rows {
		result := extgcp.RowToMap(row, account.ProjectId, "", tableConfig)
		resultMap = append(resultMap, result)
	}

	return resultMap, nil
}
