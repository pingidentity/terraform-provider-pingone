package base_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/base"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccGatewayDataSource_FindGatewayByID(t *testing.T) {
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
					resource.TestMatchResourceAttr(dataSourceFullName, "radius_davinci_policy_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullName, "radius_default_shared_secret", "sharedsecret123"),
					resource.TestCheckResourceAttr(dataSourceFullName, "radius_client.#", "2"),

					resource.TestCheckTypeSetElemNestedAttrs(dataSourceFullName, "radius_client.*", map[string]string{
						"ip":            "127.0.0.1",
						"shared_secret": "sharedsecret123-1",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(dataSourceFullName, "radius_client.*", map[string]string{
						"ip":            "127.0.0.2",
						"shared_secret": "sharedsecret123-2",
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
					//resource.TestCheckResourceAttr(dataSourceFullName, "bind_password", "dummyPasswordValue1"),
					resource.TestCheckResourceAttr(dataSourceFullName, "connection_security", "TLS"),
					resource.TestCheckResourceAttr(dataSourceFullName, "kerberos_service_account_upn", "username@domainname"),
					//resource.TestCheckResourceAttr(dataSourceFullName, "kerberos_service_account_password", "dummyKerberosPasswordValue"),
					resource.TestCheckResourceAttr(dataSourceFullName, "kerberos_retain_previous_credentials_mins", "20"),
					resource.TestCheckResourceAttr(dataSourceFullName, "servers.#", "3"),
					resource.TestCheckTypeSetElemAttr(dataSourceFullName, "servers.*", "ds2.dummyldapservice.com:636"),
					resource.TestCheckTypeSetElemAttr(dataSourceFullName, "servers.*", "ds3.dummyldapservice.com:636"),
					resource.TestCheckTypeSetElemAttr(dataSourceFullName, "servers.*", "ds1.dummyldapservice.com:636"),
					resource.TestCheckResourceAttr(dataSourceFullName, "validate_tls_certificates", "false"),
					resource.TestCheckResourceAttr(dataSourceFullName, "vendor", "Microsoft Active Directory"),
					resource.TestCheckResourceAttr(dataSourceFullName, "user_type.#", "2"),

					resource.TestCheckTypeSetElemNestedAttrs(dataSourceFullName, "user_type.*", map[string]string{
						"name":                                   "User Set 2",
						"password_authority":                     "PING_ONE",
						"search_base_dn":                         "ou=users,dc=example,dc=com",
						"user_link_attributes.#":                 "3",
						"user_link_attributes.0":                 "objectGUID",
						"user_link_attributes.1":                 "dn",
						"user_link_attributes.2":                 "objectSid",
						"user_migration.#":                       "1",
						"user_migration.0.lookup_filter_pattern": "(|(uid=${identifier})(mail=${identifier}))",
						"user_migration.0.attribute_mapping.#":   "3",
						"push_password_changes_to_ldap":          "true",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(dataSourceFullName, "user_type.*", map[string]string{
						"name":                                   "User Set 1",
						"password_authority":                     "LDAP",
						"search_base_dn":                         "ou=users1,dc=example,dc=com",
						"user_link_attributes.#":                 "2",
						"user_link_attributes.0":                 "objectGUID",
						"user_link_attributes.1":                 "objectSid",
						"user_migration.#":                       "1",
						"user_migration.0.lookup_filter_pattern": "(|(uid=${identifier})(mail=${identifier}))",
						"user_migration.0.attribute_mapping.#":   "2",
						"push_password_changes_to_ldap":          "true",
					}),
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
					resource.TestCheckResourceAttr(dataSourceFullName, "description", ""),
					resource.TestCheckResourceAttr(dataSourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "type", "LDAP"),
					resource.TestCheckResourceAttr(dataSourceFullName, "bind_dn", "ou=test1,dc=example,dc=com"),
					resource.TestCheckResourceAttr(dataSourceFullName, "bind_password", "dummyPasswordValue1"),
					resource.TestCheckResourceAttr(dataSourceFullName, "connection_security", "TLS"),
					resource.TestCheckResourceAttr(dataSourceFullName, "kerberos_service_account_upn", "username@domainname"),
					resource.TestCheckResourceAttr(dataSourceFullName, "kerberos_service_account_password", "dummyKerberosPasswordValue"),
					resource.TestCheckResourceAttr(dataSourceFullName, "kerberos_retain_previous_credentials_mins", "20"),
					resource.TestCheckResourceAttr(dataSourceFullName, "servers.#", "3"),
					resource.TestCheckTypeSetElemAttr(dataSourceFullName, "servers.*", "ds2.dummyldapservice.com:636"),
					resource.TestCheckTypeSetElemAttr(dataSourceFullName, "servers.*", "ds3.dummyldapservice.com:636"),
					resource.TestCheckTypeSetElemAttr(dataSourceFullName, "servers.*", "ds1.dummyldapservice.com:636"),
					resource.TestCheckResourceAttr(dataSourceFullName, "validate_tls_certificates", "false"),
					resource.TestCheckResourceAttr(dataSourceFullName, "vendor", "Microsoft Active Directory"),
					resource.TestCheckResourceAttr(dataSourceFullName, "user_type.#", "2"),

					resource.TestCheckTypeSetElemNestedAttrs(dataSourceFullName, "user_type.*", map[string]string{
						"name":                                   "User Set 2",
						"password_authority":                     "PING_ONE",
						"search_base_dn":                         "ou=users,dc=example,dc=com",
						"user_link_attributes.#":                 "3",
						"user_link_attributes.0":                 "objectGUID",
						"user_link_attributes.1":                 "dn",
						"user_link_attributes.2":                 "objectSid",
						"user_migration.#":                       "1",
						"user_migration.0.lookup_filter_pattern": "(|(uid=${identifier})(mail=${identifier}))",
						"user_migration.0.attribute_mapping.#":   "3",
						"push_password_changes_to_ldap":          "true",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(dataSourceFullName, "user_type.*", map[string]string{
						"name":                          "User Set 1",
						"password_authority":            "LDAP",
						"search_base_dn":                "ou=users1,dc=example,dc=com",
						"user_link_attributes.#":        "2",
						"user_link_attributes.0":        "objectGUID",
						"user_link_attributes.1":        "objectSid",
						"user_migration.#":              "0",
						"push_password_changes_to_ldap": "true",
					}),
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

func testAccGatewayDataSource_FindRADIUSGatewayByID(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_gateway" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = false
  type           = "RADIUS"

  radius_default_shared_secret = "sharedsecret123"

  radius_davinci_policy_id = "ee8470a2-8161-4d76-a7af-a8505a2da084" // dummy ID

  radius_client {
    ip            = "127.0.0.1"
    shared_secret = "sharedsecret123-1"
  }

  radius_client {
    ip            = "127.0.0.2"
    shared_secret = "sharedsecret123-2"
  }
}

data "pingone_gateway" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  gateway_id     = pingone_gateway.%[2]s.id
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

  kerberos_service_account_upn              = "username@domainname"
  kerberos_service_account_password         = "dummyKerberosPasswordValue"
  kerberos_retain_previous_credentials_mins = 20

  servers = [
    "ds1.dummyldapservice.com:636",
    "ds3.dummyldapservice.com:636",
    "ds2.dummyldapservice.com:636",
  ]

  validate_tls_certificates = false

  user_type {
    name               = "User Set 1"
    password_authority = "LDAP"
    search_base_dn     = "ou=users1,dc=example,dc=com"

    user_link_attributes = ["objectGUID", "objectSid"]

    user_migration {
      lookup_filter_pattern = "(|(uid=$${identifier})(mail=$${identifier}))"

      population_id = pingone_population.%[2]s.id

      attribute_mapping {
        name  = "username"
        value = "$${ldapAttributes.uid}"
      }

      attribute_mapping {
        name  = "email"
        value = "$${ldapAttributes.mail}"
      }
    }

    push_password_changes_to_ldap = true
  }

  user_type {
    name               = "User Set 2"
    password_authority = "PING_ONE"
    search_base_dn     = "ou=users,dc=example,dc=com"

    user_link_attributes = ["objectGUID", "dn", "objectSid"]

    user_migration {
      lookup_filter_pattern = "(|(uid=$${identifier})(mail=$${identifier}))"

      population_id = pingone_population.%[2]s.id

      attribute_mapping {
        name  = "username"
        value = "$${ldapAttributes.uid}"
      }

      attribute_mapping {
        name  = "email"
        value = "$${ldapAttributes.mail}"
      }

      attribute_mapping {
        name  = "name.family"
        value = "$${ldapAttributes.sn}"
      }
    }

    push_password_changes_to_ldap = true
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

  kerberos_service_account_upn              = "username@domainname"
  kerberos_service_account_password         = "dummyKerberosPasswordValue"
  kerberos_retain_previous_credentials_mins = 20

  servers = [
    "ds1.dummyldapservice.com:636",
    "ds3.dummyldapservice.com:636",
    "ds2.dummyldapservice.com:636",
  ]

  validate_tls_certificates = false

  user_type {
    name               = "User Set 1"
    password_authority = "LDAP"
    search_base_dn     = "ou=users1,dc=example,dc=com"

    user_link_attributes = ["objectGUID", "objectSid"]

    push_password_changes_to_ldap = true
  }

  user_type {
    name               = "User Set 2"
    password_authority = "PING_ONE"
    search_base_dn     = "ou=users,dc=example,dc=com"

    user_link_attributes = ["objectGUID", "dn", "objectSid"]

    user_migration {
      lookup_filter_pattern = "(|(uid=$${identifier})(mail=$${identifier}))"

      population_id = pingone_population.%[2]s.id

      attribute_mapping {
        name  = "username"
        value = "$${ldapAttributes.uid}"
      }

      attribute_mapping {
        name  = "email"
        value = "$${ldapAttributes.mail}"
      }

      attribute_mapping {
        name  = "name.family"
        value = "$${ldapAttributes.sn}"
      }
    }

    push_password_changes_to_ldap = true
  }
}

data "pingone_gateway" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name     = "%[3]s"
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}
