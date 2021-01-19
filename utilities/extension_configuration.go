package utilities

type ExtensionConfigurationLogging struct {
	FileName   string `json:"fileName"`
	MaxSize    int    `json:"maxSize"`
	MaxBackups int    `json:"maxBackups"`
	MaxAge     int    `json:"maxAge"`
}
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

type ExtensionConfigurationGcpAccount struct {
	KeyFile   string `json:"keyFile"`
	ProjectId string `json:"-"`
}

type ExtensionConfigurationGcp struct {
	Accounts []ExtensionConfigurationGcpAccount `json:"accounts"`
}

type ExtensionConfigurationAzureAccount struct {
	SubscriptionId string `json:"subscriptionId"`
	TenantId       string `json:"tenantId"`
	AuthFile       string `json:"authFile"`
}

type ExtensionConfigurationAzure struct {
	Accounts []ExtensionConfigurationAzureAccount `json:"accounts"`
}

type ExtensionConfiguration struct {
	ExtConfLog   ExtensionConfigurationLogging `json:"logging"`
	ExtConfAws   ExtensionConfigurationAws     `json:"aws"`
	ExtConfGcp   ExtensionConfigurationGcp     `json:"gcp"`
	ExtConfAzure ExtensionConfigurationAzure   `json:"azure"`
}
