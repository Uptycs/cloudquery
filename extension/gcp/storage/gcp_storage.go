/**
 * Copyright (c) 2020-present, The cloudquery authors
 *
 * This source code is licensed as defined by the LICENSE file found in the
 * root directory of this source tree.
 *
 * SPDX-License-Identifier: (Apache-2.0 OR GPL-2.0-only)
 */

package storage

import (
	"context"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

//type callbackDisksPages func(*storage.DiskAggregatedList) error

// GcpStorageInterface abstracts storage APIs accessed by gcp_storage_* tables
type GcpStorageInterface interface {
	NewClient(ctx context.Context, opts ...option.ClientOption) (*storage.Client, error)
	Buckets(ctx context.Context, client *storage.Client, projectID string) *storage.BucketIterator
	BucketsNewPager(itr *storage.BucketIterator, pageSize int, pageToken string) *iterator.Pager
}

// GcpStorageHandler encloses GcpStorageInterface's instance (mock or otherwise)
type GcpStorageHandler struct {
	svcInterface GcpStorageInterface
}

// NewGcpStorageHandler returns a new instance of GcpStorageHandler with provided intf
func NewGcpStorageHandler(intf GcpStorageInterface) GcpStorageHandler {
	return GcpStorageHandler{svcInterface: intf}
}
