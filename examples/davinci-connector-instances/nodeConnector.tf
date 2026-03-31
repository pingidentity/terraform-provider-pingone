resource "pingone_davinci_connector_instance" "nodeConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "nodeConnector"
  }
  name = "My awesome nodeConnector"
  property {
    name  = "inputSchema"
    type  = "string"
    value = var.nodeconnector_property_input_schema
  }
  property {
    name  = "linkStartNode"
    type  = "string"
    value = var.nodeconnector_property_link_start_node
  }
  property {
    name  = "nodeInstanceId"
    type  = "string"
    value = var.nodeconnector_property_node_instance_id
  }
  property {
    name  = "outputSchema"
    type  = "string"
    value = var.nodeconnector_property_output_schema
  }
}
