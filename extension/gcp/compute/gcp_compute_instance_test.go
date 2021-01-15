package compute

import (
	"context"
	"fmt"
	"testing"

	"github.com/kolide/osquery-go/plugin/table"
)

func TestGcpComputeInstanceGenerate(t *testing.T) {

	myGcpTest := NewGcpComputeHandler(NewGcpComputeMock())
	ctx := context.Background()
	qCtx := table.QueryContext{}

	x, err := myGcpTest.GcpComputeInstancesGenerate(ctx, qCtx)
	if err != nil {
		t.Errorf("err: %s", err.Error())
	}
	fmt.Printf("%+v\n", x)
}
