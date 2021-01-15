package compute

import (
	"context"

	"google.golang.org/api/compute/v1"
	"google.golang.org/api/option"
)

type callbackDisksPages func(*compute.DiskAggregatedList) error
type callbackInstancesPages func(*compute.InstanceAggregatedList) error

type callbackNetworksPages func(*compute.NetworkList) error

type GcpComputeInterface interface {
	NewService(context.Context, ...option.ClientOption) (*compute.Service, error)

	NewDisksService(*compute.Service) *compute.DisksService
	DisksAggregatedList(*compute.DisksService, string) *compute.DisksAggregatedListCall
	DisksPages(*compute.DisksAggregatedListCall, context.Context, callbackDisksPages) error

	NewInstancesService(*compute.Service) *compute.InstancesService
	InstancesAggregatedList(*compute.InstancesService, string) *compute.InstancesAggregatedListCall
	InstancesPages(*compute.InstancesAggregatedListCall, context.Context, callbackInstancesPages) error

	NewNetworksService(*compute.Service) *compute.NetworksService
	NetworksList(*compute.NetworksService, string) *compute.NetworksListCall
	NetworksPages(*compute.NetworksListCall, context.Context, callbackNetworksPages) error
}

type GcpComputeHandler struct {
	svcInterface GcpComputeInterface
}

func NewGcpComputeHandler(intf GcpComputeInterface) GcpComputeHandler {
	return GcpComputeHandler{svcInterface: intf}
}
