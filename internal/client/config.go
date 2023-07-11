package client

type Config struct {
	ClientID             string
	ClientSecret         string
	EnvironmentID        string
	AccessToken          string
	Region               string
	APIHostnameOverride  *string
	AuthHostnameOverride *string
	ForceDelete          bool
}
