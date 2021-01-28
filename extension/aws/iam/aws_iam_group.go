package iam

import (
	"context"
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/Uptycs/cloudquery/utilities"

	extaws "github.com/Uptycs/cloudquery/extension/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/kolide/osquery-go/plugin/table"
)

// ListGroupsColumns returns the list of columns in the table
func ListGroupsColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("account_id"),
		table.TextColumn("arn"),
		table.TextColumn("create_date"),
		//table.BigIntColumn("create_date_ext"),
		//table.TextColumn("create_date_loc"),
		//table.BigIntColumn("create_date_loc_cache_end"),
		//table.BigIntColumn("create_date_loc_cache_start"),
		//table.TextColumn("create_date_loc_cache_zone"),
		//table.TextColumn("create_date_loc_cache_zone_is_dst"),
		//table.TextColumn("create_date_loc_cache_zone_name"),
		//table.IntegerColumn("create_date_loc_cache_zone_offset"),
		//table.TextColumn("create_date_loc_name"),
		//table.TextColumn("create_date_loc_tx"),
		//table.IntegerColumn("create_date_loc_tx_index"),
		//table.TextColumn("create_date_loc_tx_isstd"),
		//table.TextColumn("create_date_loc_tx_isutc"),
		//table.BigIntColumn("create_date_loc_tx_when"),
		//table.TextColumn("create_date_loc_zone"),
		//table.TextColumn("create_date_loc_zone_is_dst"),
		//table.TextColumn("create_date_loc_zone_name"),
		//table.IntegerColumn("create_date_loc_zone_offset"),
		//table.BigIntColumn("create_date_wall"),
		table.TextColumn("group_id"),
		table.TextColumn("group_name"),
		table.TextColumn("path"),
	}
}

// ListGroupsGenerate returns the rows in the table for all configured accounts
func ListGroupsGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	if len(utilities.ExtConfiguration.ExtConfAws.Accounts) == 0 {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_iam_group",
			"account":   "default",
		}).Info("processing account")
		results, err := processAccountListGroups(nil)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, results...)
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfAws.Accounts {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_iam_group",
				"account":   account.ID,
			}).Info("processing account")
			results, err := processAccountListGroups(&account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}

	return resultMap, nil
}

func processGlobalListGroups(tableConfig *utilities.TableConfig, account *utilities.ExtensionConfigurationAwsAccount) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	sess, err := extaws.GetAwsSession(account, "aws-global")
	if err != nil {
		return resultMap, err
	}

	accountId := utilities.AwsAccountID
	if account != nil {
		accountId = account.ID
	}

	utilities.GetLogger().WithFields(log.Fields{
		"tableName": "aws_iam_group",
		"account":   accountId,
		"region":    "aws-global",
	}).Debug("processing region")

	svc := iam.New(sess)
	params := &iam.ListGroupsInput{}

	err = svc.ListGroupsPages(params,
		func(page *iam.ListGroupsOutput, lastPage bool) bool {
			byteArr, err := json.Marshal(page)
			if err != nil {
				utilities.GetLogger().WithFields(log.Fields{
					"tableName": "aws_iam_group",
					"account":   accountId,
					"region":    "aws-global",
					"errString": err.Error(),
				}).Error("failed to marshal response")
				return lastPage
			}
			table := utilities.NewTable(byteArr, tableConfig)
			for _, row := range table.Rows {
				result := extaws.RowToMap(row, accountId, "aws-global", tableConfig)
				resultMap = append(resultMap, result)
			}
			return lastPage
		})
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_iam_group",
			"account":   accountId,
			"region":    "aws-global",
			"task":      "ListGroups",
			"errString": err.Error(),
		}).Error("failed to process region")
		return resultMap, err
	}
	return resultMap, nil
}

func processAccountListGroups(account *utilities.ExtensionConfigurationAwsAccount) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	tableConfig, ok := utilities.TableConfigurationMap["aws_iam_group"]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_iam_group",
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found")
	}
	result, err := processGlobalListGroups(tableConfig, account)
	if err != nil {
		return resultMap, err
	}
	resultMap = append(resultMap, result...)
	return resultMap, nil
}
