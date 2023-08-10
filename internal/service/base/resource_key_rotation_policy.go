package base

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/utils"
)

// Types
type KeyRotationPolicyResource serviceClientType

type keyRotationPolicyResourceModel struct {
	Id                 types.String `tfsdk:"id"`
	EnvironmentId      types.String `tfsdk:"environment_id"`
	Name               types.String `tfsdk:"name"`
	Algorithm          types.String `tfsdk:"algorithm"`
	CurrentKeyId       types.String `tfsdk:"current_key_id"`
	SubjectDn          types.String `tfsdk:"subject_dn"`
	KeyLength          types.Int64  `tfsdk:"key_length"`
	NextKeyId          types.String `tfsdk:"next_key_id"`
	RotatedAt          types.String `tfsdk:"rotated_at"`
	RotationPeriod     types.Int64  `tfsdk:"rotation_period"`
	SignatureAlgorithm types.String `tfsdk:"signature_algorithm"`
	UsageType          types.String `tfsdk:"usage_type"`
	ValidityPeriod     types.Int64  `tfsdk:"validity_period"`
}

// Framework interfaces
var (
	_ resource.Resource                = &KeyRotationPolicyResource{}
	_ resource.ResourceWithConfigure   = &KeyRotationPolicyResource{}
	_ resource.ResourceWithImportState = &KeyRotationPolicyResource{}
)

// New Object
func NewKeyRotationPolicyResource() resource.Resource {
	return &KeyRotationPolicyResource{}
}

// Metadata
func (r *KeyRotationPolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_key_rotation_policy"
}

// Schema.
func (r *KeyRotationPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	const attrMinLength = 1
	const rotationPeriodMinimum = 30
	const rotationPeriodDefault = 90
	const validityPeriodDefault = 365
	var allowedKeyLengths = []int64{2048, 3072, 4096, 7680}

	algorithmDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The algorithm this key rotation policy applies to generated key rotation policy keys.",
	).AllowedValuesEnum(management.AllowedEnumKeyRotationPolicyAlgorithmEnumValues)

	subjectDnDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The DN this key rotation policy will apply to generated key rotation policy keys. The value will be applied as both issuerDN and subjectDN because generated keys are self-signed.",
	)

	keyLengthDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The number of bytes of a cryptographic key this key rotation policy will apply to generated key rotation policy keys.",
	).AllowedValues(utils.Int64SliceToAnySlice(allowedKeyLengths)...)

	rotationPeriodDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The number of days between key rotations.  The minimum value allowed is `30` days, while the maximum value allowed is 1 day less than the value set in the `validity_period` parameter.",
	).DefaultValue(rotationPeriodDefault)

	signatureAlgorithmDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The signature algorithm this key rotation policy will apply to generated key rotation policy keys.",
	).AllowedValuesEnum(management.AllowedEnumKeyRotationPolicySigAlgorithmEnumValues)

	usageTypeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"How the key rotation policy will be used, pertaining to what operations the key rotation policy supports.",
	).AllowedValuesEnum(management.AllowedEnumKeyRotationPolicyUsageTypeEnumValues)

	validityPeriodDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Controls the \"Starts At\" and \"Expires At\" fields this key rotation policy will apply to generated key rotation policy keys.",
	).DefaultValue(validityPeriodDefault)

	currentKeyIdDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The `kid` (key identifier) of the key rotation policy key designated as `CURRENT`.",
	)

	nextKeyIdDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The `kid` (key identifier) of the key rotation policy key designated as `NEXT`.",
	)

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage PingOne key rotation policies for an environment.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment to configure a key rotation policy for."),
			),

			"name": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("Human-readable name displayed in the admin console.").Description,
				Required:    true,

				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
				},
			},

			"algorithm": schema.StringAttribute{
				Description:         algorithmDescription.Description,
				MarkdownDescription: algorithmDescription.MarkdownDescription,
				Required:            true,

				Validators: []validator.String{
					stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumKeyRotationPolicyAlgorithmEnumValues)...),
				},
			},

			"subject_dn": schema.StringAttribute{
				Description:         subjectDnDescription.Description,
				MarkdownDescription: subjectDnDescription.MarkdownDescription,
				Required:            true,

				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
				},
			},

			"key_length": schema.Int64Attribute{
				Description:         keyLengthDescription.Description,
				MarkdownDescription: keyLengthDescription.MarkdownDescription,
				Required:            true,

				Validators: []validator.Int64{
					int64validator.OneOf(
						allowedKeyLengths...,
					),
				},
			},

			"rotated_at": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The last time the key rotation policy was rotated.").Description,
				Computed:    true,
			},

			"rotation_period": schema.Int64Attribute{
				Description:         rotationPeriodDescription.Description,
				MarkdownDescription: rotationPeriodDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,

				Default: int64default.StaticInt64(rotationPeriodDefault),

				Validators: []validator.Int64{
					int64validator.AtLeast(rotationPeriodMinimum),
					// todo: The maximum value is 1 day less than the `validityPeriod` value
				},
			},

			"signature_algorithm": schema.StringAttribute{
				Description:         signatureAlgorithmDescription.Description,
				MarkdownDescription: signatureAlgorithmDescription.MarkdownDescription,
				Required:            true,

				Validators: []validator.String{
					stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumKeyRotationPolicySigAlgorithmEnumValues)...),
				},
			},

			"usage_type": schema.StringAttribute{
				Description:         usageTypeDescription.Description,
				MarkdownDescription: usageTypeDescription.MarkdownDescription,
				Required:            true,

				Validators: []validator.String{
					stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumKeyRotationPolicyUsageTypeEnumValues)...),
				},
			},

			"validity_period": schema.Int64Attribute{
				Description:         validityPeriodDescription.Description,
				MarkdownDescription: validityPeriodDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,

				Default: int64default.StaticInt64(validityPeriodDefault),
			},

			"current_key_id": schema.StringAttribute{
				Description:         currentKeyIdDescription.Description,
				MarkdownDescription: currentKeyIdDescription.MarkdownDescription,
				Computed:            true,
			},

			"next_key_id": schema.StringAttribute{
				Description:         nextKeyIdDescription.Description,
				MarkdownDescription: nextKeyIdDescription.MarkdownDescription,
				Computed:            true,
			},
		},
	}
}

func (r *KeyRotationPolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *KeyRotationPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state keyRotationPolicyResourceModel

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
	keyRotationPolicy := plan.expand()

	// Run the API call
	var response *management.KeyRotationPolicy
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			return r.Client.KeyRotationPoliciesApi.CreateKeyRotationPolicy(ctx, plan.EnvironmentId.ValueString()).KeyRotationPolicy(*keyRotationPolicy).Execute()
		},
		"CreateKeyRotationPolicy",
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

func (r *KeyRotationPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *keyRotationPolicyResourceModel

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
	var response *management.KeyRotationPolicy
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			return r.Client.KeyRotationPoliciesApi.GetKeyRotationPolicy(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
		},
		"GetKeyRotationPolicy",
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

func (r *KeyRotationPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state keyRotationPolicyResourceModel

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
	keyRotationPolicy := plan.expand()

	// Run the API call
	var response *management.KeyRotationPolicy
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			return r.Client.KeyRotationPoliciesApi.UpdateKeyRotationPolicy(ctx, plan.EnvironmentId.ValueString(), plan.Id.ValueString()).KeyRotationPolicy(*keyRotationPolicy).Execute()
		},
		"UpdateKeyRotationPolicy",
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

func (r *KeyRotationPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *keyRotationPolicyResourceModel

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
			r, err := r.Client.KeyRotationPoliciesApi.DeleteKeyRotationPolicy(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return nil, r, err
		},
		"DeleteKeyRotationPolicy",
		framework.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
		nil,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *KeyRotationPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	splitLength := 2
	attributes := strings.SplitN(req.ID, "/", splitLength)

	if len(attributes) != splitLength {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("invalid id (\"%s\") specified, should be in format \"environment_id/key_rotation_policy_id\"", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("environment_id"), attributes[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), attributes[1])...)
}

func (p *keyRotationPolicyResourceModel) expand() *management.KeyRotationPolicy {
	data := management.NewKeyRotationPolicy(
		management.EnumKeyRotationPolicyAlgorithm(p.Algorithm.ValueString()),
		p.SubjectDn.ValueString(),
		int32(p.KeyLength.ValueInt64()),
		p.Name.ValueString(),
		management.EnumKeyRotationPolicySigAlgorithm(p.SignatureAlgorithm.ValueString()),
		management.EnumKeyRotationPolicyUsageType(p.UsageType.ValueString()),
	)

	if !p.RotationPeriod.IsNull() && !p.RotationPeriod.IsUnknown() {
		data.SetRotationPeriod(int32(p.RotationPeriod.ValueInt64()))
	}

	if !p.ValidityPeriod.IsNull() && !p.ValidityPeriod.IsUnknown() {
		data.SetValidityPeriod(int32(p.ValidityPeriod.ValueInt64()))
	}

	return data
}

func (p *keyRotationPolicyResourceModel) toState(apiObject *management.KeyRotationPolicy) diag.Diagnostics {
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
	p.Name = framework.StringOkToTF(apiObject.GetNameOk())
	p.Algorithm = framework.EnumOkToTF(apiObject.GetAlgorithmOk())
	p.CurrentKeyId = framework.StringOkToTF(apiObject.GetCurrentKeyIdOk())
	p.SubjectDn = framework.StringOkToTF(apiObject.GetDnOk())
	p.KeyLength = framework.Int32OkToTF(apiObject.GetKeyLengthOk())
	p.NextKeyId = framework.StringOkToTF(apiObject.GetNextKeyIdOk())
	p.RotatedAt = framework.TimeOkToTF(apiObject.GetRotatedAtOk())
	p.RotationPeriod = framework.Int32OkToTF(apiObject.GetRotationPeriodOk())
	p.SignatureAlgorithm = framework.EnumOkToTF(apiObject.GetSignatureAlgorithmOk())
	p.UsageType = framework.EnumOkToTF(apiObject.GetUsageTypeOk())
	p.ValidityPeriod = framework.Int32OkToTF(apiObject.GetValidityPeriodOk())

	return diags
}
