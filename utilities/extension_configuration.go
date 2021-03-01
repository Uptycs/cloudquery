/**
 * Copyright (c) 2020-present, The cloudquery authors
 *
 * This source code is licensed as defined by the LICENSE file found in the
 * root directory of this source tree.
 *
 * SPDX-License-Identifier: (Apache-2.0 OR GPL-2.0-only)
 */

package utilities

// ExtensionConfigurationLogging represents configuration of a logger
type ExtensionConfigurationLogging struct {
	FileName   string `json:"fileName"`
	MaxSize    int    `json:"maxSize"`
	MaxBackups int    `json:"maxBackups"`
	MaxAge     int    `json:"maxAge"`
}

type CtS3Bucket struct {
	Name   string `json:"name"`
	Region string `json:"region"`
	Prefix string `json:"prefix"`
}

// ExtensionConfigurationAwsAccount represents configuration of an AWS account
type ExtensionConfigurationAwsAccount struct {
	ID             string       `json:"id"`
	CredentialFile string       `json:"credentialFile"`
	ProfileName    string       `json:"profileName"`
	RoleArn        string       `json:"roleArn"`
	ExternalID     string       `json:"externalId"`
	CtS3Buckets    []CtS3Bucket `json:"ctS3Buckets"`
}

// ExtensionConfigurationAws holds Accounts which is a list of AWS account configurations
type ExtensionConfigurationAws struct {
	Accounts []ExtensionConfigurationAwsAccount `json:"accounts"`
}

type CloudLogStorageBucket struct {
	Name     string   `json:"name"`
	Region   string   `json:"region"`
	LogNames []string `json:"logNames"`
}

// ExtensionConfigurationGcpAccount represents configuration of a GCP account
type ExtensionConfigurationGcpAccount struct {
	KeyFile                string                  `json:"keyFile"`
	ProjectID              string                  `json:"projectId"`
	CloudLogStorageBuckets []CloudLogStorageBucket `json:"cloudLogStorageBuckets"`
}

// ExtensionConfigurationGcp holds Accounts which is a list of GCP account configurations
type ExtensionConfigurationGcp struct {
	Accounts []ExtensionConfigurationGcpAccount `json:"accounts"`
}

// ExtensionConfigurationAzureAccount represents configuration of an Azure account
type ExtensionConfigurationAzureAccount struct {
	SubscriptionID string `json:"subscriptionId"`
	TenantID       string `json:"tenantId"`
	AuthFile       string `json:"authFile"`
}

// ExtensionConfigurationAzure holds Accounts which is a list of Azure account configurations
type ExtensionConfigurationAzure struct {
	Accounts []ExtensionConfigurationAzureAccount `json:"accounts"`
}

// ExtensionConfiguration represents the configuration for cloudquery extension
type ExtensionConfiguration struct {
	ExtConfLog   ExtensionConfigurationLogging `json:"logging"`
	ExtConfAws   ExtensionConfigurationAws     `json:"aws"`
	ExtConfGcp   ExtensionConfigurationGcp     `json:"gcp"`
	ExtConfAzure ExtensionConfigurationAzure   `json:"azure"`
}
