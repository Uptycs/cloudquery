package storage

import (
	"context"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

//type callbackDisksPages func(*storage.DiskAggregatedList) error

type GcpStorageInterface interface {
	NewClient(ctx context.Context, opts ...option.ClientOption) (*storage.Client, error)
	Buckets(ctx context.Context, client *storage.Client, projectID string) *storage.BucketIterator
	BucketsNewPager(itr *storage.BucketIterator, pageSize int, pageToken string) *iterator.Pager
}

type GcpStorageHandler struct {
	svcInterface GcpStorageInterface
}

func NewGcpStorageHandler(intf GcpStorageInterface) GcpStorageHandler {
	return GcpStorageHandler{svcInterface: intf}
}
