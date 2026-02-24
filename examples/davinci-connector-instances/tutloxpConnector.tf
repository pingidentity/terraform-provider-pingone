resource "pingone_davinci_connector_instance" "tutloxpConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "tutloxpConnector"
  }
  name = "My awesome tutloxpConnector"
  properties = jsonencode({
    "apiUrl" = var.tutloxpconnector_property_api_url
    "dppaCode" = var.tutloxpconnector_property_dppa_code
    "glbCode" = var.tutloxpconnector_property_glb_code
    "password" = var.tutloxpconnector_property_password
    "username" = var.tutloxpconnector_property_username
  })
}
