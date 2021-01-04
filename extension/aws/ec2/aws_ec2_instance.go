package ec2

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/Uptycs/cloudquery/utilities"

	extaws "github.com/Uptycs/cloudquery/extension/aws"
	"github.com/kolide/osquery-go/plugin/table"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func DescribeInstancesColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("AccountId"),
		table.TextColumn("RegionCode"),
		table.TextColumn("Region"),
		table.IntegerColumn("Instances.AmiLaunchIndex"),
		table.TextColumn("Instances.Architecture"),
		table.TextColumn("Instances.BlockDeviceMappings"),
		table.TextColumn("Instances.CapacityReservationSpecification"),
		table.TextColumn("Instances.ClientToken"),
		table.TextColumn("Instances.CpuOptions"),
		table.IntegerColumn("Instances.EbsOptimized"),
		table.IntegerColumn("Instances.EnaSupport"),
		table.TextColumn("Instances.EnclaveOptions"),
		table.TextColumn("Instances.HibernationOptions"),
		table.TextColumn("Instances.Hypervisor"),
		table.TextColumn("Instances.IamInstanceProfile"),
		table.TextColumn("Instances.ImageId"),
		table.TextColumn("Instances.InstanceId"),
		table.TextColumn("Instances.InstanceLifecycle"),
		table.TextColumn("Instances.InstanceType"),
		table.TextColumn("Instances.KeyName"),
		table.TextColumn("Instances.LaunchTime"),
		table.TextColumn("Instances.MetadataOptions"),
		table.TextColumn("Instances.Monitoring"),
		table.TextColumn("Instances.NetworkInterfaces"),
		table.TextColumn("Instances.Placement"),
		table.TextColumn("Instances.PrivateDnsName"),
		table.TextColumn("Instances.PrivateIpAddress"),
		table.TextColumn("Instances.PublicDnsName"),
		table.TextColumn("Instances.PublicIpAddress"),
		table.TextColumn("Instances.RootDeviceName"),
		table.TextColumn("Instances.RootDeviceType"),
		table.TextColumn("Instances.SecurityGroups"),
		table.TextColumn("Instances.SourceDestCheck"),
		table.TextColumn("Instances.SpotInstanceRequestId"),
		table.TextColumn("Instances.State"),
		table.TextColumn("Instances.StateReason"),
		table.TextColumn("Instances.StateTransitionReason"),
		table.TextColumn("Instances.SubnetId"),
		table.TextColumn("Instances.Tags"),
		table.TextColumn("Instances.VirtualizationType"),
		table.TextColumn("Instances.VpcId"),
		table.TextColumn("OwnerId"),
		table.TextColumn("RequesterId"),
		table.TextColumn("ReservationId"),
	}
}

func DescribeInstancesGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	// TODO: Multi tenancy
	awsSession := session.Must(session.NewSession(&aws.Config{Region: aws.String("us-east-1")}))
	regions, err := extaws.FetchRegions(awsSession)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	tableConfig, ok := utilities.TableConfigurationMap["aws_ec2_instance"]
	if !ok {
		fmt.Println("DescribeInstances.Page: ", err)
		log.Fatal(err)
		return resultMap, fmt.Errorf("table configuration not found")
	}

	for _, region := range regions {
		sess := session.Must(session.NewSession(&aws.Config{
			Region: region.RegionName,
		}))

		ec2Svc := ec2.New(sess)
		params := &ec2.DescribeInstancesInput{}

		err := ec2Svc.DescribeInstancesPages(params,
			func(page *ec2.DescribeInstancesOutput, lastPage bool) bool {
				byteArr, err := json.Marshal(page)
				if err != nil {
					fmt.Fprintf(os.Stderr, "error: %v\n", err)
					os.Exit(1)
				}
				table := utilities.Table{}
				table.Init(byteArr, tableConfig.MaxLevel, tableConfig.GetParsedAttributeConfigMap())
				for _, row := range table.Rows {
					result := extaws.RowToMap(row, *region.RegionName, tableConfig)
					resultMap = append(resultMap, result)
				}
				return lastPage
			})
		if err != nil {
			fmt.Println("DescribeInstances.Page: ", err)
			log.Fatal(err)
			return resultMap, err
		}
	}
	return resultMap, nil
}
