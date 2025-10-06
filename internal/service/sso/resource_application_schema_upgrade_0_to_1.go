// Copyright © 2025 Ping Identity Corporation

package sso

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/service"
	"github.com/pingidentity/terraform-provider-pingone/internal/service/sso/helpers/beta"
)

type applicationResourceModelV0 struct {
	Id                        pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId             pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	Name                      types.String                 `tfsdk:"name"`
	Description               types.String                 `tfsdk:"description"`
	Enabled                   types.Bool                   `tfsdk:"enabled"`
	Tags                      types.Set                    `tfsdk:"tags"`
	LoginPageUrl              types.String                 `tfsdk:"login_page_url"`
	Icon                      types.List                   `tfsdk:"icon"`
	AccessControlRoleType     types.String                 `tfsdk:"access_control_role_type"`
	AccessControlGroupOptions types.List                   `tfsdk:"access_control_group_options"`
	HiddenFromAppPortal       types.Bool                   `tfsdk:"hidden_from_app_portal"`
	ExternalLinkOptions       types.List                   `tfsdk:"external_link_options"`
	OIDCOptions               types.List                   `tfsdk:"oidc_options"`
	SAMLOptions               types.List                   `tfsdk:"saml_options"`
}

type applicationAccessControlGroupOptionsResourceModelV0 applicationAccessControlGroupOptionsResourceModelV1

type applicationExternalLinkOptionsResourceModelV0 applicationExternalLinkOptionsResourceModelV1

type applicationOIDCOptionsResourceModelV0 struct {
	AdditionalRefreshTokenReplayProtectionEnabled types.Bool                   `tfsdk:"additional_refresh_token_replay_protection_enabled"`
	AllowWildcardsInRedirectUris                  types.Bool                   `tfsdk:"allow_wildcards_in_redirect_uris"`
	CertificateBasedAuthentication                types.List                   `tfsdk:"certificate_based_authentication"`
	CorsSettings                                  types.List                   `tfsdk:"cors_settings"`
	GrantTypes                                    types.Set                    `tfsdk:"grant_types"`
	HomePageUrl                                   types.String                 `tfsdk:"home_page_url"`
	InitiateLoginUri                              types.String                 `tfsdk:"initiate_login_uri"`
	Jwks                                          types.String                 `tfsdk:"jwks"`
	JwksUrl                                       types.String                 `tfsdk:"jwks_url"`
	MobileApp                                     types.List                   `tfsdk:"mobile_app"`
	ParRequirement                                types.String                 `tfsdk:"par_requirement"`
	ParTimeout                                    types.Int32                  `tfsdk:"par_timeout"`
	PKCEEnforcement                               types.String                 `tfsdk:"pkce_enforcement"`
	PostLogoutRedirectUris                        types.Set                    `tfsdk:"post_logout_redirect_uris"`
	RedirectUris                                  types.Set                    `tfsdk:"redirect_uris"`
	RefreshTokenDuration                          types.Int32                  `tfsdk:"refresh_token_duration"`
	RefreshTokenRollingDuration                   types.Int32                  `tfsdk:"refresh_token_rolling_duration"`
	RefreshTokenRollingGracePeriodDuration        types.Int32                  `tfsdk:"refresh_token_rolling_grace_period_duration"`
	RequireSignedRequestObject                    types.Bool                   `tfsdk:"require_signed_request_object"`
	ResponseTypes                                 types.Set                    `tfsdk:"response_types"`
	SupportUnsignedRequestObject                  types.Bool                   `tfsdk:"support_unsigned_request_object"`
	TargetLinkUri                                 types.String                 `tfsdk:"target_link_uri"`
	TokenEndpointAuthnMethod                      types.String                 `tfsdk:"token_endpoint_authn_method"`
	Type                                          types.String                 `tfsdk:"type"`
	ClientId                                      pingonetypes.ResourceIDValue `tfsdk:"client_id"`
	ClientSecret                                  types.String                 `tfsdk:"client_secret"`
	BundleId                                      types.String                 `tfsdk:"bundle_id"`
	PackageName                                   types.String                 `tfsdk:"package_name"`
}

type applicationCorsSettingsResourceModelV0 applicationCorsSettingsResourceModelV1

type applicationOIDCCertificateBasedAuthenticationResourceModelV0 applicationOIDCCertificateBasedAuthenticationResourceModelV1

type applicationOIDCMobileAppResourceModelV0 struct {
	BundleId               types.String `tfsdk:"bundle_id"`
	HuaweiAppId            types.String `tfsdk:"huawei_app_id"`
	HuaweiPackageName      types.String `tfsdk:"huawei_package_name"`
	IntegrityDetection     types.List   `tfsdk:"integrity_detection"`
	PackageName            types.String `tfsdk:"package_name"`
	PasscodeRefreshSeconds types.Int32  `tfsdk:"passcode_refresh_seconds"`
	UniversalAppLink       types.String `tfsdk:"universal_app_link"`
}

type applicationOIDCMobileAppIntegrityDetectionResourceModelV0 struct {
	CacheDuration     types.List `tfsdk:"cache_duration"`
	Enabled           types.Bool `tfsdk:"enabled"`
	ExcludedPlatforms types.Set  `tfsdk:"excluded_platforms"`
	GooglePlay        types.List `tfsdk:"google_play"`
}

type applicationOIDCMobileAppIntegrityDetectionCacheDurationResourceModelV0 applicationOIDCMobileAppIntegrityDetectionCacheDurationResourceModelV1

type applicationOIDCMobileAppIntegrityDetectionGooglePlayResourceModelV0 struct {
	DecryptionKey                 types.String `tfsdk:"decryption_key"`
	ServiceAccountCredentialsJson types.String `tfsdk:"service_account_credentials_json"`
	VerificationKey               types.String `tfsdk:"verification_key"`
	VerificationType              types.String `tfsdk:"verification_type"`
}

type applicationSAMLOptionsResourceModelV0 struct {
	AcsUrls                      types.Set    `tfsdk:"acs_urls"`
	AssertionDuration            types.Int32  `tfsdk:"assertion_duration"`
	AssertionSignedEnabled       types.Bool   `tfsdk:"assertion_signed_enabled"`
	CorsSettings                 types.List   `tfsdk:"cors_settings"`
	EnableRequestedAuthnContext  types.Bool   `tfsdk:"enable_requested_authn_context"`
	HomePageUrl                  types.String `tfsdk:"home_page_url"`
	IdpSigningKeyId              types.String `tfsdk:"idp_signing_key_id"`
	IdpSigningKey                types.List   `tfsdk:"idp_signing_key"`
	DefaultTargetUrl             types.String `tfsdk:"default_target_url"`
	NameIdFormat                 types.String `tfsdk:"nameid_format"`
	ResponseIsSigned             types.Bool   `tfsdk:"response_is_signed"`
	SloBinding                   types.String `tfsdk:"slo_binding"`
	SloEndpoint                  types.String `tfsdk:"slo_endpoint"`
	SloResponseEndpoint          types.String `tfsdk:"slo_response_endpoint"`
	SloWindow                    types.Int32  `tfsdk:"slo_window"`
	SpEncryption                 types.List   `tfsdk:"sp_encryption"`
	SpEntityId                   types.String `tfsdk:"sp_entity_id"`
	SpVerificationCertificateIds types.Set    `tfsdk:"sp_verification_certificate_ids"`
	SpVerification               types.List   `tfsdk:"sp_verification"`
	Type                         types.String `tfsdk:"type"`
}

type applicationSAMLOptionsIdpSigningKeyResourceModelV0 applicationOptionsIdpSigningKeyResourceModelV1

type applicationSAMLOptionsSpEncryptionResourceModelV0 struct {
	Algorithm   types.String `tfsdk:"algorithm"`
	Certificate types.List   `tfsdk:"certificate"`
}

type applicationSAMLOptionsSpEncryptionCertificateResourceModelV0 applicationSAMLOptionsSpEncryptionCertificateResourceModelV1

type applicationSAMLOptionsSpVerificationResourceModelV0 applicationSAMLOptionsSpVerificationResourceModelV1

func (r *ApplicationResource) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {

	const oidcOptionsParTimeoutDefault = 60
	const oidcOptionsRefreshTokenDurationDefault = 2592000
	const oidcOptionsRefreshTokenRollingDurationDefault = 15552000
	const oidcOptionsMobileAppPasscodeRefreshSecondsDefault = 30

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

						Default: booldefault.StaticBool(true),
					},

					"tags": schema.SetAttribute{
						Optional: true,

						ElementType: types.StringType,
					},

					"login_page_url": schema.StringAttribute{
						Optional: true,
					},

					"access_control_role_type": schema.StringAttribute{
						Optional: true,
						Computed: true,
					},

					"hidden_from_app_portal": schema.BoolAttribute{
						Optional: true,
						Computed: true,

						Default: booldefault.StaticBool(false),
					},
				},

				Blocks: map[string]schema.Block{
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

					"access_control_group_options": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"type": schema.StringAttribute{
									Required: true,
								},

								"groups": schema.SetAttribute{
									Required: true,

									ElementType: types.StringType,
								},
							},
						},
					},

					"external_link_options": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"home_page_url": schema.StringAttribute{
									Required: true,
								},
							},
						},
					},

					"oidc_options": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"type": schema.StringAttribute{
									Required: true,
								},

								"home_page_url": schema.StringAttribute{
									Optional: true,
								},

								"initiate_login_uri": schema.StringAttribute{
									Optional: true,
								},

								"jwks": schema.StringAttribute{
									Optional: true,
								},

								"jwks_url": schema.StringAttribute{
									Optional: true,
								},

								"target_link_uri": schema.StringAttribute{
									Optional: true,
								},

								"grant_types": schema.SetAttribute{
									Required: true,

									ElementType: types.StringType,
								},

								"response_types": schema.SetAttribute{
									Optional: true,

									ElementType: types.StringType,
								},

								"token_endpoint_authn_method": schema.StringAttribute{
									Required: true,
								},

								"par_requirement": schema.StringAttribute{
									Optional: true,
									Computed: true,

									Default: stringdefault.StaticString(string(management.ENUMAPPLICATIONOIDCPARREQUIREMENT_OPTIONAL)),
								},

								"par_timeout": schema.Int32Attribute{
									Optional: true,
									Computed: true,

									Default: int32default.StaticInt32(oidcOptionsParTimeoutDefault),
								},

								"pkce_enforcement": schema.StringAttribute{
									Optional: true,
									Computed: true,

									Default: stringdefault.StaticString(string(management.ENUMAPPLICATIONOIDCPKCEOPTION_OPTIONAL)),
								},

								"redirect_uris": schema.SetAttribute{
									Optional: true,

									ElementType: types.StringType,
								},

								"allow_wildcards_in_redirect_uris": schema.BoolAttribute{
									Optional: true,
									Computed: true,

									Default: booldefault.StaticBool(false),
								},

								"post_logout_redirect_uris": schema.SetAttribute{
									Optional: true,

									ElementType: types.StringType,
								},

								"refresh_token_duration": schema.Int32Attribute{
									Optional: true,
									Computed: true,

									Default: int32default.StaticInt32(oidcOptionsRefreshTokenDurationDefault),
								},

								"refresh_token_rolling_duration": schema.Int32Attribute{
									Optional: true,
									Computed: true,

									Default: int32default.StaticInt32(oidcOptionsRefreshTokenRollingDurationDefault),
								},

								"refresh_token_rolling_grace_period_duration": schema.Int32Attribute{
									Optional: true,
								},

								"additional_refresh_token_replay_protection_enabled": schema.BoolAttribute{
									Optional: true,
									Computed: true,

									Default: booldefault.StaticBool(true),
								},

								"client_id": schema.StringAttribute{
									Computed: true,

									CustomType: pingonetypes.ResourceIDType{},
								},

								"client_secret": schema.StringAttribute{
									Computed:  true,
									Sensitive: true,
								},

								"support_unsigned_request_object": schema.BoolAttribute{
									Optional: true,
									Computed: true,

									Default: booldefault.StaticBool(false),
								},

								"require_signed_request_object": schema.BoolAttribute{
									Optional: true,
									Computed: true,

									Default: booldefault.StaticBool(false),
								},

								"bundle_id": schema.StringAttribute{
									Optional: true,
									Computed: true,
								},

								"package_name": schema.StringAttribute{
									Optional: true,
									Computed: true,
								},
							},

							Blocks: map[string]schema.Block{
								"certificate_based_authentication": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"key_id": schema.StringAttribute{
												Required: true,

												CustomType: pingonetypes.ResourceIDType{},
											},
										},
									},
								},

								"cors_settings": r.resourceApplicationSchemaCorsSettingsV0(),

								"mobile_app": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"bundle_id": schema.StringAttribute{
												Optional: true,
											},

											"package_name": schema.StringAttribute{
												Optional: true,
											},

											"huawei_app_id": schema.StringAttribute{
												Optional: true,
											},

											"huawei_package_name": schema.StringAttribute{
												Optional: true,
											},

											"passcode_refresh_seconds": schema.Int32Attribute{
												Optional: true,
												Computed: true,

												Default: int32default.StaticInt32(oidcOptionsMobileAppPasscodeRefreshSecondsDefault),
											},

											"universal_app_link": schema.StringAttribute{
												Optional: true,
											},
										},

										Blocks: map[string]schema.Block{
											"integrity_detection": schema.ListNestedBlock{
												NestedObject: schema.NestedBlockObject{
													Attributes: map[string]schema.Attribute{
														"enabled": schema.BoolAttribute{
															Optional: true,
															Computed: true,

															Default: booldefault.StaticBool(false),
														},

														"excluded_platforms": schema.SetAttribute{
															Optional: true,

															ElementType: types.StringType,
														},
													},

													Blocks: map[string]schema.Block{
														"cache_duration": schema.ListNestedBlock{
															NestedObject: schema.NestedBlockObject{
																Attributes: map[string]schema.Attribute{
																	"amount": schema.Int32Attribute{
																		Required: true,
																	},

																	"units": schema.StringAttribute{
																		Optional: true,
																		Computed: true,

																		Default: stringdefault.StaticString(string(management.ENUMDURATIONUNITMINSHOURS_MINUTES)),
																	},
																},
															},
														},

														"google_play": schema.ListNestedBlock{
															NestedObject: schema.NestedBlockObject{
																Attributes: map[string]schema.Attribute{
																	"decryption_key": schema.StringAttribute{
																		Optional:  true,
																		Sensitive: true,
																	},

																	"service_account_credentials_json": schema.StringAttribute{
																		Optional:  true,
																		Sensitive: true,
																	},

																	"verification_key": schema.StringAttribute{
																		Optional:  true,
																		Sensitive: true,
																	},

																	"verification_type": schema.StringAttribute{
																		Required: true,
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},

					"saml_options": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"home_page_url": schema.StringAttribute{
									Optional: true,
								},

								"type": schema.StringAttribute{
									Optional: true,
									Computed: true,

									Default: stringdefault.StaticString(string(management.ENUMAPPLICATIONTYPE_WEB_APP)),
								},

								"acs_urls": schema.SetAttribute{
									Required: true,

									ElementType: types.StringType,
								},

								"assertion_duration": schema.Int32Attribute{
									Required: true,
								},

								"assertion_signed_enabled": schema.BoolAttribute{
									Optional: true,
									Computed: true,

									Default: booldefault.StaticBool(true),
								},

								"default_target_url": schema.StringAttribute{
									Optional: true,
								},

								"idp_signing_key_id": schema.StringAttribute{
									Optional: true,
									Computed: true,
								},

								"enable_requested_authn_context": schema.BoolAttribute{
									Optional: true,
								},

								"nameid_format": schema.StringAttribute{
									Optional: true,
								},

								"response_is_signed": schema.BoolAttribute{
									Optional: true,
									Computed: true,

									Default: booldefault.StaticBool(false),
								},

								"slo_binding": schema.StringAttribute{
									Optional: true,
									Computed: true,

									Default: stringdefault.StaticString(string(management.ENUMAPPLICATIONSAMLSLOBINDING_POST)),
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

								"sp_entity_id": schema.StringAttribute{
									Required: true,
								},

								"sp_verification_certificate_ids": schema.SetAttribute{
									Optional: true,
									Computed: true,

									ElementType: types.StringType,
								},
							},

							Blocks: map[string]schema.Block{
								"idp_signing_key": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"algorithm": schema.StringAttribute{
												Required: true,
											},

											"key_id": schema.StringAttribute{
												Required: true,

												CustomType: pingonetypes.ResourceIDType{},
											},
										},
									},
								},

								"sp_encryption": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"algorithm": schema.StringAttribute{
												Required: true,
											},
										},

										Blocks: map[string]schema.Block{
											"certificate": schema.ListNestedBlock{
												NestedObject: schema.NestedBlockObject{
													Attributes: map[string]schema.Attribute{
														"id": schema.StringAttribute{
															Required: true,

															CustomType: pingonetypes.ResourceIDType{},
														},
													},
												},
											},
										},
									},
								},

								"sp_verification": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"authn_request_signed": schema.BoolAttribute{
												Optional: true,
												Computed: true,

												Default: booldefault.StaticBool(false),
											},

											"certificate_ids": schema.SetAttribute{
												Required: true,

												ElementType: pingonetypes.ResourceIDType{},
											},
										},
									},
								},

								"cors_settings": r.resourceApplicationSchemaCorsSettingsV0(),
							},
						},
					},
				},
			},
			StateUpgrader: func(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
				var d diag.Diagnostics
				var priorStateData applicationResourceModelV0

				resp.Diagnostics.Append(req.State.Get(ctx, &priorStateData)...)

				if resp.Diagnostics.HasError() {
					return
				}

				icon, d := priorStateData.schemaUpgradeIconV0toV1(ctx)
				resp.Diagnostics.Append(d...)

				accessControlGroupOptions, d := priorStateData.schemaUpgradeAccessControlGroupOptionsV0toV1(ctx)
				resp.Diagnostics.Append(d...)

				externalLinkOptions, d := priorStateData.schemaUpgradeExternalLinkOptionsV0toV1(ctx)
				resp.Diagnostics.Append(d...)

				oIDCOptions, d := priorStateData.schemaUpgradeOIDCOptionsV0toV1(ctx)
				resp.Diagnostics.Append(d...)

				sAMLOptions, d := priorStateData.schemaUpgradeSAMLOptionsV0toV1(ctx)
				resp.Diagnostics.Append(d...)

				upgradedStateData := applicationResourceModelV1{
					Id:                        priorStateData.Id,
					EnvironmentId:             priorStateData.EnvironmentId,
					Name:                      priorStateData.Name,
					Description:               priorStateData.Description,
					Enabled:                   priorStateData.Enabled,
					Tags:                      priorStateData.Tags,
					LoginPageUrl:              priorStateData.LoginPageUrl,
					Icon:                      icon,
					AccessControlRoleType:     priorStateData.AccessControlRoleType,
					AccessControlGroupOptions: accessControlGroupOptions,
					HiddenFromAppPortal:       priorStateData.HiddenFromAppPortal,
					ExternalLinkOptions:       externalLinkOptions,
					OIDCOptions:               oIDCOptions,
					SAMLOptions:               sAMLOptions,
				}

				resp.Diagnostics.Append(resp.State.Set(ctx, upgradedStateData)...)
			},
		},
	}
}

func (r *ApplicationResource) resourceApplicationSchemaCorsSettingsV0() schema.Block {
	return schema.ListNestedBlock{
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"behavior": schema.StringAttribute{
					Required: true,
				},

				"origins": schema.SetAttribute{
					Optional: true,

					ElementType: types.StringType,
				},
			},
		},
	}
}

func (p *applicationResourceModelV0) schemaUpgradeIconV0toV1(ctx context.Context) (types.Object, diag.Diagnostics) {
	return service.ImageListToObjectSchemaUpgrade(ctx, p.Icon)
}

func (p *applicationResourceModelV0) schemaUpgradeAccessControlGroupOptionsV0toV1(ctx context.Context) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributeTypes := applicationAccessControlGroupOptionsTFObjectTypes
	planAttribute := p.AccessControlGroupOptions

	if planAttribute.IsNull() {
		return types.ObjectNull(attributeTypes), diags
	} else if planAttribute.IsUnknown() {
		return types.ObjectUnknown(attributeTypes), diags
	} else {
		var priorStateData []applicationAccessControlGroupOptionsResourceModelV0
		d := planAttribute.ElementsAs(ctx, &priorStateData, false)
		diags.Append(d...)
		if diags.HasError() {
			return types.ObjectNull(attributeTypes), diags
		}

		if len(priorStateData) == 0 {
			return types.ObjectNull(attributeTypes), diags
		}

		upgradedStateData := applicationAccessControlGroupOptionsResourceModelV1{
			Type:   priorStateData[0].Type,
			Groups: priorStateData[0].Groups,
		}

		returnVar, d := types.ObjectValueFrom(ctx, attributeTypes, upgradedStateData)
		diags.Append(d...)

		return returnVar, diags
	}
}

func (p *applicationResourceModelV0) schemaUpgradeExternalLinkOptionsV0toV1(ctx context.Context) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributeTypes := applicationExternalLinkOptionsTFObjectTypes
	planAttribute := p.ExternalLinkOptions

	if planAttribute.IsNull() {
		return types.ObjectNull(attributeTypes), diags
	} else if planAttribute.IsUnknown() {
		return types.ObjectUnknown(attributeTypes), diags
	} else {
		var priorStateData []applicationExternalLinkOptionsResourceModelV0
		d := planAttribute.ElementsAs(ctx, &priorStateData, false)
		diags.Append(d...)
		if diags.HasError() {
			return types.ObjectNull(attributeTypes), diags
		}

		if len(priorStateData) == 0 {
			return types.ObjectNull(attributeTypes), diags
		}

		upgradedStateData := applicationExternalLinkOptionsResourceModelV1{
			HomePageUrl: priorStateData[0].HomePageUrl,
		}

		returnVar, d := types.ObjectValueFrom(ctx, attributeTypes, upgradedStateData)
		diags.Append(d...)

		return returnVar, diags
	}
}

func (p *applicationResourceModelV0) schemaUpgradeOIDCOptionsV0toV1(ctx context.Context) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributeTypes := applicationOidcOptionsTFObjectTypes
	planAttribute := p.OIDCOptions

	if planAttribute.IsNull() {
		return types.ObjectNull(attributeTypes), diags
	} else if planAttribute.IsUnknown() {
		return types.ObjectUnknown(attributeTypes), diags
	} else {
		var priorStateData []applicationOIDCOptionsResourceModelV0
		d := planAttribute.ElementsAs(ctx, &priorStateData, false)
		diags.Append(d...)
		if diags.HasError() {
			return types.ObjectNull(attributeTypes), diags
		}

		if len(priorStateData) == 0 {
			return types.ObjectNull(attributeTypes), diags
		}

		certificateBasedAuthentication, d := priorStateData[0].schemaUpgradeCertificateBasedAuthenticationV0toV1(ctx)
		diags.Append(d...)

		corsSettings, d := priorStateData[0].schemaUpgradeCorsSettingsV0toV1(ctx)
		diags.Append(d...)

		mobileApp, d := priorStateData[0].schemaUpgradeMobileAppV0toV1(ctx)
		diags.Append(d...)

		const deviceTimeoutDefault = 600
		const devicePollingIntervalDefault = 5

		upgradedStateData := applicationOIDCOptionsResourceModelV1{
			ApplicationOIDCOptionsResourceModelV1:         beta.SchemaUpgradeV0toV1(priorStateData[0].ClientId),
			AdditionalRefreshTokenReplayProtectionEnabled: priorStateData[0].AdditionalRefreshTokenReplayProtectionEnabled,
			AllowWildcardsInRedirectUris:                  priorStateData[0].AllowWildcardsInRedirectUris,
			CertificateBasedAuthentication:                certificateBasedAuthentication,
			CorsSettings:                                  corsSettings,
			DevicePathId:                                  types.StringNull(),
			DeviceCustomVerificationUri:                   types.StringNull(),
			DeviceTimeout:                                 types.Int32Value(deviceTimeoutDefault),
			DevicePollingInterval:                         types.Int32Value(devicePollingIntervalDefault),
			GrantTypes:                                    priorStateData[0].GrantTypes,
			HomePageUrl:                                   priorStateData[0].HomePageUrl,
			InitiateLoginUri:                              priorStateData[0].InitiateLoginUri,
			Jwks:                                          priorStateData[0].Jwks,
			JwksUrl:                                       priorStateData[0].JwksUrl,
			MobileApp:                                     mobileApp,
			ParRequirement:                                priorStateData[0].ParRequirement,
			ParTimeout:                                    priorStateData[0].ParTimeout,
			PKCEEnforcement:                               priorStateData[0].PKCEEnforcement,
			PostLogoutRedirectUris:                        priorStateData[0].PostLogoutRedirectUris,
			RedirectUris:                                  priorStateData[0].RedirectUris,
			RefreshTokenDuration:                          priorStateData[0].RefreshTokenDuration,
			RefreshTokenRollingDuration:                   priorStateData[0].RefreshTokenRollingDuration,
			RefreshTokenRollingGracePeriodDuration:        priorStateData[0].RefreshTokenRollingGracePeriodDuration,
			RequireSignedRequestObject:                    priorStateData[0].RequireSignedRequestObject,
			ResponseTypes:                                 priorStateData[0].ResponseTypes,
			SupportUnsignedRequestObject:                  priorStateData[0].SupportUnsignedRequestObject,
			TargetLinkUri:                                 priorStateData[0].TargetLinkUri,
			TokenEndpointAuthnMethod:                      priorStateData[0].TokenEndpointAuthnMethod,
			Type:                                          priorStateData[0].Type,
		}

		returnVar, d := types.ObjectValueFrom(ctx, attributeTypes, upgradedStateData)
		diags.Append(d...)

		return returnVar, diags
	}
}

func (p *applicationOIDCOptionsResourceModelV0) schemaUpgradeCertificateBasedAuthenticationV0toV1(ctx context.Context) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributeTypes := applicationOidcOptionsCertificateAuthenticationTFObjectTypes
	planAttribute := p.CertificateBasedAuthentication

	if planAttribute.IsNull() {
		return types.ObjectNull(attributeTypes), diags
	} else if planAttribute.IsUnknown() {
		return types.ObjectUnknown(attributeTypes), diags
	} else {
		var priorStateData []applicationOIDCCertificateBasedAuthenticationResourceModelV0
		d := planAttribute.ElementsAs(ctx, &priorStateData, false)
		diags.Append(d...)
		if diags.HasError() {
			return types.ObjectNull(attributeTypes), diags
		}

		if len(priorStateData) == 0 {
			return types.ObjectNull(attributeTypes), diags
		}

		upgradedStateData := applicationOIDCCertificateBasedAuthenticationResourceModelV1{
			KeyId: priorStateData[0].KeyId,
		}

		returnVar, d := types.ObjectValueFrom(ctx, attributeTypes, upgradedStateData)
		diags.Append(d...)

		return returnVar, diags
	}
}

func (p *applicationOIDCOptionsResourceModelV0) schemaUpgradeCorsSettingsV0toV1(ctx context.Context) (types.Object, diag.Diagnostics) {
	planAttribute := p.CorsSettings
	return applicationSchemaUpgradeCorsSettingsV0toV1(ctx, planAttribute)
}

func applicationSchemaUpgradeCorsSettingsV0toV1(ctx context.Context, planAttribute types.List) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributeTypes := applicationCorsSettingsTFObjectTypes

	if planAttribute.IsNull() {
		return types.ObjectNull(attributeTypes), diags
	} else if planAttribute.IsUnknown() {
		return types.ObjectUnknown(attributeTypes), diags
	} else {
		var priorStateData []applicationCorsSettingsResourceModelV0
		d := planAttribute.ElementsAs(ctx, &priorStateData, false)
		diags.Append(d...)
		if diags.HasError() {
			return types.ObjectNull(attributeTypes), diags
		}

		if len(priorStateData) == 0 {
			return types.ObjectNull(attributeTypes), diags
		}

		upgradedStateData := applicationCorsSettingsResourceModelV1{
			Behavior: priorStateData[0].Behavior,
			Origins:  priorStateData[0].Origins,
		}

		returnVar, d := types.ObjectValueFrom(ctx, attributeTypes, upgradedStateData)
		diags.Append(d...)

		return returnVar, diags
	}
}

func (p *applicationOIDCOptionsResourceModelV0) schemaUpgradeMobileAppV0toV1(ctx context.Context) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributeTypes := applicationOidcMobileAppTFObjectTypes
	planAttribute := p.MobileApp

	if planAttribute.IsNull() {
		return types.ObjectNull(attributeTypes), diags
	} else if planAttribute.IsUnknown() {
		return types.ObjectUnknown(attributeTypes), diags
	} else {
		var priorStateData []applicationOIDCMobileAppResourceModelV0
		d := planAttribute.ElementsAs(ctx, &priorStateData, false)
		diags.Append(d...)
		if diags.HasError() {
			return types.ObjectNull(attributeTypes), diags
		}

		if len(priorStateData) == 0 {
			return types.ObjectNull(attributeTypes), diags
		}

		integrityDetection, d := priorStateData[0].schemaUpgradeIntegrityDetectionV0toV1(ctx)
		diags.Append(d...)

		upgradedStateData := applicationOIDCMobileAppResourceModelV1{
			BundleId:               priorStateData[0].BundleId,
			HuaweiAppId:            priorStateData[0].HuaweiAppId,
			HuaweiPackageName:      priorStateData[0].HuaweiPackageName,
			IntegrityDetection:     integrityDetection,
			PackageName:            priorStateData[0].PackageName,
			PasscodeRefreshSeconds: priorStateData[0].PasscodeRefreshSeconds,
			UniversalAppLink:       priorStateData[0].UniversalAppLink,
		}

		returnVar, d := types.ObjectValueFrom(ctx, attributeTypes, upgradedStateData)
		diags.Append(d...)

		return returnVar, diags
	}
}

func (p *applicationOIDCMobileAppResourceModelV0) schemaUpgradeIntegrityDetectionV0toV1(ctx context.Context) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributeTypes := applicationOidcMobileAppIntegrityDetectionTFObjectTypes
	planAttribute := p.IntegrityDetection

	if planAttribute.IsNull() {
		return types.ObjectNull(attributeTypes), diags
	} else if planAttribute.IsUnknown() {
		return types.ObjectUnknown(attributeTypes), diags
	} else {
		var priorStateData []applicationOIDCMobileAppIntegrityDetectionResourceModelV0
		d := planAttribute.ElementsAs(ctx, &priorStateData, false)
		diags.Append(d...)
		if diags.HasError() {
			return types.ObjectNull(attributeTypes), diags
		}

		if len(priorStateData) == 0 {
			return types.ObjectNull(attributeTypes), diags
		}

		cacheDuration, d := priorStateData[0].schemaUpgradeCacheDurationV0toV1(ctx)
		diags.Append(d...)

		googlePlay, d := priorStateData[0].schemaUpgradeGooglePlayV0toV1(ctx)
		diags.Append(d...)

		upgradedStateData := applicationOIDCMobileAppIntegrityDetectionResourceModelV1{
			CacheDuration:     cacheDuration,
			Enabled:           priorStateData[0].Enabled,
			ExcludedPlatforms: priorStateData[0].ExcludedPlatforms,
			GooglePlay:        googlePlay,
		}

		returnVar, d := types.ObjectValueFrom(ctx, attributeTypes, upgradedStateData)
		diags.Append(d...)

		return returnVar, diags
	}
}

func (p *applicationOIDCMobileAppIntegrityDetectionResourceModelV0) schemaUpgradeCacheDurationV0toV1(ctx context.Context) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributeTypes := applicationOidcMobileAppIntegrityDetectionCacheDurationTFObjectTypes
	planAttribute := p.CacheDuration

	if planAttribute.IsNull() {
		return types.ObjectNull(attributeTypes), diags
	} else if planAttribute.IsUnknown() {
		return types.ObjectUnknown(attributeTypes), diags
	} else {
		var priorStateData []applicationOIDCMobileAppIntegrityDetectionCacheDurationResourceModelV0
		d := planAttribute.ElementsAs(ctx, &priorStateData, false)
		diags.Append(d...)
		if diags.HasError() {
			return types.ObjectNull(attributeTypes), diags
		}

		if len(priorStateData) == 0 {
			return types.ObjectNull(attributeTypes), diags
		}

		upgradedStateData := applicationOIDCMobileAppIntegrityDetectionCacheDurationResourceModelV1{
			Amount: priorStateData[0].Amount,
			Units:  priorStateData[0].Units,
		}

		returnVar, d := types.ObjectValueFrom(ctx, attributeTypes, upgradedStateData)
		diags.Append(d...)

		return returnVar, diags
	}
}

func (p *applicationOIDCMobileAppIntegrityDetectionResourceModelV0) schemaUpgradeGooglePlayV0toV1(ctx context.Context) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributeTypes := applicationOidcMobileAppIntegrityDetectionGooglePlayTFObjectTypes
	planAttribute := p.GooglePlay

	if planAttribute.IsNull() {
		return types.ObjectNull(attributeTypes), diags
	} else if planAttribute.IsUnknown() {
		return types.ObjectUnknown(attributeTypes), diags
	} else {
		var priorStateData []applicationOIDCMobileAppIntegrityDetectionGooglePlayResourceModelV0
		d := planAttribute.ElementsAs(ctx, &priorStateData, false)
		diags.Append(d...)
		if diags.HasError() {
			return types.ObjectNull(attributeTypes), diags
		}

		if len(priorStateData) == 0 {
			return types.ObjectNull(attributeTypes), diags
		}

		upgradedStateData := applicationOIDCMobileAppIntegrityDetectionGooglePlayResourceModelV1{
			DecryptionKey:                 types.StringNull(),
			ServiceAccountCredentialsJson: jsontypes.NewNormalizedNull(),
			VerificationKey:               types.StringNull(),
			VerificationType:              priorStateData[0].VerificationType,
		}

		returnVar, d := types.ObjectValueFrom(ctx, attributeTypes, upgradedStateData)
		diags.Append(d...)

		return returnVar, diags
	}
}

func (p *applicationResourceModelV0) schemaUpgradeSAMLOptionsV0toV1(ctx context.Context) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributeTypes := applicationSamlOptionsTFObjectTypes
	planAttribute := p.SAMLOptions

	if planAttribute.IsNull() {
		return types.ObjectNull(attributeTypes), diags
	} else if planAttribute.IsUnknown() {
		return types.ObjectUnknown(attributeTypes), diags
	} else {
		var priorStateData []applicationSAMLOptionsResourceModelV0
		d := planAttribute.ElementsAs(ctx, &priorStateData, false)
		diags.Append(d...)
		if diags.HasError() {
			return types.ObjectNull(attributeTypes), diags
		}

		if len(priorStateData) == 0 {
			return types.ObjectNull(attributeTypes), diags
		}

		corsSettings, d := priorStateData[0].schemaUpgradeCorsSettingsV0toV1(ctx)
		diags.Append(d...)

		idpSigningKey, d := priorStateData[0].schemaUpgradeIdPSigningKeyV0toV1(ctx)
		diags.Append(d...)

		spEncryption, d := priorStateData[0].schemaUpgradeSpEncryptionV0toV1(ctx)
		diags.Append(d...)

		spVerification, d := priorStateData[0].schemaUpgradeSpVerificationV0toV1(ctx)
		diags.Append(d...)

		upgradedStateData := applicationSAMLOptionsResourceModelV1{
			AcsUrls:                     priorStateData[0].AcsUrls,
			AssertionDuration:           priorStateData[0].AssertionDuration,
			AssertionSignedEnabled:      priorStateData[0].AssertionSignedEnabled,
			CorsSettings:                corsSettings,
			EnableRequestedAuthnContext: priorStateData[0].EnableRequestedAuthnContext,
			HomePageUrl:                 priorStateData[0].HomePageUrl,
			IdpSigningKey:               idpSigningKey,
			DefaultTargetUrl:            priorStateData[0].DefaultTargetUrl,
			NameIdFormat:                priorStateData[0].NameIdFormat,
			ResponseIsSigned:            priorStateData[0].ResponseIsSigned,
			SloBinding:                  priorStateData[0].SloBinding,
			SloEndpoint:                 priorStateData[0].SloEndpoint,
			SloResponseEndpoint:         priorStateData[0].SloResponseEndpoint,
			SloWindow:                   priorStateData[0].SloWindow,
			SpEncryption:                spEncryption,
			SpEntityId:                  priorStateData[0].SpEntityId,
			SpVerification:              spVerification,
			Type:                        priorStateData[0].Type,
		}

		returnVar, d := types.ObjectValueFrom(ctx, attributeTypes, upgradedStateData)
		diags.Append(d...)

		return returnVar, diags
	}
}

func (p *applicationSAMLOptionsResourceModelV0) schemaUpgradeCorsSettingsV0toV1(ctx context.Context) (types.Object, diag.Diagnostics) {
	planAttribute := p.CorsSettings
	return applicationSchemaUpgradeCorsSettingsV0toV1(ctx, planAttribute)
}

func (p *applicationSAMLOptionsResourceModelV0) schemaUpgradeIdPSigningKeyV0toV1(ctx context.Context) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributeTypes := applicationIdpSigningKeyTFObjectTypes
	planAttribute := p.IdpSigningKey

	if planAttribute.IsNull() {
		return types.ObjectNull(attributeTypes), diags
	} else if planAttribute.IsUnknown() {
		return types.ObjectUnknown(attributeTypes), diags
	} else {
		var priorStateData []applicationSAMLOptionsIdpSigningKeyResourceModelV0
		d := planAttribute.ElementsAs(ctx, &priorStateData, false)
		diags.Append(d...)
		if diags.HasError() {
			return types.ObjectNull(attributeTypes), diags
		}

		if len(priorStateData) == 0 {
			return types.ObjectNull(attributeTypes), diags
		}

		upgradedStateData := applicationOptionsIdpSigningKeyResourceModelV1{
			Algorithm: priorStateData[0].Algorithm,
			KeyId:     priorStateData[0].KeyId,
		}

		returnVar, d := types.ObjectValueFrom(ctx, attributeTypes, upgradedStateData)
		diags.Append(d...)

		return returnVar, diags
	}
}

func (p *applicationSAMLOptionsResourceModelV0) schemaUpgradeSpEncryptionV0toV1(ctx context.Context) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributeTypes := applicationSamlOptionsSpEncryptionTFObjectTypes
	planAttribute := p.SpEncryption

	if planAttribute.IsNull() {
		return types.ObjectNull(attributeTypes), diags
	} else if planAttribute.IsUnknown() {
		return types.ObjectUnknown(attributeTypes), diags
	} else {
		var priorStateData []applicationSAMLOptionsSpEncryptionResourceModelV0
		d := planAttribute.ElementsAs(ctx, &priorStateData, false)
		diags.Append(d...)
		if diags.HasError() {
			return types.ObjectNull(attributeTypes), diags
		}

		if len(priorStateData) == 0 {
			return types.ObjectNull(attributeTypes), diags
		}

		certificate, d := priorStateData[0].schemaUpgradeCertificateV0toV1(ctx)
		diags.Append(d...)

		upgradedStateData := applicationSAMLOptionsSpEncryptionResourceModelV1{
			Algorithm:   priorStateData[0].Algorithm,
			Certificate: certificate,
		}

		returnVar, d := types.ObjectValueFrom(ctx, attributeTypes, upgradedStateData)
		diags.Append(d...)

		return returnVar, diags
	}
}

func (p *applicationSAMLOptionsSpEncryptionResourceModelV0) schemaUpgradeCertificateV0toV1(ctx context.Context) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributeTypes := applicationSamlOptionsSpEncryptionCertificateTFObjectTypes
	planAttribute := p.Certificate

	if planAttribute.IsNull() {
		return types.ObjectNull(attributeTypes), diags
	} else if planAttribute.IsUnknown() {
		return types.ObjectUnknown(attributeTypes), diags
	} else {
		var priorStateData []applicationSAMLOptionsSpEncryptionCertificateResourceModelV0
		d := planAttribute.ElementsAs(ctx, &priorStateData, false)
		diags.Append(d...)
		if diags.HasError() {
			return types.ObjectNull(attributeTypes), diags
		}

		if len(priorStateData) == 0 {
			return types.ObjectNull(attributeTypes), diags
		}

		upgradedStateData := applicationSAMLOptionsSpEncryptionCertificateResourceModelV1{
			Id: priorStateData[0].Id,
		}

		returnVar, d := types.ObjectValueFrom(ctx, attributeTypes, upgradedStateData)
		diags.Append(d...)

		return returnVar, diags
	}
}

func (p *applicationSAMLOptionsResourceModelV0) schemaUpgradeSpVerificationV0toV1(ctx context.Context) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributeTypes := applicationSamlOptionsSpVerificationTFObjectTypes
	planAttribute := p.SpVerification

	if planAttribute.IsNull() {
		return types.ObjectNull(attributeTypes), diags
	} else if planAttribute.IsUnknown() {
		return types.ObjectUnknown(attributeTypes), diags
	} else {
		var priorStateData []applicationSAMLOptionsSpVerificationResourceModelV0
		d := planAttribute.ElementsAs(ctx, &priorStateData, false)
		diags.Append(d...)
		if diags.HasError() {
			return types.ObjectNull(attributeTypes), diags
		}

		if len(priorStateData) == 0 {
			return types.ObjectNull(attributeTypes), diags
		}

		upgradedStateData := applicationSAMLOptionsSpVerificationResourceModelV1{
			CertificateIds:     priorStateData[0].CertificateIds,
			AuthnRequestSigned: priorStateData[0].AuthnRequestSigned,
		}

		returnVar, d := types.ObjectValueFrom(ctx, attributeTypes, upgradedStateData)
		diags.Append(d...)

		return returnVar, diags
	}
}
