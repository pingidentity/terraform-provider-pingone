// Copyright © 2025 Ping Identity Corporation

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
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/legacysdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
)

// Types
type CustomDomainVerifyResource serviceClientType

type CustomDomainVerifyResourceModel struct {
	Id             pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId  pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	CustomDomainId pingonetypes.ResourceIDValue `tfsdk:"custom_domain_id"`
	DomainName     types.String                 `tfsdk:"domain_name"`
	Status         types.String                 `tfsdk:"status"`
	Timeouts       timeouts.Value               `tfsdk:"timeouts"`
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

			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create:            true,
				CreateDescription: "A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as \"30s\" or \"2h45m\", as a time to wait for DNS record changes to propagate for validation. Valid time units are \"s\" (seconds), \"m\" (minutes), \"h\" (hours). The default is 60 minutes.",
			}),
		},
	}
}

func (r *CustomDomainVerifyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	resourceConfig, ok := req.ProviderData.(legacysdk.ResourceType)
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

func (r *CustomDomainVerifyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state CustomDomainVerifyResourceModel
	existingDomain := false

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

	defaultTimeout := 60

	timeout, d := plan.Timeouts.Create(ctx, time.Duration(defaultTimeout)*time.Minute)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *management.CustomDomain
	resp.Diagnostics.Append(legacysdk.ParseResponseWithCustomTimeout(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.CustomDomainsApi.UpdateDomain(ctx, plan.EnvironmentId.ValueString(), plan.CustomDomainId.ValueString()).ContentType(management.ENUMCUSTOMDOMAINPOSTHEADER_DOMAIN_NAME_VERIFYJSON).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdateDomain",
		func(_ *http.Response, p1Error *model.P1Error) diag.Diagnostics {
			var diags diag.Diagnostics

			if p1Error != nil {
				// Cannot validate against the authoritative name service
				if details, ok := p1Error.GetDetailsOk(); ok && details != nil && len(details) > 0 {
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

					m, _ = regexp.MatchString("^custom domain does not need to be verified", details[0].GetMessage())
					if m {
						diags.AddWarning(
							details[0].GetMessage(),
							"This is expected if the custom domain was verified previously.  The verification step has been skipped and the resource now tracks the domain's verified status.",
						)

						existingDomain = true
					}
				}
			}

			return diags
		},
		customDomainRetryConditions,
		&response,
		timeout,
	)...)

	if existingDomain {
		resp.Diagnostics.Append(legacysdk.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				fO, fR, fErr := r.Client.ManagementAPIClient.CustomDomainsApi.ReadOneDomain(ctx, plan.EnvironmentId.ValueString(), plan.CustomDomainId.ValueString()).Execute()
				return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
			},
			"ReadOneDomain",
			legacysdk.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
			&response,
		)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

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

	p.Id = framework.PingOneResourceIDOkToTF(apiObject.GetIdOk())
	p.EnvironmentId = framework.PingOneResourceIDToTF(*apiObject.GetEnvironment().Id)
	p.DomainName = framework.StringOkToTF(apiObject.GetDomainNameOk())
	p.Status = framework.EnumOkToTF(apiObject.GetStatusOk())

	return diags
}

func customDomainRetryConditions(ctx context.Context, r *http.Response, p1error *model.P1Error) bool {

	var err error

	if p1error != nil {

		// Permissions may not have propagated by this point
		if m, _ := regexp.MatchString("^The actor attempting to perform the request is not authorized.", p1error.GetMessage()); err == nil && m {
			tflog.Warn(ctx, "Insufficient PingOne privileges detected")
			return true
		}
		if err != nil {
			tflog.Warn(ctx, "Cannot match error string for retry")
			return false
		}

		// add retry time for DNS propegating
		if details, ok := p1error.GetDetailsOk(); ok && details != nil && len(details) > 0 {

			// perhaps it's the DNS authority
			if m, err := regexp.MatchString("^Error response from authoritative name servers: NXDOMAIN", details[0].GetMessage()); err == nil && m {
				tflog.Warn(ctx, fmt.Sprintf("Cannot verify the domain - %s.  Retrying...", details[0].GetMessage()))
				return true
			}
			if err != nil {
				tflog.Warn(ctx, "Cannot match error string for retry")
				return false
			}

			// perhaps it's the CNAME
			if m, err := regexp.MatchString("^No CNAME records found", details[0].GetMessage()); err == nil && m {
				tflog.Warn(ctx, fmt.Sprintf("Cannot verify the domain - %s.  Retrying...", details[0].GetMessage()))
				return true
			}
			if err != nil {
				tflog.Warn(ctx, "Cannot match error string for retry")
				return false
			}
		}

	}

	return false
}
