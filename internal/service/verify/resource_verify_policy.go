package verify

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	"github.com/patrickcping/pingone-go-sdk-v2/verify"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/utils"
)

// Types
type VerifyPolicyResource struct {
	client *verify.APIClient
	region model.RegionMapping
}

type verifyPolicyResourceModel struct {
	Id               types.String `tfsdk:"id"`
	EnvironmentId    types.String `tfsdk:"environment_id"`
	Name             types.String `tfsdk:"name"`
	Default          types.Bool   `tfsdk:"default"`
	Description      types.String `tfsdk:"description"`
	GovernmentId     types.Object `tfsdk:"government_id"`
	FacialComparison types.Object `tfsdk:"facial_comparison"`
	Liveness         types.Object `tfsdk:"liveness"`
	Email            types.Object `tfsdk:"email"`
	Phone            types.Object `tfsdk:"phone"`
	Transaction      types.Object `tfsdk:"transaction"`
	CreatedAt        types.String `tfsdk:"created_at"`
	UpdatedAt        types.String `tfsdk:"updated_at"`
}

type governmentIdModel struct {
	Verify types.String `tfsdk:"verify"`
}

type facialComparisonModel struct {
	Verify    types.String `tfsdk:"verify"`
	Threshold types.String `tfsdk:"threshold"`
}

type livenessnModel struct {
	Verify    types.String `tfsdk:"verify"`
	Threshold types.String `tfsdk:"threshold"`
}

type genericTimeoutModel struct {
	Duration types.Int64  `tfsdk:"duration"`
	TimeUnit types.String `tfsdk:"time_unit"`
}

type deviceModel struct {
	CreateMfaDevice types.Bool   `tfsdk:"create_mfa_device"`
	OTP             types.Object `tfsdk:"otp"`
	Verify          types.String `tfsdk:"verify"`
}

type otpConfigurationModel struct {
	Attempts     types.Object `tfsdk:"attempts"`
	Deliveries   types.Object `tfsdk:"deliveries"`
	LifeTime     types.Object `tfsdk:"lifetime"`
	Notification types.Object `tfsdk:"notification"`
}

type otpAttemptsModel struct {
	Count types.Int64 `tfsdk:"count"`
}

type otpDeliveriessModel struct {
	Count    types.Int64  `tfsdk:"count"`
	Cooldown types.Object `tfsdk:"cooldown"`
}

type otpDeliveriessCooldownModel struct {
	Duration types.Int64  `tfsdk:"duration"`
	TimeUnit types.String `tfsdk:"time_unit"`
}

type otpNotificationModel struct {
	TemplateName types.String `tfsdk:"template_name"`
	VariantName  types.String `tfsdk:"variant_name"`
}

type transactionModel struct {
	Timeout            types.Object `tfsdk:"timeout"`
	DataCollection     types.Object `tfsdk:"data_collection"`
	DataCollectionOnly types.Bool   `tfsdk:"data_collection_only"`
}

type transactionDataCollectionModel struct {
	Timeout types.Object `tfsdk:"timeout"`
}

var (
	genericTimeoutServiceTFObjectTypes = map[string]attr.Type{
		"duration":  types.Int64Type,
		"time_unit": types.StringType,
	}

	governmentIdServiceTFObjectTypes = map[string]attr.Type{
		"verify": types.StringType,
	}

	facialComparisonServiceTFObjectTypes = map[string]attr.Type{
		"verify":    types.StringType,
		"threshold": types.StringType,
	}

	livenessServiceTFObjectTypes = map[string]attr.Type{
		"verify":    types.StringType,
		"threshold": types.StringType,
	}

	transactionServiceTFObjectTypes = map[string]attr.Type{
		"timeout":              types.ObjectType{AttrTypes: genericTimeoutServiceTFObjectTypes},
		"data_collection":      types.ObjectType{AttrTypes: dataCollectionServiceTFObjectTypes},
		"data_collection_only": types.BoolType,
	}

	dataCollectionServiceTFObjectTypes = map[string]attr.Type{
		"timeout": types.ObjectType{AttrTypes: genericTimeoutServiceTFObjectTypes},
	}
)

// Framework interfaces
var (
	_ resource.Resource                = &VerifyPolicyResource{}
	_ resource.ResourceWithConfigure   = &VerifyPolicyResource{}
	_ resource.ResourceWithImportState = &VerifyPolicyResource{}
)

// New Object
func NewVerifyPolicyResource() resource.Resource {
	return &VerifyPolicyResource{}
}

// Metadata
func (r *VerifyPolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_verify_policy"
}

func (r *VerifyPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	// schema descriptions and validation settings
	const attrMinLength = 1
	const attrMaxLength = 1024

	nameDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Name of the verification policy displayed in PingOne Admin UI.",
	)

	descriptionDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Description of the verification policy displayed in PingOne Admin UI, 1-1024 characters.",
	)

	defaultDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Required as `true` to set the verify policy as the default policy for the environment; otherwise optional and defaults to `false`.",
	)

	verifyOptionPhraseFmt := "`REQUIRED`, `OPTIONAL`, or `DISABLED`."
	governmentIdVerifyDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("Controls if Government ID verification is %s", verifyOptionPhraseFmt),
	)

	facialComparisonVerifyDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("Controls if facial comparison verification is %s", verifyOptionPhraseFmt),
	)

	facialComparisonThresholdDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Threshold for successful facial comparison; can be `LOW`, `MEDIUM`, or `HIGH` (for which PingOne Verify uses industry and vendor recommended definitions).",
	)

	livenessVerifyDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("Controls if liveness check is %s", verifyOptionPhraseFmt),
	)

	livenessThresholdDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Threshold for successful liveness comparison; can be `LOW`, `MEDIUM`, or `HIGH` (for which PingOne Verify uses industry and vendor recommended definitions).",
	)

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to configure the requirements to verify a user, including the parameters for verification, such as the number of one-time password (OTP) attempts and OTP expiration.\n\n" +
			"A verify policy defines which of the following five checks are performed for a verification transaction and configures the parameters of each check. The checks can be either required or optional. " +
			"If a type is optional, then the transaction can be processed with or without the documents for that type. If the documents are provided for that type and the optional type verification fails, it will not cause the entire transaction to fail.\n\n" +
			"Verify policies can perform any of five checks:\n" +
			"* Government identity document - Validate a government-issued identity document, which includes a photograph." +
			"* Facial comparison - Compare a mobile phone self-image to a reference photograph, such as on a government ID or previously verified photograph." +
			"* Liveness - Inspect a mobile phone self-image for evidence that the subject is alive and not a representation, such as a photograph or mask." +
			"* Email - Receive a one-time password (OTP) on an email address and return the OTP to the service." +
			"* Phone - Receive a one-time password (OTP) on a mobile phone and return the OTP to the service.\n",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("PingOne environment identifier (UUID) in which the verify policy exists."),
			),

			"name": schema.StringAttribute{
				Description:         nameDescription.Description,
				MarkdownDescription: nameDescription.MarkdownDescription,
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
				},
			},

			"default": schema.BoolAttribute{
				Description:         defaultDescription.Description,
				MarkdownDescription: defaultDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,
			},

			"description": schema.StringAttribute{
				Description:         descriptionDescription.Description,
				MarkdownDescription: descriptionDescription.MarkdownDescription,
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(attrMinLength, attrMaxLength),
				},
			},

			"government_id": schema.SingleNestedAttribute{
				Optional: true,
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"verify": schema.StringAttribute{
						Description:         governmentIdVerifyDescription.Description,
						MarkdownDescription: governmentIdVerifyDescription.MarkdownDescription,
						Optional:            true,
						Computed:            true,
						Validators: []validator.String{
							stringvalidator.OneOf(utils.EnumSliceToStringSlice(verify.AllowedEnumVerifyEnumValues)...),
						},
					},
				},
			},

			"facial_comparison": schema.SingleNestedAttribute{
				Optional: true,
				Computed: true,

				Attributes: map[string]schema.Attribute{
					"verify": schema.StringAttribute{
						Description:         facialComparisonVerifyDescription.Description,
						MarkdownDescription: facialComparisonVerifyDescription.MarkdownDescription,
						Optional:            true,
						Computed:            true,
						Validators: []validator.String{
							stringvalidator.OneOf(utils.EnumSliceToStringSlice(verify.AllowedEnumVerifyEnumValues)...),
						},
					},
					"threshold": schema.StringAttribute{
						Description:         facialComparisonThresholdDescription.Description,
						MarkdownDescription: facialComparisonThresholdDescription.MarkdownDescription,
						Optional:            true,
						Computed:            true,
						Validators: []validator.String{
							stringvalidator.OneOf(utils.EnumSliceToStringSlice(verify.AllowedEnumThresholdEnumValues)...),
						},
					},
				},
			},

			"liveness": schema.SingleNestedAttribute{
				Optional: true,
				Computed: true,

				Attributes: map[string]schema.Attribute{
					"verify": schema.StringAttribute{
						Description:         livenessVerifyDescription.Description,
						MarkdownDescription: livenessVerifyDescription.MarkdownDescription,
						Optional:            true,
						Computed:            true,
						Validators: []validator.String{
							stringvalidator.OneOf(utils.EnumSliceToStringSlice(verify.AllowedEnumVerifyEnumValues)...),
						},
					},
					"threshold": schema.StringAttribute{
						Description:         livenessThresholdDescription.Description,
						MarkdownDescription: livenessThresholdDescription.MarkdownDescription,
						Optional:            true,
						Computed:            true,
						Validators: []validator.String{
							stringvalidator.OneOf(utils.EnumSliceToStringSlice(verify.AllowedEnumThresholdEnumValues)...),
						},
					},
				},
			},

			"email": schema.SingleNestedAttribute{
				Optional: true,

				Attributes: map[string]schema.Attribute{
					"create_mfa_device": schema.BoolAttribute{
						Description:         livenessVerifyDescription.Description,
						MarkdownDescription: livenessVerifyDescription.MarkdownDescription,
						Optional:            true,
					},
					"otp": schema.SingleNestedAttribute{
						Description: "Contains template parameters.",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"attempts": schema.SingleNestedAttribute{
								Description: "Contains template parameters.",
								Optional:    true,
								Attributes: map[string]schema.Attribute{
									"count": schema.Int64Attribute{
										Description:         livenessVerifyDescription.Description,
										MarkdownDescription: livenessVerifyDescription.MarkdownDescription,
										Optional:            true,
									},
								},
							},
							"deliveries": schema.SingleNestedAttribute{
								Description: "Contains template parameters.",
								Optional:    true,
								Attributes: map[string]schema.Attribute{
									"count": schema.Int64Attribute{
										Description:         livenessVerifyDescription.Description,
										MarkdownDescription: livenessVerifyDescription.MarkdownDescription,
										Optional:            true,
									},
									"cooldown": schema.SingleNestedAttribute{
										Description: "Contains template parameters.",
										Optional:    true,
										Attributes: map[string]schema.Attribute{
											"duration": schema.Int64Attribute{
												Description:         livenessVerifyDescription.Description,
												MarkdownDescription: livenessVerifyDescription.MarkdownDescription,
												Required:            true,
											},
											"time_unit": schema.StringAttribute{
												Description:         livenessVerifyDescription.Description,
												MarkdownDescription: livenessVerifyDescription.MarkdownDescription,
												Required:            true,
												Validators: []validator.String{
													stringvalidator.OneOf(
														string(verify.ENUMLONGTIMEUNIT_SECONDS),
														string(verify.ENUMLONGTIMEUNIT_MINUTES),
														string(verify.ENUMLONGTIMEUNIT_HOURS),
													),
												},
											},
										},
									},
								},
							},
							"lifetime": schema.SingleNestedAttribute{
								Description: "Contains template parameters.",
								Optional:    true,
								Attributes: map[string]schema.Attribute{
									"duration": schema.Int64Attribute{
										Description:         livenessVerifyDescription.Description,
										MarkdownDescription: livenessVerifyDescription.MarkdownDescription,
										Required:            true,
									},
									"time_unit": schema.StringAttribute{
										Description:         livenessVerifyDescription.Description,
										MarkdownDescription: livenessVerifyDescription.MarkdownDescription,
										Required:            true,
										Validators: []validator.String{
											stringvalidator.OneOf(
												string(verify.ENUMLONGTIMEUNIT_SECONDS),
												string(verify.ENUMLONGTIMEUNIT_MINUTES),
												string(verify.ENUMLONGTIMEUNIT_HOURS),
											),
										},
									},
								},
							},
							"notification": schema.SingleNestedAttribute{
								Description: "Contains template parameters.",
								Optional:    true,
								Attributes: map[string]schema.Attribute{
									"template_name": schema.StringAttribute{
										Description:         livenessVerifyDescription.Description,
										MarkdownDescription: livenessVerifyDescription.MarkdownDescription,
										Required:            true,
									},
									"variant_name": schema.StringAttribute{
										Description:         livenessVerifyDescription.Description,
										MarkdownDescription: livenessVerifyDescription.MarkdownDescription,
										Required:            true,
									},
								},
							},
						},
					},
					"verify": schema.StringAttribute{
						Description:         governmentIdVerifyDescription.Description,
						MarkdownDescription: governmentIdVerifyDescription.MarkdownDescription,
						Required:            true,
						Validators: []validator.String{
							stringvalidator.OneOf(utils.EnumSliceToStringSlice(verify.AllowedEnumVerifyEnumValues)...),
						},
					},
				},
			},

			"phone": schema.SingleNestedAttribute{
				Optional: true,

				Attributes: map[string]schema.Attribute{
					"create_mfa_device": schema.BoolAttribute{
						Description:         livenessVerifyDescription.Description,
						MarkdownDescription: livenessVerifyDescription.MarkdownDescription,
						Optional:            true,
					},
					"otp": schema.SingleNestedAttribute{
						Description: "Contains template parameters.",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"attempts": schema.SingleNestedAttribute{
								Description: "Contains template parameters.",
								Optional:    true,
								Attributes: map[string]schema.Attribute{
									"count": schema.Int64Attribute{
										Description:         livenessVerifyDescription.Description,
										MarkdownDescription: livenessVerifyDescription.MarkdownDescription,
										Optional:            true,
									},
								},
							},
							"deliveries": schema.SingleNestedAttribute{
								Description: "Contains template parameters.",
								Optional:    true,
								Attributes: map[string]schema.Attribute{
									"count": schema.Int64Attribute{
										Description:         livenessVerifyDescription.Description,
										MarkdownDescription: livenessVerifyDescription.MarkdownDescription,
										Optional:            true,
									},
									"cooldown": schema.SingleNestedAttribute{
										Description: "Contains template parameters.",
										Optional:    true,
										Attributes: map[string]schema.Attribute{
											"duration": schema.Int64Attribute{
												Description:         livenessVerifyDescription.Description,
												MarkdownDescription: livenessVerifyDescription.MarkdownDescription,
												Required:            true,
											},
											"time_unit": schema.StringAttribute{
												Description:         livenessVerifyDescription.Description,
												MarkdownDescription: livenessVerifyDescription.MarkdownDescription,
												Required:            true,
												Validators: []validator.String{
													stringvalidator.OneOf(
														string(verify.ENUMLONGTIMEUNIT_SECONDS),
														string(verify.ENUMLONGTIMEUNIT_MINUTES),
														string(verify.ENUMLONGTIMEUNIT_HOURS),
													),
												},
											},
										},
									},
								},
							},
							"lifetime": schema.SingleNestedAttribute{
								Description: "Contains template parameters.",
								Optional:    true,
								Attributes: map[string]schema.Attribute{
									"duration": schema.Int64Attribute{
										Description:         livenessVerifyDescription.Description,
										MarkdownDescription: livenessVerifyDescription.MarkdownDescription,
										Required:            true,
									},
									"time_unit": schema.StringAttribute{
										Description:         livenessVerifyDescription.Description,
										MarkdownDescription: livenessVerifyDescription.MarkdownDescription,
										Required:            true,
										Validators: []validator.String{
											stringvalidator.OneOf(
												string(verify.ENUMLONGTIMEUNIT_SECONDS),
												string(verify.ENUMLONGTIMEUNIT_MINUTES),
												string(verify.ENUMLONGTIMEUNIT_HOURS),
											),
										},
									},
								},
							},
							"notification": schema.SingleNestedAttribute{
								Description: "Contains template parameters.",
								Optional:    true,
								Attributes: map[string]schema.Attribute{
									"template_name": schema.StringAttribute{
										Description:         livenessVerifyDescription.Description,
										MarkdownDescription: livenessVerifyDescription.MarkdownDescription,
										Required:            true,
									},
									"variant_name": schema.StringAttribute{
										Description:         livenessVerifyDescription.Description,
										MarkdownDescription: livenessVerifyDescription.MarkdownDescription,
										Required:            true,
									},
								},
							},
						},
					},
					"verify": schema.StringAttribute{
						Description:         governmentIdVerifyDescription.Description,
						MarkdownDescription: governmentIdVerifyDescription.MarkdownDescription,
						Required:            true,
						Validators: []validator.String{
							stringvalidator.OneOf(utils.EnumSliceToStringSlice(verify.AllowedEnumVerifyEnumValues)...),
						},
					},
				},
			},

			"transaction": schema.SingleNestedAttribute{
				Optional: true,
				Computed: true,

				Attributes: map[string]schema.Attribute{
					"timeout": schema.SingleNestedAttribute{
						Description: "Contains template parameters.",
						Optional:    true,
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"duration": schema.Int64Attribute{
								Description:         livenessVerifyDescription.Description,
								MarkdownDescription: livenessVerifyDescription.MarkdownDescription,
								Optional:            true,
								Computed:            true,
							},
							"time_unit": schema.StringAttribute{
								Description:         livenessVerifyDescription.Description,
								MarkdownDescription: livenessVerifyDescription.MarkdownDescription,
								Optional:            true,
								Computed:            true,
								Validators: []validator.String{
									stringvalidator.OneOf(
										string(verify.ENUMLONGTIMEUNIT_SECONDS),
										string(verify.ENUMLONGTIMEUNIT_MINUTES),
									),
								},
							},
						},
					},
					"data_collection": schema.SingleNestedAttribute{
						Description: "Contains template parameters.",
						Optional:    true,
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"timeout": schema.SingleNestedAttribute{
								Description: "Contains template parameters.",
								Optional:    true,
								Computed:    true,
								Attributes: map[string]schema.Attribute{
									"duration": schema.Int64Attribute{
										Description:         livenessVerifyDescription.Description,
										MarkdownDescription: livenessVerifyDescription.MarkdownDescription,
										Optional:            true,
										Computed:            true,
									},
									"time_unit": schema.StringAttribute{
										Description:         livenessVerifyDescription.Description,
										MarkdownDescription: livenessVerifyDescription.MarkdownDescription,
										Optional:            true,
										Computed:            true,
										Validators: []validator.String{
											stringvalidator.OneOf(
												string(verify.ENUMLONGTIMEUNIT_SECONDS),
												string(verify.ENUMLONGTIMEUNIT_MINUTES),
											),
										},
									},
								},
							},
						},
					},
					"data_collection_only": schema.BoolAttribute{
						Description:         livenessVerifyDescription.Description,
						MarkdownDescription: livenessVerifyDescription.MarkdownDescription,
						Optional:            true,
						Computed:            true,
					},
				},
			},

			"created_at": schema.StringAttribute{
				Description: "Date and time the verify policy was created.",
				Computed:    true,
			},

			"updated_at": schema.StringAttribute{
				Description: "Date and time the verify policy was updated. Can be null.",
				Computed:    true,
			},
		},
	}
}

func (r *VerifyPolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *VerifyPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state verifyPolicyResourceModel

	if r.client == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	ctx = context.WithValue(ctx, verify.ContextServerVariables, map[string]string{
		"suffix": r.region.URLSuffix,
	})

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the model for the API
	VerifyPolicy, d := plan.expand(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	response, d := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return r.client.VerifyPoliciesApi.CreateVerifyPolicy(ctx, plan.EnvironmentId.ValueString()).VerifyPolicy(*VerifyPolicy).Execute()
		},
		"CreateVerifyPolicy",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
	)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the state to save
	state = plan

	// Save updated data into Terraform state
	resp.Diagnostics.Append(state.toState(response.(*verify.VerifyPolicy))...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *VerifyPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *verifyPolicyResourceModel

	if r.client == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	ctx = context.WithValue(ctx, verify.ContextServerVariables, map[string]string{
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
			return r.client.VerifyPoliciesApi.ReadOneVerifyPolicy(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
		},
		"ReadOneVerifyPolicy",
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
	resp.Diagnostics.Append(data.toState(response.(*verify.VerifyPolicy))...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VerifyPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state verifyPolicyResourceModel

	if r.client == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	ctx = context.WithValue(ctx, verify.ContextServerVariables, map[string]string{
		"suffix": r.region.URLSuffix,
	})

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the model for the API
	VerifyPolicy, d := plan.expand(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	response, diags := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return r.client.VerifyPoliciesApi.UpdateVerifyPolicy(ctx, plan.EnvironmentId.ValueString(), plan.Id.ValueString()).VerifyPolicy(*VerifyPolicy).Execute()
		},
		"UpdateVerifyPolicy",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
	)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update the state to save
	state = plan

	// Save updated data into Terraform state
	resp.Diagnostics.Append(state.toState(response.(*verify.VerifyPolicy))...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *VerifyPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *verifyPolicyResourceModel

	if r.client == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	ctx = context.WithValue(ctx, verify.ContextServerVariables, map[string]string{
		"suffix": r.region.URLSuffix,
	})

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	_, diags := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			r, err := r.client.VerifyPoliciesApi.DeleteVerifyPolicy(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return nil, r, err
		},
		"DeleteVerifyPolicy",
		framework.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *VerifyPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	splitLength := 2
	attributes := strings.SplitN(req.ID, "/", splitLength)

	if len(attributes) != splitLength {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("invalid id (\"%s\") specified, should be in format \"environment_id/verify_policy_id\"", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("environment_id"), attributes[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), attributes[2])...)
}

func (p *verifyPolicyResourceModel) expand(ctx context.Context) (*verify.VerifyPolicy, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := verify.NewVerifyPolicyWithDefaults()

	// Government Id Verification Object
	if !p.GovernmentId.IsNull() && !p.GovernmentId.IsUnknown() {

		var governmentId governmentIdModel
		d := p.GovernmentId.As(ctx, &governmentId, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		verifyGovernmentId, d := governmentId.expandgovernmentIdModel()
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}
		data.SetGovernmentId(*verifyGovernmentId)
	}

	// Facial Comparison Verification Object
	if !p.FacialComparison.IsNull() && !p.FacialComparison.IsUnknown() {

		var facialComparison facialComparisonModel
		d := p.FacialComparison.As(ctx, &facialComparison, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		verifyFacialComparison, d := facialComparison.expandFacialComparisonModel()
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}
		data.SetFacialComparison(*verifyFacialComparison)
	}

	// Liveness Verification Object
	if !p.Liveness.IsNull() && !p.Liveness.IsUnknown() {

		var liveness livenessnModel
		d := p.Liveness.As(ctx, &liveness, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		verifyLiveness, d := liveness.expandLivenessModel()
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}
		data.SetLiveness(*verifyLiveness)
	}

	// Transaction Object
	if !p.Transaction.IsNull() && !p.Transaction.IsUnknown() {

		var transaction transactionModel
		d := p.Transaction.As(ctx, &transaction, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		transactionSettings, d := transaction.expandTransactionModel(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		data.SetTransaction(*transactionSettings)
	}

	// Email Object
	if !p.Email.IsNull() && !p.Email.IsUnknown() {

		var emailConfiguration deviceModel
		d := p.Email.As(ctx, &emailConfiguration, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		emailSettings, d := emailConfiguration.expandDevice(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		data.SetEmail(*emailSettings)
	}

	// Phone Object
	if !p.Phone.IsNull() && !p.Phone.IsUnknown() {

		var phoneConfiguration deviceModel
		d := p.Phone.As(ctx, &phoneConfiguration, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		phoneSettings, d := phoneConfiguration.expandDevice(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		data.SetPhone(*phoneSettings)
	}

	// Top-level arguments
	data.SetId(p.Id.ValueString())

	environment := verify.NewObjectEnvironment()
	environment.SetId(p.EnvironmentId.ValueString())
	data.SetEnvironment(*environment)

	if !p.Name.IsNull() && !p.Name.IsUnknown() {
		data.SetName(p.Name.ValueString())
	}

	if !p.Default.IsNull() && !p.Default.IsUnknown() {
		data.SetDefault(p.Default.ValueBool())
	}

	if !p.Description.IsNull() && !p.Description.IsUnknown() {
		data.SetDescription(p.Description.ValueString())
	}

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
				"Verify Policy object was unexpectedly null on expansion.  Please report this to the provider maintainers.",
			)
		}
	}
	return data, diags
}

func (p *governmentIdModel) expandgovernmentIdModel() (*verify.GovernmentIdConfiguration, diag.Diagnostics) {
	var diags diag.Diagnostics

	verifyGovernmentId := verify.NewGovernmentIdConfigurationWithDefaults()
	if !p.Verify.IsNull() && !p.Verify.IsUnknown() {
		verifyGovernmentId.SetVerify(verify.EnumVerify(p.Verify.ValueString()))
	}

	if verifyGovernmentId == nil {
		diags.AddWarning(
			"Unexpected Value",
			"GovernmentId configuration object was unexpectedly null on expansion.  Please report this to the provider maintainers.",
		)
	}
	return verifyGovernmentId, diags

}

func (p *facialComparisonModel) expandFacialComparisonModel() (*verify.FacialComparisonConfiguration, diag.Diagnostics) {
	var diags diag.Diagnostics

	verifyFacialComparison := verify.NewFacialComparisonConfigurationWithDefaults()
	if !p.Verify.IsNull() && !p.Verify.IsUnknown() {
		verifyFacialComparison.SetVerify(verify.EnumVerify(p.Verify.ValueString()))
	}

	if !p.Threshold.IsNull() && !p.Threshold.IsUnknown() {
		verifyFacialComparison.SetThreshold(verify.EnumThreshold(p.Threshold.ValueString()))
	}

	if verifyFacialComparison == nil {
		diags.AddWarning(
			"Unexpected Value",
			"Facial Comparison configuration object was unexpectedly null on expansion.  Please report this to the provider maintainers.",
		)
	}
	return verifyFacialComparison, diags

}

func (p *livenessnModel) expandLivenessModel() (*verify.LivenessConfiguration, diag.Diagnostics) {
	var diags diag.Diagnostics

	verifyLiveness := verify.NewLivenessConfigurationWithDefaults()
	if !p.Verify.IsNull() && !p.Verify.IsUnknown() {
		verifyLiveness.SetVerify(verify.EnumVerify(p.Verify.ValueString()))
	}

	if !p.Threshold.IsNull() && !p.Threshold.IsUnknown() {
		verifyLiveness.SetThreshold(verify.EnumThreshold(p.Threshold.ValueString()))
	}

	if verifyLiveness == nil {
		diags.AddWarning(
			"Unexpected Value",
			"Liveness configuration object was unexpectedly null on expansion.  Please report this to the provider maintainers.",
		)
	}
	return verifyLiveness, diags

}

// i hate this function
func (p *transactionModel) expandTransactionModel(ctx context.Context) (*verify.TransactionConfiguration, diag.Diagnostics) {
	var diags diag.Diagnostics

	transactionSettings := verify.NewTransactionConfigurationWithDefaults()

	if !p.DataCollectionOnly.IsNull() && !p.DataCollectionOnly.IsUnknown() {
		transactionSettings.SetDataCollectionOnly(p.DataCollectionOnly.ValueBool())
	}

	if !p.DataCollection.IsNull() && !p.DataCollection.IsUnknown() {
		var dataCollection transactionDataCollectionModel
		d := p.DataCollection.As(ctx, &dataCollection, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		var genericTimeout genericTimeoutModel
		d = dataCollection.Timeout.As(ctx, &genericTimeout, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		dataCollectionTimeout := verify.NewTransactionConfigurationDataCollectionTimeoutWithDefaults()
		if !genericTimeout.TimeUnit.IsNull() && !genericTimeout.TimeUnit.IsUnknown() {
			dataCollectionTimeout.SetTimeUnit(verify.EnumShortTimeUnit(genericTimeout.TimeUnit.ValueString()))
		}

		if !genericTimeout.Duration.IsNull() && !genericTimeout.Duration.IsUnknown() {
			dataCollectionTimeout.SetDuration(int32(genericTimeout.Duration.ValueInt64()))
		}

		transactionDataCollection := verify.NewTransactionConfigurationDataCollection(*dataCollectionTimeout)
		transactionSettings.SetDataCollection(*transactionDataCollection)
	}

	if !p.Timeout.IsNull() && !p.Timeout.IsUnknown() {
		var genericTimeout genericTimeoutModel
		d := p.Timeout.As(ctx, &genericTimeout, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		transactionTimeout := verify.NewTransactionConfigurationTimeoutWithDefaults()
		transactionTimeout.SetTimeUnit(verify.EnumShortTimeUnit(genericTimeout.TimeUnit.ValueString()))
		transactionTimeout.SetDuration(int32(genericTimeout.Duration.ValueInt64()))

		transactionSettings.SetTimeout(*transactionTimeout)
	}

	if transactionSettings == nil {
		diags.AddWarning(
			"Unexpected Value",
			"Transaction configuration object was unexpectedly null on expansion.  Please report this to the provider maintainers.",
		)
	}
	return transactionSettings, diags

}

func (p *deviceModel) expandDevice(ctx context.Context) (*verify.EmailPhoneConfiguration, diag.Diagnostics) {
	var diags diag.Diagnostics

	deviceSettings := verify.NewEmailPhoneConfigurationWithDefaults()

	if !p.CreateMfaDevice.IsNull() && !p.CreateMfaDevice.IsUnknown() {
		deviceSettings.SetCreateMfaDevice(p.CreateMfaDevice.ValueBool())
	}

	if !p.Verify.IsNull() && !p.Verify.IsUnknown() {
		deviceSettings.SetVerify(verify.EnumVerify(p.Verify.ValueString()))
	}

	if !p.OTP.IsNull() && !p.OTP.IsUnknown() {
		var otpConfiguration otpConfigurationModel
		d := p.OTP.As(ctx, &otpConfiguration, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		// OTP Attempts
		// OTP Deliveries (also has cooldown object)

		// OTP LifeTime (generic timeout model)
		var genericTimeout genericTimeoutModel
		d = otpConfiguration.LifeTime.As(ctx, &genericTimeout, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}
		lifeTime := verify.NewEmailPhoneConfigurationOtpLifeTimeWithDefaults()
		if !genericTimeout.TimeUnit.IsNull() && !genericTimeout.TimeUnit.IsUnknown() {
			lifeTime.SetTimeUnit(verify.EnumLongTimeUnit(genericTimeout.TimeUnit.ValueString()))
		}

		if !genericTimeout.Duration.IsNull() && !genericTimeout.Duration.IsUnknown() {
			lifeTime.SetDuration(int32(genericTimeout.Duration.ValueInt64()))
		}

		// OTP Notification

		transactionDataCollection := verify.NewTransactionConfigurationDataCollection(*dataCollectionTimeout)
		transactionSettings.SetDataCollection(*transactionDataCollection)
	}

	if !p.Timeout.IsNull() && !p.Timeout.IsUnknown() {
		var genericTimeout genericTimeoutModel
		d := p.Timeout.As(ctx, &genericTimeout, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		transactionTimeout := verify.NewTransactionConfigurationTimeoutWithDefaults()
		transactionTimeout.SetTimeUnit(verify.EnumShortTimeUnit(genericTimeout.TimeUnit.ValueString()))
		transactionTimeout.SetDuration(int32(genericTimeout.Duration.ValueInt64()))

		deviceSettings.SetOtp(*transactionTimeout)
	}

	if transactionSettings == nil {
		diags.AddWarning(
			"Unexpected Value",
			"Transaction configuration object was unexpectedly null on expansion.  Please report this to the provider maintainers.",
		)
	}
	return transactionSettings, diags

}

func (p *verifyPolicyResourceModel) toState(apiObject *verify.VerifyPolicy) diag.Diagnostics {
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
	p.Default = framework.BoolOkToTF(apiObject.GetDefaultOk())
	p.Description = framework.StringOkToTF(apiObject.GetDescriptionOk())
	p.CreatedAt = framework.TimeOkToTF(apiObject.GetCreatedAtOk())
	p.UpdatedAt = framework.TimeOkToTF(apiObject.GetUpdatedAtOk())

	var d diag.Diagnostics
	p.GovernmentId, d = p.toStateGovernmentId(apiObject.GovernmentId)
	diags.Append(d...)

	p.FacialComparison, d = p.toStateFacialComparison(apiObject.FacialComparison)
	diags.Append(d...)

	p.Liveness, d = p.toStateLiveness(apiObject.Liveness)
	diags.Append(d...)

	p.Transaction, d = p.toStateTransaction(apiObject.Transaction)
	diags.Append(d...)

	return diags
}

func (p *verifyPolicyResourceModel) toStateGovernmentId(apiObject *verify.GovernmentIdConfiguration) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if apiObject == nil {
		return types.ObjectNull(governmentIdServiceTFObjectTypes), diags
	}

	objValue, d := types.ObjectValue(governmentIdServiceTFObjectTypes, map[string]attr.Value{
		"verify": framework.EnumOkToTF(apiObject.GetVerifyOk()),
	})
	diags.Append(d...)

	return objValue, diags
}

func (p *verifyPolicyResourceModel) toStateFacialComparison(apiObject *verify.FacialComparisonConfiguration) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if apiObject == nil {
		return types.ObjectNull(facialComparisonServiceTFObjectTypes), diags
	}

	objValue, d := types.ObjectValue(facialComparisonServiceTFObjectTypes, map[string]attr.Value{
		"verify":    framework.EnumOkToTF(apiObject.GetVerifyOk()),
		"threshold": framework.EnumOkToTF(apiObject.GetThresholdOk()),
	})
	diags.Append(d...)

	return objValue, diags
}

func (p *verifyPolicyResourceModel) toStateLiveness(apiObject *verify.LivenessConfiguration) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if apiObject == nil {
		return types.ObjectNull(livenessServiceTFObjectTypes), diags
	}

	objValue, d := types.ObjectValue(livenessServiceTFObjectTypes, map[string]attr.Value{
		"verify":    framework.EnumOkToTF(apiObject.GetVerifyOk()),
		"threshold": framework.EnumOkToTF(apiObject.GetThresholdOk()),
	})
	diags.Append(d...)

	return objValue, diags
}

func (p *verifyPolicyResourceModel) toStateTransaction(apiObject *verify.TransactionConfiguration) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if apiObject == nil {
		return types.ObjectNull(transactionServiceTFObjectTypes), diags
	}

	transactionTimeout := types.ObjectNull(genericTimeoutServiceTFObjectTypes)
	if v, ok := apiObject.GetTimeoutOk(); ok {
		var d diag.Diagnostics

		o := map[string]attr.Value{
			"duration":  framework.Int32OkToTF(v.GetDurationOk()),
			"time_unit": framework.EnumOkToTF(v.GetTimeUnitOk()),
		}

		objValue, d := types.ObjectValue(genericTimeoutServiceTFObjectTypes, o)
		diags.Append(d...)

		transactionTimeout = objValue
	}

	transactionDataCollection := types.ObjectNull(dataCollectionServiceTFObjectTypes)
	if v, ok := apiObject.GetDataCollectionOk(); ok {
		var d diag.Diagnostics

		transactionDataCollectionTimeout := types.ObjectNull(genericTimeoutServiceTFObjectTypes)
		if t, ok := v.GetTimeoutOk(); ok {
			o := map[string]attr.Value{
				"duration":  framework.Int32OkToTF(t.GetDurationOk()),
				"time_unit": framework.EnumOkToTF(t.GetTimeUnitOk()),
			}

			objValue, d := types.ObjectValue(genericTimeoutServiceTFObjectTypes, o)
			diags.Append(d...)

			transactionDataCollectionTimeout = objValue
		}

		o := map[string]attr.Value{
			"timeout": transactionDataCollectionTimeout,
		}

		objValue, d := types.ObjectValue(dataCollectionServiceTFObjectTypes, o)
		diags.Append(d...)

		transactionDataCollection = objValue

	}

	objValue, d := types.ObjectValue(transactionServiceTFObjectTypes, map[string]attr.Value{
		"timeout":              transactionTimeout,
		"data_collection":      transactionDataCollection,
		"data_collection_only": framework.BoolOkToTF(apiObject.GetDataCollectionOnlyOk()),
	})
	diags.Append(d...)

	return objValue, diags
}
