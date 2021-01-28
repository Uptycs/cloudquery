/**
 * Copyright (c) 2020-present, The cloudquery authors
 *
 * This source code is licensed as defined by the LICENSE file found in the
 * root directory of this source tree.
 *
 * SPDX-License-Identifier: (Apache-2.0 OR GPL-2.0-only)
 */

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

// GcpComputeInterface abstracts compute APIs accessed by gcp_compute_* tables
type GcpComputeInterface interface {
	NewService(context.Context, ...option.ClientOption) (*compute.Service, error)

	NewDisksService(*compute.Service) *compute.DisksService
	DisksAggregatedList(*compute.DisksService, string) *compute.DisksAggregatedListCall
	DisksPages(context.Context, *compute.DisksAggregatedListCall, callbackDisksPages) error

	NewInstancesService(*compute.Service) *compute.InstancesService
	InstancesAggregatedList(*compute.InstancesService, string) *compute.InstancesAggregatedListCall
	InstancesPages(context.Context, *compute.InstancesAggregatedListCall, callbackInstancesPages) error

	NewReservationsService(*compute.Service) *compute.ReservationsService
	ReservationsAggregatedList(*compute.ReservationsService, string) *compute.ReservationsAggregatedListCall
	ReservationsPages(context.Context, *compute.ReservationsAggregatedListCall, callbackReservationsPages) error

	NewRoutersService(*compute.Service) *compute.RoutersService
	RoutersAggregatedList(*compute.RoutersService, string) *compute.RoutersAggregatedListCall
	RoutersPages(context.Context, *compute.RoutersAggregatedListCall, callbackRoutersPages) error

	NewVpnTunnelsService(*compute.Service) *compute.VpnTunnelsService
	VpnTunnelsAggregatedList(*compute.VpnTunnelsService, string) *compute.VpnTunnelsAggregatedListCall
	VpnTunnelsPages(context.Context, *compute.VpnTunnelsAggregatedListCall, callbackVpnTunnelsPages) error

	NewVpnGatewaysService(*compute.Service) *compute.VpnGatewaysService
	VpnGatewaysAggregatedList(*compute.VpnGatewaysService, string) *compute.VpnGatewaysAggregatedListCall
	VpnGatewaysPages(context.Context, *compute.VpnGatewaysAggregatedListCall, callbackVpnGatewaysPages) error

	NewNetworksService(*compute.Service) *compute.NetworksService
	NetworksList(*compute.NetworksService, string) *compute.NetworksListCall
	NetworksPages(context.Context, *compute.NetworksListCall, callbackNetworksPages) error

	NewImagesService(*compute.Service) *compute.ImagesService
	ImagesList(*compute.ImagesService, string) *compute.ImagesListCall
	ImagesPages(context.Context, *compute.ImagesListCall, callbackImagesPages) error

	NewInterconnectsService(*compute.Service) *compute.InterconnectsService
	InterconnectsList(*compute.InterconnectsService, string) *compute.InterconnectsListCall
	InterconnectsPages(context.Context, *compute.InterconnectsListCall, callbackInterconnectsPages) error

	NewRoutesService(*compute.Service) *compute.RoutesService
	RoutesList(*compute.RoutesService, string) *compute.RoutesListCall
	RoutesPages(context.Context, *compute.RoutesListCall, callbackRoutesPages) error
}

// GcpComputeHandler encloses GcpComputeInterface's instance (mock or otherwise)
type GcpComputeHandler struct {
	svcInterface GcpComputeInterface
}

// NewGcpComputeHandler returns a new instance of GcpComputeHandler with provided intf
func NewGcpComputeHandler(intf GcpComputeInterface) GcpComputeHandler {
	return GcpComputeHandler{svcInterface: intf}
}
