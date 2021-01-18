package compute

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/Uptycs/cloudquery/extension/azure"
	extazure "github.com/Uptycs/cloudquery/extension/azure"

	"github.com/Uptycs/cloudquery/utilities"
	"github.com/kolide/osquery-go/plugin/table"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-01-01/network"
)

func InterfacesColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
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

func InterfacesGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	if len(utilities.ExtConfiguration.ExtConfAzure.Accounts) == 0 {
		fmt.Println("Processing default account")
		results, err := processAccountInterfaces(nil)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, results...)
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfAzure.Accounts {
			fmt.Println("Processing account:" + account.SubscriptionId)
			results, err := processAccountInterfaces(&account)
			if err != nil {
				// TODO: Continue to next account or return error ?
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
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	groups, err := azure.GetGroups(session)

	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	wg.Add(len(groups))

	tableConfig, ok := utilities.TableConfigurationMap["azure_compute_networkinterface"]
	if !ok {
		fmt.Println("getTableConfig: ", err)
		log.Fatal(err)
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
			log.Print("got error while traverising RG list: ", err)
		}

		resource := resourceItr.Value()
		byteArr, err := json.Marshal(resource)
		if err != nil {
			fmt.Println("Interfaces marshal: ", err)
			log.Fatal(err)
			continue
		}
		table := utilities.NewTable(byteArr, tableConfig)
		for _, row := range table.Rows {
			result := extazure.RowToMap(row, session.SubscriptionId, "", rg, tableConfig)
			*resultMap = append(*resultMap, result)
		}
	}
}
