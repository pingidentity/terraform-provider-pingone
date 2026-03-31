resource "pingone_davinci_connector_instance" "idemiaConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "idemiaConnector"
  }
  name = "My awesome idemiaConnector"
  property {
    name  = "apikey"
    type  = "string"
    value = var.idemiaconnector_property_apikey
  }
  property {
    name  = "baseUrl"
    type  = "string"
    value = var.idemiaconnector_property_base_url
  }
  property {
    name  = "customCSS"
    type  = "string"
    value = var.idemiaconnector_property_custom_css
  }
  property {
    name  = "customHTML"
    type  = "string"
    value = var.idemiaconnector_property_custom_html
  }
  property {
    name  = "customScript"
    type  = "string"
    value = var.idemiaconnector_property_custom_script
  }
  property {
    name  = "documentId"
    type  = "string"
    value = var.idemiaconnector_property_document_id
  }
  property {
    name  = "htmlConfig"
    type  = "string"
    value = var.idemiaconnector_property_html_config
  }
  property {
    name  = "identitiesId"
    type  = "string"
    value = var.idemiaconnector_property_identities_id
  }
  property {
    name  = "livenessMode"
    type  = "string"
    value = var.idemiaconnector_property_liveness_mode
  }
  property {
    name  = "portraitConsent"
    type  = "string"
    value = var.idemiaconnector_property_portrait_consent
  }
  property {
    name  = "portraitId"
    type  = "string"
    value = var.idemiaconnector_property_portrait_id
  }
  property {
    name  = "portraitValidDate"
    type  = "string"
    value = var.idemiaconnector_property_portrait_valid_date
  }
  property {
    name  = "useCustomScreens"
    type  = "string"
    value = var.idemiaconnector_property_use_custom_screens
  }
}
