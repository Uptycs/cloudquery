/**
 * Copyright (c) 2020-present, The cloudquery authors
 *
 * This source code is licensed as defined by the LICENSE file found in the
 * root directory of this source tree.
 *
 * SPDX-License-Identifier: (Apache-2.0 OR GPL-2.0-only)
 */

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
	mockSvc.addDisks(diskList)

	result, err := myGcpTest.GcpComputeDisksGenerate(ctx, qCtx)
	assert.Nil(t, err)

	assert.Equal(t, len(diskList), len(result))
	assert.Equal(t, diskList[0].Name, result[0]["name"])
	assert.Equal(t, strconv.FormatInt(diskList[0].SizeGb, 10), result[0]["size_gb"])

	mockSvc.clearDisks()
}
