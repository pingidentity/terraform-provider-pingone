package sso

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
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
			"pingone_sign_on_policy",
		},
	})

	resource.AddTestSweepers("pingone_sign_on_policy", &resource.Sweeper{
		Name: "pingone_sign_on_policy",
		F:    sweepSOPs,
	})
}

func sweepGroups(region string) error {

	var ctx = context.Background()

	p1Client, err := sweep.SweepClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.ManagementAPIClient

	environments, err := sweep.FetchTaggedEnvironments(ctx, apiClient)
	if err != nil {
		return err
	}

	for _, environment := range environments {

		pagedIterator := apiClient.GroupsApi.ReadAllGroups(ctx, environment.GetId()).Execute()
		for pageCursor, err := range pagedIterator {
			if err != nil {
				return fmt.Errorf("Error getting groups: %s", err)
			}

			if groups, ok := pageCursor.EntityArray.Embedded.GetGroupsOk(); ok {

				for _, group := range groups {

					_, err := apiClient.GroupsApi.DeleteGroup(ctx, environment.GetId(), group.GetId()).Execute()

					if err != nil {
						log.Printf("Error destroying group %s during sweep: %s", group.GetName(), err)
					}

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

	environments, err := sweep.FetchTaggedEnvironments(ctx, apiClient)
	if err != nil {
		return err
	}

	for _, environment := range environments {

		pagedIterator := apiClient.PopulationsApi.ReadAllPopulations(ctx, environment.GetId()).Execute()
		for pageCursor, err := range pagedIterator {
			if err != nil {
				return fmt.Errorf("Error getting populations: %s", err)
			}

			if populations, ok := pageCursor.EntityArray.Embedded.GetPopulationsOk(); ok {

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

	}
	return nil

}

func sweepSOPs(region string) error {

	var ctx = context.Background()

	p1Client, err := sweep.SweepClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.ManagementAPIClient

	environments, err := sweep.FetchTaggedEnvironments(ctx, apiClient)
	if err != nil {
		return err
	}

	for _, environment := range environments {

		pagedIterator := apiClient.SignOnPoliciesApi.ReadAllSignOnPolicies(ctx, environment.GetId()).Execute()
		for pageCursor, err := range pagedIterator {
			if err != nil {
				return fmt.Errorf("Error getting sign on policies: %s", err)
			}

			if signOnPolicies, ok := pageCursor.EntityArray.Embedded.GetSignOnPoliciesOk(); ok {

				for _, signOnPolicy := range signOnPolicies {

					_, err := apiClient.SignOnPoliciesApi.DeleteSignOnPolicy(ctx, environment.GetId(), signOnPolicy.GetId()).Execute()

					if err != nil {
						log.Printf("Error destroying sign-on policy %s during sweep: %s", signOnPolicy.GetName(), err)
					}

				}
			}
		}

	}
	return nil

}
