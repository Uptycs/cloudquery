package utilities

// ExtensionConfigurationLogging represents configuration of a logger
type ExtensionConfigurationLogging struct {
	FileName   string `json:"fileName"`
	MaxSize    int    `json:"maxSize"`
	MaxBackups int    `json:"maxBackups"`
	MaxAge     int    `json:"maxAge"`
}

// ExtensionConfigurationAwsAccount represents configuration of an AWS account
type ExtensionConfigurationAwsAccount struct {
	ID             string `json:"id"`
	CredentialFile string `json:"credentialFile"`
	ProfileName    string `json:"profileName"`
	RoleArn        string `json:"roleArn"`
	ExternalId     string `json:"externalId"`
}

// Accounts is the list of AWS account configuration
type ExtensionConfigurationAws struct {
	Accounts []ExtensionConfigurationAwsAccount `json:"accounts"`
}

// ExtensionConfigurationGcpAccount represents configuration of a GCP account
type ExtensionConfigurationGcpAccount struct {
	KeyFile   string `json:"keyFile"`
	ProjectId string `json:"-"`
}

// Accounts is the list of GCP account configuration
type ExtensionConfigurationGcp struct {
	Accounts []ExtensionConfigurationGcpAccount `json:"accounts"`
}

// ExtensionConfigurationAzureAccount represents configuration of an Azure account
type ExtensionConfigurationAzureAccount struct {
	SubscriptionId string `json:"subscriptionId"`
	TenantId       string `json:"tenantId"`
	AuthFile       string `json:"authFile"`
}

// Accounts is the list of Azure account configuration
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
