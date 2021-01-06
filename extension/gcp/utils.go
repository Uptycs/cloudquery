package gcp

import "github.com/Uptycs/cloudquery/utilities"

func FetchRegions() ([]string, error) {
	regions := make([]string, 0)
	// TODO: fetch the list from GCP
	regions = append(regions, "us-east1", "us-west1")
	return regions, nil
}

func RowToMap(row map[string]interface{}, projectId string, zone string, tableConfig *utilities.TableConfig) map[string]string {
	result := make(map[string]string)

	if len(tableConfig.Gcp.ProjectIdAttribute) != 0 {
		result[tableConfig.Gcp.ProjectIdAttribute] = projectId
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