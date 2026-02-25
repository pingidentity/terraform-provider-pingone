// Copyright © 2026 Ping Identity Corporation

package sso

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
)

type SchemaAttributeResourceModelV0 struct {
	Id               pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId    pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	Description      types.String                 `tfsdk:"description"`
	DisplayName      types.String                 `tfsdk:"display_name"`
	Enabled          types.Bool                   `tfsdk:"enabled"`
	EnumeratedValues types.Set                    `tfsdk:"enumerated_values"`
	LdapAttribute    types.String                 `tfsdk:"ldap_attribute"`
	Multivalued      types.Bool                   `tfsdk:"multivalued"`
	Name             types.String                 `tfsdk:"name"`
	RegexValidation  types.Object                 `tfsdk:"regex_validation"`
	Required         types.Bool                   `tfsdk:"required"`
	SchemaId         pingonetypes.ResourceIDValue `tfsdk:"schema_id"`
	SchemaType       types.String                 `tfsdk:"schema_type"`
	Type             types.String                 `tfsdk:"type"`
	Unique           types.Bool                   `tfsdk:"unique"`
}

func (r *SchemaAttributeResource) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{
		// State upgrade implementation from 0 (prior state version) to 1 (Schema.Version)
		0: {
			PriorSchema: &schema.Schema{
				Attributes: map[string]schema.Attribute{
					"id": framework.Attr_ID(),

					"environment_id": framework.Attr_LinkID(
						framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment to create the schema attribute in."),
					),

					"schema_id": schema.StringAttribute{
						Required: true,
					},

					"name": schema.StringAttribute{
						Required: true,
					},

					"display_name": schema.StringAttribute{
						Optional: true,
					},

					"description": schema.StringAttribute{
						Optional: true,
					},

					"enabled": schema.BoolAttribute{
						Optional: true,
						Computed: true,

						Default: booldefault.StaticBool(true),
					},

					"type": schema.StringAttribute{
						Optional: true,
						Computed: true,

						Default: stringdefault.StaticString(string(management.ENUMSCHEMAATTRIBUTETYPE_STRING)),
					},

					"unique": schema.BoolAttribute{
						Optional: true,
						Computed: true,

						Default: booldefault.StaticBool(false),
					},

					"multivalued": schema.BoolAttribute{
						Optional: true,
						Computed: true,

						Default: booldefault.StaticBool(false),
					},

					"enumerated_values": schema.SetNestedAttribute{
						Optional: true,

						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"value": schema.StringAttribute{
									Required: true,
								},

								"archived": schema.BoolAttribute{
									Optional: true,
									Computed: true,

									Default: booldefault.StaticBool(false),
								},

								"description": schema.StringAttribute{
									Optional: true,
								},
							},
						},
					},

					"regex_validation": schema.SingleNestedAttribute{
						Optional: true,

						Attributes: map[string]schema.Attribute{
							"pattern": schema.StringAttribute{
								Required: true,
							},

							"requirements": schema.StringAttribute{
								Required: true,
							},

							"values_pattern_should_match": schema.SetAttribute{
								Optional: true,

								ElementType: types.StringType,
							},

							"values_pattern_should_not_match": schema.SetAttribute{
								Optional: true,

								ElementType: types.StringType,
							},
						},
					},

					"required": schema.BoolAttribute{
						Computed: true,
					},

					"ldap_attribute": schema.StringAttribute{
						Computed: true,
					},

					"schema_type": schema.StringAttribute{
						Computed: true,
					},
				},
			},
			StateUpgrader: func(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
				var priorStateData SchemaAttributeResourceModelV0

				resp.Diagnostics.Append(req.State.Get(ctx, &priorStateData)...)

				if resp.Diagnostics.HasError() {
					return
				}

				upgradedStateData := priorStateData

				resp.Diagnostics.Append(resp.State.Set(ctx, upgradedStateData)...)
			},
		},
	}
}
