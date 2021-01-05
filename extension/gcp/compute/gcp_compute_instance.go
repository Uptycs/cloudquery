package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/kolide/osquery-go/plugin/table"
	"google.golang.org/api/option"

	compute "google.golang.org/api/compute/v1"
)

func GcpComputeInstanceColumns() []table.ColumnDefinition {
	var _, _ = strconv.Atoi("123") // Disables warning when strcov is not used
	return []table.ColumnDefinition{
		table.IntegerColumn("id"),
		table.TextColumn("hostname"),
		table.TextColumn("name"),
		table.TextColumn("kind"),
	}
}

func GcpComputeInstanceGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	var _ = queryContext
	resultMap := make([]map[string]string, 0)
	ctx, cancel := context.WithCancel(osqCtx)
	defer cancel()
	service, err := compute.NewService(ctx, option.WithCredentialsFile(*keyFile))
	if err != nil {
		fmt.Println("NewService() error: ", err)
		return resultMap, err
	}
	myApiService := compute.NewInstancesService(service)
	if myApiService == nil {
		fmt.Println("compute.NewInstancesService() returned nil")
		return resultMap, fmt.Errorf("compute.NewInstancesService() returned nil")
	}

	listCall := myApiService.List(*projectId, *zone)
	if listCall == nil {
		fmt.Println("listCall is nil")
		return resultMap, nil
	}
	if err := listCall.Pages(ctx, func(page *compute.InstanceList) error {
		for _, item := range page.Items {
			result := make(map[string]string)
			result["id"] = strconv.FormatUint(item.Id, 10)
			result["hostname"] = item.Hostname
			result["name"] = item.Name
			result["kind"] = item.Kind
			resultMap = append(resultMap, result)
		}
		return nil
	}); err != nil {
		fmt.Println("listCall.Page: ", err)
		//log.Fatal(err)
		return resultMap, err
	}
	return resultMap, nil
}
