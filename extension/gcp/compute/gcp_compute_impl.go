package compute

import (
	"context"

	compute "google.golang.org/api/compute/v1"
	"google.golang.org/api/option"
)

type GcpComputeImpl struct {
}

func NewGcpComputeImpl() *GcpComputeImpl {
	return &GcpComputeImpl{}
}

func (gcp *GcpComputeImpl) NewService(ctx context.Context, opts ...option.ClientOption) (*compute.Service, error) {
	return compute.NewService(ctx, opts...)
}

func (gcp *GcpComputeImpl) NewDisksService(svc *compute.Service) *compute.DisksService {
	return compute.NewDisksService(svc)
}

func (gcp *GcpComputeImpl) DisksAggregatedList(apiSvc *compute.DisksService, projectID string) *compute.DisksAggregatedListCall {
	return apiSvc.AggregatedList(projectID)
}

func (gcp *GcpComputeImpl) DisksPages(listCall *compute.DisksAggregatedListCall, ctx context.Context, cb callbackDisksPages) error {
	return listCall.Pages(ctx, cb)
}

func (gcp *GcpComputeImpl) NewInstancesService(svc *compute.Service) *compute.InstancesService {
	return compute.NewInstancesService(svc)
}

func (gcp *GcpComputeImpl) InstancesAggregatedList(apiSvc *compute.InstancesService, projectID string) *compute.InstancesAggregatedListCall {
	return apiSvc.AggregatedList(projectID)
}

func (gcp *GcpComputeImpl) InstancesPages(listCall *compute.InstancesAggregatedListCall, ctx context.Context, cb callbackInstancesPages) error {
	return listCall.Pages(ctx, cb)
}

func (gcp *GcpComputeImpl) NewReservationsService(svc *compute.Service) *compute.ReservationsService {
	return compute.NewReservationsService(svc)
}

func (gcp *GcpComputeImpl) ReservationsAggregatedList(apiSvc *compute.ReservationsService, projectID string) *compute.ReservationsAggregatedListCall {
	return apiSvc.AggregatedList(projectID)
}

func (gcp *GcpComputeImpl) ReservationsPages(listCall *compute.ReservationsAggregatedListCall, ctx context.Context, cb callbackReservationsPages) error {
	return listCall.Pages(ctx, cb)
}

func (gcp *GcpComputeImpl) NewRoutersService(svc *compute.Service) *compute.RoutersService {
	return compute.NewRoutersService(svc)
}

func (gcp *GcpComputeImpl) RoutersAggregatedList(apiSvc *compute.RoutersService, projectID string) *compute.RoutersAggregatedListCall {
	return apiSvc.AggregatedList(projectID)
}

func (gcp *GcpComputeImpl) RoutersPages(listCall *compute.RoutersAggregatedListCall, ctx context.Context, cb callbackRoutersPages) error {
	return listCall.Pages(ctx, cb)
}

func (gcp *GcpComputeImpl) NewVpnTunnelsService(svc *compute.Service) *compute.VpnTunnelsService {
	return compute.NewVpnTunnelsService(svc)
}

func (gcp *GcpComputeImpl) VpnTunnelsAggregatedList(apiSvc *compute.VpnTunnelsService, projectID string) *compute.VpnTunnelsAggregatedListCall {
	return apiSvc.AggregatedList(projectID)
}

func (gcp *GcpComputeImpl) VpnTunnelsPages(listCall *compute.VpnTunnelsAggregatedListCall, ctx context.Context, cb callbackVpnTunnelsPages) error {
	return listCall.Pages(ctx, cb)
}

func (gcp *GcpComputeImpl) NewVpnGatewaysService(svc *compute.Service) *compute.VpnGatewaysService {
	return compute.NewVpnGatewaysService(svc)
}

func (gcp *GcpComputeImpl) VpnGatewaysAggregatedList(apiSvc *compute.VpnGatewaysService, projectID string) *compute.VpnGatewaysAggregatedListCall {
	return apiSvc.AggregatedList(projectID)
}

func (gcp *GcpComputeImpl) VpnGatewaysPages(listCall *compute.VpnGatewaysAggregatedListCall, ctx context.Context, cb callbackVpnGatewaysPages) error {
	return listCall.Pages(ctx, cb)
}

func (gcp *GcpComputeImpl) NewNetworksService(svc *compute.Service) *compute.NetworksService {
	return compute.NewNetworksService(svc)
}

func (gcp *GcpComputeImpl) NetworksList(apiSvc *compute.NetworksService, projectID string) *compute.NetworksListCall {
	return apiSvc.List(projectID)
}

func (gcp *GcpComputeImpl) NetworksPages(listCall *compute.NetworksListCall, ctx context.Context, cb callbackNetworksPages) error {
	return listCall.Pages(ctx, cb)
}
func (gcp *GcpComputeImpl) NewImagesService(svc *compute.Service) *compute.ImagesService {
	return compute.NewImagesService(svc)
}

func (gcp *GcpComputeImpl) ImagesList(apiSvc *compute.ImagesService, projectID string) *compute.ImagesListCall {
	return apiSvc.List(projectID)
}

func (gcp *GcpComputeImpl) ImagesPages(listCall *compute.ImagesListCall, ctx context.Context, cb callbackImagesPages) error {
	return listCall.Pages(ctx, cb)
}
func (gcp *GcpComputeImpl) NewInterconnectsService(svc *compute.Service) *compute.InterconnectsService {
	return compute.NewInterconnectsService(svc)
}

func (gcp *GcpComputeImpl) InterconnectsList(apiSvc *compute.InterconnectsService, projectID string) *compute.InterconnectsListCall {
	return apiSvc.List(projectID)
}

func (gcp *GcpComputeImpl) InterconnectsPages(listCall *compute.InterconnectsListCall, ctx context.Context, cb callbackInterconnectsPages) error {
	return listCall.Pages(ctx, cb)
}
func (gcp *GcpComputeImpl) NewRoutesService(svc *compute.Service) *compute.RoutesService {
	return compute.NewRoutesService(svc)
}

func (gcp *GcpComputeImpl) RoutesList(apiSvc *compute.RoutesService, projectID string) *compute.RoutesListCall {
	return apiSvc.List(projectID)
}

func (gcp *GcpComputeImpl) RoutesPages(listCall *compute.RoutesListCall, ctx context.Context, cb callbackRoutesPages) error {
	return listCall.Pages(ctx, cb)
}
