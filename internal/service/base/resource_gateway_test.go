package base_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	pingone "github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
)

func testAccCheckGatewayDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, pingone.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_gateway" {
			continue
		}

		_, rEnv, err := apiClient.EnvironmentsApi.ReadOneEnvironment(ctx, rs.Primary.Attributes["environment_id"]).Execute()

		if rEnv.StatusCode == 404 {
			continue
		}

		if err != nil {
			return err
		}

		body, r, err := apiClient.GatewaysApi.ReadOneGateway(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.ID).Execute()

		if r.StatusCode == 404 {
			continue
		}

		if err != nil {
			tflog.Error(ctx, fmt.Sprintf("Error: %v", body))
			return err
		}

		return fmt.Errorf("PingOne Gateway Instance %s still exists", rs.Primary.ID)
	}

	return nil
}

func TestAccGateway_NewEnv(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_gateway.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckGatewayDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccGatewayConfig_NewEnv(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
				),
			},
		},
	})
}

func TestAccGateway_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_gateway.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckGatewayDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccGatewayConfig_Full(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", "My test gateway"),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "pingfederate.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "api_gateway.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "ldap.#", "0"),
				),
			},
		},
	})
}

func TestAccGateway_Minimal(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_gateway.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckGatewayDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccGatewayConfig_Minimal(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", ""),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "pingfederate.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "api_gateway.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "ldap.#", "0"),
				),
			},
		},
	})
}

func TestAccGateway_Change(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_gateway.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckGatewayDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccGatewayConfig_Full(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", "My test gateway"),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "pingfederate.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "api_gateway.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "ldap.#", "0"),
				),
			},
			{
				Config: testAccGatewayConfig_PingFederate(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", ""),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "pingfederate.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "api_gateway.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "ldap.#", "0"),
				),
			},
			{
				Config: testAccGatewayConfig_Full(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", "My test gateway"),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "pingfederate.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "api_gateway.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "ldap.#", "0"),
				),
			},
			{
				Config: testAccGatewayConfig_APIGateway(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", ""),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "pingfederate.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "api_gateway.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "ldap.#", "0"),
				),
			},
		},
	})
}

func TestAccGateway_PF(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_gateway.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckGatewayDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccGatewayConfig_PingFederate(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "pingfederate.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "api_gateway.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "ldap.#", "0"),
				),
			},
		},
	})
}

func TestAccGateway_APIG(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_gateway.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckGatewayDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccGatewayConfig_APIGateway(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "pingfederate.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "api_gateway.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "ldap.#", "0"),
				),
			},
		},
	})
}

func TestAccGateway_LDAPFull(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_gateway.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckGatewayDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccGatewayConfig_LDAPFull(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", ""),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "pingfederate.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "api_gateway.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "ldap.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "ldap.0.bind_dn", "ou=test,dc=example,dc=com"),
					resource.TestCheckResourceAttr(resourceFullName, "ldap.0.bind_password", "dummyPasswordValue"),
					resource.TestCheckResourceAttr(resourceFullName, "ldap.0.connection_security", "TLS"),
					resource.TestCheckResourceAttr(resourceFullName, "ldap.0.kerberos_service_account_upn", "upnvalue"),
					resource.TestCheckResourceAttr(resourceFullName, "ldap.0.kerberos_service_account_password", "dummyKerberosPasswordValue"),
					resource.TestCheckResourceAttr(resourceFullName, "ldap.0.kerberos_retain_previous_credentials_mins", "20"),
					resource.TestCheckResourceAttr(resourceFullName, "ldap.0.servers.#", "3"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "ldap.0.servers.*", "ds2.dummyldapservice.com"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "ldap.0.servers.*", "ds3.dummyldapservice.com"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "ldap.0.servers.*", "ds1.dummyldapservice.com"),
					resource.TestCheckResourceAttr(resourceFullName, "ldap.0.validate_tls_certificates", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "ldap.0.vendor", "PingDirectory"),
					resource.TestCheckResourceAttr(resourceFullName, "ldap.0.user_type.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "ldap.0.user_type.*", map[string]string{
						"id":                                   "id2-1234",
						"name":                                 "User Set 2",
						"password_authority":                   "PING_ONE",
						"search_base_dn":                       "ou=users,dc=example,dc=com",
						"user_link_attributes.#":               "3",
						"user_migration_lookup_filter_pattern": "((uid=$${identifier})(mail=$${identifier}))",
						"user_migration_attribute_mapping.#":   "3",
						"push_password_changes_to_ldap":        "true",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "ldap.0.user_type.*", map[string]string{
						"id":                                   "id1-1234",
						"name":                                 "User Set 1",
						"password_authority":                   "LDAP",
						"search_base_dn":                       "ou=users1,dc=example,dc=com",
						"user_link_attributes.#":               "2",
						"user_migration_lookup_filter_pattern": "((uid=$${identifier})(mail=$${identifier}))",
						"user_migration_attribute_mapping.#":   "2",
						"push_password_changes_to_ldap":        "true",
					}),
				),
			},
		},
	})
}

func TestAccGateway_LDAPMinimal(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_gateway.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckGatewayDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccGatewayConfig_LDAPMinimal(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", ""),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "pingfederate.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "api_gateway.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "ldap.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "ldap.0.bind_dn", "ou=test,dc=example,dc=com"),
					resource.TestCheckResourceAttr(resourceFullName, "ldap.0.bind_password", "dummyPasswordValue"),
					resource.TestCheckResourceAttr(resourceFullName, "ldap.0.connection_security", "NONE"),
					resource.TestCheckResourceAttr(resourceFullName, "ldap.0.kerberos_service_account_upn", ""),
					resource.TestCheckResourceAttr(resourceFullName, "ldap.0.kerberos_service_account_password", ""),
					resource.TestCheckResourceAttr(resourceFullName, "ldap.0.kerberos_retain_previous_credentials_mins", ""),
					resource.TestCheckResourceAttr(resourceFullName, "ldap.0.servers.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "ldap.0.validate_tls_certificates", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "ldap.0.vendor", "PingDirectory"),
					resource.TestCheckResourceAttr(resourceFullName, "ldap.0.user_type.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "ldap.0.user_type.*", map[string]string{
						"id":                                   "id3-1234",
						"name":                                 "User Set 3",
						"password_authority":                   "LDAP",
						"search_base_dn":                       "",
						"user_link_attributes.#":               "0",
						"user_migration_lookup_filter_pattern": "",
						"user_migration_population_id":         "",
						"user_migration_attribute_mapping.#":   "1",
						"push_password_changes_to_ldap":        "false",
					}),
				),
			},
		},
	})
}

func testAccGatewayConfig_NewEnv(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_gateway" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  name           = "%[4]s"
  enabled        = false

  pingfederate {}
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccGatewayConfig_Full(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_gateway" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "My test gateway"
  enabled        = true

  pingfederate {}
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccGatewayConfig_Minimal(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_gateway" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = false

  pingfederate {}
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccGatewayConfig_PingFederate(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_gateway" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = false

  pingfederate {}
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccGatewayConfig_APIGateway(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_gateway" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = false

  api_gateway {}
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccGatewayConfig_LDAPFull(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_population" "%[2]s" {
	environment_id = data.pingone_environment.general_test.id

	name = "Gateway Population Test"
}

resource "pingone_gateway" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = false

  ldap {

	bind_dn = "ou=test,dc=example,dc=com"
	bind_password = "dummyPasswordValue"
	connection_security = "TLS"
	vendor = "PingDirectory"

	kerberos_service_account_upn = "upnvalue"
	kerberos_service_account_password = "dummyKerberosPasswordValue"
	kerberos_retain_previous_credentials_mins = 20

	servers = [
		"ds1.dummyldapservice.com",
		"ds3.dummyldapservice.com",
		"ds2.dummyldapservice.com",
	]

	validate_tls_certificates = false

	user_type {
		id = "59e24997-f829-4206-b1b7-9b6a8a25c0b4"
		name = "User Set 1"
		password_authority = "LDAP"
		search_base_dn = "ou=users1,dc=example,dc=com"

		user_link_attributes = [ "entryUUID", "uid" ]

		user_migration_lookup_filter_pattern = "((uid=$${identifier})(mail=$${identifier}))"

		user_migration_population_id = pingone_population.%[2]s.id

		user_migration_attribute_mapping {
			name = "username"
			value = "$${ldapAttributes.uid}"
		}

		user_migration_attribute_mapping {
			name = "email"
			value = "$${ldapAttributes.mail}"
		}

		push_password_changes_to_ldap = true
	}

	user_type {
		id = "59e24997-f829-4206-b1b7-9b6a8a25c0b3"
		name = "User Set 2"
		password_authority = "PING_ONE"
		search_base_dn = "ou=users,dc=example,dc=com"

		user_link_attributes = [ "entryUUID", "dn",  "uid" ]

		user_migration_lookup_filter_pattern = "((uid=$${identifier})(mail=$${identifier}))"

		user_migration_population_id = pingone_population.%[2]s.id

		user_migration_attribute_mapping {
			name = "username"
			value = "$${ldapAttributes.uid}"
		}

		user_migration_attribute_mapping {
			name = "email"
			value = "$${ldapAttributes.mail}"
		}

		user_migration_attribute_mapping {
			name = "name.family"
			value = "$${ldapAttributes.sn}"
		}

		push_password_changes_to_ldap = true
	}

  }

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccGatewayConfig_LDAPMinimal(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_gateway" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = false

  ldap {

	bind_dn = "ou=test,dc=example,dc=com"
	bind_password = "dummyPasswordValue"
	vendor = "PingDirectory"

	user_type {
		id = "59e24997-f829-4206-b1b7-9b6a8a25c0b3"
		name = "User Set 3"
		password_authority = "LDAP"

		user_migration_attribute_mapping {
			name = "username"
			value = "$${ldapAttributes.uid}"
		}
	}

  }

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}
