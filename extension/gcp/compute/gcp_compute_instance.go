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

type myGcpComputeInstancesItemsContainer struct {
	Items []*compute.Instance `json:"items"`
}

func (handler *GcpComputeHandler) GcpComputeInstancesColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("project_id"),
		table.TextColumn("can_ip_forward"),
		//table.TextColumn("confidential_instance_config"),
		table.TextColumn("confidential_instance_config_enable_confidential_compute"),
		table.TextColumn("cpu_platform"),
		table.TextColumn("creation_timestamp"),
		table.TextColumn("deletion_protection"),
		table.TextColumn("description"),
		table.TextColumn("disks"),
		//table.TextColumn("display_device"),
		table.TextColumn("display_device_enable_display"),
		table.TextColumn("fingerprint"),
		table.TextColumn("guest_accelerators"),
		table.TextColumn("hostname"),
		table.BigIntColumn("id"),
		table.TextColumn("kind"),
		table.TextColumn("label_fingerprint"),
		table.TextColumn("labels"),
		table.TextColumn("last_start_timestamp"),
		table.TextColumn("last_stop_timestamp"),
		table.TextColumn("last_suspended_timestamp"),
		table.TextColumn("machine_type"),
		//table.TextColumn("metadata"),
		table.TextColumn("metadata_fingerprint"),
		table.TextColumn("metadata_items"),
		table.TextColumn("metadata_kind"),
		table.TextColumn("min_cpu_platform"),
		table.TextColumn("name"),
		table.TextColumn("network_interfaces"),
		table.TextColumn("private_ipv6_google_access"),
		//table.TextColumn("reservation_affinity"),
		table.TextColumn("reservation_affinity_consume_reservation_type"),
		table.TextColumn("reservation_affinity_key"),
		table.TextColumn("reservation_affinity_values"),
		table.TextColumn("resource_policies"),
		//table.TextColumn("scheduling"),
		table.TextColumn("scheduling_automatic_restart"),
		table.BigIntColumn("scheduling_min_node_cpus"),
		table.TextColumn("scheduling_node_affinities"),
		table.TextColumn("scheduling_on_host_maintenance"),
		table.TextColumn("scheduling_preemptible"),
		//table.TextColumn("self_link"),
		table.TextColumn("service_accounts"),
		//table.TextColumn("shielded_instance_config"),
		table.TextColumn("shielded_instance_config_enable_integrity_monitoring"),
		table.TextColumn("shielded_instance_config_enable_secure_boot"),
		table.TextColumn("shielded_instance_config_enable_vtpm"),
		//table.TextColumn("shielded_instance_integrity_policy"),
		table.TextColumn("shielded_instance_integrity_policy_update_auto_learn_policy"),
		table.TextColumn("start_restricted"),
		table.TextColumn("status"),
		table.TextColumn("status_message"),
		table.TextColumn("tags"),
		table.TextColumn("tags_fingerprint"),
		table.TextColumn("tags_items"),
		table.TextColumn("zone"),
	}
}

func (handler *GcpComputeHandler) GcpComputeInstancesGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	var _ = queryContext
	ctx, cancel := context.WithCancel(osqCtx)
	defer cancel()

	resultMap := make([]map[string]string, 0)

	if len(utilities.ExtConfiguration.ExtConfGcp.Accounts) == 0 {
		results, err := handler.processAccountGcpComputeInstances(ctx, nil)
		if err == nil {
			resultMap = append(resultMap, results...)
		}
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfGcp.Accounts {
			results, err := handler.processAccountGcpComputeInstances(ctx, &account)
			if err != nil {
				// TODO: Continue to next account or return error ?
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}
	return resultMap, nil
}

func (handler *GcpComputeHandler) getGcpComputeInstancesNewServiceForAccount(ctx context.Context, account *utilities.ExtensionConfigurationGcpAccount) (*compute.Service, string) {
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
		fmt.Println("NewService() error: ", err)
		return nil, ""
	}
	return service, projectID
}

func (handler *GcpComputeHandler) processAccountGcpComputeInstances(ctx context.Context,
	account *utilities.ExtensionConfigurationGcpAccount) ([]map[string]string, error) {

	resultMap := make([]map[string]string, 0)

	service, projectID := handler.getGcpComputeInstancesNewServiceForAccount(ctx, account)
	if service == nil {
		return resultMap, fmt.Errorf("failed to initialize compute.Service")
	}
	myApiService := handler.svcInterface.NewInstancesService(service)
	if myApiService == nil {
		fmt.Println("NewInstancesService() returned nil")
		return resultMap, fmt.Errorf("NewInstancesService() returned nil")
	}

	aggListCall := handler.svcInterface.InstancesAggregatedList(myApiService, projectID)
	if aggListCall == nil {
		fmt.Println("aggListCall is nil")
		return resultMap, nil
	}
	itemsContainer := myGcpComputeInstancesItemsContainer{Items: make([]*compute.Instance, 0)}
	if err := handler.svcInterface.InstancesPages(aggListCall, ctx, func(page *compute.InstanceAggregatedList) error {

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
		result := extgcp.RowToMap(row, projectID, "", tableConfig)
		resultMap = append(resultMap, result)
	}

	return resultMap, nil
}
