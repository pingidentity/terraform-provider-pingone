package base

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
)

// Types
type TrustedEmailDomainOwnershipDataSource serviceClientType

type TrustedEmailDomainOwnershipDataSourceModel struct {
	EnvironmentId        pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	TrustedEmailDomainId pingonetypes.ResourceIDValue `tfsdk:"trusted_email_domain_id"`
	Type                 types.String                 `tfsdk:"type"`
	Regions              types.Set                    `tfsdk:"regions"`
}

var (
	trustedEmailDomainOwnershipRegionTFObjectTypes = map[string]attr.Type{
		"name":   types.StringType,
		"status": types.StringType,
		"key":    types.StringType,
		"value":  types.StringType,
	}
)

// Framework interfaces
var (
	_ datasource.DataSource = &TrustedEmailDomainOwnershipDataSource{}
)

// New Object
func NewTrustedEmailDomainOwnershipDataSource() datasource.DataSource {
	return &TrustedEmailDomainOwnershipDataSource{}
}

// Metadata
func (r *TrustedEmailDomainOwnershipDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_trusted_email_domain_ownership"
}

// Schema
func (r *TrustedEmailDomainOwnershipDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	const minAttrLength = 1

	regionStatusDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The status of the email domain ownership.",
	).AllowedValuesEnum(management.AllowedEnumTrustedEmailStatusEnumValues)

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Datasource to retrieve the Trusted Email Domain ownership status for an environment.",

		Attributes: map[string]schema.Attribute{
			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment to retrieve trusted email domain ownership verification for."),
			),

			"trusted_email_domain_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the auto-generated ID of the email domain."),
			),

			"type": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the type of DNS record.").Description,
				Computed:    true,
			},

			"regions": schema.SetNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The regions collection specifies the properties for the 4 AWS SES regions that are used for sending email for the environment. The regions are determined by the geography where this environment was provisioned (North America, Canada, Europe & Asia-Pacific).").Description,
				Computed:    true,

				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the name of the region.").Description,
							Computed:    true,
						},

						"status": schema.StringAttribute{
							Description:         regionStatusDescription.Description,
							MarkdownDescription: regionStatusDescription.MarkdownDescription,
							Computed:            true,
						},

						"key": schema.StringAttribute{
							Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the record name to apply to the DNS provider.").Description,
							Computed:    true,
						},

						"value": schema.StringAttribute{
							Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the record value to apply to the DNS provider.").Description,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (r *TrustedEmailDomainOwnershipDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *TrustedEmailDomainOwnershipDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *TrustedEmailDomainOwnershipDataSourceModel

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

	var status *management.EmailDomainOwnershipStatus

	// Run the API call
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			return r.Client.ManagementAPIClient.TrustedEmailDomainsApi.ReadTrustedEmailDomainOwnershipStatus(ctx, data.EnvironmentId.ValueString(), data.TrustedEmailDomainId.ValueString()).Execute()
		},
		"ReadTrustedEmailDomainOwnershipStatus",
		framework.DefaultCustomError,
		nil,
		&status,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(status)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (p *TrustedEmailDomainOwnershipDataSourceModel) toState(v *management.EmailDomainOwnershipStatus) diag.Diagnostics {
	var diags, d diag.Diagnostics

	if v == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.Type = framework.StringOkToTF(v.GetTypeOk())

	p.Regions, d = toStateTrustedEmailDomainOwnershipRegionOkToTF(v.GetRegionsOk())
	diags.Append(d...)

	return diags
}

func toStateTrustedEmailDomainOwnershipRegionOkToTF(regions []management.EmailDomainOwnershipStatusRegionsInner, ok bool) (types.Set, diag.Diagnostics) {
	var diags diag.Diagnostics
	tfObjType := types.ObjectType{AttrTypes: trustedEmailDomainOwnershipRegionTFObjectTypes}

	if !ok {
		return types.SetNull(tfObjType), diags
	}

	flattenedList := []attr.Value{}
	for _, v := range regions {

		region := map[string]attr.Value{
			"name":   framework.StringOkToTF(v.GetNameOk()),
			"status": framework.EnumOkToTF(v.GetStatusOk()),
			"key":    framework.StringOkToTF(v.GetKeyOk()),
			"value":  framework.StringOkToTF(v.GetValueOk()),
		}

		flattenedObj, d := types.ObjectValue(trustedEmailDomainOwnershipRegionTFObjectTypes, region)
		diags.Append(d...)

		flattenedList = append(flattenedList, flattenedObj)
	}

	returnVar, d := types.SetValue(tfObjType, flattenedList)
	diags.Append(d...)

	return returnVar, diags
}
