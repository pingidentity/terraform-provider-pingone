// Copyright © 2025 Ping Identity Corporation
// Code generated by ping-terraform-plugin-framework-generator

package sso

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/filter"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
)

var (
	_ datasource.DataSource              = &populationDataSource{}
	_ datasource.DataSourceWithConfigure = &populationDataSource{}
)

func NewPopulationDataSource() datasource.DataSource {
	return &populationDataSource{}
}

type populationDataSource serviceClientType

func (r *populationDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_population"
}

func (r *populationDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.Client = resourceConfig.Client.API
	if r.Client == nil {
		resp.Diagnostics.AddError(
			"Client not initialised",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.",
		)
		return
	}
}

type populationDataSourceModel struct {
	AlternativeIdentifiers types.Set                    `tfsdk:"alternative_identifiers"`
	Default                types.Bool                   `tfsdk:"default"`
	Description            types.String                 `tfsdk:"description"`
	EnvironmentId          pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	Id                     pingonetypes.ResourceIDValue `tfsdk:"id"`
	Name                   types.String                 `tfsdk:"name"`
	PasswordPolicy         types.Object                 `tfsdk:"password_policy"`
	PasswordPolicyId       pingonetypes.ResourceIDValue `tfsdk:"password_policy_id"`
	PopulationId           pingonetypes.ResourceIDValue `tfsdk:"population_id"`
	PreferredLanguage      types.String                 `tfsdk:"preferred_language"`
	Theme                  types.Object                 `tfsdk:"theme"`
	UserCount              types.Int32                  `tfsdk:"user_count"`
}

func (r *populationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Datasource to retrieve a PingOne population in a PingOne environment, by ID or by name.",
		Attributes: map[string]schema.Attribute{
			"alternative_identifiers": schema.SetAttribute{
				ElementType:         types.StringType,
				Computed:            true,
				Description:         "Alternative identifiers that can be used to search for populations besides \"name\".",
				MarkdownDescription: "Alternative identifiers that can be used to search for populations besides `name`.",
			},
			"default": schema.BoolAttribute{
				Computed:    true,
				Description: "A boolean that indicates whether the population is the default population for the environment.",
			},
			"description": schema.StringAttribute{
				Computed:    true,
				Description: "A string that specifies the description of the population.",
			},
			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment that is configured with the population."),
			),
			"id": framework.Attr_ID(),
			"name": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "A string that specifies the name of the population to retrieve configuration for. Exactly one of \"name\" or \"population_id\" must be defined.",
				MarkdownDescription: "A string that specifies the name of the population to retrieve configuration for. Exactly one of `name` or `population_id` must be defined.",
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("population_id")),
					stringvalidator.LengthAtLeast(1),
				},
			},
			"password_policy": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed:            true,
						CustomType:          pingonetypes.ResourceIDType{},
						Description:         "The ID of the password policy that is used for this population. If absent, the environment's default is used.",
						MarkdownDescription: "The ID of the password policy that is used for this population. If absent, the environment's default is used.",
					},
				},
				Computed:    true,
				Description: "The object reference to the password policy resource applied to the population.",
			},
			"password_policy_id": schema.StringAttribute{
				Description:        framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the ID of the password policy applied to the population.").Description,
				Computed:           true,
				DeprecationMessage: "This attribute is deprecated and will be removed in a future release. Please use the `password_policy.id` attribute instead.",
				CustomType:         pingonetypes.ResourceIDType{},
			},
			"population_id": schema.StringAttribute{
				CustomType:          pingonetypes.ResourceIDType{},
				Optional:            true,
				Computed:            true,
				Description:         "A string that specifies the ID of the population to retrieve configuration for. Must be a valid PingOne resource ID. Exactly one of \"name\" or \"population_id\" must be defined.",
				MarkdownDescription: "A string that specifies the ID of the population to retrieve configuration for. Must be a valid PingOne resource ID. Exactly one of `name` or `population_id` must be defined.",
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("name")),
				},
			},
			"preferred_language": schema.StringAttribute{
				Computed:    true,
				Description: "The language locale for the population.",
			},
			"theme": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed:    true,
						CustomType:  pingonetypes.ResourceIDType{},
						Description: "The ID of the theme to use for the population.",
					},
				},
				Computed:    true,
				Description: "The object reference to the theme resource.",
			},
			"user_count": schema.Int32Attribute{
				Computed:    true,
				Description: "The number of users that belong to the population",
			},
		},
	}
}

func (state *populationDataSourceModel) readClientResponse(response *management.Population) diag.Diagnostics {
	var respDiags, diags diag.Diagnostics
	// alternative_identifiers
	state.AlternativeIdentifiers, diags = types.SetValueFrom(context.Background(), types.StringType, response.AlternativeIdentifiers)
	respDiags.Append(diags...)
	// default
	state.Default = types.BoolPointerValue(response.Default)
	// description
	state.Description = types.StringPointerValue(response.Description)
	// id
	idValue := framework.PingOneResourceIDToTF(response.GetId())
	state.Id = idValue
	// name
	nameValue := types.StringValue(response.Name)
	state.Name = nameValue
	// password_policy_id
	var passwordPolicyIdValue pingonetypes.ResourceIDValue
	if response.PasswordPolicy == nil {
		passwordPolicyIdValue = pingonetypes.NewResourceIDNull()
	} else {
		passwordPolicyIdValue = framework.PingOneResourceIDToTF(response.PasswordPolicy.Id)
	}
	state.PasswordPolicyId = passwordPolicyIdValue
	// password_policy
	passwordPolicyAttrTypes := map[string]attr.Type{
		"id": pingonetypes.ResourceIDType{},
	}
	var passwordPolicyValue types.Object
	if response.PasswordPolicy == nil {
		passwordPolicyValue = types.ObjectNull(passwordPolicyAttrTypes)
	} else {
		passwordPolicyValue, diags = types.ObjectValue(passwordPolicyAttrTypes, map[string]attr.Value{
			"id": framework.PingOneResourceIDToTF(response.PasswordPolicy.Id),
		})
		respDiags.Append(diags...)
	}
	state.PasswordPolicy = passwordPolicyValue
	// population_id
	populationIdValue := framework.PingOneResourceIDToTF(response.GetId())
	state.PopulationId = populationIdValue
	// preferred_language
	state.PreferredLanguage = types.StringPointerValue(response.PreferredLanguage)
	// theme
	themeAttrTypes := map[string]attr.Type{
		"id": pingonetypes.ResourceIDType{},
	}
	var themeValue types.Object
	if response.Theme == nil {
		themeValue = types.ObjectNull(themeAttrTypes)
	} else {
		themeValue, diags = types.ObjectValue(themeAttrTypes, map[string]attr.Value{
			"id": framework.PingOneResourceIDOkToTF(response.Theme.GetIdOk()),
		})
		respDiags.Append(diags...)
	}
	state.Theme = themeValue
	// user_count
	state.UserCount = types.Int32PointerValue(response.UserCount)
	return respDiags
}

func (r *populationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data populationDataSourceModel

	if r.Client == nil || r.Client.ManagementAPIClient == nil {
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

	// Read API call logic
	var responseData *management.Population
	var scimFilter string
	// Build scim filter
	if !data.PopulationId.IsNull() {
		scimFilter = filter.BuildScimFilter(
			append(make([]interface{}, 0), map[string]interface{}{
				"name":   "id",
				"values": []string{data.PopulationId.ValueString()},
			}), map[string]string{})
	} else if !data.Name.IsNull() {
		scimFilter = filter.BuildScimFilter(
			append(make([]interface{}, 0), map[string]interface{}{
				"name":   "name",
				"values": []string{data.Name.ValueString()},
			}), map[string]string{})
	} else {
		resp.Diagnostics.AddError(
			"Missing parameter",
			"Cannot find the requested population. population_id or name must be set.",
		)
		return
	}

	// Run the API call
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			pagedIterator := r.Client.ManagementAPIClient.PopulationsApi.ReadAllPopulations(ctx, data.EnvironmentId.ValueString()).Filter(scimFilter).Execute()

			var initialHttpResponse *http.Response

			for pageCursor, err := range pagedIterator {
				if err != nil {
					return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), nil, pageCursor.HTTPResponse, err)
				}

				if initialHttpResponse == nil {
					initialHttpResponse = pageCursor.HTTPResponse
				}

				if results, ok := pageCursor.EntityArray.Embedded.GetPopulationsOk(); ok {

					for _, resultObj := range results {
						if !data.Name.IsNull() && resultObj.GetName() == data.Name.ValueString() {
							return &resultObj, pageCursor.HTTPResponse, nil
						}
						if !data.PopulationId.IsNull() && resultObj.GetId() == data.PopulationId.ValueString() {
							return &resultObj, pageCursor.HTTPResponse, nil
						}
					}
				}
			}

			return nil, initialHttpResponse, nil
		},
		"ReadAllPopulations",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
		&responseData,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if responseData == nil {
		resp.Diagnostics.AddError(
			"Population not found",
			fmt.Sprintf("The population for environment %s cannot be found", data.EnvironmentId.ValueString()),
		)
		return
	}

	// Read response into the model
	resp.Diagnostics.Append(data.readClientResponse(responseData)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
