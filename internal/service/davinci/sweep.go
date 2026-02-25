// Copyright © 2026 Ping Identity Corporation

package davinci

import (
	"context"
	"log"
	"os"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/sweep"
)

func init() {
	resource.AddTestSweepers("pingone_environment-davinci", &resource.Sweeper{
		Name: "pingone_environment-davinci",
		F:    sweepCreateDavinciTestEnvironment,
		Dependencies: []string{
			// Only run this sweeper after other environments have been cleared
			"pingone_environment",
		},
	})
}

func sweepCreateDavinciTestEnvironment(regionString string) error {

	var ctx = context.Background()

	p1Client, err := sweep.SweepClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.ManagementAPIClient

	var region management.EnvironmentRegion
	if v := os.Getenv("PINGONE_TERRAFORM_REGION_OVERRIDE"); v != "" {
		region = management.EnvironmentRegion{
			String: &v,
		}
	} else {
		region = management.EnvironmentRegion{
			EnumRegionCode: &p1Client.API.Region.APICode,
		}
	}

	err = sweep.CreateTestEnvironment(ctx, apiClient, region, "davinci-bootstrapped-test", true)
	if err != nil {
		log.Printf("Error creating environment `davinci-bootstrapped-test` during sweep: %s", err)
	}

	return nil

}
