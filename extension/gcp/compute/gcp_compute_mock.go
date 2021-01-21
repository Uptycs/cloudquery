package compute

import (
	"context"

	"google.golang.org/api/compute/v1"
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

func (gcp *GcpComputeMock) NewImagesService(svc *gcpcompute.Service) *gcpcompute.ImagesService {
	return nil
}

func (gcp *GcpComputeMock) ImagesList(apiSvc *gcpcompute.ImagesService, projectID string) *gcpcompute.ImagesListCall {
	return nil
}

func (gcp *GcpComputeMock) ImagesPages(listCall *gcpcompute.ImagesListCall, ctx context.Context, cb callbackImagesPages) error {
	return nil
}

func (gcp *GcpComputeMock) NewInterconnectsService(svc *gcpcompute.Service) *gcpcompute.InterconnectsService {
	return nil
}

func (gcp *GcpComputeMock) InterconnectsList(apiSvc *gcpcompute.InterconnectsService, projectID string) *gcpcompute.InterconnectsListCall {
	return nil
}

func (gcp *GcpComputeMock) InterconnectsPages(listCall *gcpcompute.InterconnectsListCall, ctx context.Context, cb callbackInterconnectsPages) error {
	return nil
}

func (gcp *GcpComputeMock) NewRoutesService(svc *gcpcompute.Service) *gcpcompute.RoutesService {
	return nil
}

func (gcp *GcpComputeMock) RoutesList(apiSvc *gcpcompute.RoutesService, projectID string) *gcpcompute.RoutesListCall {
	return nil
}

func (gcp *GcpComputeMock) RoutesPages(listCall *gcpcompute.RoutesListCall, ctx context.Context, cb callbackRoutesPages) error {
	return nil
}

func (gcp *GcpComputeMock) NewReservationsService(svc *gcpcompute.Service) *gcpcompute.ReservationsService {
	return nil
}

func (gcp *GcpComputeMock) ReservationsAggregatedList(apiSvc *gcpcompute.ReservationsService, projectID string) *gcpcompute.ReservationsAggregatedListCall {
	return nil
}

func (gcp *GcpComputeMock) ReservationsPages(listCall *gcpcompute.ReservationsAggregatedListCall, ctx context.Context, cb callbackReservationsPages) error {
	return nil
}

func (gcp *GcpComputeMock) NewRoutersService(svc *gcpcompute.Service) *gcpcompute.RoutersService {
	return nil
}

func (gcp *GcpComputeMock) RoutersAggregatedList(apiSvc *gcpcompute.RoutersService, projectID string) *gcpcompute.RoutersAggregatedListCall {
	return nil
}

func (gcp *GcpComputeMock) RoutersPages(listCall *gcpcompute.RoutersAggregatedListCall, ctx context.Context, cb callbackRoutersPages) error {
	return nil
}

func (gcp *GcpComputeMock) NewVpnTunnelsService(svc *gcpcompute.Service) *gcpcompute.VpnTunnelsService {
	return nil
}

func (gcp *GcpComputeMock) VpnTunnelsAggregatedList(apiSvc *gcpcompute.VpnTunnelsService, projectID string) *gcpcompute.VpnTunnelsAggregatedListCall {
	return nil
}

func (gcp *GcpComputeMock) VpnTunnelsPages(listCall *gcpcompute.VpnTunnelsAggregatedListCall, ctx context.Context, cb callbackVpnTunnelsPages) error {
	return nil
}

func (gcp *GcpComputeMock) NewVpnGatewaysService(svc *gcpcompute.Service) *gcpcompute.VpnGatewaysService {
	return nil
}

func (gcp *GcpComputeMock) VpnGatewaysAggregatedList(apiSvc *gcpcompute.VpnGatewaysService, projectID string) *compute.VpnGatewaysAggregatedListCall {
	return nil
}

func (gcp *GcpComputeMock) VpnGatewaysPages(listCall *gcpcompute.VpnGatewaysAggregatedListCall, ctx context.Context, cb callbackVpnGatewaysPages) error {
	return nil
}
