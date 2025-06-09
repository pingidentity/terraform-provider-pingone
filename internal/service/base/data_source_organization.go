// Copyright © 2025 Ping Identity Corporation

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
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/legacysdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
)

// Types
type OrganizationDataSource serviceClientType

type OrganizationDataSourceModel struct {
	Id                   pingonetypes.ResourceIDValue `tfsdk:"id"`
	OrganizationId       pingonetypes.ResourceIDValue `tfsdk:"organization_id"`
	Name                 types.String                 `tfsdk:"name"`
	Description          types.String                 `tfsdk:"description"`
	Type                 types.String                 `tfsdk:"type"`
	BillingConnectionIds types.Set                    `tfsdk:"billing_connection_ids"`
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

				CustomType: pingonetypes.ResourceIDType{},
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
}

func (r *OrganizationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *OrganizationDataSourceModel

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

	var organization *management.Organization

	if !data.Name.IsNull() {

		// Run the API call
		resp.Diagnostics.Append(legacysdk.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				pagedIterator := r.Client.ManagementAPIClient.OrganizationsApi.ReadAllOrganizations(ctx).Execute()

				var initialHttpResponse *http.Response

				for pageCursor, err := range pagedIterator {
					if err != nil {
						return nil, pageCursor.HTTPResponse, err
					}

					if initialHttpResponse == nil {
						initialHttpResponse = pageCursor.HTTPResponse
					}

					if organizations, ok := pageCursor.EntityArray.Embedded.GetOrganizationsOk(); ok {

						for _, organizationItem := range organizations {

							if organizationItem.GetName() == data.Name.ValueString() {
								return &organizationItem, pageCursor.HTTPResponse, nil
							}
						}
					}

				}

				return nil, initialHttpResponse, nil
			},
			"ReadAllOrganizations",
			legacysdk.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
			&organization,
		)...)
		if resp.Diagnostics.HasError() {
			return
		}

		if organization == nil {
			resp.Diagnostics.AddError(
				"Cannot find organization from name",
				fmt.Sprintf("The organization %s cannot be found", data.Name.String()),
			)
			return
		}

	} else if !data.OrganizationId.IsNull() {

		// Run the API call
		resp.Diagnostics.Append(legacysdk.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				return r.Client.ManagementAPIClient.OrganizationsApi.ReadOneOrganization(ctx, data.OrganizationId.ValueString()).Execute()
			},
			"ReadOneOrganization",
			legacysdk.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
			&organization,
		)...)
		if resp.Diagnostics.HasError() {
			return
		}

	} else {
		resp.Diagnostics.AddError(
			"Missing parameter",
			"Cannot find the requested organization. organization_id or name must be set.",
		)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(organization)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (p *OrganizationDataSourceModel) toState(v *management.Organization) diag.Diagnostics {
	var diags diag.Diagnostics

	if v == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.Id = framework.PingOneResourceIDOkToTF(v.GetIdOk())
	p.OrganizationId = framework.PingOneResourceIDOkToTF(v.GetIdOk())
	p.Name = framework.StringOkToTF(v.GetNameOk())
	p.Description = framework.StringOkToTF(v.GetDescriptionOk())
	p.Type = organzationTypeEnumOkToTF(v.GetTypeOk())
	p.BillingConnectionIds = organizationBillingConnectionIdsOkToTF(v.GetBillingConnectionsOk())

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
