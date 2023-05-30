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

resource "pingone_risk_predictor" "my_awesome_geovelocity_anomaly_predictor" {
  # ...
}

resource "pingone_risk_policy" "my_awesome_scores_risk_policy" {
  environment_id = pingone_environment.my_environment.id

  name = "My Awesome Score-based Risk Policy"

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
        compact_name = pingone_risk_predictor.my_awesome_geovelocity_anomaly_predictor.compact_name
        score        = 50
      }
    ]
  }
}
```

## Example Usage - Weights based policy

```terraform
resource "pingone_risk_predictor" "my_awesome_anonymous_network_predictor" {
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
      min_score = 2
    }

    policy_threshold_high = {
      min_score = 5
    }

    predictors = [
      {
        compact_name = pingone_risk_predictor.my_awesome_anonymous_network_predictor.compact_name
        weight       = 5
      },
      {
        compact_name = pingone_risk_predictor.my_awesome_geovelocity_anomaly_predictor.compact_name
        weight       = 5
      }
    ]
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `environment_id` (String) The ID of the environment to configure the risk policy in.
- `name` (String) A string that specifies the unique, friendly name for this policy set. Valid characters consist of any Unicode letter, mark (such as, accent, umlaut), # (numeric), / (forward slash), . (period), ' (apostrophe), _ (underscore), space, or - (hyphen). Maximum size is 256 characters.

### Optional

- `default_result` (Attributes) A single nested object that specifies the default result value for the risk policy. (see [below for nested schema](#nestedatt--default_result))
- `policy_scores` (Attributes) An object that describes settings for a risk policy calculated by aggregating score values, with a final result being the sum of score values from each of the configured predictors. (see [below for nested schema](#nestedatt--policy_scores))
- `policy_weights` (Attributes) An object that describes settings for a risk policy using a weighted average calculation, with a final result being a risk score between `0` and `10`. (see [below for nested schema](#nestedatt--policy_weights))

### Read-Only

- `default` (Boolean) A boolean that indicates whether this risk policy set is the environment's default risk policy set. This is used whenever an explicit policy set ID is not specified in a risk evaluation request.
- `evaluated_predictors` (Set of String) A set of IDs for the predictors to evaluate in this policy set.  If omitted, if this property is null, all of the licensed predictors are used.
- `id` (String) The ID of this resource.

<a id="nestedatt--default_result"></a>
### Nested Schema for `default_result`

Optional:

- `level` (String) The default result level.  Options are `LOW`.

Read-Only:

- `type` (String) The default result type.  Options are `VALUE`.


<a id="nestedatt--policy_scores"></a>
### Nested Schema for `policy_scores`

Required:

- `predictors` (Attributes Set) An object that describes a predictor to apply to the risk policy and its associated high risk / true outcome score to apply to the risk calculation. (see [below for nested schema](#nestedatt--policy_scores--predictors))

Optional:

- `policy_threshold_high` (Attributes) An object that specifies the lower and upper bound threshold values that define the high risk outcome as a result of the policy evaluation. (see [below for nested schema](#nestedatt--policy_scores--policy_threshold_high))
- `policy_threshold_medium` (Attributes) An object that specifies the lower and upper bound threshold values that define the medium risk outcome as a result of the policy evaluation. (see [below for nested schema](#nestedatt--policy_scores--policy_threshold_medium))

<a id="nestedatt--policy_scores--predictors"></a>
### Nested Schema for `policy_scores.predictors`

Required:

- `compact_name` (String) A string that specifies the compact name of the predictor to apply to the risk policy.
- `score` (Number) An integer that specifies the score to apply to the High risk / true outcome of the predictor, to apply to the overall risk calculation.

Read-Only:

- `predictor_reference_value` (String) A string that specifies the attribute reference of the level to evaluate.


<a id="nestedatt--policy_scores--policy_threshold_high"></a>
### Nested Schema for `policy_scores.policy_threshold_high`

Optional:

- `min_score` (Number) An integer that specifies the minimum score to use as the lower bound value of the policy threshold.  Defaults to `75`.

Read-Only:

- `max_score` (Number) An integer that specifies the maxiumum score to use as the lower bound value of the policy threshold.


<a id="nestedatt--policy_scores--policy_threshold_medium"></a>
### Nested Schema for `policy_scores.policy_threshold_medium`

Optional:

- `min_score` (Number) An integer that specifies the minimum score to use as the lower bound value of the policy threshold.  Defaults to `40`.

Read-Only:

- `max_score` (Number) An integer that specifies the maxiumum score to use as the lower bound value of the policy threshold.



<a id="nestedatt--policy_weights"></a>
### Nested Schema for `policy_weights`

Required:

- `predictors` (Attributes Set) An object that describes a predictor to apply to the risk policy and its associated weight value for the overall weighted average risk calculation. (see [below for nested schema](#nestedatt--policy_weights--predictors))

Optional:

- `policy_threshold_high` (Attributes) An object that specifies the lower and upper bound threshold values that define the high risk outcome as a result of the policy evaluation. (see [below for nested schema](#nestedatt--policy_weights--policy_threshold_high))
- `policy_threshold_medium` (Attributes) An object that specifies the lower and upper bound threshold values that define the medium risk outcome as a result of the policy evaluation. (see [below for nested schema](#nestedatt--policy_weights--policy_threshold_medium))

<a id="nestedatt--policy_weights--predictors"></a>
### Nested Schema for `policy_weights.predictors`

Required:

- `compact_name` (String) A string that specifies the compact name of the predictor to apply to the risk policy.
- `weight` (Number) An integer that specifies the weight to apply to the predictor when calculating the overall risk score.

Read-Only:

- `predictor_reference_value` (String) A string that specifies the attribute reference of the level to evaluate.


<a id="nestedatt--policy_weights--policy_threshold_high"></a>
### Nested Schema for `policy_weights.policy_threshold_high`

Optional:

- `min_score` (Number) An integer that specifies the minimum score to use as the lower bound value of the policy threshold.  Defaults to `2`.

Read-Only:

- `max_score` (Number) An integer that specifies the maxiumum score to use as the lower bound value of the policy threshold.


<a id="nestedatt--policy_weights--policy_threshold_medium"></a>
### Nested Schema for `policy_weights.policy_threshold_medium`

Optional:

- `min_score` (Number) An integer that specifies the minimum score to use as the lower bound value of the policy threshold.  Defaults to `1`.

Read-Only:

- `max_score` (Number) An integer that specifies the maxiumum score to use as the lower bound value of the policy threshold.

## Import

Import is supported using the following syntax, where attributes in `<>` brackets are replaced with the relevant ID.  For example, `<environment_id>` should be replaced with the ID of the environment to import from.

```shell
$ terraform import pingone_risk_policy.example <environment_id>/<risk_policy_id>
```