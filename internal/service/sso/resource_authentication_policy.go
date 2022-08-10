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
)

func ResourceAuthenticationPolicy() *schema.Resource {
	return &schema.Resource{

		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage PingOne authentication policies",

		CreateContext: resourceAuthenticationPolicyCreate,
		ReadContext:   resourceAuthenticationPolicyRead,
		UpdateContext: resourceAuthenticationPolicyUpdate,
		DeleteContext: resourceAuthenticationPolicyDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceAuthenticationPolicyImport,
		},

		Schema: map[string]*schema.Schema{
			"environment_id": {
				Description:      "The ID of the environment to create the authentication policy in.",
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
			},
			"name": {
				Description:      "A string that specifies the resource name. The name must be unique within the environment, and can consist of either a string of alphanumeric letters, underscore, hyphen, period `^[a-zA-Z0-9_. -]+$` or an absolute URI if the string contains a `:` character.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9_. -]+$`))), // TODO regex
			},
			"description": {
				Description: "A string that specifies the description of the sign-on policy.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"policy_action": {
				Description: "One or more action(s) to apply to the authentication policy.",
				Type:        schema.TypeList,
				MaxItems:    10,
				Required:    true,
				ForceNew:    true,
				Elem: &schema.Resource{
					Schema: resourceAuthenticationPolicyActionSchema(),
				},
			},
		},
	}
}

func resourceAuthenticationPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	// Policy actions
	for itemIndex, item := range d.Get("policy_action").([]interface{}) {

		sopAction, diags := expandSOPAction(item, int32(itemIndex+1))

		_, r, err := apiClient.SignOnPoliciesSignOnPolicyActionsApi.CreateSignOnPolicyAction(ctx, d.Get("environment_id").(string), resp.GetId()).SignOnPolicyAction(*sopAction).Execute()
		if (err != nil) || (r.StatusCode != 201) {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Error when calling `SignOnPoliciesSignOnPolicyActionsApi.CreateSignOnPolicyAction``: %v", err),
				Detail:   fmt.Sprintf("Full HTTP response: %v\n", r.Body),
			})

			return diags
		}
	}

	return resourceAuthenticationPolicyRead(ctx, d, meta)
}

func resourceAuthenticationPolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	// Policy Actions
	respActions, r, err := apiClient.SignOnPoliciesSignOnPolicyActionsApi.ReadAllSignOnPolicyActions(ctx, d.Get("environment_id").(string), d.Id()).Execute()
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  fmt.Sprintf("Error when calling `SignOnPoliciesSignOnPolicyActionsApi.ReadAllSignOnPolicyActions``: %v", err),
			Detail:   fmt.Sprintf("Full HTTP response: %v\n", r.Body),
		})
	}

	if v, ok := respActions.Embedded.GetActionsOk(); ok {

		sopActions, err := flattenSOPActions(v)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  fmt.Sprintf("Error when parsing authentication policy actions: %v", err),
			})
			d.Set("policy_action", nil)

		} else {
			d.Set("policy_action", sopActions)

		}
	}

	return diags
}

func resourceAuthenticationPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	// Policy actions
	// policyActions := d.Get("policy_action").([]interface{})

	// respActions, r, err := apiClient.SignOnPoliciesSignOnPolicyActionsApi.ReadAllSignOnPolicyActions(ctx, d.Get("environment_id").(string), d.Id()).Execute()
	// if err != nil {
	// 	diags = append(diags, diag.Diagnostic{
	// 		Severity: diag.Warning,
	// 		Summary:  fmt.Sprintf("Error when calling `SignOnPoliciesSignOnPolicyActionsApi.ReadAllSignOnPolicyActions``: %v", err),
	// 		Detail:   fmt.Sprintf("Full HTTP response: %v\n", r.Body),
	// 	})

	// 	return diags
	// } else {

	// 	// get the list of actions
	// 	if v, ok := respActions.Embedded.GetActionsOk(); ok {

	// 		// count difference. +ve is to add to api, -ve is to delete from api
	// 		indexDiff := len(policyActions) - len(v)

	// 		// Sort by priority first, this way we can keep the ID by priority consistent
	// 		sort.Slice(v, func(i, j int) bool {
	// 			return getActionPriority(v[i]) < getActionPriority(v[j])
	// 		})

	// 		// Loop the API actions, we update as we go and remove excess if necessary
	// 		for itemIndex, item := range v {

	// 			// Get the action ID to use in further API calls
	// 			actionID := getActionID(item)
	// 			if actionID == "" {
	// 				diags = append(diags, diag.Diagnostic{
	// 					Severity: diag.Error,
	// 					Summary:  fmt.Sprintf("Cannot get action ID for action %v", item),
	// 				})

	// 				return diags
	// 			} else {

	// 				// update any api actions that have an associated configured action
	// 				if itemIndex < len(policyActions) {
	// 					sopAction, err := expandSOPAction(policyActions[itemIndex], int32(itemIndex+1))
	// 					if err != nil {
	// 						diags = append(diags, diag.Diagnostic{
	// 							Severity: diag.Warning,
	// 							Summary:  fmt.Sprintf("Cannot expand policy action: %v", err),
	// 						})
	// 					} else {

	// 						_, r, err := apiClient.SignOnPoliciesSignOnPolicyActionsApi.UpdateSignOnPolicyAction(ctx, d.Get("environment_id").(string), d.Id(), actionID).SignOnPolicyAction(*sopAction).Execute()
	// 						if (err != nil) || (r.StatusCode != 201) {
	// 							diags = append(diags, diag.Diagnostic{
	// 								Severity: diag.Warning,
	// 								Summary:  fmt.Sprintf("Error when calling `SignOnPoliciesSignOnPolicyActionsApi.UpdateSignOnPolicyAction``: %v", err),
	// 								Detail:   fmt.Sprintf("Full HTTP response: %v\n", r.Body),
	// 							})
	// 						}
	// 					}

	// 				} else { // delete any api actions that exist without a configured action

	// 					r, err := apiClient.SignOnPoliciesSignOnPolicyActionsApi.DeleteSignOnPolicyAction(ctx, d.Get("environment_id").(string), d.Id(), actionID).Execute()
	// 					if (err != nil) || (r.StatusCode != 204) {
	// 						diags = append(diags, diag.Diagnostic{
	// 							Severity: diag.Warning,
	// 							Summary:  fmt.Sprintf("Error when calling `SignOnPoliciesSignOnPolicyActionsApi.DeleteSignOnPolicyAction``: %v", err),
	// 							Detail:   fmt.Sprintf("Full HTTP response: %v\n", r.Body),
	// 						})
	// 					}

	// 				}
	// 			}

	// 		}

	// 		// add any new actions
	// 		for i := 0; i >= 0 && i < indexDiff; i += 1 {

	// 			policyActionIndex := i + len(v)

	// 			sopAction, err := expandSOPAction(policyActions[policyActionIndex], int32(policyActionIndex+1))
	// 			if err != nil {
	// 				return diag.FromErr(err)
	// 			}

	// 			_, r, err := apiClient.SignOnPoliciesSignOnPolicyActionsApi.CreateSignOnPolicyAction(ctx, d.Get("environment_id").(string), d.Id()).SignOnPolicyAction(*sopAction).Execute()
	// 			if (err != nil) || (r.StatusCode != 201) {
	// 				diags = append(diags, diag.Diagnostic{
	// 					Severity: diag.Error,
	// 					Summary:  fmt.Sprintf("Error when calling `SignOnPoliciesSignOnPolicyActionsApi.CreateSignOnPolicyAction``: %v", err),
	// 					Detail:   fmt.Sprintf("Full HTTP response: %v\n", r.Body),
	// 				})

	// 				return diags
	// 			}

	// 		}

	// 	} else {
	// 		diags = append(diags, diag.Diagnostic{
	// 			Severity: diag.Error,
	// 			Summary:  "No existing authentication policy actions found, this is a bug in the provider",
	// 			Detail:   fmt.Sprintf("Authentication policies must have at least one action configured.  Policy %s has returned no policy actions", d.Id()),
	// 		})

	// 		return diags
	// 	}
	// }

	return resourceAuthenticationPolicyRead(ctx, d, meta)
}

func resourceAuthenticationPolicyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	_, err := apiClient.SignOnPoliciesSignOnPoliciesApi.DeleteSignOnPolicy(ctx, d.Get("environment_id").(string), d.Id()).Execute()
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error when calling `SignOnPoliciesSignOnPoliciesApi.DeleteSignOnPolicy``: %v", err),
		})

		return diags
	}

	return nil
}

func resourceAuthenticationPolicyImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	attributes := strings.SplitN(d.Id(), "/", 3)

	if len(attributes) != 2 {
		return nil, fmt.Errorf("invalid id (\"%s\") specified, should be in format \"environmentID/authenticationPolicyID\"", d.Id())
	}

	environmentID, authenticationPolicyID := attributes[0], attributes[1]

	d.Set("environment_id", environmentID)

	d.SetId(authenticationPolicyID)

	resourceAuthenticationPolicyRead(ctx, d, meta)

	return []*schema.ResourceData{d}, nil
}
