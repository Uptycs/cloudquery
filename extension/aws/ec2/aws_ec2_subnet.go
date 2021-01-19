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

func DescribeSubnetsColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("account_id"),
		table.TextColumn("region_code"),
		table.TextColumn("assign_ipv6_address_on_creation"),
		table.TextColumn("availability_zone"),
		//table.TextColumn("availability_zone_id"),
		table.BigIntColumn("available_ip_address_count"),
		table.TextColumn("cidr_block"),
		table.TextColumn("customer_owned_ipv4_pool"),
		table.TextColumn("default_for_az"),
		table.TextColumn("ipv6_cidr_block_association_set"),
		//table.TextColumn("ipv6_cidr_block_association_set_association_id"),
		//table.TextColumn("ipv6_cidr_block_association_set_ipv6_cidr_block"),
		//table.TextColumn("ipv6_cidr_block_association_set_ipv6_cidr_block_state"),
		//table.TextColumn("ipv6_cidr_block_association_set_ipv6_cidr_block_state_state"),
		//table.TextColumn("ipv6_cidr_block_association_set_ipv6_cidr_block_state_status_message"),
		table.TextColumn("map_customer_owned_ip_on_launch"),
		table.TextColumn("map_public_ip_on_launch"),
		table.TextColumn("outpost_arn"),
		table.TextColumn("owner_id"),
		table.TextColumn("state"),
		table.TextColumn("subnet_arn"),
		table.TextColumn("subnet_id"),
		table.TextColumn("tags"),
		//table.TextColumn("tags_key"),
		//table.TextColumn("tags_value"),
		table.TextColumn("vpc_id"),
	}
}

func DescribeSubnetsGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	if len(utilities.ExtConfiguration.ExtConfAws.Accounts) == 0 {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_ec2_subnet",
			"account":   "default",
		}).Info("processing account")
		results, err := processAccountDescribeSubnets(nil)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, results...)
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfAws.Accounts {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_ec2_subnet",
				"account":   account.ID,
			}).Info("processing account")
			results, err := processAccountDescribeSubnets(&account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}

	return resultMap, nil
}

func processRegionDescribeSubnets(tableConfig *utilities.TableConfig, account *utilities.ExtensionConfigurationAwsAccount, region *ec2.Region) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	sess, err := extaws.GetAwsSession(account, *region.RegionName)
	if err != nil {
		return resultMap, err
	}

	accountId := utilities.AwsAccountId
	if account != nil {
		accountId = account.ID
	}

	utilities.GetLogger().WithFields(log.Fields{
		"tableName": "aws_ec2_subnet",
		"account":   accountId,
		"region":    *region.RegionName,
	}).Debug("processing region")

	svc := ec2.New(sess)
	params := &ec2.DescribeSubnetsInput{}

	err = svc.DescribeSubnetsPages(params,
		func(page *ec2.DescribeSubnetsOutput, lastPage bool) bool {
			byteArr, err := json.Marshal(page)
			if err != nil {
				utilities.GetLogger().WithFields(log.Fields{
					"tableName": "aws_ec2_subnet",
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
			"tableName": "aws_ec2_subnet",
			"account":   accountId,
			"region":    *region.RegionName,
			"task":      "DescribeSubnets",
			"errString": err.Error(),
		}).Error("failed to process region")
		return resultMap, err
	}
	return resultMap, nil
}

func processAccountDescribeSubnets(account *utilities.ExtensionConfigurationAwsAccount) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	awsSession, err := extaws.GetAwsSession(account, "us-east-1")
	if err != nil {
		return resultMap, err
	}
	regions, err := extaws.FetchRegions(awsSession)
	if err != nil {
		return resultMap, err
	}
	tableConfig, ok := utilities.TableConfigurationMap["aws_ec2_subnet"]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_ec2_subnet",
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found")
	}
	for _, region := range regions {
		result, err := processRegionDescribeSubnets(tableConfig, account, region)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, result...)
	}
	return resultMap, nil
}
