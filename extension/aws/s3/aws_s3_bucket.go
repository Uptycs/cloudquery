package s3

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/Uptycs/cloudquery/utilities"

	extaws "github.com/Uptycs/cloudquery/extension/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/kolide/osquery-go/plugin/table"
)

func ListBucketsColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("AccountId"),
		table.TextColumn("RegionCode"),
		table.TextColumn("Region"),
		table.TextColumn("Name"),
	}
}

func ListBucketsGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	if len(utilities.ExtConfiguration.ExtConfAws.Accounts) == 0 {
		fmt.Println("Processing default account")
		results, err := processAccountListBuckets(nil)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, results...)
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfAws.Accounts {
			fmt.Println("Processing account:" + account.ID)
			results, err := processAccountListBuckets(&account)
			if err != nil {
				// TODO: Continue to next account or return error ?
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}

	return resultMap, nil
}

func processRegionListBuckets(tableConfig *utilities.TableConfig, account *utilities.ExtensionConfigurationAwsAccount, region *ec2.Region) ([]map[string]string, error) {
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
	svc := s3.New(sess)
	params := &s3.ListBucketsInput{}

	output, err := svc.ListBuckets(params)
	if err != nil {
		fmt.Println("ListBuckets.Page: ", err)
		log.Fatal(err)
		return resultMap, err
	}

	byteArr, err := json.Marshal(output)
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
	return resultMap, nil
}

func processAccountListBuckets(account *utilities.ExtensionConfigurationAwsAccount) ([]map[string]string, error) {
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
	tableConfig, ok := utilities.TableConfigurationMap["aws_s3_bucket"]
	if !ok {
		fmt.Println("getTableConfig: ", err)
		log.Fatal(err)
		return resultMap, fmt.Errorf("table configuration not found")
	}
	for _, region := range regions {
		result, err := processRegionListBuckets(tableConfig, account, region)
		if err != nil {
			fmt.Println("processRegion: ", err)
			log.Fatal(err)
			return resultMap, err
		}
		resultMap = append(resultMap, result...)
	}
	return resultMap, nil
}
