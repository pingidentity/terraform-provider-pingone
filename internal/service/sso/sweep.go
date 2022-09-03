package sso

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/sweep"
)

func init() {
	resource.AddTestSweepers("pingone_group", &resource.Sweeper{
		Name: "pingone_group",
		F:    sweepGroups,
	})

	resource.AddTestSweepers("pingone_population", &resource.Sweeper{
		Name: "pingone_population",
		F:    sweepPopulations,
		Dependencies: []string{
			"pingone_group",
		},
	})
}

func sweepGroups(region string) error {

	var ctx = context.Background()

	p1Client, err := sweep.SweepClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})

	environments, err := sweep.FetchTaggedEnvironments(ctx, apiClient)
	if err != nil {
		return err
	}

	for _, environment := range environments {

		respGroupsList, _, err := apiClient.GroupsApi.ReadAllGroups(ctx, environment.GetId()).Execute()
		if err != nil {
			return fmt.Errorf("Error getting groups: %s", err)
		}

		if groups, ok := respGroupsList.Embedded.GetGroupsOk(); ok {

			for _, group := range groups {

				_, err := apiClient.GroupsApi.DeleteGroup(ctx, environment.GetId(), group.GetId()).Execute()

				if err != nil {
					log.Printf("Error destroying group %s during sweep: %s", group.GetName(), err)
				}

			}

		}

	}
	return nil

}

func sweepPopulations(region string) error {

	var ctx = context.Background()

	p1Client, err := sweep.SweepClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})

	environments, err := sweep.FetchTaggedEnvironments(ctx, apiClient)
	if err != nil {
		return err
	}

	for _, environment := range environments {

		respPopsList, _, err := apiClient.PopulationsApi.ReadAllPopulations(ctx, environment.GetId()).Execute()
		if err != nil {
			return fmt.Errorf("Error getting populations: %s", err)
		}

		if populations, ok := respPopsList.Embedded.GetPopulationsOk(); ok {

			for _, population := range populations {

				if (population.GetName() != "Default") && (strings.HasPrefix(population.GetName(), "default-")) {

					_, err := apiClient.PopulationsApi.DeletePopulation(ctx, environment.GetId(), population.GetId()).Execute()

					if err != nil {
						log.Printf("Error destroying population %s during sweep: %s", population.GetName(), err)
					}
				}

			}
		}

	}
	return nil

}
