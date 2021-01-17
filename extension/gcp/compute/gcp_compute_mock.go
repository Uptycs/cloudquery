package compute

import (
	"context"

	gcpcompute "google.golang.org/api/compute/v1"
	"google.golang.org/api/option"
)

type GcpComputeMock struct {
	svc gcpcompute.Service

	disksSvc         gcpcompute.DisksService
	disksAggList     gcpcompute.DisksAggregatedListCall
	instancesSvc     gcpcompute.InstancesService
	instancesAggList gcpcompute.InstancesAggregatedListCall
	networksSvc      gcpcompute.NetworksService
	networksList     gcpcompute.NetworksListCall

	instancesPage gcpcompute.InstanceAggregatedList
	disksPage     gcpcompute.DiskAggregatedList
	networksPage  gcpcompute.NetworkList

	itemsKey string
}

func NewGcpComputeMock() *GcpComputeMock {
	var mock = GcpComputeMock{}

	mock.itemsKey = "test"

	diskItems := make(map[string]gcpcompute.DisksScopedList)
	diskItems[mock.itemsKey] = gcpcompute.DisksScopedList{Disks: make([]*gcpcompute.Disk, 0)}
	mock.disksPage = gcpcompute.DiskAggregatedList{Items: diskItems}

	instanceItems := make(map[string]gcpcompute.InstancesScopedList)
	instanceItems[mock.itemsKey] = gcpcompute.InstancesScopedList{Instances: make([]*gcpcompute.Instance, 0)}
	mock.instancesPage = gcpcompute.InstanceAggregatedList{Items: instanceItems}

	networkItems := make([]*gcpcompute.Network, 0)
	mock.networksPage = gcpcompute.NetworkList{Items: networkItems}

	return &mock
}

func (gcp *GcpComputeMock) NewService(ctx context.Context, opts ...option.ClientOption) (*gcpcompute.Service, error) {
	return &gcp.svc, nil
}

func (gcp *GcpComputeMock) NewDisksService(svc *gcpcompute.Service) *gcpcompute.DisksService {
	return &gcp.disksSvc
}

func (gcp *GcpComputeMock) DisksAggregatedList(apiSvc *gcpcompute.DisksService, projectID string) *gcpcompute.DisksAggregatedListCall {
	return &gcp.disksAggList
}

func (gcp *GcpComputeMock) DisksPages(listCall *gcpcompute.DisksAggregatedListCall, ctx context.Context, cb callbackDisksPages) error {
	cb(&gcp.disksPage)
	return nil
}

func (gcp *GcpComputeMock) AddDisks(inList []*gcpcompute.Disk) {
	disks := gcp.disksPage.Items[gcp.itemsKey]
	disks.Disks = append(disks.Disks, inList...)
	gcp.disksPage.Items[gcp.itemsKey] = disks
}

func (gcp *GcpComputeMock) ClearDisks() {
	disks := gcp.disksPage.Items[gcp.itemsKey]
	disks.Disks = make([]*gcpcompute.Disk, 0)
	gcp.disksPage.Items[gcp.itemsKey] = disks
}

func (gcp *GcpComputeMock) NewInstancesService(svc *gcpcompute.Service) *gcpcompute.InstancesService {
	return &gcp.instancesSvc
}

func (gcp *GcpComputeMock) InstancesAggregatedList(apiSvc *gcpcompute.InstancesService, projectID string) *gcpcompute.InstancesAggregatedListCall {
	return &gcp.instancesAggList
}

func (gcp *GcpComputeMock) InstancesPages(listCall *gcpcompute.InstancesAggregatedListCall, ctx context.Context, cb callbackInstancesPages) error {
	cb(&gcp.instancesPage)
	return nil
}

func (gcp *GcpComputeMock) AddInstances(inList []*gcpcompute.Instance) {
	instances := gcp.instancesPage.Items[gcp.itemsKey]
	instances.Instances = append(instances.Instances, inList...)
	gcp.instancesPage.Items[gcp.itemsKey] = instances
}

func (gcp *GcpComputeMock) ClearInstances() {
	instances := gcp.instancesPage.Items[gcp.itemsKey]
	instances.Instances = make([]*gcpcompute.Instance, 0)
	gcp.instancesPage.Items[gcp.itemsKey] = instances
}

func (gcp *GcpComputeMock) NewNetworksService(svc *gcpcompute.Service) *gcpcompute.NetworksService {
	return &gcp.networksSvc
}

func (gcp *GcpComputeMock) NetworksList(apiSvc *gcpcompute.NetworksService, projectID string) *gcpcompute.NetworksListCall {
	return &gcp.networksList
}

func (gcp *GcpComputeMock) NetworksPages(listCall *gcpcompute.NetworksListCall, ctx context.Context, cb callbackNetworksPages) error {
	cb(&gcp.networksPage)
	return nil
}

func (gcp *GcpComputeMock) AddNetworks(inList []*gcpcompute.Network) {
	gcp.networksPage.Items = inList
}

func (gcp *GcpComputeMock) ClearNetworks() {
	gcp.networksPage.Items = make([]*gcpcompute.Network, 0)
}
