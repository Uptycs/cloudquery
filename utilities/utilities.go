/**
 * Copyright (c) 2020-present, The cloudquery authors
 *
 * This source code is licensed as defined by the LICENSE file found in the
 * root directory of this source tree.
 *
 * SPDX-License-Identifier: (Apache-2.0 OR GPL-2.0-only)
 */

package utilities

import (
	"encoding/json"
	"unicode"

	"fmt"

	log "github.com/sirupsen/logrus"
)

var (
	// TableConfigurationMap is the map of tableName->TableConfig
	TableConfigurationMap = map[string]*TableConfig{}
	// AwsAccountID is read from env variable AWS_ACCOUNT_ID
	AwsAccountID string
	// ExtConfiguration holds extension's config data including credential files` path
	ExtConfiguration ExtensionConfiguration
	// DefaultGcpProjectID is projectID read from file set in env var GOOGLE_APPLICATION_CREDENTIALS
	DefaultGcpProjectID string
)

// ReadTableConfig parses json encoded data to read list TableConfig entries
// These are available for reading from utilities.TableConfigurationMap[]
func ReadTableConfig(jsonEncoded []byte) error {
	var configurations map[string]*TableConfig
	errUnmarshal := json.Unmarshal(jsonEncoded, &configurations)
	if errUnmarshal != nil {
		return errUnmarshal
	}
	for tableName, config := range configurations {
		GetLogger().WithFields(log.Fields{
			"tableName": tableName,
		}).Debug("found table configuration")
		for _, attr := range config.ParsedAttributes {
			if attr.SourceName == "" || attr.TargetName == "" || attr.TargetType == "" {
				return fmt.Errorf("invalid parsedAttribute entry: %+v", attr)
			}
		}
		config.initParsedAttributeConfigMap()
		TableConfigurationMap[tableName] = config
	}
	return nil
}

// RowToMap converts JSON row into osquery row
func RowToMap(inMap map[string]string, row map[string]interface{}, tableConfig *TableConfig) map[string]string {
	for key, value := range tableConfig.getParsedAttributeConfigMap() {
		if row[key] != nil {
			inMap[value.TargetName] = GetStringValue(row[key])
		}
	}
	return inMap
}

// GetSnakeCase converts a string into snakecase string
func GetSnakeCase(source string) string {
	runes := []rune(source)
	length := len(runes)

	var out []rune
	for i := 0; i < length; i++ {
		if i > 0 && unicode.IsUpper(runes[i]) &&
			((i+1 < length && unicode.IsLower(runes[i+1])) || unicode.IsLower(runes[i-1])) && runes[i-1] != '_' {
			out = append(out, '_')
		}
		out = append(out, unicode.ToLower(runes[i]))
	}
	return string(out)
}
