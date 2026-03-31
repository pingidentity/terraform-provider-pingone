resource "pingone_davinci_connector_instance" "imageConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "imageConnector"
  }
  name = "My awesome imageConnector"
  property {
    name  = "imageUrl"
    type  = "string"
    value = var.imageconnector_property_image_url
  }
}
