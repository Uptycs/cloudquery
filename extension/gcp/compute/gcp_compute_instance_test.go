package compute

import (
	"context"
	"testing"

	"github.com/kolide/osquery-go/plugin/table"
	"google.golang.org/api/compute/v1"
)

func TestGcpComputeInstanceGenerate(t *testing.T) {

	mockSvc := NewGcpComputeMock()
	myGcpTest := NewGcpComputeHandler(mockSvc)
	ctx := context.Background()
	qCtx := table.QueryContext{}

	// TODO: Test more attributes
	instList := []*compute.Instance{
		{
			Name:         "Test1",
			CpuPlatform:  "Intel Haswell",
			CanIpForward: true,
		},
		{
			Name:        "Test2",
			CpuPlatform: "Intel Haswell",
		},
	}
	mockSvc.AddInstances(instList)

	result, err := myGcpTest.GcpComputeInstancesGenerate(ctx, qCtx)
	if err != nil {
		t.Errorf("err: %s", err.Error())
		return
	}

	if len(result) != len(instList) {
		t.Errorf("Unexpected result length. expected %d. got %d", len(instList), len(result))
		return
	}

	if result[0]["can_ip_forward"] != "true" {
		t.Errorf("Unexpected attribute value")
		return
	}

	mockSvc.ClearInstances()
}
