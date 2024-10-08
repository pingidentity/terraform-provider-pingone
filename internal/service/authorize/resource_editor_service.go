package authorize

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/patrickcping/pingone-go-sdk-v2/authorize"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	int32validatorinternal "github.com/pingidentity/terraform-provider-pingone/internal/framework/int32validator"
	listvalidatorinternal "github.com/pingidentity/terraform-provider-pingone/internal/framework/listvalidator"
	objectvalidatorinternal "github.com/pingidentity/terraform-provider-pingone/internal/framework/objectvalidator"
	setvalidatorinternal "github.com/pingidentity/terraform-provider-pingone/internal/framework/setvalidator"
	stringvalidatorinternal "github.com/pingidentity/terraform-provider-pingone/internal/framework/stringvalidator"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/utils"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type EditorServiceResource serviceClientType

type editorServiceResourceModel struct {
	Id              pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId   pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	Name            types.String                 `tfsdk:"name"`
	FullName        types.String                 `tfsdk:"full_name"`
	Description     types.String                 `tfsdk:"description"`
	Parent          types.Object                 `tfsdk:"parent"`
	Type            types.String                 `tfsdk:"type"`
	CacheSettings   types.Object                 `tfsdk:"cache_settings"`
	ServiceType     types.String                 `tfsdk:"service_type"`
	Version         types.String                 `tfsdk:"version"`
	Processor       types.Object                 `tfsdk:"processor"`
	ValueType       types.Object                 `tfsdk:"value_type"`
	ServiceSettings types.Object                 `tfsdk:"service_settings"`
}

type editorServiceCacheSettingsResourceModel struct {
	TtlSeconds types.Int32 `tfsdk:"ttl_seconds"`
}

type editorServiceServiceSettingsResourceModel struct {
	MaximumConcurrentRequests types.Int32   `tfsdk:"maximum_concurrent_requests"`
	MaximumRequestsPerSecond  types.Float64 `tfsdk:"maximum_requests_per_second"`

	// HTTP
	TimeoutMilliseconds types.Int32  `tfsdk:"timeout_milliseconds"`
	Url                 types.String `tfsdk:"url"`
	Verb                types.String `tfsdk:"verb"`
	Body                types.String `tfsdk:"body"`
	ContentType         types.String `tfsdk:"content_type"`
	Headers             types.Set    `tfsdk:"headers"`
	Authentication      types.Object `tfsdk:"authentication"`
	TlsSettings         types.Object `tfsdk:"tls_settings"`

	// Connector
	Channel       types.String `tfsdk:"channel"`
	Code          types.String `tfsdk:"code"`
	Capability    types.String `tfsdk:"capability"`
	SchemaVersion types.Int32  `tfsdk:"schema_version"`
	InputMappings types.List   `tfsdk:"input_mappings"`
}

type editorServiceServiceSettingsHeaderResourceModel struct {
	Key   types.String `tfsdk:"key"`
	Value types.Object `tfsdk:"value"`
}

type editorServiceServiceSettingsAuthenticationResourceModel struct {
	Type          types.String `tfsdk:"type"`
	Name          types.Object `tfsdk:"name"`
	Password      types.Object `tfsdk:"password"`
	TokenEndpoint types.String `tfsdk:"token_endpoint"`
	ClientId      types.String `tfsdk:"client_id"`
	ClientSecret  types.Object `tfsdk:"client_secret"`
	Scope         types.String `tfsdk:"scope"`
	Token         types.Object `tfsdk:"token"`
}

type editorServiceServiceSettingsTlsSettingsResourceModel struct {
	TlsValidationType types.String `tfsdk:"tls_validation_type"`
}

type editorServiceServiceSettingsInputMappingResourceModel struct {
	Property types.String `tfsdk:"property"`
	Type     types.String `tfsdk:"type"`
	ValueRef types.Object `tfsdk:"value_ref"`
	Value    types.String `tfsdk:"value"`
}

var (
	editorServiceParentTFObjectTypes = map[string]attr.Type{
		"id": pingonetypes.ResourceIDType{},
	}

	editorServiceCacheSettingsTFObjectTypes = map[string]attr.Type{
		"ttl_seconds": types.Int32Type,
	}

	editorServiceValueTypeTFObjectTypes = map[string]attr.Type{
		"type": types.StringType,
	}

	editorServiceServiceSettingsTFObjectTypes = map[string]attr.Type{
		"maximum_concurrent_requests": types.Int32Type,
		"maximum_requests_per_second": types.Float64Type,

		"timeout_milliseconds": types.Int32Type,
		"url":                  types.StringType,
		"verb":                 types.StringType,
		"body":                 types.StringType,
		"content_type":         types.StringType,
		"headers":              types.SetType{ElemType: types.ObjectType{AttrTypes: editorServiceServiceSettingsHeadersTFObjectTypes}},
		"authentication":       types.ObjectType{AttrTypes: editorServiceServiceSettingsAuthenticationTFObjectTypes},
		"tls_settings":         types.ObjectType{AttrTypes: editorServiceServiceSettingsTlsSettingsTFObjectTypes},

		"channel":        types.StringType,
		"code":           types.StringType,
		"capability":     types.StringType,
		"schema_version": types.Int32Type,
		"input_mappings": types.ListType{ElemType: types.ObjectType{AttrTypes: editorServiceServiceSettingsInputMappingsTFObjectTypes}},
	}

	editorServiceServiceSettingsHeadersTFObjectTypes = map[string]attr.Type{
		"key":   types.StringType,
		"value": types.ObjectType{AttrTypes: editorDataInputTFObjectTypes},
	}

	editorServiceServiceSettingsAuthenticationTFObjectTypes = map[string]attr.Type{
		"type":           types.StringType,
		"name":           types.ObjectType{AttrTypes: editorReferenceObjectTFObjectTypes},
		"password":       types.ObjectType{AttrTypes: editorReferenceObjectTFObjectTypes},
		"token_endpoint": types.StringType,
		"client_id":      types.StringType,
		"client_secret":  types.ObjectType{AttrTypes: editorReferenceObjectTFObjectTypes},
		"scope":          types.StringType,
		"token":          types.ObjectType{AttrTypes: editorReferenceObjectTFObjectTypes},
	}

	editorServiceServiceSettingsTlsSettingsTFObjectTypes = map[string]attr.Type{
		"tls_validation_type": types.StringType,
	}

	editorServiceServiceSettingsInputMappingsTFObjectTypes = map[string]attr.Type{
		"property":  types.StringType,
		"type":      types.StringType,
		"value_ref": types.ObjectType{AttrTypes: editorReferenceObjectTFObjectTypes},
		"value":     types.StringType,
	}
)

// Framework interfaces
var (
	_ resource.Resource                = &EditorServiceResource{}
	_ resource.ResourceWithConfigure   = &EditorServiceResource{}
	_ resource.ResourceWithImportState = &EditorServiceResource{}
)

// New Object
func NewEditorServiceResource() resource.Resource {
	return &EditorServiceResource{}
}

// Metadata
func (r *EditorServiceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_authorize_editor_service"
}

func (r *EditorServiceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	// schema descriptions and validation settings
	const attrMinLength = 1

	serviceTypeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The type of service.",
	).AllowedValuesEnum(authorize.AllowedEnumAuthorizeEditorDataDefinitionsServiceDefinitionDTOTypeEnumValues)

	processorDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"",
	).AppendMarkdownString(fmt.Sprintf("This field is optional when `service_type` is `%s` or `%s`.", string(authorize.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOTYPE_CONNECTOR), string(authorize.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOTYPE_HTTP)))

	valueTypeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"",
	).AppendMarkdownString(fmt.Sprintf("This field is required when `service_type` is `%s` or `%s`.", string(authorize.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOTYPE_CONNECTOR), string(authorize.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOTYPE_HTTP)))

	serviceSettingsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"",
	).AppendMarkdownString(fmt.Sprintf("This field is required when `service_type` is `%s` or `%s`.", string(authorize.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOTYPE_CONNECTOR), string(authorize.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOTYPE_HTTP)))

	serviceSettingsUrlDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"",
	).AppendMarkdownString(fmt.Sprintf("This field is required when `service_type` is `%s`.", string(authorize.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOTYPE_HTTP)))

	serviceSettingsVerbDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"",
	).AppendMarkdownString(fmt.Sprintf("This field is required when `service_type` is `%s`.", string(authorize.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOTYPE_HTTP)))

	serviceSettingsBodyDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"",
	).AppendMarkdownString(fmt.Sprintf("This field is optional when `service_type` is `%s`.", string(authorize.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOTYPE_HTTP)))

	serviceSettingsContentTypeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"",
	).AppendMarkdownString(fmt.Sprintf("This field is optional when `service_type` is `%s`.", string(authorize.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOTYPE_HTTP)))

	serviceSettingsHeadersDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"",
	).AppendMarkdownString(fmt.Sprintf("This field is optional when `service_type` is `%s`.", string(authorize.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOTYPE_HTTP)))

	serviceSettingsAuthenticationDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"",
	).AppendMarkdownString(fmt.Sprintf("This field is required when `service_type` is `%s`.", string(authorize.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOTYPE_HTTP)))

	serviceSettingsTlsSettingsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"",
	).AppendMarkdownString(fmt.Sprintf("This field is required when `service_type` is `%s`.", string(authorize.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOTYPE_HTTP)))

	serviceSettingsChannelDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"",
	).AppendMarkdownString(fmt.Sprintf("This field is required when `service_type` is `%s`.", string(authorize.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOTYPE_CONNECTOR)))

	serviceSettingsCodeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"",
	).AppendMarkdownString(fmt.Sprintf("This field is required when `service_type` is `%s`.", string(authorize.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOTYPE_CONNECTOR)))

	serviceSettingsCapabilityDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"",
	).AppendMarkdownString(fmt.Sprintf("This field is required when `service_type` is `%s`.", string(authorize.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOTYPE_CONNECTOR)))

	serviceSettingsSchemaVersionDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"",
	).AppendMarkdownString(fmt.Sprintf("This field is optional when `service_type` is `%s`.", string(authorize.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOTYPE_CONNECTOR)))

	serviceSettingsInputMappingsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"",
	).AppendMarkdownString(fmt.Sprintf("This field is required when `service_type` is `%s`.", string(authorize.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOTYPE_CONNECTOR)))

	serviceSettingsAuthenticationTypeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"",
	).AllowedValuesEnum(authorize.AllowedEnumAuthorizeEditorDataAuthenticationDTOTypeEnumValues)

	serviceSettingsAuthenticationNameDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"",
	).AppendMarkdownString(fmt.Sprintf("This field is required when `type` is `%s`.", string(authorize.ENUMAUTHORIZEEDITORDATAAUTHENTICATIONDTOTYPE_BASIC)))

	serviceSettingsAuthenticationPasswordDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"",
	).AppendMarkdownString(fmt.Sprintf("This field is required when `type` is `%s`.", string(authorize.ENUMAUTHORIZEEDITORDATAAUTHENTICATIONDTOTYPE_BASIC)))

	serviceSettingsAuthenticationTokenEndpointDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"",
	).AppendMarkdownString(fmt.Sprintf("This field is required when `type` is `%s`.", string(authorize.ENUMAUTHORIZEEDITORDATAAUTHENTICATIONDTOTYPE_CLIENT_CREDENTIALS)))

	serviceSettingsAuthenticationClientIdDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"",
	).AppendMarkdownString(fmt.Sprintf("This field is required when `type` is `%s`.", string(authorize.ENUMAUTHORIZEEDITORDATAAUTHENTICATIONDTOTYPE_CLIENT_CREDENTIALS)))

	serviceSettingsAuthenticationClientSecretDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"",
	).AppendMarkdownString(fmt.Sprintf("This field is required when `type` is `%s`.", string(authorize.ENUMAUTHORIZEEDITORDATAAUTHENTICATIONDTOTYPE_CLIENT_CREDENTIALS)))

	serviceSettingsAuthenticationScopeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"",
	).AppendMarkdownString(fmt.Sprintf("This field is required when `type` is `%s`.", string(authorize.ENUMAUTHORIZEEDITORDATAAUTHENTICATIONDTOTYPE_CLIENT_CREDENTIALS)))

	serviceSettingsAuthenticationTokenDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"",
	).AppendMarkdownString(fmt.Sprintf("This field is required when `type` is `%s`.", string(authorize.ENUMAUTHORIZEEDITORDATAAUTHENTICATIONDTOTYPE_TOKEN)))

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage an authorization service for the PingOne Authorize Trust Framework in a PingOne environment.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment to configure the Authorize editor service in."),
			),

			"name": schema.StringAttribute{ // DONE
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A user-friendly service name.").Description,
				Required:    true,

				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
				},
			},

			"full_name": schema.StringAttribute{ // DOC ISSUE
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A unique name generated by the system for each service resource. It is the concatenation of names in the service resource hierarchy.").Description,
				Optional:    true,
			},

			"description": schema.StringAttribute{ // DONE
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The authorization service resource's description.").Description,
				Optional:    true,
			},

			"parent": parentObjectSchema("service"),

			"type": schema.StringAttribute{ // DOC ISSUE
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The resource type.").Description,
				Computed:    true,
			},

			"cache_settings": schema.SingleNestedAttribute{ // DONE
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The service's cache settings.").Description,
				Optional:    true,

				Attributes: map[string]schema.Attribute{
					"ttl_seconds": schema.Int32Attribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("The service's time to live in seconds.").Description,
						Optional:    true,
					},
				},
			},

			"service_type": schema.StringAttribute{ // DONE
				Description:         serviceTypeDescription.Description,
				MarkdownDescription: serviceTypeDescription.MarkdownDescription,
				Required:            true,

				Validators: []validator.String{
					stringvalidator.OneOf(utils.EnumSliceToStringSlice(authorize.AllowedEnumAuthorizeEditorDataDefinitionsServiceDefinitionDTOTypeEnumValues)...),
				},
			},

			"version": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A random ID generated by the system for concurrency control purposes.").Description,
				Computed:    true,
			},

			// service_type == "CONNECTOR", service_type == "HTTP"
			"processor": schema.SingleNestedAttribute{
				Description:         processorDescription.Description,
				MarkdownDescription: processorDescription.MarkdownDescription,
				Optional:            true,

				Validators: []validator.Object{
					objectvalidatorinternal.ConflictsIfMatchesPathValue(
						types.StringValue(string(authorize.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOTYPE_NONE)),
						path.MatchRelative().AtParent().AtName("type"),
					),
				},

				Attributes: dataProcessorObjectSchemaAttributes(),
			},

			"value_type": schema.SingleNestedAttribute{
				Description:         valueTypeDescription.Description,
				MarkdownDescription: valueTypeDescription.MarkdownDescription,
				Optional:            true,

				Validators: []validator.Object{
					objectvalidatorinternal.IsRequiredIfMatchesPathValue(
						types.StringValue(string(authorize.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOTYPE_CONNECTOR)),
						path.MatchRoot("service_type"),
					),
					objectvalidatorinternal.IsRequiredIfMatchesPathValue(
						types.StringValue(string(authorize.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOTYPE_HTTP)),
						path.MatchRoot("service_type"),
					),
					objectvalidatorinternal.ConflictsIfMatchesPathValue(
						types.StringValue(string(authorize.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOTYPE_NONE)),
						path.MatchRoot("service_type"),
					),
				},

				Attributes: valueTypeObjectSchemaAttributes(),
			},

			"service_settings": schema.SingleNestedAttribute{
				Description:         serviceSettingsDescription.Description,
				MarkdownDescription: serviceSettingsDescription.MarkdownDescription,
				Optional:            true,

				Validators: []validator.Object{
					objectvalidatorinternal.IsRequiredIfMatchesPathValue(
						types.StringValue(string(authorize.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOTYPE_CONNECTOR)),
						path.MatchRoot("service_type"),
					),
					objectvalidatorinternal.IsRequiredIfMatchesPathValue(
						types.StringValue(string(authorize.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOTYPE_HTTP)),
						path.MatchRoot("service_type"),
					),
					objectvalidatorinternal.ConflictsIfMatchesPathValue(
						types.StringValue(string(authorize.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOTYPE_NONE)),
						path.MatchRoot("service_type"),
					),
				},

				Attributes: map[string]schema.Attribute{
					"maximum_concurrent_requests": schema.Int32Attribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
						Optional:    true,
					},

					"maximum_requests_per_second": schema.Float64Attribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
						Optional:    true,
					},

					"timeout_milliseconds": schema.Int32Attribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
						Optional:    true,
					},

					"url": schema.StringAttribute{
						Description:         serviceSettingsUrlDescription.Description,
						MarkdownDescription: serviceSettingsUrlDescription.MarkdownDescription,
						Optional:            true,

						Validators: []validator.String{
							stringvalidatorinternal.IsRequiredIfMatchesPathValue(
								types.StringValue(string(authorize.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOTYPE_HTTP)),
								path.MatchRoot("service_type"),
							),
							stringvalidatorinternal.ConflictsIfMatchesPathValue(
								types.StringValue(string(authorize.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOTYPE_CONNECTOR)),
								path.MatchRoot("service_type"),
							),
						},
					},

					"verb": schema.StringAttribute{
						Description:         serviceSettingsVerbDescription.Description,
						MarkdownDescription: serviceSettingsVerbDescription.MarkdownDescription,
						Optional:            true,

						Validators: []validator.String{
							stringvalidatorinternal.IsRequiredIfMatchesPathValue(
								types.StringValue(string(authorize.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOTYPE_HTTP)),
								path.MatchRoot("service_type"),
							),
							stringvalidatorinternal.ConflictsIfMatchesPathValue(
								types.StringValue(string(authorize.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOTYPE_CONNECTOR)),
								path.MatchRoot("service_type"),
							),
						},
					},

					"body": schema.StringAttribute{
						Description:         serviceSettingsBodyDescription.Description,
						MarkdownDescription: serviceSettingsBodyDescription.MarkdownDescription,
						Optional:            true,

						Validators: []validator.String{
							stringvalidatorinternal.ConflictsIfMatchesPathValue(
								types.StringValue(string(authorize.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOTYPE_CONNECTOR)),
								path.MatchRoot("service_type"),
							),
						},
					},

					"content_type": schema.StringAttribute{
						Description:         serviceSettingsContentTypeDescription.Description,
						MarkdownDescription: serviceSettingsContentTypeDescription.MarkdownDescription,
						Optional:            true,

						Validators: []validator.String{
							stringvalidatorinternal.ConflictsIfMatchesPathValue(
								types.StringValue(string(authorize.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOTYPE_CONNECTOR)),
								path.MatchRoot("service_type"),
							),
						},
					},

					"headers": schema.SetNestedAttribute{
						Description:         serviceSettingsHeadersDescription.Description,
						MarkdownDescription: serviceSettingsHeadersDescription.MarkdownDescription,
						Optional:            true,

						Validators: []validator.Set{
							setvalidatorinternal.ConflictsIfMatchesPathValue(
								types.StringValue(string(authorize.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOTYPE_CONNECTOR)),
								path.MatchRoot("service_type"),
							),
						},

						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"key": schema.StringAttribute{
									Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
									Required:    true,
								},

								"value": schema.SingleNestedAttribute{
									Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
									Optional:    true,

									Attributes: dataInputObjectSchemaAttributes(),
								},
							},
						},
					},

					"authentication": schema.SingleNestedAttribute{
						Description:         serviceSettingsAuthenticationDescription.Description,
						MarkdownDescription: serviceSettingsAuthenticationDescription.MarkdownDescription,
						Optional:            true,

						Validators: []validator.Object{
							objectvalidatorinternal.IsRequiredIfMatchesPathValue(
								types.StringValue(string(authorize.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOTYPE_HTTP)),
								path.MatchRoot("service_type"),
							),
							objectvalidatorinternal.ConflictsIfMatchesPathValue(
								types.StringValue(string(authorize.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOTYPE_CONNECTOR)),
								path.MatchRoot("service_type"),
							),
						},

						Attributes: map[string]schema.Attribute{
							"type": schema.StringAttribute{
								Description:         serviceSettingsAuthenticationTypeDescription.Description,
								MarkdownDescription: serviceSettingsAuthenticationTypeDescription.MarkdownDescription,
								Required:            true,

								Validators: []validator.String{
									stringvalidator.OneOf(utils.EnumSliceToStringSlice(authorize.AllowedEnumAuthorizeEditorDataAuthenticationDTOTypeEnumValues)...),
								},
							},

							// type == "BASIC"
							"name": schema.SingleNestedAttribute{
								Description:         serviceSettingsAuthenticationNameDescription.Description,
								MarkdownDescription: serviceSettingsAuthenticationNameDescription.MarkdownDescription,
								Optional:            true,

								Validators: []validator.Object{
									objectvalidatorinternal.IsRequiredIfMatchesPathValue(
										types.StringValue(string(authorize.ENUMAUTHORIZEEDITORDATAAUTHENTICATIONDTOTYPE_BASIC)),
										path.MatchRoot("type"),
									),
								},

								Attributes: referenceIdObjectSchemaAttributes(),
							},

							"password": schema.SingleNestedAttribute{
								Description:         serviceSettingsAuthenticationPasswordDescription.Description,
								MarkdownDescription: serviceSettingsAuthenticationPasswordDescription.MarkdownDescription,
								Optional:            true,

								Validators: []validator.Object{
									objectvalidatorinternal.IsRequiredIfMatchesPathValue(
										types.StringValue(string(authorize.ENUMAUTHORIZEEDITORDATAAUTHENTICATIONDTOTYPE_BASIC)),
										path.MatchRoot("type"),
									),
								},

								Attributes: referenceIdObjectSchemaAttributes(),
							},

							// type == "CLIENT_CREDENTIALS"
							"token_endpoint": schema.StringAttribute{
								Description:         serviceSettingsAuthenticationTokenEndpointDescription.Description,
								MarkdownDescription: serviceSettingsAuthenticationTokenEndpointDescription.MarkdownDescription,
								Optional:            true,

								Validators: []validator.String{
									stringvalidatorinternal.IsRequiredIfMatchesPathValue(
										types.StringValue(string(authorize.ENUMAUTHORIZEEDITORDATAAUTHENTICATIONDTOTYPE_CLIENT_CREDENTIALS)),
										path.MatchRoot("type"),
									),
								},
							},

							"client_id": schema.StringAttribute{
								Description:         serviceSettingsAuthenticationClientIdDescription.Description,
								MarkdownDescription: serviceSettingsAuthenticationClientIdDescription.MarkdownDescription,
								Optional:            true,

								Validators: []validator.String{
									stringvalidatorinternal.IsRequiredIfMatchesPathValue(
										types.StringValue(string(authorize.ENUMAUTHORIZEEDITORDATAAUTHENTICATIONDTOTYPE_CLIENT_CREDENTIALS)),
										path.MatchRoot("type"),
									),
								},
							},

							"client_secret": schema.SingleNestedAttribute{
								Description:         serviceSettingsAuthenticationClientSecretDescription.Description,
								MarkdownDescription: serviceSettingsAuthenticationClientSecretDescription.MarkdownDescription,
								Optional:            true,

								Validators: []validator.Object{
									objectvalidatorinternal.IsRequiredIfMatchesPathValue(
										types.StringValue(string(authorize.ENUMAUTHORIZEEDITORDATAAUTHENTICATIONDTOTYPE_CLIENT_CREDENTIALS)),
										path.MatchRoot("type"),
									),
								},

								Attributes: referenceIdObjectSchemaAttributes(),
							},

							"scope": schema.StringAttribute{
								Description:         serviceSettingsAuthenticationScopeDescription.Description,
								MarkdownDescription: serviceSettingsAuthenticationScopeDescription.MarkdownDescription,
								Optional:            true,

								Validators: []validator.String{
									stringvalidatorinternal.IsRequiredIfMatchesPathValue(
										types.StringValue(string(authorize.ENUMAUTHORIZEEDITORDATAAUTHENTICATIONDTOTYPE_CLIENT_CREDENTIALS)),
										path.MatchRoot("type"),
									),
								},
							},

							// type == "NONE"
							// (same as base model)

							// type == "TOKEN"
							"token": schema.SingleNestedAttribute{
								Description:         serviceSettingsAuthenticationTokenDescription.Description,
								MarkdownDescription: serviceSettingsAuthenticationTokenDescription.MarkdownDescription,
								Optional:            true,

								Validators: []validator.Object{
									objectvalidatorinternal.IsRequiredIfMatchesPathValue(
										types.StringValue(string(authorize.ENUMAUTHORIZEEDITORDATAAUTHENTICATIONDTOTYPE_TOKEN)),
										path.MatchRoot("type"),
									),
								},

								Attributes: referenceIdObjectSchemaAttributes(),
							},
						},
					},

					"tls_settings": schema.SingleNestedAttribute{
						Description:         serviceSettingsTlsSettingsDescription.Description,
						MarkdownDescription: serviceSettingsTlsSettingsDescription.MarkdownDescription,
						Optional:            true,

						Validators: []validator.Object{
							objectvalidatorinternal.IsRequiredIfMatchesPathValue(
								types.StringValue(string(authorize.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOTYPE_HTTP)),
								path.MatchRoot("service_type"),
							),
							objectvalidatorinternal.ConflictsIfMatchesPathValue(
								types.StringValue(string(authorize.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOTYPE_CONNECTOR)),
								path.MatchRoot("service_type"),
							),
						},

						Attributes: map[string]schema.Attribute{
							"tls_validation_type": schema.StringAttribute{
								Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
								Required:    true,
							},
						},
					},

					"channel": schema.StringAttribute{
						Description:         serviceSettingsChannelDescription.Description,
						MarkdownDescription: serviceSettingsChannelDescription.MarkdownDescription,
						Optional:            true,

						Validators: []validator.String{
							stringvalidatorinternal.IsRequiredIfMatchesPathValue(
								types.StringValue(string(authorize.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOTYPE_CONNECTOR)),
								path.MatchRoot("service_type"),
							),
							stringvalidatorinternal.ConflictsIfMatchesPathValue(
								types.StringValue(string(authorize.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOTYPE_HTTP)),
								path.MatchRoot("service_type"),
							),
						},
					},

					"code": schema.StringAttribute{
						Description:         serviceSettingsCodeDescription.Description,
						MarkdownDescription: serviceSettingsCodeDescription.MarkdownDescription,
						Optional:            true,

						Validators: []validator.String{
							stringvalidatorinternal.IsRequiredIfMatchesPathValue(
								types.StringValue(string(authorize.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOTYPE_CONNECTOR)),
								path.MatchRoot("service_type"),
							),
							stringvalidatorinternal.ConflictsIfMatchesPathValue(
								types.StringValue(string(authorize.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOTYPE_HTTP)),
								path.MatchRoot("service_type"),
							),
						},
					},

					"capability": schema.StringAttribute{
						Description:         serviceSettingsCapabilityDescription.Description,
						MarkdownDescription: serviceSettingsCapabilityDescription.MarkdownDescription,
						Optional:            true,

						Validators: []validator.String{
							stringvalidatorinternal.IsRequiredIfMatchesPathValue(
								types.StringValue(string(authorize.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOTYPE_CONNECTOR)),
								path.MatchRoot("service_type"),
							),
							stringvalidatorinternal.ConflictsIfMatchesPathValue(
								types.StringValue(string(authorize.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOTYPE_HTTP)),
								path.MatchRoot("service_type"),
							),
						},
					},

					"schema_version": schema.Int32Attribute{
						Description:         serviceSettingsSchemaVersionDescription.Description,
						MarkdownDescription: serviceSettingsSchemaVersionDescription.MarkdownDescription,
						Optional:            true,

						Validators: []validator.Int32{
							int32validatorinternal.ConflictsIfMatchesPathValue(
								types.StringValue(string(authorize.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOTYPE_HTTP)),
								path.MatchRoot("service_type"),
							),
						},
					},

					"input_mappings": schema.ListNestedAttribute{
						Description:         serviceSettingsInputMappingsDescription.Description,
						MarkdownDescription: serviceSettingsInputMappingsDescription.MarkdownDescription,
						Optional:            true,

						Validators: []validator.List{
							listvalidatorinternal.IsRequiredIfMatchesPathValue(
								types.StringValue(string(authorize.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOTYPE_CONNECTOR)),
								path.MatchRoot("service_type"),
							),
							listvalidatorinternal.ConflictsIfMatchesPathValue(
								types.StringValue(string(authorize.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOTYPE_HTTP)),
								path.MatchRoot("service_type"),
							),
						},

						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"property": schema.StringAttribute{
									Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
									Required:    true,
								},

								"type": schema.StringAttribute{
									Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
									Required:    true,
								},
							},
						},
					},
				},
			},
		},
	}
}

func (r *EditorServiceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *EditorServiceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state editorServiceResourceModel

	if r.Client == nil || r.Client.AuthorizeAPIClient == nil {
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

	// Build the model for the API
	editorService, d := plan.expand(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *authorize.AuthorizeEditorDataDefinitionsServiceDefinitionDTO
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.AuthorizeAPIClient.AuthorizeEditorServicesApi.CreateService(ctx, plan.EnvironmentId.ValueString()).AuthorizeEditorDataDefinitionsServiceDefinitionDTO(*editorService).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"CreateService",
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
	resp.Diagnostics.Append(state.toState(ctx, response)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *EditorServiceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *editorServiceResourceModel

	if r.Client == nil || r.Client.AuthorizeAPIClient == nil {
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
	var response *authorize.AuthorizeEditorDataDefinitionsServiceDefinitionDTO
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.AuthorizeAPIClient.AuthorizeEditorServicesApi.GetService(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"GetService",
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
	resp.Diagnostics.Append(data.toState(ctx, response)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *EditorServiceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state editorServiceResourceModel

	if r.Client == nil || r.Client.AuthorizeAPIClient == nil {
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

	// Build the model for the API
	editorService, d := plan.expand(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *authorize.AuthorizeEditorDataDefinitionsServiceDefinitionDTO
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.AuthorizeAPIClient.AuthorizeEditorServicesApi.UpdateService(ctx, plan.EnvironmentId.ValueString(), plan.Id.ValueString()).AuthorizeEditorDataDefinitionsServiceDefinitionDTO(*editorService).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdateService",
		framework.DefaultCustomError,
		nil,
		&response,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update the state to save
	state = plan

	// Save updated data into Terraform state
	resp.Diagnostics.Append(state.toState(ctx, response)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *EditorServiceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *editorServiceResourceModel

	if r.Client == nil || r.Client.AuthorizeAPIClient == nil {
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
			fR, fErr := r.Client.AuthorizeAPIClient.AuthorizeEditorServicesApi.DeleteService(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), nil, fR, fErr)
		},
		"DeleteService",
		framework.CustomErrorResourceNotFoundWarning,
		nil,
		nil,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *EditorServiceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	idComponents := []framework.ImportComponent{
		{
			Label:  "environment_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
		{
			Label:     "authorize_editor_service_id",
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

func (p *editorServiceResourceModel) expand(ctx context.Context) (*authorize.AuthorizeEditorDataDefinitionsServiceDefinitionDTO, diag.Diagnostics) {
	var diags, d diag.Diagnostics

	data := authorize.AuthorizeEditorDataDefinitionsServiceDefinitionDTO{}

	commonData, d := p.expandCommon(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	switch authorize.EnumAuthorizeEditorDataDefinitionsServiceDefinitionDTOType(p.ServiceType.ValueString()) {
	case authorize.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOTYPE_CONNECTOR:
		data.AuthorizeEditorDataServicesConnectorServiceDefinitionDTO, d = p.expandConnectorService(ctx, commonData)
		diags.Append(d...)
	case authorize.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOTYPE_HTTP:
		data.AuthorizeEditorDataServicesHttpServiceDefinitionDTO, d = p.expandHttpService(ctx, commonData)
		diags.Append(d...)
	case authorize.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOTYPE_NONE:
		data.AuthorizeEditorDataServicesNoneServiceDefinitionDTO, d = p.expandNoneService(commonData)
		diags.Append(d...)
	default:
		diags.AddError(
			"Invalid service type",
			fmt.Sprintf("The service type '%s' is not supported.  Please raise an issue with the provider maintainers.", p.Type.ValueString()),
		)
	}

	if diags.HasError() {
		return nil, diags
	}

	return &data, diags
}

func (p *editorServiceResourceModel) expandCommon(ctx context.Context) (*authorize.AuthorizeEditorDataDefinitionsServiceDefinitionDTOCommon, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := authorize.NewAuthorizeEditorDataDefinitionsServiceDefinitionDTOCommon(
		p.Name.ValueString(),
		p.ServiceType.ValueString(),
	)

	if !p.FullName.IsNull() && !p.FullName.IsUnknown() {
		data.SetFullName(p.FullName.ValueString())
	}

	if !p.Description.IsNull() && !p.Description.IsUnknown() {
		data.SetDescription(p.Description.ValueString())
	}

	if !p.Version.IsNull() && !p.Version.IsUnknown() {
		data.SetVersion(p.Version.ValueString())
	}

	if !p.Parent.IsNull() && !p.Parent.IsUnknown() {
		parent, d := expandEditorParent(ctx, p.Parent)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		data.SetParent(*parent)
	}

	if !p.Type.IsNull() && !p.Type.IsUnknown() {
		data.SetType(authorize.EnumAuthorizeEditorDataDefinitionsServiceDefinitionDTOType(p.Type.ValueString()))
	}

	if !p.CacheSettings.IsNull() && !p.CacheSettings.IsUnknown() {
		var plan *editorServiceCacheSettingsResourceModel
		diags.Append(p.CacheSettings.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		cacheSettings := plan.expand()

		data.SetCacheSettings(*cacheSettings)
	}

	return data, diags
}

func (p *editorServiceCacheSettingsResourceModel) expand() *authorize.AuthorizeEditorDataCacheSettingsDTO {

	data := authorize.NewAuthorizeEditorDataCacheSettingsDTO()

	if !p.TtlSeconds.IsNull() && !p.TtlSeconds.IsUnknown() {
		data.SetTtlSeconds(p.TtlSeconds.ValueInt32())
	}

	return data
}

func (p *editorServiceResourceModel) expandConnectorService(ctx context.Context, commonData *authorize.AuthorizeEditorDataDefinitionsServiceDefinitionDTOCommon) (*authorize.AuthorizeEditorDataServicesConnectorServiceDefinitionDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	if p.ServiceType.ValueString() == string(authorize.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOTYPE_CONNECTOR) {
		return nil, diags
	}

	// Use json.marshall and unmarshal to cast commonData to a AuthorizeEditorDataServicesConnectorServiceDefinitionDTO type
	bytes, err := json.Marshal(commonData)
	if err != nil {
		diags.AddError("Failed to marshal data", err.Error())
		return nil, diags
	}

	var data *authorize.AuthorizeEditorDataServicesConnectorServiceDefinitionDTO
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		diags.AddError("Failed to unmarshal data", err.Error())
		return nil, diags
	}

	valueType, d := expandEditorValueType(ctx, p.ValueType)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	var serviceSettingsPlan *editorServiceServiceSettingsResourceModel
	diags.Append(p.ServiceSettings.As(ctx, &serviceSettingsPlan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)
	if diags.HasError() {
		return nil, diags
	}

	serviceSettings, d := serviceSettingsPlan.expandConnector(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	data.SetValueType(*valueType)
	data.SetServiceSettings(*serviceSettings)

	if !p.Processor.IsNull() && !p.Processor.IsUnknown() {
		processor, d := expandEditorDataProcessor(ctx, p.Processor)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		data.SetProcessor(*processor)
	}

	return data, diags
}

func (p *editorServiceResourceModel) expandHttpService(ctx context.Context, commonData *authorize.AuthorizeEditorDataDefinitionsServiceDefinitionDTOCommon) (*authorize.AuthorizeEditorDataServicesHttpServiceDefinitionDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	if p.ServiceType.ValueString() == string(authorize.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOTYPE_HTTP) {
		return nil, diags
	}

	// Use json.marshall and unmarshal to cast commonData to a AuthorizeEditorDataServicesHttpServiceDefinitionDTO type
	bytes, err := json.Marshal(commonData)
	if err != nil {
		diags.AddError("Failed to marshal data", err.Error())
		return nil, diags
	}

	var data *authorize.AuthorizeEditorDataServicesHttpServiceDefinitionDTO
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		diags.AddError("Failed to unmarshal data", err.Error())
		return nil, diags
	}

	valueType, d := expandEditorValueType(ctx, p.ValueType)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	var serviceSettingsPlan *editorServiceServiceSettingsResourceModel
	diags.Append(p.ServiceSettings.As(ctx, &serviceSettingsPlan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)
	if diags.HasError() {
		return nil, diags
	}

	serviceSettings, d := serviceSettingsPlan.expandHttp(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	data.SetValueType(*valueType)
	data.SetServiceSettings(*serviceSettings)

	if !p.Processor.IsNull() && !p.Processor.IsUnknown() {
		processor, d := expandEditorDataProcessor(ctx, p.Processor)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		data.SetProcessor(*processor)
	}

	return data, diags
}

func (p *editorServiceResourceModel) expandNoneService(commonData *authorize.AuthorizeEditorDataDefinitionsServiceDefinitionDTOCommon) (*authorize.AuthorizeEditorDataServicesNoneServiceDefinitionDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	if p.ServiceType.ValueString() == string(authorize.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOTYPE_NONE) {
		return nil, diags
	}

	// Use json.marshall and unmarshal to cast commonData to a AuthorizeEditorDataServicesNoneServiceDefinitionDTO type
	bytes, err := json.Marshal(commonData)
	if err != nil {
		diags.AddError("Failed to marshal data", err.Error())
		return nil, diags
	}

	var data *authorize.AuthorizeEditorDataServicesNoneServiceDefinitionDTO
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		diags.AddError("Failed to unmarshal data", err.Error())
		return nil, diags
	}

	return data, diags
}

func (p *editorServiceServiceSettingsResourceModel) expandConnector(ctx context.Context) (*authorize.AuthorizeEditorDataServiceSettingsConnectorServiceSettingsDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	inputMappings := make([]authorize.AuthorizeEditorDataInputMappingDTO, 0)

	var inputMappingsPlan []editorServiceServiceSettingsInputMappingResourceModel
	diags.Append(p.InputMappings.ElementsAs(ctx, &inputMappingsPlan, false)...)
	if diags.HasError() {
		return nil, diags
	}

	for _, inputMappingPlan := range inputMappingsPlan {
		inputMapping, d := inputMappingPlan.expand(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		inputMappings = append(inputMappings, *inputMapping)
	}

	data := authorize.NewAuthorizeEditorDataServiceSettingsConnectorServiceSettingsDTO(
		authorize.EnumAuthorizeEditorDataServiceSettingsConnectorServiceSettingsDTOChannel(p.Channel.ValueString()),
		authorize.EnumAuthorizeEditorDataServiceSettingsConnectorServiceSettingsDTOCode(p.Code.ValueString()),
		p.Capability.ValueString(),
		inputMappings,
	)

	if !p.MaximumConcurrentRequests.IsNull() && !p.MaximumConcurrentRequests.IsUnknown() {
		data.SetMaximumConcurrentRequests(p.MaximumConcurrentRequests.ValueInt32())
	}

	if !p.MaximumRequestsPerSecond.IsNull() && !p.MaximumRequestsPerSecond.IsUnknown() {
		data.SetMaximumRequestsPerSecond(p.MaximumRequestsPerSecond.ValueFloat64())
	}

	if !p.SchemaVersion.IsNull() && !p.SchemaVersion.IsUnknown() {
		data.SetSchemaVersion(p.SchemaVersion.ValueInt32())
	}

	return data, diags
}

func (p *editorServiceServiceSettingsInputMappingResourceModel) expand(ctx context.Context) (*authorize.AuthorizeEditorDataInputMappingDTO, diag.Diagnostics) {
	var diags, d diag.Diagnostics

	data := authorize.AuthorizeEditorDataInputMappingDTO{}

	switch authorize.EnumAuthorizeEditorDataInputMappingDTOType(p.Type.ValueString()) {
	case authorize.ENUMAUTHORIZEEDITORDATAINPUTMAPPINGDTOTYPE_ATTRIBUTE:
		data.AuthorizeEditorDataInputMappingsAttributeInputMappingDTO, d = p.expandAttributeType(ctx)
		diags.Append(d...)
	case authorize.ENUMAUTHORIZEEDITORDATAINPUTMAPPINGDTOTYPE_INPUT:
		data.AuthorizeEditorDataInputMappingsInputInputMappingDTO = p.expandInputType()
	default:
		diags.AddError(
			"Invalid input mapping type",
			fmt.Sprintf("The input mapping type '%s' is not supported.  Please raise an issue with the provider maintainers.", p.Type.ValueString()),
		)
	}

	if diags.HasError() {
		return nil, diags
	}

	return &data, diags
}

func (p *editorServiceServiceSettingsInputMappingResourceModel) expandAttributeType(ctx context.Context) (*authorize.AuthorizeEditorDataInputMappingsAttributeInputMappingDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	valueRef, d := expandEditorReferenceData(ctx, p.ValueRef)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	data := authorize.NewAuthorizeEditorDataInputMappingsAttributeInputMappingDTO(
		p.Property.ValueString(),
		authorize.EnumAuthorizeEditorDataInputMappingDTOType(p.Type.ValueString()),
		*valueRef,
	)

	return data, diags
}

func (p *editorServiceServiceSettingsInputMappingResourceModel) expandInputType() *authorize.AuthorizeEditorDataInputMappingsInputInputMappingDTO {

	data := authorize.NewAuthorizeEditorDataInputMappingsInputInputMappingDTO(
		p.Property.ValueString(),
		authorize.EnumAuthorizeEditorDataInputMappingDTOType(p.Type.ValueString()),
		p.Value.ValueString(),
	)

	return data
}

func (p *editorServiceServiceSettingsResourceModel) expandHttp(ctx context.Context) (*authorize.AuthorizeEditorDataServiceSettingsHttpServiceSettingsDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	var authenticationPlan *editorServiceServiceSettingsAuthenticationResourceModel
	diags.Append(p.Authentication.As(ctx, &authenticationPlan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)
	if diags.HasError() {
		return nil, diags
	}

	authentication, d := authenticationPlan.expand(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	var tlsSettingsPlan *editorServiceServiceSettingsTlsSettingsResourceModel
	diags.Append(p.TlsSettings.As(ctx, &tlsSettingsPlan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)
	if diags.HasError() {
		return nil, diags
	}

	tlsSettings := tlsSettingsPlan.expand()

	data := authorize.NewAuthorizeEditorDataServiceSettingsHttpServiceSettingsDTO(
		p.Url.ValueString(),
		authorize.EnumAuthorizeEditorDataServiceSettingsHttpServiceSettingsDTOVerb(p.Verb.ValueString()),
		*authentication,
		*tlsSettings,
	)

	if !p.MaximumConcurrentRequests.IsNull() && !p.MaximumConcurrentRequests.IsUnknown() {
		data.SetMaximumConcurrentRequests(p.MaximumConcurrentRequests.ValueInt32())
	}

	if !p.MaximumRequestsPerSecond.IsNull() && !p.MaximumRequestsPerSecond.IsUnknown() {
		data.SetMaximumRequestsPerSecond(p.MaximumRequestsPerSecond.ValueFloat64())
	}

	if !p.TimeoutMilliseconds.IsNull() && !p.TimeoutMilliseconds.IsUnknown() {
		data.SetTimeoutMilliseconds(p.TimeoutMilliseconds.ValueInt32())
	}

	if !p.Body.IsNull() && !p.Body.IsUnknown() {
		data.SetBody(p.Body.ValueString())
	}

	if !p.ContentType.IsNull() && !p.ContentType.IsUnknown() {
		data.SetContentType(p.ContentType.ValueString())
	}

	if !p.Headers.IsNull() && !p.Headers.IsUnknown() {
		headers := make([]authorize.AuthorizeEditorDataHttpRequestHeaderDTO, 0)

		var plan []editorServiceServiceSettingsHeaderResourceModel
		diags.Append(p.Headers.ElementsAs(ctx, &plan, false)...)
		if diags.HasError() {
			return nil, diags
		}

		for _, headerPlan := range plan {
			header, d := headerPlan.expand(ctx)
			diags.Append(d...)
			if diags.HasError() {
				return nil, diags
			}

			headers = append(headers, *header)
		}

		data.SetHeaders(headers)
	}

	return data, diags
}

func (p *editorServiceServiceSettingsHeaderResourceModel) expand(ctx context.Context) (*authorize.AuthorizeEditorDataHttpRequestHeaderDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := authorize.NewAuthorizeEditorDataHttpRequestHeaderDTO(
		p.Key.ValueString(),
	)

	if !p.Value.IsNull() && !p.Value.IsUnknown() {
		var plan *editorDataInputResourceModel
		diags.Append(p.Value.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		value, d := plan.expand(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		data.SetValue(*value)
	}

	return data, diags
}

func (p *editorServiceServiceSettingsAuthenticationResourceModel) expand(ctx context.Context) (*authorize.AuthorizeEditorDataAuthenticationDTO, diag.Diagnostics) {
	var diags, d diag.Diagnostics

	data := authorize.AuthorizeEditorDataAuthenticationDTO{}

	switch authorize.EnumAuthorizeEditorDataAuthenticationDTOType(p.Type.ValueString()) {
	case authorize.ENUMAUTHORIZEEDITORDATAAUTHENTICATIONDTOTYPE_BASIC:
		data.AuthorizeEditorDataAuthenticationsBasicAuthenticationDTO, d = p.expandBasicAuth(ctx)
		diags.Append(d...)
	case authorize.ENUMAUTHORIZEEDITORDATAAUTHENTICATIONDTOTYPE_CLIENT_CREDENTIALS:
		data.AuthorizeEditorDataAuthenticationsClientCredentialsAuthenticationDTO, d = p.expandClientCredentialsAuth(ctx)
		diags.Append(d...)
	case authorize.ENUMAUTHORIZEEDITORDATAAUTHENTICATIONDTOTYPE_NONE:
		data.AuthorizeEditorDataAuthenticationsNoneAuthenticationDTO = p.expandNoneAuth()
	case authorize.ENUMAUTHORIZEEDITORDATAAUTHENTICATIONDTOTYPE_TOKEN:
		data.AuthorizeEditorDataAuthenticationsTokenAuthenticationDTO, d = p.expandTokenAuth(ctx)
		diags.Append(d...)
	default:
		diags.AddError(
			"Invalid service settings authentication type",
			fmt.Sprintf("The service settings authentication type '%s' is not supported.  Please raise an issue with the provider maintainers.", p.Type.ValueString()),
		)
	}

	if diags.HasError() {
		return nil, diags
	}

	return &data, diags
}

func (p *editorServiceServiceSettingsAuthenticationResourceModel) expandBasicAuth(ctx context.Context) (*authorize.AuthorizeEditorDataAuthenticationsBasicAuthenticationDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	name, d := expandEditorReferenceData(ctx, p.Name)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	password, d := expandEditorReferenceData(ctx, p.Password)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	data := authorize.NewAuthorizeEditorDataAuthenticationsBasicAuthenticationDTO(
		authorize.EnumAuthorizeEditorDataAuthenticationDTOType(p.Type.ValueString()),
		*name,
		*password,
	)

	return data, diags
}

func (p *editorServiceServiceSettingsAuthenticationResourceModel) expandClientCredentialsAuth(ctx context.Context) (*authorize.AuthorizeEditorDataAuthenticationsClientCredentialsAuthenticationDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	clientSecret, d := expandEditorReferenceData(ctx, p.ClientSecret)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	data := authorize.NewAuthorizeEditorDataAuthenticationsClientCredentialsAuthenticationDTO(
		authorize.EnumAuthorizeEditorDataAuthenticationDTOType(p.Type.ValueString()),
		p.TokenEndpoint.ValueString(),
		p.ClientId.ValueString(),
		*clientSecret,
		p.Scope.ValueString(),
	)

	return data, diags
}

func (p *editorServiceServiceSettingsAuthenticationResourceModel) expandNoneAuth() *authorize.AuthorizeEditorDataAuthenticationsNoneAuthenticationDTO {

	data := authorize.NewAuthorizeEditorDataAuthenticationsNoneAuthenticationDTO(
		authorize.EnumAuthorizeEditorDataAuthenticationDTOType(p.Type.ValueString()),
	)

	return data
}

func (p *editorServiceServiceSettingsAuthenticationResourceModel) expandTokenAuth(ctx context.Context) (*authorize.AuthorizeEditorDataAuthenticationsTokenAuthenticationDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	token, d := expandEditorReferenceData(ctx, p.Token)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	data := authorize.NewAuthorizeEditorDataAuthenticationsTokenAuthenticationDTO(
		authorize.EnumAuthorizeEditorDataAuthenticationDTOType(p.Type.ValueString()),
		*token,
	)

	return data, diags
}

func (p *editorServiceServiceSettingsTlsSettingsResourceModel) expand() *authorize.AuthorizeEditorDataTlsSettingsDTO {

	data := authorize.NewAuthorizeEditorDataTlsSettingsDTO(
		authorize.EnumAuthorizeEditorDataTlsSettingsDTOTlsValidationType(p.TlsValidationType.ValueString()),
	)

	return data
}

func (p *editorServiceResourceModel) toState(ctx context.Context, apiObject *authorize.AuthorizeEditorDataDefinitionsServiceDefinitionDTO) diag.Diagnostics {
	var diags, d diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)
		return diags
	}

	apiObjectCommon := authorize.AuthorizeEditorDataDefinitionsServiceDefinitionDTOCommon{}

	if apiObject.AuthorizeEditorDataServicesConnectorServiceDefinitionDTO != nil {
		apiObjectCommon = authorize.AuthorizeEditorDataDefinitionsServiceDefinitionDTOCommon{
			Id:            apiObject.AuthorizeEditorDataServicesConnectorServiceDefinitionDTO.Id,
			Version:       apiObject.AuthorizeEditorDataServicesConnectorServiceDefinitionDTO.Version,
			Name:          apiObject.AuthorizeEditorDataServicesConnectorServiceDefinitionDTO.Name,
			FullName:      apiObject.AuthorizeEditorDataServicesConnectorServiceDefinitionDTO.FullName,
			Description:   apiObject.AuthorizeEditorDataServicesConnectorServiceDefinitionDTO.Description,
			Parent:        apiObject.AuthorizeEditorDataServicesConnectorServiceDefinitionDTO.Parent,
			Type:          apiObject.AuthorizeEditorDataServicesConnectorServiceDefinitionDTO.Type,
			CacheSettings: apiObject.AuthorizeEditorDataServicesConnectorServiceDefinitionDTO.CacheSettings,
			ServiceType:   apiObject.AuthorizeEditorDataServicesConnectorServiceDefinitionDTO.ServiceType,
		}
	}

	if apiObject.AuthorizeEditorDataServicesHttpServiceDefinitionDTO != nil {
		apiObjectCommon = authorize.AuthorizeEditorDataDefinitionsServiceDefinitionDTOCommon{
			Id:            apiObject.AuthorizeEditorDataServicesHttpServiceDefinitionDTO.Id,
			Version:       apiObject.AuthorizeEditorDataServicesHttpServiceDefinitionDTO.Version,
			Name:          apiObject.AuthorizeEditorDataServicesHttpServiceDefinitionDTO.Name,
			FullName:      apiObject.AuthorizeEditorDataServicesHttpServiceDefinitionDTO.FullName,
			Description:   apiObject.AuthorizeEditorDataServicesHttpServiceDefinitionDTO.Description,
			Parent:        apiObject.AuthorizeEditorDataServicesHttpServiceDefinitionDTO.Parent,
			Type:          apiObject.AuthorizeEditorDataServicesHttpServiceDefinitionDTO.Type,
			CacheSettings: apiObject.AuthorizeEditorDataServicesHttpServiceDefinitionDTO.CacheSettings,
			ServiceType:   apiObject.AuthorizeEditorDataServicesHttpServiceDefinitionDTO.ServiceType,
		}
	}

	if apiObject.AuthorizeEditorDataServicesNoneServiceDefinitionDTO != nil {
		apiObjectCommon = authorize.AuthorizeEditorDataDefinitionsServiceDefinitionDTOCommon{
			Id:            apiObject.AuthorizeEditorDataServicesNoneServiceDefinitionDTO.Id,
			Version:       apiObject.AuthorizeEditorDataServicesNoneServiceDefinitionDTO.Version,
			Name:          apiObject.AuthorizeEditorDataServicesNoneServiceDefinitionDTO.Name,
			FullName:      apiObject.AuthorizeEditorDataServicesNoneServiceDefinitionDTO.FullName,
			Description:   apiObject.AuthorizeEditorDataServicesNoneServiceDefinitionDTO.Description,
			Parent:        apiObject.AuthorizeEditorDataServicesNoneServiceDefinitionDTO.Parent,
			Type:          apiObject.AuthorizeEditorDataServicesNoneServiceDefinitionDTO.Type,
			CacheSettings: apiObject.AuthorizeEditorDataServicesNoneServiceDefinitionDTO.CacheSettings,
			ServiceType:   apiObject.AuthorizeEditorDataServicesNoneServiceDefinitionDTO.ServiceType,
		}
	}

	p.Id = framework.PingOneResourceIDOkToTF(apiObjectCommon.GetIdOk())
	//p.EnvironmentId = framework.PingOneResourceIDToTF(*apiObjectCommon.GetEnvironment().Id)
	p.Name = framework.StringOkToTF(apiObjectCommon.GetNameOk())
	p.FullName = framework.StringOkToTF(apiObjectCommon.GetFullNameOk())
	p.Description = framework.StringOkToTF(apiObjectCommon.GetDescriptionOk())

	p.Parent, d = editorParentOkToTF(apiObjectCommon.GetParentOk())
	diags.Append(d...)

	p.Type = framework.EnumOkToTF(apiObjectCommon.GetTypeOk())

	p.CacheSettings, d = editorServiceCacheSettingsOkToTF(apiObjectCommon.GetCacheSettingsOk())
	diags.Append(d...)

	p.ServiceType = framework.StringOkToTF(apiObjectCommon.GetServiceTypeOk())
	p.Version = framework.StringOkToTF(apiObjectCommon.GetVersionOk())

	p.Processor = types.ObjectNull(editorDataProcessorTFObjectTypes)
	p.ValueType = types.ObjectNull(editorServiceValueTypeTFObjectTypes)
	p.ServiceSettings = types.ObjectNull(editorServiceServiceSettingsTFObjectTypes)

	if v := apiObject.AuthorizeEditorDataServicesConnectorServiceDefinitionDTO; v != nil {
		processorVal, ok := v.GetProcessorOk()
		p.Processor, d = editorDataProcessorOkToTF(ctx, processorVal, ok)
		diags.Append(d...)

		p.ValueType, d = editorValueTypeOkToTF(v.GetValueTypeOk())
		diags.Append(d...)

		serviceSettingsVal, ok := v.GetServiceSettingsOk()
		p.ServiceSettings, d = editorServiceServiceSettingsConnectorOkToTF(ctx, serviceSettingsVal, ok)
		diags.Append(d...)
	}

	if v := apiObject.AuthorizeEditorDataServicesHttpServiceDefinitionDTO; v != nil {
		processorVal, ok := v.GetProcessorOk()
		p.Processor, d = editorDataProcessorOkToTF(ctx, processorVal, ok)
		diags.Append(d...)

		p.ValueType, d = editorValueTypeOkToTF(v.GetValueTypeOk())
		diags.Append(d...)

		serviceSettingsVal, ok := v.GetServiceSettingsOk()
		p.ServiceSettings, d = editorServiceServiceSettingsHttpOkToTF(ctx, serviceSettingsVal, ok)
		diags.Append(d...)
	}

	// No implementation for "None" service

	return diags
}

func editorServiceCacheSettingsOkToTF(apiObject *authorize.AuthorizeEditorDataCacheSettingsDTO, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(editorServiceCacheSettingsTFObjectTypes), diags
	}

	objValue, d := types.ObjectValue(editorServiceCacheSettingsTFObjectTypes, map[string]attr.Value{
		"ttl_seconds": framework.Int32OkToTF(apiObject.GetTtlSecondsOk()),
	})
	diags.Append(d...)

	return objValue, diags
}

func editorServiceServiceSettingsConnectorOkToTF(ctx context.Context, apiObject *authorize.AuthorizeEditorDataServiceSettingsConnectorServiceSettingsDTO, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(editorServiceServiceSettingsTFObjectTypes), diags
	}

	inputMappingsVal, ok := apiObject.GetInputMappingsOk()
	inputMappings, d := editorServiceServiceSettingsConnectorInputMappingsOkToTF(ctx, inputMappingsVal, ok)
	diags.Append(d...)

	attributeMap := map[string]attr.Value{
		"maximum_concurrent_requests": framework.Int32OkToTF(apiObject.GetMaximumConcurrentRequestsOk()),
		"maximum_requests_per_second": framework.Float64OkToTF(apiObject.GetMaximumRequestsPerSecondOk()),
		"channel":                     framework.EnumOkToTF(apiObject.GetChannelOk()),
		"code":                        framework.EnumOkToTF(apiObject.GetCodeOk()),
		"capability":                  framework.StringOkToTF(apiObject.GetCapabilityOk()),
		"schema_version":              framework.Int32OkToTF(apiObject.GetSchemaVersionOk()),
		"input_mappings":              inputMappings,
	}

	attributeMap = editorServiceServiceSettingsConvertEmptyValuesToTFNulls(attributeMap)

	objValue, d := types.ObjectValue(editorServiceServiceSettingsTFObjectTypes, attributeMap)
	diags.Append(d...)

	return objValue, diags
}

func editorServiceServiceSettingsConnectorInputMappingsOkToTF(ctx context.Context, apiObject []authorize.AuthorizeEditorDataInputMappingDTO, ok bool) (basetypes.ListValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	tfObjType := types.ObjectType{AttrTypes: editorServiceServiceSettingsInputMappingsTFObjectTypes}

	if !ok || apiObject == nil {
		return types.ListNull(tfObjType), diags
	}

	flattenedList := []attr.Value{}
	for _, v := range apiObject {

		flattenedObj, d := editorServiceServiceSettingsConnectorInputMappingOkToTF(ctx, &v, true)
		diags.Append(d...)

		flattenedList = append(flattenedList, flattenedObj)
	}

	returnVar, d := types.ListValue(tfObjType, flattenedList)
	diags.Append(d...)

	return returnVar, diags
}

func editorServiceServiceSettingsConnectorInputMappingOkToTF(ctx context.Context, apiObject *authorize.AuthorizeEditorDataInputMappingDTO, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil || cmp.Equal(apiObject, &authorize.AuthorizeEditorDataInputMappingDTO{}) {
		return types.ObjectNull(editorServiceServiceSettingsInputMappingsTFObjectTypes), diags
	}

	attributeMap := map[string]attr.Value{}

	switch t := apiObject.GetActualInstance().(type) {
	case authorize.AuthorizeEditorDataInputMappingsAttributeInputMappingDTO:

		valueResp, ok := t.GetValueOk()
		value, d := editorDataReferenceObjectOkToTF(valueResp, ok)
		diags.Append(d...)

		attributeMap = map[string]attr.Value{
			"type":      framework.EnumOkToTF(t.GetTypeOk()),
			"property":  framework.StringOkToTF(t.GetPropertyOk()),
			"value_ref": value,
		}

	case authorize.AuthorizeEditorDataInputMappingsInputInputMappingDTO:

		attributeMap = map[string]attr.Value{
			"type":     framework.EnumOkToTF(t.GetTypeOk()),
			"property": framework.StringOkToTF(t.GetPropertyOk()),
			"value":    framework.StringOkToTF(t.GetValueOk()),
		}

	default:
		tflog.Error(ctx, "Invalid service setting connector input mapping type", map[string]interface{}{
			"service setting connector input mapping type": t,
		})
		diags.AddError(
			"Invalid service setting connector input mapping type",
			"The service setting connector input mapping type is not supported.  Please raise an issue with the provider maintainers.",
		)
	}

	attributeMap = editorServiceServiceSettingsConnectorInputMappingConvertEmptyValuesToTFNulls(attributeMap)

	objValue, d := types.ObjectValue(editorServiceServiceSettingsInputMappingsTFObjectTypes, attributeMap)
	diags.Append(d...)

	return objValue, diags
}

func editorServiceServiceSettingsConnectorInputMappingConvertEmptyValuesToTFNulls(attributeMap map[string]attr.Value) map[string]attr.Value {
	nullMap := map[string]attr.Value{
		"type":       types.StringNull(),
		"comparator": types.StringNull(),
		"value_ref":  types.ObjectNull(editorReferenceObjectTFObjectTypes),
		"value":      types.StringNull(),
	}

	for k := range nullMap {
		if attributeMap[k] == nil {
			attributeMap[k] = nullMap[k]
		}
	}

	return attributeMap
}

func editorServiceServiceSettingsHttpOkToTF(ctx context.Context, apiObject *authorize.AuthorizeEditorDataServiceSettingsHttpServiceSettingsDTO, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(editorServiceServiceSettingsTFObjectTypes), diags
	}

	headersVal, ok := apiObject.GetHeadersOk()
	headers, d := editorServiceServiceSettingsHttpHeadersOkToTF(ctx, headersVal, ok)
	diags.Append(d...)

	authenticationVal, ok := apiObject.GetAuthenticationOk()
	authentication, d := editorServiceServiceSettingsHttpAuthenticationOkToTF(ctx, authenticationVal, ok)
	diags.Append(d...)

	tlsSettings, d := editorServiceServiceSettingsHttpTlsSettingsOkToTF(apiObject.GetTlsSettingsOk())
	diags.Append(d...)

	attributeMap := map[string]attr.Value{
		"maximum_concurrent_requests": framework.Int32OkToTF(apiObject.GetMaximumConcurrentRequestsOk()),
		"maximum_requests_per_second": framework.Float64OkToTF(apiObject.GetMaximumRequestsPerSecondOk()),
		"timeout_milliseconds":        framework.Int32OkToTF(apiObject.GetTimeoutMillisecondsOk()),
		"url":                         framework.StringOkToTF(apiObject.GetUrlOk()),
		"verb":                        framework.EnumOkToTF(apiObject.GetVerbOk()),
		"body":                        framework.StringOkToTF(apiObject.GetBodyOk()),
		"content_type":                framework.StringOkToTF(apiObject.GetContentTypeOk()),
		"headers":                     headers,
		"authentication":              authentication,
		"tls_settings":                tlsSettings,
	}

	attributeMap = editorServiceServiceSettingsConvertEmptyValuesToTFNulls(attributeMap)

	objValue, d := types.ObjectValue(editorServiceServiceSettingsTFObjectTypes, attributeMap)
	diags.Append(d...)

	return objValue, diags
}

func editorServiceServiceSettingsConvertEmptyValuesToTFNulls(attributeMap map[string]attr.Value) map[string]attr.Value {
	nullMap := map[string]attr.Value{
		"maximum_concurrent_requests": types.Int32Null(),
		"maximum_requests_per_second": types.Float64Null(),
		"timeout_milliseconds":        types.Int32Null(),
		"url":                         types.StringNull(),
		"verb":                        types.StringNull(),
		"body":                        types.StringNull(),
		"content_type":                types.StringNull(),
		"headers":                     types.ObjectNull(editorServiceServiceSettingsHeadersTFObjectTypes),
		"authentication":              types.ObjectNull(editorServiceServiceSettingsAuthenticationTFObjectTypes),
		"tls_settings":                types.ObjectNull(editorServiceServiceSettingsTlsSettingsTFObjectTypes),
		"channel":                     types.StringNull(),
		"code":                        types.StringNull(),
		"capability":                  types.StringNull(),
		"schema_version":              types.Int32Null(),
		"input_mappings":              types.ListNull(types.ObjectType{AttrTypes: editorServiceServiceSettingsInputMappingsTFObjectTypes}),
	}

	for k := range nullMap {
		if attributeMap[k] == nil {
			attributeMap[k] = nullMap[k]
		}
	}

	return attributeMap
}

func editorServiceServiceSettingsHttpHeadersOkToTF(ctx context.Context, apiObject []authorize.AuthorizeEditorDataHttpRequestHeaderDTO, ok bool) (basetypes.SetValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	tfObjType := types.ObjectType{AttrTypes: editorServiceServiceSettingsHeadersTFObjectTypes}

	if !ok || apiObject == nil {
		return types.SetNull(tfObjType), diags
	}

	flattenedList := []attr.Value{}
	for _, v := range apiObject {

		dataInputVal, ok := v.GetValueOk()
		value, d := editorDataInputOkToTF(ctx, dataInputVal, ok)
		diags.Append(d...)

		flattenedObj, d := types.ObjectValue(editorServiceServiceSettingsHeadersTFObjectTypes, map[string]attr.Value{
			"key":   framework.StringOkToTF(v.GetKeyOk()),
			"value": value,
		})
		diags.Append(d...)

		flattenedList = append(flattenedList, flattenedObj)
	}

	returnVar, d := types.SetValue(tfObjType, flattenedList)
	diags.Append(d...)

	return returnVar, diags
}

func editorServiceServiceSettingsHttpAuthenticationOkToTF(ctx context.Context, apiObject *authorize.AuthorizeEditorDataAuthenticationDTO, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil || cmp.Equal(apiObject, &authorize.AuthorizeEditorDataAuthenticationDTO{}) {
		return types.ObjectNull(editorServiceServiceSettingsAuthenticationTFObjectTypes), diags
	}

	attributeMap := map[string]attr.Value{}

	switch t := apiObject.GetActualInstance().(type) {
	case authorize.AuthorizeEditorDataAuthenticationsBasicAuthenticationDTO:

		nameResp, ok := t.GetNameOk()
		name, d := editorDataReferenceObjectOkToTF(nameResp, ok)
		diags.Append(d...)

		passwordResp, ok := t.GetPasswordOk()
		password, d := editorDataReferenceObjectOkToTF(passwordResp, ok)
		diags.Append(d...)

		attributeMap = map[string]attr.Value{
			"type":     framework.EnumOkToTF(t.GetTypeOk()),
			"name":     name,
			"password": password,
		}

	case authorize.AuthorizeEditorDataAuthenticationsClientCredentialsAuthenticationDTO:

		clientSecretResp, ok := t.GetClientSecretOk()
		clientSecret, d := editorDataReferenceObjectOkToTF(clientSecretResp, ok)
		diags.Append(d...)

		attributeMap = map[string]attr.Value{
			"type":           framework.EnumOkToTF(t.GetTypeOk()),
			"token_endpoint": framework.StringOkToTF(t.GetTokenEndpointOk()),
			"client_id":      framework.StringOkToTF(t.GetClientIdOk()),
			"client_secret":  clientSecret,
			"scope":          framework.StringOkToTF(t.GetScopeOk()),
		}

	case authorize.AuthorizeEditorDataAuthenticationsNoneAuthenticationDTO:

		attributeMap = map[string]attr.Value{
			"type": framework.EnumOkToTF(t.GetTypeOk()),
		}

	case authorize.AuthorizeEditorDataAuthenticationsTokenAuthenticationDTO:

		tokenResp, ok := t.GetTokenOk()
		token, d := editorDataReferenceObjectOkToTF(tokenResp, ok)
		diags.Append(d...)

		attributeMap = map[string]attr.Value{
			"type":  framework.EnumOkToTF(t.GetTypeOk()),
			"token": token,
		}

	default:
		tflog.Error(ctx, "Invalid service settings authentication type", map[string]interface{}{
			"service settings authentication type": t,
		})
		diags.AddError(
			"Invalid service settings authentication type",
			"The service settings authentication type is not supported.  Please raise an issue with the provider maintainers.",
		)
	}

	attributeMap = editorServiceServiceSettingsAuthenticationConvertEmptyValuesToTFNulls(attributeMap)

	objValue, d := types.ObjectValue(editorServiceServiceSettingsAuthenticationTFObjectTypes, attributeMap)
	diags.Append(d...)

	return objValue, diags
}

func editorServiceServiceSettingsAuthenticationConvertEmptyValuesToTFNulls(attributeMap map[string]attr.Value) map[string]attr.Value {
	nullMap := map[string]attr.Value{
		"type":           types.StringNull(),
		"name":           types.ObjectNull(editorReferenceObjectTFObjectTypes),
		"password":       types.ObjectNull(editorReferenceObjectTFObjectTypes),
		"token_endpoint": types.StringNull(),
		"client_id":      types.StringNull(),
		"client_secret":  types.ObjectNull(editorReferenceObjectTFObjectTypes),
		"scope":          types.StringNull(),
		"token":          types.ObjectNull(editorReferenceObjectTFObjectTypes),
	}

	for k := range nullMap {
		if attributeMap[k] == nil {
			attributeMap[k] = nullMap[k]
		}
	}

	return attributeMap
}

func editorServiceServiceSettingsHttpTlsSettingsOkToTF(apiObject *authorize.AuthorizeEditorDataTlsSettingsDTO, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(editorServiceServiceSettingsTlsSettingsTFObjectTypes), diags
	}

	objValue, d := types.ObjectValue(editorServiceServiceSettingsTlsSettingsTFObjectTypes, map[string]attr.Value{
		"tls_validation_type": framework.EnumOkToTF(apiObject.GetTlsValidationTypeOk()),
	})
	diags.Append(d...)

	return objValue, diags
}
