package storage

import (
	"context"
	"encoding/json"
	"fmt"
	extgcp "github.com/Uptycs/cloudquery/extension/gcp"
	"github.com/Uptycs/cloudquery/utilities"
	"os"

	"github.com/kolide/osquery-go/plugin/table"
	"google.golang.org/api/option"

	storage "cloud.google.com/go/storage"
)

type ItemsContainer struct {
	Items []*storage.BucketAttrs `json:"items"`
}

func (handler *GcpStorageHandler) GcpStorageBucketColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("project_id"),
		table.TextColumn("name"),
		table.TextColumn("acl"),
		table.TextColumn("bucket_policy_only"),
		table.TextColumn("uniform_bucket_level_access"),
		table.TextColumn("default_object_acl"),
		table.TextColumn("default_event_based_hold"),
		table.TextColumn("predefined_acl"),
		table.TextColumn("predefined_default_object_acl"),
		table.TextColumn("location"),
		table.BigIntColumn("meta_generation"),
		table.TextColumn("storage_class"),
		table.TextColumn("created"),
		table.TextColumn("versioning_enabled"),
		table.TextColumn("labels"),
		table.TextColumn("requester_pays"),
		table.TextColumn("lifecycle"),
		//table.TextColumn("retention_policy"),
		table.BigIntColumn("retention_policy_retention_period"),
		table.TextColumn("retention_policy_effective_time"),
		table.TextColumn("retention_policy_is_locked"),
		table.TextColumn("cors"),
		//table.TextColumn("encryption"),
		//table.TextColumn("encryption_default_kms_key_name"),
		//table.TextColumn("logging"),
		table.TextColumn("logging_log_bucket"),
		table.TextColumn("logging_log_object_prefix"),
		//table.TextColumn("website"),
		table.TextColumn("website_main_page_suffix"),
		table.TextColumn("website_not_found_page"),
		table.TextColumn("etag"),
		table.TextColumn("location_type"),
	}
}

func (handler *GcpStorageHandler) GcpStorageBucketGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	var _ = queryContext
	ctx, cancel := context.WithCancel(osqCtx)
	defer cancel()

	resultMap := make([]map[string]string, 0)

	if len(utilities.ExtConfiguration.ExtConfGcp.Accounts) == 0 {
		results, err := handler.processAccountGcpStorageBucket(ctx, nil)
		if err == nil {
			resultMap = append(resultMap, results...)
		}
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfGcp.Accounts {
			results, err := handler.processAccountGcpStorageBucket(ctx, &account)
			if err != nil {
				// TODO: Continue to next account or return error ?
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}
	return resultMap, nil
}

func (handler *GcpStorageHandler) getGcpStorageBucketNewServiceForAccount(ctx context.Context, account *utilities.ExtensionConfigurationGcpAccount) (*storage.Client, string) {
	var projectID = ""
	var service *storage.Client
	var err error
	if account != nil {
		projectID = account.ProjectId
		service, err = handler.svcInterface.NewClient(ctx, option.WithCredentialsFile(account.KeyFile))
	} else {
		projectID = utilities.DefaultGcpProjectID
		service, err = handler.svcInterface.NewClient(ctx)
	}
	if err != nil {
		fmt.Println("NewClient() error: ", err)
		return nil, ""
	}
	return service, projectID
}

func (handler *GcpStorageHandler) processAccountGcpStorageBucket(ctx context.Context,
	account *utilities.ExtensionConfigurationGcpAccount) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)

	tableConfig, ok := utilities.TableConfigurationMap["gcp_storage_bucket"]
	if !ok {
		fmt.Println("getTableConfig failed for gcp_storage_bucket")
		return resultMap, fmt.Errorf("table configuration not found")
	}

	service, projectID := handler.getGcpStorageBucketNewServiceForAccount(ctx, account)
	if service == nil {
		return resultMap, fmt.Errorf("failed to initialize storage.Client")
	}
	listCall := handler.svcInterface.Buckets(ctx, service, projectID)

	if listCall == nil {
		fmt.Println("listCall is nil")
		return resultMap, nil
	}
	p := handler.svcInterface.BucketsNewPager(listCall, 10, "")
	for {
		var container = ItemsContainer{}
		pageToken, err := p.NextPage(&container.Items)
		if err != nil {
			fmt.Println("NextPage() error: ", err)
			return resultMap, err
		}

		byteArr, err := json.Marshal(&container)
		if err != nil {
			fmt.Printf("error: %v\n", err)
			os.Exit(1)
		}
		//fmt.Printf("%+v\n", string(byteArr))
		jsonTable := utilities.Table{}
		jsonTable.Init(byteArr, tableConfig.MaxLevel, tableConfig.GetParsedAttributeConfigMap())
		for _, row := range jsonTable.Rows {
			result := extgcp.RowToMap(row, projectID, "", tableConfig)
			resultMap = append(resultMap, result)
		}

		if pageToken == "" {
			break
		}
	}
	return resultMap, nil
}
