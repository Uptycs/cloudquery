package gcp

import (
	"github.com/Uptycs/cloudquery/utilities"
)

// RowToMap converts JSON row into osquery row
func RowToMap(row map[string]interface{}, projectID string, zone string, tableConfig *utilities.TableConfig) map[string]string {
	result := make(map[string]string)

	if len(tableConfig.Gcp.ProjectIdAttribute) != 0 {
		result[tableConfig.Gcp.ProjectIdAttribute] = projectID
	}
	if len(tableConfig.Gcp.ZoneAttribute) != 0 {
		result[tableConfig.Gcp.ZoneAttribute] = zone
	}
	for key, value := range tableConfig.GetParsedAttributeConfigMap() {
		if row[key] != nil {
			result[value.TargetName] = utilities.GetStringValue(row[key])
		}
	}
	return result
}
