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
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	"github.com/pingidentity/terraform-provider-pingone/internal/filter"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type UserDataSource struct {
	client *management.APIClient
	region model.RegionMapping
}

type UserDataSourceModel struct {
	Id                types.String `tfsdk:"id"`
	UserId            types.String `tfsdk:"user_id"`
	EnvironmentId     types.String `tfsdk:"environment_id"`
	Username          types.String `tfsdk:"username"`
	Email             types.String `tfsdk:"email"`
	EmailVerified     types.Bool   `tfsdk:"email_verified"`
	Status            types.String `tfsdk:"status"`
	Enabled           types.Bool   `tfsdk:"enabled"`
	PopulationId      types.String `tfsdk:"population_id"`
	Account           types.Object `tfsdk:"account"`
	Address           types.Object `tfsdk:"address"`
	ExternalId        types.String `tfsdk:"external_id"`
	IdentityProvider  types.Object `tfsdk:"identity_provider"`
	Lifecycle         types.Object `tfsdk:"user_lifecycle"`
	Locale            types.String `tfsdk:"locale"`
	MFAEnabled        types.Bool   `tfsdk:"mfa_enabled"`
	MobilePhone       types.String `tfsdk:"mobile_phone"`
	Name              types.Object `tfsdk:"name"`
	Nickname          types.String `tfsdk:"nickname"`
	Password          types.Object `tfsdk:"password"`
	Photo             types.Object `tfsdk:"photo"`
	PreferredLanguage types.String `tfsdk:"preferred_language"`
	PrimaryPhone      types.String `tfsdk:"primary_phone"`
	Timezone          types.String `tfsdk:"timezone"`
	Title             types.String `tfsdk:"title"`
	Type              types.String `tfsdk:"type"`
	VerifyStatus      types.String `tfsdk:"verify_status"`
}

var (
	userPasswordTFDSObjectTypes = map[string]attr.Type{
		"external": types.ObjectType{
			AttrTypes: userPasswordExternalTFObjectTypes,
		},
	}

	userLifecycleTFDSObjectTypes = map[string]attr.Type{
		"status": types.StringType,
	}
)

// Framework interfaces
var (
	_ datasource.DataSource = &UserDataSource{}
)

// New Object
func NewUserDataSource() datasource.DataSource {
	return &UserDataSource{}
}

// Metadata
func (r *UserDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

// Schema
func (r *UserDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	const attrMinLength = 1
	exactlyOneOfDSPaths := []string{"user_id", "username", "email"}

	userIdDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the ID of the user.  Must be a valid PingOne resource ID.",
	).ExactlyOneOf(exactlyOneOfDSPaths)

	usernameDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the user name, which is unique within an environment. The `username` must either be a well-formed email address or a string. The string can contain any letters, numbers, combining characters, math and currency symbols, dingbats and drawing characters, and invisible whitespace",
	).ExactlyOneOf(exactlyOneOfDSPaths)

	emailDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the user's email address. For more information about email address formatting, see section 3.4 of [RFC 2822, Internet Message Format](http://www.faqs.org/rfcs/rfc2822.html).",
	).ExactlyOneOf(exactlyOneOfDSPaths)

	statusDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"**Deprecation notice**: This attribute is deprecated and will be removed in a future release. Please use the `enabled` attribute instead.  The enabled status of the user.",
	).AllowedValues("ENABLED", "DISABLED")

	enabledDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether the user is enabled. This attribute is set to `true` by default when a user is created.",
	)

	accountCanAuthenticateDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether the user can authenticate. If the value is set to `false`, the account is locked or the user is disabled, and unless specified otherwise in administrative configuration, the user will be unable to authenticate.",
	)

	accountStatusDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the the account locked state.",
	).AllowedValuesEnum(management.AllowedEnumUserStatusEnumValues)

	addressCountryCodeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the country name component in [ISO 3166-1](https://www.iso.org/iso-3166-country-codes.html) \"alpha-2\" code format. For example, the country codes for the United States and Sweden are `US` and `SE`, respectively.",
	)

	identityProviderTypeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the type of identity provider used to authenticate the user.",
	).AllowedValuesEnum(management.AllowedEnumIdentityProviderEnumValues).AppendMarkdownString(
		"The default value of `PING_ONE` is set when a value for `id` was not provided when the user was originally created.",
	)

	userLifecycleDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that specifies the user's identity lifecycle information.",
	)

	userLifecycleStatusDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the status of the account lifecycle.",
	).AllowedValuesEnum(management.AllowedEnumUserLifecycleStatusEnumValues)

	localeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the user's default location as a valid language tag as defined in [RFC 5646](https://www.rfc-editor.org/rfc/rfc5646.html). The following are example tags: `fr`, `en-US`, `es-419`, `az-Arab`, `man-Nkoo-GN`. This is used for purposes of localizing such items as currency, date time format, or numerical representations.",
	)

	mfaEnabledDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether multi-factor authentication is enabled. This attribute is set to `false` by default when the user is created.",
	)

	mobilePhoneDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the user's native phone number. This might also match the `primary_phone` attribute.",
	)

	nameFamilyDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the family name of the user, or Last in most Western languages (for example, `Jensen` given the full name `Ms. Barbara J Jensen, III`).",
	)

	nameFormattedDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the fully formatted name of the user (for example `Ms. Barbara J Jensen, III`).",
	)

	nameGivenDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the given name of the user, or First in most Western languages (for example, `Barbara` given the full name `Ms. Barbara J Jensen, III`).",
	)

	nameHonorificPrefixDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the honorific prefix(es) of the user, or title in most Western languages (for example, `Ms.` given the full name `Ms. Barbara Jane Jensen, III`).",
	)

	nameHonorificSuffixDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the honorific suffix(es) of the user, or suffix in most Western languages (for example, `III` given the full name `Ms. Barbara Jane Jensen, III`).",
	)

	nameMiddleDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the middle name(s) of the user (for exmple, `Jane` given the full name `Ms. Barbara Jane Jensen, III`).",
	)

	nicknameDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the user's nickname.",
	)

	passwordExternalGatewayTypeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that indicates one of the supported gateway types.",
	).AllowedValuesEnum(management.AllowedEnumGatewayTypeEnumValues)

	photoHrefDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The URI that is a uniform resource locator (as defined in [Section 1.1.3 of RFC 3986](https://www.rfc-editor.org/rfc/rfc3986#section-1.3)) that points to a resource location representing the user's image.",
	)

	preferredLanguageDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the user's preferred written or spoken languages, as a valid language range that is the same as the HTTP `Accept-Language` header field (not including `Accept-Language:` prefix) and is specified in [Section 5.3.5 of RFC 7231](https://datatracker.ietf.org/doc/html/rfc7231#section-5.3.5). For example: `en-US`, `en-gb;q=0.8`, `en;q=0.7`.",
	)

	primaryPhoneDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the user's primary phone number. This might also match the `mobile_phone` attribute.",
	)

	timezoneDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the user's time zone, conforming with the IANA Time Zone database format [RFC 6557](https://www.rfc-editor.org/rfc/rfc6557.html), also known as the \"Olson\" time zone database format [Olson-TZ](https://www.iana.org/time-zones). For example, `America/Los_Angeles`.",
	)

	titleDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the user's title, such as `Vice President`.",
	)

	typeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the user's type.",
	)

	verifyStatusDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that indicates whether ID verification can be done for the user.",
	).AllowedValuesEnum(management.AllowedEnumUserVerifyStatusEnumValues).AppendMarkdownString(
		"If the user verification status is `DISABLED`, a new verification status cannot be created for that user until the status is changed to `ENABLED`.",
	)

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Data source to read a single PingOne user's data in an environment for a given username, email address or user ID.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment that the user has been created in."),
			),

			"user_id": schema.StringAttribute{
				Description:         userIdDescription.Description,
				MarkdownDescription: userIdDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,

				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(
						path.MatchRoot("user_id"),
						path.MatchRoot("username"),
						path.MatchRoot("email"),
					),
					verify.P1ResourceIDValidator(),
				},
			},

			"username": schema.StringAttribute{
				Description:         usernameDescription.Description,
				MarkdownDescription: usernameDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,

				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(
						path.MatchRoot("user_id"),
						path.MatchRoot("username"),
						path.MatchRoot("email"),
					),
					stringvalidator.LengthAtLeast(attrMinLength),
				},
			},

			"email": schema.StringAttribute{
				Description:         emailDescription.Description,
				MarkdownDescription: emailDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,

				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(
						path.MatchRoot("user_id"),
						path.MatchRoot("username"),
						path.MatchRoot("email"),
					),
					stringvalidator.LengthAtLeast(attrMinLength),
				},
			},

			"email_verified": schema.BoolAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown(
					"A boolean that specifies whether the user's email is verified.",
				).Description,
				Computed: true,
			},

			"status": schema.StringAttribute{
				Description:         statusDescription.Description,
				MarkdownDescription: statusDescription.MarkdownDescription,
				Computed:            true,
				DeprecationMessage:  "This attribute is deprecated and will be removed in a future release. Please use the `enabled` attribute instead.",
			},

			"enabled": schema.BoolAttribute{
				Description:         enabledDescription.Description,
				MarkdownDescription: enabledDescription.MarkdownDescription,
				Computed:            true,
			},

			"population_id": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown(
					"A PingOne resource identifier of the population resource associated with the user.",
				).Description,
				Computed: true,
			},

			"account": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown(
					"A single object that specifies the user's account information.",
				).Description,
				Computed: true,

				Attributes: map[string]schema.Attribute{
					"can_authenticate": schema.BoolAttribute{
						Description:         accountCanAuthenticateDescription.Description,
						MarkdownDescription: accountCanAuthenticateDescription.MarkdownDescription,
						Computed:            true,
					},

					"locked_at": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown(
							"The time the specified user account was locked. This property might be absent if the account is unlocked or if the account was locked out automatically by failed password attempts.",
						).Description,
						Computed: true,
					},

					"status": schema.StringAttribute{
						Description:         accountStatusDescription.Description,
						MarkdownDescription: accountStatusDescription.MarkdownDescription,
						Computed:            true,
					},
				},
			},

			"address": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown(
					"A single object that specifies the user's address information.",
				).Description,
				Computed: true,

				Attributes: map[string]schema.Attribute{
					"country_code": schema.StringAttribute{
						Description:         addressCountryCodeDescription.Description,
						MarkdownDescription: addressCountryCodeDescription.MarkdownDescription,
						Computed:            true,
					},

					"locality": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown(
							"A string that specifies the city or locality component of the address.",
						).Description,
						Computed: true,
					},

					"postal_code": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown(
							"A string that specifies the ZIP code or postal code component of the address.",
						).Description,
						Computed: true,
					},

					"region": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown(
							"A string that specifies the state, province, or region component of the address.",
						).Description,
						Computed: true,
					},

					"street_address": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown(
							"A string that specifies the full street address component, which may include house number, street name, P.O. box, and multi-line extended street address information.",
						).Description,
						Computed: true,
					},
				},
			},

			"external_id": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown(
					"A string that specifies an identifier for the user resource as defined by the provisioning client. The external id attribute simplifies the correlation of the user in PingOne with the user's account in another system of record. The platform does not use this attribute directly in any way, but it is used by Ping Identity's Data Sync product.",
				).Description,
				Computed: true,
			},

			"identity_provider": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown(
					"A single object that specifies the user's identity provider information.",
				).Description,
				Computed: true,

				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown(
							"A string that identifies the external identity provider used to authenticate the user. If not provided, PingOne is the identity provider. This attribute is required if the identity provider is authoritative for just-in-time user provisioning.",
						).Description,
						Computed: true,
					},

					"type": schema.StringAttribute{
						Description:         identityProviderTypeDescription.Description,
						MarkdownDescription: identityProviderTypeDescription.MarkdownDescription,
						Computed:            true,
					},
				},
			},

			"user_lifecycle": schema.SingleNestedAttribute{
				Description:         userLifecycleDescription.Description,
				MarkdownDescription: userLifecycleDescription.MarkdownDescription,
				Computed:            true,

				Attributes: map[string]schema.Attribute{
					"status": schema.StringAttribute{
						Description:         userLifecycleStatusDescription.Description,
						MarkdownDescription: userLifecycleStatusDescription.MarkdownDescription,
						Computed:            true,
					},
				},
			},

			"locale": schema.StringAttribute{
				Description:         localeDescription.Description,
				MarkdownDescription: localeDescription.MarkdownDescription,
				Computed:            true,
			},

			"mfa_enabled": schema.BoolAttribute{
				Description:         mfaEnabledDescription.Description,
				MarkdownDescription: mfaEnabledDescription.MarkdownDescription,
				Computed:            true,
			},

			"mobile_phone": schema.StringAttribute{
				Description:         mobilePhoneDescription.Description,
				MarkdownDescription: mobilePhoneDescription.MarkdownDescription,
				Computed:            true,
			},

			"name": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown(
					"A single object that specifies the user's name information.",
				).Description,
				Computed: true,

				Attributes: map[string]schema.Attribute{
					"family": schema.StringAttribute{
						Description:         nameFamilyDescription.Description,
						MarkdownDescription: nameFamilyDescription.MarkdownDescription,
						Computed:            true,
					},

					"formatted": schema.StringAttribute{
						Description:         nameFormattedDescription.Description,
						MarkdownDescription: nameFormattedDescription.MarkdownDescription,
						Computed:            true,
					},

					"given": schema.StringAttribute{
						Description:         nameGivenDescription.Description,
						MarkdownDescription: nameGivenDescription.MarkdownDescription,
						Computed:            true,
					},

					"honorific_prefix": schema.StringAttribute{
						Description:         nameHonorificPrefixDescription.Description,
						MarkdownDescription: nameHonorificPrefixDescription.MarkdownDescription,
						Computed:            true,
					},

					"honorific_suffix": schema.StringAttribute{
						Description:         nameHonorificSuffixDescription.Description,
						MarkdownDescription: nameHonorificSuffixDescription.MarkdownDescription,
						Computed:            true,
					},

					"middle": schema.StringAttribute{
						Description:         nameMiddleDescription.Description,
						MarkdownDescription: nameMiddleDescription.MarkdownDescription,
						Computed:            true,
					},
				},
			},

			"nickname": schema.StringAttribute{
				Description:         nicknameDescription.Description,
				MarkdownDescription: nicknameDescription.MarkdownDescription,
				Computed:            true,
			},

			"password": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown(
					"A single object that specifies the user's password information.",
				).Description,
				Computed: true,

				Attributes: map[string]schema.Attribute{
					"external": schema.SingleNestedAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown(
							"A single object that maps the information relevant to the user's password, and its association to external directories.",
						).Description,
						Computed: true,

						Attributes: map[string]schema.Attribute{
							"gateway": schema.SingleNestedAttribute{
								Description: framework.SchemaAttributeDescriptionFromMarkdown(
									"A single object that contains the external gateway properties. When this is value is specified, the user's password is managed in an external directory.",
								).Description,
								Computed: true,

								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description: framework.SchemaAttributeDescriptionFromMarkdown(
											"A string that specifies the PingOne resource ID of the linked gateway that references the remote directory.",
										).Description,
										Computed: true,
									},

									"type": schema.StringAttribute{
										Description:         passwordExternalGatewayTypeDescription.Description,
										MarkdownDescription: passwordExternalGatewayTypeDescription.MarkdownDescription,
										Computed:            true,
									},

									"user_type_id": schema.StringAttribute{
										Description: framework.SchemaAttributeDescriptionFromMarkdown(
											"A string that specifies the PingOne resource ID of a user type in the list of user types for the LDAP gateway.",
										).Description,
										Computed: true,
									},

									"correlation_attributes": schema.MapAttribute{
										Description: framework.SchemaAttributeDescriptionFromMarkdown(
											"A string map that maps the external LDAP directory attributes to PingOne attributes. PingOne uses these values to read the attributes from the external LDAP directory and map them to the corresponding PingOne attributes.",
										).Description,
										Computed: true,

										ElementType: types.StringType,
									},
								},
							},
						},
					},
				},
			},

			"photo": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown(
					"A single object that describes the user's photo information.",
				).Description,
				Computed: true,

				Attributes: map[string]schema.Attribute{
					"href": schema.StringAttribute{
						Description:         photoHrefDescription.Description,
						MarkdownDescription: photoHrefDescription.MarkdownDescription,
						Computed:            true,
					},
				},
			},

			"preferred_language": schema.StringAttribute{
				Description:         preferredLanguageDescription.Description,
				MarkdownDescription: preferredLanguageDescription.MarkdownDescription,
				Computed:            true,
			},

			"primary_phone": schema.StringAttribute{
				Description:         primaryPhoneDescription.Description,
				MarkdownDescription: primaryPhoneDescription.MarkdownDescription,
				Computed:            true,
			},

			"timezone": schema.StringAttribute{
				Description:         timezoneDescription.Description,
				MarkdownDescription: timezoneDescription.MarkdownDescription,
				Computed:            true,
			},

			"title": schema.StringAttribute{
				Description:         titleDescription.Description,
				MarkdownDescription: titleDescription.MarkdownDescription,
				Computed:            true,
			},

			"type": schema.StringAttribute{
				Description:         typeDescription.Description,
				MarkdownDescription: typeDescription.MarkdownDescription,
				Computed:            true,
			},

			"verify_status": schema.StringAttribute{
				Description:         verifyStatusDescription.Description,
				MarkdownDescription: verifyStatusDescription.MarkdownDescription,
				Computed:            true,
			},
		},
	}
}

func (r *UserDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.client = preparedClient
	r.region = resourceConfig.Client.API.Region
}

func (r *UserDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *UserDataSourceModel

	if r.client == nil {
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

	var user management.User
	var scimFilter string

	if !data.Username.IsNull() {

		scimFilter = filter.BuildScimFilter(
			append(make([]interface{}, 0), map[string]interface{}{
				"name":   "username",
				"values": []string{data.Username.ValueString()},
			}), map[string]string{})

	} else if !data.UserId.IsNull() {

		scimFilter = filter.BuildScimFilter(
			append(make([]interface{}, 0), map[string]interface{}{
				"name":   "id",
				"values": []string{data.UserId.ValueString()},
			}), map[string]string{})

	} else if !data.Email.IsNull() {

		scimFilter = filter.BuildScimFilter(
			append(make([]interface{}, 0), map[string]interface{}{
				"name":   "email",
				"values": []string{data.Email.ValueString()},
			}), map[string]string{})

	} else {
		resp.Diagnostics.AddError(
			"Missing parameter",
			"Cannot find the requested user. user_id, username or email must be set.",
		)
		return
	}

	var response *management.EntityArray
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			return r.client.UsersApi.ReadAllUsers(ctx, data.EnvironmentId.ValueString()).Filter(scimFilter).Execute()
		},
		"ReadAllUsers",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
		&response,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var responseEnabled *management.UserEnabled
	if users, ok := response.Embedded.GetUsersOk(); ok && len(users) > 0 && users[0].Id != nil {

		user = users[0]

		resp.Diagnostics.Append(framework.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				return r.client.EnableUsersApi.ReadUserEnabled(ctx, data.EnvironmentId.ValueString(), user.GetId()).Execute()
			},
			"ReadUserEnabled",
			framework.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
			&responseEnabled,
		)...)

	} else {
		resp.Diagnostics.AddError(
			"Cannot find user",
			"Cannot find the requested user from the provided values. Please check the user_id, username or email parameters.",
		)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(&user, responseEnabled)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (p *UserDataSourceModel) toState(apiObject *management.User, apiObjectEnabled *management.UserEnabled) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil || apiObjectEnabled == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.Id = framework.StringOkToTF(apiObject.GetIdOk())
	p.UserId = framework.StringOkToTF(apiObject.GetIdOk())
	p.EnvironmentId = framework.StringOkToTF(apiObject.Environment.GetIdOk())
	p.Username = framework.StringOkToTF(apiObject.GetUsernameOk())
	p.Email = framework.StringOkToTF(apiObject.GetEmailOk())

	//p.EmailVerified = framework.BoolOkToTF(apiObject.GetEmailVerifiedOk())
	p.Enabled = framework.BoolOkToTF(apiObject.GetEnabledOk())

	// deprecated start
	if v, ok := apiObject.GetEnabledOk(); ok && *v {
		p.Status = framework.StringToTF("ENABLED")
	} else {
		p.Status = framework.StringToTF("DISABLED")
	}
	// deprecated end

	var d diag.Diagnostics

	p.PopulationId = framework.StringOkToTF(apiObject.Population.GetIdOk())
	p.Account, d = p.userAccountOkToTF(apiObject.GetAccountOk())
	diags = append(diags, d...)

	p.Address, d = p.userAddressOkToTF(apiObject.GetAddressOk())
	diags = append(diags, d...)

	p.ExternalId = framework.StringOkToTF(apiObject.GetExternalIdOk())
	p.IdentityProvider, d = p.userIdentityProviderOkToTF(apiObject.GetIdentityProviderOk())
	diags = append(diags, d...)

	p.Lifecycle, d = p.userLifecycleOkToTF(apiObject.GetLifecycleOk())
	diags = append(diags, d...)

	p.Locale = framework.StringOkToTF(apiObject.GetLocaleOk())
	p.MFAEnabled = framework.BoolOkToTF(apiObject.GetMfaEnabledOk())
	p.MobilePhone = framework.StringOkToTF(apiObject.GetMobilePhoneOk())
	p.Name, d = p.userNameOkToTF(apiObject.GetNameOk())
	diags = append(diags, d...)

	p.Nickname = framework.StringOkToTF(apiObject.GetNicknameOk())

	p.Password, d = p.userPasswordOkToTF(apiObject.GetPasswordOk())
	diags = append(diags, d...)

	p.Photo, d = p.photoOkToTF(apiObject.GetPhotoOk())
	diags = append(diags, d...)

	p.PreferredLanguage = framework.StringOkToTF(apiObject.GetPreferredLanguageOk())
	p.PrimaryPhone = framework.StringOkToTF(apiObject.GetPrimaryPhoneOk())
	p.Timezone = framework.StringOkToTF(apiObject.GetTimezoneOk())
	p.Title = framework.StringOkToTF(apiObject.GetTitleOk())
	p.Type = framework.StringOkToTF(apiObject.GetTypeOk())
	p.VerifyStatus = framework.EnumOkToTF(apiObject.GetVerifyStatusOk())

	return diags
}

func (p *UserDataSourceModel) userAccountOkToTF(apiObject *management.UserAccount, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(userAccountTFObjectTypes), diags
	}

	objMap := map[string]attr.Value{
		"can_authenticate": framework.BoolOkToTF(apiObject.GetCanAuthenticateOk()),
		"locked_at":        framework.TimeOkToTF(apiObject.GetLockedAtOk()),
		"status":           framework.EnumOkToTF(apiObject.GetStatusOk()),
	}

	objValue, d := types.ObjectValue(userAccountTFObjectTypes, objMap)
	diags.Append(d...)

	return objValue, diags
}

func (p *UserDataSourceModel) userAddressOkToTF(apiObject *management.UserAddress, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(userAddressTFObjectTypes), diags
	}

	objMap := map[string]attr.Value{
		"country_code":   framework.StringOkToTF(apiObject.GetCountryCodeOk()),
		"locality":       framework.StringOkToTF(apiObject.GetLocalityOk()),
		"postal_code":    framework.StringOkToTF(apiObject.GetPostalCodeOk()),
		"region":         framework.StringOkToTF(apiObject.GetRegionOk()),
		"street_address": framework.StringOkToTF(apiObject.GetStreetAddressOk()),
	}

	objValue, d := types.ObjectValue(userAddressTFObjectTypes, objMap)
	diags.Append(d...)

	return objValue, diags
}

func (p *UserDataSourceModel) userIdentityProviderOkToTF(apiObject *management.UserIdentityProvider, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(userIdentityProviderTFObjectTypes), diags
	}

	objMap := map[string]attr.Value{
		"id":   framework.StringOkToTF(apiObject.GetIdOk()),
		"type": framework.EnumOkToTF(apiObject.GetTypeOk()),
	}

	objValue, d := types.ObjectValue(userIdentityProviderTFObjectTypes, objMap)
	diags.Append(d...)

	return objValue, diags
}

func (p *UserDataSourceModel) userLifecycleOkToTF(apiObject *management.UserLifecycle, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(userLifecycleTFDSObjectTypes), diags
	}

	objMap := map[string]attr.Value{
		"status": framework.EnumOkToTF(apiObject.GetStatusOk()),
	}

	objValue, d := types.ObjectValue(userLifecycleTFDSObjectTypes, objMap)
	diags.Append(d...)

	return objValue, diags
}

func (p *UserDataSourceModel) userNameOkToTF(apiObject *management.UserName, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(userNameTFObjectTypes), diags
	}

	objMap := map[string]attr.Value{
		"family":           framework.StringOkToTF(apiObject.GetFamilyOk()),
		"formatted":        framework.StringOkToTF(apiObject.GetFormattedOk()),
		"given":            framework.StringOkToTF(apiObject.GetGivenOk()),
		"honorific_prefix": framework.StringOkToTF(apiObject.GetHonorificPrefixOk()),
		"honorific_suffix": framework.StringOkToTF(apiObject.GetHonorificSuffixOk()),
		"middle":           framework.StringOkToTF(apiObject.GetMiddleOk()),
	}

	objValue, d := types.ObjectValue(userNameTFObjectTypes, objMap)
	diags.Append(d...)

	return objValue, diags
}

func (p *UserDataSourceModel) userPasswordOkToTF(apiObject *management.UserPassword, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	// The API object might be nil even though the plan is not.  We need to fill the state from the plan if it is
	if !ok || apiObject == nil {
		return types.ObjectNull(userPasswordTFDSObjectTypes), diags
	}

	externalObject := types.ObjectNull(userPasswordExternalTFObjectTypes)

	objMap := map[string]attr.Value{
		"external": externalObject,
	}

	objValue, d := types.ObjectValue(userPasswordTFDSObjectTypes, objMap)
	diags.Append(d...)

	return objValue, diags
}

func (p *UserDataSourceModel) photoOkToTF(apiObject *management.UserPhoto, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(userPhotoTFObjectTypes), diags
	}

	objMap := map[string]attr.Value{
		"href": framework.StringOkToTF(apiObject.GetHrefOk()),
	}

	objValue, d := types.ObjectValue(userPhotoTFObjectTypes, objMap)
	diags.Append(d...)

	return objValue, diags
}
