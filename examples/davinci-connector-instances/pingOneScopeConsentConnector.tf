resource "pingone_davinci_connector_instance" "pingOneScopeConsentConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "pingOneScopeConsentConnector"
  }
  name = "My awesome pingOneScopeConsentConnector"
  property {
    name  = "appConsentHtmlConfig"
    type  = "string"
    value = var.pingonescopeconsentconnector_property_app_consent_html_config
  }
  property {
    name  = "applicationConsentScopes"
    type  = "string"
    value = var.pingonescopeconsentconnector_property_application_consent_scopes
  }
  property {
    name  = "applicationIdentifier"
    type  = "string"
    value = var.pingonescopeconsentconnector_property_application_identifier
  }
  property {
    name  = "applicationIdentifierIDForAppConsent"
    type  = "string"
    value = var.pingonescopeconsentconnector_property_application_identifier_idfor_app_consent
  }
  property {
    name  = "applicationIdentifierNameForAppConsent"
    type  = "string"
    value = var.pingonescopeconsentconnector_property_application_identifier_name_for_app_consent
  }
  property {
    name  = "clientId"
    type  = "string"
    value = var.pingone_worker_app_client_id
  }
  property {
    name  = "clientSecret"
    type  = "string"
    value = var.pingone_worker_app_client_secret
  }
  property {
    name  = "consentIdentifier"
    type  = "string"
    value = var.pingonescopeconsentconnector_property_consent_identifier
  }
  property {
    name  = "consentResult"
    type  = "string"
    value = var.pingonescopeconsentconnector_property_consent_result
  }
  property {
    name  = "envId"
    type  = "string"
    value = var.pingone_worker_app_environment_id
  }
  property {
    name  = "matchApplicationAttribute"
    type  = "string"
    value = var.pingonescopeconsentconnector_property_match_application_attribute
  }
  property {
    name  = "matchApplicationAttributeForAppConsent"
    type  = "string"
    value = var.pingonescopeconsentconnector_property_match_application_attribute_for_app_consent
  }
  property {
    name  = "matchConsentAttribute"
    type  = "string"
    value = var.pingonescopeconsentconnector_property_match_consent_attribute
  }
  property {
    name  = "matchUserAttribute"
    type  = "string"
    value = var.pingonescopeconsentconnector_property_match_user_attribute
  }
  property {
    name  = "promptConsent"
    type  = "string"
    value = var.pingonescopeconsentconnector_property_prompt_consent
  }
  property {
    name  = "region"
    type  = "string"
    value = var.pingonescopeconsentconnector_property_region
  }
  property {
    name  = "scopes"
    type  = "string"
    value = var.pingonescopeconsentconnector_property_scopes
  }
  property {
    name  = "scopesUnconditional"
    type  = "string"
    value = var.pingonescopeconsentconnector_property_scopes_unconditional
  }
  property {
    name  = "scopesUnconditionalRequired"
    type  = "string"
    value = var.pingonescopeconsentconnector_property_scopes_unconditional_required
  }
  property {
    name  = "userIdentifier"
    type  = "string"
    value = var.pingonescopeconsentconnector_property_user_identifier
  }
}
