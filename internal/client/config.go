// Copyright © 2025 Ping Identity Corporation

// Package client provides client configuration and initialization functions for the PingOne Terraform provider.
package client

import "github.com/patrickcping/pingone-go-sdk-v2/management"

// Config represents the configuration parameters required to establish a connection to the PingOne platform.
// This structure contains authentication credentials, connection settings, and provider-specific options.
type Config struct {
	// ClientID is the OAuth 2.0 client identifier for authenticating with the PingOne platform
	ClientID string
	// ClientSecret is the OAuth 2.0 client secret for authenticating with the PingOne platform
	ClientSecret string
	// EnvironmentID is the PingOne environment identifier where resources will be managed
	EnvironmentID string
	// AccessToken is an optional pre-existing access token for API authentication (alternative to client credentials)
	AccessToken string
	// RegionCode specifies the PingOne region where the environment is hosted (e.g., NA, EU, CA, AP)
	RegionCode *management.EnumRegionCode
	// APIHostnameOverride allows customization of the API hostname for on-premises or custom deployments
	APIHostnameOverride *string
	// AuthHostnameOverride allows customization of the authentication hostname for on-premises or custom deployments
	AuthHostnameOverride *string
	// ProxyURL specifies an HTTP proxy server URL for API requests when required by network configuration
	ProxyURL *string
	// GlobalOptions contains provider-wide configuration options that affect resource behavior
	GlobalOptions *GlobalOptions
	// UserAgentAppend is an optional string to append to the HTTP User-Agent header for API requests
	UserAgentAppend *string
}

// GlobalOptions represents provider-wide configuration settings that affect resource behavior across all services.
// These options control default behaviors and global policies for resource management.
type GlobalOptions struct {
	// Environment contains environment-specific configuration options
	Environment *EnvironmentOptions
	// Population contains population-specific configuration options
	Population *PopulationOptions
}

// EnvironmentOptions represents configuration options specific to PingOne environment management.
// This structure is reserved for future environment-specific global configuration settings.
type EnvironmentOptions struct{}

// PopulationOptions represents configuration options specific to PingOne population management.
// These options control how population resources are handled by the provider.
type PopulationOptions struct {
	// ContainsUsersForceDelete determines whether populations containing users can be force deleted
	// When true, the provider will attempt to delete populations even if they contain user accounts
	ContainsUsersForceDelete bool
}
