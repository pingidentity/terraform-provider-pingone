package sso

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
type SchemaDataSource struct {
	client *management.APIClient
	region model.RegionMapping
}

type SchemaDataSourceModel struct {
	Id            types.String `tfsdk:"id"`
	EnvironmentId types.String `tfsdk:"environment_id"`
	SchemaId      types.String `tfsdk:"schema_id"`
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
}

// Framework interfaces
var (
	_ datasource.DataSource = &SchemaDataSource{}
)

// New Object
func NewSchemaDataSource() datasource.DataSource {
	return &SchemaDataSource{}
}

// Metadata
func (r *SchemaDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_schema"
}

// Schema
func (r *SchemaDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	nameLength := 1

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Datasource to read PingOne schema data.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(framework.SchemaDescription{
				Description: "The ID of the environment that is configured with the schema."},
			),

			"schema_id": schema.StringAttribute{
				Description: "The ID of the schema.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("name")),
					verify.P1ResourceIDValidator(),
				},
			},

			"name": schema.StringAttribute{
				Description: "The name of the schema.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("schema_id")),
					stringvalidator.LengthAtLeast(nameLength),
				},
			},

			"description": schema.StringAttribute{
				Description: "A description of the schema.",
				Computed:    true,
			},
		},
	}
}

func (r *SchemaDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *SchemaDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *SchemaDataSourceModel

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
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var schema management.Schema

	if !data.Name.IsNull() {

		// Run the API call
		response, diags := framework.ParseResponse(
			ctx,

			func() (interface{}, *http.Response, error) {
				return r.client.SchemasApi.ReadAllSchemas(ctx, data.EnvironmentId.ValueString()).Execute()
			},
			"ReadAllSchemas",
			framework.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
		)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		entityArray := response.(*management.EntityArray)

		if schemas, ok := entityArray.Embedded.GetSchemasOk(); ok {

			found := false
			for _, schemaItem := range schemas {

				if schemaItem.GetName() == data.Name.ValueString() {
					schema = schemaItem
					found = true
					break
				}
			}

			if !found {
				resp.Diagnostics.AddError(
					"Cannot find schema from name",
					fmt.Sprintf("The schema %s for environment %s cannot be found", data.Name.String(), data.EnvironmentId.String()),
				)
				return
			}

		}

	} else if !data.SchemaId.IsNull() {

		// Run the API call
		response, diags := framework.ParseResponse(
			ctx,

			func() (interface{}, *http.Response, error) {
				return r.client.SchemasApi.ReadOneSchema(ctx, data.EnvironmentId.ValueString(), data.SchemaId.ValueString()).Execute()
			},
			"ReadOneSchema",
			framework.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
		)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		schema = *response.(*management.Schema)
	} else {
		resp.Diagnostics.AddError(
			"Missing parameter",
			"Cannot find the requested schema. schema_id or name must be set.",
		)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(&schema)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (p *SchemaDataSourceModel) toState(apiObject *management.Schema) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.Id = framework.StringToTF(apiObject.GetId())
	p.SchemaId = framework.StringToTF(apiObject.GetId())
	p.Name = framework.StringOkToTF(apiObject.GetNameOk())
	p.Description = framework.StringOkToTF(apiObject.GetDescriptionOk())

	return diags
}
