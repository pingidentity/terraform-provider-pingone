resource "pingone_davinci_connector_instance" "adobeExperienceManagerConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "adobeExperienceManagerConnector"
  }
  name = "My awesome adobeExperienceManagerConnector"
  property {
    name  = "clientId"
    type  = "string"
    value = var.adobeexperiencemanagerconnector_property_client_id
  }
  property {
    name  = "clientSecret"
    type  = "string"
    value = var.adobeexperiencemanagerconnector_property_client_secret
  }
  property {
    name  = "nameSpace"
    type  = "string"
    value = var.adobeexperiencemanagerconnector_property_name_space
  }
  property {
    name  = "orgId"
    type  = "string"
    value = var.adobeexperiencemanagerconnector_property_org_id
  }
  property {
    name  = "requestId"
    type  = "string"
    value = var.adobeexperiencemanagerconnector_property_request_id
  }
  property {
    name  = "sandboxName"
    type  = "string"
    value = var.adobeexperiencemanagerconnector_property_sandbox_name
  }
  property {
    name  = "transactionId"
    type  = "string"
    value = var.adobeexperiencemanagerconnector_property_transaction_id
  }
  property {
    name  = "value"
    type  = "string"
    value = var.adobeexperiencemanagerconnector_property_value
  }
}
