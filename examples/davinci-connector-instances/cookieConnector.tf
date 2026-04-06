resource "pingone_davinci_connector_instance" "cookieConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "cookieConnector"
  }
  name = "My awesome cookieConnector"
  property {
    name  = "claimsNameValuePairsSessionCookie"
    type  = "string"
    value = var.cookieconnector_property_claims_name_value_pairs_session_cookie
  }
  property {
    name  = "cookieDomain"
    type  = "string"
    value = var.cookieconnector_property_cookie_domain
  }
  property {
    name  = "cookieExpiresInSeconds"
    type  = "string"
    value = var.cookieconnector_property_cookie_expires_in_seconds
  }
  property {
    name  = "cookieName"
    type  = "string"
    value = var.cookieconnector_property_cookie_name
  }
  property {
    name  = "cookiePath"
    type  = "string"
    value = var.cookieconnector_property_cookie_path
  }
  property {
    name  = "cookieSameSite"
    type  = "string"
    value = var.cookieconnector_property_cookie_same_site
  }
  property {
    name  = "cookieValue"
    type  = "string"
    value = var.cookieconnector_property_cookie_value
  }
  property {
    name  = "enforceClientIP"
    type  = "string"
    value = var.cookieconnector_property_enforce_client_ip
  }
  property {
    name  = "enforceFlowIdMatch"
    type  = "string"
    value = var.cookieconnector_property_enforce_flow_id_match
  }
  property {
    name  = "hmacSigningKey"
    type  = "string"
    value = var.cookieconnector_property_hmac_signing_key
  }
  property {
    name  = "resolveToUser"
    type  = "string"
    value = var.cookieconnector_property_resolve_to_user
  }
  property {
    name  = "sessionToken"
    type  = "string"
    value = var.cookieconnector_property_session_token
  }
  property {
    name  = "setCookieClientSide"
    type  = "string"
    value = var.cookieconnector_property_set_cookie_client_side
  }
  property {
    name  = "signCookie"
    type  = "string"
    value = var.cookieconnector_property_sign_cookie
  }
  property {
    name  = "useHttpOnlyCookie"
    type  = "string"
    value = var.cookieconnector_property_use_http_only_cookie
  }
  property {
    name  = "useSecureCookie"
    type  = "string"
    value = var.cookieconnector_property_use_secure_cookie
  }
  property {
    name  = "useSessionTokenFlag"
    type  = "string"
    value = var.cookieconnector_property_use_session_token_flag
  }
}
