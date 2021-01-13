package gcp

import (
	"encoding/json"
	"fmt"
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
			},
			{
				"sourceName": "ID",
				"targetName": "id",
				"targetType": "INTEGER",
				"enabled": true
			}
		]
	}
}`

type rowToMapTestInputType struct {
	Src string
	Dst string
	Val interface{}
}

var rowToMapTestIput = []rowToMapTestInputType{{"Description", "description", "testDesc"}, {"Name", "name", "testName"}, {"ID", "id", 1234}}

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
		var valStr string
		valStr = fmt.Sprintf("%v", entry.Val)
		if outRow[entry.Dst] != valStr {
			t.Errorf("%+v != %+v", outRow[entry.Dst], entry.Val)
		}
	}
}
