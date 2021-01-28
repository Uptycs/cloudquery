package compute

import (
	"context"

	"google.golang.org/api/compute/v1"
	gcpcompute "google.golang.org/api/compute/v1"
	"google.golang.org/api/option"
)

// GcpComputeMock routes gcp_compute_* table's API invocation to a mock implementation to enable unit tests
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

// NewGcpComputeMock returns new instance of GcpComputeMock
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

// NewService returns compute.Service or error
func (gcp *GcpComputeMock) NewService(ctx context.Context, opts ...option.ClientOption) (*gcpcompute.Service, error) {
	return &gcp.svc, nil
}

// NewDisksService returns *compute.DisksService
func (gcp *GcpComputeMock) NewDisksService(svc *gcpcompute.Service) *gcpcompute.DisksService {
	return &gcp.disksSvc
}

// DisksAggregatedList returns *compute.DisksAggregatedListCall
func (gcp *GcpComputeMock) DisksAggregatedList(apiSvc *gcpcompute.DisksService, projectID string) *gcpcompute.DisksAggregatedListCall {
	return &gcp.disksAggList
}

// DisksPages invokes cb for each page of results.
// Returns error on failure
func (gcp *GcpComputeMock) DisksPages(ctx context.Context, listCall *gcpcompute.DisksAggregatedListCall, cb callbackDisksPages) error {
	cb(&gcp.disksPage)
	return nil
}

func (gcp *GcpComputeMock) addDisks(inList []*gcpcompute.Disk) {
	disks := gcp.disksPage.Items[gcp.itemsKey]
	disks.Disks = append(disks.Disks, inList...)
	gcp.disksPage.Items[gcp.itemsKey] = disks
}

func (gcp *GcpComputeMock) clearDisks() {
	disks := gcp.disksPage.Items[gcp.itemsKey]
	disks.Disks = make([]*gcpcompute.Disk, 0)
	gcp.disksPage.Items[gcp.itemsKey] = disks
}

// NewInstancesService returns *compute.InstancesService
func (gcp *GcpComputeMock) NewInstancesService(svc *gcpcompute.Service) *gcpcompute.InstancesService {
	return &gcp.instancesSvc
}

// InstancesAggregatedList returns *compute.InstancesAggregatedListCall
func (gcp *GcpComputeMock) InstancesAggregatedList(apiSvc *gcpcompute.InstancesService, projectID string) *gcpcompute.InstancesAggregatedListCall {
	return &gcp.instancesAggList
}

// InstancesPages invokes cb for each page of results.
// Returns error on failure
func (gcp *GcpComputeMock) InstancesPages(ctx context.Context, listCall *gcpcompute.InstancesAggregatedListCall, cb callbackInstancesPages) error {
	cb(&gcp.instancesPage)
	return nil
}

func (gcp *GcpComputeMock) addInstances(inList []*gcpcompute.Instance) {
	instances := gcp.instancesPage.Items[gcp.itemsKey]
	instances.Instances = append(instances.Instances, inList...)
	gcp.instancesPage.Items[gcp.itemsKey] = instances
}

func (gcp *GcpComputeMock) clearInstances() {
	instances := gcp.instancesPage.Items[gcp.itemsKey]
	instances.Instances = make([]*gcpcompute.Instance, 0)
	gcp.instancesPage.Items[gcp.itemsKey] = instances
}

// NewNetworksService returns *compute.NetworksService
func (gcp *GcpComputeMock) NewNetworksService(svc *gcpcompute.Service) *gcpcompute.NetworksService {
	return &gcp.networksSvc
}

// NetworksList returns *compute.NetworksListCall
func (gcp *GcpComputeMock) NetworksList(apiSvc *gcpcompute.NetworksService, projectID string) *gcpcompute.NetworksListCall {
	return &gcp.networksList
}

// NetworksPages invokes cb for each page of results.
// Returns error on failure
func (gcp *GcpComputeMock) NetworksPages(ctx context.Context, listCall *gcpcompute.NetworksListCall, cb callbackNetworksPages) error {
	cb(&gcp.networksPage)
	return nil
}

func (gcp *GcpComputeMock) addNetworks(inList []*gcpcompute.Network) {
	gcp.networksPage.Items = inList
}

func (gcp *GcpComputeMock) clearNetworks() {
	gcp.networksPage.Items = make([]*gcpcompute.Network, 0)
}

// NewImagesService returns *compute.ImagesService
func (gcp *GcpComputeMock) NewImagesService(svc *gcpcompute.Service) *gcpcompute.ImagesService {
	return nil
}

// ImagesList returns *compute.ImagesListCall
func (gcp *GcpComputeMock) ImagesList(apiSvc *gcpcompute.ImagesService, projectID string) *gcpcompute.ImagesListCall {
	return nil
}

// ImagesPages invokes cb for each page of results.
// Returns error on failure
func (gcp *GcpComputeMock) ImagesPages(ctx context.Context, listCall *gcpcompute.ImagesListCall, cb callbackImagesPages) error {
	return nil
}

// NewInterconnectsService returns *compute.InterconnectsService
func (gcp *GcpComputeMock) NewInterconnectsService(svc *gcpcompute.Service) *gcpcompute.InterconnectsService {
	return nil
}

// InterconnectsList returns *compute.InterconnectsListCall
func (gcp *GcpComputeMock) InterconnectsList(apiSvc *gcpcompute.InterconnectsService, projectID string) *gcpcompute.InterconnectsListCall {
	return nil
}

// InterconnectsPages invokes cb for each page of results.
// Returns error on failure
func (gcp *GcpComputeMock) InterconnectsPages(ctx context.Context, listCall *gcpcompute.InterconnectsListCall, cb callbackInterconnectsPages) error {
	return nil
}

// NewRoutesService returns *compute.RoutesService
func (gcp *GcpComputeMock) NewRoutesService(svc *gcpcompute.Service) *gcpcompute.RoutesService {
	return nil
}

// RoutesList returns *compute.RoutesListCall
func (gcp *GcpComputeMock) RoutesList(apiSvc *gcpcompute.RoutesService, projectID string) *gcpcompute.RoutesListCall {
	return nil
}

// RoutesPages invokes cb for each page of results.
// Returns error on failure
func (gcp *GcpComputeMock) RoutesPages(ctx context.Context, listCall *gcpcompute.RoutesListCall, cb callbackRoutesPages) error {
	return nil
}

// NewReservationsService returns *compute.ReservationsService
func (gcp *GcpComputeMock) NewReservationsService(svc *gcpcompute.Service) *gcpcompute.ReservationsService {
	return nil
}

// ReservationsAggregatedList returns *compute.ReservationsAggregatedListCall
func (gcp *GcpComputeMock) ReservationsAggregatedList(apiSvc *gcpcompute.ReservationsService, projectID string) *gcpcompute.ReservationsAggregatedListCall {
	return nil
}

// ReservationsPages invokes cb for each page of results.
// Returns error on failure
func (gcp *GcpComputeMock) ReservationsPages(ctx context.Context, listCall *gcpcompute.ReservationsAggregatedListCall, cb callbackReservationsPages) error {
	return nil
}

// NewRoutersService returns *compute.RoutersService
func (gcp *GcpComputeMock) NewRoutersService(svc *gcpcompute.Service) *gcpcompute.RoutersService {
	return nil
}

// RoutersAggregatedList returns *compute.RoutersAggregatedListCall
func (gcp *GcpComputeMock) RoutersAggregatedList(apiSvc *gcpcompute.RoutersService, projectID string) *gcpcompute.RoutersAggregatedListCall {
	return nil
}

// RoutersPages invokes cb for each page of results.
// Returns error on failure
func (gcp *GcpComputeMock) RoutersPages(ctx context.Context, listCall *gcpcompute.RoutersAggregatedListCall, cb callbackRoutersPages) error {
	return nil
}

// NewVpnTunnelsService returns *compute.VpnTunnelsService
func (gcp *GcpComputeMock) NewVpnTunnelsService(svc *gcpcompute.Service) *gcpcompute.VpnTunnelsService {
	return nil
}

// VpnTunnelsAggregatedList returns *compute.VpnTunnelsAggregatedListCall
func (gcp *GcpComputeMock) VpnTunnelsAggregatedList(apiSvc *gcpcompute.VpnTunnelsService, projectID string) *gcpcompute.VpnTunnelsAggregatedListCall {
	return nil
}

// VpnTunnelsPages invokes cb for each page of results.
// Returns error on failure
func (gcp *GcpComputeMock) VpnTunnelsPages(ctx context.Context, listCall *gcpcompute.VpnTunnelsAggregatedListCall, cb callbackVpnTunnelsPages) error {
	return nil
}

// NewVpnGatewaysService returns *compute.VpnGatewaysService
func (gcp *GcpComputeMock) NewVpnGatewaysService(svc *gcpcompute.Service) *gcpcompute.VpnGatewaysService {
	return nil
}

// VpnGatewaysAggregatedList returns *compute.VpnGatewaysAggregatedListCall
func (gcp *GcpComputeMock) VpnGatewaysAggregatedList(apiSvc *gcpcompute.VpnGatewaysService, projectID string) *compute.VpnGatewaysAggregatedListCall {
	return nil
}

// VpnGatewaysPages invokes cb for each page of results.
// Returns error on failure
func (gcp *GcpComputeMock) VpnGatewaysPages(ctx context.Context, listCall *gcpcompute.VpnGatewaysAggregatedListCall, cb callbackVpnGatewaysPages) error {
	return nil
}
