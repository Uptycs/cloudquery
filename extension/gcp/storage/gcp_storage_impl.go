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

// GcpStorageImpl routes gcp_storage_* table's API invocation to appropriate google cloud's API
type GcpStorageImpl struct {
}

// NewGcpStorageImpl returns new instance of GcpStorageImpl
func NewGcpStorageImpl() *GcpStorageImpl {
	return &GcpStorageImpl{}
}

// NewClient returns storage.Client or error
func (gcp *GcpStorageImpl) NewClient(ctx context.Context, opts ...option.ClientOption) (*storage.Client, error) {
	return storage.NewClient(ctx, opts...)
}

// Buckets returns BucketIterator for given projectID
func (gcp *GcpStorageImpl) Buckets(ctx context.Context, client *storage.Client, projectID string) *storage.BucketIterator {
	return client.Buckets(ctx, projectID)
}

// BucketsNewPager returns Pager with Buckets.
// Refer to iterator.NewPager() to see how pageSize and pageToken can be set for pagination
func (gcp *GcpStorageImpl) BucketsNewPager(itr *storage.BucketIterator, pageSize int, pageToken string) *iterator.Pager {
	return iterator.NewPager(itr, pageSize, pageToken)
}
