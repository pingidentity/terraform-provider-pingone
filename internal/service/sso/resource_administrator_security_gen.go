// Copyright © 2025 Ping Identity Corporation
// Code generated by ping-terraform-plugin-framework-generator

package sso

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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

var (
	_ resource.Resource                = &administratorSecurityResource{}
	_ resource.ResourceWithConfigure   = &administratorSecurityResource{}
	_ resource.ResourceWithImportState = &administratorSecurityResource{}
)

func NewAdministratorSecurityResource() resource.Resource {
	return &administratorSecurityResource{}
}

type administratorSecurityResource serviceClientType

func (r *administratorSecurityResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_administrator_security"
}

func (r *administratorSecurityResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

type administratorSecurityResourceModel struct {
	AllowedMethods       types.Object                 `tfsdk:"allowed_methods"`
	AuthenticationMethod types.String                 `tfsdk:"authentication_method"`
	EnvironmentId        pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	HasFido2Capabilities types.Bool                   `tfsdk:"has_fido2_capabilities"`
	Id                   pingonetypes.ResourceIDValue `tfsdk:"id"`
	IsPingIdinBom        types.Bool                   `tfsdk:"is_pingid_in_bom"`
	MfaStatus            types.String                 `tfsdk:"mfa_status"`
	IdentityProvider     types.Object                 `tfsdk:"identity_provider"`
	Recovery             types.Bool                   `tfsdk:"recovery"`
}

func (r *administratorSecurityResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Resource to create and manage environment administrator sign-on settings.",
		Attributes: map[string]schema.Attribute{
			"allowed_methods": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"email": schema.StringAttribute{
						Required:    true,
						Description: "Indicates whether to enable email for sign-on.",
						Validators: []validator.String{
							stringvalidator.OneOf(
								`{"enabled":true}`,
								`{"enabled":false}`,
							),
						},
					},
					"fido2": schema.StringAttribute{
						Required:    true,
						Description: "Indicates whether to enable FIDO2 for sign-on.",
						Validators: []validator.String{
							stringvalidator.OneOf(
								`{"enabled":true}`,
								`{"enabled":false}`,
							),
						},
					},
					"totp": schema.StringAttribute{
						Required:    true,
						Description: "Indicates whether to enable TOTP for sign-on.",
						Validators: []validator.String{
							stringvalidator.OneOf(
								`{"enabled":true}`,
								`{"enabled":false}`,
							),
						},
					},
				},
				Optional: true,
				Computed: true,
				Default: objectdefault.StaticValue(types.ObjectValueMust(map[string]attr.Type{
					"email": types.StringType,
					"fido2": types.StringType,
					"totp":  types.StringType,
				}, map[string]attr.Value{
					"email": types.StringValue(`{"enabled":true}`),
					"fido2": types.StringValue(`{"enabled":true}`),
					"totp":  types.StringValue(`{"enabled":true}`),
				})),
				Description:         "Indicates the methods to enable or disable for admin sign-on. Required properties are \"TOTP\" (temporary one-time password), \"FIDO2\", or \"EMAIL\".",
				MarkdownDescription: "Indicates the methods to enable or disable for admin sign-on. Required properties are `TOTP` (temporary one-time password), `FIDO2`, or `EMAIL`.",
			},
			"authentication_method": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Indicates whether to use PingOne MFA, an external IdP, or a combination of both for admin sign-on. Possible values are \"PINGONE\", \"EXTERNAL\", or \"HYBRID\". The default is \"PINGONE\".",
				MarkdownDescription: "Indicates whether to use PingOne MFA, an external IdP, or a combination of both for admin sign-on. Possible values are `PINGONE`, `EXTERNAL`, or `HYBRID`. The default is `PINGONE`.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"PINGONE",
						"EXTERNAL",
						"HYBRID",
					),
				},
				Default: stringdefault.StaticString("PINGONE"),
			},
			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment to create and manage the administrator_security in."),
			),
			"has_fido2_capabilities": schema.BoolAttribute{
				Computed:    true,
				Description: "Indicates whether the environment supports FIDO2 passkeys for MFA.",
			},
			"id": framework.Attr_ID(),
			"is_pingid_in_bom": schema.BoolAttribute{
				Computed:    true,
				Description: "Indicates whether the environment supports FIDO2 passkeys for MFA.",
			},
			"mfa_status": schema.StringAttribute{
				Required:            true,
				Description:         "This property must be set to \"ENFORCE\" as MFA is required for administrator sign-ons. This property applies only to the specified environment.",
				MarkdownDescription: "This property must be set to `ENFORCE` as MFA is required for administrator sign-ons. This property applies only to the specified environment.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"ENFORCE",
					),
				},
			},
			"identity_provider": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Required:    true,
						Description: "The UUID of the external IdP, if applicable.",
						CustomType:  pingonetypes.ResourceIDType{},
					},
				},
				Optional:            true,
				Description:         "The external IdP, if applicable. Required when the authentication_method is set to \"EXTERNAL\" or \"HYBRID\", otherwise should not be set.",
				MarkdownDescription: "The external IdP, if applicable. Required when the authentication_method is set to `EXTERNAL` or `HYBRID`, otherwise should not be set.",
			},
			"recovery": schema.BoolAttribute{
				Required:    true,
				Description: "Indicates whether to allow account recovery within the admin policy.",
			},
		},
	}
}

func (model *administratorSecurityResourceModel) buildClientStruct() (*management.AdministratorSecurity, diag.Diagnostics) {
	result := &management.AdministratorSecurity{}
	var respDiags diag.Diagnostics
	// allowed_methods
	if !model.AllowedMethods.IsNull() {
		allowedMethodsValue := &management.AdministratorSecurityAllowedMethods{}
		allowedMethodsAttrs := model.AllowedMethods.Attributes()
		allowedMethodsValue.EMAIL = allowedMethodsAttrs["email"].(types.String).ValueString()
		allowedMethodsValue.FIDO2 = allowedMethodsAttrs["fido2"].(types.String).ValueString()
		allowedMethodsValue.TOTP = allowedMethodsAttrs["totp"].(types.String).ValueString()
		result.AllowedMethods = allowedMethodsValue
	}

	// authentication_method
	if !model.AuthenticationMethod.IsNull() {
		authenticationMethodValue, err := management.NewEnumAdministratorSecurityAuthenticationMethodFromValue(model.AuthenticationMethod.ValueString())
		if err != nil {
			respDiags.AddAttributeError(
				path.Root("authentication_method"),
				"Provided value is not valid",
				fmt.Sprintf("The value provided for authentication_method is not valid: %s", err.Error()),
			)
		} else {
			result.AuthenticationMethod = *authenticationMethodValue
		}
	}

	// mfa_status
	mfaStatusValue, err := management.NewEnumAdministratorSecurityMfaStatusFromValue(model.MfaStatus.ValueString())
	if err != nil {
		respDiags.AddAttributeError(
			path.Root("mfa_status"),
			"Provided value is not valid",
			fmt.Sprintf("The value provided for mfa_status is not valid: %s", err.Error()),
		)
	} else {
		result.MfaStatus = *mfaStatusValue
	}

	// identity_provider
	if !model.IdentityProvider.IsNull() {
		providerValue := &management.AdministratorSecurityProvider{}
		providerAttrs := model.IdentityProvider.Attributes()
		providerValue.Id = providerAttrs["id"].(pingonetypes.ResourceIDValue).ValueString()
		result.Provider = providerValue
	}

	// recovery
	result.Recovery = model.Recovery.ValueBool()
	return result, respDiags
}

// Build a default client struct to reset the resource to its default state
// If necessary, update this function to set any other values that should be present in the default state of the resource
func (model *administratorSecurityResource) buildDefaultClientStruct() *management.AdministratorSecurity {
	result := &management.AdministratorSecurity{}
	result.AuthenticationMethod = management.EnumAdministratorSecurityAuthenticationMethod("PINGONE")
	result.MfaStatus = management.EnumAdministratorSecurityMfaStatus("ENFORCE")
	result.Recovery = true
	return result
}

func (state *administratorSecurityResourceModel) readClientResponse(response *management.AdministratorSecurity) diag.Diagnostics {
	var respDiags, diags diag.Diagnostics
	// allowed_methods
	allowedMethodsAttrTypes := map[string]attr.Type{
		"email": types.StringType,
		"fido2": types.StringType,
		"totp":  types.StringType,
	}
	var allowedMethodsValue types.Object
	if response.AllowedMethods == nil {
		allowedMethodsValue = types.ObjectNull(allowedMethodsAttrTypes)
	} else {
		allowedMethodsValue, diags = types.ObjectValue(allowedMethodsAttrTypes, map[string]attr.Value{
			"email": types.StringValue(response.AllowedMethods.EMAIL),
			"fido2": types.StringValue(response.AllowedMethods.FIDO2),
			"totp":  types.StringValue(response.AllowedMethods.TOTP),
		})
		respDiags.Append(diags...)
	}
	state.AllowedMethods = allowedMethodsValue
	// authentication_method
	authenticationMethodValue := types.StringValue(string(response.AuthenticationMethod))
	state.AuthenticationMethod = authenticationMethodValue
	// has_fido2_capabilities
	state.HasFido2Capabilities = types.BoolPointerValue(response.HasFido2Capabilities)
	// id
	state.Id = framework.PingOneResourceIDToTF(*response.GetEnvironment().Id)
	// is_pingid_in_bom
	state.IsPingIdinBom = types.BoolPointerValue(response.IsPingIDInBOM)
	// mfa_status
	mfaStatusValue := types.StringValue(string(response.MfaStatus))
	state.MfaStatus = mfaStatusValue
	// identity_provider
	providerAttrTypes := map[string]attr.Type{
		"id": types.StringType,
	}
	var providerValue types.Object
	if response.Provider == nil {
		providerValue = types.ObjectNull(providerAttrTypes)
	} else {
		providerValue, diags = types.ObjectValue(providerAttrTypes, map[string]attr.Value{
			"id": types.StringValue(response.Provider.Id),
		})
		respDiags.Append(diags...)
	}
	state.IdentityProvider = providerValue
	// recovery
	state.Recovery = types.BoolValue(response.Recovery)
	return respDiags
}

func (r *administratorSecurityResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data administratorSecurityResourceModel

	if r.Client == nil || r.Client.ManagementAPIClient == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update API call logic, since this is a singleton resource
	clientData, diags := data.buildClientStruct()
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var responseData *management.AdministratorSecurity
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.AdministratorSecurityApi.UpdateAdministratorSecurity(ctx, data.EnvironmentId.ValueString()).AdministratorSecurity(*clientData).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdateAdministratorSecurity-Create",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
		&responseData,
	)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read response into the model
	resp.Diagnostics.Append(data.readClientResponse(responseData)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *administratorSecurityResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data administratorSecurityResourceModel

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

	// Read API call logic
	var responseData *management.AdministratorSecurity
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.AdministratorSecurityApi.ReadAdministratorSecurity(ctx, data.EnvironmentId.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"ReadAdministratorSecurity",
		framework.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
		&responseData,
	)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Remove from state if resource is not found
	if responseData == nil {
		resp.State.RemoveResource(ctx)
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

func (r *administratorSecurityResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data administratorSecurityResourceModel

	if r.Client == nil || r.Client.ManagementAPIClient == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update API call logic
	clientData, diags := data.buildClientStruct()
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var responseData *management.AdministratorSecurity
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.AdministratorSecurityApi.UpdateAdministratorSecurity(ctx, data.EnvironmentId.ValueString()).AdministratorSecurity(*clientData).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdateAdministratorSecurity-Update",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
		&responseData,
	)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read response into the model
	resp.Diagnostics.Append(data.readClientResponse(responseData)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *administratorSecurityResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// This resource is singleton, so it can't be deleted from the service.
	// Instead this delete method will attempt to set the resource to its default state on the service.
	var data administratorSecurityResourceModel

	if r.Client == nil || r.Client.ManagementAPIClient == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	// Read Terraform state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update API call logic to reset to default
	clientData := r.buildDefaultClientStruct()
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.AdministratorSecurityApi.UpdateAdministratorSecurity(ctx, data.EnvironmentId.ValueString()).AdministratorSecurity(*clientData).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdateAdministratorSecurity",
		framework.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
		nil,
	)...)
}

func (r *administratorSecurityResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("environment_id"), attributes["environment_id"])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), attributes["environment_id"])...)
}
