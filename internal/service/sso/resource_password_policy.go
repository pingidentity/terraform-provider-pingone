package sso

import (
	"context"
	"fmt"
	"net/http"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type PasswordPolicyResource serviceClientType

type PasswordPolicyResourceModel struct {
	Id                            pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId                 pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	Name                          types.String                 `tfsdk:"name"`
	Description                   types.String                 `tfsdk:"description"`
	Default                       types.Bool                   `tfsdk:"default"`
	ExcludesCommonlyUsedPasswords types.Bool                   `tfsdk:"excludes_commonly_used_passwords"`
	ExcludesProfileData           types.Bool                   `tfsdk:"excludes_profile_data"`
	History                       types.Object                 `tfsdk:"history"`
	Length                        types.Object                 `tfsdk:"length"`
	Lockout                       types.Object                 `tfsdk:"lockout"`
	MinCharacters                 types.Object                 `tfsdk:"min_characters"`
	PasswordAgeMax                types.Int64                  `tfsdk:"password_age_max"`
	PasswordAgeMin                types.Int64                  `tfsdk:"password_age_min"`
	MaxRepeatedCharacters         types.Int64                  `tfsdk:"max_repeated_characters"`
	MinComplexity                 types.Int64                  `tfsdk:"min_complexity"`
	MinUniqueCharacters           types.Int64                  `tfsdk:"min_unique_characters"`
	NotSimilarToCurrent           types.Bool                   `tfsdk:"not_similar_to_current"`
	PopulationCount               types.Int64                  `tfsdk:"population_count"`
}

type PasswordPolicyPasswordHistoryResourceModel struct {
	Count         types.Int64 `tfsdk:"count"`
	RetentionDays types.Int64 `tfsdk:"retention_days"`
}

type PasswordPolicyPasswordLengthResourceModel struct {
	Max types.Int64 `tfsdk:"max"`
	Min types.Int64 `tfsdk:"min"`
}

type PasswordPolicyAccountLockoutResourceModel struct {
	DurationSeconds types.Int64 `tfsdk:"duration_seconds"`
	FailureCount    types.Int64 `tfsdk:"failure_count"`
}

type PasswordPolicyMinCharactersResourceModel struct {
	AlphabeticalUppercase types.Int64 `tfsdk:"alphabetical_uppercase"`
	AlphabeticalLowercase types.Int64 `tfsdk:"alphabetical_lowercase"`
	Numeric               types.Int64 `tfsdk:"numeric"`
	SpecialCharacters     types.Int64 `tfsdk:"special_characters"`
}

var (
	passwordPolicyHistoryTFObjectTypes = map[string]attr.Type{
		"count":          types.Int64Type,
		"retention_days": types.Int64Type,
	}

	passwordPolicyLengthTFObjectTypes = map[string]attr.Type{
		"max": types.Int64Type,
		"min": types.Int64Type,
	}

	passwordPolicyLockoutTFObjectTypes = map[string]attr.Type{
		"duration_seconds": types.Int64Type,
		"failure_count":    types.Int64Type,
	}

	passwordPolicyMinCharactersTFObjectTypes = map[string]attr.Type{
		"alphabetical_uppercase": types.Int64Type,
		"alphabetical_lowercase": types.Int64Type,
		"numeric":                types.Int64Type,
		"special_characters":     types.Int64Type,
	}
)

// Framework interfaces
var (
	_ resource.Resource                = &PasswordPolicyResource{}
	_ resource.ResourceWithConfigure   = &PasswordPolicyResource{}
	_ resource.ResourceWithImportState = &PasswordPolicyResource{}
)

// New Object
func NewPasswordPolicyResource() resource.Resource {
	return &PasswordPolicyResource{}
}

// Metadata
func (r *PasswordPolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_password_policy"
}

// Schema.
func (r *PasswordPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	const attrMinLength = 1
	const passwordLengthMax = 255
	const passwordLengthMinMin = 8
	const passwordLengthMinMax = 32
	const minCharactersFixedValue = 1
	const maxRepeatedCharactersFixedValue = 2
	const minComplexityFixedValue = 7
	const minUniqueCharactersFixedValue = 5

	defaultDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether this password policy is enforced as the default within the environment. When set to `true`, all other password policies are set to `false`.",
	).DefaultValue(false)

	excludeCommonlyUsedPasswordsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether to ensure the password is not one of the commonly used passwords.",
	).DefaultValue(false)

	excludeProfileDataDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether to ensure the password is not an exact match for the value of any attribute in the user's profile, such as name, phone number, or address.",
	).DefaultValue(false)

	passwordLengthMaxDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An integer that specifies the maximum number of characters allowed for the password. This property is not enforced when not present.",
	).DefaultValue(passwordLengthMax).FixedValue(passwordLengthMax)

	passwordLengthMinDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("An integer that specifies the minimum number of characters required for the password. This can be from `%d` to `%d` (inclusive). This property is not enforced when not present.", passwordLengthMinMin, passwordLengthMinMax),
	)

	minCharactersDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that specifies sets of characters that can be included, and the value is the minimum number of times one of the characters must appear in the user's password. The only allowed key values are `ABCDEFGHIJKLMNOPQRSTUVWXYZ`, `abcdefghijklmnopqrstuvwxyz`, `0123456789`, and `~!@#$%^&*()-_=+[]{}\\|;:,.<>/?`. This property is not enforced when not present.",
	)

	minCharactersAlphabeticalUppercaseDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An integer that specifies the count of alphabetical uppercase characters (`ABCDEFGHIJKLMNOPQRSTUVWXYZ`) that should feature in the user's password.",
	).DefaultValue(minCharactersFixedValue).FixedValue(minCharactersFixedValue)

	minCharactersAlphabeticalLowercaseDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An integer that specifies the count of alphabetical uppercase characters (`abcdefghijklmnopqrstuvwxyz`) that should feature in the user's password.",
	).DefaultValue(minCharactersFixedValue).FixedValue(minCharactersFixedValue)

	minCharactersNumericDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An integer that specifies the count of numeric characters (`0123456789`) that should feature in the user's password.",
	).DefaultValue(minCharactersFixedValue).FixedValue(minCharactersFixedValue)

	minCharactersSpecialCharactersDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An integer that specifies the count of special characters (`~!@#$%^&*()-_=+[]{}\\|;:,.<>/?`) that should feature in the user's password.",
	).DefaultValue(minCharactersFixedValue).FixedValue(minCharactersFixedValue)

	passwordAgeMaxDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An integer that specifies the maximum number of days the same password can be used before it must be changed. The value must be a positive, non-zero integer.  The value must be greater than the sum of `min` (if set) + 21 (the expiration warning interval for passwords).",
	)

	maxRepeatedCharactersDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An integer that specifies the maximum number of repeated characters allowed. This property is not enforced when not present.",
	).FixedValue(maxRepeatedCharactersFixedValue)

	minComplexityDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An integer that specifies the minimum complexity of the password based on the concept of password haystacks. The value is the number of days required to exhaust the entire search space during a brute force attack. This property is not enforced when not present.",
	).FixedValue(minComplexityFixedValue)

	minUniqueCharactersDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An integer that specifies the minimum number of unique characters required. This property is not enforced when not present.",
	).FixedValue(minUniqueCharactersFixedValue)

	notSimilarToCurrentDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that, when set to `true`, ensures that the proposed password is not too similar to the user's current password based on the Levenshtein distance algorithm. The value of this parameter is evaluated only for password change actions in which the user enters both the current and the new password. By design, PingOne does not know the user's current password.",
	).DefaultValue(false)

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage PingOne password policies in an environment.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment to manage the password policy in."),
			),

			"name": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the name of the password policy.").Description,
				Required:    true,

				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
				},
			},

			"description": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the description to apply to the password policy.").Description,
				Optional:    true,
			},

			"default": schema.BoolAttribute{
				Description:         defaultDescription.Description,
				MarkdownDescription: defaultDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,

				Default: booldefault.StaticBool(false),
			},

			"excludes_commonly_used_passwords": schema.BoolAttribute{
				Description:         excludeCommonlyUsedPasswordsDescription.Description,
				MarkdownDescription: excludeCommonlyUsedPasswordsDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,

				Default: booldefault.StaticBool(false),
			},

			"excludes_profile_data": schema.BoolAttribute{
				Description:         excludeProfileDataDescription.Description,
				MarkdownDescription: excludeProfileDataDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,

				Default: booldefault.StaticBool(false),
			},

			"history": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that specifies settings to control the user's password history.").Description,
				Optional:    true,

				Attributes: map[string]schema.Attribute{
					"count": schema.Int64Attribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("An integer that specifies the number of prior passwords to keep for prevention of password re-use. The value must be a positive, non-zero integer.").Description,
						Required:    true,

						Validators: []validator.Int64{
							int64validator.AtLeast(attrMinLength),
						},
					},

					"retention_days": schema.Int64Attribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("An integer that specifies the length of time to keep recent passwords for prevention of password re-use. The value must be a positive, non-zero integer.").Description,
						Required:    true,

						Validators: []validator.Int64{
							int64validator.AtLeast(attrMinLength),
						},
					},
				},
			},

			"length": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that specifies settings to control the user's password length.").Description,
				Optional:    true,

				Attributes: map[string]schema.Attribute{
					"max": schema.Int64Attribute{
						Description:         passwordLengthMaxDescription.Description,
						MarkdownDescription: passwordLengthMaxDescription.MarkdownDescription,
						Optional:            true,
						Computed:            true,

						Default: int64default.StaticInt64(passwordLengthMax),

						Validators: []validator.Int64{
							int64validator.Between(passwordLengthMax, passwordLengthMax),
						},
					},

					"min": schema.Int64Attribute{
						Description:         passwordLengthMinDescription.Description,
						MarkdownDescription: passwordLengthMinDescription.MarkdownDescription,
						Required:            true,

						Validators: []validator.Int64{
							int64validator.Between(passwordLengthMinMin, passwordLengthMinMax),
						},
					},
				},
			},

			"lockout": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that specifies settings to control the user's lockout on unsuccessful authentication attempts.").Description,
				Optional:    true,

				Attributes: map[string]schema.Attribute{
					"duration_seconds": schema.Int64Attribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("An integer that specifies the length of time before a password is automatically moved out of the lock out state. The value must be a positive, non-zero integer.").Description,
						Required:    true,

						Validators: []validator.Int64{
							int64validator.AtLeast(attrMinLength),
						},
					},

					"failure_count": schema.Int64Attribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("An integer that specifies the number of tries before a password is placed in the lockout state. The value must be a positive, non-zero integer.").Description,
						Required:    true,

						Validators: []validator.Int64{
							int64validator.AtLeast(attrMinLength),
						},
					},
				},
			},

			"min_characters": schema.SingleNestedAttribute{
				Description:         minCharactersDescription.Description,
				MarkdownDescription: minCharactersDescription.MarkdownDescription,
				Optional:            true,

				Attributes: map[string]schema.Attribute{
					"alphabetical_uppercase": schema.Int64Attribute{
						Description:         minCharactersAlphabeticalUppercaseDescription.Description,
						MarkdownDescription: minCharactersAlphabeticalUppercaseDescription.MarkdownDescription,
						Optional:            true,
						Computed:            true,

						Default: int64default.StaticInt64(minCharactersFixedValue),

						Validators: []validator.Int64{
							int64validator.Between(minCharactersFixedValue, minCharactersFixedValue),
						},
					},

					"alphabetical_lowercase": schema.Int64Attribute{
						Description:         minCharactersAlphabeticalLowercaseDescription.Description,
						MarkdownDescription: minCharactersAlphabeticalLowercaseDescription.MarkdownDescription,
						Optional:            true,
						Computed:            true,

						Default: int64default.StaticInt64(minCharactersFixedValue),

						Validators: []validator.Int64{
							int64validator.Between(minCharactersFixedValue, minCharactersFixedValue),
						},
					},

					"numeric": schema.Int64Attribute{
						Description:         minCharactersNumericDescription.Description,
						MarkdownDescription: minCharactersNumericDescription.MarkdownDescription,
						Optional:            true,
						Computed:            true,

						Default: int64default.StaticInt64(minCharactersFixedValue),

						Validators: []validator.Int64{
							int64validator.Between(minCharactersFixedValue, minCharactersFixedValue),
						},
					},

					"special_characters": schema.Int64Attribute{
						Description:         minCharactersSpecialCharactersDescription.Description,
						MarkdownDescription: minCharactersSpecialCharactersDescription.MarkdownDescription,
						Optional:            true,
						Computed:            true,

						Default: int64default.StaticInt64(minCharactersFixedValue),

						Validators: []validator.Int64{
							int64validator.Between(minCharactersFixedValue, minCharactersFixedValue),
						},
					},
				},
			},

			"password_age_max": schema.Int64Attribute{
				Description:         passwordAgeMaxDescription.Description,
				MarkdownDescription: passwordAgeMaxDescription.MarkdownDescription,
				Optional:            true,

				Validators: []validator.Int64{
					int64validator.AtLeast(attrMinLength),
				},
			},

			"password_age_min": schema.Int64Attribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("An integer that specifies the minimum number of days a password must be used before changing. The value must be a positive, non-zero integer. This property is not enforced when not present.").Description,
				Optional:    true,

				Validators: []validator.Int64{
					int64validator.AtLeast(attrMinLength),
				},
			},

			"max_repeated_characters": schema.Int64Attribute{
				Description:         maxRepeatedCharactersDescription.Description,
				MarkdownDescription: maxRepeatedCharactersDescription.MarkdownDescription,
				Optional:            true,

				Validators: []validator.Int64{
					int64validator.Between(maxRepeatedCharactersFixedValue, maxRepeatedCharactersFixedValue),
				},
			},

			"min_complexity": schema.Int64Attribute{
				Description:         minComplexityDescription.Description,
				MarkdownDescription: minComplexityDescription.MarkdownDescription,
				Optional:            true,

				Validators: []validator.Int64{
					int64validator.Between(minComplexityFixedValue, minComplexityFixedValue),
				},
			},

			"min_unique_characters": schema.Int64Attribute{
				Description:         minUniqueCharactersDescription.Description,
				MarkdownDescription: minUniqueCharactersDescription.MarkdownDescription,
				Optional:            true,

				Validators: []validator.Int64{
					int64validator.Between(minUniqueCharactersFixedValue, minUniqueCharactersFixedValue),
				},
			},

			"not_similar_to_current": schema.BoolAttribute{
				Description:         notSimilarToCurrentDescription.Description,
				MarkdownDescription: notSimilarToCurrentDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,

				Default: booldefault.StaticBool(false),
			},

			"population_count": schema.Int64Attribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("An integer that specifies the number of populations associated with the password policy.").Description,
				Computed:    true,
			},
		},
	}
}

func (r *PasswordPolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *PasswordPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state PasswordPolicyResourceModel

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
	passwordPolicy, d := plan.expand(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *management.PasswordPolicy
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.PasswordPoliciesApi.CreatePasswordPolicy(ctx, plan.EnvironmentId.ValueString()).PasswordPolicy(*passwordPolicy).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"CreatePasswordPolicy",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
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

func (r *PasswordPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *PasswordPolicyResourceModel

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
	var response *management.PasswordPolicy
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.PasswordPoliciesApi.ReadOnePasswordPolicy(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"ReadOnePasswordPolicy",
		framework.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
		&response,
	)...)
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

func (r *PasswordPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state PasswordPolicyResourceModel

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
	passwordPolicy, d := plan.expand(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *management.PasswordPolicy
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.PasswordPoliciesApi.UpdatePasswordPolicy(ctx, plan.EnvironmentId.ValueString(), plan.Id.ValueString()).PasswordPolicy(*passwordPolicy).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdatePasswordPolicy",
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

func (r *PasswordPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *PasswordPolicyResourceModel

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
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fR, fErr := r.Client.ManagementAPIClient.PasswordPoliciesApi.DeletePasswordPolicy(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), nil, fR, fErr)
		},
		"DeletePasswordPolicy",
		passwordPolicyDeleteCustomError,
		nil,
		nil,
	)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

var passwordPolicyDeleteCustomError = func(p1Error model.P1Error) diag.Diagnostics {
	var diags diag.Diagnostics

	// Undeletable default risk policy
	if v, ok := p1Error.GetDetailsOk(); ok && v != nil && len(v) > 0 {
		if v[0].GetCode() == "CONSTRAINT_VIOLATION" {
			if match, _ := regexp.MatchString("Default Password Policy can not be deleted", v[0].GetMessage()); match {

				diags.AddWarning("Cannot delete the default password policy", "Due to API restrictions, the provider cannot delete the default password policy for an environment.  The policy has been removed from Terraform state but has been left in place in the PingOne service.")

				return diags
			}
		}
	}

	return framework.CustomErrorResourceNotFoundWarning(p1Error)
}

func (r *PasswordPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	idComponents := []framework.ImportComponent{
		{
			Label:  "environment_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
		{
			Label:     "password_policy_id",
			Regexp:    verify.P1ResourceIDRegexp,
			PrimaryID: true,
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

func (p *PasswordPolicyResourceModel) expand(ctx context.Context) (*management.PasswordPolicy, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := management.NewPasswordPolicy(
		p.ExcludesCommonlyUsedPasswords.ValueBool(),
		p.ExcludesProfileData.ValueBool(),
		p.Name.ValueString(),
		p.NotSimilarToCurrent.ValueBool(),
	)

	if !p.Description.IsNull() && !p.Description.IsUnknown() {
		data.SetDescription(p.Description.ValueString())
	}

	if !p.Default.IsNull() && !p.Default.IsUnknown() {
		data.SetDefault(p.Default.ValueBool())
	} else {
		data.SetDefault(false)
	}

	if !p.History.IsNull() && !p.History.IsUnknown() {
		var plan PasswordPolicyPasswordHistoryResourceModel
		diags.Append(p.History.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		history := management.NewPasswordPolicyHistory()

		if !plan.Count.IsNull() && !plan.Count.IsUnknown() {
			history.SetCount(int32(plan.Count.ValueInt64()))
		}

		if !plan.RetentionDays.IsNull() && !plan.RetentionDays.IsUnknown() {
			history.SetRetentionDays(int32(plan.RetentionDays.ValueInt64()))
		}

		data.SetHistory(*history)
	}

	if !p.Length.IsNull() && !p.Length.IsUnknown() {
		var plan PasswordPolicyPasswordLengthResourceModel
		diags.Append(p.Length.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		length := management.NewPasswordPolicyLength()

		if !plan.Max.IsNull() && !plan.Max.IsUnknown() {
			length.SetMax(int32(plan.Max.ValueInt64()))
		}

		if !plan.Min.IsNull() && !plan.Min.IsUnknown() {
			length.SetMin(int32(plan.Min.ValueInt64()))
		}

		data.SetLength(*length)
	}

	if !p.Lockout.IsNull() && !p.Lockout.IsUnknown() {
		var plan PasswordPolicyAccountLockoutResourceModel
		diags.Append(p.Lockout.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		lockout := management.NewPasswordPolicyLockout()

		if !plan.DurationSeconds.IsNull() && !plan.DurationSeconds.IsUnknown() {
			lockout.SetDurationSeconds(int32(plan.DurationSeconds.ValueInt64()))
		}

		if !plan.FailureCount.IsNull() && !plan.FailureCount.IsUnknown() {
			lockout.SetFailureCount(int32(plan.FailureCount.ValueInt64()))
		}

		data.SetLockout(*lockout)
	}

	if !p.MinCharacters.IsNull() && !p.MinCharacters.IsUnknown() {
		var plan PasswordPolicyMinCharactersResourceModel
		diags.Append(p.MinCharacters.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		minCharacters := management.NewPasswordPolicyMinCharacters()

		if !plan.AlphabeticalUppercase.IsNull() && !plan.AlphabeticalUppercase.IsUnknown() {
			minCharacters.SetABCDEFGHIJKLMNOPQRSTUVWXYZ(int32(plan.AlphabeticalUppercase.ValueInt64()))
		}

		if !plan.AlphabeticalLowercase.IsNull() && !plan.AlphabeticalLowercase.IsUnknown() {
			minCharacters.SetAbcdefghijklmnopqrstuvwxyz(int32(plan.AlphabeticalLowercase.ValueInt64()))
		}

		if !plan.Numeric.IsNull() && !plan.Numeric.IsUnknown() {
			minCharacters.SetVar0123456789(int32(plan.Numeric.ValueInt64()))
		}

		if !plan.SpecialCharacters.IsNull() && !plan.SpecialCharacters.IsUnknown() {
			minCharacters.SetSpecialChar(int32(plan.SpecialCharacters.ValueInt64()))
		}

		data.SetMinCharacters(*minCharacters)
	}

	if !p.PasswordAgeMax.IsNull() && !p.PasswordAgeMax.IsUnknown() {
		data.SetMaxAgeDays(int32(p.PasswordAgeMax.ValueInt64()))
	}

	if !p.PasswordAgeMin.IsNull() && !p.PasswordAgeMin.IsUnknown() {
		data.SetMinAgeDays(int32(p.PasswordAgeMin.ValueInt64()))
	}

	if !p.MaxRepeatedCharacters.IsNull() && !p.MaxRepeatedCharacters.IsUnknown() {
		data.SetMaxRepeatedCharacters(int32(p.MaxRepeatedCharacters.ValueInt64()))
	}

	if !p.MinComplexity.IsNull() && !p.MinComplexity.IsUnknown() {
		data.SetMinComplexity(int32(p.MinComplexity.ValueInt64()))
	}

	if !p.MinUniqueCharacters.IsNull() && !p.MinUniqueCharacters.IsUnknown() {
		data.SetMinUniqueCharacters(int32(p.MinUniqueCharacters.ValueInt64()))
	}

	if !p.MinUniqueCharacters.IsNull() && !p.MinUniqueCharacters.IsUnknown() {
		data.SetMinUniqueCharacters(int32(p.MinUniqueCharacters.ValueInt64()))
	}

	return data, diags
}

func (p *PasswordPolicyResourceModel) toState(apiObject *management.PasswordPolicy) diag.Diagnostics {
	var diags, d diag.Diagnostics

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
	p.Default = framework.BoolOkToTF(apiObject.GetDefaultOk())
	p.ExcludesCommonlyUsedPasswords = framework.BoolOkToTF(apiObject.GetExcludesCommonlyUsedOk())
	p.ExcludesProfileData = framework.BoolOkToTF(apiObject.GetExcludesProfileDataOk())

	p.History, d = passwordPolicyHistoryOkToTF(apiObject.GetHistoryOk())
	diags.Append(d...)

	p.Length, d = passwordPolicyLengthOkToTF(apiObject.GetLengthOk())
	diags.Append(d...)

	p.Lockout, d = passwordPolicyLockoutOkToTF(apiObject.GetLockoutOk())
	diags.Append(d...)

	p.MinCharacters, d = passwordPolicyMinCharactersOkToTF(apiObject.GetMinCharactersOk())
	diags.Append(d...)

	p.PasswordAgeMax = framework.Int32OkToTF(apiObject.GetMaxAgeDaysOk())
	p.PasswordAgeMin = framework.Int32OkToTF(apiObject.GetMinAgeDaysOk())
	p.MaxRepeatedCharacters = framework.Int32OkToTF(apiObject.GetMaxRepeatedCharactersOk())
	p.MinComplexity = framework.Int32OkToTF(apiObject.GetMinComplexityOk())
	p.MinUniqueCharacters = framework.Int32OkToTF(apiObject.GetMinUniqueCharactersOk())
	p.NotSimilarToCurrent = framework.BoolOkToTF(apiObject.GetNotSimilarToCurrentOk())
	p.PopulationCount = framework.Int32OkToTF(apiObject.GetPopulationCountOk())

	return diags
}

func passwordPolicyHistoryOkToTF(apiObject *management.PasswordPolicyHistory, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(passwordPolicyHistoryTFObjectTypes), diags
	}

	o := map[string]attr.Value{
		"count":          framework.Int32OkToTF(apiObject.GetCountOk()),
		"retention_days": framework.Int32OkToTF(apiObject.GetRetentionDaysOk()),
	}

	returnVar, d := types.ObjectValue(passwordPolicyHistoryTFObjectTypes, o)
	diags.Append(d...)

	return returnVar, diags
}

func passwordPolicyLengthOkToTF(apiObject *management.PasswordPolicyLength, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(passwordPolicyLengthTFObjectTypes), diags
	}

	o := map[string]attr.Value{
		"max": framework.Int32OkToTF(apiObject.GetMaxOk()),
		"min": framework.Int32OkToTF(apiObject.GetMinOk()),
	}

	returnVar, d := types.ObjectValue(passwordPolicyLengthTFObjectTypes, o)
	diags.Append(d...)

	return returnVar, diags
}

func passwordPolicyLockoutOkToTF(apiObject *management.PasswordPolicyLockout, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(passwordPolicyLockoutTFObjectTypes), diags
	}

	o := map[string]attr.Value{
		"duration_seconds": framework.Int32OkToTF(apiObject.GetDurationSecondsOk()),
		"failure_count":    framework.Int32OkToTF(apiObject.GetFailureCountOk()),
	}

	returnVar, d := types.ObjectValue(passwordPolicyLockoutTFObjectTypes, o)
	diags.Append(d...)

	return returnVar, diags
}

func passwordPolicyMinCharactersOkToTF(apiObject *management.PasswordPolicyMinCharacters, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(passwordPolicyMinCharactersTFObjectTypes), diags
	}

	o := map[string]attr.Value{
		"alphabetical_uppercase": framework.Int32OkToTF(apiObject.GetABCDEFGHIJKLMNOPQRSTUVWXYZOk()),
		"alphabetical_lowercase": framework.Int32OkToTF(apiObject.GetAbcdefghijklmnopqrstuvwxyzOk()),
		"numeric":                framework.Int32OkToTF(apiObject.GetVar0123456789Ok()),
		"special_characters":     framework.Int32OkToTF(apiObject.GetSpecialCharOk()),
	}

	returnVar, d := types.ObjectValue(passwordPolicyMinCharactersTFObjectTypes, o)
	diags.Append(d...)

	return returnVar, diags
}
