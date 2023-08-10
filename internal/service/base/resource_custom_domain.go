package base

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type CustomDomainResource serviceClientType

type CustomDomainResourceModel struct {
	Id                   types.String `tfsdk:"id"`
	EnvironmentId        types.String `tfsdk:"environment_id"`
	DomainName           types.String `tfsdk:"domain_name"`
	Status               types.String `tfsdk:"status"`
	CanonicalName        types.String `tfsdk:"canonical_name"`
	CertificateExpiresAt types.String `tfsdk:"certificate_expires_at"`
}

// Framework interfaces
var (
	_ resource.Resource                = &CustomDomainResource{}
	_ resource.ResourceWithConfigure   = &CustomDomainResource{}
	_ resource.ResourceWithImportState = &CustomDomainResource{}
)

// New Object
func NewCustomDomainResource() resource.Resource {
	return &CustomDomainResource{}
}

// Metadata
func (r *CustomDomainResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_custom_domain"
}

// Schema
func (r *CustomDomainResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	domainNameDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the domain name to use, which must be provided and must be unique within an environment (for example, `demo.bxretail.org`).",
	)

	statusDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the status of the custom domain.",
	).AllowedValuesEnum(management.AllowedEnumCustomDomainStatusEnumValues)

	const attrMinLength = 2

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: framework.SchemaDescriptionFromMarkdown("Resource to create and manage PingOne Custom Domains.").Description,

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment to associate the custom domain with."),
			),

			"domain_name": schema.StringAttribute{
				MarkdownDescription: domainNameDescription.MarkdownDescription,
				Description:         domainNameDescription.Description,
				Required:            true,

				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},

				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
					stringvalidator.RegexMatches(verify.IsHostname, "Value must be a valid domain name"),
				},
			},

			"status": schema.StringAttribute{
				MarkdownDescription: statusDescription.MarkdownDescription,
				Description:         statusDescription.Description,
				Computed:            true,
			},

			"canonical_name": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the domain name that should be used as the value of the CNAME record in the customer's DNS.").Description,
				Computed:    true,
			},

			"certificate_expires_at": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The time when the certificate expires.  If this property is not present, it indicates that an SSL certificate has not been setup for this custom domain.").Description,
				Computed:    true,
			},
		},
	}
}

func (r *CustomDomainResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.Client = preparedClient
}

func (r *CustomDomainResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state CustomDomainResourceModel

	if r.Client == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the model for the API
	customDomain := plan.expand()

	// Run the API call
	var response *management.CustomDomain
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			return r.Client.CustomDomainsApi.CreateDomain(ctx, plan.EnvironmentId.ValueString()).CustomDomain(*customDomain).Execute()
		},
		"CreateDomain",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
		&response,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the state to save
	state = plan

	// Save updated data into Terraform state
	resp.Diagnostics.Append(state.toState(response)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *CustomDomainResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *CustomDomainResourceModel

	if r.Client == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *management.CustomDomain
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			return r.Client.CustomDomainsApi.ReadOneDomain(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
		},
		"ReadOneDomain",
		framework.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
		&response,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Remove from state if resource is not found
	if response == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(response)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CustomDomainResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *CustomDomainResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *CustomDomainResourceModel

	if r.Client == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			r, err := r.Client.CustomDomainsApi.DeleteDomain(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return nil, r, err
		},
		"DeleteDomain",
		framework.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
		nil,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *CustomDomainResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	splitLength := 2
	attributes := strings.SplitN(req.ID, "/", splitLength)

	if len(attributes) != splitLength {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("invalid id (\"%s\") specified, should be in format \"environment_id/custom_domain_id\"", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("environment_id"), attributes[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), attributes[1])...)
}

func (p *CustomDomainResourceModel) expand() *management.CustomDomain {
	data := management.NewCustomDomain(p.DomainName.ValueString())

	return data
}

func (p *CustomDomainResourceModel) toState(apiObject *management.CustomDomain) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.Id = framework.StringOkToTF(apiObject.GetIdOk())
	p.EnvironmentId = framework.StringToTF(*apiObject.GetEnvironment().Id)
	p.DomainName = framework.StringOkToTF(apiObject.GetDomainNameOk())
	p.Status = framework.EnumOkToTF(apiObject.GetStatusOk())
	p.CanonicalName = framework.StringOkToTF(apiObject.GetCanonicalNameOk())

	if v, ok := apiObject.GetCertificateOk(); ok {
		p.CertificateExpiresAt = framework.TimeOkToTF(v.GetExpiresAtOk())
	} else {
		p.CertificateExpiresAt = types.StringNull()
	}

	return diags
}
