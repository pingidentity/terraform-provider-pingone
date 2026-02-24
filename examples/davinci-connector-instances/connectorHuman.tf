resource "pingone_davinci_connector_instance" "connectorHuman" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorHuman"
  }
  name = "My awesome connectorHuman"
  properties = jsonencode({
    "humanAuthenticationToken" = var.connectorhuman_property_human_authentication_token
    "humanCustomerID" = var.human_customer_id
    "humanPolicyName" = var.connectorhuman_property_human_policy_name
  })
}
