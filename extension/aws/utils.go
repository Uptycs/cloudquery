package aws

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/credentials"

	"github.com/Uptycs/cloudquery/utilities"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func GetAwsSession(account *utilities.ExtensionConfigurationAwsAccount, regionCode string) (*session.Session, error) {
	if account == nil {
		fmt.Println("Fetching default aws session")
		return getDefaultAwsSession(regionCode)
	}

	if len(account.ProfileName) != 0 {
		var enable bool = true
		sess, err := session.NewSession(&aws.Config{
			EnableEndpointDiscovery: &enable,
			Region:      aws.String(regionCode),
			Credentials: credentials.NewSharedCredentials(account.CredentialFile, account.ProfileName),
		})
		if err != nil {
			fmt.Printf("Failed to create AWS Session. Error:%v\n", err)
			return nil, err
		}
		return sess, nil
	} else if len(account.RoleArn) != 0 {
		// TODO: Get token from STS
		return nil, fmt.Errorf("role arn is not yet supported")
	}
	return nil, nil
}

func getDefaultAwsSession(regionCode string) (*session.Session, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(regionCode),
	})
	if err != nil {
		fmt.Printf("Failed to create AWS Session. Error:%v\n", err)
		return nil, err
	}
	return sess, nil
}

func FetchRegions(awsSession *session.Session) ([]*ec2.Region, error) {

	// awsSession := session.Must(session.NewSession(&aws.Config{Region: aws.String("us-east-1")})) //Credentials: credentials.NewSharedCredentials("/home/apatil/.aws/credentials", "uptycs-dev")

	// awsSession := session.Must(session.NewSession(&aws.Config{}))
	svc := ec2.New(awsSession)
	awsRegions, err := svc.DescribeRegions(&ec2.DescribeRegionsInput{})
	if err != nil {
		return nil, err
	}
	return awsRegions.Regions, nil
}

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
	for key, value := range tableConfig.GetParsedAttributeConfigMap() {
		if row[key] != nil {
			result[value.TargetName] = utilities.GetStringValue(row[key])
		}
	}
	return result
}
