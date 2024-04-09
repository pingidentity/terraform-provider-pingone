package risk

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/patrickcping/pingone-go-sdk-v2/risk"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
)

// Data source TODO

func riskPredictorFetchIDsFromCompactNames(ctx context.Context, apiClient *risk.APIClient, managementApiClient *management.APIClient, environmentID string, predictorCompactNames []string) ([]string, diag.Diagnostics) {
	var diags diag.Diagnostics

	var entityArray *risk.EntityArray
	diags.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := apiClient.RiskAdvancedPredictorsApi.ReadAllRiskPredictors(ctx, environmentID).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, managementApiClient, environmentID, fO, fR, fErr)
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

			// Get the ID and compact name of the risk predictor
			var predictor struct {
				ID          string `json:"id"`
				CompactName string `json:"compactName"`
			}

			predictorBytes, err := json.Marshal(riskPredictorActualInstance)
			if err != nil {
				diags.AddError(
					"Cannot marshal risk predictor",
					fmt.Sprintf("Error marshalling risk predictor: %s", err),
				)
				return nil, diags
			}

			err = json.Unmarshal(predictorBytes, &predictor)
			if err != nil {
				diags.AddError(
					"Cannot unmarshal risk predictor",
					fmt.Sprintf("Error unmarshalling risk predictor: %s", err),
				)
				return nil, diags
			}

			// Add the ID to the map of all risk predictors
			if predictor.ID != "" && predictor.CompactName != "" {
				riskPredictorIDs[predictor.CompactName] = predictor.ID
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
