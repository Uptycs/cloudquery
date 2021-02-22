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
func ListCertificatesColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("account_id"),
		table.TextColumn("region_code"),
		table.TextColumn("certificate_arn"),
		table.TextColumn("domain_name"),
		//table.TextColumn("values"),

	}
}

// ListCertificatesGenerate returns the rows in the table for all configured accounts
func ListCertificatesGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	if len(utilities.ExtConfiguration.ExtConfAws.Accounts) == 0 {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_acm_certificate",
			"account":   "default",
		}).Info("processing account")
		results, err := processAccountListCertificates(nil)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, results...)
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfAws.Accounts {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_acm_certificate",
				"account":   account.ID,
			}).Info("processing account")
			results, err := processAccountListCertificates(&account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}

	return resultMap, nil
}

func processRegionListCertificates(tableConfig *utilities.TableConfig, account *utilities.ExtensionConfigurationAwsAccount, region types.Region) ([]map[string]string, error) {
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
		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_acm_certificate",
				"account":   accountId,
				"region":    *region.RegionName,
				"task":      "ListCertificates",
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
				"task":      "ListCertificates",
				"errString": err.Error(),
			}).Error("failed to marshal response")
			return nil, err
		}
		table := utilities.NewTable(byteArr, tableConfig)
		for _, row := range table.Rows {
			result := extaws.RowToMap(row, accountId, *region.RegionName, tableConfig)
			resultMap = append(resultMap, result)
		}
		if !paginator.HasMorePages() {
			break
		}
	}
	return resultMap, nil
}

func processAccountListCertificates(account *utilities.ExtensionConfigurationAwsAccount) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	awsSession, err := extaws.GetAwsConfig(account, "us-east-1")
	if err != nil {
		return resultMap, err
	}
	regions, err := extaws.FetchRegions(context.TODO(), awsSession)
	if err != nil {
		return resultMap, err
	}
	tableConfig, ok := utilities.TableConfigurationMap["aws_acm_certificate"]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_acm_certificate",
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found")
	}
	for _, region := range regions {
		result, err := processRegionListCertificates(tableConfig, account, region)
		if err != nil {
			continue
		}
		resultMap = append(resultMap, result...)
	}
	return resultMap, nil
}
