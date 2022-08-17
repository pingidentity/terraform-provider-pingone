package sso

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
)

func ResourceSignOnPolicyAction() *schema.Resource {
	return &schema.Resource{

		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage PingOne sign on policy actions.",

		CreateContext: resourceSignOnPolicyActionCreate,
		ReadContext:   resourceSignOnPolicyActionRead,
		UpdateContext: resourceSignOnPolicyActionUpdate,
		DeleteContext: resourceSignOnPolicyActionDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceSignOnPolicyActionImport,
		},

		Schema: resourceSignOnPolicyActionSchema(),
	}
}

func resourceSignOnPolicyActionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	signOnPolicyAction, diags := expandSOPAction(d)
	if diags.HasError() {
		return diags
	}

	resp, r, err := apiClient.SignOnPoliciesSignOnPolicyActionsApi.CreateSignOnPolicyAction(ctx, d.Get("environment_id").(string), d.Get("sign_on_policy_id").(string)).SignOnPolicyAction(*signOnPolicyAction).Execute()
	if (err != nil) || (r.StatusCode != 201) {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error when calling `SignOnPoliciesSignOnPolicyActionsApi.CreateSignOnPolicyAction``: %v", err),
			Detail:   fmt.Sprintf("Full HTTP response: %v\n", r.Body),
		})

		return diags
	}

	d.SetId(getActionID(*resp))

	return resourceSignOnPolicyActionRead(ctx, d, meta)
}

func resourceSignOnPolicyActionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	resp, r, err := apiClient.SignOnPoliciesSignOnPolicyActionsApi.ReadOneSignOnPolicyAction(ctx, d.Get("environment_id").(string), d.Get("sign_on_policy_id").(string), d.Id()).Execute()
	if err != nil {

		if r.StatusCode == 404 {
			log.Printf("[INFO] PingOne Sign on policy action %s no longer exists", d.Id())
			d.SetId("")
			return nil
		}
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error when calling `SignOnPoliciesSignOnPolicyActionsApi.ReadOneSignOnPolicyAction``: %v", err),
			Detail:   fmt.Sprintf("Full HTTP response: %v\n", r.Body),
		})

		return diags
	}

	values := map[string]interface{}{
		"priority":                               nil,
		"conditions":                             nil,
		"registration_external_href":             nil,
		"registration_local_population_id":       nil,
		"social_provider_ids":                    nil,
		"registration_confirm_user_attributes":   nil,
		"enforce_lockout_for_identity_providers": nil,
		"agreement":                              nil,
		"identifier_first":                       nil,
		"identity_provider":                      nil,
		"login":                                  nil,
		"mfa":                                    nil,
		"progressive_profiling":                  nil,
	}

	switch resp.GetActualInstance().(type) {
	case *management.SignOnPolicyActionLogin:

		values["priority"] = resp.SignOnPolicyActionLogin.GetPriority()

		if v, ok := resp.SignOnPolicyActionLogin.GetConditionOk(); ok {
			var conditions interface{}
			conditions, diags = flattenConditions(*v)
			if diags.HasError() {
				return diags
			}
			values["conditions"] = conditions
		}

		if v, ok := resp.SignOnPolicyActionLogin.GetRegistrationOk(); ok {
			if v1, ok := v.GetExternalOk(); ok {
				values["registration_external_href"] = v1.GetHref()
			}

			if v1, ok := v.GetPopulationOk(); ok {
				values["registration_local_population_id"] = v1.GetId()
			}

			if v1, ok := v.GetConfirmIdentityProviderAttributesOk(); ok {
				values["registration_confirm_user_attributes"] = v1
			}
		}

		if v, ok := resp.SignOnPolicyActionLogin.GetSocialProvidersOk(); ok {
			values["social_provider_ids"] = flattenActionSocialProvidersInner(v)
		}

		if v, ok := resp.SignOnPolicyActionLogin.GetEnforceLockoutForIdentityProvidersOk(); ok {
			values["enforce_lockout_for_identity_providers"] = v
		}

		values["login"] = flattenActionLogin(resp.SignOnPolicyActionLogin)

	case *management.SignOnPolicyActionAgreement:

		values["priority"] = resp.SignOnPolicyActionAgreement.GetPriority()

		if v, ok := resp.SignOnPolicyActionAgreement.GetConditionOk(); ok {
			var conditions interface{}
			conditions, diags = flattenConditions(*v)
			if diags.HasError() {
				return diags
			}
			values["conditions"] = conditions
		}

		values["agreement"] = flattenActionAgreement(resp.SignOnPolicyActionAgreement)

	case *management.SignOnPolicyActionIDFirst:

		values["priority"] = resp.SignOnPolicyActionIDFirst.GetPriority()

		if v, ok := resp.SignOnPolicyActionIDFirst.GetConditionOk(); ok {
			var conditions interface{}
			conditions, diags = flattenConditions(*v)
			if diags.HasError() {
				return diags
			}
			values["conditions"] = conditions
		}

		if v, ok := resp.SignOnPolicyActionIDFirst.GetRegistrationOk(); ok {
			if v1, ok := v.GetExternalOk(); ok {
				values["registration_external_href"] = v1.GetHref()
			}

			if v1, ok := v.GetPopulationOk(); ok {
				values["registration_local_population_id"] = v1.GetId()
			}

			if v1, ok := v.GetConfirmIdentityProviderAttributesOk(); ok {
				values["registration_confirm_user_attributes"] = v1
			}
		}

		if v, ok := resp.SignOnPolicyActionIDFirst.GetSocialProvidersOk(); ok {
			values["social_provider_ids"] = flattenActionSocialProvidersInner(v)
		}

		if v, ok := resp.SignOnPolicyActionIDFirst.GetEnforceLockoutForIdentityProvidersOk(); ok {
			values["enforce_lockout_for_identity_providers"] = v
		}

		var idFirst []interface{}
		idFirst, diags = flattenActionIDFirst(resp.SignOnPolicyActionIDFirst)
		if diags.HasError() {
			return diags
		}
		values["identifier_first"] = idFirst

	case *management.SignOnPolicyActionIDP:

		values["priority"] = resp.SignOnPolicyActionIDP.GetPriority()

		if v, ok := resp.SignOnPolicyActionIDP.GetConditionOk(); ok {
			var conditions interface{}
			conditions, diags = flattenConditions(*v)
			if diags.HasError() {
				return diags
			}
			values["conditions"] = conditions
		}

		if v, ok := resp.SignOnPolicyActionIDP.GetRegistrationOk(); ok {
			if v1, ok := v.GetPopulationOk(); ok {
				values["registration_local_population_id"] = v1.GetId()
			}

			if v1, ok := v.GetConfirmIdentityProviderAttributesOk(); ok {
				values["registration_confirm_user_attributes"] = v1
			}
		}

		values["identity_provider"] = flattenActionIDP(resp.SignOnPolicyActionIDP)

	case *management.SignOnPolicyActionProgressiveProfiling:

		values["priority"] = resp.SignOnPolicyActionProgressiveProfiling.GetPriority()

		if v, ok := resp.SignOnPolicyActionProgressiveProfiling.GetConditionOk(); ok {
			var conditions interface{}
			conditions, diags = flattenConditions(*v)
			if diags.HasError() {
				return diags
			}
			values["conditions"] = conditions
		}

		values["progressive_profiling"] = flattenActionProgressiveProfiling(resp.SignOnPolicyActionProgressiveProfiling)

	case *management.SignOnPolicyActionMFA:

		values["priority"] = resp.SignOnPolicyActionMFA.GetPriority()

		if v, ok := resp.SignOnPolicyActionMFA.GetConditionOk(); ok {
			var conditions interface{}
			conditions, diags = flattenConditions(*v)
			if diags.HasError() {
				return diags
			}
			values["conditions"] = conditions
		}

		values["mfa"] = flattenActionMFA(resp.SignOnPolicyActionMFA)

	}

	d.Set("priority", values["priority"])
	d.Set("conditions", values["conditions"])

	d.Set("registration_external_href", values["registration_external_href"])
	d.Set("registration_local_population_id", values["registration_local_population_id"])
	d.Set("registration_confirm_user_attributes", values["registration_confirm_user_attributes"])

	d.Set("social_provider_ids", values["social_provider_ids"])
	d.Set("enforce_lockout_for_identity_providers", values["enforce_lockout_for_identity_providers"])

	d.Set("agreement", values["agreement"])
	d.Set("identifier_first", values["identifier_first"])
	d.Set("identity_provider", values["identity_provider"])
	d.Set("login", values["login"])
	d.Set("mfa", values["mfa"])
	d.Set("progressive_profiling", values["progressive_profiling"])

	return diags
}

func resourceSignOnPolicyActionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	signOnPolicyAction, diags := expandSOPAction(d)
	if diags.HasError() {
		return diags
	}

	_, r, err := apiClient.SignOnPoliciesSignOnPolicyActionsApi.UpdateSignOnPolicyAction(ctx, d.Get("environment_id").(string), d.Get("sign_on_policy_id").(string), d.Id()).SignOnPolicyAction(*signOnPolicyAction).Execute()
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error when calling `SignOnPoliciesSignOnPolicyActionsApi.UpdateSignOnPolicyAction``: %v", err),
			Detail:   fmt.Sprintf("Full HTTP response: %v\n", r.Body),
		})

		return diags
	}

	return resourceSignOnPolicyActionRead(ctx, d, meta)
}

func resourceSignOnPolicyActionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	r, err := apiClient.SignOnPoliciesSignOnPolicyActionsApi.DeleteSignOnPolicyAction(ctx, d.Get("environment_id").(string), d.Get("sign_on_policy_id").(string), d.Id()).Execute()
	if err != nil || r.StatusCode != 204 {
		response := &management.P1Error{}
		errDecode := json.NewDecoder(r.Body).Decode(response)
		if errDecode == nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  fmt.Sprintf("Cannot decode error response: %v", errDecode),
				Detail:   fmt.Sprintf("Full HTTP response: %v\n", r.Body),
			})
		}

		if r.StatusCode == 400 && response.GetDetails()[0].GetCode() == "CONSTRAINT_VIOLATION" {
			if match, _ := regexp.MatchString("Cannot delete last action from the policy", response.GetDetails()[0].GetMessage()); match {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Cannot delete last action from the policy.  The remaining policy action is left with no state.",
				})

				return diags
			}
		}

		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error when calling `SignOnPoliciesSignOnPolicyActionsApi.DeleteSignOnPolicyAction``: %v", err),
			Detail:   fmt.Sprintf("Full HTTP response: %v\n", r.Body),
		})

		return diags
	}

	return nil
}

func resourceSignOnPolicyActionImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	attributes := strings.SplitN(d.Id(), "/", 3)

	if len(attributes) != 2 {
		return nil, fmt.Errorf("invalid id (\"%s\") specified, should be in format \"environmentID/signOnPolicyID/policyActionID\"", d.Id())
	}

	environmentID, signOnPolicyID, policyActionID := attributes[0], attributes[1], attributes[2]

	d.Set("environment_id", environmentID)
	d.Set("sign_on_policy_id", signOnPolicyID)

	d.SetId(policyActionID)

	resourceSignOnPolicyActionRead(ctx, d, meta)

	return []*schema.ResourceData{d}, nil
}

func expandSOPAction(d *schema.ResourceData) (*management.SignOnPolicyAction, diag.Diagnostics) {

	signOnPolicyAction := &management.SignOnPolicyAction{}
	var diags diag.Diagnostics

	sopPriority := int32(d.Get("priority").(int))

	processedCount := 0

	if _, ok := d.GetOk("agreement"); ok {
		signOnPolicyAction.SignOnPolicyActionAgreement, diags = expandSOPActionAgreement(d, sopPriority)
		processedCount += 1
	}

	if _, ok := d.GetOk("identifier_first"); ok {
		signOnPolicyAction.SignOnPolicyActionIDFirst, diags = expandSOPActionIDFirst(d, sopPriority)
		processedCount += 1
	}

	if _, ok := d.GetOk("identity_provider"); ok {
		signOnPolicyAction.SignOnPolicyActionIDP, diags = expandSOPActionIDP(d, sopPriority)
		processedCount += 1
	}

	if _, ok := d.GetOk("login"); ok {
		signOnPolicyAction.SignOnPolicyActionLogin, diags = expandSOPActionLogin(d, sopPriority)
		processedCount += 1
	}

	if _, ok := d.GetOk("mfa"); ok {
		signOnPolicyAction.SignOnPolicyActionMFA, diags = expandSOPActionMFA(d, sopPriority)
		processedCount += 1
	}

	if _, ok := d.GetOk("progressive_profiling"); ok {
		signOnPolicyAction.SignOnPolicyActionProgressiveProfiling, diags = expandSOPActionProgressiveProfiling(d, sopPriority)
		processedCount += 1
	}

	if processedCount > 1 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "More than one policy action type configured.  This is not supported.",
		})
		return nil, diags
	} else if processedCount == 0 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "No policy action types configured.  This is not supported.",
		})
		return nil, diags
	}

	return signOnPolicyAction, diags
}

func expandSOPActionAgreement(d *schema.ResourceData, sopPriority int32) (*management.SignOnPolicyActionAgreement, diag.Diagnostics) {
	var diags diag.Diagnostics

	if v, ok := d.Get("agreement").([]interface{}); ok && v != nil && len(v) > 0 && v[0] != nil {
		vp := v[0].(map[string]interface{})

		sopActionType := management.NewSignOnPolicyActionAgreement(
			sopPriority,
			management.ENUMSIGNONPOLICYTYPE_AGREEMENT,
			*management.NewSignOnPolicyActionAgreementAllOfAgreement(vp["agreement_id"].(string)),
		)

		if vc, ok := d.GetOk("conditions"); ok {
			if vc1, ok := vc.([]interface{}); ok && vc1 != nil && len(vc1) > 0 && vc1[0] != nil {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Block `conditions` has no effect when using the agreement action type",
				})
			}
		}

		if vd, ok := vp["show_decline_option"].(bool); ok {
			sopActionType.SetDisableDeclineOption(!vd)
		}

		return sopActionType, diags

	}

	diags = append(diags, diag.Diagnostic{
		Severity: diag.Error,
		Summary:  "Block `agreement` with `agreement_id` must be defined when using the agreement action type",
	})

	return nil, diags
}

func expandSOPActionIDFirst(d *schema.ResourceData, sopPriority int32) (*management.SignOnPolicyActionIDFirst, diag.Diagnostics) {

	var diags diag.Diagnostics

	sopActionType := management.NewSignOnPolicyActionIDFirst(
		sopPriority,
		management.ENUMSIGNONPOLICYTYPE_IDENTIFIER_FIRST,
	)

	if v, ok := d.GetOk("conditions"); ok {
		if vc, ok := v.([]interface{}); ok && vc != nil && len(vc) > 0 && vc[0] != nil {
			var conditions *management.SignOnPolicyActionCommonConditionOrOrInner
			conditions, diags = expandSOPActionCondition(vc[0], management.ENUMSIGNONPOLICYTYPE_IDENTIFIER_FIRST, sopPriority)
			sopActionType.SetCondition(*conditions)
		}
	}

	if v, ok := d.GetOk("registration_external_href"); ok && v != "" {
		obj := *management.NewSignOnPolicyActionLoginAllOfRegistration(false)
		obj.SetExternal(*management.NewSignOnPolicyActionLoginAllOfRegistrationExternal(v.(string)))
		sopActionType.SetRegistration(obj)
	}

	if v, ok := d.GetOk("registration_local_population_id"); ok && v != "" {
		obj := *management.NewSignOnPolicyActionLoginAllOfRegistration(true)
		obj.SetPopulation(*management.NewSignOnPolicyActionLoginAllOfRegistrationPopulation(v.(string)))

		if v1, ok := d.GetOk("registration_confirm_user_attributes"); ok {
			obj.SetConfirmIdentityProviderAttributes(v1.(bool))
		}

		sopActionType.SetRegistration(obj)
	}

	socialIDSet := false
	if v, ok := d.GetOk("social_provider_ids"); ok {
		if vc, ok := v.([]interface{}); ok && vc != nil && len(vc) > 0 && vc[0] != "" {
			obj := make([]string, 0)
			for _, str := range vc {
				obj = append(obj, str.(string))
			}
			sopActionType.SetSocialProviders(expandSOPActionSocialProviders(obj))
			socialIDSet = true
		}
	}

	if v, ok := d.GetOk("enforce_lockout_for_identity_providers"); ok {
		if socialIDSet {
			sopActionType.SetEnforceLockoutForIdentityProviders(v.(bool))
		} else {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "`enforce_lockout_for_identity_providers`, where `social_provider_ids` is not set, has no effect.",
			})
		}
	}

	if v, ok := d.Get("identifier_first").([]interface{}); ok && v != nil && len(v) > 0 && v[0] != nil {
		vp := v[0].(map[string]interface{})

		if v1, ok := vp["recovery_enabled"].(bool); ok {
			sopActionType.SetRecovery(*management.NewSignOnPolicyActionLoginAllOfRecovery(v1))
		}

		if v1, ok := vp["discovery_rule"].([]interface{}); ok && v1 != nil && len(v1) > 0 && v1[0] != nil {
			sopActionType.SetDiscoveryRules(expandSOPActionDiscoveryRules(v1))
		}

	}

	return sopActionType, diags

}

func expandSOPActionIDP(d *schema.ResourceData, sopPriority int32) (*management.SignOnPolicyActionIDP, diag.Diagnostics) {
	var diags diag.Diagnostics

	if v, ok := d.Get("identity_provider").([]interface{}); ok && v != nil && len(v) > 0 && v[0] != nil {
		vp := v[0].(map[string]interface{})

		sopActionType := management.NewSignOnPolicyActionIDP(
			sopPriority,
			management.ENUMSIGNONPOLICYTYPE_IDENTITY_PROVIDER,
			*management.NewSignOnPolicyActionIDPAllOfIdentityProvider(vp["identity_provider_id"].(string)),
		)

		if v1, ok := d.GetOk("conditions"); ok {
			if vc, ok := v1.([]interface{}); ok && vc != nil && len(vc) > 0 && vc[0] != nil {
				var conditions *management.SignOnPolicyActionCommonConditionOrOrInner
				conditions, diags = expandSOPActionCondition(vc[0], management.ENUMSIGNONPOLICYTYPE_IDENTITY_PROVIDER, sopPriority)
				sopActionType.SetCondition(*conditions)
			}
		}

		if v1, ok := d.GetOk("registration_local_population_id"); ok && v1 != "" {
			obj := *management.NewSignOnPolicyActionIDPAllOfRegistration(true)
			obj.SetPopulation(*management.NewSignOnPolicyActionLoginAllOfRegistrationPopulation(v1.(string)))
			sopActionType.SetRegistration(obj)
		}

		if v1, ok := vp["acr_values"].(string); ok && v1 != "" {
			sopActionType.SetAcrValues(v1)
		}

		if v1, ok := vp["pass_user_context"].(bool); ok {
			sopActionType.SetPassUserContext(v1)
		}

		return sopActionType, diags

	}

	diags = append(diags, diag.Diagnostic{
		Severity: diag.Error,
		Summary:  "Block `identity_provider` with `identity_provider_id` must be defined when using the identity provider action type",
	})

	return nil, diags
}

func expandSOPActionLogin(d *schema.ResourceData, sopPriority int32) (*management.SignOnPolicyActionLogin, diag.Diagnostics) {
	var diags diag.Diagnostics

	sopActionType := management.NewSignOnPolicyActionLogin(
		sopPriority,
		management.ENUMSIGNONPOLICYTYPE_LOGIN,
	)

	if v, ok := d.GetOk("conditions"); ok {
		if v1, ok := v.([]interface{}); ok && v1 != nil && len(v1) > 0 && v1[0] != nil {
			var conditions *management.SignOnPolicyActionCommonConditionOrOrInner
			conditions, diags = expandSOPActionCondition(v1[0], management.ENUMSIGNONPOLICYTYPE_LOGIN, sopPriority)
			sopActionType.SetCondition(*conditions)
		}
	}

	if v, ok := d.GetOk("registration_external_href"); ok && v != "" {
		obj := *management.NewSignOnPolicyActionLoginAllOfRegistration(false)
		obj.SetExternal(*management.NewSignOnPolicyActionLoginAllOfRegistrationExternal(v.(string)))
		sopActionType.SetRegistration(obj)
	}

	if v, ok := d.GetOk("registration_local_population_id"); ok && v != "" {
		obj := *management.NewSignOnPolicyActionLoginAllOfRegistration(true)
		obj.SetPopulation(*management.NewSignOnPolicyActionLoginAllOfRegistrationPopulation(v.(string)))
		sopActionType.SetRegistration(obj)
	}

	socialIDSet := false
	if v, ok := d.GetOk("social_provider_ids"); ok {
		if vc, ok := v.([]interface{}); ok && vc != nil && len(vc) > 0 && vc[0] != "" {
			obj := make([]string, 0)
			for _, str := range vc {
				obj = append(obj, str.(string))
			}
			sopActionType.SetSocialProviders(expandSOPActionSocialProviders(obj))
			socialIDSet = true
		}
	}

	if v, ok := d.GetOk("enforce_lockout_for_identity_providers"); ok {
		if socialIDSet {
			sopActionType.SetEnforceLockoutForIdentityProviders(v.(bool))
		} else {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "`enforce_lockout_for_identity_providers`, where `social_provider_ids` is not set, has no effect.",
			})
		}
	}

	// block is optional
	if v, ok := d.Get("login").([]interface{}); ok && v != nil && len(v) > 0 && v[0] != nil {
		vp := v[0].(map[string]interface{})

		if v1, okJ := vp["recovery_enabled"].(bool); okJ {
			sopActionType.SetRecovery(*management.NewSignOnPolicyActionLoginAllOfRecovery(v1))
		} else {
			sopActionType.SetRecovery(*management.NewSignOnPolicyActionLoginAllOfRecovery(true))
		}

	}

	return sopActionType, diags
}

func expandSOPActionMFA(d *schema.ResourceData, sopPriority int32) (*management.SignOnPolicyActionMFA, diag.Diagnostics) {
	var diags diag.Diagnostics

	sopActionType := management.NewSignOnPolicyActionMFA(
		sopPriority,
		management.ENUMSIGNONPOLICYTYPE_MULTI_FACTOR_AUTHENTICATION,
	)

	if v, ok := d.GetOk("conditions"); ok {
		if vc, ok := v.([]interface{}); ok && vc != nil && len(vc) > 0 && vc[0] != nil {
			var conditions *management.SignOnPolicyActionCommonConditionOrOrInner
			conditions, diags = expandSOPActionCondition(vc[0], management.ENUMSIGNONPOLICYTYPE_MULTI_FACTOR_AUTHENTICATION, sopPriority)
			sopActionType.SetCondition(*conditions)
		}
	}

	if v, ok := d.Get("mfa").([]interface{}); ok && v != nil && len(v) > 0 && v[0] != nil {
		vp := v[0].(map[string]interface{})

		if v1, ok := vp["device_sign_on_policy_id"].(string); ok && v1 != "" {
			sopActionType.SetDeviceAuthenticationPolicy(*management.NewSignOnPolicyActionMFAAllOfDeviceAuthenticationPolicy(v1))
		}

		if v1, ok := vp["no_device_mode"].(string); ok && v1 != "" {
			sopActionType.SetNoDevicesMode(management.EnumSignOnPolicyNoDeviceMode(v1))
		}

	}

	return sopActionType, diags
}

func expandSOPActionProgressiveProfiling(d *schema.ResourceData, sopPriority int32) (*management.SignOnPolicyActionProgressiveProfiling, diag.Diagnostics) {
	var diags diag.Diagnostics

	if v, ok := d.Get("progressive_profiling").([]interface{}); ok && v != nil && len(v) > 0 && v[0] != nil {
		vp := v[0].(map[string]interface{})

		sopActionType := management.NewSignOnPolicyActionProgressiveProfiling(
			sopPriority,
			management.ENUMSIGNONPOLICYTYPE_PROGRESSIVE_PROFILING,
			expandSOPActionAttributes(vp["attribute"].(*schema.Set)),
			vp["prevent_multiple_prompts_per_flow"].(bool),
			int32(vp["prompt_interval_seconds"].(int)),
			vp["prompt_text"].(string),
		)

		if v, ok := d.GetOk("conditions"); ok {
			if vc, ok := v.([]interface{}); ok && vc != nil && len(vc) > 0 && vc[0] != nil {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Block `conditions` has no effect when using the progressive profiling action type",
				})
			}
		}

		return sopActionType, diags

	}

	diags = append(diags, diag.Diagnostic{
		Severity: diag.Error,
		Summary:  "Block `progressive_profiling` with `prompt_text` must be defined when using the progressive profiling action type",
	})

	return nil, diags
}

func expandSOPActionDiscoveryRules(items []interface{}) []management.SignOnPolicyActionIDFirstAllOfDiscoveryRules {

	var rules []management.SignOnPolicyActionIDFirstAllOfDiscoveryRules

	for _, item := range items {

		condition := item.(map[string]interface{})["attribute_contains_text"]
		conditionObj := *management.NewSignOnPolicyActionIDFirstAllOfCondition(condition.(string), "${identifier}")

		identityProviderObj := *management.NewSignOnPolicyActionIDFirstAllOfIdentityProvider(item.(map[string]interface{})["identity_provider_id"].(string))

		rules = append(rules, *management.NewSignOnPolicyActionIDFirstAllOfDiscoveryRules(conditionObj, identityProviderObj))

	}

	return rules

}

func expandSOPActionAttributes(items *schema.Set) []management.SignOnPolicyActionProgressiveProfilingAllOfAttributes {

	var attributes []management.SignOnPolicyActionProgressiveProfilingAllOfAttributes

	for _, item := range items.List() {

		attributes = append(attributes, *management.NewSignOnPolicyActionProgressiveProfilingAllOfAttributes(item.(map[string]interface{})["name"].(string), item.(map[string]interface{})["required"].(bool)))

	}

	return attributes

}

func expandSOPActionSocialProviders(items []string) []management.SignOnPolicyActionLoginAllOfSocialProviders {

	var socialProviders []management.SignOnPolicyActionLoginAllOfSocialProviders

	for _, item := range items {

		socialProviders = append(socialProviders, *management.NewSignOnPolicyActionLoginAllOfSocialProviders(item))

	}

	return socialProviders

}

func expandSOPActionCondition(condition interface{}, actionType management.EnumSignOnPolicyType, sopPriority int32) (*management.SignOnPolicyActionCommonConditionOrOrInner, diag.Diagnostics) {

	var sopConditions *management.SignOnPolicyActionCommonConditionOrOrInner
	var diags diag.Diagnostics

	switch actionType {
	case management.ENUMSIGNONPOLICYTYPE_IDENTIFIER_FIRST:
		sopConditions, diags = expandSOPActionConditionIDFirst(condition, sopPriority)
	case management.ENUMSIGNONPOLICYTYPE_IDENTITY_PROVIDER:
		sopConditions = expandSOPActionConditionIDP(condition)
	case management.ENUMSIGNONPOLICYTYPE_LOGIN:
		sopConditions, diags = expandSOPActionConditionLogin(condition, sopPriority)
	case management.ENUMSIGNONPOLICYTYPE_MULTI_FACTOR_AUTHENTICATION:
		sopConditions, diags = expandSOPActionConditionMFA(condition, sopPriority)
	default:
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Policy action %s not supported when evaluating condition block", actionType),
		})
		return nil, diags
	}

	return sopConditions, diags
}

func buildSignOnOlderThan(v int32) management.SignOnPolicyActionCommonConditionOrOrInner {
	return management.SignOnPolicyActionCommonConditionOrOrInner{
		SignOnPolicyActionCommonConditionAggregate: &management.SignOnPolicyActionCommonConditionAggregate{
			SignOnPolicyActionCommonConditionGreater: management.NewSignOnPolicyActionCommonConditionGreater(v, "${session.lastSignOn.withAuthenticator.pwd.at}"),
		},
	}
}

type attributeEquality struct {
	attributeReference string
	attributeValue     string
}

func buildUserMemberOfPopulation(v []interface{}, sopPriority int32) (*management.SignOnPolicyActionCommonConditionOrOrInner, diag.Diagnostics) {

	attributeList := make([]attributeEquality, 0)

	for _, population := range v {
		attributeList = append(attributeList, attributeEquality{
			attributeReference: "${user.population.id}",
			attributeValue:     population.(string),
		})
	}

	return buildAttributeEqualsCondition(attributeList, "user_is_member_of_any_population_id", sopPriority)

}

func buildUserAttributes(v []interface{}, sopPriority int32) (*management.SignOnPolicyActionCommonConditionOrOrInner, diag.Diagnostics) {
	attributeList := make([]attributeEquality, 0)

	for _, attributeMap := range v {
		attributeList = append(attributeList, attributeEquality{
			attributeReference: attributeMap.(map[string]interface{})["attribute_reference"].(string),
			attributeValue:     attributeMap.(map[string]interface{})["value"].(string),
		})
	}

	return buildAttributeEqualsCondition(attributeList, "user_attribute_equals", sopPriority)

}

func buildAttributeEqualsCondition(v []attributeEquality, tfSchemaAttribute string, sopPriority int32) (*management.SignOnPolicyActionCommonConditionOrOrInner, diag.Diagnostics) {
	var diags diag.Diagnostics

	if sopPriority < 2 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  fmt.Sprintf("Condition `%s` is defined but has no effect where the action priority is 1.", tfSchemaAttribute),
		})
	} else {

		if len(v) > 1 {

			conditionList := make([]management.SignOnPolicyActionCommonConditionOrOrInner, 0)

			for _, attribute := range v {

				conditionList = append(conditionList, management.SignOnPolicyActionCommonConditionOrOrInner{
					SignOnPolicyActionCommonConditionAggregate: &management.SignOnPolicyActionCommonConditionAggregate{
						SignOnPolicyActionCommonConditionEquals: management.NewSignOnPolicyActionCommonConditionEquals(attribute.attributeReference, attribute.attributeValue),
					},
				})

			}

			return &management.SignOnPolicyActionCommonConditionOrOrInner{
				SignOnPolicyActionCommonConditionOr: &management.SignOnPolicyActionCommonConditionOr{
					Or: conditionList,
				},
			}, nil

		} else {

			return &management.SignOnPolicyActionCommonConditionOrOrInner{
				SignOnPolicyActionCommonConditionAggregate: &management.SignOnPolicyActionCommonConditionAggregate{
					SignOnPolicyActionCommonConditionEquals: management.NewSignOnPolicyActionCommonConditionEquals(v[0].attributeReference, v[0].attributeValue),
				},
			}, nil
		}
	}

	return nil, diags

}

func expandSOPActionConditionIDFirstAndLogin(condition interface{}, sopPriority int32) (*management.SignOnPolicyActionCommonConditionOrOrInner, diag.Diagnostics) {
	var diags diag.Diagnostics

	var conditionStructList = make([]management.SignOnPolicyActionCommonConditionOrOrInner, 0)

	if v, ok := condition.(map[string]interface{})["last_sign_on_older_than_seconds"].(int); ok {

		conditionStructList = append(conditionStructList, buildSignOnOlderThan(int32(v)))
	}

	if v, ok := condition.(map[string]interface{})["user_is_member_of_any_population_id"].([]interface{}); ok && v != nil && len(v) > 0 && v[0] != "" {

		var populationStructs *management.SignOnPolicyActionCommonConditionOrOrInner
		populationStructs, diags = buildUserMemberOfPopulation(v, sopPriority)
		if diags.HasError() {
			return nil, diags
		}

		conditionStructList = append(conditionStructList, *populationStructs)

	}

	if v, ok := condition.(map[string]interface{})["user_attribute_equals"].(*schema.Set); ok && v != nil && len(v.List()) > 0 && v.List()[0] != nil {

		var attributeStructs *management.SignOnPolicyActionCommonConditionOrOrInner
		attributeStructs, diags = buildUserAttributes(v.List(), sopPriority)
		if diags.HasError() {
			return nil, diags
		}

		conditionStructList = append(conditionStructList, *attributeStructs)

	}

	var conditionStruct = &management.SignOnPolicyActionCommonConditionOrOrInner{}

	if len(conditionStructList) > 1 {

		conditionStruct = &management.SignOnPolicyActionCommonConditionOrOrInner{
			SignOnPolicyActionCommonConditionOr: &management.SignOnPolicyActionCommonConditionOr{
				Or: conditionStructList,
			},
		}

	} else if len(conditionStructList) == 1 {
		conditionStruct = &conditionStructList[0]
	}

	return conditionStruct, diags
}

func expandSOPActionConditionIDFirst(condition interface{}, sopPriority int32) (*management.SignOnPolicyActionCommonConditionOrOrInner, diag.Diagnostics) {
	return expandSOPActionConditionIDFirstAndLogin(condition, sopPriority)
}

func expandSOPActionConditionIDP(condition interface{}) *management.SignOnPolicyActionCommonConditionOrOrInner {

	conditionStruct := &management.SignOnPolicyActionCommonConditionOrOrInner{}

	if v, ok := condition.(map[string]interface{})["last_sign_on_older_than_seconds"].(int); ok {

		conditionStruct.SignOnPolicyActionCommonConditionAggregate = buildSignOnOlderThan(int32(v)).SignOnPolicyActionCommonConditionAggregate
	}

	return conditionStruct

}

func expandSOPActionConditionLogin(condition interface{}, sopPriority int32) (*management.SignOnPolicyActionCommonConditionOrOrInner, diag.Diagnostics) {
	return expandSOPActionConditionIDFirstAndLogin(condition, sopPriority)
}

func expandSOPActionConditionMFA(condition interface{}, sopPriority int32) (*management.SignOnPolicyActionCommonConditionOrOrInner, diag.Diagnostics) {
	var diags diag.Diagnostics

	// heres johnny
	var conditionStructList = make([]management.SignOnPolicyActionCommonConditionOrOrInner, 0)

	if v, ok := condition.(map[string]interface{})["last_sign_on_older_than_seconds"].(int); ok {

		conditionStructList = append(conditionStructList, buildSignOnOlderThan(int32(v)))
	}

	if v, ok := condition.(map[string]interface{})["ip_out_of_range_cidr"].([]interface{}); ok && v != nil && len(v) > 0 && v[0] != "" {

		obj := make([]string, 0)
		for _, str := range v {
			obj = append(obj, str.(string))
		}

		conditionStructList = append(conditionStructList, management.SignOnPolicyActionCommonConditionOrOrInner{
			SignOnPolicyActionCommonConditionNot: &management.SignOnPolicyActionCommonConditionNot{
				Not: &management.SignOnPolicyActionCommonConditionAggregate{
					SignOnPolicyActionCommonConditionIPRange: management.NewSignOnPolicyActionCommonConditionIPRange("${flow.request.http.remoteIp}", obj),
				},
			},
		})

	}

	if v, ok := condition.(map[string]interface{})["ip_reputation_high_risk"].(bool); ok && v {

		min := 80
		max := 100
		ipRisk := *management.NewSignOnPolicyActionCommonConditionIPRiskIpRisk(int32(min), int32(max))

		conditionStructList = append(conditionStructList, management.SignOnPolicyActionCommonConditionOrOrInner{
			SignOnPolicyActionCommonConditionAggregate: &management.SignOnPolicyActionCommonConditionAggregate{
				SignOnPolicyActionCommonConditionIPRisk: management.NewSignOnPolicyActionCommonConditionIPRisk(ipRisk, "${flow.request.http.remoteIp}"),
			},
		})
	}

	if v, ok := condition.(map[string]interface{})["geovelocity_anomaly_detected"].(bool); ok && v {

		validObj := *management.NewSignOnPolicyActionCommonConditionGeovelocityValid("${user.lastSignOn.remoteIp}", "${user.lastSignOn.at}")

		conditionStructList = append(conditionStructList, management.SignOnPolicyActionCommonConditionOrOrInner{
			SignOnPolicyActionCommonConditionAggregate: &management.SignOnPolicyActionCommonConditionAggregate{
				SignOnPolicyActionCommonConditionGeovelocity: management.NewSignOnPolicyActionCommonConditionGeovelocity("${flow.request.http.remoteIp}", validObj),
			},
		})
	}

	if v, ok := condition.(map[string]interface{})["anonymous_network_detected"].(bool); ok && v {

		anonymousNetworkList := make([]string, 0)

		if j, ok := condition.(map[string]interface{})["anonymous_network_detected_allowed_cidr"].([]interface{}); ok && j != nil && len(j) > 0 && j[0] != "" {
			obj := make([]string, 0)
			for _, str := range j {
				obj = append(obj, str.(string))
			}
			anonymousNetworkList = obj
		}

		conditionStructList = append(conditionStructList, management.SignOnPolicyActionCommonConditionOrOrInner{
			SignOnPolicyActionCommonConditionAggregate: &management.SignOnPolicyActionCommonConditionAggregate{
				SignOnPolicyActionCommonConditionAnonymousNetwork: management.NewSignOnPolicyActionCommonConditionAnonymousNetwork(anonymousNetworkList, "${flow.request.http.remoteIp}"),
			},
		})
	}

	if v, ok := condition.(map[string]interface{})["user_is_member_of_any_population_id"].([]interface{}); ok && v != nil && len(v) > 0 && v[0] != "" {

		var populationStructs *management.SignOnPolicyActionCommonConditionOrOrInner
		populationStructs, diags = buildUserMemberOfPopulation(v, sopPriority)
		if diags.HasError() {
			return nil, diags
		}

		conditionStructList = append(conditionStructList, *populationStructs)
	}

	if v, ok := condition.(map[string]interface{})["user_attribute_equals"].(*schema.Set); ok && v != nil && len(v.List()) > 0 && v.List()[0] != nil {

		var attributeStructs *management.SignOnPolicyActionCommonConditionOrOrInner
		attributeStructs, diags = buildUserAttributes(v.List(), sopPriority)
		if diags.HasError() {
			return nil, diags
		}

		conditionStructList = append(conditionStructList, *attributeStructs)
	}

	var conditionStruct = &management.SignOnPolicyActionCommonConditionOrOrInner{}

	if len(conditionStructList) > 1 {

		conditionStruct = &management.SignOnPolicyActionCommonConditionOrOrInner{
			SignOnPolicyActionCommonConditionOr: &management.SignOnPolicyActionCommonConditionOr{
				Or: conditionStructList,
			},
		}

	} else if len(conditionStructList) == 1 {
		conditionStruct = &conditionStructList[0]
	}

	return conditionStruct, diags
}

type flattenedConditions struct {
	last_sign_on_older_than_seconds         *int32
	user_is_member_of_any_population_id     []string
	user_attribute_equals                   []attributeEquality
	ip_out_of_range_cidr                    []string
	ip_reputation_high_risk                 *bool
	geovelocity_anomaly_detected            *bool
	anonymous_network_detected              *bool
	anonymous_network_detected_allowed_cidr []string
}

func processConditions(conditions *flattenedConditions, v management.SignOnPolicyActionCommonConditionOrOrInner) (*flattenedConditions, diag.Diagnostics) {
	var diags diag.Diagnostics

	returnCondition := conditions

	if v.SignOnPolicyActionCommonConditionAnd != nil {
		// AND doesn't feature in the conditions set by the UI
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unexpected condition `AND` while processing conditions. This is not supported in the provider.  Please raise an issue.",
		})

		return nil, diags
	}

	if j := v.SignOnPolicyActionCommonConditionNot; j != nil {
		// This is only the IP CIDR range rule

		if vc := j.GetNot().SignOnPolicyActionCommonConditionIPRange; vc != nil {

			if condition, ok := vc.GetContainsOk(); ok {
				if *condition != "${flow.request.http.remoteIp}" {
					diags = append(diags, diag.Diagnostic{
						Severity: diag.Error,
						Summary:  fmt.Sprintf("Condition `ip_out_of_range_cidr` has unknown field %s, but expecting value ${flow.request.http.remoteIp}.  This is not supported in the provider.  Please raise an issue.", *condition),
					})

					return nil, diags
				}
			}

			if returnCondition.ip_out_of_range_cidr != nil {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Condition `ip_out_of_range_cidr` has has multiple nested values.  This is not supported in the provider.  Please raise an issue.",
				})

				return nil, diags
			}
			returnCondition.ip_out_of_range_cidr = vc.GetIpRange()
		}

	}

	if j := v.SignOnPolicyActionCommonConditionOr; j != nil {
		for _, orCondition := range j.GetOr() {
			returnCondition, diags = processConditions(returnCondition, orCondition)
			if diags.HasError() {
				return nil, diags
			}
		}
	}

	if j := v.SignOnPolicyActionCommonConditionAggregate; j != nil {
		if vc := j.SignOnPolicyActionCommonConditionGreater; vc != nil {

			if condition, ok := vc.GetSecondsSinceOk(); ok {
				if *condition != "${session.lastSignOn.withAuthenticator.pwd.at}" {
					diags = append(diags, diag.Diagnostic{
						Severity: diag.Error,
						Summary:  fmt.Sprintf("Condition `last_sign_on_older_than_seconds` has unknown field %s, but expecting value ${session.lastSignOn.withAuthenticator.pwd.at}.  This is not supported in the provider.  Please raise an issue.", *condition),
					})

					return nil, diags
				}

				if returnCondition.last_sign_on_older_than_seconds != nil {
					diags = append(diags, diag.Diagnostic{
						Severity: diag.Error,
						Summary:  "Condition `last_sign_on_older_than_seconds` has has multiple nested values.  This is not supported in the provider.  Please raise an issue.",
					})

					return nil, diags
				}
				returnCondition.last_sign_on_older_than_seconds = func() *int32 { b := vc.GetGreater(); return &b }()
			}
		}

		if vc := j.SignOnPolicyActionCommonConditionIPRange; vc != nil {

			if condition, ok := vc.GetContainsOk(); ok {
				if *condition != "${flow.request.http.remoteIp}" {
					diags = append(diags, diag.Diagnostic{
						Severity: diag.Error,
						Summary:  fmt.Sprintf("Condition `ip_out_of_range_cidr` has unknown field %s, but expecting value ${flow.request.http.remoteIp}.  This is not supported in the provider.  Please raise an issue.", *condition),
					})

					return nil, diags
				}

				if returnCondition.ip_out_of_range_cidr != nil {
					diags = append(diags, diag.Diagnostic{
						Severity: diag.Error,
						Summary:  "Condition `ip_out_of_range_cidr` has has multiple nested values.  This is not supported in the provider.  Please raise an issue.",
					})

					return nil, diags
				}

				returnCondition.ip_out_of_range_cidr = vc.GetIpRange()
			}
		}

		if vc := j.SignOnPolicyActionCommonConditionIPRisk; vc != nil {

			if condition, ok := vc.GetValidOk(); ok {
				if *condition != "${flow.request.http.remoteIp}" {
					diags = append(diags, diag.Diagnostic{
						Severity: diag.Error,
						Summary:  fmt.Sprintf("Condition `ip_reputation_high_risk` has unknown field %s, but expecting value ${flow.request.http.remoteIp}.  This is not supported in the provider.  Please raise an issue.", *condition),
					})

					return nil, diags
				}

				if vc.GetIpRisk().MaxScore != 100 || vc.GetIpRisk().MinScore != 80 {
					diags = append(diags, diag.Diagnostic{
						Severity: diag.Error,
						Summary:  fmt.Sprintf("Condition `ip_reputation_high_risk` has unknown min and max scores of %d and %d, but expecting values 80 and 100.  This is not supported in the provider.  Please raise an issue.", vc.GetIpRisk().MinScore, vc.GetIpRisk().MaxScore),
					})

					return nil, diags
				}

				if returnCondition.ip_reputation_high_risk != nil {
					diags = append(diags, diag.Diagnostic{
						Severity: diag.Error,
						Summary:  "Condition `ip_reputation_high_risk` has has multiple nested values.  This is not supported in the provider.  Please raise an issue.",
					})

					return nil, diags
				}

				returnCondition.ip_reputation_high_risk = func() *bool { b := true; return &b }()

			}
		}

		if vc := j.SignOnPolicyActionCommonConditionGeovelocity; vc != nil {

			if condition, ok := vc.GetGeoVelocityOk(); ok {
				if *condition != "${flow.request.http.remoteIp}" {
					diags = append(diags, diag.Diagnostic{
						Severity: diag.Error,
						Summary:  fmt.Sprintf("Condition `geovelocity_anomaly_detected` has unknown field %s, but expecting value ${flow.request.http.remoteIp}.  This is not supported in the provider.  Please raise an issue.", *condition),
					})

					return nil, diags
				}

				if vc.GetValid().PreviousSuccessfulAuthenticationIp != "${user.lastSignOn.remoteIp}" {
					diags = append(diags, diag.Diagnostic{
						Severity: diag.Error,
						Summary:  fmt.Sprintf("Condition `geovelocity_anomaly_detected` has unknown field %s, but expecting value ${user.lastSignOn.remoteIp}.  This is not supported in the provider.  Please raise an issue.", *condition),
					})

					return nil, diags
				}

				if vc.GetValid().PreviousSuccessfulAuthenticationTime != "${user.lastSignOn.at}" {
					diags = append(diags, diag.Diagnostic{
						Severity: diag.Error,
						Summary:  fmt.Sprintf("Condition `geovelocity_anomaly_detected` has unknown field %s, but expecting value ${user.lastSignOn.at}.  This is not supported in the provider.  Please raise an issue.", *condition),
					})

					return nil, diags
				}

				if returnCondition.geovelocity_anomaly_detected != nil {
					diags = append(diags, diag.Diagnostic{
						Severity: diag.Error,
						Summary:  "Condition `geovelocity_anomaly_detected` has has multiple nested values.  This is not supported in the provider.  Please raise an issue.",
					})

					return nil, diags
				}

				returnCondition.geovelocity_anomaly_detected = func() *bool { b := true; return &b }()
			}
		}

		if vc := j.SignOnPolicyActionCommonConditionAnonymousNetwork; vc != nil {

			if condition, ok := vc.GetAnonymousNetworkOk(); ok {
				returnCondition.anonymous_network_detected_allowed_cidr = condition
			}

			if condition, ok := vc.GetValidOk(); ok {
				if *condition != "${flow.request.http.remoteIp}" {
					diags = append(diags, diag.Diagnostic{
						Severity: diag.Error,
						Summary:  fmt.Sprintf("Condition `anonymous_network_detected` has unknown field %s, but expecting value ${flow.request.http.remoteIp}.  This is not supported in the provider.  Please raise an issue.", *condition),
					})

					return nil, diags
				}

				if returnCondition.anonymous_network_detected != nil {
					diags = append(diags, diag.Diagnostic{
						Severity: diag.Error,
						Summary:  "Condition `anonymous_network_detected` has has multiple nested values.  This is not supported in the provider.  Please raise an issue.",
					})

					return nil, diags
				}

				returnCondition.anonymous_network_detected = func() *bool { b := true; return &b }()
			}
		}

		if vc := j.SignOnPolicyActionCommonConditionEquals; vc != nil {

			if populationField := vc.GetValue(); populationField == "${user.population.id}" {
				returnCondition.user_is_member_of_any_population_id = append(returnCondition.user_is_member_of_any_population_id, vc.GetEquals())

			} else {
				condition := attributeEquality{
					attributeValue:     vc.GetEquals(),
					attributeReference: vc.GetValue(),
				}

				returnCondition.user_attribute_equals = append(returnCondition.user_attribute_equals, condition)

			}

		}
	}

	return returnCondition, diags

}

func flattenConditions(signOnPolicyActionCommonConditions management.SignOnPolicyActionCommonConditionOrOrInner) ([]interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	var conditionStruct *management.SignOnPolicyActionCommonConditionOrOrInner

	// If condition is in the aggregate, then add the only condition to the or list to process uniformly later
	if v := signOnPolicyActionCommonConditions.SignOnPolicyActionCommonConditionAggregate; v != nil {

		conditionList := make([]management.SignOnPolicyActionCommonConditionOrOrInner, 0, 1)

		conditionStruct = &management.SignOnPolicyActionCommonConditionOrOrInner{
			SignOnPolicyActionCommonConditionOr: &management.SignOnPolicyActionCommonConditionOr{
				Or: append(conditionList, signOnPolicyActionCommonConditions),
			},
		}
	} else {
		conditionStruct = &signOnPolicyActionCommonConditions
	}

	conditions, diags := processConditions(&flattenedConditions{}, *conditionStruct)

	flattenedConditions := map[string]interface{}{
		"last_sign_on_older_than_seconds":         nil,
		"user_is_member_of_any_population_id":     nil,
		"user_attribute_equals":                   nil,
		"ip_out_of_range_cidr":                    nil,
		"ip_reputation_high_risk":                 nil,
		"geovelocity_anomaly_detected":            nil,
		"anonymous_network_detected":              nil,
		"anonymous_network_detected_allowed_cidr": nil,
	}

	if conditions.last_sign_on_older_than_seconds != nil {
		flattenedConditions["last_sign_on_older_than_seconds"] = *conditions.last_sign_on_older_than_seconds
	}

	if conditions.user_is_member_of_any_population_id != nil {
		flattenedConditions["user_is_member_of_any_population_id"] = conditions.user_is_member_of_any_population_id
	}

	if conditions.user_attribute_equals != nil {
		attributeList := make([]map[string]interface{}, 0)
		for _, attributeStruct := range conditions.user_attribute_equals {
			attributeList = append(attributeList, map[string]interface{}{
				"attribute_reference": attributeStruct.attributeReference,
				"value":               attributeStruct.attributeValue,
			})
		}
		flattenedConditions["user_attribute_equals"] = attributeList
	}

	if conditions.ip_out_of_range_cidr != nil {
		flattenedConditions["ip_out_of_range_cidr"] = conditions.ip_out_of_range_cidr
	}

	if conditions.ip_reputation_high_risk != nil {
		flattenedConditions["ip_reputation_high_risk"] = *conditions.ip_reputation_high_risk
	}

	if conditions.geovelocity_anomaly_detected != nil {
		flattenedConditions["geovelocity_anomaly_detected"] = *conditions.geovelocity_anomaly_detected
	}

	if conditions.anonymous_network_detected != nil {
		flattenedConditions["anonymous_network_detected"] = *conditions.anonymous_network_detected
	}

	if conditions.anonymous_network_detected_allowed_cidr != nil {
		flattenedConditions["anonymous_network_detected_allowed_cidr"] = conditions.anonymous_network_detected_allowed_cidr
	}

	conditionsList := make([]interface{}, 0, 1)
	return append(conditionsList, flattenedConditions), diags
}

func flattenActionProgressiveProfiling(signOnPolicyActionProgressiveProfiling *management.SignOnPolicyActionProgressiveProfiling) []interface{} {
	actionList := make([]interface{}, 0, 1)

	action := map[string]interface{}{
		"attribute":                         flattenActionProgressiveProfilingAttributes(signOnPolicyActionProgressiveProfiling.GetAttributes()),
		"prevent_multiple_prompts_per_flow": signOnPolicyActionProgressiveProfiling.GetPreventMultiplePromptsPerFlow(),
		"prompt_interval_seconds":           signOnPolicyActionProgressiveProfiling.GetPromptIntervalSeconds(),
		"prompt_text":                       signOnPolicyActionProgressiveProfiling.GetPromptText(),
	}

	return append(actionList, action)
}

func flattenActionProgressiveProfilingAttributes(signOnPolicyActionProgressiveProfilingAllOfAttributes []management.SignOnPolicyActionProgressiveProfilingAllOfAttributes) []interface{} {
	attributes := make([]interface{}, 0, len(signOnPolicyActionProgressiveProfilingAllOfAttributes))
	for _, attribute := range signOnPolicyActionProgressiveProfilingAllOfAttributes {

		name := attribute.GetName()
		required := attribute.GetRequired()

		attributes = append(attributes, map[string]interface{}{
			"name":     name,
			"required": required,
		})
	}
	return attributes
}

func flattenActionMFA(signOnPolicyActionMFA *management.SignOnPolicyActionMFA) []interface{} {
	actionList := make([]interface{}, 0, 1)

	action := map[string]interface{}{
		"device_sign_on_policy_id": signOnPolicyActionMFA.DeviceAuthenticationPolicy.GetId(),
	}

	if v, ok := signOnPolicyActionMFA.GetNoDevicesModeOk(); ok {
		action["no_device_mode"] = v
	}

	return append(actionList, action)
}

func flattenActionLogin(signOnPolicyActionLogin *management.SignOnPolicyActionLogin) []interface{} {
	actionList := make([]interface{}, 0, 1)

	action := map[string]interface{}{}

	if v, ok := signOnPolicyActionLogin.GetRecoveryOk(); ok {
		action["recovery_enabled"] = v.GetEnabled()
	} else {
		action["recovery_enabled"] = nil
	}

	return append(actionList, action)
}

func flattenActionIDP(signOnPolicyActionIDP *management.SignOnPolicyActionIDP) []interface{} {
	actionList := make([]interface{}, 0, 1)

	action := map[string]interface{}{
		"identity_provider_id": signOnPolicyActionIDP.IdentityProvider.GetId(),
	}

	if v, ok := signOnPolicyActionIDP.GetAcrValuesOk(); ok {
		action["acr_values"] = v
	}

	if v, ok := signOnPolicyActionIDP.GetPassUserContextOk(); ok {
		action["pass_user_context"] = v
	}

	return append(actionList, action)
}

func flattenActionIDFirst(signOnPolicyActionIDFirst *management.SignOnPolicyActionIDFirst) ([]interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	actionList := make([]interface{}, 0, 1)

	action := map[string]interface{}{}

	if v, ok := signOnPolicyActionIDFirst.GetDiscoveryRulesOk(); ok {
		var discoveryRules []interface{}
		discoveryRules, diags = flattenDiscoveryRulesInner(v)
		if diags.HasError() {
			return nil, diags
		}
		action["discovery_rule"] = discoveryRules
	}

	if v, ok := signOnPolicyActionIDFirst.GetRecoveryOk(); ok {
		action["recovery_enabled"] = v.GetEnabled()
	}

	return append(actionList, action), diags
}

func flattenActionAgreement(signOnPolicyActionAgreement *management.SignOnPolicyActionAgreement) []interface{} {
	actionList := make([]interface{}, 0, 1)

	action := map[string]interface{}{
		"agreement_id": signOnPolicyActionAgreement.Agreement.GetId(),
	}

	if v, ok := signOnPolicyActionAgreement.GetDisableDeclineOptionOk(); ok {
		action["show_decline_option"] = !*v
	}

	return append(actionList, action)
}

func flattenDiscoveryRulesInner(signOnPolicyActionIDFirstAllOfDiscoveryRules []management.SignOnPolicyActionIDFirstAllOfDiscoveryRules) ([]interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	rules := make([]interface{}, 0, len(signOnPolicyActionIDFirstAllOfDiscoveryRules))
	for _, rule := range signOnPolicyActionIDFirstAllOfDiscoveryRules {

		condition := rule.GetCondition()

		if condition.GetValue() != "${identifier}" {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "`discovery_rule` has unknown field %s, but expecting value ${identifier}.  This is not supported in the provider.  Please raise an issue.",
			})

			return nil, diags
		}

		rules = append(rules, map[string]interface{}{
			"attribute_contains_text": condition.GetContains(),
			"identity_provider_id":    rule.GetIdentityProvider().Id,
		})
	}
	return rules, diags
}

func flattenActionSocialProvidersInner(signOnPolicyActionLoginAllOfSocialProviders []management.SignOnPolicyActionLoginAllOfSocialProviders) []string {
	providerList := make([]string, 0, len(signOnPolicyActionLoginAllOfSocialProviders))

	for _, provider := range signOnPolicyActionLoginAllOfSocialProviders {
		providerList = append(providerList, provider.GetId())
	}

	if len(providerList) == 0 {
		providerList = nil
	}

	return providerList
}

func getActionID(instance management.SignOnPolicyAction) string {
	var actionID string
	switch instance.GetActualInstance().(type) {
	case *management.SignOnPolicyActionLogin:
		actionID = instance.SignOnPolicyActionLogin.GetId()
	case *management.SignOnPolicyActionAgreement:
		actionID = instance.SignOnPolicyActionAgreement.GetId()
	case *management.SignOnPolicyActionIDFirst:
		actionID = instance.SignOnPolicyActionIDFirst.GetId()
	case *management.SignOnPolicyActionIDP:
		actionID = instance.SignOnPolicyActionIDP.GetId()
	case *management.SignOnPolicyActionProgressiveProfiling:
		actionID = instance.SignOnPolicyActionProgressiveProfiling.GetId()
	case *management.SignOnPolicyActionMFA:
		actionID = instance.SignOnPolicyActionMFA.GetId()
	}

	return actionID
}
