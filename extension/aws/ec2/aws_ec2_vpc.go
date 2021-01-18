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
		table.TextColumn("account_id"),
		table.TextColumn("region_code"),
		table.TextColumn("cidr_block"),
		table.TextColumn("cidr_block_association_set"),
		//table.TextColumn("cidr_block_association_set_association_id"),
		//table.TextColumn("cidr_block_association_set_cidr_block"),
		//table.TextColumn("cidr_block_association_set_cidr_block_state"),
		//table.TextColumn("cidr_block_association_set_cidr_block_state_state"),
		//table.TextColumn("cidr_block_association_set_cidr_block_state_status_message"),
		table.TextColumn("dhcp_options_id"),
		table.TextColumn("instance_tenancy"),
		table.TextColumn("ipv6_cidr_block_association_set"),
		//table.TextColumn("ipv6_cidr_block_association_set_association_id"),
		//table.TextColumn("ipv6_cidr_block_association_set_ipv6_cidr_block"),
		//table.TextColumn("ipv6_cidr_block_association_set_ipv6_cidr_block_state"),
		//table.TextColumn("ipv6_cidr_block_association_set_ipv6_cidr_block_state_state"),
		//table.TextColumn("ipv6_cidr_block_association_set_ipv6_cidr_block_state_status_message"),
		//table.TextColumn("ipv6_cidr_block_association_set_ipv6_pool"),
		//table.TextColumn("ipv6_cidr_block_association_set_network_border_group"),
		table.TextColumn("is_default"),
		table.TextColumn("owner_id"),
		table.TextColumn("state"),
		table.TextColumn("tags"),
		table.TextColumn("tags_key"),
		table.TextColumn("tags_value"),
		table.TextColumn("vpc_id"),
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
			table := utilities.NewTable(byteArr, tableConfig)
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
