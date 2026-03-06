resource "pingone_davinci_connector_instance" "flowConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "flowConnector"
  }
  name = "My awesome flowConnector"
  properties = jsonencode({
    "enforcedSignedToken" = var.flowconnector_property_enforced_signed_token
    "inputSchema" = var.flowconnector_property_input_schema
    "pemPublicKey" = var.flowconnector_property_pem_public_key
  })
}
