package authorize

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/patrickcping/pingone-go-sdk-v2/authorize"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
)

func managedEntityObjectSchemaAttributes() (attributes map[string]schema.Attribute) {

	attributes = map[string]schema.Attribute{
		"owner": schema.SingleNestedAttribute{
			Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
			Computed:    true,

			Attributes: map[string]schema.Attribute{
				"service": schema.SingleNestedAttribute{
					Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
					Computed:    true,

					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
							Computed:    true,
						},
					},
				},
			},
		},

		"reference": schema.SingleNestedAttribute{
			Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
			Computed:    true,

			Attributes: map[string]schema.Attribute{
				"id": schema.StringAttribute{
					Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
					Computed:    true,
				},

				"type": schema.StringAttribute{
					Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
					Computed:    true,
				},

				"name": schema.StringAttribute{
					Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
					Computed:    true,
				},

				"ui_deep_link": schema.StringAttribute{
					Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
					Computed:    true,
				},
			},
		},

		"restrictions": schema.SingleNestedAttribute{
			Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
			Computed:    true,

			Attributes: map[string]schema.Attribute{
				"read_only": schema.BoolAttribute{
					Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
					Computed:    true,
				},

				"disallow_children": schema.BoolAttribute{
					Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
					Computed:    true,
				},
			},
		},
	}

	return attributes
}

func dataSourceManagedEntityObjectSchemaAttributes() (attributes map[string]dsschema.Attribute) {

	attributes = map[string]dsschema.Attribute{
		"owner": schema.SingleNestedAttribute{
			Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
			Computed:    true,

			Attributes: map[string]schema.Attribute{
				"service": schema.SingleNestedAttribute{
					Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
					Computed:    true,

					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
							Computed:    true,
						},
					},
				},
			},
		},

		"reference": schema.SingleNestedAttribute{
			Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
			Computed:    true,

			Attributes: map[string]schema.Attribute{
				"id": schema.StringAttribute{
					Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
					Computed:    true,
				},

				"type": schema.StringAttribute{
					Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
					Computed:    true,
				},

				"name": schema.StringAttribute{
					Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
					Computed:    true,
				},

				"ui_deep_link": schema.StringAttribute{
					Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
					Computed:    true,
				},
			},
		},

		"restrictions": schema.SingleNestedAttribute{
			Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
			Computed:    true,

			Attributes: map[string]schema.Attribute{
				"read_only": schema.BoolAttribute{
					Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
					Computed:    true,
				},

				"disallow_children": schema.BoolAttribute{
					Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
					Computed:    true,
				},
			},
		},
	}

	return attributes
}

type editorManagedEntityResourceModel struct {
	Owner        types.Object `tfsdk:"owner"`
	Reference    types.Object `tfsdk:"reference"`
	Restrictions types.Object `tfsdk:"restrictions"`
}

type editorManagedEntityOwnerResourceModel struct {
	Service types.Object `tfsdk:"service"`
}

type editorManagedEntityOwnerServiceResourceModel struct {
	Name types.String `tfsdk:"name"`
}

type editorManagedEntityReferenceResourceModel struct {
	Id         types.String `tfsdk:"id"`
	Type       types.String `tfsdk:"type"`
	Name       types.String `tfsdk:"name"`
	UiDeepLink types.String `tfsdk:"ui_deep_link"`
}

type editorManagedEntityRestrictionsResourceModel struct {
	ReadOnly         types.Bool `tfsdk:"read_only"`
	DisallowChildren types.Bool `tfsdk:"disallow_children"`
}

var (
	editorManagedEntityTFObjectTypes = map[string]attr.Type{
		"owner":        types.ObjectType{AttrTypes: editorManagedEntityOwnerTFObjectTypes},
		"reference":    types.ObjectType{AttrTypes: editorManagedEntityReferenceTFObjectTypes},
		"restrictions": types.ObjectType{AttrTypes: editorManagedEntityRestrictionsTFObjectTypes},
	}

	editorManagedEntityOwnerTFObjectTypes = map[string]attr.Type{
		"service": types.ObjectType{AttrTypes: editorManagedEntityOwnerServiceTFObjectTypes},
	}

	editorManagedEntityOwnerServiceTFObjectTypes = map[string]attr.Type{
		"name": types.StringType,
	}

	editorManagedEntityReferenceTFObjectTypes = map[string]attr.Type{
		"id":           types.StringType,
		"type":         types.StringType,
		"name":         types.StringType,
		"ui_deep_link": types.StringType,
	}

	editorManagedEntityRestrictionsTFObjectTypes = map[string]attr.Type{
		"read_only":         types.BoolType,
		"disallow_children": types.BoolType,
	}
)

func expandEditorManagedEntity(ctx context.Context, managedEntity basetypes.ObjectValue) (managedEntityObject *authorize.AuthorizeEditorDataManagedEntityDTO, diags diag.Diagnostics) {
	var plan *editorManagedEntityResourceModel
	diags.Append(managedEntity.As(ctx, &plan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)
	if diags.HasError() {
		return nil, diags
	}

	managedEntityObject, d := plan.expand(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	return
}

func (p *editorManagedEntityResourceModel) expand(ctx context.Context) (*authorize.AuthorizeEditorDataManagedEntityDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	var ownerPlan *editorManagedEntityOwnerResourceModel
	diags.Append(p.Owner.As(ctx, &ownerPlan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)
	if diags.HasError() {
		return nil, diags
	}

	owner, d := ownerPlan.expand(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	data := authorize.NewAuthorizeEditorDataManagedEntityDTO(*owner)

	if !p.Reference.IsNull() && !p.Reference.IsUnknown() {
		var plan *editorManagedEntityReferenceResourceModel
		diags.Append(p.Reference.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		reference := plan.expand()

		data.SetReference(*reference)
	}

	if !p.Restrictions.IsNull() && !p.Restrictions.IsUnknown() {
		var plan *editorManagedEntityRestrictionsResourceModel
		diags.Append(p.Restrictions.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		restrictions := plan.expand()

		data.SetRestrictions(*restrictions)
	}

	return data, diags
}

func (p *editorManagedEntityOwnerResourceModel) expand(ctx context.Context) (*authorize.AuthorizeEditorDataManagedEntityOwnerDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	var servicePlan *editorManagedEntityOwnerServiceResourceModel
	diags.Append(p.Service.As(ctx, &servicePlan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)
	if diags.HasError() {
		return nil, diags
	}

	service := servicePlan.expand()

	data := authorize.NewAuthorizeEditorDataManagedEntityOwnerDTO(*service)

	return data, diags
}

func (p *editorManagedEntityOwnerServiceResourceModel) expand() *authorize.AuthorizeEditorDataServiceObjectDTO {

	data := authorize.NewAuthorizeEditorDataServiceObjectDTO(
		p.Name.ValueString(),
	)

	return data
}

func (p *editorManagedEntityReferenceResourceModel) expand() *authorize.AuthorizeEditorDataManagedEntityManagedEntityReferenceDTO {

	data := authorize.NewAuthorizeEditorDataManagedEntityManagedEntityReferenceDTO()

	if !p.Id.IsNull() && !p.Id.IsUnknown() {
		data.SetId(p.Id.ValueString())
	}

	if !p.Type.IsNull() && !p.Type.IsUnknown() {
		data.SetType(p.Type.ValueString())
	}

	if !p.Name.IsNull() && !p.Name.IsUnknown() {
		data.SetName(p.Name.ValueString())
	}

	if !p.UiDeepLink.IsNull() && !p.UiDeepLink.IsUnknown() {
		data.SetUiDeepLink(p.UiDeepLink.ValueString())
	}

	return data
}

func (p *editorManagedEntityRestrictionsResourceModel) expand() *authorize.AuthorizeEditorDataManagedEntityRestrictionsDTO {

	data := authorize.NewAuthorizeEditorDataManagedEntityRestrictionsDTO()

	if !p.ReadOnly.IsNull() && !p.ReadOnly.IsUnknown() {
		data.SetReadOnly(p.ReadOnly.ValueBool())
	}

	if !p.DisallowChildren.IsNull() && !p.DisallowChildren.IsUnknown() {
		data.SetDisallowChildren(p.DisallowChildren.ValueBool())
	}

	return data
}

func editorManagedEntityOkToTF(apiObject *authorize.AuthorizeEditorDataManagedEntityDTO, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags, d diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(editorManagedEntityTFObjectTypes), diags
	}

	owner, d := editorManagedEntityOwnerOkToTF(apiObject.GetOwnerOk())
	diags.Append(d...)

	reference, d := editorManagedEntityReferenceOkToTF(apiObject.GetReferenceOk())
	diags.Append(d...)

	restrictions, d := editorManagedEntityRestrictionsOkToTF(apiObject.GetRestrictionsOk())
	diags.Append(d...)

	objValue, d := types.ObjectValue(editorManagedEntityTFObjectTypes, map[string]attr.Value{
		"owner":        owner,
		"reference":    reference,
		"restrictions": restrictions,
	})
	diags.Append(d...)

	return objValue, diags
}

func editorManagedEntityOwnerOkToTF(apiObject *authorize.AuthorizeEditorDataManagedEntityOwnerDTO, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(editorManagedEntityOwnerTFObjectTypes), diags
	}

	service, d := editorManagedEntityOwnerServiceOkToTF(apiObject.GetServiceOk())
	diags.Append(d...)

	objValue, d := types.ObjectValue(editorManagedEntityOwnerTFObjectTypes, map[string]attr.Value{
		"service": service,
	})
	diags.Append(d...)

	return objValue, diags
}

func editorManagedEntityOwnerServiceOkToTF(apiObject *authorize.AuthorizeEditorDataServiceObjectDTO, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(editorManagedEntityOwnerServiceTFObjectTypes), diags
	}

	objValue, d := types.ObjectValue(editorManagedEntityOwnerServiceTFObjectTypes, map[string]attr.Value{
		"name": framework.StringOkToTF(apiObject.GetNameOk()),
	})
	diags.Append(d...)

	return objValue, diags
}

func editorManagedEntityReferenceOkToTF(apiObject *authorize.AuthorizeEditorDataManagedEntityManagedEntityReferenceDTO, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(editorManagedEntityReferenceTFObjectTypes), diags
	}

	objValue, d := types.ObjectValue(editorManagedEntityReferenceTFObjectTypes, map[string]attr.Value{
		"id":           framework.StringOkToTF(apiObject.GetIdOk()),
		"type":         framework.StringOkToTF(apiObject.GetTypeOk()),
		"name":         framework.StringOkToTF(apiObject.GetNameOk()),
		"ui_deep_link": framework.StringOkToTF(apiObject.GetUiDeepLinkOk()),
	})
	diags.Append(d...)

	return objValue, diags
}

func editorManagedEntityRestrictionsOkToTF(apiObject *authorize.AuthorizeEditorDataManagedEntityRestrictionsDTO, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(editorManagedEntityRestrictionsTFObjectTypes), diags
	}

	objValue, d := types.ObjectValue(editorManagedEntityRestrictionsTFObjectTypes, map[string]attr.Value{
		"read_only":         framework.BoolOkToTF(apiObject.GetReadOnlyOk()),
		"disallow_children": framework.BoolOkToTF(apiObject.GetDisallowChildrenOk()),
	})
	diags.Append(d...)

	return objValue, diags
}
