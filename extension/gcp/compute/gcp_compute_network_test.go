package compute

import (
	"context"
	"testing"

	"github.com/kolide/osquery-go/plugin/table"
	"google.golang.org/api/compute/v1"
)

func TestGcpComputeNetworkGenerate(t *testing.T) {

	mockSvc := NewGcpComputeMock()
	myGcpTest := NewGcpComputeHandler(mockSvc)
	ctx := context.Background()
	qCtx := table.QueryContext{}

	// TODO: Test more attributes
	nwkList := []*compute.Network{
		{
			Name: "Test1",
		},
		{
			Name: "Test2",
		},
	}
	mockSvc.AddNetworks(nwkList)

	result, err := myGcpTest.GcpComputeNetworksGenerate(ctx, qCtx)
	if err != nil {
		t.Errorf("err: %s", err.Error())
		return
	}

	if len(result) != len(nwkList) {
		t.Errorf("Unexpected result length. expected %d. got %d", len(nwkList), len(result))
		return
	}

	if result[0]["name"] != nwkList[0].Name {
		t.Errorf("Unexpected attribute value: %s != %s", nwkList[0].Name, result[0]["name"])
		return
	}

	mockSvc.ClearInstances()
}
