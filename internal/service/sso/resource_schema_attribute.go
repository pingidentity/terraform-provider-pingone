package sso

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	pingone "github.com/patrickcping/pingone-go-sdk-v2/management"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
)

func ResourceSchemaAttribute() *schema.Resource {
	return &schema.Resource{

		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage PingOne schema attributes",

		CreateContext: resourceSchemaAttributeCreate,
		ReadContext:   resourceSchemaAttributeRead,
		UpdateContext: resourceSchemaAttributeUpdate,
		DeleteContext: resourceSchemaAttributeDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceSchemaAttributeImport,
		},

		Schema: map[string]*schema.Schema{
			"environment_id": {
				Description:      "The ID of the environment to create the schema attribute in.",
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
			},
			"schema_id": {
				Description:      "The ID of the schema to apply the schema attribute to.",
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
			},
			"name": {
				Description:      "The system name of the schema attribute.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
			},
			"display_name": {
				Description: "The display name of the attribute such as 'T-shirt size’. If provided, it must not be an empty string. Valid characters consist of any Unicode letter, mark (for example, accent or umlaut), numeric character, forward slash, dot, apostrophe, underscore, space, or hyphen.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"description": {
				Description: "A description of the attribute. If provided, it must not be an empty string. Valid characters consists of any Unicode letter, mark (for example, accent or umlaut), numeric character, punctuation character, or space.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"enabled": {
				Description: "Indicates whether or not the attribute is enabled.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"type": {
				Description:  "The type of the attribute. This can be `STRING`, `JSON`, `BOOLEAN`, or `COMPLEX`. `COMPLEX` and `BOOLEAN` attributes cannot be created, but standard attributes of those types may be updated. `JSON` attributes are limited by size (total size must not exceed 16KB).",
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "STRING",
				ValidateFunc: validation.StringInSlice([]string{"STRING", "JSON", "BOOLEAN", "COMPLEX"}, false),
			},
			"unique": {
				Description: "Indicates whether or not the attribute must have a unique value within the PingOne environment.",
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     false,
			},
			"multivalued": {
				Description: "Indicates whether the attribute has multiple values or a single one.  Maximum number of values stored is 1,000.",
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     false,
			},
			"required": {
				Description: "Indicates whether or not the attribute is required.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"ldap_attribute": {
				Description: "The unique identifier for the LDAP attribute.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"schema_type": {
				Description: "The schema type of the attribute. This can be `CORE`, `STANDARD` or `CUSTOM`. `CORE` and `STANDARD` attributes are supplied by default. `CORE` attributes cannot be updated or deleted. `STANDARD` attributes cannot be deleted, but their mutable properties can be updated. `CUSTOM` attributes can be deleted, and their mutable properties can be updated. New attributes are created with a schema type of `CUSTOM`.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func resourceSchemaAttributeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API
	ctx = context.WithValue(ctx, pingone.ContextServerVariables, map[string]string{
		"suffix": p1Client.RegionSuffix,
	})
	var diags diag.Diagnostics

	schemaAttribute, err := buildSchemaAttribute(d, "CREATE")
	if err != nil {
		return diag.FromErr(err)
	}

	resp, r, err := apiClient.SchemasApi.CreateAttribute(ctx, d.Get("environment_id").(string), d.Get("schema_id").(string)).SchemaAttribute(schemaAttribute.(pingone.SchemaAttribute)).Execute()
	if (err != nil) || (r.StatusCode != 201) {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error when calling `SchemasApi.CreateAttribute``: %v", err),
			Detail:   fmt.Sprintf("Full HTTP response: %v\n", r.Body),
		})

		return diags
	}

	d.SetId(resp.GetId())

	return resourceSchemaAttributeRead(ctx, d, meta)
}

func resourceSchemaAttributeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API
	ctx = context.WithValue(ctx, pingone.ContextServerVariables, map[string]string{
		"suffix": p1Client.RegionSuffix,
	})
	var diags diag.Diagnostics

	resp, r, err := apiClient.SchemasApi.ReadOneAttribute(ctx, d.Get("environment_id").(string), d.Get("schema_id").(string), d.Id()).Execute()
	if err != nil {

		if r.StatusCode == 404 {
			log.Printf("[INFO] PingOne Schema Attribute %s no longer exists", d.Id())
			d.SetId("")
			return nil
		}
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error when calling `SchemasApi.ReadOneSchemaAttribute``: %v", err),
			Detail:   fmt.Sprintf("Full HTTP response: %v\n", r.Body),
		})

		return diags
	}

	log.Printf("resp.GetDisplayName(): %s", resp.GetDisplayName())

	d.Set("name", resp.GetName())

	if v, ok := resp.GetDisplayNameOk(); ok {
		d.Set("display_name", v)
	} else {
		d.Set("display_name", nil)
	}

	if v, ok := resp.GetDescriptionOk(); ok {
		d.Set("description", v)
	} else {
		d.Set("description", nil)
	}

	d.Set("enabled", resp.GetEnabled())
	d.Set("type", resp.GetType())
	d.Set("unique", resp.GetUnique())
	d.Set("multivalued", resp.GetMultiValued())
	d.Set("ldap_attribute", resp.GetLdapAttribute())
	d.Set("required", resp.GetRequired())
	d.Set("schema_type", resp.GetSchemaType())

	return diags
}

func resourceSchemaAttributeUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API
	ctx = context.WithValue(ctx, pingone.ContextServerVariables, map[string]string{
		"suffix": p1Client.RegionSuffix,
	})
	var diags diag.Diagnostics

	schemaAttribute, err := buildSchemaAttribute(d, "UPDATE")
	if err != nil {
		return diag.FromErr(err)
	}

	_, r, err := apiClient.SchemasApi.UpdateAttributePatch(ctx, d.Get("environment_id").(string), d.Get("schema_id").(string), d.Id()).SchemaAttribute(schemaAttribute.(pingone.SchemaAttribute)).Execute()
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error when calling `SchemasApi.UpdateAttributePatch``: %v", err),
			Detail:   fmt.Sprintf("Full HTTP response: %v\n", r.Body),
		})

		return diags
	}

	return resourceSchemaAttributeRead(ctx, d, meta)
}

func resourceSchemaAttributeDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API
	ctx = context.WithValue(ctx, pingone.ContextServerVariables, map[string]string{
		"suffix": p1Client.RegionSuffix,
	})
	var diags diag.Diagnostics

	_, err := apiClient.SchemasApi.DeleteAttribute(ctx, d.Get("environment_id").(string), d.Get("schema_id").(string), d.Id()).Execute()
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error when calling `SchemasApi.DeleteAttribute``: %v", err),
		})

		return diags
	}

	return nil
}

func resourceSchemaAttributeImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	attributes := strings.SplitN(d.Id(), "/", 3)

	if len(attributes) != 2 {
		return nil, fmt.Errorf("invalid id (\"%s\") specified, should be in format \"environmentID/schemaID/attributeID\"", d.Id())
	}

	environmentID, schemaID, attributeID := attributes[0], attributes[1], attributes[2]

	d.Set("environment_id", environmentID)
	d.Set("schema_id", schemaID)

	d.SetId(attributeID)

	resourceSchemaAttributeRead(ctx, d, meta)

	return []*schema.ResourceData{d}, nil
}

func buildSchemaAttribute(d *schema.ResourceData, action string) (interface{}, error) {

	attrType := d.Get("type").(string)

	if (attrType == "BOOLEAN" || attrType == "COMPLEX") && action == "CREATE" {
		return nil, fmt.Errorf("Cannot create attributes of type BOOLEAN or COMPLEX.  Custom attributes must be either STRING or JSON.  Attribute type found: %s", attrType)
	}

	schemaAttribute := *pingone.NewSchemaAttribute(d.Get("enabled").(bool), d.Get("name").(string), attrType) // SchemaAttribute |  (optional)

	if v, ok := d.GetOk("display_name"); ok {
		schemaAttribute.SetDisplayName(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		schemaAttribute.SetDescription(v.(string))
	}

	attrUnique := d.Get("unique").(bool)

	if attrUnique && attrType != "STRING" {
		return nil, fmt.Errorf("Cannot set attribute unique parameter when the attribute type is not STRING.  Attribute type found: %s", attrType)
	}

	schemaAttribute.SetUnique(attrUnique)

	schemaAttribute.SetMultiValued(d.Get("multivalued").(bool))

	schemaAttribute.SetRequired(d.Get("required").(bool))

	return schemaAttribute, nil
}
