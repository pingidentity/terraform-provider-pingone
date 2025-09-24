// Copyright © 2025 Ping Identity Corporation

package sso

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
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
)

// Types
type ResourceDataSource serviceClientType

type ResourceDataSourceModel struct {
	Id                             pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId                  pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	ResourceId                     pingonetypes.ResourceIDValue `tfsdk:"resource_id"`
	Name                           types.String                 `tfsdk:"name"`
	Description                    types.String                 `tfsdk:"description"`
	Type                           types.String                 `tfsdk:"type"`
	Audience                       types.String                 `tfsdk:"audience"`
	AccessTokenValiditySeconds     types.Int32                  `tfsdk:"access_token_validity_seconds"`
	ApplicationPermissionsSettings types.Object                 `tfsdk:"application_permissions_settings"`
	IntrospectEndpointAuthMethod   types.String                 `tfsdk:"introspect_endpoint_auth_method"`
}

// Framework interfaces
var (
	_ datasource.DataSource = &ResourceDataSource{}
)

// New Object
func NewResourceDataSource() datasource.DataSource {
	return &ResourceDataSource{}
}

// Metadata
func (r *ResourceDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_resource"
}

// Schema
func (r *ResourceDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	nameLength := 1

	resourceIdDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The ID of the resource.",
	).ExactlyOneOf([]string{"resource_id", "name"}).AppendMarkdownString("Must be a valid PingOne resource ID.")

	nameDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The name of the resource.",
	).ExactlyOneOf([]string{"resource_id", "name"})

	typeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the type of resource.",
	).AllowedValuesComplex(map[string]string{
		string(management.ENUMRESOURCETYPE_OPENID_CONNECT): "specifies the built-in platform resource for OpenID Connect",
		string(management.ENUMRESOURCETYPE_PINGONE_API):    "specifies the built-in platform resource for PingOne",
		string(management.ENUMRESOURCETYPE_CUSTOM):         "specifies the a resource that has been created by admin",
	})

	audienceDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies a URL without a fragment or `@ObjectName` and must not contain `pingone` or `pingidentity` (for example, `https://api.myresource.com`). If a URL is not specified, the resource name is used.",
	)

	introspectEndpointAuthMethodDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The client authentication methods supported by the token endpoint",
	).AllowedValuesEnum(management.AllowedEnumResourceIntrospectEndpointAuthMethodEnumValues)

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Datasource to read PingOne OAuth 2.0 resource data.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment that is configured with the resource."),
			),

			"resource_id": schema.StringAttribute{
				Description:         resourceIdDescription.Description,
				MarkdownDescription: resourceIdDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,

				CustomType: pingonetypes.ResourceIDType{},

				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("name")),
				},
			},

			"name": schema.StringAttribute{
				Description:         nameDescription.Description,
				MarkdownDescription: nameDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("resource_id")),
					stringvalidator.LengthAtLeast(nameLength),
				},
			},

			"description": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A description of the resource.").Description,
				Computed:    true,
			},

			"type": schema.StringAttribute{
				Description:         typeDescription.Description,
				MarkdownDescription: typeDescription.MarkdownDescription,
				Computed:            true,
			},

			"audience": schema.StringAttribute{
				Description:         audienceDescription.Description,
				MarkdownDescription: audienceDescription.MarkdownDescription,
				Computed:            true,
			},

			"access_token_validity_seconds": schema.Int32Attribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("An integer that specifies the number of seconds that the access token is valid.").Description,
				Computed:    true,
			},

			"introspect_endpoint_auth_method": schema.StringAttribute{
				Description:         introspectEndpointAuthMethodDescription.Description,
				MarkdownDescription: introspectEndpointAuthMethodDescription.MarkdownDescription,
				Computed:            true,
			},

			"application_permissions_settings": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that specifies whether application permissions are added to access tokens generated by PingOne.").Description,
				Computed:    true,

				Attributes: map[string]schema.Attribute{
					"claim_enabled": schema.BoolAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A boolean setting to enable application permission claims in the access token.").Description,
						Computed:    true,
					},
				},
			},
		},
	}
}

func (r *ResourceDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *ResourceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *ResourceDataSourceModel

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

	var resource *management.Resource

	if !data.Name.IsNull() {

		var d diag.Diagnostics
		resource, d = fetchResourceFromName(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), data.Name.ValueString(), false)
		resp.Diagnostics.Append(d...)

	} else if !data.ResourceId.IsNull() {

		var d diag.Diagnostics
		resource, d = fetchResourceFromID(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), data.ResourceId.ValueString(), false)
		resp.Diagnostics.Append(d...)

	} else {
		resp.Diagnostics.AddError(
			"Missing parameter",
			"Cannot find the requested resource. resource_id or name must be set.",
		)
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(resource)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (p *ResourceDataSourceModel) toState(apiObject *management.Resource) diag.Diagnostics {
	var diags, d diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.Id = framework.PingOneResourceIDToTF(apiObject.GetId())
	p.ResourceId = framework.PingOneResourceIDToTF(apiObject.GetId())
	p.Name = framework.StringOkToTF(apiObject.GetNameOk())
	p.Description = framework.StringOkToTF(apiObject.GetDescriptionOk())
	p.Type = framework.EnumOkToTF(apiObject.GetTypeOk())
	p.Audience = framework.StringOkToTF(apiObject.GetAudienceOk())
	p.AccessTokenValiditySeconds = framework.Int32OkToTF(apiObject.GetAccessTokenValiditySecondsOk())

	p.ApplicationPermissionsSettings, d = resourceApplicationPermissionsSettingsOk(apiObject.GetApplicationPermissionsSettingsOk())
	diags.Append(d...)

	p.IntrospectEndpointAuthMethod = framework.EnumOkToTF(apiObject.GetIntrospectEndpointAuthMethodOk())

	return diags
}
