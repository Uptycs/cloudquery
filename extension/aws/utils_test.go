/**
 * Copyright (c) 2020-present, The cloudquery authors
 *
 * This source code is licensed as defined by the LICENSE file found in the
 * root directory of this source tree.
 *
 * SPDX-License-Identifier: (Apache-2.0 OR GPL-2.0-only)
 */

package aws

import (
	"os"
	"testing"

	"github.com/Uptycs/cloudquery/utilities"
	"github.com/stretchr/testify/assert"
)

var tableConfigJSON = `
{
	"test_table_1": {
    	"aws": {
			"regionCodeAttribute": "region_code",
			"accountIdAttribute": "account_id"	  
		},
		"gcp": {},
		"azure": {},
    	"parsedAttributes": []
	}
}`

func TestMain(m *testing.M) {
	utilities.CreateLogger(true, 20, 1, 30)
	os.Exit(m.Run())
}

func TestRowToMap(t *testing.T) {
	err := utilities.ReadTableConfig([]byte(tableConfigJSON))
	assert.Nil(t, err)

	acntID, region := "test-account", "us-east4"
	inRow := make(map[string]interface{})
	tabConfig := utilities.TableConfigurationMap["test_table_1"]
	outRow := RowToMap(inRow, acntID, region, tabConfig)

	assert.Equal(t, acntID, outRow["account_id"])
	assert.Equal(t, region, outRow["region_code"])
}
