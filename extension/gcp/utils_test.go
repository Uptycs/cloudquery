package gcp

import (
	"os"
	"testing"

	"github.com/Uptycs/cloudquery/utilities"
	"github.com/stretchr/testify/assert"
)

var tableConfigJSON = `
{
	"test_table_1": {
    	"aws": {},
    	"gcp": {
      		"projectIdAttribute": "project_id"
    	},
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

	projName, zone := "test-project", "us-east4-zone1"
	inRow := make(map[string]interface{})
	tabConfig := utilities.TableConfigurationMap["test_table_1"]
	outRow := RowToMap(inRow, projName, zone, tabConfig)

	assert.Equal(t, projName, outRow["project_id"])
	assert.Equal(t, "", outRow["zone"])
}
