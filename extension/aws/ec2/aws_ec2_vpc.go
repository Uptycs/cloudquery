package ec2

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/Uptycs/cloudquery/utilities"

	extaws "github.com/Uptycs/cloudquery/extension/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/kolide/osquery-go/plugin/table"
)

func DescribeVpcsColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("AccountId"),
		table.TextColumn("RegionCode"),
		table.TextColumn("Region"),
		table.TextColumn("CidrBlock"),
		table.TextColumn("DhcpOptionsId"),
		table.TextColumn("InstanceTenancy"),
		table.TextColumn("OwnerId"),
		table.TextColumn("State"),
		table.TextColumn("Id"),
	}
}

func DescribeVpcsGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	if len(utilities.ExtConfiguration.ExtConfAws.Accounts) == 0 {
		fmt.Println("Processing default account")
		results, err := processAccountDescribeVpcs(nil)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, results...)
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfAws.Accounts {
			fmt.Println("Processing account:" + account.ID)
			results, err := processAccountDescribeVpcs(&account)
			if err != nil {
				// TODO: Continue to next account or return error ?
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}

	return resultMap, nil
}

func processRegionDescribeVpcs(tableConfig *utilities.TableConfig, account *utilities.ExtensionConfigurationAwsAccount, region *ec2.Region) ([]map[string]string, error) {
	fmt.Println("Processing region:" + *region.RegionName + ", EndPoint:" + *region.Endpoint)
	resultMap := make([]map[string]string, 0)
	sess, err := extaws.GetAwsSession(account, *region.RegionName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	accountId := utilities.AwsAccountId
	if account != nil {
		accountId = account.ID
	}
	svc := ec2.New(sess)
	params := &ec2.DescribeVpcsInput{}

	err = svc.DescribeVpcsPages(params,
		func(page *ec2.DescribeVpcsOutput, lastPage bool) bool {
			byteArr, err := json.Marshal(page)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				os.Exit(1)
			}
			table := utilities.Table{}
			table.Init(byteArr, tableConfig.MaxLevel, tableConfig.GetParsedAttributeConfigMap())
			for _, row := range table.Rows {
				result := extaws.RowToMap(row, accountId, *region.RegionName, tableConfig)
				resultMap = append(resultMap, result)
			}
			return lastPage
		})
	if err != nil {
		fmt.Println("processRegion : DescribeVpcs: ", err)
		log.Fatal(err)
		return resultMap, err
	}
	return resultMap, nil
}

func processAccountDescribeVpcs(account *utilities.ExtensionConfigurationAwsAccount) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	awsSession, err := extaws.GetAwsSession(account, "us-east-1")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	regions, err := extaws.FetchRegions(awsSession)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	tableConfig, ok := utilities.TableConfigurationMap["aws_ec2_vpc"]
	if !ok {
		fmt.Println("getTableConfig: ", err)
		log.Fatal(err)
		return resultMap, fmt.Errorf("table configuration not found")
	}
	for _, region := range regions {
		result, err := processRegionDescribeVpcs(tableConfig, account, region)
		if err != nil {
			fmt.Println("processRegion: ", err)
			log.Fatal(err)
			return resultMap, err
		}
		resultMap = append(resultMap, result...)
	}
	return resultMap, nil
}
