// Copyright © 2025 Ping Identity Corporation

package base

import (
	"context"
	"fmt"
	"net/http"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/legacysdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
)

// Types
type CustomDomainSSLResource serviceClientType

type CustomDomainSSLResourceModel struct {
	Id                              pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId                   pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	CustomDomainId                  pingonetypes.ResourceIDValue `tfsdk:"custom_domain_id"`
	CerificatePemFile               types.String                 `tfsdk:"certificate_pem_file"`
	IntermediateCertificatesPemFile types.String                 `tfsdk:"intermediate_certificates_pem_file"`
	PrivateKeyPemFile               types.String                 `tfsdk:"private_key_pem_file"`
	DomainName                      types.String                 `tfsdk:"domain_name"`
	Status                          types.String                 `tfsdk:"status"`
	CertificateExpiresAt            timetypes.RFC3339            `tfsdk:"certificate_expires_at"`
}

// Framework interfaces
var (
	_ resource.Resource              = &CustomDomainSSLResource{}
	_ resource.ResourceWithConfigure = &CustomDomainSSLResource{}
)

// New Object
func NewCustomDomainSSLResource() resource.Resource {
	return &CustomDomainSSLResource{}
}

// Metadata
func (r *CustomDomainSSLResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_custom_domain_ssl"
}

// Schema
func (r *CustomDomainSSLResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	certificatePemFileDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the PEM-encoded certificate to import. The certificate must not be expired, must not be self signed and the domain must match one of the subject alternative name (SAN) values on the certificate.",
	).RequiresReplace()

	intermediateCertificatesPemFileDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the PEM-encoded certificate chain.",
	).RequiresReplace()

	privateKeyPemFileDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the PEM-encoded, unencrypted private key that matches the certificate's public key.",
	).RequiresReplace()

	statusDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the status of the custom domain.",
	).AllowedValuesEnum(management.AllowedEnumCustomDomainStatusEnumValues)

	const attrMinLength = 2

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: framework.SchemaDescriptionFromMarkdown("Resource to create and manage PingOne Custom Domain SSL settings.").Description,

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment to configure custom domain SSL in."),
			),

			"custom_domain_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the custom domain to set SSL settings for."),
			),

			"certificate_pem_file": schema.StringAttribute{
				Description:         certificatePemFileDescription.Description,
				MarkdownDescription: certificatePemFileDescription.MarkdownDescription,
				Required:            true,

				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},

				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
				},
			},

			"intermediate_certificates_pem_file": schema.StringAttribute{
				Description:         intermediateCertificatesPemFileDescription.Description,
				MarkdownDescription: intermediateCertificatesPemFileDescription.MarkdownDescription,
				Optional:            true,

				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},

				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
				},
			},

			"private_key_pem_file": schema.StringAttribute{
				Description:         privateKeyPemFileDescription.Description,
				MarkdownDescription: privateKeyPemFileDescription.MarkdownDescription,
				Required:            true,
				Sensitive:           true,

				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},

				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
				},
			},

			"domain_name": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the domain name in use.").Description,
				Computed:    true,
			},

			"status": schema.StringAttribute{
				MarkdownDescription: statusDescription.MarkdownDescription,
				Description:         statusDescription.Description,
				Computed:            true,
			},

			"certificate_expires_at": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The time when the certificate expires.  If this property is not present, it indicates that an SSL certificate has not been setup for this custom domain.").Description,
				Computed:    true,

				CustomType: timetypes.RFC3339Type{},
			},
		},
	}
}

func (r *CustomDomainSSLResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *CustomDomainSSLResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state CustomDomainSSLResourceModel

	if r.Client == nil || r.Client.ManagementAPIClient == nil {
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
	customDomainSSL := plan.expand()

	// Run the API call
	var response *management.CustomDomain
	resp.Diagnostics.Append(legacysdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.CustomDomainsApi.UpdateDomain(ctx, plan.EnvironmentId.ValueString(), plan.CustomDomainId.ValueString()).ContentType(management.ENUMCUSTOMDOMAINPOSTHEADER_CERTIFICATE_IMPORTJSON).CustomDomainCertificateRequest(*customDomainSSL).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdateDomain",
		func(_ *http.Response, p1Error *model.P1Error) diag.Diagnostics {
			var diags diag.Diagnostics

			if p1Error != nil {
				// Cannot validate against the authoritative name service
				if details, ok := p1Error.GetDetailsOk(); ok && details != nil && len(details) > 0 {
					m, _ := regexp.MatchString("^Custom domain status must be 'SSL_CERTIFICATE_REQUIRED' or 'ACTIVE' in order to import a certificate", details[0].GetMessage())
					if m {
						diags.AddError(
							fmt.Sprintf("Cannot add SSL certificate settings to the custom domain - %s", details[0].GetMessage()),
							`Please verify the domain first (hint: use the "pingone_custom_domain_verify" resource).)`,
						)

						return diags
					}
				}
			}
			return diags
		},
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

func (r *CustomDomainSSLResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *CustomDomainSSLResourceModel

	if r.Client == nil || r.Client.ManagementAPIClient == nil {
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
	resp.Diagnostics.Append(legacysdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.CustomDomainsApi.ReadOneDomain(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"ReadOneDomain",
		legacysdk.CustomErrorResourceNotFoundWarning,
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

func (r *CustomDomainSSLResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *CustomDomainSSLResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (p *CustomDomainSSLResourceModel) expand() *management.CustomDomainCertificateRequest {
	data := management.NewCustomDomainCertificateRequest(p.CerificatePemFile.ValueString(), p.PrivateKeyPemFile.ValueString())

	if !p.IntermediateCertificatesPemFile.IsNull() && !p.IntermediateCertificatesPemFile.IsUnknown() {
		data.SetIntermediateCertificates(p.IntermediateCertificatesPemFile.ValueString())
	}

	return data
}

func (p *CustomDomainSSLResourceModel) toState(apiObject *management.CustomDomain) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.Id = framework.PingOneResourceIDOkToTF(apiObject.GetIdOk())
	p.EnvironmentId = framework.PingOneResourceIDToTF(*apiObject.GetEnvironment().Id)
	p.DomainName = framework.StringOkToTF(apiObject.GetDomainNameOk())
	p.Status = framework.EnumOkToTF(apiObject.GetStatusOk())

	if v, ok := apiObject.GetCertificateOk(); ok {
		p.CertificateExpiresAt = framework.TimeOkToTF(v.GetExpiresAtOk())
	} else {
		p.CertificateExpiresAt = timetypes.NewRFC3339Null()
	}

	return diags
}
