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

func DescribeInstancesColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("account_id"),
		table.TextColumn("region_code"),
		table.TextColumn("region"),
		table.IntegerColumn("ami_launch_index"),
		table.TextColumn("architecture"),
		table.TextColumn("block_device_mappings"),
		table.TextColumn("capacity_reservation_specification"),
		table.TextColumn("client_token"),
		table.TextColumn("cpu_options"),
		table.IntegerColumn("ebs_optimized"),
		table.IntegerColumn("ena_support"),
		table.TextColumn("enclave_options"),
		table.TextColumn("hibernation_options"),
		table.TextColumn("hypervisor"),
		table.TextColumn("iam_instance_profile"),
		table.TextColumn("image_id"),
		table.TextColumn("instance_id"),
		table.TextColumn("instance_lifecycle"),
		table.TextColumn("instance_type"),
		table.TextColumn("key_name"),
		table.TextColumn("launch_time"),
		table.TextColumn("metadata_options"),
		table.TextColumn("monitoring"),
		table.TextColumn("network_interfaces"),
		table.TextColumn("placement"),
		table.TextColumn("private_dns_name"),
		table.TextColumn("private_ip_address"),
		table.TextColumn("public_dns_name"),
		table.TextColumn("public_ip_address"),
		table.TextColumn("root_device_name"),
		table.TextColumn("root_device_type"),
		table.TextColumn("security_groups"),
		table.TextColumn("source_dest_check"),
		table.TextColumn("spot_instance_request_id"),
		table.TextColumn("state"),
		table.TextColumn("state_reason"),
		table.TextColumn("state_transition_reason"),
		table.TextColumn("subnet_id"),
		table.TextColumn("tags"),
		table.TextColumn("virtualization_type"),
		table.TextColumn("vpc_id"),
		table.TextColumn("owner_id"),
		table.TextColumn("requester_id"),
		table.TextColumn("reservation_id"),
	}
}

func DescribeInstancesGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	if len(utilities.ExtConfiguration.ExtConfAws.Accounts) == 0 {
		fmt.Println("Processing default account")
		results, err := processAccountDescribeInstances(nil)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, results...)
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfAws.Accounts {
			fmt.Println("Processing account:" + account.ID)
			results, err := processAccountDescribeInstances(&account)
			if err != nil {
				// TODO: Continue to next account or return error ?
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}

	return resultMap, nil
}

func processRegionDescribeInstances(tableConfig *utilities.TableConfig, account *utilities.ExtensionConfigurationAwsAccount, region *ec2.Region) ([]map[string]string, error) {
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
	params := &ec2.DescribeInstancesInput{}

	err = svc.DescribeInstancesPages(params,
		func(page *ec2.DescribeInstancesOutput, lastPage bool) bool {
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
		fmt.Println("processRegion : DescribeInstances: ", err)
		log.Fatal(err)
		return resultMap, err
	}
	return resultMap, nil
}

func processAccountDescribeInstances(account *utilities.ExtensionConfigurationAwsAccount) ([]map[string]string, error) {
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
	tableConfig, ok := utilities.TableConfigurationMap["aws_ec2_instance"]
	if !ok {
		fmt.Println("getTableConfig: ", err)
		log.Fatal(err)
		return resultMap, fmt.Errorf("table configuration not found")
	}
	for _, region := range regions {
		result, err := processRegionDescribeInstances(tableConfig, account, region)
		if err != nil {
			fmt.Println("processRegion: ", err)
			log.Fatal(err)
			return resultMap, err
		}
		resultMap = append(resultMap, result...)
	}
	return resultMap, nil
}
