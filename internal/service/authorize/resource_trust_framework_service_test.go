package authorize_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/authorize"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/base"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccTrustFrameworkService_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authorize_trust_framework_service.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	var serviceID, environmentID string

	var p1Client *client.Client
	var ctx = context.Background()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckFeatureFlag(t, acctest.ENUMFEATUREFLAG_AUTHORIZEPMTF)

			p1Client = acctest.PreCheckTestClient(ctx, t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             authorize.TrustFrameworkService_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccTrustFrameworkServiceConfig_Minimal(resourceName, name),
				Check:  authorize.TrustFrameworkService_GetIDs(resourceFullName, &environmentID, &serviceID),
			},
			{
				PreConfig: func() {
					authorize.TrustFrameworkService_RemovalDrift_PreConfig(ctx, p1Client.API.AuthorizeAPIClient, t, environmentID, serviceID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			// Test removal of the environment
			{
				Config: testAccTrustFrameworkServiceConfig_NewEnv(environmentName, licenseID, resourceName, name),
				Check:  authorize.TrustFrameworkService_GetIDs(resourceFullName, &environmentID, &serviceID),
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

func TestAccTrustFrameworkService_NewEnv(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authorize_trust_framework_service.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckFeatureFlag(t, acctest.ENUMFEATUREFLAG_AUTHORIZEPMTF)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             authorize.TrustFrameworkService_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccTrustFrameworkServiceConfig_NewEnv(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
				),
			},
		},
	})
}

func TestAccTrustFrameworkService_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authorize_trust_framework_service.%s", resourceName)

	name := resourceName

	fullCheck := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckResourceAttr(resourceFullName, "name", name),
		resource.TestCheckResourceAttr(resourceFullName, "description", "Test application service full"),
		resource.TestCheckResourceAttr(resourceFullName, "full_name", fmt.Sprintf("%[1]s-parent.%[1]s", name)),
		resource.TestMatchResourceAttr(resourceFullName, "parent.id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckResourceAttr(resourceFullName, "type", "SERVICE"),
		resource.TestCheckResourceAttr(resourceFullName, "service_type", "NONE"),
		resource.TestCheckResourceAttr(resourceFullName, "cache_settings.ttl_seconds", "300"),
		resource.TestMatchResourceAttr(resourceFullName, "version", verify.P1ResourceIDRegexpFullString),
	)

	minimalCheck := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckResourceAttr(resourceFullName, "name", name),
		resource.TestCheckNoResourceAttr(resourceFullName, "description"),
		resource.TestCheckResourceAttr(resourceFullName, "full_name", name),
		resource.TestCheckNoResourceAttr(resourceFullName, "parent"),
		resource.TestCheckResourceAttr(resourceFullName, "type", "SERVICE"),
		resource.TestCheckResourceAttr(resourceFullName, "service_type", "NONE"),
		resource.TestCheckNoResourceAttr(resourceFullName, "cache_settings"),
		resource.TestMatchResourceAttr(resourceFullName, "version", verify.P1ResourceIDRegexpFullString),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckFeatureFlag(t, acctest.ENUMFEATUREFLAG_AUTHORIZEPMTF)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             authorize.TrustFrameworkService_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			{
				Config: testAccTrustFrameworkServiceConfig_Full(resourceName, name),
				Check:  fullCheck,
			},
			{
				Config:  testAccTrustFrameworkServiceConfig_Full(resourceName, name),
				Destroy: true,
			},
			// Minimal
			{
				Config: testAccTrustFrameworkServiceConfig_Minimal(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config:  testAccTrustFrameworkServiceConfig_Minimal(resourceName, name),
				Destroy: true,
			},
			// Change
			{
				Config: testAccTrustFrameworkServiceConfig_Full(resourceName, name),
				Check:  fullCheck,
			},
			{
				Config: testAccTrustFrameworkServiceConfig_Minimal(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config: testAccTrustFrameworkServiceConfig_Full(resourceName, name),
				Check:  fullCheck,
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

func TestAccTrustFrameworkService_Service_HTTP(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authorize_trust_framework_service.%s", resourceName)

	name := resourceName

	fullCheck1 := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "service_type", "HTTP"),
		resource.TestCheckResourceAttr(resourceFullName, "cache_settings.ttl_seconds", "300"),
		resource.TestCheckResourceAttr(resourceFullName, "processor.name", fmt.Sprintf("%s Test chain processor", name)),
		resource.TestCheckResourceAttr(resourceFullName, "processor.type", "CHAIN"),
		resource.TestCheckResourceAttr(resourceFullName, "processor.processors.#", "3"), // processors tested in processors_test
		resource.TestCheckResourceAttr(resourceFullName, "service_settings.authentication.type", "BASIC"),
		resource.TestCheckNoResourceAttr(resourceFullName, "service_settings.authentication.client_id"),
		resource.TestCheckNoResourceAttr(resourceFullName, "service_settings.authentication.client_secret"),
		resource.TestMatchResourceAttr(resourceFullName, "service_settings.authentication.name.id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "service_settings.authentication.password.id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckNoResourceAttr(resourceFullName, "service_settings.authentication.scope"),
		resource.TestCheckNoResourceAttr(resourceFullName, "service_settings.authentication.token"),
		resource.TestCheckNoResourceAttr(resourceFullName, "service_settings.authentication.token_endpoint"),
		resource.TestCheckResourceAttr(resourceFullName, "service_settings.body", "{\"data\": {\"item1\": \"value1\", \"item2\": \"value2\", \"item3\": true}}"),
		resource.TestCheckNoResourceAttr(resourceFullName, "service_settings.capability"),
		resource.TestCheckNoResourceAttr(resourceFullName, "service_settings.channel"),
		resource.TestCheckNoResourceAttr(resourceFullName, "service_settings.code"),
		resource.TestCheckResourceAttr(resourceFullName, "service_settings.content_type", "application/json"),
		resource.TestCheckResourceAttr(resourceFullName, "service_settings.headers.#", "3"),
		resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "service_settings.headers.*", map[string]string{
			"key":         "my_custom_header",
			"value.type":  "CONSTANT",
			"value.value": "my_custom_value",
		}),
		resource.TestMatchTypeSetElemNestedAttrs(resourceFullName, "service_settings.headers.*", map[string]*regexp.Regexp{
			"key":                regexp.MustCompile(`^my_custom_header2$`),
			"value.type":         regexp.MustCompile(`^ATTRIBUTE$`),
			"value.attribute.id": verify.P1ResourceIDRegexpFullString,
		}),
		resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "service_settings.headers.*", map[string]string{
			"key":         "my_custom_header3",
			"value.type":  "CONSTANT",
			"value.value": "my_custom_value3",
		}),
		resource.TestCheckNoResourceAttr(resourceFullName, "service_settings.input_mappings"),
		resource.TestCheckResourceAttr(resourceFullName, "service_settings.maximum_concurrent_requests", "6"),
		resource.TestCheckResourceAttr(resourceFullName, "service_settings.maximum_requests_per_second", "10"),
		resource.TestCheckNoResourceAttr(resourceFullName, "service_settings.schema_version"),
		resource.TestCheckResourceAttr(resourceFullName, "service_settings.timeout_milliseconds", "2000"),
		resource.TestCheckResourceAttr(resourceFullName, "service_settings.tls_settings.tls_validation_type", "NONE"),
		resource.TestCheckResourceAttr(resourceFullName, "service_settings.url", "https://pingidentity.com"),
		resource.TestCheckResourceAttr(resourceFullName, "service_settings.verb", "POST"),
		resource.TestCheckResourceAttr(resourceFullName, "value_type.type", "JSON"),
	)

	fullCheck2 := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "service_type", "HTTP"),
		resource.TestCheckResourceAttr(resourceFullName, "cache_settings.ttl_seconds", "400"),
		resource.TestCheckResourceAttr(resourceFullName, "processor.name", fmt.Sprintf("%s Test chain processor 1", name)),
		resource.TestCheckResourceAttr(resourceFullName, "processor.type", "JSON_PATH"),
		resource.TestCheckResourceAttr(resourceFullName, "service_settings.authentication.type", "CLIENT_CREDENTIALS"),
		resource.TestCheckResourceAttr(resourceFullName, "service_settings.authentication.client_id", "test_client_id"),
		resource.TestMatchResourceAttr(resourceFullName, "service_settings.authentication.client_secret.id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckNoResourceAttr(resourceFullName, "service_settings.authentication.name"),
		resource.TestCheckNoResourceAttr(resourceFullName, "service_settings.authentication.password"),
		resource.TestCheckResourceAttr(resourceFullName, "service_settings.authentication.scope", "scope1 scope2 scope3"),
		resource.TestCheckNoResourceAttr(resourceFullName, "service_settings.authentication.token"),
		resource.TestCheckResourceAttr(resourceFullName, "service_settings.authentication.token_endpoint", "https://auth.pingidentity.com/example"),
		resource.TestCheckResourceAttr(resourceFullName, "service_settings.body", "test body"),
		resource.TestCheckNoResourceAttr(resourceFullName, "service_settings.capability"),
		resource.TestCheckNoResourceAttr(resourceFullName, "service_settings.channel"),
		resource.TestCheckNoResourceAttr(resourceFullName, "service_settings.code"),
		resource.TestCheckResourceAttr(resourceFullName, "service_settings.content_type", "application/text"),
		resource.TestCheckResourceAttr(resourceFullName, "service_settings.headers.#", "2"),
		resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "service_settings.headers.*", map[string]string{
			"key":         "my_custom_header",
			"value.type":  "CONSTANT",
			"value.value": "my_custom_value",
		}),
		resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "service_settings.headers.*", map[string]string{
			"key":         "my_custom_header3",
			"value.type":  "CONSTANT",
			"value.value": "my_custom_value3",
		}),
		resource.TestCheckNoResourceAttr(resourceFullName, "service_settings.input_mappings"),
		resource.TestCheckResourceAttr(resourceFullName, "service_settings.maximum_concurrent_requests", "4"),
		resource.TestCheckResourceAttr(resourceFullName, "service_settings.maximum_requests_per_second", "11"),
		resource.TestCheckNoResourceAttr(resourceFullName, "service_settings.schema_version"),
		resource.TestCheckResourceAttr(resourceFullName, "service_settings.timeout_milliseconds", "2500"),
		resource.TestCheckResourceAttr(resourceFullName, "service_settings.tls_settings.tls_validation_type", "DEFAULT"),
		resource.TestCheckResourceAttr(resourceFullName, "service_settings.url", "https://pingidentity.com/test"),
		resource.TestCheckResourceAttr(resourceFullName, "service_settings.verb", "PUT"),
		resource.TestCheckResourceAttr(resourceFullName, "value_type.type", "XML"),
	)

	minimalCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "service_type", "HTTP"),
		resource.TestCheckNoResourceAttr(resourceFullName, "cache_settings"),
		resource.TestCheckNoResourceAttr(resourceFullName, "processor"),
		resource.TestCheckResourceAttr(resourceFullName, "service_settings.authentication.type", "NONE"),
		resource.TestCheckNoResourceAttr(resourceFullName, "service_settings.authentication.client_id"),
		resource.TestCheckNoResourceAttr(resourceFullName, "service_settings.authentication.client_secret"),
		resource.TestCheckNoResourceAttr(resourceFullName, "service_settings.authentication.username"),
		resource.TestCheckNoResourceAttr(resourceFullName, "service_settings.authentication.password"),
		resource.TestCheckNoResourceAttr(resourceFullName, "service_settings.authentication.scope"),
		resource.TestCheckNoResourceAttr(resourceFullName, "service_settings.authentication.token"),
		resource.TestCheckNoResourceAttr(resourceFullName, "service_settings.authentication.token_endpoint"),
		resource.TestCheckNoResourceAttr(resourceFullName, "service_settings.body"),
		resource.TestCheckNoResourceAttr(resourceFullName, "service_settings.capability"),
		resource.TestCheckNoResourceAttr(resourceFullName, "service_settings.channel"),
		resource.TestCheckNoResourceAttr(resourceFullName, "service_settings.code"),
		resource.TestCheckNoResourceAttr(resourceFullName, "service_settings.content_type"),
		resource.TestCheckNoResourceAttr(resourceFullName, "service_settings.headers"),
		resource.TestCheckNoResourceAttr(resourceFullName, "service_settings.input_mappings"),
		resource.TestCheckResourceAttr(resourceFullName, "service_settings.maximum_concurrent_requests", "6"),
		resource.TestCheckResourceAttr(resourceFullName, "service_settings.maximum_requests_per_second", "10"),
		resource.TestCheckNoResourceAttr(resourceFullName, "service_settings.schema_version"),
		resource.TestCheckResourceAttr(resourceFullName, "service_settings.timeout_milliseconds", "2000"),
		resource.TestCheckResourceAttr(resourceFullName, "service_settings.tls_settings.tls_validation_type", "DEFAULT"),
		resource.TestCheckResourceAttr(resourceFullName, "service_settings.url", "https://pingidentity.com"),
		resource.TestCheckResourceAttr(resourceFullName, "service_settings.verb", "GET"),
		resource.TestCheckResourceAttr(resourceFullName, "value_type.type", "STRING"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckFeatureFlag(t, acctest.ENUMFEATUREFLAG_AUTHORIZEPMTF)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             authorize.TrustFrameworkService_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			{
				Config: testAccTrustFrameworkServiceConfig_Service_HTTP_Full1(resourceName, name),
				Check:  fullCheck1,
			},
			{
				Config:  testAccTrustFrameworkServiceConfig_Service_HTTP_Full1(resourceName, name),
				Destroy: true,
			},
			// Minimal
			{
				Config: testAccTrustFrameworkServiceConfig_Service_HTTP_Minimal(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config:  testAccTrustFrameworkServiceConfig_Service_HTTP_Minimal(resourceName, name),
				Destroy: true,
			},
			// Change
			{
				Config: testAccTrustFrameworkServiceConfig_Service_HTTP_Full1(resourceName, name),
				Check:  fullCheck1,
			},
			{
				Config: testAccTrustFrameworkServiceConfig_Service_HTTP_Full2(resourceName, name),
				Check:  fullCheck2,
			},
			{
				Config: testAccTrustFrameworkServiceConfig_Service_HTTP_Minimal(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config: testAccTrustFrameworkServiceConfig_Service_HTTP_Full1(resourceName, name),
				Check:  fullCheck1,
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

func TestAccTrustFrameworkService_Service_Connector(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authorize_trust_framework_service.%s", resourceName)

	name := resourceName

	fullCheck1 := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "service_type", "CONNECTOR"),
		resource.TestCheckResourceAttr(resourceFullName, "cache_settings.ttl_seconds", "300"),
		resource.TestCheckResourceAttr(resourceFullName, "processor.name", fmt.Sprintf("%s Test chain processor", name)),
		resource.TestCheckResourceAttr(resourceFullName, "processor.type", "CHAIN"),
		resource.TestCheckResourceAttr(resourceFullName, "processor.processors.#", "3"), // processors tested in processors_test
		resource.TestCheckNoResourceAttr(resourceFullName, "service_settings.authentication"),
		resource.TestCheckNoResourceAttr(resourceFullName, "service_settings.body"),
		resource.TestCheckResourceAttr(resourceFullName, "service_settings.capability", "createRiskEvaluation"),
		resource.TestCheckResourceAttr(resourceFullName, "service_settings.channel", "AUTHORIZE"),
		resource.TestCheckResourceAttr(resourceFullName, "service_settings.code", "P1_RISK"),
		resource.TestCheckNoResourceAttr(resourceFullName, "service_settings.content_type"),
		resource.TestCheckNoResourceAttr(resourceFullName, "service_settings.headers"),
		resource.TestCheckResourceAttr(resourceFullName, "service_settings.input_mappings.#", "8"),
		resource.TestMatchTypeSetElemNestedAttrs(resourceFullName, "service_settings.input_mappings.*", map[string]*regexp.Regexp{
			"property":     regexp.MustCompile(`^userId$`),
			"type":         regexp.MustCompile(`^ATTRIBUTE$`),
			"value_ref.id": verify.P1ResourceIDRegexpFullString,
		}),
		resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "service_settings.input_mappings.*", map[string]string{
			"property": "userType",
			"type":     "INPUT",
			"value":    "test2",
		}),
		resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "service_settings.input_mappings.*", map[string]string{
			"property": "ipAddress",
			"type":     "INPUT",
			"value":    "test3",
		}),
		resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "service_settings.input_mappings.*", map[string]string{
			"property": "applicationId",
			"type":     "INPUT",
			"value":    "test4",
		}),
		resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "service_settings.input_mappings.*", map[string]string{
			"property": "applicationName",
			"type":     "INPUT",
			"value":    "test5",
		}),
		resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "service_settings.input_mappings.*", map[string]string{
			"property": "sessionId",
			"type":     "INPUT",
			"value":    "test6",
		}),
		resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "service_settings.input_mappings.*", map[string]string{
			"property": "browserData",
			"type":     "INPUT",
			"value":    "test7",
		}),
		resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "service_settings.input_mappings.*", map[string]string{
			"property": "riskPolicySetId",
			"type":     "INPUT",
			"value":    "test8",
		}),
		resource.TestCheckResourceAttr(resourceFullName, "service_settings.maximum_concurrent_requests", "6"),
		resource.TestCheckResourceAttr(resourceFullName, "service_settings.maximum_requests_per_second", "10"),
		resource.TestCheckResourceAttr(resourceFullName, "service_settings.schema_version", "1"),
		resource.TestCheckResourceAttr(resourceFullName, "service_settings.timeout_milliseconds", "2000"),
		resource.TestCheckNoResourceAttr(resourceFullName, "service_settings.tls_settings"),
		resource.TestCheckNoResourceAttr(resourceFullName, "service_settings.url"),
		resource.TestCheckNoResourceAttr(resourceFullName, "service_settings.verb"),
		resource.TestCheckResourceAttr(resourceFullName, "value_type.type", "JSON"),
	)

	fullCheck2 := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "service_type", "CONNECTOR"),
		resource.TestCheckResourceAttr(resourceFullName, "cache_settings.ttl_seconds", "400"),
		resource.TestCheckResourceAttr(resourceFullName, "processor.name", fmt.Sprintf("%s Test chain processor 1", name)),
		resource.TestCheckResourceAttr(resourceFullName, "processor.type", "JSON_PATH"),
		resource.TestCheckNoResourceAttr(resourceFullName, "service_settings.authentication"),
		resource.TestCheckNoResourceAttr(resourceFullName, "service_settings.body"),
		resource.TestCheckResourceAttr(resourceFullName, "service_settings.capability", "updateRiskEvaluation"),
		resource.TestCheckResourceAttr(resourceFullName, "service_settings.channel", "AUTHORIZE"),
		resource.TestCheckResourceAttr(resourceFullName, "service_settings.code", "P1_RISK"),
		resource.TestCheckNoResourceAttr(resourceFullName, "service_settings.content_type"),
		resource.TestCheckNoResourceAttr(resourceFullName, "service_settings.headers"),
		resource.TestCheckResourceAttr(resourceFullName, "service_settings.input_mappings.#", "2"),
		resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "service_settings.input_mappings.*", map[string]string{
			"property": "completionStatus",
			"type":     "INPUT",
			"value":    "completed",
		}),
		resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "service_settings.input_mappings.*", map[string]string{
			"property": "riskId",
			"type":     "INPUT",
			"value":    "test",
		}),
		resource.TestCheckResourceAttr(resourceFullName, "service_settings.maximum_concurrent_requests", "4"),
		resource.TestCheckResourceAttr(resourceFullName, "service_settings.maximum_requests_per_second", "11"),
		resource.TestCheckResourceAttr(resourceFullName, "service_settings.schema_version", "2"),
		resource.TestCheckResourceAttr(resourceFullName, "service_settings.timeout_milliseconds", "2500"),
		resource.TestCheckNoResourceAttr(resourceFullName, "service_settings.tls_settings"),
		resource.TestCheckNoResourceAttr(resourceFullName, "service_settings.url"),
		resource.TestCheckNoResourceAttr(resourceFullName, "service_settings.verb"),
		resource.TestCheckResourceAttr(resourceFullName, "value_type.type", "XML"),
	)

	minimalCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "service_type", "CONNECTOR"),
		resource.TestCheckNoResourceAttr(resourceFullName, "cache_settings"),
		resource.TestCheckNoResourceAttr(resourceFullName, "processor"),
		resource.TestCheckNoResourceAttr(resourceFullName, "service_settings.authentication"),
		resource.TestCheckNoResourceAttr(resourceFullName, "service_settings.body"),
		resource.TestCheckResourceAttr(resourceFullName, "service_settings.capability", "createRiskEvaluation"),
		resource.TestCheckResourceAttr(resourceFullName, "service_settings.channel", "AUTHORIZE"),
		resource.TestCheckResourceAttr(resourceFullName, "service_settings.code", "P1_RISK"),
		resource.TestCheckNoResourceAttr(resourceFullName, "service_settings.content_type"),
		resource.TestCheckNoResourceAttr(resourceFullName, "service_settings.headers"),
		resource.TestCheckResourceAttr(resourceFullName, "service_settings.input_mappings.#", "3"),
		resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "service_settings.input_mappings.*", map[string]string{
			"property": "ipAddress",
			"type":     "INPUT",
			"value":    "192.168.0.1",
		}),
		resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "service_settings.input_mappings.*", map[string]string{
			"property": "userType",
			"type":     "INPUT",
			"value":    "person",
		}),
		resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "service_settings.input_mappings.*", map[string]string{
			"property": "userId",
			"type":     "INPUT",
			"value":    "test1",
		}),
		resource.TestCheckResourceAttr(resourceFullName, "service_settings.maximum_concurrent_requests", "6"),
		resource.TestCheckResourceAttr(resourceFullName, "service_settings.maximum_requests_per_second", "10"),
		resource.TestCheckResourceAttr(resourceFullName, "service_settings.schema_version", "1"),
		resource.TestCheckResourceAttr(resourceFullName, "service_settings.timeout_milliseconds", "2000"),
		resource.TestCheckNoResourceAttr(resourceFullName, "service_settings.tls_settings"),
		resource.TestCheckNoResourceAttr(resourceFullName, "service_settings.url"),
		resource.TestCheckNoResourceAttr(resourceFullName, "service_settings.verb"),
		resource.TestCheckResourceAttr(resourceFullName, "value_type.type", "STRING"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckFeatureFlag(t, acctest.ENUMFEATUREFLAG_AUTHORIZEPMTF)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             authorize.TrustFrameworkService_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			{
				Config: testAccTrustFrameworkServiceConfig_Service_Connector_Full1(resourceName, name),
				Check:  fullCheck1,
			},
			{
				Config:  testAccTrustFrameworkServiceConfig_Service_Connector_Full1(resourceName, name),
				Destroy: true,
			},
			// Minimal
			{
				Config: testAccTrustFrameworkServiceConfig_Service_Connector_Minimal(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config:  testAccTrustFrameworkServiceConfig_Service_Connector_Minimal(resourceName, name),
				Destroy: true,
			},
			// Change
			{
				Config: testAccTrustFrameworkServiceConfig_Service_Connector_Full1(resourceName, name),
				Check:  fullCheck1,
			},
			{
				Config: testAccTrustFrameworkServiceConfig_Service_Connector_Full2(resourceName, name),
				Check:  fullCheck2,
			},
			{
				Config: testAccTrustFrameworkServiceConfig_Service_Connector_Minimal(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config: testAccTrustFrameworkServiceConfig_Service_Connector_Full1(resourceName, name),
				Check:  fullCheck1,
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

func TestAccTrustFrameworkService_Service_None(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authorize_trust_framework_service.%s", resourceName)

	name := resourceName

	fullCheck1 := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "service_type", "NONE"),
		resource.TestCheckResourceAttr(resourceFullName, "cache_settings.ttl_seconds", "300"),
		resource.TestCheckNoResourceAttr(resourceFullName, "processor"),
		resource.TestCheckNoResourceAttr(resourceFullName, "service_settings"),
		resource.TestCheckNoResourceAttr(resourceFullName, "value_type"),
	)

	minimalCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "service_type", "NONE"),
		resource.TestCheckNoResourceAttr(resourceFullName, "cache_settings"),
		resource.TestCheckNoResourceAttr(resourceFullName, "processor"),
		resource.TestCheckNoResourceAttr(resourceFullName, "service_settings"),
		resource.TestCheckNoResourceAttr(resourceFullName, "value_type"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckFeatureFlag(t, acctest.ENUMFEATUREFLAG_AUTHORIZEPMTF)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             authorize.TrustFrameworkService_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			{
				Config: testAccTrustFrameworkServiceConfig_Service_None_Full(resourceName, name),
				Check:  fullCheck1,
			},
			{
				Config:  testAccTrustFrameworkServiceConfig_Service_None_Full(resourceName, name),
				Destroy: true,
			},
			// Minimal
			{
				Config: testAccTrustFrameworkServiceConfig_Service_None_Minimal(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config:  testAccTrustFrameworkServiceConfig_Service_None_Minimal(resourceName, name),
				Destroy: true,
			},
			// Change
			{
				Config: testAccTrustFrameworkServiceConfig_Service_None_Full(resourceName, name),
				Check:  fullCheck1,
			},
			{
				Config: testAccTrustFrameworkServiceConfig_Service_None_Minimal(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config: testAccTrustFrameworkServiceConfig_Service_None_Full(resourceName, name),
				Check:  fullCheck1,
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

func TestAccTrustFrameworkService_Service_Change(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authorize_trust_framework_service.%s", resourceName)

	name := resourceName

	check1 := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "service_type", "NONE"),
		resource.TestCheckResourceAttr(resourceFullName, "cache_settings.ttl_seconds", "300"),
		resource.TestCheckNoResourceAttr(resourceFullName, "processor"),
		resource.TestCheckNoResourceAttr(resourceFullName, "service_settings"),
		resource.TestCheckNoResourceAttr(resourceFullName, "value_type"),
	)

	check2 := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "service_type", "CONNECTOR"),
		resource.TestCheckResourceAttr(resourceFullName, "cache_settings.ttl_seconds", "300"),
		resource.TestCheckResourceAttr(resourceFullName, "processor.name", fmt.Sprintf("%s Test chain processor", name)),
		resource.TestCheckResourceAttr(resourceFullName, "processor.type", "CHAIN"),
		resource.TestCheckResourceAttr(resourceFullName, "processor.processors.#", "3"), // processors tested in processors_test
		resource.TestCheckNoResourceAttr(resourceFullName, "service_settings.authentication"),
		resource.TestCheckNoResourceAttr(resourceFullName, "service_settings.body"),
		resource.TestCheckResourceAttr(resourceFullName, "service_settings.capability", "createRiskEvaluation"),
		resource.TestCheckResourceAttr(resourceFullName, "service_settings.channel", "AUTHORIZE"),
		resource.TestCheckResourceAttr(resourceFullName, "service_settings.code", "P1_RISK"),
		resource.TestCheckNoResourceAttr(resourceFullName, "service_settings.content_type"),
		resource.TestCheckNoResourceAttr(resourceFullName, "service_settings.headers"),
		resource.TestCheckResourceAttr(resourceFullName, "service_settings.input_mappings.#", "8"),
		resource.TestMatchTypeSetElemNestedAttrs(resourceFullName, "service_settings.input_mappings.*", map[string]*regexp.Regexp{
			"property":     regexp.MustCompile(`^userId$`),
			"type":         regexp.MustCompile(`^ATTRIBUTE$`),
			"value_ref.id": verify.P1ResourceIDRegexpFullString,
		}),
		resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "service_settings.input_mappings.*", map[string]string{
			"property": "userType",
			"type":     "INPUT",
			"value":    "test2",
		}),
		resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "service_settings.input_mappings.*", map[string]string{
			"property": "ipAddress",
			"type":     "INPUT",
			"value":    "test3",
		}),
		resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "service_settings.input_mappings.*", map[string]string{
			"property": "applicationId",
			"type":     "INPUT",
			"value":    "test4",
		}),
		resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "service_settings.input_mappings.*", map[string]string{
			"property": "applicationName",
			"type":     "INPUT",
			"value":    "test5",
		}),
		resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "service_settings.input_mappings.*", map[string]string{
			"property": "sessionId",
			"type":     "INPUT",
			"value":    "test6",
		}),
		resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "service_settings.input_mappings.*", map[string]string{
			"property": "browserData",
			"type":     "INPUT",
			"value":    "test7",
		}),
		resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "service_settings.input_mappings.*", map[string]string{
			"property": "riskPolicySetId",
			"type":     "INPUT",
			"value":    "test8",
		}),
		resource.TestCheckResourceAttr(resourceFullName, "service_settings.maximum_concurrent_requests", "6"),
		resource.TestCheckResourceAttr(resourceFullName, "service_settings.maximum_requests_per_second", "10"),
		resource.TestCheckResourceAttr(resourceFullName, "service_settings.schema_version", "1"),
		resource.TestCheckResourceAttr(resourceFullName, "service_settings.timeout_milliseconds", "2000"),
		resource.TestCheckNoResourceAttr(resourceFullName, "service_settings.tls_settings"),
		resource.TestCheckNoResourceAttr(resourceFullName, "service_settings.url"),
		resource.TestCheckNoResourceAttr(resourceFullName, "service_settings.verb"),
		resource.TestCheckResourceAttr(resourceFullName, "value_type.type", "JSON"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckFeatureFlag(t, acctest.ENUMFEATUREFLAG_AUTHORIZEPMTF)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             authorize.TrustFrameworkService_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccTrustFrameworkServiceConfig_Service_None_Full(resourceName, name),
				Check:  check1,
			},
			{
				Config: testAccTrustFrameworkServiceConfig_Service_Connector_Full1(resourceName, name),
				Check:  check2,
			},
			{
				Config: testAccTrustFrameworkServiceConfig_Service_None_Full(resourceName, name),
				Check:  check1,
			},
		},
	})
}

func TestAccTrustFrameworkService_BadParameters(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authorize_trust_framework_service.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckFeatureFlag(t, acctest.ENUMFEATUREFLAG_AUTHORIZEPMTF)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             authorize.TrustFrameworkService_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccTrustFrameworkServiceConfig_Minimal(resourceName, name),
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

func testAccTrustFrameworkServiceConfig_NewEnv(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_authorize_trust_framework_service" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  name           = "%[3]s"
  description    = "Test application service"

  service_type = "NONE"
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccTrustFrameworkServiceConfig_Full(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_authorize_trust_framework_service" "%[2]s-parent" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s-parent"
  description    = "Test application service"

  service_type = "NONE"
}

resource "pingone_authorize_trust_framework_service" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "Test application service full"

  parent = {
    id = pingone_authorize_trust_framework_service.%[2]s-parent.id
  }

  service_type = "NONE"

  cache_settings = {
    ttl_seconds = 300
  }
}`, acctest.AuthorizePMTFSandboxEnvironment(), resourceName, name)
}

func testAccTrustFrameworkServiceConfig_Minimal(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_authorize_trust_framework_service" "%[2]s-parent" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s-parent"
  description    = "Test application service"

  service_type = "NONE"
}

resource "pingone_authorize_trust_framework_service" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  service_type = "NONE"
}`, acctest.AuthorizePMTFSandboxEnvironment(), resourceName, name)
}

func testAccTrustFrameworkServiceConfig_Service_HTTP_Full1(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_authorize_trust_framework_attribute" "%[2]s-header" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s-header"
  description    = "Test application service"

  resolvers = [
    {
      type = "CURRENT_USER_ID"
    }
  ]

  value_type = {
    type = "STRING"
  }
}

resource "pingone_authorize_trust_framework_attribute" "%[2]s-basic-auth-username" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s-basic-auth-username"
  description    = "Test application service"

  resolvers = [
    {
      type = "CONSTANT"
      value_type = {
        type = "STRING"
      }
      value = "username"
    }
  ]

  value_type = {
    type = "STRING"
  }
}

resource "pingone_authorize_trust_framework_attribute" "%[2]s-basic-auth-password" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s-basic-auth-password"
  description    = "Test application service"

  resolvers = [
    {
      type = "CONSTANT"
      value_type = {
        type = "STRING"
      }
      value = "password"
    }
  ]

  value_type = {
    type = "STRING"
  }
}

resource "pingone_authorize_trust_framework_attribute" "%[2]s-client-secret" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s-client-secret"
  description    = "Test application service"

  resolvers = [
    {
      type = "CONSTANT"
      value_type = {
        type = "STRING"
      }
      value = "clientsecret"
    }
  ]

  value_type = {
    type = "STRING"
  }
}

resource "pingone_authorize_trust_framework_service" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "Test application service"

  service_type = "HTTP"

  cache_settings = {
    ttl_seconds = 300
  }

  processor = {
    name = "%[3]s Test chain processor"
    type = "CHAIN"

    processors = [
      {
        name = "%[3]s Test chain processor 1"
        type = "JSON_PATH"

        expression = "$.data.item1"
        value_type = {
          type = "STRING"
        }
      },
      {
        name = "%[3]s Test chain processor 3"
        type = "JSON_PATH"

        expression = "$.data.item3"
        value_type = {
          type = "JSON"
        }
      },
      {
        name = "%[3]s Test chain processor 2"
        type = "JSON_PATH"

        expression = "$.data.item2"
        value_type = {
          type = "STRING"
        }
      },
    ],
  }

  service_settings = {
    authentication = {
      type = "BASIC"

      name = {
        id = pingone_authorize_trust_framework_attribute.%[2]s-basic-auth-username.id
      }

      password = {
        id = pingone_authorize_trust_framework_attribute.%[2]s-basic-auth-password.id
      }
    }

    body         = "{\"data\": {\"item1\": \"value1\", \"item2\": \"value2\", \"item3\": true}}"
    content_type = "application/json"

    headers = [
      {
        key = "my_custom_header",
        value = {
          type  = "CONSTANT"
          value = "my_custom_value"
        }
      },
      {
        key = "my_custom_header2",
        value = {
          type = "ATTRIBUTE"
          attribute = {
            id = pingone_authorize_trust_framework_attribute.%[2]s-header.id
          }
        }
      },
      {
        key = "my_custom_header3",
        value = {
          type  = "CONSTANT"
          value = "my_custom_value3"
        }
      }
    ]

    maximum_concurrent_requests = 6
    maximum_requests_per_second = 10
    timeout_milliseconds        = 2000

    tls_settings = {
      tls_validation_type = "NONE"
    }

    url  = "https://pingidentity.com"
    verb = "POST"
  }

  value_type = {
    type = "JSON"
  }
}`, acctest.AuthorizePMTFSandboxEnvironment(), resourceName, name)
}

func testAccTrustFrameworkServiceConfig_Service_HTTP_Full2(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_authorize_trust_framework_attribute" "%[2]s-header" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s-header"
  description    = "Test application service"

  resolvers = [
    {
      type = "CURRENT_USER_ID"
    }
  ]

  value_type = {
    type = "STRING"
  }
}

resource "pingone_authorize_trust_framework_attribute" "%[2]s-basic-auth-username" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s-basic-auth-username"
  description    = "Test application service"

  resolvers = [
    {
      type = "CONSTANT"
      value_type = {
        type = "STRING"
      }
      value = "username"
    }
  ]

  value_type = {
    type = "STRING"
  }
}

resource "pingone_authorize_trust_framework_attribute" "%[2]s-basic-auth-password" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s-basic-auth-password"
  description    = "Test application service"

  resolvers = [
    {
      type = "CONSTANT"
      value_type = {
        type = "STRING"
      }
      value = "password"
    }
  ]

  value_type = {
    type = "STRING"
  }
}

resource "pingone_authorize_trust_framework_attribute" "%[2]s-client-secret" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s-client-secret"
  description    = "Test application service"

  resolvers = [
    {
      type = "CONSTANT"
      value_type = {
        type = "STRING"
      }
      value = "clientsecret"
    }
  ]

  value_type = {
    type = "STRING"
  }
}

resource "pingone_authorize_trust_framework_service" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "Test application service"

  service_type = "HTTP"

  cache_settings = {
    ttl_seconds = 400
  }

  processor = {
    name = "%[3]s Test chain processor 1"
    type = "JSON_PATH"

    expression = "$.data.item1"
    value_type = {
      type = "STRING"
    }
  }

  service_settings = {
    authentication = {
      type = "CLIENT_CREDENTIALS"

      client_id = "test_client_id"

      client_secret = {
        id = pingone_authorize_trust_framework_attribute.%[2]s-client-secret.id
      }

      scope          = "scope1 scope2 scope3"
      token_endpoint = "https://auth.pingidentity.com/example"
    }

    body         = "test body"
    content_type = "application/text"

    headers = [
      {
        key = "my_custom_header",
        value = {
          type  = "CONSTANT"
          value = "my_custom_value"
        }
      },
      {
        key = "my_custom_header3",
        value = {
          type  = "CONSTANT"
          value = "my_custom_value3"
        }
      }
    ]

    maximum_concurrent_requests = 4
    maximum_requests_per_second = 11
    timeout_milliseconds        = 2500

    tls_settings = {
      tls_validation_type = "DEFAULT"
    }

    url  = "https://pingidentity.com/test"
    verb = "PUT"
  }

  value_type = {
    type = "XML"
  }
}`, acctest.AuthorizePMTFSandboxEnvironment(), resourceName, name)
}

func testAccTrustFrameworkServiceConfig_Service_HTTP_Minimal(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_authorize_trust_framework_attribute" "%[2]s-header" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s-header"
  description    = "Test application service"

  resolvers = [
    {
      type = "CURRENT_USER_ID"
    }
  ]

  value_type = {
    type = "STRING"
  }
}

resource "pingone_authorize_trust_framework_attribute" "%[2]s-basic-auth-username" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s-basic-auth-username"
  description    = "Test application service"

  resolvers = [
    {
      type = "CONSTANT"
      value_type = {
        type = "STRING"
      }
      value = "username"
    }
  ]

  value_type = {
    type = "STRING"
  }
}

resource "pingone_authorize_trust_framework_attribute" "%[2]s-basic-auth-password" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s-basic-auth-password"
  description    = "Test application service"

  resolvers = [
    {
      type = "CONSTANT"
      value_type = {
        type = "STRING"
      }
      value = "password"
    }
  ]

  value_type = {
    type = "STRING"
  }
}

resource "pingone_authorize_trust_framework_attribute" "%[2]s-client-secret" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s-client-secret"
  description    = "Test application service"

  resolvers = [
    {
      type = "CONSTANT"
      value_type = {
        type = "STRING"
      }
      value = "clientsecret"
    }
  ]

  value_type = {
    type = "STRING"
  }
}

resource "pingone_authorize_trust_framework_service" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "Test application service"

  service_type = "HTTP"

  service_settings = {
    authentication = {
      type = "NONE"
    }

    maximum_concurrent_requests = 6
    maximum_requests_per_second = 10
    timeout_milliseconds        = 2000

    tls_settings = {
      tls_validation_type = "DEFAULT"
    }

    url  = "https://pingidentity.com"
    verb = "GET"
  }

  value_type = {
    type = "STRING"
  }
}`, acctest.AuthorizePMTFSandboxEnvironment(), resourceName, name)
}

func testAccTrustFrameworkServiceConfig_Service_Connector_Full1(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_authorize_trust_framework_attribute" "%[2]s-user" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "Test application service"

  resolvers = [
    {
      type = "CURRENT_USER_ID"
    }
  ]

  value_type = {
    type = "JSON"
  }
}

resource "pingone_authorize_trust_framework_service" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "Test application service"

  service_type = "CONNECTOR"

  cache_settings = {
    ttl_seconds = 300
  }

  processor = {
    name = "%[3]s Test chain processor"
    type = "CHAIN"

    processors = [
      {
        name = "%[3]s Test chain processor 1"
        type = "JSON_PATH"

        expression = "$.data.item1"
        value_type = {
          type = "STRING"
        }
      },
      {
        name = "%[3]s Test chain processor 3"
        type = "JSON_PATH"

        expression = "$.data.item3"
        value_type = {
          type = "JSON"
        }
      },
      {
        name = "%[3]s Test chain processor 2"
        type = "JSON_PATH"

        expression = "$.data.item2"
        value_type = {
          type = "STRING"
        }
      },
    ],
  }

  service_settings = {
    capability = "createRiskEvaluation"
    channel    = "AUTHORIZE"
    code       = "P1_RISK"

    input_mappings = [
      {
        type     = "ATTRIBUTE"
        property = "userId"
        value_ref = {
          id = pingone_authorize_trust_framework_attribute.%[2]s-user.id
        }
      },
      {
        type     = "INPUT"
        property = "userType"
        value    = "test2"
      },
      {
        type     = "INPUT"
        property = "ipAddress"
        value    = "test3"
      },
      {
        type     = "INPUT"
        property = "applicationId"
        value    = "test4"
      },
      {
        type     = "INPUT"
        property = "applicationName"
        value    = "test5"
      },
      {
        type     = "INPUT"
        property = "sessionId"
        value    = "test6"
      },
      {
        type     = "INPUT"
        property = "browserData"
        value    = "test7"
      },
      {
        type     = "INPUT"
        property = "riskPolicySetId"
        value    = "test8"
      }
    ]

    maximum_concurrent_requests = 6
    maximum_requests_per_second = 10
    schema_version              = "1"
    timeout_milliseconds        = 2000
  }

  value_type = {
    type = "JSON"
  }
}`, acctest.AuthorizePMTFSandboxEnvironment(), resourceName, name)
}

func testAccTrustFrameworkServiceConfig_Service_Connector_Full2(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_authorize_trust_framework_attribute" "%[2]s-user" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "Test application service"

  resolvers = [
    {
      type = "CURRENT_USER_ID"
    }
  ]

  value_type = {
    type = "JSON"
  }
}

resource "pingone_authorize_trust_framework_service" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "Test application service"

  service_type = "CONNECTOR"

  cache_settings = {
    ttl_seconds = 400
  }

  processor = {
    name = "%[3]s Test chain processor 1"
    type = "JSON_PATH"

    expression = "$.data.item1"
    value_type = {
      type = "STRING"
    }
  }

  service_settings = {
    capability = "updateRiskEvaluation"
    channel    = "AUTHORIZE"
    code       = "P1_RISK"

    input_mappings = [
      {
        property = "completionStatus"
        type     = "INPUT"

        value = "completed"
      },
      {
        property = "riskId"
        type     = "INPUT"

        value = "test"
      },
    ]

    maximum_concurrent_requests = 4
    maximum_requests_per_second = 11
    schema_version              = "2"
    timeout_milliseconds        = 2500
  }

  value_type = {
    type = "XML"
  }
}`, acctest.AuthorizePMTFSandboxEnvironment(), resourceName, name)
}

func testAccTrustFrameworkServiceConfig_Service_Connector_Minimal(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_authorize_trust_framework_attribute" "%[2]s-user" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "Test application service"

  resolvers = [
    {
      type = "CURRENT_USER_ID"
    }
  ]

  value_type = {
    type = "JSON"
  }
}

resource "pingone_authorize_trust_framework_service" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "Test application service"

  service_type = "CONNECTOR"

  service_settings = {
    capability = "createRiskEvaluation"
    channel    = "AUTHORIZE"
    code       = "P1_RISK"

    input_mappings = [
      {
        property = "ipAddress"
        type     = "INPUT"

        value = "192.168.0.1"
      },
      {
        property = "userType"
        type     = "INPUT"

        value = "person"
      },
      {
        property = "userId"
        type     = "INPUT"

        value = "test1"
      },
    ]

    maximum_concurrent_requests = 6
    maximum_requests_per_second = 10
    schema_version              = "1"
    timeout_milliseconds        = 2000
  }

  value_type = {
    type = "STRING"
  }
}`, acctest.AuthorizePMTFSandboxEnvironment(), resourceName, name)
}

func testAccTrustFrameworkServiceConfig_Service_None_Full(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_authorize_trust_framework_attribute" "%[2]s-user" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "Test application service"

  resolvers = [
    {
      type = "CURRENT_USER_ID"
    }
  ]

  value_type = {
    type = "JSON"
  }
}

resource "pingone_authorize_trust_framework_service" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "Test application service"

  service_type = "NONE"

  cache_settings = {
    ttl_seconds = 300
  }
}`, acctest.AuthorizePMTFSandboxEnvironment(), resourceName, name)
}

func testAccTrustFrameworkServiceConfig_Service_None_Minimal(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_authorize_trust_framework_service" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "Test application service"

  service_type = "NONE"
}`, acctest.AuthorizePMTFSandboxEnvironment(), resourceName, name)
}
