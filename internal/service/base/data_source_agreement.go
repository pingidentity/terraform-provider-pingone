package base

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
)

// Types
type AgreementDataSource serviceClientType

type AgreementDataSourceModel struct {
	Id                    pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId         pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	AgreementId           pingonetypes.ResourceIDValue `tfsdk:"agreement_id"`
	Name                  types.String                 `tfsdk:"name"`
	Enabled               types.Bool                   `tfsdk:"enabled"`
	Description           types.String                 `tfsdk:"description"`
	ReconsentPeriodDays   types.Float64                `tfsdk:"reconsent_period_days"`
	TotalUserConsents     types.Int64                  `tfsdk:"total_user_consent_count"`
	ExpiredUserConsents   types.Int64                  `tfsdk:"expired_user_consent_count"`
	ConsentCountsUpdateAt timetypes.RFC3339            `tfsdk:"consent_counts_updated_at"`
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

				CustomType: pingonetypes.ResourceIDType{},

				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("name")),
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

				CustomType: timetypes.RFC3339Type{},
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

	r.Client = resourceConfig.Client.API
	if r.Client == nil {
		resp.Diagnostics.AddError(
			"Client not initialised",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.",
		)
		return
	}
}

func (r *AgreementDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *AgreementDataSourceModel

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

	var agreement *management.Agreement

	if !data.Name.IsNull() {

		// Run the API call
		resp.Diagnostics.Append(framework.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				pagedIterator := r.Client.ManagementAPIClient.AgreementsResourcesApi.ReadAllAgreements(ctx, data.EnvironmentId.ValueString()).Execute()
				var initialHttpResponse *http.Response

				for pageCursor, err := range pagedIterator {
					if err != nil {
						return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), nil, pageCursor.HTTPResponse, err)
					}

					if initialHttpResponse == nil {
						initialHttpResponse = pageCursor.HTTPResponse
					}

					if agreements, ok := pageCursor.EntityArray.Embedded.GetAgreementsOk(); ok {

						for _, agreementItem := range agreements {
							if agreementItem.GetName() == data.Name.ValueString() {
								return &agreementItem, pageCursor.HTTPResponse, nil
							}
						}
					}
				}

				return nil, initialHttpResponse, nil
			},
			"ReadAllAgreements",
			framework.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
			&agreement,
		)...)
		if resp.Diagnostics.HasError() {
			return
		}

		if agreement == nil {
			resp.Diagnostics.AddError(
				"Cannot find agreement from name",
				fmt.Sprintf("The agreement %s for environment %s cannot be found", data.Name.String(), data.EnvironmentId.String()),
			)
			return
		}

	} else if !data.AgreementId.IsNull() {

		// Run the API call
		resp.Diagnostics.Append(framework.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				fO, fR, fErr := r.Client.ManagementAPIClient.AgreementsResourcesApi.ReadOneAgreement(ctx, data.EnvironmentId.ValueString(), data.AgreementId.ValueString()).Execute()
				return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
			},
			"ReadOneAgreement",
			framework.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
			&agreement,
		)...)
		if resp.Diagnostics.HasError() {
			return
		}

	} else {
		resp.Diagnostics.AddError(
			"Missing parameter",
			"Cannot find the requested agreement. agreement_id or name must be set.",
		)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(agreement)...)
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

	p.Id = framework.PingOneResourceIDToTF(apiObject.GetId())
	p.EnvironmentId = framework.PingOneResourceIDToTF(*apiObject.GetEnvironment().Id)
	p.AgreementId = framework.PingOneResourceIDToTF(apiObject.GetId())
	p.Name = framework.StringOkToTF(apiObject.GetNameOk())
	p.Enabled = framework.BoolOkToTF(apiObject.GetEnabledOk())
	p.Description = framework.StringOkToTF(apiObject.GetDescriptionOk())
	p.ReconsentPeriodDays = framework.Float32OkToTF(apiObject.GetReconsentPeriodDaysOk())
	p.TotalUserConsents = framework.Int32OkToTF(apiObject.GetTotalConsentsOk())
	p.ExpiredUserConsents = framework.Int32OkToTF(apiObject.GetTotalExpiredConsentsOk())
	p.ConsentCountsUpdateAt = framework.TimeOkToTF(apiObject.GetConsentsAggregatedAtOk())

	return diags
}
