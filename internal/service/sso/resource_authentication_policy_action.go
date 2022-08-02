package sso

import (
	"fmt"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
)

func expandSOPAction(d interface{}, sopPriority int32) (*management.SignOnPolicyAction, diag.Diagnostics) {

	signOnPolicyAction := &management.SignOnPolicyAction{}
	var diags diag.Diagnostics

	actionType := d.(map[string]interface{})["action_type"].(string)

	switch actionType {
	case string(management.ENUMSIGNONPOLICYTYPE_AGREEMENT):
		signOnPolicyAction.SignOnPolicyActionAgreement, diags = expandSOPActionAgreement(d, sopPriority)
	case string(management.ENUMSIGNONPOLICYTYPE_IDENTIFIER_FIRST):
		signOnPolicyAction.SignOnPolicyActionIDFirst, diags = expandSOPActionIDFirst(d, sopPriority)
	case string(management.ENUMSIGNONPOLICYTYPE_IDENTITY_PROVIDER):
		signOnPolicyAction.SignOnPolicyActionIDP, diags = expandSOPActionIDP(d, sopPriority)
	case string(management.ENUMSIGNONPOLICYTYPE_LOGIN):
		signOnPolicyAction.SignOnPolicyActionLogin, diags = expandSOPActionLogin(d, sopPriority)
	case string(management.ENUMSIGNONPOLICYTYPE_MULTI_FACTOR_AUTHENTICATION):
		signOnPolicyAction.SignOnPolicyActionMFA, diags = expandSOPActionMFA(d, sopPriority)
	case string(management.ENUMSIGNONPOLICYTYPE_PROGRESSIVE_PROFILING):
		signOnPolicyAction.SignOnPolicyActionProgressiveProfiling, diags = expandSOPActionProgressiveProfiling(d, sopPriority)
	default:
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Policy action %s not supported", actionType),
		})
		return nil, diags
	}

	return signOnPolicyAction, diags
}

func expandSOPActionAgreement(d interface{}, sopPriority int32) (*management.SignOnPolicyActionAgreement, diag.Diagnostics) {
	var diags diag.Diagnostics

	if v, ok := d.(map[string]interface{})["agreement_options"].([]interface{}); ok && v != nil && len(v) > 0 && v[0] != nil {
		vp := v[0].(map[string]interface{})

		sopActionType := management.NewSignOnPolicyActionAgreement(
			sopPriority,
			management.ENUMSIGNONPOLICYTYPE_AGREEMENT,
			*management.NewSignOnPolicyActionAgreementAllOfAgreement(vp["agreement_id"].(string)),
		)

		if vc, ok := d.(map[string]interface{})["conditions"].([]interface{}); ok && vc != nil && len(vc) > 0 && vc[0] != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Block `conditions` has no effect when using the agreement action type",
			})
		}

		if vd, ok := vp["show_decline_option"].(bool); ok {
			sopActionType.SetDisableDeclineOption(!vd)
		}

		return sopActionType, nil

	}

	diags = append(diags, diag.Diagnostic{
		Severity: diag.Error,
		Summary:  "Block `agreement` with `agreement_id` must be defined when using the agreement action type",
	})

	return nil, diags
}

func expandSOPActionIDFirst(d interface{}, sopPriority int32) (*management.SignOnPolicyActionIDFirst, diag.Diagnostics) {

	var diags diag.Diagnostics

	sopActionType := management.NewSignOnPolicyActionIDFirst(
		sopPriority,
		management.ENUMSIGNONPOLICYTYPE_IDENTIFIER_FIRST,
	)

	if vc, ok := d.(map[string]interface{})["conditions"].([]interface{}); ok && vc != nil && len(vc) > 0 && vc[0] != nil {
		var conditions *management.SignOnPolicyActionCommonConditionOrOrInner
		conditions, diags = expandSOPActionCondition(vc[0], management.ENUMSIGNONPOLICYTYPE_IDENTIFIER_FIRST, sopPriority)
		sopActionType.SetCondition(*conditions)
	}

	if v, ok := d.(map[string]interface{})["identifier_first_options"].([]interface{}); ok && v != nil && len(v) > 0 && v[0] != nil {
		vp := v[0].(map[string]interface{})

		idpAttributes := vp["confirm_identity_provider_attributes"]
		if idpAttributes != nil {
			sopActionType.SetConfirmIdentityProviderAttributes(idpAttributes.(bool))
		}

		idpLockout := vp["enforce_lockout_for_identity_providers"]
		if idpLockout != nil {
			sopActionType.SetEnforceLockoutForIdentityProviders(idpLockout.(bool))
		}

		if v1, ok := vp["recovery"].([]interface{}); ok && v1 != nil && len(v1) > 0 && v1[0] != nil {

			recoveryObj := *management.NewSignOnPolicyActionLoginAllOfRecovery(v1[0].(map[string]interface{})["enabled"].(bool))
			sopActionType.SetRecovery(recoveryObj)

		}

		if v1, ok := vp["registration"].([]interface{}); ok && v1 != nil && len(v1) > 0 && v1[0] != nil {

			registrationObj := *management.NewSignOnPolicyActionLoginAllOfRegistration(v1[0].(map[string]interface{})["enabled"].(bool))

			externalHref := v1[0].(map[string]interface{})["external_href"]
			if externalHref != nil {

				registrationExternalObj := *management.NewSignOnPolicyActionLoginAllOfRegistrationExternal(externalHref.(string))
				registrationObj.SetExternal(registrationExternalObj)

			}

			populationID := v1[0].(map[string]interface{})["population_id"]
			if populationID != nil {

				registrationPopulationObj := *management.NewSignOnPolicyActionLoginAllOfRegistrationPopulation(populationID.(string))
				registrationObj.SetPopulation(registrationPopulationObj)

			}

			sopActionType.SetRegistration(registrationObj)

		}

		if v1, ok := vp["social_providers"].([]interface{}); ok && v1 != nil && len(v1) > 0 && v1[0] != nil {
			sopActionType.SetSocialProviders(expandSOPActionSocialProviders(v1))
		}

		if v1, ok := vp["discovery_rule"].([]interface{}); ok && v1 != nil && len(v1) > 0 && v1[0] != nil {
			sopActionType.SetDiscoveryRules(expandSOPActionDiscoveryRules(v1))
		}

	}

	return sopActionType, diags

}

func expandSOPActionIDP(d interface{}, sopPriority int32) (*management.SignOnPolicyActionIDP, diag.Diagnostics) {
	var diags diag.Diagnostics

	if v, ok := d.(map[string]interface{})["identity_provider_options"].([]interface{}); ok && v != nil && len(v) > 0 && v[0] != nil {
		vp := v[0].(map[string]interface{})

		sopActionType := management.NewSignOnPolicyActionIDP(
			sopPriority,
			management.ENUMSIGNONPOLICYTYPE_IDENTITY_PROVIDER,
			*management.NewSignOnPolicyActionIDPAllOfIdentityProvider(vp["identity_provider_id"].(string)),
		)

		if vc, ok := d.(map[string]interface{})["conditions"].([]interface{}); ok && vc != nil && len(vc) > 0 && vc[0] != nil {
			var conditions *management.SignOnPolicyActionCommonConditionOrOrInner
			conditions, diags = expandSOPActionCondition(vc[0], management.ENUMSIGNONPOLICYTYPE_IDENTITY_PROVIDER, sopPriority)
			sopActionType.SetCondition(*conditions)
		}

		acrValues := vp["acr_values"]
		if acrValues != nil {
			sopActionType.SetAcrValues(acrValues.(string))
		}

		passUserCtx := vp["pass_user_context"]
		if passUserCtx != nil {
			sopActionType.SetPassUserContext(passUserCtx.(bool))
		}

		registration := vp["registration"]
		if registration != nil {

			registrationObj := *management.NewSignOnPolicyActionIDPAllOfRegistration(registration.(map[string]interface{})["enabled"].(bool))

			sopActionType.SetRegistration(registrationObj)

		}

		return sopActionType, nil

	}

	diags = append(diags, diag.Diagnostic{
		Severity: diag.Error,
		Summary:  "Block `identity_provider` with `identity_provider_id` must be defined when using the identity provider action type",
	})

	return nil, diags
}

func expandSOPActionLogin(d interface{}, sopPriority int32) (*management.SignOnPolicyActionLogin, diag.Diagnostics) {
	var diags diag.Diagnostics

	sopActionType := management.NewSignOnPolicyActionLogin(
		sopPriority,
		management.ENUMSIGNONPOLICYTYPE_LOGIN,
	)

	if v, ok := d.(map[string]interface{})["conditions"].([]interface{}); ok && v != nil && len(v) > 0 && v[0] != nil {
		var conditions *management.SignOnPolicyActionCommonConditionOrOrInner
		conditions, diags = expandSOPActionCondition(v[0], management.ENUMSIGNONPOLICYTYPE_LOGIN, sopPriority)
		sopActionType.SetCondition(*conditions)
	}

	// block is optional
	if v, ok := d.(map[string]interface{})["login_options"].([]interface{}); ok && v != nil && len(v) > 0 && v[0] != nil {
		vp := v[0].(map[string]interface{})

		idpAttributes := vp["confirm_identity_provider_attributes"]
		if idpAttributes != nil {
			sopActionType.SetConfirmIdentityProviderAttributes(idpAttributes.(bool))
		}

		idpLockout := vp["enforce_lockout_for_identity_providers"]
		if idpLockout != nil {
			sopActionType.SetEnforceLockoutForIdentityProviders(idpLockout.(bool))
		}

		if v1, ok := vp["recovery"].([]interface{}); ok && v1 != nil && len(v1) > 0 && v1[0] != nil {

			recoveryObj := *management.NewSignOnPolicyActionLoginAllOfRecovery(v1[0].(map[string]interface{})["enabled"].(bool))
			sopActionType.SetRecovery(recoveryObj)

		}

		if v1, ok := vp["registration"].([]interface{}); ok && v1 != nil && len(v1) > 0 && v1[0] != nil {

			registrationObj := *management.NewSignOnPolicyActionLoginAllOfRegistration(v1[0].(map[string]interface{})["enabled"].(bool))

			externalHref := v1[0].(map[string]interface{})["external_href"]
			if externalHref != nil {

				registrationExternalObj := *management.NewSignOnPolicyActionLoginAllOfRegistrationExternal(externalHref.(string))
				registrationObj.SetExternal(registrationExternalObj)

			}

			populationID := v1[0].(map[string]interface{})["population_id"]
			if populationID != nil {

				registrationPopulationObj := *management.NewSignOnPolicyActionLoginAllOfRegistrationPopulation(populationID.(string))
				registrationObj.SetPopulation(registrationPopulationObj)

			}

			sopActionType.SetRegistration(registrationObj)

		}

		if v1, ok := vp["social_providers"].([]interface{}); ok && v1 != nil && len(v1) > 0 && v1[0] != nil {
			sopActionType.SetSocialProviders(expandSOPActionSocialProviders(v1))
		}

	}

	return sopActionType, diags
}

func expandSOPActionMFA(d interface{}, sopPriority int32) (*management.SignOnPolicyActionMFA, diag.Diagnostics) {
	var diags diag.Diagnostics

	sopActionType := management.NewSignOnPolicyActionMFA(
		sopPriority,
		management.ENUMSIGNONPOLICYTYPE_MULTI_FACTOR_AUTHENTICATION,
	)

	if vc, ok := d.(map[string]interface{})["conditions"].([]interface{}); ok && vc != nil && len(vc) > 0 && vc[0] != nil {
		var conditions *management.SignOnPolicyActionCommonConditionOrOrInner
		conditions, diags = expandSOPActionCondition(vc[0], management.ENUMSIGNONPOLICYTYPE_MULTI_FACTOR_AUTHENTICATION, sopPriority)
		sopActionType.SetCondition(*conditions)
	}

	if v, ok := d.(map[string]interface{})["mfa_options"].([]interface{}); ok && v != nil && len(v) > 0 && v[0] != nil {
		vp := v[0].(map[string]interface{})

		devicePolicyID := vp["device_authentication_policy_id"]
		if devicePolicyID != nil {
			sopDevicePolicy := *management.NewSignOnPolicyActionMFAAllOfDeviceAuthenticationPolicy(devicePolicyID.(string))
			sopActionType.SetDeviceAuthenticationPolicy(sopDevicePolicy)
		}

		noDeviceMode := vp["no_device_mode"]
		if noDeviceMode != nil {
			sopActionType.SetNoDeviceMode(management.EnumSignOnPolicyNoDeviceMode(noDeviceMode.(string)))
		}

	}

	return sopActionType, diags
}

func expandSOPActionProgressiveProfiling(d interface{}, sopPriority int32) (*management.SignOnPolicyActionProgressiveProfiling, diag.Diagnostics) {
	var diags diag.Diagnostics

	if v, ok := d.(map[string]interface{})["progressive_profiling_options"].([]interface{}); ok && v != nil && len(v) > 0 && v[0] != nil {
		vp := v[0].(map[string]interface{})

		sopActionType := management.NewSignOnPolicyActionProgressiveProfiling(
			sopPriority,
			management.ENUMSIGNONPOLICYTYPE_PROGRESSIVE_PROFILING,
			expandSOPActionAttributes(vp["attribute"].(*schema.Set)),
			vp["prevent_multiple_prompts_per_flow"].(bool),
			int32(vp["prompt_interval_seconds"].(int)),
			vp["prompt_text"].(string),
		)

		if vc, ok := d.(map[string]interface{})["conditions"].([]interface{}); ok && vc != nil && len(vc) > 0 && vc[0] != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Block `conditions` has no effect when using the progressive profiling action type",
			})
		}

		return sopActionType, nil

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

		condition := item.(map[string]interface{})["condition"]
		conditionObj := *management.NewSignOnPolicyActionIDFirstAllOfCondition(condition.(map[string]interface{})["contains"].(string), condition.(map[string]interface{})["value"].(string))

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

func expandSOPActionSocialProviders(items []interface{}) []management.SignOnPolicyActionLoginAllOfSocialProviders {

	var socialProviders []management.SignOnPolicyActionLoginAllOfSocialProviders

	for _, item := range items {

		socialProviders = append(socialProviders, *management.NewSignOnPolicyActionLoginAllOfSocialProviders(item.(string)))

	}

	return socialProviders

}

func expandSOPActionCondition(condition interface{}, actionType management.EnumSignOnPolicyType, sopPriority int32) (*management.SignOnPolicyActionCommonConditionOrOrInner, diag.Diagnostics) {

	sopConditions := &management.SignOnPolicyActionCommonConditionOrOrInner{}
	var diags diag.Diagnostics

	switch actionType {
	case management.ENUMSIGNONPOLICYTYPE_IDENTIFIER_FIRST:
		sopConditions, diags = expandSOPActionConditionIDFirst(condition, sopPriority)
	case management.ENUMSIGNONPOLICYTYPE_IDENTITY_PROVIDER:
		sopConditions, diags = expandSOPActionConditionIDP(condition, sopPriority)
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

func expandSOPActionConditionIDFirstAndLogin(condition interface{}, sopPriority int32) (*management.SignOnPolicyActionCommonConditionOrOrInner, diag.Diagnostics) {
	var diags diag.Diagnostics

	conditionStruct := &management.SignOnPolicyActionCommonConditionOrOrInner{}

	var lastSignOnCond *management.SignOnPolicyActionCommonConditionGreater
	var populationMembership *management.SignOnPolicyActionCommonConditionOr
	var userAttributes *management.SignOnPolicyActionCommonConditionOr

	processed := 0

	if v, ok := condition.(map[string]interface{})["last_sign_on_older_than"].(int); ok {
		lastSignOnCond = management.NewSignOnPolicyActionCommonConditionGreater(int32(v), "${session.lastSignOn.withAuthenticator.pwd.at}")
		processed += 1
	}

	if v, ok := condition.(map[string]interface{})["user_is_member_of_any_population_id"].([]string); ok {

		if sopPriority < 2 {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Condition `user_is_member_of_any_population_id` is defined but has no effect where the action priority is 1.",
			})
		} else {

			//loop the things

			processed += 1
		}
	}

	if v, ok := condition.(map[string]interface{})["user_attribute_equals"].(*schema.Set); ok {

		if sopPriority < 2 {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Condition `user_attribute_equals` is defined but has no effect where the action priority is 1.",
			})
		} else {

			//loop the things

			processed += 1
		}
	}

	if processed > 1 {

	} else {

	}

	return conditionStruct, diags
}

func expandSOPActionConditionIDFirst(condition interface{}, sopPriority int32) (*management.SignOnPolicyActionCommonConditionOrOrInner, diag.Diagnostics) {
	return expandSOPActionConditionIDFirstAndLogin(condition, sopPriority)
}

func expandSOPActionConditionIDP(condition interface{}, sopPriority int32) (*management.SignOnPolicyActionCommonConditionOrOrInner, diag.Diagnostics) {
	var diags diag.Diagnostics

	conditionStruct := &management.SignOnPolicyActionCommonConditionOrOrInner{}

	if v, ok := condition.(map[string]interface{})["last_sign_on_older_than"].(int); ok {

		conditionAggregateStruct := &management.SignOnPolicyActionCommonConditionAggregate{
			SignOnPolicyActionCommonConditionGreater: management.NewSignOnPolicyActionCommonConditionGreater(int32(v), "${session.lastSignOn.withAuthenticator.pwd.at}"),
		}

		conditionStruct.SignOnPolicyActionCommonConditionAggregate = conditionAggregateStruct
	}

	return conditionStruct, diags

}

func expandSOPActionConditionLogin(condition interface{}, sopPriority int32) (*management.SignOnPolicyActionCommonConditionOrOrInner, diag.Diagnostics) {
	return expandSOPActionConditionIDFirstAndLogin(condition, sopPriority)
}

func expandSOPActionConditionMFA(condition interface{}, sopPriority int32) (*management.SignOnPolicyActionCommonConditionOrOrInner, diag.Diagnostics) {
	var diags diag.Diagnostics

	conditionStruct := &management.SignOnPolicyActionCommonConditionOrOrInner{}

	if v, ok := condition.(map[string]interface{})["last_sign_on_older_than"].(int); ok {
		condition := management.NewSignOnPolicyActionCommonConditionGreater(int32(v), "${session.lastSignOn.withAuthenticator.pwd.at}")
	}

	ip_out_of_range_cidr.([]string)
	ip_reputation_high_risk.(bool)
	geovelocity_anomaly_detected.(bool)
	anonymous_network_detected.(bool)
	anonymous_network_detected_allowed_cidr.([]string)

	user_is_member_of_any_population_id.([]string)
	user_attribute_equals.(*schema.Set)

	return conditionStruct, diags
}

func flattenSOPActions(actions []management.SignOnPolicyAction) ([]interface{}, error) {

	// Sort by priorty first
	sort.Slice(actions, func(i, j int) bool {
		return getActionPriority(actions[i]) < getActionPriority(actions[j])
	})

	// Set the return var
	sopActions := make([]interface{}, 0)

	// Loop the ordered slice
	for _, action := range actions {

		actionMap := map[string]interface{}{}

		match := 0

		if action.SignOnPolicyActionAgreement != nil {

			actionMap["action_type"] = string(management.ENUMSIGNONPOLICYTYPE_AGREEMENT)

			if v, ok := action.SignOnPolicyActionAgreement.GetConditionOk(); ok {
				actionMap["conditions"] = flattenConditions(*v)
			}
			actionMap["agreement_options"] = flattenActionAgreement(action.SignOnPolicyActionAgreement)

			match++
		}

		if action.SignOnPolicyActionIDFirst != nil {

			actionMap["action_type"] = string(management.ENUMSIGNONPOLICYTYPE_IDENTIFIER_FIRST)

			if v, ok := action.SignOnPolicyActionIDFirst.GetConditionOk(); ok {
				actionMap["conditions"] = flattenConditions(*v)
			}
			actionMap["identifier_first_options"] = flattenActionIDFirst(action.SignOnPolicyActionIDFirst)

			match++
		}

		if action.SignOnPolicyActionIDP != nil {

			actionMap["action_type"] = string(management.ENUMSIGNONPOLICYTYPE_IDENTITY_PROVIDER)

			if v, ok := action.SignOnPolicyActionIDP.GetConditionOk(); ok {
				actionMap["conditions"] = flattenConditions(*v)
			}
			actionMap["identity_provider_options"] = flattenActionIDP(action.SignOnPolicyActionIDP)

			match++
		}

		if action.SignOnPolicyActionLogin != nil {

			actionMap["action_type"] = string(management.ENUMSIGNONPOLICYTYPE_LOGIN)

			if v, ok := action.SignOnPolicyActionLogin.GetConditionOk(); ok {
				actionMap["conditions"] = flattenConditions(*v)
			}
			actionMap["login_options"] = flattenActionLogin(action.SignOnPolicyActionLogin)

			match++
		}

		if action.SignOnPolicyActionMFA != nil {

			actionMap["action_type"] = string(management.ENUMSIGNONPOLICYTYPE_MULTI_FACTOR_AUTHENTICATION)

			if v, ok := action.SignOnPolicyActionMFA.GetConditionOk(); ok {
				actionMap["conditions"] = flattenConditions(*v)
			}
			actionMap["mfa_options"] = flattenActionMFA(action.SignOnPolicyActionMFA)

			match++
		}

		if action.SignOnPolicyActionProgressiveProfiling != nil {

			actionMap["action_type"] = string(management.ENUMSIGNONPOLICYTYPE_PROGRESSIVE_PROFILING)

			if v, ok := action.SignOnPolicyActionProgressiveProfiling.GetConditionOk(); ok {
				actionMap["conditions"] = flattenConditions(*v)
			}
			actionMap["progressive_profiling_options"] = flattenActionProgressiveProfiling(action.SignOnPolicyActionProgressiveProfiling)

			match++
		}

		if match == 1 { // One match
			sopActions = append(sopActions, actionMap)

		} else if match > 1 { // More than one shouldn't happen
			return nil, fmt.Errorf("More than one action type exists for a single authentication action, this is not supported")

		} else if match < 1 { // None shouldn't happen
			return nil, fmt.Errorf("No action type exists for a single authentication action, this is not supported")
		}

	}

	return sopActions, nil

}

func flattenConditions(signOnPolicyActionCommonConditions management.SignOnPolicyActionCommonConditionOrOrInner) []interface{} {

	// TODO: this whole thing
	conditionsList := make([]interface{}, 0, 1)

	// conditions := map[string]interface{}{}

	// if v, ok := signOnPolicyActionCommonConditions.GetIpRangeOk(); ok {
	// 	conditions["ip_range"] = v
	// } else {
	// 	conditions["ip_range"] = nil
	// }

	// if v, ok := signOnPolicyActionCommonConditions.GetSecondsSinceOk(); ok {
	// 	conditions["action_session_length_mins"] = v
	// } else {
	// 	conditions["action_session_length_mins"] = nil
	// }

	// conditionsList = append(conditionsList, conditions)
	return conditionsList
}

func flattenActionProgressiveProfiling(signOnPolicyActionProgressiveProfiling *management.SignOnPolicyActionProgressiveProfiling) []interface{} {
	actionList := make([]interface{}, 0, 1)

	action := map[string]interface{}{}

	action["attribute"] = flattenActionProgressiveProfilingAttributes(signOnPolicyActionProgressiveProfiling.GetAttributes())
	action["prevent_multiple_prompts_per_flow"] = signOnPolicyActionProgressiveProfiling.GetPreventMultiplePromptsPerFlow()
	action["prompt_interval_seconds"] = signOnPolicyActionProgressiveProfiling.GetPromptIntervalSeconds()
	action["prompt_text"] = signOnPolicyActionProgressiveProfiling.GetPromptText()

	actionList = append(actionList, action)
	return actionList
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

	action := map[string]interface{}{}

	action["device_authentication_policy_id"] = signOnPolicyActionMFA.DeviceAuthenticationPolicy.GetId()

	if v, ok := signOnPolicyActionMFA.GetNoDeviceModeOk(); ok {
		action["no_device_mode"] = v
	}

	actionList = append(actionList, action)
	return actionList
}

func flattenActionLogin(signOnPolicyActionLogin *management.SignOnPolicyActionLogin) []interface{} {
	actionList := make([]interface{}, 0, 1)

	action := map[string]interface{}{}

	if v, ok := signOnPolicyActionLogin.GetConfirmIdentityProviderAttributesOk(); ok {
		action["confirm_identity_provider_attributes"] = v
	}

	if v, ok := signOnPolicyActionLogin.GetEnforceLockoutForIdentityProvidersOk(); ok {
		action["enforce_lockout_for_identity_providers"] = v
	}

	if v, ok := signOnPolicyActionLogin.GetRecoveryOk(); ok {
		action["recovery"] = flattenActionRecoveryInner(v)
	}

	if v, ok := signOnPolicyActionLogin.GetRegistrationOk(); ok {
		action["registration"] = flattenActionRegistrationInner(v)
	}

	if v, ok := signOnPolicyActionLogin.GetSocialProvidersOk(); ok {
		action["social_providers"] = flattenActionSocialProvidersInner(v)
	}

	actionList = append(actionList, action)
	return actionList
}

func flattenActionIDP(signOnPolicyActionIDP *management.SignOnPolicyActionIDP) []interface{} {
	actionList := make([]interface{}, 0, 1)

	action := map[string]interface{}{}

	action["identity_provider_id"] = signOnPolicyActionIDP.IdentityProvider.GetId()

	if v, ok := signOnPolicyActionIDP.GetAcrValuesOk(); ok {
		action["acr_values"] = v
	}

	if v, ok := signOnPolicyActionIDP.GetPassUserContextOk(); ok {
		action["pass_user_context"] = v
	}

	if v, ok := signOnPolicyActionIDP.GetRegistrationOk(); ok {
		action["registration"] = flattenActionIDPRegistrationInner(v)
	}

	actionList = append(actionList, action)
	return actionList
}

func flattenActionIDFirst(signOnPolicyActionIDFirst *management.SignOnPolicyActionIDFirst) []interface{} {
	actionList := make([]interface{}, 0, 1)

	action := map[string]interface{}{}

	if v, ok := signOnPolicyActionIDFirst.GetConfirmIdentityProviderAttributesOk(); ok {
		action["confirm_identity_provider_attributes"] = v
	}

	if v, ok := signOnPolicyActionIDFirst.GetDiscoveryRulesOk(); ok {
		action["discovery_rule"] = flattenDiscoveryRulesInner(v)
	}

	if v, ok := signOnPolicyActionIDFirst.GetEnforceLockoutForIdentityProvidersOk(); ok {
		action["enforce_lockout_for_identity_providers"] = v
	}

	if v, ok := signOnPolicyActionIDFirst.GetRecoveryOk(); ok {
		action["recovery"] = flattenActionRecoveryInner(v)
	}

	if v, ok := signOnPolicyActionIDFirst.GetRegistrationOk(); ok {
		action["registration"] = flattenActionRegistrationInner(v)
	}

	if v, ok := signOnPolicyActionIDFirst.GetSocialProvidersOk(); ok {
		action["social_providers"] = flattenActionSocialProvidersInner(v)
	}

	actionList = append(actionList, action)
	return actionList
}

func flattenActionAgreement(signOnPolicyActionAgreement *management.SignOnPolicyActionAgreement) []interface{} {
	actionList := make([]interface{}, 0, 1)

	action := map[string]interface{}{}

	action["agreement_id"] = signOnPolicyActionAgreement.Agreement.GetId()

	if v, ok := signOnPolicyActionAgreement.GetDisableDeclineOptionOk(); ok {
		action["show_decline_option"] = !*v
	}

	actionList = append(actionList, action)
	return actionList
}

func flattenActionRegistrationInner(signOnPolicyActionLoginAllOfRegistration *management.SignOnPolicyActionLoginAllOfRegistration) []interface{} {
	actionList := make([]interface{}, 0, 1)

	action := map[string]interface{}{}

	action["enabled"] = signOnPolicyActionLoginAllOfRegistration.GetEnabled()

	if v, ok := signOnPolicyActionLoginAllOfRegistration.GetExternalOk(); ok {
		action["external_href"] = v.GetHref()
	}

	if v, ok := signOnPolicyActionLoginAllOfRegistration.GetPopulationOk(); ok {
		action["population_id"] = v.GetId()
	}

	actionList = append(actionList, action)
	return actionList
}

func flattenActionIDPRegistrationInner(signOnPolicyActionIDPAllOfRegistration *management.SignOnPolicyActionIDPAllOfRegistration) []interface{} {
	actionList := make([]interface{}, 0, 1)

	action := map[string]interface{}{}

	action["enabled"] = signOnPolicyActionIDPAllOfRegistration.GetEnabled()

	if v, ok := signOnPolicyActionIDPAllOfRegistration.GetPopulationOk(); ok {
		action["population_id"] = v.GetId()
	}

	actionList = append(actionList, action)
	return actionList
}

func flattenDiscoveryRulesInner(signOnPolicyActionIDFirstAllOfDiscoveryRules []management.SignOnPolicyActionIDFirstAllOfDiscoveryRules) []interface{} {
	rules := make([]interface{}, 0, len(signOnPolicyActionIDFirstAllOfDiscoveryRules))
	for _, rule := range signOnPolicyActionIDFirstAllOfDiscoveryRules {

		condition := flattenDiscoveryRulesInnerCondition(rule.GetCondition())
		idpID := rule.IdentityProvider.GetId()

		rules = append(rules, map[string]interface{}{
			"condition":            condition,
			"identity_provider_id": idpID,
		})
	}
	return rules
}

func flattenDiscoveryRulesInnerCondition(signOnPolicyActionIDFirstAllOfCondition management.SignOnPolicyActionIDFirstAllOfCondition) []interface{} {
	conditionList := make([]interface{}, 0, 1)

	condition := map[string]interface{}{}

	condition["contains"] = signOnPolicyActionIDFirstAllOfCondition.GetContains()
	condition["value"] = signOnPolicyActionIDFirstAllOfCondition.GetValue()

	conditionList = append(conditionList, condition)
	return conditionList
}

func flattenActionRecoveryInner(signOnPolicyActionLoginAllOfRecovery *management.SignOnPolicyActionLoginAllOfRecovery) []interface{} {
	recoveryList := make([]interface{}, 0, 1)

	recovery := map[string]interface{}{}

	recovery["enabled"] = signOnPolicyActionLoginAllOfRecovery.GetEnabled()

	recoveryList = append(recoveryList, recovery)
	return recoveryList
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

func getActionPriority(instance management.SignOnPolicyAction) int32 {
	var priority int32
	switch instance.GetActualInstance().(type) {
	case *management.SignOnPolicyActionLogin:
		priority = instance.SignOnPolicyActionLogin.GetPriority()
	case *management.SignOnPolicyActionAgreement:
		priority = instance.SignOnPolicyActionAgreement.GetPriority()
	case *management.SignOnPolicyActionIDFirst:
		priority = instance.SignOnPolicyActionIDFirst.GetPriority()
	case *management.SignOnPolicyActionIDP:
		priority = instance.SignOnPolicyActionIDP.GetPriority()
	case *management.SignOnPolicyActionProgressiveProfiling:
		priority = instance.SignOnPolicyActionProgressiveProfiling.GetPriority()
	case *management.SignOnPolicyActionMFA:
		priority = instance.SignOnPolicyActionMFA.GetPriority()
	}

	return priority
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
