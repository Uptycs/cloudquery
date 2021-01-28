package file

import (
	"context"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"

	extgcp "github.com/Uptycs/cloudquery/extension/gcp"
	"github.com/Uptycs/cloudquery/utilities"
	"github.com/kolide/osquery-go/plugin/table"

	"google.golang.org/api/option"

	gcpfile "google.golang.org/api/file/v1beta1"
)

type myGcpFileInstancesItemsContainer struct {
	Items []*gcpfile.Instance `json:"items"`
}

// GcpFileInstancesColumns returns the list of columns for gcp_file_instance
func GcpFileInstancesColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("project_id"),
		table.TextColumn("create_time"),
		table.TextColumn("description"),
		table.TextColumn("etag"),
		table.TextColumn("file_shares"),
		//table.BigIntColumn("file_shares_capacity_gb"),
		//table.TextColumn("file_shares_name"),
		//table.TextColumn("file_shares_nfs_export_options"),
		//table.TextColumn("file_shares_nfs_export_options_access_mode"),
		//table.BigIntColumn("file_shares_nfs_export_options_anon_gid"),
		//table.BigIntColumn("file_shares_nfs_export_options_anon_uid"),
		//table.TextColumn("file_shares_nfs_export_options_ip_ranges"),
		//table.TextColumn("file_shares_nfs_export_options_squash_mode"),
		//table.TextColumn("file_shares_source_backup"),
		table.TextColumn("labels"),
		table.TextColumn("name"),
		table.TextColumn("networks"),
		//table.TextColumn("networks_ip_addresses"),
		//table.TextColumn("networks_modes"),
		//table.TextColumn("networks_network"),
		//table.TextColumn("networks_reserved_ip_range"),
		table.TextColumn("state"),
		table.TextColumn("status_message"),
		table.TextColumn("tier"),
	}
}

// GcpFileInstancesGenerate returns the rows in the table for all configured accounts
func GcpFileInstancesGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	var _ = queryContext
	ctx, cancel := context.WithCancel(osqCtx)
	defer cancel()

	resultMap := make([]map[string]string, 0)

	if len(utilities.ExtConfiguration.ExtConfGcp.Accounts) == 0 {
		results, err := processAccountGcpFileInstances(ctx, nil)
		if err == nil {
			resultMap = append(resultMap, results...)
		}
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfGcp.Accounts {
			results, err := processAccountGcpFileInstances(ctx, &account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}
	return resultMap, nil
}

func getGcpFileInstancesNewServiceForAccount(ctx context.Context, account *utilities.ExtensionConfigurationGcpAccount) (*gcpfile.Service, string) {
	var projectID = ""
	var service *gcpfile.Service
	var err error
	if account != nil {
		projectID = account.ProjectID
		service, err = gcpfile.NewService(ctx, option.WithCredentialsFile(account.KeyFile))
	} else {
		projectID = utilities.DefaultGcpProjectID
		service, err = gcpfile.NewService(ctx)
	}
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_file_instance",
			"projectId": projectID,
			"errString": err.Error(),
		}).Error("failed to create service")
		return nil, ""
	}
	return service, projectID
}

func processAccountGcpFileInstances(ctx context.Context,
	account *utilities.ExtensionConfigurationGcpAccount) ([]map[string]string, error) {

	resultMap := make([]map[string]string, 0)

	service, projectID := getGcpFileInstancesNewServiceForAccount(ctx, account)
	if service == nil {
		return resultMap, fmt.Errorf("failed to initialize gcpfile.Service")
	}

	listCall := service.Projects.Locations.Instances.List("projects/" + projectID + "/locations/-")
	if listCall == nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_file_instance",
			"projectId": projectID,
		}).Debug("list call is nil")
		return resultMap, nil
	}
	itemsContainer := myGcpFileInstancesItemsContainer{Items: make([]*gcpfile.Instance, 0)}
	if err := listCall.Pages(ctx, func(page *gcpfile.ListInstancesResponse) error {

		itemsContainer.Items = append(itemsContainer.Items, page.Instances...)

		return nil
	}); err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_file_instance",
			"projectId": projectID,
			"errString": err.Error(),
		}).Error("failed to get aggregate list page")
		return resultMap, nil
	}

	byteArr, err := json.Marshal(itemsContainer)
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_file_instance",
			"errString": err.Error(),
		}).Error("failed to marshal response")
		return resultMap, err
	}
	tableConfig, ok := utilities.TableConfigurationMap["gcp_file_instance"]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_file_instance",
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found for \"gcp_file_instance\"")
	}
	jsonTable := utilities.NewTable(byteArr, tableConfig)
	for _, row := range jsonTable.Rows {
		result := extgcp.RowToMap(row, projectID, "", tableConfig)
		resultMap = append(resultMap, result)
	}

	return resultMap, nil
}
