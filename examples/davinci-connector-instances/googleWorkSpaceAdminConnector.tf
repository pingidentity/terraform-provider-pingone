resource "pingone_davinci_connector_instance" "googleWorkSpaceAdminConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "googleWorkSpaceAdminConnector"
  }
  name = "My awesome googleWorkSpaceAdminConnector"
  properties = jsonencode({
    "iss" = var.googleworkspaceadminconnector_property_iss
    "privateKey" = var.googleworkspaceadminconnector_property_private_key
    "sub" = var.googleworkspaceadminconnector_property_sub
  })
}
