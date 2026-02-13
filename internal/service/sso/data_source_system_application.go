// Copyright © 2026 Ping Identity Corporation

package sso

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/legacysdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/service"
)

// Types
type SystemApplicationDataSource serviceClientType

type systemApplicationDataSourceModel struct {
	Id                        pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId             pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	ApplicationId             pingonetypes.ResourceIDValue `tfsdk:"application_id"`
	Name                      types.String                 `tfsdk:"name"`
	Description               types.String                 `tfsdk:"description"`
	Type                      types.String                 `tfsdk:"type"`
	Protocol                  types.String                 `tfsdk:"protocol"`
	Enabled                   types.Bool                   `tfsdk:"enabled"`
	HiddenFromAppPortal       types.Bool                   `tfsdk:"hidden_from_app_portal"`
	Icon                      types.Object                 `tfsdk:"icon"`
	AccessControlRoleType     types.String                 `tfsdk:"access_control_role_type"`
	AccessControlGroupOptions types.Object                 `tfsdk:"access_control_group_options"`
	ClientId                  types.String                 `tfsdk:"client_id"`
	PkceEnforcement           types.String                 `tfsdk:"pkce_enforcement"`
	TokenEndpointAuthMethod   types.String                 `tfsdk:"token_endpoint_auth_method"`
	ApplyDefaultTheme         types.Bool                   `tfsdk:"apply_default_theme"`
	EnableDefaultThemeFooter  types.Bool                   `tfsdk:"enable_default_theme_footer"`
}

var (
	systemApplicationAccessControlGroupOptionsTFObjectTypes = map[string]attr.Type{
		"type":   types.StringType,
		"groups": types.SetType{ElemType: pingonetypes.ResourceIDType{}},
	}
)

// Framework interfaces
var (
	_ datasource.DataSource = &SystemApplicationDataSource{}
)

// New Object
func NewSystemApplicationDataSource() datasource.DataSource {
	return &SystemApplicationDataSource{}
}

// Metadata
func (r *SystemApplicationDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_system_application"
}

func (r *SystemApplicationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	// schema descriptions and validation settings
	const attrMinLength = 1

	applicationIdDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The identifier (UUID) of the system application.",
	).ExactlyOneOf([]string{"application_id", "name"}).AppendMarkdownString("Must be a valid PingOne resource ID.")

	nameDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The name of the system application.",
	).ExactlyOneOf([]string{"application_id", "name"})

	typeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the type of system application.",
	).AllowedValues(
		string(management.ENUMAPPLICATIONTYPE_PING_ONE_PORTAL),
		string(management.ENUMAPPLICATIONTYPE_PING_ONE_SELF_SERVICE),
	)

	accessControlRoleTypeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the user role required to access the application. A user is an admin user if they have one or more of the following roles assigned: `Organization Admin`, `Environment Admin`, `Identity Data Admin`, or `Client Application Developer`.",
	)

	accessControlGroupOptionsTypeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the group type required to access the application.",
	).AllowedValuesComplex(map[string]string{
		"ANY_GROUP":  "the actor must belong to at least one group listed in the `groups` property",
		"ALL_GROUPS": "the actor must belong to all groups listed in the `groups` property",
	})

	applyDefaultThemeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether to apply the default theme to the Self-Service or PingOne Portal application.",
	)

	enableDefaultThemeFooterDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether to show the default theme footer on the self-service application. Configurable only when the `type` is `PING_ONE_SELF_SERVICE` and `apply_default_theme` is also `true`.",
	)

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Data source to retrieve information about a built-in system application (PingOne Self-Service or PingOne Portal) in PingOne.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment that contains the system application."),
			),

			"application_id": schema.StringAttribute{
				Description:         applicationIdDescription.Description,
				MarkdownDescription: applicationIdDescription.MarkdownDescription,
				Optional:            true,

				CustomType: pingonetypes.ResourceIDType{},

				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(
						path.MatchRelative().AtParent().AtName("name"),
					),
				},
			},

			"name": schema.StringAttribute{
				Description:         nameDescription.Description,
				MarkdownDescription: nameDescription.MarkdownDescription,
				Optional:            true,

				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(
						path.MatchRelative().AtParent().AtName("application_id"),
					),
				},
			},
			"description": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the description of the application.").Description,
				Computed:    true,
			},
			"type": schema.StringAttribute{
				Description:         typeDescription.Description,
				MarkdownDescription: typeDescription.MarkdownDescription,
				Computed:            true,
			},

			"protocol": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the protocol used by the application.").Description,
				Computed:    true,
			},

			"enabled": schema.BoolAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A boolean that specifies the enabled/disabled status of the application.").Description,
				Computed:    true,
			},

			"hidden_from_app_portal": schema.BoolAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A boolean to specify whether the application is hidden in the application portal despite the configured group access policy.").Description,
				Computed:    true,
			},

			"icon": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The HREF and the ID for the application icon.").Description,
				Computed:    true,

				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("The ID for the application icon.").Description,
						Computed:    true,
					},
					"href": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("The HREF for the application icon.").Description,
						Computed:    true,
					},
				},
			},

			"access_control_role_type": schema.StringAttribute{
				Description:         accessControlRoleTypeDescription.Description,
				MarkdownDescription: accessControlRoleTypeDescription.MarkdownDescription,
				Computed:            true,
			},

			"access_control_group_options": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("Group access control settings.").Description,
				Computed:    true,

				Attributes: map[string]schema.Attribute{
					"groups": schema.SetAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A set that specifies the group IDs for the groups the actor must belong to for access to the application.").Description,
						Computed:    true,

						ElementType: pingonetypes.ResourceIDType{},
					},

					"type": schema.StringAttribute{
						Description:         accessControlGroupOptionsTypeDescription.Description,
						MarkdownDescription: accessControlGroupOptionsTypeDescription.MarkdownDescription,
						Computed:            true,
					},
				},
			},

			"client_id": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the application ID used to authenticate to the authorization server.").Description,
				Computed:    true,
			},

			"pkce_enforcement": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies how PKCE request parameters are handled on the authorize request.").Description,
				Computed:    true,
			},

			"token_endpoint_auth_method": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the client authentication methods supported by the token endpoint.").Description,
				Computed:    true,
			},

			"apply_default_theme": schema.BoolAttribute{
				Description:         applyDefaultThemeDescription.Description,
				MarkdownDescription: applyDefaultThemeDescription.MarkdownDescription,
				Computed:            true,
			},

			"enable_default_theme_footer": schema.BoolAttribute{
				Description:         enableDefaultThemeFooterDescription.Description,
				MarkdownDescription: enableDefaultThemeFooterDescription.MarkdownDescription,
				Computed:            true,
			},
		},
	}
}

func (r *SystemApplicationDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	resourceConfig, ok := req.ProviderData.(legacysdk.ResourceType)
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

func (r *SystemApplicationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *systemApplicationDataSourceModel

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

	var application *management.ReadOneApplication200Response

	// Application API does not support SCIM filtering
	if !data.ApplicationId.IsNull() {
		// Run the API call
		resp.Diagnostics.Append(legacysdk.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				fO, fR, fErr := r.Client.ManagementAPIClient.ApplicationsApi.ReadOneApplication(ctx, data.EnvironmentId.ValueString(), data.ApplicationId.ValueString()).Execute()
				return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
			},
			"ReadOneApplication",
			legacysdk.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
			&application,
		)...)
		if resp.Diagnostics.HasError() {
			return
		}

	} else if !data.Name.IsNull() {
		// Run the API call
		resp.Diagnostics.Append(legacysdk.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				pagedIterator := r.Client.ManagementAPIClient.ApplicationsApi.ReadAllApplications(ctx, data.EnvironmentId.ValueString()).Execute()

				var initialHttpResponse *http.Response

				for pageCursor, err := range pagedIterator {
					if err != nil {
						return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), nil, pageCursor.HTTPResponse, err)
					}

					if initialHttpResponse == nil {
						initialHttpResponse = pageCursor.HTTPResponse
					}

					if applications, ok := pageCursor.EntityArray.Embedded.GetApplicationsOk(); ok {

						var applicationObj management.ReadOneApplication200Response
						for _, applicationObj = range applications {
							applicationInstance := applicationObj.GetActualInstance()

							applicationName := ""

							switch v := applicationInstance.(type) {
							case *management.ApplicationPingOnePortal:
								applicationName = v.GetName()

							case *management.ApplicationPingOneSelfService:
								applicationName = v.GetName()
							}

							if applicationName != "" && strings.EqualFold(applicationName, data.Name.ValueString()) {
								return &applicationObj, pageCursor.HTTPResponse, nil
							}
						}
					}
				}

				return nil, initialHttpResponse, nil
			},
			"ReadAllApplications",
			legacysdk.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
			&application,
		)...)
		if resp.Diagnostics.HasError() {
			return
		}

		if application == nil {
			resp.Diagnostics.AddError(
				"Cannot find the system application from name",
				fmt.Sprintf("The system application name %s for environment %s cannot be found. Only system application types (%s, %s) are retrievable by this data source.",
					data.Name.String(),
					data.EnvironmentId.String(),
					string(management.ENUMAPPLICATIONTYPE_PING_ONE_PORTAL),
					string(management.ENUMAPPLICATIONTYPE_PING_ONE_SELF_SERVICE)),
			)
			return
		}

	} else {
		resp.Diagnostics.AddError(
			"Missing parameter",
			"Cannot find the requested PingOne System Application: application_id or name argument must be set.",
		)
		return
	}

	// Validate that the application is a system application
	if !isSystemApplication(application) {
		resp.Diagnostics.AddError(
			"Application is not a system application",
			fmt.Sprintf("The requested application is not a system application. Only system application types (%s, %s) are retrievable by this data source.",
				string(management.ENUMAPPLICATIONTYPE_PING_ONE_PORTAL),
				string(management.ENUMAPPLICATIONTYPE_PING_ONE_SELF_SERVICE)),
		)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(application)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (p *systemApplicationDataSourceModel) toState(apiObject *management.ReadOneApplication200Response) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	apiObjectCommon := management.Application{}

	if apiObject.ApplicationPingOnePortal != nil {
		apiObjectCommon = management.Application{
			Id:                  apiObject.ApplicationPingOnePortal.Id,
			Name:                apiObject.ApplicationPingOnePortal.Name,
			Environment:         apiObject.ApplicationPingOnePortal.Environment,
			Type:                apiObject.ApplicationPingOnePortal.Type,
			Protocol:            apiObject.ApplicationPingOnePortal.Protocol,
			Enabled:             apiObject.ApplicationPingOnePortal.Enabled,
			Description:         apiObject.ApplicationPingOnePortal.Description,
			HiddenFromAppPortal: apiObject.ApplicationPingOnePortal.HiddenFromAppPortal,
			Icon:                apiObject.ApplicationPingOnePortal.Icon,
			AccessControl:       apiObject.ApplicationPingOnePortal.AccessControl,
		}

		p.PkceEnforcement = framework.EnumOkToTF(apiObject.ApplicationPingOnePortal.GetPkceEnforcementOk())
		p.TokenEndpointAuthMethod = framework.EnumOkToTF(apiObject.ApplicationPingOnePortal.GetTokenEndpointAuthMethodOk())
		p.ApplyDefaultTheme = framework.BoolOkToTF(apiObject.ApplicationPingOnePortal.GetApplyDefaultThemeOk())
		p.EnableDefaultThemeFooter = types.BoolNull()
	}

	if apiObject.ApplicationPingOneSelfService != nil {
		apiObjectCommon = management.Application{
			Id:                  apiObject.ApplicationPingOneSelfService.Id,
			Name:                apiObject.ApplicationPingOneSelfService.Name,
			Environment:         apiObject.ApplicationPingOneSelfService.Environment,
			Type:                apiObject.ApplicationPingOneSelfService.Type,
			Protocol:            apiObject.ApplicationPingOneSelfService.Protocol,
			Enabled:             apiObject.ApplicationPingOneSelfService.Enabled,
			Description:         apiObject.ApplicationPingOneSelfService.Description,
			HiddenFromAppPortal: apiObject.ApplicationPingOneSelfService.HiddenFromAppPortal,
			Icon:                apiObject.ApplicationPingOneSelfService.Icon,
			AccessControl:       apiObject.ApplicationPingOneSelfService.AccessControl,
		}

		p.PkceEnforcement = framework.EnumOkToTF(apiObject.ApplicationPingOneSelfService.GetPkceEnforcementOk())
		p.TokenEndpointAuthMethod = framework.EnumOkToTF(apiObject.ApplicationPingOneSelfService.GetTokenEndpointAuthMethodOk())
		p.ApplyDefaultTheme = framework.BoolOkToTF(apiObject.ApplicationPingOneSelfService.GetApplyDefaultThemeOk())
		p.EnableDefaultThemeFooter = framework.BoolOkToTF(apiObject.ApplicationPingOneSelfService.GetEnableDefaultThemeFooterOk())
	}

	p.Id = framework.PingOneResourceIDToTF(apiObjectCommon.GetId())
	p.EnvironmentId = framework.PingOneResourceIDToTF(*apiObjectCommon.GetEnvironment().Id)
	p.Type = framework.EnumOkToTF(apiObjectCommon.GetTypeOk())
	p.Protocol = framework.EnumOkToTF(apiObjectCommon.GetProtocolOk())
	p.Name = framework.StringOkToTF(apiObjectCommon.GetNameOk())
	p.Description = framework.StringOkToTF(apiObjectCommon.GetDescriptionOk())
	p.Enabled = framework.BoolOkToTF(apiObjectCommon.GetEnabledOk())
	p.HiddenFromAppPortal = framework.BoolOkToTF(apiObjectCommon.GetHiddenFromAppPortalOk())

	var d diag.Diagnostics
	p.Icon, d = service.ImageOkToTF(apiObjectCommon.GetIconOk())
	diags.Append(d...)

	// Client ID is the same as the application ID for system applications
	p.ClientId = framework.StringOkToTF(apiObjectCommon.GetIdOk())
	if v, ok := apiObjectCommon.GetAccessControlOk(); ok {
		if v1, ok := v.GetRoleOk(); ok {
			p.AccessControlRoleType = framework.EnumOkToTF(v1.GetTypeOk())
		}

		if v1, ok := v.GetGroupOk(); ok {

			if v2, ok := v1.GetGroupsOk(); ok {

				groupsSlice := make([]string, 0)

				for _, group := range v2 {
					groupsSlice = append(groupsSlice, group.GetId())
				}

				tfGroupsSlice := framework.PingOneResourceIDSetToTF(groupsSlice)

				objValue, d := types.ObjectValue(systemApplicationAccessControlGroupOptionsTFObjectTypes, map[string]attr.Value{
					"groups": tfGroupsSlice,
					"type":   framework.EnumOkToTF(v1.GetTypeOk()),
				})
				diags.Append(d...)

				p.AccessControlGroupOptions = objValue
			}
		}
	}

	return diags
}

// isSystemApplication checks if the response contains a system application type
func isSystemApplication(application *management.ReadOneApplication200Response) bool {
	if application == nil {
		return false
	}

	return application.ApplicationPingOnePortal != nil ||
		application.ApplicationPingOneSelfService != nil
}
