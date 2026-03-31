resource "pingone_davinci_connector_instance" "connectorRandomUserMe" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorRandomUserMe"
  }
  name = "My awesome connectorRandomUserMe"
  property {
    name  = "gender"
    type  = "string"
    value = var.connectorrandomuserme_property_gender
  }
  property {
    name  = "seed"
    type  = "string"
    value = var.connectorrandomuserme_property_seed
  }
}
