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

func (gcp *GcpComputeImpl) NewNetworksService(svc *compute.Service) *compute.NetworksService {
	return compute.NewNetworksService(svc)
}

func (gcp *GcpComputeImpl) NetworksList(apiSvc *compute.NetworksService, projectID string) *compute.NetworksListCall {
	return apiSvc.List(projectID)
}

func (gcp *GcpComputeImpl) NetworksPages(listCall *compute.NetworksListCall, ctx context.Context, cb callbackNetworksPages) error {
	return listCall.Pages(ctx, cb)
}
