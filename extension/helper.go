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

func readProjectIDFromCredentialFile(filePath string) string {
	reader, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("error reading %s. err:%s\n", filePath, err.Error())
		return ""
	}
	var jsonObj map[string]interface{}
	errUnmarshal := json.Unmarshal(reader, &jsonObj)
	if errUnmarshal != nil {
		fmt.Printf("error unmarshaling json in %s. err:%s\n", filePath, errUnmarshal.Error())
		return ""
	}

	if idIntfc, found := jsonObj["project_id"]; found {
		return idIntfc.(string)
	}

	fmt.Printf("cannot find \"project_id\" property in file %s. \n", filePath)
	return ""
}

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

	// Set projectID for GCP accounts
	for idx := range utilities.ExtConfiguration.ExtConfGcp.Accounts {
		keyFilePath := utilities.ExtConfiguration.ExtConfGcp.Accounts[idx].KeyFile
		projectID := readProjectIDFromCredentialFile(keyFilePath)
		utilities.ExtConfiguration.ExtConfGcp.Accounts[idx].ProjectId = projectID
	}

	// Read project ID from ADC
	adcFilePath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if adcFilePath != "" {
		utilities.DefaultGcpProjectID = readProjectIDFromCredentialFile(adcFilePath)
	}

	if len(utilities.ExtConfiguration.ExtConfGcp.Accounts) == 0 {
		if adcFilePath == "" {
			fmt.Println("missing env GOOGLE_APPLICATION_CREDENTIALS")
		} else if utilities.DefaultGcpProjectID == "" {
			fmt.Println("missing Default Project ID for GCP")
		} else {
			fmt.Printf("Gcp accounts not found in extension_config. Falling back to ADC\n")
		}
	}

	return nil
}

func readTableConfig(filePath string) error {
	reader, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	var configurations map[string]*utilities.TableConfig
	errUnmarshal := json.Unmarshal(reader, &configurations)
	if errUnmarshal != nil {
		return errUnmarshal
	}
	for tableName, config := range configurations {
		fmt.Println("Found configuration for table:" + tableName)
		config.InitParsedAttributeConfigMap()
		utilities.TableConfigurationMap[tableName] = config
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

var gcpComputeHandler = compute.NewGcpComputeHandler(compute.NewGcpComputeImpl())
var gcpStorageHandler = storage.NewGcpStorageHandler(storage.NewGcpStorageImpl())

func registerPlugins(server *osquery.ExtensionManagerServer) {
	server.RegisterPlugin(table.NewPlugin("aws_ec2_instance", ec2.DescribeInstancesColumns(), ec2.DescribeInstancesGenerate))
	server.RegisterPlugin(table.NewPlugin("aws_ec2_vpc", ec2.DescribeVpcsColumns(), ec2.DescribeVpcsGenerate))
	server.RegisterPlugin(table.NewPlugin("aws_ec2_subnet", ec2.DescribeSubnetsColumns(), ec2.DescribeSubnetsGenerate))
	server.RegisterPlugin(table.NewPlugin("aws_ec2_image", ec2.DescribeImagesColumns(), ec2.DescribeImagesGenerate))
	server.RegisterPlugin(table.NewPlugin("aws_s3_bucket", s3.ListBucketsColumns(), s3.ListBucketsGenerate))
	server.RegisterPlugin(table.NewPlugin("gcp_compute_instance", gcpComputeHandler.GcpComputeInstancesColumns(), gcpComputeHandler.GcpComputeInstancesGenerate))
	server.RegisterPlugin(table.NewPlugin("gcp_compute_network", gcpComputeHandler.GcpComputeNetworksColumns(), gcpComputeHandler.GcpComputeNetworksGenerate))
	server.RegisterPlugin(table.NewPlugin("gcp_compute_disk", gcpComputeHandler.GcpComputeDisksColumns(), gcpComputeHandler.GcpComputeDisksGenerate))
	server.RegisterPlugin(table.NewPlugin("gcp_storage_bucket", gcpStorageHandler.GcpStorageBucketColumns(), gcpStorageHandler.GcpStorageBucketGenerate))
	server.RegisterPlugin(table.NewPlugin("azure_compute_vm", azurecompute.VirtualMachinesColumns(), azurecompute.VirtualMachinesGenerate))
	server.RegisterPlugin(table.NewPlugin("azure_compute_networkinterface", azurecompute.InterfacesColumns(), azurecompute.InterfacesGenerate))
}
