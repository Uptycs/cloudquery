package compute

import (
	"context"
	"strconv"
	"testing"

	"github.com/kolide/osquery-go/plugin/table"
	"github.com/stretchr/testify/assert"
	"google.golang.org/api/compute/v1"
)

func TestGcpComputeDiskGenerate(t *testing.T) {

	mockSvc := NewGcpComputeMock()
	myGcpTest := NewGcpComputeHandler(mockSvc)
	ctx := context.Background()
	qCtx := table.QueryContext{}

	// TODO: Test more attributes
	diskList := []*compute.Disk{
		{
			Description: "Test1",
			SizeGb:      20,
		},
		{
			Description: "Test2",
		},
	}
	mockSvc.AddDisks(diskList)

	result, err := myGcpTest.GcpComputeDisksGenerate(ctx, qCtx)
	assert.Nil(t, err)

	assert.Equal(t, len(diskList), len(result))
	assert.Equal(t, diskList[0].Name, result[0]["name"])
	assert.Equal(t, strconv.FormatInt(diskList[0].SizeGb, 10), result[0]["size_gb"])

	mockSvc.ClearDisks()
}
