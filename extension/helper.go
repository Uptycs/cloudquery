package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/Uptycs/cloudquery/extension/aws/s3"

	"github.com/Uptycs/cloudquery/utilities"

	"github.com/Uptycs/cloudquery/extension/aws/ec2"
	azurecompute "github.com/Uptycs/cloudquery/extension/azure/compute"
	"github.com/Uptycs/cloudquery/extension/gcp/compute"
	"github.com/Uptycs/cloudquery/extension/gcp/storage"

	"github.com/kolide/osquery-go"
	"github.com/kolide/osquery-go/plugin/table"
)

func readExtensionConfigurations(filePath string) error {
	utilities.AwsAccountId = os.Getenv("AWS_ACCOUNT_ID")
	reader, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	extConfig := utilities.ExtensionConfiguration{}
	errUnmarshal := json.Unmarshal(reader, &extConfig)
	if errUnmarshal != nil {
		return errUnmarshal
	}
	//fmt.Printf("Config:%v\n", extConfig)
	utilities.ExtConfiguration = extConfig
	return nil
}

func readTableConfig(filePath string) error {
	reader, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	var configurations map[string]interface{}
	errUnmarshal := json.Unmarshal(reader, &configurations)
	if errUnmarshal != nil {
		return errUnmarshal
	}
	for tableName, config := range configurations {
		fmt.Println("Found configuration for table:" + tableName)
		configAsMap := config.(map[string]interface{})
		tableConfig := utilities.TableConfig{}
		tableConfig.Init(configAsMap)
		//fmt.Printf("OrigTableConfig:%v\n", tableConfig)
		//fmt.Printf("OrigTableConfigMap:%v\n", configAsMap)
		utilities.TableConfigurationMap[tableName] = &tableConfig

		//fmt.Printf("So far Read config for %d tables\n", len(utilities.TableConfigurationMap))
	}
	return nil
}

func readTableConfigurations() {
	var awsConfigFileList = []string{"aws/ec2/table_config.json", "aws/s3/table_config.json"}
	var gcpConfigFileList = []string{"gcp/compute/table_config.json", "gcp/storage/table_config.json"}
	var azureConfigFileList = []string{"azure/compute/table_config.json"}
	var configFileList = append(awsConfigFileList, gcpConfigFileList...)
	configFileList = append(configFileList, azureConfigFileList...)

	for _, fileName := range configFileList {
		fmt.Println("Reading config file:" + *homeDirectory + string(os.PathSeparator) + fileName)
		readTableConfig(*homeDirectory + string(os.PathSeparator) + fileName)
	}
	fmt.Printf("Read config for total %d tables\n", len(utilities.TableConfigurationMap))
}

func registerPlugins(server *osquery.ExtensionManagerServer) {
	server.RegisterPlugin(table.NewPlugin("aws_ec2_instance", ec2.DescribeInstancesColumns(), ec2.DescribeInstancesGenerate))
	server.RegisterPlugin(table.NewPlugin("aws_ec2_vpc", ec2.DescribeVpcsColumns(), ec2.DescribeVpcsGenerate))
	server.RegisterPlugin(table.NewPlugin("aws_ec2_subnet", ec2.DescribeSubnetsColumns(), ec2.DescribeSubnetsGenerate))
	server.RegisterPlugin(table.NewPlugin("aws_ec2_image", ec2.DescribeImagesColumns(), ec2.DescribeImagesGenerate))
	server.RegisterPlugin(table.NewPlugin("aws_s3_bucket", s3.ListBucketsColumns(), s3.ListBucketsGenerate))
	server.RegisterPlugin(table.NewPlugin("gcp_compute_instance", compute.GcpComputeInstanceColumns(), compute.GcpComputeInstanceGenerate))
	server.RegisterPlugin(table.NewPlugin("gcp_compute_network", compute.GcpComputeNetworkColumns(), compute.GcpComputeNetworkGenerate))
	server.RegisterPlugin(table.NewPlugin("gcp_compute_disk", compute.GcpComputeDiskColumns(), compute.GcpComputeDiskGenerate))
	server.RegisterPlugin(table.NewPlugin("gcp_storage_bucket", storage.GcpStorageBucketColumns(), storage.GcpStorageBucketGenerate))
	server.RegisterPlugin(table.NewPlugin("azure_compute_vm", azurecompute.ComputeVmColumns(), azurecompute.ComputeVmGenerate))
}
