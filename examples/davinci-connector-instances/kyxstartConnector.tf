resource "pingone_davinci_connector_instance" "kyxstartConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "kyxstartConnector"
  }
  name = "My awesome kyxstartConnector"
  property {
    name  = "body"
    type  = "string"
    value = var.kyxstartconnector_property_body
  }
  property {
    name  = "clientId"
    type  = "string"
    value = var.kyxstartconnector_property_client_id
  }
  property {
    name  = "clientSecret"
    type  = "string"
    value = var.kyxstartconnector_property_client_secret
  }
  property {
    name  = "tenantName"
    type  = "string"
    value = var.kyxstartconnector_property_tenant_name
  }
}
