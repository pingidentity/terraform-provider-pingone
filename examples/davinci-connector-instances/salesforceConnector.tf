resource "pingone_davinci_connector_instance" "salesforceConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "salesforceConnector"
  }
  name = "My awesome salesforceConnector"
  properties = jsonencode({
    "adminUsername" = var.salesforceconnector_property_admin_username
    "consumerKey" = var.salesforceconnector_property_consumer_key
    "domainName" = var.salesforceconnector_property_domain_name
    "environment" = var.salesforceconnector_property_environment
    "privateKey" = var.salesforceconnector_property_private_key
  })
}
