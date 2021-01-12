package gcp

import (
	"encoding/json"
	"testing"

	"github.com/Uptycs/cloudquery/utilities"
)

var tableConfigJSON = `
{
	"test_table_1": {
    	"aws": {},
    	"gcp": {
      		"projectIdAttribute": "project_id"
    	},
    	"parsedAttributes": [
			{
				"sourceName": "Description",
				"targetName": "description",
				"targetType": "TEXT",
				"enabled": true
			},
			{
				"sourceName": "Name",
				"targetName": "name",
				"targetType": "TEXT",
				"enabled": false
			}
		]
	}
}`

type rowToMapTestInputType struct {
	Src string
	Dst string
	Val string
}

var rowToMapTestIput = []rowToMapTestInputType{{"Description", "description", "testDesc"}, {"Name", "name", "testName"}}

func getTableConfig() *utilities.TableConfig {
	var configs map[string]utilities.TableConfig
	if err := json.Unmarshal([]byte(tableConfigJSON), &configs); err != nil {
		panic(err)
	}
	cfg := configs["test_table_1"]
	cfg.InitParsedAttributeConfigMap()
	return &cfg
}

func TestRowToMap(t *testing.T) {
	inRow := make(map[string]interface{})
	for _, entry := range rowToMapTestIput {
		inRow[entry.Src] = entry.Val
	}
	tabConfig := getTableConfig()
	outRow := RowToMap(inRow, "test-project", "us-east4", tabConfig)
	for _, entry := range rowToMapTestIput {
		if outRow[entry.Dst] != entry.Val {
			t.Errorf("%+v != %+v", outRow[entry.Dst], entry.Val)
		}
	}
}
