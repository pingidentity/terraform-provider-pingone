resource "pingone_davinci_connector_instance" "iproovV2Connector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "iproovV2Connector"
  }
  name = "My awesome iproovV2Connector"
  property {
    name  = "activate"
    type  = "string"
    value = var.iproovv2connector_property_activate
  }
  property {
    name  = "apiKey"
    type  = "string"
    value = var.iproovv2connector_property_api_key
  }
  property {
    name  = "assuranceType"
    type  = "string"
    value = var.iproovv2connector_property_assurance_type
  }
  property {
    name  = "client"
    type  = "string"
    value = var.iproovv2connector_property_client
  }
  property {
    name  = "image"
    type  = "string"
    value = var.iproovv2connector_property_image
  }
  property {
    name  = "resource"
    type  = "string"
    value = var.iproovv2connector_property_resource
  }
  property {
    name  = "riskProfile"
    type  = "string"
    value = var.iproovv2connector_property_risk_profile
  }
  property {
    name  = "secret"
    type  = "string"
    value = var.iproovv2connector_property_secret
  }
  property {
    name  = "source"
    type  = "string"
    value = var.iproovv2connector_property_source
  }
  property {
    name  = "tenant"
    type  = "string"
    value = var.iproovv2connector_property_tenant
  }
  property {
    name  = "token"
    type  = "string"
    value = var.iproovv2connector_property_token
  }
  property {
    name  = "userId"
    type  = "string"
    value = var.iproovv2connector_property_user_id
  }
}
