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
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
)

// Types
type PasswordPolicyDataSource serviceClientType

type PasswordPolicyDataSourceModel struct {
	Id                            pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId                 pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	PasswordPolicyId              pingonetypes.ResourceIDValue `tfsdk:"password_policy_id"`
	Name                          types.String                 `tfsdk:"name"`
	Description                   types.String                 `tfsdk:"description"`
	Default                       types.Bool                   `tfsdk:"default"`
	ExcludesCommonlyUsedPasswords types.Bool                   `tfsdk:"excludes_commonly_used_passwords"`
	ExcludesProfileData           types.Bool                   `tfsdk:"excludes_profile_data"`
	History                       types.Object                 `tfsdk:"history"`
	Length                        types.Object                 `tfsdk:"length"`
	Lockout                       types.Object                 `tfsdk:"lockout"`
	MinCharacters                 types.Object                 `tfsdk:"min_characters"`
	PasswordAgeMax                types.Int32                  `tfsdk:"password_age_max"`
	PasswordAgeMin                types.Int32                  `tfsdk:"password_age_min"`
	MaxRepeatedCharacters         types.Int32                  `tfsdk:"max_repeated_characters"`
	MinComplexity                 types.Int32                  `tfsdk:"min_complexity"`
	MinUniqueCharacters           types.Int32                  `tfsdk:"min_unique_characters"`
	NotSimilarToCurrent           types.Bool                   `tfsdk:"not_similar_to_current"`
	PopulationCount               types.Int32                  `tfsdk:"population_count"`
}

// Framework interfaces
var (
	_ datasource.DataSource = &PasswordPolicyDataSource{}
)

// New Object
func NewPasswordPolicyDataSource() datasource.DataSource {
	return &PasswordPolicyDataSource{}
}

// Metadata
func (r *PasswordPolicyDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_password_policy"
}

// Schema
func (r *PasswordPolicyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	nameLength := 1

	passwordPolicyIdDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the ID of the password policy to retrieve configuration for.  Must be a valid PingOne resource ID.",
	).ExactlyOneOf([]string{"password_policy_id", "name"})

	nameDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the name of the password policy to retrieve configuration for.",
	).ExactlyOneOf([]string{"password_policy_id", "name"})

	defaultDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether this password policy is enforced as the default within the environment. When set to `true`, all other password policies are set to `false`.",
	)

	excludeCommonlyUsedPasswordsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether to ensure the password is not one of the commonly used passwords.",
	)

	excludeProfileDataDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether to ensure the password is not an exact match for the value of any attribute in the user's profile, such as name, phone number, or address.",
	)

	passwordLengthMaxDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An integer that specifies the maximum number of characters allowed for the password. This property is not enforced when not present.",
	)

	passwordLengthMinDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An integer that specifies the minimum number of characters required for the password. This can be from `8` to `32` (inclusive). This property is not enforced when not present.",
	)

	minCharactersDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that specifies sets of characters that can be included, and the value is the minimum number of times one of the characters must appear in the user's password. The only allowed key values are `ABCDEFGHIJKLMNOPQRSTUVWXYZ`, `abcdefghijklmnopqrstuvwxyz`, `0123456789`, and `~!@#$%^&*()-_=+[]{}\\|;:,.<>/?`. This property is not enforced when not present.",
	)

	minCharactersAlphabeticalUppercaseDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An integer that specifies the count of alphabetical uppercase characters (`ABCDEFGHIJKLMNOPQRSTUVWXYZ`) that should feature in the user's password.  Fixed value of 1.",
	)

	minCharactersAlphabeticalLowercaseDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An integer that specifies the count of alphabetical uppercase characters (`abcdefghijklmnopqrstuvwxyz`) that should feature in the user's password.  Fixed value of 1.",
	)

	minCharactersNumericDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An integer that specifies the count of numeric characters (`0123456789`) that should feature in the user's password.  Fixed value of 1.",
	)

	minCharactersSpecialCharactersDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An integer that specifies the count of special characters (`~!@#$%^&*()-_=+[]{}\\|;:,.<>/?`) that should feature in the user's password.  Fixed value of 1.",
	)

	passwordAgeMaxDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An integer that specifies the maximum number of days the same password can be used before it must be changed. The value must be a positive, non-zero integer.  The value must be greater than the sum of `min` (if set) + 21 (the expiration warning interval for passwords).",
	)

	notSimilarToCurrentDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that, when set to `true`, ensures that the proposed password is not too similar to the user's current password based on the Levenshtein distance algorithm. The value of this parameter is evaluated only for password change actions in which the user enters both the current and the new password. By design, PingOne does not know the user's current password.",
	).DefaultValue(false)

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Datasource to retrieve a PingOne password policy in an environment by ID or by name.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment that is configured with the password policy.  Must be a valid PingOne resource ID.").Description,
				Required:    true,

				CustomType: pingonetypes.ResourceIDType{},
			},

			"password_policy_id": schema.StringAttribute{
				Description:         passwordPolicyIdDescription.Description,
				MarkdownDescription: passwordPolicyIdDescription.MarkdownDescription,
				Optional:            true,

				CustomType: pingonetypes.ResourceIDType{},

				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("name")),
				},
			},

			"name": schema.StringAttribute{
				Description:         nameDescription.Description,
				MarkdownDescription: nameDescription.MarkdownDescription,
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("password_policy_id")),
					stringvalidator.LengthAtLeast(nameLength),
				},
			},

			"description": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the description to apply to the password policy.").Description,
				Computed:    true,
			},

			"default": schema.BoolAttribute{
				Description:         defaultDescription.Description,
				MarkdownDescription: defaultDescription.MarkdownDescription,
				Computed:            true,
			},

			"excludes_commonly_used_passwords": schema.BoolAttribute{
				Description:         excludeCommonlyUsedPasswordsDescription.Description,
				MarkdownDescription: excludeCommonlyUsedPasswordsDescription.MarkdownDescription,
				Computed:            true,
			},

			"excludes_profile_data": schema.BoolAttribute{
				Description:         excludeProfileDataDescription.Description,
				MarkdownDescription: excludeProfileDataDescription.MarkdownDescription,
				Computed:            true,
			},

			"history": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that specifies settings to control the user's password history.").Description,
				Computed:    true,

				Attributes: map[string]schema.Attribute{
					"count": schema.Int32Attribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("An integer that specifies the number of prior passwords to keep for prevention of password re-use. The value must be a positive, non-zero integer.").Description,
						Computed:    true,
					},

					"retention_days": schema.Int32Attribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("An integer that specifies the length of time to keep recent passwords for prevention of password re-use. The value must be a positive, non-zero integer.").Description,
						Computed:    true,
					},
				},
			},

			"length": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that specifies settings to control the user's password length.").Description,
				Computed:    true,

				Attributes: map[string]schema.Attribute{
					"max": schema.Int32Attribute{
						Description:         passwordLengthMaxDescription.Description,
						MarkdownDescription: passwordLengthMaxDescription.MarkdownDescription,
						Computed:            true,
					},

					"min": schema.Int32Attribute{
						Description:         passwordLengthMinDescription.Description,
						MarkdownDescription: passwordLengthMinDescription.MarkdownDescription,
						Computed:            true,
					},
				},
			},

			"lockout": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that specifies settings to control the user's lockout on unsuccessful authentication attempts.").Description,
				Computed:    true,

				Attributes: map[string]schema.Attribute{
					"duration_seconds": schema.Int32Attribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("An integer that specifies the length of time before a password is automatically moved out of the lock out state. The value must be a positive, non-zero integer.").Description,
						Computed:    true,
					},

					"failure_count": schema.Int32Attribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("An integer that specifies the number of tries before a password is placed in the lockout state. The value must be a positive, non-zero integer.").Description,
						Computed:    true,
					},
				},
			},

			"min_characters": schema.SingleNestedAttribute{
				Description:         minCharactersDescription.Description,
				MarkdownDescription: minCharactersDescription.MarkdownDescription,
				Computed:            true,

				Attributes: map[string]schema.Attribute{
					"alphabetical_uppercase": schema.Int32Attribute{
						Description:         minCharactersAlphabeticalUppercaseDescription.Description,
						MarkdownDescription: minCharactersAlphabeticalUppercaseDescription.MarkdownDescription,
						Computed:            true,
					},

					"alphabetical_lowercase": schema.Int32Attribute{
						Description:         minCharactersAlphabeticalLowercaseDescription.Description,
						MarkdownDescription: minCharactersAlphabeticalLowercaseDescription.MarkdownDescription,
						Computed:            true,
					},

					"numeric": schema.Int32Attribute{
						Description:         minCharactersNumericDescription.Description,
						MarkdownDescription: minCharactersNumericDescription.MarkdownDescription,
						Computed:            true,
					},

					"special_characters": schema.Int32Attribute{
						Description:         minCharactersSpecialCharactersDescription.Description,
						MarkdownDescription: minCharactersSpecialCharactersDescription.MarkdownDescription,
						Computed:            true,
					},
				},
			},

			"password_age_max": schema.Int32Attribute{
				Description:         passwordAgeMaxDescription.Description,
				MarkdownDescription: passwordAgeMaxDescription.MarkdownDescription,
				Computed:            true,
			},

			"password_age_min": schema.Int32Attribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("An integer that specifies the minimum number of days a password must be used before changing. The value must be a positive, non-zero integer. This property is not enforced when not present.").Description,
				Computed:    true,
			},

			"max_repeated_characters": schema.Int32Attribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("An integer that specifies the maximum number of repeated characters allowed. This property is not enforced when not present.").Description,
				Computed:    true,
			},

			"min_complexity": schema.Int32Attribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("An integer that specifies the minimum complexity of the password based on the concept of password haystacks. The value is the number of days required to exhaust the entire search space during a brute force attack. This property is not enforced when not present.").Description,
				Computed:    true,
			},

			"min_unique_characters": schema.Int32Attribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("An integer that specifies the minimum number of unique characters required. This property is not enforced when not present.").Description,
				Computed:    true,
			},

			"not_similar_to_current": schema.BoolAttribute{
				Description:         notSimilarToCurrentDescription.Description,
				MarkdownDescription: notSimilarToCurrentDescription.MarkdownDescription,
				Computed:            true,
			},

			"population_count": schema.Int32Attribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("An integer that specifies the number of populations associated with the password policy.").Description,
				Computed:    true,
			},
		},
	}
}

func (r *PasswordPolicyDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *PasswordPolicyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *PasswordPolicyDataSourceModel

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

	var policyInstance *management.PasswordPolicy

	// Gateway API does not support SCIM filtering
	if !data.PasswordPolicyId.IsNull() {
		// Run the API call
		resp.Diagnostics.Append(framework.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				fO, fR, fErr := r.Client.ManagementAPIClient.PasswordPoliciesApi.ReadOnePasswordPolicy(ctx, data.EnvironmentId.ValueString(), data.PasswordPolicyId.ValueString()).Execute()
				return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
			},
			"ReadOnePasswordPolicy",
			framework.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
			&policyInstance,
		)...)
		if resp.Diagnostics.HasError() {
			return
		}

	} else if !data.Name.IsNull() {
		// Run the API call
		resp.Diagnostics.Append(framework.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				pagedIterator := r.Client.ManagementAPIClient.PasswordPoliciesApi.ReadAllPasswordPolicies(ctx, data.EnvironmentId.ValueString()).Execute()

				var initialHttpResponse *http.Response

				for pageCursor, err := range pagedIterator {
					if err != nil {
						return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), nil, pageCursor.HTTPResponse, err)
					}

					if initialHttpResponse == nil {
						initialHttpResponse = pageCursor.HTTPResponse
					}

					if passwordPolicies, ok := pageCursor.EntityArray.Embedded.GetPasswordPoliciesOk(); ok {

						for _, passwordPolicyObject := range passwordPolicies {
							if passwordPolicyObject.GetId() != "" && passwordPolicyObject.GetName() == data.Name.ValueString() {
								return &passwordPolicyObject, pageCursor.HTTPResponse, nil
							}
						}

					}
				}

				return nil, initialHttpResponse, nil
			},
			"ReadAllPasswordPolicies",
			framework.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
			&policyInstance,
		)...)
		if resp.Diagnostics.HasError() {
			return
		}

		if policyInstance == nil {
			resp.Diagnostics.AddError(
				"Cannot find the password policy from name",
				fmt.Sprintf("The password policy name %s for environment %s cannot be found", data.Name.String(), data.EnvironmentId.String()),
			)
			return
		}

	} else {
		resp.Diagnostics.AddError(
			"Missing parameter",
			"Cannot find the requested PingOne Password policy: password_policy_id or name argument must be set.",
		)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(policyInstance)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (p *PasswordPolicyDataSourceModel) toState(apiObject *management.PasswordPolicy) diag.Diagnostics {
	var diags, d diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.Id = framework.PingOneResourceIDOkToTF(apiObject.GetIdOk())
	p.PasswordPolicyId = framework.PingOneResourceIDToTF(apiObject.GetId())
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
