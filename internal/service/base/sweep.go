// Copyright © 2025 Ping Identity Corporation

package base

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/sweep"
)

func init() {
	resource.AddTestSweepers("pingone_environment", &resource.Sweeper{
		Name: "pingone_environment",
		F:    sweepEnvironments,
		Dependencies: []string{
			"pingone_group",
			"pingone_population",
		},
	})
}

func sweepEnvironments(regionString string) error {

	var ctx = context.Background()

	p1Client, err := sweep.SweepClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.ManagementAPIClient

	environments, err := sweep.FetchTaggedEnvironmentsByPrefix(ctx, apiClient, fmt.Sprintf("%sdynamic-", sweep.EnvironmentNamePrefix))
	if err != nil {
		return err
	}

	for _, environment := range environments {

		// Reset back to sandbox
		if environment.GetType() == "PRODUCTION" {
			updateEnvironmentTypeRequest := *management.NewUpdateEnvironmentTypeRequest()
			updateEnvironmentTypeRequest.SetType("SANDBOX")
			_, _, err := apiClient.EnvironmentsApi.UpdateEnvironmentType(ctx, environment.GetId()).UpdateEnvironmentTypeRequest(updateEnvironmentTypeRequest).Execute()

			if err != nil {
				log.Printf("Error setting environment %s of type PRODUCTION to SANDBOX during sweep: %s", environment.GetName(), err)
			}
		}

		// Delete the environment
		_, err := apiClient.EnvironmentsApi.DeleteEnvironment(ctx, environment.GetId()).Execute()

		if err != nil {
			log.Printf("Error destroying environment %s during sweep: %s", environment.GetName(), err)
		}

	}

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

	err = sweep.CreateTestEnvironment(ctx, apiClient, region, "general-test")
	if err != nil {
		log.Printf("Error creating environment `general-test` during sweep: %s", err)
	}

	return nil

}
