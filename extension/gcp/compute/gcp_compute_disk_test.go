package compute

import (
	"context"
	"strconv"
	"testing"

	"github.com/kolide/osquery-go/plugin/table"
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
			Name:   "Test1",
			SizeGb: 20,
		},
		{
			Name: "Test2",
		},
	}
	mockSvc.AddDisks(diskList)

	result, err := myGcpTest.GcpComputeDisksGenerate(ctx, qCtx)
	if err != nil {
		t.Errorf("err: %s", err.Error())
		return
	}

	if len(result) != len(diskList) {
		t.Errorf("Unexpected result length. expected %d. got %d", len(diskList), len(result))
		return
	}

	if result[0]["name"] != diskList[0].Name {
		t.Errorf("Unexpected attribute value: %s != %s", diskList[0].Name, result[0]["name"])
		return
	}

	if result[0]["size_gb"] != strconv.FormatInt(diskList[0].SizeGb, 10) {
		t.Errorf("Unexpected attribute value")
		return
	}

	mockSvc.ClearInstances()
}
