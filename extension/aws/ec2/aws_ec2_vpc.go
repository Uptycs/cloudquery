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
	// TODO: Multi tenancy
	awsSession := session.Must(session.NewSession(&aws.Config{Region: aws.String("us-east-1")}))
	regions, err := extaws.FetchRegions(awsSession)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	tableConfig, ok := utilities.TableConfigurationMap["aws_ec2_vpc"]
	if !ok {
		fmt.Println("DescribeVpcs.Page: ", err)
		log.Fatal(err)
		return resultMap, fmt.Errorf("table configuration not found")
	}

	for _, region := range regions {
		sess := session.Must(session.NewSession(&aws.Config{
			Region: region.RegionName,
		}))

		ec2Svc := ec2.New(sess)
		params := &ec2.DescribeVpcsInput{}

		err := ec2Svc.DescribeVpcsPages(params,
			func(page *ec2.DescribeVpcsOutput, lastPage bool) bool {
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
			fmt.Println("DescribeVpcs.Page: ", err)
			log.Fatal(err)
			return resultMap, err
		}
	}
	return resultMap, nil
}
