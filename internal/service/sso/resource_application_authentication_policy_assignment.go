package sso

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
)

func ResourceApplicationAuthenticationPolicyAssignment() *schema.Resource {
	return &schema.Resource{

		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage an authentication policy assignment for applications configured in PingOne.",

		CreateContext: resourcePingOneApplicationAuthenticationPolicyAssignmentCreate,
		ReadContext:   resourcePingOneApplicationAuthenticationPolicyAssignmentRead,
		UpdateContext: resourcePingOneApplicationAuthenticationPolicyAssignmentUpdate,
		DeleteContext: resourcePingOneApplicationAuthenticationPolicyAssignmentDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourcePingOneApplicationAuthenticationPolicyAssignmentImport,
		},

		Schema: map[string]*schema.Schema{
			"environment_id": {
				Description:      "The ID of the environment to create the application authentication policy assignment in.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
				ForceNew:         true,
			},
			"application_id": {
				Description:      "The ID of the application to create the authentication policy assignment for.",
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
			},
			"authentication_policy_id": {
				Description:      "The ID of the authentication policy resource to associate.",
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
			},
			"priority": {
				Description:      "The order in which the policy referenced by this assignment is evaluated during an authentication flow relative to other policies. An assignment with a lower priority will be evaluated first.",
				Type:             schema.TypeInt,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(1)),
			},
		},
	}
}

func resourcePingOneApplicationAuthenticationPolicyAssignmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	signOnPolicy := *management.NewSignOnPolicyActionCommonSignOnPolicy(d.Get("authentication_policy_id").(string))
	applicationAuthenticationPolicyAssignment := *management.NewSignOnPolicyAssignment(int32(d.Get("priority").(int)), signOnPolicy)

	resp, r, err := apiClient.ApplicationsApplicationSignOnPolicyAssignmentsApi.CreateSignOnPolicyAssignment(ctx, d.Get("environment_id").(string), d.Get("application_id").(string)).SignOnPolicyAssignment(applicationAuthenticationPolicyAssignment).Execute()
	if (err != nil) || (r.StatusCode != 201) {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error when calling `ApplicationsApplicationSignOnPolicyAssignmentsApi.CreateSignOnPolicyAssignment``: %v", err),
			Detail:   fmt.Sprintf("Full HTTP response: %v\n", r.Body),
		})

		return diags
	}

	d.SetId(resp.GetId())

	return resourcePingOneApplicationAuthenticationPolicyAssignmentRead(ctx, d, meta)
}

func resourcePingOneApplicationAuthenticationPolicyAssignmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	resp, r, err := apiClient.ApplicationsApplicationSignOnPolicyAssignmentsApi.ReadOneSignOnPolicyAssignment(ctx, d.Get("environment_id").(string), d.Get("application_id").(string), d.Id()).Execute()
	if err != nil {

		if r.StatusCode == 404 {
			log.Printf("[INFO] PingOne Application Sign on Policy Mapping %s no longer exists", d.Id())
			d.SetId("")
			return nil
		}
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error when calling `ApplicationsApplicationSignOnPolicyAssignmentsApi.ReadOneSignOnPolicyAssignment``: %v", err),
			Detail:   fmt.Sprintf("Full HTTP response: %v\n", r.Body),
		})

		return diags
	}

	d.Set("priority", resp.GetPriority())

	return diags
}

func resourcePingOneApplicationAuthenticationPolicyAssignmentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	signOnPolicy := *management.NewSignOnPolicyActionCommonSignOnPolicy(d.Get("authentication_policy_id").(string))
	applicationAuthenticationPolicyAssignment := *management.NewSignOnPolicyAssignment(int32(d.Get("priority").(int)), signOnPolicy)

	_, r, err := apiClient.ApplicationsApplicationSignOnPolicyAssignmentsApi.UpdateSignOnPolicyAssignment(ctx, d.Get("environment_id").(string), d.Get("application_id").(string), d.Id()).SignOnPolicyAssignment(applicationAuthenticationPolicyAssignment).Execute()
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error when calling `ApplicationsApplicationSignOnPolicyAssignmentsApi.UpdateSignOnPolicyAssignment``: %v", err),
			Detail:   fmt.Sprintf("Full HTTP response: %v\n", r.Body),
		})

		return diags
	}

	return resourcePingOneApplicationAuthenticationPolicyAssignmentRead(ctx, d, meta)
}

func resourcePingOneApplicationAuthenticationPolicyAssignmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	_, err := apiClient.ApplicationsApplicationSignOnPolicyAssignmentsApi.DeleteSignOnPolicyAssignment(ctx, d.Get("environment_id").(string), d.Get("application_id").(string), d.Id()).Execute()
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error when calling `ApplicationsApplicationSignOnPolicyAssignmentsApi.DeleteSignOnPolicyAssignment``: %v", err),
		})

		return diags
	}

	return nil
}

func resourcePingOneApplicationAuthenticationPolicyAssignmentImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	attributes := strings.SplitN(d.Id(), "/", 2)

	if len(attributes) != 2 {
		return nil, fmt.Errorf("invalid id (\"%s\") specified, should be in format \"environmentID/applicationID/authenticationPolicyAssignmentID\"", d.Id())
	}

	environmentID, applicationID, authenticationPolicyAssignmentID := attributes[0], attributes[1], attributes[2]

	d.Set("environment_id", environmentID)
	d.Set("application_id", applicationID)
	d.SetId(authenticationPolicyAssignmentID)

	resourcePingOneApplicationAuthenticationPolicyAssignmentRead(ctx, d, meta)

	return []*schema.ResourceData{d}, nil
}
