package utilities

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var tableConfigJSON = `
{
	"test_table_1": {
    	"aws": {},
		"gcp": {},
		"azure": {},
    	"parsedAttributes": [
			{
				"sourceName": "Description",
				"targetName": "description",
				"targetType": "TEXT",
				"enabled": true
			},
			{
				"sourceName": "Item_Object_Name",
				"targetName": "name",
				"targetType": "TEXT",
				"enabled": true
			},
			{
				"sourceName": "ID",
				"targetName": "id",
				"targetType": "INTEGER",
				"enabled": true
			},
			{
				"sourceName": "Item_NotNeeded_OtherObject_Prop1",
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
	Val interface{}
}

var rowToMapTestIput = []rowToMapTestInputType{
	{"Description", "description", "testDesc"},
	{"Item_Object_Name", "name", "testName"},
	{"ID", "id", 1234},
}

func TestReadTableConfig(t *testing.T) {
	readErr := ReadTableConfig([]byte(tableConfigJSON))
	assert.Nil(t, readErr)

	myTable, found := TableConfigurationMap["test_table_1"]
	assert.True(t, found)

	assert.Equal(t, 4, len(myTable.ParsedAttributes))
	assert.Equal(t, 4, len(myTable.getParsedAttributeConfigMap()))
	// Col "Item_Object_Name" is deepest enabled attributes with level 2
	assert.Equal(t, 2, myTable.MaxLevel)
}

func TestRowToMap(t *testing.T) {
	readErr := ReadTableConfig([]byte(tableConfigJSON))
	assert.Nil(t, readErr)

	tabConfig, found := TableConfigurationMap["test_table_1"]
	assert.True(t, found)

	inRow := make(map[string]interface{})
	for _, entry := range rowToMapTestIput {
		inRow[entry.Src] = entry.Val
	}
	outRow := make(map[string]string)
	outRow = RowToMap(outRow, inRow, tabConfig)
	for _, entry := range rowToMapTestIput {
		var valStr string
		valStr = fmt.Sprintf("%v", entry.Val)
		assert.Equal(t, valStr, outRow[entry.Dst])
	}
}
