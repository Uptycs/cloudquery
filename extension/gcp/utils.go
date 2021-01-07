package gcp

import (
	"context"
	"fmt"
	"github.com/Uptycs/cloudquery/utilities"
	"google.golang.org/api/compute/v1"
	"strings"
)

func FetchRegions() ([]string, error) {
	regions := make([]string, 0)
	// TODO: fetch the list from GCP
	regions = append(regions, "us-east1", "us-west1")
	return regions, nil
}

func GetZones(ctx context.Context, myApiService *compute.InstancesService, projectId string) []string {
	myZonesMap := make(map[string]bool, 0)
	aggListCall := myApiService.AggregatedList(projectId)
	if aggListCall == nil {
		fmt.Println("aggListCall is nil")
		return nil
	}
	if err := aggListCall.Pages(ctx, func(page *compute.InstanceAggregatedList) error {
		for _, item := range page.Items {
			for _, inst := range item.Instances {
				zonePathSplit := strings.Split(inst.Zone, "/")
				myZonesMap[zonePathSplit[len(zonePathSplit) - 1]] = true
			}
		}
		return nil
	}); err != nil {
		fmt.Println("aggListCal.Page: ", err)
		return nil
	}
	myZones := make([]string, 0)
	for k, _ := range myZonesMap {
		myZones = append(myZones, k)
	}
	return myZones
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