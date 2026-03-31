resource "pingone_davinci_connector_instance" "pingOneAuthenticationConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "pingOneAuthenticationConnector"
  }
  name = "My awesome pingOneAuthenticationConnector"
  property {
    name  = "accessTokenClaims"
    type  = "string"
    value = var.pingoneauthenticationconnector_property_access_token_claims
  }
  property {
    name  = "acrValues"
    type  = "string"
    value = var.pingoneauthenticationconnector_property_acr_values
  }
  property {
    name  = "additionalProperties"
    type  = "string"
    value = var.pingoneauthenticationconnector_property_additional_properties
  }
  property {
    name  = "additionalPropertiesName"
    type  = "string"
    value = var.pingoneauthenticationconnector_property_additional_properties_name
  }
  property {
    name  = "application"
    type  = "string"
    value = var.pingoneauthenticationconnector_property_application
  }
  property {
    name  = "applicationId"
    type  = "string"
    value = var.pingoneauthenticationconnector_property_application_id
  }
  property {
    name  = "authType"
    type  = "string"
    value = var.pingoneauthenticationconnector_property_auth_type
  }
  property {
    name  = "authenticationContextReference"
    type  = "string"
    value = var.pingoneauthenticationconnector_property_authentication_context_reference
  }
  property {
    name  = "authenticationMethodLastUsedIn"
    type  = "string"
    value = var.pingoneauthenticationconnector_property_authentication_method_last_used_in
  }
  property {
    name  = "authenticationMethods"
    type  = "string"
    value = var.pingoneauthenticationconnector_property_authentication_methods
  }
  property {
    name  = "button"
    type  = "string"
    value = var.pingoneauthenticationconnector_property_button
  }
  property {
    name  = "checkSessionAuthenticator"
    type  = "string"
    value = var.pingoneauthenticationconnector_property_check_session_authenticator
  }
  property {
    name  = "customAuthenticationMethods"
    type  = "string"
    value = var.pingoneauthenticationconnector_property_custom_authentication_methods
  }
  property {
    name  = "customErrorFlag"
    type  = "string"
    value = var.pingoneauthenticationconnector_property_custom_error_flag
  }
  property {
    name  = "errorCode"
    type  = "string"
    value = var.pingoneauthenticationconnector_property_error_code
  }
  property {
    name  = "errorDescription"
    type  = "string"
    value = var.pingoneauthenticationconnector_property_error_description
  }
  property {
    name  = "errorMessage"
    type  = "string"
    value = var.pingoneauthenticationconnector_property_error_message
  }
  property {
    name  = "errorReason"
    type  = "string"
    value = var.pingoneauthenticationconnector_property_error_reason
  }
  property {
    name  = "idTokenClaims"
    type  = "string"
    value = var.pingoneauthenticationconnector_property_id_token_claims
  }
  property {
    name  = "idTokenHint"
    type  = "string"
    value = var.pingoneauthenticationconnector_property_id_token_hint
  }
  property {
    name  = "identifiedDeviceId"
    type  = "string"
    value = var.pingoneauthenticationconnector_property_identified_device_id
  }
  property {
    name  = "identityProvider"
    type  = "string"
    value = var.pingoneauthenticationconnector_property_identity_provider
  }
  property {
    name  = "identityProviderId"
    type  = "string"
    value = var.pingoneauthenticationconnector_property_identity_provider_id
  }
  property {
    name  = "idleTimeout"
    type  = "string"
    value = var.pingoneauthenticationconnector_property_idle_timeout
  }
  property {
    name  = "linkWithP1User"
    type  = "string"
    value = var.pingoneauthenticationconnector_property_link_with_p1_user
  }
  property {
    name  = "loginHint"
    type  = "string"
    value = var.pingoneauthenticationconnector_property_login_hint
  }
  property {
    name  = "overrideAuthState"
    type  = "string"
    value = var.pingoneauthenticationconnector_property_override_auth_state
  }
  property {
    name  = "policyPurpose"
    type  = "string"
    value = var.pingoneauthenticationconnector_property_policy_purpose
  }
  property {
    name  = "population"
    type  = "string"
    value = var.pingoneauthenticationconnector_property_population
  }
  property {
    name  = "populationId"
    type  = "string"
    value = var.pingoneauthenticationconnector_property_population_id
  }
  property {
    name  = "requestedAuthenticationContext"
    type  = "string"
    value = var.pingoneauthenticationconnector_property_requested_authentication_context
  }
  property {
    name  = "returnUrl"
    type  = "string"
    value = var.pingoneauthenticationconnector_property_return_url
  }
  property {
    name  = "scopes"
    type  = "string"
    value = var.pingoneauthenticationconnector_property_scopes
  }
  property {
    name  = "showPoweredBy"
    type  = "string"
    value = var.pingoneauthenticationconnector_property_show_powered_by
  }
  property {
    name  = "skipButtonPress"
    type  = "string"
    value = var.pingoneauthenticationconnector_property_skip_button_press
  }
  property {
    name  = "softDelete"
    type  = "string"
    value = var.pingoneauthenticationconnector_property_soft_delete
  }
  property {
    name  = "userCode"
    type  = "string"
    value = var.pingoneauthenticationconnector_property_user_code
  }
  property {
    name  = "userId"
    type  = "string"
    value = var.pingoneauthenticationconnector_property_user_id
  }
  property {
    name  = "widgetScopes"
    type  = "string"
    value = var.pingoneauthenticationconnector_property_widget_scopes
  }
}
