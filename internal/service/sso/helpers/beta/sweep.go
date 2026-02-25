// Copyright © 2026 Ping Identity Corporation

//go:build beta

package beta

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pingidentity/terraform-provider-pingone/internal/sweep"
)

func init() {
	resource.AddTestSweepers("pingone_application_ff_app_import", &resource.Sweeper{
		Name: "pingone_application_ff_app_import",
		F:    sweepFFAppImportApplications,
	})
}

func sweepFFAppImportApplications(region string) error {

	var ctx = context.Background()

	p1Client, err := sweep.SweepClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.ManagementAPIClient

	scimFilter := fmt.Sprintf("name sw \"%s\"", appImportFFSandboxEnvironmentName)

	pagedIterator := apiClient.EnvironmentsApi.ReadAllEnvironments(ctx).Filter(scimFilter).Execute()

	for pageCursor, err := range pagedIterator {
		if err != nil {
			return err
		}

		if environments, ok := pageCursor.EntityArray.Embedded.GetEnvironmentsOk(); ok {
			for _, environment := range environments {

				if environment.GetName() == appImportFFSandboxEnvironmentName {
					pagedIterator := apiClient.ApplicationsApi.ReadAllApplications(ctx, environment.GetId()).Execute()
					for pageCursor, err := range pagedIterator {
						if err != nil {
							return fmt.Errorf("error getting applications: %s", err)
						}

						if applications, ok := pageCursor.EntityArray.Embedded.GetApplicationsOk(); ok {

							for _, application := range applications {

								if v := application.ApplicationOIDC; v != nil {

									_, err := apiClient.ApplicationsApi.DeleteApplication(ctx, environment.GetId(), v.GetId()).Execute()

									if err != nil {
										log.Printf("Error destroying application %s during sweep: %s", v.GetName(), err)
									}
								}
							}
						}
					}
				}
			}
		}
	}

	return nil

}
