resource "pingone_davinci_connector_instance" "connectorPaloAltoPrisma" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorPaloAltoPrisma"
  }
  name = "My awesome connectorPaloAltoPrisma"
  property {
    name  = "accessKeysAllowed"
    type  = "string"
    value = var.connectorpaloaltoprisma_property_access_keys_allowed
  }
  property {
    name  = "baseURL"
    type  = "string"
    value = var.base_url
  }
  property {
    name  = "defaultRoleId"
    type  = "string"
    value = var.connectorpaloaltoprisma_property_default_role_id
  }
  property {
    name  = "email"
    type  = "string"
    value = var.connectorpaloaltoprisma_property_email
  }
  property {
    name  = "firstName"
    type  = "string"
    value = var.connectorpaloaltoprisma_property_first_name
  }
  property {
    name  = "id"
    type  = "string"
    value = var.connectorpaloaltoprisma_property_id
  }
  property {
    name  = "lastName"
    type  = "string"
    value = var.connectorpaloaltoprisma_property_last_name
  }
  property {
    name  = "prismaPassword"
    type  = "string"
    value = var.connectorpaloaltoprisma_property_prisma_password
  }
  property {
    name  = "prismaUsername"
    type  = "string"
    value = var.connectorpaloaltoprisma_property_prisma_username
  }
  property {
    name  = "roleIds"
    type  = "string"
    value = var.connectorpaloaltoprisma_property_role_ids
  }
  property {
    name  = "timeZone"
    type  = "string"
    value = var.connectorpaloaltoprisma_property_time_zone
  }
}
