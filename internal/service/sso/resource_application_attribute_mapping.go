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
)

func ResourceApplicationAttributeMapping() *schema.Resource {
	return &schema.Resource{

		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage an attribute mapping for applications configured in PingOne.",

		CreateContext: resourcePingOneApplicationAttributeMappingCreate,
		ReadContext:   resourcePingOneApplicationAttributeMappingRead,
		UpdateContext: resourcePingOneApplicationAttributeMappingUpdate,
		DeleteContext: resourcePingOneApplicationAttributeMappingDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourcePingOneApplicationAttributeMappingImport,
		},

		Schema: map[string]*schema.Schema{
			"environment_id": {
				Description:      "The ID of the environment to create the application attribute mapping in.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
				ForceNew:         true,
			},
			"application_id": {
				Description:      "The ID of the application to create the attribute mapping for.",
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
			},
			"name": {
				Description:      "A string that specifies the name of attribute and must be unique within an application. For SAML applications, the `samlAssertion.subject` name is a reserved case-insensitive name which indicates the mapping to be used for the subject in an assertion. For OpenID Connect applications, the following names are reserved and cannot be used `acr`, `amr`, `at_hash`, `aud`, `auth_time`, `azp`, `client_id`, `exp`, `iat`, `iss`, `jti`, `nbf`, `nonce`, `org`, `scope`, `sid`, `sub`.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
			},
			"required": {
				Description: "A boolean to specify whether a mapping value is required for this attribute. If true, a value must be set and a non-empty value must be available in the SAML assertion or ID token.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"value": {
				Description:      "A string that specifies the string constants or expression for mapping the attribute path against a specific source. The expression format is `${<source>.<attribute_path>}`. The only supported source is user (for example, `${user.id}`).",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
			},
			"mapping_type": {
				Description: "A string that specifies the mapping type of the attribute. Options are `CORE`, `SCOPE`, and `CUSTOM`.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func resourcePingOneApplicationAttributeMappingCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	applicationAttributeMapping := *management.NewApplicationAttributeMapping(d.Get("name").(string), d.Get("required").(bool), d.Get("value").(string))

	resp, r, err := apiClient.ApplicationsApplicationAttributeMappingApi.CreateApplicationAttributeMapping(ctx, d.Get("environment_id").(string), d.Get("application_id").(string)).ApplicationAttributeMapping(applicationAttributeMapping).Execute()
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
			Summary:  fmt.Sprintf("Error when calling `ApplicationsApplicationAttributeMappingApi.CreateApplicationAttributeMapping``: %v", err),
			Detail:   fmt.Sprintf("Full HTTP response: %v\n", r.Body),
		})

		return diags
	}

	d.SetId(resp.GetId())

	return resourcePingOneApplicationAttributeMappingRead(ctx, d, meta)
}

func resourcePingOneApplicationAttributeMappingRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	resp, r, err := apiClient.ApplicationsApplicationAttributeMappingApi.ReadOneApplicationAttributeMapping(ctx, d.Get("environment_id").(string), d.Get("application_id").(string), d.Id()).Execute()
	if err != nil {

		if r.StatusCode == 404 {
			log.Printf("[INFO] PingOne Application Attribute Mapping %s no longer exists", d.Id())
			d.SetId("")
			return nil
		}
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error when calling `ApplicationsApplicationAttributeMappingApi.ReadOneApplicationAttributeMapping``: %v", err),
			Detail:   fmt.Sprintf("Full HTTP response: %v\n", r.Body),
		})

		return diags
	}

	d.Set("name", resp.GetName())
	d.Set("required", resp.GetRequired())
	d.Set("value", resp.GetValue())
	d.Set("mapping_type", resp.GetMappingType())

	return diags
}

func resourcePingOneApplicationAttributeMappingUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	applicationAttributeMapping := *management.NewApplicationAttributeMapping(d.Get("name").(string), d.Get("required").(bool), d.Get("value").(string))

	_, r, err := apiClient.ApplicationsApplicationAttributeMappingApi.UpdateApplicationAttributeMapping(ctx, d.Get("environment_id").(string), d.Get("application_id").(string), d.Id()).ApplicationAttributeMapping(applicationAttributeMapping).Execute()
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
			Summary:  fmt.Sprintf("Error when calling `ApplicationsApplicationAttributeMappingApi.UpdateApplicationAttributeMapping``: %v", err),
			Detail:   fmt.Sprintf("Full HTTP response: %v\n", r.Body),
		})

		return diags
	}

	return resourcePingOneApplicationAttributeMappingRead(ctx, d, meta)
}

func resourcePingOneApplicationAttributeMappingDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	_, err := apiClient.ApplicationsApplicationAttributeMappingApi.DeleteApplicationAttributeMapping(ctx, d.Get("environment_id").(string), d.Get("application_id").(string), d.Id()).Execute()
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error when calling `ApplicationsApplicationAttributeMappingApi.DeleteApplicationAttributeMapping``: %v", err),
		})

		return diags
	}

	return nil
}

func resourcePingOneApplicationAttributeMappingImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	attributes := strings.SplitN(d.Id(), "/", 2)

	if len(attributes) != 2 {
		return nil, fmt.Errorf("invalid id (\"%s\") specified, should be in format \"environmentID/applicationID/attributeMappingID\"", d.Id())
	}

	environmentID, applicationID, attributeMappingID := attributes[0], attributes[1], attributes[2]

	d.Set("environment_id", environmentID)
	d.Set("application_id", applicationID)
	d.SetId(attributeMappingID)

	resourcePingOneApplicationAttributeMappingRead(ctx, d, meta)

	return []*schema.ResourceData{d}, nil
}
