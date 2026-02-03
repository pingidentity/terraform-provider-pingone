// Copyright © 2025 Ping Identity Corporation

package sso

import (
	"context"
	"fmt"
	"net/http"
	"strings"

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
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/legacysdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
)

// Types
type SchemaAttributeDataSource serviceClientType

type SchemaAttributeDataSourceModel struct {
	SchemaAttributeResourceModelV1
	AttributeId pingonetypes.ResourceIDValue `tfsdk:"attribute_id"`
}

// Framework interfaces
var (
	_ datasource.DataSource = &SchemaAttributeDataSource{}
)

// New Object
func NewSchemaAttributeDataSource() datasource.DataSource {
	return &SchemaAttributeDataSource{}
}

// Metadata
func (r *SchemaAttributeDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_schema_attribute"
}

// Schema
func (r *SchemaAttributeDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	nameLength := 1

	attributeIdDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The ID of the schema attribute.",
	).ExactlyOneOf([]string{"attribute_id", "name"}).AppendMarkdownString("Must be a valid PingOne resource ID.")

	environmentIdDescription := framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment.").AppendMarkdownString("Must be a valid PingOne resource ID.")

	nameDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The system name of the schema attribute.",
	).ExactlyOneOf([]string{"name", "attribute_id"})

	schemaIdDescription := framework.SchemaAttributeDescriptionFromMarkdown("The ID of the schema the schema attribute belongs to.").AppendMarkdownString("Must be a valid PingOne resource ID.")

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Data source to read PingOne schema attribute data.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": schema.StringAttribute{
				Description:         environmentIdDescription.Description,
				MarkdownDescription: environmentIdDescription.MarkdownDescription,
				Required:            true,
				CustomType:          pingonetypes.ResourceIDType{},
			},

			"schema_id": schema.StringAttribute{
				Description:         schemaIdDescription.Description,
				MarkdownDescription: schemaIdDescription.MarkdownDescription,
				Required:            true,
				CustomType:          pingonetypes.ResourceIDType{},
			},

			"attribute_id": schema.StringAttribute{
				Description:         attributeIdDescription.Description,
				MarkdownDescription: attributeIdDescription.MarkdownDescription,
				Optional:            true,

				CustomType: pingonetypes.ResourceIDType{},

				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("name")),
				},
			},

			"name": schema.StringAttribute{
				Description:         nameDescription.Description,
				MarkdownDescription: nameDescription.MarkdownDescription,
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("attribute_id")),
					stringvalidator.LengthAtLeast(nameLength),
				},
			},

			"display_name": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The display name of the attribute such as 'T-shirt size'.").Description,
				Computed:    true,
			},

			"description": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A description of the attribute.").Description,
				Computed:    true,
			},

			"enabled": schema.BoolAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("Indicates whether or not the attribute is enabled.").Description,
				Computed:    true,
			},

			"type": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The type of the attribute.").Description,
				Computed:    true,
			},

			"unique": schema.BoolAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("Indicates whether or not the attribute must have a unique value within the PingOne environment.").Description,
				Computed:    true,
			},

			"multivalued": schema.BoolAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("Indicates whether the attribute has multiple values or a single one. Maximum number of values stored is 1,000.").Description,
				Computed:    true,
			},

			"enumerated_values": schema.SetNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A set of one or more enumerated values for the attribute.").Description,
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"value": schema.StringAttribute{
							Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the value of the enumerated value item.").Description,
							Computed:    true,
						},
						"archived": schema.BoolAttribute{
							Description: framework.SchemaAttributeDescriptionFromMarkdown("A boolean that specifies whether the enumerated value is archived. Archived values cannot be added to a user, but existing archived values are preserved. This allows clients that read the schema to know all possible values of an attribute.").Description,
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the description of the enumerated value.").Description,
							Computed:    true,
						},
					},
				},
			},

			"regex_validation": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object representation of the optional regular expression representation of this attribute.").Description,
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"pattern": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the regular expression to which the attribute must conform.").Description,
						Computed:    true,
					},
					"requirements": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies a developer friendly description of the regular expression requirements.").Description,
						Computed:    true,
					},
					"values_pattern_should_match": schema.SetAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A set of one or more strings matching the regular expression.").Description,
						Computed:    true,
						ElementType: types.StringType,
					},
					"values_pattern_should_not_match": schema.SetAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A set of one or more strings that do not match the regular expression.").Description,
						Computed:    true,
						ElementType: types.StringType,
					},
				},
			},

			"required": schema.BoolAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("Indicates whether or not the attribute is required.").Description,
				Computed:    true,
			},

			"ldap_attribute": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The unique identifier for the LDAP attribute.").Description,
				Computed:    true,
			},

			"schema_type": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The schema type of the attribute.").Description,
				Computed:    true,
			},
		},
	}
}

func (r *SchemaAttributeDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *SchemaAttributeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *SchemaAttributeDataSourceModel

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

	var schemaAttribute *management.SchemaAttribute

	if !data.Name.IsNull() {

		var d diag.Diagnostics
		schemaAttribute, d = fetchSchemaAttributeFromName(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), data.SchemaId.ValueString(), data.Name.ValueString())
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}

	} else if !data.AttributeId.IsNull() {

		// Run the API call
		resp.Diagnostics.Append(legacysdk.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				fO, fR, fErr := r.Client.ManagementAPIClient.SchemasApi.ReadOneAttribute(ctx, data.EnvironmentId.ValueString(), data.SchemaId.ValueString(), data.AttributeId.ValueString()).Execute()
				return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
			},
			"ReadOneAttribute",
			legacysdk.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
			&schemaAttribute,
		)...)
		if resp.Diagnostics.HasError() {
			return
		}

	} else {
		resp.Diagnostics.AddError(
			"Missing attribute_id or name",
			"One of attribute_id or name must be provided.",
		)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(schemaAttribute)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (p *SchemaAttributeDataSourceModel) toState(apiObject *management.SchemaAttribute) diag.Diagnostics {
	diags := p.SchemaAttributeResourceModelV1.toState(apiObject)

	if diags.HasError() {
		return diags
	}

	p.AttributeId = p.Id

	return diags
}

func fetchSchemaAttributeFromName(ctx context.Context, apiClient *management.APIClient, environmentID, schemaID, name string) (*management.SchemaAttribute, diag.Diagnostics) {
	var diags diag.Diagnostics

	var schemaAttribute *management.SchemaAttribute
	diags.Append(legacysdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			pagedIterator := apiClient.SchemasApi.ReadAllSchemaAttributes(ctx, environmentID, schemaID).Execute()

			var initialHttpResponse *http.Response

			for pageCursor, err := range pagedIterator {
				if err != nil {
					return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, apiClient, environmentID, nil, pageCursor.HTTPResponse, err)
				}

				if initialHttpResponse == nil {
					initialHttpResponse = pageCursor.HTTPResponse
				}

				if schemaAttributes, ok := pageCursor.EntityArray.Embedded.GetAttributesOk(); ok {

					for _, schemaAttribute := range schemaAttributes {

						if schemaAttribute.SchemaAttribute != nil && strings.EqualFold(schemaAttribute.SchemaAttribute.GetName(), name) {
							return schemaAttribute.SchemaAttribute, pageCursor.HTTPResponse, nil
						}
					}

				}
			}

			return nil, initialHttpResponse, nil
		},
		"ReadAllSchemaAttributes",
		legacysdk.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
		&schemaAttribute,
	)...)
	if diags.HasError() {
		return nil, diags
	}

	if schemaAttribute == nil {
		diags.AddError(
			"Cannot find schema attribute from name",
			fmt.Sprintf("The schema attribute %s for schema %s in environment %s cannot be found", name, schemaID, environmentID),
		)

		return nil, diags
	}

	return schemaAttribute, diags
}
