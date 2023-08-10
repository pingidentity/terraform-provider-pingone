package base

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
)

// Types
type CustomDomainVerifyResource struct {
	client *management.APIClient
	region model.RegionMapping
}

type CustomDomainVerifyResourceModel struct {
	Id             types.String   `tfsdk:"id"`
	EnvironmentId  types.String   `tfsdk:"environment_id"`
	CustomDomainId types.String   `tfsdk:"custom_domain_id"`
	DomainName     types.String   `tfsdk:"domain_name"`
	Status         types.String   `tfsdk:"status"`
	Timeouts       timeouts.Value `tfsdk:"timeouts"`
}

// Framework interfaces
var (
	_ resource.Resource              = &CustomDomainVerifyResource{}
	_ resource.ResourceWithConfigure = &CustomDomainVerifyResource{}
)

// New Object
func NewCustomDomainVerifyResource() resource.Resource {
	return &CustomDomainVerifyResource{}
}

// Metadata
func (r *CustomDomainVerifyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_custom_domain_verify"
}

// Schema
func (r *CustomDomainVerifyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	statusDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the status of the custom domain.",
	).AllowedValuesEnum(management.AllowedEnumCustomDomainStatusEnumValues)

	const attrMinLength = 2

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: framework.SchemaDescriptionFromMarkdown("Resource to create and manage PingOne Custom Domain verification.").Description,

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment to verify the custom domain in."),
			),

			"custom_domain_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the custom domain to verify."),
			),

			"domain_name": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the domain name in use.").Description,
				Computed:    true,
			},

			"status": schema.StringAttribute{
				MarkdownDescription: statusDescription.MarkdownDescription,
				Description:         statusDescription.Description,
				Computed:            true,
			},
		},

		Blocks: map[string]schema.Block{
			"timeouts": timeouts.Block(ctx, timeouts.Opts{
				Create: true,
			}),
		},
	}
}

func (r *CustomDomainVerifyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *CustomDomainVerifyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state CustomDomainVerifyResourceModel

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

	timeoutValue := 60

	// Run the API call
	var response *management.CustomDomain
	resp.Diagnostics.Append(framework.ParseResponseWithCustomTimeout(
		ctx,

		func() (any, *http.Response, error) {
			return r.client.CustomDomainsApi.UpdateDomain(ctx, plan.EnvironmentId.ValueString(), plan.CustomDomainId.ValueString()).ContentType(management.ENUMCUSTOMDOMAINPOSTHEADER_DOMAIN_NAME_VERIFYJSON).Execute()
		},
		"UpdateDomain",
		func(error model.P1Error) diag.Diagnostics {
			var diags diag.Diagnostics

			// Cannot validate against the authoritative name service
			if details, ok := error.GetDetailsOk(); ok && details != nil && len(details) > 0 {
				m, _ := regexp.MatchString("^Error response from authoritative name servers: NXDOMAIN", details[0].GetMessage())
				if m {
					diags.AddError(
						fmt.Sprintf("Cannot verify the domain - %s", details[0].GetMessage()),
						"Please check the domain authority exists or is reachable.",
					)

					return diags
				}

				m, _ = regexp.MatchString("^No CNAME records found", details[0].GetMessage())
				if m {
					diags.AddError(
						fmt.Sprintf("Cannot verify the domain - %s", details[0].GetMessage()),
						"Please check the domain authority has the correct CNAME value set (hint: if using the \"pingone_custom_domain\" resource, the CNAME value to use is returned in the \"canonical_name\" attribute.)",
					)

					return diags
				}
			}

			return nil
		},
		sdk.DefaultCreateReadRetryable,
		&response,
		time.Duration(timeoutValue)*time.Minute, // 60 mins
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

func (r *CustomDomainVerifyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *CustomDomainVerifyResourceModel

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
	var response *management.CustomDomain
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			return r.client.CustomDomainsApi.ReadOneDomain(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
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

func (r *CustomDomainVerifyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *CustomDomainVerifyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (p *CustomDomainVerifyResourceModel) toState(apiObject *management.CustomDomain) diag.Diagnostics {
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

	return diags
}
