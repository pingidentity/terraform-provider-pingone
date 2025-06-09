// Copyright © 2025 Ping Identity Corporation

package sso

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/legacysdk"
)

// Types
type ResourceSecretDataSource serviceClientType

type ResourceSecretDataSourceModel struct {
	EnvironmentId pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	ResourceId    pingonetypes.ResourceIDValue `tfsdk:"resource_id"`
	Previous      types.Object                 `tfsdk:"previous"`
	Secret        types.String                 `tfsdk:"secret"`
}

// Framework interfaces
var (
	_ datasource.DataSource = &ResourceSecretDataSource{}
)

// New Object
func NewResourceSecretDataSource() datasource.DataSource {
	return &ResourceSecretDataSource{}
}

// Metadata
func (r *ResourceSecretDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_resource_secret"
}

// Schema
func (r *ResourceSecretDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Datasource to retrieve the currently active secret, and the active previous secret for a PingOne resource in an environment.",

		Attributes: map[string]schema.Attribute{
			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("PingOne environment identifier (UUID) in which the resource exists."),
			),

			"resource_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("PingOne resource identifier (UUID) for which to retrieve the resource secret.  The resource must be an OpenID Connect type."),
			),

			"previous": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("An object that specifies the previous secret, when it expires, and when it was last used.").Description,
				Computed:    true,

				Attributes: map[string]schema.Attribute{
					"secret": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the previous resource secret. This property is returned in the response if the previous secret is not expired.").Description,
						Computed:    true,
						Sensitive:   true,
					},

					"expires_at": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A timestamp that specifies how long this secret is saved (and can be used) before it expires. Supported time range is 1 minute to 30 days.").Description,
						Computed:    true,
					},

					"last_used": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A timestamp that specifies when the previous secret was last used.").Description,
						Computed:    true,
					},
				},
			},

			"secret": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The resource secret ID used to authenticate to the authorization server. The secret has a minimum length of 64 characters per SHA-512 requirements when using the HS512 algorithm to sign ID tokens using the secret as the key.").Description,
				Computed:    true,
				Sensitive:   true,
			},
		},
	}
}

func (r *ResourceSecretDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *ResourceSecretDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *ResourceSecretDataSourceModel

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

	var secretResponse *management.ResourceSecret
	resp.Diagnostics.Append(legacysdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.ResourceClientSecretApi.ReadResourceSecret(ctx, data.EnvironmentId.ValueString(), data.ResourceId.ValueString()).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"ReadResourceSecret",
		legacysdk.DefaultCustomError,
		resourceOIDCSecretDataSourceRetryConditions,
		&secretResponse,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(secretResponse)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (p *ResourceSecretDataSourceModel) toState(apiObject *management.ResourceSecret) diag.Diagnostics {
	var diags, d diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.Previous, d = resourceSecretPreviousOkToTF(apiObject.GetPreviousOk())
	diags.Append(d...)

	p.Secret = framework.StringOkToTF(apiObject.GetSecretOk())

	return diags
}
