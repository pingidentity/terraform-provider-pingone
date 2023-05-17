---
page_title: "pingone_risk_predictor Resource - terraform-provider-pingone"
subcategory: "Risk"
description: |-
  Resource to manage risk predictors in a PingOne environment.
---

# pingone_risk_predictor (Resource)

Resource to manage risk predictors in a PingOne environment.

## Example Usage - Anonymous Network Predictor

```terraform
resource "pingone_risk_predictor" "my_awesome_anonymous_network_predictor" {
  environment_id = pingone_environment.my_environment.id
  name           = "My Awesome Anonymous Network Predictor"
  compact_name   = "my_awesome_anonymous_network_predictor"

  default = {
    result = {
      level = "MEDIUM"
    }
  }

  predictor_anonymous_network = {
    allowed_cidr_list = ["10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/24"]
  }
}
```

## Example Usage - Composite Predictor

```terraform
resource "pingone_risk_predictor" "my_awesome_anonymous_network_predictor" {
  environment_id = pingone_environment.my_environment.id
  name           = "My Awesome Anonymous Network Predictor"
  compact_name   = "my_awesome_anonymous_network_predictor"

  predictor_anonymous_network = {}
}

resource "pingone_risk_predictor" "my_awesome_geovelocity_predictor" {
  environment_id = pingone_environment.my_environment.id
  name           = "My Awesome Geovelocity Predictor"
  compact_name   = "my_awesome_geovelocity_predictor"

  predictor_geovelocity = {}
}

resource "pingone_risk_predictor" "my_awesome_composite_predictor" {
  environment_id = pingone_environment.my_environment.id
  name           = "My Awesome Composite Predictor"
  compact_name   = "my_awesome_composite_predictor"

  predictor_composite = {
    composition = {
      level = "HIGH"

      condition_json = jsonencode({
        "not" : {
          "or" : [{
            "equals" : 0,
            "value" : "$${details.counters.predictorLevels.medium}",
            "type" : "VALUE_COMPARISON"
            }, {
            "equals" : "High",
            "value" : "$${details.${pingone_risk_predictor.my_awesome_geovelocity_predictor.compact_name}.level}",
            "type" : "VALUE_COMPARISON"
            }, {
            "and" : [{
              "equals" : "High",
              "value" : "$${details.${pingone_risk_predictor.my_awesome_anonymous_network_predictor.compact_name}.level}",
              "type" : "VALUE_COMPARISON"
            }],
            "type" : "AND"
          }],
          "type" : "OR"
        },
        "type" : "NOT"
      })
    }
  }
}
```

## Example Usage - Custom Predictor

```terraform
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
        max_score = 6
        min_score = 5
      }

      medium = {
        max_score = 4
        min_score = 3
      }

      low = {
        max_score = 2
        min_score = 1
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
        values = ["HIGH", "HIGH321", "HIGH123"]
      }

      medium = {
        values = ["MEDIUM", "MED321", "MED123"]
      }

      low = {
        values = ["LOW"]
      }
    }
  }
}
```

## Example Usage - Geovelocity Predictor

```terraform
resource "pingone_risk_predictor" "my_awesome_geovelocity_predictor" {
  environment_id = pingone_environment.my_environment.id
  name           = "My Awesome Geovelocity Predictor"
  compact_name   = "my_awesome_geovelocity_predictor"

  default = {
    result = {
      level = "MEDIUM"
    }
  }

  predictor_geovelocity = {
    allowed_cidr_list = ["10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/24"]
  }
}
```

## Example Usage - IP Reputation Predictor

```terraform
resource "pingone_risk_predictor" "my_awesome_ip_reputation_predictor" {
  environment_id = pingone_environment.my_environment.id
  name           = "My Awesome IP Reputation Predictor"
  compact_name   = "my_awesome_ip_reputation_predictor"

  default = {
    result = {
      level = "MEDIUM"
    }
  }

  predictor_ip_reputation = {
    allowed_cidr_list = ["10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/24"]
  }
}
```

## Example Usage - New Device Predictor

```terraform
resource "pingone_risk_predictor" "my_awesome_new_device_predictor" {
  environment_id = pingone_environment.my_environment.id
  name           = "My Awesome New Device Predictor"
  compact_name   = "my_awesome_new_device_predictor"

  default = {
    result = {
      level = "MEDIUM"
    }
  }

  predictor_device = {
    detect        = "NEW_DEVICE"
    activation_at = "2023-05-01T00:00:00Z"
  }
}
```

## Example Usage - User Location Predictor

```terraform
resource "pingone_risk_predictor" "my_awesome_user_location_predictor" {
  environment_id = pingone_environment.my_environment.id
  name           = "My Awesome User Location Predictor"
  compact_name   = "my_awesome_user_location_predictor"

  default = {
    result = {
      level = "MEDIUM"
    }
  }

  predictor_user_location_anomaly = {
    radius = {
      distance = 100
      unit     = "miles"
    }
  }
}
```

## Example Usage - User Risk Behavior Predictor

```terraform
resource "pingone_risk_predictor" "my_awesome_user_risk_behavior_predictor_by_user" {
  environment_id = pingone_environment.my_environment.id
  name           = "My Awesome User Risk Behavior Predictor By User"
  compact_name   = "my_awesome_user_risk_behavior_predictor_by_user"

  default = {
    result = {
      level = "MEDIUM"
    }
  }

  predictor_user_risk_behavior = {
    prediction_model = {
      name = "points"
    }
  }
}

resource "pingone_risk_predictor" "my_awesome_user_risk_behavior_predictor_by_organization" {
  environment_id = pingone_environment.my_environment.id
  name           = "My Awesome User Risk Behavior Predictor By Organization"
  compact_name   = "my_awesome_user_risk_behavior_predictor_by_organization"

  default = {
    result = {
      level = "MEDIUM"
    }
  }

  predictor_user_risk_behavior = {
    prediction_model = {
      name = "login_anomaly_statistic"
    }
  }
}
```

## Example Usage - Velocity Predictor

```terraform
resource "pingone_risk_predictor" "my_awesome_velocity_predictor_by_ip" {
  environment_id = pingone_environment.my_environment.id
  name           = "My Awesome Velocity Predictor By IP"
  compact_name   = "my_awesome_velocity_predictor_by_ip"

  default = {
    result = {
      level = "MEDIUM"
    }
  }

  predictor_velocity = {
    of = "$${event.user.id}"
  }
}

resource "pingone_risk_predictor" "my_awesome_velocity_predictor_by_user" {
  environment_id = pingone_environment.my_environment.id
  name           = "My Awesome Velocity Predictor By User"
  compact_name   = "my_awesome_velocity_predictor_by_user"

  default = {
    result = {
      level = "MEDIUM"
    }
  }

  predictor_velocity = {
    of = "$${event.ip}"
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `compact_name` (String) A string that specifies the unique name for the predictor for use in risk evaluation request/response payloads. This property is immutable; it cannot be modified after initial creation. The value must be alpha-numeric, with no special characters or spaces. This name is used in the API both for policy configuration, and in the Risk Evaluation response (under details).
- `environment_id` (String) The ID of the environment to configure the risk predictor in.
- `name` (String) A string that specifies the unique, friendly name for the predictor. This name is displayed in the Risk Policies UI, when the admin is asked to define the overrides and weights in policy configuration.

### Optional

- `default` (Attributes) (see [below for nested schema](#nestedatt--default))
- `description` (String) A string that specifies the description of the risk predictor. Maximum length is 1024 characters.
- `predictor_anonymous_network` (Attributes) A single nested attribute that specifies options for the Anonymous Network predictor. (see [below for nested schema](#nestedatt--predictor_anonymous_network))
- `predictor_composite` (Attributes) A single nested attribute that specifies options for the Composite predictor. (see [below for nested schema](#nestedatt--predictor_composite))
- `predictor_custom_map` (Attributes) A single nested attribute that specifies options for the Composite predictor. (see [below for nested schema](#nestedatt--predictor_custom_map))
- `predictor_device` (Attributes) A single nested attribute that specifies options for the Device predictor. (see [below for nested schema](#nestedatt--predictor_device))
- `predictor_geovelocity` (Attributes) A single nested attribute that specifies options for the Geovelocity predictor. (see [below for nested schema](#nestedatt--predictor_geovelocity))
- `predictor_ip_reputation` (Attributes) A single nested attribute that specifies options for the IP reputation predictor. (see [below for nested schema](#nestedatt--predictor_ip_reputation))
- `predictor_user_location_anomaly` (Attributes) A single nested attribute that specifies options for the User Location Anomaly predictor. (see [below for nested schema](#nestedatt--predictor_user_location_anomaly))
- `predictor_user_risk_behavior` (Attributes) A single nested attribute that specifies options for the User Risk Behavior predictor. (see [below for nested schema](#nestedatt--predictor_user_risk_behavior))
- `predictor_velocity` (Attributes) A single nested attribute that specifies options for the Velocity predictor. (see [below for nested schema](#nestedatt--predictor_velocity))

### Read-Only

- `deletable` (Boolean) A boolean that indicates the PingOne Risk predictor can be deleted or not.
- `id` (String) The ID of this resource.
- `licensed` (Boolean) A boolean that indicates whether PingOne Risk is licensed for the environment.
- `type` (String) A string that specifies the type of the risk predictor.  This can be either `ANONYMOUS_NETWORK`, `COMPOSITE`, `GEO_VELOCITY`, `IP_REPUTATION`, `MAP`, `DEVICE`, `USER_LOCATION_ANOMALY`, `USER_RISK_BEHAVIOR` or `VELOCITY`.

<a id="nestedatt--default"></a>
### Nested Schema for `default`

Optional:

- `result` (Attributes) (see [below for nested schema](#nestedatt--default--result))
- `weight` (Number) A number that specifies the default weight for the risk predictor. This value is used when the risk predictor is not explicitly configured in a policy.

<a id="nestedatt--default--result"></a>
### Nested Schema for `default.result`

Optional:

- `level` (String) A string that specifies the type of the risk predictor.  This can be either `ANONYMOUS_NETWORK`, `COMPOSITE`, `GEO_VELOCITY`, `IP_REPUTATION`, `MAP`, `DEVICE`, `USER_LOCATION_ANOMALY`, `USER_RISK_BEHAVIOR` or `VELOCITY`.
- `type` (String) A string that specifies the type of the risk predictor.  This can be either `ANONYMOUS_NETWORK`, `COMPOSITE`, `GEO_VELOCITY`, `IP_REPUTATION`, `MAP`, `DEVICE`, `USER_LOCATION_ANOMALY`, `USER_RISK_BEHAVIOR` or `VELOCITY`.



<a id="nestedatt--predictor_anonymous_network"></a>
### Nested Schema for `predictor_anonymous_network`

Optional:

- `allowed_cidr_list` (Set of String)


<a id="nestedatt--predictor_composite"></a>
### Nested Schema for `predictor_composite`

Required:

- `composition` (Attributes) (see [below for nested schema](#nestedatt--predictor_composite--composition))

<a id="nestedatt--predictor_composite--composition"></a>
### Nested Schema for `predictor_composite.composition`

Required:

- `condition_json` (String) A string that specifies the condition for the composite risk predictor. The value must be a valid JSON string.
- `level` (String) A string that specifies the risk level for the composite risk predictor. The value must be one of the following: LOW, MEDIUM, HIGH.

Read-Only:

- `condition` (String) A string that specifies the condition for the composite risk predictor. The value must be a valid JSON string.



<a id="nestedatt--predictor_custom_map"></a>
### Nested Schema for `predictor_custom_map`

Required:

- `contains` (String) A string that specifies the value to match in the custom map. Maximum length is 1024 characters.

Optional:

- `between_ranges` (Attributes) (see [below for nested schema](#nestedatt--predictor_custom_map--between_ranges))
- `ip_ranges` (Attributes) (see [below for nested schema](#nestedatt--predictor_custom_map--ip_ranges))
- `string_list` (Attributes) (see [below for nested schema](#nestedatt--predictor_custom_map--string_list))

Read-Only:

- `type` (String) A string that specifies the type of the risk predictor.  This can be either "ANONYMOUS_NETWORK", "COMPOSITE", "GEO_VELOCITY", "IP_REPUTATION", "MAP", "DEVICE", "USER_LOCATION_ANOMALY", "USER_RISK_BEHAVIOR" or "VELOCITY".

<a id="nestedatt--predictor_custom_map--between_ranges"></a>
### Nested Schema for `predictor_custom_map.between_ranges`

Optional:

- `high` (Attributes) (see [below for nested schema](#nestedatt--predictor_custom_map--between_ranges--high))
- `low` (Attributes) (see [below for nested schema](#nestedatt--predictor_custom_map--between_ranges--low))
- `medium` (Attributes) (see [below for nested schema](#nestedatt--predictor_custom_map--between_ranges--medium))

<a id="nestedatt--predictor_custom_map--between_ranges--high"></a>
### Nested Schema for `predictor_custom_map.between_ranges.high`

Required:

- `max_score` (Number) A number that specifies the maximum score for the risk predictor. This value is used when the risk predictor is not explicitly configured in a policy.
- `min_score` (Number) A number that specifies the minimum score for the risk predictor. This value is used when the risk predictor is not explicitly configured in a policy.


<a id="nestedatt--predictor_custom_map--between_ranges--low"></a>
### Nested Schema for `predictor_custom_map.between_ranges.low`

Required:

- `max_score` (Number) A number that specifies the maximum score for the risk predictor. This value is used when the risk predictor is not explicitly configured in a policy.
- `min_score` (Number) A number that specifies the minimum score for the risk predictor. This value is used when the risk predictor is not explicitly configured in a policy.


<a id="nestedatt--predictor_custom_map--between_ranges--medium"></a>
### Nested Schema for `predictor_custom_map.between_ranges.medium`

Required:

- `max_score` (Number) A number that specifies the maximum score for the risk predictor. This value is used when the risk predictor is not explicitly configured in a policy.
- `min_score` (Number) A number that specifies the minimum score for the risk predictor. This value is used when the risk predictor is not explicitly configured in a policy.



<a id="nestedatt--predictor_custom_map--ip_ranges"></a>
### Nested Schema for `predictor_custom_map.ip_ranges`

Optional:

- `high` (Attributes) (see [below for nested schema](#nestedatt--predictor_custom_map--ip_ranges--high))
- `low` (Attributes) (see [below for nested schema](#nestedatt--predictor_custom_map--ip_ranges--low))
- `medium` (Attributes) (see [below for nested schema](#nestedatt--predictor_custom_map--ip_ranges--medium))

<a id="nestedatt--predictor_custom_map--ip_ranges--high"></a>
### Nested Schema for `predictor_custom_map.ip_ranges.high`

Optional:

- `values` (Set of String)


<a id="nestedatt--predictor_custom_map--ip_ranges--low"></a>
### Nested Schema for `predictor_custom_map.ip_ranges.low`

Optional:

- `values` (Set of String)


<a id="nestedatt--predictor_custom_map--ip_ranges--medium"></a>
### Nested Schema for `predictor_custom_map.ip_ranges.medium`

Optional:

- `values` (Set of String)



<a id="nestedatt--predictor_custom_map--string_list"></a>
### Nested Schema for `predictor_custom_map.string_list`

Optional:

- `high` (Attributes) (see [below for nested schema](#nestedatt--predictor_custom_map--string_list--high))
- `low` (Attributes) (see [below for nested schema](#nestedatt--predictor_custom_map--string_list--low))
- `medium` (Attributes) (see [below for nested schema](#nestedatt--predictor_custom_map--string_list--medium))

<a id="nestedatt--predictor_custom_map--string_list--high"></a>
### Nested Schema for `predictor_custom_map.string_list.high`

Optional:

- `values` (Set of String)


<a id="nestedatt--predictor_custom_map--string_list--low"></a>
### Nested Schema for `predictor_custom_map.string_list.low`

Optional:

- `values` (Set of String)


<a id="nestedatt--predictor_custom_map--string_list--medium"></a>
### Nested Schema for `predictor_custom_map.string_list.medium`

Optional:

- `values` (Set of String)




<a id="nestedatt--predictor_device"></a>
### Nested Schema for `predictor_device`

Optional:

- `activation_at` (String) You can use the `activation_at` parameter to specify a date on which the learning process for the predictor should be restarted. This can be used in conjunction with the fallback setting (`default.result.level`) to force strong authentication when moving the predictor to production. The date should be in an RFC3339 format. Note that activation date uses UTC time.
- `detect` (String)


<a id="nestedatt--predictor_geovelocity"></a>
### Nested Schema for `predictor_geovelocity`

Optional:

- `allowed_cidr_list` (Set of String)


<a id="nestedatt--predictor_ip_reputation"></a>
### Nested Schema for `predictor_ip_reputation`

Optional:

- `allowed_cidr_list` (Set of String)


<a id="nestedatt--predictor_user_location_anomaly"></a>
### Nested Schema for `predictor_user_location_anomaly`

Optional:

- `days` (Number)
- `radius` (Attributes) (see [below for nested schema](#nestedatt--predictor_user_location_anomaly--radius))

<a id="nestedatt--predictor_user_location_anomaly--radius"></a>
### Nested Schema for `predictor_user_location_anomaly.radius`

Required:

- `distance` (Number)

Optional:

- `unit` (String)



<a id="nestedatt--predictor_user_risk_behavior"></a>
### Nested Schema for `predictor_user_risk_behavior`

Required:

- `prediction_model` (Attributes) (see [below for nested schema](#nestedatt--predictor_user_risk_behavior--prediction_model))

<a id="nestedatt--predictor_user_risk_behavior--prediction_model"></a>
### Nested Schema for `predictor_user_risk_behavior.prediction_model`

Required:

- `name` (String)



<a id="nestedatt--predictor_velocity"></a>
### Nested Schema for `predictor_velocity`

Optional:

- `measure` (String)
- `of` (String)

Read-Only:

- `by` (Set of String)
- `every` (Attributes) An object that contains configuration values for the every risk predictor type. (see [below for nested schema](#nestedatt--predictor_velocity--every))
- `fallback` (Attributes) An object that contains configuration values for the fallback risk predictor type. (see [below for nested schema](#nestedatt--predictor_velocity--fallback))
- `sliding_window` (Attributes) An object that contains configuration values for the sliding window risk predictor type. (see [below for nested schema](#nestedatt--predictor_velocity--sliding_window))
- `use` (Attributes) (see [below for nested schema](#nestedatt--predictor_velocity--use))

<a id="nestedatt--predictor_velocity--every"></a>
### Nested Schema for `predictor_velocity.every`

Read-Only:

- `min_sample` (Number) The minimum number of samples to use for the risk predictor.
- `quantity` (Number) The number of `unit` intervals to use for the risk predictor.
- `unit` (String) The unit of measurement for the `interval` parameter.


<a id="nestedatt--predictor_velocity--fallback"></a>
### Nested Schema for `predictor_velocity.fallback`

Read-Only:

- `high` (Number) The high risk level.
- `medium` (Number) The medium risk level.
- `strategy` (String) The strategy to use when the risk predictor is not able to determine a risk level.


<a id="nestedatt--predictor_velocity--sliding_window"></a>
### Nested Schema for `predictor_velocity.sliding_window`

Read-Only:

- `min_sample` (Number) The minimum number of samples to use for the risk predictor.
- `quantity` (Number) The number of `unit` intervals to use for the risk predictor.
- `unit` (String) The unit of measurement for the `interval` parameter.


<a id="nestedatt--predictor_velocity--use"></a>
### Nested Schema for `predictor_velocity.use`

Read-Only:

- `high` (Number) The high risk level.
- `medium` (Number) The medium risk level.
- `type` (String) The type of the risk predictor.

## Import

Import is supported using the following syntax, where attributes in `<>` brackets are replaced with the relevant ID.  For example, `<environment_id>` should be replaced with the ID of the environment to import from.

```shell
$ terraform import pingone_risk_predictor.example <environment_id>/<risk_predictor_id>
```
