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
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
)

// Types
type TrustedEmailAddressResource struct {
	client *management.APIClient
	region model.RegionMapping
}

type TrustedEmailAddressResourceModel struct {
	EmailDomainId types.String `tfsdk:"email_domain_id"`
	EmailAddress  types.String `tfsdk:"email_address"`
	EnvironmentId types.String `tfsdk:"environment_id"`
	Id            types.String `tfsdk:"id"`
	Status        types.String `tfsdk:"status"`
}

// Framework interfaces
var (
	_ resource.Resource                = &TrustedEmailAddressResource{}
	_ resource.ResourceWithConfigure   = &TrustedEmailAddressResource{}
	_ resource.ResourceWithImportState = &TrustedEmailAddressResource{}
)

// New Object
func NewTrustedEmailAddressResource() resource.Resource {
	return &TrustedEmailAddressResource{}
}

// Metadata
func (r *TrustedEmailAddressResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_trusted_email_address"
}

// Schema
func (r *TrustedEmailAddressResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	resourceDescriptionFmt := "Resource to create and manage trusted email addresses in PingOne.  PingOne supports the ability to configure up to 10 trusted email addresses for an existing trusted email domain. See %s.  Once configured and if the email address has not been previously verified, a verification email is sent."
	providerDescription := framework.SchemaAttributeDescription{
		MarkdownDescription: fmt.Sprintf(resourceDescriptionFmt, "[Trusted email domains](https://apidocs.pingidentity.com/pingone/platform/v1/api/#trusted-email-domains)"),
		Description:         fmt.Sprintf(resourceDescriptionFmt, "Trusted email domains (https://apidocs.pingidentity.com/pingone/platform/v1/api/#trusted-email-domains)"),
	}

	emailAddressDescriptionFmt := "The trusted email address, for example %s."
	emailAddressDescription := framework.SchemaAttributeDescription{
		MarkdownDescription: fmt.Sprintf(emailAddressDescriptionFmt, "`john.smith@bxretail.org`"),
		Description:         fmt.Sprintf(emailAddressDescriptionFmt, "\"john.smith@bxretail.org\""),
	}

	statusDescriptionFmt := "The status of the trusted email address.  Possible values are %s."
	statusDescription := framework.SchemaAttributeDescription{
		MarkdownDescription: fmt.Sprintf(statusDescriptionFmt, "`ACTIVE` and `VERIFICATION_REQUIRED`"),
		Description:         fmt.Sprintf(statusDescriptionFmt, "\"ACTIVE\" and \"VERIFICATION_REQUIRED\""),
	}

	const emailAddressMaxLength = 5

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: providerDescription.MarkdownDescription,
		Description:         providerDescription.Description,

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment to associate the trusted email address with."),
			),

			"email_domain_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the email domain to associate the trusted email address with."),
			),

			"email_address": schema.StringAttribute{
				MarkdownDescription: emailAddressDescription.MarkdownDescription,
				Description:         emailAddressDescription.Description,
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(emailAddressMaxLength),
				},
			},

			"status": schema.StringAttribute{
				MarkdownDescription: statusDescription.MarkdownDescription,
				Description:         statusDescription.Description,
				Computed:            true,
			},
		},
	}
}

func (r *TrustedEmailAddressResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	preparedClient, err := PrepareClient(ctx, resourceConfig)
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

func (r *TrustedEmailAddressResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state TrustedEmailAddressResourceModel

	if r.client == nil {
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
	emailDomainTrustedEmail := plan.expand()

	// Run the API call
	var response *management.EmailDomainTrustedEmail
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			return r.client.TrustedEmailAddressesApi.CreateTrustedEmailAddress(ctx, plan.EnvironmentId.ValueString(), plan.EmailDomainId.ValueString()).EmailDomainTrustedEmail(*emailDomainTrustedEmail).Execute()
		},
		"CreateTrustedEmailAddress",
		trustedEmailAddressAPIErrors,
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

func (r *TrustedEmailAddressResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *TrustedEmailAddressResourceModel

	if r.client == nil {
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
	var response *management.EmailDomainTrustedEmail
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			return r.client.TrustedEmailAddressesApi.ReadOneTrustedEmailAddress(ctx, data.EnvironmentId.ValueString(), data.EmailDomainId.ValueString(), data.Id.ValueString()).Execute()
		},
		"ReadOneTrustedEmailAddress",
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

func (r *TrustedEmailAddressResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *TrustedEmailAddressResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *TrustedEmailAddressResourceModel

	if r.client == nil {
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
			r, err := r.client.TrustedEmailAddressesApi.DeleteTrustedEmailAddress(ctx, data.EnvironmentId.ValueString(), data.EmailDomainId.ValueString(), data.Id.ValueString()).Execute()
			return nil, r, err
		},
		"DeleteTrustedEmailAddress",
		framework.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
		nil,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *TrustedEmailAddressResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	splitLength := 3
	attributes := strings.SplitN(req.ID, "/", splitLength)

	if len(attributes) != splitLength {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("invalid id (\"%s\") specified, should be in format \"environment_id/email_domain_id/trusted_email_address_id\"", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("environment_id"), attributes[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("email_domain_id"), attributes[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), attributes[2])...)
}

func (p *TrustedEmailAddressResourceModel) expand() *management.EmailDomainTrustedEmail {
	data := management.NewEmailDomainTrustedEmail(p.EmailAddress.ValueString())

	return data
}

func (p *TrustedEmailAddressResourceModel) toState(apiObject *management.EmailDomainTrustedEmail) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.Id = framework.StringToTF(apiObject.GetId())
	p.EnvironmentId = framework.StringToTF(*apiObject.GetEnvironment().Id)
	p.EmailDomainId = framework.StringToTF(*apiObject.GetDomain().Id)
	p.EmailAddress = framework.StringOkToTF(apiObject.GetEmailAddressOk())

	if v, ok := apiObject.GetStatusOk(); ok {
		p.Status = framework.StringToTF(string(*v))
	} else {
		p.Status = types.StringNull()
	}

	return diags
}

func trustedEmailAddressAPIErrors(error model.P1Error) diag.Diagnostics {
	var diags diag.Diagnostics

	// Domain not verified
	if details, ok := error.GetDetailsOk(); ok && details != nil && len(details) > 0 {
		if code, ok := details[0].GetCodeOk(); ok && *code == "INVALID_VALUE" {
			if target, ok := details[0].GetTargetOk(); ok && *target == "trustedEmail" {
				diags.AddError(
					"The domain of the given email address is not verified",
					"Ensure that the domain of the given trusted email address has been verified first.  This can be configured with the `pingone_trusted_email_domain` resource.",
				)

				return diags
			}
		}
	}
	return nil
}
