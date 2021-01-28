package utilities

import (
	"encoding/json"

	"fmt"

	log "github.com/sirupsen/logrus"
)

var (
	// TableConfigurationMap is the map of tableName->TableConfig
	TableConfigurationMap = map[string]*TableConfig{}
	AwsAccountId          string
	ExtConfiguration      ExtensionConfiguration
	DefaultGcpProjectID   string
)

// ReadTableConfig parses json encoded data to read list TableConfig entries
// These are available for reading from utilities.TableConfigurationMap[]
func ReadTableConfig(jsonEncoded []byte) error {
	var configurations map[string]*TableConfig
	errUnmarshal := json.Unmarshal(jsonEncoded, &configurations)
	if errUnmarshal != nil {
		return errUnmarshal
	}
	for tableName, config := range configurations {
		GetLogger().WithFields(log.Fields{
			"tableName": tableName,
		}).Debug("found table configuration")
		for _, attr := range config.ParsedAttributes {
			if attr.SourceName == "" || attr.TargetName == "" || attr.TargetType == "" {
				return fmt.Errorf("invalid parsedAttribute entry: %+v", attr)
			}
		}
		config.initParsedAttributeConfigMap()
		TableConfigurationMap[tableName] = config
	}
	return nil
}

// RowToMap converts JSON row into osquery row
func RowToMap(inMap map[string]string, row map[string]interface{}, tableConfig *TableConfig) map[string]string {
	for key, value := range tableConfig.getParsedAttributeConfigMap() {
		if row[key] != nil {
			inMap[value.TargetName] = getStringValue(row[key])
		}
	}
	return inMap
}
