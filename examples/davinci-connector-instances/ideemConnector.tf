resource "pingone_davinci_connector_instance" "ideemConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "ideemConnector"
  }
  name = "My awesome ideemConnector"
  property {
    name  = "applicationEnvironment"
    type  = "string"
    value = var.ideemconnector_property_application_environment
  }
  property {
    name  = "applicationId"
    type  = "string"
    value = var.ideemconnector_property_application_id
  }
  property {
    name  = "button"
    type  = "string"
    value = var.ideemconnector_property_button
  }
  property {
    name  = "companyLogo"
    type  = "string"
    value = var.ideemconnector_property_company_logo
  }
  property {
    name  = "companyName"
    type  = "string"
    value = var.ideemconnector_property_company_name
  }
  property {
    name  = "hostURL"
    type  = "string"
    value = var.ideemconnector_property_host_url
  }
  property {
    name  = "pollingText"
    type  = "string"
    value = var.ideemconnector_property_polling_text
  }
  property {
    name  = "timeoutMs"
    type  = "string"
    value = var.ideemconnector_property_timeout_ms
  }
  property {
    name  = "tokenForm"
    type  = "string"
    value = var.ideemconnector_property_token_form
  }
  property {
    name  = "userIdentifier"
    type  = "string"
    value = var.ideemconnector_property_user_identifier
  }
  property {
    name  = "validateTokenApiKey"
    type  = "string"
    value = var.ideemconnector_property_validate_token_api_key
  }
  property {
    name  = "zsmClientSdkApiKey"
    type  = "string"
    value = var.ideemconnector_property_zsm_client_sdk_api_key
  }
}
