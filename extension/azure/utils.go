/**
 * Copyright (c) 2020-present, The cloudquery authors
 *
 * This source code is licensed as defined by the LICENSE file found in the
 * root directory of this source tree.
 *
 * SPDX-License-Identifier: (Apache-2.0 OR GPL-2.0-only)
 */

package azure

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-02-01/resources"
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

// GetAuthSession creates an authorizer for the given account
// If account is nil, it creates an authorizer for the default account,
// by locating the auth file by reading "AZURE_AUTH_LOCATION" env variable
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

// RowToMap converts JSON row into osquery row.
// If configured it will copy some metadata vaues into appropriate columns
func RowToMap(row map[string]interface{}, subscriptionId string, tenantId string, resourceGroup string, tableConfig *utilities.TableConfig) map[string]string {
	result := make(map[string]string)
	if len(tableConfig.Azure.SubscriptionIDAttribute) != 0 {
		result[tableConfig.Azure.SubscriptionIDAttribute] = subscriptionId
	}
	if len(tableConfig.Azure.TenantIDAttribute) != 0 {
		result[tableConfig.Azure.TenantIDAttribute] = tenantId
	}
	if len(tableConfig.Azure.ResourceGroupAttribute) != 0 {
		result[tableConfig.Azure.ResourceGroupAttribute] = resourceGroup
	}

	result = utilities.RowToMap(result, row, tableConfig)
	return result
}

// GetGroups returns the list of resource groups for given azure session
func GetGroups(session *AzureSession) ([]string, error) {
	tab := make([]string, 0)
	var err error

	grClient := resources.NewGroupsClient(session.SubscriptionId)
	grClient.Authorizer = session.Authorizer

	for list, err := grClient.ListComplete(context.Background(), "", nil); list.NotDone(); err = list.Next() {
		if err != nil {
			return nil, errors.Wrap(err, "error traverising resource group list")
		}
		rgName := *list.Value().Name
		tab = append(tab, rgName)
	}
	return tab, err
}
