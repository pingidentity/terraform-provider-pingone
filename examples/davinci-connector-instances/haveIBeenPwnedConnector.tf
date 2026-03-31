resource "pingone_davinci_connector_instance" "haveIBeenPwnedConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "haveIBeenPwnedConnector"
  }
  name = "My awesome haveIBeenPwnedConnector"
  property {
    name  = "apiKey"
    type  = "string"
    value = var.haveibeenpwnedconnector_property_api_key
  }
  property {
    name  = "apiUrl"
    type  = "string"
    value = var.haveibeenpwnedconnector_property_api_url
  }
  property {
    name  = "domain"
    type  = "string"
    value = var.haveibeenpwnedconnector_property_domain
  }
  property {
    name  = "email"
    type  = "string"
    value = var.haveibeenpwnedconnector_property_email
  }
  property {
    name  = "includeUnverified"
    type  = "string"
    value = var.haveibeenpwnedconnector_property_include_unverified
  }
  property {
    name  = "password"
    type  = "string"
    value = var.haveibeenpwnedconnector_property_password
  }
  property {
    name  = "site"
    type  = "string"
    value = var.haveibeenpwnedconnector_property_site
  }
  property {
    name  = "truncateResponse"
    type  = "string"
    value = var.haveibeenpwnedconnector_property_truncate_response
  }
  property {
    name  = "userAgent"
    type  = "string"
    value = var.haveibeenpwnedconnector_property_user_agent
  }
}
