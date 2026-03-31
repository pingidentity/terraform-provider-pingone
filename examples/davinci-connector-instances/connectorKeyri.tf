resource "pingone_davinci_connector_instance" "connectorKeyri" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorKeyri"
  }
  name = "My awesome connectorKeyri"
  property {
    name  = "userData"
    type  = "string"
    value = var.connectorkeyri_property_user_data
  }
  property {
    name  = "userPublicKey"
    type  = "string"
    value = var.connectorkeyri_property_user_public_key
  }
  property {
    name  = "userSignature"
    type  = "string"
    value = var.connectorkeyri_property_user_signature
  }
}
