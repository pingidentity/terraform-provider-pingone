resource "pingone_davinci_connector_instance" "jiraConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "jiraConnector"
  }
  name = "My awesome jiraConnector"
  properties = jsonencode({
    "apiKey" = var.jiraconnector_property_api_key
    "apiUrl" = var.jiraconnector_property_api_url
    "email" = var.jiraconnector_property_email
  })
}
