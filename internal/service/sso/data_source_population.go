package sso

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type PopulationDataSource struct {
	client *management.APIClient
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

			"environment_id": framework.Attr_EnvironmentID(framework.SchemaDescription{
				Description: "The ID of the environment that is configured with the population."},
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

	preparedClient, err := prepareClient(ctx, resourceConfig)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			err.Error(),
		)

		return
	}

	r.client = preparedClient
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

	if !data.Name.IsNull() {

		// Run the API call
		response, diags := framework.ParseResponse(
			ctx,

			func() (interface{}, *http.Response, error) {
				return r.client.PopulationsApi.ReadAllPopulations(ctx, data.EnvironmentId.ValueString()).Execute()
			},
			"ReadAllPopulations",
			framework.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
		)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		entityArray := response.(*management.EntityArray)

		if populations, ok := entityArray.Embedded.GetPopulationsOk(); ok {

			found := false
			for _, populationItem := range populations {

				if populationItem.GetName() == data.Name.ValueString() {
					population = populationItem
					found = true
					break
				}
			}

			if !found {
				resp.Diagnostics.AddError(
					"Cannot find population from name",
					fmt.Sprintf("The population %s for environment %s cannot be found", data.Name.String(), data.EnvironmentId.String()),
				)
				return
			}

		}

	} else if !data.PopulationId.IsNull() {

		// Run the API call
		response, diags := framework.ParseResponse(
			ctx,

			func() (interface{}, *http.Response, error) {
				return r.client.PopulationsApi.ReadOnePopulation(ctx, data.EnvironmentId.ValueString(), data.PopulationId.ValueString()).Execute()
			},
			"ReadOnePopulation",
			framework.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
		)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		population = *response.(*management.Population)
	} else {
		resp.Diagnostics.AddError(
			"Missing parameter",
			"Cannot find the requested population. population_id or name must be set.",
		)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(&population)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (p *PopulationDataSourceModel) toState(v *management.Population) diag.Diagnostics {
	var diags diag.Diagnostics

	if v == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.Id = framework.StringToTF(v.GetId())
	p.PopulationId = framework.StringToTF(v.GetId())
	p.Name = framework.StringToTF(v.GetName())
	p.Description = framework.StringToTF(v.GetDescription())
	p.PasswordPolicyId = framework.StringToTF(v.GetPasswordPolicy().Id)

	return diags
}
