resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_form" "my_awesome_form" {
  environment_id = pingone_environment.my_environment.id

  name        = "My Awesome Form"
  description = "This is my awesome form"

  mark_required = false
  mark_optional = true

  cols = 4

  language_bundle = {
    "button.text"                              = "Submit",
    "fields.user.email.label"                  = "Email Address",
    "fields.user.password.label"               = "Password"
    "fields.user.password.labelPasswordVerify" = "Verify Password",
    "fields.user.username.label"               = "Username",
  }

  components = {
    fields = [
      {
        position = {
          row = 0,
          col = 0,

          field_slate_textblob = {
            content = jsonencode(
              [
                {
                  "children" : [
                    {
                      "text" : "My awesome form!"
                    }
                  ],
                  "type" : "heading-1"
                },
                {
                  "children" : [
                    {
                      "text" : "Enter the required information below"
                    }
                  ]
                },
                {
                  "type" : "divider",
                  "children" : [
                    {
                      "text" : ""
                    }
                  ]
                },
                {
                  "type" : "paragraph",
                  "children" : [
                    {
                      "text" : ""
                    }
                  ]
                }
            ])
          }
        }
      },
      {
        position = {
          row = 1,
          col = 0,

          field_error_display = {}
        }
      }
    ]
  }
}
