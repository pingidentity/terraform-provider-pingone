package sso

import (
	"context"
	"fmt"
	"net/http"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type ResourceDataSource struct {
	client *management.APIClient
	region model.RegionMapping
}

type ResourceDataSourceModel struct {
	Id                           types.String `tfsdk:"id"`
	EnvironmentId                types.String `tfsdk:"environment_id"`
	ResourceId                   types.String `tfsdk:"resource_id"`
	Name                         types.String `tfsdk:"name"`
	Description                  types.String `tfsdk:"description"`
	Type                         types.String `tfsdk:"type"`
	Audience                     types.String `tfsdk:"audience"`
	AccessTokenValiditySeconds   types.Int64  `tfsdk:"access_token_validity_seconds"`
	IntrospectEndpointAuthMethod types.String `tfsdk:"introspect_endpoint_auth_method"`
	ClientSecret                 types.String `tfsdk:"client_secret"`
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
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("name")),
					verify.P1ResourceIDValidator(),
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
				Description: "A description of the resource.",
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

			"access_token_validity_seconds": schema.Int64Attribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("An integer that specifies the number of seconds that the access token is valid.").Description,
				Computed:    true,
			},

			"introspect_endpoint_auth_method": schema.StringAttribute{
				Description:         introspectEndpointAuthMethodDescription.Description,
				MarkdownDescription: introspectEndpointAuthMethodDescription.MarkdownDescription,
				Computed:            true,
			},

			"client_secret": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("An auto-generated resource client secret.").Description,
				Computed:    true,
				Sensitive:   true,
			},
		},
	}
}

func (r *ResourceDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	preparedClient, err := prepareClient(ctx, resourceConfig)
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

func (r *ResourceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *ResourceDataSourceModel

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

	var resource *management.Resource

	if !data.Name.IsNull() {

		var d diag.Diagnostics
		resource, d = fetchResourceFromName(ctx, r.client, data.EnvironmentId.ValueString(), data.Name.ValueString())
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}

	} else if !data.ResourceId.IsNull() {

		var d diag.Diagnostics
		resource, d = fetchResourceFromID(ctx, r.client, data.EnvironmentId.ValueString(), data.ResourceId.ValueString())
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}

	} else {
		resp.Diagnostics.AddError(
			"Missing parameter",
			"Cannot find the requested resource. resource_id or name must be set.",
		)
		return
	}

	var resourceClientSecret *management.ResourceSecret
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			return r.client.ResourceClientSecretApi.ReadResourceSecret(ctx, data.EnvironmentId.ValueString(), resource.GetId()).Execute()
		},
		"ReadResourceSecret",
		framework.CustomErrorResourceNotFoundWarning,
		func(ctx context.Context, r *http.Response, p1error *model.P1Error) bool {

			// The secret may take a short time to propagate
			if r.StatusCode == 404 && resource.GetType() == management.ENUMRESOURCETYPE_CUSTOM {
				tflog.Warn(ctx, "Resource secret not found, available for retry")
				return true
			}

			if p1error != nil {
				var err error

				// Permissions may not have propagated by this point
				if m, err := regexp.MatchString("^The actor attempting to perform the request is not authorized.", p1error.GetMessage()); err == nil && m {
					tflog.Warn(ctx, "Insufficient PingOne privileges detected")
					return true
				}
				if err != nil {
					tflog.Warn(ctx, "Cannot match error string for retry")
					return false
				}

			}

			return false
		},
		&resourceClientSecret,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(resource, resourceClientSecret)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (p *ResourceDataSourceModel) toState(apiObject *management.Resource, apiObjectSecret *management.ResourceSecret) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.Id = framework.StringToTF(apiObject.GetId())
	p.ResourceId = framework.StringToTF(apiObject.GetId())
	p.Name = framework.StringOkToTF(apiObject.GetNameOk())
	p.Description = framework.StringOkToTF(apiObject.GetDescriptionOk())
	p.Type = framework.EnumOkToTF(apiObject.GetTypeOk())
	p.Audience = framework.StringOkToTF(apiObject.GetAudienceOk())
	p.AccessTokenValiditySeconds = framework.Int32OkToTF(apiObject.GetAccessTokenValiditySecondsOk())
	p.IntrospectEndpointAuthMethod = framework.EnumOkToTF(apiObject.GetIntrospectEndpointAuthMethodOk())

	if apiObjectSecret == nil {
		p.ClientSecret = types.StringNull()
	} else {
		p.ClientSecret = framework.StringOkToTF(apiObjectSecret.GetSecretOk())
	}

	return diags
}

func fetchResourceFromID(ctx context.Context, apiClient *management.APIClient, environmentId, resourceId string) (*management.Resource, diag.Diagnostics) {
	var diags diag.Diagnostics

	var resource *management.Resource
	diags.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			return apiClient.ResourcesApi.ReadOneResource(ctx, environmentId, resourceId).Execute()
		},
		"ReadOneResource",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
		&resource,
	)...)
	if diags.HasError() {
		return nil, diags
	}

	return resource, diags
}

func fetchResourceFromName(ctx context.Context, apiClient *management.APIClient, environmentId string, resourceName string) (*management.Resource, diag.Diagnostics) {
	var diags diag.Diagnostics

	var resource management.Resource

	// Run the API call
	var entityArray *management.EntityArray
	diags.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			return apiClient.ResourcesApi.ReadAllResources(ctx, environmentId).Execute()
		},
		"ReadAllResources",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
		&entityArray,
	)...)
	if diags.HasError() {
		return nil, diags
	}

	if resources, ok := entityArray.Embedded.GetResourcesOk(); ok {

		found := false
		for _, resourceItem := range resources {

			if resourceItem.GetName() == resourceName {
				resource = resourceItem
				found = true
				break
			}
		}

		if !found {
			diags.AddError(
				"Cannot find resource from name",
				fmt.Sprintf("The resource %s for environment %s cannot be found", resourceName, environmentId),
			)
			return nil, diags
		}

	}

	return &resource, diags
}
