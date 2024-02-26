package base

import (
	"context"
	"encoding/base64"
	"fmt"
	"math/big"
	"net/http"
	"regexp"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	stringvalidatorinternal "github.com/pingidentity/terraform-provider-pingone/internal/framework/stringvalidator"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/utils"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type KeyResource serviceClientType

type keyResourceModel struct {
	Id                 types.String `tfsdk:"id"`
	EnvironmentId      types.String `tfsdk:"environment_id"`
	Name               types.String `tfsdk:"name"`
	Algorithm          types.String `tfsdk:"algorithm"`
	Default            types.Bool   `tfsdk:"default"`
	ExpiresAt          types.String `tfsdk:"expires_at"`
	IssuerDn           types.String `tfsdk:"issuer_dn"`
	KeyLength          types.Int64  `tfsdk:"key_length"`
	SerialNumber       types.String `tfsdk:"serial_number"`
	SignatureAlgorithm types.String `tfsdk:"signature_algorithm"`
	StartsAt           types.String `tfsdk:"starts_at"`
	Status             types.String `tfsdk:"status"`
	SubjectDn          types.String `tfsdk:"subject_dn"`
	UsageType          types.String `tfsdk:"usage_type"`
	ValidityPeriod     types.Int64  `tfsdk:"validity_period"`
	CustomCrl          types.String `tfsdk:"custom_crl"`
	PKCS12FileBase64   types.String `tfsdk:"pkcs12_file_base64"`
	PKCS12FilePassword types.String `tfsdk:"pkcs12_file_password"`
}

// Framework interfaces
var (
	_ resource.Resource                   = &KeyResource{}
	_ resource.ResourceWithConfigure      = &KeyResource{}
	_ resource.ResourceWithValidateConfig = &KeyResource{}
	_ resource.ResourceWithImportState    = &KeyResource{}
)

var (
	allowedKeyLengthsRSA = []int64{2048, 3072, 4096, 7680}
	allowedKeyLengthsEC  = []int64{224, 256, 384, 521}
)

// New Object
func NewKeyResource() resource.Resource {
	return &KeyResource{}
}

// Metadata
func (r *KeyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_key"
}

// Schema.
func (r *KeyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	const attrMinLength = 1
	const rotationPeriodMinimum = 30
	const rotationPeriodDefault = 90
	const validityPeriodDefault = 365

	nameDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the system name of the key.",
	).ConflictsWith([]string{"pkcs12_file_base64", "pkcs12_file_password"}).RequiresReplace()

	algorithmDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the key algorithm.",
	).AllowedValuesEnum(management.AllowedEnumCertificateKeyAlgorithmEnumValues).ConflictsWith([]string{"pkcs12_file_base64"}).RequiresReplace()

	defaultDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether this is the default key for the specified environment.",
	).DefaultValue(false)

	issuerDnDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the distinguished name of the certificate issuer.",
	).ConflictsWith([]string{"pkcs12_file_base64", "pkcs12_file_password"}).RequiresReplace()

	keyLengthDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An integer that specifies the key length. For RSA keys, options are `2048`, `3072`, `4096` and `7680`. For elliptical curve (EC) keys, options are `224`, `256`, `384` and `521`.",
	).ConflictsWith([]string{"pkcs12_file_base64", "pkcs12_file_password"}).RequiresReplace()

	serialNumberDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An integer (in string data type) that specifies the serial number of the key or certificate.",
	).ConflictsWith([]string{"pkcs12_file_base64", "pkcs12_file_password"}).RequiresReplace()

	signatureAlgorithmDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("A string that specifies the signature algorithm of the key. For RSA keys, options are `%s`, `%s` and `%s`. For elliptical curve (EC) keys, options are `%s`, `%s` and `%s`.", string(management.ENUMCERTIFICATEKEYSIGNAGUREALGORITHM_SHA256WITH_RSA), string(management.ENUMCERTIFICATEKEYSIGNAGUREALGORITHM_SHA384WITH_RSA), string(management.ENUMCERTIFICATEKEYSIGNAGUREALGORITHM_SHA512WITH_RSA), string(management.ENUMCERTIFICATEKEYSIGNAGUREALGORITHM_SHA256WITH_ECDSA), string(management.ENUMCERTIFICATEKEYSIGNAGUREALGORITHM_SHA384WITH_ECDSA), string(management.ENUMCERTIFICATEKEYSIGNAGUREALGORITHM_SHA512WITH_ECDSA)),
	).ConflictsWith([]string{"pkcs12_file_base64", "pkcs12_file_password"}).RequiresReplace()

	statusDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the status of the key.",
	).AllowedValuesEnum(management.AllowedEnumCertificateKeyStatusEnumValues)

	subjectDnDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the distinguished name of the subject being secured.",
	).ConflictsWith([]string{"pkcs12_file_base64", "pkcs12_file_password"}).RequiresReplace()

	usageTypeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies how the certificate is used.",
	).AllowedValuesEnum(management.AllowedEnumCertificateKeyUsageTypeEnumValues)

	validityPeriodDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An integer that specifies the number of days the key is valid.",
	).ConflictsWith([]string{"pkcs12_file_base64", "pkcs12_file_password"}).RequiresReplace()

	customCrlDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A URL string of a custom Certificate Revokation List endpoint.  Used for certificates of type `ISSUANCE`.",
	)

	pkcs12FileBase64Description := framework.SchemaAttributeDescriptionFromMarkdown(
		"A base64 encoded PKCS12 file to import.",
	).ConflictsWith([]string{"name", "algorithm", "issuer_dn", "key_length", "serial_number", "signature_algorithm", "subject_dn", "validity_period", "custom_crl"}).RequiresReplace()

	pkcs12FilePasswordDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the password to decrypt the PKCS12 file, if it is encrypted.  Optional if `pkcs12_file_base64` is defined.",
	).ConflictsWith([]string{"name", "algorithm", "issuer_dn", "key_length", "serial_number", "signature_algorithm", "subject_dn", "validity_period", "custom_crl"}).RequiresReplace()

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage PingOne keys for an environment.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment to manage the key in."),
			),

			"name": schema.StringAttribute{
				Description:         nameDescription.Description,
				MarkdownDescription: nameDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,

				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},

				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
					stringvalidator.ConflictsWith(
						path.MatchRelative().AtParent().AtName("pkcs12_file_base64"),
						path.MatchRelative().AtParent().AtName("pkcs12_file_password"),
					),
					stringvalidator.AlsoRequires(
						path.MatchRelative().AtParent().AtName("name"),
						path.MatchRelative().AtParent().AtName("algorithm"),
						path.MatchRelative().AtParent().AtName("key_length"),
						path.MatchRelative().AtParent().AtName("signature_algorithm"),
						path.MatchRelative().AtParent().AtName("subject_dn"),
						path.MatchRelative().AtParent().AtName("validity_period"),
					),
				},
			},

			"algorithm": schema.StringAttribute{
				Description:         algorithmDescription.Description,
				MarkdownDescription: algorithmDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,

				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},

				Validators: []validator.String{
					stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumCertificateKeyAlgorithmEnumValues)...),
					stringvalidator.ConflictsWith(
						path.MatchRelative().AtParent().AtName("pkcs12_file_base64"),
						path.MatchRelative().AtParent().AtName("pkcs12_file_password"),
					),
					stringvalidator.AlsoRequires(
						path.MatchRelative().AtParent().AtName("name"),
						path.MatchRelative().AtParent().AtName("algorithm"),
						path.MatchRelative().AtParent().AtName("key_length"),
						path.MatchRelative().AtParent().AtName("signature_algorithm"),
						path.MatchRelative().AtParent().AtName("subject_dn"),
						path.MatchRelative().AtParent().AtName("validity_period"),
					),
				},
			},

			"default": schema.BoolAttribute{
				Description:         defaultDescription.Description,
				MarkdownDescription: defaultDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,

				Default: booldefault.StaticBool(false),
			},

			"expires_at": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the date and time the key resource expires.").Description,
				Computed:    true,

				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"issuer_dn": schema.StringAttribute{
				Description:         issuerDnDescription.Description,
				MarkdownDescription: issuerDnDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,

				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},

				Validators: []validator.String{
					stringvalidator.ConflictsWith(
						path.MatchRelative().AtParent().AtName("pkcs12_file_base64"),
						path.MatchRelative().AtParent().AtName("pkcs12_file_password"),
					),
				},
			},

			"key_length": schema.Int64Attribute{
				Description:         keyLengthDescription.Description,
				MarkdownDescription: keyLengthDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,

				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},

				Validators: []validator.Int64{
					int64validator.Any(
						int64validator.OneOf(
							allowedKeyLengthsRSA...,
						),
						int64validator.OneOf(
							allowedKeyLengthsEC...,
						),
					),
					int64validator.ConflictsWith(
						path.MatchRelative().AtParent().AtName("pkcs12_file_base64"),
						path.MatchRelative().AtParent().AtName("pkcs12_file_password"),
					),
					int64validator.AlsoRequires(
						path.MatchRelative().AtParent().AtName("name"),
						path.MatchRelative().AtParent().AtName("algorithm"),
						path.MatchRelative().AtParent().AtName("key_length"),
						path.MatchRelative().AtParent().AtName("signature_algorithm"),
						path.MatchRelative().AtParent().AtName("subject_dn"),
						path.MatchRelative().AtParent().AtName("validity_period"),
					),
				},
			},

			"serial_number": schema.StringAttribute{
				Description:         serialNumberDescription.Description,
				MarkdownDescription: serialNumberDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,

				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},

				Validators: []validator.String{
					stringvalidator.ConflictsWith(
						path.MatchRelative().AtParent().AtName("pkcs12_file_base64"),
						path.MatchRelative().AtParent().AtName("pkcs12_file_password"),
					),
				},
			},

			"signature_algorithm": schema.StringAttribute{
				Description:         signatureAlgorithmDescription.Description,
				MarkdownDescription: signatureAlgorithmDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,

				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},

				Validators: []validator.String{
					stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumCertificateKeySignagureAlgorithmEnumValues)...),
					stringvalidator.ConflictsWith(
						path.MatchRelative().AtParent().AtName("pkcs12_file_base64"),
						path.MatchRelative().AtParent().AtName("pkcs12_file_password"),
					),
					stringvalidator.AlsoRequires(
						path.MatchRelative().AtParent().AtName("name"),
						path.MatchRelative().AtParent().AtName("algorithm"),
						path.MatchRelative().AtParent().AtName("key_length"),
						path.MatchRelative().AtParent().AtName("signature_algorithm"),
						path.MatchRelative().AtParent().AtName("subject_dn"),
						path.MatchRelative().AtParent().AtName("validity_period"),
					),
				},
			},

			"starts_at": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the date and time the validity period starts.").Description,
				Computed:    true,

				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"status": schema.StringAttribute{
				Description:         statusDescription.Description,
				MarkdownDescription: statusDescription.MarkdownDescription,
				Computed:            true,

				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"subject_dn": schema.StringAttribute{
				Description:         subjectDnDescription.Description,
				MarkdownDescription: subjectDnDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,

				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},

				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
					stringvalidator.ConflictsWith(
						path.MatchRelative().AtParent().AtName("pkcs12_file_base64"),
						path.MatchRelative().AtParent().AtName("pkcs12_file_password"),
					),
					stringvalidator.AlsoRequires(
						path.MatchRelative().AtParent().AtName("name"),
						path.MatchRelative().AtParent().AtName("algorithm"),
						path.MatchRelative().AtParent().AtName("key_length"),
						path.MatchRelative().AtParent().AtName("signature_algorithm"),
						path.MatchRelative().AtParent().AtName("subject_dn"),
						path.MatchRelative().AtParent().AtName("validity_period"),
					),
				},
			},

			"usage_type": schema.StringAttribute{
				Description:         usageTypeDescription.Description,
				MarkdownDescription: usageTypeDescription.MarkdownDescription,
				Required:            true,

				Validators: []validator.String{
					stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumCertificateKeyUsageTypeEnumValues)...),
				},
			},

			"validity_period": schema.Int64Attribute{
				Description:         validityPeriodDescription.Description,
				MarkdownDescription: validityPeriodDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,

				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},

				Validators: []validator.Int64{
					int64validator.AtLeast(attrMinLength),
					int64validator.ConflictsWith(
						path.MatchRelative().AtParent().AtName("pkcs12_file_base64"),
						path.MatchRelative().AtParent().AtName("pkcs12_file_password"),
					),
					int64validator.AlsoRequires(
						path.MatchRelative().AtParent().AtName("name"),
						path.MatchRelative().AtParent().AtName("algorithm"),
						path.MatchRelative().AtParent().AtName("key_length"),
						path.MatchRelative().AtParent().AtName("signature_algorithm"),
						path.MatchRelative().AtParent().AtName("subject_dn"),
						path.MatchRelative().AtParent().AtName("validity_period"),
					),
				},
			},

			"custom_crl": schema.StringAttribute{
				Description:         customCrlDescription.Description,
				MarkdownDescription: customCrlDescription.MarkdownDescription,
				Optional:            true,

				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(`^http:\/\/[a-zA-Z0-9.-\/]*$`), "`custom_crl` must be a `http://` URL endpoint."),
					stringvalidator.ConflictsWith(
						path.MatchRelative().AtParent().AtName("pkcs12_file_base64"),
						path.MatchRelative().AtParent().AtName("pkcs12_file_password"),
					),
				},
			},

			"pkcs12_file_base64": schema.StringAttribute{
				Description:         pkcs12FileBase64Description.Description,
				MarkdownDescription: pkcs12FileBase64Description.MarkdownDescription,
				Optional:            true,
				Sensitive:           true,

				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},

				Validators: []validator.String{
					stringvalidatorinternal.IsBase64Encoded(),
					stringvalidator.ConflictsWith(
						path.MatchRelative().AtParent().AtName("name"),
						path.MatchRelative().AtParent().AtName("algorithm"),
						path.MatchRelative().AtParent().AtName("issuer_dn"),
						path.MatchRelative().AtParent().AtName("key_length"),
						path.MatchRelative().AtParent().AtName("serial_number"),
						path.MatchRelative().AtParent().AtName("signature_algorithm"),
						path.MatchRelative().AtParent().AtName("subject_dn"),
						path.MatchRelative().AtParent().AtName("validity_period"),
						path.MatchRelative().AtParent().AtName("custom_crl"),
					),
				},
			},

			"pkcs12_file_password": schema.StringAttribute{
				Description:         pkcs12FilePasswordDescription.Description,
				MarkdownDescription: pkcs12FilePasswordDescription.MarkdownDescription,
				Optional:            true,
				Sensitive:           true,

				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},

				Validators: []validator.String{
					stringvalidator.AlsoRequires(
						path.MatchRelative().AtParent().AtName("pkcs12_file_base64"),
					),
					stringvalidator.ConflictsWith(
						path.MatchRelative().AtParent().AtName("name"),
						path.MatchRelative().AtParent().AtName("algorithm"),
						path.MatchRelative().AtParent().AtName("issuer_dn"),
						path.MatchRelative().AtParent().AtName("key_length"),
						path.MatchRelative().AtParent().AtName("serial_number"),
						path.MatchRelative().AtParent().AtName("signature_algorithm"),
						path.MatchRelative().AtParent().AtName("subject_dn"),
						path.MatchRelative().AtParent().AtName("validity_period"),
						path.MatchRelative().AtParent().AtName("custom_crl"),
					),
				},
			},
		},
	}
}

func (p *KeyResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data keyResourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var validKeyLengths string
	keyLengthValidationError := false
	if data.Algorithm.Equal(types.StringValue(string(management.ENUMCERTIFICATEKEYALGORITHM_RSA))) && !slices.Contains(allowedKeyLengthsRSA, data.KeyLength.ValueInt64()) {

		keyLengthStrs := make([]string, len(allowedKeyLengthsRSA))
		for i, v := range allowedKeyLengthsRSA {
			keyLengthStrs[i] = fmt.Sprintf("%d", v)
		}

		keyLengthValidationError = true
		validKeyLengths = strings.Join(keyLengthStrs, "`, `")
	}

	if data.Algorithm.Equal(types.StringValue(string(management.ENUMCERTIFICATEKEYALGORITHM_EC))) && !slices.Contains(allowedKeyLengthsEC, data.KeyLength.ValueInt64()) {

		keyLengthStrs := make([]string, len(allowedKeyLengthsEC))
		for i, v := range allowedKeyLengthsEC {
			keyLengthStrs[i] = fmt.Sprintf("%d", v)
		}

		keyLengthValidationError = true
		validKeyLengths = strings.Join(keyLengthStrs, "`, `")
	}

	if keyLengthValidationError {
		resp.Diagnostics.AddAttributeError(
			path.Root("key_length"),
			"Invalid attribute combination",
			fmt.Sprintf("When using an `algorithm` value of `%s`, only the following key lengths are valid: `%s`.", data.Algorithm.ValueString(), validKeyLengths),
		)
	}

	if !data.CustomCrl.IsNull() && !data.UsageType.Equal(types.StringValue(string(management.ENUMCERTIFICATEKEYUSAGETYPE_ISSUANCE))) {
		resp.Diagnostics.AddAttributeError(
			path.Root("custom_crl"),
			"Invalid attribute combination",
			"`custom_crl` can only be set for keys that have a `type` value of `ISSUANCE`.",
		)
	}
}

func (r *KeyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *KeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state keyResourceModel

	if r.Client.ManagementAPIClient == nil {
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

	var response *management.Certificate
	if !plan.PKCS12FileBase64.IsNull() && !plan.PKCS12FileBase64.IsUnknown() {

		archive, err := base64.StdEncoding.DecodeString(plan.PKCS12FileBase64.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Cannot base64 decode provided PKCS12 key file.",
				fmt.Sprintf("Please ensure the PKCS12 key file is base64 encoded.  Error: %s", err.Error()),
			)

			return
		}

		var archivePassword *string

		if !plan.PKCS12FilePassword.IsNull() && !plan.PKCS12FilePassword.IsUnknown() {
			archivePassword = plan.PKCS12FilePassword.ValueStringPointer()
		}

		// Run the API call
		resp.Diagnostics.Append(framework.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				request := r.Client.ManagementAPIClient.CertificateManagementApi.CreateKey(ctx, plan.EnvironmentId.ValueString()).ContentType("multipart/form-data").UsageType(plan.UsageType.ValueString()).File(&archive)

				if archivePassword != nil {
					request = request.Password(*archivePassword)
				}

				fO, fR, fErr := request.Execute()
				return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
			},
			"CreateKey",
			framework.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
			&response,
		)...)
	} else {
		// Validate
		resp.Diagnostics.Append(plan.validate()...)
		if resp.Diagnostics.HasError() {
			return
		}

		// Build the model for the API
		certificateKey := plan.expand()

		// Run the API call
		resp.Diagnostics.Append(framework.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				fO, fR, fErr := r.Client.ManagementAPIClient.CertificateManagementApi.CreateKey(ctx, plan.EnvironmentId.ValueString()).Certificate(*certificateKey).Execute()
				return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
			},
			"CreateKey",
			framework.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
			&response,
		)...)
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

func (r *KeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *keyResourceModel

	if r.Client.ManagementAPIClient == nil {
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
	var response *management.Certificate
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.CertificateManagementApi.GetKey(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"GetKey",
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

func (r *KeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state keyResourceModel

	if r.Client.ManagementAPIClient == nil {
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

	// Validate
	resp.Diagnostics.Append(plan.validate()...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the model for the API
	certificateKey := plan.expandUpdate()

	// Run the API call
	var response *management.Certificate
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.CertificateManagementApi.UpdateKey(ctx, plan.EnvironmentId.ValueString(), plan.Id.ValueString()).CertificateKeyUpdate(*certificateKey).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdateKey",
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

func (r *KeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *keyResourceModel

	if r.Client.ManagementAPIClient == nil {
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
			fR, fErr := r.Client.ManagementAPIClient.CertificateManagementApi.DeleteKey(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), nil, fR, fErr)
		},
		"DeleteKey",
		framework.CustomErrorResourceNotFoundWarning,
		func(ctx context.Context, r *http.Response, p1error *model.P1Error) bool {

			if p1error != nil {
				var err error

				// It seems the key might not release itself immediately
				if m, err := regexp.MatchString("The Key must not be in use", p1error.GetMessage()); err == nil && m {
					tflog.Warn(ctx, "Key in use detected")
					return true
				}
				if err != nil {
					tflog.Warn(ctx, "Cannot match error string for retry (DeleteKey)")
					return false
				}

			}

			return false
		},
		nil,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *KeyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	idComponents := []framework.ImportComponent{
		{
			Label:  "environment_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
		{
			Label:     "key_id",
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

func (p *keyResourceModel) expand() *management.Certificate {

	usageType := management.EnumCertificateKeyUsageType(p.UsageType.ValueString())

	data := management.NewCertificate(
		management.EnumCertificateKeyAlgorithm(p.Algorithm.ValueString()),
		int32(p.KeyLength.ValueInt64()),
		p.Name.ValueString(),
		management.EnumCertificateKeySignagureAlgorithm(p.SignatureAlgorithm.ValueString()),
		p.SubjectDn.ValueString(),
		usageType,
		int32(p.ValidityPeriod.ValueInt64()),
	)

	if !p.CustomCrl.IsNull() && !p.CustomCrl.IsUnknown() {
		data.SetCustomCRL(p.CustomCrl.ValueString())
	}

	if !p.Default.IsNull() && !p.Default.IsUnknown() {
		data.SetDefault(p.Default.ValueBool())
	}

	if !p.IssuerDn.IsNull() && !p.IssuerDn.IsUnknown() {
		data.SetIssuerDN(p.IssuerDn.ValueString())
	}

	if !p.SerialNumber.IsNull() && !p.SerialNumber.IsUnknown() {
		if j, ok := new(big.Int).SetString(p.SerialNumber.ValueString(), 0); ok {
			data.SetSerialNumber(*j)
		}
	}

	return data
}

func (p *keyResourceModel) expandUpdate() *management.CertificateKeyUpdate {

	data := management.NewCertificateKeyUpdate(p.Default.ValueBool(), management.EnumCertificateKeyUsageType(p.UsageType.ValueString()))

	if !p.IssuerDn.IsNull() && !p.IssuerDn.IsUnknown() {
		data.SetIssuerDN(p.IssuerDn.ValueString())
	}

	return data
}

func (p *keyResourceModel) validate() diag.Diagnostics {
	var diags diag.Diagnostics

	if !p.CustomCrl.IsNull() && !p.CustomCrl.IsUnknown() && !p.UsageType.Equal(types.StringValue(string(management.ENUMCERTIFICATEKEYUSAGETYPE_ISSUANCE))) {
		diags.AddAttributeError(
			path.Root("custom_crl"),
			"Invalid attribute combination",
			"`custom_crl` can only be set for keys that have a `type` value of `ISSUANCE`.",
		)
	}

	return diags
}

func (p *keyResourceModel) toState(apiObject *management.Certificate) diag.Diagnostics {
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
	p.Default = framework.BoolOkToTF(apiObject.GetDefaultOk())
	p.ExpiresAt = framework.TimeOkToTF(apiObject.GetExpiresAtOk())
	p.IssuerDn = framework.StringOkToTF(apiObject.GetIssuerDNOk())
	p.KeyLength = framework.Int32OkToTF(apiObject.GetKeyLengthOk())

	if v, ok := apiObject.GetSerialNumberOk(); ok {
		p.SerialNumber = framework.StringToTF(v.String())
	} else {
		p.SerialNumber = types.StringNull()
	}

	p.SignatureAlgorithm = framework.EnumOkToTF(apiObject.GetSignatureAlgorithmOk())
	p.StartsAt = framework.TimeOkToTF(apiObject.GetStartsAtOk())
	p.Status = framework.EnumOkToTF(apiObject.GetStatusOk())
	p.SubjectDn = framework.StringOkToTF(apiObject.GetSubjectDNOk())
	p.UsageType = framework.EnumOkToTF(apiObject.GetUsageTypeOk())
	p.ValidityPeriod = framework.Int32OkToTF(apiObject.GetValidityPeriodOk())
	p.CustomCrl = framework.StringOkToTF(apiObject.GetCustomCRLOk())

	return diags
}
