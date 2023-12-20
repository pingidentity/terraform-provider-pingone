resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_form" "my_awesome_form" {
  environment_id = pingone_environment.my_environment.id

  name        = "My Awesome Sign On Form"
  description = "This is my awesome form with fields for sign on"

  mark_required = false
  mark_optional = true

  cols = 4

  components = {
    fields = [
      {
        type = "TEXTBLOB"

        position = {
          row = 0
          col = 0
        }

        content = "<h2>Sign On</h2><hr>"
      },
      {
        type = "ERROR_DISPLAY"

        position = {
          row = 1
          col = 0
        }
      },
      {
        type = "TEXT"

        position = {
          row = 2
          col = 0
        }

        key = "user.username"
        label = jsonencode(
          [
            {
              "type" = "paragraph",
              "children" = [
                {
                  "text" = ""
                },
                {
                  "type"               = "i18n",
                  "key"                = "fields.user.username.label",
                  "defaultTranslation" = "Username",
                  "inline"             = true,
                  "children" = [
                    {
                      "text" = ""
                    }
                  ]
                },
                {
                  "text" = ""
                }
              ]
            }
          ]
        )

        required = true

        validation = {
          type = "NONE"
        }
      },
      {
        type = "PASSWORD"

        position = {
          row = 3
          col = 0
        }

        key = "user.password"
        label = jsonencode(
          [
            {
              "type" = "paragraph",
              "children" = [
                {
                  "text" = ""
                },
                {
                  "type"               = "i18n",
                  "key"                = "fields.user.password.label",
                  "defaultTranslation" = "Password",
                  "inline"             = true,
                  "children" = [
                    {
                      "text" = ""
                    }
                  ]
                },
                {
                  "text" = ""
                }
              ]
            }
          ]
        )

        required = true
      },
      {
        type = "SUBMIT_BUTTON"

        position = {
          row = 4
          col = 0
        }

        label = jsonencode(
          [
            {
              "type" = "paragraph",
              "children" = [
                {
                  "text" = ""
                },
                {
                  "type"               = "i18n",
                  "key"                = "button.text.signOn",
                  "defaultTranslation" = "Sign On",
                  "inline"             = true,
                  "children" = [
                    {
                      "text" = ""
                    }
                  ]
                },
                {
                  "text" = ""
                }
              ]
            }
          ]
        )
      }
    ]
  }
}
