// Copyright © 2025 Ping Identity Corporation

package credentials

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/patrickcping/pingone-go-sdk-v2/credentials"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/legacysdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type CredentialIssuerProfileResource serviceClientType

type CredentialIssuerProfileResourceModel struct {
	Id                    pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId         pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	ApplicationInstanceId pingonetypes.ResourceIDValue `tfsdk:"application_instance_id"`
	CreatedAt             timetypes.RFC3339            `tfsdk:"created_at"`
	UpdatedAt             timetypes.RFC3339            `tfsdk:"updated_at"`
	Name                  types.String                 `tfsdk:"name"`
	Timeouts              timeouts.Value               `tfsdk:"timeouts"`
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

				CustomType: pingonetypes.ResourceIDType{},
			},

			"created_at": schema.StringAttribute{
				Description: "Date and time the issuer profile was created.",
				Computed:    true,

				CustomType: timetypes.RFC3339Type{},

				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"updated_at": schema.StringAttribute{
				Description: "Date and time the issuer profile was last updated.",
				Computed:    true,

				CustomType: timetypes.RFC3339Type{},
			},

			"name": schema.StringAttribute{
				Description: "The name of the credential issuer. The name is included in the metadata of an issued verifiable credential.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(attrMinLength, attrMaxLength),
				},
			},

			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
			}),
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

	r.Client = resourceConfig.Client.API
	if r.Client == nil {
		resp.Diagnostics.AddError(
			"Client not initialised",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.",
		)
		return
	}
}

func (r *CredentialIssuerProfileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state CredentialIssuerProfileResourceModel

	if r.Client == nil || r.Client.CredentialsAPIClient == nil {
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

	defaultTimeout := 10

	timeout, d := plan.Timeouts.Create(ctx, time.Duration(defaultTimeout)*time.Minute)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	readStateConf := &retry.StateChangeConf{
		Pending: []string{
			"404",
			"403",
		},
		Target: []string{
			"200",
		},
		Refresh: func() (interface{}, string, error) {
			base := 10

			fO, fR, fErr := r.Client.CredentialsAPIClient.CredentialIssuersApi.ReadCredentialIssuerProfile(ctx, plan.EnvironmentId.ValueString()).Execute()
			resp, r, err := legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)

			if err != nil {
				return nil, strconv.FormatInt(int64(r.StatusCode), base), err
			}

			return resp, strconv.FormatInt(int64(r.StatusCode), base), nil
		},
		Timeout:                   timeout,
		Delay:                     5 * time.Second,
		MinTimeout:                10 * time.Second,
		ContinuousTargetOccurence: 3,
	}
	readIssuerProfileResponse, err := readStateConf.WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Credential Issuer Profile Create Error",
			fmt.Sprintf("Error waiting for credential issuer profile (environment: %s) to be created: %s", plan.EnvironmentId.ValueString(), err),
		)

		return
	}

	// Build the model for the Create API call
	credentialIssuerProfile := plan.expand()

	// Execute a Create or Update depending on existence of credential issuer profile
	var response *credentials.CredentialIssuerProfile
	if readIssuerProfileResponse == nil {
		// create the issuer profile
		resp.Diagnostics.Append(legacysdk.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				fO, fR, fErr := r.Client.CredentialsAPIClient.CredentialIssuersApi.CreateCredentialIssuerProfile(ctx, plan.EnvironmentId.ValueString()).CredentialIssuerProfile(*credentialIssuerProfile).Execute()
				return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
			},
			"CreateCredentialIssuerProfile",
			legacysdk.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
			&response,
		)...)
		if resp.Diagnostics.HasError() {
			return
		}
	} else {
		// update existing issuer profile
		resp.Diagnostics.Append(legacysdk.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				fO, fR, fErr := r.Client.CredentialsAPIClient.CredentialIssuersApi.UpdateCredentialIssuerProfile(ctx, plan.EnvironmentId.ValueString()).CredentialIssuerProfile(*credentialIssuerProfile).Execute()
				return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
			},
			"UpdateCredentialIssuerProfile",
			legacysdk.DefaultCustomError,
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

	if r.Client == nil || r.Client.CredentialsAPIClient == nil {
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
	resp.Diagnostics.Append(legacysdk.ParseResponseWithCustomTimeout(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.CredentialsAPIClient.CredentialIssuersApi.ReadCredentialIssuerProfile(ctx, data.EnvironmentId.ValueString()).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)

		},
		"ReadCredentialIssuerProfile",
		legacysdk.CustomErrorResourceNotFoundWarning,
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

	if r.Client == nil || r.Client.CredentialsAPIClient == nil {
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
	CredentialIssuerProfile := plan.expand()

	// Run the API call
	var response *credentials.CredentialIssuerProfile
	resp.Diagnostics.Append(legacysdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.CredentialsAPIClient.CredentialIssuersApi.UpdateCredentialIssuerProfile(ctx, plan.EnvironmentId.ValueString()).CredentialIssuerProfile(*CredentialIssuerProfile).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdateCredentialIssuerProfile",
		legacysdk.DefaultCustomError,
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

	idComponents := []framework.ImportComponent{
		{
			Label:  "environment_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
		{
			Label:     "credential_issuer_profile_id",
			Regexp:    verify.P1ResourceIDRegexp,
			PrimaryID: true,
		},
	}

	attributes, err := framework.ParseImportID(req.ID, idComponents...)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			err.Error(),
		)
		return
	}

	for _, idComponent := range idComponents {
		pathKey := idComponent.Label

		if idComponent.PrimaryID {
			pathKey = "id"
		}

		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root(pathKey), attributes[idComponent.Label])...)
	}
}

func (p *CredentialIssuerProfileResourceModel) expand() *credentials.CredentialIssuerProfile {

	data := credentials.NewCredentialIssuerProfile(p.Name.ValueString())

	applicationInstanceId := credentials.NewCredentialIssuerProfileApplicationInstance()
	applicationInstanceId.SetId(p.ApplicationInstanceId.ValueString())

	data.SetApplicationInstance(*applicationInstanceId)

	return data
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

	p.Id = framework.PingOneResourceIDOkToTF(apiObject.GetIdOk())
	p.EnvironmentId = framework.PingOneResourceIDToTF(*apiObject.GetEnvironment().Id)
	p.ApplicationInstanceId = framework.PingOneResourceIDToTF(*apiObject.GetApplicationInstance().Id)
	p.CreatedAt = framework.TimeOkToTF(apiObject.GetCreatedAtOk())
	p.UpdatedAt = framework.TimeOkToTF(apiObject.GetUpdatedAtOk())
	p.Name = framework.StringOkToTF(apiObject.GetNameOk())

	return diags
}

func credentialIssuerRetryConditions(ctx context.Context, r *http.Response, p1error *model.P1Error) bool {

	if p1error != nil {

		// Credential Issuer Profile's keys may not have propagated after initial environment setup.
		// A read operation is performed to determine if a create or update is necessary when the profile is added to TF state.
		// The read may return nil if the profile has not fully propagated, which woiuld trigger a create instead of an update.
		// Rare, but possible.
		if details, ok := p1error.GetDetailsOk(); ok && details != nil && len(details) > 0 {

			m, err := regexp.MatchString("^A resource with the specified name already exists", details[0].GetMessage())
			if err == nil && m {
				tflog.Warn(ctx, fmt.Sprintf("IssuerProfile (prerequisite) has not finished provisioning - %s.  Retrying...", details[0].GetMessage()))
				return true
			}
			if err != nil {
				tflog.Warn(ctx, "Cannot match error string for retry")
				return false
			}

		}

		// detected credentials service not fully deployed yet
		m, err := regexp.MatchString("^The actor attempting to perform the request is not authorized.", p1error.GetMessage())
		if err == nil && m {
			tflog.Warn(ctx, "Insufficient PingOne privileges detected. Retrying...")
			return true
		}
		if err != nil {
			tflog.Warn(ctx, "Cannot match error string for retry")
			return false
		}

		// issuer not found could be the caused by delayed credential issuer
		m, err = regexp.MatchString("^The requested resource object cannot be found.", p1error.GetMessage())
		if err == nil && m {
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
