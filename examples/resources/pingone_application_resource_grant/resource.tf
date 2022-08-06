resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_application" "my_awesome_spa" {
	environment_id  = pingone_environment.my_environment.id
	name            = "My Awesome Single Page App"

	oidc_options {
		type                        = "SINGLE_PAGE_APP"
		grant_types                 = [ "AUTHORIZATION_CODE" ]
    response_types              = [ "CODE" ]
    pkce_enforcement            = "S256_REQUIRED"
	  token_endpoint_authn_method = "NONE"
    redirect_uris               = [ "https://my-website.com" ]
  }
}

data "pingone_resource_scope" "openid_profile" {
	environment_id = pingone_environment.my_environment.id
  resource_id    = data.pingone_resource.openid_resource.id

  name = "profile"
}

data "pingone_resource_scope" "openid_email" {
	environment_id = pingone_environment.my_environment.id
  resource_id    = data.pingone_resource.openid_resource.id

  name = "email"

}

resource "pingone_application_resource_grant" "foo" {
	environment_id = pingone_environment.my_environment.id
	application_id = pingone_application.my_awesome_spa

  resource_id = data.pingone_resource.openid_resource.id

  scopes = [
    data.pingone_resource_scope.openid_profile.id,
    data.pingone_resource_scope.openid_email.id
  ]
}
