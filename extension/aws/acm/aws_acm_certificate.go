/**
 * Copyright (c) 2020-present, The cloudquery authors
 *
 * This source code is licensed as defined by the LICENSE file found in the
 * root directory of this source tree.
 *
 * SPDX-License-Identifier: (Apache-2.0 OR GPL-2.0-only)
 */

package acm

import (
	"context"
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/Uptycs/cloudquery/utilities"

	"github.com/Uptycs/basequery-go/plugin/table"
	extaws "github.com/Uptycs/cloudquery/extension/aws"
	"github.com/aws/aws-sdk-go-v2/service/acm"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

// ListCertificatesColumns returns the list of columns in the table
func DescribeCertificateColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("account_id"),
		table.TextColumn("region_code"),
		table.TextColumn("certificate_arn"),
		table.TextColumn("certificate_authority_arn"),
		table.TextColumn("created_at"),
		table.TextColumn("domain_name"),
		table.TextColumn("domain_validation_options"),
		table.TextColumn("extended_key_usages"),
		table.TextColumn("failure_reason"),
		table.TextColumn("imported_at"),
		table.TextColumn("in_use_by"),
		table.TextColumn("issued_at"),
		table.TextColumn("issuer"),
		table.TextColumn("key_algorithm"),
		table.TextColumn("key_usages"),
		table.TextColumn("not_after"),
		table.TextColumn("not_before"),
		table.TextColumn("options"),
		table.TextColumn("renewal_eligibility"),
		table.TextColumn("renewal_summary"),
		table.TextColumn("revocation_reason"),
		table.TextColumn("revoked_at"),
		table.TextColumn("serial"),
		table.TextColumn("signature_algorithm"),
		table.TextColumn("status"),
		table.TextColumn("subject"),
		table.TextColumn("type"),
	}
}

// ListTagsForCertificateColumns returns the list of columns in the table
func ListTagsForCertificateColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("account_id"),
		table.TextColumn("region_code"),
		table.TextColumn("key"),
		table.TextColumn("value"),
	}
}

// DescribeCertificatesGenerate returns the rows in the table for all configured accounts
func DescribeCertificateGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	tableConfig, ok := utilities.TableConfigurationMap["aws_acm_certificate"]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_acm_certificate",
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found")
	}
	if len(utilities.ExtConfiguration.ExtConfAws.Accounts) == 0 && extaws.ShouldProcessAccount("aws_acm_certificate", utilities.AwsAccountID) {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_acm_certificate",
			"account":   "default",
		}).Info("processing account")
		results, err := processAccountListCertificates(osqCtx, queryContext, nil, "aws_acm_certificate", tableConfig)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, results...)
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfAws.Accounts {
			if !extaws.ShouldProcessAccount("aws_acm_certificate", account.ID) {
				continue
			}
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_acm_certificate",
				"account":   account.ID,
			}).Info("processing account")
			results, err := processAccountListCertificates(osqCtx, queryContext, &account, "aws_acm_certificate", tableConfig)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}

	return resultMap, nil
}

// ListTagsForCertificateGenerate returns the rows in the table for all configured accounts
func ListTagsForCertificateGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	tableConfig, ok := utilities.TableConfigurationMap["aws_acm_tag"]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_acm_tag",
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found")
	}
	if len(utilities.ExtConfiguration.ExtConfAws.Accounts) == 0 && extaws.ShouldProcessAccount("aws_acm_tag", utilities.AwsAccountID) {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_acm_tag",
			"account":   "default",
		}).Info("processing account")
		results, err := processAccountListCertificates(osqCtx, queryContext, nil, "aws_acm_tag", tableConfig)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, results...)
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfAws.Accounts {
			if !extaws.ShouldProcessAccount("aws_acm_tag", account.ID) {
				continue
			}
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_acm_tag",
				"account":   account.ID,
			}).Info("processing account")
			results, err := processAccountListCertificates(osqCtx, queryContext, &account, "aws_acm_tag", tableConfig)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}
	return resultMap, nil
}

func processAccountListCertificates(osqCtx context.Context, queryContext table.QueryContext, account *utilities.ExtensionConfigurationAwsAccount, tableName string, tableConfig *utilities.TableConfig) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	awsSession, err := extaws.GetAwsConfig(account, "us-east-1")
	if err != nil {
		return resultMap, err
	}
	regions, err := extaws.FetchRegions(osqCtx, awsSession)
	if err != nil {
		return resultMap, err
	}
	for _, region := range regions {
		accountId := utilities.AwsAccountID
		if account != nil {
			accountId = account.ID
		}
		if !extaws.ShouldProcessRegion(tableName, accountId, *region.RegionName) {
			continue
		}
		result, err := processRegionListCertificates(osqCtx, queryContext, tableConfig, account, region, tableName)
		if err != nil {
			continue
		}
		resultMap = append(resultMap, result...)
	}
	return resultMap, nil
}

func processRegionListCertificates(osqCtx context.Context, queryContext table.QueryContext, tableConfig *utilities.TableConfig, account *utilities.ExtensionConfigurationAwsAccount, region types.Region, tableName string) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	sess, err := extaws.GetAwsConfig(account, *region.RegionName)
	if err != nil {
		return resultMap, err
	}

	accountId := utilities.AwsAccountID
	if account != nil {
		accountId = account.ID
	}

	utilities.GetLogger().WithFields(log.Fields{
		"tableName": "aws_acm_certificate",
		"account":   accountId,
		"region":    *region.RegionName,
	}).Debug("processing region")

	svc := acm.NewFromConfig(*sess)
	params := &acm.ListCertificatesInput{}

	paginator := acm.NewListCertificatesPaginator(svc, params)

	for {
		page, err := paginator.NextPage(osqCtx)
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": tableName,
				"account":   accountId,
				"region":    *region.RegionName,
				"task":      "ListCertificates",
				"errString": err.Error(),
			}).Error("failed to process region")
			return resultMap, err
		}
		
		for _, certificateSummary := range page.CertificateSummaryList {
			
			if tableName == "aws_acm_certificate" {
				result, err := processDescribeCertificate(osqCtx, queryContext, tableConfig, account, region, *certificateSummary.CertificateArn)
			
				if err != nil {
					continue
				}
				resultMap = append(resultMap, result...)
			}else{
				result, err := processListTagsForCertificate(osqCtx, queryContext, tableConfig, account, region, *certificateSummary.CertificateArn)
			
				if err != nil {
					continue
				}
				resultMap = append(resultMap, result...)
			}
		}
	
		if !paginator.HasMorePages() {
			break
		}
	}
	return resultMap, nil
}

func processDescribeCertificate(osqCtx context.Context, queryContext table.QueryContext, tableConfig *utilities.TableConfig, account *utilities.ExtensionConfigurationAwsAccount, region types.Region, CertificateArn string) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	sess, err := extaws.GetAwsConfig(account, *region.RegionName)
	if err != nil {
		return resultMap, err
	}
	accountId := utilities.AwsAccountID
	if account != nil {
		accountId = account.ID
	}

	utilities.GetLogger().WithFields(log.Fields{
		"tableName": "aws_acm_certificate",
		"account":   accountId,
		//"region":    *region.RegionName,
	}).Debug("processing region")

	svc := acm.NewFromConfig(*sess)
	params := &acm.DescribeCertificateInput{
		CertificateArn: &CertificateArn,
	}
	page, err := svc.DescribeCertificate(osqCtx, params)
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_acm_certificate",
			"account":   accountId,
			"region":    *region.RegionName,
			"task":      "DescribeCertificate",
			"errString": err.Error(),
		}).Error("failed to process region")
		return resultMap, err
	}

	byteArr, err := json.Marshal(page)
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_acm_certificate",
			"account":   accountId,
			"region":    *region.RegionName,
			"errString": err.Error(),
		}).Error("failed to marshal response")
		return resultMap, err
	}
	table := utilities.NewTable(byteArr, tableConfig)
	for _, row := range table.Rows {
		if !extaws.ShouldProcessRow(osqCtx, queryContext, "aws_acm_certificate", accountId, *region.RegionName, row) {
			continue
		}
		result := extaws.RowToMap(row, accountId, *region.RegionName, tableConfig)
		resultMap = append(resultMap, result)
	}
	return resultMap, nil
}

func processListTagsForCertificate(osqCtx context.Context, queryContext table.QueryContext, tableConfig *utilities.TableConfig, account *utilities.ExtensionConfigurationAwsAccount, region types.Region, CertificateArn string) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	sess, err := extaws.GetAwsConfig(account, *region.RegionName)
	if err != nil {
		return resultMap, err
	}
	accountId := utilities.AwsAccountID
	if account != nil {
		accountId = account.ID
	}

	utilities.GetLogger().WithFields(log.Fields{
		"tableName": "aws_acm_tag",
		"account":   accountId,
		//"region":    *region.RegionName,
	}).Debug("processing region")

	svc := acm.NewFromConfig(*sess)
	params := &acm.ListTagsForCertificateInput{
		CertificateArn: &CertificateArn,
	}
	page, err := svc.ListTagsForCertificate(osqCtx, params)
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_acm_tag",
			"account":   accountId,
			"region":    *region.RegionName,
			"task":      "ListTagsForCertificate",
			"errString": err.Error(),
		}).Error("failed to process region")
		return resultMap, err
	}

	byteArr, err := json.Marshal(page)
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_acm_tag",
			"account":   accountId,
			"region":    *region.RegionName,
			"errString": err.Error(),
		}).Error("failed to marshal response")
		return resultMap, err
	}
	table := utilities.NewTable(byteArr, tableConfig)
	for _, row := range table.Rows {
		if !extaws.ShouldProcessRow(osqCtx, queryContext, "aws_acm_tag", accountId, *region.RegionName, row) {
			continue
		}
		result := extaws.RowToMap(row, accountId, *region.RegionName, tableConfig)
		resultMap = append(resultMap, result)
	}
	return resultMap, nil
}

