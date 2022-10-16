resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_image" "company_logo" {
  environment_id = pingone_environment.my_environment.id

  image_file_base64 = filebase64("../path/to/image.jpg")
}

resource "pingone_branding_settings" "branding" {
  environment_id = pingone_environment.my_environment.id

  company_name = "BXRetail"

  logo_image {
    id   = pingone_image.company_logo.id
    href = pingone_image.company_logo.uploaded_image[0].href
  }
}
