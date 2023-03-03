package base

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
)

// Types
type AgreementEnableResource struct {
	client *management.APIClient
	region model.RegionMapping
}

type AgreementEnableResourceModel struct {
	Id            types.String `tfsdk:"id"`
	EnvironmentId types.String `tfsdk:"environment_id"`
	AgreementId   types.String `tfsdk:"agreement_id"`
	Enabled       types.Bool   `tfsdk:"enabled"`
}

// Framework interfaces
var (
	_ resource.Resource                = &AgreementEnableResource{}
	_ resource.ResourceWithConfigure   = &AgreementEnableResource{}
	_ resource.ResourceWithImportState = &AgreementEnableResource{}
)

// New Object
func NewAgreementEnableResource() resource.Resource {
	return &AgreementEnableResource{}
}

// Metadata
func (r *AgreementEnableResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_agreement_enable"
}

// Schema.
func (r *AgreementEnableResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	const attrMinLength = 1

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage agreements in a PingOne environment.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(framework.SchemaDescription{
				Description: "The ID of the environment to associate the agreement with."},
			),

			"agreement_id": framework.Attr_LinkID(framework.SchemaDescription{
				Description: "The ID of the agreement to enable."},
			),

			"enabled": schema.BoolAttribute{
				Description: "A boolean that specifies the current enabled state of the agreement. The agreement must support the default language to be enabled. It cannot be disabled if it is referenced by a sign-on policy action. When an agreement is disabled, it is not used anywhere that it is configured across PingOne.",
				Required:    true,
			},
		},
	}
}

func (r *AgreementEnableResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *AgreementEnableResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state AgreementEnableResourceModel

	if r.client == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": r.region.URLSuffix,
	})

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	agreementResponse, diags := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return r.client.AgreementsResourcesApi.ReadOneAgreement(ctx, plan.EnvironmentId.ValueString(), plan.AgreementId.ValueString()).Execute()
		},
		"ReadOneAgreement",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the model for the API
	agreementEnable := plan.expand(agreementResponse.(*management.Agreement))

	// Run the API call
	response, d := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return r.client.AgreementsResourcesApi.UpdateAgreement(ctx, plan.EnvironmentId.ValueString(), plan.AgreementId.ValueString()).Agreement(*agreementEnable).Execute()
		},
		"UpdateAgreement",
		agreementEnableUpdateCustomErrorHandler,
		sdk.DefaultCreateReadRetryable,
	)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the state to save
	state = plan

	// Save updated data into Terraform state
	resp.Diagnostics.Append(state.toState(response.(*management.Agreement))...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *AgreementEnableResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *AgreementEnableResourceModel

	if r.client == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": r.region.URLSuffix,
	})

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	response, diags := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return r.client.AgreementsResourcesApi.ReadOneAgreement(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
		},
		"ReadOneAgreement",
		framework.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Remove from state if resource is not found
	if response == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(response.(*management.Agreement))...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AgreementEnableResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state AgreementEnableResourceModel

	if r.client == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": r.region.URLSuffix,
	})

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	agreementResponse, diags := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return r.client.AgreementsResourcesApi.ReadOneAgreement(ctx, plan.EnvironmentId.ValueString(), plan.AgreementId.ValueString()).Execute()
		},
		"ReadOneAgreement",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the model for the API
	agreementEnable := plan.expand(agreementResponse.(*management.Agreement))

	// Run the API call
	response, d := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return r.client.AgreementsResourcesApi.UpdateAgreement(ctx, plan.EnvironmentId.ValueString(), plan.Id.ValueString()).Agreement(*agreementEnable).Execute()
		},
		"UpdateAgreement",
		agreementEnableUpdateCustomErrorHandler,
		sdk.DefaultCreateReadRetryable,
	)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the state to save
	state = plan

	// Save updated data into Terraform state
	resp.Diagnostics.Append(state.toState(response.(*management.Agreement))...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *AgreementEnableResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *AgreementEnableResourceModel

	if r.client == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": r.region.URLSuffix,
	})

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	agreementResponse, diags := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return r.client.AgreementsResourcesApi.ReadOneAgreement(ctx, data.EnvironmentId.ValueString(), data.AgreementId.ValueString()).Execute()
		},
		"ReadOneAgreement",
		framework.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the model for the API
	data.Enabled = types.BoolValue(false)
	agreementDisable := data.expand(agreementResponse.(*management.Agreement))

	// Run the API call
	_, d := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return r.client.AgreementsResourcesApi.UpdateAgreement(ctx, data.EnvironmentId.ValueString(), data.AgreementId.ValueString()).Agreement(*agreementDisable).Execute()
		},
		"UpdateAgreement",
		framework.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
	)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *AgreementEnableResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	splitLength := 2
	attributes := strings.SplitN(req.ID, "/", splitLength)

	if len(attributes) != splitLength {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("invalid id (\"%s\") specified, should be in format \"environment_id/agreement_id\"", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("environment_id"), attributes[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), attributes[1])...)
}

func (p *AgreementEnableResourceModel) expand(existingObject *management.Agreement) *management.Agreement {

	data := management.NewAgreement(p.Enabled.ValueBool(), existingObject.GetName())

	if v, ok := existingObject.GetDescriptionOk(); ok {
		data.SetDescription(*v)
	}

	if v, ok := existingObject.GetReconsentPeriodDaysOk(); ok {
		data.SetReconsentPeriodDays(*v)
	}

	return data
}

func (p *AgreementEnableResourceModel) toState(apiObject *management.Agreement) diag.Diagnostics {
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
	p.AgreementId = framework.StringToTF(apiObject.GetId())
	p.Enabled = framework.BoolOkToTF(apiObject.GetEnabledOk())

	return diags
}

func agreementEnableUpdateCustomErrorHandler(error model.P1Error) diag.Diagnostics {
	var diags diag.Diagnostics

	if v, ok := error.GetDetailsOk(); ok && v != nil && len(v) > 0 {
		if v[0].GetCode() == "CONSTRAINT_VIOLATION" {
			if match, _ := regexp.MatchString("The agreement can not be enabled without supporting the default language configured for the environment.", v[0].GetMessage()); match {
				diags.AddError(
					v[0].GetMessage(),
					"The agreement must have an enabled agreement localization for the default language of the environment.  Ensure that a `pingone_agreement_localization`, `pingone_agreement_localization_revision` and `pingone_agreement_localization_enable` resource exist for the default langauge, or the environment's default language is re-configured using the `pingone_language_update` resource.",
				)

				return diags
			}
		}
	}

	return nil
}
