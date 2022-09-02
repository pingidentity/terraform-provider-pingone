resource "pingone_sign_on_policy_action" "my_policy_progressive_profiling" {
  environment_id    = pingone_environment.my_environment.id
  sign_on_policy_id = pingone_sign_on_policy.my_policy.id

  priority = 3

  progressive_profiling {

    attribute {
      name     = "name.given"
      required = false
    }

    attribute {
      name     = "name.family"
      required = true
    }

    prompt_text = "For the best experience, we need a couple things from you."

  }
}