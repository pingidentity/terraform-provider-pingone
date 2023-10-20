package sso

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	validation "github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type ApplicationDataSource serviceClientType

// todo review - will i need to update set/list to objects; double check
type applicationDataSourceModel struct {
	Id                        types.String `tfsdk:"id"`
	EnvironmentId             types.String `tfsdk:"environment_id"`
	ApplicationId             types.String `tfsdk:"application_id"`
	Name                      types.String `tfsdk:"name"`
	Description               types.String `tfsdk:"description"`
	Enabled                   types.Bool   `tfsdk:"enabled"`
	Tags                      types.Set    `tfsdk:"tags"`
	LoginPageUrl              types.String `tfsdk:"login_page_url"`
	Icon                      types.List   `tfsdk:"icon"`
	AccessControlRoleType     types.String `tfsdk:"access_control_role_type"`
	AccessControlGroupOptions types.List   `tfsdk:"access_control_group_options"`
	HiddenFromAppPortal       types.Bool   `tfsdk:"hidden_from_app_portal"`
	ExternalLinkOptions       types.List   `tfsdk:"external_link_options"`
	OIDCOptions               types.List   `tfsdk:"oidc_options"`
	SAMLOptions               types.List   `tfsdk:"saml_options"`
}

var (
	applicationOIDCTFObjectTypes = map[string]attr.Type{
		"name":                         types.StringType,
		"enabled":                      types.StringType,
		"description":                  types.StringType,
		"tags":                         types.SetType{ElemType: types.StringType},
		"login_page_url":               types.StringType,
		"icon":                         types.ListType{ElemType: types.ObjectType{AttrTypes: iconTFObjectTypes}},
		"access_control_role_type":     types.StringType,
		"access_control_group_options": types.ListType{ElemType: types.ObjectType{AttrTypes: accessControlGroupOptionsTFObjectTypes}},
		"hidden_from_app_portal":       types.BoolType,
		"oidc_options":                 types.ListType{ElemType: types.ObjectType{AttrTypes: oidcOptionsTFObjectTypes}},
	}

	oidcOptionsTFObjectTypes = map[string]attr.Type{
		"type":                                        types.StringType,
		"home_page_url":                               types.StringType,
		"initiate_login_uri":                          types.StringType,
		"target_link_uri":                             types.StringType,
		"grant_types":                                 types.SetType{ElemType: types.StringType},
		"response_types":                              types.SetType{ElemType: types.StringType},
		"token_endpoint_authn_method":                 types.StringType,
		"par_requirement":                             types.StringType,
		"par_timeout":                                 types.Int64Type,
		"pkce_enforcement":                            types.StringType,
		"redirect_uris":                               types.SetType{ElemType: types.StringType},
		"allow_wildcards_in_redirect_uris":            types.BoolType,
		"post_logout_redirect_uris":                   types.SetType{ElemType: types.StringType},
		"refresh_token_duration":                      types.Int64Type,
		"refresh_token_rolling_duration":              types.Int64Type,
		"refresh_token_rolling_grace_period_duration": types.Int64Type,
		"additional_refresh_token_replay_protection_enabled": types.BoolType,
		"client_id":                        types.StringType,
		"client_secret":                    types.StringType,
		"certificate_based_authentication": types.ListType{ElemType: types.ObjectType{AttrTypes: certificateAuthenticationTFObjectTypes}},
		"support_unsigned_request_object":  types.BoolType,
		"require_signed_request_object":    types.BoolType,
		"mobile_app":                       types.ListType{ElemType: types.ObjectType{AttrTypes: mobileAppTFObjectTypes}},
		"bundle_id":                        types.StringType,
		"package_name":                     types.StringType,
	}

	mobileAppTFObjectTypes = map[string]attr.Type{
		"bundle_id":                types.StringType,
		"package_name":             types.StringType,
		"huawei_app_id":            types.StringType,
		"huawei_package_name":      types.StringType,
		"passcode_refresh_seconds": types.Int64Type,
		"universal_app_link":       types.StringType,
		"integrity_detection":      types.ListType{ElemType: types.ObjectType{AttrTypes: integrityDetectionTFObjectTypes}},
	}

	integrityDetectionTFObjectTypes = map[string]attr.Type{
		"enabled":            types.BoolType,
		"excluded_platforms": types.ListType{ElemType: types.StringType},
		"cache_duration":     types.ListType{ElemType: types.ObjectType{AttrTypes: cacheDurationTFObjectTypes}},
		"google_play":        types.ListType{ElemType: types.ObjectType{AttrTypes: googlePlayTFObjectTypes}},
	}

	cacheDurationTFObjectTypes = map[string]attr.Type{
		"amount": types.Int64Type,
		"units":  types.StringType,
	}

	googlePlayTFObjectTypes = map[string]attr.Type{
		"decryption_key":                   types.StringType,
		"service_account_credentials_json": types.StringType,
		"verification_key":                 types.StringType,
		"verification_type":                types.StringType,
	}

	certificateAuthenticationTFObjectTypes = map[string]attr.Type{
		"key_id": types.StringType,
	}

	applicationSAMLTFObjectTypes = map[string]attr.Type{
		"name":                         types.StringType,
		"enabled":                      types.StringType,
		"description":                  types.StringType,
		"login_page_url":               types.StringType,
		"icon":                         types.ListType{ElemType: types.ObjectType{AttrTypes: iconTFObjectTypes}},
		"access_control_role_type":     types.StringType,
		"access_control_group_options": types.ListType{ElemType: types.ObjectType{AttrTypes: accessControlGroupOptionsTFObjectTypes}},
		"hidden_from_app_portal":       types.BoolType,
		"saml_options":                 types.ListType{ElemType: types.ObjectType{AttrTypes: samlOptionsTFObjectTypes}},
	}

	samlOptionsTFObjectTypes = map[string]attr.Type{
		"home_page_url":                   types.StringType,
		"type":                            types.StringType,
		"acs_urls":                        types.SetType{ElemType: types.StringType},
		"assertion_duration":              types.Int64Type,
		"assertion_signed_enabled":        types.BoolType,
		"idp_signing_key_id":              types.StringType,
		"idp_signing_key":                 types.ListType{ElemType: types.ObjectType{AttrTypes: idpSigningKeyTFObjectTypes}},
		"enable_requested_authn_context":  types.BoolType,
		"nameid_format":                   types.StringType,
		"response_is_signed":              types.BoolType,
		"slo_binding":                     types.StringType,
		"slo_endpoint":                    types.StringType,
		"slo_response_endpoint":           types.StringType,
		"slo_window":                      types.Int64Type,
		"sp_entity_id":                    types.StringType,
		"sp_verification_certificate_ids": types.SetType{ElemType: types.StringType},
	}

	idpSigningKeyTFObjectTypes = map[string]attr.Type{
		"algorithm": types.StringType,
		"key_id":    types.StringType,
	}

	applicationExternalLinkTFObjectTypes = map[string]attr.Type{
		"name":                         types.StringType,
		"enabled":                      types.StringType,
		"description":                  types.StringType,
		"login_page_url":               types.StringType,
		"icon":                         types.ListType{ElemType: types.ObjectType{AttrTypes: iconTFObjectTypes}},
		"access_control_group_options": types.ListType{ElemType: types.ObjectType{AttrTypes: accessControlGroupOptionsTFObjectTypes}},
		"hidden_from_app_portal":       types.BoolType,
		"external_link_options":        types.ListType{ElemType: types.ObjectType{AttrTypes: externalLinkOptionsTFObjectTypes}},
	}

	externalLinkOptionsTFObjectTypes = map[string]attr.Type{
		"home_page_url": types.StringType,
	}

	iconTFObjectTypes = map[string]attr.Type{
		"id":   types.StringType,
		"href": types.StringType,
	}

	accessControlGroupOptionsTFObjectTypes = map[string]attr.Type{
		"type":   types.StringType,
		"groups": types.SetType{ElemType: types.StringType},
	}
)

// Framework interfaces
var (
	_ datasource.DataSource = &ApplicationDataSource{}
)

// New Object
func NewApplicationDataSource() datasource.DataSource {
	return &ApplicationDataSource{}
}

// Metadata
func (r *ApplicationDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_application"
}

func (r *ApplicationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	// schema descriptions and validation settings
	const attrMinLength = 1

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Data source to retrieve a PingOne application.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),
			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("PingOne environment identifier (UUID) in which the application exists."),
			),
			"application_id": schema.StringAttribute{
				Description: "A string that specifies the identifier (UUID) of the application.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(
						path.MatchRelative().AtParent().AtName("name"),
					),
					validation.P1ResourceIDValidator(),
				},
			},
			"name": schema.StringAttribute{
				Description: "A string that specifies the name of the application.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(
						path.MatchRelative().AtParent().AtName("application_id"),
					),
					stringvalidator.LengthAtLeast(attrMinLength),
				},
			},
			"description": schema.StringAttribute{
				Description: "A string that specifies the description of the application.",
				Computed:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "A boolean that specifies whether the application is enabled in the environment.",
				Computed:    true,
			},
			"tags": schema.SetAttribute{
				Description: "An array that specifies the list of labels associated with the application.",
				ElementType: types.StringType,
				Computed:    true,
			},
			"login_page_url": schema.StringAttribute{
				Description: "A string that specifies the custom login page URL for the application.",
				Computed:    true,
			},
			"icon": schema.ListNestedAttribute{
				Description: "The HREF and the ID for the application icon.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "The ID for the application icon.",
							Computed:    true,
						},
						"href": schema.StringAttribute{
							Description: "The HREF for the application icon.",
							Computed:    true,
						},
					},
				},
			},
			"access_control_role_type": schema.StringAttribute{
				Description: "A string that specifies the user role required to access the application. A user is an admin user if the user has one or more of the following roles: Organization Admin, Environment Admin, Identity Data Admin, or Client Application Developer.",
				Computed:    true,
			},
			"access_control_group_options": schema.ListNestedAttribute{
				Description: "Group access control settings.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Description: "A string that specifies the group type required to access the application.",
							Computed:    true,
						},
						"groups": schema.SetAttribute{
							Description: "A set that specifies the group IDs for the groups the actor must belong to for access to the application.",
							ElementType: types.StringType,
							Computed:    true,
						},
					},
				},
			},
			"hidden_from_app_portal": schema.BoolAttribute{
				Description: "A boolean to specify whether the application is hidden in the application portal despite the configured group access policy.",
				Computed:    true,
			},
			"external_link_options": schema.ListNestedAttribute{
				Description: "External link application specific settings.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"home_page_url": schema.StringAttribute{
							Description: "A string that specifies the custom home page URL for the application.  Both `http://` and `https://` URLs are permitted.",
							Computed:    true,
						},
					},
				},
			},
			"oidc_options": schema.ListNestedAttribute{
				Description: "OIDC/OAuth application specific settings.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Description: "A string that specifies the type associated with the application.",
							Computed:    true,
						},
						"home_page_url": schema.StringAttribute{
							Description: "A string that specifies the custom home page URL for the application.",
							Computed:    true,
						},
						"initiate_login_uri": schema.StringAttribute{
							Description: "A string that specifies the URI to use for third-parties to begin the sign-on process for the application.",
							Computed:    true,
						},
						"target_link_uri": schema.StringAttribute{
							Description: "The URI for the application.",
							Computed:    true,
						},
						"grant_types": schema.SetAttribute{
							Description: "A list that specifies the grant type for the authorization request.",
							ElementType: types.StringType,
							Computed:    true,
						},
						"response_types": schema.SetAttribute{
							Description: "A list that specifies the code or token type returned by an authorization request.",
							ElementType: types.StringType,
							Computed:    true,
						},
						"token_endpoint_authn_method": schema.StringAttribute{
							Description: "A string that specifies the client authentication methods supported by the token endpoint.",
							Computed:    true,
						},
						"par_requirement": schema.StringAttribute{
							Description: "A string that specifies whether pushed authorization requests (PAR) are required.",
							Computed:    true,
						},
						"par_timeout": schema.Int64Attribute{
							Description: "An integer that specifies the pushed authorization request (PAR) timeout in seconds.",
							Computed:    true,
						},
						"pkce_enforcement": schema.StringAttribute{
							Description: "A string that specifies how `PKCE` request parameters are handled on the authorize request.",
							Computed:    true,
						},
						"redirect_uris": schema.SetAttribute{
							Description: "A list of strings that specifies the allowed callback URIs for the authentication response.",
							ElementType: types.StringType,
							Computed:    true,
						},
						"allow_wildcards_in_redirect_uris": schema.BoolAttribute{
							Description: "A boolean to specify whether wildcards are allowed in redirect URIs. For more information, see [Wildcards in Redirect URIs](https://docs.pingidentity.com/csh?context=p1_c_wildcard_redirect_uri).",
							Computed:    true,
						},
						"post_logout_redirect_uris": schema.SetAttribute{
							Description: "A list of strings that specifies the URLs that the browser can be redirected to after logout.  The provided URLs are expected to use the `https://`, `http://` schema, or a custom mobile native schema (e.g., `org.bxretail.app://logout`).",
							ElementType: types.StringType,
							Computed:    true,
						},
						"refresh_token_duration": schema.Int64Attribute{
							Description: "An integer that specifies the lifetime in seconds of the refresh token.",
							Computed:    true,
						},
						"refresh_token_rolling_duration": schema.Int64Attribute{
							Description: "An integer that specifies the number of seconds a refresh token can be exchanged before re-authentication is required.",
							Computed:    true,
						},
						"refresh_token_rolling_grace_period_duration": schema.Int64Attribute{
							Description: "The number of seconds that a refresh token may be reused after having been exchanged for a new set of tokens.",
							Computed:    true,
						},
						"additional_refresh_token_replay_protection_enabled": schema.BoolAttribute{
							Description: "A boolean that, when set to `true` (the default), if you attempt to reuse the refresh token, the authorization server immediately revokes the reused refresh token, as well as all descendant tokens.",
							Computed:    true,
						},
						"client_id": schema.StringAttribute{
							Description: "A string that specifies the application ID used to authenticate to the authorization server.",
							Computed:    true,
						},
						"client_secret": schema.StringAttribute{
							Description: "A string that specifies the application secret ID used to authenticate to the authorization server.",
							Computed:    true,
						},
						"certificate_based_authentication": schema.ListNestedAttribute{
							Description: "Certificate based authentication settings.",
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"key_id": schema.StringAttribute{
										Description: "A string that represents a PingOne ID for the issuance certificate key.",
										Computed:    true,
									},
								},
							},
						},
						"support_unsigned_request_object": schema.BoolAttribute{
							Description: "A boolean that specifies whether the request query parameter JWT is allowed to be unsigned.",
							Computed:    true,
						},
						"require_signed_request_object": schema.BoolAttribute{
							Description: "A boolean that indicates that the Java Web Token (JWT) for the [request query](https://openid.net/specs/openid-connect-core-1_0.html#RequestObject) parameter is required to be signed. If `false` or null (default), a signed request object is not required. Both `support_unsigned_request_object` and this property cannot be set to `true`.",
							Computed:    true,
						},
						"mobile_app": schema.ListNestedAttribute{
							Description: "Mobile application integration settings.",
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"bundle_id": schema.StringAttribute{
										Description: "A string that specifies the bundle associated with the application, for push notifications in native apps.",
										Computed:    true,
									},
									"package_name": schema.StringAttribute{
										Description: "A string that specifies the package name associated with the application, for push notifications in native apps.",
										Computed:    true,
									},
									"huawei_app_id": schema.StringAttribute{
										Description: "The unique identifier for the app on the device and in the Huawei Mobile Service AppGallery.",
										Computed:    true,
									},
									"huawei_package_name": schema.StringAttribute{
										Description: "The package name associated with the application, for push notifications in native apps.",
										Computed:    true,
									},
									"passcode_refresh_seconds": schema.Int64Attribute{
										Description: "The amount of time a passcode should be displayed before being replaced with a new passcode.",
										Computed:    true,
									},
									"universal_app_link": schema.StringAttribute{
										Description: "A string that specifies a URI prefix that enables direct triggering of the mobile application when scanning a QR code.",
										Computed:    true,
									},
									"integrity_detection": schema.ListNestedAttribute{
										Description: "Mobile application integrity detection settings.",
										Computed:    true,
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"enabled": schema.BoolAttribute{
													Description: "A boolean that specifies whether device integrity detection takes place on mobile devices.",
													Computed:    true,
												},
												"excluded_platforms": schema.ListAttribute{
													Description: "Indicates OS excluded from device integrity checking.",
													ElementType: types.StringType,
													Computed:    true,
												},
												"cache_duration": schema.SingleNestedAttribute{
													Description: "Indicates the caching duration of successful integrity detection calls.",
													Computed:    true,

													Attributes: map[string]schema.Attribute{
														"amount": schema.Int64Attribute{
															Description: "An integer that specifies the number of minutes or hours that specify the duration between successful integrity detection calls.",
															Computed:    true,
														},
														"units": schema.StringAttribute{
															Description: "A string that specifies the cache duration time units.",
															Computed:    true,
														},
													},
												},
												"google_play": schema.SingleNestedAttribute{
													Description: "A single block that describes Google Play Integrity API credential settings for Android device integrity detection.",
													Computed:    true,

													Attributes: map[string]schema.Attribute{
														"decryption_key": schema.StringAttribute{
															Description: "Play Integrity verdict decryption key from your Google Play Services account. This parameter must be provided if you have set `verification_type` to `INTERNAL`.  Cannot be set with `service_account_credentials_json`.",
															Computed:    true,
														},
														"service_account_credentials_json": schema.StringAttribute{
															Description: "Contents of the JSON file that represents your Service Account Credentials.",
															Computed:    true,
														},
														"verification_key": schema.StringAttribute{
															Description: "Play Integrity verdict signature verification key from your Google Play Services account.",
															Computed:    true,
														},
														"verification_type": schema.StringAttribute{
															Description: "The type of verification.",
															Computed:    true,
														},
													},
												},
											},
										},
									},
								},
							},
						},
						"bundle_id": schema.StringAttribute{
							Description: "**Deprecation Notice** This field is deprecated and will be removed in a future release. Use `oidc_options.mobile_app.bundle_id` instead. A string that specifies the bundle associated with the application, for push notifications in native apps.",
							Computed:    true,
						},
						"package_name": schema.StringAttribute{
							Description: "**Deprecation Notice** This field is deprecated and will be removed in a future release. Use `oidc_options.mobile_app.package_name` instead. A string that specifies the package name associated with the application, for push notifications in native apps.",
							Computed:    true,
						},
					},
				},
			},
			"saml_options": schema.ListNestedAttribute{
				Description: "SAML application specific settings.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"home_page_url": schema.StringAttribute{
							Description: "A string that specifies the custom home page URL for the application.",
							Computed:    true,
						},
						"type": schema.StringAttribute{
							Description: "A string that specifies the type associated with the application.",
							Computed:    true,
						},
						"acs_urls": schema.SetAttribute{
							Description: "A list of string that specifies the Assertion Consumer Service URLs. The first URL in the list is used as default (there must be at least one URL).",
							ElementType: types.StringType,
							Computed:    true,
						},
						"assertion_duration": schema.Int64Attribute{
							Description: "An integer that specifies the assertion validity duration in seconds.",
							Computed:    true,
						},
						"assertion_signed_enabled": schema.BoolAttribute{
							Description: "A boolean that specifies whether the SAML assertion itself should be signed.",
							Computed:    true,
						},
						"idp_signing_key_id": schema.StringAttribute{
							Description: "An ID for the certificate key pair to be used by the identity provider to sign assertions and responses.",
							Computed:    true,
						},
						"idp_signing_key": schema.SingleNestedAttribute{
							Description: "SAML application assertion/response signing key settings.",
							Computed:    true,

							Attributes: map[string]schema.Attribute{
								"algorithm": schema.StringAttribute{
									Description: "A string that specifies the signature algorithm of the key.",
									Computed:    true,
								},
								"key_id": schema.StringAttribute{
									Description: "An ID for the certificate key pair to be used by the identity provider to sign assertions and responses.",
									Computed:    true,
								},
							},
						},
						"enable_requested_authn_context": schema.BoolAttribute{
							Description: "A boolean that specifies whether `requestedAuthnContext` is taken into account in policy decision-making.",
							Computed:    true,
						},
						"nameid_format": schema.StringAttribute{
							Description: "A string that specifies the format of the Subject NameID attibute in the SAML assertion.",
							Computed:    true,
						},
						"response_is_signed": schema.BoolAttribute{
							Description: "A boolean that specifies whether the SAML assertion response itself should be signed.",
							Computed:    true,
						},
						"slo_binding": schema.StringAttribute{
							Description: "A string that specifies the binding protocol to be used for the logout response.",
							Computed:    true,
						},
						"slo_endpoint": schema.StringAttribute{
							Description: "A string that specifies the logout endpoint URL.",
							Computed:    true,
						},
						"slo_response_endpoint": schema.StringAttribute{
							Description: "A string that specifies the endpoint URL to submit the logout response.",
							Computed:    true,
						},
						"slo_window": schema.Int64Attribute{
							Description: "An integer that defines how long (hours) PingOne can exchange logout messages with the application, specifically a logout request from the application, since the initial request.",
							Computed:    true,
						},
						"sp_entity_id": schema.StringAttribute{
							Description: "A string that specifies the service provider entity ID used to lookup the application. This is a required property and is unique within the environment.",
							Computed:    true,
						},
						"sp_verification_certificate_ids": schema.SetAttribute{
							Description: "A list that specifies the certificate IDs used to verify the service provider signature.",
							ElementType: types.StringType,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (r *ApplicationDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	preparedClient, err := PrepareClient(ctx, resourceConfig)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			err.Error(),
		)

		return
	}

	r.Client = preparedClient
}

func (r *ApplicationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *applicationDataSourceModel

	if r.Client == nil {
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

	var application management.ReadOneApplication200Response

	// Application API does not support SCIM filtering
	if !data.ApplicationId.IsNull() {
		// Run the API call
		var response *management.ReadOneApplication200Response
		resp.Diagnostics.Append(framework.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				return r.Client.ApplicationsApi.ReadOneApplication(ctx, data.EnvironmentId.ValueString(), data.ApplicationId.ValueString()).Execute()
			},
			"ReadOneApplication",
			framework.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
			&response,
		)...)
		if resp.Diagnostics.HasError() {
			return
		}

		application = *response

	} else if !data.Name.IsNull() {
		// Run the API call
		var entityArray *management.EntityArray
		resp.Diagnostics.Append(framework.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				return r.Client.ApplicationsApi.ReadAllApplications(ctx, data.EnvironmentId.ValueString()).Execute()
			},
			"ReadAllApplications",
			framework.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
			&entityArray,
		)...)
		if resp.Diagnostics.HasError() {
			return
		}

		if applications, ok := entityArray.Embedded.GetApplicationsOk(); ok {
			found := false

			for _, applicationObject := range applications {
				applicationInstance := applicationObject.GetActualInstance()

				applicationName := ""

				switch v := applicationInstance.(type) {
				case *management.ApplicationExternalLink:
					applicationName = v.GetName()

				case *management.ApplicationOIDC:
					applicationName = v.GetName()

				case *management.ApplicationSAML:
					applicationName = v.GetName()
				}
				if applicationName == data.Name.ValueString() {
					application = applicationObject
					found = true
					break
				}
			}

			if !found {
				resp.Diagnostics.AddError(
					"Cannot find the application from name",
					fmt.Sprintf("The application name %s for environment %s cannot be found", data.Name.String(), data.EnvironmentId.String()),
				)
				return
			}

		}
	} else {
		resp.Diagnostics.AddError(
			"Missing parameter",
			"Cannot find the requested PingOne Application: application_id or name argument must be set.",
		)
		return
	}

	// todo: Talk to Patrick about this secondary call for the secret if don't find in api docs what/why
	if application.ApplicationOIDC != nil && application.ApplicationOIDC.GetId() != "" {

	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(&application)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (p *applicationDataSourceModel) toState(apiObject *management.ReadOneApplication200Response) diag.Diagnostics {
	var diags diag.Diagnostics
	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	applicationInstance := apiObject.GetActualInstance()
	switch v := applicationInstance.(type) {
	case *management.ApplicationExternalLink:
		p.Id = framework.StringOkToTF(v.GetIdOk())
		p.EnvironmentId = framework.StringToTF(*v.GetEnvironment().Id)
		p.ApplicationId = framework.StringOkToTF(v.GetIdOk())
		p.Name = framework.StringOkToTF(v.GetNameOk())
		p.Description = framework.StringOkToTF(v.GetDescriptionOk())
		p.HiddenFromAppPortal = framework.BoolOkToTF(v.GetHiddenFromAppPortalOk())

		var d diag.Diagnostics

		p.Icon, d = p.toStateIcon(v.GetIconOk())
		diags.Append(d...)

		p.AccessControlRoleType, p.AccessControlGroupOptions, d = p.toStateAccessControl(v.GetAccessControlOk())
		diags.Append(d...)

		p.ExternalLinkOptions, d = p.toStateExternalLinkOptions(v.GetHomePageUrlOk())
		diags.Append(d...)

	case *management.ApplicationOIDC:
		p.Id = framework.StringOkToTF(v.GetIdOk())
		p.EnvironmentId = framework.StringToTF(*v.GetEnvironment().Id)
		p.ApplicationId = framework.StringOkToTF(v.GetIdOk())
		p.Name = framework.StringOkToTF(v.GetNameOk())
		p.Description = framework.StringOkToTF(v.GetDescriptionOk())
		p.Enabled = framework.BoolOkToTF(v.GetEnabledOk())
		p.Tags = framework.EnumSetOkToTF(v.GetTagsOk())
		p.LoginPageUrl = framework.StringOkToTF(v.GetLoginPageUrlOk())
		p.HiddenFromAppPortal = framework.BoolOkToTF(v.GetHiddenFromAppPortalOk())

		var d diag.Diagnostics

		p.Icon, d = p.toStateIcon(v.GetIconOk())
		diags.Append(d...)

		p.AccessControlRoleType, p.AccessControlGroupOptions, d = p.toStateAccessControl(v.GetAccessControlOk())
		diags.Append(d...)

		p.OIDCOptions, d = p.toStateOIDCOptions(v)
		diags.Append(d...)

	case *management.ApplicationSAML:
		p.Id = framework.StringOkToTF(v.GetIdOk())
		p.EnvironmentId = framework.StringToTF(*v.GetEnvironment().Id)
		p.ApplicationId = framework.StringOkToTF(v.GetIdOk())
		p.Name = framework.StringOkToTF(v.GetNameOk())
		p.Description = framework.StringOkToTF(v.GetDescriptionOk())
		p.Enabled = framework.BoolOkToTF(v.GetEnabledOk())
		p.LoginPageUrl = framework.StringOkToTF(v.GetLoginPageUrlOk())
		p.HiddenFromAppPortal = framework.BoolOkToTF(v.GetHiddenFromAppPortalOk())

		var d diag.Diagnostics

		p.Icon, d = p.toStateIcon(v.GetIconOk())
		diags.Append(d...)

		p.AccessControlRoleType, p.AccessControlGroupOptions, d = p.toStateAccessControl(v.GetAccessControlOk())
		diags.Append(d...)

		//toStateSAMLOptions
	}

	return diags
}

func (p *applicationDataSourceModel) toStateExternalLinkOptions(apiObject *string, ok bool) (basetypes.ListValue, diag.Diagnostics) {
	var diags diag.Diagnostics
	tfObjType := types.ObjectType{AttrTypes: externalLinkOptionsTFObjectTypes}

	if !ok || apiObject == nil {
		return types.ListNull(tfObjType), diags
	}

	externalLinkOptions := map[string]attr.Value{
		"home_page_url": framework.StringToTF(*apiObject),
	}

	flattenedObj, d := types.ObjectValue(externalLinkOptionsTFObjectTypes, externalLinkOptions)
	diags.Append(d...)

	returnVar, d := types.ListValue(tfObjType, append([]attr.Value{}, flattenedObj))
	diags.Append(d...)

	return returnVar, diags
}

func (p *applicationDataSourceModel) toStateIcon(apiObject *management.ApplicationIcon, ok bool) (basetypes.ListValue, diag.Diagnostics) {
	var diags diag.Diagnostics
	tfObjType := types.ObjectType{AttrTypes: iconTFObjectTypes}

	if !ok || apiObject == nil {
		return types.ListNull(tfObjType), diags
	}

	applicationIcon := map[string]attr.Value{
		"id":   framework.StringOkToTF(apiObject.GetIdOk()),
		"href": framework.StringOkToTF(apiObject.GetHrefOk()),
	}

	flattenedObj, d := types.ObjectValue(iconTFObjectTypes, applicationIcon)
	diags.Append(d...)

	returnVar, d := types.ListValue(tfObjType, append([]attr.Value{}, flattenedObj))
	diags.Append(d...)

	return returnVar, diags
}

func (p *applicationDataSourceModel) toStateAccessControl(apiObject *management.ApplicationAccessControl, ok bool) (basetypes.StringValue, basetypes.ListValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	objAccessControlGroupOptions := types.ListUnknown(types.ObjectType{AttrTypes: accessControlGroupOptionsTFObjectTypes})
	accessControlRoleType := types.StringNull()
	accessControlGroupOption := []attr.Value{}

	if !ok || apiObject == nil {
		return accessControlRoleType, objAccessControlGroupOptions, diags
	}

	accessControlRoleType = framework.EnumOkToTF(apiObject.Role.GetTypeOk())

	if group, ok := apiObject.GetGroupOk(); ok {
		var d diag.Diagnostics
		groups := make([]string, 0)
		for _, v := range group.GetGroups() {
			groups = append(groups, v.GetId())
		}

		groupObj := map[string]attr.Value{
			"type":   framework.StringOkToTF(group.GetTypeOk()),
			"groups": framework.StringSetToTF(groups),
		}

		groupsObj, d := types.ObjectValue(accessControlGroupOptionsTFObjectTypes, groupObj)
		diags.Append(d...)

		accessControlGroupOption = append(accessControlGroupOption, groupsObj)

		objAccessControlGroupOptions, d = types.ListValue(types.ObjectType{AttrTypes: accessControlGroupOptionsTFObjectTypes}, accessControlGroupOption)
		diags.Append(d...)
	}

	return accessControlRoleType, objAccessControlGroupOptions, diags

}

// todo - fix this
func (p *applicationDataSourceModel) toStateOIDCOptions(apiObject *management.ApplicationOIDC) (basetypes.ListValue, diag.Diagnostics) {
	var diags diag.Diagnostics
	objOIDCOptions := []attr.Value{}

	// todo - need to ensure nested objects are also null
	if apiObject == nil {
		return types.ListUnknown(types.ObjectType{AttrTypes: oidcOptionsTFObjectTypes}), diags
	}

	oidcOptions := map[string]attr.Value{
		"client_id": framework.StringOkToTF(apiObject.GetIdOk()),
		"type":      framework.EnumOkToTF(apiObject.GetTypeOk()),
		//"grant_types":                 // toStateGrantTypes
		"token_endpoint_authn_method": framework.EnumOkToTF(apiObject.GetTokenEndpointAuthMethodOk()),
	}

	oidcObject, d := types.ObjectValue(accessControlGroupOptionsTFObjectTypes, oidcOptions)
	objOIDCOptions = append(objOIDCOptions, oidcObject)
	returnVar, d := types.ListValue(types.ObjectType{AttrTypes: oidcOptionsTFObjectTypes}, objOIDCOptions)

	diags.Append(d...)

	return returnVar, diags
}
