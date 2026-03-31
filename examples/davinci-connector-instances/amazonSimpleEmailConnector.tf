resource "pingone_davinci_connector_instance" "amazonSimpleEmailConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "amazonSimpleEmailConnector"
  }
  name = "My awesome amazonSimpleEmailConnector"
  property {
    name  = "awsAccessKey"
    type  = "string"
    value = var.amazonsimpleemailconnector_property_aws_access_key
  }
  property {
    name  = "awsAccessSecret"
    type  = "string"
    value = var.amazonsimpleemailconnector_property_aws_access_secret
  }
  property {
    name  = "awsRegion"
    type  = "string"
    value = "eu-west-1"
  }
  property {
    name  = "body"
    type  = "string"
    value = var.amazonsimpleemailconnector_property_body
  }
  property {
    name  = "continueFlowLinkEnabled"
    type  = "string"
    value = var.amazonsimpleemailconnector_property_continue_flow_link_enabled
  }
  property {
    name  = "from"
    type  = "string"
    value = "support@bxretail.org"
  }
  property {
    name  = "subject"
    type  = "string"
    value = var.amazonsimpleemailconnector_property_subject
  }
  property {
    name  = "to"
    type  = "string"
    value = var.amazonsimpleemailconnector_property_to
  }
}
