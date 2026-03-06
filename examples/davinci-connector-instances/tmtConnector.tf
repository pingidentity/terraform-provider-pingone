resource "pingone_davinci_connector_instance" "tmtConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "tmtConnector"
  }
  name = "My awesome tmtConnector"
  properties = jsonencode({
    "apiKey" = var.tmtconnector_property_api_key
    "apiSecret" = var.tmtconnector_property_api_secret
    "apiUrl" = var.tmtconnector_property_api_url
  })
}
