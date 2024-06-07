package client

import "github.com/patrickcping/pingone-go-sdk-v2/management"

type Config struct {
	ClientID             string
	ClientSecret         string
	EnvironmentID        string
	AccessToken          string
	RegionCode           *management.EnumRegionCode
	APIHostnameOverride  *string
	AuthHostnameOverride *string
	ProxyURL             *string
	GlobalOptions        *GlobalOptions
}

type GlobalOptions struct {
	Environment *EnvironmentOptions
	Population  *PopulationOptions
}

type EnvironmentOptions struct{}

type PopulationOptions struct {
	ContainsUsersForceDelete bool
}
