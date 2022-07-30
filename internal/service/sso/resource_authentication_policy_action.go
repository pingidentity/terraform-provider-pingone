package sso

import (
	"fmt"
	"log"
	"sort"

	"github.com/patrickcping/pingone-go-sdk-v2/management"
)

func expandSOPAction(d interface{}, sopPriority int32) (*management.SignOnPolicyAction, error) {

	signOnPolicyAction := &management.SignOnPolicyAction{}

	actionType := d.(map[string]interface{})["action_type"].(string)

	var err error

	log.Printf("Action type: %s", actionType)
	switch actionType {
	case string(management.ENUMSIGNONPOLICYTYPE_AGREEMENT):
		signOnPolicyAction.SignOnPolicyActionAgreement, err = expandSOPActionAgreement(d, sopPriority)
	case string(management.ENUMSIGNONPOLICYTYPE_IDENTIFIER_FIRST):
		signOnPolicyAction.SignOnPolicyActionIDFirst = expandSOPActionIDFirst(d, sopPriority)
	case string(management.ENUMSIGNONPOLICYTYPE_IDENTITY_PROVIDER):
		signOnPolicyAction.SignOnPolicyActionIDP, err = expandSOPActionIDP(d, sopPriority)
	case string(management.ENUMSIGNONPOLICYTYPE_LOGIN):
		signOnPolicyAction.SignOnPolicyActionLogin = expandSOPActionLogin(d, sopPriority)
	case string(management.ENUMSIGNONPOLICYTYPE_MULTI_FACTOR_AUTHENTICATION):
		signOnPolicyAction.SignOnPolicyActionMFA = expandSOPActionMFA(d, sopPriority)
	case string(management.ENUMSIGNONPOLICYTYPE_PROGRESSIVE_PROFILING):
		signOnPolicyAction.SignOnPolicyActionProgressiveProfiling, err = expandSOPActionProgressiveProfiling(d, sopPriority)
	default:
		return nil, fmt.Errorf("Policy action %s not supported", actionType)
	}

	if err != nil {
		return nil, fmt.Errorf("Error parsing Authentication policy action: %v", err)
	}

	return signOnPolicyAction, nil
}

func expandSOPActionProgressiveProfiling(d interface{}, sopPriority int32) (*management.SignOnPolicyActionProgressiveProfiling, error) {

	if v, ok := d.(map[string]interface{})["progressive_profiling"].([]interface{}); ok && v != nil && len(v) > 0 && v[0] != nil {
		vp := v[0].(map[string]interface{})

		sopActionType := management.NewSignOnPolicyActionProgressiveProfiling(
			sopPriority,
			management.ENUMSIGNONPOLICYTYPE_PROGRESSIVE_PROFILING,
			expandSOPActionAttributes(v[0].([]interface{})),
			vp["prevent_multiple_prompts_per_flow"].(bool),
			int32(vp["prompt_interval_seconds"].(int)),
			vp["prompt_text"].(string),
		)

		if vc, ok := d.(map[string]interface{})["conditions"].([]interface{}); ok && vc != nil && len(vc) > 0 && vc[0] != nil {
			sopActionType.SetConditions(expandSOPActionConditions(vc))
		}

		return sopActionType, nil

	}

	return nil, fmt.Errorf("Block `progressive_profiling` with `prevent_multiple_prompts_per_flow`, `prompt_interval_seconds` and `prompt_text` must be defined when using the progressive profiling action type")
}

func expandSOPActionMFA(d interface{}, sopPriority int32) *management.SignOnPolicyActionMFA {

	sopActionType := management.NewSignOnPolicyActionMFA(
		sopPriority,
		management.ENUMSIGNONPOLICYTYPE_MULTI_FACTOR_AUTHENTICATION,
	)

	if vc, ok := d.(map[string]interface{})["conditions"].([]interface{}); ok && vc != nil && len(vc) > 0 && vc[0] != nil {
		sopActionType.SetConditions(expandSOPActionConditions(vc))
	}

	if v, ok := d.(map[string]interface{})["mfa"].([]interface{}); ok && v != nil && len(v) > 0 && v[0] != nil {
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

	return sopActionType
}

func expandSOPActionLogin(d interface{}, sopPriority int32) *management.SignOnPolicyActionLogin {

	sopActionType := management.NewSignOnPolicyActionLogin(
		sopPriority,
		management.ENUMSIGNONPOLICYTYPE_LOGIN,
	)

	if v, ok := d.(map[string]interface{})["conditions"].([]interface{}); ok && v != nil && len(v) > 0 && v[0] != nil {
		sopActionType.SetConditions(expandSOPActionConditions(v))
	}

	// block is optional
	if v, ok := d.(map[string]interface{})["login"].([]interface{}); ok && v != nil && len(v) > 0 && v[0] != nil {
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
				registrationObj.External.SetHref(externalHref.(string))
			}

			sopActionType.SetRegistration(registrationObj)

		}

		if v1, ok := vp["social_providers"].([]interface{}); ok && v1 != nil && len(v1) > 0 && v1[0] != nil {
			sopActionType.SetSocialProviders(expandSOPActionSocialProviders(v1))
		}

	}

	return sopActionType
}

func expandSOPActionIDP(d interface{}, sopPriority int32) (*management.SignOnPolicyActionIDP, error) {

	if v, ok := d.(map[string]interface{})["identity_provider"].([]interface{}); ok && v != nil && len(v) > 0 && v[0] != nil {
		vp := v[0].(map[string]interface{})

		sopActionType := management.NewSignOnPolicyActionIDP(
			sopPriority,
			management.ENUMSIGNONPOLICYTYPE_IDENTITY_PROVIDER,
			*management.NewSignOnPolicyActionIDPAllOfIdentityProvider(vp["identity_provider_id"].(string)),
		)

		if vc, ok := d.(map[string]interface{})["conditions"].([]interface{}); ok && vc != nil && len(vc) > 0 && vc[0] != nil {
			sopActionType.SetConditions(expandSOPActionConditions(vc))
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

	return nil, fmt.Errorf("Block `identity_provider` with `identity_provider_id` must be defined when using the identity provider action type")
}

func expandSOPActionIDFirst(d interface{}, sopPriority int32) *management.SignOnPolicyActionIDFirst {

	sopActionType := management.NewSignOnPolicyActionIDFirst(
		sopPriority,
		management.ENUMSIGNONPOLICYTYPE_IDENTIFIER_FIRST,
	)

	if vc, ok := d.(map[string]interface{})["conditions"].([]interface{}); ok && vc != nil && len(vc) > 0 && vc[0] != nil {
		sopActionType.SetConditions(expandSOPActionConditions(vc))
	}

	if v, ok := d.(map[string]interface{})["identifier_first"].([]interface{}); ok && v != nil && len(v) > 0 && v[0] != nil {
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
				registrationObj.External.SetHref(externalHref.(string))
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

	return sopActionType

}

func expandSOPActionAgreement(d interface{}, sopPriority int32) (*management.SignOnPolicyActionAgreement, error) {

	if v, ok := d.(map[string]interface{})["agreement"].([]interface{}); ok && v != nil && len(v) > 0 && v[0] != nil {
		vp := v[0].(map[string]interface{})

		sopActionType := management.NewSignOnPolicyActionAgreement(
			sopPriority,
			management.ENUMSIGNONPOLICYTYPE_AGREEMENT,
			*management.NewSignOnPolicyActionAgreementAllOfAgreement(vp["agreement_id"].(string)),
		)

		if vc, ok := d.(map[string]interface{})["conditions"].([]interface{}); ok && vc != nil && len(vc) > 0 && vc[0] != nil {
			sopActionType.SetConditions(expandSOPActionConditions(vc))
		}

		return sopActionType, nil

	}

	return nil, fmt.Errorf("Block `agreement` with `agreement_id` must be defined when using the agreement action type")
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

func expandSOPActionAttributes(items []interface{}) []management.SignOnPolicyActionProgressiveProfilingAllOfAttributes {

	var attributes []management.SignOnPolicyActionProgressiveProfilingAllOfAttributes

	for _, item := range items {

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

func expandSOPActionConditions(condition []interface{}) management.SignOnPolicyActionCommonConditions {

	sopConditions := *management.NewSignOnPolicyActionCommonConditions()

	ipRange := condition[0].(map[string]interface{})["ip_range"]
	if ipRange != nil {
		sopConditions.SetIpRange(ipRange.(string))
	}

	sessionLength := condition[0].(map[string]interface{})["action_session_length_mins"]
	if sessionLength != nil {
		sopConditions.SetSecondsSince(int32(sessionLength.(int)))
	}

	return sopConditions
}

func flattenSOPActions(actions []management.SignOnPolicyAction) ([]interface{}, error) {

	// Sort by priorty first
	sort.Slice(actions, func(i, j int) bool {
		return actions[i].GetActualInstance().(map[string]interface{})["priority"].(int32) < actions[j].GetActualInstance().(map[string]interface{})["priority"].(int32)
	})

	// Set the return var
	sopActions := make([]interface{}, 0)

	// Loop the ordered slice
	for _, action := range actions {

		actionMap := map[string]interface{}{}

		match := 0

		if action.SignOnPolicyActionAgreement != nil {

			actionMap["action_type"] = string(management.ENUMSIGNONPOLICYTYPE_AGREEMENT)

			if v, ok := action.SignOnPolicyActionAgreement.GetConditionsOk(); ok {
				actionMap["conditions"] = flattenConditions(*v)
			}
			actionMap["agreement"] = flattenActionAgreement(action.SignOnPolicyActionAgreement)

			match++
		}

		if action.SignOnPolicyActionIDFirst != nil {

			actionMap["action_type"] = string(management.ENUMSIGNONPOLICYTYPE_IDENTIFIER_FIRST)

			if v, ok := action.SignOnPolicyActionIDFirst.GetConditionsOk(); ok {
				actionMap["conditions"] = flattenConditions(*v)
			}
			actionMap["identifier_first"] = flattenActionIDFirst(action.SignOnPolicyActionIDFirst)

			match++
		}

		if action.SignOnPolicyActionIDP != nil {

			actionMap["action_type"] = string(management.ENUMSIGNONPOLICYTYPE_IDENTITY_PROVIDER)

			if v, ok := action.SignOnPolicyActionIDP.GetConditionsOk(); ok {
				actionMap["conditions"] = flattenConditions(*v)
			}
			actionMap["identity_provider"] = flattenActionIDP(action.SignOnPolicyActionIDP)

			match++
		}

		if action.SignOnPolicyActionLogin != nil {

			actionMap["action_type"] = string(management.ENUMSIGNONPOLICYTYPE_LOGIN)

			if v, ok := action.SignOnPolicyActionLogin.GetConditionsOk(); ok {
				actionMap["conditions"] = flattenConditions(*v)
			}
			actionMap["login"] = flattenActionLogin(action.SignOnPolicyActionLogin)

			match++
		}

		if action.SignOnPolicyActionMFA != nil {

			actionMap["action_type"] = string(management.ENUMSIGNONPOLICYTYPE_MULTI_FACTOR_AUTHENTICATION)

			if v, ok := action.SignOnPolicyActionMFA.GetConditionsOk(); ok {
				actionMap["conditions"] = flattenConditions(*v)
			}
			actionMap["mfa"] = flattenActionMFA(action.SignOnPolicyActionMFA)

			match++
		}

		if action.SignOnPolicyActionProgressiveProfiling != nil {

			actionMap["action_type"] = string(management.ENUMSIGNONPOLICYTYPE_PROGRESSIVE_PROFILING)

			if v, ok := action.SignOnPolicyActionProgressiveProfiling.GetConditionsOk(); ok {
				actionMap["conditions"] = flattenConditions(*v)
			}
			actionMap["progressive_profiling"] = flattenActionProgressiveProfiling(action.SignOnPolicyActionProgressiveProfiling)

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

func flattenConditions(signOnPolicyActionCommonConditions management.SignOnPolicyActionCommonConditions) []interface{} {

	conditionsList := make([]interface{}, 0, 1)

	conditions := map[string]interface{}{}

	if v, ok := signOnPolicyActionCommonConditions.GetIpRangeOk(); ok {
		conditions["ip_range"] = v
	} else {
		conditions["ip_range"] = nil
	}

	if v, ok := signOnPolicyActionCommonConditions.GetSecondsSinceOk(); ok {
		conditions["action_session_length_mins"] = v
	} else {
		conditions["action_session_length_mins"] = nil
	}

	conditionsList = append(conditionsList, conditions)
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
