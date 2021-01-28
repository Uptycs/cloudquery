package compute

import (
	"context"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"strings"

	extgcp "github.com/Uptycs/cloudquery/extension/gcp"
	"github.com/Uptycs/cloudquery/utilities"
	"github.com/kolide/osquery-go/plugin/table"

	"google.golang.org/api/option"

	compute "google.golang.org/api/compute/v1"
)

type myGcpComputeReservationsItemsContainer struct {
	Items []*compute.Reservation `json:"items"`
}

// GcpComputeReservationsColumns returns the list of columns for gcp_compute_reservation
func (handler *GcpComputeHandler) GcpComputeReservationsColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("project_id"),
		table.TextColumn("commitment"),
		table.TextColumn("creation_timestamp"),
		table.TextColumn("description"),
		table.BigIntColumn("id"),
		table.TextColumn("kind"),
		table.TextColumn("name"),
		//table.TextColumn("self_link"),
		table.TextColumn("specific_reservation"),
		table.TextColumn("specific_reservation_required"),
		//table.BigIntColumn("specific_reservation_count"),
		//table.BigIntColumn("specific_reservation_in_use_count"),
		//table.TextColumn("specific_reservation_instance_properties"),
		//table.TextColumn("specific_reservation_instance_properties_guest_accelerators"),
		//table.BigIntColumn("specific_reservation_instance_properties_guest_accelerators_accelerator_count"),
		//table.TextColumn("specific_reservation_instance_properties_guest_accelerators_accelerator_type"),
		//table.TextColumn("specific_reservation_instance_properties_local_ssds"),
		//table.BigIntColumn("specific_reservation_instance_properties_local_ssds_disk_size_gb"),
		//table.TextColumn("specific_reservation_instance_properties_local_ssds_interface"),
		//table.TextColumn("specific_reservation_instance_properties_machine_type"),
		//table.TextColumn("specific_reservation_instance_properties_min_cpu_platform"),
		table.TextColumn("status"),
		table.TextColumn("zone"),
	}
}

// GcpComputeReservationsGenerate returns the rows in the table for all configured accounts
func (handler *GcpComputeHandler) GcpComputeReservationsGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	var _ = queryContext
	ctx, cancel := context.WithCancel(osqCtx)
	defer cancel()

	resultMap := make([]map[string]string, 0)

	if len(utilities.ExtConfiguration.ExtConfGcp.Accounts) == 0 {
		results, err := handler.processAccountGcpComputeReservations(ctx, nil)
		if err == nil {
			resultMap = append(resultMap, results...)
		}
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfGcp.Accounts {
			results, err := handler.processAccountGcpComputeReservations(ctx, &account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}
	return resultMap, nil
}

func (handler *GcpComputeHandler) getGcpComputeReservationsNewServiceForAccount(ctx context.Context, account *utilities.ExtensionConfigurationGcpAccount) (*compute.Service, string) {
	var projectID = ""
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
			"tableName": "gcp_compute_reservation",
			"projectId": projectID,
			"errString": err.Error(),
		}).Error("failed to create service")
		return nil, ""
	}
	return service, projectID
}

func (handler *GcpComputeHandler) processAccountGcpComputeReservations(ctx context.Context,
	account *utilities.ExtensionConfigurationGcpAccount) ([]map[string]string, error) {

	resultMap := make([]map[string]string, 0)

	service, projectID := handler.getGcpComputeReservationsNewServiceForAccount(ctx, account)
	if service == nil {
		return resultMap, fmt.Errorf("failed to initialize compute.Service")
	}
	myAPIService := handler.svcInterface.NewReservationsService(service)
	if myAPIService == nil {
		return resultMap, fmt.Errorf("NewReservationsService() returned nil")
	}

	aggListCall := handler.svcInterface.ReservationsAggregatedList(myAPIService, projectID)
	if aggListCall == nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_compute_reservation",
			"projectId": projectID,
		}).Debug("aggregate list call is nil")
		return resultMap, nil
	}
	itemsContainer := myGcpComputeReservationsItemsContainer{Items: make([]*compute.Reservation, 0)}
	if err := handler.svcInterface.ReservationsPages(ctx, aggListCall, func(page *compute.ReservationAggregatedList) error {

		for _, item := range page.Items {
			for _, inst := range item.Reservations {
				zonePathSplit := strings.Split(inst.Zone, "/")
				inst.Zone = zonePathSplit[len(zonePathSplit)-1]
			}
			itemsContainer.Items = append(itemsContainer.Items, item.Reservations...)
		}

		return nil
	}); err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_compute_reservation",
			"projectId": projectID,
			"errString": err.Error(),
		}).Error("failed to get aggregate list page")
		return resultMap, nil
	}

	byteArr, err := json.Marshal(itemsContainer)
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_compute_reservation",
			"errString": err.Error(),
		}).Error("failed to marshal response")
		return resultMap, err
	}
	tableConfig, ok := utilities.TableConfigurationMap["gcp_compute_reservation"]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_compute_reservation",
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found for \"gcp_compute_reservation\"")
	}
	jsonTable := utilities.NewTable(byteArr, tableConfig)
	for _, row := range jsonTable.Rows {
		result := extgcp.RowToMap(row, projectID, "", tableConfig)
		resultMap = append(resultMap, result)
	}

	return resultMap, nil
}
