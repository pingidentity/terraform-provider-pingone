package base

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	pingone "github.com/patrickcping/pingone-go/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/sweep"
)

func init() {
	resource.AddTestSweepers("pingone_environment", &resource.Sweeper{
		Name: "pingone_environment",
		F:    sweepEnvironments,
	})
}

func sweepEnvironments(region string) error {

	var ctx = context.Background()

	p1Client, err := sweep.SweepClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API
	ctx = context.WithValue(ctx, pingone.ContextServerVariables, map[string]string{
		"suffix": p1Client.RegionSuffix,
	})

	respList, _, err := apiClient.EnvironmentsApi.ReadAllEnvironments(ctx).Execute()
	if err != nil {
		return fmt.Errorf("Error getting environments: %s", err)
	}

	if environments, ok := respList.Embedded.GetEnvironmentsOk(); ok {

		for _, environment := range environments {

			if (environment.GetName() != "Administrators") && (strings.HasPrefix(environment.GetName(), "tf-testacc-")) {

				// Reset back to sandbox
				if environment.GetType() == "PRODUCTION" {
					updateEnvironmentTypeRequest := *pingone.NewUpdateEnvironmentTypeRequest()
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
		}

	}
	return nil

}
