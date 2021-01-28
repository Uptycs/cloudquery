package utilities

import (
	"strings"
)

// ParsedAttributeConfig represents the attributes for a table
type ParsedAttributeConfig struct {
	SourceName string `json:"sourceName"`
	TargetName string `json:"targetName"`
	TargetType string `json:"targetType"`
	Enabled    bool   `json:"enabled"`
}

// AwsConfig represents the additional attributes for AWS table
type AwsConfig struct {
	RegionAttribute     string `json:"regionAttribute"`
	RegionCodeAttribute string `json:"regionCodeAttribute"`
	AccountIDAttribute  string `json:"accountIdAttribute"`
}

// GcpConfig represents the additional attributes for GCP table
type GcpConfig struct {
	ProjectIDAttribute string `json:"projectIdAttribute,omitempty"`
	ZoneAttribute      string `json:"zoneAttribute,omitempty"`
}

// AzureConfig represents the additional attributes for Azure table
type AzureConfig struct {
	SubscriptionIDAttribute string `json:"subscriptionIdAttribute,omitempty"`
	TenantIDAttribute       string `json:"tenantIdAttribute,omitempty"`
	ResourceGroupAttribute  string `json:"resourceGroupAttribute,omitempty"`
}

// TableConfig represents the configuration of a table
type TableConfig struct {
	Imports          []string                `json:"imports"`
	MaxLevel         int                     `json:"maxLevel"`
	API              string                  `json:"api"`
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
