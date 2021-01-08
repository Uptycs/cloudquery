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
	iterator "google.golang.org/api/iterator"
)

type ItemsContainer struct {
	Items []*storage.BucketAttrs
}

func GcpStorageBucketColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("project_id"),
		table.TextColumn("kind"),
		table.TextColumn("id"),
		table.TextColumn("project_number"),
		table.TextColumn("name"),
		table.TextColumn("time_created"),
		table.TextColumn("updated"),
		table.TextColumn("default_event_based_hold"),
		table.TextColumn("retention_policy_retention_period"),
		table.TextColumn("retention_policy_effective_time"),
		table.TextColumn("retention_policy_is_locked"),
		table.TextColumn("metageneration"),
		table.TextColumn("iam_configuration_uniform_bucket_level_access_enabled"),
		table.TextColumn("iam_configuration_uniform_bucket_level_access_locked_time"),
		table.TextColumn("encryption_default_kms_key_name"),
		table.TextColumn("owner_entity"),
		table.TextColumn("owner_entity_id"),
		table.TextColumn("location"),
		table.TextColumn("location_type"),
		table.TextColumn("website_main_page_suffix"),
		table.TextColumn("website_not_found_page"),
		table.TextColumn("logging_log_bucket"),
		table.TextColumn("logging_log_object_prefix"),
		table.TextColumn("versioning_enabled"),
		table.TextColumn("labels"),
		table.TextColumn("storage_class"),
		table.TextColumn("billing_requester_pays"),
		table.TextColumn("etag"),
	}
}

func GcpStorageBucketGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	var _ = queryContext
	ctx, cancel := context.WithCancel(osqCtx)
	defer cancel()

	resultMap := make([]map[string]string, 0)

	for _, account := range utilities.ExtConfiguration.ExtConfGcp.Accounts {
		results, err := processAccountGcpStorageBucket(ctx, &account)
		if err != nil {
			// TODO: Continue to next account or return error ?
			continue
		}
		resultMap = append(resultMap, results...)
	}
	return resultMap, nil
}

// TODO: Remove this
var sample1 = storage.BucketAttrs{Name: "Test1", Location: "TestLocation1"}
var sample2 = storage.BucketAttrs{Name: "Test2", Location: "TestLocation2"}

func processAccountGcpStorageBucket(ctx context.Context,
	account *utilities.ExtensionConfigurationGcpAccount) ([]map[string]string, error) {

	resultMap := make([]map[string]string, 0)

	service, err := storage.NewClient(ctx, option.WithCredentialsFile(account.KeyFile))
	if err != nil {
		fmt.Println("storage.NewClient() error: ", err)
		return resultMap, err
	}

	tableConfig, ok := utilities.TableConfigurationMap["gcp_storage_bucket"]
	if !ok {
		fmt.Println("getTableConfig: ", err)
		return resultMap, fmt.Errorf("table configuration not found")
	}

	listCall := service.Buckets(ctx, account.ProjectId)

	if listCall == nil {
		fmt.Println("listCall is nil")
		return resultMap, nil
	}
	p := iterator.NewPager(listCall, 10, "") // TODO:: fix me
	for {
		var container = ItemsContainer{}
		pageToken, err := p.NextPage(&container.Items)
		if err != nil {
			fmt.Println("NextPage() error: ", err)
			return resultMap, err
		}
		if len(container.Items) == 0 { // TODO: remove it
			container.Items = append(container.Items, &sample1)
			container.Items = append(container.Items, &sample2)
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
			result := extgcp.RowToMap(row, account.ProjectId, "", tableConfig)
			resultMap = append(resultMap, result)
		}

		if pageToken == "" {
			break
		}
	}
	return resultMap, nil
}
