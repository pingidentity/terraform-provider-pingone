resource "pingone_davinci_connector_instance" "microsoftDynamicsCustomerInsightsConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "microsoftDynamicsCustomerInsightsConnector"
  }
  name = "My awesome microsoftDynamicsCustomerInsightsConnector"
  properties = jsonencode({
    "baseURL" = var.microsoftdynamicscustomerinsightsconnector_property_base_u_r_l
    "clientId" = var.microsoftdynamicscustomerinsightsconnector_property_client_id
    "clientSecret" = var.microsoftdynamicscustomerinsightsconnector_property_client_secret
    "environmentName" = var.microsoftdynamicscustomerinsightsconnector_property_environment_name
    "grantType" = var.microsoftdynamicscustomerinsightsconnector_property_grant_type
    "tenant" = var.microsoftdynamicscustomerinsightsconnector_property_tenant
    "version" = var.microsoftdynamicscustomerinsightsconnector_property_version
  })
}
