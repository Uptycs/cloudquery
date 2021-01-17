package utilities

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

var (
	TableConfigurationMap = map[string]*TableConfig{}
	AwsAccountId          string
	ExtConfiguration      ExtensionConfiguration
	DefaultGcpProjectID   string
)

func readTableConfig(filePath string) error {
	reader, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	var configurations map[string]*TableConfig
	errUnmarshal := json.Unmarshal(reader, &configurations)
	if errUnmarshal != nil {
		return errUnmarshal
	}
	for tableName, config := range configurations {
		fmt.Println("Found configuration for table:" + tableName)
		config.InitParsedAttributeConfigMap()
		TableConfigurationMap[tableName] = config
		//fmt.Printf("So far Read config for %d tables\n", len(utilities.TableConfigurationMap))
	}
	return nil
}

func ReadTableConfigurations(extensionHomeDir string) {
	var awsConfigFileList = []string{"aws/ec2/table_config.json", "aws/s3/table_config.json"}
	var gcpConfigFileList = []string{"gcp/compute/table_config.json", "gcp/storage/table_config.json"}
	var azureConfigFileList = []string{"azure/compute/table_config.json"}
	var configFileList = append(awsConfigFileList, gcpConfigFileList...)
	configFileList = append(configFileList, azureConfigFileList...)

	for _, fileName := range configFileList {
		fmt.Println("Reading config file:" + extensionHomeDir + string(os.PathSeparator) + fileName)
		readTableConfig(extensionHomeDir + string(os.PathSeparator) + fileName)
	}
	fmt.Printf("Read config for total %d tables\n", len(TableConfigurationMap))
}
