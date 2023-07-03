package mfa

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/patrickcping/pingone-go-sdk-v2/mfa"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/utils"
)

// Types
type FIDO2PolicyResource struct {
	client *mfa.APIClient
	region model.RegionMapping
}

type FIDO2PolicyResourceModel struct {
	Id                            types.String `tfsdk:"id"`
	EnvironmentId                 types.String `tfsdk:"environment_id"`
	Name                          types.String `tfsdk:"name"`
	Description                   types.String `tfsdk:"description"`
	Default                       types.Bool   `tfsdk:"default"`
	AttestationRequirements       types.String `tfsdk:"attestation_requirements"`
	AuthenticatorAttachment       types.String `tfsdk:"authenticator_attachment"`
	BackupEligibility             types.Object `tfsdk:"backup_eligibility"`
	DeviceDisplayName             types.String `tfsdk:"device_display_name"`
	DiscoverableCredentials       types.String `tfsdk:"discoverable_credentials"`
	MdsAuthenticatorsRequirements types.Object `tfsdk:"mds_authenticators_requirements"`
	RelyingPartyId                types.String `tfsdk:"relying_party_id"`
	UserDisplayNameAttributes     types.Object `tfsdk:"user_display_name_attributes"`
	UserVerification              types.Object `tfsdk:"user_verification"`
}

type FIDO2PolicyBackupEligibilityResourceModel struct {
	Allow                       types.Bool `tfsdk:"allow"`
	EnforceDuringAuthentication types.Bool `tfsdk:"enforce_during_authentication"`
}

type FIDO2PolicyMdsAuthenticatorsRequirementsResourceModel struct {
	AllowedAuthenticatorIDs     types.Set    `tfsdk:"allowed_authenticator_ids"`
	EnforceDuringAuthentication types.Bool   `tfsdk:"enforce_during_authentication"`
	Option                      types.String `tfsdk:"option"`
}

type FIDO2PolicyUserDisplayNameAttributesResourceModel struct {
	Attributes types.Set `tfsdk:"attributes"`
}

type FIDO2PolicyUserDisplayNameAttributesAttributesResourceModel struct {
	Name          types.String `tfsdk:"name"`
	SubAttributes types.Set    `tfsdk:"sub_attributes"`
}

type FIDO2PolicyUserDisplayNameAttributesAttributesSubAttributesResourceModel struct {
	Name types.String `tfsdk:"name"`
}

type FIDO2PolicyUserVerificationResourceModel struct {
	EnforceDuringAuthentication types.Bool   `tfsdk:"enforce_during_authentication"`
	Option                      types.String `tfsdk:"option"`
}

var (
	fido2PolicyBackupEligibilityTFObjectTypes = map[string]attr.Type{
		"allow":                         types.BoolType,
		"enforce_during_authentication": types.BoolType,
	}

	fido2PolicyMdsAuthenticatorRequirementsTFObjectTypes = map[string]attr.Type{
		"allowed_authenticator_ids":     types.SetType{ElemType: types.StringType},
		"enforce_during_authentication": types.BoolType,
		"option":                        types.StringType,
	}

	fido2PolicyUserDisplayNameAttributesTFObjectTypes = map[string]attr.Type{
		"attributes": types.SetType{ElemType: types.ObjectType{AttrTypes: fido2PolicyUserDisplayNameAttributesAttributesTFObjectTypes}},
	}

	fido2PolicyUserDisplayNameAttributesAttributesTFObjectTypes = map[string]attr.Type{
		"sub_attributes": types.SetType{ElemType: types.ObjectType{AttrTypes: fido2PolicyUserDisplayNameAttributesAttributesSubAttributesTFObjectTypes}},
		"name":           types.StringType,
	}

	fido2PolicyUserDisplayNameAttributesAttributesSubAttributesTFObjectTypes = map[string]attr.Type{
		"name": types.StringType,
	}

	fido2PolicyUserVerificationTFObjectTypes = map[string]attr.Type{
		"enforce_during_authentication": types.BoolType,
		"option":                        types.StringType,
	}
)

// Framework interfaces
var (
	_ resource.Resource                = &FIDO2PolicyResource{}
	_ resource.ResourceWithConfigure   = &FIDO2PolicyResource{}
	_ resource.ResourceWithImportState = &FIDO2PolicyResource{}
)

// New Object
func NewFIDO2PolicyResource() resource.Resource {
	return &FIDO2PolicyResource{}
}

// Metadata
func (r *FIDO2PolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mfa_fido2_policy"
}

func (r *FIDO2PolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	// schema descriptions and validation settings
	const attrMinLength = 1
	const attrMinColumns = 1
	const attrMaxColumns = 3
	const attrDefaultVersion = 5
	const attrMinPercent = 0
	const attrMaxPercent = 100
	const imageMaxSize = 50000

	const attrNameMaxLength = 256
	const attrDeviceDisplayNameMaxLength = 100

	attestationRequirementsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"",
	).AllowedValuesEnum(mfa.AllowedEnumFIDO2PolicyAttestationRequirementsEnumValues)

	authenticatorAttachmentDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"",
	).AllowedValuesEnum(mfa.AllowedEnumFIDO2PolicyAuthenticatorAttachmentEnumValues)

	backupEligibilityAllowDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"",
	)

	backupEligibilityEnforceDuringAuthnDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"",
	)

	deviceDisplayNameDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"",
	)

	discoverableCredentialsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"",
	).AllowedValuesEnum(mfa.AllowedEnumFIDO2PolicyDiscoverableCredentialsEnumValues)

	mdsAuthenticatorRequirementsAllowedAuthenticatorIDsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"",
	)

	mdsAuthenticatorRequirementsEnforceDuringAuthnDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"",
	)

	mdsAuthenticatorRequirementsOptionDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"",
	).AllowedValuesEnum(mfa.AllowedEnumFIDO2PolicyMDSAuthenticatorOptionEnumValues)

	relyingPartyIDDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"",
	)

	userDisplayNameAttributesAttributesDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"",
	)

	userDisplayNameAttributesAttributesNameDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"",
	)

	userDisplayNameAttributesAttributesSubAttributesDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"",
	)

	userDisplayNameAttributesAttributesSubAttributesNameDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"",
	)

	userVerificationEnforceDuringAuthnDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"",
	)

	userVerificationOptionDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"",
	).AllowedValuesEnum(mfa.AllowedEnumFIDO2PolicyUserVerificationOptionEnumValues)

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage FIDO2 policies in a PingOne environment.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment to configure the FIDO2 policy in."),
			),

			"name": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the unique, friendly name for this FIDO2 policy.").Description,
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
					stringvalidator.LengthAtMost(attrNameMaxLength),
				},
			},

			"description": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the description of the FIDO2 policy.").Description,
				Optional:    true,
			},

			"default": schema.BoolAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A boolean that describes whether this policy should serve as the default FIDO policy.").Description,
				Computed:    true,
			},

			"attestation_requirements": schema.StringAttribute{
				Description:         attestationRequirementsDescription.Description,
				MarkdownDescription: attestationRequirementsDescription.MarkdownDescription,
				Required:            true,

				Validators: []validator.String{
					stringvalidator.OneOf(utils.EnumSliceToStringSlice(mfa.AllowedEnumFIDO2PolicyAttestationRequirementsEnumValues)...),
				},
			},

			"authenticator_attachment": schema.StringAttribute{
				Description:         authenticatorAttachmentDescription.Description,
				MarkdownDescription: authenticatorAttachmentDescription.MarkdownDescription,
				Required:            true,

				Validators: []validator.String{
					stringvalidator.OneOf(utils.EnumSliceToStringSlice(mfa.AllowedEnumFIDO2PolicyAuthenticatorAttachmentEnumValues)...),
				},
			},

			"backup_eligibility": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A single nested object that specifies the backup eligibility of FIDO2 devices.").Description,
				Required:    true,

				Attributes: map[string]schema.Attribute{
					"allow": schema.BoolAttribute{
						Description:         backupEligibilityAllowDescription.Description,
						MarkdownDescription: backupEligibilityAllowDescription.MarkdownDescription,
						Required:            true,
					},

					"enforce_during_authentication": schema.BoolAttribute{
						Description:         backupEligibilityEnforceDuringAuthnDescription.Description,
						MarkdownDescription: backupEligibilityEnforceDuringAuthnDescription.MarkdownDescription,
						Required:            true,
					},
				},
			},

			"device_display_name": schema.StringAttribute{
				Description:         deviceDisplayNameDescription.Description,
				MarkdownDescription: deviceDisplayNameDescription.MarkdownDescription,
				Required:            true,

				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
					stringvalidator.LengthAtMost(attrDeviceDisplayNameMaxLength),
				},
			},

			"discoverable_credentials": schema.StringAttribute{
				Description:         discoverableCredentialsDescription.Description,
				MarkdownDescription: discoverableCredentialsDescription.MarkdownDescription,
				Required:            true,

				Validators: []validator.String{
					stringvalidator.OneOf(utils.EnumSliceToStringSlice(mfa.AllowedEnumFIDO2PolicyDiscoverableCredentialsEnumValues)...),
				},
			},

			"mds_authenticators_requirements": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A single nested object that specifies MDS authenticator requirements.").Description,
				Required:    true,

				Attributes: map[string]schema.Attribute{
					"allowed_authenticator_ids": schema.SetAttribute{
						Description:         mdsAuthenticatorRequirementsAllowedAuthenticatorIDsDescription.Description,
						MarkdownDescription: mdsAuthenticatorRequirementsAllowedAuthenticatorIDsDescription.MarkdownDescription,
						Optional:            true,

						ElementType: types.StringType,

						Validators: []validator.Set{
							setvalidator.SizeAtLeast(1),
						},
					},

					"enforce_during_authentication": schema.BoolAttribute{
						Description:         mdsAuthenticatorRequirementsEnforceDuringAuthnDescription.Description,
						MarkdownDescription: mdsAuthenticatorRequirementsEnforceDuringAuthnDescription.MarkdownDescription,
						Required:            true,
					},

					"option": schema.StringAttribute{
						Description:         mdsAuthenticatorRequirementsOptionDescription.Description,
						MarkdownDescription: mdsAuthenticatorRequirementsOptionDescription.MarkdownDescription,
						Required:            true,

						Validators: []validator.String{
							stringvalidator.OneOf(utils.EnumSliceToStringSlice(mfa.AllowedEnumFIDO2PolicyMDSAuthenticatorOptionEnumValues)...),
						},
					},
				},
			},

			"relying_party_id": schema.StringAttribute{
				Description:         relyingPartyIDDescription.Description,
				MarkdownDescription: relyingPartyIDDescription.MarkdownDescription,
				Required:            true,

				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
				},
			},

			"user_display_name_attributes": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A single nested object that specifies user display name attributes.").Description,
				Required:    true,

				Attributes: map[string]schema.Attribute{
					"attributes": schema.SetNestedAttribute{
						Description:         userDisplayNameAttributesAttributesDescription.Description,
						MarkdownDescription: userDisplayNameAttributesAttributesDescription.MarkdownDescription,
						Required:            true,

						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Description:         userDisplayNameAttributesAttributesNameDescription.Description,
									MarkdownDescription: userDisplayNameAttributesAttributesNameDescription.MarkdownDescription,
									Required:            true,

									Validators: []validator.String{
										stringvalidator.LengthAtLeast(attrMinLength),
									},
								},

								"sub_attributes": schema.SetNestedAttribute{
									Description:         userDisplayNameAttributesAttributesSubAttributesDescription.Description,
									MarkdownDescription: userDisplayNameAttributesAttributesSubAttributesDescription.MarkdownDescription,
									Optional:            true,

									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"name": schema.StringAttribute{
												Description:         userDisplayNameAttributesAttributesSubAttributesNameDescription.Description,
												MarkdownDescription: userDisplayNameAttributesAttributesSubAttributesNameDescription.MarkdownDescription,
												Required:            true,

												Validators: []validator.String{
													stringvalidator.LengthAtLeast(attrMinLength),
												},
											},
										},
									},

									Validators: []validator.Set{
										setvalidator.SizeAtLeast(1),
									},
								},
							},
						},

						Validators: []validator.Set{
							setvalidator.SizeAtLeast(1),
						},
					},
				},
			},

			"user_verification": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A single nested object that specifies user verification settings.").Description,
				Required:    true,

				Attributes: map[string]schema.Attribute{
					"enforce_during_authentication": schema.BoolAttribute{
						Description:         userVerificationEnforceDuringAuthnDescription.Description,
						MarkdownDescription: userVerificationEnforceDuringAuthnDescription.MarkdownDescription,
						Required:            true,
					},

					"option": schema.StringAttribute{
						Description:         userVerificationOptionDescription.Description,
						MarkdownDescription: userVerificationOptionDescription.MarkdownDescription,
						Required:            true,

						Validators: []validator.String{
							stringvalidator.OneOf(utils.EnumSliceToStringSlice(mfa.AllowedEnumFIDO2PolicyUserVerificationOptionEnumValues)...),
						},
					},
				},
			},
		},
	}
}

func (r *FIDO2PolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
	r.region = resourceConfig.Client.API.Region
}

func (r *FIDO2PolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state FIDO2PolicyResourceModel

	if r.client == nil {
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
	fido2Policy, d := plan.expand(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	response, d := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return r.client.FIDO2PolicyApi.CreateFIDO2Policy(ctx, plan.EnvironmentId.ValueString()).FIDO2Policy(*fido2Policy).Execute()
		},
		"CreateFIDO2Policy",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
	)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the state to save
	state = plan

	// Save updated data into Terraform state
	resp.Diagnostics.Append(state.toState(response.(*mfa.FIDO2Policy))...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *FIDO2PolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *FIDO2PolicyResourceModel

	if r.client == nil {
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
	response, diags := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return r.client.FIDO2PolicyApi.ReadOneFIDO2Policy(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
		},
		"ReadOneFIDO2Policy",
		framework.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Remove from state if resource is not found
	if response == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(response.(*mfa.FIDO2Policy))...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FIDO2PolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state FIDO2PolicyResourceModel

	if r.client == nil {
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
	fido2Policy, d := plan.expand(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	response, diags := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return r.client.FIDO2PolicyApi.UpdateFIDO2Policy(ctx, plan.EnvironmentId.ValueString(), plan.Id.ValueString()).FIDO2Policy(*fido2Policy).Execute()
		},
		"UpdateFIDO2Policy",
		framework.DefaultCustomError,
		nil,
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update the state to save
	state = plan

	// Save updated data into Terraform state
	resp.Diagnostics.Append(state.toState(response.(*mfa.FIDO2Policy))...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *FIDO2PolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *FIDO2PolicyResourceModel

	if r.client == nil {
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
	_, diags := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			r, err := r.client.FIDO2PolicyApi.DeleteFIDO2Policy(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return nil, r, err
		},
		"DeleteFIDO2Policy",
		framework.CustomErrorResourceNotFoundWarning,
		nil,
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *FIDO2PolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	splitLength := 2
	attributes := strings.SplitN(req.ID, "/", splitLength)

	if len(attributes) != splitLength {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("invalid id (\"%s\") specified, should be in format \"environment_id/fido2_policy_id/\"", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("environment_id"), attributes[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), attributes[1])...)
}

func (p *FIDO2PolicyResourceModel) expand(ctx context.Context) (*mfa.FIDO2Policy, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Backup eligibility
	var backupEligibilityPlan FIDO2PolicyBackupEligibilityResourceModel
	diags.Append(p.BackupEligibility.As(ctx, &backupEligibilityPlan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)
	if diags.HasError() {
		return nil, diags
	}
	backupEligibility := mfa.NewFIDO2PolicyBackupEligibility(
		backupEligibilityPlan.Allow.ValueBool(),
		backupEligibilityPlan.EnforceDuringAuthentication.ValueBool(),
	)

	// MDS Authenticator Requirements
	var mdsAuthenticatorRequirementsPlan FIDO2PolicyMdsAuthenticatorsRequirementsResourceModel
	diags.Append(p.MdsAuthenticatorsRequirements.As(ctx, &mdsAuthenticatorRequirementsPlan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)
	if diags.HasError() {
		return nil, diags
	}

	mdsAuthenticatorsRequirements := mfa.NewFIDO2PolicyMdsAuthenticatorsRequirements(
		mdsAuthenticatorRequirementsPlan.EnforceDuringAuthentication.ValueBool(),
		mfa.EnumFIDO2PolicyMDSAuthenticatorOption(mdsAuthenticatorRequirementsPlan.Option.ValueString()),
	)

	if !mdsAuthenticatorRequirementsPlan.AllowedAuthenticatorIDs.IsNull() && !mdsAuthenticatorRequirementsPlan.AllowedAuthenticatorIDs.IsUnknown() {
		allowedAuthenticators := make([]mfa.FIDO2PolicyMdsAuthenticatorsRequirementsAllowedAuthenticatorsInner, 0)

		var allowedAuthenticatorIDsPlan []string
		diags.Append(mdsAuthenticatorRequirementsPlan.AllowedAuthenticatorIDs.ElementsAs(ctx, &allowedAuthenticatorIDsPlan, false)...)
		if diags.HasError() {
			return nil, diags
		}

		for _, allowedAuthenticatorIDPlan := range allowedAuthenticatorIDsPlan {

			allowedAuthenticator := *mfa.NewFIDO2PolicyMdsAuthenticatorsRequirementsAllowedAuthenticatorsInner(
				allowedAuthenticatorIDPlan,
			)

			allowedAuthenticators = append(allowedAuthenticators, allowedAuthenticator)
		}

		mdsAuthenticatorsRequirements.SetAllowedAuthenticators(allowedAuthenticators)
	}

	// User display name attributes
	var userDisplayNameAttributesPlan FIDO2PolicyUserDisplayNameAttributesResourceModel
	diags.Append(p.UserDisplayNameAttributes.As(ctx, &userDisplayNameAttributesPlan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)
	if diags.HasError() {
		return nil, diags
	}

	attributes := make([]mfa.FIDO2PolicyUserDisplayNameAttributesAttributesInner, 0)
	if !userDisplayNameAttributesPlan.Attributes.IsNull() && !userDisplayNameAttributesPlan.Attributes.IsUnknown() {

		var userDisplayNameAttributesAttributesPlan []FIDO2PolicyUserDisplayNameAttributesAttributesResourceModel
		diags.Append(userDisplayNameAttributesPlan.Attributes.ElementsAs(ctx, &userDisplayNameAttributesAttributesPlan, false)...)
		if diags.HasError() {
			return nil, diags
		}

		for _, attributePlan := range userDisplayNameAttributesAttributesPlan {

			attribute := *mfa.NewFIDO2PolicyUserDisplayNameAttributesAttributesInner(
				attributePlan.Name.ValueString(),
			)

			if !attributePlan.SubAttributes.IsNull() && !attributePlan.SubAttributes.IsUnknown() {
				var userDisplayNameAttributesAttributesSubAttributesPlan []FIDO2PolicyUserDisplayNameAttributesAttributesSubAttributesResourceModel
				diags.Append(attributePlan.SubAttributes.ElementsAs(ctx, &userDisplayNameAttributesAttributesSubAttributesPlan, false)...)
				if diags.HasError() {
					return nil, diags
				}

				subAttributes := make([]mfa.FIDO2PolicyUserDisplayNameAttributesAttributesInnerSubAttributesInner, 0)

				for _, subAttributePlan := range userDisplayNameAttributesAttributesSubAttributesPlan {
					subAttributes = append(subAttributes, *mfa.NewFIDO2PolicyUserDisplayNameAttributesAttributesInnerSubAttributesInner(
						subAttributePlan.Name.ValueString(),
					))
				}

				attribute.SetSubAttributes(subAttributes)
			}

			attributes = append(attributes, attribute)
		}
	}

	userDisplayNameAttributes := mfa.NewFIDO2PolicyUserDisplayNameAttributes(
		attributes,
	)

	// User verification
	var userVerificationPlan FIDO2PolicyUserVerificationResourceModel
	diags.Append(p.UserVerification.As(ctx, &userVerificationPlan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)
	if diags.HasError() {
		return nil, diags
	}
	userVerification := mfa.NewFIDO2PolicyUserVerification(
		userVerificationPlan.EnforceDuringAuthentication.ValueBool(),
		mfa.EnumFIDO2PolicyUserVerificationOption(userVerificationPlan.Option.ValueString()),
	)

	// Main object
	data := mfa.NewFIDO2Policy(
		mfa.EnumFIDO2PolicyAttestationRequirements(p.AttestationRequirements.ValueString()),
		mfa.EnumFIDO2PolicyAuthenticatorAttachment(p.AuthenticatorAttachment.ValueString()),
		*backupEligibility,
		p.DeviceDisplayName.ValueString(),
		mfa.EnumFIDO2PolicyDiscoverableCredentials(p.DiscoverableCredentials.ValueString()),
		*mdsAuthenticatorsRequirements,
		p.Name.ValueString(),
		p.RelyingPartyId.ValueString(),
		*userDisplayNameAttributes,
		*userVerification,
	)

	if !p.Description.IsNull() && !p.Description.IsUnknown() {
		data.SetDescription(p.Description.ValueString())
	}

	data.SetDefault(false)

	return data, diags
}

func (p *FIDO2PolicyResourceModel) toState(apiObject *mfa.FIDO2Policy) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)
		return diags
	}

	var d diag.Diagnostics

	p.Id = framework.StringOkToTF(apiObject.GetIdOk())
	p.EnvironmentId = framework.StringToTF(*apiObject.GetEnvironment().Id)
	p.Name = framework.StringOkToTF(apiObject.GetNameOk())
	p.Description = framework.StringOkToTF(apiObject.GetDescriptionOk())
	p.Default = framework.BoolOkToTF(apiObject.GetDefaultOk())
	p.AttestationRequirements = framework.EnumOkToTF(apiObject.GetAttestationRequirementsOk())
	p.AuthenticatorAttachment = framework.EnumOkToTF(apiObject.GetAuthenticatorAttachmentOk())

	p.BackupEligibility, d = toStateBackupEligibility(apiObject.GetBackupEligibilityOk())
	diags.Append(d...)

	p.DeviceDisplayName = framework.StringOkToTF(apiObject.GetDeviceDisplayNameOk())
	p.DiscoverableCredentials = framework.EnumOkToTF(apiObject.GetDiscoverableCredentialsOk())

	p.MdsAuthenticatorsRequirements, d = toStateMdsAuthenticatorsRequirements(apiObject.GetMdsAuthenticatorsRequirementsOk())
	diags.Append(d...)

	p.RelyingPartyId = framework.StringOkToTF(apiObject.GetRelyingPartyIdOk())

	p.UserDisplayNameAttributes, d = toStateUserDisplayNameAttributes(apiObject.GetUserDisplayNameAttributesOk())
	diags.Append(d...)

	p.UserVerification, d = toStateUserVerification(apiObject.GetUserVerificationOk())
	diags.Append(d...)

	return diags
}

func toStateBackupEligibility(apiObject *mfa.FIDO2PolicyBackupEligibility, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(fido2PolicyBackupEligibilityTFObjectTypes), nil
	}

	o := map[string]attr.Value{
		"allow":                         framework.BoolOkToTF(apiObject.GetAllowOk()),
		"enforce_during_authentication": framework.BoolOkToTF(apiObject.GetEnforceDuringAuthenticationOk()),
	}

	objValue, d := types.ObjectValue(fido2PolicyBackupEligibilityTFObjectTypes, o)
	diags.Append(d...)

	return objValue, diags
}

func toStateMdsAuthenticatorsRequirements(apiObject *mfa.FIDO2PolicyMdsAuthenticatorsRequirements, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(fido2PolicyMdsAuthenticatorRequirementsTFObjectTypes), nil
	}

	o := map[string]attr.Value{
		"allowed_authenticator_ids":     types.SetNull(types.StringType),
		"enforce_during_authentication": framework.BoolOkToTF(apiObject.GetEnforceDuringAuthenticationOk()),
		"option":                        framework.EnumOkToTF(apiObject.GetOptionOk()),
	}

	if v, ok := apiObject.GetAllowedAuthenticatorsOk(); ok {
		allowedAuthenticatorsList := make([]string, 0)
		for _, item := range v {
			if id, ok := item.GetIdOk(); ok {
				allowedAuthenticatorsList = append(allowedAuthenticatorsList, *id)
			}
		}

		o["allowed_authenticator_ids"] = framework.StringSetToTF(allowedAuthenticatorsList)
	}

	objValue, d := types.ObjectValue(fido2PolicyMdsAuthenticatorRequirementsTFObjectTypes, o)
	diags.Append(d...)

	return objValue, diags
}

func toStateUserDisplayNameAttributes(apiObject *mfa.FIDO2PolicyUserDisplayNameAttributes, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(fido2PolicyUserDisplayNameAttributesTFObjectTypes), nil
	}

	attributesSet, d := toStateUserDisplayNameAttributesAttributes(apiObject.GetAttributesOk())
	diags.Append(d...)

	o := map[string]attr.Value{
		"attributes": attributesSet,
	}

	objValue, d := types.ObjectValue(fido2PolicyUserDisplayNameAttributesTFObjectTypes, o)
	diags.Append(d...)

	return objValue, diags
}

func toStateUserDisplayNameAttributesAttributes(apiObject []mfa.FIDO2PolicyUserDisplayNameAttributesAttributesInner, ok bool) (types.Set, diag.Diagnostics) {
	var diags diag.Diagnostics
	tfObjType := types.ObjectType{AttrTypes: fido2PolicyUserDisplayNameAttributesAttributesTFObjectTypes}

	if !ok || len(apiObject) == 0 {
		return types.SetNull(tfObjType), diags
	}

	flattenedList := []attr.Value{}
	for _, v := range apiObject {

		subAttributes, d := toStateUserDisplayNameAttributesAttributesSubAttributes(v.GetSubAttributesOk())
		diags.Append(d...)

		objMap := map[string]attr.Value{
			"sub_attributes": subAttributes,
			"name":           framework.StringOkToTF(v.GetNameOk()),
		}

		flattenedObj, d := types.ObjectValue(fido2PolicyUserDisplayNameAttributesAttributesTFObjectTypes, objMap)
		diags.Append(d...)

		flattenedList = append(flattenedList, flattenedObj)
	}

	returnVar, d := types.SetValue(tfObjType, flattenedList)
	diags.Append(d...)

	return returnVar, diags
}

func toStateUserDisplayNameAttributesAttributesSubAttributes(apiObject []mfa.FIDO2PolicyUserDisplayNameAttributesAttributesInnerSubAttributesInner, ok bool) (types.Set, diag.Diagnostics) {
	var diags diag.Diagnostics
	tfObjType := types.ObjectType{AttrTypes: fido2PolicyUserDisplayNameAttributesAttributesSubAttributesTFObjectTypes}

	if !ok || len(apiObject) == 0 {
		return types.SetNull(tfObjType), diags
	}

	flattenedList := []attr.Value{}
	for _, v := range apiObject {

		objMap := map[string]attr.Value{
			"name": framework.StringOkToTF(v.GetNameOk()),
		}

		flattenedObj, d := types.ObjectValue(fido2PolicyUserDisplayNameAttributesAttributesSubAttributesTFObjectTypes, objMap)
		diags.Append(d...)

		flattenedList = append(flattenedList, flattenedObj)
	}

	returnVar, d := types.SetValue(tfObjType, flattenedList)
	diags.Append(d...)

	return returnVar, diags
}

func toStateUserVerification(apiObject *mfa.FIDO2PolicyUserVerification, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(fido2PolicyUserVerificationTFObjectTypes), nil
	}

	o := map[string]attr.Value{
		"enforce_during_authentication": framework.BoolOkToTF(apiObject.GetEnforceDuringAuthenticationOk()),
		"option":                        framework.EnumOkToTF(apiObject.GetOptionOk()),
	}

	objValue, d := types.ObjectValue(fido2PolicyUserVerificationTFObjectTypes, o)
	diags.Append(d...)

	return objValue, diags
}
