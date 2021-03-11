/**
 * Copyright (c) 2020-present, The cloudquery authors
 *
 * This source code is licensed as defined by the LICENSE file found in the
 * root directory of this source tree.
 *
 * SPDX-License-Identifier: (Apache-2.0 OR GPL-2.0-only)
 */

package cloudfront

import (
	"context"
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/Uptycs/cloudquery/utilities"

	"github.com/Uptycs/basequery-go/plugin/table"
	extaws "github.com/Uptycs/cloudquery/extension/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
)

// ListDistributionsColumns returns the list of columns in the table
func ListDistributionsColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("account_id"),
		table.TextColumn("is_truncated"),
		//table.TextColumn("items"),
		table.TextColumn("items_arn"),
		table.TextColumn("items_alias_icp_recordals"),
		//table.TextColumn("items_alias_icp_recordals_cname"),
		//table.TextColumn("items_alias_icp_recordals_icp_recordal_status"),
		table.TextColumn("items_aliases"),
		//table.TextColumn("items_aliases_items"),
		//table.IntegerColumn("items_aliases_quantity"),
		table.TextColumn("items_cache_behaviors"),
		//table.TextColumn("items_cache_behaviors_items"),
		//table.TextColumn("items_cache_behaviors_items_allowed_methods"),
		//table.TextColumn("items_cache_behaviors_items_allowed_methods_cached_methods"),
		//table.TextColumn("items_cache_behaviors_items_allowed_methods_cached_methods_items"),
		//table.IntegerColumn("items_cache_behaviors_items_allowed_methods_cached_methods_quantity"),
		//table.TextColumn("items_cache_behaviors_items_allowed_methods_items"),
		//table.IntegerColumn("items_cache_behaviors_items_allowed_methods_quantity"),
		//table.TextColumn("items_cache_behaviors_items_cache_policy_id"),
		//table.TextColumn("items_cache_behaviors_items_compress"),
		//table.BigIntColumn("items_cache_behaviors_items_default_ttl"),
		//table.TextColumn("items_cache_behaviors_items_field_level_encryption_id"),
		//table.TextColumn("items_cache_behaviors_items_forwarded_values"),
		//table.TextColumn("items_cache_behaviors_items_forwarded_values_cookies"),
		//table.TextColumn("items_cache_behaviors_items_forwarded_values_cookies_forward"),
		//table.TextColumn("items_cache_behaviors_items_forwarded_values_cookies_whitelisted_names"),
		//table.TextColumn("items_cache_behaviors_items_forwarded_values_cookies_whitelisted_names_items"),
		//table.IntegerColumn("items_cache_behaviors_items_forwarded_values_cookies_whitelisted_names_quantity"),
		//table.TextColumn("items_cache_behaviors_items_forwarded_values_headers"),
		//table.TextColumn("items_cache_behaviors_items_forwarded_values_headers_items"),
		//table.IntegerColumn("items_cache_behaviors_items_forwarded_values_headers_quantity"),
		//table.TextColumn("items_cache_behaviors_items_forwarded_values_query_string"),
		//table.TextColumn("items_cache_behaviors_items_forwarded_values_query_string_cache_keys"),
		//table.TextColumn("items_cache_behaviors_items_forwarded_values_query_string_cache_keys_items"),
		//table.IntegerColumn("items_cache_behaviors_items_forwarded_values_query_string_cache_keys_quantity"),
		//table.TextColumn("items_cache_behaviors_items_lambda_function_associations"),
		//table.TextColumn("items_cache_behaviors_items_lambda_function_associations_items"),
		//table.TextColumn("items_cache_behaviors_items_lambda_function_associations_items_event_type"),
		//table.TextColumn("items_cache_behaviors_items_lambda_function_associations_items_include_body"),
		//table.TextColumn("items_cache_behaviors_items_lambda_function_associations_items_lambda_function_arn"),
		//table.IntegerColumn("items_cache_behaviors_items_lambda_function_associations_quantity"),
		//table.BigIntColumn("items_cache_behaviors_items_max_ttl"),
		//table.BigIntColumn("items_cache_behaviors_items_min_ttl"),
		//table.TextColumn("items_cache_behaviors_items_origin_request_policy_id"),
		//table.TextColumn("items_cache_behaviors_items_path_pattern"),
		//table.TextColumn("items_cache_behaviors_items_realtime_log_config_arn"),
		//table.TextColumn("items_cache_behaviors_items_smooth_streaming"),
		//table.TextColumn("items_cache_behaviors_items_target_origin_id"),
		//table.TextColumn("items_cache_behaviors_items_trusted_key_groups"),
		//table.TextColumn("items_cache_behaviors_items_trusted_key_groups_enabled"),
		//table.TextColumn("items_cache_behaviors_items_trusted_key_groups_items"),
		//table.IntegerColumn("items_cache_behaviors_items_trusted_key_groups_quantity"),
		//table.TextColumn("items_cache_behaviors_items_trusted_signers"),
		//table.TextColumn("items_cache_behaviors_items_trusted_signers_enabled"),
		//table.TextColumn("items_cache_behaviors_items_trusted_signers_items"),
		//table.IntegerColumn("items_cache_behaviors_items_trusted_signers_quantity"),
		//table.TextColumn("items_cache_behaviors_items_viewer_protocol_policy"),
		//table.IntegerColumn("items_cache_behaviors_quantity"),
		table.TextColumn("items_comment"),
		table.TextColumn("items_custom_error_responses"),
		//table.TextColumn("items_custom_error_responses_items"),
		//table.BigIntColumn("items_custom_error_responses_items_error_caching_min_ttl"),
		//table.IntegerColumn("items_custom_error_responses_items_error_code"),
		//table.TextColumn("items_custom_error_responses_items_response_code"),
		//table.TextColumn("items_custom_error_responses_items_response_page_path"),
		//table.IntegerColumn("items_custom_error_responses_quantity"),
		table.TextColumn("items_default_cache_behavior"),
		//table.TextColumn("items_default_cache_behavior_allowed_methods"),
		//table.TextColumn("items_default_cache_behavior_allowed_methods_cached_methods"),
		//table.TextColumn("items_default_cache_behavior_allowed_methods_cached_methods_items"),
		//table.IntegerColumn("items_default_cache_behavior_allowed_methods_cached_methods_quantity"),
		//table.TextColumn("items_default_cache_behavior_allowed_methods_items"),
		//table.IntegerColumn("items_default_cache_behavior_allowed_methods_quantity"),
		//table.TextColumn("items_default_cache_behavior_cache_policy_id"),
		//table.TextColumn("items_default_cache_behavior_compress"),
		//table.BigIntColumn("items_default_cache_behavior_default_ttl"),
		//table.TextColumn("items_default_cache_behavior_field_level_encryption_id"),
		//table.TextColumn("items_default_cache_behavior_forwarded_values"),
		//table.TextColumn("items_default_cache_behavior_forwarded_values_cookies"),
		//table.TextColumn("items_default_cache_behavior_forwarded_values_cookies_forward"),
		//table.TextColumn("items_default_cache_behavior_forwarded_values_cookies_whitelisted_names"),
		//table.TextColumn("items_default_cache_behavior_forwarded_values_cookies_whitelisted_names_items"),
		//table.IntegerColumn("items_default_cache_behavior_forwarded_values_cookies_whitelisted_names_quantity"),
		//table.TextColumn("items_default_cache_behavior_forwarded_values_headers"),
		//table.TextColumn("items_default_cache_behavior_forwarded_values_headers_items"),
		//table.IntegerColumn("items_default_cache_behavior_forwarded_values_headers_quantity"),
		//table.TextColumn("items_default_cache_behavior_forwarded_values_query_string"),
		//table.TextColumn("items_default_cache_behavior_forwarded_values_query_string_cache_keys"),
		//table.TextColumn("items_default_cache_behavior_forwarded_values_query_string_cache_keys_items"),
		//table.IntegerColumn("items_default_cache_behavior_forwarded_values_query_string_cache_keys_quantity"),
		//table.TextColumn("items_default_cache_behavior_lambda_function_associations"),
		//table.TextColumn("items_default_cache_behavior_lambda_function_associations_items"),
		//table.TextColumn("items_default_cache_behavior_lambda_function_associations_items_event_type"),
		//table.TextColumn("items_default_cache_behavior_lambda_function_associations_items_include_body"),
		//table.TextColumn("items_default_cache_behavior_lambda_function_associations_items_lambda_function_arn"),
		//table.IntegerColumn("items_default_cache_behavior_lambda_function_associations_quantity"),
		//table.BigIntColumn("items_default_cache_behavior_max_ttl"),
		//table.BigIntColumn("items_default_cache_behavior_min_ttl"),
		//table.TextColumn("items_default_cache_behavior_origin_request_policy_id"),
		//table.TextColumn("items_default_cache_behavior_realtime_log_config_arn"),
		//table.TextColumn("items_default_cache_behavior_smooth_streaming"),
		//table.TextColumn("items_default_cache_behavior_target_origin_id"),
		//table.TextColumn("items_default_cache_behavior_trusted_key_groups"),
		//table.TextColumn("items_default_cache_behavior_trusted_key_groups_enabled"),
		//table.TextColumn("items_default_cache_behavior_trusted_key_groups_items"),
		//table.IntegerColumn("items_default_cache_behavior_trusted_key_groups_quantity"),
		//table.TextColumn("items_default_cache_behavior_trusted_signers"),
		//table.TextColumn("items_default_cache_behavior_trusted_signers_enabled"),
		//table.TextColumn("items_default_cache_behavior_trusted_signers_items"),
		//table.IntegerColumn("items_default_cache_behavior_trusted_signers_quantity"),
		//table.TextColumn("items_default_cache_behavior_viewer_protocol_policy"),
		table.TextColumn("items_domain_name"),
		table.TextColumn("items_enabled"),
		table.TextColumn("items_http_version"),
		table.TextColumn("items_id"),
		table.TextColumn("items_is_ipv6_enabled"),
		table.TextColumn("items_last_modified_time"),
		//table.BigIntColumn("items_last_modified_time_ext"),
		//table.TextColumn("items_last_modified_time_loc"),
		//table.BigIntColumn("items_last_modified_time_loc_cache_end"),
		//table.BigIntColumn("items_last_modified_time_loc_cache_start"),
		//table.TextColumn("items_last_modified_time_loc_cache_zone"),
		//table.TextColumn("items_last_modified_time_loc_cache_zone_is_dst"),
		//table.TextColumn("items_last_modified_time_loc_cache_zone_name"),
		//table.IntegerColumn("items_last_modified_time_loc_cache_zone_offset"),
		//table.TextColumn("items_last_modified_time_loc_extend"),
		//table.TextColumn("items_last_modified_time_loc_name"),
		//table.TextColumn("items_last_modified_time_loc_tx"),
		//table.IntegerColumn("items_last_modified_time_loc_tx_index"),
		//table.TextColumn("items_last_modified_time_loc_tx_isstd"),
		//table.TextColumn("items_last_modified_time_loc_tx_isutc"),
		//table.BigIntColumn("items_last_modified_time_loc_tx_when"),
		//table.TextColumn("items_last_modified_time_loc_zone"),
		//table.TextColumn("items_last_modified_time_loc_zone_is_dst"),
		//table.TextColumn("items_last_modified_time_loc_zone_name"),
		//table.IntegerColumn("items_last_modified_time_loc_zone_offset"),
		//table.BigIntColumn("items_last_modified_time_wall"),
		table.TextColumn("items_origin_groups"),
		//table.TextColumn("items_origin_groups_items"),
		//table.TextColumn("items_origin_groups_items_failover_criteria"),
		//table.TextColumn("items_origin_groups_items_failover_criteria_status_codes"),
		//table.TextColumn("items_origin_groups_items_failover_criteria_status_codes_items"),
		//table.IntegerColumn("items_origin_groups_items_failover_criteria_status_codes_quantity"),
		//table.TextColumn("items_origin_groups_items_id"),
		//table.TextColumn("items_origin_groups_items_members"),
		//table.TextColumn("items_origin_groups_items_members_items"),
		//table.TextColumn("items_origin_groups_items_members_items_origin_id"),
		//table.IntegerColumn("items_origin_groups_items_members_quantity"),
		//table.IntegerColumn("items_origin_groups_quantity"),
		table.TextColumn("items_origins"),
		//table.TextColumn("items_origins_items"),
		//table.IntegerColumn("items_origins_items_connection_attempts"),
		//table.IntegerColumn("items_origins_items_connection_timeout"),
		//table.TextColumn("items_origins_items_custom_headers"),
		//table.TextColumn("items_origins_items_custom_headers_items"),
		//table.TextColumn("items_origins_items_custom_headers_items_header_name"),
		//table.TextColumn("items_origins_items_custom_headers_items_header_value"),
		//table.IntegerColumn("items_origins_items_custom_headers_quantity"),
		//table.TextColumn("items_origins_items_custom_origin_config"),
		//table.IntegerColumn("items_origins_items_custom_origin_config_http_port"),
		//table.IntegerColumn("items_origins_items_custom_origin_config_https_port"),
		//table.IntegerColumn("items_origins_items_custom_origin_config_origin_keepalive_timeout"),
		//table.TextColumn("items_origins_items_custom_origin_config_origin_protocol_policy"),
		//table.IntegerColumn("items_origins_items_custom_origin_config_origin_read_timeout"),
		//table.TextColumn("items_origins_items_custom_origin_config_origin_ssl_protocols"),
		//table.TextColumn("items_origins_items_custom_origin_config_origin_ssl_protocols_items"),
		//table.IntegerColumn("items_origins_items_custom_origin_config_origin_ssl_protocols_quantity"),
		//table.TextColumn("items_origins_items_domain_name"),
		//table.TextColumn("items_origins_items_id"),
		//table.TextColumn("items_origins_items_origin_path"),
		//table.TextColumn("items_origins_items_origin_shield"),
		//table.TextColumn("items_origins_items_origin_shield_enabled"),
		//table.TextColumn("items_origins_items_origin_shield_origin_shield_region"),
		//table.TextColumn("items_origins_items_s3_origin_config"),
		//table.TextColumn("items_origins_items_s3_origin_config_origin_access_identity"),
		//table.IntegerColumn("items_origins_quantity"),
		table.TextColumn("items_price_class"),
		table.TextColumn("items_restrictions"),
		//table.TextColumn("items_restrictions_geo_restriction"),
		//table.TextColumn("items_restrictions_geo_restriction_items"),
		//table.IntegerColumn("items_restrictions_geo_restriction_quantity"),
		//table.TextColumn("items_restrictions_geo_restriction_restriction_type"),
		table.TextColumn("items_status"),
		table.TextColumn("items_viewer_certificate"),
		//table.TextColumn("items_viewer_certificate_acm_certificate_arn"),
		//table.TextColumn("items_viewer_certificate_certificate"),
		//table.TextColumn("items_viewer_certificate_certificate_source"),
		//table.TextColumn("items_viewer_certificate_cloud_front_default_certificate"),
		//table.TextColumn("items_viewer_certificate_iam_certificate_id"),
		//table.TextColumn("items_viewer_certificate_minimum_protocol_version"),
		//table.TextColumn("items_viewer_certificate_ssl_support_method"),
		table.TextColumn("items_web_acl_id"),
		//table.TextColumn("marker"),
		//table.IntegerColumn("max_items"),
		//table.TextColumn("next_marker"),
		table.IntegerColumn("quantity"),
		//table.TextColumn("values"),

	}
}

// ListDistributionsGenerate returns the rows in the table for all configured accounts
func ListDistributionsGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	if len(utilities.ExtConfiguration.ExtConfAws.Accounts) == 0 && extaws.ShouldProcessAccount("aws_cloudfront_distribution", utilities.AwsAccountID) {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_cloudfront_distribution",
			"account":   "default",
		}).Info("processing account")
		results, err := processAccountListDistributions(osqCtx, queryContext, nil)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, results...)
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfAws.Accounts {
			if !extaws.ShouldProcessAccount("aws_cloudfront_distribution", account.ID) {
				continue
			}
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_cloudfront_distribution",
				"account":   account.ID,
			}).Info("processing account")
			results, err := processAccountListDistributions(osqCtx, queryContext, &account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}

	return resultMap, nil
}

func processGlobalListDistributions(osqCtx context.Context, queryContext table.QueryContext, tableConfig *utilities.TableConfig, account *utilities.ExtensionConfigurationAwsAccount) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	sess, err := extaws.GetAwsConfig(account, "aws-global")
	if err != nil {
		return resultMap, err
	}

	accountId := utilities.AwsAccountID
	if account != nil {
		accountId = account.ID
	}

	utilities.GetLogger().WithFields(log.Fields{
		"tableName": "aws_cloudfront_distribution",
		"account":   accountId,
		"region":    "aws-global",
	}).Debug("processing region")

	svc := cloudfront.NewFromConfig(*sess)
	params := &cloudfront.ListDistributionsInput{}

	paginator := cloudfront.NewListDistributionsPaginator(svc, params)

	for {
		page, err := paginator.NextPage(osqCtx)
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_cloudfront_distribution",
				"account":   accountId,
				"region":    "aws-global",
				"task":      "ListDistributions",
				"errString": err.Error(),
			}).Error("failed to process region")
			return resultMap, err
		}
		byteArr, err := json.Marshal(page)
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_cloudfront_distribution",
				"account":   accountId,
				"region":    "aws-global",
				"task":      "ListDistributions",
				"errString": err.Error(),
			}).Error("failed to marshal response")
			return nil, err
		}
		table := utilities.NewTable(byteArr, tableConfig)
		for _, row := range table.Rows {
			if !extaws.ShouldProcessRow(osqCtx, queryContext, "aws_cloudfront_distribution", accountId, "aws-global", row) {
				continue
			}
			result := extaws.RowToMap(row, accountId, "aws-global", tableConfig)
			resultMap = append(resultMap, result)
		}
		if !paginator.HasMorePages() {
			break
		}
	}
	return resultMap, nil
}

func processAccountListDistributions(osqCtx context.Context, queryContext table.QueryContext, account *utilities.ExtensionConfigurationAwsAccount) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	tableConfig, ok := utilities.TableConfigurationMap["aws_cloudfront_distribution"]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_cloudfront_distribution",
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found")
	}
	result, err := processGlobalListDistributions(osqCtx, queryContext, tableConfig, account)
	if err != nil {
		return resultMap, err
	}
	resultMap = append(resultMap, result...)
	return resultMap, nil
}
