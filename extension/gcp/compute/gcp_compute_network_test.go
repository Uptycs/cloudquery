package compute

import (
	"context"
	"strings"
	"testing"

	"github.com/kolide/osquery-go/plugin/table"
	"github.com/stretchr/testify/assert"
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
			Name:              "Test1",
			CreationTimestamp: "2020-11-29T22:13:42.629-08:00",
			Subnetworks: []string{
				"https://www.googleapis.com/compute/v1/projects/testProject/regions/europe-north1/subnetworks/default",
				"https://www.googleapis.com/compute/v1/projects/testProject/regions/us-east1/subnetworks/default",
				"https://www.googleapis.com/compute/v1/projects/testProject/regions/northamerica-northeast1/subnetworks/default",
			},
		},
		{
			Name: "Test2",
		},
	}
	mockSvc.addNetworks(nwkList)

	result, err := myGcpTest.GcpComputeNetworksGenerate(ctx, qCtx)
	assert.Nil(t, err)

	assert.Equal(t, len(nwkList), len(result))
	assert.Equal(t, nwkList[0].Name, result[0]["name"])
	assert.Equal(t, "", result[0]["description"])

	expectedSubNetworksVal := "[\"" + strings.Join(nwkList[0].Subnetworks, "\",\"") + "\"]"
	assert.Equal(t, expectedSubNetworksVal, result[0]["subnetworks"])

	mockSvc.clearNetworks()
}
