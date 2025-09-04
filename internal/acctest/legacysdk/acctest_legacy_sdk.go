// Copyright © 2025 Ping Identity Corporation

package legacysdk

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	legacyclient "github.com/pingidentity/terraform-provider-pingone/internal/client"
)

func TestClient(ctx context.Context) (*legacyclient.Client, error) {

	regionCode := management.EnumRegionCode(os.Getenv("PINGONE_REGION_CODE"))

	config := &legacyclient.Config{
		ClientID:      os.Getenv("PINGONE_CLIENT_ID"),
		ClientSecret:  os.Getenv("PINGONE_CLIENT_SECRET"),
		EnvironmentID: os.Getenv("PINGONE_ENVIRONMENT_ID"),
		RegionCode:    &regionCode,
		GlobalOptions: &legacyclient.GlobalOptions{
			Population: &legacyclient.PopulationOptions{
				ContainsUsersForceDelete: false,
			},
		},
	}

	return config.APIClient(ctx, acctest.GetProviderTestingVersion())

}

func PreCheckTestClient(ctx context.Context, t *testing.T) *legacyclient.Client {
	p1Client, err := TestClient(ctx)

	if err != nil {
		t.Fatalf("Failed to get API client: %v", err)
	}

	return p1Client
}

func CheckParentEnvironmentDestroy(ctx context.Context, apiClient *management.APIClient, environmentID string) (bool, error) {
	environment, r, err := apiClient.EnvironmentsApi.ReadOneEnvironment(ctx, environmentID).Execute()

	destroyed, err := acctest.CheckForResourceDestroy(r, err)
	if err != nil {
		return destroyed, err
	}

	if destroyed {
		return destroyed, nil
	} else {
		if environment != nil && environment.Type == management.ENUMENVIRONMENTTYPE_PRODUCTION {
			return true, nil
		} else {
			return false, nil
		}
	}
}

func MinimalSandboxDaVinciEnvironmentNoPopulation(resourceName, licenseID string) string {
	return MinimalDaVinciEnvironmentNoPopulation(resourceName, licenseID, management.ENUMENVIRONMENTTYPE_SANDBOX)
}

func MinimalDaVinciEnvironmentNoPopulation(resourceName, licenseID string, environmentType management.EnumEnvironmentType) string {
	return fmt.Sprintf(`
	resource "pingone_environment" "%[1]s" {
		name = "%[1]s"
		license_id = "%[2]s"
		type = "%[3]s"

	services = [
		{
			type = "SSO"
		},
		{
			type = "MFA"
		},
		{
			type = "Risk"
		},
		{
			type = "Credentials"
		},
		{
			type = "Verify"
		},
		{
		    type = "DaVinci"
		}
	]
}
`, resourceName, licenseID, string(environmentType))
}

func MinimalSandboxEnvironmentNoPopulation(resourceName, licenseID string) string {
	return MinimalEnvironmentNoPopulation(resourceName, licenseID, management.ENUMENVIRONMENTTYPE_SANDBOX)
}

func MinimalEnvironmentNoPopulation(resourceName, licenseID string, environmentType management.EnumEnvironmentType) string {
	return fmt.Sprintf(`
	resource "pingone_environment" "%[1]s" {
		name = "%[1]s"
		license_id = "%[2]s"
		type = "%[3]s"

	services = [
		{
			type = "SSO"
		},
		{
			type = "MFA"
		},
		{
			type = "Risk"
		},
		{
			type = "Credentials"
		},
		{
			type = "Verify"
		}
	]
}
`, resourceName, licenseID, string(environmentType))
}

func MinimalSandboxEnvironment(resourceName, licenseID string) string {
	return fmt.Sprintf(`
	%[1]s
		
	resource "pingone_population_default" "%[2]s" {
		environment_id = pingone_environment.%[2]s.id

		name = "%[2]s"
	}
`, MinimalSandboxEnvironmentNoPopulation(resourceName, licenseID), resourceName)
}

func MinimalSandboxDaVinciEnvironment(resourceName, licenseID string) string {
	return fmt.Sprintf(`
	%[1]s
		
	resource "pingone_population_default" "%[2]s" {
		environment_id = pingone_environment.%[2]s.id

		name = "%[2]s"
	}
`, MinimalSandboxDaVinciEnvironmentNoPopulation(resourceName, licenseID), resourceName)
}

func CheckParentUserDestroy(ctx context.Context, apiClient *management.APIClient, environmentID, userID string) (bool, error) {
	_, r, err := apiClient.UsersApi.ReadUser(ctx, environmentID, userID).Execute()

	return acctest.CheckForResourceDestroy(r, err)
}
