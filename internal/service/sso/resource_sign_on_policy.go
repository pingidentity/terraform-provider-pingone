package sso

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func ResourceSignOnPolicy() *schema.Resource {
	return &schema.Resource{

		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage PingOne sign on policies",

		CreateContext: resourceSignOnPolicyCreate,
		ReadContext:   resourceSignOnPolicyRead,
		UpdateContext: resourceSignOnPolicyUpdate,
		DeleteContext: resourceSignOnPolicyDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceSignOnPolicyImport,
		},

		Schema: map[string]*schema.Schema{
			"environment_id": {
				Description:      "The ID of the environment to create the sign on policy in.",
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
			},
			"name": {
				Description:      "A string that specifies the resource name. The name must be unique within the environment, and can consist of either a string of alphanumeric letters, underscore, hyphen, period `^[a-zA-Z0-9_. -]+$` or an absolute URI if the string contains a `:` character.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9_. -]+$`), "Names must consist of either a string of alphanumeric letters, underscore, hyphen, period `^[a-zA-Z0-9_. -]+$`")), // TODO regex
			},
			"description": {
				Description: "A string that specifies the description of the sign-on policy.",
				Type:        schema.TypeString,
				Optional:    true,
			},
		},
	}
}

func resourceSignOnPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	signOnPolicy := *management.NewSignOnPolicy(d.Get("name").(string))

	if v, ok := d.GetOk("description"); ok {
		signOnPolicy.SetDescription(v.(string))
	}

	resp, r, err := apiClient.SignOnPoliciesSignOnPoliciesApi.CreateSignOnPolicy(ctx, d.Get("environment_id").(string)).SignOnPolicy(signOnPolicy).Execute()
	if (err != nil) || (r.StatusCode != 201) {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error when calling `SignOnPoliciesSignOnPoliciesApi.CreateSignOnPolicy``: %v", err),
			Detail:   fmt.Sprintf("Full HTTP response: %v\n", r.Body),
		})

		return diags
	}

	d.SetId(resp.GetId())

	return resourceSignOnPolicyRead(ctx, d, meta)
}

func resourceSignOnPolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	resp, r, err := apiClient.SignOnPoliciesSignOnPoliciesApi.ReadOneSignOnPolicy(ctx, d.Get("environment_id").(string), d.Id()).Execute()
	if err != nil {

		if r.StatusCode == 404 {
			log.Printf("[INFO] PingOne Sign on policy %s no longer exists", d.Id())
			d.SetId("")
			return nil
		}
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error when calling `SignOnPoliciesSignOnPoliciesApi.ReadOneSignOnPolicy``: %v", err),
			Detail:   fmt.Sprintf("Full HTTP response: %v\n", r.Body),
		})

		return diags
	}

	d.Set("name", resp.GetName())

	if v, ok := resp.GetDescriptionOk(); ok {
		d.Set("description", v)
	} else {
		d.Set("description", nil)
	}

	return diags
}

func resourceSignOnPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	signOnPolicy := *management.NewSignOnPolicy(d.Get("name").(string))

	if v, ok := d.GetOk("description"); ok {
		signOnPolicy.SetDescription(v.(string))
	}

	_, r, err := apiClient.SignOnPoliciesSignOnPoliciesApi.UpdateSignOnPolicy(ctx, d.Get("environment_id").(string), d.Id()).SignOnPolicy(signOnPolicy).Execute()
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error when calling `SignOnPoliciesSignOnPoliciesApi.UpdateSignOnPolicy``: %v", err),
			Detail:   fmt.Sprintf("Full HTTP response: %v\n", r.Body),
		})

		return diags
	}

	return resourceSignOnPolicyRead(ctx, d, meta)
}

func resourceSignOnPolicyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	r, err := apiClient.SignOnPoliciesSignOnPoliciesApi.DeleteSignOnPolicy(ctx, d.Get("environment_id").(string), d.Id()).Execute()
	if err != nil {

		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error when calling `SignOnPoliciesSignOnPoliciesApi.DeleteSignOnPolicy``: %v", err),
			Detail:   fmt.Sprintf("Full HTTP response: %v\n", r.Body),
		})

		return diags
	}

	return nil
}

func resourceSignOnPolicyImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	attributes := strings.SplitN(d.Id(), "/", 3)

	if len(attributes) != 2 {
		return nil, fmt.Errorf("invalid id (\"%s\") specified, should be in format \"environmentID/signOnPolicyID\"", d.Id())
	}

	environmentID, signOnPolicyID := attributes[0], attributes[1]

	d.Set("environment_id", environmentID)

	d.SetId(signOnPolicyID)

	resourceSignOnPolicyRead(ctx, d, meta)

	return []*schema.ResourceData{d}, nil
}
