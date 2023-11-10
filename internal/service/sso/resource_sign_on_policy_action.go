package sso

import (
	"context"
	"fmt"
	"net/http"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
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

	var diags diag.Diagnostics

	signOnPolicyAction, diags := expandSOPAction(d)
	if diags.HasError() {
		return diags
	}

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := apiClient.SignOnPolicyActionsApi.CreateSignOnPolicyAction(ctx, d.Get("environment_id").(string), d.Get("sign_on_policy_id").(string)).SignOnPolicyAction(*signOnPolicyAction).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, apiClient, d.Get("environment_id").(string), fO, fR, fErr)
		},
		"CreateSignOnPolicyAction",
		customErrorSignOnPolicyActionCreateUpdate,
		sdk.DefaultCreateReadRetryable,
	)
	if diags.HasError() {
		return diags
	}

	respObject := resp.(*management.SignOnPolicyAction)

	d.SetId(getActionID(*respObject))

	return resourceSignOnPolicyActionRead(ctx, d, meta)
}

func resourceSignOnPolicyActionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient

	var diags diag.Diagnostics

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := apiClient.SignOnPolicyActionsApi.ReadOneSignOnPolicyAction(ctx, d.Get("environment_id").(string), d.Get("sign_on_policy_id").(string), d.Id()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, apiClient, d.Get("environment_id").(string), fO, fR, fErr)
		},
		"ReadOneSignOnPolicyAction",
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

	respObject := resp.(*management.SignOnPolicyAction)

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
		"pingid":                                 nil,
		"pingid_windows_login_passwordless":      nil,
	}

	switch respObject.GetActualInstance().(type) {
	case *management.SignOnPolicyActionLogin:

		values["priority"] = respObject.SignOnPolicyActionLogin.GetPriority()

		if v, ok := respObject.SignOnPolicyActionLogin.GetConditionOk(); ok {
			var conditions interface{}
			conditions, diags = flattenConditions(*v)
			if diags.HasError() {
				return diags
			}
			values["conditions"] = conditions
		}

		if v, ok := respObject.SignOnPolicyActionLogin.GetRegistrationOk(); ok {
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

		if v, ok := respObject.SignOnPolicyActionLogin.GetSocialProvidersOk(); ok {
			values["social_provider_ids"] = flattenActionSocialProvidersInner(v)
		}

		if v, ok := respObject.SignOnPolicyActionLogin.GetEnforceLockoutForIdentityProvidersOk(); ok {
			values["enforce_lockout_for_identity_providers"] = v
		}

		values["login"] = flattenActionLogin(respObject.SignOnPolicyActionLogin)

	case *management.SignOnPolicyActionAgreement:

		values["priority"] = respObject.SignOnPolicyActionAgreement.GetPriority()

		if v, ok := respObject.SignOnPolicyActionAgreement.GetConditionOk(); ok {
			var conditions interface{}
			conditions, diags = flattenConditions(*v)
			if diags.HasError() {
				return diags
			}
			values["conditions"] = conditions
		}

		values["agreement"] = flattenActionAgreement(respObject.SignOnPolicyActionAgreement)

	case *management.SignOnPolicyActionIDFirst:

		values["priority"] = respObject.SignOnPolicyActionIDFirst.GetPriority()

		if v, ok := respObject.SignOnPolicyActionIDFirst.GetConditionOk(); ok {
			var conditions interface{}
			conditions, diags = flattenConditions(*v)
			if diags.HasError() {
				return diags
			}
			values["conditions"] = conditions
		}

		if v, ok := respObject.SignOnPolicyActionIDFirst.GetRegistrationOk(); ok {
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

		if v, ok := respObject.SignOnPolicyActionIDFirst.GetSocialProvidersOk(); ok {
			values["social_provider_ids"] = flattenActionSocialProvidersInner(v)
		}

		if v, ok := respObject.SignOnPolicyActionIDFirst.GetEnforceLockoutForIdentityProvidersOk(); ok {
			values["enforce_lockout_for_identity_providers"] = v
		}

		var idFirst []interface{}
		idFirst, diags = flattenActionIDFirst(respObject.SignOnPolicyActionIDFirst)
		if diags.HasError() {
			return diags
		}
		values["identifier_first"] = idFirst

	case *management.SignOnPolicyActionIDP:

		values["priority"] = respObject.SignOnPolicyActionIDP.GetPriority()

		if v, ok := respObject.SignOnPolicyActionIDP.GetConditionOk(); ok {
			var conditions interface{}
			conditions, diags = flattenConditions(*v)
			if diags.HasError() {
				return diags
			}
			values["conditions"] = conditions
		}

		if v, ok := respObject.SignOnPolicyActionIDP.GetRegistrationOk(); ok {
			if v1, ok := v.GetPopulationOk(); ok {
				values["registration_local_population_id"] = v1.GetId()
			}

			if v1, ok := v.GetConfirmIdentityProviderAttributesOk(); ok {
				values["registration_confirm_user_attributes"] = v1
			}
		}

		values["identity_provider"] = flattenActionIDP(respObject.SignOnPolicyActionIDP)

	case *management.SignOnPolicyActionProgressiveProfiling:

		values["priority"] = respObject.SignOnPolicyActionProgressiveProfiling.GetPriority()

		if v, ok := respObject.SignOnPolicyActionProgressiveProfiling.GetConditionOk(); ok {
			var conditions interface{}
			conditions, diags = flattenConditions(*v)
			if diags.HasError() {
				return diags
			}
			values["conditions"] = conditions
		}

		values["progressive_profiling"] = flattenActionProgressiveProfiling(respObject.SignOnPolicyActionProgressiveProfiling)

	case *management.SignOnPolicyActionMFA:

		values["priority"] = respObject.SignOnPolicyActionMFA.GetPriority()

		if v, ok := respObject.SignOnPolicyActionMFA.GetConditionOk(); ok {
			var conditions interface{}
			conditions, diags = flattenConditions(*v)
			if diags.HasError() {
				return diags
			}
			values["conditions"] = conditions
		}

		values["mfa"] = flattenActionMFA(respObject.SignOnPolicyActionMFA)

	case *management.SignOnPolicyActionCommon:

		values["priority"] = respObject.SignOnPolicyActionCommon.GetPriority()

		values["pingid"] = make([]interface{}, 1)

	case *management.SignOnPolicyActionPingIDWinLoginPasswordless:

		values["priority"] = respObject.SignOnPolicyActionPingIDWinLoginPasswordless.GetPriority()

		values["pingid_windows_login_passwordless"] = flattenActionPingIDWinLoginPasswordless(respObject.SignOnPolicyActionPingIDWinLoginPasswordless)

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
	d.Set("pingid", values["pingid"])
	d.Set("pingid_windows_login_passwordless", values["pingid_windows_login_passwordless"])

	return diags
}

func resourceSignOnPolicyActionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient

	var diags diag.Diagnostics

	signOnPolicyAction, diags := expandSOPAction(d)
	if diags.HasError() {
		return diags
	}

	_, diags = sdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := apiClient.SignOnPolicyActionsApi.UpdateSignOnPolicyAction(ctx, d.Get("environment_id").(string), d.Get("sign_on_policy_id").(string), d.Id()).SignOnPolicyAction(*signOnPolicyAction).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, apiClient, d.Get("environment_id").(string), fO, fR, fErr)
		},
		"UpdateSignOnPolicyAction",
		customErrorSignOnPolicyActionCreateUpdate,
		nil,
	)
	if diags.HasError() {
		return diags
	}

	return resourceSignOnPolicyActionRead(ctx, d, meta)
}

func resourceSignOnPolicyActionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient

	var diags diag.Diagnostics

	_, diags = sdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fR, fErr := apiClient.SignOnPolicyActionsApi.DeleteSignOnPolicyAction(ctx, d.Get("environment_id").(string), d.Get("sign_on_policy_id").(string), d.Id()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, apiClient, d.Get("environment_id").(string), nil, fR, fErr)
		},
		"DeleteSignOnPolicyAction",
		sdk.CustomErrorResourceNotFoundWarning,
		func(ctx context.Context, r *http.Response, p1Error *model.P1Error) bool {

			if p1Error != nil {
				// Last action in the policy
				if v, ok := p1Error.GetDetailsOk(); ok && v != nil && len(v) > 0 {
					if v[0].GetCode() == "CONSTRAINT_VIOLATION" {
						if match, _ := regexp.MatchString("Cannot delete last action from the policy", v[0].GetMessage()); match {

							// create a generic SOP action
							sopPriority := 1

							signOnPolicyAction := &management.SignOnPolicyAction{}
							signOnPolicyAction.SignOnPolicyActionLogin = management.NewSignOnPolicyActionLogin(
								int32(sopPriority),
								management.ENUMSIGNONPOLICYTYPE_LOGIN,
							)

							_, innerDiags := sdk.ParseResponse(
								ctx,

								func() (any, *http.Response, error) {
									fO, fR, fErr := apiClient.SignOnPolicyActionsApi.CreateSignOnPolicyAction(ctx, d.Get("environment_id").(string), d.Get("sign_on_policy_id").(string)).SignOnPolicyAction(*signOnPolicyAction).Execute()
									return framework.CheckEnvironmentExistsOnPermissionsError(ctx, apiClient, d.Get("environment_id").(string), fO, fR, fErr)
								},
								"CreateSignOnPolicyAction-Delete",
								customErrorSignOnPolicyActionCreateUpdate,
								sdk.DefaultCreateReadRetryable,
							)

							diags = append(diags, innerDiags...)
							diags = append(diags, diag.Diagnostic{
								Severity: diag.Warning,
								Summary:  fmt.Sprintf("Cannot delete last action from the sign-on policy %s.  A generic policy action will be left in place but is not managed by the provider. This warning can be safely ignored if the sign-on policy %s was also destroyed.", d.Get("sign_on_policy_id").(string), d.Get("sign_on_policy_id").(string)),
								Detail:   "For more details about this warning, please see https://github.com/pingidentity/terraform-provider-pingone/issues/68",
							})

							if diags.HasError() {
								return false
							}

							return true
						}
					}
				}

			}

			return false
		},
	)

	return diags
}

func resourceSignOnPolicyActionImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {

	idComponents := []framework.ImportComponent{
		{
			Label:  "environment_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
		{
			Label:  "sign_on_policy_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
		{
			Label:  "sign_on_policy_action_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
	}

	attributes, err := framework.ParseImportID(d.Id(), idComponents...)
	if err != nil {
		return nil, err
	}

	d.Set("environment_id", attributes["environment_id"])
	d.Set("sign_on_policy_id", attributes["sign_on_policy_id"])
	d.SetId(attributes["sign_on_policy_action_id"])

	resourceSignOnPolicyActionRead(ctx, d, meta)

	return []*schema.ResourceData{d}, nil
}

var (
	customErrorSignOnPolicyActionCreateUpdate = func(error model.P1Error) diag.Diagnostics {
		var diags diag.Diagnostics

		// Value not allowed
		if details, ok := error.GetDetailsOk(); ok && details != nil && len(details) > 0 {
			if target, ok := details[0].GetTargetOk(); ok && details[0].GetCode() == "INVALID_VALUE" && *target == "newUserProvisioning.gateways" {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Only 'LDAP' type gateways are supported for new user provisioning.",
					Detail:   "The \"new_user_provisioning.gateway\" provided for the login sign on policy action is not supported by the PingOne platform.  Please ensure gateways are of type 'LDAP'.",
				})

				return diags
			}
		}

		return nil
	}
)

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

	if _, ok := d.GetOk("pingid"); ok {
		signOnPolicyAction.SignOnPolicyActionCommon = expandSOPActionPingID(d, sopPriority)
		processedCount += 1
	}

	if _, ok := d.GetOk("pingid_windows_login_passwordless"); ok {
		signOnPolicyAction.SignOnPolicyActionPingIDWinLoginPasswordless = expandSOPActionPingIDWinLoginPasswordless(d, sopPriority)
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
			if diags.HasError() {
				return nil, diags
			}
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
		if vc, ok := v.(*schema.Set); ok && vc != nil && len(vc.List()) > 0 && vc.List()[0] != "" {
			obj := make([]string, 0)
			for _, str := range vc.List() {
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

		if v1, ok := vp["discovery_rule"].(*schema.Set); ok && v1 != nil && len(v1.List()) > 0 && v1.List()[0] != nil {
			sopActionType.SetDiscoveryRules(expandSOPActionDiscoveryRules(v1.List()))
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
				if diags.HasError() {
					return nil, diags
				}
				sopActionType.SetCondition(*conditions)
			}
		}

		if v1, ok := d.GetOk("registration_local_population_id"); ok && v1 != "" {
			obj := *management.NewSignOnPolicyActionIDPAllOfRegistration(true)
			obj.SetPopulation(*management.NewSignOnPolicyActionLoginAllOfRegistrationPopulation(v1.(string)))

			if v2, ok := d.GetOk("registration_confirm_user_attributes"); ok {
				obj.SetConfirmIdentityProviderAttributes(v2.(bool))
			}

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
			if diags.HasError() {
				return nil, diags
			}
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

		if v2, ok := d.GetOk("registration_confirm_user_attributes"); ok {
			obj.SetConfirmIdentityProviderAttributes(v2.(bool))
		}

		sopActionType.SetRegistration(obj)
	}

	socialIDSet := false
	if v, ok := d.GetOk("social_provider_ids"); ok {
		if vc, ok := v.(*schema.Set); ok && vc != nil && len(vc.List()) > 0 && vc.List()[0] != "" {
			obj := make([]string, 0)
			for _, str := range vc.List() {
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

		if newUserProvisioningPlan, ok := vp["new_user_provisioning"].([]interface{}); ok && newUserProvisioningPlan != nil && len(newUserProvisioningPlan) > 0 && newUserProvisioningPlan[0] != nil {
			newUserProvisioningPlanMap := newUserProvisioningPlan[0].(map[string]interface{})

			if gatewaysPlan, ok := newUserProvisioningPlanMap["gateway"].(*schema.Set); ok && gatewaysPlan != nil && len(gatewaysPlan.List()) > 0 && gatewaysPlan.List()[0] != "" {
				gateways := make([]management.SignOnPolicyActionLoginAllOfNewUserProvisioningGateways, 0)

				for _, gatewayPlan := range gatewaysPlan.List() {
					gatewayPlanMap := gatewayPlan.(map[string]interface{})

					gateways = append(gateways, *management.NewSignOnPolicyActionLoginAllOfNewUserProvisioningGateways(
						gatewayPlanMap["id"].(string),
						management.EnumSignOnPolicyActionLoginNewUserProvisioningGatewayType(gatewayPlanMap["type"].(string)),
						*management.NewSignOnPolicyActionLoginAllOfNewUserProvisioningUserType(gatewayPlanMap["user_type_id"].(string)),
					))
				}

				sopActionType.SetNewUserProvisioning(
					*management.NewSignOnPolicyActionLoginAllOfNewUserProvisioning(
						gateways,
					),
				)
			}
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
			if diags.HasError() {
				return nil, diags
			}
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

func expandSOPActionPingID(d *schema.ResourceData, sopPriority int32) *management.SignOnPolicyActionCommon {

	if v, ok := d.Get("pingid").([]interface{}); ok && v != nil && len(v) > 0 {

		sopActionType := management.NewSignOnPolicyActionCommon(
			sopPriority,
			management.ENUMSIGNONPOLICYTYPE_PINGID_AUTHENTICATION,
		)

		return sopActionType

	}

	return nil
}

func expandSOPActionPingIDWinLoginPasswordless(d *schema.ResourceData, sopPriority int32) *management.SignOnPolicyActionPingIDWinLoginPasswordless {

	if v, ok := d.Get("pingid_windows_login_passwordless").([]interface{}); ok && v != nil && len(v) > 0 && v[0] != nil {
		vp := v[0].(map[string]interface{})

		sopActionType := management.NewSignOnPolicyActionPingIDWinLoginPasswordless(
			sopPriority,
			management.ENUMSIGNONPOLICYTYPE_PINGID_WINLOGIN_PASSWORDLESS_AUTHENTICATION,
			*management.NewSignOnPolicyActionPingIDWinLoginPasswordlessAllOfUniqueUserAttribute(vp["unique_user_attribute_name"].(string)),
			*management.NewSignOnPolicyActionPingIDWinLoginPasswordlessAllOfOfflineMode(vp["offline_mode_enabled"].(bool)),
		)

		return sopActionType

	}

	return nil
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

func buildSignOnOlderThanPwd(v int32) management.SignOnPolicyActionCommonConditionOrOrInner {
	return buildSignOnOlderThan(v, "pwd")
}

func buildSignOnOlderThanMfa(v int32) management.SignOnPolicyActionCommonConditionOrOrInner {
	return buildSignOnOlderThan(v, "mfa")
}

func buildSignOnOlderThan(v int32, lastSignOnContext string) management.SignOnPolicyActionCommonConditionOrOrInner {
	return management.SignOnPolicyActionCommonConditionOrOrInner{
		SignOnPolicyActionCommonConditionAggregate: &management.SignOnPolicyActionCommonConditionAggregate{
			SignOnPolicyActionCommonConditionGreater: management.NewSignOnPolicyActionCommonConditionGreater(v, getLastSignOnContextFull(lastSignOnContext)),
		},
	}
}

type attributeEquality struct {
	attributeReference   string
	attributeValueString *string
	attributeValueBool   *bool
}

func buildUserMemberOfPopulation(v []interface{}, sopPriority int32) (*management.SignOnPolicyActionCommonConditionOrOrInner, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributeList := make([]attributeEquality, 0)

	for _, population := range v {

		attribute := attributeEquality{
			attributeReference: "${user.population.id}",
		}

		if v, ok := population.(string); ok && v != "" {
			attribute.attributeValueString = &v
		} else {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Found Population ID to be invalid, found a non-string type or blank value",
				Detail:   "A population ID defined in `conditions.0.user_is_member_of_any_population_id` is either not a string type or is a blank value. This is most likely an issue with the provider itself. Please raise an issue with the provider maintainers.",
			})
			return nil, diags
		}

		attributeList = append(attributeList, attribute)
	}

	return buildAttributeEqualsCondition(attributeList, "user_is_member_of_any_population_id", sopPriority)

}

func buildUserAttributes(v []interface{}, sopPriority int32) (*management.SignOnPolicyActionCommonConditionOrOrInner, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributeList := make([]attributeEquality, 0)

	for _, attributeMap := range v {

		attribute := attributeEquality{
			attributeReference: attributeMap.(map[string]interface{})["attribute_reference"].(string),
		}

		valueSet := false

		if v, ok := attributeMap.(map[string]interface{})["value"].(string); ok && v != "" && !valueSet {
			attribute.attributeValueString = &v
			valueSet = true
		}

		if v, ok := attributeMap.(map[string]interface{})["value_boolean"].(bool); ok && !valueSet {
			attribute.attributeValueBool = &v
			valueSet = true
		}

		if !valueSet {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "One of `user_attribute_equals.value` or `user_attribute_equals.value_boolean` is required to be set.",
			})
			return nil, diags
		}

		attributeList = append(attributeList, attribute)
	}

	return buildAttributeEqualsCondition(attributeList, "user_attribute_equals", sopPriority)

}

func buildAttributeEqualsCondition(v []attributeEquality, tfSchemaAttribute string, sopPriority int32) (*management.SignOnPolicyActionCommonConditionOrOrInner, diag.Diagnostics) {
	var diags diag.Diagnostics

	if sopPriority < 2 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Condition `%s` is defined cannot be set when the policy action priority is 1.", tfSchemaAttribute),
		})
	} else {

		if len(v) > 1 {

			conditionList := make([]management.SignOnPolicyActionCommonConditionOrOrInner, 0)

			for _, attribute := range v {

				conditionList = append(conditionList, management.SignOnPolicyActionCommonConditionOrOrInner{
					SignOnPolicyActionCommonConditionAggregate: &management.SignOnPolicyActionCommonConditionAggregate{
						SignOnPolicyActionCommonConditionEquals: management.NewSignOnPolicyActionCommonConditionEquals(attribute.attributeReference, management.SignOnPolicyActionCommonConditionEqualsEquals{
							Bool:   attribute.attributeValueBool,
							String: attribute.attributeValueString,
						}),
					},
				})

			}

			return &management.SignOnPolicyActionCommonConditionOrOrInner{
				SignOnPolicyActionCommonConditionOr: &management.SignOnPolicyActionCommonConditionOr{
					Or: conditionList,
				},
			}, diags

		} else if len(v) == 1 {

			return &management.SignOnPolicyActionCommonConditionOrOrInner{
				SignOnPolicyActionCommonConditionAggregate: &management.SignOnPolicyActionCommonConditionAggregate{
					SignOnPolicyActionCommonConditionEquals: management.NewSignOnPolicyActionCommonConditionEquals(v[0].attributeReference, management.SignOnPolicyActionCommonConditionEqualsEquals{
						Bool:   v[0].attributeValueBool,
						String: v[0].attributeValueString,
					}),
				},
			}, diags
		} else {
			return nil, diags
		}
	}

	return nil, diags

}

func expandSOPActionConditionIDFirstAndLogin(condition interface{}, sopPriority int32) (*management.SignOnPolicyActionCommonConditionOrOrInner, diag.Diagnostics) {
	var diags diag.Diagnostics

	var conditionStructList = make([]management.SignOnPolicyActionCommonConditionOrOrInner, 0)

	if v, ok := condition.(map[string]interface{})["last_sign_on_older_than_seconds"].(int); ok && v > 0 {

		conditionStructList = append(conditionStructList, buildSignOnOlderThanPwd(int32(v)))
	}

	if v, ok := condition.(map[string]interface{})["user_is_member_of_any_population_id"].(*schema.Set); ok && v != nil && len(v.List()) > 0 && v.List()[0] != "" {

		var populationStructs *management.SignOnPolicyActionCommonConditionOrOrInner
		populationStructs, diags = buildUserMemberOfPopulation(v.List(), sopPriority)
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

	if v, ok := condition.(map[string]interface{})["last_sign_on_older_than_seconds"].(int); ok && v > 0 {

		conditionStruct.SignOnPolicyActionCommonConditionAggregate = buildSignOnOlderThanPwd(int32(v)).SignOnPolicyActionCommonConditionAggregate
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

	if v, ok := condition.(map[string]interface{})["last_sign_on_older_than_seconds"].(int); ok && v > 0 {

		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Argument is deprecated",
			Detail:   "Parameter `conditions.last_sign_on_older_than_seconds` has been deprecated when configured with an MFA Sign-on policy action.  Please use `last_sign_on_older_than_seconds_mfa` instead.",
		})

		conditionStructList = append(conditionStructList, buildSignOnOlderThanPwd(int32(v)))
	}

	if v, ok := condition.(map[string]interface{})["last_sign_on_older_than_seconds_mfa"].(int); ok && v > 0 {

		conditionStructList = append(conditionStructList, buildSignOnOlderThanMfa(int32(v)))
	}

	if v, ok := condition.(map[string]interface{})["ip_out_of_range_cidr"].(*schema.Set); ok && v != nil && len(v.List()) > 0 && v.List()[0] != "" {

		obj := make([]string, 0)
		for _, str := range v.List() {
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

		if j, ok := condition.(map[string]interface{})["anonymous_network_detected_allowed_cidr"].(*schema.Set); ok && j != nil && len(j.List()) > 0 && j.List()[0] != "" {
			obj := make([]string, 0)
			for _, str := range j.List() {
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

	if v, ok := condition.(map[string]interface{})["user_is_member_of_any_population_id"].(*schema.Set); ok && v != nil && len(v.List()) > 0 && v.List()[0] != "" {

		var populationStructs *management.SignOnPolicyActionCommonConditionOrOrInner
		populationStructs, diags = buildUserMemberOfPopulation(v.List(), sopPriority)
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
	last_sign_on_older_than_seconds_mfa     *int32
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

				if returnCondition.last_sign_on_older_than_seconds != nil || returnCondition.last_sign_on_older_than_seconds_mfa != nil {
					diags = append(diags, diag.Diagnostic{
						Severity: diag.Error,
						Summary:  "Condition `last_sign_on_older_than_seconds` or `last_sign_on_older_than_seconds_mfa` has has multiple nested values.  This is not supported in the provider.  Please raise an issue.",
					})

					return nil, diags
				}

				greaterThanValue := func() *int32 { b := vc.GetGreater(); return &b }()

				switch *condition {
				case "${session.lastSignOn.withAuthenticator.pwd.at}":
					returnCondition.last_sign_on_older_than_seconds = greaterThanValue
				case "${session.lastSignOn.withAuthenticator.mfa.at}":
					returnCondition.last_sign_on_older_than_seconds_mfa = greaterThanValue
				}

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
				returnCondition.user_is_member_of_any_population_id = append(returnCondition.user_is_member_of_any_population_id, *vc.GetEquals().String)

			} else {
				condition := attributeEquality{
					attributeReference:   vc.GetValue(),
					attributeValueString: vc.GetEquals().String,
					attributeValueBool:   vc.GetEquals().Bool,
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
	if diags.HasError() {
		return nil, diags
	}

	flattenedConditions := map[string]interface{}{
		"last_sign_on_older_than_seconds":         nil,
		"last_sign_on_older_than_seconds_mfa":     nil,
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

	if conditions.last_sign_on_older_than_seconds_mfa != nil {
		flattenedConditions["last_sign_on_older_than_seconds_mfa"] = *conditions.last_sign_on_older_than_seconds_mfa
	}

	if conditions.user_is_member_of_any_population_id != nil {
		flattenedConditions["user_is_member_of_any_population_id"] = conditions.user_is_member_of_any_population_id
	}

	if conditions.user_attribute_equals != nil {
		attributeList := make([]map[string]interface{}, 0)
		for _, attributeStruct := range conditions.user_attribute_equals {

			attribute := map[string]interface{}{
				"attribute_reference": attributeStruct.attributeReference,
			}

			if v := attributeStruct.attributeValueString; v != nil {
				attribute["value"] = *v
			} else {
				attribute["value"] = ""
			}

			if v := attributeStruct.attributeValueBool; v != nil {
				attribute["value_boolean"] = *v
			} else {
				attribute["value_boolean"] = nil
			}

			attributeList = append(attributeList, attribute)
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

	action["new_user_provisioning"] = nil
	if v, ok := signOnPolicyActionLogin.GetNewUserProvisioningOk(); ok {

		newUserProvisioningList := make([]interface{}, 0, 1)
		newUserProvisioning := map[string]interface{}{}

		if gateways, ok := v.GetGatewaysOk(); ok && len(gateways) > 0 {
			gatewaysMap := make([]interface{}, 0, 1)

			for _, gateway := range gateways {

				gatewayMap := map[string]interface{}{}

				if c, ok := gateway.GetIdOk(); ok && *c != "" {
					gatewayMap["id"] = *c
				} else {
					gatewayMap["id"] = nil
				}

				if c, ok := gateway.GetTypeOk(); ok && *c != "" {
					gatewayMap["type"] = *c
				} else {
					gatewayMap["type"] = nil
				}

				if c, ok := gateway.GetUserTypeOk(); ok && c != nil {
					gatewayMap["user_type_id"] = c.GetId()
				} else {
					gatewayMap["user_type_id"] = nil
				}

				gatewaysMap = append(gatewaysMap, gatewayMap)
			}

			newUserProvisioning["gateway"] = gatewaysMap
		} else {
			newUserProvisioning["gateway"] = nil
		}

		action["new_user_provisioning"] = append(newUserProvisioningList, newUserProvisioning)
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

func flattenActionPingIDWinLoginPasswordless(signOnPolicyActionPingIDWinLoginPasswordless *management.SignOnPolicyActionPingIDWinLoginPasswordless) []interface{} {
	actionList := make([]interface{}, 0, 1)

	return append(actionList, map[string]interface{}{
		"unique_user_attribute_name": signOnPolicyActionPingIDWinLoginPasswordless.GetUniqueUserAttribute().Name,
		"offline_mode_enabled":       signOnPolicyActionPingIDWinLoginPasswordless.GetOfflineMode().Enabled,
	})
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
	case *management.SignOnPolicyActionCommon:
		actionID = instance.SignOnPolicyActionCommon.GetId()
	case *management.SignOnPolicyActionPingIDWinLoginPasswordless:
		actionID = instance.SignOnPolicyActionPingIDWinLoginPasswordless.GetId()
	}

	return actionID
}

func getLastSignOnContextFull(lastSignOnContext string) string {
	lastSignOnContextFull := ""
	switch lastSignOnContext {
	case "pwd":
		lastSignOnContextFull = "${session.lastSignOn.withAuthenticator.pwd.at}"
	case "mfa":
		lastSignOnContextFull = "${session.lastSignOn.withAuthenticator.mfa.at}"
	default:
		lastSignOnContextFull = "${session.lastSignOn.withAuthenticator.pwd.at}"
	}
	return lastSignOnContextFull
}
