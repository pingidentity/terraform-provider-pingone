package sso

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
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

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.SignOnPoliciesSignOnPoliciesApi.CreateSignOnPolicy(ctx, d.Get("environment_id").(string)).SignOnPolicy(signOnPolicy).Execute()
		},
		"CreateSignOnPolicy",
		sdk.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
	)
	if diags.HasError() {
		return diags
	}

	respObject := resp.(*management.SignOnPolicy)

	d.SetId(respObject.GetId())

	return resourceSignOnPolicyRead(ctx, d, meta)
}

func resourceSignOnPolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.SignOnPoliciesSignOnPoliciesApi.ReadOneSignOnPolicy(ctx, d.Get("environment_id").(string), d.Id()).Execute()
		},
		"ReadOneSignOnPolicy",
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

	respObject := resp.(*management.SignOnPolicy)

	d.Set("name", respObject.GetName())

	if v, ok := respObject.GetDescriptionOk(); ok {
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

	_, diags = sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.SignOnPoliciesSignOnPoliciesApi.UpdateSignOnPolicy(ctx, d.Get("environment_id").(string), d.Id()).SignOnPolicy(signOnPolicy).Execute()
		},
		"UpdateSignOnPolicy",
		sdk.DefaultCustomError,
		sdk.DefaultRetryable,
	)
	if diags.HasError() {
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

	_, diags = sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			r, err := apiClient.SignOnPoliciesSignOnPoliciesApi.DeleteSignOnPolicy(ctx, d.Get("environment_id").(string), d.Id()).Execute()
			return nil, r, err
		},
		"DeleteSignOnPolicy",
		sdk.CustomErrorResourceNotFoundWarning,
		sdk.DefaultRetryable,
	)
	if diags.HasError() {
		return diags
	}

	return diags
}

func resourceSignOnPolicyImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	attributes := strings.SplitN(d.Id(), "/", 2)

	if len(attributes) != 2 {
		return nil, fmt.Errorf("invalid id (\"%s\") specified, should be in format \"environmentID/signOnPolicyID\"", d.Id())
	}

	environmentID, signOnPolicyID := attributes[0], attributes[1]

	d.Set("environment_id", environmentID)

	d.SetId(signOnPolicyID)

	resourceSignOnPolicyRead(ctx, d, meta)

	return []*schema.ResourceData{d}, nil
}
