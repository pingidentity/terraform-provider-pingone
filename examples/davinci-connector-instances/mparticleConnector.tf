resource "pingone_davinci_connector_instance" "mparticleConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "mparticleConnector"
  }
  name = "My awesome mparticleConnector"
  properties = jsonencode({
    "clientID" = var.mparticleconnector_property_client_i_d
    "clientSecret" = var.mparticleconnector_property_client_secret
    "pod" = var.mparticleconnector_property_pod
  })
}
