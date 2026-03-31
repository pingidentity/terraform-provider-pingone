resource "pingone_davinci_connector_instance" "mparticleConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "mparticleConnector"
  }
  name = "My awesome mparticleConnector"
  property {
    name  = "accountId"
    type  = "string"
    value = var.mparticleconnector_property_account_id
  }
  property {
    name  = "clientID"
    type  = "string"
    value = var.mparticleconnector_property_client_id
  }
  property {
    name  = "clientSecret"
    type  = "string"
    value = var.mparticleconnector_property_client_secret
  }
  property {
    name  = "displayName"
    type  = "string"
    value = var.mparticleconnector_property_display_name
  }
  property {
    name  = "jsonBody"
    type  = "string"
    value = var.mparticleconnector_property_json_body
  }
  property {
    name  = "nameOptional"
    type  = "string"
    value = var.mparticleconnector_property_name_optional
  }
  property {
    name  = "pod"
    type  = "string"
    value = var.mparticleconnector_property_pod
  }
  property {
    name  = "updatedName"
    type  = "string"
    value = var.mparticleconnector_property_updated_name
  }
  property {
    name  = "workspace"
    type  = "string"
    value = var.mparticleconnector_property_workspace
  }
}
