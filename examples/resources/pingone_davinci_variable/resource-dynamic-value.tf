resource "pingone_davinci_flow" "my_awesome_main_flow" {
  environment_id = var.environment_id
  name           = "My Awesome Main Flow"

  graph_data = {
    // edges, nodes, etc.
  }
}

resource "pingone_davinci_variable" "my_awesome_usercontext_variable" {
  environment_id = var.environment_id

  context   = "flow"
  data_type = "string"
  mutable   = true
  name      = "usercontext"

  display_name = "User Context Variable"

  flow = {
    id = pingone_davinci_flow.my_awesome_main_flow.id
  }
}