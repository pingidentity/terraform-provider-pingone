resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_image" "foo" {
  environment_id = pingone_environment.my_environment.id

  image_file_base64 = filebase64("../path/to/image.jpg")
}