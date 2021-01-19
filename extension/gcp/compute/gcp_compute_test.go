package compute

import (
	"os"
	"testing"

	"github.com/Uptycs/cloudquery/utilities"
)

var tableConfigJSON = `
{
	"gcp_compute_disk": {
	  "aws": {},
	  "gcp": {
		"projectIdAttribute": "project_id"
	  },
	  "azure": {},
	  "parsedAttributes": [
		{
		  "sourceName": "items_sizeGb",
		  "targetName": "size_gb",
		  "targetType": "INTEGER",
		  "enabled": true
		},
		{
		  "sourceName": "items_description",
		  "targetName": "description",
		  "targetType": "TEXT",
		  "enabled": true
		}
	  ]
	},
	"gcp_compute_instance": {
	  "aws": {},
	  "gcp": {
		"projectIdAttribute": "project_id"
	  },
	  "azure": {},
	  "parsedAttributes": [
		{
		  "sourceName": "items_canIpForward",
		  "targetName": "can_ip_forward",
		  "targetType": "TEXT",
		  "enabled": true
		},
		{
			"sourceName": "items_name",
			"targetName": "name",
			"targetType": "TEXT",
			"enabled": true
		}
	  ]
	},
	"gcp_compute_network": {
	  "aws": {},
	  "gcp": {
		"projectIdAttribute": "project_id"
	  },
	  "azure": {},
	  "parsedAttributes": [
		{
			"sourceName": "items_name",
			"targetName": "name",
			"targetType": "TEXT",
			"enabled": true
		},
		{
			"sourceName": "items_creationTimestamp",
			"targetName": "creation_timestamp",
			"targetType": "TEXT",
			"enabled": true
		},
		{
			"sourceName": "items_subnetworks",
			"targetName": "subnetworks",
			"targetType": "TEXT",
			"enabled": true
		}
	  ]
	}
  }
`

func TestMain(m *testing.M) {
	readErr := utilities.ReadTableConfig([]byte(tableConfigJSON))
	if readErr != nil {
		os.Exit(1)
	}
	os.Exit(m.Run())
}
