---
page_title: "pingone_risk_predictor Resource - terraform-provider-pingone"
subcategory: "Protect"
description: |-
  Resource to manage Risk predictors in a PingOne environment.
---

# pingone_risk_predictor (Resource)

Resource to manage Risk predictors in a PingOne environment.

## Example Usage - Adversary-In-The-Middle

```terraform
resource "pingone_risk_predictor" "my_awesome_adversary_in_the_middle_predictor" {
  environment_id = pingone_environment.my_environment.id
  name           = "My Awesome Adversary In The Middle Predictor"
  compact_name   = "myAwesomeAitMPredictor"

  default = {
    result = {
      level = "MEDIUM"
    }
  }

  predictor_adversary_in_the_middle = {
    allowed_domain_list = ["domain1.com", "domain2.com", "domain3.com"]
  }
}
```

## Example Usage - Anonymous Network Predictor

```terraform
resource "pingone_risk_predictor" "my_awesome_anonymous_network_predictor" {
  environment_id = pingone_environment.my_environment.id
  name           = "My Awesome Anonymous Network Predictor"
  compact_name   = "myAwesomeAnonymousNetworkPredictor"

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

## Example Usage - Bot Detection Predictor

```terraform
resource "pingone_risk_predictor" "my_awesome_bot_detection_predictor" {
  environment_id = pingone_environment.my_environment.id
  name           = "My Awesome Bot Detection Predictor"
  compact_name   = "myAwesomeBotDetectionPredictor"

  default = {
    result = {
      level = "MEDIUM"
    }
  }

  predictor_bot_detection = {}
}
```

## Example Usage - Composite Predictor

```terraform
resource "pingone_risk_predictor" "my_awesome_anonymous_network_predictor" {
  environment_id = pingone_environment.my_environment.id
  name           = "My Awesome Anonymous Network Predictor"
  compact_name   = "myAwesomeAnonymousNetworkPredictor"

  predictor_anonymous_network = {}
}

resource "pingone_risk_predictor" "my_awesome_geovelocity_anomaly_predictor" {
  environment_id = pingone_environment.my_environment.id
  name           = "My Awesome Geovelocity Predictor"
  compact_name   = "myAwesomeGeovelocityPredictor"

  predictor_geovelocity = {}
}

resource "pingone_risk_predictor" "my_awesome_composite_predictor" {
  environment_id = pingone_environment.my_environment.id
  name           = "My Awesome Composite Predictor"
  compact_name   = "myAwesomeCompositePredictor"

  predictor_composite = {
    compositions = [
      {
        level = "HIGH"

        condition_json = jsonencode({
          "not" : {
            "or" : [{
              "equals" : 0,
              "value" : "$${details.counters.predictorLevels.medium}",
              "type" : "VALUE_COMPARISON"
              }, {
              "equals" : "High",
              "value" : "$${details.${pingone_risk_predictor.my_awesome_geovelocity_anomaly_predictor.compact_name}.level}",
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
    ]
  }
}
```

## Example Usage - Custom Predictor

```terraform
resource "pingone_risk_predictor" "my_awesome_custom_predictor_between_ranges" {
  environment_id = pingone_environment.my_environment.id
  name           = "My Awesome Custom Predictor Between Ranges"
  compact_name   = "myAwesomeCustomPredictorBetweenRanges"

  default = {
    result = {
      level = "MEDIUM"
    }
  }

  predictor_custom_map = {
    contains = "$${event.myCustomField}"

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
  compact_name   = "myAwesomeCustomPredictorIpRanges"

  default = {
    result = {
      level = "MEDIUM"
    }
  }

  predictor_custom_map = {
    contains = "$${event.myCustomField}"

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
  compact_name   = "myAwesomeCustomPredictorList"

  default = {
    result = {
      level = "MEDIUM"
    }
  }

  predictor_custom_map = {
    contains = "$${event.myCustomField}"

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
```

## Example Usage - Email Reputation Predictor

```terraform
resource "pingone_risk_predictor" "my_awesome_email_reputation_predictor" {
  environment_id = pingone_environment.my_environment.id
  name           = "My Awesome Email Reputation Predictor"
  compact_name   = "myAwesomeEmailReputationPredictor"

  default = {
    result = {
      level = "MEDIUM"
    }
  }

  predictor_email_reputation = {}
}
```

## Example Usage - Geovelocity Predictor

```terraform
resource "pingone_risk_predictor" "my_awesome_geovelocity_anomaly_predictor" {
  environment_id = pingone_environment.my_environment.id
  name           = "My Awesome Geovelocity Predictor"
  compact_name   = "myAwesomeGeovelocityPredictor"

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
  compact_name   = "myAwesomeIpReputationPredictor"

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
resource "time_static" "my_awesome_new_device_predictor_activation" {}

resource "pingone_risk_predictor" "my_awesome_new_device_predictor" {
  environment_id = pingone_environment.my_environment.id
  name           = "My Awesome New Device Predictor"
  compact_name   = "myAwesomeNewDevicePredictor"

  default = {
    result = {
      level = "MEDIUM"
    }
  }

  predictor_device = {
    detect        = "NEW_DEVICE"
    activation_at = format("%sT00:00:00Z", formatdate("YYYY-MM-DD", time_static.my_awesome_new_device_predictor_activation.rfc3339))
  }
}
```

## Example Usage - Suspicious Device Predictor

```terraform
resource "pingone_risk_predictor" "my_awesome_suspicious_device_predictor" {
  environment_id = pingone_environment.my_environment.id
  name           = "My Awesome Suspicious Device Predictor"
  compact_name   = "myAwesomeSuspiciousDevicePredictor"

  default = {
    result = {
      level = "MEDIUM"
    }
  }

  predictor_device = {
    detect = "SUSPICIOUS_DEVICE"
  }
}
```

## Example Usage - User Location Predictor

```terraform
resource "pingone_risk_predictor" "my_awesome_user_location_predictor" {
  environment_id = pingone_environment.my_environment.id
  name           = "My Awesome User Location Predictor"
  compact_name   = "myAwesomeUserLocationPredictor"

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
  compact_name   = "myAwesomeUserRiskBehaviorPredictorByUser"

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
  compact_name   = "myAwesomeUserRiskBehaviorPredictorByOrganization"

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
  compact_name   = "myAwesomeVelocityPredictorByIp"

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
  compact_name   = "myAwesomeVelocityPredictorByUser"

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

- `compact_name` (String) A string that specifies the unique name for the predictor for use in risk evaluation request/response payloads. The value must be alpha-numeric, with no special characters or spaces. This name is used in the API both for policy configuration, and in the Risk Evaluation response (under `details`).  If the value used for `compact_name` relates to a built-in predictor (a predictor that cannot be deleted), then this resource will attempt to overwrite the predictor's configuration.  This field is immutable and will trigger a replace plan if changed.
- `environment_id` (String) The ID of the environment to configure the risk predictor in.  Must be a valid PingOne resource ID.  This field is immutable and will trigger a replace plan if changed.
- `name` (String) A string that specifies the unique, friendly name for the predictor. This name is displayed in the Risk Policies UI, when the admin is asked to define the overrides and weights in policy configuration and is unique per environment.

### Optional

- `default` (Attributes) A single nested object that specifies the default configuration values for the risk predictor. (see [below for nested schema](#nestedatt--default))
- `description` (String) A string that specifies the description of the risk predictor. Maximum length is 1024 characters.
- `predictor_adversary_in_the_middle` (Attributes) A single nested object that specifies options for the Adversary-In-The-Middle (AitM) predictor.  Exactly one of the following must be defined: `predictor_adversary_in_the_middle`, `predictor_anonymous_network`, `predictor_bot_detection`, `predictor_composite`, `predictor_custom_map`, `predictor_device`, `predictor_email_reputation`, `predictor_geovelocity`, `predictor_ip_reputation`, `predictor_traffic_anomaly`, `predictor_user_location_anomaly`, `predictor_user_risk_behavior`, `predictor_velocity`. (see [below for nested schema](#nestedatt--predictor_adversary_in_the_middle))
- `predictor_anonymous_network` (Attributes) A single nested object that specifies options for the Anonymous Network predictor.  Exactly one of the following must be defined: `predictor_adversary_in_the_middle`, `predictor_anonymous_network`, `predictor_bot_detection`, `predictor_composite`, `predictor_custom_map`, `predictor_device`, `predictor_email_reputation`, `predictor_geovelocity`, `predictor_ip_reputation`, `predictor_traffic_anomaly`, `predictor_user_location_anomaly`, `predictor_user_risk_behavior`, `predictor_velocity`. (see [below for nested schema](#nestedatt--predictor_anonymous_network))
- `predictor_bot_detection` (Attributes) A single nested object that specifies options for the Bot Detection predictor.  Exactly one of the following must be defined: `predictor_adversary_in_the_middle`, `predictor_anonymous_network`, `predictor_bot_detection`, `predictor_composite`, `predictor_custom_map`, `predictor_device`, `predictor_email_reputation`, `predictor_geovelocity`, `predictor_ip_reputation`, `predictor_traffic_anomaly`, `predictor_user_location_anomaly`, `predictor_user_risk_behavior`, `predictor_velocity`. (see [below for nested schema](#nestedatt--predictor_bot_detection))
- `predictor_composite` (Attributes) A single nested object that specifies options for the Composite predictor.  Exactly one of the following must be defined: `predictor_adversary_in_the_middle`, `predictor_anonymous_network`, `predictor_bot_detection`, `predictor_composite`, `predictor_custom_map`, `predictor_device`, `predictor_email_reputation`, `predictor_geovelocity`, `predictor_ip_reputation`, `predictor_traffic_anomaly`, `predictor_user_location_anomaly`, `predictor_user_risk_behavior`, `predictor_velocity`. (see [below for nested schema](#nestedatt--predictor_composite))
- `predictor_custom_map` (Attributes) A single nested object that specifies options for the Custom Map predictor.  Exactly one of the following must be defined: `predictor_adversary_in_the_middle`, `predictor_anonymous_network`, `predictor_bot_detection`, `predictor_composite`, `predictor_custom_map`, `predictor_device`, `predictor_email_reputation`, `predictor_geovelocity`, `predictor_ip_reputation`, `predictor_traffic_anomaly`, `predictor_user_location_anomaly`, `predictor_user_risk_behavior`, `predictor_velocity`. (see [below for nested schema](#nestedatt--predictor_custom_map))
- `predictor_device` (Attributes) A single nested object that specifies options for the Device predictor.  Exactly one of the following must be defined: `predictor_adversary_in_the_middle`, `predictor_anonymous_network`, `predictor_bot_detection`, `predictor_composite`, `predictor_custom_map`, `predictor_device`, `predictor_email_reputation`, `predictor_geovelocity`, `predictor_ip_reputation`, `predictor_traffic_anomaly`, `predictor_user_location_anomaly`, `predictor_user_risk_behavior`, `predictor_velocity`. (see [below for nested schema](#nestedatt--predictor_device))
- `predictor_email_reputation` (Attributes) A single nested object that specifies options for the Email reputation predictor.  Exactly one of the following must be defined: `predictor_adversary_in_the_middle`, `predictor_anonymous_network`, `predictor_bot_detection`, `predictor_composite`, `predictor_custom_map`, `predictor_device`, `predictor_email_reputation`, `predictor_geovelocity`, `predictor_ip_reputation`, `predictor_traffic_anomaly`, `predictor_user_location_anomaly`, `predictor_user_risk_behavior`, `predictor_velocity`. (see [below for nested schema](#nestedatt--predictor_email_reputation))
- `predictor_geovelocity` (Attributes) A single nested object that specifies options for the Geovelocity predictor.  Exactly one of the following must be defined: `predictor_adversary_in_the_middle`, `predictor_anonymous_network`, `predictor_bot_detection`, `predictor_composite`, `predictor_custom_map`, `predictor_device`, `predictor_email_reputation`, `predictor_geovelocity`, `predictor_ip_reputation`, `predictor_traffic_anomaly`, `predictor_user_location_anomaly`, `predictor_user_risk_behavior`, `predictor_velocity`. (see [below for nested schema](#nestedatt--predictor_geovelocity))
- `predictor_ip_reputation` (Attributes) A single nested object that specifies options for the IP reputation predictor.  Exactly one of the following must be defined: `predictor_adversary_in_the_middle`, `predictor_anonymous_network`, `predictor_bot_detection`, `predictor_composite`, `predictor_custom_map`, `predictor_device`, `predictor_email_reputation`, `predictor_geovelocity`, `predictor_ip_reputation`, `predictor_traffic_anomaly`, `predictor_user_location_anomaly`, `predictor_user_risk_behavior`, `predictor_velocity`. (see [below for nested schema](#nestedatt--predictor_ip_reputation))
- `predictor_traffic_anomaly` (Attributes) A single nested object that specifies options for the Traffic Anomaly predictor.  Exactly one of the following must be defined: `predictor_adversary_in_the_middle`, `predictor_anonymous_network`, `predictor_bot_detection`, `predictor_composite`, `predictor_custom_map`, `predictor_device`, `predictor_email_reputation`, `predictor_geovelocity`, `predictor_ip_reputation`, `predictor_traffic_anomaly`, `predictor_user_location_anomaly`, `predictor_user_risk_behavior`, `predictor_velocity`. (see [below for nested schema](#nestedatt--predictor_traffic_anomaly))
- `predictor_user_location_anomaly` (Attributes) A single nested object that specifies options for the User Location Anomaly predictor.  Exactly one of the following must be defined: `predictor_adversary_in_the_middle`, `predictor_anonymous_network`, `predictor_bot_detection`, `predictor_composite`, `predictor_custom_map`, `predictor_device`, `predictor_email_reputation`, `predictor_geovelocity`, `predictor_ip_reputation`, `predictor_traffic_anomaly`, `predictor_user_location_anomaly`, `predictor_user_risk_behavior`, `predictor_velocity`. (see [below for nested schema](#nestedatt--predictor_user_location_anomaly))
- `predictor_user_risk_behavior` (Attributes) A single nested object that specifies options for the User Risk Behavior predictor.  Exactly one of the following must be defined: `predictor_adversary_in_the_middle`, `predictor_anonymous_network`, `predictor_bot_detection`, `predictor_composite`, `predictor_custom_map`, `predictor_device`, `predictor_email_reputation`, `predictor_geovelocity`, `predictor_ip_reputation`, `predictor_traffic_anomaly`, `predictor_user_location_anomaly`, `predictor_user_risk_behavior`, `predictor_velocity`. (see [below for nested schema](#nestedatt--predictor_user_risk_behavior))
- `predictor_velocity` (Attributes) A single nested object that specifies options for the Velocity predictor.  Exactly one of the following must be defined: `predictor_adversary_in_the_middle`, `predictor_anonymous_network`, `predictor_bot_detection`, `predictor_composite`, `predictor_custom_map`, `predictor_device`, `predictor_email_reputation`, `predictor_geovelocity`, `predictor_ip_reputation`, `predictor_traffic_anomaly`, `predictor_user_location_anomaly`, `predictor_user_risk_behavior`, `predictor_velocity`. (see [below for nested schema](#nestedatt--predictor_velocity))

### Read-Only

- `deletable` (Boolean) A boolean that indicates the PingOne Risk predictor can be deleted or not.
- `id` (String) The ID of this resource.
- `licensed` (Boolean) A boolean that indicates whether PingOne Risk is licensed for the environment.
- `type` (String) A string that specifies the type of the risk predictor.  Options are `ADVERSARY_IN_THE_MIDDLE`, `ANONYMOUS_NETWORK`, `BOT`, `COMPOSITE`, `DEVICE`, `EMAIL_REPUTATION`, `GEO_VELOCITY`, `IP_REPUTATION`, `MAP`, `TRAFFIC_ANOMALY`, `USER_LOCATION_ANOMALY`, `USER_RISK_BEHAVIOR`, `VELOCITY`.

<a id="nestedatt--default"></a>
### Nested Schema for `default`

Optional:

- `result` (Attributes) A single nested object that contains the result assigned to the predictor if the predictor could not be calculated during the risk evaluation. If this field is not provided, and the predictor could not be calculated during risk evaluation, the behavior is: 1) If the predictor is used in an override, the override is skipped; 2) In the weighted policy, the predictor will have a `weight` of `0`. (see [below for nested schema](#nestedatt--default--result))
- `weight` (Number) A number that specifies the default weight for the risk predictor. This value is used when the risk predictor is not explicitly configured in a policy.  Defaults to `5`.

<a id="nestedatt--default--result"></a>
### Nested Schema for `default.result`

Optional:

- `level` (String) The default result level.  Options are `HIGH`, `LOW`, `MEDIUM`.

Read-Only:

- `type` (String) The default result type.  Options are `VALUE` (any custom attribute value that's defined).



<a id="nestedatt--predictor_adversary_in_the_middle"></a>
### Nested Schema for `predictor_adversary_in_the_middle`

Optional:

- `allowed_domain_list` (Set of String) A set of domains that are ignored for the predictor results.


<a id="nestedatt--predictor_anonymous_network"></a>
### Nested Schema for `predictor_anonymous_network`

Optional:

- `allowed_cidr_list` (Set of String) A set of IP addresses (CIDRs) that are ignored for the predictor results. The list can include IPs in IPv4 format and IPs in IPv6 format.


<a id="nestedatt--predictor_bot_detection"></a>
### Nested Schema for `predictor_bot_detection`

Optional:

- `include_repeated_events_without_sdk` (Boolean) A boolean that specifies whether to expand the range of bot activity that PingOne Protect can detect.


<a id="nestedatt--predictor_composite"></a>
### Nested Schema for `predictor_composite`

Required:

- `compositions` (Attributes List) A list of compositions of risk factors you want to use, and the condition logic that determines when or whether a risk factor is applied.  The minimum number of compositions is 1 and the maximum number of compositions is 3. (see [below for nested schema](#nestedatt--predictor_composite--compositions))

<a id="nestedatt--predictor_composite--compositions"></a>
### Nested Schema for `predictor_composite.compositions`

Required:

- `condition_json` (String) A string that specifies the condition logic for the composite risk predictor. The value must be a valid JSON string.
- `level` (String) A string that specifies the risk level for the composite risk predictor.  Options are `HIGH`, `LOW`, `MEDIUM`.

Read-Only:

- `condition` (String) A string that specifies the condition logic for the composite risk predictor as applied to the service.



<a id="nestedatt--predictor_custom_map"></a>
### Nested Schema for `predictor_custom_map`

Required:

- `contains` (String) A string that specifies the attribute reference that contains the value to match in the custom map.  The attribute reference should come from either the incoming event (`${event.*}`) or the evaluation details (`${details.*}`).  When defining attribute references in Terraform, the leading `$` needs to be escaped with an additional `$` character, e.g. `contains = "$${event.myattribute}"`.

Optional:

- `between_ranges` (Attributes) A single nested object that describes the upper and lower bounds of ranges of values that apply to the attribute reference in `predictor_custom_map.contains`, that map to high, medium or low risk results. (see [below for nested schema](#nestedatt--predictor_custom_map--between_ranges))
- `ip_ranges` (Attributes) A single nested object that describes IP CIDR ranges of values that apply to the attribute reference in `predictor_custom_map.contains`, that map to high, medium or low risk results. (see [below for nested schema](#nestedatt--predictor_custom_map--ip_ranges))
- `string_list` (Attributes) A single nested object that describes the string values that apply to the attribute reference in `predictor_custom_map.contains`, that map to high, medium or low risk results. (see [below for nested schema](#nestedatt--predictor_custom_map--string_list))

Read-Only:

- `type` (String) A string that specifies the type of custom map predictor.

<a id="nestedatt--predictor_custom_map--between_ranges"></a>
### Nested Schema for `predictor_custom_map.between_ranges`

Optional:

- `high` (Attributes) A single nested object that describes the upper and lower bounds of ranges that map to a high risk result. (see [below for nested schema](#nestedatt--predictor_custom_map--between_ranges--high))
- `low` (Attributes) A single nested object that describes the upper and lower bounds of ranges that map to a low risk result. (see [below for nested schema](#nestedatt--predictor_custom_map--between_ranges--low))
- `medium` (Attributes) A single nested object that describes the upper and lower bounds of ranges that map to a medium risk result. (see [below for nested schema](#nestedatt--predictor_custom_map--between_ranges--medium))

<a id="nestedatt--predictor_custom_map--between_ranges--high"></a>
### Nested Schema for `predictor_custom_map.between_ranges.high`

Required:

- `max_value` (Number) A number that specifies the minimum value of the attribute named in `predictor_custom_map.contains`.  This represents the lower bound of this risk result range.
- `min_value` (Number) A number that specifies the minimum value of the attribute named in `predictor_custom_map.contains`.  This represents the lower bound of this risk result range.


<a id="nestedatt--predictor_custom_map--between_ranges--low"></a>
### Nested Schema for `predictor_custom_map.between_ranges.low`

Required:

- `max_value` (Number) A number that specifies the minimum value of the attribute named in `predictor_custom_map.contains`.  This represents the lower bound of this risk result range.
- `min_value` (Number) A number that specifies the minimum value of the attribute named in `predictor_custom_map.contains`.  This represents the lower bound of this risk result range.


<a id="nestedatt--predictor_custom_map--between_ranges--medium"></a>
### Nested Schema for `predictor_custom_map.between_ranges.medium`

Required:

- `max_value` (Number) A number that specifies the minimum value of the attribute named in `predictor_custom_map.contains`.  This represents the lower bound of this risk result range.
- `min_value` (Number) A number that specifies the minimum value of the attribute named in `predictor_custom_map.contains`.  This represents the lower bound of this risk result range.



<a id="nestedatt--predictor_custom_map--ip_ranges"></a>
### Nested Schema for `predictor_custom_map.ip_ranges`

Optional:

- `high` (Attributes) A single nested object that describes the IP CIDR ranges that map to a high risk result. (see [below for nested schema](#nestedatt--predictor_custom_map--ip_ranges--high))
- `low` (Attributes) A single nested object that describes the IP CIDR ranges that map to a low risk result. (see [below for nested schema](#nestedatt--predictor_custom_map--ip_ranges--low))
- `medium` (Attributes) A single nested object that describes the IP CIDR ranges that map to a medium risk result. (see [below for nested schema](#nestedatt--predictor_custom_map--ip_ranges--medium))

<a id="nestedatt--predictor_custom_map--ip_ranges--high"></a>
### Nested Schema for `predictor_custom_map.ip_ranges.high`

Optional:

- `values` (Set of String) A set of strings, in CIDR format, that describe the CIDR ranges that should evaluate against the value of the attribute named in `predictor_custom_map.contains` for this risk result.


<a id="nestedatt--predictor_custom_map--ip_ranges--low"></a>
### Nested Schema for `predictor_custom_map.ip_ranges.low`

Optional:

- `values` (Set of String) A set of strings, in CIDR format, that describe the CIDR ranges that should evaluate against the value of the attribute named in `predictor_custom_map.contains` for this risk result.


<a id="nestedatt--predictor_custom_map--ip_ranges--medium"></a>
### Nested Schema for `predictor_custom_map.ip_ranges.medium`

Optional:

- `values` (Set of String) A set of strings, in CIDR format, that describe the CIDR ranges that should evaluate against the value of the attribute named in `predictor_custom_map.contains` for this risk result.



<a id="nestedatt--predictor_custom_map--string_list"></a>
### Nested Schema for `predictor_custom_map.string_list`

Optional:

- `high` (Attributes) A single nested object that describes the string values that map to a high risk result. (see [below for nested schema](#nestedatt--predictor_custom_map--string_list--high))
- `low` (Attributes) A single nested object that describes the string values that map to a low risk result. (see [below for nested schema](#nestedatt--predictor_custom_map--string_list--low))
- `medium` (Attributes) A single nested object that describes the string values that map to a medium risk result. (see [below for nested schema](#nestedatt--predictor_custom_map--string_list--medium))

<a id="nestedatt--predictor_custom_map--string_list--high"></a>
### Nested Schema for `predictor_custom_map.string_list.high`

Optional:

- `values` (Set of String) A set of strings that should evaluate against the value of the attribute named in `predictor_custom_map.contains` for this risk result.


<a id="nestedatt--predictor_custom_map--string_list--low"></a>
### Nested Schema for `predictor_custom_map.string_list.low`

Optional:

- `values` (Set of String) A set of strings that should evaluate against the value of the attribute named in `predictor_custom_map.contains` for this risk result.


<a id="nestedatt--predictor_custom_map--string_list--medium"></a>
### Nested Schema for `predictor_custom_map.string_list.medium`

Optional:

- `values` (Set of String) A set of strings that should evaluate against the value of the attribute named in `predictor_custom_map.contains` for this risk result.




<a id="nestedatt--predictor_device"></a>
### Nested Schema for `predictor_device`

Optional:

- `activation_at` (String) A string that represents a date on which the learning process for the device predictor should be restarted.  Can only be configured where the `detect` parameter is `NEW_DEVICE`. This can be used in conjunction with the fallback setting (`default.result.level`) to force strong authentication when moving the predictor to production. The date should be in an RFC3339 format. Note that activation date uses UTC time.
- `detect` (String) A string that represents the type of device detection to use.  Options are `NEW_DEVICE` (to configure a model based on new devices), `SUSPICIOUS_DEVICE` (to configure a model based on detection of suspicious devices).  Defaults to `NEW_DEVICE`.
- `should_validate_payload_signature` (Boolean) Relevant only for Suspicious Device predictors. A boolean that, if set to `true`, then any risk policies that include this predictor will require that the Signals SDK payload be provided as a signed JWT whose signature will be verified before proceeding with risk evaluation. You instruct the Signals SDK to provide the payload as a signed JWT by using the `universalDeviceIdentification` flag during initialization of the SDK, or by selecting the relevant setting for the `skrisk` component in DaVinci flows.


<a id="nestedatt--predictor_email_reputation"></a>
### Nested Schema for `predictor_email_reputation`


<a id="nestedatt--predictor_geovelocity"></a>
### Nested Schema for `predictor_geovelocity`

Optional:

- `allowed_cidr_list` (Set of String) A set of IP addresses (CIDRs) that are ignored for the predictor results. The list can include IPs in IPv4 format and IPs in IPv6 format.


<a id="nestedatt--predictor_ip_reputation"></a>
### Nested Schema for `predictor_ip_reputation`

Optional:

- `allowed_cidr_list` (Set of String) A set of IP addresses (CIDRs) that are ignored for the predictor results. The list can include IPs in IPv4 format and IPs in IPv6 format.


<a id="nestedatt--predictor_traffic_anomaly"></a>
### Nested Schema for `predictor_traffic_anomaly`

Required:

- `rules` (Attributes List) A list collection of rules to use for this traffic anomaly predictor. (see [below for nested schema](#nestedatt--predictor_traffic_anomaly--rules))

<a id="nestedatt--predictor_traffic_anomaly--rules"></a>
### Nested Schema for `predictor_traffic_anomaly.rules`

Required:

- `enabled` (Boolean) A boolean to use the defined rule in the predictor.
- `interval` (Attributes) A single nested object that contains the fields used to define the timeframe to consider. The timeframe can be between 1 hour and 14 days. (see [below for nested schema](#nestedatt--predictor_traffic_anomaly--rules--interval))
- `threshold` (Attributes) A single nested object that contains the fields used to define the risk thresholds. (see [below for nested schema](#nestedatt--predictor_traffic_anomaly--rules--threshold))
- `type` (String) A string that specifies the type of velocity algorithm to use.  Options are `UNIQUE_USERS_PER_DEVICE`.

<a id="nestedatt--predictor_traffic_anomaly--rules--interval"></a>
### Nested Schema for `predictor_traffic_anomaly.rules.interval`

Required:

- `quantity` (Number) An integer that specifies the number of days or hours for the timeframe for tracking number of users on the device.
- `unit` (String) A string that specifies time unit for defining the timeframe for tracking number of users on the device.  Options are `DAY`, `HOUR`.


<a id="nestedatt--predictor_traffic_anomaly--rules--threshold"></a>
### Nested Schema for `predictor_traffic_anomaly.rules.threshold`

Required:

- `high` (Number) A float that specifies the number of users during the defined timeframe that will be considered High risk.
- `medium` (Number) A float that specifies the number of users during the defined timeframe that will be considered Medium risk.




<a id="nestedatt--predictor_user_location_anomaly"></a>
### Nested Schema for `predictor_user_location_anomaly`

Optional:

- `radius` (Attributes) A single nested object that specifies options for the radius to apply to the predictor evaluation (see [below for nested schema](#nestedatt--predictor_user_location_anomaly--radius))

Read-Only:

- `days` (Number) An integer that specifies the number of days to apply to the predictor evaluation.

<a id="nestedatt--predictor_user_location_anomaly--radius"></a>
### Nested Schema for `predictor_user_location_anomaly.radius`

Required:

- `distance` (Number) An integer that specifies the distance to apply to the predictor evaluation.

Optional:

- `unit` (String) A string that specifies the unit of distance to apply to the predictor distance.  Options are `kilometers`, `miles`.  Defaults to `kilometers`.



<a id="nestedatt--predictor_user_risk_behavior"></a>
### Nested Schema for `predictor_user_risk_behavior`

Required:

- `prediction_model` (Attributes) A single nested object that specifies options for the prediction model to apply to the predictor evaluation. (see [below for nested schema](#nestedatt--predictor_user_risk_behavior--prediction_model))

<a id="nestedatt--predictor_user_risk_behavior--prediction_model"></a>
### Nested Schema for `predictor_user_risk_behavior.prediction_model`

Required:

- `name` (String) A string that specifies the name of the prediction model to apply to the predictor evaluation.  Options are `login_anomaly_statistic` (to configure the organisation based risk model), `points` (to configure the user-based risk model).



<a id="nestedatt--predictor_velocity"></a>
### Nested Schema for `predictor_velocity`

Required:

- `of` (String) A string value that specifies the attribute reference for the value to aggregate when calculating velocity metrics.  Options are `${event.ip}` (to configure IP address velocity by user ID), `${event.user.id}` (to configure user velocity by IP address).  When defining attribute references in Terraform, the leading `$` needs to be escaped with an additional `$` character, e.g. `of = "$${event.ip}"`.

Optional:

- `measure` (String) A string value that specifies the type of measure to use for the predictor.  Options are `DISTINCT_COUNT`.  Defaults to `DISTINCT_COUNT`.

Read-Only:

- `by` (Set of String) A set of string values that specifies the attribute references that denote the subject of the velocity metric.  Options are `${event.ip}` (denotes the velocity metric is calculated by IP address), `${event.user.id}` (denotes the velocity metric is calculated by user ID).
- `every` (Attributes) A single nested object that specifies options for the granularlity of data sampling. (see [below for nested schema](#nestedatt--predictor_velocity--every))
- `fallback` (Attributes) A single nested object that specifies options for the predictor fallback strategy. (see [below for nested schema](#nestedatt--predictor_velocity--fallback))
- `sliding_window` (Attributes) A single nested object that specifies options for the distribution of data that is compared against to detect anomaly. (see [below for nested schema](#nestedatt--predictor_velocity--sliding_window))
- `use` (Attributes) A single nested object that specifies options for the velocity algorithm. (see [below for nested schema](#nestedatt--predictor_velocity--use))

<a id="nestedatt--predictor_velocity--every"></a>
### Nested Schema for `predictor_velocity.every`

Read-Only:

- `min_sample` (Number) An integer that denotes the minimum sample of data to use for the velocity algorithm.
- `quantity` (Number) An integer that denotes the quantity of unit intervals to use for the velocity algorithm.
- `unit` (String) A string value that specifies the time unit to use when sampling data.  Options are `DAY`, `HOUR`.


<a id="nestedatt--predictor_velocity--fallback"></a>
### Nested Schema for `predictor_velocity.fallback`

Read-Only:

- `high` (Number) A floating point value that specifies a high risk threshold for the fallback strategy.
- `medium` (Number) A floating point value that specifies a medium risk threshold for the fallback strategy.
- `strategy` (String) A string value that specifies the type of fallback strategy algorithm to use.  Options are `ENVIRONMENT_MAX`.


<a id="nestedatt--predictor_velocity--sliding_window"></a>
### Nested Schema for `predictor_velocity.sliding_window`

Read-Only:

- `min_sample` (Number) An integer that denotes the minimum sample of data to use for the velocity algorithm.
- `quantity` (Number) An integer that denotes the quantity of unit intervals to use for the velocity algorithm.
- `unit` (String) A string value that specifies the time unit to use when sampling data over time.  Options are `DAY`, `HOUR`.


<a id="nestedatt--predictor_velocity--use"></a>
### Nested Schema for `predictor_velocity.use`

Read-Only:

- `high` (Number) A floating point value that specifies a high risk threshold for the velocity algorithm.
- `medium` (Number) A floating point value that specifies a medium risk threshold for the velocity algorithm.
- `type` (String) A string value that specifies the type of velocity algorithm to use.  Options are `POISSON_WITH_MAX`.

## Import

Import is supported using the following syntax, where attributes in `<>` brackets are replaced with the relevant ID.  For example, `<environment_id>` should be replaced with the ID of the environment to import from.

```shell
terraform import pingone_risk_predictor.example <environment_id>/<risk_predictor_id>
```
