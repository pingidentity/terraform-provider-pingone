package authorize

import (
	"context"
	"fmt"
	"slices"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/patrickcping/pingone-go-sdk-v2/authorize"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	listvalidatorinternal "github.com/pingidentity/terraform-provider-pingone/internal/framework/listvalidator"
	objectvalidatorinternal "github.com/pingidentity/terraform-provider-pingone/internal/framework/objectvalidator"
	stringvalidatorinternal "github.com/pingidentity/terraform-provider-pingone/internal/framework/stringvalidator"
	"github.com/pingidentity/terraform-provider-pingone/internal/utils"
)

const processorNestedIterationMaxDepth = 2

var leafProcessorTypes = []authorize.EnumAuthorizeEditorDataProcessorDTOType{
	"JSON_PATH",
	"REFERENCE",
	"SPEL",
	"XPATH",
}

func dataProcessorObjectSchemaAttributes() (attributes map[string]schema.Attribute) {
	const initialIteration = 1
	return dataProcessorObjectSchemaAttributesIteration(initialIteration)
}
func dataProcessorObjectSchemaAttributesIteration(iteration int32) (attributes map[string]schema.Attribute) {

	supportedTypes := authorize.AllowedEnumAuthorizeEditorDataProcessorDTOTypeEnumValues

	if iteration >= processorNestedIterationMaxDepth {
		supportedTypes = leafProcessorTypes
	}

	typeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the processor type.",
	).AllowedValuesEnum(supportedTypes)

	processorsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The list of processors to apply in the given order.",
	).AppendMarkdownString(fmt.Sprintf("This field is required when `type` is `%s`.", string(authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_CHAIN)))

	predicateDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"",
	).AppendMarkdownString(fmt.Sprintf("This field is required when `type` is `%s`.", string(authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_COLLECTION_FILTER)))

	processorDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"",
	).AppendMarkdownString(fmt.Sprintf("This field is required when `type` is `%s`.", string(authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_COLLECTION_TRANSFORM)))

	expressionDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("A string that specifies the expression to use. If the `type` is `%s`, the expression should be a valid JSON path expression, if the `type` is `%s`, the expression should be a valid SpEL expression and if the `type` is `%s`, the expression should be a valid XPath expression.", string(authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_JSON_PATH), string(authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_SPEL), string(authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_XPATH)),
	).AppendMarkdownString(fmt.Sprintf("This field is required when `type` is `%s`, `%s` or `%s`.", string(authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_JSON_PATH), string(authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_SPEL), string(authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_XPATH)))

	valueTypeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An object that specifies the output type of the value.",
	).AppendMarkdownString(fmt.Sprintf("This field is required when `type` is `%s`, `%s` or `%s`.", string(authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_JSON_PATH), string(authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_SPEL), string(authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_XPATH)))

	valueTypeValueCollectionFilterDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("Must be `%s` when the processor type (`processor.type`) is `%s`.", string(authorize.ENUMAUTHORIZEEDITORDATAVALUETYPEDTO_BOOLEAN), string(authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_COLLECTION_FILTER)),
	)

	processorRefDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An object that specifies configuration settings for the authorization processor to reference.",
	).AppendMarkdownString(fmt.Sprintf("This field is required when `type` is `%s`.", string(authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_REFERENCE)))

	attributes = map[string]schema.Attribute{
		"name": schema.StringAttribute{
			Description: framework.SchemaAttributeDescriptionFromMarkdown("A user-friendly authorization processor name. The value must be unique.").Description,
			Required:    true,
		},

		"type": schema.StringAttribute{
			Description:         typeDescription.Description,
			MarkdownDescription: typeDescription.MarkdownDescription,
			Required:            true,

			Validators: []validator.String{
				stringvalidator.OneOf(utils.EnumSliceToStringSlice(supportedTypes)...),
			},
		},
	}

	// type == "CHAIN"
	if slices.Contains(supportedTypes, authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_CHAIN) {
		attributes["processors"] = schema.ListNestedAttribute{
			Description:         processorsDescription.Description,
			MarkdownDescription: processorsDescription.MarkdownDescription,
			Optional:            true,

			Validators: []validator.List{
				listvalidatorinternal.IsRequiredIfMatchesPathValue(
					types.StringValue(string(authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_CHAIN)),
					path.MatchRelative().AtParent().AtName("type"),
				),
			},

			NestedObject: schema.NestedAttributeObject{
				Attributes: dataProcessorObjectSchemaAttributesIteration(iteration + 1),
			},
		}
	}

	// type == "COLLECTION_FILTER"
	if slices.Contains(supportedTypes, authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_COLLECTION_FILTER) {
		attributes["predicate"] = schema.SingleNestedAttribute{
			Description:         predicateDescription.Description,
			MarkdownDescription: predicateDescription.MarkdownDescription,
			Optional:            true,

			Validators: []validator.Object{
				objectvalidatorinternal.IsRequiredIfMatchesPathValue(
					types.StringValue(string(authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_COLLECTION_FILTER)),
					path.MatchRelative().AtParent().AtName("type"),
				),
			},

			Attributes: dataProcessorObjectSchemaAttributesIteration(iteration + 1),
		}
	}

	// type == "COLLECTION_TRANSFORM"
	if slices.Contains(supportedTypes, authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_COLLECTION_TRANSFORM) {
		attributes["processor"] = schema.SingleNestedAttribute{
			Description:         processorDescription.Description,
			MarkdownDescription: processorDescription.MarkdownDescription,
			Optional:            true,

			Validators: []validator.Object{
				objectvalidatorinternal.IsRequiredIfMatchesPathValue(
					types.StringValue(string(authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_COLLECTION_TRANSFORM)),
					path.MatchRelative().AtParent().AtName("type"),
				),
			},

			Attributes: dataProcessorObjectSchemaAttributesIteration(iteration + 1),
		}
	}

	// type == "JSON_PATH", type == "SPEL", type == "XPATH"
	if slices.Contains(supportedTypes, authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_JSON_PATH) ||
		slices.Contains(supportedTypes, authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_SPEL) ||
		slices.Contains(supportedTypes, authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_XPATH) {
		attributes["expression"] = schema.StringAttribute{
			Description:         expressionDescription.Description,
			MarkdownDescription: expressionDescription.MarkdownDescription,
			Optional:            true,

			Validators: []validator.String{
				stringvalidatorinternal.IsRequiredIfMatchesPathValue(
					types.StringValue(string(authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_JSON_PATH)),
					path.MatchRelative().AtParent().AtName("type"),
				),
				stringvalidatorinternal.IsRequiredIfMatchesPathValue(
					types.StringValue(string(authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_SPEL)),
					path.MatchRelative().AtParent().AtName("type"),
				),
				stringvalidatorinternal.IsRequiredIfMatchesPathValue(
					types.StringValue(string(authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_XPATH)),
					path.MatchRelative().AtParent().AtName("type"),
				),
			},
		}

		attributes["value_type"] = schema.SingleNestedAttribute{
			Description:         valueTypeDescription.Description,
			MarkdownDescription: valueTypeDescription.MarkdownDescription,
			Optional:            true,

			Validators: []validator.Object{
				objectvalidatorinternal.IsRequiredIfMatchesPathValue(
					types.StringValue(string(authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_JSON_PATH)),
					path.MatchRelative().AtParent().AtName("type"),
				),
				objectvalidatorinternal.IsRequiredIfMatchesPathValue(
					types.StringValue(string(authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_SPEL)),
					path.MatchRelative().AtParent().AtName("type"),
				),
				objectvalidatorinternal.IsRequiredIfMatchesPathValue(
					types.StringValue(string(authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_XPATH)),
					path.MatchRelative().AtParent().AtName("type"),
				),
			},

			Attributes: valueTypeObjectSchemaAttributes(
				// If processor type is collection filter, then the value type value must be `BOOLEAN`
				stringvalidatorinternal.CustomStringValidatorModel{
					Description: valueTypeValueCollectionFilterDescription,
					Validators: []validator.String{
						stringvalidatorinternal.ShouldBeDefinedValueIfPathMatchesValue(
							types.StringValue(string(authorize.ENUMAUTHORIZEEDITORDATAVALUETYPEDTO_BOOLEAN)),
							types.StringValue(string(authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_COLLECTION_FILTER)),
							path.MatchRoot("processor").AtName("type"),
						),
					},
				},
			),
		}
	}

	// type == "REFERENCE"
	if slices.Contains(supportedTypes, authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_REFERENCE) {
		attributes["processor_ref"] = schema.SingleNestedAttribute{
			Description:         processorRefDescription.Description,
			MarkdownDescription: processorRefDescription.MarkdownDescription,
			Optional:            true,

			Validators: []validator.Object{
				objectvalidatorinternal.IsRequiredIfMatchesPathValue(
					types.StringValue(string(authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_REFERENCE)),
					path.MatchRelative().AtParent().AtName("type"),
				),
			},

			Attributes: referenceIdObjectSchemaAttributes(framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the ID of the authorization processor in the trust framework.")),
		}
	}

	return attributes
}

func dataSourceDataProcessorObjectSchemaAttributes() (attributes map[string]dsschema.Attribute) {
	const initialIteration = 1
	return dataSourceDataProcessorObjectSchemaAttributesIteration(initialIteration)
}
func dataSourceDataProcessorObjectSchemaAttributesIteration(iteration int32) (attributes map[string]dsschema.Attribute) {

	supportedTypes := authorize.AllowedEnumAuthorizeEditorDataProcessorDTOTypeEnumValues

	if iteration >= processorNestedIterationMaxDepth {
		supportedTypes = leafProcessorTypes
	}

	typeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the processor type.",
	).AllowedValuesEnum(supportedTypes)

	processorsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The list of processors to apply in the given order.",
	).AppendMarkdownString(fmt.Sprintf("This field is required when `type` is `%s`.", string(authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_CHAIN)))

	predicateDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"",
	).AppendMarkdownString(fmt.Sprintf("This field is required when `type` is `%s`.", string(authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_COLLECTION_FILTER)))

	processorDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"",
	).AppendMarkdownString(fmt.Sprintf("This field is required when `type` is `%s`.", string(authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_COLLECTION_TRANSFORM)))

	expressionDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("A string that specifies the expression to use. If the `type` is `%s`, the expression should be a valid JSON path expression, if the `type` is `%s`, the expression should be a valid SpEL expression and if the `type` is `%s`, the expression should be a valid XPath expression.", string(authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_JSON_PATH), string(authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_SPEL), string(authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_XPATH)),
	).AppendMarkdownString(fmt.Sprintf("This field is required when `type` is `%s`, `%s` or `%s`.", string(authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_JSON_PATH), string(authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_SPEL), string(authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_XPATH)))

	valueTypeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An object that specifies the output type of the value.",
	).AppendMarkdownString(fmt.Sprintf("This field is required when `type` is `%s`, `%s` or `%s`.", string(authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_JSON_PATH), string(authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_SPEL), string(authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_XPATH)))

	valueTypeValueCollectionFilterDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("Must be `%s` when the processor type (`processor.type`) is `%s`.", string(authorize.ENUMAUTHORIZEEDITORDATAVALUETYPEDTO_BOOLEAN), string(authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_COLLECTION_FILTER)),
	)

	processorRefDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An object that specifies configuration settings for the authorization processor to reference.",
	).AppendMarkdownString(fmt.Sprintf("This field is required when `type` is `%s`.", string(authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_REFERENCE)))

	attributes = map[string]dsschema.Attribute{
		"name": schema.StringAttribute{
			Description: framework.SchemaAttributeDescriptionFromMarkdown("A user-friendly authorization processor name. The value must be unique.").Description,
			Computed:    true,
		},

		"type": schema.StringAttribute{
			Description:         typeDescription.Description,
			MarkdownDescription: typeDescription.MarkdownDescription,
			Computed:            true,
		},
	}

	// type == "CHAIN"
	if slices.Contains(supportedTypes, authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_CHAIN) {
		attributes["processors"] = schema.ListNestedAttribute{
			Description:         processorsDescription.Description,
			MarkdownDescription: processorsDescription.MarkdownDescription,
			Computed:            true,

			NestedObject: schema.NestedAttributeObject{
				Attributes: dataProcessorObjectSchemaAttributesIteration(iteration + 1),
			},
		}
	}

	// type == "COLLECTION_FILTER"
	if slices.Contains(supportedTypes, authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_COLLECTION_FILTER) {
		attributes["predicate"] = schema.SingleNestedAttribute{
			Description:         predicateDescription.Description,
			MarkdownDescription: predicateDescription.MarkdownDescription,
			Computed:            true,

			Attributes: dataProcessorObjectSchemaAttributesIteration(iteration + 1),
		}
	}

	// type == "COLLECTION_TRANSFORM"
	if slices.Contains(supportedTypes, authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_COLLECTION_TRANSFORM) {
		attributes["processor"] = schema.SingleNestedAttribute{
			Description:         processorDescription.Description,
			MarkdownDescription: processorDescription.MarkdownDescription,
			Computed:            true,

			Attributes: dataProcessorObjectSchemaAttributesIteration(iteration + 1),
		}
	}

	// type == "JSON_PATH", type == "SPEL", type == "XPATH"
	if slices.Contains(supportedTypes, authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_JSON_PATH) ||
		slices.Contains(supportedTypes, authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_SPEL) ||
		slices.Contains(supportedTypes, authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_XPATH) {
		attributes["expression"] = schema.StringAttribute{
			Description:         expressionDescription.Description,
			MarkdownDescription: expressionDescription.MarkdownDescription,
			Computed:            true,
		}

		attributes["value_type"] = schema.SingleNestedAttribute{
			Description:         valueTypeDescription.Description,
			MarkdownDescription: valueTypeDescription.MarkdownDescription,
			Computed:            true,

			Attributes: valueTypeObjectSchemaAttributes(
				// If processor type is collection filter, then the value type value must be `BOOLEAN`
				stringvalidatorinternal.CustomStringValidatorModel{
					Description: valueTypeValueCollectionFilterDescription,
				},
			),
		}
	}

	// type == "REFERENCE"
	if slices.Contains(supportedTypes, authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_REFERENCE) {
		attributes["processor_ref"] = schema.SingleNestedAttribute{
			Description:         processorRefDescription.Description,
			MarkdownDescription: processorRefDescription.MarkdownDescription,
			Computed:            true,

			Attributes: referenceIdObjectSchemaAttributes(framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the ID of the authorization processor in the trust framework.")),
		}
	}

	return attributes
}

type editorDataProcessorLeafResourceModel struct {
	Name         types.String `tfsdk:"name"`
	Type         types.String `tfsdk:"type"`
	Expression   types.String `tfsdk:"expression"`
	ValueType    types.Object `tfsdk:"value_type"`
	ProcessorRef types.Object `tfsdk:"processor_ref"`
}

type editorDataProcessorResourceModel struct {
	Name         types.String `tfsdk:"name"`
	Type         types.String `tfsdk:"type"`
	Expression   types.String `tfsdk:"expression"`
	ValueType    types.Object `tfsdk:"value_type"`
	ProcessorRef types.Object `tfsdk:"processor_ref"`
	Processors   types.List   `tfsdk:"processors"`
	Predicate    types.Object `tfsdk:"predicate"`
	Processor    types.Object `tfsdk:"processor"`
}

var editorDataProcessorTFObjectTypes = initializeEditorDataProcessorTFObjectTypes(1)

func initializeEditorDataProcessorTFObjectTypes(iteration int32) map[string]attr.Type {

	supportedTypes := authorize.AllowedEnumAuthorizeEditorDataProcessorDTOTypeEnumValues

	if iteration >= processorNestedIterationMaxDepth {
		supportedTypes = leafProcessorTypes
	}

	attrMap := map[string]attr.Type{
		"name": types.StringType,
		"type": types.StringType,
	}

	if slices.Contains(supportedTypes, authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_CHAIN) {
		attrMap["processors"] = types.ListType{
			ElemType: types.ObjectType{AttrTypes: initializeEditorDataProcessorTFObjectTypes(iteration + 1)},
		}
	}

	if slices.Contains(supportedTypes, authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_COLLECTION_FILTER) {
		attrMap["predicate"] = types.ObjectType{AttrTypes: initializeEditorDataProcessorTFObjectTypes(iteration + 1)}
	}

	if slices.Contains(supportedTypes, authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_COLLECTION_TRANSFORM) {
		attrMap["processor"] = types.ObjectType{AttrTypes: initializeEditorDataProcessorTFObjectTypes(iteration + 1)}
	}

	if slices.Contains(supportedTypes, authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_JSON_PATH) ||
		slices.Contains(supportedTypes, authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_SPEL) ||
		slices.Contains(supportedTypes, authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_XPATH) {
		attrMap["expression"] = types.StringType
		attrMap["value_type"] = types.ObjectType{AttrTypes: editorValueTypeTFObjectTypes}
	}

	if slices.Contains(supportedTypes, authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_REFERENCE) {
		attrMap["processor_ref"] = types.ObjectType{AttrTypes: editorReferenceObjectTFObjectTypes}
	}

	return attrMap
}

func expandEditorDataProcessor(ctx context.Context, processor basetypes.ObjectValue) (processorObject *authorize.AuthorizeEditorDataProcessorDTO, diags diag.Diagnostics) {
	const initialIteration = 1
	return expandEditorDataProcessorIteration(ctx, processor, initialIteration)
}

func expandEditorDataProcessorIteration(ctx context.Context, processor basetypes.ObjectValue, iteration int32) (processorObject *authorize.AuthorizeEditorDataProcessorDTO, diags diag.Diagnostics) {
	var d diag.Diagnostics

	if processor.IsNull() {
		return nil, diags
	}

	leaf := iteration >= processorNestedIterationMaxDepth

	if leaf {
		var plan *editorDataProcessorLeafResourceModel
		diags.Append(processor.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return
		}

		processorObject, d = plan.expand(ctx)
		diags.Append(d...)
	} else {
		var plan *editorDataProcessorResourceModel
		diags.Append(processor.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return
		}

		processorObject, d = plan.expand(ctx, iteration)
		diags.Append(d...)
	}
	if diags.HasError() {
		return
	}

	return
}

func expandEditorDataProcessors(ctx context.Context, processor basetypes.ListValue) (processorObject []authorize.AuthorizeEditorDataProcessorDTO, diags diag.Diagnostics) {
	const initialIteration = 1
	return expandEditorDataProcessorsIteration(ctx, processor, initialIteration)
}

func expandEditorDataProcessorsIteration(ctx context.Context, processors basetypes.ListValue, iteration int32) (processorObjects []authorize.AuthorizeEditorDataProcessorDTO, diags diag.Diagnostics) {

	leaf := iteration >= processorNestedIterationMaxDepth

	if leaf {

		var plan []editorDataProcessorLeafResourceModel
		diags.Append(processors.ElementsAs(ctx, &plan, false)...)
		if diags.HasError() {
			return nil, diags
		}

		processorObjects = make([]authorize.AuthorizeEditorDataProcessorDTO, 0, len(plan))
		for _, processorPlan := range plan {

			processorObject, d := processorPlan.expand(ctx)
			diags.Append(d...)
			if diags.HasError() {
				return nil, diags
			}

			processorObjects = append(processorObjects, *processorObject)
		}

	} else {

		var plan []editorDataProcessorResourceModel
		diags.Append(processors.ElementsAs(ctx, &plan, false)...)
		if diags.HasError() {
			return nil, diags
		}

		processorObjects = make([]authorize.AuthorizeEditorDataProcessorDTO, 0, len(plan))
		for _, processorPlan := range plan {

			processorObject, d := processorPlan.expand(ctx, iteration)
			diags.Append(d...)
			if diags.HasError() {
				return nil, diags
			}

			processorObjects = append(processorObjects, *processorObject)
		}

	}

	return

}

func (p *editorDataProcessorResourceModel) expand(ctx context.Context, iteration int32) (*authorize.AuthorizeEditorDataProcessorDTO, diag.Diagnostics) {
	var diags, d diag.Diagnostics

	data := authorize.AuthorizeEditorDataProcessorDTO{}

	switch authorize.EnumAuthorizeEditorDataProcessorDTOType(p.Type.ValueString()) {
	case authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_CHAIN:
		data.AuthorizeEditorDataProcessorsChainProcessorDTO, d = p.expandChainProcessor(ctx, iteration)
		diags.Append(d...)
	case authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_COLLECTION_FILTER:
		data.AuthorizeEditorDataProcessorsCollectionFilterProcessorDTO, d = p.expandCollectionFilterProcessor(ctx, iteration)
		diags.Append(d...)
	case authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_COLLECTION_TRANSFORM:
		data.AuthorizeEditorDataProcessorsCollectionTransformProcessorDTO, d = p.expandCollectionTransformProcessor(ctx, iteration)
		diags.Append(d...)
	case authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_JSON_PATH:
		data.AuthorizeEditorDataProcessorsJsonPathProcessorDTO, d = p.expandJsonPathProcessor(ctx)
		diags.Append(d...)
	case authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_REFERENCE:
		data.AuthorizeEditorDataProcessorsReferenceProcessorDTO, d = p.expandReferenceProcessor(ctx)
		diags.Append(d...)
	case authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_SPEL:
		data.AuthorizeEditorDataProcessorsSpelProcessorDTO, d = p.expandSPELProcessor(ctx)
		diags.Append(d...)
	case authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_XPATH:
		data.AuthorizeEditorDataProcessorsXPathProcessorDTO, d = p.expandXPATHProcessor(ctx)
		diags.Append(d...)
	default:
		diags.AddError(
			"Invalid processor type",
			fmt.Sprintf("The processor type '%s' is not supported.  Please raise an issue with the provider maintainers.", p.Type.ValueString()),
		)
	}

	if diags.HasError() {
		return nil, diags
	}

	return &data, diags
}

func (p *editorDataProcessorLeafResourceModel) expand(ctx context.Context) (*authorize.AuthorizeEditorDataProcessorDTO, diag.Diagnostics) {
	var diags, d diag.Diagnostics

	data := authorize.AuthorizeEditorDataProcessorDTO{}

	switch authorize.EnumAuthorizeEditorDataProcessorDTOType(p.Type.ValueString()) {
	case authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_CHAIN:
		diags.AddError(
			"Invalid leaf processor type",
			fmt.Sprintf("The leaf processor type '%s' is not supported.  Please raise an issue with the provider maintainers.", p.Type.ValueString()),
		)
	case authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_COLLECTION_FILTER:
		diags.AddError(
			"Invalid leaf processor type",
			fmt.Sprintf("The leaf processor type '%s' is not supported.  Please raise an issue with the provider maintainers.", p.Type.ValueString()),
		)
	case authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_COLLECTION_TRANSFORM:
		diags.AddError(
			"Invalid leaf processor type",
			fmt.Sprintf("The leaf processor type '%s' is not supported.  Please raise an issue with the provider maintainers.", p.Type.ValueString()),
		)
	case authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_JSON_PATH:
		data.AuthorizeEditorDataProcessorsJsonPathProcessorDTO, d = p.expandJsonPathProcessor(ctx)
		diags.Append(d...)
	case authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_REFERENCE:
		data.AuthorizeEditorDataProcessorsReferenceProcessorDTO, d = p.expandReferenceProcessor(ctx)
		diags.Append(d...)
	case authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_SPEL:
		data.AuthorizeEditorDataProcessorsSpelProcessorDTO, d = p.expandSPELProcessor(ctx)
		diags.Append(d...)
	case authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_XPATH:
		data.AuthorizeEditorDataProcessorsXPathProcessorDTO, d = p.expandXPATHProcessor(ctx)
		diags.Append(d...)
	default:
		diags.AddError(
			"Invalid processor type",
			fmt.Sprintf("The processor type '%s' is not supported.  Please raise an issue with the provider maintainers.", p.Type.ValueString()),
		)
	}

	if diags.HasError() {
		return nil, diags
	}

	return &data, diags
}

func (p *editorDataProcessorResourceModel) expandChainProcessor(ctx context.Context, iteration int32) (*authorize.AuthorizeEditorDataProcessorsChainProcessorDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	processors, d := expandEditorDataProcessorsIteration(ctx, p.Processors, iteration+1)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	data := authorize.NewAuthorizeEditorDataProcessorsChainProcessorDTO(
		p.Name.ValueString(),
		authorize.EnumAuthorizeEditorDataProcessorDTOType(p.Type.ValueString()),
		processors,
	)

	return data, diags
}

func (p *editorDataProcessorResourceModel) expandCollectionFilterProcessor(ctx context.Context, iteration int32) (*authorize.AuthorizeEditorDataProcessorsCollectionFilterProcessorDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	predicate, d := expandEditorDataProcessorIteration(ctx, p.Predicate, iteration+1)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	data := authorize.NewAuthorizeEditorDataProcessorsCollectionFilterProcessorDTO(
		p.Name.ValueString(),
		authorize.EnumAuthorizeEditorDataProcessorDTOType(p.Type.ValueString()),
		*predicate,
	)

	return data, diags
}

func (p *editorDataProcessorResourceModel) expandCollectionTransformProcessor(ctx context.Context, iteration int32) (*authorize.AuthorizeEditorDataProcessorsCollectionTransformProcessorDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	processor, d := expandEditorDataProcessorIteration(ctx, p.Processor, iteration+1)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	data := authorize.NewAuthorizeEditorDataProcessorsCollectionTransformProcessorDTO(
		p.Name.ValueString(),
		authorize.EnumAuthorizeEditorDataProcessorDTOType(p.Type.ValueString()),
		*processor,
	)

	return data, diags
}

func (p *editorDataProcessorResourceModel) expandJsonPathProcessor(ctx context.Context) (*authorize.AuthorizeEditorDataProcessorsJsonPathProcessorDTO, diag.Diagnostics) {
	return expandJsonPathProcessor(ctx, p.ValueType, p.Name, p.Expression)
}

func (p *editorDataProcessorResourceModel) expandReferenceProcessor(ctx context.Context) (*authorize.AuthorizeEditorDataProcessorsReferenceProcessorDTO, diag.Diagnostics) {
	return expandReferenceProcessor(ctx, p.ProcessorRef, p.Name)
}

func (p *editorDataProcessorResourceModel) expandSPELProcessor(ctx context.Context) (*authorize.AuthorizeEditorDataProcessorsSpelProcessorDTO, diag.Diagnostics) {
	return expandSPELProcessor(ctx, p.ValueType, p.Name, p.Expression)
}

func (p *editorDataProcessorResourceModel) expandXPATHProcessor(ctx context.Context) (*authorize.AuthorizeEditorDataProcessorsXPathProcessorDTO, diag.Diagnostics) {
	return expandXPATHProcessor(ctx, p.ValueType, p.Name, p.Expression)
}

func (p *editorDataProcessorLeafResourceModel) expandJsonPathProcessor(ctx context.Context) (*authorize.AuthorizeEditorDataProcessorsJsonPathProcessorDTO, diag.Diagnostics) {
	return expandJsonPathProcessor(ctx, p.ValueType, p.Name, p.Expression)
}

func (p *editorDataProcessorLeafResourceModel) expandReferenceProcessor(ctx context.Context) (*authorize.AuthorizeEditorDataProcessorsReferenceProcessorDTO, diag.Diagnostics) {
	return expandReferenceProcessor(ctx, p.ProcessorRef, p.Name)
}

func (p *editorDataProcessorLeafResourceModel) expandSPELProcessor(ctx context.Context) (*authorize.AuthorizeEditorDataProcessorsSpelProcessorDTO, diag.Diagnostics) {
	return expandSPELProcessor(ctx, p.ValueType, p.Name, p.Expression)
}

func (p *editorDataProcessorLeafResourceModel) expandXPATHProcessor(ctx context.Context) (*authorize.AuthorizeEditorDataProcessorsXPathProcessorDTO, diag.Diagnostics) {
	return expandXPATHProcessor(ctx, p.ValueType, p.Name, p.Expression)
}

func expandJsonPathProcessor(ctx context.Context, valueTypeP basetypes.ObjectValue, name, expression basetypes.StringValue) (*authorize.AuthorizeEditorDataProcessorsJsonPathProcessorDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	valueType, d := expandEditorValueType(ctx, valueTypeP)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	data := authorize.NewAuthorizeEditorDataProcessorsJsonPathProcessorDTO(
		name.ValueString(),
		authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_JSON_PATH,
		expression.ValueString(),
		*valueType,
	)

	return data, diags
}

func expandReferenceProcessor(ctx context.Context, referenceData basetypes.ObjectValue, name basetypes.StringValue) (*authorize.AuthorizeEditorDataProcessorsReferenceProcessorDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	processorRef, d := expandEditorReferenceData(ctx, referenceData)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	data := authorize.NewAuthorizeEditorDataProcessorsReferenceProcessorDTO(
		name.ValueString(),
		authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_REFERENCE,
		*processorRef,
	)

	return data, diags
}

func expandSPELProcessor(ctx context.Context, valueTypeP basetypes.ObjectValue, name, expression basetypes.StringValue) (*authorize.AuthorizeEditorDataProcessorsSpelProcessorDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	valueType, d := expandEditorValueType(ctx, valueTypeP)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	data := authorize.NewAuthorizeEditorDataProcessorsSpelProcessorDTO(
		name.ValueString(),
		authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_SPEL,
		expression.ValueString(),
		*valueType,
	)

	return data, diags
}

func expandXPATHProcessor(ctx context.Context, valueTypeP basetypes.ObjectValue, name, expression basetypes.StringValue) (*authorize.AuthorizeEditorDataProcessorsXPathProcessorDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	valueType, d := expandEditorValueType(ctx, valueTypeP)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	data := authorize.NewAuthorizeEditorDataProcessorsXPathProcessorDTO(
		name.ValueString(),
		authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_XPATH,
		expression.ValueString(),
		*valueType,
	)

	return data, diags
}

func editorDataProcessorOkToTF(ctx context.Context, apiObject *authorize.AuthorizeEditorDataProcessorDTO, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	const initialIteration = 1
	return editorDataProcessorOkToTFIteration(ctx, initialIteration, apiObject, ok)
}

func editorDataProcessorOkToTFIteration(ctx context.Context, iteration int32, apiObject *authorize.AuthorizeEditorDataProcessorDTO, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil || cmp.Equal(apiObject, &authorize.AuthorizeEditorDataProcessorDTO{}) {
		return types.ObjectNull(initializeEditorDataProcessorTFObjectTypes(iteration)), diags
	}

	attributeMap := map[string]attr.Value{}

	switch t := apiObject.GetActualInstance().(type) {
	case *authorize.AuthorizeEditorDataProcessorsChainProcessorDTO:

		processorsResp, ok := t.GetProcessorsOk()
		processors, d := editorDataProcessorsOkToListTFIteration(ctx, iteration+1, processorsResp, ok)
		diags.Append(d...)

		attributeMap["name"] = framework.StringOkToTF(t.GetNameOk())
		attributeMap["type"] = framework.EnumOkToTF(t.GetTypeOk())
		attributeMap["processors"] = processors

	case *authorize.AuthorizeEditorDataProcessorsCollectionFilterProcessorDTO:

		predicateResp, ok := t.GetPredicateOk()
		predicate, d := editorDataProcessorOkToTFIteration(ctx, iteration+1, predicateResp, ok)
		diags.Append(d...)

		attributeMap["name"] = framework.StringOkToTF(t.GetNameOk())
		attributeMap["type"] = framework.EnumOkToTF(t.GetTypeOk())
		attributeMap["predicate"] = predicate

	case *authorize.AuthorizeEditorDataProcessorsCollectionTransformProcessorDTO:

		processorResp, ok := t.GetProcessorOk()
		processor, d := editorDataProcessorOkToTFIteration(ctx, iteration+1, processorResp, ok)
		diags.Append(d...)

		attributeMap["name"] = framework.StringOkToTF(t.GetNameOk())
		attributeMap["type"] = framework.EnumOkToTF(t.GetTypeOk())
		attributeMap["processor"] = processor

	case *authorize.AuthorizeEditorDataProcessorsJsonPathProcessorDTO:

		valueType, d := editorValueTypeOkToTF(t.GetValueTypeOk())
		diags.Append(d...)

		attributeMap["name"] = framework.StringOkToTF(t.GetNameOk())
		attributeMap["type"] = framework.EnumOkToTF(t.GetTypeOk())
		attributeMap["expression"] = framework.StringOkToTF(t.GetExpressionOk())
		attributeMap["value_type"] = valueType

	case *authorize.AuthorizeEditorDataProcessorsReferenceProcessorDTO:

		processorRef, d := editorDataReferenceObjectOkToTF(t.GetProcessorOk())
		diags.Append(d...)

		attributeMap["name"] = framework.StringOkToTF(t.GetNameOk())
		attributeMap["type"] = framework.EnumOkToTF(t.GetTypeOk())
		attributeMap["processor_ref"] = processorRef

	case *authorize.AuthorizeEditorDataProcessorsSpelProcessorDTO:

		valueType, d := editorValueTypeOkToTF(t.GetValueTypeOk())
		diags.Append(d...)

		attributeMap["name"] = framework.StringOkToTF(t.GetNameOk())
		attributeMap["type"] = framework.EnumOkToTF(t.GetTypeOk())
		attributeMap["expression"] = framework.StringOkToTF(t.GetExpressionOk())
		attributeMap["value_type"] = valueType

	case *authorize.AuthorizeEditorDataProcessorsXPathProcessorDTO:

		valueType, d := editorValueTypeOkToTF(t.GetValueTypeOk())
		diags.Append(d...)

		attributeMap["name"] = framework.StringOkToTF(t.GetNameOk())
		attributeMap["type"] = framework.EnumOkToTF(t.GetTypeOk())
		attributeMap["expression"] = framework.StringOkToTF(t.GetExpressionOk())
		attributeMap["value_type"] = valueType

	default:
		tflog.Error(ctx, "Invalid processor type", map[string]interface{}{
			"processor type": t,
		})
		diags.AddError(
			"Invalid processor type",
			"The processor type is not supported.  Please raise an issue with the provider maintainers.",
		)
		return types.ObjectNull(initializeEditorDataProcessorTFObjectTypes(iteration)), diags
	}

	attributeMap = editorDataProcessorConvertEmptyValuesToTFNulls(attributeMap, iteration)

	objValue, d := types.ObjectValue(initializeEditorDataProcessorTFObjectTypes(iteration), attributeMap)
	diags.Append(d...)

	return objValue, diags
}

func editorDataProcessorsOkToListTFIteration(ctx context.Context, iteration int32, apiObject []authorize.AuthorizeEditorDataProcessorDTO, ok bool) (basetypes.ListValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	tfObjType := types.ObjectType{AttrTypes: initializeEditorDataProcessorTFObjectTypes(iteration)}

	if !ok || apiObject == nil {
		return types.ListNull(tfObjType), diags
	}

	flattenedList := []attr.Value{}
	for _, v := range apiObject {

		flattenedObj, d := editorDataProcessorOkToTFIteration(ctx, iteration, &v, true)
		diags.Append(d...)
		if diags.HasError() {
			return types.ListNull(tfObjType), diags
		}

		flattenedList = append(flattenedList, flattenedObj)
	}

	returnVar, d := types.ListValue(tfObjType, flattenedList)
	diags.Append(d...)

	return returnVar, diags
}

func editorDataProcessorConvertEmptyValuesToTFNulls(attributeMap map[string]attr.Value, iteration int32) map[string]attr.Value {

	supportedTypes := authorize.AllowedEnumAuthorizeEditorDataProcessorDTOTypeEnumValues

	if iteration >= processorNestedIterationMaxDepth {
		supportedTypes = leafProcessorTypes
	}

	nullMap := map[string]attr.Value{
		"name": types.StringNull(),
		"type": types.StringNull(),
	}

	if slices.Contains(supportedTypes, authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_CHAIN) {
		nullMap["processors"] = types.ListNull(types.ObjectType{AttrTypes: initializeEditorDataProcessorTFObjectTypes(iteration + 1)})
	}

	if slices.Contains(supportedTypes, authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_COLLECTION_FILTER) {
		nullMap["predicate"] = types.ObjectNull(initializeEditorDataProcessorTFObjectTypes(iteration + 1))
	}

	if slices.Contains(supportedTypes, authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_COLLECTION_TRANSFORM) {
		nullMap["processor"] = types.ObjectNull(initializeEditorDataProcessorTFObjectTypes(iteration + 1))
	}

	if slices.Contains(supportedTypes, authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_JSON_PATH) ||
		slices.Contains(supportedTypes, authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_SPEL) ||
		slices.Contains(supportedTypes, authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_XPATH) {
		nullMap["expression"] = types.StringNull()
		nullMap["value_type"] = types.ObjectNull(editorValueTypeTFObjectTypes)
	}

	if slices.Contains(supportedTypes, authorize.ENUMAUTHORIZEEDITORDATAPROCESSORDTOTYPE_REFERENCE) {
		nullMap["processor_ref"] = types.ObjectNull(editorReferenceObjectTFObjectTypes)
	}

	for k := range nullMap {
		if attributeMap[k] == nil {
			attributeMap[k] = nullMap[k]
		}
	}

	return attributeMap
}
