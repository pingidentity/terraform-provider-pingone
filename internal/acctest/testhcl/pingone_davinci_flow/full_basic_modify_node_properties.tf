		%[1]s

resource "pingone_davinci_connector_instance" "%[2]s-http" {
  environment_id = data.pingone_environment.general_test.id
  connector = {
    id = "httpConnector"
  }
  name = "%[2]s-http"
}

resource "pingone_davinci_connector_instance" "%[2]s-functions" {
  environment_id = data.pingone_environment.general_test.id
  connector = {
    id = "functionsConnector"
  }
  name = "%[2]s-functions"
}

resource "pingone_davinci_connector_instance" "%[2]s-errors" {
  environment_id = data.pingone_environment.general_test.id
  connector = {
    id = "errorConnector"
  }
  name = "%[2]s-errors"
}

resource "pingone_davinci_connector_instance" "%[2]s-flow" {
  environment_id = data.pingone_environment.general_test.id
  connector = {
    id = "flowConnector"
  }
  name = "%[2]s-flow"
}

resource "pingone_davinci_connector_instance" "%[2]s-variables" {
  environment_id = data.pingone_environment.general_test.id
  connector = {
    id = "variablesConnector"
  }
  name = "%[2]s-variables"
}

resource "pingone_davinci_flow" "%[2]s-subflow1" {
  environment_id = data.pingone_environment.general_test.id
  name = "subflow 1"
  description = "subflow 1 desc"
  color = "#AFD5FF"
  
  graph_data = {
    elements = {
      nodes = {
        "9awrr4q360" = {
          data = {
            id = "9awrr4q360"
            node_type = "CONNECTION"
            connection_id = pingone_davinci_connector_instance.%[2]s-http.id
            connector_id = "httpConnector"
            name = "Http"
            label = "Http"
            status = "configured"
            capability_name = "customHtmlMessage"
            type = "trigger"
            properties = jsonencode({
              "message": {
                "value": "[\n  {\n    \"children\": [\n      {\n        \"text\": \"Subflow 1\"\n      }\n    ]\n  }\n]"
              }
            })
          }
          position = {
            x = 277
            y = 236
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
        "rbi38g672i" = {
          data = {
            id = "rbi38g672i"
            node_type = "EVAL"
            label = "Evaluator"
          }
          position = {
            x = 394
            y = 237.25
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
        "exljnczoqz" = {
          data = {
            id = "exljnczoqz"
            node_type = "CONNECTION"
            connection_id = pingone_davinci_connector_instance.%[2]s-http.id
            connector_id = "httpConnector"
            name = "Http"
            label = "HTTP"
            status = "configured"
            capability_name = "createSuccessResponse"
            type = "action"
          }
          position = {
            x = 511
            y = 238.5
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
      edges = {
        "jv7enynltp|9awrr4q360|rbi38g672i" = {

          data = {
            id = "jv7enynltp"
            source = "9awrr4q360"
            target = "rbi38g672i"
          }
          position = {
            x = 0
            y = 0
          }
          group = "edges"
          removed = false
          selected = false
          selectable = true
          locked = false
          grabbable = true
          pannable = true
          classes = ""
        
        }
        "bn6hy8ycra|rbi38g672i|exljnczoqz" = {

          data = {
            id = "bn6hy8ycra"
            source = "rbi38g672i"
            target = "exljnczoqz"
          }
          position = {
            x = 0
            y = 0
          }
          group = "edges"
          removed = false
          selected = false
          selectable = true
          locked = false
          grabbable = true
          pannable = true
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


resource "pingone_davinci_flow" "%[2]s-subflow2" {
  environment_id = data.pingone_environment.general_test.id
  name = "subflow 2"
  description = "Cloned on Wed Jan 31 2024 13:43:43 GMT+0000 (Coordinated Universal Time). \n"
  color = "#AFD5FF"
  
  graph_data = {
    elements = {
      nodes = {
        "9awrr4q360" = {
          data = {
            id = "9awrr4q360"
            node_type = "CONNECTION"
            connection_id = pingone_davinci_connector_instance.%[2]s-http.id
            connector_id = "httpConnector"
            name = "Http"
            label = "Http"
            status = "configured"
            capability_name = "customHtmlMessage"
            type = "trigger"
            properties = jsonencode({
              "message": {
                "value": "[\n  {\n    \"children\": [\n      {\n        \"text\": \"Subflow 2\"\n      }\n    ]\n  }\n]"
              }
            })
          }
          position = {
            x = 277
            y = 236
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
        "rbi38g672i" = {
          data = {
            id = "rbi38g672i"
            node_type = "EVAL"
            label = "Evaluator"
          }
          position = {
            x = 394
            y = 237.25
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
        "exljnczoqz" = {
          data = {
            id = "exljnczoqz"
            node_type = "CONNECTION"
            connection_id = pingone_davinci_connector_instance.%[2]s-http.id
            connector_id = "httpConnector"
            name = "Http"
            label = "HTTP"
            status = "configured"
            capability_name = "createSuccessResponse"
            type = "action"
          }
          position = {
            x = 511
            y = 238.5
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
      edges = {
        "jv7enynltp|9awrr4q360|rbi38g672i" = {

          data = {
            id = "jv7enynltp"
            source = "9awrr4q360"
            target = "rbi38g672i"
          }
          position = {
            x = 0
            y = 0
          }
          group = "edges"
          removed = false
          selected = false
          selectable = true
          locked = false
          grabbable = true
          pannable = true
          classes = ""
        
        }
        "bn6hy8ycra|rbi38g672i|exljnczoqz" = {

          data = {
            id = "bn6hy8ycra"
            source = "rbi38g672i"
            target = "exljnczoqz"
          }
          position = {
            x = 0
            y = 0
          }
          group = "edges"
          removed = false
          selected = false
          selectable = true
          locked = false
          grabbable = true
          pannable = true
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

resource "pingone_davinci_flow" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name = "full-basic"
  color = "#E3F0FF"
  
  graph_data = {
    elements = {
      nodes = {
        "1u2m5vzr49" = {
          data = {
            id = "1u2m5vzr49"
            node_type = "CONNECTION"
            connection_id = pingone_davinci_connector_instance.%[2]s-http.id
            connector_id = "httpConnector"
            name = "%[2]s-http"
            label = "Http"
            status = "configured"
            capability_name = "customHtmlMessage"
            type = "trigger"
            properties = jsonencode({
              "message": {
                "value": "[\n  {\n    \"children\": [\n      {\n        \"text\": \"Hello, world?\"\n      }\n    ]\n  }\n]"
              }
            })
          }
          position = {
            x = 277
            y = 336
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
        "8fvg7tfr8j" = {
          data = {
            id = "8fvg7tfr8j"
            node_type = "EVAL"
            label = "Evaluator"
          }
          position = {
            x = 426.5
            y = 337.25
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
        "nx0o1b2cmw" = {
          data = {
            id = "nx0o1b2cmw"
            node_type = "CONNECTION"
            connection_id = pingone_davinci_connector_instance.%[2]s-functions.id
            connector_id = "functionsConnector"
            name = "%[2]s-functions"
            label = "Functions"
            status = "configured"
            capability_name = "AEqualsB"
            type = "trigger"
            properties = jsonencode({
              "leftValueA": {
                "value": "[\n  {\n    \"children\": [\n      {\n        \"text\": \"1\"\n      }\n    ]\n  }\n]"
              },
              "rightValueB": {
                "value": "[\n  {\n    \"children\": [\n      {\n        \"text\": \"1\"\n      }\n    ]\n  }\n]"
              }
            })
          }
          position = {
            x = 576
            y = 338.5
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
        "cdcw8k7dnx" = {
          data = {
            id = "cdcw8k7dnx"
            node_type = "EVAL"
            label = "Evaluator"
            properties = jsonencode({
              "vsp1ewtr9m": {
                "value": "allTriggersFalse"
              },
              "xb74p6rkd8": {
                "value": "anyTriggersFalse"
              }
            })
          }
          position = {
            x = 717
            y = 326
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
        "ikt13crnhy" = {
          data = {
            id = "ikt13crnhy"
            node_type = "CONNECTION"
            connection_id = pingone_davinci_connector_instance.%[2]s-http.id
            connector_id = "httpConnector"
            name = "%[2]s-http"
            label = "Http"
            status = "configured"
            capability_name = "createSuccessResponse"
            type = "action"
          }
          position = {
            x = 1197
            y = 266
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
        "vsp1ewtr9m" = {
          data = {
            id = "vsp1ewtr9m"
            node_type = "CONNECTION"
            connection_id = pingone_davinci_connector_instance.%[2]s-errors.id
            connector_id = "errorConnector"
            label = "Error Message"
            status = "configured"
            capability_name = "customErrorMessage"
            type = "action"
            properties = jsonencode({
              "errorMessage": {
                "value": "[\n  {\n    \"children\": [\n      {\n        \"text\": \"Error\"\n      }\n    ]\n  }\n]"
              }
            })
          }
          position = {
            x = 1197
            y = 416
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
        "xb74p6rkd8" = {
          data = {
            id = "xb74p6rkd8"
            node_type = "CONNECTION"
            connection_id = pingone_davinci_connector_instance.%[2]s-flow.id
            connector_id = "flowConnector"
            name = "%[2]s-flow"
            label = "Flow Conductor"
            status = "configured"
            capability_name = "startUiSubFlow"
            type = "trigger"
            properties = jsonencode({
              "subFlowId": {
                "value": {
                  "label": "subflow 2",
                  "value": "${pingone_davinci_flow.%[2]s-subflow2.id}"
                }
              },
              "subFlowVersionId": {
                "value": -1
              }
            })
          }
          position = {
            x = 867
            y = 446
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
        "kq5ybvwvro" = {
          data = {
            id = "kq5ybvwvro"
            node_type = "CONNECTION"
            connection_id = pingone_davinci_connector_instance.%[2]s-flow.id
            connector_id = "flowConnector"
            name = "%[2]s-flow"
            label = "Flow Conductor"
            status = "configured"
            capability_name = "startUiSubFlow"
            type = "trigger"
            properties = jsonencode({
              "subFlowId": {
                "value": {
                  "label": "subflow 1",
                  "value": "${pingone_davinci_flow.%[2]s-subflow1.id}"
                }
              },
              "subFlowVersionId": {
                "value": -1
              }
            })
          }
          position = {
            x = 867
            y = 236
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
        "j74pmg6577" = {
          data = {
            id = "j74pmg6577"
            node_type = "EVAL"
          }
          position = {
            x = 1017
            y = 236
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
        "pensvkew7y" = {
          data = {
            id = "pensvkew7y"
            node_type = "EVAL"
          }
          position = {
            x = 1032
            y = 431
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
        "3zvjdgdljx" = {
          data = {
            id = "3zvjdgdljx"
            node_type = "CONNECTION"
            connection_id = pingone_davinci_connector_instance.%[2]s-variables.id
            connector_id = "variablesConnector"
            name = "%[2]s-variables"
            label = "Variables"
            status = "configured"
            capability_name = "saveFlowValue"
            type = "trigger"
            properties = jsonencode({
              "saveFlowVariables": {
                "value": [
                  {
                    "name": "test123",
                    "value": "[\n  {\n    \"children\": [\n      {\n        \"text\": \"test456\"\n      }\n    ]\n  }\n]",
                    "key": 0.379286774724122,
                    "label": "test123 (string - flow)",
                    "type": "string"
                  },
                  {
                    "name": "test123NEW",
                    "value": "[\n  {\n    \"children\": [\n      {\n        \"text\": \"25\"\n      }\n    ]\n  }\n]",
                    "key": 0.379286774724122,
                    "label": "test123NEW (number - flow)",
                    "type": "number"
                  }
                ]
              }
            })
          }
          position = {
            x = 270
            y = 180
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
        "bbemfztdyk" = {
          data = {
            id = "bbemfztdyk"
            node_type = "EVAL"
          }
          position = {
            x = 273.5
            y = 258
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
      edges = {
        "hseww5vtf0|1u2m5vzr49|8fvg7tfr8j" = {

          data = {
            id = "hseww5vtf0"
            source = "1u2m5vzr49"
            target = "8fvg7tfr8j"
          }
          position = {
            x = 0
            y = 0
          }
          group = "edges"
          removed = false
          selected = false
          selectable = true
          locked = false
          grabbable = true
          pannable = true
          classes = ""
        
        }
        "ljavni2nky|8fvg7tfr8j|nx0o1b2cmw" = {

          data = {
            id = "ljavni2nky"
            source = "8fvg7tfr8j"
            target = "nx0o1b2cmw"
          }
          position = {
            x = 0
            y = 0
          }
          group = "edges"
          removed = false
          selected = false
          selectable = true
          locked = false
          grabbable = true
          pannable = true
          classes = ""
        
        }
        "0o2fqy3mf3|nx0o1b2cmw|cdcw8k7dnx" = {

          data = {
            id = "0o2fqy3mf3"
            source = "nx0o1b2cmw"
            target = "cdcw8k7dnx"
          }
          position = {
            x = 0
            y = 0
          }
          group = "edges"
          removed = false
          selected = false
          selectable = true
          locked = false
          grabbable = true
          pannable = true
          classes = ""
        
        }
        "493yd0jbi6|cdcw8k7dnx|kq5ybvwvro" = {

          data = {
            id = "493yd0jbi6"
            source = "cdcw8k7dnx"
            target = "kq5ybvwvro"
          }
          position = {
            x = 0
            y = 0
          }
          group = "edges"
          removed = false
          selected = false
          selectable = true
          locked = false
          grabbable = true
          pannable = true
          classes = ""
        
        }
        "pn2kixnzms|j74pmg6577|ikt13crnhy" = {

          data = {
            id = "pn2kixnzms"
            source = "j74pmg6577"
            target = "ikt13crnhy"
          }
          position = {
            x = 0
            y = 0
          }
          group = "edges"
          removed = false
          selected = false
          selectable = true
          locked = false
          grabbable = true
          pannable = true
          classes = ""
        
        }
        "0sb4quzlgx|kq5ybvwvro|j74pmg6577" = {

          data = {
            id = "0sb4quzlgx"
            source = "kq5ybvwvro"
            target = "j74pmg6577"
          }
          position = {
            x = 0
            y = 0
          }
          group = "edges"
          removed = false
          selected = false
          selectable = true
          locked = false
          grabbable = true
          pannable = true
          classes = ""
        
        }
        "v5p4i55lt9|cdcw8k7dnx|xb74p6rkd8" = {

          data = {
            id = "v5p4i55lt9"
            source = "cdcw8k7dnx"
            target = "xb74p6rkd8"
          }
          position = {
            x = 0
            y = 0
          }
          group = "edges"
          removed = false
          selected = false
          selectable = true
          locked = false
          grabbable = true
          pannable = true
          classes = ""
        
        }
        "k0trrhjqt6|xb74p6rkd8|pensvkew7y" = {

          data = {
            id = "k0trrhjqt6"
            source = "xb74p6rkd8"
            target = "pensvkew7y"
          }
          position = {
            x = 0
            y = 0
          }
          group = "edges"
          removed = false
          selected = false
          selectable = true
          locked = false
          grabbable = true
          pannable = true
          classes = ""
        
        }
        "2g0chago4l|pensvkew7y|vsp1ewtr9m" = {

          data = {
            id = "2g0chago4l"
            source = "pensvkew7y"
            target = "vsp1ewtr9m"
          }
          position = {
            x = 0
            y = 0
          }
          group = "edges"
          removed = false
          selected = false
          selectable = true
          locked = false
          grabbable = true
          pannable = true
          classes = ""
        
        }
        "gs1fx4x303|3zvjdgdljx|bbemfztdyk" = {

          data = {
            id = "gs1fx4x303"
            source = "3zvjdgdljx"
            target = "bbemfztdyk"
          }
          position = {
            x = 0
            y = 0
          }
          group = "edges"
          removed = false
          selected = false
          selectable = true
          locked = false
          grabbable = true
          pannable = true
          classes = ""
        
        }
        "cum544luro|bbemfztdyk|1u2m5vzr49" = {

          data = {
            id = "cum544luro"
            source = "bbemfztdyk"
            target = "1u2m5vzr49"
          }
          position = {
            x = 0
            y = 0
          }
          group = "edges"
          removed = false
          selected = false
          selectable = true
          locked = false
          grabbable = true
          pannable = true
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
