package compute

import (
	"context"
	"fmt"

	compute "google.golang.org/api/compute/v1"
	"google.golang.org/api/option"
)

type GcpComputeMock struct {
	svc compute.Service

	disksSvc         compute.DisksService
	disksAggList     compute.DisksAggregatedListCall
	instancesSvc     compute.InstancesService
	instancesAggList compute.InstancesAggregatedListCall
	networksSvc      compute.NetworksService
	networksList     compute.NetworksListCall

	instancesPage compute.InstanceAggregatedList
}

func NewGcpComputeMock() *GcpComputeMock {
	var mock = GcpComputeMock{}

	instances := make([]*compute.Instance, 0)
	inst1 := compute.Instance{Name: "MockInstance1"}
	instances = append(instances, &inst1)

	instanceItems := make(map[string]compute.InstancesScopedList)
	instanceItems["test"] = compute.InstancesScopedList{Instances: instances}
	mock.instancesPage = compute.InstanceAggregatedList{Items: instanceItems}
	return &mock
}

func (gcp *GcpComputeMock) NewService(ctx context.Context, opts ...option.ClientOption) (*compute.Service, error) {
	return &gcp.svc, nil
}

func (gcp *GcpComputeMock) NewDisksService(svc *compute.Service) *compute.DisksService {
	return &gcp.disksSvc
}

func (gcp *GcpComputeMock) DisksAggregatedList(apiSvc *compute.DisksService, projectID string) *compute.DisksAggregatedListCall {
	return &gcp.disksAggList
}

func (gcp *GcpComputeMock) DisksPages(listCall *compute.DisksAggregatedListCall, ctx context.Context, cb callbackDisksPages) error {
	fmt.Printf("No support yet!")
	return nil
}

func (gcp *GcpComputeMock) NewInstancesService(svc *compute.Service) *compute.InstancesService {
	return &gcp.instancesSvc
}

func (gcp *GcpComputeMock) InstancesAggregatedList(apiSvc *compute.InstancesService, projectID string) *compute.InstancesAggregatedListCall {
	return &gcp.instancesAggList
}

func (gcp *GcpComputeMock) InstancesPages(listCall *compute.InstancesAggregatedListCall, ctx context.Context, cb callbackInstancesPages) error {
	cb(&gcp.instancesPage)
	return nil
}

func (gcp *GcpComputeMock) NewNetworksService(svc *compute.Service) *compute.NetworksService {
	return &gcp.networksSvc
}

func (gcp *GcpComputeMock) NetworksList(apiSvc *compute.NetworksService, projectID string) *compute.NetworksListCall {
	return &gcp.networksList
}

func (gcp *GcpComputeMock) NetworksPages(listCall *compute.NetworksListCall, ctx context.Context, cb callbackNetworksPages) error {
	fmt.Printf("No support yet!")
	return nil
}
