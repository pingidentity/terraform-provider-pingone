// Copyright © 2025 Ping Identity Corporation

// Package client provides client configuration and initialization functions for the PingOne Terraform provider.
// This package contains the core client configuration structures and methods for connecting to the PingOne platform API.
package client

import (
	"context"
	"fmt"

	"github.com/patrickcping/pingone-go-sdk-v2/pingone"
)

// Client represents a configured PingOne API client with additional global options.
// This structure wraps the PingOne SDK client and provides access to provider-specific configuration.
type Client struct {
	// API is the initialized PingOne SDK client used for making API calls to the PingOne platform
	API *pingone.Client
	// GlobalOptions contains provider-wide configuration options that affect client behavior
	GlobalOptions *GlobalOptions
}

// APIClient creates and configures a new PingOne API client using the provided configuration.
// It returns a configured Client instance that can be used to interact with the PingOne platform.
// The version parameter is used to construct the user agent string for API requests.
// The context parameter is required for client initialization and authentication flows.
// Returns an error if client configuration or initialization fails.
func (c *Config) APIClient(ctx context.Context, version string) (*Client, error) {

	userAgent := fmt.Sprintf("terraform-provider-pingone/%s", version)

	if v := c.UserAgentAppend; v != nil && *v != "" {
		userAgent += fmt.Sprintf(" %s", *v)
	}

	config := &pingone.Config{
		ClientID:             &c.ClientID,
		ClientSecret:         &c.ClientSecret,
		EnvironmentID:        &c.EnvironmentID,
		AccessToken:          &c.AccessToken,
		RegionCode:           c.RegionCode,
		APIHostnameOverride:  c.APIHostnameOverride,
		AuthHostnameOverride: c.AuthHostnameOverride,
		UserAgentSuffix:      &userAgent,
		ProxyURL:             c.ProxyURL,
	}

	client, err := config.APIClient(ctx)
	if err != nil {
		return nil, err
	}

	tfClient := &Client{
		API:           client,
		GlobalOptions: c.GlobalOptions,
	}

	return tfClient, nil
}
