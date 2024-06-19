package client

import (
	"context"
	"fmt"

	"github.com/patrickcping/pingone-go-sdk-v2/pingone"
)

type Client struct {
	API           *pingone.Client
	GlobalOptions *GlobalOptions
}

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
