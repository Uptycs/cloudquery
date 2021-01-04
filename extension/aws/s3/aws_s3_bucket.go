package s3

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
	"github.com/aws/aws-sdk-go/service/s3"
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
	// TODO: Multi tenancy
	awsSession := session.Must(session.NewSession(&aws.Config{Region: aws.String("us-east-1")}))
	regions, err := extaws.FetchRegions(awsSession)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	tableConfig, ok := utilities.TableConfigurationMap["aws_s3_bucket"]
	if !ok {
		fmt.Println("ListBuckets.Page: ", err)
		log.Fatal(err)
		return resultMap, fmt.Errorf("table configuration not found")
	}

	for _, region := range regions {
		sess := session.Must(session.NewSession(&aws.Config{
			Region: region.RegionName,
		}))

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
			result := extaws.RowToMap(row, *region.RegionName, tableConfig)
			resultMap = append(resultMap, result)
		}
	}
	return resultMap, nil
}
