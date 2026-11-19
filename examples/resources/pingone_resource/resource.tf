resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_resource" "my_resource" {
  environment_id = pingone_environment.my_environment.id

  name        = "My resource"
  description = "My new Resource"

  audience                      = "https://api.myresource.com"
  access_token_validity_seconds = 3600

  introspect_endpoint_auth_method = "PRIVATE_KEY_JWT"

  jwks = jsonencode({
    keys = [
      {
        kty = "RSA"
        kid = "my-resource-signing-key"
        alg = "RS256"
        n   = "0vx7agoebGcQSuuPiLJXZptN9nndrQmbXEps2aiAFbWhM78LhWx4cbbfAAtb0CmzzzyYV8dBnrKamTmULCcBLine1V2d70uH7jXGNKLBaaNjZjxBUxQXKFCpbmHBcjW1z5IbjgkG9ESj8STbLIimSbvW7wPE_Td0dzRRPUHMYu7YmC3RLWkv12M7O4MTFGM9vLhUdAD0msh0uDdrEBFQ0dO7YbXCHJ1FFG6gJXf9cmU3g8PLAETWdanT6wT36B0-HmF4FCYUktLXBGikmLQyqxCYBNfPTGqbUc0ETzAToZg2yx3TFW1hQxTVBwOq0ephT6who7IW38SyLv0WaBSG27Z9c9TH0mp2OqSjIiFfr7XcvKX8Q"
        e   = "AQAB"
      },
    ]
  })
}
