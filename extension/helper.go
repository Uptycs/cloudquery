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
	log "github.com/sirupsen/logrus"
)

func initializeLogger() {
	utilities.CreateLogger(*verbose, utilities.ExtConfiguration.ExtConfLog.MaxSize,
		utilities.ExtConfiguration.ExtConfLog.MaxBackups, utilities.ExtConfiguration.ExtConfLog.MaxAge,
		utilities.ExtConfiguration.ExtConfLog.FileName)
}

func readProjectIDFromCredentialFile(filePath string) string {
	reader, err := ioutil.ReadFile(filePath)
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"fileName":  filePath,
			"errString": err.Error(),
		}).Info("failed to read default gcp credentials file")
		return ""
	}
	var jsonObj map[string]interface{}
	errUnmarshal := json.Unmarshal(reader, &jsonObj)
	if errUnmarshal != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"fileName":  filePath,
			"errString": errUnmarshal.Error(),
		}).Error("failed to unmarshal json")
		return ""
	}

	if idIntfc, found := jsonObj["project_id"]; found {
		return idIntfc.(string)
	}

	utilities.GetLogger().WithFields(log.Fields{
		"fileName": filePath,
	}).Error("failed to find project_id")
	return ""
}

func readExtensionConfigurations(filePath string) error {
	utilities.AwsAccountId = os.Getenv("AWS_ACCOUNT_ID")
	reader, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("failed to read configuration file %s. err:%v\n", filePath, err)
		return err
	}
	extConfig := utilities.ExtensionConfiguration{}
	errUnmarshal := json.Unmarshal(reader, &extConfig)
	if errUnmarshal != nil {
		return errUnmarshal
	}
	utilities.ExtConfiguration = extConfig

	initializeLogger()
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
			utilities.GetLogger().Warn("missing env GOOGLE_APPLICATION_CREDENTIALS")
		} else if utilities.DefaultGcpProjectID == "" {
			utilities.GetLogger().Warn("missing Default Project ID for GCP")
		} else {
			utilities.GetLogger().Warn("Gcp accounts not found in extension_config. Falling back to ADC\n")
		}
	}

	return nil
}

func readTableConfigurations(homeDir string) {
	var awsConfigFileList = []string{"aws/ec2/table_config.json", "aws/s3/table_config.json"}
	var gcpConfigFileList = []string{"gcp/compute/table_config.json", "gcp/storage/table_config.json"}
	var azureConfigFileList = []string{"azure/compute/table_config.json"}
	var configFileList = append(awsConfigFileList, gcpConfigFileList...)
	configFileList = append(configFileList, azureConfigFileList...)

	for _, fileName := range configFileList {
		fmt.Println("Reading config file:" + homeDir + string(os.PathSeparator) + fileName)
		filePath := homeDir + string(os.PathSeparator) + fileName
		jsonEncoded, err := ioutil.ReadFile(filePath)
		if err != nil {
			fmt.Println("error reading config file:" + homeDir + string(os.PathSeparator) + fileName)
			continue
		}
		readErr := utilities.ReadTableConfig(jsonEncoded)
		if readErr != nil {
			fmt.Println("error parsing json from file:" + homeDir + string(os.PathSeparator) + fileName)
			continue
		}
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
