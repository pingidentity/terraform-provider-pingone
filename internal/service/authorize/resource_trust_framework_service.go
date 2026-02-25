// Copyright © 2026 Ping Identity Corporation

//go:build beta

package authorize

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework-validators/float64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/patrickcping/pingone-go-sdk-v2/authorizeeditor"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	int32validatorinternal "github.com/pingidentity/terraform-provider-pingone/internal/framework/int32validator"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/legacysdk"
	listvalidatorinternal "github.com/pingidentity/terraform-provider-pingone/internal/framework/listvalidator"
	objectvalidatorinternal "github.com/pingidentity/terraform-provider-pingone/internal/framework/objectvalidator"
	setvalidatorinternal "github.com/pingidentity/terraform-provider-pingone/internal/framework/setvalidator"
	stringvalidatorinternal "github.com/pingidentity/terraform-provider-pingone/internal/framework/stringvalidator"
	"github.com/pingidentity/terraform-provider-pingone/internal/utils"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type TrustFrameworkServiceResource serviceClientType

type trustFrameworkServiceResourceModel struct {
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

type trustFrameworkServiceCacheSettingsResourceModel struct {
	TtlSeconds types.Int32 `tfsdk:"ttl_seconds"`
}

type trustFrameworkServiceServiceSettingsResourceModel struct {
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

type trustFrameworkServiceServiceSettingsHeaderResourceModel struct {
	Key   types.String `tfsdk:"key"`
	Value types.Object `tfsdk:"value"`
}

type trustFrameworkServiceServiceSettingsAuthenticationResourceModel struct {
	Type          types.String `tfsdk:"type"`
	Name          types.Object `tfsdk:"name"`
	Password      types.Object `tfsdk:"password"`
	TokenEndpoint types.String `tfsdk:"token_endpoint"`
	ClientId      types.String `tfsdk:"client_id"`
	ClientSecret  types.Object `tfsdk:"client_secret"`
	Scope         types.String `tfsdk:"scope"`
	Token         types.Object `tfsdk:"token"`
}

type trustFrameworkServiceServiceSettingsTlsSettingsResourceModel struct {
	TlsValidationType types.String `tfsdk:"tls_validation_type"`
}

type trustFrameworkServiceServiceSettingsInputMappingResourceModel struct {
	Property types.String `tfsdk:"property"`
	Type     types.String `tfsdk:"type"`
	ValueRef types.Object `tfsdk:"value_ref"`
	Value    types.String `tfsdk:"value"`
}

var (
	trustFrameworkServiceParentTFObjectTypes = map[string]attr.Type{
		"id": pingonetypes.ResourceIDType{},
	}

	trustFrameworkServiceCacheSettingsTFObjectTypes = map[string]attr.Type{
		"ttl_seconds": types.Int32Type,
	}

	trustFrameworkServiceValueTypeTFObjectTypes = map[string]attr.Type{
		"type": types.StringType,
	}

	trustFrameworkServiceServiceSettingsTFObjectTypes = map[string]attr.Type{
		"maximum_concurrent_requests": types.Int32Type,
		"maximum_requests_per_second": types.Float64Type,

		"timeout_milliseconds": types.Int32Type,
		"url":                  types.StringType,
		"verb":                 types.StringType,
		"body":                 types.StringType,
		"content_type":         types.StringType,
		"headers":              types.SetType{ElemType: types.ObjectType{AttrTypes: trustFrameworkServiceServiceSettingsHeadersTFObjectTypes}},
		"authentication":       types.ObjectType{AttrTypes: trustFrameworkServiceServiceSettingsAuthenticationTFObjectTypes},
		"tls_settings":         types.ObjectType{AttrTypes: trustFrameworkServiceServiceSettingsTlsSettingsTFObjectTypes},

		"channel":        types.StringType,
		"code":           types.StringType,
		"capability":     types.StringType,
		"schema_version": types.Int32Type,
		"input_mappings": types.ListType{ElemType: types.ObjectType{AttrTypes: trustFrameworkServiceServiceSettingsInputMappingsTFObjectTypes}},
	}

	trustFrameworkServiceServiceSettingsHeadersTFObjectTypes = map[string]attr.Type{
		"key":   types.StringType,
		"value": types.ObjectType{AttrTypes: editorDataInputTFObjectTypes},
	}

	trustFrameworkServiceServiceSettingsAuthenticationTFObjectTypes = map[string]attr.Type{
		"type":           types.StringType,
		"name":           types.ObjectType{AttrTypes: editorReferenceObjectTFObjectTypes},
		"password":       types.ObjectType{AttrTypes: editorReferenceObjectTFObjectTypes},
		"token_endpoint": types.StringType,
		"client_id":      types.StringType,
		"client_secret":  types.ObjectType{AttrTypes: editorReferenceObjectTFObjectTypes},
		"scope":          types.StringType,
		"token":          types.ObjectType{AttrTypes: editorReferenceObjectTFObjectTypes},
	}

	trustFrameworkServiceServiceSettingsTlsSettingsTFObjectTypes = map[string]attr.Type{
		"tls_validation_type": types.StringType,
	}

	trustFrameworkServiceServiceSettingsInputMappingsTFObjectTypes = map[string]attr.Type{
		"property":  types.StringType,
		"type":      types.StringType,
		"value_ref": types.ObjectType{AttrTypes: editorReferenceObjectTFObjectTypes},
		"value":     types.StringType,
	}
)

// Framework interfaces
var (
	_ resource.Resource                = &TrustFrameworkServiceResource{}
	_ resource.ResourceWithConfigure   = &TrustFrameworkServiceResource{}
	_ resource.ResourceWithImportState = &TrustFrameworkServiceResource{}
)

// New Object
func NewTrustFrameworkServiceResource() resource.Resource {
	return &TrustFrameworkServiceResource{}
}

// Metadata
func (r *TrustFrameworkServiceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_authorize_trust_framework_service"
}

func (r *TrustFrameworkServiceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	// schema descriptions and validation settings
	const attrMinLength = 1
	const attrMinMaxRequestsPerSecond = 0.1
	const attrMinTimeoutMilliseconds = 0
	const attrMaxTimeoutMilliseconds = 3000

	typeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that describes the resource type.",
	).AllowedValuesEnum(authorizeeditor.AllowedEnumAuthorizeEditorDataDefinitionsServiceDefinitionDTOTypeEnumValues)

	serviceTypeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the type of service.",
	).AllowedValuesEnum(authorizeeditor.AllowedEnumAuthorizeEditorDataDefinitionsServiceDefinitionDTOServiceTypeEnumValues)

	processorDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An object that specifies configuration settings for the processor to transform the value returned from the resolver.",
	).AppendMarkdownString(fmt.Sprintf("This field is optional when `service_type` is `%s` or `%s`.", string(authorizeeditor.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOSERVICETYPE_CONNECTOR), string(authorizeeditor.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOSERVICETYPE_HTTP)))

	valueTypeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An object that specifies configuration settings for the final output type of the service.",
	).AppendMarkdownString(fmt.Sprintf("This field is required when `service_type` is `%s` or `%s`.", string(authorizeeditor.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOSERVICETYPE_CONNECTOR), string(authorizeeditor.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOSERVICETYPE_HTTP)))

	serviceSettingsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An object that specifies configuration settings for the service connection.",
	).AppendMarkdownString(fmt.Sprintf("This field is required when `service_type` is `%s` or `%s`.", string(authorizeeditor.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOSERVICETYPE_CONNECTOR), string(authorizeeditor.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOSERVICETYPE_HTTP)))

	serviceSettingsMaximumConcurrentRequestsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An integer that specifies the number of maximum concurrent requests to the service. The value must be greater than or equal to `1`.",
	)

	serviceSettingsMaximumRequestsPerSecondDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A floating point number that specifies the number of maximum requests per second to the service. The value must be greater than `0`.",
	)

	serviceSettingsTimeoutMillisecondsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An integer that specifies the timeout, in milliseconds, when attempting connection to the service. The value must be between `0` and `3000`.",
	)

	serviceSettingsUrlDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the URL of the HTTP service.",
	).AppendMarkdownString(fmt.Sprintf("This field is required when `service_type` is `%s`.", string(authorizeeditor.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOSERVICETYPE_HTTP)))

	serviceSettingsVerbDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the HTTP method to use.",
	).AppendMarkdownString(fmt.Sprintf("This field is required when `service_type` is `%s`.", string(authorizeeditor.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOSERVICETYPE_HTTP))).AllowedValuesEnum(authorizeeditor.AllowedEnumAuthorizeEditorDataServiceSettingsHttpServiceSettingsDTOVerbEnumValues)

	serviceSettingsBodyDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the body of the HTTP request.",
	).AppendMarkdownString(fmt.Sprintf("This field is optional when `service_type` is `%s`.", string(authorizeeditor.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOSERVICETYPE_HTTP)))

	serviceSettingsContentTypeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the content type of the HTTP request.  The service will use the value of this field to set the `Content-Type` header.",
	).AppendMarkdownString(fmt.Sprintf("This field is optional when `service_type` is `%s`.", string(authorizeeditor.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOSERVICETYPE_HTTP)))

	serviceSettingsHeadersDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A set of objects that specify the headers to include in the HTTP request.",
	).AppendMarkdownString(fmt.Sprintf("This field is optional when `service_type` is `%s`.", string(authorizeeditor.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOSERVICETYPE_HTTP)))

	serviceSettingsAuthenticationDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An object that specifies configuration settings for authenticating to the service.",
	).AppendMarkdownString(fmt.Sprintf("This field is required when `service_type` is `%s`.", string(authorizeeditor.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOSERVICETYPE_HTTP)))

	serviceSettingsAuthenticationTypeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the type of service authentication to use.",
	).AllowedValuesEnum(authorizeeditor.AllowedEnumAuthorizeEditorDataAuthenticationDTOTypeEnumValues)

	serviceSettingsAuthenticationNameDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An object that specifies Trust Framework authorization attribute that contains the user name to use for basic authentication.",
	).AppendMarkdownString(fmt.Sprintf("This field is required when `type` is `%s`.", string(authorizeeditor.ENUMAUTHORIZEEDITORDATAAUTHENTICATIONDTOTYPE_BASIC)))

	serviceSettingsAuthenticationPasswordDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An object that specifies the Trust Framework authorization attribute that contains the user password to use for basic authentication.",
	).AppendMarkdownString(fmt.Sprintf("This field is required when `type` is `%s`.", string(authorizeeditor.ENUMAUTHORIZEEDITORDATAAUTHENTICATIONDTOTYPE_BASIC)))

	serviceSettingsAuthenticationTokenEndpointDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the HTTPS token endpoint to use for authentication.",
	).AppendMarkdownString(fmt.Sprintf("This field is required when `type` is `%s`.", string(authorizeeditor.ENUMAUTHORIZEEDITORDATAAUTHENTICATIONDTOTYPE_CLIENT_CREDENTIALS)))

	serviceSettingsAuthenticationClientIdDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the client ID to use for client credentials authentication.",
	).AppendMarkdownString(fmt.Sprintf("This field is required when `type` is `%s`.", string(authorizeeditor.ENUMAUTHORIZEEDITORDATAAUTHENTICATIONDTOTYPE_CLIENT_CREDENTIALS)))

	serviceSettingsAuthenticationClientSecretDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An object that specifies the Trust Framework authorization attribute that contains the client secret to use for client credentials authentication.",
	).AppendMarkdownString(fmt.Sprintf("This field is required when `type` is `%s`.", string(authorizeeditor.ENUMAUTHORIZEEDITORDATAAUTHENTICATIONDTOTYPE_CLIENT_CREDENTIALS)))

	serviceSettingsAuthenticationScopeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the scope(s) to request from the token endpoint during client credentials authentication.",
	).AppendMarkdownString(fmt.Sprintf("This field is required when `type` is `%s`.", string(authorizeeditor.ENUMAUTHORIZEEDITORDATAAUTHENTICATIONDTOTYPE_CLIENT_CREDENTIALS)))

	serviceSettingsAuthenticationTokenDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An object that specifies configuration settings for the token value to use for static token authentication.",
	).AppendMarkdownString(fmt.Sprintf("This field is required when `type` is `%s`.", string(authorizeeditor.ENUMAUTHORIZEEDITORDATAAUTHENTICATIONDTOTYPE_TOKEN)))

	serviceSettingsTlsSettingsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An object that specifies configuration settings when connecting to the service using TLS.",
	).AppendMarkdownString(fmt.Sprintf("This field is required when `service_type` is `%s`.", string(authorizeeditor.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOSERVICETYPE_HTTP)))

	serviceSettingsTlsSettingsTlsValidationTypeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the TLS validation type.",
	).AllowedValuesEnum(authorizeeditor.AllowedEnumAuthorizeEditorDataTlsSettingsDTOTlsValidationTypeEnumValues)

	serviceSettingsChannelDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the connector channel to use for the service.",
	).AppendMarkdownString(fmt.Sprintf("This field is required when `service_type` is `%s`.", string(authorizeeditor.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOSERVICETYPE_CONNECTOR))).AllowedValuesEnum(authorizeeditor.AllowedEnumAuthorizeEditorDataServiceSettingsConnectorServiceSettingsDTOChannelEnumValues)

	serviceSettingsCodeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the connector code to use for the service.",
	).AppendMarkdownString(fmt.Sprintf("This field is required when `service_type` is `%s`.", string(authorizeeditor.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOSERVICETYPE_CONNECTOR))).AllowedValuesEnum(authorizeeditor.AllowedEnumAuthorizeEditorDataServiceSettingsConnectorServiceSettingsDTOCodeEnumValues)

	serviceSettingsCapabilityDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the connector capability associated with the connector code and channel.",
	).AppendMarkdownString(fmt.Sprintf("This field is required when `service_type` is `%s`.", string(authorizeeditor.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOSERVICETYPE_CONNECTOR)))

	serviceSettingsSchemaVersionDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An integer that specifies the schema version of the connector template.",
	).AppendMarkdownString(fmt.Sprintf("This field is optional when `service_type` is `%s`.", string(authorizeeditor.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOSERVICETYPE_CONNECTOR)))

	serviceSettingsInputMappingsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A list of objects that specify configuration settings for the input mappings to use for the service.  Input mappings may be attribute based, or input based.",
	).AppendMarkdownString(fmt.Sprintf("This field is required when `service_type` is `%s`.", string(authorizeeditor.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOSERVICETYPE_CONNECTOR)))

	serviceSettingsInputMappingTypeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the input mapping type.",
	).AllowedValuesEnum(authorizeeditor.AllowedEnumAuthorizeEditorDataInputMappingDTOTypeEnumValues)

	serviceSettingsInputMappingValueRefDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An object that specifies configuration settings for the trust framework attribute to use as an input mapping.",
	).AppendMarkdownString(fmt.Sprintf("This field is required when `type` is `%s`.", string(authorizeeditor.ENUMAUTHORIZEEDITORDATAINPUTMAPPINGDTOTYPE_ATTRIBUTE)))

	serviceSettingsInputMappingValueDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the value for the input mapping.",
	).AppendMarkdownString(fmt.Sprintf("This field is required when `type` is `%s`.", string(authorizeeditor.ENUMAUTHORIZEEDITORDATAINPUTMAPPINGDTOTYPE_INPUT)))

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage an authorization service for the PingOne Authorize Trust Framework in a PingOne environment.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment to configure the Authorize editor service in."),
			),

			"name": schema.StringAttribute{ // DONE
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies a user-friendly service name.").Description,
				Required:    true,

				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
				},
			},

			"full_name": schema.StringAttribute{ // DOC ISSUE
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies a unique name generated by the system for each service resource. It is the concatenation of names in the service resource hierarchy.").Description,
				Computed:    true,
			},

			"description": schema.StringAttribute{ // DONE
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the authorization service resource's description.").Description,
				Optional:    true,
			},

			"parent": parentObjectSchema("service"),

			"type": schema.StringAttribute{ // DOC ISSUE
				Description:         typeDescription.Description,
				MarkdownDescription: typeDescription.MarkdownDescription,
				Computed:            true,

				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"cache_settings": schema.SingleNestedAttribute{ // DONE
				Description: framework.SchemaAttributeDescriptionFromMarkdown("An object that specifies configuration settings for cache settings to apply to the service responses.").Description,
				Optional:    true,

				Attributes: map[string]schema.Attribute{
					"ttl_seconds": schema.Int32Attribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("An integer that specifies the time to live (in seconds) for the service cache.").Description,
						Optional:    true,
					},
				},
			},

			"service_type": schema.StringAttribute{ // DONE
				Description:         serviceTypeDescription.Description,
				MarkdownDescription: serviceTypeDescription.MarkdownDescription,
				Required:            true,

				Validators: []validator.String{
					stringvalidator.OneOf(utils.EnumSliceToStringSlice(authorizeeditor.AllowedEnumAuthorizeEditorDataDefinitionsServiceDefinitionDTOServiceTypeEnumValues)...),
				},
			},

			"version": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that describes a random ID generated by the system for concurrency control purposes.").Description,
				Computed:    true,
			},

			// service_type == "CONNECTOR", service_type == "HTTP"
			"processor": schema.SingleNestedAttribute{
				Description:         processorDescription.Description,
				MarkdownDescription: processorDescription.MarkdownDescription,
				Optional:            true,

				Validators: []validator.Object{
					objectvalidatorinternal.ConflictsIfMatchesPathValue(
						types.StringValue(string(authorizeeditor.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOSERVICETYPE_NONE)),
						path.MatchRoot("service_type"),
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
						types.StringValue(string(authorizeeditor.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOSERVICETYPE_CONNECTOR)),
						path.MatchRoot("service_type"),
					),
					objectvalidatorinternal.IsRequiredIfMatchesPathValue(
						types.StringValue(string(authorizeeditor.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOSERVICETYPE_HTTP)),
						path.MatchRoot("service_type"),
					),
					objectvalidatorinternal.ConflictsIfMatchesPathValue(
						types.StringValue(string(authorizeeditor.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOSERVICETYPE_NONE)),
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
						types.StringValue(string(authorizeeditor.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOSERVICETYPE_CONNECTOR)),
						path.MatchRoot("service_type"),
					),
					objectvalidatorinternal.IsRequiredIfMatchesPathValue(
						types.StringValue(string(authorizeeditor.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOSERVICETYPE_HTTP)),
						path.MatchRoot("service_type"),
					),
					objectvalidatorinternal.ConflictsIfMatchesPathValue(
						types.StringValue(string(authorizeeditor.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOSERVICETYPE_NONE)),
						path.MatchRoot("service_type"),
					),
				},

				Attributes: map[string]schema.Attribute{
					"maximum_concurrent_requests": schema.Int32Attribute{
						Description:         serviceSettingsMaximumConcurrentRequestsDescription.Description,
						MarkdownDescription: serviceSettingsMaximumConcurrentRequestsDescription.MarkdownDescription,
						Required:            true,

						Validators: []validator.Int32{
							int32validator.AtLeast(1),
						},
					},

					"maximum_requests_per_second": schema.Float64Attribute{
						Description:         serviceSettingsMaximumRequestsPerSecondDescription.Description,
						MarkdownDescription: serviceSettingsMaximumRequestsPerSecondDescription.MarkdownDescription,
						Required:            true,

						Validators: []validator.Float64{
							float64validator.AtLeast(attrMinMaxRequestsPerSecond),
						},
					},

					"timeout_milliseconds": schema.Int32Attribute{
						Description:         serviceSettingsTimeoutMillisecondsDescription.Description,
						MarkdownDescription: serviceSettingsTimeoutMillisecondsDescription.MarkdownDescription,
						Required:            true,

						Validators: []validator.Int32{
							int32validator.Between(attrMinTimeoutMilliseconds, attrMaxTimeoutMilliseconds),
						},
					},

					"url": schema.StringAttribute{
						Description:         serviceSettingsUrlDescription.Description,
						MarkdownDescription: serviceSettingsUrlDescription.MarkdownDescription,
						Optional:            true,

						Validators: []validator.String{
							stringvalidatorinternal.IsRequiredIfMatchesPathValue(
								types.StringValue(string(authorizeeditor.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOSERVICETYPE_HTTP)),
								path.MatchRoot("service_type"),
							),
							stringvalidatorinternal.ConflictsIfMatchesPathValue(
								types.StringValue(string(authorizeeditor.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOSERVICETYPE_CONNECTOR)),
								path.MatchRoot("service_type"),
							),
						},
					},

					"verb": schema.StringAttribute{
						Description:         serviceSettingsVerbDescription.Description,
						MarkdownDescription: serviceSettingsVerbDescription.MarkdownDescription,
						Optional:            true,

						Validators: []validator.String{
							stringvalidator.OneOf(utils.EnumSliceToStringSlice(authorizeeditor.AllowedEnumAuthorizeEditorDataServiceSettingsHttpServiceSettingsDTOVerbEnumValues)...),
							stringvalidatorinternal.IsRequiredIfMatchesPathValue(
								types.StringValue(string(authorizeeditor.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOSERVICETYPE_HTTP)),
								path.MatchRoot("service_type"),
							),
							stringvalidatorinternal.ConflictsIfMatchesPathValue(
								types.StringValue(string(authorizeeditor.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOSERVICETYPE_CONNECTOR)),
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
								types.StringValue(string(authorizeeditor.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOSERVICETYPE_CONNECTOR)),
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
								types.StringValue(string(authorizeeditor.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOSERVICETYPE_CONNECTOR)),
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
								types.StringValue(string(authorizeeditor.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOSERVICETYPE_CONNECTOR)),
								path.MatchRoot("service_type"),
							),
						},

						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"key": schema.StringAttribute{
									Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the header key.").Description,
									Required:    true,
								},

								"value": schema.SingleNestedAttribute{
									Description: framework.SchemaAttributeDescriptionFromMarkdown("An object that specifies configuration settings for the header value.  The header value may be configured as an authorization attribute, or a constant value.").Description,
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
								types.StringValue(string(authorizeeditor.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOSERVICETYPE_HTTP)),
								path.MatchRoot("service_type"),
							),
							objectvalidatorinternal.ConflictsIfMatchesPathValue(
								types.StringValue(string(authorizeeditor.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOSERVICETYPE_CONNECTOR)),
								path.MatchRoot("service_type"),
							),
						},

						Attributes: map[string]schema.Attribute{
							"type": schema.StringAttribute{
								Description:         serviceSettingsAuthenticationTypeDescription.Description,
								MarkdownDescription: serviceSettingsAuthenticationTypeDescription.MarkdownDescription,
								Required:            true,

								Validators: []validator.String{
									stringvalidator.OneOf(utils.EnumSliceToStringSlice(authorizeeditor.AllowedEnumAuthorizeEditorDataAuthenticationDTOTypeEnumValues)...),
								},
							},

							// type == "BASIC"
							"name": schema.SingleNestedAttribute{
								Description:         serviceSettingsAuthenticationNameDescription.Description,
								MarkdownDescription: serviceSettingsAuthenticationNameDescription.MarkdownDescription,
								Optional:            true,

								Validators: []validator.Object{
									objectvalidatorinternal.IsRequiredIfMatchesPathValue(
										types.StringValue(string(authorizeeditor.ENUMAUTHORIZEEDITORDATAAUTHENTICATIONDTOTYPE_BASIC)),
										path.MatchRelative().AtParent().AtName("type"),
									),
								},

								Attributes: referenceIdObjectSchemaAttributes(framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the ID of the user name reference used to authenticate the service.")),
							},

							"password": schema.SingleNestedAttribute{
								Description:         serviceSettingsAuthenticationPasswordDescription.Description,
								MarkdownDescription: serviceSettingsAuthenticationPasswordDescription.MarkdownDescription,
								Optional:            true,

								Validators: []validator.Object{
									objectvalidatorinternal.IsRequiredIfMatchesPathValue(
										types.StringValue(string(authorizeeditor.ENUMAUTHORIZEEDITORDATAAUTHENTICATIONDTOTYPE_BASIC)),
										path.MatchRelative().AtParent().AtName("type"),
									),
								},

								Attributes: referenceIdObjectSchemaAttributes(framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the ID of the password reference used to authenticate the service.")),
							},

							// type == "CLIENT_CREDENTIALS"
							"token_endpoint": schema.StringAttribute{
								Description:         serviceSettingsAuthenticationTokenEndpointDescription.Description,
								MarkdownDescription: serviceSettingsAuthenticationTokenEndpointDescription.MarkdownDescription,
								Optional:            true,

								Validators: []validator.String{
									stringvalidatorinternal.IsRequiredIfMatchesPathValue(
										types.StringValue(string(authorizeeditor.ENUMAUTHORIZEEDITORDATAAUTHENTICATIONDTOTYPE_CLIENT_CREDENTIALS)),
										path.MatchRelative().AtParent().AtName("type"),
									),
								},
							},

							"client_id": schema.StringAttribute{
								Description:         serviceSettingsAuthenticationClientIdDescription.Description,
								MarkdownDescription: serviceSettingsAuthenticationClientIdDescription.MarkdownDescription,
								Optional:            true,

								Validators: []validator.String{
									stringvalidatorinternal.IsRequiredIfMatchesPathValue(
										types.StringValue(string(authorizeeditor.ENUMAUTHORIZEEDITORDATAAUTHENTICATIONDTOTYPE_CLIENT_CREDENTIALS)),
										path.MatchRelative().AtParent().AtName("type"),
									),
								},
							},

							"client_secret": schema.SingleNestedAttribute{
								Description:         serviceSettingsAuthenticationClientSecretDescription.Description,
								MarkdownDescription: serviceSettingsAuthenticationClientSecretDescription.MarkdownDescription,
								Optional:            true,

								Validators: []validator.Object{
									objectvalidatorinternal.IsRequiredIfMatchesPathValue(
										types.StringValue(string(authorizeeditor.ENUMAUTHORIZEEDITORDATAAUTHENTICATIONDTOTYPE_CLIENT_CREDENTIALS)),
										path.MatchRelative().AtParent().AtName("type"),
									),
								},

								Attributes: referenceIdObjectSchemaAttributes(framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the ID of the client secret reference used to authenticate the service.")),
							},

							"scope": schema.StringAttribute{
								Description:         serviceSettingsAuthenticationScopeDescription.Description,
								MarkdownDescription: serviceSettingsAuthenticationScopeDescription.MarkdownDescription,
								Optional:            true,

								Validators: []validator.String{
									stringvalidatorinternal.IsRequiredIfMatchesPathValue(
										types.StringValue(string(authorizeeditor.ENUMAUTHORIZEEDITORDATAAUTHENTICATIONDTOTYPE_CLIENT_CREDENTIALS)),
										path.MatchRelative().AtParent().AtName("type"),
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
										types.StringValue(string(authorizeeditor.ENUMAUTHORIZEEDITORDATAAUTHENTICATIONDTOTYPE_TOKEN)),
										path.MatchRelative().AtParent().AtName("type"),
									),
								},

								Attributes: referenceIdObjectSchemaAttributes(framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the ID of the token reference used to authenticate the service.")),
							},
						},
					},

					"tls_settings": schema.SingleNestedAttribute{
						Description:         serviceSettingsTlsSettingsDescription.Description,
						MarkdownDescription: serviceSettingsTlsSettingsDescription.MarkdownDescription,
						Optional:            true,

						Validators: []validator.Object{
							objectvalidatorinternal.IsRequiredIfMatchesPathValue(
								types.StringValue(string(authorizeeditor.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOSERVICETYPE_HTTP)),
								path.MatchRoot("service_type"),
							),
							objectvalidatorinternal.ConflictsIfMatchesPathValue(
								types.StringValue(string(authorizeeditor.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOSERVICETYPE_CONNECTOR)),
								path.MatchRoot("service_type"),
							),
						},

						Attributes: map[string]schema.Attribute{
							"tls_validation_type": schema.StringAttribute{
								Description:         serviceSettingsTlsSettingsTlsValidationTypeDescription.Description,
								MarkdownDescription: serviceSettingsTlsSettingsTlsValidationTypeDescription.MarkdownDescription,
								Required:            true,

								Validators: []validator.String{
									stringvalidator.OneOf(utils.EnumSliceToStringSlice(authorizeeditor.AllowedEnumAuthorizeEditorDataTlsSettingsDTOTlsValidationTypeEnumValues)...),
								},
							},
						},
					},

					"channel": schema.StringAttribute{
						Description:         serviceSettingsChannelDescription.Description,
						MarkdownDescription: serviceSettingsChannelDescription.MarkdownDescription,
						Optional:            true,

						Validators: []validator.String{
							stringvalidator.OneOf(utils.EnumSliceToStringSlice(authorizeeditor.AllowedEnumAuthorizeEditorDataServiceSettingsConnectorServiceSettingsDTOChannelEnumValues)...),
							stringvalidatorinternal.IsRequiredIfMatchesPathValue(
								types.StringValue(string(authorizeeditor.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOSERVICETYPE_CONNECTOR)),
								path.MatchRoot("service_type"),
							),
							stringvalidatorinternal.ConflictsIfMatchesPathValue(
								types.StringValue(string(authorizeeditor.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOSERVICETYPE_HTTP)),
								path.MatchRoot("service_type"),
							),
						},
					},

					"code": schema.StringAttribute{
						Description:         serviceSettingsCodeDescription.Description,
						MarkdownDescription: serviceSettingsCodeDescription.MarkdownDescription,
						Optional:            true,

						Validators: []validator.String{
							stringvalidator.OneOf(utils.EnumSliceToStringSlice(authorizeeditor.AllowedEnumAuthorizeEditorDataServiceSettingsConnectorServiceSettingsDTOCodeEnumValues)...),
							stringvalidatorinternal.IsRequiredIfMatchesPathValue(
								types.StringValue(string(authorizeeditor.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOSERVICETYPE_CONNECTOR)),
								path.MatchRoot("service_type"),
							),
							stringvalidatorinternal.ConflictsIfMatchesPathValue(
								types.StringValue(string(authorizeeditor.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOSERVICETYPE_HTTP)),
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
								types.StringValue(string(authorizeeditor.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOSERVICETYPE_CONNECTOR)),
								path.MatchRoot("service_type"),
							),
							stringvalidatorinternal.ConflictsIfMatchesPathValue(
								types.StringValue(string(authorizeeditor.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOSERVICETYPE_HTTP)),
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
								types.StringValue(string(authorizeeditor.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOSERVICETYPE_HTTP)),
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
								types.StringValue(string(authorizeeditor.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOSERVICETYPE_CONNECTOR)),
								path.MatchRoot("service_type"),
							),
							listvalidatorinternal.ConflictsIfMatchesPathValue(
								types.StringValue(string(authorizeeditor.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOSERVICETYPE_HTTP)),
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
									Description:         serviceSettingsInputMappingTypeDescription.Description,
									MarkdownDescription: serviceSettingsInputMappingTypeDescription.MarkdownDescription,
									Required:            true,

									Validators: []validator.String{
										stringvalidator.OneOf(utils.EnumSliceToStringSlice(authorizeeditor.AllowedEnumAuthorizeEditorDataInputMappingDTOTypeEnumValues)...),
									},
								},

								"value_ref": schema.SingleNestedAttribute{
									Description:         serviceSettingsInputMappingValueRefDescription.Description,
									MarkdownDescription: serviceSettingsInputMappingValueRefDescription.MarkdownDescription,
									Optional:            true,

									Validators: []validator.Object{
										objectvalidatorinternal.IsRequiredIfMatchesPathValue(
											types.StringValue(string(authorizeeditor.ENUMAUTHORIZEEDITORDATAINPUTMAPPINGDTOTYPE_ATTRIBUTE)),
											path.MatchRelative().AtParent().AtName("type"),
										),
										objectvalidatorinternal.ConflictsIfMatchesPathValue(
											types.StringValue(string(authorizeeditor.ENUMAUTHORIZEEDITORDATAINPUTMAPPINGDTOTYPE_INPUT)),
											path.MatchRelative().AtParent().AtName("type"),
										),
									},

									Attributes: referenceIdObjectSchemaAttributes(framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the ID of the authorization attribute in the trust framework.")),
								},

								"value": schema.StringAttribute{
									Description:         serviceSettingsInputMappingValueDescription.Description,
									MarkdownDescription: serviceSettingsInputMappingValueDescription.MarkdownDescription,
									Optional:            true,

									Validators: []validator.String{
										stringvalidatorinternal.IsRequiredIfMatchesPathValue(
											types.StringValue(string(authorizeeditor.ENUMAUTHORIZEEDITORDATAINPUTMAPPINGDTOTYPE_INPUT)),
											path.MatchRelative().AtParent().AtName("type"),
										),
										stringvalidatorinternal.ConflictsIfMatchesPathValue(
											types.StringValue(string(authorizeeditor.ENUMAUTHORIZEEDITORDATAINPUTMAPPINGDTOTYPE_ATTRIBUTE)),
											path.MatchRelative().AtParent().AtName("type"),
										),
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (r *TrustFrameworkServiceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *TrustFrameworkServiceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state trustFrameworkServiceResourceModel

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
	trustFrameworkService, d := plan.expand(ctx, nil)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *authorizeeditor.AuthorizeEditorDataDefinitionsServiceDefinitionDTO
	resp.Diagnostics.Append(legacysdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.BetaAPIClients.AuthorizeEditorAPIClient.AuthorizeEditorServicesApi.CreateService(ctx, plan.EnvironmentId.ValueString()).AuthorizeEditorDataDefinitionsServiceDefinitionDTO(*trustFrameworkService).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"CreateService",
		legacysdk.DefaultCustomError,
		retryAuthorizeEditorCreateUpdate,
		&response,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the state to save
	state = plan

	// Save updated data into Terraform state
	resp.Diagnostics.Append(state.toState(ctx, response)...)
	if !resp.Diagnostics.HasError() {
		resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
	}
}

func (r *TrustFrameworkServiceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *trustFrameworkServiceResourceModel

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
	var response *authorizeeditor.AuthorizeEditorDataDefinitionsServiceDefinitionDTO
	resp.Diagnostics.Append(legacysdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.BetaAPIClients.AuthorizeEditorAPIClient.AuthorizeEditorServicesApi.GetService(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"GetService",
		legacysdk.CustomErrorResourceNotFoundWarning,
		nil,
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
	if !resp.Diagnostics.HasError() {
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	}
}

func (r *TrustFrameworkServiceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state trustFrameworkServiceResourceModel

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

	// Run the API call
	var getResponse *authorizeeditor.AuthorizeEditorDataDefinitionsServiceDefinitionDTO
	resp.Diagnostics.Append(legacysdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.BetaAPIClients.AuthorizeEditorAPIClient.AuthorizeEditorServicesApi.GetService(ctx, plan.EnvironmentId.ValueString(), plan.Id.ValueString()).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"GetService-Update",
		legacysdk.DefaultCustomError,
		retryAuthorizeEditorCreateUpdate,
		&getResponse,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceObj := getResponse.GetActualInstance()

	var version string

	switch t := serviceObj.(type) {
	case *authorizeeditor.AuthorizeEditorDataServicesConnectorServiceDefinitionDTO:
		version = t.GetVersion()
	case *authorizeeditor.AuthorizeEditorDataServicesHttpServiceDefinitionDTO:
		version = t.GetVersion()
	case *authorizeeditor.AuthorizeEditorDataServicesNoneServiceDefinitionDTO:
		version = t.GetVersion()
	default:
		tflog.Error(
			ctx,
			"Service type not supported",
			map[string]interface{}{
				"service type": t,
			},
		)
		resp.Diagnostics.AddError(
			"Service type not supported",
			"The service type is not supported.  Please report this issue to the provider maintainers.",
		)
	}

	if version == "" {
		resp.Diagnostics.AddError(
			"Version not found",
			"Expected the version to be set, got empty.  Please report this issue to the provider maintainers.",
		)
		return
	}

	// Build the model for the API
	trustFrameworkService, d := plan.expand(ctx, &version)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *authorizeeditor.AuthorizeEditorDataDefinitionsServiceDefinitionDTO
	resp.Diagnostics.Append(legacysdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.BetaAPIClients.AuthorizeEditorAPIClient.AuthorizeEditorServicesApi.UpdateService(ctx, plan.EnvironmentId.ValueString(), plan.Id.ValueString()).AuthorizeEditorDataDefinitionsServiceDefinitionDTO(*trustFrameworkService).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdateService",
		legacysdk.DefaultCustomError,
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
	if !resp.Diagnostics.HasError() {
		resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
	}
}

func (r *TrustFrameworkServiceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *trustFrameworkServiceResourceModel

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

	deleteStateConf := &retry.StateChangeConf{
		Pending: []string{
			"200",
		},
		Target: []string{
			"404",
			"ERROR",
		},
		Refresh: func() (interface{}, string, error) {
			// Run the API call
			resp.Diagnostics.Append(legacysdk.ParseResponse(
				ctx,

				func() (any, *http.Response, error) {
					fR, fErr := r.Client.BetaAPIClients.AuthorizeEditorAPIClient.AuthorizeEditorServicesApi.DeleteService(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
					return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), nil, fR, fErr)
				},
				"DeleteService",
				legacysdk.CustomErrorResourceNotFoundWarning,
				nil,
				nil,
			)...)
			if resp.Diagnostics.HasError() {
				return nil, "ERROR", fmt.Errorf("Error deleting authorize service (%s)", data.Id.ValueString())
			}

			fO, fR, fErr := r.Client.BetaAPIClients.AuthorizeEditorAPIClient.AuthorizeEditorServicesApi.GetService(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			getResp, r, err := legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)

			if err != nil || r == nil {
				return getResp, "ERROR", err
			}

			base := 10
			return getResp, strconv.FormatInt(int64(r.StatusCode), base), nil
		},
		Timeout:                   20 * time.Minute,
		Delay:                     1 * time.Second,
		MinTimeout:                500 * time.Millisecond,
		ContinuousTargetOccurence: 2,
	}
	_, err := deleteStateConf.WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddWarning(
			"Authorize Service Delete Timeout",
			fmt.Sprintf("Error waiting for authorize service (%s) to be deleted: %s", data.Id.ValueString(), err),
		)

		return
	}
}

func (r *TrustFrameworkServiceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	idComponents := []framework.ImportComponent{
		{
			Label:  "environment_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
		{
			Label:     "authorize_trust_framework_service_id",
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

func (p *trustFrameworkServiceResourceModel) expand(ctx context.Context, updateVersionId *string) (*authorizeeditor.AuthorizeEditorDataDefinitionsServiceDefinitionDTO, diag.Diagnostics) {
	var diags, d diag.Diagnostics

	data := authorizeeditor.AuthorizeEditorDataDefinitionsServiceDefinitionDTO{}

	commonData, d := p.expandCommon(ctx, updateVersionId)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	switch authorizeeditor.EnumAuthorizeEditorDataDefinitionsServiceDefinitionDTOServiceType(p.ServiceType.ValueString()) {
	case authorizeeditor.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOSERVICETYPE_CONNECTOR:
		data.AuthorizeEditorDataServicesConnectorServiceDefinitionDTO, d = p.expandConnectorService(ctx, commonData)
		diags.Append(d...)
	case authorizeeditor.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOSERVICETYPE_HTTP:
		data.AuthorizeEditorDataServicesHttpServiceDefinitionDTO, d = p.expandHttpService(ctx, commonData)
		diags.Append(d...)
	case authorizeeditor.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOSERVICETYPE_NONE:
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

func (p *trustFrameworkServiceResourceModel) expandCommon(ctx context.Context, updateVersionId *string) (*authorizeeditor.AuthorizeEditorDataDefinitionsServiceDefinitionDTOCommon, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := authorizeeditor.NewAuthorizeEditorDataDefinitionsServiceDefinitionDTOCommon(
		p.Name.ValueString(),
		authorizeeditor.EnumAuthorizeEditorDataDefinitionsServiceDefinitionDTOServiceType(p.ServiceType.ValueString()),
	)

	if !p.Description.IsNull() && !p.Description.IsUnknown() {
		data.SetDescription(p.Description.ValueString())
	}

	if !p.Parent.IsNull() && !p.Parent.IsUnknown() {
		parent, d := expandEditorParent(ctx, p.Parent)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		data.SetParent(*parent)
	}

	if !p.CacheSettings.IsNull() && !p.CacheSettings.IsUnknown() {
		var plan *trustFrameworkServiceCacheSettingsResourceModel
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

	if updateVersionId != nil {
		data.SetVersion(*updateVersionId)

		if !p.Id.IsNull() && !p.Id.IsUnknown() {
			data.SetId(p.Id.ValueString())
		}
	}

	return data, diags
}

func (p *trustFrameworkServiceCacheSettingsResourceModel) expand() *authorizeeditor.AuthorizeEditorDataCacheSettingsDTO {

	data := authorizeeditor.NewAuthorizeEditorDataCacheSettingsDTO()

	if !p.TtlSeconds.IsNull() && !p.TtlSeconds.IsUnknown() {
		data.SetTtlSeconds(p.TtlSeconds.ValueInt32())
	}

	return data
}

func (p *trustFrameworkServiceResourceModel) expandConnectorService(ctx context.Context, commonData *authorizeeditor.AuthorizeEditorDataDefinitionsServiceDefinitionDTOCommon) (*authorizeeditor.AuthorizeEditorDataServicesConnectorServiceDefinitionDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	if p.ServiceType.ValueString() != string(authorizeeditor.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOSERVICETYPE_CONNECTOR) {
		return nil, diags
	}

	// Use json.marshall and unmarshal to cast commonData to a AuthorizeEditorDataServicesConnectorServiceDefinitionDTO type
	bytes, err := json.Marshal(commonData)
	if err != nil {
		diags.AddError("Failed to marshal data", err.Error())
		return nil, diags
	}

	var data *authorizeeditor.AuthorizeEditorDataServicesConnectorServiceDefinitionDTO
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

	var serviceSettingsPlan *trustFrameworkServiceServiceSettingsResourceModel
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

func (p *trustFrameworkServiceResourceModel) expandHttpService(ctx context.Context, commonData *authorizeeditor.AuthorizeEditorDataDefinitionsServiceDefinitionDTOCommon) (*authorizeeditor.AuthorizeEditorDataServicesHttpServiceDefinitionDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	if p.ServiceType.ValueString() != string(authorizeeditor.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOSERVICETYPE_HTTP) {
		return nil, diags
	}

	// Use json.marshall and unmarshal to cast commonData to a AuthorizeEditorDataServicesHttpServiceDefinitionDTO type
	bytes, err := json.Marshal(commonData)
	if err != nil {
		diags.AddError("Failed to marshal data", err.Error())
		return nil, diags
	}

	var data *authorizeeditor.AuthorizeEditorDataServicesHttpServiceDefinitionDTO
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

	var serviceSettingsPlan *trustFrameworkServiceServiceSettingsResourceModel
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

func (p *trustFrameworkServiceResourceModel) expandNoneService(commonData *authorizeeditor.AuthorizeEditorDataDefinitionsServiceDefinitionDTOCommon) (*authorizeeditor.AuthorizeEditorDataServicesNoneServiceDefinitionDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	if p.ServiceType.ValueString() != string(authorizeeditor.ENUMAUTHORIZEEDITORDATADEFINITIONSSERVICEDEFINITIONDTOSERVICETYPE_NONE) {
		return nil, diags
	}

	// Use json.marshall and unmarshal to cast commonData to a AuthorizeEditorDataServicesNoneServiceDefinitionDTO type
	bytes, err := json.Marshal(commonData)
	if err != nil {
		diags.AddError("Failed to marshal data", err.Error())
		return nil, diags
	}

	var data *authorizeeditor.AuthorizeEditorDataServicesNoneServiceDefinitionDTO
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		diags.AddError("Failed to unmarshal data", err.Error())
		return nil, diags
	}

	return data, diags
}

func (p *trustFrameworkServiceServiceSettingsResourceModel) expandConnector(ctx context.Context) (*authorizeeditor.AuthorizeEditorDataServiceSettingsConnectorServiceSettingsDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	inputMappings := make([]authorizeeditor.AuthorizeEditorDataInputMappingDTO, 0)

	var inputMappingsPlan []trustFrameworkServiceServiceSettingsInputMappingResourceModel
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

	data := authorizeeditor.NewAuthorizeEditorDataServiceSettingsConnectorServiceSettingsDTO(
		authorizeeditor.EnumAuthorizeEditorDataServiceSettingsConnectorServiceSettingsDTOChannel(p.Channel.ValueString()),
		authorizeeditor.EnumAuthorizeEditorDataServiceSettingsConnectorServiceSettingsDTOCode(p.Code.ValueString()),
		p.Capability.ValueString(),
		inputMappings,
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

	if !p.SchemaVersion.IsNull() && !p.SchemaVersion.IsUnknown() {
		data.SetSchemaVersion(p.SchemaVersion.ValueInt32())
	}

	return data, diags
}

func (p *trustFrameworkServiceServiceSettingsInputMappingResourceModel) expand(ctx context.Context) (*authorizeeditor.AuthorizeEditorDataInputMappingDTO, diag.Diagnostics) {
	var diags, d diag.Diagnostics

	data := authorizeeditor.AuthorizeEditorDataInputMappingDTO{}

	switch authorizeeditor.EnumAuthorizeEditorDataInputMappingDTOType(p.Type.ValueString()) {
	case authorizeeditor.ENUMAUTHORIZEEDITORDATAINPUTMAPPINGDTOTYPE_ATTRIBUTE:
		data.AuthorizeEditorDataInputMappingsAttributeInputMappingDTO, d = p.expandAttributeType(ctx)
		diags.Append(d...)
	case authorizeeditor.ENUMAUTHORIZEEDITORDATAINPUTMAPPINGDTOTYPE_INPUT:
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

func (p *trustFrameworkServiceServiceSettingsInputMappingResourceModel) expandAttributeType(ctx context.Context) (*authorizeeditor.AuthorizeEditorDataInputMappingsAttributeInputMappingDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	valueRef, d := expandEditorReferenceData(ctx, p.ValueRef)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	data := authorizeeditor.NewAuthorizeEditorDataInputMappingsAttributeInputMappingDTO(
		p.Property.ValueString(),
		authorizeeditor.EnumAuthorizeEditorDataInputMappingDTOType(p.Type.ValueString()),
		*valueRef,
	)

	return data, diags
}

func (p *trustFrameworkServiceServiceSettingsInputMappingResourceModel) expandInputType() *authorizeeditor.AuthorizeEditorDataInputMappingsInputInputMappingDTO {

	data := authorizeeditor.NewAuthorizeEditorDataInputMappingsInputInputMappingDTO(
		p.Property.ValueString(),
		authorizeeditor.EnumAuthorizeEditorDataInputMappingDTOType(p.Type.ValueString()),
		p.Value.ValueString(),
	)

	return data
}

func (p *trustFrameworkServiceServiceSettingsResourceModel) expandHttp(ctx context.Context) (*authorizeeditor.AuthorizeEditorDataServiceSettingsHttpServiceSettingsDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	var authenticationPlan *trustFrameworkServiceServiceSettingsAuthenticationResourceModel
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

	var tlsSettingsPlan *trustFrameworkServiceServiceSettingsTlsSettingsResourceModel
	diags.Append(p.TlsSettings.As(ctx, &tlsSettingsPlan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)
	if diags.HasError() {
		return nil, diags
	}

	tlsSettings := tlsSettingsPlan.expand()

	data := authorizeeditor.NewAuthorizeEditorDataServiceSettingsHttpServiceSettingsDTO(
		p.Url.ValueString(),
		authorizeeditor.EnumAuthorizeEditorDataServiceSettingsHttpServiceSettingsDTOVerb(p.Verb.ValueString()),
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
		headers := make([]authorizeeditor.AuthorizeEditorDataHttpRequestHeaderDTO, 0)

		var plan []trustFrameworkServiceServiceSettingsHeaderResourceModel
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

func (p *trustFrameworkServiceServiceSettingsHeaderResourceModel) expand(ctx context.Context) (*authorizeeditor.AuthorizeEditorDataHttpRequestHeaderDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := authorizeeditor.NewAuthorizeEditorDataHttpRequestHeaderDTO(
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

func (p *trustFrameworkServiceServiceSettingsAuthenticationResourceModel) expand(ctx context.Context) (*authorizeeditor.AuthorizeEditorDataAuthenticationDTO, diag.Diagnostics) {
	var diags, d diag.Diagnostics

	data := authorizeeditor.AuthorizeEditorDataAuthenticationDTO{}

	switch authorizeeditor.EnumAuthorizeEditorDataAuthenticationDTOType(p.Type.ValueString()) {
	case authorizeeditor.ENUMAUTHORIZEEDITORDATAAUTHENTICATIONDTOTYPE_BASIC:
		data.AuthorizeEditorDataAuthenticationsBasicAuthenticationDTO, d = p.expandBasicAuth(ctx)
		diags.Append(d...)
	case authorizeeditor.ENUMAUTHORIZEEDITORDATAAUTHENTICATIONDTOTYPE_CLIENT_CREDENTIALS:
		data.AuthorizeEditorDataAuthenticationsClientCredentialsAuthenticationDTO, d = p.expandClientCredentialsAuth(ctx)
		diags.Append(d...)
	case authorizeeditor.ENUMAUTHORIZEEDITORDATAAUTHENTICATIONDTOTYPE_NONE:
		data.AuthorizeEditorDataAuthenticationsNoneAuthenticationDTO = p.expandNoneAuth()
	case authorizeeditor.ENUMAUTHORIZEEDITORDATAAUTHENTICATIONDTOTYPE_TOKEN:
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

func (p *trustFrameworkServiceServiceSettingsAuthenticationResourceModel) expandBasicAuth(ctx context.Context) (*authorizeeditor.AuthorizeEditorDataAuthenticationsBasicAuthenticationDTO, diag.Diagnostics) {
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

	data := authorizeeditor.NewAuthorizeEditorDataAuthenticationsBasicAuthenticationDTO(
		authorizeeditor.EnumAuthorizeEditorDataAuthenticationDTOType(p.Type.ValueString()),
		*name,
		*password,
	)

	return data, diags
}

func (p *trustFrameworkServiceServiceSettingsAuthenticationResourceModel) expandClientCredentialsAuth(ctx context.Context) (*authorizeeditor.AuthorizeEditorDataAuthenticationsClientCredentialsAuthenticationDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	clientSecret, d := expandEditorReferenceData(ctx, p.ClientSecret)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	data := authorizeeditor.NewAuthorizeEditorDataAuthenticationsClientCredentialsAuthenticationDTO(
		authorizeeditor.EnumAuthorizeEditorDataAuthenticationDTOType(p.Type.ValueString()),
		p.TokenEndpoint.ValueString(),
		p.ClientId.ValueString(),
		*clientSecret,
		p.Scope.ValueString(),
	)

	return data, diags
}

func (p *trustFrameworkServiceServiceSettingsAuthenticationResourceModel) expandNoneAuth() *authorizeeditor.AuthorizeEditorDataAuthenticationsNoneAuthenticationDTO {

	data := authorizeeditor.NewAuthorizeEditorDataAuthenticationsNoneAuthenticationDTO(
		authorizeeditor.EnumAuthorizeEditorDataAuthenticationDTOType(p.Type.ValueString()),
	)

	return data
}

func (p *trustFrameworkServiceServiceSettingsAuthenticationResourceModel) expandTokenAuth(ctx context.Context) (*authorizeeditor.AuthorizeEditorDataAuthenticationsTokenAuthenticationDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	token, d := expandEditorReferenceData(ctx, p.Token)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	data := authorizeeditor.NewAuthorizeEditorDataAuthenticationsTokenAuthenticationDTO(
		authorizeeditor.EnumAuthorizeEditorDataAuthenticationDTOType(p.Type.ValueString()),
		*token,
	)

	return data, diags
}

func (p *trustFrameworkServiceServiceSettingsTlsSettingsResourceModel) expand() *authorizeeditor.AuthorizeEditorDataTlsSettingsDTO {

	data := authorizeeditor.NewAuthorizeEditorDataTlsSettingsDTO(
		authorizeeditor.EnumAuthorizeEditorDataTlsSettingsDTOTlsValidationType(p.TlsValidationType.ValueString()),
	)

	return data
}

func (p *trustFrameworkServiceResourceModel) toState(ctx context.Context, apiObject *authorizeeditor.AuthorizeEditorDataDefinitionsServiceDefinitionDTO) diag.Diagnostics {
	var diags, d diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)
		return diags
	}

	apiObjectCommon := authorizeeditor.AuthorizeEditorDataDefinitionsServiceDefinitionDTOCommon{}

	if apiObject.AuthorizeEditorDataServicesConnectorServiceDefinitionDTO != nil {
		apiObjectCommon = authorizeeditor.AuthorizeEditorDataDefinitionsServiceDefinitionDTOCommon{
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
		apiObjectCommon = authorizeeditor.AuthorizeEditorDataDefinitionsServiceDefinitionDTOCommon{
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
		apiObjectCommon = authorizeeditor.AuthorizeEditorDataDefinitionsServiceDefinitionDTOCommon{
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

	p.CacheSettings, d = trustFrameworkServiceCacheSettingsOkToTF(apiObjectCommon.GetCacheSettingsOk())
	diags.Append(d...)

	p.ServiceType = framework.EnumOkToTF(apiObjectCommon.GetServiceTypeOk())
	p.Version = framework.StringOkToTF(apiObjectCommon.GetVersionOk())

	p.Processor = types.ObjectNull(editorDataProcessorTFObjectTypes)
	p.ValueType = types.ObjectNull(trustFrameworkServiceValueTypeTFObjectTypes)
	p.ServiceSettings = types.ObjectNull(trustFrameworkServiceServiceSettingsTFObjectTypes)

	if v := apiObject.AuthorizeEditorDataServicesConnectorServiceDefinitionDTO; v != nil {
		processorVal, ok := v.GetProcessorOk()
		p.Processor, d = editorDataProcessorOkToTF(ctx, processorVal, ok)
		diags.Append(d...)

		p.ValueType, d = editorValueTypeOkToTF(v.GetValueTypeOk())
		diags.Append(d...)

		serviceSettingsVal, ok := v.GetServiceSettingsOk()
		p.ServiceSettings, d = trustFrameworkServiceServiceSettingsConnectorOkToTF(ctx, serviceSettingsVal, ok)
		diags.Append(d...)
	}

	if v := apiObject.AuthorizeEditorDataServicesHttpServiceDefinitionDTO; v != nil {
		processorVal, ok := v.GetProcessorOk()
		p.Processor, d = editorDataProcessorOkToTF(ctx, processorVal, ok)
		diags.Append(d...)

		p.ValueType, d = editorValueTypeOkToTF(v.GetValueTypeOk())
		diags.Append(d...)

		serviceSettingsVal, ok := v.GetServiceSettingsOk()
		p.ServiceSettings, d = trustFrameworkServiceServiceSettingsHttpOkToTF(ctx, serviceSettingsVal, ok)
		diags.Append(d...)
	}

	// No implementation for "None" service

	return diags
}

func trustFrameworkServiceCacheSettingsOkToTF(apiObject *authorizeeditor.AuthorizeEditorDataCacheSettingsDTO, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(trustFrameworkServiceCacheSettingsTFObjectTypes), diags
	}

	objValue, d := types.ObjectValue(trustFrameworkServiceCacheSettingsTFObjectTypes, map[string]attr.Value{
		"ttl_seconds": framework.Int32OkToTF(apiObject.GetTtlSecondsOk()),
	})
	diags.Append(d...)

	return objValue, diags
}

func trustFrameworkServiceServiceSettingsConnectorOkToTF(ctx context.Context, apiObject *authorizeeditor.AuthorizeEditorDataServiceSettingsConnectorServiceSettingsDTO, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(trustFrameworkServiceServiceSettingsTFObjectTypes), diags
	}

	inputMappingsVal, ok := apiObject.GetInputMappingsOk()
	inputMappings, d := trustFrameworkServiceServiceSettingsConnectorInputMappingsOkToTF(ctx, inputMappingsVal, ok)
	diags.Append(d...)

	attributeMap := map[string]attr.Value{
		"maximum_concurrent_requests": framework.Int32OkToTF(apiObject.GetMaximumConcurrentRequestsOk()),
		"maximum_requests_per_second": framework.Float64OkToTF(apiObject.GetMaximumRequestsPerSecondOk()),
		"timeout_milliseconds":        framework.Int32OkToTF(apiObject.GetTimeoutMillisecondsOk()),
		"channel":                     framework.EnumOkToTF(apiObject.GetChannelOk()),
		"code":                        framework.EnumOkToTF(apiObject.GetCodeOk()),
		"capability":                  framework.StringOkToTF(apiObject.GetCapabilityOk()),
		"schema_version":              framework.Int32OkToTF(apiObject.GetSchemaVersionOk()),
		"input_mappings":              inputMappings,
	}

	attributeMap = trustFrameworkServiceServiceSettingsConvertEmptyValuesToTFNulls(attributeMap)

	objValue, d := types.ObjectValue(trustFrameworkServiceServiceSettingsTFObjectTypes, attributeMap)
	diags.Append(d...)

	return objValue, diags
}

func trustFrameworkServiceServiceSettingsConnectorInputMappingsOkToTF(ctx context.Context, apiObject []authorizeeditor.AuthorizeEditorDataInputMappingDTO, ok bool) (basetypes.ListValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	tfObjType := types.ObjectType{AttrTypes: trustFrameworkServiceServiceSettingsInputMappingsTFObjectTypes}

	if !ok || apiObject == nil {
		return types.ListNull(tfObjType), diags
	}

	flattenedList := []attr.Value{}
	for _, v := range apiObject {

		flattenedObj, d := trustFrameworkServiceServiceSettingsConnectorInputMappingOkToTF(ctx, &v, true)
		diags.Append(d...)

		flattenedList = append(flattenedList, flattenedObj)
	}

	returnVar, d := types.ListValue(tfObjType, flattenedList)
	diags.Append(d...)

	return returnVar, diags
}

func trustFrameworkServiceServiceSettingsConnectorInputMappingOkToTF(ctx context.Context, apiObject *authorizeeditor.AuthorizeEditorDataInputMappingDTO, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil || cmp.Equal(apiObject, &authorizeeditor.AuthorizeEditorDataInputMappingDTO{}) {
		return types.ObjectNull(trustFrameworkServiceServiceSettingsInputMappingsTFObjectTypes), diags
	}

	attributeMap := map[string]attr.Value{}

	switch t := apiObject.GetActualInstance().(type) {
	case *authorizeeditor.AuthorizeEditorDataInputMappingsAttributeInputMappingDTO:

		valueResp, ok := t.GetValueOk()
		value, d := editorDataReferenceObjectOkToTF(valueResp, ok)
		diags.Append(d...)

		attributeMap["type"] = framework.EnumOkToTF(t.GetTypeOk())
		attributeMap["property"] = framework.StringOkToTF(t.GetPropertyOk())
		attributeMap["value_ref"] = value

	case *authorizeeditor.AuthorizeEditorDataInputMappingsInputInputMappingDTO:

		attributeMap["type"] = framework.EnumOkToTF(t.GetTypeOk())
		attributeMap["property"] = framework.StringOkToTF(t.GetPropertyOk())
		attributeMap["value"] = framework.StringOkToTF(t.GetValueOk())

	default:
		tflog.Error(ctx, "Invalid service setting connector input mapping type", map[string]interface{}{
			"service setting connector input mapping type": t,
		})
		diags.AddError(
			"Invalid service setting connector input mapping type",
			"The service setting connector input mapping type is not supported.  Please raise an issue with the provider maintainers.",
		)
		return types.ObjectNull(trustFrameworkServiceServiceSettingsInputMappingsTFObjectTypes), diags
	}

	attributeMap = trustFrameworkServiceServiceSettingsConnectorInputMappingConvertEmptyValuesToTFNulls(attributeMap)

	objValue, d := types.ObjectValue(trustFrameworkServiceServiceSettingsInputMappingsTFObjectTypes, attributeMap)
	diags.Append(d...)

	return objValue, diags
}

func trustFrameworkServiceServiceSettingsConnectorInputMappingConvertEmptyValuesToTFNulls(attributeMap map[string]attr.Value) map[string]attr.Value {
	nullMap := map[string]attr.Value{
		"type":      types.StringNull(),
		"property":  types.StringNull(),
		"value_ref": types.ObjectNull(editorReferenceObjectTFObjectTypes),
		"value":     types.StringNull(),
	}

	for k := range nullMap {
		if attributeMap[k] == nil {
			attributeMap[k] = nullMap[k]
		}
	}

	return attributeMap
}

func trustFrameworkServiceServiceSettingsHttpOkToTF(ctx context.Context, apiObject *authorizeeditor.AuthorizeEditorDataServiceSettingsHttpServiceSettingsDTO, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(trustFrameworkServiceServiceSettingsTFObjectTypes), diags
	}

	headersVal, ok := apiObject.GetHeadersOk()
	headers, d := trustFrameworkServiceServiceSettingsHttpHeadersOkToTF(ctx, headersVal, ok)
	diags.Append(d...)

	authenticationVal, ok := apiObject.GetAuthenticationOk()
	authentication, d := trustFrameworkServiceServiceSettingsHttpAuthenticationOkToTF(ctx, authenticationVal, ok)
	diags.Append(d...)

	tlsSettings, d := trustFrameworkServiceServiceSettingsHttpTlsSettingsOkToTF(apiObject.GetTlsSettingsOk())
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

	attributeMap = trustFrameworkServiceServiceSettingsConvertEmptyValuesToTFNulls(attributeMap)

	objValue, d := types.ObjectValue(trustFrameworkServiceServiceSettingsTFObjectTypes, attributeMap)
	diags.Append(d...)

	return objValue, diags
}

func trustFrameworkServiceServiceSettingsConvertEmptyValuesToTFNulls(attributeMap map[string]attr.Value) map[string]attr.Value {
	nullMap := map[string]attr.Value{
		"maximum_concurrent_requests": types.Int32Null(),
		"maximum_requests_per_second": types.Float64Null(),
		"timeout_milliseconds":        types.Int32Null(),
		"url":                         types.StringNull(),
		"verb":                        types.StringNull(),
		"body":                        types.StringNull(),
		"content_type":                types.StringNull(),
		"headers":                     types.SetNull(types.ObjectType{AttrTypes: trustFrameworkServiceServiceSettingsHeadersTFObjectTypes}),
		"authentication":              types.ObjectNull(trustFrameworkServiceServiceSettingsAuthenticationTFObjectTypes),
		"tls_settings":                types.ObjectNull(trustFrameworkServiceServiceSettingsTlsSettingsTFObjectTypes),
		"channel":                     types.StringNull(),
		"code":                        types.StringNull(),
		"capability":                  types.StringNull(),
		"schema_version":              types.Int32Null(),
		"input_mappings":              types.ListNull(types.ObjectType{AttrTypes: trustFrameworkServiceServiceSettingsInputMappingsTFObjectTypes}),
	}

	for k := range nullMap {
		if attributeMap[k] == nil {
			attributeMap[k] = nullMap[k]
		}
	}

	return attributeMap
}

func trustFrameworkServiceServiceSettingsHttpHeadersOkToTF(ctx context.Context, apiObject []authorizeeditor.AuthorizeEditorDataHttpRequestHeaderDTO, ok bool) (basetypes.SetValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	tfObjType := types.ObjectType{AttrTypes: trustFrameworkServiceServiceSettingsHeadersTFObjectTypes}

	if !ok || apiObject == nil {
		return types.SetNull(tfObjType), diags
	}

	flattenedList := []attr.Value{}
	for _, v := range apiObject {

		dataInputVal, ok := v.GetValueOk()
		value, d := editorDataInputOkToTF(ctx, dataInputVal, ok)
		diags.Append(d...)

		flattenedObj, d := types.ObjectValue(trustFrameworkServiceServiceSettingsHeadersTFObjectTypes, map[string]attr.Value{
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

func trustFrameworkServiceServiceSettingsHttpAuthenticationOkToTF(ctx context.Context, apiObject *authorizeeditor.AuthorizeEditorDataAuthenticationDTO, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil || cmp.Equal(apiObject, &authorizeeditor.AuthorizeEditorDataAuthenticationDTO{}) {
		return types.ObjectNull(trustFrameworkServiceServiceSettingsAuthenticationTFObjectTypes), diags
	}

	attributeMap := map[string]attr.Value{}

	switch t := apiObject.GetActualInstance().(type) {
	case *authorizeeditor.AuthorizeEditorDataAuthenticationsBasicAuthenticationDTO:

		nameResp, ok := t.GetNameOk()
		name, d := editorDataReferenceObjectOkToTF(nameResp, ok)
		diags.Append(d...)

		passwordResp, ok := t.GetPasswordOk()
		password, d := editorDataReferenceObjectOkToTF(passwordResp, ok)
		diags.Append(d...)

		attributeMap["type"] = framework.EnumOkToTF(t.GetTypeOk())
		attributeMap["name"] = name
		attributeMap["password"] = password

	case *authorizeeditor.AuthorizeEditorDataAuthenticationsClientCredentialsAuthenticationDTO:

		clientSecretResp, ok := t.GetClientSecretOk()
		clientSecret, d := editorDataReferenceObjectOkToTF(clientSecretResp, ok)
		diags.Append(d...)

		attributeMap["type"] = framework.EnumOkToTF(t.GetTypeOk())
		attributeMap["token_endpoint"] = framework.StringOkToTF(t.GetTokenEndpointOk())
		attributeMap["client_id"] = framework.StringOkToTF(t.GetClientIdOk())
		attributeMap["client_secret"] = clientSecret
		attributeMap["scope"] = framework.StringOkToTF(t.GetScopeOk())

	case *authorizeeditor.AuthorizeEditorDataAuthenticationsNoneAuthenticationDTO:

		attributeMap["type"] = framework.EnumOkToTF(t.GetTypeOk())

	case *authorizeeditor.AuthorizeEditorDataAuthenticationsTokenAuthenticationDTO:

		tokenResp, ok := t.GetTokenOk()
		token, d := editorDataReferenceObjectOkToTF(tokenResp, ok)
		diags.Append(d...)

		attributeMap["type"] = framework.EnumOkToTF(t.GetTypeOk())
		attributeMap["token"] = token

	default:
		tflog.Error(ctx, "Invalid service settings authentication type", map[string]interface{}{
			"service settings authentication type": t,
		})
		diags.AddError(
			"Invalid service settings authentication type",
			"The service settings authentication type is not supported.  Please raise an issue with the provider maintainers.",
		)
		return types.ObjectNull(trustFrameworkServiceServiceSettingsAuthenticationTFObjectTypes), diags
	}

	attributeMap = trustFrameworkServiceServiceSettingsAuthenticationConvertEmptyValuesToTFNulls(attributeMap)

	objValue, d := types.ObjectValue(trustFrameworkServiceServiceSettingsAuthenticationTFObjectTypes, attributeMap)
	diags.Append(d...)

	return objValue, diags
}

func trustFrameworkServiceServiceSettingsAuthenticationConvertEmptyValuesToTFNulls(attributeMap map[string]attr.Value) map[string]attr.Value {
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

func trustFrameworkServiceServiceSettingsHttpTlsSettingsOkToTF(apiObject *authorizeeditor.AuthorizeEditorDataTlsSettingsDTO, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(trustFrameworkServiceServiceSettingsTlsSettingsTFObjectTypes), diags
	}

	objValue, d := types.ObjectValue(trustFrameworkServiceServiceSettingsTlsSettingsTFObjectTypes, map[string]attr.Value{
		"tls_validation_type": framework.EnumOkToTF(apiObject.GetTlsValidationTypeOk()),
	})
	diags.Append(d...)

	return objValue, diags
}
