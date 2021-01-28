package ec2

import (
	"context"
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/Uptycs/cloudquery/utilities"

	extaws "github.com/Uptycs/cloudquery/extension/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/kolide/osquery-go/plugin/table"
)

// DescribeInstancesColumns returns the list of columns in the table
func DescribeInstancesColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("account_id"),
		table.TextColumn("region_code"),
		table.TextColumn("groups"),
		//table.TextColumn("groups_group_id"),
		//table.TextColumn("groups_group_name"),
		//table.TextColumn("instances"),
		table.TextColumn("instances_root_device_type"),
		table.TextColumn("instances_spot_instance_request_id"),
		table.TextColumn("instances_metadata_options"),
		//table.TextColumn("instances_metadata_options_state"),
		//table.TextColumn("instances_metadata_options_http_tokens"),
		//table.IntegerColumn("instances_metadata_options_http_put_response_hop_limit"),
		//table.TextColumn("instances_metadata_options_http_endpoint"),
		table.IntegerColumn("instances_ami_launch_index"),
		table.TextColumn("instances_kernel_id"),
		table.TextColumn("instances_launch_time"),
		table.TextColumn("instances_placement"),
		//table.TextColumn("instances_placement_tenancy"),
		//table.TextColumn("instances_placement_spread_domain"),
		//table.TextColumn("instances_placement_host_resource_group_arn"),
		//table.TextColumn("instances_placement_availability_zone"),
		//table.TextColumn("instances_placement_affinity"),
		//table.TextColumn("instances_placement_group_name"),
		//table.IntegerColumn("instances_placement_partition_number"),
		//table.TextColumn("instances_placement_host_id"),
		table.TextColumn("instances_ramdisk_id"),
		table.TextColumn("instances_key_name"),
		table.TextColumn("instances_private_ip_address"),
		table.TextColumn("instances_ebs_optimized"),
		table.TextColumn("instances_capacity_reservation_id"),
		table.TextColumn("instances_network_interfaces"),
		//table.TextColumn("instances_network_interfaces_source_dest_check"),
		//table.TextColumn("instances_network_interfaces_vpc_id"),
		//table.TextColumn("instances_network_interfaces_association"),
		//table.TextColumn("instances_network_interfaces_association_ip_owner_id"),
		//table.TextColumn("instances_network_interfaces_association_public_dns_name"),
		//table.TextColumn("instances_network_interfaces_association_public_ip"),
		//table.TextColumn("instances_network_interfaces_association_carrier_ip"),
		//table.TextColumn("instances_network_interfaces_attachment"),
		//table.TextColumn("instances_network_interfaces_attachment_delete_on_termination"),
		//table.IntegerColumn("instances_network_interfaces_attachment_device_index"),
		//table.TextColumn("instances_network_interfaces_attachment_status"),
		//table.IntegerColumn("instances_network_interfaces_attachment_network_card_index"),
		//table.TextColumn("instances_network_interfaces_attachment_attach_time"),
		//table.TextColumn("instances_network_interfaces_attachment_attachment_id"),
		//table.TextColumn("instances_network_interfaces_private_dns_name"),
		//table.TextColumn("instances_network_interfaces_owner_id"),
		//table.TextColumn("instances_network_interfaces_mac_address"),
		//table.TextColumn("instances_network_interfaces_network_interface_id"),
		//table.TextColumn("instances_network_interfaces_status"),
		//table.TextColumn("instances_network_interfaces_subnet_id"),
		//table.TextColumn("instances_network_interfaces_description"),
		//table.TextColumn("instances_network_interfaces_groups"),
		//table.TextColumn("instances_network_interfaces_groups_group_name"),
		//table.TextColumn("instances_network_interfaces_groups_group_id"),
		//table.TextColumn("instances_network_interfaces_ipv6_addresses"),
		//table.TextColumn("instances_network_interfaces_ipv6_addresses_ipv6_address"),
		//table.TextColumn("instances_network_interfaces_private_ip_address"),
		//table.TextColumn("instances_network_interfaces_private_ip_addresses"),
		//table.TextColumn("instances_network_interfaces_private_ip_addresses_private_ip_address"),
		//table.TextColumn("instances_network_interfaces_private_ip_addresses_association"),
		//table.TextColumn("instances_network_interfaces_private_ip_addresses_association_public_ip"),
		//table.TextColumn("instances_network_interfaces_private_ip_addresses_association_carrier_ip"),
		//table.TextColumn("instances_network_interfaces_private_ip_addresses_association_ip_owner_id"),
		//table.TextColumn("instances_network_interfaces_private_ip_addresses_association_public_dns_name"),
		//table.TextColumn("instances_network_interfaces_private_ip_addresses_primary"),
		//table.TextColumn("instances_network_interfaces_private_ip_addresses_private_dns_name"),
		//table.TextColumn("instances_network_interfaces_interface_type"),
		table.TextColumn("instances_root_device_name"),
		table.TextColumn("instances_cpu_options"),
		//table.IntegerColumn("instances_cpu_options_core_count"),
		//table.IntegerColumn("instances_cpu_options_threads_per_core"),
		table.TextColumn("instances_tags"),
		//table.TextColumn("instances_tags_key"),
		//table.TextColumn("instances_tags_value"),
		table.TextColumn("instances_hibernation_options"),
		//table.TextColumn("instances_hibernation_options_configured"),
		table.TextColumn("instances_enclave_options"),
		//table.TextColumn("instances_enclave_options_enabled"),
		table.TextColumn("instances_private_dns_name"),
		table.TextColumn("instances_product_codes"),
		//table.TextColumn("instances_product_codes_product_code_type"),
		//table.TextColumn("instances_product_codes_product_code_id"),
		table.TextColumn("instances_public_dns_name"),
		table.TextColumn("instances_state_transition_reason"),
		table.TextColumn("instances_outpost_arn"),
		table.TextColumn("instances_licenses"),
		table.TextColumn("instances_licenses_license_configuration_arn"),
		table.TextColumn("instances_state"),
		//table.IntegerColumn("instances_state_code"),
		//table.TextColumn("instances_state_name"),
		table.TextColumn("instances_elastic_inference_accelerator_associations"),
		//table.TextColumn("instances_elastic_inference_accelerator_associations_elastic_inference_accelerator_arn"),
		//table.TextColumn("instances_elastic_inference_accelerator_associations_elastic_inference_accelerator_association_id"),
		//table.TextColumn("instances_elastic_inference_accelerator_associations_elastic_inference_accelerator_association_state"),
		//table.TextColumn("instances_elastic_inference_accelerator_associations_elastic_inference_accelerator_association_time"),
		table.TextColumn("instances_security_groups"),
		//table.TextColumn("instances_security_groups_group_name"),
		//table.TextColumn("instances_security_groups_group_id"),
		table.TextColumn("instances_source_dest_check"),
		table.TextColumn("instances_virtualization_type"),
		table.TextColumn("instances_ena_support"),
		table.TextColumn("instances_instance_lifecycle"),
		table.TextColumn("instances_elastic_gpu_associations"),
		//table.TextColumn("instances_elastic_gpu_associations_elastic_gpu_association_time"),
		//table.TextColumn("instances_elastic_gpu_associations_elastic_gpu_id"),
		//table.TextColumn("instances_elastic_gpu_associations_elastic_gpu_association_id"),
		//table.TextColumn("instances_elastic_gpu_associations_elastic_gpu_association_state"),
		table.TextColumn("instances_instance_id"),
		table.TextColumn("instances_public_ip_address"),
		table.TextColumn("instances_subnet_id"),
		table.TextColumn("instances_block_device_mappings"),
		//table.TextColumn("instances_block_device_mappings_device_name"),
		//table.TextColumn("instances_block_device_mappings_ebs"),
		//table.TextColumn("instances_block_device_mappings_ebs_status"),
		//table.TextColumn("instances_block_device_mappings_ebs_volume_id"),
		//table.TextColumn("instances_block_device_mappings_ebs_attach_time"),
		//table.TextColumn("instances_block_device_mappings_ebs_delete_on_termination"),
		table.TextColumn("instances_client_token"),
		table.TextColumn("instances_state_reason"),
		//table.TextColumn("instances_state_reason_message"),
		//table.TextColumn("instances_state_reason_code"),
		table.TextColumn("instances_platform"),
		table.TextColumn("instances_vpc_id"),
		table.TextColumn("instances_iam_instance_profile"),
		//table.TextColumn("instances_iam_instance_profile_arn"),
		//table.TextColumn("instances_iam_instance_profile_id"),
		table.TextColumn("instances_sriov_net_support"),
		table.TextColumn("instances_capacity_reservation_specification"),
		//table.TextColumn("instances_capacity_reservation_specification_capacity_reservation_preference"),
		//table.TextColumn("instances_capacity_reservation_specification_capacity_reservation_target"),
		//table.TextColumn("instances_capacity_reservation_specification_capacity_reservation_target_capacity_reservation_id"),
		//table.TextColumn("instances_capacity_reservation_specification_capacity_reservation_target_capacity_reservation_resource_group_arn"),
		table.TextColumn("instances_image_id"),
		table.TextColumn("instances_instance_type"),
		table.TextColumn("instances_monitoring"),
		//table.TextColumn("instances_monitoring_state"),
		table.TextColumn("instances_architecture"),
		table.TextColumn("instances_hypervisor"),
		table.TextColumn("owner_id"),
		table.TextColumn("requester_id"),
		table.TextColumn("reservation_id"),
	}
}

// DescribeInstancesGenerate returns the rows in the table for all configured accounts
func DescribeInstancesGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	if len(utilities.ExtConfiguration.ExtConfAws.Accounts) == 0 {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_ec2_instance",
			"account":   "default",
		}).Info("processing account")
		results, err := processAccountDescribeInstances(nil)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, results...)
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfAws.Accounts {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_ec2_instance",
				"account":   account.ID,
			}).Info("processing account")
			results, err := processAccountDescribeInstances(&account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}

	return resultMap, nil
}

func processRegionDescribeInstances(tableConfig *utilities.TableConfig, account *utilities.ExtensionConfigurationAwsAccount, region *ec2.Region) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	sess, err := extaws.GetAwsSession(account, *region.RegionName)
	if err != nil {
		return resultMap, err
	}

	accountId := utilities.AwsAccountId
	if account != nil {
		accountId = account.ID
	}

	utilities.GetLogger().WithFields(log.Fields{
		"tableName": "aws_ec2_instance",
		"account":   accountId,
		"region":    *region.RegionName,
	}).Debug("processing region")

	svc := ec2.New(sess)
	params := &ec2.DescribeInstancesInput{}

	err = svc.DescribeInstancesPages(params,
		func(page *ec2.DescribeInstancesOutput, lastPage bool) bool {
			byteArr, err := json.Marshal(page)
			if err != nil {
				utilities.GetLogger().WithFields(log.Fields{
					"tableName": "aws_ec2_instance",
					"account":   accountId,
					"region":    *region.RegionName,
					"errString": err.Error(),
				}).Error("failed to marshal response")
				return lastPage
			}
			table := utilities.NewTable(byteArr, tableConfig)
			for _, row := range table.Rows {
				result := extaws.RowToMap(row, accountId, *region.RegionName, tableConfig)
				resultMap = append(resultMap, result)
			}
			return lastPage
		})
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_ec2_instance",
			"account":   accountId,
			"region":    *region.RegionName,
			"task":      "DescribeInstances",
			"errString": err.Error(),
		}).Error("failed to process region")
		return resultMap, err
	}
	return resultMap, nil
}

func processAccountDescribeInstances(account *utilities.ExtensionConfigurationAwsAccount) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	awsSession, err := extaws.GetAwsSession(account, "us-east-1")
	if err != nil {
		return resultMap, err
	}
	regions, err := extaws.FetchRegions(awsSession)
	if err != nil {
		return resultMap, err
	}
	tableConfig, ok := utilities.TableConfigurationMap["aws_ec2_instance"]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_ec2_instance",
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found")
	}
	for _, region := range regions {
		result, err := processRegionDescribeInstances(tableConfig, account, region)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, result...)
	}
	return resultMap, nil
}
