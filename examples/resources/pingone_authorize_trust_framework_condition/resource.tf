resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_authorize_trust_framework_condition" "my_awesome_condition" {
  environment_id = pingone_environment.my_environment.id
  name           = "Compare Values"
  description    = "My awesome compound condition."

  condition = {
    type = "OR"

    conditions = [
      {
        type       = "COMPARISON"
        comparator = "EQUALS"

        left = {
          type = "ATTRIBUTE"
          id   = pingone_authorize_trust_framework_attribute.my_awesome_attribute.id
        }

        right = {
          type  = "CONSTANT"
          value = "my_awesome_value"
        }
      },
      {
        type = "NOT"

        condition = {
          type       = "COMPARISON"
          comparator = "EQUALS"

          left = {
            type  = "ATTRIBUTE"
            value = pingone_authorize_trust_framework_attribute.my_awesome_attribute.id
          }

          right = {
            type  = "CONSTANT"
            value = "my_not_so_awesome_value"
          }
        }
      }
    ]
  }
}
