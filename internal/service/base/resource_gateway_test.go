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
			// Replan after removal preconfig
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
					resource.TestCheckResourceAttr(resourceFullName, "description", ""),
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
					resource.TestCheckResourceAttr(resourceFullName, "description", ""),
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
					resource.TestCheckResourceAttr(resourceFullName, "description", ""),
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
				ImportStateVerifyIgnore: []string{
					"connection_security",
					"validate_tls_certificates",
				},
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
				ImportStateVerifyIgnore: []string{
					"connection_security",
					"validate_tls_certificates",
				},
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
				ImportStateVerifyIgnore: []string{
					"connection_security",
					"validate_tls_certificates",
				},
			},
		},
	})
}

func TestAccGateway_LDAP(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_gateway.%s", resourceName)

	name := resourceName

	fullStep := resource.TestStep{
		Config: testAccGatewayConfig_LDAPFull(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "name", name),
			resource.TestCheckResourceAttr(resourceFullName, "description", ""),
			resource.TestCheckResourceAttr(resourceFullName, "enabled", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "type", "LDAP"),
			resource.TestCheckResourceAttr(resourceFullName, "bind_dn", "ou=test,dc=example,dc=com"),
			resource.TestCheckResourceAttr(resourceFullName, "bind_password", "dummyPasswordValue"),
			resource.TestCheckResourceAttr(resourceFullName, "connection_security", "TLS"),
			resource.TestCheckResourceAttr(resourceFullName, "kerberos_service_account_upn", "username@domainname"),
			resource.TestCheckResourceAttr(resourceFullName, "kerberos_service_account_password", "dummyKerberosPasswordValue"),
			resource.TestCheckResourceAttr(resourceFullName, "kerberos_retain_previous_credentials_mins", "20"),
			resource.TestCheckResourceAttr(resourceFullName, "servers.#", "3"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "servers.*", "ds2.dummyldapservice.com:636"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "servers.*", "ds3.dummyldapservice.com:636"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "servers.*", "ds1.dummyldapservice.com:636"),
			resource.TestCheckResourceAttr(resourceFullName, "validate_tls_certificates", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "vendor", "Microsoft Active Directory"),
			resource.TestCheckResourceAttr(resourceFullName, "user_type.#", "2"),

			resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "user_type.*", map[string]string{
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

			/*
				resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "user_type.0.user_migration.0.attribute_mapping.*", map[string]string{
					"name":  "username",
					"value": "${ldapAttributes.uid}",
				}),
				resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "user_type.0.user_migration.0.attribute_mapping.*", map[string]string{
					"name":  "email",
					"value": "${ldapAttributes.mail}",
				}),
				resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "user_type.0.user_migration.0.attribute_mapping.*", map[string]string{
					"name":  "name.family",
					"value": "${ldapAttributes.sn}",
				}),
			*/
			resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "user_type.*", map[string]string{
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
			/*
				resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "user_type.1.user_migration.0.attribute_mapping.*", map[string]string{
					"name":  "username",
					"value": "${ldapAttributes.uid}",
				}),
				resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "user_type.1.user_migration.0.attribute_mapping.*", map[string]string{
					"name":  "email",
					"value": "${ldapAttributes.mail}",
				}),
			*/
		),
	}

	minimalStep := resource.TestStep{
		Config: testAccGatewayConfig_LDAPMinimal(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "name", name),
			resource.TestCheckResourceAttr(resourceFullName, "description", ""),
			resource.TestCheckResourceAttr(resourceFullName, "enabled", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "type", "LDAP"),
			resource.TestCheckResourceAttr(resourceFullName, "bind_dn", "ou=test,dc=example,dc=com"),
			resource.TestCheckResourceAttr(resourceFullName, "bind_password", "dummyPasswordValue"),
			resource.TestCheckResourceAttr(resourceFullName, "connection_security", "None"),
			resource.TestCheckResourceAttr(resourceFullName, "kerberos_service_account_upn", ""),
			resource.TestCheckResourceAttr(resourceFullName, "kerberos_service_account_password", ""),
			resource.TestCheckResourceAttr(resourceFullName, "kerberos_retain_previous_credentials_mins", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "servers.#", "3"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "servers.*", "ds2.dummyldapservice.com:389"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "servers.*", "ds3.dummyldapservice.com:389"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "servers.*", "ds1.dummyldapservice.com:389"),
			resource.TestCheckResourceAttr(resourceFullName, "validate_tls_certificates", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "vendor", "PingDirectory"),
			resource.TestCheckResourceAttr(resourceFullName, "user_type.#", "0"),
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
				Config:  testAccGatewayConfig_LDAPFull(resourceName, name),
				Destroy: true,
			},
			// Minimal
			minimalStep,
			{
				Config:  testAccGatewayConfig_LDAPMinimal(resourceName, name),
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
				ImportStateVerifyIgnore: []string{
					"bind_password",
					"kerberos_service_account_password",
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
			resource.TestCheckResourceAttr(resourceFullName, "description", ""),
			resource.TestCheckResourceAttr(resourceFullName, "enabled", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "type", "RADIUS"),
			resource.TestMatchResourceAttr(resourceFullName, "radius_davinci_policy_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "radius_default_shared_secret", "sharedsecret123"),
			resource.TestCheckResourceAttr(resourceFullName, "radius_client.#", "2"),

			resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "radius_client.*", map[string]string{
				"ip":            "127.0.0.1",
				"shared_secret": "sharedsecret123-1",
			}),
			resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "radius_client.*", map[string]string{
				"ip":            "127.0.0.2",
				"shared_secret": "sharedsecret123-2",
			}),
		),
	}

	minimalStep := resource.TestStep{
		Config: testAccGatewayConfig_RADIUSDefaultSharedSecret(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "name", name),
			resource.TestCheckResourceAttr(resourceFullName, "description", ""),
			resource.TestCheckResourceAttr(resourceFullName, "enabled", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "type", "RADIUS"),
			resource.TestMatchResourceAttr(resourceFullName, "radius_davinci_policy_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "radius_default_shared_secret", "sharedsecret123"),
			resource.TestCheckResourceAttr(resourceFullName, "radius_client.#", "1"),

			resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "radius_client.*", map[string]string{
				"ip":            "127.0.0.3",
				"shared_secret": "",
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
				ImportStateVerifyIgnore: []string{
					"connection_security",
					"validate_tls_certificates",
				},
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
			resource.TestCheckResourceAttr(resourceFullName, "radius_client.#", "1"),

			resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "radius_client.*", map[string]string{
				"ip":            "127.0.0.3",
				"shared_secret": "",
			}),
		),
	}

	perClientSecretStep := resource.TestStep{
		Config: testAccGatewayConfig_RADIUSSharedSecretPerClient(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "radius_default_shared_secret", ""),
			resource.TestCheckResourceAttr(resourceFullName, "radius_client.#", "1"),

			resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "radius_client.*", map[string]string{
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
				ExpectError: regexp.MustCompile(`RadiusClient\[127\.0\.0\.3\] shared secret cannot be empty, if default shared secret is empty\.`),
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
				ExpectError: regexp.MustCompile("Unexpected parameter bind_dn for PING_FEDERATE gateway type"),
			},
			// Configure for import testing
			{
				Config: testAccGatewayConfig_LDAPMinimal(resourceName, name),
			},
			// Errors
			{
				ResourceName: resourceFullName,
				ImportState:  true,
				ExpectError:  regexp.MustCompile(`Invalid import ID specified \(".*"\).  The ID should be in the format "environment_id/gateway_id" and must match regex: .*`),
			},
			{
				ResourceName:  resourceFullName,
				ImportStateId: "/",
				ImportState:   true,
				ExpectError:   regexp.MustCompile(`Invalid import ID specified \(".*"\).  The ID should be in the format "environment_id/gateway_id" and must match regex: .*`),
			},
			{
				ResourceName:  resourceFullName,
				ImportStateId: "badformat/badformat",
				ImportState:   true,
				ExpectError:   regexp.MustCompile(`Invalid import ID specified \(".*"\).  The ID should be in the format "environment_id/gateway_id" and must match regex: .*`),
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

func testAccGatewayConfig_LDAPFull(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_population" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "pingone_gateway" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = false
  type           = "LDAP"

  bind_dn       = "ou=test,dc=example,dc=com"
  bind_password = "dummyPasswordValue"

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
    // id = "59e24997-f829-4206-b1b7-9b6a8a25c0b4"
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
    // id = "59e24997-f829-4206-b1b7-9b6a8a25c0b3"
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

  vendor = "PingDirectory"

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

  radius_client {
    ip            = "127.0.0.1"
    shared_secret = "sharedsecret123-1"
  }

  radius_client {
    ip            = "127.0.0.2"
    shared_secret = "sharedsecret123-2"
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

  radius_client {
    ip = "127.0.0.3"
  }

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

  radius_client {
    ip            = "127.0.0.3"
    shared_secret = "sharedsecret123-3"
  }

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

  radius_client {
    ip = "127.0.0.3"
  }

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
