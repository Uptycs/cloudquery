package utilities

type ExtensionConfigurationAwsAccount struct {
	ID             string `json:"id"`
	CredentialFile string `json:"credentialFile"`
	ProfileName    string `json:"profileName"`
	RoleArn        string `json:"roleArn"`
	ExternalId     string `json:"externalId"`
}

type ExtensionConfigurationAws struct {
	Accounts []ExtensionConfigurationAwsAccount `json:"accounts"`
}
type ExtensionConfiguration struct {
	ExtConfAws ExtensionConfigurationAws `json:"aws"`
}
