package risk

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/patrickcping/pingone-go-sdk-v2/risk"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
)

// Data source TODO

func predictorFetchIDsFromCompactNames(ctx context.Context, apiClient *risk.APIClient, environmentID string, predictorCompactNames []string) ([]string, diag.Diagnostics) {
	var diags diag.Diagnostics

	var entityArray *risk.EntityArray
	diags.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			return apiClient.RiskAdvancedPredictorsApi.ReadAllRiskPredictors(ctx, environmentID).Execute()
		},
		"ReadAllRiskPredictors",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
		&entityArray,
	)...)
	if diags.HasError() {
		return nil, diags
	}

	riskPredictorIDs := make(map[string]string)

	if riskPredictors, ok := entityArray.Embedded.GetRiskPredictorsOk(); ok {

		for _, riskPredictor := range riskPredictors {

			riskPredictorActualInstance := riskPredictor.GetActualInstance()

			var id string
			var compactName string
			switch v := riskPredictorActualInstance.(type) {
			case *risk.RiskPredictorAnonymousNetwork:
				id = v.GetId()
				compactName = v.GetCompactName()
			case *risk.RiskPredictorComposite:
				id = v.GetId()
				compactName = v.GetCompactName()
			case *risk.RiskPredictorCustom:
				id = v.GetId()
				compactName = v.GetCompactName()
			case *risk.RiskPredictorDevice:
				id = v.GetId()
				compactName = v.GetCompactName()
			case *risk.RiskPredictorGeovelocity:
				id = v.GetId()
				compactName = v.GetCompactName()
			case *risk.RiskPredictorIPReputation:
				id = v.GetId()
				compactName = v.GetCompactName()
			case *risk.RiskPredictorUserLocationAnomaly:
				id = v.GetId()
				compactName = v.GetCompactName()
			case *risk.RiskPredictorUserRiskBehavior:
				id = v.GetId()
				compactName = v.GetCompactName()
			case *risk.RiskPredictorVelocity:
				id = v.GetId()
				compactName = v.GetCompactName()
			}

			// Add the ID to the map of all risk predictors
			if id != "" && compactName != "" {
				riskPredictorIDs[compactName] = id
			}
		}
	}

	// Check that all the input risk predictors were found
	returnVar := make([]string, 0)

	for _, predictorCompactName := range predictorCompactNames {
		if _, ok := riskPredictorIDs[predictorCompactName]; ok {
			returnVar = append(returnVar, riskPredictorIDs[predictorCompactName])
		} else {
			diags.AddError(
				"Cannot find risk predictor from compact name",
				fmt.Sprintf("The risk predictor \"%s\" cannot be found in the environment ID \"%s\".  Please check input parameters and retry.", predictorCompactName, environmentID),
			)
		}
	}

	return returnVar, diags

}
