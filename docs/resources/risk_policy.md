---
page_title: "pingone_risk_policy Resource - terraform-provider-pingone"
subcategory: "Risk"
description: |-
  Resource to manage Risk policies in a PingOne environment.
---

# pingone_risk_policy (Resource)

Resource to manage Risk policies in a PingOne environment.

## Example Usage - Scores based policy

```terraform
resource "pingone_risk_predictor" "my_awesome_anonymous_network_predictor" {
  # ...
}

resource "pingone_risk_predictor" "my_awesome_user_location_predictor" {
  # ...
}

resource "pingone_risk_predictor" "my_awesome_geovelocity_anomaly_predictor" {
  # ...
}

resource "pingone_risk_policy" "my_awesome_scores_risk_policy" {
  environment_id = pingone_environment.my_environment.id

  name = "My Awesome Scores-based Risk Policy"

  policy_scores = {
    policy_threshold_medium = {
      min_score = 40
    }

    policy_threshold_high = {
      min_score = 75
    }

    predictors = [
      {
        compact_name = pingone_risk_predictor.my_awesome_anonymous_network_predictor.compact_name
        score        = 50
      },
      {
        compact_name = pingone_risk_predictor.my_awesome_user_location_predictor.compact_name
        score        = 50
      }
    ]
  }

  overrides = [
    {
      result = {
        level = "LOW"
      }

      condition = {
        type = "IP_RANGE"
        ip_range = [
          "10.0.0.0/8",
        ]
      }
    },

    {
      result = {
        level = "HIGH"
      }

      condition = {
        type         = "VALUE_COMPARISON"
        compact_name = pingone_risk_predictor.my_awesome_geovelocity_anomaly_predictor.compact_name
        equals       = "HIGH"
      }
    }
  ]
}
```

## Example Usage - Weights based policy

```terraform
resource "pingone_risk_predictor" "my_awesome_anonymous_network_predictor" {
  # ...
}

resource "pingone_risk_predictor" "my_awesome_user_location_predictor" {
  # ...
}

resource "pingone_risk_predictor" "my_awesome_geovelocity_anomaly_predictor" {
  # ...
}

resource "pingone_risk_policy" "my_awesome_weights_risk_policy" {
  environment_id = pingone_environment.my_environment.id

  name = "My Awesome Weights-based Risk Policy"

  policy_weights = {
    policy_threshold_medium = {
      min_score = 50
    }

    policy_threshold_high = {
      min_score = 60
    }

    predictors = [
      {
        compact_name = pingone_risk_predictor.my_awesome_anonymous_network_predictor.compact_name
        weight       = 5
      },
      {
        compact_name = pingone_risk_predictor.my_awesome_user_location_predictor.compact_name
        weight       = 5
      }
    ]
  }

  overrides = [
    {
      result = {
        level = "LOW"
      }

      condition = {
        type = "IP_RANGE"
        ip_range = [
          "10.0.0.0/8",
        ]
      }
    },

    {
      result = {
        level = "HIGH"
      }

      condition = {
        type         = "VALUE_COMPARISON"
        compact_name = pingone_risk_predictor.my_awesome_geovelocity_anomaly_predictor.compact_name
        equals       = "HIGH"
      }
    }
  ]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `environment_id` (String) The ID of the environment to configure the risk policy in.
- `name` (String) A string that specifies the unique, friendly name for this policy set. Valid characters consist of any Unicode letter, mark (such as, accent, umlaut), # (numeric), / (forward slash), . (period), ' (apostrophe), _ (underscore), space, or - (hyphen). Maximum size is 256 characters.

### Optional

- `default_result` (Attributes) A single nested object that specifies the default result value for the risk policy. (see [below for nested schema](#nestedatt--default_result))
- `evaluated_predictors` (Set of String) A set of IDs for the predictors to evaluate in this policy set.  If omitted, if this property is null, all of the licensed predictors are used.
- `overrides` (Attributes List) An ordered list of policy overrides to apply to the policy.  The ordering of the overrides is important as it determines the priority of the policy override during policy evaluation. (see [below for nested schema](#nestedatt--overrides))
- `policy_scores` (Attributes) An object that describes settings for a risk policy calculated by aggregating score values, with a final result being the sum of score values from each of the configured predictors.  At least one of the following must be defined: `policy_weights`, `policy_scores`. (see [below for nested schema](#nestedatt--policy_scores))
- `policy_weights` (Attributes) An object that describes settings for a risk policy using a weighted average calculation, with a final result being a risk score between `0` and `10`.  At least one of the following must be defined: `policy_weights`, `policy_scores`. (see [below for nested schema](#nestedatt--policy_weights))

### Read-Only

- `default` (Boolean) A boolean that indicates whether this risk policy set is the environment's default risk policy set. This is used whenever an explicit policy set ID is not specified in a risk evaluation request.
- `id` (String) The ID of this resource.

<a id="nestedatt--default_result"></a>
### Nested Schema for `default_result`

Required:

- `level` (String) The default result level.  Options are `LOW`.

Read-Only:

- `type` (String) The default result type.  Options are `VALUE`.


<a id="nestedatt--overrides"></a>
### Nested Schema for `overrides`

Required:

- `condition` (Attributes) A single object that contains the conditions to evaluate that determine whether the override result will be applied to the risk policy evaluation. (see [below for nested schema](#nestedatt--overrides--condition))
- `result` (Attributes) A single object that contains the risk result that should be applied to the policy evaluation result when the override condition is met. (see [below for nested schema](#nestedatt--overrides--result))

Optional:

- `name` (String) A string that represents the name of the overriding risk policy in the set.

Read-Only:

- `priority` (Number) An integer that indicates the order in which the override is applied during risk policy evaluation.  The lower the value, the higher the priority.  The priority is determined by the order in which the overrides are defined in HCL.

<a id="nestedatt--overrides--condition"></a>
### Nested Schema for `overrides.condition`

Required:

- `type` (String) A string that specifies the type of the override condition to evaluate.  Options are `IP_RANGE`, `VALUE_COMPARISON`.

Optional:

- `compact_name` (String) Required when `equals` is set to `VALUE_COMPARISON`.  A string that specifies the compact name of the predictor to apply to the override condition.
- `equals` (String) Required when `equals` is set to `VALUE_COMPARISON`.  A string that specifies the value of the `predictor_reference_value` that must be matched for the override result to be applied to the policy evaluation.
- `ip_range` (Set of String) Required when `equals` is set to `IP_RANGE`.  A set of strings that specifies the CIDR ranges that should be evaluated against the value of the `predictor_reference_contains` attribute, that must be matched for the override result to be applied to the policy evaluation.  Values must be valid IPv4 or IPv6 CIDR ranges.

Read-Only:

- `predictor_reference_contains` (String) A string that specifies the attribute reference of the collection to evaluate.
- `predictor_reference_value` (String) A string that specifies the attribute reference of the value to evaluate.


<a id="nestedatt--overrides--result"></a>
### Nested Schema for `overrides.result`

Required:

- `level` (String) A string that specifies the risk level that should be applied to the policy evalution result when the override condition is met.  Options are `HIGH`, `LOW`, `MEDIUM`.

Optional:

- `type` (String) A string that specifies the type of the risk result should be applied to the policy evalution result when the override condition is met.  Options are `VALUE`.  Defaults to `VALUE`.
- `value` (String) An administrator defined string value that is applied to the policy evaluation result when the override condition is met.



<a id="nestedatt--policy_scores"></a>
### Nested Schema for `policy_scores`

Required:

- `policy_threshold_high` (Attributes) An object that specifies the lower and upper bound threshold values that define the high risk outcome as a result of the policy evaluation. (see [below for nested schema](#nestedatt--policy_scores--policy_threshold_high))
- `policy_threshold_medium` (Attributes) An object that specifies the lower and upper bound threshold values that define the medium risk outcome as a result of the policy evaluation. (see [below for nested schema](#nestedatt--policy_scores--policy_threshold_medium))
- `predictors` (Attributes Set) An object that describes a predictor to apply to the risk policy and its associated high risk / true outcome score to apply to the risk calculation. (see [below for nested schema](#nestedatt--policy_scores--predictors))

<a id="nestedatt--policy_scores--policy_threshold_high"></a>
### Nested Schema for `policy_scores.policy_threshold_high`

Required:

- `min_score` (Number) An integer that specifies the minimum score to use as the lower bound value of the policy threshold.  Maximum value allowed is `1000`

Read-Only:

- `max_score` (Number) An integer that specifies the maxiumum score to use as the lower bound value of the policy threshold.


<a id="nestedatt--policy_scores--policy_threshold_medium"></a>
### Nested Schema for `policy_scores.policy_threshold_medium`

Required:

- `min_score` (Number) An integer that specifies the minimum score to use as the lower bound value of the policy threshold.  Maximum value allowed is `1000`

Read-Only:

- `max_score` (Number) An integer that specifies the maxiumum score to use as the lower bound value of the policy threshold.


<a id="nestedatt--policy_scores--predictors"></a>
### Nested Schema for `policy_scores.predictors`

Required:

- `compact_name` (String) A string that specifies the compact name of the predictor to apply to the risk policy.
- `score` (Number) An integer that specifies the score to apply to the High risk / true outcome of the predictor, to apply to the overall risk calculation.

Read-Only:

- `predictor_reference_value` (String) A string that specifies the attribute reference of the level to evaluate.



<a id="nestedatt--policy_weights"></a>
### Nested Schema for `policy_weights`

Required:

- `policy_threshold_high` (Attributes) An object that specifies the lower and upper bound threshold score values that define the high risk outcome as a result of the policy evaluation. (see [below for nested schema](#nestedatt--policy_weights--policy_threshold_high))
- `policy_threshold_medium` (Attributes) An object that specifies the lower and upper bound threshold score values that define the medium risk outcome as a result of the policy evaluation. (see [below for nested schema](#nestedatt--policy_weights--policy_threshold_medium))
- `predictors` (Attributes Set) An object that describes a predictor to apply to the risk policy and its associated weight value for the overall weighted average risk calculation. (see [below for nested schema](#nestedatt--policy_weights--predictors))

<a id="nestedatt--policy_weights--policy_threshold_high"></a>
### Nested Schema for `policy_weights.policy_threshold_high`

Required:

- `min_score` (Number) An integer that specifies the minimum score to use as the lower bound value of the policy threshold.  For weights policies, the score values should be 10x the desired risk value in the console. For example, a risk score of `5` in the console should be entered as `50`.  The provided score must be exactly divisible by 10.  Maximum value allowed is `100`

Read-Only:

- `max_score` (Number) An integer that specifies the maxiumum score to use as the lower bound value of the policy threshold.


<a id="nestedatt--policy_weights--policy_threshold_medium"></a>
### Nested Schema for `policy_weights.policy_threshold_medium`

Required:

- `min_score` (Number) An integer that specifies the minimum score to use as the lower bound value of the policy threshold.  For weights policies, the score values should be 10x the desired risk value in the console. For example, a risk score of `5` in the console should be entered as `50`.  The provided score must be exactly divisible by 10.  Maximum value allowed is `100`

Read-Only:

- `max_score` (Number) An integer that specifies the maxiumum score to use as the lower bound value of the policy threshold.


<a id="nestedatt--policy_weights--predictors"></a>
### Nested Schema for `policy_weights.predictors`

Required:

- `compact_name` (String) A string that specifies the compact name of the predictor to apply to the risk policy.
- `weight` (Number) An integer that specifies the weight to apply to the predictor when calculating the overall risk score.

Read-Only:

- `predictor_reference_value` (String) A string that specifies the attribute reference of the level to evaluate.

## Import

Import is supported using the following syntax, where attributes in `<>` brackets are replaced with the relevant ID.  For example, `<environment_id>` should be replaced with the ID of the environment to import from.

```shell
$ terraform import pingone_risk_policy.example <environment_id>/<risk_policy_id>
```