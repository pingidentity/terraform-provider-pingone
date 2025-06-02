// Copyright © 2025 Ping Identity Corporation

package mfa

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/patrickcping/pingone-go-sdk-v2/mfa"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/utils"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type MFASettingsResource serviceClientType

type mFASettingsResourceModelV1 struct {
	EnvironmentId   pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	Lockout         types.Object                 `tfsdk:"lockout"`
	Pairing         types.Object                 `tfsdk:"pairing"`
	PhoneExtensions types.Object                 `tfsdk:"phone_extensions"`
	Users           types.Object                 `tfsdk:"users"`
}

type mFASettingsLockoutResourceModelV1 struct {
	FailureCount    types.Int32 `tfsdk:"failure_count"`
	DurationSeconds types.Int32 `tfsdk:"duration_seconds"`
}

type mFASettingsPairingResourceModelV1 struct {
	MaxAllowedDevices types.Int32  `tfsdk:"max_allowed_devices"`
	PairingKeyFormat  types.String `tfsdk:"pairing_key_format"`
}

type mFASettingsPhoneExtensionsResourceModelV1 struct {
	Enabled types.Bool `tfsdk:"enabled"`
}

type mFASettingsUsersResourceModelV1 struct {
	MFAEnabled types.Bool `tfsdk:"mfa_enabled"`
}

var (
	MFASettingsLockoutTFObjectTypes = map[string]attr.Type{
		"failure_count":    types.Int32Type,
		"duration_seconds": types.Int32Type,
	}

	MFASettingsPairingTFObjectTypes = map[string]attr.Type{
		"max_allowed_devices": types.Int32Type,
		"pairing_key_format":  types.StringType,
	}

	MFASettingsPhoneExtensionsTFObjectTypes = map[string]attr.Type{
		"enabled": types.BoolType,
	}

	MFASettingsUsersTFObjectTypes = map[string]attr.Type{
		"mfa_enabled": types.BoolType,
	}
)

// Framework interfaces
var (
	_ resource.Resource                 = &MFASettingsResource{}
	_ resource.ResourceWithConfigure    = &MFASettingsResource{}
	_ resource.ResourceWithImportState  = &MFASettingsResource{}
	_ resource.ResourceWithUpgradeState = &MFASettingsResource{}
)

// New Object
func NewMFASettingsResource() resource.Resource {
	return &MFASettingsResource{}
}

// Metadata
func (r *MFASettingsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mfa_settings"
}

func (r *MFASettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	// schema descriptions and validation settings
	const maxAllowedDevicesDefault = 5
	const maxAllowedDevicesMin = 1
	const maxAllowedDevicesMax = 15

	pairingMaxAllowedDevicesDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("An integer that defines the maximum number of MFA devices each user can have. This can be any number from `%d` to `%d`. All devices that are Active or Blocked are subject to this limit.", maxAllowedDevicesMin, maxAllowedDevicesMax),
	).DefaultValue(maxAllowedDevicesDefault)

	pairingPairingKeyFormatDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that controls the type of pairing key issued.",
	).AllowedValuesComplex(map[string]string{
		string(mfa.ENUMMFASETTINGSPAIRINGKEYFORMAT_NUMERIC):      "12-digit key",
		string(mfa.ENUMMFASETTINGSPAIRINGKEYFORMAT_ALPHANUMERIC): "16-character alphanumeric key",
	})

	phoneExtensionsEnabledDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean when set to `true` to allow one-time passwords to be delivered via voice to phone numbers that include extensions. Set to `false` to disable support for phone numbers with extensions. By default, support for extensions is disabled.",
	)

	usersMfaEnabledDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that, when set to `true`, will enable MFA by default for new users.",
	)

	resp.Schema = schema.Schema{

		Version: 1,

		// This description is used by the documentation generator and the language server.
		Description: "Resource to manage the MFA settings for a PingOne environment.",

		Attributes: map[string]schema.Attribute{
			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment to manage MFA settings for."),
			),

			"lockout": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that contains information about the MFA policy lockout settings.").Description,
				Optional:    true,

				Attributes: map[string]schema.Attribute{
					"failure_count": schema.Int32Attribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("An integer that defines the maximum number of incorrect authentication attempts before the account is locked.").Description,
						Required:    true,

						Validators: []validator.Int32{
							int32validator.AtLeast(0),
						},
					},

					"duration_seconds": schema.Int32Attribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("An integer that defines the number of seconds to keep the account in a locked state").Description,
						Optional:    true,

						Validators: []validator.Int32{
							int32validator.AtLeast(0),
						},
					},
				},
			},

			"pairing": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that contains information about the MFA policy device pairing settings.").Description,
				Required:    true,

				Attributes: map[string]schema.Attribute{
					"max_allowed_devices": schema.Int32Attribute{
						Description:         pairingMaxAllowedDevicesDescription.Description,
						MarkdownDescription: pairingMaxAllowedDevicesDescription.MarkdownDescription,
						Optional:            true,
						Computed:            true,

						Default: int32default.StaticInt32(maxAllowedDevicesDefault),

						Validators: []validator.Int32{
							int32validator.Between(maxAllowedDevicesMin, maxAllowedDevicesMax),
						},
					},

					"pairing_key_format": schema.StringAttribute{
						Description:         pairingPairingKeyFormatDescription.Description,
						MarkdownDescription: pairingPairingKeyFormatDescription.MarkdownDescription,
						Required:            true,

						Validators: []validator.String{
							stringvalidator.OneOf(utils.EnumSliceToStringSlice(mfa.AllowedEnumMFASettingsPairingKeyFormatEnumValues)...),
						},
					},
				},
			},

			"phone_extensions": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that contains settings for phone extension support.").Description,
				Optional:    true,
				Computed:    true,

				Default: objectdefault.StaticValue(types.ObjectValueMust(
					MFASettingsPhoneExtensionsTFObjectTypes,
					map[string]attr.Value{
						"enabled": types.BoolValue(false),
					},
				)),

				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						Description:         phoneExtensionsEnabledDescription.Description,
						MarkdownDescription: phoneExtensionsEnabledDescription.MarkdownDescription,
						Required:            true,
					},
				},
			},

			"users": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that contains information about the default settings for new users.").Description,
				Optional:    true,
				Computed:    true,

				Default: objectdefault.StaticValue(types.ObjectValueMust(
					MFASettingsUsersTFObjectTypes,
					map[string]attr.Value{
						"mfa_enabled": types.BoolValue(true),
					},
				)),

				Attributes: map[string]schema.Attribute{
					"mfa_enabled": schema.BoolAttribute{
						Description:         usersMfaEnabledDescription.Description,
						MarkdownDescription: usersMfaEnabledDescription.MarkdownDescription,
						Required:            true,
					},
				},
			},
		},
	}
}

func (r *MFASettingsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *MFASettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state mFASettingsResourceModelV1

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
	mFASettings, d := plan.expand(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *mfa.MFASettings
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.MFAAPIClient.MFASettingsApi.UpdateMFASettings(ctx, plan.EnvironmentId.ValueString()).MFASettings(*mFASettings).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdateMFASettings",
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

func (r *MFASettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *mFASettingsResourceModelV1

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
	var response *mfa.MFASettings
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.MFAAPIClient.MFASettingsApi.ReadMFASettings(ctx, data.EnvironmentId.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"ReadMFASettings",
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

func (r *MFASettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state mFASettingsResourceModelV1

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
	mFASettings, d := plan.expand(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *mfa.MFASettings
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.MFAAPIClient.MFASettingsApi.UpdateMFASettings(ctx, plan.EnvironmentId.ValueString()).MFASettings(*mFASettings).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdateMFASettings",
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

func (r *MFASettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *mFASettingsResourceModelV1

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
			f0, fR, fErr := r.Client.MFAAPIClient.MFASettingsApi.ResetMFASettings(ctx, data.EnvironmentId.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), f0, fR, fErr)
		},
		"ResetMFASettings",
		framework.CustomErrorResourceNotFoundWarning,
		nil,
		nil,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *MFASettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	idComponents := []framework.ImportComponent{
		{
			Label:     "environment_id",
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
			pathKey = "environment_id"
		}

		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root(pathKey), attributes[idComponent.Label])...)
	}
}

func (p *mFASettingsResourceModelV1) expand(ctx context.Context) (*mfa.MFASettings, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Pairing
	var pairingPlan mFASettingsPairingResourceModelV1
	diags.Append(p.Pairing.As(ctx, &pairingPlan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)
	if diags.HasError() {
		return nil, diags
	}
	pairing := mfa.NewMFASettingsPairing(
		pairingPlan.MaxAllowedDevices.ValueInt32(),
		mfa.EnumMFASettingsPairingKeyFormat(pairingPlan.PairingKeyFormat.ValueString()),
	)

	// Main object
	data := mfa.NewMFASettings(
		*pairing,
	)

	// Lockout
	if !p.Lockout.IsNull() && !p.Lockout.IsUnknown() {
		var lockoutPlan mFASettingsLockoutResourceModelV1
		diags.Append(p.Lockout.As(ctx, &lockoutPlan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}
		lockout := mfa.NewMFASettingsLockout(
			lockoutPlan.FailureCount.ValueInt32(),
		)

		if !lockoutPlan.DurationSeconds.IsNull() && !lockoutPlan.DurationSeconds.IsUnknown() {
			lockout.SetDurationSeconds(lockoutPlan.DurationSeconds.ValueInt32())
		}

		data.SetLockout(*lockout)
	}

	// Phone Extensions
	if !p.PhoneExtensions.IsNull() && !p.PhoneExtensions.IsUnknown() {
		var phoneExtensionsPlan mFASettingsPhoneExtensionsResourceModelV1
		diags.Append(p.PhoneExtensions.As(ctx, &phoneExtensionsPlan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}
		phoneExtensions := mfa.NewMFASettingsPhoneExtensions()

		if !phoneExtensionsPlan.Enabled.IsNull() && !phoneExtensionsPlan.Enabled.IsUnknown() {
			phoneExtensions.SetEnabled(phoneExtensionsPlan.Enabled.ValueBool())
		}

		data.SetPhoneExtensions(*phoneExtensions)
	}

	// Users
	if !p.Users.IsNull() && !p.Users.IsUnknown() {
		var usersPlan mFASettingsUsersResourceModelV1
		diags.Append(p.Users.As(ctx, &usersPlan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}
		users := mfa.NewMFASettingsUsers()

		if !usersPlan.MFAEnabled.IsNull() && !usersPlan.MFAEnabled.IsUnknown() {
			users.SetMfaEnabled(usersPlan.MFAEnabled.ValueBool())
		}

		data.SetUsers(*users)
	}

	return data, diags
}

func (p *mFASettingsResourceModelV1) toState(apiObject *mfa.MFASettings) diag.Diagnostics {
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

func toStateLockout(apiObject *mfa.MFASettingsLockout, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(MFASettingsLockoutTFObjectTypes), nil
	}

	o := map[string]attr.Value{
		"failure_count":    framework.Int32OkToTF(apiObject.GetFailureCountOk()),
		"duration_seconds": framework.Int32OkToTF(apiObject.GetDurationSecondsOk()),
	}

	objValue, d := types.ObjectValue(MFASettingsLockoutTFObjectTypes, o)
	diags.Append(d...)

	return objValue, diags
}

func toStatePairing(apiObject *mfa.MFASettingsPairing, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(MFASettingsPairingTFObjectTypes), nil
	}

	o := map[string]attr.Value{
		"max_allowed_devices": framework.Int32OkToTF(apiObject.GetMaxAllowedDevicesOk()),
		"pairing_key_format":  framework.EnumOkToTF(apiObject.GetPairingKeyFormatOk()),
	}

	objValue, d := types.ObjectValue(MFASettingsPairingTFObjectTypes, o)
	diags.Append(d...)

	return objValue, diags
}

func toStatePhoneExtensions(apiObject *mfa.MFASettingsPhoneExtensions, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(MFASettingsPhoneExtensionsTFObjectTypes), nil
	}

	o := map[string]attr.Value{
		"enabled": framework.BoolOkToTF(apiObject.GetEnabledOk()),
	}

	objValue, d := types.ObjectValue(MFASettingsPhoneExtensionsTFObjectTypes, o)
	diags.Append(d...)

	return objValue, diags
}

func toStateUsers(apiObject *mfa.MFASettingsUsers, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(MFASettingsUsersTFObjectTypes), nil
	}

	o := map[string]attr.Value{
		"mfa_enabled": framework.BoolOkToTF(apiObject.GetMfaEnabledOk()),
	}

	objValue, d := types.ObjectValue(MFASettingsUsersTFObjectTypes, o)
	diags.Append(d...)

	return objValue, diags
}
