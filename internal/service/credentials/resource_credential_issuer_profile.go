package credentials

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/patrickcping/pingone-go-sdk-v2/credentials"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
)

// Types
type CredentialIssuerProfileResource struct {
	client *credentials.APIClient
	region model.RegionMapping
}

type CredentialIssuerProfileResourceModel struct {
	Id                    types.String `tfsdk:"id"`
	EnvironmentId         types.String `tfsdk:"environment_id"`
	ApplicationInstanceId types.String `tfsdk:"application_instance_id"`
	CreatedAt             types.String `tfsdk:"created_at"`
	UpdatedAt             types.String `tfsdk:"updated_at"`
	Name                  types.String `tfsdk:"name"`
}

// Framework interfaces
var (
	_ resource.Resource                = &CredentialIssuerProfileResource{}
	_ resource.ResourceWithConfigure   = &CredentialIssuerProfileResource{}
	_ resource.ResourceWithImportState = &CredentialIssuerProfileResource{}
)

// New Object
func NewCredentialIssuerProfileResource() resource.Resource {
	return &CredentialIssuerProfileResource{}
}

// Metadata
func (r *CredentialIssuerProfileResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_credential_issuer_profile"
}

// Schema
func (r *CredentialIssuerProfileResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	// schema descriptions and validation settings
	const attrMinLength = 1
	const attrMaxLength = 256

	// schema definition
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to retrieve or update the credential issuer information.\n\n" +
			"A credential issuer profile, which enables issuance of credentials, is automatically created when the credential service is added to an environment. This resource is typically only required to update the credential issuer name.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment to create the credential issuer in."),
			),

			"application_instance_id": schema.StringAttribute{
				Description: "Identifier (UUID) of the application instance registered with the PingOne platform service. This enables the client to send messages to the service.",
				Computed:    true,
			},

			"created_at": schema.StringAttribute{
				Description: "Date and time the issuer profile was created.",
				Computed:    true,
			},

			"updated_at": schema.StringAttribute{
				Description: "Date and time the issuer profile was last updated.",
				Computed:    true,
			},

			"name": schema.StringAttribute{
				Description: "The name of the credential issuer. The name is included in the metadata of an issued verifiable credential.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(attrMinLength, attrMaxLength),
				},
			},
		},
	}
}

func (r *CredentialIssuerProfileResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *CredentialIssuerProfileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state CredentialIssuerProfileResourceModel

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

	// Historical:  Pre-EA and initial-EA environments required creation of the issuer profile. Environments created after 2023.05.01 no longer have this requirement.
	// On 'create' [adding to state], check to see if the profile exists, and if not, create it.  Otherwise, only update the profile, while still adding to TF state.
	timeoutValue := 5

	var readIssuerProfileResponse *credentials.CredentialIssuerProfile
	resp.Diagnostics.Append(framework.ParseResponseWithCustomTimeout(
		ctx,

		func() (any, *http.Response, error) {
			return r.client.CredentialIssuersApi.ReadCredentialIssuerProfile(ctx, plan.EnvironmentId.ValueString()).Execute()
		},
		"ReadCredentialIssuerProfile",
		framework.CustomErrorResourceNotFoundWarning,
		credentialIssuerRetryConditions,
		&readIssuerProfileResponse,
		time.Duration(timeoutValue)*time.Minute, // 5 mins
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the model for the Create API call
	CredentialIssuerProfile, d := plan.expand()
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Execute a Create or Update depending on existence of credential issuer profile
	var response *credentials.CredentialIssuerProfile
	if readIssuerProfileResponse == nil {
		// create the issuer profile
		resp.Diagnostics.Append(framework.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				return r.client.CredentialIssuersApi.CreateCredentialIssuerProfile(ctx, plan.EnvironmentId.ValueString()).CredentialIssuerProfile(*CredentialIssuerProfile).Execute()
			},
			"CreateCredentialIssuerProfile",
			framework.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
			&response,
		)...)
		if resp.Diagnostics.HasError() {
			return
		}
	} else {
		// update existing issuer profile
		resp.Diagnostics.Append(framework.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				return r.client.CredentialIssuersApi.UpdateCredentialIssuerProfile(ctx, plan.EnvironmentId.ValueString()).CredentialIssuerProfile(*CredentialIssuerProfile).Execute()
			},
			"UpdateCredentialIssuerProfile",
			framework.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
			&response,
		)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	// Create the state to save
	state = plan

	// Save updated data into Terraform state
	resp.Diagnostics.Append(state.toState(response)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *CredentialIssuerProfileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *CredentialIssuerProfileResourceModel

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
	timeoutValue := 5

	var response *credentials.CredentialIssuerProfile
	resp.Diagnostics.Append(framework.ParseResponseWithCustomTimeout(
		ctx,

		func() (any, *http.Response, error) {
			return r.client.CredentialIssuersApi.ReadCredentialIssuerProfile(ctx, data.EnvironmentId.ValueString()).Execute()

		},
		"ReadCredentialIssuerProfile",
		framework.CustomErrorResourceNotFoundWarning,
		credentialIssuerRetryConditions,
		&response,
		time.Duration(timeoutValue)*time.Minute, // 5 mins
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

func (r *CredentialIssuerProfileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state CredentialIssuerProfileResourceModel

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
	CredentialIssuerProfile, d := plan.expand()
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *credentials.CredentialIssuerProfile
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			return r.client.CredentialIssuersApi.UpdateCredentialIssuerProfile(ctx, plan.EnvironmentId.ValueString()).CredentialIssuerProfile(*CredentialIssuerProfile).Execute()
		},
		"UpdateCredentialIssuerProfile",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
		&response,
	)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update the state to save
	state = plan

	// Save updated data into Terraform state
	resp.Diagnostics.Append(state.toState(response)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *CredentialIssuerProfileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Deletion of a credential issuer profile is not allowed, and there is not an associated API.
}

func (r *CredentialIssuerProfileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	splitLength := 2
	attributes := strings.SplitN(req.ID, "/", splitLength)

	if len(attributes) != splitLength {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("invalid id (\"%s\") specified, should be in format \"environment_id/credential_issuer_profile_id/\"", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("environment_id"), attributes[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), attributes[1])...)
}

func (p *CredentialIssuerProfileResourceModel) expand() (*credentials.CredentialIssuerProfile, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := credentials.NewCredentialIssuerProfile(p.Name.ValueString())

	applicationInstanceId := credentials.NewCredentialIssuerProfileApplicationInstance()
	applicationInstanceId.SetId(p.ApplicationInstanceId.ValueString())

	data.SetApplicationInstance(*applicationInstanceId)

	if !p.CreatedAt.IsNull() && !p.CreatedAt.IsUnknown() {
		createdAt, err := time.Parse(time.RFC3339, p.CreatedAt.ValueString())
		if err != nil {
			diags.AddWarning(
				"Unexpected Value",
				fmt.Sprintf("Unexpected createdAt value: %s.  Please report this to the provider maintainers.", err.Error()),
			)
		}
		data.SetCreatedAt(createdAt)
	}

	if !p.UpdatedAt.IsNull() && !p.UpdatedAt.IsUnknown() {
		updatedAt, err := time.Parse(time.RFC3339, p.UpdatedAt.ValueString())
		if err != nil {
			diags.AddWarning(
				"Unexpected Value",
				fmt.Sprintf("Unexpected updatedAt value: %s.  Please report this to the provider maintainers.", err.Error()),
			)
		}
		data.SetUpdatedAt(updatedAt)

		if data == nil {
			diags.AddWarning(
				"Unexpected Value",
				"Credential Issuer Profile object was unexpectedly null on expansion.  Please report this to the provider maintainers.",
			)
		}
	}

	return data, diags
}

func (p *CredentialIssuerProfileResourceModel) toState(apiObject *credentials.CredentialIssuerProfile) diag.Diagnostics {
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
	p.ApplicationInstanceId = framework.StringToTF(*apiObject.GetApplicationInstance().Id)
	p.CreatedAt = framework.TimeOkToTF(apiObject.GetUpdatedAtOk())
	p.UpdatedAt = framework.TimeOkToTF(apiObject.GetUpdatedAtOk())
	p.Name = framework.StringOkToTF(apiObject.GetNameOk())

	return diags
}

func credentialIssuerRetryConditions(ctx context.Context, r *http.Response, p1error *model.P1Error) bool {

	var err error

	if p1error != nil {

		// Credential Issuer Profile's keys may not have propagated after initial environment setup.
		// A read operation is performed to determine if a create or update is necessary when the profile is added to TF state.
		// The read may return nil if the profile has not fully propagated, which woiuld trigger a create instead of an update.
		// Rare, but possible.
		if details, ok := p1error.GetDetailsOk(); ok && details != nil && len(details) > 0 {

			if m, err := regexp.MatchString("^A resource with the specified name already exists", details[0].GetMessage()); err == nil && m {
				tflog.Warn(ctx, fmt.Sprintf("IssuerProfile (prerequisite) has not finished provisioning - %s.  Retrying...", details[0].GetMessage()))
				return true
			}
			if err != nil {
				tflog.Warn(ctx, "Cannot match error string for retry")
				return false
			}

		}

		// detected credentials service not fully deployed yet
		if m, _ := regexp.MatchString("^The actor attempting to perform the request is not authorized.", p1error.GetMessage()); err == nil && m {
			tflog.Warn(ctx, "Insufficient PingOne privileges detected. Retrying...")
			return true
		}
		if err != nil {
			tflog.Warn(ctx, "Cannot match error string for retry")
			return false
		}

		// issuer not found could be the caused by delayed credential issuer
		if m, _ := regexp.MatchString("^The requested resource object cannot be found.", p1error.GetMessage()); err == nil && m {
			tflog.Warn(ctx, "Credential Issuer not found. Retrying...")
			return true
		}
		if err != nil {
			tflog.Warn(ctx, "Cannot match error string for retry")
			return false
		}

	}

	return false
}
