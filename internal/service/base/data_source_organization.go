package base

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
	"github.com/patrickcping/pingone-go-sdk-v2/pingone"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
)

// Types
type OrganizationDataSource struct {
	Client *pingone.Client
	region model.RegionMapping
}

type OrganizationDataSourceModel struct {
	Id                   types.String `tfsdk:"id"`
	OrganizationId       types.String `tfsdk:"organization_id"`
	Name                 types.String `tfsdk:"name"`
	Description          types.String `tfsdk:"description"`
	Type                 types.String `tfsdk:"type"`
	BillingConnectionIds types.Set    `tfsdk:"billing_connection_ids"`
	BaseUrlAPI           types.String `tfsdk:"base_url_api"`
	BaseUrlAuth          types.String `tfsdk:"base_url_auth"`
	BaseUrlOrchestrate   types.String `tfsdk:"base_url_orchestrate"`
	BaseUrlAgreementMgmt types.String `tfsdk:"base_url_agreement_management"`
	BaseUrlConsole       types.String `tfsdk:"base_url_console"`
	BaseUrlApps          types.String `tfsdk:"base_url_apps"`
}

// Framework interfaces
var (
	_ datasource.DataSource = &OrganizationDataSource{}
)

// New Object
func NewOrganizationDataSource() datasource.DataSource {
	return &OrganizationDataSource{}
}

// Metadata
func (r *OrganizationDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization"
}

// Schema
func (r *OrganizationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	typeDescripton := framework.SchemaAttributeDescriptionFromMarkdown(
		"The organization type. If the organization has any paid licenses, the type property value is set to `PAID`. Otherwise, the property value is set to `TRIAL`.  Internal organizations have a property value of `INTERNAL`.",
	)

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Datasource to retrieve details of the active PingOne organization.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"organization_id": schema.StringAttribute{
				Description: "The ID of the organization to retrieve.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("name")),
				},
			},

			"name": schema.StringAttribute{
				Description: "The name of the organization to retrieve.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("organization_id")),
				},
			},

			"description": schema.StringAttribute{
				Description: "The description of the organization.",
				Computed:    true,
			},

			"type": schema.StringAttribute{
				Description:         typeDescripton.Description,
				MarkdownDescription: typeDescripton.MarkdownDescription,
				Computed:            true,
			},

			"billing_connection_ids": schema.SetAttribute{
				Description: "The list of the BillingConnection resource IDs for the organization.",
				ElementType: types.StringType,
				Computed:    true,
			},

			"base_url_api": schema.StringAttribute{
				Description:        "**Deprecation message**.  This attribute is deprecated and will be removed in a future release.  Please review published modules for the PingOne provider on the Terraform Registry to gain equivalent functionality.  Helper attribute that provides an indication of the hostname of the API endpoint.  This attribute does not update if a non-production PingOne organization is used, nor if a custom domain is configured in any environment.",
				Computed:           true,
				DeprecationMessage: "This attribute is deprecated and will be removed in a future release.  Please review published modules for the PingOne provider on the Terraform Registry to gain equivalent functionality.",
			},

			"base_url_auth": schema.StringAttribute{
				Description:        "**Deprecation message**.  This attribute is deprecated and will be removed in a future release.  Please review published modules for the PingOne provider on the Terraform Registry to gain equivalent functionality.  Helper attribute that provides an indication of the hostname of the Authentication endpoint.  This attribute does not update if a non-production PingOne organization is used, nor if a custom domain is configured in any environment.",
				Computed:           true,
				DeprecationMessage: "This attribute is deprecated and will be removed in a future release.  Please review published modules for the PingOne provider on the Terraform Registry to gain equivalent functionality.",
			},

			"base_url_orchestrate": schema.StringAttribute{
				Description:        "**Deprecation message**.  This attribute is deprecated and will be removed in a future release.  Please review published modules for the PingOne provider on the Terraform Registry to gain equivalent functionality.  Helper attribute that provides an indication of the hostname of the Orchestration endpoint.  This attribute does not update if a non-production PingOne organization is used, nor if a custom domain is configured in any environment.",
				Computed:           true,
				DeprecationMessage: "This attribute is deprecated and will be removed in a future release.  Please review published modules for the PingOne provider on the Terraform Registry to gain equivalent functionality.",
			},

			"base_url_agreement_management": schema.StringAttribute{
				Description:        "**Deprecation message**.  This attribute is deprecated and will be removed in a future release.  Please review published modules for the PingOne provider on the Terraform Registry to gain equivalent functionality.  Helper attribute that provides an indication of the hostname of the Agreement Management endpoint.  This attribute does not update if a non-production PingOne organization is used, nor if a custom domain is configured in any environment.",
				Computed:           true,
				DeprecationMessage: "This attribute is deprecated and will be removed in a future release.  Please review published modules for the PingOne provider on the Terraform Registry to gain equivalent functionality.",
			},

			"base_url_console": schema.StringAttribute{
				Description:        "**Deprecation message**.  This attribute is deprecated and will be removed in a future release.  Please review published modules for the PingOne provider on the Terraform Registry to gain equivalent functionality.  Helper attribute that provides an indication of the hostname of the Console endpoint.  This attribute does not update if a non-production PingOne organization is used, nor if a custom domain is configured in any environment.",
				Computed:           true,
				DeprecationMessage: "This attribute is deprecated and will be removed in a future release.  Please review published modules for the PingOne provider on the Terraform Registry to gain equivalent functionality.",
			},

			"base_url_apps": schema.StringAttribute{
				Description:        "**Deprecation message**.  This attribute is deprecated and will be removed in a future release.  Please review published modules for the PingOne provider on the Terraform Registry to gain equivalent functionality.  Helper attribute that provides an indication of the hostname of the Applications endpoint.  This attribute does not update if a non-production PingOne organization is used, nor if a custom domain is configured in any environment.",
				Computed:           true,
				DeprecationMessage: "This attribute is deprecated and will be removed in a future release.  Please review published modules for the PingOne provider on the Terraform Registry to gain equivalent functionality.",
			},
		},
	}
}

func (r *OrganizationDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
	r.region = resourceConfig.Client.API.Region
}

func (r *OrganizationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *OrganizationDataSourceModel

	if r.Client.ManagementAPIClient == nil {
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

	var organization management.Organization

	if !data.Name.IsNull() {

		// Run the API call
		var entityArray *management.EntityArray
		resp.Diagnostics.Append(framework.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				return r.Client.ManagementAPIClient.OrganizationsApi.ReadAllOrganizations(ctx).Execute()
			},
			"ReadAllOrganizations",
			framework.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
			&entityArray,
		)...)
		if resp.Diagnostics.HasError() {
			return
		}

		if organizations, ok := entityArray.Embedded.GetOrganizationsOk(); ok {

			found := false
			for _, organizationItem := range organizations {

				if organizationItem.GetName() == data.Name.ValueString() {
					organization = organizationItem
					found = true
					break
				}
			}

			if !found {
				resp.Diagnostics.AddError(
					"Cannot find organization from name",
					fmt.Sprintf("The organization %s cannot be found", data.Name.String()),
				)
				return
			}

		}

	} else if !data.OrganizationId.IsNull() {

		// Run the API call
		var response *management.Organization
		resp.Diagnostics.Append(framework.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				return r.Client.ManagementAPIClient.OrganizationsApi.ReadOneOrganization(ctx, data.OrganizationId.ValueString()).Execute()
			},
			"ReadOneOrganization",
			framework.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
			&response,
		)...)
		if resp.Diagnostics.HasError() {
			return
		}

		organization = *response
	} else {
		resp.Diagnostics.AddError(
			"Missing parameter",
			"Cannot find the requested organization. organization_id or name must be set.",
		)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(&organization, r.region)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (p *OrganizationDataSourceModel) toState(v *management.Organization, region model.RegionMapping) diag.Diagnostics {
	var diags diag.Diagnostics

	if v == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.Id = framework.StringOkToTF(v.GetIdOk())
	p.OrganizationId = framework.StringOkToTF(v.GetIdOk())
	p.Name = framework.StringOkToTF(v.GetNameOk())
	p.Description = framework.StringOkToTF(v.GetDescriptionOk())
	p.Type = organzationTypeEnumOkToTF(v.GetTypeOk())
	p.BillingConnectionIds = organizationBillingConnectionIdsOkToTF(v.GetBillingConnectionsOk())

	p.BaseUrlAPI = types.StringValue(fmt.Sprintf("api.pingone.%s", region.URLSuffix))
	p.BaseUrlAuth = types.StringValue(fmt.Sprintf("auth.pingone.%s", region.URLSuffix))
	p.BaseUrlOrchestrate = types.StringValue(fmt.Sprintf("orchestrate-api.pingone.%s", region.URLSuffix))
	p.BaseUrlAgreementMgmt = types.StringValue(fmt.Sprintf("agreement-mgmt.pingone.%s", region.URLSuffix))
	p.BaseUrlConsole = types.StringValue(fmt.Sprintf("console.pingone.%s", region.URLSuffix))
	p.BaseUrlApps = types.StringValue(fmt.Sprintf("apps.pingone.%s", region.URLSuffix))

	return diags
}

func organzationTypeEnumOkToTF(v *management.EnumOrganizationType, ok bool) basetypes.StringValue {
	if !ok || v == nil {
		return types.StringNull()
	} else {
		return types.StringValue(string(*v))
	}
}

func organizationBillingConnectionIdsOkToTF(v []management.OrganizationBillingConnectionsInner, ok bool) basetypes.SetValue {
	if !ok || v == nil {
		return types.SetNull(types.StringType)
	} else {

		list := make([]attr.Value, 0)
		for _, item := range v {
			if i, ok := item.GetIdOk(); ok {
				list = append(list, types.StringValue(*i))
			}
		}

		return types.SetValueMust(types.StringType, list)
	}
}
