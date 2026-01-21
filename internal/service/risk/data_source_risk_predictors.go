// Copyright © 2026 Ping Identity Corporation

package risk

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/patrickcping/pingone-go-sdk-v2/risk"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/legacysdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
)

// Types
type RiskPredictorsDataSource serviceClientType

type RiskPredictorsDataSourceModel struct {
	Id            pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	Ids           types.List                   `tfsdk:"ids"`
}

// Framework interfaces
var (
	_ datasource.DataSource              = &RiskPredictorsDataSource{}
	_ datasource.DataSourceWithConfigure = &RiskPredictorsDataSource{}
)

func NewRiskPredictorsDataSource() datasource.DataSource {
	return &RiskPredictorsDataSource{}
}

// Metadata
func (r *RiskPredictorsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_risk_predictors"
}

// Schema
func (r *RiskPredictorsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Datasource to retrieve the IDs of multiple PingOne Risk Predictors.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": schema.StringAttribute{
				Description: "The ID of the environment to retrieve the risk predictors from.",
				Required:    true,
				CustomType:  pingonetypes.ResourceIDType{},
			},

			"ids": framework.Attr_DataSourceReturnIDs(framework.SchemaAttributeDescriptionFromMarkdown(
				"The list of resulting IDs of risk predictors that have been successfully retrieved.",
			)),
		},
	}
}

func (r *RiskPredictorsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	resourceConfig, ok := req.ProviderData.(legacysdk.ResourceType)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected the provider client, got: %T. Please report this to the provider maintainers.", req.ProviderData),
		)

		return
	}

	r.Client = resourceConfig.Client.API
	if r.Client == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this to the provider maintainers.",
		)
		return
	}
}

func (r *RiskPredictorsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data RiskPredictorsDataSourceModel

	if r.Client == nil || r.Client.RiskAPIClient == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this to the provider maintainers.",
		)
		return
	}

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var riskPredictorIDs []string
	resp.Diagnostics.Append(legacysdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			pagedIterator := r.Client.RiskAPIClient.RiskAdvancedPredictorsApi.ReadAllRiskPredictors(ctx, data.EnvironmentId.ValueString()).Execute()

			var initialHttpResponse *http.Response

			for pageCursor, err := range pagedIterator {
				if err != nil {
					return nil, pageCursor.HTTPResponse, err
				}

				if initialHttpResponse == nil {
					initialHttpResponse = pageCursor.HTTPResponse
				}

				if riskPredictors, ok := pageCursor.EntityArray.Embedded.GetRiskPredictorsOk(); ok {
					for _, riskPredictor := range riskPredictors {

						// Get the ID of the risk predictor
						var predictor struct {
							ID string `json:"id"`
						}

						riskPredictorActualInstance := riskPredictor.GetActualInstance()
						predictorBytes, err := json.Marshal(riskPredictorActualInstance)
						if err != nil {
							return nil, pageCursor.HTTPResponse, err
						}

						err = json.Unmarshal(predictorBytes, &predictor)
						if err != nil {
							return nil, pageCursor.HTTPResponse, err
						}

						if predictor.ID != "" {
							riskPredictorIDs = append(riskPredictorIDs, predictor.ID)
						}
					}
				}
			}

			return riskPredictorIDs, initialHttpResponse, nil
		},
		"ReadAllRiskPredictors",
		legacysdk.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
		&riskPredictorIDs,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	data.Id = data.EnvironmentId
	data.Ids, resp.Diagnostics = framework.StringSliceToTF(riskPredictorIDs)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.State.Set(ctx, &data)
}

func riskPredictorFetchIDsFromCompactNames(ctx context.Context, apiClient *risk.APIClient, managementApiClient *management.APIClient, environmentID string, predictorCompactNames []string) ([]string, diag.Diagnostics) {
	var diags diag.Diagnostics

	riskPredictorIDsMap := make(map[string]string)
	diags.Append(legacysdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			pagedIterator := apiClient.RiskAdvancedPredictorsApi.ReadAllRiskPredictors(ctx, environmentID).Execute()

			var initialHttpResponse *http.Response

			riskPredictorIDsMap := make(map[string]string)

			for pageCursor, err := range pagedIterator {
				if err != nil {
					return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, managementApiClient, environmentID, nil, pageCursor.HTTPResponse, err)
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
		legacysdk.DefaultCustomError,
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
