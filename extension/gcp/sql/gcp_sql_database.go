package sql

import (
	"context"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"

	extgcp "github.com/Uptycs/cloudquery/extension/gcp"
	"github.com/Uptycs/cloudquery/utilities"
	"github.com/kolide/osquery-go/plugin/table"

	"google.golang.org/api/option"

	gcpsql "google.golang.org/api/sqladmin/v1beta4"
)

type myGcpSQLDatabasesItemsContainer struct {
	Items []*gcpsql.Database `json:"items"`
}

// GcpSQLDatabasesColumns returns the list of columns for gcp_sql_database
func GcpSQLDatabasesColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("project_id"),
		table.TextColumn("charset"),
		table.TextColumn("collation"),
		table.TextColumn("etag"),
		table.TextColumn("instance"),
		table.TextColumn("kind"),
		table.TextColumn("name"),
		table.TextColumn("project"),
		//table.TextColumn("self_link"),
		table.TextColumn("sqlserver_database_details"),
		//table.BigIntColumn("sqlserver_database_details_compatibility_level"),
		//table.TextColumn("sqlserver_database_details_recovery_model"),

	}
}

// GcpSQLDatabasesGenerate returns the rows in the table for all configured accounts
func GcpSQLDatabasesGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	var _ = queryContext
	ctx, cancel := context.WithCancel(osqCtx)
	defer cancel()

	resultMap := make([]map[string]string, 0)

	if len(utilities.ExtConfiguration.ExtConfGcp.Accounts) == 0 {
		results, err := processAccountGcpSQLDatabases(ctx, nil)
		if err == nil {
			resultMap = append(resultMap, results...)
		}
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfGcp.Accounts {
			results, err := processAccountGcpSQLDatabases(ctx, &account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}
	return resultMap, nil
}

func getGcpSQLDatabasesNewServiceForAccount(ctx context.Context, account *utilities.ExtensionConfigurationGcpAccount) (*gcpsql.Service, string) {
	var projectID = ""
	var service *gcpsql.Service
	var err error
	if account != nil {
		projectID = account.ProjectID
		service, err = gcpsql.NewService(ctx, option.WithCredentialsFile(account.KeyFile))
	} else {
		projectID = utilities.DefaultGcpProjectID
		service, err = gcpsql.NewService(ctx)
	}
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_sql_database",
			"projectId": projectID,
			"errString": err.Error(),
		}).Error("failed to create service")
		return nil, ""
	}
	return service, projectID
}

func getGcpSQLDatabasesKeys(service *gcpsql.Service, projectID string) ([]string, error) {
	resultList := make([]string, 0)
	listCall := service.Instances.List(projectID)
	if listCall == nil {
		utilities.GetLogger().WithFields(log.Fields{
			"projectId": projectID,
		}).Debug("Instances.List returned nil")
		return resultList, nil
	}

	rsp, doErr := listCall.Do()
	if doErr != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"projectId": projectID,
			"errString": doErr.Error(),
		}).Error("failed to get list.Do")
		return resultList, nil
	}

	for _, inst := range rsp.Items {
		resultList = append(resultList, inst.Name)
	}
	return resultList, nil
}

func processAccountGcpSQLDatabases(ctx context.Context,
	account *utilities.ExtensionConfigurationGcpAccount) ([]map[string]string, error) {

	resultMap := make([]map[string]string, 0)

	service, projectID := getGcpSQLDatabasesNewServiceForAccount(ctx, account)
	if service == nil {
		return resultMap, fmt.Errorf("failed to initialize gcpsql.Service")
	}

	keys, err := getGcpSQLDatabasesKeys(service, projectID)
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_sql_database",
			"projectId": projectID,
			"errString": err.Error(),
		}).Error("failed to get keys for Databases")
		return resultMap, nil
	}

	for _, key := range keys {
		getGcpSQLDatabasesRowsForKey(ctx, resultMap, service, projectID, key)
	}

	return resultMap, nil
}

func getGcpSQLDatabasesRowsForKey(ctx context.Context, resultMap []map[string]string, service *gcpsql.Service, projectID string, key string) ([]map[string]string, error) {
	listCall := service.Databases.List(projectID, key)
	if listCall == nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_sql_database",
			"projectId": projectID,
			"key":       key,
		}).Debug("list call is nil")
		return resultMap, nil
	}
	rsp, err := listCall.Do()
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_sql_database",
			"projectId": projectID,
			"key":       key,
			"errString": err.Error(),
		}).Error("failed to List.Do()")
		return resultMap, nil
	}
	itemsContainer := myGcpSQLDatabasesItemsContainer{Items: rsp.Items}
	byteArr, err := json.Marshal(itemsContainer)
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_sql_database",
			"errString": err.Error(),
		}).Error("failed to marshal response")
		return resultMap, err
	}
	tableConfig, ok := utilities.TableConfigurationMap["gcp_sql_database"]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_sql_database",
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found for \"gcp_sql_database\"")
	}
	jsonTable := utilities.NewTable(byteArr, tableConfig)
	for _, row := range jsonTable.Rows {
		result := extgcp.RowToMap(row, projectID, "", tableConfig)
		resultMap = append(resultMap, result)
	}

	return resultMap, nil
}
