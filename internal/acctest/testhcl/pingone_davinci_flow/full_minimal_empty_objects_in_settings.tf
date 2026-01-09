
		%[1]s

resource "pingone_davinci_connector_instance" "%[2]s-errors" {
  environment_id = data.pingone_environment.general_test.id
  connector = {
    id = "errorConnector"
  }
  name = "%[2]s-errors"
}

resource "pingone_davinci_flow" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  description = "empty objects in settings test"
  name = "simple"
  settings = {
    csp = "worker-src 'self' blob:; script-src 'self' https://cdn.jsdelivr.net https://code.jquery.com https://devsdk.singularkey.com http://cdnjs.cloudflare.com 'unsafe-inline' 'unsafe-eval';"
    intermediate_loading_screen_css = "{}"
    intermediate_loading_screen_html = "{}"
    flow_http_timeout_in_seconds = 180
    log_level = 3
    use_custom_css = false
    css_links = []
    use_custom_script = false
    js_links = []
    flow_timeout_in_seconds = 0
    require_authentication_to_initiate = false
  }
  color = "#FFC8C1"
  graph_data = {
    elements = {
      nodes = {
        "2pzouq7el7" = {
          data = {
            id = "2pzouq7el7"
            node_type = "CONNECTION"
            connection_id = pingone_davinci_connector_instance.%[2]s-errors.id
            connector_id = "errorConnector"
            label = "Error Message"
            status = "configured"
            capability_name = "customErrorMessage"
            type = "action"
            properties = jsonencode({
              "error_message": {
                "value": "[\n  {\n    \"children\": [\n      {\n        \"text\": \"This is an error\"\n      }\n    ]\n  }\n]"
              },
              "error_description": {
                "value": "[\n  {\n    \"children\": [\n      {\n        \"text\": \"This is an error, really\"\n      }\n    ]\n  }\n]"
              }
            })
          }
          position = {
            x = 400
            y = 400
          }
          group = "nodes"
          removed = false
          selected = false
          selectable = true
          locked = false
          grabbable = true
          pannable = false
          classes = ""
        }
      }
    }
    
    zooming_enabled = true
    user_zooming_enabled = true
    zoom = 1
    min_zoom = 1e-50
    max_zoom = 1e+50
    panning_enabled = true
    user_panning_enabled = true
    pan = {
      x = 0
      y = 0
    }
    box_selection_enabled = true
    renderer = jsonencode({
      "name": "null"
    })
  }
}
