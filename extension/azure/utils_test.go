package azure

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
		"gcp": {},
		"azure": {
			"subscriptionIdAttribute": "subscription_id",
			"tenantIdAttribute": "abc"
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

	subID, tenantID, rscGroup := "test-account", "us-east4", ""
	inRow := make(map[string]interface{})
	tabConfig := utilities.TableConfigurationMap["test_table_1"]
	outRow := RowToMap(inRow, subID, tenantID, rscGroup, tabConfig)

	assert.Equal(t, subID, outRow["subscription_id"])
	assert.Equal(t, tenantID, outRow["abc"])
}
