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
		table.TextColumn("AccountId"),
		table.TextColumn("RegionCode"),
		table.TextColumn("Region"),
		table.IntegerColumn("Instances_AmiLaunchIndex"),
		table.TextColumn("Instances_Architecture"),
		table.TextColumn("Instances_BlockDeviceMappings"),
		table.TextColumn("Instances_CapacityReservationSpecification"),
		table.TextColumn("Instances_ClientToken"),
		table.TextColumn("Instances_CpuOptions"),
		table.IntegerColumn("Instances_EbsOptimized"),
		table.IntegerColumn("Instances_EnaSupport"),
		table.TextColumn("Instances_EnclaveOptions"),
		table.TextColumn("Instances_HibernationOptions"),
		table.TextColumn("Instances_Hypervisor"),
		table.TextColumn("Instances_IamInstanceProfile"),
		table.TextColumn("Instances_ImageId"),
		table.TextColumn("Instances_InstanceId"),
		table.TextColumn("Instances_InstanceLifecycle"),
		table.TextColumn("Instances_InstanceType"),
		table.TextColumn("Instances_KeyName"),
		table.TextColumn("Instances_LaunchTime"),
		table.TextColumn("Instances_MetadataOptions"),
		table.TextColumn("Instances_Monitoring"),
		table.TextColumn("Instances_NetworkInterfaces"),
		table.TextColumn("Instances_Placement"),
		table.TextColumn("Instances_PrivateDnsName"),
		table.TextColumn("Instances_PrivateIpAddress"),
		table.TextColumn("Instances_PublicDnsName"),
		table.TextColumn("Instances_PublicIpAddress"),
		table.TextColumn("Instances_RootDeviceName"),
		table.TextColumn("Instances_RootDeviceType"),
		table.TextColumn("Instances_SecurityGroups"),
		table.TextColumn("Instances_SourceDestCheck"),
		table.TextColumn("Instances_SpotInstanceRequestId"),
		table.TextColumn("Instances_State"),
		table.TextColumn("Instances_StateReason"),
		table.TextColumn("Instances_StateTransitionReason"),
		table.TextColumn("Instances_SubnetId"),
		table.TextColumn("Instances_Tags"),
		table.TextColumn("Instances_VirtualizationType"),
		table.TextColumn("Instances_VpcId"),
		table.TextColumn("OwnerId"),
		table.TextColumn("RequesterId"),
		table.TextColumn("ReservationId"),
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
