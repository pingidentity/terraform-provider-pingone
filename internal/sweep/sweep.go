// Copyright © 2025 Ping Identity Corporation

// Package sweep provides utilities for cleaning up test resources during acceptance testing.
// This package contains functions for creating test environments, fetching tagged test resources,
// and configuring clients for resource cleanup operations.
package sweep

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"regexp"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
)

var (
	// EnvironmentNamePrefix is the standard prefix used for test environment names.
	// This prefix helps identify environments created during acceptance testing
	// so they can be properly cleaned up after test completion.
	EnvironmentNamePrefix = "tf-testacc-"
)

// SweepClient creates and configures a PingOne API client for resource cleanup operations.
// It returns a configured client instance that can be used to identify and clean up test resources.
// The client is configured using environment variables and includes global options suitable for testing,
// such as enabling force deletion of populations that contain users.
// This function is typically used in acceptance test sweep functions to obtain an API client for cleanup.
func SweepClient(ctx context.Context) (*client.Client, error) {

	regionCode := management.EnumRegionCode(os.Getenv("PINGONE_REGION_CODE"))

	config := &client.Config{
		ClientID:      os.Getenv("PINGONE_CLIENT_ID"),
		ClientSecret:  os.Getenv("PINGONE_CLIENT_SECRET"),
		EnvironmentID: os.Getenv("PINGONE_ENVIRONMENT_ID"),
		RegionCode:    &regionCode,
		GlobalOptions: &client.GlobalOptions{
			Population: &client.PopulationOptions{
				ContainsUsersForceDelete: true,
			},
		},
	}

	return config.APIClient(ctx, getProviderTestingVersion())

}

// getProviderTestingVersion returns the provider version string for testing purposes.
// It checks for the PINGONE_TESTING_PROVIDER_VERSION environment variable and returns
// its value if set, otherwise defaults to "dev".
// This version string is used in User-Agent headers for API requests during testing.
func getProviderTestingVersion() string {
	returnVar := "dev"
	if v := os.Getenv("PINGONE_TESTING_PROVIDER_VERSION"); v != "" {
		returnVar = v
	}
	return returnVar
}

// FetchTaggedEnvironments retrieves all environments with names starting with the standard test prefix.
// It returns a slice of environments that were created for testing purposes and may need cleanup.
// This function uses the default EnvironmentNamePrefix to identify test environments.
// The returned environments can be used in sweep operations to clean up test resources.
func FetchTaggedEnvironments(ctx context.Context, apiClient *management.APIClient) ([]management.Environment, error) {
	return FetchTaggedEnvironmentsByPrefix(ctx, apiClient, EnvironmentNamePrefix)
}

// FetchTaggedEnvironmentsByPrefix retrieves all environments with names starting with the specified prefix.
// It returns a slice of environments that match the prefix pattern and may need cleanup.
// The prefix parameter allows filtering for specific test environment naming patterns.
// This function includes safety checks to prevent accidental deletion of the "Administrators" environment.
// The function uses pagination to retrieve all matching environments and includes retry logic for permission issues.
func FetchTaggedEnvironmentsByPrefix(ctx context.Context, apiClient *management.APIClient, prefix string) ([]management.Environment, error) {

	filter := fmt.Sprintf("name sw \"%s\"", prefix)

	resp, diags := sdk.ParseResponse(
		ctx,
		func() (any, *http.Response, error) {
			pagedIterator := apiClient.EnvironmentsApi.ReadAllEnvironments(ctx).Filter(filter).Execute()

			returnEnvironments := make([]management.Environment, 0)

			var initialHttpResponse *http.Response

			for pageCursor, err := range pagedIterator {
				if err != nil {
					return nil, pageCursor.HTTPResponse, err
				}

				if initialHttpResponse == nil {
					initialHttpResponse = pageCursor.HTTPResponse
				}

				if environments, ok := pageCursor.EntityArray.Embedded.GetEnvironmentsOk(); ok {

					for _, environment := range environments {
						if environment.GetName() == "Administrators" {
							return nil, nil, fmt.Errorf("Unsafe filter, Administrators environment present: %s", filter)
						}
					}

					returnEnvironments = append(returnEnvironments, environments...)
				}
			}

			return returnEnvironments, initialHttpResponse, nil
		},
		"ReadAllEnvironments",
		sdk.CustomErrorResourceNotFoundWarning,
		func(ctx context.Context, r *http.Response, p1error *model.P1Error) bool {

			if p1error != nil {
				var err error

				// Permissions may not have propagated by this point
				m, err := regexp.MatchString("^The request could not be completed. You do not have access to this resource.", p1error.GetMessage())
				if err == nil && m {
					tflog.Warn(ctx, "Insufficient PingOne privileges detected")
					return true
				}
				if err != nil {
					tflog.Warn(ctx, "Cannot match error string for retry")
					return false
				}

			}

			return false
		},
	)
	if diags.HasError() {
		return nil, fmt.Errorf("Error getting environments for sweep")
	}

	respList := resp.([]management.Environment)

	return respList, nil
}

// CreateTestEnvironment creates a new PingOne environment configured for testing purposes.
// It returns an error if environment creation fails for any reason.
// The apiClient parameter is used to make API calls to the PingOne platform.
// The region parameter specifies the geographic region where the environment should be created.
// The index parameter is used to generate a unique environment name for parallel test execution.
// The created environment includes all PingOne services enabled and a default population for testing.
// Environment credentials are obtained from the PINGONE_LICENSE_ID environment variable.
func CreateTestEnvironment(ctx context.Context, apiClient *management.APIClient, region management.EnvironmentRegion, index string) error {

	environmentLicense := os.Getenv("PINGONE_LICENSE_ID")

	environment := *management.NewEnvironment(
		*management.NewEnvironmentLicense(environmentLicense),
		fmt.Sprintf("%sdynamic-%s", EnvironmentNamePrefix, index),
		region,
		management.ENUMENVIRONMENTTYPE_SANDBOX,
	)

	productBOMItems := make([]management.BillOfMaterialsProductsInner, 0)

	daVinciService := management.NewBillOfMaterialsProductsInner(management.ENUMPRODUCTTYPE_ONE_DAVINCI)
	daVinciService.SetTags([]management.EnumBillOfMaterialsProductTags{management.ENUMBILLOFMATERIALSPRODUCTTAGS_DAVINCI_MINIMAL})

	productBOMItems = append(productBOMItems, *management.NewBillOfMaterialsProductsInner(management.ENUMPRODUCTTYPE_ONE_AUTHORIZE))
	productBOMItems = append(productBOMItems, *management.NewBillOfMaterialsProductsInner(management.ENUMPRODUCTTYPE_ONE_BASE))
	productBOMItems = append(productBOMItems, *management.NewBillOfMaterialsProductsInner(management.ENUMPRODUCTTYPE_ONE_CREDENTIALS))
	productBOMItems = append(productBOMItems, *daVinciService)
	productBOMItems = append(productBOMItems, *management.NewBillOfMaterialsProductsInner(management.ENUMPRODUCTTYPE_ONE_MFA))
	productBOMItems = append(productBOMItems, *management.NewBillOfMaterialsProductsInner(management.ENUMPRODUCTTYPE_ONE_RISK))
	productBOMItems = append(productBOMItems, *management.NewBillOfMaterialsProductsInner(management.ENUMPRODUCTTYPE_ONE_VERIFY))

	environment.SetBillOfMaterials(*management.NewBillOfMaterials(productBOMItems))

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			return apiClient.EnvironmentsApi.CreateEnvironmentActiveLicense(ctx).Environment(environment).Execute()
		},
		"CreateEnvironmentActiveLicense",
		func(error model.P1Error) diag.Diagnostics {

			// Invalid region
			if details, ok := error.GetDetailsOk(); ok && details != nil && len(details) > 0 {
				if target, ok := details[0].GetTargetOk(); ok && *target == "region" {
					diags := diag.FromErr(fmt.Errorf("Incompatible environment region for the organization tenant.  Expecting regions %v, region provided: %+v", details[0].GetInnerError().AllowedValues, region))

					return diags
				}
			}

			return nil
		},
		nil,
	)
	if diags.HasError() {
		return fmt.Errorf("Cannot create environment `%s`", environment.GetName())
	}

	environmentID := resp.(*management.Environment).GetId()

	// A population, because we must have one

	population := *management.NewPopulation("Default")

	_, diags = sdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := apiClient.PopulationsApi.CreatePopulation(ctx, environmentID).Population(population).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, apiClient, environmentID, fO, fR, fErr)
		},
		"CreatePopulation",
		sdk.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
	)
	if diags.HasError() {
		return fmt.Errorf("Cannot create population for environment `%s`", environment.GetName())
	}

	return nil

}
