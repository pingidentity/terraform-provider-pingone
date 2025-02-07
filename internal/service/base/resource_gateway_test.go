// Copyright © 2025 Ping Identity Corporation

package base_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/base"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccGateway_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_gateway.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	var gatewayID, environmentID string

	var p1Client *client.Client
	var ctx = context.Background()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)

			p1Client = acctest.PreCheckTestClient(ctx, t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Gateway_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccGatewayConfig_Minimal(resourceName, name),
				Check:  base.Gateway_GetIDs(resourceFullName, &environmentID, &gatewayID),
			},
			{
				PreConfig: func() {
					base.Gateway_RemovalDrift_PreConfig(ctx, p1Client.API.ManagementAPIClient, t, environmentID, gatewayID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			// Test removal of the environment
			{
				Config: testAccGatewayConfig_NewEnv(environmentName, licenseID, resourceName, name),
				Check:  base.Gateway_GetIDs(resourceFullName, &environmentID, &gatewayID),
			},
			{
				PreConfig: func() {
					base.Environment_RemovalDrift_PreConfig(ctx, p1Client.API.ManagementAPIClient, t, environmentID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccGateway_NewEnv(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_gateway.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Gateway_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
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
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Gateway_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccGatewayConfig_Full(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", "My test gateway"),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "type", "PING_FEDERATE"),
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
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Gateway_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccGatewayConfig_Minimal(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckNoResourceAttr(resourceFullName, "description"),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "type", "PING_FEDERATE"),
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
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Gateway_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccGatewayConfig_Full(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", "My test gateway"),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "type", "PING_FEDERATE"),
				),
			},
			{
				Config: testAccGatewayConfig_PingFederate(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckNoResourceAttr(resourceFullName, "description"),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "type", "PING_FEDERATE"),
				),
			},
			{
				Config: testAccGatewayConfig_Full(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", "My test gateway"),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "type", "PING_FEDERATE"),
				),
			},
			{
				Config: testAccGatewayConfig_APIGateway(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckNoResourceAttr(resourceFullName, "description"),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "type", "API_GATEWAY_INTEGRATION"),
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
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Gateway_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccGatewayConfig_PingFederate(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "type", "PING_FEDERATE"),
				),
			},
			// Test importing the resource
			{
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("Resource Not found: %s", resourceFullName)
						}

						return fmt.Sprintf("%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
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
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Gateway_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccGatewayConfig_APIGateway(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "type", "API_GATEWAY_INTEGRATION"),
				),
			},
			// Test importing the resource
			{
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("Resource Not found: %s", resourceFullName)
						}

						return fmt.Sprintf("%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccGateway_Intelligence(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_gateway.%s", resourceName)

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
				Config: testAccGatewayConfig_Intelligence(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "type", "PING_INTELLIGENCE"),
				),
			},
			// Test importing the resource
			{
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("Resource Not found: %s", resourceFullName)
						}

						return fmt.Sprintf("%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccGateway_LDAP(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_gateway.%s", resourceName)

	name := resourceName

	fullStep1 := resource.TestStep{
		Config: testAccGatewayConfig_LDAPFull1(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "name", name),
			resource.TestCheckNoResourceAttr(resourceFullName, "description"),
			resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "type", "LDAP"),
			resource.TestCheckResourceAttr(resourceFullName, "bind_dn", "ou=test1,dc=example,dc=com"),
			resource.TestCheckResourceAttr(resourceFullName, "bind_password", "dummyPasswordValue1"),
			resource.TestCheckResourceAttr(resourceFullName, "connection_security", "TLS"),
			resource.TestCheckResourceAttr(resourceFullName, "kerberos.service_account_upn", "username@domainname"),
			resource.TestCheckResourceAttr(resourceFullName, "kerberos.service_account_password", "dummyKerberosPasswordValue"),
			resource.TestCheckResourceAttr(resourceFullName, "kerberos.retain_previous_credentials_mins", "20"),
			resource.TestCheckResourceAttr(resourceFullName, "servers.#", "3"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "servers.*", "ds2.dummyldapservice.com:636"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "servers.*", "ds3.dummyldapservice.com:636"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "servers.*", "ds1.dummyldapservice.com:636"),
			resource.TestCheckResourceAttr(resourceFullName, "validate_tls_certificates", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "vendor", "Microsoft Active Directory"),
			resource.TestCheckResourceAttr(resourceFullName, "user_types.%", "2"),

			resource.TestCheckResourceAttr(resourceFullName, "user_types.User Set 2.password_authority", "PING_ONE"),
			resource.TestCheckResourceAttr(resourceFullName, "user_types.User Set 2.search_base_dn", "ou=users,dc=example,dc=com"),
			resource.TestCheckResourceAttr(resourceFullName, "user_types.User Set 2.user_link_attributes.#", "3"),
			resource.TestCheckResourceAttr(resourceFullName, "user_types.User Set 2.user_link_attributes.0", "objectGUID"),
			resource.TestCheckResourceAttr(resourceFullName, "user_types.User Set 2.user_link_attributes.1", "dn"),
			resource.TestCheckResourceAttr(resourceFullName, "user_types.User Set 2.user_link_attributes.2", "objectSid"),
			resource.TestCheckResourceAttr(resourceFullName, "user_types.User Set 2.new_user_lookup.ldap_filter_pattern", "(|(uid=${identifier})(mail=${identifier}))"),
			resource.TestCheckResourceAttr(resourceFullName, "user_types.User Set 2.new_user_lookup.attribute_mappings.#", "3"),
			resource.TestCheckResourceAttr(resourceFullName, "user_types.User Set 2.allow_password_changes", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "user_types.User Set 2.update_user_on_successful_authentication", "false"),

			resource.TestCheckResourceAttr(resourceFullName, "user_types.User Set 1.password_authority", "LDAP"),
			resource.TestCheckResourceAttr(resourceFullName, "user_types.User Set 1.search_base_dn", "ou=users1,dc=example,dc=com"),
			resource.TestCheckResourceAttr(resourceFullName, "user_types.User Set 1.user_link_attributes.#", "2"),
			resource.TestCheckResourceAttr(resourceFullName, "user_types.User Set 1.user_link_attributes.0", "objectGUID"),
			resource.TestCheckResourceAttr(resourceFullName, "user_types.User Set 1.user_link_attributes.1", "objectSid"),
			resource.TestCheckResourceAttr(resourceFullName, "user_types.User Set 1.new_user_lookup.ldap_filter_pattern", "(|(uid=${identifier})(mail=${identifier}))"),
			resource.TestCheckResourceAttr(resourceFullName, "user_types.User Set 1.new_user_lookup.attribute_mappings.#", "2"),
			resource.TestCheckResourceAttr(resourceFullName, "user_types.User Set 1.allow_password_changes", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "user_types.User Set 1.update_user_on_successful_authentication", "true"),
		),
	}

	fullStep2 := resource.TestStep{
		Config: testAccGatewayConfig_LDAPFull2(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "name", name),
			resource.TestCheckNoResourceAttr(resourceFullName, "description"),
			resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "type", "LDAP"),
			resource.TestCheckResourceAttr(resourceFullName, "bind_dn", "ou=test1,dc=example,dc=com"),
			resource.TestCheckResourceAttr(resourceFullName, "bind_password", "dummyPasswordValue1"),
			resource.TestCheckResourceAttr(resourceFullName, "connection_security", "TLS"),
			resource.TestCheckResourceAttr(resourceFullName, "kerberos.service_account_upn", "username@domainname"),
			resource.TestCheckResourceAttr(resourceFullName, "kerberos.service_account_password", "dummyKerberosPasswordValue"),
			resource.TestCheckResourceAttr(resourceFullName, "kerberos.retain_previous_credentials_mins", "20"),
			resource.TestCheckResourceAttr(resourceFullName, "servers.#", "3"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "servers.*", "ds2.dummyldapservice.com:636"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "servers.*", "ds3.dummyldapservice.com:636"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "servers.*", "ds1.dummyldapservice.com:636"),
			resource.TestCheckResourceAttr(resourceFullName, "validate_tls_certificates", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "vendor", "Microsoft Active Directory"),
			resource.TestCheckResourceAttr(resourceFullName, "user_types.%", "2"),

			resource.TestCheckResourceAttr(resourceFullName, "user_types.User Set 2.password_authority", "PING_ONE"),
			resource.TestCheckResourceAttr(resourceFullName, "user_types.User Set 2.search_base_dn", "ou=users,dc=example,dc=com"),
			resource.TestCheckResourceAttr(resourceFullName, "user_types.User Set 2.user_link_attributes.#", "3"),
			resource.TestCheckResourceAttr(resourceFullName, "user_types.User Set 2.user_link_attributes.0", "objectGUID"),
			resource.TestCheckResourceAttr(resourceFullName, "user_types.User Set 2.user_link_attributes.1", "dn"),
			resource.TestCheckResourceAttr(resourceFullName, "user_types.User Set 2.user_link_attributes.2", "objectSid"),
			resource.TestCheckResourceAttr(resourceFullName, "user_types.User Set 2.new_user_lookup.ldap_filter_pattern", "(|(uid=${identifier})(mail=${identifier}))"),
			resource.TestCheckResourceAttr(resourceFullName, "user_types.User Set 2.new_user_lookup.attribute_mappings.#", "3"),
			resource.TestCheckResourceAttr(resourceFullName, "user_types.User Set 2.allow_password_changes", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "user_types.User Set 2.update_user_on_successful_authentication", "false"),

			resource.TestCheckResourceAttr(resourceFullName, "user_types.User Set 1.password_authority", "LDAP"),
			resource.TestCheckResourceAttr(resourceFullName, "user_types.User Set 1.search_base_dn", "ou=users1,dc=example,dc=com"),
			resource.TestCheckResourceAttr(resourceFullName, "user_types.User Set 1.user_link_attributes.#", "2"),
			resource.TestCheckResourceAttr(resourceFullName, "user_types.User Set 1.user_link_attributes.0", "objectGUID"),
			resource.TestCheckResourceAttr(resourceFullName, "user_types.User Set 1.user_link_attributes.1", "objectSid"),
			resource.TestCheckNoResourceAttr(resourceFullName, "user_types.User Set 1.new_user_lookup"),
			resource.TestCheckResourceAttr(resourceFullName, "user_types.User Set 1.allow_password_changes", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "user_types.User Set 1.update_user_on_successful_authentication", "false"),
		),
	}

	minimalStep := resource.TestStep{
		Config: testAccGatewayConfig_LDAPMinimal(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "name", name),
			resource.TestCheckNoResourceAttr(resourceFullName, "description"),
			resource.TestCheckResourceAttr(resourceFullName, "enabled", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "type", "LDAP"),
			resource.TestCheckResourceAttr(resourceFullName, "bind_dn", "ou=test,dc=example,dc=com"),
			resource.TestCheckResourceAttr(resourceFullName, "bind_password", "dummyPasswordValue"),
			resource.TestCheckResourceAttr(resourceFullName, "connection_security", "None"),
			resource.TestCheckNoResourceAttr(resourceFullName, "kerberos.service_account_upn"),
			resource.TestCheckNoResourceAttr(resourceFullName, "kerberos.service_account_password"),
			resource.TestCheckNoResourceAttr(resourceFullName, "kerberos.retain_previous_credentials_mins"),
			resource.TestCheckResourceAttr(resourceFullName, "servers.#", "3"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "servers.*", "ds2.dummyldapservice.com:389"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "servers.*", "ds3.dummyldapservice.com:389"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "servers.*", "ds1.dummyldapservice.com:389"),
			resource.TestCheckResourceAttr(resourceFullName, "validate_tls_certificates", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "vendor", "Microsoft Active Directory"),
			resource.TestCheckResourceAttr(resourceFullName, "user_types.%", "0"),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Gateway_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			fullStep1,
			{
				Config:  testAccGatewayConfig_LDAPFull1(resourceName, name),
				Destroy: true,
			},
			// Minimal
			minimalStep,
			{
				Config:  testAccGatewayConfig_LDAPMinimal(resourceName, name),
				Destroy: true,
			},
			// Full Change
			fullStep1,
			fullStep2,
			fullStep1,
			// Change
			fullStep1,
			minimalStep,
			fullStep1,
			// Test importing the resource
			{
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("Resource Not found: %s", resourceFullName)
						}

						return fmt.Sprintf("%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"bind_password",
					"kerberos.service_account_password",
				},
			},
		},
	})
}

func TestAccGateway_RADIUS(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_gateway.%s", resourceName)

	name := resourceName

	fullStep := resource.TestStep{
		Config: testAccGatewayConfig_RADIUSFull(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "name", name),
			resource.TestCheckNoResourceAttr(resourceFullName, "description"),
			resource.TestCheckResourceAttr(resourceFullName, "enabled", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "type", "RADIUS"),
			resource.TestMatchResourceAttr(resourceFullName, "radius_davinci_policy_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "radius_default_shared_secret", "sharedsecret123"),
			resource.TestCheckResourceAttr(resourceFullName, "radius_clients.#", "2"),
			resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "radius_clients.*", map[string]string{
				"ip":            "127.0.0.1",
				"shared_secret": "sharedsecret123-1",
			}),
			resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "radius_clients.*", map[string]string{
				"ip":            "127.0.0.2",
				"shared_secret": "sharedsecret123-2",
			}),
			resource.TestCheckResourceAttr(resourceFullName, "radius_network_policy_server.ip", "10.1.1.1"),
			resource.TestCheckResourceAttr(resourceFullName, "radius_network_policy_server.port", "5000"),
		),
	}

	minimalStep := resource.TestStep{
		Config: testAccGatewayConfig_RADIUSDefaultSharedSecret(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "name", name),
			resource.TestCheckNoResourceAttr(resourceFullName, "description"),
			resource.TestCheckResourceAttr(resourceFullName, "enabled", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "type", "RADIUS"),
			resource.TestMatchResourceAttr(resourceFullName, "radius_davinci_policy_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "radius_default_shared_secret", "sharedsecret123"),
			resource.TestCheckResourceAttr(resourceFullName, "radius_clients.#", "1"),

			resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "radius_clients.*", map[string]string{
				"ip":            "127.0.0.3",
				"shared_secret": "",
			}),
			resource.TestCheckNoResourceAttr(resourceFullName, "radius_network_policy_server.ip"),
			resource.TestCheckNoResourceAttr(resourceFullName, "radius_network_policy_server.port"),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Gateway_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			fullStep,
			{
				Config:  testAccGatewayConfig_RADIUSFull(resourceName, name),
				Destroy: true,
			},
			// Minimal
			minimalStep,
			{
				Config:  testAccGatewayConfig_RADIUSDefaultSharedSecret(resourceName, name),
				Destroy: true,
			},
			// Change
			fullStep,
			minimalStep,
			fullStep,
			// Test importing the resource
			{
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("Resource Not found: %s", resourceFullName)
						}

						return fmt.Sprintf("%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccGateway_RADIUSSharedSecrets(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_gateway.%s", resourceName)

	name := resourceName

	defaultSecretStep := resource.TestStep{
		Config: testAccGatewayConfig_RADIUSDefaultSharedSecret(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "radius_default_shared_secret", "sharedsecret123"),
			resource.TestCheckResourceAttr(resourceFullName, "radius_clients.#", "1"),

			resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "radius_clients.*", map[string]string{
				"ip": "127.0.0.3",
			}),
		),
	}

	perClientSecretStep := resource.TestStep{
		Config: testAccGatewayConfig_RADIUSSharedSecretPerClient(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckNoResourceAttr(resourceFullName, "radius_default_shared_secret"),
			resource.TestCheckResourceAttr(resourceFullName, "radius_clients.#", "1"),

			resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "radius_clients.*", map[string]string{
				"ip":            "127.0.0.3",
				"shared_secret": "sharedsecret123-3",
			}),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Gateway_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			defaultSecretStep,
			{
				Config:  testAccGatewayConfig_RADIUSSharedSecretPerClient(resourceName, name),
				Destroy: true,
			},
			// Minimal
			perClientSecretStep,
			{
				Config:  testAccGatewayConfig_RADIUSDefaultSharedSecret(resourceName, name),
				Destroy: true,
			},
			// Change
			defaultSecretStep,
			perClientSecretStep,
			defaultSecretStep,
			// Invalid shared secret
			{
				Config:      testAccGatewayConfig_RADIUSInvalidSecretCombination(resourceName, name),
				ExpectError: regexp.MustCompile(`Invalid Value`),
			},
		},
	})
}

func TestAccGateway_BadParameter(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_gateway.%s", resourceName)

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
				Config:      testAccGatewayConfig_BadParameter(resourceName, name),
				ExpectError: regexp.MustCompile("Invalid argument combination"),
			},
			// Configure for import testing
			{
				Config: testAccGatewayConfig_LDAPMinimal(resourceName, name),
			},
			// Errors
			{
				ResourceName: resourceFullName,
				ImportState:  true,
				ExpectError:  regexp.MustCompile(`Unexpected Import Identifier`),
			},
			{
				ResourceName:  resourceFullName,
				ImportStateId: "/",
				ImportState:   true,
				ExpectError:   regexp.MustCompile(`Unexpected Import Identifier`),
			},
			{
				ResourceName:  resourceFullName,
				ImportStateId: "badformat/badformat",
				ImportState:   true,
				ExpectError:   regexp.MustCompile(`Unexpected Import Identifier`),
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

  type = "PING_FEDERATE"
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

  type = "PING_FEDERATE"
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccGatewayConfig_Minimal(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_gateway" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = false

  type = "PING_FEDERATE"
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccGatewayConfig_PingFederate(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_gateway" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = false

  type = "PING_FEDERATE"
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccGatewayConfig_APIGateway(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_gateway" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = false

  type = "API_GATEWAY_INTEGRATION"
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccGatewayConfig_Intelligence(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_gateway" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = false

  type = "PING_INTELLIGENCE"
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccGatewayConfig_LDAPFull1(resourceName, name string) string {
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
      update_user_on_successful_authentication = true

      user_link_attributes = ["objectGUID", "objectSid"]

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

      allow_password_changes                   = true
      update_user_on_successful_authentication = false

      user_link_attributes = ["objectGUID", "dn", "objectSid"]

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
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccGatewayConfig_LDAPFull2(resourceName, name string) string {
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
    },
    "User Set 2" = {
      password_authority = "PING_ONE"
      search_base_dn     = "ou=users,dc=example,dc=com"

      user_link_attributes = ["objectGUID", "dn", "objectSid"]

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
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccGatewayConfig_LDAPMinimal(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_gateway" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = false

  type = "LDAP"

  bind_dn       = "ou=test,dc=example,dc=com"
  bind_password = "dummyPasswordValue"

  vendor = "Microsoft Active Directory"

  servers = [
    "ds1.dummyldapservice.com:389",
    "ds3.dummyldapservice.com:389",
    "ds2.dummyldapservice.com:389",
  ]

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccGatewayConfig_RADIUSFull(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_gateway" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = false
  type           = "RADIUS"

  radius_default_shared_secret = "sharedsecret123"

  radius_davinci_policy_id = "ee8470a2-8161-4d76-a7af-a8505a2da084" // dummy ID

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

  radius_network_policy_server = {
    ip   = "10.1.1.1"
    port = 5000
  }

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccGatewayConfig_RADIUSDefaultSharedSecret(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_gateway" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = false
  type           = "RADIUS"

  radius_default_shared_secret = "sharedsecret123"

  radius_davinci_policy_id = "ee8470a2-8161-4d76-a7af-a8505a2da085" // dummy ID

  radius_clients = [
    {
      ip = "127.0.0.3"
    }
  ]

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccGatewayConfig_RADIUSSharedSecretPerClient(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_gateway" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = false
  type           = "RADIUS"

  radius_davinci_policy_id = "ee8470a2-8161-4d76-a7af-a8505a2da085" // dummy ID

  radius_clients = [
    {
      ip            = "127.0.0.3"
      shared_secret = "sharedsecret123-3"
    }
  ]

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccGatewayConfig_RADIUSInvalidSecretCombination(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_gateway" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = false
  type           = "RADIUS"

  radius_davinci_policy_id = "ee8470a2-8161-4d76-a7af-a8505a2da085" // dummy ID

  radius_clients = [
    {
      ip = "127.0.0.3"
    }
  ]

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccGatewayConfig_BadParameter(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_gateway" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = false

  type = "PING_FEDERATE"

  bind_dn       = "ou=test,dc=example,dc=com"
  bind_password = "dummyPasswordValue"

  vendor = "Microsoft Active Directory"

  servers = [
    "ds1.dummyldapservice.com:636",
    "ds3.dummyldapservice.com:636",
    "ds2.dummyldapservice.com:636",
  ]
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}
