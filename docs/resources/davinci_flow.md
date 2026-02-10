---
page_title: "pingone_davinci_flow Resource - terraform-provider-pingone"
subcategory: "DaVinci"
description: |-
  Resource to create and manage a davinci flow.
---

# pingone_davinci_flow (Resource)

Resource to create and manage a davinci flow.

## Example Usage

```terraform
resource "pingone_davinci_connector_instance" "example-http" {
  environment_id = var.environment_id
  connector = {
    id = "httpConnector"
  }
  name = "example-http"
}

resource "pingone_davinci_connector_instance" "example-functions" {
  environment_id = var.environment_id
  connector = {
    id = "functionsConnector"
  }
  name = "example-functions"
}

resource "pingone_davinci_connector_instance" "example-errors" {
  environment_id = var.environment_id
  connector = {
    id = "errorConnector"
  }
  name = "example-errors"
}

resource "pingone_davinci_connector_instance" "example-flow" {
  environment_id = var.environment_id
  connector = {
    id = "flowConnector"
  }
  name = "example-flow"
}

resource "pingone_davinci_connector_instance" "example-variables" {
  environment_id = var.environment_id
  connector = {
    id = "variablesConnector"
  }
  name = "example-variables"
}

resource "pingone_davinci_variable" "example-variable1" {
  environment_id = var.environment_id
  name           = "testVariable"
  context        = "company"
  display_name   = "Test Variable"
  data_type      = "string"
  mutable        = true
}

resource "pingone_davinci_variable" "example-variable2" {
  environment_id = var.environment_id
  name           = "testVariable2"
  context        = "company"
  display_name   = "Test Variable"
  data_type      = "string"
  mutable        = true
}

resource "pingone_davinci_variable" "example-flowInstanceVariable1" {
  environment_id = var.environment_id
  name           = "flowInstanceVariable1"
  context        = "flowInstance"
  data_type      = "string"
  mutable        = true
}

resource "pingone_davinci_variable" "example-flowVariable" {
  environment_id = var.environment_id
  name           = "testuser"
  context        = "flow"
  flow = {
    id = pingone_davinci_flow.example.id
  }
  data_type = "string"
  mutable   = true
}

resource "pingone_davinci_flow" "example-subflow1" {
  environment_id = var.environment_id
  name           = "subflow 1"
  description    = "subflow 1 desc"
  color          = "#AFD5FF"

  graph_data = {
    elements = {
      nodes = {
        "9awrr4q360" = {

          data = {
            node_type       = "CONNECTION"
            connection_id   = pingone_davinci_connector_instance.example-http.id
            connector_id    = "httpConnector"
            name            = "Http"
            label           = "Http"
            status          = "configured"
            capability_name = "customHtmlMessage"
            type            = "trigger"
            properties = jsonencode({
              "message" : {
                "value" : "[\n  {\n    \"children\": [\n      {\n        \"text\": \"Subflow 1\"\n      }\n    ]\n  }\n]"
              }
            })
          }
          position = {
            x = 277
            y = 236
          }
          group      = "nodes"
          removed    = false
          selected   = false
          selectable = true
          locked     = false
          grabbable  = true
          pannable   = false


        }
        "rbi38g672i" = {

          data = {
            node_type = "EVAL"
            label     = "Evaluator"

          }
          position = {
            x = 394
            y = 237.25
          }
          group      = "nodes"
          removed    = false
          selected   = false
          selectable = true
          locked     = false
          grabbable  = true
          pannable   = false


        }
        "exljnczoqz" = {

          data = {
            node_type       = "CONNECTION"
            connection_id   = pingone_davinci_connector_instance.example-http.id
            connector_id    = "httpConnector"
            name            = "Http"
            label           = "HTTP"
            status          = "configured"
            capability_name = "createSuccessResponse"
            type            = "action"

          }
          position = {
            x = 511
            y = 238.5
          }
          group      = "nodes"
          removed    = false
          selected   = false
          selectable = true
          locked     = false
          grabbable  = true
          pannable   = false


        }

      }
      edges = {
        "jv7enynltp" = {

          data = {
            source = "9awrr4q360"
            target = "rbi38g672i"
          }
          position = {
            x = 0
            y = 0
          }
          group      = "edges"
          removed    = false
          selected   = false
          selectable = true
          locked     = false
          grabbable  = true
          pannable   = true


        }
        "bn6hy8ycra" = {

          data = {
            source = "rbi38g672i"
            target = "exljnczoqz"
          }
          position = {
            x = 0
            y = 0
          }
          group      = "edges"
          removed    = false
          selected   = false
          selectable = true
          locked     = false
          grabbable  = true
          pannable   = true


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


resource "pingone_davinci_flow" "example-subflow2" {
  environment_id = var.environment_id
  name           = "subflow 2"
  description    = "subflow 2 desc"
  color          = "#AFD5FF"

  graph_data = {
    elements = {
      nodes = {
        "9awrr4q360" = {

          data = {
            node_type       = "CONNECTION"
            connection_id   = pingone_davinci_connector_instance.example-http.id
            connector_id    = "httpConnector"
            name            = "Http"
            label           = "Http"
            status          = "configured"
            capability_name = "customHtmlMessage"
            type            = "trigger"
            properties = jsonencode({
              "message" : {
                "value" : "[\n  {\n    \"children\": [\n      {\n        \"text\": \"Subflow 2\"\n      }\n    ]\n  }\n]"
              }
            })
          }
          position = {
            x = 277
            y = 236
          }
          group      = "nodes"
          removed    = false
          selected   = false
          selectable = true
          locked     = false
          grabbable  = true
          pannable   = false


        }
        "rbi38g672i" = {

          data = {
            node_type = "EVAL"
            label     = "Evaluator"

          }
          position = {
            x = 394
            y = 237.25
          }
          group      = "nodes"
          removed    = false
          selected   = false
          selectable = true
          locked     = false
          grabbable  = true
          pannable   = false


        }
        "exljnczoqz" = {

          data = {
            node_type       = "CONNECTION"
            connection_id   = pingone_davinci_connector_instance.example-http.id
            connector_id    = "httpConnector"
            name            = "Http"
            label           = "HTTP"
            status          = "configured"
            capability_name = "createSuccessResponse"
            type            = "action"

          }
          position = {
            x = 511
            y = 238.5
          }
          group      = "nodes"
          removed    = false
          selected   = false
          selectable = true
          locked     = false
          grabbable  = true
          pannable   = false


        }

      }
      edges = {
        "jv7enynltp" = {

          data = {
            source = "9awrr4q360"
            target = "rbi38g672i"
          }
          position = {
            x = 0
            y = 0
          }
          group      = "edges"
          removed    = false
          selected   = false
          selectable = true
          locked     = false
          grabbable  = true
          pannable   = true


        }
        "bn6hy8ycra" = {

          data = {
            source = "rbi38g672i"
            target = "exljnczoqz"
          }
          position = {
            x = 0
            y = 0
          }
          group      = "edges"
          removed    = false
          selected   = false
          selectable = true
          locked     = false
          grabbable  = true
          pannable   = true


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

resource "pingone_davinci_flow" "example" {
  environment_id = var.environment_id
  name           = "example"
  description    = "Main flow description"
  color          = "#E3F0FF"

  graph_data = {
    elements = {
      nodes = {
        "1u2m5vzr49" = {

          data = {
            node_type       = "CONNECTION"
            connection_id   = pingone_davinci_connector_instance.example-http.id
            connector_id    = "httpConnector"
            name            = "example-http"
            label           = "Http"
            status          = "configured"
            capability_name = "customHtmlMessage"
            type            = "trigger"
            properties = jsonencode({
              "message" : {
                "value" : "[\n  {\n    \"children\": [\n      {\n        \"text\": \"Hello, world?\"\n      },\n      {\n        \"text\": \"\"\n      },\n      {\n        \"type\": \"link\",\n        \"src\": \"variable.svg\",\n        \"url\": \"${pingone_davinci_variable.example-variable1.name}\",\n        \"data\": \"{{global.company.variables.${pingone_davinci_variable.example-variable1.name}}}\",\n        \"tooltip\": \"{{global.company.variables.${pingone_davinci_variable.example-variable1.name}}}\",\n        \"children\": [\n          {\n            \"text\": \"${pingone_davinci_variable.example-variable1.name}\"\n          }\n        ]\n      },\n      {\n        \"text\": \"\"\n      },\n      {\n        \"type\": \"link\",\n        \"src\": \"variable.svg\",\n        \"url\": \"${pingone_davinci_variable.example-variable2.name}\",\n        \"data\": \"{{global.company.variables.${pingone_davinci_variable.example-variable2.name}}}\",\n        \"tooltip\": \"{{global.company.variables.${pingone_davinci_variable.example-variable2.name}}}\",\n        \"children\": [\n          {\n            \"text\": \"${pingone_davinci_variable.example-variable2.name}\"\n          }\n        ]\n      }\n    ]\n  }\n]"
              }
            })
          }
          position = {
            x = 284
            y = 392
          }
          group      = "nodes"
          removed    = false
          selected   = false
          selectable = true
          locked     = false
          grabbable  = true
          pannable   = false


        }
        "8fvg7tfr8j" = {

          data = {
            node_type = "EVAL"
            label     = "Evaluator"

          }
          position = {
            x = 433.5
            y = 393.25
          }
          group      = "nodes"
          removed    = false
          selected   = false
          selectable = true
          locked     = false
          grabbable  = true
          pannable   = false


        }
        "nx0o1b2cmw" = {

          data = {
            node_type       = "CONNECTION"
            connection_id   = pingone_davinci_connector_instance.example-functions.id
            connector_id    = "functionsConnector"
            name            = "example-functions"
            label           = "Functions"
            status          = "configured"
            capability_name = "AEqualsB"
            type            = "trigger"
            properties = jsonencode({
              "leftValueA" : {
                "value" : "[\n  {\n    \"children\": [\n      {\n        \"text\": \"1\"\n      }\n    ]\n  }\n]"
              },
              "rightValueB" : {
                "value" : "[\n  {\n    \"children\": [\n      {\n        \"text\": \"1\"\n      }\n    ]\n  }\n]"
              }
            })
          }
          position = {
            x = 583
            y = 394.5
          }
          group      = "nodes"
          removed    = false
          selected   = false
          selectable = true
          locked     = false
          grabbable  = true
          pannable   = false


        }
        "cdcw8k7dnx" = {

          data = {
            node_type = "EVAL"
            label     = "Evaluator"
            properties = jsonencode({
              "vsp1ewtr9m" : {
                "value" : "allTriggersFalse"
              },
              "xb74p6rkd8" : {
                "value" : "anyTriggersFalse"
              }
            })
          }
          position = {
            x = 724
            y = 382
          }
          group      = "nodes"
          removed    = false
          selected   = false
          selectable = true
          locked     = false
          grabbable  = true
          pannable   = false


        }
        "ikt13crnhy" = {

          data = {
            node_type       = "CONNECTION"
            connection_id   = pingone_davinci_connector_instance.example-http.id
            connector_id    = "httpConnector"
            name            = "example-http"
            label           = "Http"
            status          = "configured"
            capability_name = "createSuccessResponse"
            type            = "action"

          }
          position = {
            x = 1204
            y = 322
          }
          group      = "nodes"
          removed    = false
          selected   = false
          selectable = true
          locked     = false
          grabbable  = true
          pannable   = false


        }
        "vsp1ewtr9m" = {

          data = {
            node_type       = "CONNECTION"
            connection_id   = pingone_davinci_connector_instance.example-errors.id
            connector_id    = "errorConnector"
            name            = "example-error"
            label           = "Error Message"
            status          = "configured"
            capability_name = "customErrorMessage"
            type            = "action"
            properties = jsonencode({
              "errorMessage" : {
                "value" : "[\n  {\n    \"children\": [\n      {\n        \"text\": \"Error\"\n      }\n    ]\n  }\n]"
              }
            })
          }
          position = {
            x = 1204
            y = 472
          }
          group      = "nodes"
          removed    = false
          selected   = false
          selectable = true
          locked     = false
          grabbable  = true
          pannable   = false


        }
        "xb74p6rkd8" = {

          data = {
            node_type       = "CONNECTION"
            connection_id   = pingone_davinci_connector_instance.example-flow.id
            connector_id    = "flowConnector"
            name            = "example-flow"
            label           = "Flow Conductor"
            status          = "configured"
            capability_name = "startUiSubFlow"
            type            = "trigger"
            properties = jsonencode({
              "subFlowId" : {
                "value" : {
                  "label" : "subflow 2",
                  "value" : "${pingone_davinci_flow.example-subflow2.id}"
                }
              },
              "subFlowVersionId" : {
                "value" : -1
              }
            })
          }
          position = {
            x = 874
            y = 502
          }
          group      = "nodes"
          removed    = false
          selected   = false
          selectable = true
          locked     = false
          grabbable  = true
          pannable   = false


        }
        "kq5ybvwvro" = {

          data = {
            node_type       = "CONNECTION"
            connection_id   = pingone_davinci_connector_instance.example-flow.id
            connector_id    = "flowConnector"
            name            = "example-flow"
            label           = "Flow Conductor"
            status          = "configured"
            capability_name = "startUiSubFlow"
            type            = "trigger"
            properties = jsonencode({
              "subFlowId" : {
                "value" : {
                  "label" : "subflow 1",
                  "value" : "${pingone_davinci_flow.example-subflow1.id}"
                }
              },
              "subFlowVersionId" : {
                "value" : -1
              }
            })
          }
          position = {
            x = 874
            y = 292
          }
          group      = "nodes"
          removed    = false
          selected   = false
          selectable = true
          locked     = false
          grabbable  = true
          pannable   = false


        }
        "j74pmg6577" = {

          data = {
            node_type = "EVAL"

          }
          position = {
            x = 1024
            y = 292
          }
          group      = "nodes"
          removed    = false
          selected   = false
          selectable = true
          locked     = false
          grabbable  = true
          pannable   = false


        }
        "pensvkew7y" = {

          data = {
            node_type = "EVAL"

          }
          position = {
            x = 1039
            y = 487
          }
          group      = "nodes"
          removed    = false
          selected   = false
          selectable = true
          locked     = false
          grabbable  = true
          pannable   = false


        }
        "3zvjdgdljx" = {

          data = {
            node_type       = "CONNECTION"
            connection_id   = pingone_davinci_connector_instance.example-variables.id
            connector_id    = "variablesConnector"
            name            = "example-variables"
            label           = "Variables"
            status          = "configured"
            capability_name = "saveFlowValue"
            type            = "trigger"
            properties = jsonencode({
              "saveFlowVariables" : {
                "value" : [
                  {
                    "name" : "fdgdfgfdg",
                    "key" : 0.8936786494474329,
                    "label" : "fdgdfgfdg (string - flow)",
                    "type" : "string",
                    "value" : "fdgdfgfdgValue"
                  },
                  {
                    "name" : "fdgdfgfdgNEW",
                    "key" : 0.8936786494474330,
                    "label" : "fdgdfgfdgNEW (string - flow)",
                    "type" : "string",
                    "value" : "fdgdfgfdgNEWValue"
                  },
                  {
                    "name" : "test123",
                    "key" : 0.379286774724122,
                    "label" : "test123 (number - flow)",
                    "type" : "number",
                    "value" : 5
                  }
                ]
              }
            })
          }
          position = {
            x = 277
            y = 236
          }
          group      = "nodes"
          removed    = false
          selected   = false
          selectable = true
          locked     = false
          grabbable  = true
          pannable   = false


        }
        "bbemfztdyk" = {

          data = {
            node_type = "EVAL"

          }
          position = {
            x = 280.5
            y = 314
          }
          group      = "nodes"
          removed    = false
          selected   = false
          selectable = true
          locked     = false
          grabbable  = true
          pannable   = false


        }
        "0cj7n971ix" = {

          data = {
            node_type       = "CONNECTION"
            connection_id   = pingone_davinci_connector_instance.example-variables.id
            connector_id    = "variablesConnector"
            name            = "example-variables"
            label           = "example-variables"
            status          = "configured"
            capability_name = "saveValue"
            type            = "trigger"
            properties = jsonencode({
              "saveVariables" : {
                "value" : [
                  {
                    "name" : "${pingone_davinci_variable.example-variable1.name}",
                    "key" : 0.09068454768967449,
                    "label" : "flowInstanceVariable1 (string - flowInstance)",
                    "type" : "${pingone_davinci_variable.example-variable1.data_type}",
                    "value" : "flowInstanceVariable1Value"
                  }
                ]
              }
            })
          }
          position = {
            x = 270
            y = 120
          }
          group      = "nodes"
          removed    = false
          selected   = false
          selectable = true
          locked     = false
          grabbable  = true
          pannable   = false


        }
        "esg7oyahen" = {

          data = {
            node_type = "EVAL"

          }
          position = {
            x = 273.5
            y = 178
          }
          group      = "nodes"
          removed    = false
          selected   = false
          selectable = true
          locked     = false
          grabbable  = true
          pannable   = false


        }
        "j3j8fmgc9q" = {

          data = {
            node_type       = "CONNECTION"
            connection_id   = pingone_davinci_connector_instance.example-variables.id
            connector_id    = "variablesConnector"
            name            = "example-variables"
            label           = "example-variables"
            status          = "configured"
            capability_name = "saveValueUserInfo"
            type            = "trigger"
            properties = jsonencode({
              "saveVariables" : {
                "value" : [
                  {
                    "name" : "testuser",
                    "key" : 0.9814043007447408,
                    "label" : "testuser (string - flow)",
                    "type" : "string",
                    "value" : "testuserValue"
                  }
                ]
              }
            })
          }
          position = {
            x = 90
            y = 60
          }
          group      = "nodes"
          removed    = false
          selected   = false
          selectable = true
          locked     = false
          grabbable  = true
          pannable   = false


        }
        "1uu35lv024" = {

          data = {
            node_type = "EVAL"

          }
          position = {
            x = 180
            y = 90
          }
          group      = "nodes"
          removed    = false
          selected   = false
          selectable = true
          locked     = false
          grabbable  = true
          pannable   = false


        }

      }
      edges = {
        "hseww5vtf0" = {

          data = {
            source = "1u2m5vzr49"
            target = "8fvg7tfr8j"

          }
          position = {
            x = 0
            y = 0
          }
          group      = "edges"
          removed    = false
          selected   = false
          selectable = true
          locked     = false
          grabbable  = true
          pannable   = true


        }
        "ljavni2nky" = {

          data = {
            source = "8fvg7tfr8j"
            target = "nx0o1b2cmw"

          }
          position = {
            x = 0
            y = 0
          }
          group      = "edges"
          removed    = false
          selected   = false
          selectable = true
          locked     = false
          grabbable  = true
          pannable   = true


        }
        "0o2fqy3mf3" = {

          data = {
            source = "nx0o1b2cmw"
            target = "cdcw8k7dnx"

          }
          position = {
            x = 0
            y = 0
          }
          group      = "edges"
          removed    = false
          selected   = false
          selectable = true
          locked     = false
          grabbable  = true
          pannable   = true


        }
        "493yd0jbi6" = {

          data = {
            source = "cdcw8k7dnx"
            target = "kq5ybvwvro"

          }
          position = {
            x = 0
            y = 0
          }
          group      = "edges"
          removed    = false
          selected   = false
          selectable = true
          locked     = false
          grabbable  = true
          pannable   = true


        }
        "pn2kixnzms" = {

          data = {
            source = "j74pmg6577"
            target = "ikt13crnhy"

          }
          position = {
            x = 0
            y = 0
          }
          group      = "edges"
          removed    = false
          selected   = false
          selectable = true
          locked     = false
          grabbable  = true
          pannable   = true


        }
        "0sb4quzlgx" = {

          data = {
            source = "kq5ybvwvro"
            target = "j74pmg6577"

          }
          position = {
            x = 0
            y = 0
          }
          group      = "edges"
          removed    = false
          selected   = false
          selectable = true
          locked     = false
          grabbable  = true
          pannable   = true


        }
        "v5p4i55lt9" = {

          data = {
            source = "cdcw8k7dnx"
            target = "xb74p6rkd8"

          }
          position = {
            x = 0
            y = 0
          }
          group      = "edges"
          removed    = false
          selected   = false
          selectable = true
          locked     = false
          grabbable  = true
          pannable   = true


        }
        "k0trrhjqt6" = {

          data = {
            source = "xb74p6rkd8"
            target = "pensvkew7y"

          }
          position = {
            x = 0
            y = 0
          }
          group      = "edges"
          removed    = false
          selected   = false
          selectable = true
          locked     = false
          grabbable  = true
          pannable   = true


        }
        "2g0chago4l" = {

          data = {
            source = "pensvkew7y"
            target = "vsp1ewtr9m"

          }
          position = {
            x = 0
            y = 0
          }
          group      = "edges"
          removed    = false
          selected   = false
          selectable = true
          locked     = false
          grabbable  = true
          pannable   = true


        }
        "gs1fx4x303" = {

          data = {
            source = "3zvjdgdljx"
            target = "bbemfztdyk"

          }
          position = {
            x = 0
            y = 0
          }
          group      = "edges"
          removed    = false
          selected   = false
          selectable = true
          locked     = false
          grabbable  = true
          pannable   = true


        }
        "cum544luro" = {

          data = {
            source = "bbemfztdyk"
            target = "1u2m5vzr49"

          }
          position = {
            x = 0
            y = 0
          }
          group      = "edges"
          removed    = false
          selected   = false
          selectable = true
          locked     = false
          grabbable  = true
          pannable   = true


        }
        "mz51tp7j0f" = {

          data = {
            source = "0cj7n971ix"
            target = "esg7oyahen"

          }
          position = {
            x = 0
            y = 0
          }
          group      = "edges"
          removed    = false
          selected   = false
          selectable = true
          locked     = false
          grabbable  = true
          pannable   = true


        }
        "as3c6w9yus" = {

          data = {
            source = "esg7oyahen"
            target = "3zvjdgdljx"

          }
          position = {
            x = 0
            y = 0
          }
          group      = "edges"
          removed    = false
          selected   = false
          selectable = true
          locked     = false
          grabbable  = true
          pannable   = true


        }
        "hdak1wwkml" = {

          data = {
            source = "j3j8fmgc9q"
            target = "1uu35lv024"

          }
          position = {
            x = 0
            y = 0
          }
          group      = "edges"
          removed    = false
          selected   = false
          selectable = true
          locked     = false
          grabbable  = true
          pannable   = true


        }
        "dv5jn5u6e7" = {

          data = {
            source = "1uu35lv024"
            target = "0cj7n971ix"

          }
          position = {
            x = 0
            y = 0
          }
          group      = "edges"
          removed    = false
          selected   = false
          selectable = true
          locked     = false
          grabbable  = true
          pannable   = true


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
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `environment_id` (String) The ID of the environment to create and manage the flow in. Must be a valid PingOne resource ID. This field is immutable and will trigger a replace plan if changed.
- `name` (String)

### Optional

- `color` (String)
- `description` (String)
- `graph_data` (Attributes) (see [below for nested schema](#nestedatt--graph_data))
- `input_schema` (Attributes List) (see [below for nested schema](#nestedatt--input_schema))
- `output_schema` (Attributes) (see [below for nested schema](#nestedatt--output_schema))
- `settings` (Attributes) (see [below for nested schema](#nestedatt--settings))
- `trigger` (Attributes) (see [below for nested schema](#nestedatt--trigger))

### Read-Only

- `connectors` (Attributes Set) (see [below for nested schema](#nestedatt--connectors))
- `current_version` (Number)
- `enabled` (Boolean)
- `id` (String) The ID of this resource.
- `published_version` (Number)

<a id="nestedatt--graph_data"></a>
### Nested Schema for `graph_data`

Optional:

- `box_selection_enabled` (Boolean)
- `data` (String)
- `elements` (Attributes) (see [below for nested schema](#nestedatt--graph_data--elements))
- `max_zoom` (Number)
- `min_zoom` (Number)
- `pan` (Attributes) (see [below for nested schema](#nestedatt--graph_data--pan))
- `panning_enabled` (Boolean)
- `renderer` (String)
- `user_panning_enabled` (Boolean)
- `user_zooming_enabled` (Boolean)
- `zoom` (Number)
- `zooming_enabled` (Boolean)

<a id="nestedatt--graph_data--elements"></a>
### Nested Schema for `graph_data.elements`

Optional:

- `edges` (Attributes Map) (see [below for nested schema](#nestedatt--graph_data--elements--edges))
- `nodes` (Attributes Map) (see [below for nested schema](#nestedatt--graph_data--elements--nodes))

<a id="nestedatt--graph_data--elements--edges"></a>
### Nested Schema for `graph_data.elements.edges`

Required:

- `data` (Attributes) (see [below for nested schema](#nestedatt--graph_data--elements--edges--data))

Optional:

- `classes` (String)
- `grabbable` (Boolean)
- `group` (String)
- `locked` (Boolean)
- `pannable` (Boolean)
- `position` (Attributes) (see [below for nested schema](#nestedatt--graph_data--elements--edges--position))
- `removed` (Boolean)
- `selectable` (Boolean)
- `selected` (Boolean)

<a id="nestedatt--graph_data--elements--edges--data"></a>
### Nested Schema for `graph_data.elements.edges.data`

Required:

- `source` (String)
- `target` (String)


<a id="nestedatt--graph_data--elements--edges--position"></a>
### Nested Schema for `graph_data.elements.edges.position`

Required:

- `x` (Number)
- `y` (Number)



<a id="nestedatt--graph_data--elements--nodes"></a>
### Nested Schema for `graph_data.elements.nodes`

Required:

- `data` (Attributes) (see [below for nested schema](#nestedatt--graph_data--elements--nodes--data))

Optional:

- `classes` (String)
- `grabbable` (Boolean)
- `group` (String)
- `locked` (Boolean)
- `pannable` (Boolean)
- `position` (Attributes) (see [below for nested schema](#nestedatt--graph_data--elements--nodes--position))
- `removed` (Boolean)
- `selectable` (Boolean)
- `selected` (Boolean)

<a id="nestedatt--graph_data--elements--nodes--data"></a>
### Nested Schema for `graph_data.elements.nodes.data`

Required:

- `node_type` (String)

Optional:

- `capability_name` (String)
- `connection_id` (String)
- `connector_id` (String)
- `id_unique` (String)
- `label` (String)
- `name` (String)
- `properties` (String, Sensitive)
- `status` (String)
- `type` (String)


<a id="nestedatt--graph_data--elements--nodes--position"></a>
### Nested Schema for `graph_data.elements.nodes.position`

Required:

- `x` (Number)
- `y` (Number)




<a id="nestedatt--graph_data--pan"></a>
### Nested Schema for `graph_data.pan`

Required:

- `x` (Number)
- `y` (Number)



<a id="nestedatt--input_schema"></a>
### Nested Schema for `input_schema`

Optional:

- `description` (String)
- `is_expanded` (Boolean)
- `preferred_control_type` (String) Options are `button`, `colorPicker`, `contentEditableTextArea`, `cssArea`, `dropDown`, `dropDownMultiSelect`, `dropDownMultiSelect2`, `dropDownWithCreate`, `functionArgumentList`, `keyValueList`, `label`, `radioSelect`, `textArea`, `textField`, `textFieldArrayView`, `toggleSwitch`.
- `preferred_data_type` (String) Options are `array`, `boolean`, `number`, `object`, `string`.
- `property_name` (String)
- `required` (Boolean)


<a id="nestedatt--output_schema"></a>
### Nested Schema for `output_schema`

Required:

- `output` (Attributes) (see [below for nested schema](#nestedatt--output_schema--output))

<a id="nestedatt--output_schema--output"></a>
### Nested Schema for `output_schema.output`

Required:

- `type` (String)

Optional:

- `additional_properties` (Boolean)
- `properties` (String)



<a id="nestedatt--settings"></a>
### Nested Schema for `settings`

Optional:

- `csp` (String)
- `css` (String)
- `css_links` (Set of String)
- `custom_error_screen_brand_logo_url` (String)
- `custom_error_show_footer` (Boolean)
- `custom_favicon_link` (String)
- `custom_logo_urlselection` (Number)
- `custom_title` (String)
- `default_error_screen_brand_logo` (Boolean)
- `flow_http_timeout_in_seconds` (Number)
- `flow_timeout_in_seconds` (Number)
- `intermediate_loading_screen_css` (String)
- `intermediate_loading_screen_html` (String)
- `js_custom_flow_player` (String)
- `js_links` (Attributes Set) (see [below for nested schema](#nestedatt--settings--js_links))
- `log_level` (Number)
- `require_authentication_to_initiate` (Boolean)
- `scrub_sensitive_info` (Boolean)
- `sensitive_info_fields` (Set of String)
- `use_csp` (Boolean)
- `use_custom_css` (Boolean)
- `use_custom_flow_player` (Boolean)
- `use_custom_script` (Boolean)
- `use_intermediate_loading_screen` (Boolean)
- `validate_on_save` (Boolean)

<a id="nestedatt--settings--js_links"></a>
### Nested Schema for `settings.js_links`

Required:

- `crossorigin` (String)
- `defer` (Boolean)
- `integrity` (String)
- `label` (String)
- `referrerpolicy` (String)
- `type` (String)
- `value` (String)



<a id="nestedatt--trigger"></a>
### Nested Schema for `trigger`

Required:

- `type` (String)

Optional:

- `configuration` (Attributes) (see [below for nested schema](#nestedatt--trigger--configuration))

<a id="nestedatt--trigger--configuration"></a>
### Nested Schema for `trigger.configuration`

Required:

- `mfa` (Attributes) (see [below for nested schema](#nestedatt--trigger--configuration--mfa))
- `pwd` (Attributes) (see [below for nested schema](#nestedatt--trigger--configuration--pwd))

<a id="nestedatt--trigger--configuration--mfa"></a>
### Nested Schema for `trigger.configuration.mfa`

Required:

- `enabled` (Boolean)
- `time` (Number)
- `time_format` (String)


<a id="nestedatt--trigger--configuration--pwd"></a>
### Nested Schema for `trigger.configuration.pwd`

Required:

- `enabled` (Boolean)
- `time` (Number)
- `time_format` (String)




<a id="nestedatt--connectors"></a>
### Nested Schema for `connectors`

Read-Only:

- `id` (String)

## Import

Import is supported using the following syntax, where attributes in `<>` brackets are replaced with the relevant ID.  For example, `<environment_id>` should be replaced with the ID of the environment to import from.

```shell
terraform import pingone_davinci_flow.example <environment_id>/<flow_id>
```