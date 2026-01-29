resource "pingone_davinci_connector_instance" "example-errors" {
  environment_id = var.environment_id
  connector = {
    id = "errorConnector"
  }
  name = "example-errors"
}

resource "pingone_davinci_flow" "example" {
  environment_id = var.environment_id
  description    = "A simple flow that shows an error message"
  name           = "simple"
  settings = {
    csp                              = "worker-src 'self' blob:; script-src 'self' https://cdn.jsdelivr.net https://code.jquery.com https://devsdk.singularkey.com http://cdnjs.cloudflare.com 'unsafe-inline' 'unsafe-eval';"
    intermediate_loading_screen_css  = ""
    intermediate_loading_screen_html = ""
    flow_http_timeout_in_seconds     = 300
    log_level                        = 1
    use_custom_css                   = true
  }
  color = "#FFC8C1"
  graph_data = {
    elements = {
      nodes = {
        "2pzouq7el7" = {
          data = {
            node_type       = "CONNECTION"
            connection_id   = pingone_davinci_connector_instance.example-errors.id
            connector_id    = "errorConnector"
            label           = "Error Message"
            status          = "configured"
            capability_name = "customErrorMessage"
            type            = "action"
            properties = jsonencode({
              "error_message" : {
                "value" : "[\n  {\n    \"children\": [\n      {\n        \"text\": \"This is an error\"\n      }\n    ]\n  }\n]"
              },
              "error_description" : {
                "value" : "[\n  {\n    \"children\": [\n      {\n        \"text\": \"This is an error, really\"\n      }\n    ]\n  }\n]"
              }
            })
          }
          position = {
            x = 400
            y = 400
          }
          group      = "nodes"
          removed    = false
          selected   = false
          selectable = true
          locked     = false
          grabbable  = true
          pannable   = false
          classes    = ""
        }
      }
    }

    zooming_enabled      = true
    user_zooming_enabled = true
    zoom                 = 1
    min_zoom             = 1e-50
    max_zoom             = 1e+50
    panning_enabled      = true
    user_panning_enabled = true
    pan = {
      x = 0
      y = 0
    }
    box_selection_enabled = true
    renderer = jsonencode({
      "name" : "null"
    })
  }
}

resource "pingone_davinci_flow_enabled" "example" {
  environment_id = var.environment_id
  flow_id        = pingone_davinci_flow.example.id
  enabled        = true
}