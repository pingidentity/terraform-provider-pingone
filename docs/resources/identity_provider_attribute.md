---
page_title: "pingone_identity_provider_attribute Resource - terraform-provider-pingone"
subcategory: "SSO"
description: |-
  Resource to create and manage an attribute mapping for identity providers configured in PingOne.
---

# pingone_identity_provider_attribute (Resource)

Resource to create and manage an attribute mapping for identity providers configured in PingOne.

## Identity Provider Attribute Reference

PingOne supports several external IdPs. IdP resources in PingOne configure the external IdP settings, which include the type of provider and the user attributes from the external IdP that are mapped to PingOne user attributes. These attributes might have one or many values assigned to them. As you might expect, mapping a single-value IdP attribute to a single-value PingOne attribute results in a PingOne attribute having the same value as the IdP attribute. Similarly, if the IdP attribute is also multi-valued, the PingOne attribute value will be an array of the IdP attribute values. If the attributes are not the same format, then the following rules apply:

* If the IdP attribute is single-value and the PingOne attribute is multi-valued, then the PingOne attribute will be a single-element array containing the value of the IdP attribute.
* If the IdP attribute is multi-valued and the PingOne attribute is single-value, then the PingOne attribute will use the first element in the IdP attribute as its value.

The mapping attribute placeholder value must be expressed using the following syntax in the request body in the platform:

`${providerAttributes.<IdP attribute name>}`

Terraform HCL expects the attribute placeholder (used in the `value` argument of this `pingone_identity_provider_attribute` resource) to be prefixed with an additional `$` (dollar) sign:

```
...
  value = "$${providerAttributes.<IdP attribute name>}"
...
```

The following are IdP attributes expected per identity provider:

### Amazon
#### Core attributes
| Property  | Description                                                                                                                                           |
|-----------|-------------------------------------------------------------------------------------------------------------------------------------------------------|
| `user_id` | A string that specifies the core Amazon attribute. The default value is `${providerAttributes.user_id}` and the default update value is `EMPTY_ONLY`. |

#### Provider attributes
| Permission    | Provider attributes                     |
|---------------|-----------------------------------------|
| `profile`     | Options are: `user_id`, `email`, `name` |
| `postal_code` | Options are: `postal_code`              |

### Apple
#### Core attributes
| Property | Description                                                                                                                                      |
|----------|--------------------------------------------------------------------------------------------------------------------------------------------------|
| `sub`    | A string that specifies the core Apple attribute. The default value is `${providerAttributes.sub}` and the default update value is `EMPTY_ONLY`. |

#### Provider attributes
| Permission | Provider attributes                                                         |
|------------|-----------------------------------------------------------------------------|
| `name`     | Options are: `sub`, `iss`, `iat`, `expt`, `aud`, `nonce`, `nonce_supported` |
| `email`    | Options are: `email`, `email_verified`                                      |


### Facebook
#### Core attributes
| Property   | Description                                                                                                                                           |
|------------|-------------------------------------------------------------------------------------------------------------------------------------------------------|
| `username` | A string that specifies the core Facebook attribute. The default value is `${providerAttributes.email}` and the default update value is `EMPTY_ONLY`. |

#### Provider attributes
| Permission       | Provider attributes                                                                              |
|------------------|--------------------------------------------------------------------------------------------------|
| `<default>`      | Options are: `id`, `first_name`, `last_name`, `middle_name`, `name`, `name_format`, and `email`. |
| `USER_AGE_RANGE` | Options are: `age_range`.                                                                        |
| `USER_BIRTHDAY`  | Options are: `birthday`.                                                                         |
| `USER_GENDER`    | Options are: `gender`.                                                                           |

### Github
#### Core attributes
| Property | Description                                                                                                                                      |
|----------|--------------------------------------------------------------------------------------------------------------------------------------------------|
| `id`     | A string that specifies the core Github attribute. The default value is `${providerAttributes.id}` and the default update value is `EMPTY_ONLY`. |

#### Provider attributes
| Permission   | Provider attributes                                                                                                                          |
|--------------|----------------------------------------------------------------------------------------------------------------------------------------------|
| `read:user`  | Options are: `email`, `login`, `id`, `node_id`, `avatar_url`, `url`, `html_url`, `type`, `site_admin`, `name`, `company`, `blog`, `location` |
| `user:email` | Options are: `email`.                                                                                                                        |

### Google
#### Core attributes
| Property   | Description                                                                                                                                                      |
|------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `username` | A string that specifies the core Google attribute. The default value is `${providerAttributes.emailAddress.value}` and the default update value is `EMPTY_ONLY`. |

#### Provider attributes
| Permission                                               | Provider attributes                                                                                                                                                                                                    |
|----------------------------------------------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `profile, email`                                         | Options are: `resourceName`, `etag`, `emailAddress.value`, `name.displayName`, `name.familyName`, `name.givenName`, `name.middleName`, `nickname.value`, `nickname.type`, `gender.value`, and `gender.formattedValue`. |
| `https://www.googleapis.com/auth/profile.agerange.read`  | Options are: `ageRange.ageRange`.                                                                                                                                                                                      |
| `https://www.googleapis.com/auth/profile.language.read`  | Options are: `locale.value`.                                                                                                                                                                                           |
| `https://www.googleapis.com/auth/user.birthday.read`     | Options are: `birthday.date.month`, `birthday.date.day`, `birthday.date.year`, and `birthday.text`.                                                                                                                    |
| `https://www.googleapis.com/auth/user.phonenumbers.read` | Options are: `phoneNumber.value`.                                                                                                                                                                                      |

### LinkedIn
#### Core attributes
| Property   | Description                                                                                                                                                  |
|------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `username` | A string that specifies the core LinkedIn attribute. The default value is `${providerAttributes.emailAddress}` and the default update value is `EMPTY_ONLY`. |

#### Provider attributes
| Permission       | Provider attributes                         |
|------------------|---------------------------------------------|
| `r_liteprofile`  | Options are: `id`, `firstName`, `lastName`. |
| `r_emailaddress` | Options are: `emailAddress`.                |

### Microsoft
#### Core attributes
| Property | Description                                                                                                                                         |
|----------|-----------------------------------------------------------------------------------------------------------------------------------------------------|
| `id`     | A string that specifies the core Microsoft attribute. The default value is `${providerAttributes.id}` and the default update value is `EMPTY_ONLY`. |

#### Provider attributes
| Permission                               | Provider attributes                                                                                                                                                 |
|------------------------------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| OpenID Connect scopes: `openid`, `email` | `email`                                                                                                                                                             |
| `User:Read`                              | Options are: `displayName`, `surname`, `givenName`, `id`, `userPrincipalName`, `businessPhones`, `jobTitle`, `mail`, `officeLocation`, `postalCode`, `mainNickname` |

### OpenID Connect (Generic)
#### Core attributes
| Property   | Description                                                                                                                                               |
|------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------|
| `username` | A string that specifies the core OpenID Connect attribute. The default value is `${providerAttributes.sub}` and the default update value is `EMPTY_ONLY`. |

#### Provider attributes
| Permission | Provider attributes |
|------------|---------------------|
| `openid`   | `sub`               |

### Paypal
#### Core attributes
| Property  | Description                                                                                                                                           |
|-----------|-------------------------------------------------------------------------------------------------------------------------------------------------------|
| `user_id` | A string that specifies the core PayPal attribute. The default value is `${providerAttributes.user_id}` and the default update value is `EMPTY_ONLY`. |

#### Provider attributes
| Permission                                          | Provider attributes                                                                                                   |
|-----------------------------------------------------|-----------------------------------------------------------------------------------------------------------------------|
| OpenID Connect scopes: `openid`, `profile`, `email` | Options are: `user_id`, `name`, `email`                                                                               |
| `address`                                           | Options are: `address.street_address`, `address.locality`, `address.region`, `address.postal_code`, `address.country` |
| `paypalattributes`                                  | Options are: `payer_id`, `verified_account`                                                                           |	

### SAML (Generic)
#### Core attributes
| Property   | Description                                                                                                                                    |
|------------|------------------------------------------------------------------------------------------------------------------------------------------------|
| `username` | A string that specifies the core SAML attribute. The default value is `${samlAssertion.subject}` and the default update value is `EMPTY_ONLY`. |

### Twitter
#### Core attributes
| Property | Description                                                                                                                                       |
|----------|---------------------------------------------------------------------------------------------------------------------------------------------------|
| `id`     | A string that specifies the core Twitter attribute. The default value is `${providerAttributes.id}` and the default update value is `EMPTY_ONLY`. |

#### Provider attributes
| Permission                | Provider attributes                                                                                                                                                                                                             |
|---------------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `{no defined permission}` | Options are: `id`, `email`, `name`, `screen_name`, `created_at`, `statuses_count`, `favourites_count`, `friends_count`, `followers_count`, `verified`, `protected`, `description`, `url`, `location`, `profile_image_url_https` |

### Yahoo
#### Core attributes
| Property   | Description                                                                                                                                      |
|------------|--------------------------------------------------------------------------------------------------------------------------------------------------|
| `sub`      | A string that specifies the core Yahoo attribute. The default value is `${providerAttributes.sub}` and the default update value is `EMPTY_ONLY`. |

#### Provider attributes
| Permission | Provider attributes                                                               |
|------------|-----------------------------------------------------------------------------------|
| `openid`   | `sub`                                                                             |
| `email`    | `email`                                                                           |
| `profile`  | Options are: `name`, `given_name`, `family_name`, `picture`, `nickname`, `locale` |

## Example Usage

```terraform
resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_identity_provider" "apple" {
  environment_id = pingone_environment.my_environment.id

  name    = "Apple"
  enabled = true

  apple {
    client_id                 = var.apple_client_id
    client_secret_signing_key = var.apple_client_secret_signing_key
    key_id                    = var.apple_key_id
    team_id                   = var.apple_team_id
  }
}

resource "pingone_identity_provider_attribute" "apple_email" {
  environment_id       = pingone_environment.my_environment.id
  identity_provider_id = pingone_identity_provider.apple.id

  name   = "email"
  update = "EMPTY_ONLY"
  value  = "$${providerAttributes.user.email}"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `environment_id` (String) The ID of the environment to create the identity provider attribute in.
- `identity_provider_id` (String) The ID of the identity provider to create the attribute mapping for.
- `name` (String) The user attribute, which is unique per provider. The attribute must not be defined as read only from the user schema or of type `COMPLEX` based on the user schema. Valid examples `username`, and `name.first`. The following attributes may not be used `account`, `id`, `created`, `updated`, `lifecycle`, `mfaEnabled`, and `enabled`.
- `update` (String) Indicates whether to update the user attribute in the directory with the non-empty mapped value from the IdP. Options are `EMPTY_ONLY` (only update the user attribute if it has an empty value); `ALWAYS` (always update the user attribute value).
- `value` (String) A placeholder referring to the attribute (or attributes) from the provider. Placeholders must be valid for the attributes returned by the IdP type and use the `${}` syntax (for example, `${email}`). For SAML, any placeholder is acceptable, and it is mapped against the attributes available in the SAML assertion after authentication. The `${samlAssertion.subject}` placeholder is a special reserved placeholder used to refer to the subject name ID in the SAML assertion response.

### Read-Only

- `id` (String) The ID of this resource.
- `mapping_type` (String) The mapping type. Options are `CORE` (This attribute is required by the schema and cannot be removed. The name and update properties cannot be changed.) or `CUSTOM` (All user-created attributes are of this type.)

## Import

Import is supported using the following syntax, where attributes in `<>` brackets are replaced with the relevant ID.  For example, `<environment_id>` should be replaced with the ID of the environment to import from.

```shell
$ terraform import pingone_identity_provider_attribute.example <environment_id>/<identity_provider_id>/<identity_provider_attribute_id>
```
