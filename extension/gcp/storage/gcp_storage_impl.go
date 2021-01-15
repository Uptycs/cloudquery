package storage

import (
	"context"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type GcpStorageImpl struct {
}

func NewGcpStorageImpl() *GcpStorageImpl {
	return &GcpStorageImpl{}
}

func (gcp *GcpStorageImpl) NewClient(ctx context.Context, opts ...option.ClientOption) (*storage.Client, error) {
	return storage.NewClient(ctx, opts...)
}

func (gcp *GcpStorageImpl) Buckets(ctx context.Context, client *storage.Client, projectID string) *storage.BucketIterator {
	return client.Buckets(ctx, projectID)
}

func (gcp *GcpStorageImpl) BucketsNewPager(itr *storage.BucketIterator, pageSize int, pageToken string) *iterator.Pager {
	return iterator.NewPager(itr, pageSize, pageToken)
}
