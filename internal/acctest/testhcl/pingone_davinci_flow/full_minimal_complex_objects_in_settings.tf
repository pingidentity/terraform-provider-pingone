
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
    intermediate_loading_screen_css = ""
    intermediate_loading_screen_html = ""
    flow_http_timeout_in_seconds = 180
    log_level = 3
    use_custom_css = true
    css_links = [
        "https://assets.pingone.com/ux/end-user-nano/0.1.0-alpha.9/end-user-nano.css",
        "https://assets.pingone.com/ux/astro-nano/0.1.0-alpha.11/icons.css",
        "https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.0.0-beta3/css/all.min.css"
      ]
    use_custom_script = true
    js_links = [
        {
          crossorigin = "test",
          defer = false,
          integrity = "test",
          label = "https://cdnjs.cloudflare.com/ajax/libs/html2pdf.js/0.9.2/html2pdf.bundle.min.js",
          referrerpolicy = "test",
          type = "test",
          value = "https://cdnjs.cloudflare.com/ajax/libs/html2pdf.js/0.9.2/html2pdf.bundle.min.js"
        },
        {
          crossorigin = "anonymous",
          defer = false,
          integrity = "test",
          label = "https://ajax.googleapis.com/ajax/libs/jquery/2.2.4/jquery.min.js",
          referrerpolicy = "test",
          type = "test",
          value = "https://ajax.googleapis.com/ajax/libs/jquery/2.2.4/jquery.min.js"
        }
      ]
    flow_timeout_in_seconds = 0
    require_authentication_to_initiate = false
  }
  color = "#FFC8C1"
  graph_data = {
    elements = {
      nodes = [
        {
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
      ]
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
