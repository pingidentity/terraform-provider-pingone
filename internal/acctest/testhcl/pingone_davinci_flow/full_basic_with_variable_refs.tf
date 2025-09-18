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

resource "pingone_davinci_variable" "%[2]s-variable1" {
  environment_id = data.pingone_environment.general_test.id
  name = "testVariable"
  context = "company"
  display_name = "Test Variable"
  data_type = "string"
  mutable = true
}

resource "pingone_davinci_variable" "%[2]s-variable2" {
  environment_id = data.pingone_environment.general_test.id
  name = "testVariable2"
  context = "company"
  display_name = "Test Variable"
  data_type = "string"
  mutable = true
}

resource "pingone_davinci_variable" "%[2]s-flowInstanceVariable1" {
  environment_id = data.pingone_environment.general_test.id
  name = "flowInstanceVariable1"
  context = "flowInstance"
  data_type = "string"
  mutable = true
}

resource "pingone_davinci_variable" "%[2]s-flowVariable" {
  environment_id = data.pingone_environment.general_test.id
  name = "testuser"
  context = "flow"
  flow = {
    id = pingone_davinci_flow.%[2]s.id
  }
  data_type = "string"
  mutable = true
}

resource "pingone_davinci_flow" "%[2]s-subflow1" {
  environment_id = data.pingone_environment.general_test.id
  name = "subflow 1"
  description = "subflow 1 desc"
  color = "#AFD5FF"
  
  graph_data = {
    elements = {
      nodes = [
        {
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
        },
        {
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
        },
        {
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
      ]
      edges = [
        {
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
        },
        {
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


resource "pingone_davinci_flow" "%[2]s-subflow2" {
  environment_id = data.pingone_environment.general_test.id
  name = "subflow 2"
  description = "Cloned on Wed Jan 31 2024 13:43:43 GMT+0000 (Coordinated Universal Time). \n"
  color = "#AFD5FF"
  
  graph_data = {
    elements = {
      nodes = [
        {
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
        },
        {
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
        },
        {
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
      ]
      edges = [
        {
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
        },
        {
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

resource "pingone_davinci_flow" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name = "%[2]s"
  description = "Imported on Fri Aug 16 2024 09:44:27 GMT+0000 (Coordinated Universal Time)"
  color = "#E3F0FF"
  
  graph_data = {
    elements = {
      nodes = [
        {
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
                "value": "[\n  {\n    \"children\": [\n      {\n        \"text\": \"Hello, world?\"\n      },\n      {\n        \"text\": \"\"\n      },\n      {\n        \"type\": \"link\",\n        \"src\": \"variable.svg\",\n        \"url\": \"${pingone_davinci_variable.%[2]s-variable1.name}\",\n        \"data\": \"{{global.company.variables.${pingone_davinci_variable.%[2]s-variable1.name}}}\",\n        \"tooltip\": \"{{global.company.variables.${pingone_davinci_variable.%[2]s-variable1.name}}}\",\n        \"children\": [\n          {\n            \"text\": \"${pingone_davinci_variable.%[2]s-variable1.name}\"\n          }\n        ]\n      },\n      {\n        \"text\": \"\"\n      },\n      {\n        \"type\": \"link\",\n        \"src\": \"variable.svg\",\n        \"url\": \"${pingone_davinci_variable.%[2]s-variable2.name}\",\n        \"data\": \"{{global.company.variables.${pingone_davinci_variable.%[2]s-variable2.name}}}\",\n        \"tooltip\": \"{{global.company.variables.${pingone_davinci_variable.%[2]s-variable2.name}}}\",\n        \"children\": [\n          {\n            \"text\": \"${pingone_davinci_variable.%[2]s-variable2.name}\"\n          }\n        ]\n      }\n    ]\n  }\n]"
              }
            })
          }
          position = {
            x = 284
            y = 392
          }
          group = "nodes"
          removed = false
          selected = false
          selectable = true
          locked = false
          grabbable = true
          pannable = false
          classes = ""
        },
        {
          data = {
            id = "8fvg7tfr8j"
            node_type = "EVAL"
            label = "Evaluator"
            
          }
          position = {
            x = 433.5
            y = 393.25
          }
          group = "nodes"
          removed = false
          selected = false
          selectable = true
          locked = false
          grabbable = true
          pannable = false
          classes = ""
        },
        {
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
            x = 583
            y = 394.5
          }
          group = "nodes"
          removed = false
          selected = false
          selectable = true
          locked = false
          grabbable = true
          pannable = false
          classes = ""
        },
        {
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
            x = 724
            y = 382
          }
          group = "nodes"
          removed = false
          selected = false
          selectable = true
          locked = false
          grabbable = true
          pannable = false
          classes = ""
        },
        {
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
            x = 1204
            y = 322
          }
          group = "nodes"
          removed = false
          selected = false
          selectable = true
          locked = false
          grabbable = true
          pannable = false
          classes = ""
        },
        {
          data = {
            id = "vsp1ewtr9m"
            node_type = "CONNECTION"
            connection_id = pingone_davinci_connector_instance.%[2]s-errors.id
            connector_id = "errorConnector"
            name = "%[2]s-error"
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
            x = 1204
            y = 472
          }
          group = "nodes"
          removed = false
          selected = false
          selectable = true
          locked = false
          grabbable = true
          pannable = false
          classes = ""
        },
        {
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
            x = 874
            y = 502
          }
          group = "nodes"
          removed = false
          selected = false
          selectable = true
          locked = false
          grabbable = true
          pannable = false
          classes = ""
        },
        {
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
            x = 874
            y = 292
          }
          group = "nodes"
          removed = false
          selected = false
          selectable = true
          locked = false
          grabbable = true
          pannable = false
          classes = ""
        },
        {
          data = {
            id = "j74pmg6577"
            node_type = "EVAL"
            
          }
          position = {
            x = 1024
            y = 292
          }
          group = "nodes"
          removed = false
          selected = false
          selectable = true
          locked = false
          grabbable = true
          pannable = false
          classes = ""
        },
        {
          data = {
            id = "pensvkew7y"
            node_type = "EVAL"
            
          }
          position = {
            x = 1039
            y = 487
          }
          group = "nodes"
          removed = false
          selected = false
          selectable = true
          locked = false
          grabbable = true
          pannable = false
          classes = ""
        },
        {
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
                    "name": "fdgdfgfdg",
                    "key": 0.8936786494474329,
                    "label": "fdgdfgfdg (string - flow)",
                    "type": "string",
                    "value": "fdgdfgfdgValue"
                  },
                  {
                    "name": "fdgdfgfdgNEW",
                    "key": 0.8936786494474330,
                    "label": "fdgdfgfdgNEW (string - flow)",
                    "type": "string",
                    "value": "fdgdfgfdgNEWValue"
                  },
                  {
                    "name": "test123",
                    "key": 0.379286774724122,
                    "label": "test123 (number - flow)",
                    "type": "number",
                    "value": 5
                  }
                ]
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
        },
        {
          data = {
            id = "bbemfztdyk"
            node_type = "EVAL"
            
          }
          position = {
            x = 280.5
            y = 314
          }
          group = "nodes"
          removed = false
          selected = false
          selectable = true
          locked = false
          grabbable = true
          pannable = false
          classes = ""
        },
        {
          data = {
            id = "0cj7n971ix"
            node_type = "CONNECTION"
            connection_id = pingone_davinci_connector_instance.%[2]s-variables.id
            connector_id = "variablesConnector"
            name = "%[2]s-variables"
            label = "%[2]s-variables"
            status = "configured"
            capability_name = "saveValue"
            type = "trigger"
            properties = jsonencode({
              "saveVariables": {
                "value": [
                  {
                    "name": "${pingone_davinci_variable.%[2]s-variable1.name}",
                    "key": 0.09068454768967449,
                    "label": "flowInstanceVariable1 (string - flowInstance)",
                    "type": "${pingone_davinci_variable.%[2]s-variable1.data_type}",
                    "value": "flowInstanceVariable1Value"
                  }
                ]
              }
            })
          }
          position = {
            x = 270
            y = 120
          }
          group = "nodes"
          removed = false
          selected = false
          selectable = true
          locked = false
          grabbable = true
          pannable = false
          classes = ""
        },
        {
          data = {
            id = "esg7oyahen"
            node_type = "EVAL"
            
          }
          position = {
            x = 273.5
            y = 178
          }
          group = "nodes"
          removed = false
          selected = false
          selectable = true
          locked = false
          grabbable = true
          pannable = false
          classes = ""
        },
        {
          data = {
            id = "j3j8fmgc9q"
            node_type = "CONNECTION"
            connection_id = pingone_davinci_connector_instance.%[2]s-variables.id
            connector_id = "variablesConnector"
            name = "%[2]s-variables"
            label = "%[2]s-variables"
            status = "configured"
            capability_name = "saveValueUserInfo"
            type = "trigger"
            properties = jsonencode({
              "saveVariables": {
                "value": [
                  {
                    "name": "testuser",
                    "key": 0.9814043007447408,
                    "label": "testuser (string - flow)",
                    "type": "string",
                    "value": "testuserValue"
                  }
                ]
              }
            })
          }
          position = {
            x = 90
            y = 60
          }
          group = "nodes"
          removed = false
          selected = false
          selectable = true
          locked = false
          grabbable = true
          pannable = false
          classes = ""
        },
        {
          data = {
            id = "1uu35lv024"
            node_type = "EVAL"
            
          }
          position = {
            x = 180
            y = 90
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
      edges = [
        {
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
        },
        {
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
        },
        {
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
        },
        {
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
        },
        {
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
        },
        {
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
        },
        {
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
        },
        {
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
        },
        {
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
        },
        {
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
        },
        {
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
        },
        {
          data = {
            id = "mz51tp7j0f"
            source = "0cj7n971ix"
            target = "esg7oyahen"
            
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
        },
        {
          data = {
            id = "as3c6w9yus"
            source = "esg7oyahen"
            target = "3zvjdgdljx"
            
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
        },
        {
          data = {
            id = "hdak1wwkml"
            source = "j3j8fmgc9q"
            target = "1uu35lv024"
            
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
        },
        {
          data = {
            id = "dv5jn5u6e7"
            source = "1uu35lv024"
            target = "0cj7n971ix"
            
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
