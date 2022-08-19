package sso

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
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

	resp, r, err := apiClient.IdentityProviderManagementIdentityProviderAttributesApi.CreateIdentityProviderAttribute(ctx, d.Get("environment_id").(string), d.Get("identity_provider_id").(string)).IdentityProviderAttribute(idpAttributeMapping).Execute()
	if (err != nil) || (r.StatusCode != 201) {
		response := &management.P1Error{}
		errDecode := json.NewDecoder(r.Body).Decode(response)
		if errDecode == nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  fmt.Sprintf("Cannot decode error response: %v", errDecode),
				Detail:   fmt.Sprintf("Full HTTP response: %v\n", r.Body),
			})
		}

		if r.StatusCode == 400 && response.GetDetails()[0].GetCode() == "INVALID_VALUE" && response.GetDetails()[0].GetTarget() == "name" {
			diags = diag.FromErr(fmt.Errorf(response.GetDetails()[0].GetMessage()))

			return diags
		}

		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error when calling `IdentityProviderManagementIdentityProviderAttributesApi.CreateIdentityProviderAttribute``: %v", err),
			Detail:   fmt.Sprintf("Full HTTP response: %v\n", r.Body),
		})

		return diags
	}

	d.SetId(resp.GetId())

	return resourcePingOneIdentityProviderAttributeRead(ctx, d, meta)
}

func resourcePingOneIdentityProviderAttributeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	resp, r, err := apiClient.IdentityProviderManagementIdentityProviderAttributesApi.ReadOneIdentityProviderAttribute(ctx, d.Get("environment_id").(string), d.Get("identity_provider_id").(string), d.Id()).Execute()
	if err != nil {

		if r.StatusCode == 404 {
			log.Printf("[INFO] PingOne Identity Provider Attribute %s no longer exists", d.Id())
			d.SetId("")
			return nil
		}
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error when calling `IdentityProviderManagementIdentityProviderAttributesApi.ReadOneIdentityProviderAttribute``: %v", err),
			Detail:   fmt.Sprintf("Full HTTP response: %v\n", r.Body),
		})

		return diags
	}

	d.Set("name", resp.GetName())
	d.Set("value", resp.GetValue())
	d.Set("update", resp.GetUpdate())
	d.Set("mapping_type", resp.GetMappingType())

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

	_, r, err := apiClient.IdentityProviderManagementIdentityProviderAttributesApi.UpdateIdentityProviderAttribute(ctx, d.Get("environment_id").(string), d.Get("identity_provider_id").(string), d.Id()).IdentityProviderAttribute(idpAttributeMapping).Execute()
	if err != nil {
		response := &management.P1Error{}
		errDecode := json.NewDecoder(r.Body).Decode(response)
		if errDecode == nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  fmt.Sprintf("Cannot decode error response: %v", errDecode),
				Detail:   fmt.Sprintf("Full HTTP response: %v\n", r.Body),
			})
		}

		if r.StatusCode == 400 && response.GetDetails()[0].GetCode() == "INVALID_VALUE" && response.GetDetails()[0].GetTarget() == "name" {
			diags = diag.FromErr(fmt.Errorf(response.GetDetails()[0].GetMessage()))

			return diags
		}

		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error when calling `IdentityProviderManagementIdentityProviderAttributesApi.UpdateIdentityProviderAttribute``: %v", err),
			Detail:   fmt.Sprintf("Full HTTP response: %v\n", r.Body),
		})

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

	_, err := apiClient.IdentityProviderManagementIdentityProviderAttributesApi.DeleteIdentityProviderAttribute(ctx, d.Get("environment_id").(string), d.Get("identity_provider_id").(string), d.Id()).Execute()
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error when calling `IdentityProviderManagementIdentityProviderAttributesApi.DeleteIdentityProviderAttribute``: %v", err),
		})

		return diags
	}

	return nil
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
