// Copyright © 2026 Ping Identity Corporation

package base

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/pingidentity/pingone-go-sdk-legacy/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/legacysdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/utils"
)

// Types
type KeyRotationPolicyDataSource serviceClientType

type KeyRotationPolicyDataSourceModel struct {
	Id                  pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId       pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	KeyRotationPolicyId pingonetypes.ResourceIDValue `tfsdk:"key_rotation_policy_id"`
	Name                types.String                 `tfsdk:"name"`
	Algorithm           types.String                 `tfsdk:"algorithm"`
	CurrentKeyId        pingonetypes.ResourceIDValue `tfsdk:"current_key_id"`
	SubjectDn           types.String                 `tfsdk:"subject_dn"`
	KeyLength           types.Int32                  `tfsdk:"key_length"`
	NextKeyId           pingonetypes.ResourceIDValue `tfsdk:"next_key_id"`
	RotatedAt           timetypes.RFC3339            `tfsdk:"rotated_at"`
	RotationPeriod      types.Int32                  `tfsdk:"rotation_period"`
	SignatureAlgorithm  types.String                 `tfsdk:"signature_algorithm"`
	UsageType           types.String                 `tfsdk:"usage_type"`
	ValidityPeriod      types.Int32                  `tfsdk:"validity_period"`
}

// Framework interfaces
var (
	_ datasource.DataSource = &KeyRotationPolicyDataSource{}
)

// New Object
func NewKeyRotationPolicyDataSource() datasource.DataSource {
	return &KeyRotationPolicyDataSource{}
}

// Metadata
func (r *KeyRotationPolicyDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_key_rotation_policy"
}

// Schema
func (r *KeyRotationPolicyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	const attrMinLength = 1

	algorithmDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The algorithm this key rotation policy applies to generated key rotation policy keys.",
	).AllowedValuesEnum(management.AllowedEnumKeyRotationPolicyAlgorithmEnumValues)

	subjectDnDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The DN this key rotation policy will apply to generated key rotation policy keys. The value will be applied as both issuerDN and subjectDN because generated keys are self-signed.",
	)

	keyLengthDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The number of bytes of a cryptographic key this key rotation policy will apply to generated key rotation policy keys.",
	).AllowedValues(utils.Int32SliceToAnySlice([]int32{2048, 3072, 4096, 7680})...)

	rotationPeriodDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The number of days between key rotations.  The minimum value allowed is `30` days, while the maximum value allowed is 1 day less than the value set in the `validity_period` parameter.",
	)

	signatureAlgorithmDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The signature algorithm this key rotation policy will apply to generated key rotation policy keys.",
	).AllowedValuesEnum(management.AllowedEnumKeyRotationPolicySigAlgorithmEnumValues)

	usageTypeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"How the key rotation policy will be used, pertaining to what operations the key rotation policy supports.",
	).AllowedValuesEnum(management.AllowedEnumKeyRotationPolicyUsageTypeEnumValues)

	validityPeriodDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Controls the \"Starts At\" and \"Expires At\" fields this key rotation policy will apply to generated key rotation policy keys.",
	)

	currentKeyIdDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The `kid` (key identifier) of the key rotation policy key designated as `CURRENT`.",
	)

	nextKeyIdDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The `kid` (key identifier) of the key rotation policy key designated as `NEXT`.",
	)

	rotatedAtDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The last time the key rotation policy was rotated.",
	)

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Datasource to retrieve details of a key rotation policy configured in a PingOne environment.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment that is configured with the key rotation policy."),
			),

			"key_rotation_policy_id": schema.StringAttribute{
				Description: "The ID of the key rotation policy to retrieve. Either `key_rotation_policy_id`, or `name` can be used to retrieve the key rotation policy, but cannot be set together.",
				Optional:    true,

				CustomType: pingonetypes.ResourceIDType{},

				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("name")),
				},
			},

			"name": schema.StringAttribute{
				Description: "A string that specifies the name of the key rotation policy to retrieve. Either `key_rotation_policy_id`, or `name` can be used to retrieve the key rotation policy, but cannot be set together.",
				Optional:    true,

				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("key_rotation_policy_id")),
					stringvalidator.LengthAtLeast(attrMinLength),
				},
			},

			"algorithm": schema.StringAttribute{
				Description:         algorithmDescription.Description,
				MarkdownDescription: algorithmDescription.MarkdownDescription,
				Computed:            true,
			},

			"subject_dn": schema.StringAttribute{
				Description:         subjectDnDescription.Description,
				MarkdownDescription: subjectDnDescription.MarkdownDescription,
				Computed:            true,
			},

			"key_length": schema.Int32Attribute{
				Description:         keyLengthDescription.Description,
				MarkdownDescription: keyLengthDescription.MarkdownDescription,
				Computed:            true,
			},

			"rotated_at": schema.StringAttribute{
				Description:         rotatedAtDescription.Description,
				MarkdownDescription: rotatedAtDescription.MarkdownDescription,
				Computed:            true,

				CustomType: timetypes.RFC3339Type{},
			},

			"rotation_period": schema.Int32Attribute{
				Description:         rotationPeriodDescription.Description,
				MarkdownDescription: rotationPeriodDescription.MarkdownDescription,
				Computed:            true,
			},

			"signature_algorithm": schema.StringAttribute{
				Description:         signatureAlgorithmDescription.Description,
				MarkdownDescription: signatureAlgorithmDescription.MarkdownDescription,
				Computed:            true,
			},

			"usage_type": schema.StringAttribute{
				Description:         usageTypeDescription.Description,
				MarkdownDescription: usageTypeDescription.MarkdownDescription,
				Computed:            true,
			},

			"validity_period": schema.Int32Attribute{
				Description:         validityPeriodDescription.Description,
				MarkdownDescription: validityPeriodDescription.MarkdownDescription,
				Computed:            true,
			},

			"current_key_id": schema.StringAttribute{
				Description:         currentKeyIdDescription.Description,
				MarkdownDescription: currentKeyIdDescription.MarkdownDescription,
				Computed:            true,

				CustomType: pingonetypes.ResourceIDType{},
			},

			"next_key_id": schema.StringAttribute{
				Description:         nextKeyIdDescription.Description,
				MarkdownDescription: nextKeyIdDescription.MarkdownDescription,
				Computed:            true,

				CustomType: pingonetypes.ResourceIDType{},
			},
		},
	}
}

func (r *KeyRotationPolicyDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *KeyRotationPolicyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *KeyRotationPolicyDataSourceModel

	if r.Client == nil || r.Client.ManagementAPIClient == nil {
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

	var keyRotationPolicy *management.KeyRotationPolicy

	if !data.Name.IsNull() {

		// Run the API call
		resp.Diagnostics.Append(legacysdk.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				pagedIterator := r.Client.ManagementAPIClient.KeyRotationPoliciesApi.GetKeyRotationPolicies(ctx, data.EnvironmentId.ValueString()).Execute()
				var initialHttpResponse *http.Response

				for pageCursor, err := range pagedIterator {
					if err != nil {
						return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), nil, pageCursor.HTTPResponse, err)
					}

					if initialHttpResponse == nil {
						initialHttpResponse = pageCursor.HTTPResponse
					}

					if keyRotationPolicies, ok := pageCursor.EntityArray.Embedded.GetKeyRotationPoliciesOk(); ok {

						for _, keyRotationPolicyItem := range keyRotationPolicies {
							if strings.EqualFold(keyRotationPolicyItem.GetName(), data.Name.ValueString()) {
								return &keyRotationPolicyItem, pageCursor.HTTPResponse, nil
							}
						}
					}
				}

				return nil, initialHttpResponse, nil
			},
			"GetKeyRotationPolicies",
			legacysdk.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
			&keyRotationPolicy,
		)...)
		if resp.Diagnostics.HasError() {
			return
		}

		if keyRotationPolicy == nil {
			resp.Diagnostics.AddError(
				"Cannot find key rotation policy from name",
				fmt.Sprintf("The key rotation policy %s for environment %s cannot be found", data.Name.String(), data.EnvironmentId.String()),
			)
			return
		}

	} else if !data.KeyRotationPolicyId.IsNull() {

		// Run the API call
		resp.Diagnostics.Append(legacysdk.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				fO, fR, fErr := r.Client.ManagementAPIClient.KeyRotationPoliciesApi.GetKeyRotationPolicy(ctx, data.EnvironmentId.ValueString(), data.KeyRotationPolicyId.ValueString()).Execute()
				return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
			},
			"GetKeyRotationPolicy",
			legacysdk.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
			&keyRotationPolicy,
		)...)
		if resp.Diagnostics.HasError() {
			return
		}

	} else {
		resp.Diagnostics.AddError(
			"Missing parameter",
			"Cannot find the requested key rotation policy. key_rotation_policy_id or name must be set.",
		)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(keyRotationPolicy)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (p *KeyRotationPolicyDataSourceModel) toState(apiObject *management.KeyRotationPolicy) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	environment := apiObject.GetEnvironment()
	p.Id = framework.PingOneResourceIDOkToTF(apiObject.GetIdOk())
	p.EnvironmentId = framework.PingOneResourceIDOkToTF(environment.GetIdOk())
	p.KeyRotationPolicyId = framework.PingOneResourceIDOkToTF(apiObject.GetIdOk())
	p.Name = framework.StringOkToTF(apiObject.GetNameOk())
	p.Algorithm = framework.EnumOkToTF(apiObject.GetAlgorithmOk())
	p.CurrentKeyId = framework.PingOneResourceIDOkToTF(apiObject.GetCurrentKeyIdOk())
	p.SubjectDn = framework.StringOkToTF(apiObject.GetDnOk())
	p.KeyLength = framework.Int32OkToTF(apiObject.GetKeyLengthOk())
	p.NextKeyId = framework.PingOneResourceIDOkToTF(apiObject.GetNextKeyIdOk())
	p.RotatedAt = framework.TimeOkToTF(apiObject.GetRotatedAtOk())
	p.RotationPeriod = framework.Int32OkToTF(apiObject.GetRotationPeriodOk())
	p.SignatureAlgorithm = framework.EnumOkToTF(apiObject.GetSignatureAlgorithmOk())
	p.UsageType = framework.EnumOkToTF(apiObject.GetUsageTypeOk())
	p.ValidityPeriod = framework.Int32OkToTF(apiObject.GetValidityPeriodOk())

	return diags
}
