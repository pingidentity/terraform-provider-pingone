package base

import (
	"context"
	"fmt"
	"net/http"

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
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
)

// Types
type TrustedEmailDomainDataSource serviceClientType

type TrustedEmailDomainDataSourceModel struct {
	DomainName    types.String                 `tfsdk:"domain_name"`
	EmailDomainId pingonetypes.ResourceIDValue `tfsdk:"trusted_email_domain_id"`
	EnvironmentId pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	Id            pingonetypes.ResourceIDValue `tfsdk:"id"`
}

// Framework interfaces
var (
	_ datasource.DataSource = &TrustedEmailDomainDataSource{}
)

// New Object
func NewTrustedEmailDomainDataSource() datasource.DataSource {
	return &TrustedEmailDomainDataSource{}
}

// Metadata
func (r *TrustedEmailDomainDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_trusted_email_domain"
}

// Schema
func (r *TrustedEmailDomainDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	primaryLookupFieldExactlyOneOf := []string{"trusted_email_domain_id", "domain_name"}

	trustedEmailDomainIdDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The ID of the trusted email domain.  Must be a valid PingOne resource ID.",
	).ExactlyOneOf(primaryLookupFieldExactlyOneOf)

	domainDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the domain name to use, which must be provided and must be unique within an environment (for example, `demo.bxretail.org`).",
	).ExactlyOneOf(primaryLookupFieldExactlyOneOf)

	const emailAddressMaxLength = 5

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Datasource to retrieve a Trusted Email Domain.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment that is configured with the trusted email domain."),
			),

			"trusted_email_domain_id": schema.StringAttribute{
				MarkdownDescription: trustedEmailDomainIdDescription.MarkdownDescription,
				Description:         trustedEmailDomainIdDescription.Description,
				Optional:            true,

				CustomType: pingonetypes.ResourceIDType{},

				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("domain_name")),
				},
			},

			"domain_name": schema.StringAttribute{
				MarkdownDescription: domainDescription.MarkdownDescription,
				Description:         domainDescription.Description,
				Optional:            true,

				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("trusted_email_domain_id")),
				},
			},
		},
	}
}

func (r *TrustedEmailDomainDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *TrustedEmailDomainDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *TrustedEmailDomainDataSourceModel

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

	var emailDomain management.EmailDomain

	if !data.DomainName.IsNull() {

		// Run the API call
		var entityArray *management.EntityArray
		resp.Diagnostics.Append(framework.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				fO, fR, fErr := r.Client.ManagementAPIClient.TrustedEmailDomainsApi.ReadAllTrustedEmailDomains(ctx, data.EnvironmentId.ValueString()).Execute()
				return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
			},
			"ReadAllTrustedEmailDomains",
			framework.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
			&entityArray,
		)...)
		if resp.Diagnostics.HasError() {
			return
		}

		if emailDomains, ok := entityArray.Embedded.GetEmailDomainsOk(); ok {

			found := false
			for _, emailDomainItem := range emailDomains {

				if emailDomainItem.GetDomainName() == data.DomainName.ValueString() {
					emailDomain = emailDomainItem
					found = true
					break
				}
			}

			if !found {
				resp.Diagnostics.AddError(
					"Cannot find trusted email domain from domain_name",
					fmt.Sprintf("The trusted email domain %s for environment %s cannot be found", data.DomainName.String(), data.EnvironmentId.String()),
				)
				return
			}

		}

	} else if !data.EmailDomainId.IsNull() {

		// Run the API call
		var response *management.EmailDomain
		resp.Diagnostics.Append(framework.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				fO, fR, fErr := r.Client.ManagementAPIClient.TrustedEmailDomainsApi.ReadOneTrustedEmailDomain(ctx, data.EnvironmentId.ValueString(), data.EmailDomainId.ValueString()).Execute()
				return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
			},
			"ReadOneTrustedEmailDomain",
			framework.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
			&response,
		)...)
		if resp.Diagnostics.HasError() {
			return
		}

		emailDomain = *response
	} else {
		resp.Diagnostics.AddError(
			"Missing parameter",
			"Cannot find the requested trusted email domain. trusted_email_domain_id or domain_name must be set.",
		)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(&emailDomain)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (p *TrustedEmailDomainDataSourceModel) toState(v *management.EmailDomain) diag.Diagnostics {
	var diags diag.Diagnostics

	if v == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.Id = framework.PingOneResourceIDOkToTF(v.GetIdOk())

	if e, ok := v.GetEnvironmentOk(); ok {
		p.EnvironmentId = framework.PingOneResourceIDOkToTF(e.GetIdOk())
	} else {
		p.EnvironmentId = pingonetypes.NewResourceIDNull()
	}

	p.EmailDomainId = framework.PingOneResourceIDOkToTF(v.GetIdOk())
	p.DomainName = framework.StringOkToTF(v.GetDomainNameOk())

	return diags
}
