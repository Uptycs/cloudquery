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

	compute "google.golang.org/api/compute/v1"
	"google.golang.org/api/option"
)

// GcpComputeImpl routes gcp_compute_* table's API invocation to appropriate google cloud's API
type GcpComputeImpl struct {
}

// NewGcpComputeImpl returns new instance of GcpComputeImpl
func NewGcpComputeImpl() *GcpComputeImpl {
	return &GcpComputeImpl{}
}

// NewService returns compute.Service or error
func (gcp *GcpComputeImpl) NewService(ctx context.Context, opts ...option.ClientOption) (*compute.Service, error) {
	return compute.NewService(ctx, opts...)
}

// NewDisksService returns *compute.DisksService
func (gcp *GcpComputeImpl) NewDisksService(svc *compute.Service) *compute.DisksService {
	return compute.NewDisksService(svc)
}

// DisksAggregatedList returns *compute.DisksAggregatedListCall
func (gcp *GcpComputeImpl) DisksAggregatedList(apiSvc *compute.DisksService, projectID string) *compute.DisksAggregatedListCall {
	return apiSvc.AggregatedList(projectID)
}

// DisksPages invokes cb for each page of results.
// Returns error on failure
func (gcp *GcpComputeImpl) DisksPages(ctx context.Context, listCall *compute.DisksAggregatedListCall, cb callbackDisksPages) error {
	return listCall.Pages(ctx, cb)
}

// NewInstancesService returns *compute.InstancesService
func (gcp *GcpComputeImpl) NewInstancesService(svc *compute.Service) *compute.InstancesService {
	return compute.NewInstancesService(svc)
}

// InstancesAggregatedList returns *compute.InstancesAggregatedListCall
func (gcp *GcpComputeImpl) InstancesAggregatedList(apiSvc *compute.InstancesService, projectID string) *compute.InstancesAggregatedListCall {
	return apiSvc.AggregatedList(projectID)
}

// InstancesPages invokes cb for each page of results.
// Returns error on failure
func (gcp *GcpComputeImpl) InstancesPages(ctx context.Context, listCall *compute.InstancesAggregatedListCall, cb callbackInstancesPages) error {
	return listCall.Pages(ctx, cb)
}

// NewReservationsService returns *compute.ReservationsService
func (gcp *GcpComputeImpl) NewReservationsService(svc *compute.Service) *compute.ReservationsService {
	return compute.NewReservationsService(svc)
}

// ReservationsAggregatedList returns *compute.ReservationsAggregatedListCall
func (gcp *GcpComputeImpl) ReservationsAggregatedList(apiSvc *compute.ReservationsService, projectID string) *compute.ReservationsAggregatedListCall {
	return apiSvc.AggregatedList(projectID)
}

// ReservationsPages invokes cb for each page of results.
// Returns error on failure
func (gcp *GcpComputeImpl) ReservationsPages(ctx context.Context, listCall *compute.ReservationsAggregatedListCall, cb callbackReservationsPages) error {
	return listCall.Pages(ctx, cb)
}

// NewRoutersService returns *compute.RoutersService
func (gcp *GcpComputeImpl) NewRoutersService(svc *compute.Service) *compute.RoutersService {
	return compute.NewRoutersService(svc)
}

// RoutersAggregatedList returns *compute.RoutersAggregatedListCall
func (gcp *GcpComputeImpl) RoutersAggregatedList(apiSvc *compute.RoutersService, projectID string) *compute.RoutersAggregatedListCall {
	return apiSvc.AggregatedList(projectID)
}

// RoutersPages invokes cb for each page of results.
// Returns error on failure
func (gcp *GcpComputeImpl) RoutersPages(ctx context.Context, listCall *compute.RoutersAggregatedListCall, cb callbackRoutersPages) error {
	return listCall.Pages(ctx, cb)
}

// NewVpnTunnelsService returns *compute.VpnTunnelsService
func (gcp *GcpComputeImpl) NewVpnTunnelsService(svc *compute.Service) *compute.VpnTunnelsService {
	return compute.NewVpnTunnelsService(svc)
}

// VpnTunnelsAggregatedList returns *compute.VpnTunnelsAggregatedListCall
func (gcp *GcpComputeImpl) VpnTunnelsAggregatedList(apiSvc *compute.VpnTunnelsService, projectID string) *compute.VpnTunnelsAggregatedListCall {
	return apiSvc.AggregatedList(projectID)
}

// VpnTunnelsPages invokes cb for each page of results.
// Returns error on failure
func (gcp *GcpComputeImpl) VpnTunnelsPages(ctx context.Context, listCall *compute.VpnTunnelsAggregatedListCall, cb callbackVpnTunnelsPages) error {
	return listCall.Pages(ctx, cb)
}

// NewVpnGatewaysService returns *compute.VpnGatewaysService
func (gcp *GcpComputeImpl) NewVpnGatewaysService(svc *compute.Service) *compute.VpnGatewaysService {
	return compute.NewVpnGatewaysService(svc)
}

// VpnGatewaysAggregatedList returns *compute.VpnGatewaysAggregatedListCall
func (gcp *GcpComputeImpl) VpnGatewaysAggregatedList(apiSvc *compute.VpnGatewaysService, projectID string) *compute.VpnGatewaysAggregatedListCall {
	return apiSvc.AggregatedList(projectID)
}

// VpnGatewaysPages invokes cb for each page of results.
// Returns error on failure
func (gcp *GcpComputeImpl) VpnGatewaysPages(ctx context.Context, listCall *compute.VpnGatewaysAggregatedListCall, cb callbackVpnGatewaysPages) error {
	return listCall.Pages(ctx, cb)
}

// NewNetworksService returns *compute.NetworksService
func (gcp *GcpComputeImpl) NewNetworksService(svc *compute.Service) *compute.NetworksService {
	return compute.NewNetworksService(svc)
}

// NetworksList returns *compute.NetworksListCall
func (gcp *GcpComputeImpl) NetworksList(apiSvc *compute.NetworksService, projectID string) *compute.NetworksListCall {
	return apiSvc.List(projectID)
}

// NetworksPages invokes cb for each page of results.
// Returns error on failure
func (gcp *GcpComputeImpl) NetworksPages(ctx context.Context, listCall *compute.NetworksListCall, cb callbackNetworksPages) error {
	return listCall.Pages(ctx, cb)
}

// NewImagesService returns *compute.ImagesService
func (gcp *GcpComputeImpl) NewImagesService(svc *compute.Service) *compute.ImagesService {
	return compute.NewImagesService(svc)
}

// ImagesList returns *compute.ImagesListCall
func (gcp *GcpComputeImpl) ImagesList(apiSvc *compute.ImagesService, projectID string) *compute.ImagesListCall {
	return apiSvc.List(projectID)
}

// ImagesPages invokes cb for each page of results.
// Returns error on failure
func (gcp *GcpComputeImpl) ImagesPages(ctx context.Context, listCall *compute.ImagesListCall, cb callbackImagesPages) error {
	return listCall.Pages(ctx, cb)
}

// NewInterconnectsService returns *compute.InterconnectsService
func (gcp *GcpComputeImpl) NewInterconnectsService(svc *compute.Service) *compute.InterconnectsService {
	return compute.NewInterconnectsService(svc)
}

// InterconnectsList returns *compute.InterconnectsListCall
func (gcp *GcpComputeImpl) InterconnectsList(apiSvc *compute.InterconnectsService, projectID string) *compute.InterconnectsListCall {
	return apiSvc.List(projectID)
}

// InterconnectsPages invokes cb for each page of results.
// Returns error on failure
func (gcp *GcpComputeImpl) InterconnectsPages(ctx context.Context, listCall *compute.InterconnectsListCall, cb callbackInterconnectsPages) error {
	return listCall.Pages(ctx, cb)
}

// NewRoutesService returns *compute.RoutesService
func (gcp *GcpComputeImpl) NewRoutesService(svc *compute.Service) *compute.RoutesService {
	return compute.NewRoutesService(svc)
}

// RoutesList returns *compute.RoutesListCall
func (gcp *GcpComputeImpl) RoutesList(apiSvc *compute.RoutesService, projectID string) *compute.RoutesListCall {
	return apiSvc.List(projectID)
}

// RoutesPages invokes cb for each page of results.
// Returns error on failure
func (gcp *GcpComputeImpl) RoutesPages(ctx context.Context, listCall *compute.RoutesListCall, cb callbackRoutesPages) error {
	return listCall.Pages(ctx, cb)
}
