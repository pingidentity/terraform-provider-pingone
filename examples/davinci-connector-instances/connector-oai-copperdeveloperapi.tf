resource "pingone_davinci_connector_instance" "connector-oai-copperdeveloperapi" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connector-oai-copperdeveloperapi"
  }
  name = "My awesome connector-oai-copperdeveloperapi"
  properties = jsonencode({
    "basePath" = var.connector-oai-copperdeveloperapi_property_base_path
    "contentType" = var.connector-oai-copperdeveloperapi_property_content_type
    "xPWAccessToken" = var.connector-oai-copperdeveloperapi_property_x_p_w_access_token
    "xPWApplication" = var.connector-oai-copperdeveloperapi_property_x_p_w_application
    "xPWUserEmail" = var.connector-oai-copperdeveloperapi_property_x_p_w_user_email
  })
}
