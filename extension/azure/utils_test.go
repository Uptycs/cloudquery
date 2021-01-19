package azure

import (
	"testing"

	"github.com/Uptycs/cloudquery/utilities"
	"github.com/stretchr/testify/assert"
)

var tableConfigJSON = `
{
	"test_table_1": {
    	"aws": {},
		"gcp": {},
		"azure": {
			"subscriptionIdAttribute": "subscription_id",
			"tenantIdAttribute": "abc"
		},
    	"parsedAttributes": []
	}
}`

func TestRowToMap(t *testing.T) {
	err := utilities.ReadTableConfig([]byte(tableConfigJSON))
	assert.Nil(t, err)

	subID, tenantID, rscGroup := "test-account", "us-east4", ""
	inRow := make(map[string]interface{})
	tabConfig := utilities.TableConfigurationMap["test_table_1"]
	outRow := RowToMap(inRow, subID, tenantID, rscGroup, tabConfig)

	assert.Equal(t, subID, outRow["subscription_id"])
	assert.Equal(t, tenantID, outRow["abc"])
}
