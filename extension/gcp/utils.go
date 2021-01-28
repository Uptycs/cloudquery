package gcp

import (
	"github.com/Uptycs/cloudquery/utilities"
)

// RowToMap converts JSON row into osquery row
// If configured it will copy some metadata values into appropriate columns
func RowToMap(row map[string]interface{}, projectID string, zone string, tableConfig *utilities.TableConfig) map[string]string {
	result := make(map[string]string)

	if len(tableConfig.Gcp.ProjectIDAttribute) != 0 {
		result[tableConfig.Gcp.ProjectIDAttribute] = projectID
	}
	if len(tableConfig.Gcp.ZoneAttribute) != 0 {
		result[tableConfig.Gcp.ZoneAttribute] = zone
	}

	result = utilities.RowToMap(result, row, tableConfig)
	return result
}
