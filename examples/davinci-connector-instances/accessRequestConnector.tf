resource "pingone_davinci_connector_instance" "accessRequestConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "accessRequestConnector"
  }
  name = "My awesome accessRequestConnector"
  properties = jsonencode({
    "baseURL" = var.accessrequestconnector_property_base_u_r_l
    "endUserClientId" = var.accessrequestconnector_property_end_user_client_id
    "endUserClientPrivateKey" = var.accessrequestconnector_property_end_user_client_private_key
    "realm" = var.accessrequestconnector_property_realm
    "serviceAccountId" = var.accessrequestconnector_property_service_account_id
    "serviceAccountPrivateKey" = var.accessrequestconnector_property_service_account_private_key
  })
}
