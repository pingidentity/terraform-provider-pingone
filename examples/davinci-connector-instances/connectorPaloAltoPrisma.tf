resource "pingone_davinci_connector_instance" "connectorPaloAltoPrisma" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorPaloAltoPrisma"
  }
  name = "My awesome connectorPaloAltoPrisma"
  properties = jsonencode({
    "baseURL" = var.base_url
    "prismaPassword" = var.connectorpaloaltoprisma_property_prisma_password
    "prismaUsername" = var.connectorpaloaltoprisma_property_prisma_username
  })
}
