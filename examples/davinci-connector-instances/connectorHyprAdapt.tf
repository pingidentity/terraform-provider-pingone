resource "pingone_davinci_connector_instance" "connectorHyprAdapt" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorHyprAdapt"
  }
  name = "My awesome connectorHyprAdapt"
  property {
    name  = "accessToken"
    type  = "string"
    value = var.connectorhypradapt_property_access_token
  }
  property {
    name  = "domain"
    type  = "string"
    value = var.connectorhypradapt_property_domain
  }
  property {
    name  = "dynamicPolicyData"
    type  = "string"
    value = var.connectorhypradapt_property_dynamic_policy_data
  }
  property {
    name  = "policyContent"
    type  = "string"
    value = var.connectorhypradapt_property_policy_content
  }
  property {
    name  = "policyData"
    type  = "string"
    value = var.connectorhypradapt_property_policy_data
  }
  property {
    name  = "policyId"
    type  = "string"
    value = var.connectorhypradapt_property_policy_id
  }
  property {
    name  = "username"
    type  = "string"
    value = var.connectorhypradapt_property_username
  }
}
