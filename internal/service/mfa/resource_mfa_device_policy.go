package mfa

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/patrickcping/pingone-go-sdk-v2/mfa"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/utils"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type MFADevicePolicyResource serviceClientType

type MFADevicePolicyResourceModel struct {
	Id                    pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId         pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	Name                  types.String                 `tfsdk:"name"`
	Authentication        types.Object                 `tfsdk:"authentication"`
	NewDeviceNotification types.String                 `tfsdk:"new_device_notification"`
	Sms                   types.Object                 `tfsdk:"sms"`
	Voice                 types.Object                 `tfsdk:"voice"`
	Email                 types.Object                 `tfsdk:"email"`
	Mobile                types.Object                 `tfsdk:"mobile"`
	Totp                  types.Object                 `tfsdk:"totp"`
	Fido2                 types.Object                 `tfsdk:"fido2"`
}

type MFADevicePolicyAuthenticationResourceModel struct {
	DeviceSelection types.String `tfsdk:"device_selection"`
}

type MFADevicePolicySmsResourceModel MFADevicePolicyOfflineDeviceResourceModel
type MFADevicePolicyVoiceResourceModel MFADevicePolicyOfflineDeviceResourceModel
type MFADevicePolicyEmailResourceModel MFADevicePolicyOfflineDeviceResourceModel
type MFADevicePolicyTotpResourceModel MFADevicePolicyOfflineDeviceResourceModel

type MFADevicePolicyOfflineDeviceResourceModel struct {
	Enabled         types.Bool   `tfsdk:"enabled"`
	Otp             types.Object `tfsdk:"otp"`
	PairingDisabled types.Bool   `tfsdk:"pairing_disabled"`
}

type MFADevicePolicyOfflineDeviceOtpResourceModel struct {
	MFADevicePolicyOtpResourceModel
	Lifetime types.Object `tfsdk:"lifetime"`
}

type MFADevicePolicyOtpResourceModel struct {
	Failure types.Object `tfsdk:"failure"`
}

type MFADevicePolicyFailureResourceModel struct {
	CoolDown types.Object `tfsdk:"cool_down"`
	Count    types.Int64  `tfsdk:"count"`
}

type MFADevicePolicyCooldownResourceModel MFADevicePolicyTimePeriodResourceModel
type MFADevicePolicyPushTimeoutResourceModel MFADevicePolicyTimePeriodResourceModel
type MFADevicePolicyLockDurationResourceModel MFADevicePolicyTimePeriodResourceModel
type MFADevicePolicyPairingKeyLifetimeResourceModel MFADevicePolicyTimePeriodResourceModel
type MFADevicePolicyTimePeriodResourceModel struct {
	Duration types.Int64  `tfsdk:"duration"`
	TimeUnit types.String `tfsdk:"time_unit"`
}

type MFADevicePolicyFido2ResourceModel struct {
	Enabled         types.Bool                   `tfsdk:"enabled"`
	Fido2PolicyId   pingonetypes.ResourceIDValue `tfsdk:"fido2_policy_id"`
	PairingDisabled types.Bool                   `tfsdk:"pairing_disabled"`
}

type MFADevicePolicyMobileResourceModel struct {
	Applications types.Map    `tfsdk:"applications"`
	Enabled      types.Bool   `tfsdk:"enabled"`
	Otp          types.Object `tfsdk:"otp"`
}

type MFADevicePolicyMobileApplicationResourceModel struct {
	AutoEnrolment       types.Object `tfsdk:"auto_enrollment"`
	DeviceAuthorization types.Object `tfsdk:"device_authorization"`
	IntegrityDetection  types.Object `tfsdk:"integrity_detection"`
	Otp                 types.Object `tfsdk:"otp"`
	PairingDisabled     types.Bool   `tfsdk:"pairing_disabled"`
	PairingKeyLifetime  types.Object `tfsdk:"pairing_key_lifetime"`
	Push                types.Object `tfsdk:"push"`
	PushLimit           types.Object `tfsdk:"push_limit"`
	PushTimeout         types.Object `tfsdk:"push_timeout"`
}

type MFADevicePolicyMobileApplicationAutoEnrolmentResourceModel struct {
	Enabled types.Bool `tfsdk:"enabled"`
}

type MFADevicePolicyMobileApplicationDeviceAuthorizationResourceModel struct {
	Enabled           types.Bool   `tfsdk:"enabled"`
	ExtraVerification types.String `tfsdk:"extra_verification"`
}

type MFADevicePolicyMobileApplicationOtpResourceModel MFADevicePolicyEnabledResourceModel
type MFADevicePolicyMobileApplicationPushResourceModel MFADevicePolicyEnabledResourceModel
type MFADevicePolicyEnabledResourceModel struct {
	Enabled types.Bool `tfsdk:"enabled"`
}

type MFADevicePolicyPushLimitResourceModel struct {
	Count        types.Int64  `tfsdk:"count"`
	LockDuration types.Object `tfsdk:"lock_duration"`
	TimePeriod   types.Object `tfsdk:"time_period"`
}

var (
	MFADevicePolicyAuthenticationTFObjectTypes = map[string]attr.Type{
		"device_selection": types.StringType,
	}

	MFADevicePolicyOfflineDeviceTFObjectTypes = map[string]attr.Type{
		"enabled":          types.Int64Type,
		"otp":              types.ObjectType{},
		"pairing_disabled": types.BoolType,
	}

	MFADevicePolicyOfflineDeviceOtpTFObjectTypes = map[string]attr.Type{
		"failure":  types.ObjectType{AttrTypes: MFADevicePolicyFailureTFObjectTypes},
		"lifetime": types.ObjectType{},
	}

	MFADevicePolicyFailureTFObjectTypes = map[string]attr.Type{
		"cool_down": types.ObjectType{AttrTypes: MFADevicePolicyTimePeriodTFObjectTypes},
		"count":     types.Int64Type,
	}

	MFADevicePolicyTimePeriodTFObjectTypes = map[string]attr.Type{
		"duration":  types.Int64Type,
		"time_unit": types.StringType,
	}

	MFADevicePolicyFido2TFObjectTypes = map[string]attr.Type{
		"enabled":          types.BoolType,
		"fido2_policy_id":  pingonetypes.ResourceIDType{},
		"pairing_disabled": types.BoolType,
	}

	MFADevicePolicyMobileTFObjectTypes = map[string]attr.Type{
		"applications": types.MapType{},
		"enabled":      types.BoolType,
		"otp":          types.ObjectType{},
	}
)

// Framework interfaces
var (
	_ resource.Resource                = &MFADevicePolicyResource{}
	_ resource.ResourceWithConfigure   = &MFADevicePolicyResource{}
	_ resource.ResourceWithImportState = &MFADevicePolicyResource{}
)

// New Object
func NewMFADevicePolicyResource() resource.Resource {
	return &MFADevicePolicyResource{}
}

// Metadata
func (r *MFADevicePolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mfa_device_policy"
}

func (r *MFADevicePolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	// schema descriptions and validation settings
	deviceSelectionDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that defines the device selection method.",
	).AllowedValuesEnum(mfa.AllowedEnumMFADevicePolicySelectionEnumValues).DefaultValue(string(mfa.ENUMMFADEVICEPOLICYSELECTION_DEFAULT_TO_FIRST))

	newDeviceNotificationDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that defines whether a user should be notified if a new authentication method has been added to their account.",
	).AllowedValuesEnum(mfa.AllowedEnumMFADevicePolicyNewDeviceNotificationEnumValues).DefaultValue(string(mfa.ENUMMFADEVICEPOLICYNEWDEVICENOTIFICATION_NONE))

	totpPairingDisabledDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that, when set to `true`, prevents users from pairing new devices with the TOTP method, though keeping it active in the policy for existing users. You can use this option if you want to phase out an existing authentication method but want to allow users to continue using the method for authentication for existing devices.",
	).DefaultValue(false)

	totpOtpFailureCoolDownDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the type of time unit for `duration`.",
	).AllowedValuesEnum(mfa.AllowedEnumTimeUnitEnumValues)

	fido2PairingDisabledDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that, when set to `true`, prevents users from pairing new devices with the FIDO2 method, though keeping it active in the policy for existing users. You can use this option if you want to phase out an existing authentication method but want to allow users to continue using the method for authentication for existing devices.",
	).DefaultValue(false)

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage MFA device policies for a PingOne environment.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment that contains the MFA device policy to manage."),
			),

			"name": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the MFA policy's unique name within the environment.").Description,
				Required:    true,
			},

			"device_selection": schema.StringAttribute{
				Description:         deviceSelectionDescription.Description,
				MarkdownDescription: deviceSelectionDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,

				Default: stringdefault.StaticString(string(mfa.ENUMMFADEVICEPOLICYSELECTION_DEFAULT_TO_FIRST)),

				Validators: []validator.String{
					stringvalidator.OneOf(utils.EnumSliceToStringSlice(mfa.AllowedEnumMFADevicePolicySelectionEnumValues)...),
				},
			},

			"new_device_notification": schema.StringAttribute{
				Description:         newDeviceNotificationDescription.Description,
				MarkdownDescription: newDeviceNotificationDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,

				Default: stringdefault.StaticString(string(mfa.ENUMMFADEVICEPOLICYNEWDEVICENOTIFICATION_NONE)),

				Validators: []validator.String{
					stringvalidator.OneOf(utils.EnumSliceToStringSlice(mfa.AllowedEnumMFADevicePolicyNewDeviceNotificationEnumValues)...),
				},
			},

			"sms": r.devicePolicyOfflineDeviceSchemaAttribute("SMS OTP"),

			"voice": r.devicePolicyOfflineDeviceSchemaAttribute("voice OTP"),

			"email": r.devicePolicyOfflineDeviceSchemaAttribute("email OTP"),

			"mobile": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that allows configuration of mobile push/OTP device authentication policy settings.").Description,
				Optional:    true,

				Attributes: map[string]schema.Attribute{
					"mfa_enabled": schema.BoolAttribute{
						Description:         usersMfaEnabledDescription.Description,
						MarkdownDescription: usersMfaEnabledDescription.MarkdownDescription,
						Required:            true,
					},
				},
			},

			"totp": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
				Optional:    true,

				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A boolean that specifies whether the TOTP method is enabled or disabled in the policy.").Description,
						Required:    true,
					},

					"pairing_disabled": schema.BoolAttribute{
						Description:         totpPairingDisabledDescription.Description,
						MarkdownDescription: totpPairingDisabledDescription.MarkdownDescription,
						Optional:            true,
						Computed:            true,

						Default: booldefault.StaticBool(false),
					},

					"otp": schema.SingleNestedAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that allows configuration of TOTP OTP settings.").Description,
						Optional:    true,

						Attributes: map[string]schema.Attribute{
							"failure": schema.SingleNestedAttribute{
								Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that allows configuration of TOTP OTP failure settings.").Description,
								Optional:    true,

								Attributes: map[string]schema.Attribute{
									"cool_down": schema.SingleNestedAttribute{
										Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that allows configuration of TOTP OTP failure cool down settings.").Description,
										Optional:    true,

										Attributes: map[string]schema.Attribute{
											"duration": schema.Int64Attribute{
												Description: framework.SchemaAttributeDescriptionFromMarkdown("An integer that defines the duration (number of time units) the user is blocked after reaching the maximum number of passcode failures.").Description,
												Required:    true,
											},

											"time_unit": schema.StringAttribute{
												Description:         totpOtpFailureCoolDownDescription.Description,
												MarkdownDescription: totpOtpFailureCoolDownDescription.MarkdownDescription,
												Required:            true,

												Validators: []validator.String{
													stringvalidator.OneOf(utils.EnumSliceToStringSlice(mfa.AllowedEnumTimeUnitEnumValues)...),
												},
											},
										},
									},

									"count": schema.Int64Attribute{
										Description: framework.SchemaAttributeDescriptionFromMarkdown("An integer that defines the maximum number of times that the OTP entry can fail for a user, before they are blocked.").Description,
										Required:    true,
									},
								},
							},
						},
					},
				},
			},

			"fido2": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that allows configuration of FIDO2 device authentication policy settings.").Description,
				Optional:    true,

				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A boolean that specifies whether the FIDO2 method is enabled or disabled in the policy.").Description,
						Required:    true,
					},

					"pairing_disabled": schema.BoolAttribute{
						Description:         fido2PairingDisabledDescription.Description,
						MarkdownDescription: fido2PairingDisabledDescription.MarkdownDescription,
						Optional:            true,
						Computed:            true,

						Default: booldefault.StaticBool(false),
					},

					"fido2_policy_id": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the resource UUID that represents the FIDO2 policy in PingOne. This property can be null / left undefined. When null, the environment's default FIDO2 Policy is used.  Must be a valid PingOne resource ID.").Description,
						Optional:    true,

						CustomType: pingonetypes.ResourceIDType{},
					},
				},
			},
		},
	}
}

func (r *MFADevicePolicyResource) devicePolicyOfflineDeviceSchemaAttribute(descriptionMethod string) schema.SingleNestedAttribute {

	pairingDisabledDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("A boolean that, when set to `true`, prevents users from pairing new devices with the %s method, though keeping it active in the policy for existing users. You can use this option if you want to phase out an existing authentication method but want to allow users to continue using the method for authentication for existing devices.", descriptionMethod),
	).DefaultValue(false)

	otpCoolDownDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the type of time unit for `duration`.",
	).AllowedValuesEnum(mfa.AllowedEnumTimeUnitEnumValues)

	return schema.SingleNestedAttribute{
		Description: framework.SchemaAttributeDescriptionFromMarkdown(fmt.Sprintf("A single object that allows configuration of %s device authentication policy settings.", descriptionMethod)).Description,
		Optional:    true,

		Attributes: map[string]schema.Attribute{
			"enabled": schema.BoolAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown(fmt.Sprintf("A boolean that specifies whether the %s method is enabled or disabled in the policy.", descriptionMethod)).Description,
				Required:    true,
			},

			"pairing_disabled": schema.BoolAttribute{
				Description:         pairingDisabledDescription.Description,
				MarkdownDescription: pairingDisabledDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,

				Default: booldefault.StaticBool(false),
			},

			"otp": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown(fmt.Sprintf("A single object that allows configuration of %s settings.", descriptionMethod)).Description,
				Optional:    true,

				Attributes: map[string]schema.Attribute{
					"failure": schema.SingleNestedAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown(fmt.Sprintf("A single object that allows configuration of %s failure settings.", descriptionMethod)).Description,
						Optional:    true,

						Attributes: map[string]schema.Attribute{
							"cool_down": schema.SingleNestedAttribute{
								Description: framework.SchemaAttributeDescriptionFromMarkdown(fmt.Sprintf("A single object that allows configuration of %s failure cool down settings.", descriptionMethod)).Description,
								Optional:    true,

								Attributes: map[string]schema.Attribute{
									"duration": schema.Int64Attribute{
										Description: framework.SchemaAttributeDescriptionFromMarkdown("An integer that defines the duration (number of time units) the user is blocked after reaching the maximum number of passcode failures.").Description,
										Required:    true,
									},

									"time_unit": schema.StringAttribute{
										Description:         otpCoolDownDescription.Description,
										MarkdownDescription: otpCoolDownDescription.MarkdownDescription,
										Required:            true,

										Validators: []validator.String{
											stringvalidator.OneOf(utils.EnumSliceToStringSlice(mfa.AllowedEnumTimeUnitEnumValues)...),
										},
									},
								},
							},

							"count": schema.Int64Attribute{
								Description: framework.SchemaAttributeDescriptionFromMarkdown("An integer that defines the maximum number of times that the OTP entry can fail for a user, before they are blocked.").Description,
								Required:    true,
							},
						},
					},

					"lifetime": schema.SingleNestedAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown(fmt.Sprintf("A single object that allows configuration of %s lifetime settings.", descriptionMethod)).Description,
						Optional:    true,

						Attributes: map[string]schema.Attribute{
							"duration": schema.Int64Attribute{
								Description: framework.SchemaAttributeDescriptionFromMarkdown("An integer that defines the duration (number of time units) that the passcode is valid before it expires.").Description,
								Required:    true,
							},

							"time_unit": schema.StringAttribute{
								Description:         otpCoolDownDescription.Description,
								MarkdownDescription: otpCoolDownDescription.MarkdownDescription,
								Required:            true,

								Validators: []validator.String{
									stringvalidator.OneOf(utils.EnumSliceToStringSlice(mfa.AllowedEnumTimeUnitEnumValues)...),
								},
							},
						},
					},
				},
			},
		},
	}
}

func (r *MFADevicePolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *MFADevicePolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state MFADevicePolicyResourceModel

	if r.Client.MFAAPIClient == nil {
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
	mFADevicePolicy, d := plan.expand(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *mfa.DeviceAuthenticationPolicyPostResponse
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.MFAAPIClient.DeviceAuthenticationPolicyApi.CreateDeviceAuthenticationPolicies(ctx, plan.EnvironmentId.ValueString()).DeviceAuthenticationPolicyPost(*mFADevicePolicy).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"CreateDeviceAuthenticationPolicies",
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
	resp.Diagnostics.Append(state.toStateCreate(response)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *MFADevicePolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *MFADevicePolicyResourceModel

	if r.Client.MFAAPIClient == nil {
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
	var response *mfa.DeviceAuthenticationPolicy
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.MFAAPIClient.DeviceAuthenticationPolicyApi.ReadOneDeviceAuthenticationPolicy(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"ReadOneDeviceAuthenticationPolicy",
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

func (r *MFADevicePolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state MFADevicePolicyResourceModel

	if r.Client.MFAAPIClient == nil {
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
	mFADevicePolicy, d := plan.expand(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *mfa.DeviceAuthenticationPolicy
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.MFAAPIClient.DeviceAuthenticationPolicyApi.UpdateDeviceAuthenticationPolicy(ctx, plan.EnvironmentId.ValueString(), plan.Id.ValueString()).DeviceAuthenticationPolicy(*mFADevicePolicy).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdateDeviceAuthenticationPolicy",
		framework.DefaultCustomError,
		nil,
		&response,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update the state to save
	state = plan

	// Save updated data into Terraform state
	resp.Diagnostics.Append(state.toState(response)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *MFADevicePolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *MFADevicePolicyResourceModel

	if r.Client.MFAAPIClient == nil {
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
			fR, fErr := r.Client.MFAAPIClient.DeviceAuthenticationPolicyApi.DeleteDeviceAuthenticationPolicy(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), nil, fR, fErr)
		},
		"DeleteDeviceAuthenticationPolicy",
		framework.CustomErrorResourceNotFoundWarning,
		nil,
		nil,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *MFADevicePolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	idComponents := []framework.ImportComponent{
		{
			Label:  "environment_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
		{
			Label:     "mfa_device_policy_id",
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

func (p *MFADevicePolicyResourceModel) expand(ctx context.Context) (*mfa.MFADevicePolicy, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Pairing
	var pairingPlan MFADevicePolicyPairingResourceModel
	diags.Append(p.Pairing.As(ctx, &pairingPlan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)
	if diags.HasError() {
		return nil, diags
	}
	pairing := mfa.NewMFADevicePolicyPairing(
		int32(pairingPlan.MaxAllowedDevices.ValueInt64()),
		mfa.EnumMFADevicePolicyPairingKeyFormat(pairingPlan.PairingKeyFormat.ValueString()),
	)

	// Main object
	data := mfa.NewMFADevicePolicy(
		*pairing,
	)

	// Lockout
	if !p.Lockout.IsNull() && !p.Lockout.IsUnknown() {
		var lockoutPlan MFADevicePolicyLockoutResourceModel
		diags.Append(p.Lockout.As(ctx, &lockoutPlan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}
		lockout := mfa.NewMFADevicePolicyLockout(
			int32(lockoutPlan.FailureCount.ValueInt64()),
		)

		if !lockoutPlan.DurationSeconds.IsNull() && !lockoutPlan.DurationSeconds.IsUnknown() {
			lockout.SetDurationSeconds(int32(lockoutPlan.DurationSeconds.ValueInt64()))
		}

		data.SetLockout(*lockout)
	}

	// Phone Extensions
	if !p.PhoneExtensions.IsNull() && !p.PhoneExtensions.IsUnknown() {
		var phoneExtensionsPlan MFADevicePolicyPhoneExtensionsResourceModel
		diags.Append(p.PhoneExtensions.As(ctx, &phoneExtensionsPlan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}
		phoneExtensions := mfa.NewMFADevicePolicyPhoneExtensions()

		if !phoneExtensionsPlan.Enabled.IsNull() && !phoneExtensionsPlan.Enabled.IsUnknown() {
			phoneExtensions.SetEnabled(phoneExtensionsPlan.Enabled.ValueBool())
		}

		data.SetPhoneExtensions(*phoneExtensions)
	}

	// Users
	if !p.Users.IsNull() && !p.Users.IsUnknown() {
		var usersPlan MFADevicePolicyUsersResourceModel
		diags.Append(p.Users.As(ctx, &usersPlan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}
		users := mfa.NewMFADevicePolicyUsers()

		if !usersPlan.MFAEnabled.IsNull() && !usersPlan.MFAEnabled.IsUnknown() {
			users.SetMfaEnabled(usersPlan.MFAEnabled.ValueBool())
		}

		data.SetUsers(*users)
	}

	return data, diags
}

func (p *MFADevicePolicyResourceModel) toState(apiObject *mfa.DeviceAuthenticationPolicy) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)
		return diags
	}

	var d diag.Diagnostics

	p.EnvironmentId = framework.PingOneResourceIDToTF(*apiObject.GetEnvironment().Id)

	p.Lockout, d = toStateLockout(apiObject.GetLockoutOk())
	diags.Append(d...)

	p.Pairing, d = toStatePairing(apiObject.GetPairingOk())
	diags.Append(d...)

	p.PhoneExtensions, d = toStatePhoneExtensions(apiObject.GetPhoneExtensionsOk())
	diags.Append(d...)

	p.Users, d = toStateUsers(apiObject.GetUsersOk())
	diags.Append(d...)

	return diags
}

func (p *MFADevicePolicyResourceModel) toStateCreate(apiObject *mfa.DeviceAuthenticationPolicyPostResponse) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)
		return diags
	}

	return p.toState(apiObject.DeviceAuthenticationPolicy)
}

func toStateLockout(apiObject *mfa.MFADevicePolicyLockout, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(MFADevicePolicyLockoutTFObjectTypes), nil
	}

	o := map[string]attr.Value{
		"failure_count":    framework.Int32OkToTF(apiObject.GetFailureCountOk()),
		"duration_seconds": framework.Int32OkToTF(apiObject.GetDurationSecondsOk()),
	}

	objValue, d := types.ObjectValue(MFADevicePolicyLockoutTFObjectTypes, o)
	diags.Append(d...)

	return objValue, diags
}

func toStatePairing(apiObject *mfa.MFADevicePolicyPairing, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(MFADevicePolicyPairingTFObjectTypes), nil
	}

	o := map[string]attr.Value{
		"max_allowed_devices": framework.Int32OkToTF(apiObject.GetMaxAllowedDevicesOk()),
		"pairing_key_format":  framework.EnumOkToTF(apiObject.GetPairingKeyFormatOk()),
	}

	objValue, d := types.ObjectValue(MFADevicePolicyPairingTFObjectTypes, o)
	diags.Append(d...)

	return objValue, diags
}

func toStatePhoneExtensions(apiObject *mfa.MFADevicePolicyPhoneExtensions, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(MFADevicePolicyPhoneExtensionsTFObjectTypes), nil
	}

	o := map[string]attr.Value{
		"enabled": framework.BoolOkToTF(apiObject.GetEnabledOk()),
	}

	objValue, d := types.ObjectValue(MFADevicePolicyPhoneExtensionsTFObjectTypes, o)
	diags.Append(d...)

	return objValue, diags
}

func toStateUsers(apiObject *mfa.MFADevicePolicyUsers, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(MFADevicePolicyUsersTFObjectTypes), nil
	}

	o := map[string]attr.Value{
		"mfa_enabled": framework.BoolOkToTF(apiObject.GetMfaEnabledOk()),
	}

	objValue, d := types.ObjectValue(MFADevicePolicyUsersTFObjectTypes, o)
	diags.Append(d...)

	return objValue, diags
}
