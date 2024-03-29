resource "pingone_schema_attribute" "my_awesome_regex_attribute" {
  environment_id = pingone_environment.my_environment.id

  name = "awesomeRegexValidatedAttribute"
  type = "STRING"

  regex_validation = {
    pattern      = "^[a-zA-Z0-9]*$",
    requirements = "Lowercase, uppercase and numeric."

    values_pattern_should_match = [
      "Up123",
      "Down456"
    ]

    values_pattern_should_not_match = [
      "Charm123!",
      "Strange456!"
    ]
  }

  lifecycle {
    # change the `prevent_destroy` parameter value to `true` to prevent this data carrying resource from being destroyed
    prevent_destroy = false
  }
}