resource "pingone_schema_attribute" "my_awesome_regex_attribute" {
  environment_id = pingone_environment.my_environment.id
  schema_id      = data.pingone_schema.users.id

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
}