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
	ProjectIdAttribute string `json:"projectIdAttribute,omitempty"`
	ZoneAttribute      string `json:"zoneAttribute,omitempty"`
}

type AzureConfig struct {
	SubscriptionIdAttribute string `json:"subscriptionIdAttribute,omitempty"`
	TenantIdAttribute       string `json:"tenantIdAttribute,omitempty"`
	ResourceGroupAttribute  string `json:"resourceGroupAttribute,omitempty"`
}

type TableConfig struct {
	Imports          []string                `json:"imports"`
	MaxLevel         int                     `json:"maxLevel"`
	Api              string                  `json:"api"`
	Paginated        bool                    `json:"paginated"`
	TemplateFile     string                  `json:"templateFile"`
	Aws              AwsConfig               `json:"aws"`
	Gcp              GcpConfig               `json:"gcp"`
	Azure            AzureConfig             `json:"azure"`
	ParsedAttributes []ParsedAttributeConfig `json:"parsedAttributes"`

	parsedAttributeConfigMap map[string]ParsedAttributeConfig
}

func (tableConfig *TableConfig) initParsedAttributeConfigMap() {
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

func (tableConfig *TableConfig) getParsedAttributeConfigMap() map[string]ParsedAttributeConfig {
	return tableConfig.parsedAttributeConfigMap
}
