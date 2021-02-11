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
	"io/ioutil"
	"os"

	"github.com/Uptycs/cloudquery/extension/aws/acm"

	"github.com/Uptycs/cloudquery/extension/aws/s3"

	"github.com/Uptycs/cloudquery/utilities"

	"github.com/Uptycs/cloudquery/extension/aws/ec2"
	"github.com/Uptycs/cloudquery/extension/aws/iam"
	azurecompute "github.com/Uptycs/cloudquery/extension/azure/compute"
	"github.com/Uptycs/cloudquery/extension/gcp/compute"
	gcpcontainer "github.com/Uptycs/cloudquery/extension/gcp/container"
	gcpdns "github.com/Uptycs/cloudquery/extension/gcp/dns"
	gcpfile "github.com/Uptycs/cloudquery/extension/gcp/file"
	gcpfunction "github.com/Uptycs/cloudquery/extension/gcp/function"
	gcpiam "github.com/Uptycs/cloudquery/extension/gcp/iam"
	gcprun "github.com/Uptycs/cloudquery/extension/gcp/run"
	gcpsql "github.com/Uptycs/cloudquery/extension/gcp/sql"
	"github.com/Uptycs/cloudquery/extension/gcp/storage"

	"github.com/kolide/osquery-go"
	"github.com/kolide/osquery-go/plugin/table"
	log "github.com/sirupsen/logrus"
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
		projectID := readProjectIDFromCredentialFile(keyFilePath)
		utilities.ExtConfiguration.ExtConfGcp.Accounts[idx].ProjectID = projectID
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

// ReadTableConfigurations TODO
func ReadTableConfigurations(homeDir string) {
	var awsConfigFileList = []string{"aws/acm/table_config.json", "aws/ec2/table_config.json", "aws/s3/table_config.json", "aws/iam/table_config.json"}
	var gcpConfigFileList = []string{
		"gcp/compute/table_config.json",
		"gcp/storage/table_config.json",
		"gcp/iam/table_config.json",
		"gcp/sql/table_config.json",
		"gcp/dns/table_config.json",
		"gcp/file/table_config.json",
		"gcp/container/table_config.json",
		"gcp/function/table_config.json",
		"gcp/run/table_config.json",
	}
	var azureConfigFileList = []string{"azure/compute/table_config.json"}
	var configFileList = append(awsConfigFileList, gcpConfigFileList...)
	configFileList = append(configFileList, azureConfigFileList...)

	for _, fileName := range configFileList {
		utilities.GetLogger().WithFields(log.Fields{
			"fileName": homeDir + string(os.PathSeparator) + fileName,
		}).Info("reading config file")
		filePath := homeDir + string(os.PathSeparator) + fileName
		jsonEncoded, err := ioutil.ReadFile(filePath)
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"fileName":  homeDir + string(os.PathSeparator) + fileName,
				"errString": err.Error(),
			}).Error("failed to read config file")
			continue
		}
		readErr := utilities.ReadTableConfig(jsonEncoded)
		if readErr != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"fileName":  homeDir + string(os.PathSeparator) + fileName,
				"errString": readErr.Error(),
			}).Error("failed to parse config file")
			continue
		}
	}
	utilities.GetLogger().WithFields(log.Fields{
		"totalTables": len(utilities.TableConfigurationMap),
	}).Info("read all config files")
}

var gcpComputeHandler = compute.NewGcpComputeHandler(compute.NewGcpComputeImpl())
var gcpStorageHandler = storage.NewGcpStorageHandler(storage.NewGcpStorageImpl())

// RegisterPlugins TODO
func RegisterPlugins(server *osquery.ExtensionManagerServer) {
	// AWS ACM
	server.RegisterPlugin(table.NewPlugin("aws_acm_certificate", acm.ListCertificatesColumns(), acm.ListCertificatesGenerate))
	// AWS EC2
	server.RegisterPlugin(table.NewPlugin("aws_ec2_instance", ec2.DescribeInstancesColumns(), ec2.DescribeInstancesGenerate))
	server.RegisterPlugin(table.NewPlugin("aws_ec2_vpc", ec2.DescribeVpcsColumns(), ec2.DescribeVpcsGenerate))
	server.RegisterPlugin(table.NewPlugin("aws_ec2_subnet", ec2.DescribeSubnetsColumns(), ec2.DescribeSubnetsGenerate))
	server.RegisterPlugin(table.NewPlugin("aws_ec2_image", ec2.DescribeImagesColumns(), ec2.DescribeImagesGenerate))
	server.RegisterPlugin(table.NewPlugin("aws_ec2_egress_only_internet_gateway", ec2.DescribeEgressOnlyInternetGatewaysColumns(), ec2.DescribeEgressOnlyInternetGatewaysGenerate))
	server.RegisterPlugin(table.NewPlugin("aws_ec2_internet_gateway", ec2.DescribeInternetGatewaysColumns(), ec2.DescribeInternetGatewaysGenerate))
	server.RegisterPlugin(table.NewPlugin("aws_ec2_nat_gateway", ec2.DescribeNatGatewaysColumns(), ec2.DescribeNatGatewaysGenerate))
	server.RegisterPlugin(table.NewPlugin("aws_ec2_network_acl", ec2.DescribeNetworkAclsColumns(), ec2.DescribeNetworkAclsGenerate))
	server.RegisterPlugin(table.NewPlugin("aws_ec2_route_table", ec2.DescribeRouteTablesColumns(), ec2.DescribeRouteTablesGenerate))
	server.RegisterPlugin(table.NewPlugin("aws_ec2_security_group", ec2.DescribeSecurityGroupsColumns(), ec2.DescribeSecurityGroupsGenerate))
	server.RegisterPlugin(table.NewPlugin("aws_ec2_tag", ec2.DescribeTagsColumns(), ec2.DescribeTagsGenerate))
	server.RegisterPlugin(table.NewPlugin("aws_ec2_address", ec2.DescribeAddressesColumns(), ec2.DescribeAddressesGenerate))
	server.RegisterPlugin(table.NewPlugin("aws_ec2_flowlog", ec2.DescribeFlowLogsColumns(), ec2.DescribeFlowLogsGenerate))
	server.RegisterPlugin(table.NewPlugin("aws_ec2_keypair", ec2.DescribeKeyPairsColumns(), ec2.DescribeKeyPairsGenerate))
	server.RegisterPlugin(table.NewPlugin("aws_ec2_snapshot", ec2.DescribeSnapshotsColumns(), ec2.DescribeSnapshotsGenerate))
	server.RegisterPlugin(table.NewPlugin("aws_ec2_volume", ec2.DescribeVolumesColumns(), ec2.DescribeVolumesGenerate))
	// AWS S3
	server.RegisterPlugin(table.NewPlugin("aws_s3_bucket", s3.ListBucketsColumns(), s3.ListBucketsGenerate))
	// AWS IAM
	server.RegisterPlugin(table.NewPlugin("aws_iam_user", iam.ListUsersColumns(), iam.ListUsersGenerate))
	server.RegisterPlugin(table.NewPlugin("aws_iam_role", iam.ListRolesColumns(), iam.ListRolesGenerate))
	server.RegisterPlugin(table.NewPlugin("aws_iam_group", iam.ListGroupsColumns(), iam.ListGroupsGenerate))
	server.RegisterPlugin(table.NewPlugin("aws_iam_policy", iam.ListPoliciesColumns(), iam.ListPoliciesGenerate))
	server.RegisterPlugin(table.NewPlugin("aws_iam_account_password_policy", iam.GetAccountPasswordPolicyColumns(), iam.GetAccountPasswordPolicyGenerate))
	// GCP Compute
	server.RegisterPlugin(table.NewPlugin("gcp_compute_instance", gcpComputeHandler.GcpComputeInstancesColumns(), gcpComputeHandler.GcpComputeInstancesGenerate))
	server.RegisterPlugin(table.NewPlugin("gcp_compute_network", gcpComputeHandler.GcpComputeNetworksColumns(), gcpComputeHandler.GcpComputeNetworksGenerate))
	server.RegisterPlugin(table.NewPlugin("gcp_compute_disk", gcpComputeHandler.GcpComputeDisksColumns(), gcpComputeHandler.GcpComputeDisksGenerate))
	server.RegisterPlugin(table.NewPlugin("gcp_compute_image", gcpComputeHandler.GcpComputeImagesColumns(), gcpComputeHandler.GcpComputeImagesGenerate))
	server.RegisterPlugin(table.NewPlugin("gcp_compute_interconnect", gcpComputeHandler.GcpComputeInterconnectsColumns(), gcpComputeHandler.GcpComputeInterconnectsGenerate))
	server.RegisterPlugin(table.NewPlugin("gcp_compute_route", gcpComputeHandler.GcpComputeRoutesColumns(), gcpComputeHandler.GcpComputeRoutesGenerate))
	server.RegisterPlugin(table.NewPlugin("gcp_compute_reservation", gcpComputeHandler.GcpComputeReservationsColumns(), gcpComputeHandler.GcpComputeReservationsGenerate))
	server.RegisterPlugin(table.NewPlugin("gcp_compute_router", gcpComputeHandler.GcpComputeRoutersColumns(), gcpComputeHandler.GcpComputeRoutersGenerate))
	server.RegisterPlugin(table.NewPlugin("gcp_compute_vpn_tunnel", gcpComputeHandler.GcpComputeVpnTunnelsColumns(), gcpComputeHandler.GcpComputeVpnTunnelsGenerate))
	server.RegisterPlugin(table.NewPlugin("gcp_compute_vpn_gateway", gcpComputeHandler.GcpComputeVpnGatewaysColumns(), gcpComputeHandler.GcpComputeVpnGatewaysGenerate))
	// GCP Storage
	server.RegisterPlugin(table.NewPlugin("gcp_storage_bucket", gcpStorageHandler.GcpStorageBucketColumns(), gcpStorageHandler.GcpStorageBucketGenerate))
	// GCP IAM
	server.RegisterPlugin(table.NewPlugin("gcp_iam_role", gcpiam.GcpIamRolesColumns(), gcpiam.GcpIamRolesGenerate))
	server.RegisterPlugin(table.NewPlugin("gcp_iam_service_account", gcpiam.GcpIamServiceAccountsColumns(), gcpiam.GcpIamServiceAccountsGenerate))
	// GCP SQL
	server.RegisterPlugin(table.NewPlugin("gcp_sql_instance", gcpsql.GcpSQLInstancesColumns(), gcpsql.GcpSQLInstancesGenerate))
	server.RegisterPlugin(table.NewPlugin("gcp_sql_database", gcpsql.GcpSQLDatabasesColumns(), gcpsql.GcpSQLDatabasesGenerate))
	// GCP DNS
	server.RegisterPlugin(table.NewPlugin("gcp_dns_managed_zone", gcpdns.GcpDNSManagedZonesColumns(), gcpdns.GcpDNSManagedZonesGenerate))
	server.RegisterPlugin(table.NewPlugin("gcp_dns_policy", gcpdns.GcpDNSPoliciesColumns(), gcpdns.GcpDNSPoliciesGenerate))
	// GCP File
	server.RegisterPlugin(table.NewPlugin("gcp_file_instance", gcpfile.GcpFileInstancesColumns(), gcpfile.GcpFileInstancesGenerate))
	server.RegisterPlugin(table.NewPlugin("gcp_file_backup", gcpfile.GcpFileBackupsColumns(), gcpfile.GcpFileBackupsGenerate))
	// GCP Container
	server.RegisterPlugin(table.NewPlugin("gcp_container_cluster", gcpcontainer.GcpContainerClustersColumns(), gcpcontainer.GcpContainerClustersGenerate))
	// GCP Cloud Function
	server.RegisterPlugin(table.NewPlugin("gcp_cloud_function", gcpfunction.GcpCloudFunctionsColumns(), gcpfunction.GcpCloudFunctionsGenerate))
	// GCP Cloud Run
	server.RegisterPlugin(table.NewPlugin("gcp_cloud_run_service", gcprun.GcpCloudRunServicesColumns(), gcprun.GcpCloudRunServicesGenerate))
	server.RegisterPlugin(table.NewPlugin("gcp_cloud_run_revision", gcprun.GcpCloudRunRevisionsColumns(), gcprun.GcpCloudRunRevisionsGenerate))
	// Azure Compute
	server.RegisterPlugin(table.NewPlugin("azure_compute_vm", azurecompute.VirtualMachinesColumns(), azurecompute.VirtualMachinesGenerate))
	server.RegisterPlugin(table.NewPlugin("azure_compute_networkinterface", azurecompute.InterfacesColumns(), azurecompute.InterfacesGenerate))
}
