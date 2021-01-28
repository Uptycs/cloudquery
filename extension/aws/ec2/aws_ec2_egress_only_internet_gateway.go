package ec2

import (
	"context"
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/Uptycs/cloudquery/utilities"

	extaws "github.com/Uptycs/cloudquery/extension/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/kolide/osquery-go/plugin/table"
)

// DescribeEgressOnlyInternetGatewaysColumns returns the list of columns in the table
func DescribeEgressOnlyInternetGatewaysColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("account_id"),
		table.TextColumn("region_code"),
		table.TextColumn("attachments"),
		//table.TextColumn("attachments_state"),
		//table.TextColumn("attachments_vpc_id"),
		table.TextColumn("egress_only_internet_gateway_id"),
		table.TextColumn("tags"),
		//table.TextColumn("tags_key"),
		//table.TextColumn("tags_value"),

	}
}

// DescribeEgressOnlyInternetGatewaysGenerate returns the rows in the table for all configured accounts
func DescribeEgressOnlyInternetGatewaysGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	if len(utilities.ExtConfiguration.ExtConfAws.Accounts) == 0 {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_ec2_egress_only_internet_gateway",
			"account":   "default",
		}).Info("processing account")
		results, err := processAccountDescribeEgressOnlyInternetGateways(nil)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, results...)
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfAws.Accounts {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_ec2_egress_only_internet_gateway",
				"account":   account.ID,
			}).Info("processing account")
			results, err := processAccountDescribeEgressOnlyInternetGateways(&account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}

	return resultMap, nil
}

func processRegionDescribeEgressOnlyInternetGateways(tableConfig *utilities.TableConfig, account *utilities.ExtensionConfigurationAwsAccount, region *ec2.Region) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	sess, err := extaws.GetAwsSession(account, *region.RegionName)
	if err != nil {
		return resultMap, err
	}

	accountId := utilities.AwsAccountID
	if account != nil {
		accountId = account.ID
	}

	utilities.GetLogger().WithFields(log.Fields{
		"tableName": "aws_ec2_egress_only_internet_gateway",
		"account":   accountId,
		"region":    *region.RegionName,
	}).Debug("processing region")

	svc := ec2.New(sess)
	params := &ec2.DescribeEgressOnlyInternetGatewaysInput{}

	err = svc.DescribeEgressOnlyInternetGatewaysPages(params,
		func(page *ec2.DescribeEgressOnlyInternetGatewaysOutput, lastPage bool) bool {
			byteArr, err := json.Marshal(page)
			if err != nil {
				utilities.GetLogger().WithFields(log.Fields{
					"tableName": "aws_ec2_egress_only_internet_gateway",
					"account":   accountId,
					"region":    *region.RegionName,
					"errString": err.Error(),
				}).Error("failed to marshal response")
				return lastPage
			}
			table := utilities.NewTable(byteArr, tableConfig)
			for _, row := range table.Rows {
				result := extaws.RowToMap(row, accountId, *region.RegionName, tableConfig)
				resultMap = append(resultMap, result)
			}
			return lastPage
		})
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_ec2_egress_only_internet_gateway",
			"account":   accountId,
			"region":    *region.RegionName,
			"task":      "DescribeEgressOnlyInternetGateways",
			"errString": err.Error(),
		}).Error("failed to process region")
		return resultMap, err
	}
	return resultMap, nil
}

func processAccountDescribeEgressOnlyInternetGateways(account *utilities.ExtensionConfigurationAwsAccount) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	awsSession, err := extaws.GetAwsSession(account, "us-east-1")
	if err != nil {
		return resultMap, err
	}
	regions, err := extaws.FetchRegions(awsSession)
	if err != nil {
		return resultMap, err
	}
	tableConfig, ok := utilities.TableConfigurationMap["aws_ec2_egress_only_internet_gateway"]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_ec2_egress_only_internet_gateway",
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found")
	}
	for _, region := range regions {
		result, err := processRegionDescribeEgressOnlyInternetGateways(tableConfig, account, region)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, result...)
	}
	return resultMap, nil
}
