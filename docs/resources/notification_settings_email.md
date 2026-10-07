---
page_title: "pingone_notification_settings_email Resource - terraform-provider-pingone"
subcategory: "Platform"
description: |-
  Resource to manage the email sender settings in a PingOne environment.
---

# pingone_notification_settings_email (Resource)

Resource to manage the email sender settings in a PingOne environment.

~> Only one `pingone_notification_settings_email` resource should be configured for an environment.  If multiple `pingone_notification_settings_email` resource definitions exist in HCL code, these are likely to conflict with each other on apply. The `pingone_notification_settings` resource should be used if using the Ping-hosted SMTP email service.

## Example Usage - SMTP

```terraform
resource "pingone_environment" "my_environment" {
  # ...
}

# Example custom SMTP server configuration
resource "pingone_notification_settings_email" "my_awesome_smtp_settings" {
  environment_id = pingone_environment.my_environment.id

  host     = "smtp-example.bxretail.org"
  port     = 25
  username = var.smtp_server_username
  password = var.smtp_server_password

  from = {
    email_address = "services@bxretail.org"
    name          = "Customer Services"
  }
}
```

## Example Usage - Custom Provider

```terraform
resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_notification_settings_email" "my_awesome_custom_provider_settings" {
  environment_id       = pingone_environment.my_environment.id
  custom_provider_name = "My Custom Email Provider"

  username = var.custom_provider_username
  password = var.custom_provider_password
  protocol = "HTTP"

  from = {
    name          = "From Services"
    email_address = "noreply@bxretail.org"
  }

  reply_to = {
    name          = "Reply To Services"
    email_address = "services@bxretail.org"
  }

  requests = [
    {
      method = "POST"
      headers = {
        "Content-Type" = "application/x-www-form-urlencoded"
        "subject" : "$${subject}"
        "reply-to" : "$${reply_to}"
        "from" : "$${from}"
      }
      body = "to=$${to}&message=$${message}"
      url  = "https://api.bxretail.org/send-email"
    }
  ]
}
```

## Example Usage - Custom Provider with OAuth 2.0 Client Credentials

```terraform
resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_notification_settings_email" "my_awesome_custom_provider_settings" {
  environment_id       = pingone_environment.my_environment.id
  custom_provider_name = "My Custom Email Provider"

  auth_url      = var.custom_provider_auth_url
  grant_type    = "CLIENT_CREDENTIALS"
  client_id     = var.custom_provider_client_id
  client_secret = var.custom_provider_client_secret
  scopes        = ["mail.send"]
  protocol      = "HTTP"

  from = {
    name          = "From Services"
    email_address = "noreply@bxretail.org"
  }

  reply_to = {
    name          = "Reply To Services"
    email_address = "services@bxretail.org"
  }

  requests = [
    {
      method = "POST"
      headers = {
        "Content-Type" = "application/x-www-form-urlencoded"
        "subject" : "$${subject}"
        "reply-to" : "$${reply_to}"
        "from" : "$${from}"
      }
      body = "to=$${to}&message=$${message}"
      url  = "https://api.bxretail.org/send-email"
    }
  ]
}
```

## Example Usage - Custom Provider with OAuth 2.0 JWT Bearer

```terraform
resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_notification_settings_email" "my_awesome_custom_provider_settings" {
  environment_id       = pingone_environment.my_environment.id
  custom_provider_name = "My Custom Email Provider"

  auth_url   = var.custom_provider_auth_url
  grant_type = "JWT_BEARER"
  assertion  = var.custom_provider_jwt_assertion
  scopes     = ["mail.send"]
  protocol   = "HTTP"

  from = {
    name          = "From Services"
    email_address = "noreply@bxretail.org"
  }

  reply_to = {
    name          = "Reply To Services"
    email_address = "services@bxretail.org"
  }

  requests = [
    {
      method = "POST"
      headers = {
        "Content-Type" = "application/x-www-form-urlencoded"
        "subject" : "$${subject}"
        "reply-to" : "$${reply_to}"
        "from" : "$${from}"
      }
      body = "to=$${to}&message=$${message}"
      url  = "https://api.bxretail.org/send-email"
    }
  ]
}
```

## Example Usage - Custom Provider with Custom Header Authentication

```terraform
resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_notification_settings_email" "my_awesome_custom_provider_settings" {
  environment_id       = pingone_environment.my_environment.id
  custom_provider_name = "My Custom Email Provider"

  header_name  = "X-Api-Key"
  header_value = var.custom_provider_header_value
  protocol     = "HTTP"

  from = {
    name          = "From Services"
    email_address = "noreply@bxretail.org"
  }

  reply_to = {
    name          = "Reply To Services"
    email_address = "services@bxretail.org"
  }

  requests = [
    {
      method = "POST"
      headers = {
        "Content-Type" = "application/x-www-form-urlencoded"
        "subject" : "$${subject}"
        "reply-to" : "$${reply_to}"
        "from" : "$${from}"
      }
      body = "to=$${to}&message=$${message}"
      url  = "https://api.bxretail.org/send-email"
    }
  ]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `environment_id` (String) The ID of the environment to configure email settings in.  Must be a valid PingOne resource ID.  This field is immutable and will trigger a replace plan if changed.
- `from` (Attributes) A single block that specifies the email sender's "from" name and email address. (see [below for nested schema](#nestedatt--from))

### Optional

- `assertion` (String, Sensitive) A string that specifies the JWT assertion used to request the access token from the authorization server, when the Custom Provider is authenticated via OAuth 2.0.  Must be a valid JWT.  Required when `grant_type` is `JWT_BEARER`.
- `auth_token` (String, Sensitive) A string that specifies the authentication token when using a Custom Provider.
- `auth_url` (String) A string that specifies the URL of the authorization server that issues the access token for the email provider, when the Custom Provider is authenticated via OAuth 2.0.
- `client_id` (String) A string that specifies the client ID used to request the access token from the authorization server, when the Custom Provider is authenticated via OAuth 2.0.  Required when `grant_type` is `CLIENT_CREDENTIALS`.
- `client_secret` (String, Sensitive) A string that specifies the client secret used to request the access token from the authorization server, when the Custom Provider is authenticated via OAuth 2.0.  Required when `grant_type` is `CLIENT_CREDENTIALS`.
- `custom_provider_name` (String) A string to use to identify the provider.
- `grant_type` (String) A string that specifies the grant type used to request the access token from the authorization server, when the Custom Provider is authenticated via OAuth 2.0.  Options are `CLIENT_CREDENTIALS`, `JWT_BEARER`.  Defaults to `CLIENT_CREDENTIALS`.
- `header_name` (String) A string that specifies the name of the custom header used to authenticate requests to the email provider, when the Custom Provider is authenticated via a custom header.
- `header_value` (String, Sensitive) A string that specifies the value of the custom header used to authenticate requests to the email provider, when the Custom Provider is authenticated via a custom header.
- `host` (String) A string that specifies the organization's SMTP server.
- `password` (String, Sensitive) A string that specifies the organization's server's password.
- `port` (Number) An integer that specifies the port used by the organization's SMTP server to send emails (default: `465`). Note that the protocol used depends upon the port specified. If you specify port `25`, `587`, or `2525`, SMTP with `STARTTLS` is used. Otherwise, `SMTPS` is used.
- `protocol` (String) A string that specifies the current protocol in use.
- `reply_to` (Attributes) A single block that specifies the email sender's "reply to" name and email address. (see [below for nested schema](#nestedatt--reply_to))
- `requests` (Attributes Set) A list of objects that is used to configure the API requests sent to the custom email provider. (see [below for nested schema](#nestedatt--requests))
- `scopes` (Set of String) A set of strings that specifies the scopes to request in the access token from the authorization server (for example, `mail.send`), when the Custom Provider is authenticated via OAuth 2.0.
- `username` (String) A string that specifies the organization's server's username.

### Read-Only

- `client_authentication_method` (String) A string that specifies the method the service uses to send the OAuth 2.0 client credentials to the authorization server.  Returned by the service when the Custom Provider is authenticated via OAuth 2.0.
- `id` (String) The ID of this resource.
- `provider_type` (String) A string that specifies the provider type.

<a id="nestedatt--from"></a>
### Nested Schema for `from`

Required:

- `email_address` (String) A string that specifies the email sender's "from" email address.

Optional:

- `name` (String) A string that specifies the email sender's "from" name.


<a id="nestedatt--reply_to"></a>
### Nested Schema for `reply_to`

Required:

- `email_address` (String) A string that specifies the email sender's "reply to" email address.

Optional:

- `name` (String) A string that specifies the email sender's "reply to" name.


<a id="nestedatt--requests"></a>
### Nested Schema for `requests`

Required:

- `method` (String) Use method to specify the type of API request the email provider requires. Valid values are `GET` and `POST`.
- `url` (String) A string that specifies the endpoint for your email provider.

Optional:

- `body` (String) Required if method is set to `POST`. Use body to provide the content of the body for the request sent to the email provider.
- `headers` (Map of String) A map of key-value pairs to specify the headers that your email provider's API expects.

Read-Only:

- `delivery_method` (String) A string that specifies the delivery method for the request.

## Import

Import is supported using the following syntax, where attributes in `<>` brackets are replaced with the relevant ID.  For example, `<environment_id>` should be replaced with the ID of the environment to import from.

```shell
terraform import pingone_notification_settings_email.example <environment_id>
```
