package client

type Config struct {
	ClientID             string
	ClientSecret         string
	EnvironmentID        string
	AccessToken          string
	Region               string
	APIHostnameOverride  *string
	AuthHostnameOverride *string
	ProxyURL             *string
	GlobalOptions        *GlobalOptions
}

type GlobalOptions struct {
	Environment *EnvironmentOptions
	Population  *PopulationOptions
}

type EnvironmentOptions struct {
	ProductionTypeForceDelete bool
}

type PopulationOptions struct {
	ContainsUsersForceDelete bool
}
