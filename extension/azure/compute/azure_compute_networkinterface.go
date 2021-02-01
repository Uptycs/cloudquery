/**
 * Copyright (c) 2020-present, The cloudquery authors
 *
 * This source code is licensed as defined by the LICENSE file found in the
 * root directory of this source tree.
 *
 * SPDX-License-Identifier: (Apache-2.0 OR GPL-2.0-only)
 */

package compute

import (
	"context"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"sync"

	"github.com/Uptycs/cloudquery/extension/azure"
	extazure "github.com/Uptycs/cloudquery/extension/azure"

	"github.com/Uptycs/cloudquery/utilities"
	"github.com/kolide/osquery-go/plugin/table"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-01-01/network"
)

// InterfacesColumns returns the list of columns in the table
func InterfacesColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("subscription_id"),
		table.TextColumn("etag"),
		table.TextColumn("id"),
		table.TextColumn("location"),
		table.TextColumn("name"),
		//table.TextColumn("properties"),
		table.TextColumn("dns_settings"),
		//table.TextColumn("dns_settings_applied_dns_servers"),
		//table.TextColumn("dns_settings_dns_servers"),
		//table.TextColumn("dns_settings_internal_dns_name_label"),
		//table.TextColumn("dns_settings_internal_domain_name_suffix"),
		//table.TextColumn("dns_settings_internal_fqdn"),
		table.TextColumn("enable_accelerated_networking"),
		table.TextColumn("enable_ip_forwarding"),
		table.TextColumn("ip_configurations"),
		//table.TextColumn("ip_configurations_etag"),
		//table.TextColumn("ip_configurations_id"),
		//table.TextColumn("ip_configurations_name"),
		table.TextColumn("mac_address"),
		table.TextColumn("network_security_group"),
		//table.TextColumn("network_security_group_etag"),
		//table.TextColumn("network_security_group_id"),
		//table.TextColumn("network_security_group_location"),
		//table.TextColumn("network_security_group_name"),
		//table.TextColumn("network_security_group_tags"),
		//table.TextColumn("network_security_group_type"),
		table.TextColumn("primary"),
		table.TextColumn("provisioning_state"),
		table.TextColumn("resource_guid"),
		table.TextColumn("virtual_machine"),
		//table.TextColumn("virtual_machine_id"),
		table.TextColumn("tags"),
		table.TextColumn("type"),
	}
}

// InterfacesGenerate returns the rows in the table for all configured accounts
func InterfacesGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	if len(utilities.ExtConfiguration.ExtConfAzure.Accounts) == 0 {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "azure_compute_networkinterface",
			"account":   "default",
		}).Info("processing account")
		results, err := processAccountInterfaces(nil)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, results...)
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfAzure.Accounts {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "azure_compute_networkinterface",
				"account":   account.SubscriptionID,
			}).Info("processing account")
			results, err := processAccountInterfaces(&account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}

	return resultMap, nil
}

func processAccountInterfaces(account *utilities.ExtensionConfigurationAzureAccount) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	var wg sync.WaitGroup
	session, err := azure.GetAuthSession(account)
	if err != nil {
		return resultMap, err
	}
	groups, err := azure.GetGroups(session)

	if err != nil {
		return resultMap, err
	}

	wg.Add(len(groups))

	tableConfig, ok := utilities.TableConfigurationMap["azure_compute_networkinterface"]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "azure_compute_networkinterface",
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found")
	}

	for _, group := range groups {
		go getInterfaces(session, group, &wg, &resultMap, tableConfig)
	}
	wg.Wait()
	return resultMap, nil
}

func getInterfaces(session *azure.AzureSession, rg string, wg *sync.WaitGroup, resultMap *[]map[string]string, tableConfig *utilities.TableConfig) {
	defer wg.Done()

	svcClient := network.NewInterfacesClient(session.SubscriptionId)
	svcClient.Authorizer = session.Authorizer

	for resourceItr, err := svcClient.ListComplete(context.Background(), rg); resourceItr.NotDone(); err = resourceItr.Next() {
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName":     "azure_compute_networkinterface",
				"resourceGroup": rg,
				"errString":     err.Error(),
			}).Error("failed to get resource list")
			continue
		}

		resource := resourceItr.Value()
		byteArr, err := json.Marshal(resource)
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName":     "azure_compute_networkinterface",
				"resourceGroup": rg,
				"errString":     err.Error(),
			}).Error("failed to marshal response")
			continue
		}
		table := utilities.NewTable(byteArr, tableConfig)
		for _, row := range table.Rows {
			result := extazure.RowToMap(row, session.SubscriptionId, "", rg, tableConfig)
			*resultMap = append(*resultMap, result)
		}
	}
}
