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
	},
	"test_table_2": {
    	"aws": {
			"regionAttribute": "region"
		},
		"gcp": {
		},
		"azure": {},
		"parsedAttributes": []
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

	myTable1, found := TableConfigurationMap["test_table_1"]
	assert.True(t, found)

	assert.Equal(t, 4, len(myTable1.ParsedAttributes))
	assert.Equal(t, 4, len(myTable1.getParsedAttributeConfigMap()))
	// Col "Item_Object_Name" is deepest enabled attributes with level 2
	assert.Equal(t, 2, myTable1.MaxLevel)

	for _, v := range TableConfigurationMap {
		assert.Equal(t, len(v.parsedAttributeConfigMap), len(v.ParsedAttributes))
	}

	assert.Equal(t, 2, len(TableConfigurationMap))
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

var tableConfigJSONBadList = []string{
	`{
		"test_table_missing_source_name": {
    		"aws": {},
			"gcp": {},
			"azure": {},
    		"parsedAttributes": [
				{
					"targetName": "description",
					"targetType": "TEXT",
					"enabled": true
				}
			]
		}
	}`,
	`{
		"test_table_missing_target_name": {
    		"aws": {},
			"gcp": {},
			"azure": {},
    		"parsedAttributes": [
				{
					"sourceName": "description",
					"targetType": "TEXT",
					"enabled": true
				}
			]
		}
	}`,
	`{
		"test_table_missing_target_type": {
    		"aws": {},
			"gcp": {},
			"azure": {},
    		"parsedAttributes": [
				{
					"sourceName": "description",
					"targetName": "description",
					"enabled": true
				}
			]
		}
	}`,
	`{
		"test_table_bad_target_type_val": {
    		"aws": {},
			"gcp": {},
			"azure": {},
    		"parsedAttributes": [
				{
					"sourceName": "description",
					"targetName": "description",
					"targetType": 123,
					"enabled": true
				}
			]
		}
	}`}

func TestReadTableConfig_missingAttrProperties(t *testing.T) {
	for _, testJSON := range tableConfigJSONBadList {
		readErr := ReadTableConfig([]byte(testJSON))
		assert.NotNil(t, readErr)
	}
}
