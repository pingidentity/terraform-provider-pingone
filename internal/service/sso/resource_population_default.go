// Copyright © 2025 Ping Identity Corporation

package sso

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type PopulationDefaultResource serviceClientType

type PopulationDefaultResourceModel struct {
	Id                     pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId          pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	Name                   types.String                 `tfsdk:"name"`
	Description            types.String                 `tfsdk:"description"`
	PasswordPolicyId       pingonetypes.ResourceIDValue `tfsdk:"password_policy_id"`
	PasswordPolicy         types.Object                 `tfsdk:"password_policy"`
	AlternativeIdentifiers types.Set                    `tfsdk:"alternative_identifiers"`
	PreferredLanguage      types.String                 `tfsdk:"preferred_language"`
	Theme                  types.Object                 `tfsdk:"theme"`
}

// Framework interfaces
var (
	_ resource.Resource                = &PopulationDefaultResource{}
	_ resource.ResourceWithConfigure   = &PopulationDefaultResource{}
	_ resource.ResourceWithModifyPlan  = &PopulationDefaultResource{}
	_ resource.ResourceWithImportState = &PopulationDefaultResource{}
)

// New Object
func NewPopulationDefaultResource() resource.Resource {
	return &PopulationDefaultResource{}
}

// Metadata
func (r *PopulationDefaultResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_population_default"
}

// Schema.
func (r *PopulationDefaultResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	const attrMinLength = 1

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to overwrite the default PingOne population, or create it if it doesn't already exist.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment to manage the default population in."),
			),

			"name": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The name to apply to the default population.").Description,
				Required:    true,

				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
				},
			},

			"description": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A description to apply to the default population.").Description,
				Optional:    true,
			},

			"password_policy": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Required:    true,
						CustomType:  pingonetypes.ResourceIDType{},
						Description: "The ID of the password policy that is used for this population. If absent, the environment's default is used. Must be a valid PingOne resource ID.",
					},
				},
				Optional:    true,
				Description: "The object reference to the password policy resource. This is an optional property. Conflicts with `password_policy_id`.",
				Validators: []validator.Object{
					objectvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("password_policy_id")),
				},
			},

			"password_policy_id": schema.StringAttribute{
				Description:        framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the ID of a password policy to assign to the population.  Must be a valid PingOne resource ID. The `password_policy.id` attribute should be used instead of this attribute.").ConflictsWith([]string{"password_policy"}).Description,
				DeprecationMessage: "This attribute is deprecated and will be removed in a future release. Please use the `password_policy.id` attribute instead.",
				Optional:           true,
				CustomType:         pingonetypes.ResourceIDType{},
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("password_policy")),
				},
			},

			"alternative_identifiers": schema.SetAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				Description:         "Alternative identifiers that can be used to search for populations besides \"name\".",
				MarkdownDescription: "Alternative identifiers that can be used to search for populations besides `name`.",
			},
			"preferred_language": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "The language locale for the population. If absent, the environment default is used.",
			},
			"theme": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Required:    true,
						CustomType:  pingonetypes.ResourceIDType{},
						Description: "The ID of the theme to use for the population. If absent, the environment's default is used. Must be a valid PingOne resource ID.",
					},
				},
				Optional:    true,
				Computed:    true,
				Description: "The object reference to the theme resource.",
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (p *PopulationDefaultResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// Destruction plan
	if req.Plan.Raw.IsNull() {
		resp.Diagnostics.AddWarning(
			"State change warning",
			"A destroy plan has been detected for the \"pingone_population_default\" resource.  The default population will be reset to it's original configuration, and then removed from Terraform's state.  The population itself (and any user data contained in the population) will not be removed from the PingOne service.",
		)
	}
}

func (r *PopulationDefaultResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.Client = resourceConfig.Client.API
	if r.Client == nil {
		resp.Diagnostics.AddError(
			"Client not initialised",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.",
		)
		return
	}
}

func (r *PopulationDefaultResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state PopulationDefaultResourceModel

	if r.Client == nil || r.Client.ManagementAPIClient == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the model for the API
	population := plan.expand()

	// Run the API call
	readResponse, d := FetchDefaultPopulation(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), false)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	var response *management.Population
	if readResponse == nil {
		resp.Diagnostics.Append(framework.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				fO, fR, fErr := r.Client.ManagementAPIClient.PopulationsApi.CreatePopulation(ctx, plan.EnvironmentId.ValueString()).Population(*population).Execute()
				return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
			},
			"CreatePopulation-Default",
			framework.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
			&response,
		)...)
		if resp.Diagnostics.HasError() {
			return
		}

		if response.Theme == nil || response.Theme.Id == nil {
			responseWithTheme, diags := populationWaitForAssignedThemeId(ctx, r.Client, plan.EnvironmentId.ValueString(), *response.Id)
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}
			if responseWithTheme != nil {
				response = responseWithTheme
			}
		}
	} else {
		resp.Diagnostics.Append(framework.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				fO, fR, fErr := r.Client.ManagementAPIClient.PopulationsApi.UpdatePopulation(ctx, plan.EnvironmentId.ValueString(), readResponse.GetId()).Population(*population).Execute()
				return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
			},
			"UpdatePopulation-Default",
			framework.DefaultCustomError,
			nil,
			&response,
		)...)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the state to save
	state = plan

	// Save updated data into Terraform state
	resp.Diagnostics.Append(state.toState(response)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *PopulationDefaultResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *PopulationDefaultResourceModel

	if r.Client == nil || r.Client.ManagementAPIClient == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	response, d := FetchDefaultPopulation(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), true)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Remove from state if resource is not found
	if response == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(response)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PopulationDefaultResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state PopulationDefaultResourceModel

	if r.Client == nil || r.Client.ManagementAPIClient == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the model for the API
	population := plan.expand()

	// Run the API call
	readResponse, d := FetchDefaultPopulation(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), false)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	var response *management.Population
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.PopulationsApi.UpdatePopulation(ctx, plan.EnvironmentId.ValueString(), readResponse.GetId()).Population(*population).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdatePopulation-Default",
		framework.DefaultCustomError,
		nil,
		&response,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the state to save
	state = plan

	// Save updated data into Terraform state
	resp.Diagnostics.Append(state.toState(response)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *PopulationDefaultResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data PopulationDefaultResourceModel

	if r.Client == nil || r.Client.ManagementAPIClient == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the model for the API
	population := management.NewPopulation("Default")
	population.SetDescription("Automatically created population.")
	population.SetDefault(true)

	// Run the API call
	readResponse, d := FetchDefaultPopulation(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), true)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	var response *management.Population
	if readResponse != nil {
		resp.Diagnostics.Append(framework.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				fO, fR, fErr := r.Client.ManagementAPIClient.PopulationsApi.UpdatePopulation(ctx, data.EnvironmentId.ValueString(), readResponse.GetId()).Population(*population).Execute()
				return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
			},
			"UpdatePopulation-DeleteDefault",
			framework.CustomErrorResourceNotFoundWarning,
			nil,
			&response,
		)...)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	if response != nil {
		resp.Diagnostics.AddWarning(
			"State change warning",
			"The \"pingone_population_default\" resource has been destroyed.  The default population has been reset to it's original configuration, and removed from Terraform's state.  The population itself (and any user data contained in the population) has not been removed from the PingOne service.",
		)
	}
}

func (r *PopulationDefaultResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	idComponents := []framework.ImportComponent{
		{
			Label:  "environment_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
	}

	attributes, err := framework.ParseImportID(req.ID, idComponents...)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			err.Error(),
		)
		return
	}

	for _, idComponent := range idComponents {
		pathKey := idComponent.Label

		if idComponent.PrimaryID {
			pathKey = "id"
		}

		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root(pathKey), attributes[idComponent.Label])...)
	}
}

func (p *PopulationDefaultResourceModel) expand() *management.Population {
	data := management.NewPopulation(p.Name.ValueString())
	data.SetDefault(true)

	if !p.Description.IsNull() && !p.Description.IsUnknown() {
		data.SetDescription(p.Description.ValueString())
	}

	// password_policy
	if !p.PasswordPolicy.IsNull() && !p.PasswordPolicy.IsUnknown() {
		passwordPolicyValue := &management.PopulationPasswordPolicy{}
		passwordPolicyAttrs := p.PasswordPolicy.Attributes()
		passwordPolicyValue.Id = passwordPolicyAttrs["id"].(pingonetypes.ResourceIDValue).ValueString()
		data.PasswordPolicy = passwordPolicyValue
	} else if !p.PasswordPolicyId.IsNull() && !p.PasswordPolicyId.IsUnknown() {
		// password_policy_id
		data.PasswordPolicy = &management.PopulationPasswordPolicy{
			Id: p.PasswordPolicyId.ValueString(),
		}
	}

	// alternative_identifiers
	if !p.AlternativeIdentifiers.IsNull() && !p.AlternativeIdentifiers.IsUnknown() {
		altIds := []string{}
		for _, elem := range p.AlternativeIdentifiers.Elements() {
			altIds = append(altIds, elem.(types.String).ValueString())
		}
		data.SetAlternativeIdentifiers(altIds)
	}

	// preferred_language
	if !p.PreferredLanguage.IsNull() && !p.PreferredLanguage.IsUnknown() {
		data.PreferredLanguage = p.PreferredLanguage.ValueStringPointer()
	}

	// theme
	if !p.Theme.IsNull() && !p.Theme.IsUnknown() {
		themeValue := &management.PopulationTheme{}
		themeAttrs := p.Theme.Attributes()
		themeValue.Id = themeAttrs["id"].(pingonetypes.ResourceIDValue).ValueStringPointer()
		data.Theme = themeValue
	}

	return data
}

func (p *PopulationDefaultResourceModel) toState(apiObject *management.Population) diag.Diagnostics {
	var diags, buildDiags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)
		return diags
	}

	p.Id = framework.PingOneResourceIDOkToTF(apiObject.GetIdOk())
	p.EnvironmentId = framework.PingOneResourceIDOkToTF(apiObject.Environment.GetIdOk())
	p.Name = framework.StringOkToTF(apiObject.GetNameOk())
	p.Description = framework.StringOkToTF(apiObject.GetDescriptionOk())

	// password_policy_id
	if !p.PasswordPolicyId.IsNull() {
		var passwordPolicyIdValue pingonetypes.ResourceIDValue
		if apiObject.PasswordPolicy == nil {
			passwordPolicyIdValue = pingonetypes.NewResourceIDNull()
		} else {
			passwordPolicyIdValue = framework.PingOneResourceIDToTF(apiObject.PasswordPolicy.Id)
		}
		p.PasswordPolicyId = passwordPolicyIdValue
	} else {
		// password_policy
		passwordPolicyAttrTypes := map[string]attr.Type{
			"id": pingonetypes.ResourceIDType{},
		}
		var passwordPolicyValue types.Object
		if apiObject.PasswordPolicy == nil {
			passwordPolicyValue = types.ObjectNull(passwordPolicyAttrTypes)
		} else {
			passwordPolicyValue, buildDiags = types.ObjectValue(passwordPolicyAttrTypes, map[string]attr.Value{
				"id": framework.PingOneResourceIDToTF(apiObject.PasswordPolicy.Id),
			})
			diags.Append(buildDiags...)
		}
		p.PasswordPolicy = passwordPolicyValue
	}

	// alternative_identifiers
	p.AlternativeIdentifiers, buildDiags = types.SetValueFrom(context.Background(), types.StringType, apiObject.AlternativeIdentifiers)
	diags.Append(buildDiags...)

	// preferred_language
	p.PreferredLanguage = framework.StringOkToTF(apiObject.GetPreferredLanguageOk())

	// theme
	themeAttrTypes := map[string]attr.Type{
		"id": pingonetypes.ResourceIDType{},
	}
	var themeValue types.Object
	if apiObject.Theme == nil {
		themeValue = types.ObjectNull(themeAttrTypes)
	} else {
		var themeDiags diag.Diagnostics
		themeValue, themeDiags = types.ObjectValue(themeAttrTypes, map[string]attr.Value{
			"id": framework.PingOneResourceIDOkToTF(apiObject.Theme.GetIdOk()),
		})
		diags.Append(themeDiags...)
	}
	p.Theme = themeValue

	return diags
}

func FetchDefaultPopulation(ctx context.Context, apiClient *management.APIClient, environmentID string, warnOnNotFound bool) (*management.Population, diag.Diagnostics) {
	defaultTimeout := 30 * time.Second
	return FetchDefaultPopulationWithTimeout(ctx, apiClient, environmentID, warnOnNotFound, defaultTimeout)
}

func FetchDefaultPopulationWithTimeout(ctx context.Context, apiClient *management.APIClient, environmentID string, warnOnNotFound bool, timeout time.Duration) (*management.Population, diag.Diagnostics) {
	var diags diag.Diagnostics

	errorFunction := framework.DefaultCustomError
	if warnOnNotFound {
		errorFunction = framework.CustomErrorResourceNotFoundWarning
	}

	stateConf := &retry.StateChangeConf{
		Pending: []string{
			"false",
		},
		Target: []string{
			"true",
			"err",
		},
		Refresh: func() (interface{}, string, error) {

			// Run the API call
			var defaultPopulation *management.Population
			diags.Append(framework.ParseResponse(
				ctx,

				func() (any, *http.Response, error) {
					pagedIterator := apiClient.PopulationsApi.ReadAllPopulations(ctx, environmentID).Execute()

					var initialHttpResponse *http.Response

					for pageCursor, err := range pagedIterator {
						if err != nil {
							return framework.CheckEnvironmentExistsOnPermissionsError(ctx, apiClient, environmentID, nil, pageCursor.HTTPResponse, err)
						}

						if initialHttpResponse == nil {
							initialHttpResponse = pageCursor.HTTPResponse
						}

						if populations, ok := pageCursor.EntityArray.Embedded.GetPopulationsOk(); ok {

							for _, populationItem := range populations {

								if populationItem.GetDefault() {
									return &populationItem, pageCursor.HTTPResponse, nil
								}
							}
						}
					}

					return nil, initialHttpResponse, nil
				},
				"ReadAllPopulations-FetchDefaultPopulation",
				errorFunction,
				sdk.DefaultCreateReadRetryable,
				&defaultPopulation,
			)...)
			if diags.HasError() {
				return nil, "err", fmt.Errorf("Error reading populations")
			}

			tflog.Debug(ctx, "Find default population attempt", map[string]interface{}{
				"population": defaultPopulation,
				"result":     strings.ToLower(strconv.FormatBool(defaultPopulation != nil)),
			})

			return defaultPopulation, strings.ToLower(strconv.FormatBool(defaultPopulation != nil)), nil
		},
		Timeout:                   timeout,
		Delay:                     1 * time.Second,
		MinTimeout:                1 * time.Second,
		ContinuousTargetOccurence: 2,
	}
	population, err := stateConf.WaitForStateContext(ctx)

	if err != nil {
		tflog.Warn(ctx, "Cannot find default population for the environment", map[string]interface{}{
			"environment": environmentID,
			"err":         err,
		})

		return nil, diags
	}

	returnVar := population.(*management.Population)

	return returnVar, diags
}
