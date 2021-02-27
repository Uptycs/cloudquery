/**
 * Copyright (c) 2020-present, The cloudquery authors
 *
 * This source code is licensed as defined by the LICENSE file found in the
 * root directory of this source tree.
 *
 * SPDX-License-Identifier: (Apache-2.0 OR GPL-2.0-only)
 */

package extension

import (
	"encoding/json"
	"fmt"
	"github.com/Uptycs/cloudquery/utilities"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)

// InitializeLogger TODO
func InitializeLogger(verbose bool) {
	utilities.CreateLogger(verbose, utilities.ExtConfiguration.ExtConfLog.MaxSize,
		utilities.ExtConfiguration.ExtConfLog.MaxBackups, utilities.ExtConfiguration.ExtConfLog.MaxAge,
		utilities.ExtConfiguration.ExtConfLog.FileName)
}

func readProjectIDFromCredentialFile(filePath string) string {
	reader, err := ioutil.ReadFile(filePath)
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"fileName":  filePath,
			"errString": err.Error(),
		}).Info("failed to read default gcp credentials file")
		return ""
	}
	var jsonObj map[string]interface{}
	errUnmarshal := json.Unmarshal(reader, &jsonObj)
	if errUnmarshal != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"fileName":  filePath,
			"errString": errUnmarshal.Error(),
		}).Error("failed to unmarshal json")
		return ""
	}

	if idIntfc, found := jsonObj["project_id"]; found {
		return idIntfc.(string)
	}

	utilities.GetLogger().WithFields(log.Fields{
		"fileName": filePath,
	}).Error("failed to find project_id")
	return ""
}

// ReadExtensionConfigurations TODO
func ReadExtensionConfigurations(filePath string, verbose bool) error {
	utilities.AwsAccountID = os.Getenv("AWS_ACCOUNT_ID")
	reader, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("failed to read configuration file %s. err:%v\n", filePath, err)
		return err
	}
	extConfig := utilities.ExtensionConfiguration{}
	errUnmarshal := json.Unmarshal(reader, &extConfig)
	if errUnmarshal != nil {
		return errUnmarshal
	}
	utilities.ExtConfiguration = extConfig

	// Log config is read. Init the logger now.
	InitializeLogger(verbose)

	// Set projectID for GCP accounts
	for idx := range utilities.ExtConfiguration.ExtConfGcp.Accounts {
		keyFilePath := utilities.ExtConfiguration.ExtConfGcp.Accounts[idx].KeyFile
		if keyFilePath != "" {
			projectID := readProjectIDFromCredentialFile(keyFilePath)
			// Read ProjectID from keyFile
			utilities.ExtConfiguration.ExtConfGcp.Accounts[idx].ProjectID = projectID
		} else {
			// This is case where we are not using shared credentials.
			// ProjectID must be set in config.
			if utilities.ExtConfiguration.ExtConfGcp.Accounts[idx].ProjectID == "" {
				utilities.GetLogger().Error("GCP account is missing projectId setting")
			}
		}
	}

	// Read project ID from ADC
	adcFilePath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if adcFilePath != "" {
		utilities.DefaultGcpProjectID = readProjectIDFromCredentialFile(adcFilePath)
	}

	if len(utilities.ExtConfiguration.ExtConfGcp.Accounts) == 0 {
		if adcFilePath == "" {
			utilities.GetLogger().Warn("missing env GOOGLE_APPLICATION_CREDENTIALS")
		} else if utilities.DefaultGcpProjectID == "" {
			utilities.GetLogger().Warn("missing Default Project ID for GCP")
		} else {
			utilities.GetLogger().Warn("Gcp accounts not found in extension_config. Falling back to ADC\n")
		}
	}

	return nil
}
