resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_credential_type" "verifiedemployee" {
  environment_id       = pingone_environment.my_environment.id
  title                = "VerifiedEmployee"
  description          = "Demo Proof of Employment"
  card_type            = "VerifiedEmployee"
  card_design_template = "<svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 740 480\"><rect fill=\"none\" width=\"736\" height=\"476\" stroke=\"#CACED3\" stroke-width=\"3\" rx=\"10\" ry=\"10\" x=\"2\" y=\"2\"></rect><rect fill=\"$${cardColor}\" height=\"476\" rx=\"10\" ry=\"10\" width=\"736\" x=\"2\" y=\"2\" opacity=\"$${bgOpacityPercent}\"></rect><image href=\"$${backgroundImage}\" opacity=\"$${bgOpacityPercent}\" height=\"301\" rx=\"10\" ry=\"10\" width=\"589\" x=\"75\" y=\"160\"></image><image href=\"$${logoImage}\" x=\"42\" y=\"43\" height=\"90px\" width=\"90px\"></image><line y2=\"160\" x2=\"695\" y1=\"160\" x1=\"42.5\" stroke=\"$${textColor}\"></line><text fill=\"$${textColor}\" font-weight=\"450\" font-size=\"30\" x=\"160\" y=\"90\">$${cardTitle}</text><text fill=\"$${textColor}\" font-size=\"25\" font-weight=\"300\" x=\"160\" y=\"130\">$${cardSubtitle}</text></svg>"
  revoke_on_delete     = false

  metadata = {
    name               = "VerifiedEmployee"
    description        = "Demo Proof of Employment"
    bg_opacity_percent = 100

    background_image = "data:image/png;base64,${filebase64("./images/verifiedemployee_background.png")}"
    logo_image       = "data:image/png;base64,${filebase64("./images/verifiedemployee_logo.png")}"

    card_color = "#ffffff"
    text_color = "#000000"

    fields = [
      {
        type       = "Directory Attribute"
        title      = "givenName"
        attribute  = "name.given"
        is_visible = false
      },
      {
        type       = "Directory Attribute"
        title      = "surname"
        attribute  = "name.family"
        is_visible = false
      },
      {
        type       = "Directory Attribute"
        title      = "jobTitle"
        attribute  = "title"
        is_visible = false
      },
      {
        type       = "Directory Attribute"
        title      = "displayName"
        attribute  = "displayName"
        is_visible = false
      },
      {
        type       = "Directory Attribute"
        title      = "mail"
        attribute  = "email"
        is_visible = false
      },
      {
        type       = "Directory Attribute"
        title      = "preferredLanguage"
        attribute  = "preferredLanguage"
        is_visible = false
      },
      {
        type       = "Directory Attribute"
        title      = "id"
        attribute  = "id"
        is_visible = false
      }
    ]
  }

}