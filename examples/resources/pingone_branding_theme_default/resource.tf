resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_branding_theme" "my_awesome_theme" {
  environment_id = pingone_environment.my_environment.id

  name     = "My Awesome Theme"
  template = "split"

  use_default_background = true

  button_text_color  = "#FFFFFF"
  heading_text_color = "#686F77"
  card_color         = "#FCFCFC"
  body_text_color    = "#263956"
  link_text_color    = "#263956"
  button_color       = "#263956"
}

resource "pingone_branding_theme_default" "my_awesome_theme_active" {
  environment_id = pingone_environment.my_environment.id

  branding_theme_id = pingone_branding_theme.my_awesome_theme.id
}