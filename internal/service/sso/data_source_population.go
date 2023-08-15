package sso

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	sdkv2resource "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	"github.com/pingidentity/terraform-provider-pingone/internal/filter"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type PopulationDataSource struct {
	client *management.APIClient
	region model.RegionMapping
}

type PopulationDataSourceModel struct {
	Description      types.String `tfsdk:"description"`
	EnvironmentId    types.String `tfsdk:"environment_id"`
	Id               types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	PasswordPolicyId types.String `tfsdk:"password_policy_id"`
	PopulationId     types.String `tfsdk:"population_id"`
}

// Framework interfaces
var (
	_ datasource.DataSource = &PopulationDataSource{}
)

// New Object
func NewPopulationDataSource() datasource.DataSource {
	return &PopulationDataSource{}
}

// Metadata
func (r *PopulationDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_population"
}

// Schema
func (r *PopulationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	nameLength := 1

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Datasource to retrieve a PingOne population.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment that is configured with the population."),
			),

			"population_id": schema.StringAttribute{
				Description: "The ID of the population.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("name")),
					verify.P1ResourceIDValidator(),
				},
			},

			"name": schema.StringAttribute{
				Description: "The name of the population.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("population_id")),
					stringvalidator.LengthAtLeast(nameLength),
				},
			},

			"description": schema.StringAttribute{
				Description: "The description applied to the population.",
				Computed:    true,
			},

			"password_policy_id": schema.StringAttribute{
				Description: "The ID of the password policy applied to the population.",
				Computed:    true,
			},
		},
	}
}

func (r *PopulationDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	resourceConfig, ok := req.ProviderData.(framework.ResourceType)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected the provider client, got: %T. Please report this issue to the provider maintainers.", req.ProviderData),
		)

		return
	}

	preparedClient, err := PrepareClient(ctx, resourceConfig)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			err.Error(),
		)

		return
	}

	r.client = preparedClient
	r.region = resourceConfig.Client.API.Region
}

func (r *PopulationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *PopulationDataSourceModel

	if r.client == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var population management.Population
	var scimFilter string

	if !data.Name.IsNull() {

		scimFilter = filter.BuildScimFilter(
			append(make([]interface{}, 0), map[string]interface{}{
				"name":   "name",
				"values": []string{data.Name.ValueString()},
			}), map[string]string{})

	} else if !data.PopulationId.IsNull() {

		scimFilter = filter.BuildScimFilter(
			append(make([]interface{}, 0), map[string]interface{}{
				"name":   "id",
				"values": []string{data.PopulationId.ValueString()},
			}), map[string]string{})

	} else {
		resp.Diagnostics.AddError(
			"Missing parameter",
			"Cannot find the requested population. population_id or name must be set.",
		)
		return
	}

	// Run the API call
	var entityArray *management.EntityArray
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			return r.client.PopulationsApi.ReadAllPopulations(ctx, data.EnvironmentId.ValueString()).Filter(scimFilter).Execute()
		},
		"ReadAllPopulations",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
		&entityArray,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if populations, ok := entityArray.Embedded.GetPopulationsOk(); ok && len(populations) > 0 && populations[0].Id != nil {

		population = populations[0]

	} else {
		resp.Diagnostics.AddError(
			"Population not found",
			fmt.Sprintf("The population with the specified population_id or name cannot be found in environment %s.", data.EnvironmentId.String()),
		)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(&population)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (p *PopulationDataSourceModel) toState(apiObject *management.Population) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.Id = framework.StringToTF(apiObject.GetId())
	p.PopulationId = framework.StringToTF(apiObject.GetId())
	p.Name = framework.StringOkToTF(apiObject.GetNameOk())
	p.Description = framework.StringOkToTF(apiObject.GetDescriptionOk())

	if v, ok := apiObject.GetPasswordPolicyOk(); ok && v != nil {
		p.PasswordPolicyId = framework.StringOkToTF(v.GetIdOk())
	} else {
		p.PasswordPolicyId = types.StringNull()
	}

	return diags
}

func FetchDefaultPopulation(ctx context.Context, apiClient *management.APIClient, environmentID string) (*management.Population, diag.Diagnostics) {
	defaultTimeout := 30 * time.Second
	return FetchDefaultPopulationWithTimeout(ctx, apiClient, environmentID, defaultTimeout)
}

func FetchDefaultPopulationWithTimeout(ctx context.Context, apiClient *management.APIClient, environmentID string, timeout time.Duration) (*management.Population, diag.Diagnostics) {
	var diags diag.Diagnostics

	stateConf := &sdkv2resource.StateChangeConf{
		Pending: []string{
			"false",
		},
		Target: []string{
			"true",
			"err",
		},
		Refresh: func() (interface{}, string, error) {

			// Run the API call
			var entityArray *management.EntityArray
			diags.Append(framework.ParseResponse(
				ctx,

				func() (any, *http.Response, error) {
					return apiClient.PopulationsApi.ReadAllPopulations(ctx, environmentID).Execute()
				},
				"ReadAllPopulations-FetchDefaultPopulation",
				framework.DefaultCustomError,
				sdk.DefaultCreateReadRetryable,
				&entityArray,
			)...)
			if diags.HasError() {
				return nil, "err", fmt.Errorf("Error reading populations")
			}

			found := false

			var population management.Population

			if populations, ok := entityArray.Embedded.GetPopulationsOk(); ok {

				for _, populationItem := range populations {

					if populationItem.GetDefault() {
						population = populationItem
						found = true
						break
					}
				}
			}

			tflog.Debug(ctx, "Find default population attempt", map[string]interface{}{
				"population": population,
				"result":     strings.ToLower(strconv.FormatBool(found)),
			})

			return population, strings.ToLower(strconv.FormatBool(found)), nil
		},
		Timeout:                   timeout,
		Delay:                     1 * time.Second,
		MinTimeout:                5 * time.Second,
		ContinuousTargetOccurence: 2,
	}
	population, err := stateConf.WaitForState()

	if err != nil {
		diags.AddWarning(
			"Cannot find default population",
			fmt.Sprintf("The default population for environment %s cannot be found: %s", environmentID, err),
		)

		return nil, diags
	}

	returnVar := population.(management.Population)

	return &returnVar, diags

}
