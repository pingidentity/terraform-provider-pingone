resource "pingone_davinci_connector_instance" "tmtConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "tmtConnector"
  }
  name = "My awesome tmtConnector"
  property {
    name  = "apiKey"
    type  = "string"
    value = var.tmtconnector_property_api_key
  }
  property {
    name  = "apiSecret"
    type  = "string"
    value = var.tmtconnector_property_api_secret
  }
  property {
    name  = "apiUrl"
    type  = "string"
    value = var.tmtconnector_property_api_url
  }
  property {
    name  = "city"
    type  = "string"
    value = var.tmtconnector_property_city
  }
  property {
    name  = "country"
    type  = "string"
    value = var.tmtconnector_property_country
  }
  property {
    name  = "dataPoints"
    type  = "string"
    value = var.tmtconnector_property_data_points
  }
  property {
    name  = "day"
    type  = "number"
    value = var.tmtconnector_property_day
  }
  property {
    name  = "email"
    type  = "string"
    value = var.tmtconnector_property_email
  }
  property {
    name  = "first_name"
    type  = "string"
    value = var.tmtconnector_property_first_name
  }
  property {
    name  = "last_name"
    type  = "string"
    value = var.tmtconnector_property_last_name
  }
  property {
    name  = "middle_name"
    type  = "string"
    value = var.tmtconnector_property_middle_name
  }
  property {
    name  = "month"
    type  = "number"
    value = var.tmtconnector_property_month
  }
  property {
    name  = "phoneNumber"
    type  = "string"
    value = var.tmtconnector_property_phone_number
  }
  property {
    name  = "postcode"
    type  = "string"
    value = var.tmtconnector_property_postcode
  }
  property {
    name  = "province"
    type  = "string"
    value = var.tmtconnector_property_province
  }
  property {
    name  = "street"
    type  = "string"
    value = var.tmtconnector_property_street
  }
  property {
    name  = "street_no"
    type  = "string"
    value = var.tmtconnector_property_street_no
  }
  property {
    name  = "unit_name"
    type  = "string"
    value = var.tmtconnector_property_unit_name
  }
  property {
    name  = "unit_no"
    type  = "string"
    value = var.tmtconnector_property_unit_no
  }
  property {
    name  = "year"
    type  = "number"
    value = var.tmtconnector_property_year
  }
}
