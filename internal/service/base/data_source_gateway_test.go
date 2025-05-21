// Copyright © 2025 Ping Identity Corporation

package base_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/base"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccGatewayDataSource_FindGatewayAll(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	dataSourceFullName := fmt.Sprintf("data.pingone_gateway.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Gateway_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccGatewayDataSource_FindGatewayByID(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "gateway_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullName, "name", name),
					resource.TestCheckResourceAttr(dataSourceFullName, "description", name),
					resource.TestCheckResourceAttr(dataSourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "type", "PING_FEDERATE"),
				),
			},
			{
				Config: testAccGatewayDataSource_FindGatewayByName(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "gateway_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullName, "name", name),

					resource.TestCheckNoResourceAttr(dataSourceFullName, "description"),

					resource.TestCheckResourceAttr(dataSourceFullName, "enabled", "false"),
					resource.TestCheckResourceAttr(dataSourceFullName, "type", "API_GATEWAY_INTEGRATION"),
				),
			},
		},
	})
}

func TestAccGatewayDataSource_FindRADIUSGatewayByID(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	dataSourceFullName := fmt.Sprintf("data.pingone_gateway.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Gateway_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccGatewayDataSource_FindRADIUSGatewayByID(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullName, "name", name),

					resource.TestCheckNoResourceAttr(dataSourceFullName, "description"),

					resource.TestCheckResourceAttr(dataSourceFullName, "enabled", "false"),
					resource.TestCheckResourceAttr(dataSourceFullName, "type", "RADIUS"),
					resource.TestMatchResourceAttr(dataSourceFullName, "radius_davinci_policy_id", verify.P1DVResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullName, "radius_default_shared_secret", "sharedsecret123"),
					resource.TestCheckResourceAttr(dataSourceFullName, "radius_clients.#", "2"),

					resource.TestCheckTypeSetElemNestedAttrs(dataSourceFullName, "radius_clients.*", map[string]string{
						"ip":            "127.0.0.1",
						"shared_secret": "sharedsecret123-1",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(dataSourceFullName, "radius_clients.*", map[string]string{
						"ip":            "127.0.0.2",
						"shared_secret": "sharedsecret123-2",
					}),
				),
			},
		},
	})
}

func TestAccGatewayDataSource_FindRADIUSGatewayByName(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	dataSourceFullName := fmt.Sprintf("data.pingone_gateway.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Gateway_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccGatewayDataSource_FindRADIUSGatewayByName(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullName, "name", name),
					resource.TestCheckNoResourceAttr(dataSourceFullName, "description"),
					resource.TestCheckResourceAttr(dataSourceFullName, "enabled", "false"),
					resource.TestCheckResourceAttr(dataSourceFullName, "type", "RADIUS"),
					resource.TestMatchResourceAttr(dataSourceFullName, "radius_davinci_policy_id", verify.P1DVResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullName, "radius_default_shared_secret", "sharedsecret123"),
					resource.TestCheckResourceAttr(dataSourceFullName, "radius_clients.#", "1"),

					resource.TestCheckTypeSetElemNestedAttrs(dataSourceFullName, "radius_clients.*", map[string]string{
						"ip":            "127.0.0.3",
						"shared_secret": "",
					}),
				),
			},
		},
	})
}

func TestAccGatewayDataSource_FindLDAPGatewayByID(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	dataSourceFullName := fmt.Sprintf("data.pingone_gateway.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Gateway_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccGatewayDataSource_FindLDAPGatewayByID(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullName, "name", name),

					resource.TestCheckNoResourceAttr(dataSourceFullName, "description"),

					resource.TestCheckResourceAttr(dataSourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "type", "LDAP"),
					resource.TestCheckResourceAttr(dataSourceFullName, "bind_dn", "ou=test1,dc=example,dc=com"),
					resource.TestCheckResourceAttr(dataSourceFullName, "connection_security", "TLS"),
					resource.TestCheckResourceAttr(dataSourceFullName, "kerberos.service_account_upn", "username@domainname"),
					resource.TestCheckResourceAttr(dataSourceFullName, "kerberos.retain_previous_credentials_mins", "20"),
					resource.TestCheckResourceAttr(dataSourceFullName, "servers.#", "3"),
					resource.TestCheckTypeSetElemAttr(dataSourceFullName, "servers.*", "ds2.dummyldapservice.com:636"),
					resource.TestCheckTypeSetElemAttr(dataSourceFullName, "servers.*", "ds3.dummyldapservice.com:636"),
					resource.TestCheckTypeSetElemAttr(dataSourceFullName, "servers.*", "ds1.dummyldapservice.com:636"),
					resource.TestCheckResourceAttr(dataSourceFullName, "validate_tls_certificates", "false"),
					resource.TestCheckResourceAttr(dataSourceFullName, "vendor", "Microsoft Active Directory"),
					resource.TestCheckResourceAttr(dataSourceFullName, "user_types.%", "2"),

					resource.TestCheckNoResourceAttr(dataSourceFullName, "bind_password"),
					resource.TestCheckNoResourceAttr(dataSourceFullName, "kerberos.service_account_password"),

					resource.TestCheckResourceAttr(dataSourceFullName, "user_types.User Set 2.password_authority", "PING_ONE"),
					resource.TestCheckResourceAttr(dataSourceFullName, "user_types.User Set 2.search_base_dn", "ou=users,dc=example,dc=com"),
					resource.TestCheckResourceAttr(dataSourceFullName, "user_types.User Set 2.user_link_attributes.#", "3"),
					resource.TestCheckResourceAttr(dataSourceFullName, "user_types.User Set 2.user_link_attributes.0", "objectGUID"),
					resource.TestCheckResourceAttr(dataSourceFullName, "user_types.User Set 2.user_link_attributes.1", "dn"),
					resource.TestCheckResourceAttr(dataSourceFullName, "user_types.User Set 2.user_link_attributes.2", "objectSid"),
					resource.TestCheckResourceAttr(dataSourceFullName, "user_types.User Set 2.new_user_lookup.ldap_filter_pattern", "(|(uid=${identifier})(mail=${identifier}))"),
					resource.TestCheckResourceAttr(dataSourceFullName, "user_types.User Set 2.new_user_lookup.attribute_mappings.#", "3"),
					resource.TestCheckResourceAttr(dataSourceFullName, "user_types.User Set 2.allow_password_changes", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "user_types.User Set 2.update_user_on_successful_authentication", "false"),

					resource.TestCheckResourceAttr(dataSourceFullName, "user_types.User Set 1.password_authority", "LDAP"),
					resource.TestCheckResourceAttr(dataSourceFullName, "user_types.User Set 1.search_base_dn", "ou=users1,dc=example,dc=com"),
					resource.TestCheckResourceAttr(dataSourceFullName, "user_types.User Set 1.user_link_attributes.#", "2"),
					resource.TestCheckResourceAttr(dataSourceFullName, "user_types.User Set 1.user_link_attributes.0", "objectGUID"),
					resource.TestCheckResourceAttr(dataSourceFullName, "user_types.User Set 1.user_link_attributes.1", "objectSid"),
					resource.TestCheckResourceAttr(dataSourceFullName, "user_types.User Set 1.new_user_lookup.ldap_filter_pattern", "(|(uid=${identifier})(mail=${identifier}))"),
					resource.TestCheckResourceAttr(dataSourceFullName, "user_types.User Set 1.new_user_lookup.attribute_mappings.#", "2"),
					resource.TestCheckResourceAttr(dataSourceFullName, "user_types.User Set 1.allow_password_changes", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "user_types.User Set 1.update_user_on_successful_authentication", "true"),
				),
			},
		},
	})
}

func TestAccGatewayDataSource_FindLDAPGatewayByName(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	dataSourceFullName := fmt.Sprintf("data.pingone_gateway.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Gateway_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccGatewayDataSource_FindLDAPGatewayByName(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullName, "name", name),
					resource.TestCheckNoResourceAttr(dataSourceFullName, "description"),
					resource.TestCheckResourceAttr(dataSourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "type", "LDAP"),
					resource.TestCheckResourceAttr(dataSourceFullName, "bind_dn", "ou=test1,dc=example,dc=com"),
					resource.TestCheckResourceAttr(dataSourceFullName, "connection_security", "TLS"),
					resource.TestCheckResourceAttr(dataSourceFullName, "kerberos.service_account_upn", "username@domainname"),
					resource.TestCheckResourceAttr(dataSourceFullName, "kerberos.retain_previous_credentials_mins", "20"),
					resource.TestCheckResourceAttr(dataSourceFullName, "servers.#", "3"),
					resource.TestCheckTypeSetElemAttr(dataSourceFullName, "servers.*", "ds2.dummyldapservice.com:636"),
					resource.TestCheckTypeSetElemAttr(dataSourceFullName, "servers.*", "ds3.dummyldapservice.com:636"),
					resource.TestCheckTypeSetElemAttr(dataSourceFullName, "servers.*", "ds1.dummyldapservice.com:636"),
					resource.TestCheckResourceAttr(dataSourceFullName, "validate_tls_certificates", "false"),
					resource.TestCheckResourceAttr(dataSourceFullName, "vendor", "Microsoft Active Directory"),
					resource.TestCheckResourceAttr(dataSourceFullName, "user_types.%", "2"),

					resource.TestCheckNoResourceAttr(dataSourceFullName, "bind_password"),
					resource.TestCheckNoResourceAttr(dataSourceFullName, "kerberos_service_account_password"),

					resource.TestCheckResourceAttr(dataSourceFullName, "user_types.User Set 2.password_authority", "PING_ONE"),
					resource.TestCheckResourceAttr(dataSourceFullName, "user_types.User Set 2.search_base_dn", "ou=users,dc=example,dc=com"),
					resource.TestCheckResourceAttr(dataSourceFullName, "user_types.User Set 2.user_link_attributes.#", "3"),
					resource.TestCheckResourceAttr(dataSourceFullName, "user_types.User Set 2.user_link_attributes.0", "objectGUID"),
					resource.TestCheckResourceAttr(dataSourceFullName, "user_types.User Set 2.user_link_attributes.1", "dn"),
					resource.TestCheckResourceAttr(dataSourceFullName, "user_types.User Set 2.user_link_attributes.2", "objectSid"),
					resource.TestCheckResourceAttr(dataSourceFullName, "user_types.User Set 2.new_user_lookup.ldap_filter_pattern", "(|(uid=${identifier})(mail=${identifier}))"),
					resource.TestCheckResourceAttr(dataSourceFullName, "user_types.User Set 2.new_user_lookup.attribute_mappings.#", "3"),
					resource.TestCheckResourceAttr(dataSourceFullName, "user_types.User Set 2.allow_password_changes", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "user_types.User Set 2.update_user_on_successful_authentication", "false"),

					resource.TestCheckResourceAttr(dataSourceFullName, "user_types.User Set 1.password_authority", "LDAP"),
					resource.TestCheckResourceAttr(dataSourceFullName, "user_types.User Set 1.search_base_dn", "ou=users1,dc=example,dc=com"),
					resource.TestCheckResourceAttr(dataSourceFullName, "user_types.User Set 1.user_link_attributes.#", "2"),
					resource.TestCheckResourceAttr(dataSourceFullName, "user_types.User Set 1.user_link_attributes.0", "objectGUID"),
					resource.TestCheckResourceAttr(dataSourceFullName, "user_types.User Set 1.user_link_attributes.1", "objectSid"),
					resource.TestCheckResourceAttr(dataSourceFullName, "user_types.User Set 1.allow_password_changes", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "user_types.User Set 1.update_user_on_successful_authentication", "false"),
				),
			},
		},
	})
}

func testAccGatewayDataSource_FindGatewayByID(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_gateway" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "%[3]s"
  enabled        = true

  type = "PING_FEDERATE"
}

data "pingone_gateway" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  gateway_id     = pingone_gateway.%[2]s.id
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccGatewayDataSource_FindGatewayByName(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_gateway" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = false

  type = "API_GATEWAY_INTEGRATION"
}

data "pingone_gateway" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  depends_on = [pingone_gateway.%[2]s]
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccGatewayDataSource_FindRADIUSGatewayByID(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_gateway" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = false
  type           = "RADIUS"

  radius_default_shared_secret = "sharedsecret123"

  radius_davinci_policy_id = "ee8470a281614d76a7afa8505a2da084" // dummy ID

  radius_clients = [
    {
      ip            = "127.0.0.1"
      shared_secret = "sharedsecret123-1"
    },
    {
      ip            = "127.0.0.2"
      shared_secret = "sharedsecret123-2"
    }
  ]
}

data "pingone_gateway" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  gateway_id     = pingone_gateway.%[2]s.id
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccGatewayDataSource_FindRADIUSGatewayByName(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_gateway" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = false
  type           = "RADIUS"

  radius_default_shared_secret = "sharedsecret123"

  radius_davinci_policy_id = "ee8470a281614d76a7afa8505a2da085" // dummy ID

  radius_clients = [
    {
      ip = "127.0.0.3"
    }
  ]
}

data "pingone_gateway" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  depends_on = [pingone_gateway.%[2]s]
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccGatewayDataSource_FindLDAPGatewayByID(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_population" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "pingone_gateway" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = true
  type           = "LDAP"

  bind_dn       = "ou=test1,dc=example,dc=com"
  bind_password = "dummyPasswordValue1"

  connection_security = "TLS"
  vendor              = "Microsoft Active Directory"

  kerberos = {
    service_account_upn              = "username@domainname"
    service_account_password         = "dummyKerberosPasswordValue"
    retain_previous_credentials_mins = 20
  }

  servers = [
    "ds1.dummyldapservice.com:636",
    "ds3.dummyldapservice.com:636",
    "ds2.dummyldapservice.com:636",
  ]

  validate_tls_certificates = false

  user_types = {
    "User Set 1" = {
      password_authority = "LDAP"
      search_base_dn     = "ou=users1,dc=example,dc=com"

      user_link_attributes = ["objectGUID", "objectSid"]

      allow_password_changes                   = true
      update_user_on_successful_authentication = true

      new_user_lookup = {
        ldap_filter_pattern = "(|(uid=$${identifier})(mail=$${identifier}))"

        population_id = pingone_population.%[2]s.id

        attribute_mappings = [
          {
            name  = "username"
            value = "$${ldapAttributes.uid}"
          },
          {
            name  = "email"
            value = "$${ldapAttributes.mail}"
          }
        ]
      }
    },
    "User Set 2" = {
      password_authority = "PING_ONE"
      search_base_dn     = "ou=users,dc=example,dc=com"

      user_link_attributes = ["objectGUID", "dn", "objectSid"]

      allow_password_changes                   = true
      update_user_on_successful_authentication = false

      new_user_lookup = {
        ldap_filter_pattern = "(|(uid=$${identifier})(mail=$${identifier}))"

        population_id = pingone_population.%[2]s.id

        attribute_mappings = [
          {
            name  = "username"
            value = "$${ldapAttributes.uid}"
          },
          {
            name  = "email"
            value = "$${ldapAttributes.mail}"
          },
          {
            name  = "name.family"
            value = "$${ldapAttributes.sn}"
          }
        ]
      }
    }
  }
}

data "pingone_gateway" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  gateway_id     = pingone_gateway.%[2]s.id
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccGatewayDataSource_FindLDAPGatewayByName(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_population" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "pingone_gateway" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = true
  type           = "LDAP"

  bind_dn       = "ou=test1,dc=example,dc=com"
  bind_password = "dummyPasswordValue1"

  connection_security = "TLS"
  vendor              = "Microsoft Active Directory"

  kerberos = {
    service_account_upn              = "username@domainname"
    service_account_password         = "dummyKerberosPasswordValue"
    retain_previous_credentials_mins = 20
  }

  servers = [
    "ds1.dummyldapservice.com:636",
    "ds3.dummyldapservice.com:636",
    "ds2.dummyldapservice.com:636",
  ]

  validate_tls_certificates = false

  user_types = {
    "User Set 1" = {
      password_authority = "LDAP"
      search_base_dn     = "ou=users1,dc=example,dc=com"

      allow_password_changes                   = true
      update_user_on_successful_authentication = false

      user_link_attributes = ["objectGUID", "objectSid"]
    },
    "User Set 2" = {
      password_authority = "PING_ONE"
      search_base_dn     = "ou=users,dc=example,dc=com"

      user_link_attributes = ["objectGUID", "dn", "objectSid"]

      allow_password_changes                   = true
      update_user_on_successful_authentication = false

      new_user_lookup = {
        ldap_filter_pattern = "(|(uid=$${identifier})(mail=$${identifier}))"

        population_id = pingone_population.%[2]s.id

        attribute_mappings = [
          {
            name  = "username"
            value = "$${ldapAttributes.uid}"
          },
          {
            name  = "email"
            value = "$${ldapAttributes.mail}"
          },
          {
            name  = "name.family"
            value = "$${ldapAttributes.sn}"
          }
        ]
      }
    }
  }
}

data "pingone_gateway" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  depends_on = [pingone_gateway.%[2]s]
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}
