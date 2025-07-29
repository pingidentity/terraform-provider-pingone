// Copyright © 2025 Ping Identity Corporation

package davinci

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/pingidentity/pingone-go-client/pingone"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
)

// Build the PUT client struct to be used after initial creation
func (model *davinciApplicationResourceModel) buildClientStructPutAfterCreate(createResponse *pingone.DaVinciApplication) (*pingone.DaVinciApplicationReplaceRequest, diag.Diagnostics) {
	result := &pingone.DaVinciApplicationReplaceRequest{}
	var respDiags diag.Diagnostics
	// First copy over values from the create response
	result.ApiKeyEnabled = &createResponse.ApiKey.Enabled
	var grantTypes []pingone.DaVinciApplicationReplaceRequestOauthGrantTypes
	for _, grantType := range createResponse.Oauth.GrantTypes {
		grantTypes = append(grantTypes, pingone.DaVinciApplicationReplaceRequestOauthGrantTypes(grantType))
	}
	var scopes []pingone.DaVinciApplicationReplaceRequestOauthScopes
	for _, scope := range createResponse.Oauth.Scopes {
		scopes = append(scopes, pingone.DaVinciApplicationReplaceRequestOauthScopes(scope))
	}
	result.Oauth = &pingone.DaVinciApplicationReplaceRequestOauth{
		EnforceSignedRequestOpenid: createResponse.Oauth.EnforceSignedRequestOpenid,
		GrantTypes:                 grantTypes,
		LogoutUris:                 createResponse.Oauth.LogoutUris,
		RedirectUris:               createResponse.Oauth.RedirectUris,
		Scopes:                     scopes,
		SpJwksOpenid:               createResponse.Oauth.SpJwksOpenid,
		SpjwksUrl:                  createResponse.Oauth.SpjwksUrl,
	}

	// Then overwrite with anything specified in the plan model
	result.Name = model.Name.ValueString()
	if !model.ApiKey.IsNull() && !model.ApiKey.IsUnknown() {
		if !model.ApiKey.Attributes()["enabled"].IsNull() && !model.ApiKey.Attributes()["enabled"].IsUnknown() {
			result.ApiKeyEnabled = model.ApiKey.Attributes()["enabled"].(types.Bool).ValueBoolPointer()
		}
	}
	if !model.Oauth.IsNull() && !model.Oauth.IsUnknown() {
		oauthAttrs := model.Oauth.Attributes()
		if !oauthAttrs["enforce_signed_request_openid"].IsNull() && !oauthAttrs["enforce_signed_request_openid"].IsUnknown() {
			result.Oauth.EnforceSignedRequestOpenid = oauthAttrs["enforce_signed_request_openid"].(types.Bool).ValueBoolPointer()
		}
		if !oauthAttrs["grant_types"].IsNull() && !oauthAttrs["grant_types"].IsUnknown() {
			result.Oauth.GrantTypes = []pingone.DaVinciApplicationReplaceRequestOauthGrantTypes{}
			for _, grantTypesElement := range oauthAttrs["grant_types"].(types.Set).Elements() {
				var grantTypesValue pingone.DaVinciApplicationReplaceRequestOauthGrantTypes
				grantTypesEnumValue, err := pingone.NewDaVinciApplicationReplaceRequestOauthGrantTypesFromValue(grantTypesElement.(types.String).ValueString())
				if err != nil {
					respDiags.AddAttributeError(
						path.Root("grant_types"),
						"Provided value is not valid",
						fmt.Sprintf("The value provided for grant_types is not valid: %s", err.Error()),
					)
				} else {
					grantTypesValue = *grantTypesEnumValue
				}
				result.Oauth.GrantTypes = append(result.Oauth.GrantTypes, grantTypesValue)
			}
		}
		if !oauthAttrs["logout_uris"].IsNull() && !oauthAttrs["logout_uris"].IsUnknown() {
			result.Oauth.LogoutUris = []string{}
			for _, logoutUrisElement := range oauthAttrs["logout_uris"].(types.Set).Elements() {
				result.Oauth.LogoutUris = append(result.Oauth.LogoutUris, logoutUrisElement.(types.String).ValueString())
			}
		}
		if !oauthAttrs["redirect_uris"].IsNull() && !oauthAttrs["redirect_uris"].IsUnknown() {
			result.Oauth.RedirectUris = []string{}
			for _, redirectUrisElement := range oauthAttrs["redirect_uris"].(types.Set).Elements() {
				result.Oauth.RedirectUris = append(result.Oauth.RedirectUris, redirectUrisElement.(types.String).ValueString())
			}
		}
		if !oauthAttrs["scopes"].IsNull() && !oauthAttrs["scopes"].IsUnknown() {
			result.Oauth.Scopes = []pingone.DaVinciApplicationReplaceRequestOauthScopes{}
			for _, scopesElement := range oauthAttrs["scopes"].(types.Set).Elements() {
				var scopesValue pingone.DaVinciApplicationReplaceRequestOauthScopes
				scopesEnumValue, err := pingone.NewDaVinciApplicationReplaceRequestOauthScopesFromValue(scopesElement.(types.String).ValueString())
				if err != nil {
					respDiags.AddAttributeError(
						path.Root("scopes"),
						"Provided value is not valid",
						fmt.Sprintf("The value provided for scopes is not valid: %s", err.Error()),
					)
				} else {
					scopesValue = *scopesEnumValue
				}
				result.Oauth.Scopes = append(result.Oauth.Scopes, scopesValue)
			}
		}
		if !oauthAttrs["sp_jwks_openid"].IsNull() && !oauthAttrs["sp_jwks_openid"].IsUnknown() {
			result.Oauth.SpJwksOpenid = oauthAttrs["sp_jwks_openid"].(types.String).ValueStringPointer()
		}
		if !oauthAttrs["sp_jwks_url"].IsNull() && !oauthAttrs["sp_jwks_url"].IsUnknown() {
			result.Oauth.SpjwksUrl = oauthAttrs["sp_jwks_url"].(types.String).ValueStringPointer()
		}
	}

	return result, respDiags
}

// Application Creates have to run a POST followed by a PUT to set fields other than name.
func (r *davinciApplicationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data davinciApplicationResourceModel

	if r.Client == nil {
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

	// Create API call logic
	clientData, diags := data.buildClientStructPost()
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	environmentIdUuid, err := uuid.Parse(data.EnvironmentId.ValueString())
	if err != nil {
		resp.Diagnostics.AddAttributeError(
			path.Root("environment_id"),
			"Attribute Validation Error",
			fmt.Sprintf("The value '%s' for attribute '%s' is not a valid UUID: %s", data.EnvironmentId.ValueString(), "EnvironmentId", err.Error()),
		)
		return
	}
	var createResponseData *pingone.DaVinciApplication
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.DaVinciApplicationApi.CreateDavinciApplication(ctx, environmentIdUuid).DaVinciApplicationCreateRequest(*clientData).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"CreateDavinciApplication",
		framework.DefaultCustomError,
		&createResponseData,
	)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// In order to update fields, a second call to the PUT endpoint has to be made, because only name can be set on creation.
	// Update API call logic
	updateClientData, diags := data.buildClientStructPutAfterCreate(createResponseData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var responseData *pingone.DaVinciApplication
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.DaVinciApplicationApi.ReplaceDavinciApplicationById(ctx, environmentIdUuid, createResponseData.Id).DaVinciApplicationReplaceRequest(*updateClientData).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"ReplaceDavinciApplicationById-Create",
		framework.DefaultCustomError,
		&responseData,
	)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read update response into the model
	resp.Diagnostics.Append(data.readClientResponse(responseData)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
