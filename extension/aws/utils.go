package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/credentials"

	"github.com/Uptycs/cloudquery/utilities"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	log "github.com/sirupsen/logrus"
)

func GetAwsSession(account *utilities.ExtensionConfigurationAwsAccount, regionCode string) (*session.Session, error) {
	if account == nil {
		utilities.GetLogger().Debug("creating default session")
		return getDefaultAwsSession(regionCode)
	}

	if len(account.ProfileName) != 0 {
		utilities.GetLogger().WithFields(log.Fields{
			"account": account.ID,
			"region":  regionCode,
			"profile": account.ProfileName,
		}).Debug("creating session")
		var enable bool = true
		sess, err := session.NewSession(&aws.Config{
			EnableEndpointDiscovery: &enable,
			Region:                  aws.String(regionCode),
			Credentials:             credentials.NewSharedCredentials(account.CredentialFile, account.ProfileName),
		})
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"account":   account.ID,
				"profile":   account.ProfileName,
				"errString": err.Error(),
			}).Error("failed to create session")
			return nil, err
		}
		return sess, nil
	} else if len(account.RoleArn) != 0 {
		// TODO: Get token from STS
		utilities.GetLogger().WithFields(log.Fields{
			"account":   account.ID,
			"profile":   account.ProfileName,
			"errString": "role arn is not yet supported",
		}).Error("failed to create session")
		return nil, fmt.Errorf("role arn is not yet supported")
	}
	return nil, nil
}

func getDefaultAwsSession(regionCode string) (*session.Session, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(regionCode),
	})
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"account":   "default",
			"region":    regionCode,
			"errString": err.Error(),
		}).Error("failed to create session")
		return nil, err
	}
	return sess, nil
}

func FetchRegions(awsSession *session.Session) ([]*ec2.Region, error) {
	svc := ec2.New(awsSession)
	awsRegions, err := svc.DescribeRegions(&ec2.DescribeRegionsInput{})
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"errString": err.Error(),
		}).Error("failed to get regions")
		return nil, err
	}
	return awsRegions.Regions, nil
}

// RowToMap converts JSON row into osquery row.
// If configured it will copy some metadata vaues into appropriate columns
func RowToMap(row map[string]interface{}, accountId string, region string, tableConfig *utilities.TableConfig) map[string]string {
	result := make(map[string]string)

	if len(tableConfig.Aws.AccountIdAttribute) != 0 {
		result[tableConfig.Aws.AccountIdAttribute] = accountId
	}
	if len(tableConfig.Aws.RegionCodeAttribute) != 0 {
		result[tableConfig.Aws.RegionCodeAttribute] = region
	}
	if len(tableConfig.Aws.RegionAttribute) != 0 {
		result[tableConfig.Aws.RegionAttribute] = region // TODO: Fix it
	}

	result = utilities.RowToMap(result, row, tableConfig)
	return result
}
