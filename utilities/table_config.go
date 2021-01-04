package utilities

import (
	"encoding/json"
	"io/ioutil"
)

type ParsedAttributeConfig struct {
	SourceName string `json:"sourceName"`
	TargetName string `json:"targetName"`
	TargetType string `json:"targetType"`
}

type AwsConfig struct {
	RegionAttribute     string `json:"regionAttribute"`
	RegionCodeAttribute string `json:"regionCodeAttribute"`
	AccountIdAttribute  string `json:"accountIdAttribute"`
}
type TableConfig struct {
	Imports          []string                `json:"imports"`
	MaxLevel         int                     `json:"maxLevel"`
	Api              string                  `json:"api"`
	Paginated        bool                    `json:"paginated"`
	TemplateFile     string                  `json:"templateFile"`
	Aws              AwsConfig               `json:"aws"`
	ParsedAttributes []ParsedAttributeConfig `json:"parsedAttributes"`

	parsedAttributeConfigMap map[string]ParsedAttributeConfig
}

func ReadTableConfig(configPath string) (map[string]TableConfig, error) {
	reader, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	var result map[string]TableConfig
	errUnmarshal := json.Unmarshal(reader, &result)
	if errUnmarshal != nil {
		return nil, errUnmarshal
	}
	return result, nil
}

func (tableConfig *TableConfig) ParseAwsConfig(configInterface interface{}) {
	config := configInterface.(map[string]interface{})
	if value, ok := config["regionAttribute"]; ok {
		tableConfig.Aws.RegionAttribute = GetStringValue(value)
	}
	if value, ok := config["regionCodeAttribute"]; ok {
		tableConfig.Aws.RegionCodeAttribute = GetStringValue(value)
	}
	if value, ok := config["accountIdAttribute"]; ok {
		tableConfig.Aws.AccountIdAttribute = GetStringValue(value)
	}
}

func (tableConfig *TableConfig) ParseAttributeConfigs(configInterface interface{}) {
	tableConfig.parsedAttributeConfigMap = make(map[string]ParsedAttributeConfig)
	configArr := configInterface.([]interface{})
	for _, attrConfigInterface := range configArr {
		attrConfig := attrConfigInterface.(map[string]interface{})
		parsedAttrConfig := ParsedAttributeConfig{}
		if value, ok := attrConfig["sourceName"]; ok {
			parsedAttrConfig.SourceName = GetStringValue(value)
		} else {
			// Error
			continue
		}
		if value, ok := attrConfig["targetName"]; ok {
			parsedAttrConfig.TargetName = GetStringValue(value)
		} else {
			// Error
			continue
		}
		if value, ok := attrConfig["targetType"]; ok {
			parsedAttrConfig.TargetType = GetStringValue(value)
		} else {
			// Error
			continue
		}
		tableConfig.parsedAttributeConfigMap[parsedAttrConfig.SourceName] = parsedAttrConfig
	}
}

func (tableConfig *TableConfig) Init(config map[string]interface{}) error {
	if value, ok := config["maxLevel"]; ok {
		tableConfig.MaxLevel = GetIntegerValue(value)
	}
	if value, ok := config["aws"]; ok {
		tableConfig.ParseAwsConfig(value)
	}
	if value, ok := config["parsedAttributes"]; ok {
		tableConfig.ParseAttributeConfigs(value)
	}
	return nil
}

func (tableConfig *TableConfig) GetAttributeConfig(attrName string) *ParsedAttributeConfig {
	if value, ok := tableConfig.parsedAttributeConfigMap[attrName]; ok {
		return &value
	} else {
		return nil
	}
}

func (tableConfig *TableConfig) GetParsedAttributeConfigMap() map[string]ParsedAttributeConfig {
	return tableConfig.parsedAttributeConfigMap
}
