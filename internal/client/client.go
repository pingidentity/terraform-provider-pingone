package client

import (
	"context"
	"fmt"

	"github.com/patrickcping/pingone-go-sdk-v2/pingone"
)

type Client struct {
	API         *pingone.Client
	ForceDelete bool
}

func (c *Config) APIClient(ctx context.Context, version string) (*Client, error) {

	userAgent := fmt.Sprintf("terraform-provider-pingone/%s/go", version)

	config := &pingone.Config{
		ClientID:             &c.ClientID,
		ClientSecret:         &c.ClientSecret,
		EnvironmentID:        &c.EnvironmentID,
		AccessToken:          &c.AccessToken,
		Region:               c.Region,
		APIHostnameOverride:  c.APIHostnameOverride,
		AuthHostnameOverride: c.AuthHostnameOverride,
		UserAgentOverride:    &userAgent,
		ProxyURL:             c.ProxyURL,
	}

	client, err := config.APIClient(ctx)
	if err != nil {
		return nil, err
	}

	tfClient := &Client{
		API:         client,
		ForceDelete: c.ForceDelete,
	}

	return tfClient, nil
}
