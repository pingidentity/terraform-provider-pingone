// Copyright © 2025 Ping Identity Corporation

package sso

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/service"
)

type identityProviderResourceModelV0 struct {
	Id                       pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId            pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	Name                     types.String                 `tfsdk:"name"`
	Description              types.String                 `tfsdk:"description"`
	Enabled                  types.Bool                   `tfsdk:"enabled"`
	RegistrationPopulationId pingonetypes.ResourceIDValue `tfsdk:"registration_population_id"`
	LoginButtonIcon          types.List                   `tfsdk:"login_button_icon"`
	Icon                     types.List                   `tfsdk:"icon"`
	Facebook                 types.List                   `tfsdk:"facebook"`
	Google                   types.List                   `tfsdk:"google"`
	LinkedIn                 types.List                   `tfsdk:"linkedin"`
	Yahoo                    types.List                   `tfsdk:"yahoo"`
	Amazon                   types.List                   `tfsdk:"amazon"`
	Twitter                  types.List                   `tfsdk:"twitter"`
	Apple                    types.List                   `tfsdk:"apple"`
	Paypal                   types.List                   `tfsdk:"paypal"`
	Microsoft                types.List                   `tfsdk:"microsoft"`
	Github                   types.List                   `tfsdk:"github"`
	OpenIDConnect            types.List                   `tfsdk:"openid_connect"`
	Saml                     types.List                   `tfsdk:"saml"`
}

type identityProviderClientIdClientSecretResourceModelV0 identityProviderClientIdClientSecretResourceModelV1

type identityProviderFacebookResourceModelV0 identityProviderFacebookResourceModelV1

type identityProviderAppleResourceModelV0 identityProviderAppleResourceModelV1

type identityProviderPaypalResourceModelV0 identityProviderPaypalResourceModelV1

type identityProviderOIDCResourceModelV0 struct {
	AuthorizationEndpoint   types.String `tfsdk:"authorization_endpoint"`
	ClientId                types.String `tfsdk:"client_id"`
	ClientSecret            types.String `tfsdk:"client_secret"`
	DiscoveryEndpoint       types.String `tfsdk:"discovery_endpoint"`
	Issuer                  types.String `tfsdk:"issuer"`
	JwksEndpoint            types.String `tfsdk:"jwks_endpoint"`
	Scopes                  types.Set    `tfsdk:"scopes"`
	TokenEndpoint           types.String `tfsdk:"token_endpoint"`
	TokenEndpointAuthMethod types.String `tfsdk:"token_endpoint_auth_method"`
	UserinfoEndpoint        types.String `tfsdk:"userinfo_endpoint"`
}

type identityProviderSAMLResourceModelV0 struct {
	AuthenticationRequestSigned   types.Bool                   `tfsdk:"authentication_request_signed"`
	IdpEntityId                   types.String                 `tfsdk:"idp_entity_id"`
	SpEntityId                    types.String                 `tfsdk:"sp_entity_id"`
	IdpVerificationCertificateIds types.Set                    `tfsdk:"idp_verification_certificate_ids"`
	SpSigningKeyId                pingonetypes.ResourceIDValue `tfsdk:"sp_signing_key_id"`
	SsoBinding                    types.String                 `tfsdk:"sso_binding"`
	SsoEndpoint                   types.String                 `tfsdk:"sso_endpoint"`
	SloBinding                    types.String                 `tfsdk:"slo_binding"`
	SloEndpoint                   types.String                 `tfsdk:"slo_endpoint"`
	SloResponseEndpoint           types.String                 `tfsdk:"slo_response_endpoint"`
	SloWindow                     types.Int32                  `tfsdk:"slo_window"`
}

func (r *IdentityProviderResource) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{
		// State upgrade implementation from 0 (prior state version) to 1 (Schema.Version)
		0: {
			PriorSchema: &schema.Schema{
				Attributes: map[string]schema.Attribute{
					"id": framework.Attr_ID(),

					"environment_id": framework.Attr_LinkID(
						framework.SchemaAttributeDescriptionFromMarkdown(""),
					),

					"name": schema.StringAttribute{
						Required: true,
					},

					"description": schema.StringAttribute{
						Optional: true,
					},

					"enabled": schema.BoolAttribute{
						Optional: true,
						Computed: true,

						Default: booldefault.StaticBool(false),
					},

					"registration_population_id": schema.StringAttribute{
						Optional: true,

						CustomType: pingonetypes.ResourceIDType{},
					},
				},

				Blocks: map[string]schema.Block{
					"login_button_icon": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									Required: true,

									CustomType: pingonetypes.ResourceIDType{},
								},

								"href": schema.StringAttribute{
									Required: true,
								},
							},
						},
					},

					"icon": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									Required: true,

									CustomType: pingonetypes.ResourceIDType{},
								},

								"href": schema.StringAttribute{
									Required: true,
								},
							},
						},
					},

					// The providers
					"facebook": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"app_id": schema.StringAttribute{
									Required: true,
								},

								"app_secret": schema.StringAttribute{
									Required:  true,
									Sensitive: true,
								},
							},
						},
					},

					"google": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"client_id": schema.StringAttribute{
									Required: true,
								},

								"client_secret": schema.StringAttribute{
									Required:  true,
									Sensitive: true,
								},
							},
						},
					},

					"linkedin": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"client_id": schema.StringAttribute{
									Required: true,
								},

								"client_secret": schema.StringAttribute{
									Required:  true,
									Sensitive: true,
								},
							},
						},
					},

					"yahoo": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"client_id": schema.StringAttribute{
									Required: true,
								},

								"client_secret": schema.StringAttribute{
									Required:  true,
									Sensitive: true,
								},
							},
						},
					},

					"amazon": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"client_id": schema.StringAttribute{
									Required: true,
								},

								"client_secret": schema.StringAttribute{
									Required:  true,
									Sensitive: true,
								},
							},
						},
					},

					"twitter": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"client_id": schema.StringAttribute{
									Required: true,
								},

								"client_secret": schema.StringAttribute{
									Required:  true,
									Sensitive: true,
								},
							},
						},
					},

					"apple": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"client_id": schema.StringAttribute{
									Required: true,
								},

								"client_secret_signing_key": schema.StringAttribute{
									Required:  true,
									Sensitive: true,
								},

								"key_id": schema.StringAttribute{
									Required: true,
								},

								"team_id": schema.StringAttribute{
									Required: true,
								},
							},
						},
					},

					"paypal": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"client_id": schema.StringAttribute{
									Required: true,
								},

								"client_secret": schema.StringAttribute{
									Required:  true,
									Sensitive: true,
								},

								"client_environment": schema.StringAttribute{
									Required: true,
								},
							},
						},
					},

					"microsoft": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"client_id": schema.StringAttribute{
									Required: true,
								},

								"client_secret": schema.StringAttribute{
									Required:  true,
									Sensitive: true,
								},
							},
						},
					},

					"github": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"client_id": schema.StringAttribute{
									Required: true,
								},

								"client_secret": schema.StringAttribute{
									Required:  true,
									Sensitive: true,
								},
							},
						},
					},

					"openid_connect": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"authorization_endpoint": schema.StringAttribute{
									Required: true,
								},

								"client_id": schema.StringAttribute{
									Required: true,
								},

								"client_secret": schema.StringAttribute{
									Required:  true,
									Sensitive: true,
								},

								"discovery_endpoint": schema.StringAttribute{
									Optional: true,
								},

								"issuer": schema.StringAttribute{
									Required: true,
								},

								"jwks_endpoint": schema.StringAttribute{
									Required: true,
								},

								"scopes": schema.SetAttribute{
									Required: true,

									ElementType: types.StringType,
								},

								"token_endpoint": schema.StringAttribute{
									Required: true,
								},

								"token_endpoint_auth_method": schema.StringAttribute{
									Optional: true,
									Computed: true,

									Default: stringdefault.StaticString(string(management.ENUMIDENTITYPROVIDEROIDCTOKENAUTHMETHOD_CLIENT_SECRET_BASIC)),
								},

								"userinfo_endpoint": schema.StringAttribute{
									Optional: true,
								},
							},
						},
					},

					"saml": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"authentication_request_signed": schema.BoolAttribute{
									Optional: true,
									Computed: true,

									Default: booldefault.StaticBool(false),
								},

								"idp_entity_id": schema.StringAttribute{
									Required: true,
								},

								"sp_entity_id": schema.StringAttribute{
									Required: true,
								},

								"idp_verification_certificate_ids": schema.SetAttribute{
									Required: true,

									ElementType: pingonetypes.ResourceIDType{},
								},

								"sp_signing_key_id": schema.StringAttribute{
									Optional: true,

									CustomType: pingonetypes.ResourceIDType{},
								},

								"sso_binding": schema.StringAttribute{
									Required: true,
								},

								"sso_endpoint": schema.StringAttribute{
									Required: true,
								},

								"slo_binding": schema.StringAttribute{
									Optional: true,
									Computed: true,

									Default: stringdefault.StaticString(string(management.ENUMIDENTITYPROVIDERSAMLSLOBINDING_POST)),
								},

								"slo_endpoint": schema.StringAttribute{
									Optional: true,
								},

								"slo_response_endpoint": schema.StringAttribute{
									Optional: true,
								},

								"slo_window": schema.Int32Attribute{
									Optional: true,
								},
							},
						},
					},
				},
			},
			StateUpgrader: func(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
				var d diag.Diagnostics
				var priorStateData identityProviderResourceModelV0

				resp.Diagnostics.Append(req.State.Get(ctx, &priorStateData)...)

				if resp.Diagnostics.HasError() {
					return
				}

				icon, d := priorStateData.schemaUpgradeIconV0toV1(ctx)
				resp.Diagnostics.Append(d...)

				loginButtonIcon, d := priorStateData.schemaUpgradeLoginButtonIconV0toV1(ctx)
				resp.Diagnostics.Append(d...)

				facebook, d := priorStateData.schemaUpgradeFacebookV0toV1(ctx)
				resp.Diagnostics.Append(d...)

				google, d := priorStateData.schemaUpgradeGoogleV0toV1(ctx)
				resp.Diagnostics.Append(d...)

				linkedIn, d := priorStateData.schemaUpgradeLinkedInV0toV1(ctx)
				resp.Diagnostics.Append(d...)

				yahoo, d := priorStateData.schemaUpgradeYahooV0toV1(ctx)
				resp.Diagnostics.Append(d...)

				amazon, d := priorStateData.schemaUpgradeAmazonV0toV1(ctx)
				resp.Diagnostics.Append(d...)

				twitter, d := priorStateData.schemaUpgradeTwitterV0toV1(ctx)
				resp.Diagnostics.Append(d...)

				apple, d := priorStateData.schemaUpgradeAppleV0toV1(ctx)
				resp.Diagnostics.Append(d...)

				paypal, d := priorStateData.schemaUpgradePaypalV0toV1(ctx)
				resp.Diagnostics.Append(d...)

				microsoft, d := priorStateData.schemaUpgradeMicrosoftV0toV1(ctx)
				resp.Diagnostics.Append(d...)

				github, d := priorStateData.schemaUpgradeGithubV0toV1(ctx)
				resp.Diagnostics.Append(d...)

				openIDConnect, d := priorStateData.schemaUpgradeOpenIDConnectV0toV1(ctx)
				resp.Diagnostics.Append(d...)

				sAMLOptions, d := priorStateData.schemaUpgradeSAMLV0toV1(ctx)
				resp.Diagnostics.Append(d...)

				upgradedStateData := identityProviderResourceModelV1{
					Id:                       priorStateData.Id,
					EnvironmentId:            priorStateData.EnvironmentId,
					Name:                     priorStateData.Name,
					Description:              priorStateData.Description,
					Enabled:                  priorStateData.Enabled,
					RegistrationPopulationId: priorStateData.RegistrationPopulationId,
					LoginButtonIcon:          loginButtonIcon,
					Icon:                     icon,
					Facebook:                 facebook,
					Google:                   google,
					LinkedIn:                 linkedIn,
					Yahoo:                    yahoo,
					Amazon:                   amazon,
					Twitter:                  twitter,
					Apple:                    apple,
					Paypal:                   paypal,
					Microsoft:                microsoft,
					Github:                   github,
					OpenIDConnect:            openIDConnect,
					Saml:                     sAMLOptions,
				}

				resp.Diagnostics.Append(resp.State.Set(ctx, upgradedStateData)...)
			},
		},
	}
}

func (p *identityProviderResourceModelV0) schemaUpgradeIconV0toV1(ctx context.Context) (types.Object, diag.Diagnostics) {
	return service.ImageListToObjectSchemaUpgrade(ctx, p.Icon)
}

func (p *identityProviderResourceModelV0) schemaUpgradeLoginButtonIconV0toV1(ctx context.Context) (types.Object, diag.Diagnostics) {
	return service.ImageListToObjectSchemaUpgrade(ctx, p.Icon)
}

func (p *identityProviderResourceModelV0) schemaUpgradeFacebookV0toV1(ctx context.Context) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributeTypes := identityProviderFacebookTFObjectTypes
	planAttribute := p.Facebook

	if planAttribute.IsNull() {
		return types.ObjectNull(attributeTypes), diags
	} else if planAttribute.IsUnknown() {
		return types.ObjectUnknown(attributeTypes), diags
	} else {
		var priorStateData []identityProviderFacebookResourceModelV0
		d := planAttribute.ElementsAs(ctx, &priorStateData, false)
		diags.Append(d...)
		if diags.HasError() {
			return types.ObjectNull(attributeTypes), diags
		}

		if len(priorStateData) == 0 {
			return types.ObjectNull(attributeTypes), diags
		}

		upgradedStateData := priorStateData[0]

		returnVar, d := types.ObjectValueFrom(ctx, attributeTypes, upgradedStateData)
		diags.Append(d...)

		return returnVar, diags
	}
}

func (p *identityProviderResourceModelV0) schemaUpgradeClientIDClientSecretV0toV1(ctx context.Context, planAttribute types.List) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributeTypes := identityProviderClientIDClientSecretTFObjectTypes

	if planAttribute.IsNull() {
		return types.ObjectNull(attributeTypes), diags
	} else if planAttribute.IsUnknown() {
		return types.ObjectUnknown(attributeTypes), diags
	} else {
		var priorStateData []identityProviderClientIdClientSecretResourceModelV0
		d := planAttribute.ElementsAs(ctx, &priorStateData, false)
		diags.Append(d...)
		if diags.HasError() {
			return types.ObjectNull(attributeTypes), diags
		}

		if len(priorStateData) == 0 {
			return types.ObjectNull(attributeTypes), diags
		}

		upgradedStateData := priorStateData[0]

		returnVar, d := types.ObjectValueFrom(ctx, attributeTypes, upgradedStateData)
		diags.Append(d...)

		return returnVar, diags
	}
}

func (p *identityProviderResourceModelV0) schemaUpgradeGoogleV0toV1(ctx context.Context) (types.Object, diag.Diagnostics) {
	return p.schemaUpgradeClientIDClientSecretV0toV1(ctx, p.Google)
}

func (p *identityProviderResourceModelV0) schemaUpgradeLinkedInV0toV1(ctx context.Context) (types.Object, diag.Diagnostics) {
	return p.schemaUpgradeClientIDClientSecretV0toV1(ctx, p.LinkedIn)
}

func (p *identityProviderResourceModelV0) schemaUpgradeYahooV0toV1(ctx context.Context) (types.Object, diag.Diagnostics) {
	return p.schemaUpgradeClientIDClientSecretV0toV1(ctx, p.Yahoo)
}

func (p *identityProviderResourceModelV0) schemaUpgradeAmazonV0toV1(ctx context.Context) (types.Object, diag.Diagnostics) {
	return p.schemaUpgradeClientIDClientSecretV0toV1(ctx, p.Amazon)
}

func (p *identityProviderResourceModelV0) schemaUpgradeTwitterV0toV1(ctx context.Context) (types.Object, diag.Diagnostics) {
	return p.schemaUpgradeClientIDClientSecretV0toV1(ctx, p.Twitter)
}

func (p *identityProviderResourceModelV0) schemaUpgradeAppleV0toV1(ctx context.Context) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributeTypes := identityProviderAppleTFObjectTypes
	planAttribute := p.Apple

	if planAttribute.IsNull() {
		return types.ObjectNull(attributeTypes), diags
	} else if planAttribute.IsUnknown() {
		return types.ObjectUnknown(attributeTypes), diags
	} else {
		var priorStateData []identityProviderAppleResourceModelV0
		d := planAttribute.ElementsAs(ctx, &priorStateData, false)
		diags.Append(d...)
		if diags.HasError() {
			return types.ObjectNull(attributeTypes), diags
		}

		if len(priorStateData) == 0 {
			return types.ObjectNull(attributeTypes), diags
		}

		upgradedStateData := priorStateData[0]

		returnVar, d := types.ObjectValueFrom(ctx, attributeTypes, upgradedStateData)
		diags.Append(d...)

		return returnVar, diags
	}
}

func (p *identityProviderResourceModelV0) schemaUpgradePaypalV0toV1(ctx context.Context) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributeTypes := identityProviderPaypalTFObjectTypes
	planAttribute := p.Paypal

	if planAttribute.IsNull() {
		return types.ObjectNull(attributeTypes), diags
	} else if planAttribute.IsUnknown() {
		return types.ObjectUnknown(attributeTypes), diags
	} else {
		var priorStateData []identityProviderPaypalResourceModelV0
		d := planAttribute.ElementsAs(ctx, &priorStateData, false)
		diags.Append(d...)
		if diags.HasError() {
			return types.ObjectNull(attributeTypes), diags
		}

		if len(priorStateData) == 0 {
			return types.ObjectNull(attributeTypes), diags
		}

		upgradedStateData := priorStateData[0]

		returnVar, d := types.ObjectValueFrom(ctx, attributeTypes, upgradedStateData)
		diags.Append(d...)

		return returnVar, diags
	}
}

func (p *identityProviderResourceModelV0) schemaUpgradeMicrosoftV0toV1(ctx context.Context) (types.Object, diag.Diagnostics) {
	return p.schemaUpgradeClientIDClientSecretV0toV1(ctx, p.Microsoft)
}

func (p *identityProviderResourceModelV0) schemaUpgradeGithubV0toV1(ctx context.Context) (types.Object, diag.Diagnostics) {
	return p.schemaUpgradeClientIDClientSecretV0toV1(ctx, p.Github)
}

func (p *identityProviderResourceModelV0) schemaUpgradeOpenIDConnectV0toV1(ctx context.Context) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributeTypes := identityProviderOIDCTFObjectTypes
	planAttribute := p.OpenIDConnect

	if planAttribute.IsNull() {
		return types.ObjectNull(attributeTypes), diags
	} else if planAttribute.IsUnknown() {
		return types.ObjectUnknown(attributeTypes), diags
	} else {
		var priorStateData []identityProviderOIDCResourceModelV0
		d := planAttribute.ElementsAs(ctx, &priorStateData, false)
		diags.Append(d...)
		if diags.HasError() {
			return types.ObjectNull(attributeTypes), diags
		}

		if len(priorStateData) == 0 {
			return types.ObjectNull(attributeTypes), diags
		}

		upgradedStateData := identityProviderOIDCResourceModelV1{
			AuthorizationEndpoint:   priorStateData[0].AuthorizationEndpoint,
			ClientId:                priorStateData[0].ClientId,
			ClientSecret:            priorStateData[0].ClientSecret,
			DiscoveryEndpoint:       priorStateData[0].DiscoveryEndpoint,
			Issuer:                  priorStateData[0].Issuer,
			PkceMethod:              types.StringValue(string(management.ENUMIDENTITYPROVIDERPKCEMETHOD_NONE)),
			JwksEndpoint:            priorStateData[0].JwksEndpoint,
			Scopes:                  priorStateData[0].Scopes,
			TokenEndpoint:           priorStateData[0].TokenEndpoint,
			TokenEndpointAuthMethod: priorStateData[0].TokenEndpointAuthMethod,
			UserinfoEndpoint:        priorStateData[0].UserinfoEndpoint,
		}

		returnVar, d := types.ObjectValueFrom(ctx, attributeTypes, upgradedStateData)
		diags.Append(d...)

		return returnVar, diags
	}
}

func (p *identityProviderResourceModelV0) schemaUpgradeSAMLV0toV1(ctx context.Context) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributeTypes := identityProviderSAMLTFObjectTypes
	planAttribute := p.Saml

	if planAttribute.IsNull() {
		return types.ObjectNull(attributeTypes), diags
	} else if planAttribute.IsUnknown() {
		return types.ObjectUnknown(attributeTypes), diags
	} else {
		var priorStateData []identityProviderSAMLResourceModelV0
		d := planAttribute.ElementsAs(ctx, &priorStateData, false)
		diags.Append(d...)
		if diags.HasError() {
			return types.ObjectNull(attributeTypes), diags
		}

		if len(priorStateData) == 0 {
			return types.ObjectNull(attributeTypes), diags
		}

		idpVerification, d := priorStateData[0].schemaUpgradeIdpVerificationV0toV1(ctx)
		diags.Append(d...)

		spSigning, d := priorStateData[0].schemaUpgradeSpSigningV0toV1(ctx)
		diags.Append(d...)

		upgradedStateData := identityProviderSAMLResourceModelV1{
			AuthenticationRequestSigned: priorStateData[0].AuthenticationRequestSigned,
			IdpEntityId:                 priorStateData[0].IdpEntityId,
			SpEntityId:                  priorStateData[0].SpEntityId,
			IdpVerification:             idpVerification,
			SpSigning:                   spSigning,
			SsoBinding:                  priorStateData[0].SsoBinding,
			SsoEndpoint:                 priorStateData[0].SsoEndpoint,
			SloBinding:                  priorStateData[0].SloBinding,
			SloEndpoint:                 priorStateData[0].SloEndpoint,
			SloResponseEndpoint:         priorStateData[0].SloResponseEndpoint,
			SloWindow:                   priorStateData[0].SloWindow,
		}

		returnVar, d := types.ObjectValueFrom(ctx, attributeTypes, upgradedStateData)
		diags.Append(d...)

		return returnVar, diags
	}
}

func (p *identityProviderSAMLResourceModelV0) schemaUpgradeIdpVerificationV0toV1(ctx context.Context) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributeTypes := identityProviderSAMLIdPVerificationTFObjectTypes
	planAttribute := p.IdpVerificationCertificateIds

	if planAttribute.IsNull() {
		return types.ObjectNull(attributeTypes), diags
	} else if planAttribute.IsUnknown() {
		return types.ObjectUnknown(attributeTypes), diags
	} else {
		var priorStateData []types.String
		d := planAttribute.ElementsAs(ctx, &priorStateData, false)
		diags.Append(d...)
		if diags.HasError() {
			return types.ObjectNull(attributeTypes), diags
		}

		if len(priorStateData) == 0 {
			return types.ObjectNull(attributeTypes), diags
		}

		certificatesArr := make([]identityProviderSAMLResourceIdPVerificationCertificatesModelV1, len(priorStateData))

		for i, cert := range priorStateData {
			certificatesArr[i] = identityProviderSAMLResourceIdPVerificationCertificatesModelV1{
				Id: framework.PingOneResourceIDToTF(cert.ValueString()),
			}
		}

		certificates, d := types.SetValueFrom(ctx, types.ObjectType{AttrTypes: identityProviderSAMLIdPVerificationCertificateTFObjectTypes}, certificatesArr)
		diags.Append(d...)

		upgradedStateData := identityProviderSAMLResourceIdPVerificationModelV1{
			Certificates: certificates,
		}

		returnVar, d := types.ObjectValueFrom(ctx, attributeTypes, upgradedStateData)
		diags.Append(d...)

		return returnVar, diags
	}
}

func (p *identityProviderSAMLResourceModelV0) schemaUpgradeSpSigningV0toV1(ctx context.Context) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributeTypes := identityProviderSAMLSpSigningTFObjectTypes
	planAttribute := p.SpSigningKeyId

	if planAttribute.IsNull() {
		return types.ObjectNull(attributeTypes), diags
	} else if planAttribute.IsUnknown() {
		return types.ObjectUnknown(attributeTypes), diags
	} else {

		upgradedStateIdData := identityProviderSAMLResourceSpSigningKeyModelV1{
			Id: planAttribute,
		}

		key, d := types.ObjectValueFrom(ctx, identityProviderSAMLSpSigningKeyTFObjectTypes, upgradedStateIdData)
		diags.Append(d...)

		upgradedStateData := identityProviderSAMLResourceSpSigningModelV1{
			Key:       key,
			Algorithm: types.StringNull(),
		}

		returnVar, d := types.ObjectValueFrom(ctx, attributeTypes, upgradedStateData)
		diags.Append(d...)

		return returnVar, diags
	}
}
