// Copyright © 2025 Ping Identity Corporation

package sso

import (
	"context"
	"fmt"
	"net/http"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func populationDeleteCustomErrorHandler(r *http.Response, p1Error *model.P1Error) diag.Diagnostics {
	var diags diag.Diagnostics

	if p1Error != nil {
		// Env must contain at least one population
		if details, ok := p1Error.GetDetailsOk(); ok && details != nil && len(details) > 0 {
			if code, ok := details[0].GetCodeOk(); ok && *code == "CONSTRAINT_VIOLATION" {
				if message, ok := details[0].GetMessageOk(); ok {
					m, err := regexp.MatchString(`must contain at least one population`, *message)
					if err == nil && m {
						diags.AddWarning(
							"Constraint violation",
							fmt.Sprintf("A constraint violation error was encountered: %s\n\nThe population has been removed from Terraform state, but has been left in place in the environment.", p1Error.GetMessage()),
						)

						return diags
					}
				}
			}
		}
	}

	return diags
}

func (r *populationResource) hasUsersAssigned(ctx context.Context, environmentID, populationID string) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	users, d := r.readUsers(ctx, environmentID, populationID)
	diags.Append(d...)
	if diags.HasError() {
		return false, diags
	}

	if len(users) > 0 {
		return true, diags
	}

	return false, diags
}

func (r *populationResource) readUsers(ctx context.Context, environmentID, populationID string) ([]management.User, diag.Diagnostics) {
	var diags diag.Diagnostics

	m, err := regexp.MatchString(verify.P1ResourceIDRegexpFullString.String(), populationID)
	if err != nil {
		diags.AddError(
			"Population ID validation",
			fmt.Sprintf("An error occurred while validating the population ID: %s", err.Error()),
		)
		return nil, diags
	}

	if m {

		scimFilter := fmt.Sprintf(`population.id eq "%s"`, populationID)

		// Run the API call
		var users []management.User
		diags.Append(framework.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				pagedIterator := r.Client.ManagementAPIClient.UsersApi.ReadAllUsers(ctx, environmentID).Filter(scimFilter).Execute()

				var initialHttpResponse *http.Response

				foundUsers := make([]management.User, 0)

				for pageCursor, err := range pagedIterator {
					if err != nil {
						return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, environmentID, nil, pageCursor.HTTPResponse, err)
					}

					if initialHttpResponse == nil {
						initialHttpResponse = pageCursor.HTTPResponse
					}

					if pageCursor.EntityArray.Embedded != nil && pageCursor.EntityArray.Embedded.Users != nil {
						foundUsers = append(foundUsers, pageCursor.EntityArray.Embedded.GetUsers()...)
					}
				}

				return foundUsers, initialHttpResponse, nil
			},
			"ReadAllUsers",
			framework.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
			&users,
		)...)

		if diags.HasError() {
			return nil, diags
		}

		return users, nil
	}

	if r.options.Population.ContainsUsersForceDelete {
		diags.AddError(
			"Data protection notice",
			fmt.Sprintf("For data protection reasons, it could not be determined whether users exist in the population %[2]s in environment %[1]s. Any users in this population will not be deleted.", environmentID, populationID),
		)
	}

	return nil, diags
}

func (r *populationResource) checkEnvironmentControls(ctx context.Context, environmentID, populationID string) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	if r.options.Population.ContainsUsersForceDelete {
		// Check if the environment is a sandbox type.  We'll only delete users in sandbox environments
		var environmentResponse *management.Environment
		diags.Append(framework.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				fO, fR, fErr := r.Client.ManagementAPIClient.EnvironmentsApi.ReadOneEnvironment(ctx, environmentID).Execute()
				return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, environmentID, fO, fR, fErr)
			},
			"ReadOneEnvironment-DeletePopulation",
			framework.DefaultCustomError,
			nil,
			&environmentResponse,
		)...)
		if diags.HasError() {
			return false, diags
		}

		if v, ok := environmentResponse.GetTypeOk(); ok && *v == management.ENUMENVIRONMENTTYPE_SANDBOX {
			return true, diags
		} else {
			diags.AddWarning(
				"Data protection notice",
				fmt.Sprintf("For data protection reasons, the provider configuration `global_options.population.contains_users_force_delete` has no effect on environment ID %[1]s as it has a type set to `PRODUCTION`.  Users in this population will not be deleted.\n"+
					"If you wish to force delete population %[2]s in environment %[1]s, please review and remove user data manually.", environmentID, populationID),
			)
		}
	}

	return false, diags
}

func (r *populationResource) checkControlsAndDeletePopulationUsers(ctx context.Context, environmentID, populationID string) diag.Diagnostics {
	var diags diag.Diagnostics

	environmentControlsOk, d := r.checkEnvironmentControls(ctx, environmentID, populationID)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}

	if environmentControlsOk {

		loopCounter := 1
		for loopCounter > 0 {

			users, d := r.readUsers(ctx, environmentID, populationID)
			diags.Append(d...)
			if diags.HasError() {
				return diags
			}

			// DELETE USERS
			if len(users) == 0 {
				break
			} else {
				for _, user := range users {
					var entityArray *management.EntityArray
					diags.Append(framework.ParseResponse(
						ctx,

						func() (any, *http.Response, error) {
							fR, fErr := r.Client.ManagementAPIClient.UsersApi.DeleteUser(ctx, environmentID, user.GetId()).Execute()
							return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, environmentID, nil, fR, fErr)
						},
						"DeleteUser-DeletePopulation",
						framework.DefaultCustomError,
						nil,
						&entityArray,
					)...)

					if diags.HasError() {
						return diags
					}
				}
			}
		}
	}

	return diags
}
