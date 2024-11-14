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

	riskPredictorIDsMap := make(map[string]string)
	diags.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			pagedIterator := apiClient.RiskAdvancedPredictorsApi.ReadAllRiskPredictors(ctx, environmentID).Execute()

			var initialHttpResponse *http.Response

			riskPredictorIDsMap := make(map[string]string)

			for pageCursor, err := range pagedIterator {
				if err != nil {
					return framework.CheckEnvironmentExistsOnPermissionsError(ctx, managementApiClient, environmentID, nil, pageCursor.HTTPResponse, err)
				}

				if initialHttpResponse == nil {
					initialHttpResponse = pageCursor.HTTPResponse
				}

				if riskPredictors, ok := pageCursor.EntityArray.Embedded.GetRiskPredictorsOk(); ok {

					for _, riskPredictor := range riskPredictors {

						riskPredictorActualInstance := riskPredictor.GetActualInstance()

						// Get the ID and compact name of the risk predictor
						var predictor struct {
							ID          string `json:"id"`
							CompactName string `json:"compactName"`
						}

						predictorBytes, err := json.Marshal(riskPredictorActualInstance)
						if err != nil {
							return nil, pageCursor.HTTPResponse, err
						}

						err = json.Unmarshal(predictorBytes, &predictor)
						if err != nil {
							return nil, pageCursor.HTTPResponse, err
						}

						// Add the ID to the map of all risk predictors
						if predictor.ID != "" && predictor.CompactName != "" {
							riskPredictorIDsMap[predictor.CompactName] = predictor.ID
						}
					}
				}
			}

			return riskPredictorIDsMap, initialHttpResponse, nil
		},
		"ReadAllRiskPredictors",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
		&riskPredictorIDsMap,
	)...)
	if diags.HasError() {
		return nil, diags
	}

	// Check that all the input risk predictors were found
	returnVar := make([]string, 0)

	for _, predictorCompactName := range predictorCompactNames {
		if _, ok := riskPredictorIDsMap[predictorCompactName]; ok {
			returnVar = append(returnVar, riskPredictorIDsMap[predictorCompactName])
		} else {
			diags.AddError(
				"Cannot find risk predictor from compact name",
				fmt.Sprintf("The risk predictor \"%s\" cannot be found in the environment ID \"%s\".  Please check input parameters and retry.", predictorCompactName, environmentID),
			)
		}
	}

	return returnVar, diags

}
