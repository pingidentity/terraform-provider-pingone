resource "pingone_davinci_connector_instance" "skOpenIdConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "skOpenIdConnector"
  }
  name = "My awesome skOpenIdConnector"
  property {
    name  = "accessTokenExpiry"
    type  = "string"
    value = var.skopenidconnector_property_access_token_expiry
  }
  property {
    name  = "claimsNameValuePairs"
    type  = "string"
    value = var.skopenidconnector_property_claims_name_value_pairs
  }
  property {
    name  = "claimsNameValuePairsAccessToken"
    type  = "string"
    value = var.skopenidconnector_property_claims_name_value_pairs_access_token
  }
  property {
    name  = "claimsNameValuePairsSessionToken"
    type  = "string"
    value = var.skopenidconnector_property_claims_name_value_pairs_session_token
  }
  property {
    name  = "createAccessTokenFlag"
    type  = "string"
    value = var.skopenidconnector_property_create_access_token_flag
  }
  property {
    name  = "createIdTokenFlag"
    type  = "string"
    value = var.skopenidconnector_property_create_id_token_flag
  }
  property {
    name  = "createSessionTokenFlag"
    type  = "string"
    value = var.skopenidconnector_property_create_session_token_flag
  }
  property {
    name  = "customErrorFlag"
    type  = "string"
    value = var.skopenidconnector_property_custom_error_flag
  }
  property {
    name  = "customScopes"
    type  = "string"
    value = var.skopenidconnector_property_custom_scopes
  }
  property {
    name  = "customScopesFlag"
    type  = "string"
    value = var.skopenidconnector_property_custom_scopes_flag
  }
  property {
    name  = "customScopesSeparateField"
    type  = "string"
    value = var.skopenidconnector_property_custom_scopes_separate_field
  }
  property {
    name  = "customScopesSeparateFieldName"
    type  = "string"
    value = var.skopenidconnector_property_custom_scopes_separate_field_name
  }
  property {
    name  = "errorCode"
    type  = "string"
    value = var.skopenidconnector_property_error_code
  }
  property {
    name  = "errorDescription"
    type  = "string"
    value = var.skopenidconnector_property_error_description
  }
  property {
    name  = "errorMessage"
    type  = "string"
    value = var.skopenidconnector_property_error_message
  }
  property {
    name  = "errorOnExpiry"
    type  = "string"
    value = var.skopenidconnector_property_error_on_expiry
  }
  property {
    name  = "errorReason"
    type  = "string"
    value = var.skopenidconnector_property_error_reason
  }
  property {
    name  = "genericToken"
    type  = "string"
    value = var.skopenidconnector_property_generic_token
  }
  property {
    name  = "idTokenExpiry"
    type  = "string"
    value = var.skopenidconnector_property_id_token_expiry
  }
  property {
    name  = "publicKeyId"
    type  = "string"
    value = var.skopenidconnector_property_public_key_id
  }
  property {
    name  = "publicKeyJWTEndpointURL"
    type  = "string"
    value = var.skopenidconnector_property_public_key_jwtendpoint_url
  }
  property {
    name  = "publicKeyPEMContents"
    type  = "string"
    value = var.skopenidconnector_property_public_key_pemcontents
  }
  property {
    name  = "publicKeyType"
    type  = "string"
    value = var.skopenidconnector_property_public_key_type
  }
  property {
    name  = "resolveToUser"
    type  = "string"
    value = var.skopenidconnector_property_resolve_to_user
  }
  property {
    name  = "secretKey"
    type  = "string"
    value = var.skopenidconnector_property_secret_key
  }
  property {
    name  = "sessionToken"
    type  = "string"
    value = var.skopenidconnector_property_session_token
  }
  property {
    name  = "sessionTokenExpiry"
    type  = "string"
    value = var.skopenidconnector_property_session_token_expiry
  }
  property {
    name  = "sessionTokenLocation"
    type  = "string"
    value = var.skopenidconnector_property_session_token_location
  }
  property {
    name  = "sessionTokenName"
    type  = "string"
    value = var.skopenidconnector_property_session_token_name
  }
  property {
    name  = "shadowUserNotPresentFlag"
    type  = "string"
    value = var.skopenidconnector_property_shadow_user_not_present_flag
  }
  property {
    name  = "validAlgorithms"
    type  = "string"
    value = var.skopenidconnector_property_valid_algorithms
  }
  property {
    name  = "validAudiences"
    type  = "string"
    value = var.skopenidconnector_property_valid_audiences
  }
  property {
    name  = "validIssuers"
    type  = "string"
    value = var.skopenidconnector_property_valid_issuers
  }
  property {
    name  = "validSubjects"
    type  = "string"
    value = var.skopenidconnector_property_valid_subjects
  }
}
