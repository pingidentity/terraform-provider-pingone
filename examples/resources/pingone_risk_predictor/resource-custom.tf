resource "pingone_risk_predictor" "my_awesome_custom_predictor_between_ranges" {
  environment_id = pingone_environment.my_environment.id
  name           = "My Awesome Custom Predictor Between Ranges"
  compact_name   = "my_awesome_custom_predictor_between_ranges"

  default = {
    result = {
      level = "MEDIUM"
    }
  }

  predictor_custom_map = {
    contains = "$${event.my_custom_field}"

    between_ranges = {
      high = {
        max_value = 6
        min_value = 5
      }

      medium = {
        max_value = 4
        min_value = 3
      }

      low = {
        max_value = 2
        min_value = 1
      }
    }
  }
}

resource "pingone_risk_predictor" "my_awesome_custom_predictor_ip_ranges" {
  environment_id = pingone_environment.my_environment.id
  name           = "My Awesome Custom Predictor IP Ranges"
  compact_name   = "my_awesome_custom_predictor_ip_ranges"

  default = {
    result = {
      level = "MEDIUM"
    }
  }

  predictor_custom_map = {
    contains = "$${event.my_custom_field}"

    ip_ranges = {
      high = {
        values = ["192.168.0.0/24", "10.0.0.0/8", "172.16.0.0/12"
        ]
      }

      medium = {
        values = ["192.0.2.0/24", "192.168.1.0/26", "10.10.0.0/16"]
      }

      low = {
        values = [
          "172.16.0.0/16"
        ]
      }
    }
  }
}

resource "pingone_risk_predictor" "my_awesome_custom_predictor_list" {
  environment_id = pingone_environment.my_environment.id
  name           = "My Awesome Custom Predictor List"
  compact_name   = "my_awesome_custom_predictor_list"

  default = {
    result = {
      level = "MEDIUM"
    }
  }

  predictor_custom_map = {
    contains = "$${event.my_custom_field}"

    string_list = {
      high = {
        values = ["HIGH"]
      }

      medium = {
        values = ["MEDIUM"]
      }

      low = {
        values = ["LOW"]
      }
    }
  }
}