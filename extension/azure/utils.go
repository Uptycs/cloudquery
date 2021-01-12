package azure

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/Uptycs/cloudquery/utilities"
	"github.com/pkg/errors"
)

// AzureSession is an object representing session for subscription
type AzureSession struct {
	SubscriptionId string
	Authorizer     autorest.Authorizer
}

var (
	authGeneratorMutex sync.Mutex
)

func readJSON(path string) (*map[string]interface{}, error) {
	data, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, errors.Wrap(err, "Can't open the file")
	}

	contents := make(map[string]interface{})
	err = json.Unmarshal(data, &contents)

	if err != nil {
		err = errors.Wrap(err, "Can't unmarshal file")
	}

	return &contents, err
}

func GetAuthSession(account *utilities.ExtensionConfigurationAzureAccount) (*AzureSession, error) {
	authGeneratorMutex.Lock()
	defer authGeneratorMutex.Unlock()

	if account != nil {
		os.Setenv("AZURE_AUTH_LOCATION", account.AuthFile)
	}
	authorizer, err := auth.NewAuthorizerFromFile(azure.PublicCloud.ResourceManagerEndpoint)
	if err != nil {
		return nil, errors.Wrap(err, "Can't initialize authorizer")
	}
	authInfo, err := readJSON(os.Getenv("AZURE_AUTH_LOCATION"))
	if err != nil {
		return nil, errors.Wrap(err, "Can't get authinfo")
	}
	session := AzureSession{
		SubscriptionId: (*authInfo)["subscriptionId"].(string),
		Authorizer:     authorizer,
	}

	return &session, nil
}

func RowToMap(row map[string]interface{}, subscriptionId string, tenantId string, resourceGroup string, tableConfig *utilities.TableConfig) map[string]string {
	result := make(map[string]string)
	if len(tableConfig.Azure.SubscriptionIdAttribute) != 0 {
		result[tableConfig.Azure.SubscriptionIdAttribute] = subscriptionId
	}
	if len(tableConfig.Azure.TenantIdAttribute) != 0 {
		result[tableConfig.Azure.TenantIdAttribute] = tenantId
	}
	if len(tableConfig.Azure.ResourceGroupAttribute) != 0 {
		result[tableConfig.Azure.ResourceGroupAttribute] = resourceGroup
	}
	for key, value := range tableConfig.GetParsedAttributeConfigMap() {
		if row[key] != nil {
			result[value.TargetName] = utilities.GetStringValue(row[key])
		}
	}
	return result
}
