package compute

import (
	"context"
	"strconv"
	"testing"

	"github.com/kolide/osquery-go/plugin/table"
	"github.com/stretchr/testify/assert"
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
			Name:         "Test2",
			CpuPlatform:  "Intel Haswell",
			CanIpForward: false,
		},
	}
	mockSvc.AddInstances(instList)

	result, err := myGcpTest.GcpComputeInstancesGenerate(ctx, qCtx)
	assert.Nil(t, err)

	assert.Equal(t, len(instList), len(result))
	assert.Equal(t, instList[0].Name, result[0]["name"])
	assert.Equal(t, "", result[0]["cpu_platform"])
	assert.Equal(t, strconv.FormatBool(instList[0].CanIpForward), result[0]["can_ip_forward"])
	//assert.Equal(t, strconv.FormatBool(instList[1].CanIpForward), result[1]["can_ip_forward"])

	mockSvc.ClearInstances()
}
