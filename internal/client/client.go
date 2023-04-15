package client

import (
	"context"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/mjspi/pingone-neo-go-sdk/pingoneneo"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone"
)

type Client struct {
	API         *pingone.Client
	ForceDelete bool
}

func (c *Config) APIClient(ctx context.Context) (*Client, error) {

	config := &pingone.Config{
		ClientID:      c.ClientID,
		ClientSecret:  c.ClientSecret,
		EnvironmentID: c.EnvironmentID,
		AccessToken:   c.AccessToken,
		Region:        c.Region,
	}

	var client *pingone.Client

	defaultTimeout := 30

	err := resource.RetryContext(ctx, time.Duration(defaultTimeout)*time.Second, func() *resource.RetryError {
		var err error

		client, err = config.APIClient(ctx)

		if err != nil {

			if isClientRetryable(ctx, err) {
				tflog.Warn(ctx, "Client Retrying ... ")
				return resource.RetryableError(err)
			}

			return resource.NonRetryableError(err)

		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	tfClient := &Client{
		API:         client,
		ForceDelete: c.ForceDelete,
	}

	return tfClient, nil
}

// Neo SDK is currently separate from the primary PingOne SDK - defining as a separate client until the SDKs are combined.
type NeoClient struct {
	API         *pingoneneo.Client
	ForceDelete bool
}

func (c *Config) NeoAPIClient(ctx context.Context) (*NeoClient, error) {

	config := &pingoneneo.Config{
		ClientID:      c.ClientID,
		ClientSecret:  c.ClientSecret,
		EnvironmentID: c.EnvironmentID,
		AccessToken:   c.AccessToken,
		Region:        c.Region,
	}

	var client *pingoneneo.Client

	defaultTimeout := 30

	err := resource.RetryContext(ctx, time.Duration(defaultTimeout)*time.Second, func() *resource.RetryError {
		var err error

		client, err = config.APIClient(ctx)

		if err != nil {

			if isClientRetryable(ctx, err) {
				tflog.Warn(ctx, "Client Retrying ... ")
				return resource.RetryableError(err)
			}

			return resource.NonRetryableError(err)

		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	tfClient := &NeoClient{
		API:         client,
		ForceDelete: c.ForceDelete,
	}

	return tfClient, nil
}

var (
	isClientRetryable = func(ctx context.Context, err error) bool {

		// Gateway errors
		if m, mErr := regexp.MatchString("502 Bad Gateway", err.Error()); mErr == nil && m {
			tflog.Warn(ctx, "Gateway error detected on retrieving client token, available for retry")
			return true
		}

		if m, mErr := regexp.MatchString("503 Service Unavailable", err.Error()); mErr == nil && m {
			tflog.Warn(ctx, "Service error detected on retrieving client token, available for retry")
			return true
		}

		if m, mErr := regexp.MatchString("504 Gateway Timeout", err.Error()); mErr == nil && m {
			tflog.Warn(ctx, "Gateway error detected on retrieving client token, available for retry")
			return true
		}

		return false
	}
)
