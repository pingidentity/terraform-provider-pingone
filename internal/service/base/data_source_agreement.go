package base

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type AgreementDataSource struct {
	client *management.APIClient
	region model.RegionMapping
}

type AgreementDataSourceModel struct {
	Id                    types.String  `tfsdk:"id"`
	EnvironmentId         types.String  `tfsdk:"environment_id"`
	AgreementId           types.String  `tfsdk:"agreement_id"`
	Name                  types.String  `tfsdk:"name"`
	Enabled               types.Bool    `tfsdk:"enabled"`
	Description           types.String  `tfsdk:"description"`
	ReconsentPeriodDays   types.Float64 `tfsdk:"reconsent_period_days"`
	TotalUserConsents     types.Int64   `tfsdk:"total_user_consent_count"`
	ExpiredUserConsents   types.Int64   `tfsdk:"expired_user_consent_count"`
	ConsentCountsUpdateAt types.String  `tfsdk:"consent_counts_updated_at"`
}

// Framework interfaces
var (
	_ datasource.DataSource = &AgreementDataSource{}
)

// New Object
func NewAgreementDataSource() datasource.DataSource {
	return &AgreementDataSource{}
}

// Metadata
func (r *AgreementDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_agreement"
}

// Schema
func (r *AgreementDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	nameLength := 1

	totalUserCountDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The total number of users who have consented to the agreement. This value is last calculated at the `consent_counts_updated_at` time.",
	)

	expiredUserCountDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The number of users who have consented to the agreement, but their consent has expired. This value is last calculated at the `consent_counts_updated_at` time.",
	)

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Datasource to retrieve details of an agreement configured in a PingOne environment.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment that is configured with the agreement."),
			),

			"agreement_id": schema.StringAttribute{
				Description: "The ID of the agreement to retrieve. Either `agreement_id`, or `name` can be used to retrieve the agreement localization, but cannot be set together.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("name")),
					verify.P1ResourceIDValidator(),
				},
			},

			"name": schema.StringAttribute{
				Description: "A string that specifies the name of the agreement to retrieve. Either `agreement_id`, or `name` can be used to retrieve the agreement localization, but cannot be set together.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("agreement_id")),
					stringvalidator.LengthAtLeast(nameLength),
				},
			},

			"description": schema.StringAttribute{
				Description: "A string that specifies the description of the agreement.",
				Computed:    true,
			},

			"enabled": schema.BoolAttribute{
				Description: "The current enabled state of the agreement.",
				Computed:    true,
			},

			"reconsent_period_days": schema.Float64Attribute{
				Description: "A number that specifies the number of days until a consent to this agreement expires.",
				Computed:    true,
			},

			"total_user_consent_count": schema.Int64Attribute{
				Description:         totalUserCountDescription.Description,
				MarkdownDescription: totalUserCountDescription.MarkdownDescription,
				Computed:            true,
			},

			"expired_user_consent_count": schema.Int64Attribute{
				Description:         expiredUserCountDescription.Description,
				MarkdownDescription: expiredUserCountDescription.MarkdownDescription,
				Computed:            true,
			},

			"consent_counts_updated_at": schema.StringAttribute{
				Description: "The date and time the consent user count metrics were last updated. This value is typically updated once every 24 hours.",
				Computed:    true,
			},
		},
	}
}

func (r *AgreementDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *AgreementDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *AgreementDataSourceModel

	if r.client == nil {
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

	var agreement management.Agreement

	if !data.Name.IsNull() {

		// Run the API call
		response, diags := framework.ParseResponse(
			ctx,

			func() (interface{}, *http.Response, error) {
				return r.client.AgreementsResourcesApi.ReadAllAgreements(ctx, data.EnvironmentId.ValueString()).Execute()
			},
			"ReadAllAgreements",
			framework.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
		)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		entityArray := response.(*management.EntityArray)

		if agreements, ok := entityArray.Embedded.GetAgreementsOk(); ok {

			found := false
			for _, agreementItem := range agreements {

				if agreementItem.GetName() == data.Name.ValueString() {
					agreement = agreementItem
					found = true
					break
				}
			}

			if !found {
				resp.Diagnostics.AddError(
					"Cannot find agreement from name",
					fmt.Sprintf("The agreement %s for environment %s cannot be found", data.Name.String(), data.EnvironmentId.String()),
				)
				return
			}

		}

	} else if !data.AgreementId.IsNull() {

		// Run the API call
		response, diags := framework.ParseResponse(
			ctx,

			func() (interface{}, *http.Response, error) {
				return r.client.AgreementsResourcesApi.ReadOneAgreement(ctx, data.EnvironmentId.ValueString(), data.AgreementId.ValueString()).Execute()
			},
			"ReadOneAgreement",
			framework.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
		)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		agreement = *response.(*management.Agreement)
	} else {
		resp.Diagnostics.AddError(
			"Missing parameter",
			"Cannot find the requested agreement. agreement_id or name must be set.",
		)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(&agreement)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (p *AgreementDataSourceModel) toState(apiObject *management.Agreement) diag.Diagnostics {
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
	p.Name = framework.StringOkToTF(apiObject.GetNameOk())
	p.Enabled = framework.BoolOkToTF(apiObject.GetEnabledOk())
	p.Description = framework.StringOkToTF(apiObject.GetDescriptionOk())
	p.ReconsentPeriodDays = framework.Float32OkToTF(apiObject.GetReconsentPeriodDaysOk())
	p.TotalUserConsents = framework.Int32OkToTF(apiObject.GetTotalConsentsOk())
	p.ExpiredUserConsents = framework.Int32OkToTF(apiObject.GetTotalExpiredConsentsOk())
	p.ConsentCountsUpdateAt = framework.TimeOkToTF(apiObject.GetConsentsAggregatedAtOk())

	return diags
}
