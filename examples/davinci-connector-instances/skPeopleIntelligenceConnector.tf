resource "pingone_davinci_connector_instance" "skPeopleIntelligenceConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "skPeopleIntelligenceConnector"
  }
  name = "My awesome skPeopleIntelligenceConnector"
  properties = jsonencode({
    "authUrl" = var.skpeopleintelligenceconnector_property_auth_url
    "clientId" = var.skpeopleintelligenceconnector_property_client_id
    "clientSecret" = var.skpeopleintelligenceconnector_property_client_secret
    "dppa" = var.skpeopleintelligenceconnector_property_dppa
    "glba" = var.skpeopleintelligenceconnector_property_glba
    "searchUrl" = var.skpeopleintelligenceconnector_property_search_url
  })
}
