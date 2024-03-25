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
	EnvironmentNamePrefix = "tf-testacc-"
)

func SweepClient(ctx context.Context) (*client.Client, error) {

	config := &client.Config{
		ClientID:      os.Getenv("PINGONE_CLIENT_ID"),
		ClientSecret:  os.Getenv("PINGONE_CLIENT_SECRET"),
		EnvironmentID: os.Getenv("PINGONE_ENVIRONMENT_ID"),
		Region:        os.Getenv("PINGONE_REGION"),
		GlobalOptions: &client.GlobalOptions{
			Environment: &client.EnvironmentOptions{
				ProductionTypeForceDelete: true,
			},
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
			return apiClient.EnvironmentsApi.ReadAllEnvironments(ctx).Filter(filter).Execute()
		},
		"ReadAllEnvironments",
		sdk.CustomErrorResourceNotFoundWarning,
		func(ctx context.Context, r *http.Response, p1error *model.P1Error) bool {

			if p1error != nil {
				var err error

				// Permissions may not have propagated by this point
				if m, err := regexp.MatchString("^The request could not be completed. You do not have access to this resource.", p1error.GetMessage()); err == nil && m {
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

	respList := resp.(*management.EntityArray)

	if environments, ok := respList.Embedded.GetEnvironmentsOk(); ok {

		for _, environment := range environments {
			if environment.GetName() == "Administrators" {
				return nil, fmt.Errorf("Unsafe filter, Administrators environment present: %s", filter)
			}
		}
		return environments, nil
	} else {
		return make([]management.Environment, 0), nil
	}

}

func CreateTestEnvironment(ctx context.Context, apiClient *management.APIClient, region management.EnvironmentRegion, index string) error {

	environmentLicense := os.Getenv("PINGONE_LICENSE_ID")

	environment := *management.NewEnvironment(
		*management.NewEnvironmentLicense(environmentLicense),
		fmt.Sprintf("%sdynamic-%s", EnvironmentNamePrefix, index),
		region,
		management.ENUMENVIRONMENTTYPE_SANDBOX,
	)

	productBOMItems := make([]management.BillOfMaterialsProductsInner, 0)

	productBOMItems = append(productBOMItems, *management.NewBillOfMaterialsProductsInner(management.ENUMPRODUCTTYPE_ONE_BASE))
	productBOMItems = append(productBOMItems, *management.NewBillOfMaterialsProductsInner(management.ENUMPRODUCTTYPE_ONE_MFA))
	productBOMItems = append(productBOMItems, *management.NewBillOfMaterialsProductsInner(management.ENUMPRODUCTTYPE_ONE_RISK))
	productBOMItems = append(productBOMItems, *management.NewBillOfMaterialsProductsInner(management.ENUMPRODUCTTYPE_ONE_AUTHORIZE))
	productBOMItems = append(productBOMItems, *management.NewBillOfMaterialsProductsInner(management.ENUMPRODUCTTYPE_ONE_CREDENTIALS))
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
					allowedRegions := make([]string, 0)
					for _, allowedRegion := range details[0].GetInnerError().AllowedValues {
						allowedRegions = append(allowedRegions, model.FindRegionByAPICode(management.EnumRegionCode(allowedRegion)).Region)
					}
					diags := diag.FromErr(fmt.Errorf("Incompatible environment region for the organization tenant.  Expecting regions %v, region provided: %+v", allowedRegions, region))

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
