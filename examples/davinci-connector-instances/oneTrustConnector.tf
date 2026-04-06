resource "pingone_davinci_connector_instance" "oneTrustConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "oneTrustConnector"
  }
  name = "My awesome oneTrustConnector"
  property {
    name  = "apiToken"
    type  = "string"
    value = var.onetrustconnector_property_api_token
  }
  property {
    name  = "applicationDomain"
    type  = "string"
    value = var.onetrustconnector_property_application_domain
  }
  property {
    name  = "clientId"
    type  = "string"
    value = var.onetrustconnector_property_client_id
  }
  property {
    name  = "clientSecret"
    type  = "string"
    value = var.onetrustconnector_property_client_secret
  }
  property {
    name  = "dataElements"
    type  = "string"
    value = var.onetrustconnector_property_data_elements
  }
  property {
    name  = "includeNotGiven"
    type  = "string"
    value = var.onetrustconnector_property_include_not_given
  }
  property {
    name  = "privacyPortalDomain"
    type  = "string"
    value = var.onetrustconnector_property_privacy_portal_domain
  }
  property {
    name  = "purposes"
    type  = "string"
    value = var.onetrustconnector_property_purposes
  }
  property {
    name  = "receiptId"
    type  = "string"
    value = var.onetrustconnector_property_receipt_id
  }
  property {
    name  = "userIdentifier"
    type  = "string"
    value = var.onetrustconnector_property_user_identifier
  }
}
