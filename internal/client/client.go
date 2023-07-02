package client

import (
	"context"

	"github.com/patrickcping/pingone-go-sdk-v2/pingone"
)

type Client struct {
	API         *pingone.Client
	ForceDelete bool
}

func (c *Config) APIClient(ctx context.Context) (*Client, error) {

	config := &pingone.Config{
		ClientID:                      c.ClientID,
		ClientSecret:                  c.ClientSecret,
		EnvironmentID:                 c.EnvironmentID,
		AccessToken:                   c.AccessToken,
		Region:                        c.Region,
		APIHostnameOverride:           c.APIHostnameOverride,
		AgreementMgmtHostnameOverride: c.AgreementMgmtHostnameOverride,
		AuthHostnameOverride:          c.AuthHostnameOverride,
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
