// Copyright © 2025 Ping Identity Corporation

package sso

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/filter"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
)

// Types
type GroupDataSource serviceClientType

type GroupDataSourceModel struct {
	Id            pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	GroupId       pingonetypes.ResourceIDValue `tfsdk:"group_id"`
	Name          types.String                 `tfsdk:"name"`
	Description   types.String                 `tfsdk:"description"`
	PopulationId  pingonetypes.ResourceIDValue `tfsdk:"population_id"`
	UserFilter    types.String                 `tfsdk:"user_filter"`
	ExternalId    types.String                 `tfsdk:"external_id"`
	CustomData    jsontypes.Normalized         `tfsdk:"custom_data"`
}

// Framework interfaces
var (
	_ datasource.DataSource = &GroupDataSource{}
)

// New Object
func NewGroupDataSource() datasource.DataSource {
	return &GroupDataSource{}
}

// Metadata
func (r *GroupDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group"
}

// Schema
func (r *GroupDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	nameLength := 1

	groupIdDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the ID of the group to retrieve configuration for.  Must be a valid PingOne resource ID.",
	).ExactlyOneOf([]string{"group_id", "name"})

	groupNameDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the name of the group to retrieve configuration for.",
	).ExactlyOneOf([]string{"group_id", "name"})

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Datasource to retrieve a PingOne group in an environment by ID or by name.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment that is configured with the group.  Must be a valid PingOne resource ID.").Description,
				Required:    true,

				CustomType: pingonetypes.ResourceIDType{},
			},

			"group_id": schema.StringAttribute{
				Description:         groupIdDescription.Description,
				MarkdownDescription: groupIdDescription.MarkdownDescription,
				Optional:            true,

				CustomType: pingonetypes.ResourceIDType{},

				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("name")),
				},
			},

			"name": schema.StringAttribute{
				Description:         groupNameDescription.Description,
				MarkdownDescription: groupNameDescription.MarkdownDescription,
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("group_id")),
					stringvalidator.LengthAtLeast(nameLength),
				},
			},

			"description": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the description applied to the group.").Description,
				Computed:    true,
			},

			"population_id": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the ID of the population that the group is assigned to.").Description,
				Computed:    true,

				CustomType: pingonetypes.ResourceIDType{},
			},

			"user_filter": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the SCIM filter applied to dynamically assign users to the group.").Description,
				Computed:    true,
			},

			"external_id": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies a user defined ID that represents the counterpart group in an external system.").Description,
				Computed:    true,
			},

			"custom_data": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A JSON string that specifies user-defined custom data.").Description,
				Computed:    true,

				CustomType: jsontypes.NormalizedType{},
			},
		},
	}
}

func (r *GroupDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *GroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *GroupDataSourceModel

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

	var group *management.Group
	var scimFilter string

	if !data.Name.IsNull() {

		scimFilter = filter.BuildScimFilter(
			append(make([]interface{}, 0), map[string]interface{}{
				"name":   "name",
				"values": []string{data.Name.ValueString()},
			}), map[string]string{})

	} else if !data.GroupId.IsNull() {

		scimFilter = filter.BuildScimFilter(
			append(make([]interface{}, 0), map[string]interface{}{
				"name":   "id",
				"values": []string{data.GroupId.ValueString()},
			}), map[string]string{})

	} else {
		resp.Diagnostics.AddError(
			"Missing parameter",
			"Cannot find the requested group. group_id or name must be set.",
		)
		return
	}

	// Run the API call
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			pagedIterator := r.Client.ManagementAPIClient.GroupsApi.ReadAllGroups(ctx, data.EnvironmentId.ValueString()).Filter(scimFilter).Execute()

			var initialHttpResponse *http.Response

			for pageCursor, err := range pagedIterator {
				if err != nil {
					return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), nil, pageCursor.HTTPResponse, err)
				}

				if initialHttpResponse == nil {
					initialHttpResponse = pageCursor.HTTPResponse
				}

				if groups, ok := pageCursor.EntityArray.Embedded.GetGroupsOk(); ok {
					for _, g := range groups {

						if !data.Name.IsNull() && g.GetName() == data.Name.ValueString() {
							return &g, pageCursor.HTTPResponse, nil
						}

						if !data.GroupId.IsNull() && g.GetId() == data.GroupId.ValueString() {
							return &g, pageCursor.HTTPResponse, nil
						}
					}
				}
			}

			return nil, initialHttpResponse, nil
		},
		"ReadAllGroups",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
		&group,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if group == nil {
		resp.Diagnostics.AddError(
			"Group not found",
			fmt.Sprintf("The group with the specified group_id or name cannot be found in environment %s.", data.EnvironmentId.String()),
		)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(group)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (p *GroupDataSourceModel) toState(apiObject *management.Group) diag.Diagnostics {
	var diags, d diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.Id = framework.PingOneResourceIDToTF(apiObject.GetId())
	p.GroupId = framework.PingOneResourceIDToTF(apiObject.GetId())
	p.Name = framework.StringOkToTF(apiObject.GetNameOk())
	p.Description = framework.StringOkToTF(apiObject.GetDescriptionOk())

	if v, ok := apiObject.GetPopulationOk(); ok && v != nil {
		p.PopulationId = framework.PingOneResourceIDOkToTF(v.GetIdOk())
	} else {
		p.PopulationId = pingonetypes.NewResourceIDNull()
	}

	p.UserFilter = framework.StringOkToTF(apiObject.GetUserFilterOk())
	p.ExternalId = framework.StringOkToTF(apiObject.GetExternalIdOk())
	p.CustomData, d = framework.JSONNormalizedOkToTF(apiObject.GetCustomDataOk())
	diags.Append(d...)

	return diags
}
