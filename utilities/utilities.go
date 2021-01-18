package utilities

import (
	"encoding/json"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
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
		GetLogger().WithFields(log.Fields{
			"tableName": tableName,
		}).Debug("found table configuration")

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
		GetLogger().WithFields(log.Fields{
			"fileName": extensionHomeDir + string(os.PathSeparator) + fileName,
		}).Debug("reading config file")

		readTableConfig(extensionHomeDir + string(os.PathSeparator) + fileName)
	}
	GetLogger().Info("read table configurations for ", len(TableConfigurationMap), " tables")
}
