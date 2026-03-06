resource "pingone_davinci_connector_instance" "idmissionOidcConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "idmissionOidcConnector"
  }
  name = "My awesome idmissionOidcConnector"
  properties = jsonencode({
    "customAuth" = var.idmissionoidcconnector_property_custom_auth
  })
}
