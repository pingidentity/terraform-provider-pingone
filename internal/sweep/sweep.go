// Copyright © 2026 Ping Identity Corporation

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
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/legacysdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
)

var (
	EnvironmentNamePrefix = "tf-testacc-"
)

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

func getProviderTestingVersion() string {
	returnVar := "dev"
	if v := os.Getenv("PINGONE_TESTING_PROVIDER_VERSION"); v != "" {
		returnVar = v
	}
	return returnVar
}

func FetchTaggedEnvironments(ctx context.Context, apiClient *management.APIClient) ([]management.Environment, error) {
	return FetchTaggedEnvironmentsByPrefix(ctx, apiClient, EnvironmentNamePrefix)
}

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
							return nil, nil, fmt.Errorf("unsafe filter, Administrators environment present: %s", filter)
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
		return nil, fmt.Errorf("error getting environments for sweep")
	}

	respList := resp.([]management.Environment)

	return respList, nil
}

func CreateTestEnvironment(ctx context.Context, apiClient *management.APIClient, region management.EnvironmentRegion, index string, davinciBootstrapped bool) error {

	environmentLicense := os.Getenv("PINGONE_LICENSE_ID")

	environment := *management.NewEnvironment(
		*management.NewEnvironmentLicense(environmentLicense),
		fmt.Sprintf("%sdynamic-%s", EnvironmentNamePrefix, index),
		region,
		management.ENUMENVIRONMENTTYPE_SANDBOX,
	)

	productBOMItems := make([]management.BillOfMaterialsProductsInner, 0)

	daVinciService := management.NewBillOfMaterialsProductsInner(management.ENUMPRODUCTTYPE_ONE_DAVINCI)
	if !davinciBootstrapped {
		daVinciService.SetTags([]management.EnumBillOfMaterialsProductTags{management.ENUMBILLOFMATERIALSPRODUCTTAGS_DAVINCI_MINIMAL})
	}

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
					diags := diag.FromErr(fmt.Errorf("incompatible environment region for the organization tenant.  Expecting regions %v, region provided: %+v", details[0].GetInnerError().AllowedValues, region))

					return diags
				}
			}

			return nil
		},
		nil,
	)
	if diags.HasError() {
		return fmt.Errorf("cannot create environment `%s`", environment.GetName())
	}

	environmentID := resp.(*management.Environment).GetId()

	// A population, because we must have one

	population := *management.NewPopulation("Default")

	_, diags = sdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := apiClient.PopulationsApi.CreatePopulation(ctx, environmentID).Population(population).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, apiClient, environmentID, fO, fR, fErr)
		},
		"CreatePopulation",
		sdk.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
	)
	if diags.HasError() {
		return fmt.Errorf("cannot create population for environment `%s`", environment.GetName())
	}

	return nil

}
