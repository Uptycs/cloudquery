package utilities

import (
	"strings"
)

type ParsedAttributeConfig struct {
	SourceName string `json:"sourceName"`
	TargetName string `json:"targetName"`
	TargetType string `json:"targetType"`
	Enabled    bool   `json:"enabled"`
}

type AwsConfig struct {
	RegionAttribute     string `json:"regionAttribute"`
	RegionCodeAttribute string `json:"regionCodeAttribute"`
	AccountIdAttribute  string `json:"accountIdAttribute"`
}

type GcpConfig struct {
	ProjectIdAttribute string   `json:"projectIdAttribute,omitempty"`
	ZoneAttribute      string   `json:"zoneAttribute,omitempty"`
	Zones              []string `json:"zones"`
}

type TableConfig struct {
	Imports          []string                `json:"imports"`
	MaxLevel         int                     `json:"maxLevel"`
	Api              string                  `json:"api"`
	Paginated        bool                    `json:"paginated"`
	TemplateFile     string                  `json:"templateFile"`
	Aws              AwsConfig               `json:"aws"`
	Gcp              GcpConfig               `json:"gcp"`
	ParsedAttributes []ParsedAttributeConfig `json:"parsedAttributes"`

	parsedAttributeConfigMap map[string]ParsedAttributeConfig
}

func (tableConfig *TableConfig) InitParsedAttributeConfigMap() {
	tableConfig.parsedAttributeConfigMap = make(map[string]ParsedAttributeConfig)
	for _, attr := range tableConfig.ParsedAttributes {
		if attr.Enabled {
			level := strings.Count(attr.SourceName, "_")
			if level > tableConfig.MaxLevel {
				tableConfig.MaxLevel = level
			}
		}
		tableConfig.parsedAttributeConfigMap[attr.SourceName] = attr
	}
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
