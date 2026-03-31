resource "pingone_davinci_connector_instance" "rsaConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "rsaConnector"
  }
  name = "My awesome rsaConnector"
  property {
    name  = "accessId"
    type  = "string"
    value = var.rsaconnector_property_access_id
  }
  property {
    name  = "accessKey"
    type  = "string"
    value = var.rsaconnector_property_access_key
  }
  property {
    name  = "baseUrl"
    type  = "string"
    value = var.rsaconnector_property_base_url
  }
  property {
    name  = "email"
    type  = "string"
    value = var.rsaconnector_property_email
  }
}
