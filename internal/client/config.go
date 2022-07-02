package config

import (
	"context"
	"fmt"
	"log"

	pingone "github.com/patrickcping/pingone-go/management"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

type Config struct {
	ClientID      string
	ClientSecret  string
	EnvironmentID string
	Region        string
	ForceDelete   bool
}

type Client struct {
	API          *pingone.APIClient
	RegionSuffix string
	ForceDelete  bool
}

func (c *Config) APIClient(ctx context.Context) (*Client, error) {

	var client *pingone.APIClient

	var regionSuffix string
	switch p1Region := c.Region; p1Region {
	case "EU":
		regionSuffix = "eu"
	case "US":
		regionSuffix = "com"
	case "ASIA":
		regionSuffix = "asia"
	case "CA":
		regionSuffix = "ca"
	default:
		regionSuffix = "com"
	}

	token, err := getToken(ctx, c, regionSuffix)
	if err != nil {
		return nil, err
	}

	clientcfg := pingone.NewConfiguration()
	clientcfg.AddDefaultHeader("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
	client = pingone.NewAPIClient(clientcfg)

	log.Printf("[INFO] PingOne Client using region suffix %s", regionSuffix)

	apiClient := &Client{
		API:          client,
		RegionSuffix: regionSuffix,
		ForceDelete:  c.ForceDelete,
	}

	log.Printf("[INFO] PingOne Client configured")
	return apiClient, nil
}

func getToken(ctx context.Context, c *Config, regionSuffix string) (*oauth2.Token, error) {

	//Get URL from SDK
	authURL := fmt.Sprintf("https://auth.pingone.%s", regionSuffix)
	log.Printf("[INFO] Getting token from %s", authURL)

	//OAuth 2.0 config for client creds
	config := clientcredentials.Config{
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		TokenURL:     fmt.Sprintf("%s/%s/as/token", authURL, c.EnvironmentID),
		AuthStyle:    oauth2.AuthStyleAutoDetect,
	}

	token, err := config.Token(ctx)
	if err != nil {
		return nil, err
	}
	log.Printf("[INFO] Token retrieved")

	return token, nil
}
