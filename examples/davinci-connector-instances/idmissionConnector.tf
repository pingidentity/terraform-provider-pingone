resource "pingone_davinci_connector_instance" "idmissionConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "idmissionConnector"
  }
  name = "My awesome idmissionConnector"
  properties = jsonencode({
    "authDescription" = var.idmissionconnector_property_auth_description
    "connectorName" = var.idmissionconnector_property_connector_name
    "description" = var.idmissionconnector_property_description
    "details1" = var.idmissionconnector_property_details1
    "details2" = var.idmissionconnector_property_details2
    "iconUrl" = var.idmissionconnector_property_icon_url
    "iconUrlPng" = var.idmissionconnector_property_icon_url_png
    "loginId" = var.idmissionconnector_property_login_id
    "merchantId" = var.idmissionconnector_property_merchant_id
    "password" = var.idmissionconnector_property_password
    "productId" = var.idmissionconnector_property_product_id
    "productName" = var.idmissionconnector_property_product_name
    "showCredAddedOn" = var.idmissionconnector_property_show_cred_added_on
    "showCredAddedVia" = var.idmissionconnector_property_show_cred_added_via
    "title" = var.idmissionconnector_property_title
    "toolTip" = var.idmissionconnector_property_tool_tip
    "url" = var.idmissionconnector_property_url
  })
}
