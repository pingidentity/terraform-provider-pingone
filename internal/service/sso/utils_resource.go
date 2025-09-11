// Copyright © 2025 Ping Identity Corporation

// Package sso provides utility functions for managing resources in PingOne SSO service configurations.
// This file contains functions for fetching and managing resource data from the PingOne API.
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

// fetchResourceFromID retrieves a specific resource by its ID from the PingOne API.
// It returns the resource object and any diagnostics encountered during the fetch operation.
// The ctx parameter provides context for the API call.
// The apiClient parameter is the PingOne Management API client instance.
// The environmentId parameter specifies the environment containing the resource.
// The resourceId parameter specifies the ID of the resource to fetch.
// The warnIfNotFound parameter controls whether a missing resource generates a warning or error.
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

// fetchResourceFromName retrieves a specific resource by its name from the PingOne API.
// It returns the resource object and any diagnostics encountered during the fetch operation.
// The ctx parameter provides context for the API call.
// The apiClient parameter is the PingOne Management API client instance.
// The environmentId parameter specifies the environment containing the resource.
// The resourceName parameter specifies the name of the resource to fetch.
// The warnIfNotFound parameter controls whether a missing resource generates a warning or error.
func fetchResourceFromName(ctx context.Context, apiClient *management.APIClient, environmentId string, resourceName string, warnIfNotFound bool) (*management.Resource, diag.Diagnostics) {
	var diags diag.Diagnostics

	var resource management.Resource

	resources, d := fetchResources(ctx, apiClient, environmentId, warnIfNotFound)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	if len(resources) > 0 {

		found := false
		for _, resourceItem := range resources {
			if resourceItem.Resource != nil && resourceItem.Resource.GetName() == resourceName {
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

// fetchResourceByType retrieves a specific resource by its type from the PingOne API.
// It returns the resource object and any diagnostics encountered during the fetch operation.
// The ctx parameter provides context for the API call.
// The apiClient parameter is the PingOne Management API client instance.
// The environmentId parameter specifies the environment containing the resource.
// The resourceType parameter specifies the type of resource to fetch (cannot be CUSTOM type).
// The warnIfNotFound parameter controls whether a missing resource generates a warning or error.
func fetchResourceByType(ctx context.Context, apiClient *management.APIClient, environmentId string, resourceType management.EnumResourceType, warnIfNotFound bool) (*management.Resource, diag.Diagnostics) {
	var diags diag.Diagnostics

	if resourceType == management.ENUMRESOURCETYPE_CUSTOM {
		diags.AddError("Invalid resource type", "Cannot find a resource by custom type.")
		return nil, diags
	}

	var resource management.Resource

	resources, d := fetchResources(ctx, apiClient, environmentId, warnIfNotFound)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	if len(resources) > 0 {

		found := false
		for _, resourceItem := range resources {
			if resourceItem.Resource != nil && resourceItem.Resource.GetType() == resourceType {
				resource = *resourceItem.Resource
				found = true
				break
			}
		}

		if !found {
			if warnIfNotFound {
				diags.AddWarning(
					"Cannot find resource from type",
					fmt.Sprintf("The resource %s for environment %s cannot be found", resourceType, environmentId),
				)
			} else {
				diags.AddError(
					"Cannot find resource from type",
					fmt.Sprintf("The resource %s for environment %s cannot be found", resourceType, environmentId),
				)
			}
			return nil, diags
		}

	}

	return &resource, diags
}

// fetchResources retrieves all resources from a specific environment using the PingOne API.
// It returns a slice of resource objects and any diagnostics encountered during the fetch operation.
// The ctx parameter provides context for the API call.
// The apiClient parameter is the PingOne Management API client instance.
// The environmentId parameter specifies the environment containing the resources.
// The warnIfNotFound parameter controls whether a missing environment generates a warning or error.
// This function handles paginated responses to retrieve all available resources.
func fetchResources(ctx context.Context, apiClient *management.APIClient, environmentId string, warnIfNotFound bool) ([]management.EntityArrayEmbeddedResourcesInner, diag.Diagnostics) {
	var diags diag.Diagnostics

	errorFunction := framework.DefaultCustomError
	if warnIfNotFound {
		errorFunction = framework.CustomErrorResourceNotFoundWarning
	}

	// Run the API call
	var resources []management.EntityArrayEmbeddedResourcesInner
	diags.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			pagedIterator := apiClient.ResourcesApi.ReadAllResources(ctx, environmentId).Execute()

			var initialHttpResponse *http.Response

			foundResources := make([]management.EntityArrayEmbeddedResourcesInner, 0)

			for pageCursor, err := range pagedIterator {
				if err != nil {
					return framework.CheckEnvironmentExistsOnPermissionsError(ctx, apiClient, environmentId, nil, pageCursor.HTTPResponse, err)
				}

				if initialHttpResponse == nil {
					initialHttpResponse = pageCursor.HTTPResponse
				}

				if pageCursor.EntityArray.Embedded != nil && pageCursor.EntityArray.Embedded.Resources != nil {
					foundResources = append(foundResources, pageCursor.EntityArray.Embedded.GetResources()...)
				}
			}

			return foundResources, initialHttpResponse, nil
		},
		"ReadAllResources",
		errorFunction,
		sdk.DefaultCreateReadRetryable,
		&resources,
	)...)
	if diags.HasError() {
		return nil, diags
	}

	if resources == nil {
		if warnIfNotFound {
			diags.AddWarning(
				"Environment cannot be found",
				fmt.Sprintf("The environment %s cannot be found when attempting to find resource", environmentId),
			)
		} else {
			diags.AddError(
				"Environment cannot be found",
				fmt.Sprintf("The environment %s cannot be found when attempting to find resource", environmentId),
			)
		}
		return nil, diags
	}

	return resources, diags
}
