package compute

import (
	"context"

	"google.golang.org/api/compute/v1"
	"google.golang.org/api/option"
)

type callbackDisksPages func(*compute.DiskAggregatedList) error
type callbackInstancesPages func(*compute.InstanceAggregatedList) error
type callbackReservationsPages func(*compute.ReservationAggregatedList) error
type callbackRoutersPages func(*compute.RouterAggregatedList) error
type callbackVpnTunnelsPages func(*compute.VpnTunnelAggregatedList) error
type callbackVpnGatewaysPages func(*compute.VpnGatewayAggregatedList) error

type callbackNetworksPages func(*compute.NetworkList) error
type callbackImagesPages func(*compute.ImageList) error
type callbackInterconnectsPages func(*compute.InterconnectList) error
type callbackRoutesPages func(*compute.RouteList) error

type GcpComputeInterface interface {
	NewService(context.Context, ...option.ClientOption) (*compute.Service, error)

	NewDisksService(*compute.Service) *compute.DisksService
	DisksAggregatedList(*compute.DisksService, string) *compute.DisksAggregatedListCall
	DisksPages(*compute.DisksAggregatedListCall, context.Context, callbackDisksPages) error

	NewInstancesService(*compute.Service) *compute.InstancesService
	InstancesAggregatedList(*compute.InstancesService, string) *compute.InstancesAggregatedListCall
	InstancesPages(*compute.InstancesAggregatedListCall, context.Context, callbackInstancesPages) error

	NewReservationsService(*compute.Service) *compute.ReservationsService
	ReservationsAggregatedList(*compute.ReservationsService, string) *compute.ReservationsAggregatedListCall
	ReservationsPages(*compute.ReservationsAggregatedListCall, context.Context, callbackReservationsPages) error

	NewRoutersService(*compute.Service) *compute.RoutersService
	RoutersAggregatedList(*compute.RoutersService, string) *compute.RoutersAggregatedListCall
	RoutersPages(*compute.RoutersAggregatedListCall, context.Context, callbackRoutersPages) error

	NewVpnTunnelsService(*compute.Service) *compute.VpnTunnelsService
	VpnTunnelsAggregatedList(*compute.VpnTunnelsService, string) *compute.VpnTunnelsAggregatedListCall
	VpnTunnelsPages(*compute.VpnTunnelsAggregatedListCall, context.Context, callbackVpnTunnelsPages) error

	NewVpnGatewaysService(*compute.Service) *compute.VpnGatewaysService
	VpnGatewaysAggregatedList(*compute.VpnGatewaysService, string) *compute.VpnGatewaysAggregatedListCall
	VpnGatewaysPages(*compute.VpnGatewaysAggregatedListCall, context.Context, callbackVpnGatewaysPages) error

	NewNetworksService(*compute.Service) *compute.NetworksService
	NetworksList(*compute.NetworksService, string) *compute.NetworksListCall
	NetworksPages(*compute.NetworksListCall, context.Context, callbackNetworksPages) error

	NewImagesService(*compute.Service) *compute.ImagesService
	ImagesList(*compute.ImagesService, string) *compute.ImagesListCall
	ImagesPages(*compute.ImagesListCall, context.Context, callbackImagesPages) error

	NewInterconnectsService(*compute.Service) *compute.InterconnectsService
	InterconnectsList(*compute.InterconnectsService, string) *compute.InterconnectsListCall
	InterconnectsPages(*compute.InterconnectsListCall, context.Context, callbackInterconnectsPages) error

	NewRoutesService(*compute.Service) *compute.RoutesService
	RoutesList(*compute.RoutesService, string) *compute.RoutesListCall
	RoutesPages(*compute.RoutesListCall, context.Context, callbackRoutesPages) error
}

type GcpComputeHandler struct {
	svcInterface GcpComputeInterface
}

func NewGcpComputeHandler(intf GcpComputeInterface) GcpComputeHandler {
	return GcpComputeHandler{svcInterface: intf}
}
