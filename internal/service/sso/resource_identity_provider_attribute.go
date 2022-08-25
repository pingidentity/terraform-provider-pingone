package sso

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func ResourceIdentityProviderAttribute() *schema.Resource {
	return &schema.Resource{

		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage an attribute mapping for identity providers configured in PingOne.",

		CreateContext: resourcePingOneIdentityProviderAttributeCreate,
		ReadContext:   resourcePingOneIdentityProviderAttributeRead,
		UpdateContext: resourcePingOneIdentityProviderAttributeUpdate,
		DeleteContext: resourcePingOneIdentityProviderAttributeDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourcePingOneIdentityProviderAttributeImport,
		},

		Schema: map[string]*schema.Schema{
			"environment_id": {
				Description:      "The ID of the environment to create the identity provider attribute in.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
				ForceNew:         true,
			},
			"identity_provider_id": {
				Description:      "The ID of the identity provider to create the attribute mapping for.",
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
			},
			"name": {
				Description:      "The user attribute, which is unique per provider. The attribute must not be defined as read only from the user schema or of type `COMPLEX` based on the user schema. Valid examples `username`, and `name.first`. The following attributes may not be used `account`, `id`, `created`, `updated`, `lifecycle`, `mfaEnabled`, and `enabled`.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringNotInSlice([]string{"account", "id", "created", "updated", "lifecycle", "mfaEnabled", "enabled"}, false)),
			},
			"update": {
				Description:      fmt.Sprintf("Indicates whether to update the user attribute in the directory with the non-empty mapped value from the IdP. Options are `%s` (only update the user attribute if it has an empty value); `%s` (always update the user attribute value).", string(management.ENUMIDENTITYPROVIDERATTRIBUTEMAPPINGUPDATE_EMPTY_ONLY), string(management.ENUMIDENTITYPROVIDERATTRIBUTEMAPPINGUPDATE_ALWAYS)),
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{string(management.ENUMIDENTITYPROVIDERATTRIBUTEMAPPINGUPDATE_EMPTY_ONLY), string(management.ENUMIDENTITYPROVIDERATTRIBUTEMAPPINGUPDATE_ALWAYS)}, false)),
			},
			"value": {
				Description:      "A placeholder referring to the attribute (or attributes) from the provider. Placeholders must be valid for the attributes returned by the IdP type and use the ${} syntax (for example, `${email}`). For SAML, any placeholder is acceptable, and it is mapped against the attributes available in the SAML assertion after authentication. The `${samlAssertion.subject}` placeholder is a special reserved placeholder used to refer to the subject name ID in the SAML assertion response.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
			},
			"mapping_type": {
				Description: fmt.Sprintf("The mapping type. Options are `%s` (This attribute is required by the schema and cannot be removed. The name and update properties cannot be changed.) or `%s` (All user-created attributes are of this type.)", string(management.ENUMIDENTITYPROVIDERATTRIBUTEMAPPINGTYPE_CORE), string(management.ENUMIDENTITYPROVIDERATTRIBUTEMAPPINGTYPE_CUSTOM)),
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func resourcePingOneIdentityProviderAttributeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	idpAttributeMapping := *management.NewIdentityProviderAttribute(d.Get("name").(string), d.Get("value").(string), management.EnumIdentityProviderAttributeMappingType(d.Get("update").(string)))

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.IdentityProviderManagementIdentityProviderAttributesApi.CreateIdentityProviderAttribute(ctx, d.Get("environment_id").(string), d.Get("identity_provider_id").(string)).IdentityProviderAttribute(idpAttributeMapping).Execute()
		},
		"CreateIdentityProviderAttribute",
		sdk.CustomErrorInvalidValue,
		sdk.DefaultCreateReadRetryable,
	)
	if diags.HasError() {
		return diags
	}

	respObject := resp.(*management.IdentityProviderAttribute)

	d.SetId(respObject.GetId())

	return resourcePingOneIdentityProviderAttributeRead(ctx, d, meta)
}

func resourcePingOneIdentityProviderAttributeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.IdentityProviderManagementIdentityProviderAttributesApi.ReadOneIdentityProviderAttribute(ctx, d.Get("environment_id").(string), d.Get("identity_provider_id").(string), d.Id()).Execute()
		},
		"ReadOneIdentityProviderAttribute",
		sdk.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
	)
	if diags.HasError() {
		return diags
	}

	if resp == nil {
		d.SetId("")
		return nil
	}

	respObject := resp.(*management.IdentityProviderAttribute)

	d.Set("name", respObject.GetName())
	d.Set("value", respObject.GetValue())
	d.Set("update", respObject.GetUpdate())
	d.Set("mapping_type", respObject.GetMappingType())

	return diags
}

func resourcePingOneIdentityProviderAttributeUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	idpAttributeMapping := *management.NewIdentityProviderAttribute(d.Get("name").(string), d.Get("value").(string), management.EnumIdentityProviderAttributeMappingType(d.Get("update").(string)))

	_, diags = sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.IdentityProviderManagementIdentityProviderAttributesApi.UpdateIdentityProviderAttribute(ctx, d.Get("environment_id").(string), d.Get("identity_provider_id").(string), d.Id()).IdentityProviderAttribute(idpAttributeMapping).Execute()
		},
		"UpdateIdentityProviderAttribute",
		sdk.CustomErrorInvalidValue,
		sdk.DefaultRetryable,
	)
	if diags.HasError() {
		return diags
	}

	return resourcePingOneIdentityProviderAttributeRead(ctx, d, meta)
}

func resourcePingOneIdentityProviderAttributeDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	_, diags = sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			r, err := apiClient.IdentityProviderManagementIdentityProviderAttributesApi.DeleteIdentityProviderAttribute(ctx, d.Get("environment_id").(string), d.Get("identity_provider_id").(string), d.Id()).Execute()
			return nil, r, err
		},
		"DeleteIdentityProviderAttribute",
		sdk.CustomErrorResourceNotFoundWarning,
		sdk.DefaultRetryable,
	)
	if diags.HasError() {
		return diags
	}

	return diags
}

func resourcePingOneIdentityProviderAttributeImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	attributes := strings.SplitN(d.Id(), "/", 3)

	if len(attributes) != 2 {
		return nil, fmt.Errorf("invalid id (\"%s\") specified, should be in format \"environmentID/applicationID/identityProviderAttributeID\"", d.Id())
	}

	environmentID, applicationID, identityProviderAttributeID := attributes[0], attributes[1], attributes[2]

	d.Set("environment_id", environmentID)
	d.Set("identity_provider_id", applicationID)
	d.SetId(identityProviderAttributeID)

	resourcePingOneIdentityProviderAttributeRead(ctx, d, meta)

	return []*schema.ResourceData{d}, nil
}
