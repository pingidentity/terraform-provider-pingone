package client

import (
	"fmt"
	"strings"

	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	"golang.org/x/exp/slices"
)

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

func (c *Config) Validate() error {

	err := fmt.Errorf("Must provide the region parameter and either client_id, client_secret and environment_id, or api_access_token.")

	// Missing attrs check
	if c.Region == "" {
		return err
	}

	if (c.ClientID == "" || c.ClientSecret == "" || c.EnvironmentID == "") && c.AccessToken == "" {
		return err
	}

	// Conflicting attrs check
	if (c.ClientID != "" || c.ClientSecret != "" || c.EnvironmentID != "") && c.AccessToken != "" {
		return fmt.Errorf("api_access_token cannot be set with client_id, client_secret or environment_id")
	}

	// Region data
	if !slices.Contains(model.RegionsAvailableList(), c.Region) {
		return fmt.Errorf("Invalid region value.  The region parameter is case sensitive and must be one of the following values: %s", strings.Join(model.RegionsAvailableList(), ", "))
	}

	return nil
}
