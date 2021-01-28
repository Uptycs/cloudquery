/**
 * Copyright (c) 2020-present, The cloudquery authors
 *
 * This source code is licensed as defined by the LICENSE file found in the
 * root directory of this source tree.
 *
 * SPDX-License-Identifier: (Apache-2.0 OR GPL-2.0-only)
 */

package storage

import (
	"context"
	"encoding/json"
	"fmt"
	extgcp "github.com/Uptycs/cloudquery/extension/gcp"
	"github.com/Uptycs/cloudquery/utilities"
	log "github.com/sirupsen/logrus"

	"github.com/kolide/osquery-go/plugin/table"
	"google.golang.org/api/option"

	storage "cloud.google.com/go/storage"
)

type myGcpStorageBucketItemsContainer struct {
	Items []*storage.BucketAttrs `json:"items"`
}

// GcpStorageBucketColumns returns the list of columns for gcp_storage_bucket
func (handler *GcpStorageHandler) GcpStorageBucketColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("project_id"),
		table.TextColumn("acl"),
		//table.TextColumn("acl_domain"),
		//table.TextColumn("acl_email"),
		//table.TextColumn("acl_entity"),
		//table.TextColumn("acl_entity_id"),
		//table.TextColumn("acl_project_team"),
		//table.TextColumn("acl_project_team_project_number"),
		//table.TextColumn("acl_project_team_team"),
		//table.TextColumn("acl_role"),
		table.TextColumn("bucket_policy_only"),
		//table.TextColumn("bucket_policy_only_enabled"),
		//table.TextColumn("bucket_policy_only_locked_time"),
		//table.BigIntColumn("bucket_policy_only_locked_time_ext"),
		//table.TextColumn("bucket_policy_only_locked_time_loc"),
		//table.BigIntColumn("bucket_policy_only_locked_time_loc_cache_end"),
		//table.BigIntColumn("bucket_policy_only_locked_time_loc_cache_start"),
		//table.TextColumn("bucket_policy_only_locked_time_loc_cache_zone"),
		//table.TextColumn("bucket_policy_only_locked_time_loc_cache_zone_is_dst"),
		//table.TextColumn("bucket_policy_only_locked_time_loc_cache_zone_name"),
		//table.IntegerColumn("bucket_policy_only_locked_time_loc_cache_zone_offset"),
		//table.TextColumn("bucket_policy_only_locked_time_loc_name"),
		//table.TextColumn("bucket_policy_only_locked_time_loc_tx"),
		//table.IntegerColumn("bucket_policy_only_locked_time_loc_tx_index"),
		//table.TextColumn("bucket_policy_only_locked_time_loc_tx_isstd"),
		//table.TextColumn("bucket_policy_only_locked_time_loc_tx_isutc"),
		//table.BigIntColumn("bucket_policy_only_locked_time_loc_tx_when"),
		//table.TextColumn("bucket_policy_only_locked_time_loc_zone"),
		//table.TextColumn("bucket_policy_only_locked_time_loc_zone_is_dst"),
		//table.TextColumn("bucket_policy_only_locked_time_loc_zone_name"),
		//table.IntegerColumn("bucket_policy_only_locked_time_loc_zone_offset"),
		//table.BigIntColumn("bucket_policy_only_locked_time_wall"),
		table.TextColumn("cors"),
		//table.BigIntColumn("cors_max_age"),
		//table.TextColumn("cors_methods"),
		//table.TextColumn("cors_origins"),
		//table.TextColumn("cors_response_headers"),
		table.TextColumn("created"),
		//table.BigIntColumn("created_ext"),
		//table.TextColumn("created_loc"),
		//table.BigIntColumn("created_loc_cache_end"),
		//table.BigIntColumn("created_loc_cache_start"),
		//table.TextColumn("created_loc_cache_zone"),
		//table.TextColumn("created_loc_cache_zone_is_dst"),
		//table.TextColumn("created_loc_cache_zone_name"),
		//table.IntegerColumn("created_loc_cache_zone_offset"),
		//table.TextColumn("created_loc_name"),
		//table.TextColumn("created_loc_tx"),
		//table.IntegerColumn("created_loc_tx_index"),
		//table.TextColumn("created_loc_tx_isstd"),
		//table.TextColumn("created_loc_tx_isutc"),
		//table.BigIntColumn("created_loc_tx_when"),
		//table.TextColumn("created_loc_zone"),
		//table.TextColumn("created_loc_zone_is_dst"),
		//table.TextColumn("created_loc_zone_name"),
		//table.IntegerColumn("created_loc_zone_offset"),
		//table.BigIntColumn("created_wall"),
		table.TextColumn("default_event_based_hold"),
		table.TextColumn("default_object_acl"),
		//table.TextColumn("default_object_acl_domain"),
		//table.TextColumn("default_object_acl_email"),
		//table.TextColumn("default_object_acl_entity"),
		//table.TextColumn("default_object_acl_entity_id"),
		//table.TextColumn("default_object_acl_project_team"),
		//table.TextColumn("default_object_acl_project_team_project_number"),
		//table.TextColumn("default_object_acl_project_team_team"),
		//table.TextColumn("default_object_acl_role"),
		table.TextColumn("encryption"),
		//table.TextColumn("encryption_default_kms_key_name"),
		table.TextColumn("etag"),
		table.TextColumn("labels"),
		table.TextColumn("lifecycle"),
		//table.TextColumn("lifecycle_rules"),
		//table.TextColumn("lifecycle_rules_action"),
		//table.TextColumn("lifecycle_rules_action_storage_class"),
		//table.TextColumn("lifecycle_rules_action_type"),
		//table.TextColumn("lifecycle_rules_condition"),
		//table.BigIntColumn("lifecycle_rules_condition_age_in_days"),
		//table.TextColumn("lifecycle_rules_condition_created_before"),
		//table.BigIntColumn("lifecycle_rules_condition_created_before_ext"),
		//table.TextColumn("lifecycle_rules_condition_created_before_loc"),
		//table.BigIntColumn("lifecycle_rules_condition_created_before_loc_cache_end"),
		//table.BigIntColumn("lifecycle_rules_condition_created_before_loc_cache_start"),
		//table.TextColumn("lifecycle_rules_condition_created_before_loc_cache_zone"),
		//table.TextColumn("lifecycle_rules_condition_created_before_loc_cache_zone_is_dst"),
		//table.TextColumn("lifecycle_rules_condition_created_before_loc_cache_zone_name"),
		//table.IntegerColumn("lifecycle_rules_condition_created_before_loc_cache_zone_offset"),
		//table.TextColumn("lifecycle_rules_condition_created_before_loc_name"),
		//table.TextColumn("lifecycle_rules_condition_created_before_loc_tx"),
		//table.IntegerColumn("lifecycle_rules_condition_created_before_loc_tx_index"),
		//table.TextColumn("lifecycle_rules_condition_created_before_loc_tx_isstd"),
		//table.TextColumn("lifecycle_rules_condition_created_before_loc_tx_isutc"),
		//table.BigIntColumn("lifecycle_rules_condition_created_before_loc_tx_when"),
		//table.TextColumn("lifecycle_rules_condition_created_before_loc_zone"),
		//table.TextColumn("lifecycle_rules_condition_created_before_loc_zone_is_dst"),
		//table.TextColumn("lifecycle_rules_condition_created_before_loc_zone_name"),
		//table.IntegerColumn("lifecycle_rules_condition_created_before_loc_zone_offset"),
		//table.BigIntColumn("lifecycle_rules_condition_created_before_wall"),
		//table.IntegerColumn("lifecycle_rules_condition_liveness"),
		//table.TextColumn("lifecycle_rules_condition_matches_storage_classes"),
		//table.BigIntColumn("lifecycle_rules_condition_num_newer_versions"),
		table.TextColumn("location"),
		table.TextColumn("location_type"),
		table.TextColumn("logging"),
		//table.TextColumn("logging_log_bucket"),
		//table.TextColumn("logging_log_object_prefix"),
		table.BigIntColumn("meta_generation"),
		table.TextColumn("name"),
		table.TextColumn("predefined_acl"),
		table.TextColumn("predefined_default_object_acl"),
		table.TextColumn("requester_pays"),
		table.TextColumn("retention_policy"),
		//table.TextColumn("retention_policy_effective_time"),
		//table.BigIntColumn("retention_policy_effective_time_ext"),
		//table.TextColumn("retention_policy_effective_time_loc"),
		//table.BigIntColumn("retention_policy_effective_time_loc_cache_end"),
		//table.BigIntColumn("retention_policy_effective_time_loc_cache_start"),
		//table.TextColumn("retention_policy_effective_time_loc_cache_zone"),
		//table.TextColumn("retention_policy_effective_time_loc_cache_zone_is_dst"),
		//table.TextColumn("retention_policy_effective_time_loc_cache_zone_name"),
		//table.IntegerColumn("retention_policy_effective_time_loc_cache_zone_offset"),
		//table.TextColumn("retention_policy_effective_time_loc_name"),
		//table.TextColumn("retention_policy_effective_time_loc_tx"),
		//table.IntegerColumn("retention_policy_effective_time_loc_tx_index"),
		//table.TextColumn("retention_policy_effective_time_loc_tx_isstd"),
		//table.TextColumn("retention_policy_effective_time_loc_tx_isutc"),
		//table.BigIntColumn("retention_policy_effective_time_loc_tx_when"),
		//table.TextColumn("retention_policy_effective_time_loc_zone"),
		//table.TextColumn("retention_policy_effective_time_loc_zone_is_dst"),
		//table.TextColumn("retention_policy_effective_time_loc_zone_name"),
		//table.IntegerColumn("retention_policy_effective_time_loc_zone_offset"),
		//table.BigIntColumn("retention_policy_effective_time_wall"),
		//table.TextColumn("retention_policy_is_locked"),
		//table.BigIntColumn("retention_policy_retention_period"),
		table.TextColumn("storage_class"),
		table.TextColumn("uniform_bucket_level_access"),
		//table.TextColumn("uniform_bucket_level_access_enabled"),
		//table.TextColumn("uniform_bucket_level_access_locked_time"),
		//table.BigIntColumn("uniform_bucket_level_access_locked_time_ext"),
		//table.TextColumn("uniform_bucket_level_access_locked_time_loc"),
		//table.BigIntColumn("uniform_bucket_level_access_locked_time_loc_cache_end"),
		//table.BigIntColumn("uniform_bucket_level_access_locked_time_loc_cache_start"),
		//table.TextColumn("uniform_bucket_level_access_locked_time_loc_cache_zone"),
		//table.TextColumn("uniform_bucket_level_access_locked_time_loc_cache_zone_is_dst"),
		//table.TextColumn("uniform_bucket_level_access_locked_time_loc_cache_zone_name"),
		//table.IntegerColumn("uniform_bucket_level_access_locked_time_loc_cache_zone_offset"),
		//table.TextColumn("uniform_bucket_level_access_locked_time_loc_name"),
		//table.TextColumn("uniform_bucket_level_access_locked_time_loc_tx"),
		//table.IntegerColumn("uniform_bucket_level_access_locked_time_loc_tx_index"),
		//table.TextColumn("uniform_bucket_level_access_locked_time_loc_tx_isstd"),
		//table.TextColumn("uniform_bucket_level_access_locked_time_loc_tx_isutc"),
		//table.BigIntColumn("uniform_bucket_level_access_locked_time_loc_tx_when"),
		//table.TextColumn("uniform_bucket_level_access_locked_time_loc_zone"),
		//table.TextColumn("uniform_bucket_level_access_locked_time_loc_zone_is_dst"),
		//table.TextColumn("uniform_bucket_level_access_locked_time_loc_zone_name"),
		//table.IntegerColumn("uniform_bucket_level_access_locked_time_loc_zone_offset"),
		//table.BigIntColumn("uniform_bucket_level_access_locked_time_wall"),
		table.TextColumn("versioning_enabled"),
		table.TextColumn("website"),
		//table.TextColumn("website_main_page_suffix"),
		//table.TextColumn("website_not_found_page"),

	}
}

// GcpStorageBucketGenerate returns the rows in the table for all configured accounts
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
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}
	return resultMap, nil
}

func (handler *GcpStorageHandler) getGcpStorageBucketNewServiceForAccount(ctx context.Context, account *utilities.ExtensionConfigurationGcpAccount) (*storage.Client, string) {
	var projectID string
	var service *storage.Client
	var err error
	if account != nil {
		projectID = account.ProjectID
		service, err = handler.svcInterface.NewClient(ctx, option.WithCredentialsFile(account.KeyFile))
	} else {
		projectID = utilities.DefaultGcpProjectID
		service, err = handler.svcInterface.NewClient(ctx)
	}
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_storage_bucket",
			"projectId": projectID,
			"errString": err.Error(),
		}).Error("failed to create client")
		return nil, ""
	}
	return service, projectID
}

func (handler *GcpStorageHandler) processAccountGcpStorageBucket(ctx context.Context,
	account *utilities.ExtensionConfigurationGcpAccount) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)

	tableConfig, ok := utilities.TableConfigurationMap["gcp_storage_bucket"]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_storage_bucket",
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found for \"gcp_storage_bucket\"")
	}

	service, projectID := handler.getGcpStorageBucketNewServiceForAccount(ctx, account)
	if service == nil {
		return resultMap, fmt.Errorf("failed to initialize storage.Client")
	}
	listCall := handler.svcInterface.Buckets(ctx, service, projectID)

	if listCall == nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "gcp_storage_bucket",
			"projectId": projectID,
		}).Debug("listCall is nil")
		return resultMap, nil
	}
	p := handler.svcInterface.BucketsNewPager(listCall, 10, "")
	for {
		var container = myGcpStorageBucketItemsContainer{}
		pageToken, err := p.NextPage(&container.Items)
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "gcp_storage_bucket",
				"projectId": projectID,
				"errString": err.Error(),
			}).Error("failed to get next page")
			return resultMap, err
		}

		byteArr, err := json.Marshal(&container)
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "gcp_storage_bucket",
				"errString": err.Error(),
			}).Error("failed to marshal response")
		}
		jsonTable := utilities.NewTable(byteArr, tableConfig)
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
