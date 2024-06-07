package sso

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
)

func fetchResourceFromID(ctx context.Context, apiClient *management.APIClient, environmentId, resourceId string, warnIfNotFound bool) (*management.Resource, diag.Diagnostics) {
	var diags diag.Diagnostics

	errorFunction := framework.DefaultCustomError

	if warnIfNotFound {
		errorFunction = framework.CustomErrorResourceNotFoundWarning
	}

	var resource *management.Resource
	diags.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := apiClient.ResourcesApi.ReadOneResource(ctx, environmentId, resourceId).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, apiClient, environmentId, fO, fR, fErr)
		},
		"ReadOneResource",
		errorFunction,
		sdk.DefaultCreateReadRetryable,
		&resource,
	)...)
	if diags.HasError() {
		return nil, diags
	}

	return resource, diags
}

func fetchResourceFromName(ctx context.Context, apiClient *management.APIClient, environmentId string, resourceName string, warnIfNotFound bool) (*management.Resource, diag.Diagnostics) {
	var diags diag.Diagnostics

	var resource management.Resource

	errorFunction := framework.DefaultCustomError
	if warnIfNotFound {
		errorFunction = framework.CustomErrorResourceNotFoundWarning
	}

	// Run the API call
	var entityArray *management.EntityArray
	diags.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := apiClient.ResourcesApi.ReadAllResources(ctx, environmentId).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, apiClient, environmentId, fO, fR, fErr)
		},
		"ReadAllResources",
		errorFunction,
		sdk.DefaultCreateReadRetryable,
		&entityArray,
	)...)
	if diags.HasError() {
		return nil, diags
	}

	if entityArray == nil {
		if warnIfNotFound {
			diags.AddWarning(
				"Environment cannot be found",
				fmt.Sprintf("The environment %s cannot be found when finding resource %s by name", environmentId, resourceName),
			)
		} else {
			diags.AddError(
				"Environment cannot be found",
				fmt.Sprintf("The environment %s cannot be found when finding resource %s by name", environmentId, resourceName),
			)
		}
		return nil, diags
	}

	if resources, ok := entityArray.Embedded.GetResourcesOk(); ok {

		found := false
		for _, resourceItem := range resources {

			if resourceItem.Resource.GetName() == resourceName {
				resource = *resourceItem.Resource
				found = true
				break
			}
		}

		if !found {
			if warnIfNotFound {
				diags.AddWarning(
					"Cannot find resource from name",
					fmt.Sprintf("The resource %s for environment %s cannot be found", resourceName, environmentId),
				)
			} else {
				diags.AddError(
					"Cannot find resource from name",
					fmt.Sprintf("The resource %s for environment %s cannot be found", resourceName, environmentId),
				)
			}
			return nil, diags
		}

	}

	return &resource, diags
}
