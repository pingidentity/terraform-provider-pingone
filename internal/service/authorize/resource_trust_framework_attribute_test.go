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

func TestAccTrustFrameworkAttribute_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authorize_trust_framework_attribute.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	var attributeID, environmentID string

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
		CheckDestroy:             authorize.TrustFrameworkAttribute_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccTrustFrameworkAttributeConfig_Minimal(resourceName, name),
				Check:  authorize.TrustFrameworkAttribute_GetIDs(resourceFullName, &environmentID, &attributeID),
			},
			{
				PreConfig: func() {
					authorize.TrustFrameworkAttribute_RemovalDrift_PreConfig(ctx, p1Client.API.AuthorizeAPIClient, t, environmentID, attributeID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			// Test removal of the environment
			{
				Config: testAccTrustFrameworkAttributeConfig_NewEnv(environmentName, licenseID, resourceName, name),
				Check:  authorize.TrustFrameworkAttribute_GetIDs(resourceFullName, &environmentID, &attributeID),
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

func TestAccTrustFrameworkAttribute_NewEnv(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authorize_trust_framework_attribute.%s", resourceName)

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
		CheckDestroy:             authorize.TrustFrameworkAttribute_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccTrustFrameworkAttributeConfig_NewEnv(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
				),
			},
		},
	})
}

func TestAccTrustFrameworkAttribute_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authorize_trust_framework_attribute.%s", resourceName)

	name := resourceName

	fullCheck := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckResourceAttr(resourceFullName, "name", name),
		resource.TestCheckResourceAttr(resourceFullName, "description", "Test attribute full"),
		resource.TestCheckResourceAttr(resourceFullName, "full_name", fmt.Sprintf("%[1]s-parent.%[1]s", name)),
		resource.TestMatchResourceAttr(resourceFullName, "parent.id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckResourceAttr(resourceFullName, "default_value", "test"),
		resource.TestCheckResourceAttr(resourceFullName, "processor.name", fmt.Sprintf("%s Test processor", name)),
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.#", "1"),
		resource.TestCheckNoResourceAttr(resourceFullName, "managed_entity"),
		resource.TestMatchResourceAttr(resourceFullName, "repetition_source.id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckResourceAttr(resourceFullName, "type", "ATTRIBUTE"),
		resource.TestCheckResourceAttr(resourceFullName, "value_type.type", "STRING"),
		resource.TestMatchResourceAttr(resourceFullName, "version", verify.P1ResourceIDRegexpFullString),
	)

	minimalCheck := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckResourceAttr(resourceFullName, "name", name),
		resource.TestCheckResourceAttr(resourceFullName, "description", "Test attribute"),
		resource.TestCheckResourceAttr(resourceFullName, "full_name", name),
		resource.TestCheckNoResourceAttr(resourceFullName, "parent"),
		resource.TestCheckNoResourceAttr(resourceFullName, "default_value"),
		resource.TestCheckNoResourceAttr(resourceFullName, "processor"),
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.#", "0"),
		resource.TestCheckNoResourceAttr(resourceFullName, "managed_entity"),
		resource.TestCheckNoResourceAttr(resourceFullName, "repetition_source"),
		resource.TestCheckResourceAttr(resourceFullName, "type", "ATTRIBUTE"),
		resource.TestCheckResourceAttr(resourceFullName, "value_type.type", "STRING"),
		resource.TestMatchResourceAttr(resourceFullName, "version", verify.P1ResourceIDRegexpFullString),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckFeatureFlag(t, acctest.ENUMFEATUREFLAG_AUTHORIZEPMTF)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             authorize.TrustFrameworkAttribute_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			{
				Config: testAccTrustFrameworkAttributeConfig_Full(resourceName, name),
				Check:  fullCheck,
			},
			{
				Config:  testAccTrustFrameworkAttributeConfig_Full(resourceName, name),
				Destroy: true,
			},
			// Minimal
			{
				Config: testAccTrustFrameworkAttributeConfig_Minimal(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config:  testAccTrustFrameworkAttributeConfig_Minimal(resourceName, name),
				Destroy: true,
			},
			// Change
			{
				Config: testAccTrustFrameworkAttributeConfig_Full(resourceName, name),
				Check:  fullCheck,
			},
			{
				Config: testAccTrustFrameworkAttributeConfig_Minimal(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config: testAccTrustFrameworkAttributeConfig_Full(resourceName, name),
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

func TestAccTrustFrameworkAttribute_ComplexAttributes(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authorize_trust_framework_attribute.%s", resourceName)

	name := resourceName

	attribute1Check := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "default_value", "{\"foo\":\"bar\"}"),
		resource.TestCheckResourceAttr(resourceFullName, "processor.name", fmt.Sprintf("%s Test chain processor", name)),
		resource.TestCheckResourceAttr(resourceFullName, "processor.type", "CHAIN"),
		resource.TestCheckResourceAttr(resourceFullName, "processor.processors.0.name", fmt.Sprintf("%s Test chain processor 1", name)),
		resource.TestCheckResourceAttr(resourceFullName, "processor.processors.0.type", "JSON_PATH"),
		resource.TestCheckResourceAttr(resourceFullName, "processor.processors.0.expression", "$.data.item1"),
		resource.TestCheckResourceAttr(resourceFullName, "processor.processors.0.value_type.type", "STRING"),
		resource.TestCheckResourceAttr(resourceFullName, "processor.processors.1.name", fmt.Sprintf("%s Test chain processor 2", name)),
		resource.TestCheckResourceAttr(resourceFullName, "processor.processors.1.type", "JSON_PATH"),
		resource.TestCheckResourceAttr(resourceFullName, "processor.processors.1.expression", "$.data.item2"),
		resource.TestCheckResourceAttr(resourceFullName, "processor.processors.1.value_type.type", "STRING"),
		resource.TestCheckResourceAttr(resourceFullName, "processor.processors.2.name", fmt.Sprintf("%s Test chain processor 3", name)),
		resource.TestCheckResourceAttr(resourceFullName, "processor.processors.2.type", "JSON_PATH"),
		resource.TestCheckResourceAttr(resourceFullName, "processor.processors.2.expression", "$.data.item3"),
		resource.TestCheckResourceAttr(resourceFullName, "processor.processors.2.value_type.type", "JSON"),
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.#", "3"),
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.0.name", fmt.Sprintf("%s Test Attribute Resolver", name)),
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.1.name", fmt.Sprintf("%s Test Constant Resolver", name)),
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.2.name", fmt.Sprintf("%s Test Current User ID Resolver", name)),
		resource.TestCheckResourceAttr(resourceFullName, "value_type.type", "JSON"),
	)

	attribute2Check := resource.ComposeTestCheckFunc(
		resource.TestCheckNoResourceAttr(resourceFullName, "default_value"),
		resource.TestCheckResourceAttr(resourceFullName, "processor.name", fmt.Sprintf("%s Test processor", name)),
		resource.TestCheckResourceAttr(resourceFullName, "processor.type", "JSON_PATH"),
		resource.TestCheckResourceAttr(resourceFullName, "processor.expression", "$.data.item.item1"),
		resource.TestCheckResourceAttr(resourceFullName, "processor.value_type.type", "STRING"),
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.#", "3"),
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.0.name", fmt.Sprintf("%s Test Current User ID Resolver", name)),
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.1.name", fmt.Sprintf("%s Test Attribute Resolver", name)),
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.2.name", fmt.Sprintf("%s Test Constant Resolver", name)),
		resource.TestCheckResourceAttr(resourceFullName, "value_type.type", "DURATION"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckFeatureFlag(t, acctest.ENUMFEATUREFLAG_AUTHORIZEPMTF)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             authorize.TrustFrameworkAttribute_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// ComplexAttribute1
			{
				Config: testAccTrustFrameworkAttributeConfig_ComplexAttribute1(resourceName, name),
				Check:  attribute1Check,
			},
			{
				Config:  testAccTrustFrameworkAttributeConfig_ComplexAttribute1(resourceName, name),
				Destroy: true,
			},
			// ComplexAttribute2
			{
				Config: testAccTrustFrameworkAttributeConfig_ComplexAttribute2(resourceName, name),
				Check:  attribute2Check,
			},
			{
				Config:  testAccTrustFrameworkAttributeConfig_ComplexAttribute2(resourceName, name),
				Destroy: true,
			},
			// Change
			{
				Config: testAccTrustFrameworkAttributeConfig_ComplexAttribute1(resourceName, name),
				Check:  attribute1Check,
			},
			{
				Config: testAccTrustFrameworkAttributeConfig_ComplexAttribute2(resourceName, name),
				Check:  attribute2Check,
			},
			{
				Config: testAccTrustFrameworkAttributeConfig_ComplexAttribute1(resourceName, name),
				Check:  attribute1Check,
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

func TestAccTrustFrameworkAttribute_Resolver_Attribute(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authorize_trust_framework_attribute.%s", resourceName)

	name := resourceName

	fullCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.#", "1"),
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.0.condition.type", "AND"),
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.0.condition.conditions.#", "3"),
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.0.name", fmt.Sprintf("%s Test Attribute Resolver", name)),
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.0.processor.name", fmt.Sprintf("%s Test processor 1", name)),
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.0.processor.type", "JSON_PATH"),
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.0.type", "ATTRIBUTE"),
		resource.TestMatchResourceAttr(resourceFullName, "resolvers.0.value_ref.id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckNoResourceAttr(resourceFullName, "resolvers.0.value_type"),
		resource.TestCheckNoResourceAttr(resourceFullName, "resolvers.0.value"),
		resource.TestCheckNoResourceAttr(resourceFullName, "resolvers.0.query"),
	)

	minimalCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.#", "1"),
		resource.TestCheckNoResourceAttr(resourceFullName, "resolvers.0.condition"),
		resource.TestCheckNoResourceAttr(resourceFullName, "resolvers.0.name"),
		resource.TestCheckNoResourceAttr(resourceFullName, "resolvers.0.processor"),
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.0.type", "ATTRIBUTE"),
		resource.TestMatchResourceAttr(resourceFullName, "resolvers.0.value_ref.id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckNoResourceAttr(resourceFullName, "resolvers.0.value_type"),
		resource.TestCheckNoResourceAttr(resourceFullName, "resolvers.0.value"),
		resource.TestCheckNoResourceAttr(resourceFullName, "resolvers.0.query"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckFeatureFlag(t, acctest.ENUMFEATUREFLAG_AUTHORIZEPMTF)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             authorize.TrustFrameworkAttribute_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// fullCheck
			{
				Config: testAccTrustFrameworkAttributeConfig_Resolver_Attribute_Full(resourceName, name),
				Check:  fullCheck,
			},
			{
				Config:  testAccTrustFrameworkAttributeConfig_Resolver_Attribute_Full(resourceName, name),
				Destroy: true,
			},
			// minimalCheck
			{
				Config: testAccTrustFrameworkAttributeConfig_Resolver_Attribute_Min(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config:  testAccTrustFrameworkAttributeConfig_Resolver_Attribute_Min(resourceName, name),
				Destroy: true,
			},
			// Change
			{
				Config: testAccTrustFrameworkAttributeConfig_Resolver_Attribute_Full(resourceName, name),
				Check:  fullCheck,
			},
			{
				Config: testAccTrustFrameworkAttributeConfig_Resolver_Attribute_Min(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config: testAccTrustFrameworkAttributeConfig_Resolver_Attribute_Full(resourceName, name),
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

func TestAccTrustFrameworkAttribute_Resolver_Constant(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authorize_trust_framework_attribute.%s", resourceName)

	name := resourceName

	fullCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.#", "1"),
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.0.condition.type", "AND"),
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.0.condition.conditions.#", "3"),
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.0.name", fmt.Sprintf("%s Test Constant Resolver", name)),
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.0.processor.name", fmt.Sprintf("%s Test processor 1", name)),
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.0.processor.type", "JSON_PATH"),
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.0.type", "CONSTANT"),
		resource.TestCheckNoResourceAttr(resourceFullName, "resolvers.0.value_ref"),
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.0.value_type.type", "JSON"),
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.0.value", "{\"foo\":\"bar\"}"),
		resource.TestCheckNoResourceAttr(resourceFullName, "resolvers.0.query"),
	)

	minimalCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.#", "1"),
		resource.TestCheckNoResourceAttr(resourceFullName, "resolvers.0.condition"),
		resource.TestCheckNoResourceAttr(resourceFullName, "resolvers.0.name"),
		resource.TestCheckNoResourceAttr(resourceFullName, "resolvers.0.processor"),
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.0.type", "CONSTANT"),
		resource.TestCheckNoResourceAttr(resourceFullName, "resolvers.0.value_ref"),
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.0.value_type.type", "STRING"),
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.0.value", "test"),
		resource.TestCheckNoResourceAttr(resourceFullName, "resolvers.0.query"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckFeatureFlag(t, acctest.ENUMFEATUREFLAG_AUTHORIZEPMTF)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             authorize.TrustFrameworkAttribute_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// fullCheck
			{
				Config: testAccTrustFrameworkAttributeConfig_Resolver_Constant_Full(resourceName, name),
				Check:  fullCheck,
			},
			{
				Config:  testAccTrustFrameworkAttributeConfig_Resolver_Constant_Full(resourceName, name),
				Destroy: true,
			},
			// minimalCheck
			{
				Config: testAccTrustFrameworkAttributeConfig_Resolver_Constant_Min(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config:  testAccTrustFrameworkAttributeConfig_Resolver_Constant_Min(resourceName, name),
				Destroy: true,
			},
			// Change
			{
				Config: testAccTrustFrameworkAttributeConfig_Resolver_Constant_Full(resourceName, name),
				Check:  fullCheck,
			},
			{
				Config: testAccTrustFrameworkAttributeConfig_Resolver_Constant_Min(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config: testAccTrustFrameworkAttributeConfig_Resolver_Constant_Full(resourceName, name),
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

func TestAccTrustFrameworkAttribute_Resolver_CurrentRepetitionValue(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authorize_trust_framework_attribute.%s", resourceName)

	name := resourceName

	fullCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.#", "1"),
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.0.condition.type", "AND"),
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.0.condition.conditions.#", "3"),
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.0.name", fmt.Sprintf("%s Test Current Repetition Resolver", name)),
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.0.processor.name", fmt.Sprintf("%s Test processor 1", name)),
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.0.processor.type", "JSON_PATH"),
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.0.type", "CURRENT_REPETITION_VALUE"),
		resource.TestCheckNoResourceAttr(resourceFullName, "resolvers.0.value_ref"),
		resource.TestCheckNoResourceAttr(resourceFullName, "resolvers.0.value_type"),
		resource.TestCheckNoResourceAttr(resourceFullName, "resolvers.0.value"),
		resource.TestCheckNoResourceAttr(resourceFullName, "resolvers.0.query"),
	)

	minimalCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.#", "1"),
		resource.TestCheckNoResourceAttr(resourceFullName, "resolvers.0.condition"),
		resource.TestCheckNoResourceAttr(resourceFullName, "resolvers.0.name"),
		resource.TestCheckNoResourceAttr(resourceFullName, "resolvers.0.processor"),
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.0.type", "CURRENT_REPETITION_VALUE"),
		resource.TestCheckNoResourceAttr(resourceFullName, "resolvers.0.value_ref"),
		resource.TestCheckNoResourceAttr(resourceFullName, "resolvers.0.value_type"),
		resource.TestCheckNoResourceAttr(resourceFullName, "resolvers.0.value"),
		resource.TestCheckNoResourceAttr(resourceFullName, "resolvers.0.query"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckFeatureFlag(t, acctest.ENUMFEATUREFLAG_AUTHORIZEPMTF)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             authorize.TrustFrameworkAttribute_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// fullCheck
			{
				Config: testAccTrustFrameworkAttributeConfig_Resolver_CurrentRepetitionValue_Full(resourceName, name),
				Check:  fullCheck,
			},
			{
				Config:  testAccTrustFrameworkAttributeConfig_Resolver_CurrentRepetitionValue_Full(resourceName, name),
				Destroy: true,
			},
			// minimalCheck
			{
				Config: testAccTrustFrameworkAttributeConfig_Resolver_CurrentRepetitionValue_Min(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config:  testAccTrustFrameworkAttributeConfig_Resolver_CurrentRepetitionValue_Min(resourceName, name),
				Destroy: true,
			},
			// Change
			{
				Config: testAccTrustFrameworkAttributeConfig_Resolver_CurrentRepetitionValue_Full(resourceName, name),
				Check:  fullCheck,
			},
			{
				Config: testAccTrustFrameworkAttributeConfig_Resolver_CurrentRepetitionValue_Min(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config: testAccTrustFrameworkAttributeConfig_Resolver_CurrentRepetitionValue_Full(resourceName, name),
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

func TestAccTrustFrameworkAttribute_Resolver_CurrentUserId(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authorize_trust_framework_attribute.%s", resourceName)

	name := resourceName

	fullCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.#", "1"),
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.0.condition.type", "AND"),
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.0.condition.conditions.#", "3"),
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.0.name", fmt.Sprintf("%s Test Current User Id Resolver", name)),
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.0.processor.name", fmt.Sprintf("%s Test processor 1", name)),
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.0.processor.type", "JSON_PATH"),
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.0.type", "CURRENT_USER_ID"),
		resource.TestCheckNoResourceAttr(resourceFullName, "resolvers.0.value_ref"),
		resource.TestCheckNoResourceAttr(resourceFullName, "resolvers.0.value_type"),
		resource.TestCheckNoResourceAttr(resourceFullName, "resolvers.0.value"),
		resource.TestCheckNoResourceAttr(resourceFullName, "resolvers.0.query"),
	)

	minimalCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.#", "1"),
		resource.TestCheckNoResourceAttr(resourceFullName, "resolvers.0.condition"),
		resource.TestCheckNoResourceAttr(resourceFullName, "resolvers.0.name"),
		resource.TestCheckNoResourceAttr(resourceFullName, "resolvers.0.processor"),
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.0.type", "CURRENT_USER_ID"),
		resource.TestCheckNoResourceAttr(resourceFullName, "resolvers.0.value_ref"),
		resource.TestCheckNoResourceAttr(resourceFullName, "resolvers.0.value_type"),
		resource.TestCheckNoResourceAttr(resourceFullName, "resolvers.0.value"),
		resource.TestCheckNoResourceAttr(resourceFullName, "resolvers.0.query"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckFeatureFlag(t, acctest.ENUMFEATUREFLAG_AUTHORIZEPMTF)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             authorize.TrustFrameworkAttribute_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// fullCheck
			{
				Config: testAccTrustFrameworkAttributeConfig_Resolver_CurrentUserId_Full(resourceName, name),
				Check:  fullCheck,
			},
			{
				Config:  testAccTrustFrameworkAttributeConfig_Resolver_CurrentUserId_Full(resourceName, name),
				Destroy: true,
			},
			// minimalCheck
			{
				Config: testAccTrustFrameworkAttributeConfig_Resolver_CurrentUserId_Min(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config:  testAccTrustFrameworkAttributeConfig_Resolver_CurrentUserId_Min(resourceName, name),
				Destroy: true,
			},
			// Change
			{
				Config: testAccTrustFrameworkAttributeConfig_Resolver_CurrentUserId_Full(resourceName, name),
				Check:  fullCheck,
			},
			{
				Config: testAccTrustFrameworkAttributeConfig_Resolver_CurrentUserId_Min(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config: testAccTrustFrameworkAttributeConfig_Resolver_CurrentUserId_Full(resourceName, name),
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

func TestAccTrustFrameworkAttribute_Resolver_Request(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authorize_trust_framework_attribute.%s", resourceName)

	name := resourceName

	fullCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.#", "1"),
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.0.condition.type", "AND"),
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.0.condition.conditions.#", "3"),
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.0.name", fmt.Sprintf("%s Test Request Resolver", name)),
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.0.processor.name", fmt.Sprintf("%s Test processor 1", name)),
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.0.processor.type", "JSON_PATH"),
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.0.type", "REQUEST"),
		resource.TestCheckNoResourceAttr(resourceFullName, "resolvers.0.value_ref"),
		resource.TestCheckNoResourceAttr(resourceFullName, "resolvers.0.value_type"),
		resource.TestCheckNoResourceAttr(resourceFullName, "resolvers.0.value"),
		resource.TestCheckNoResourceAttr(resourceFullName, "resolvers.0.query"),
	)

	minimalCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.#", "1"),
		resource.TestCheckNoResourceAttr(resourceFullName, "resolvers.0.condition"),
		resource.TestCheckNoResourceAttr(resourceFullName, "resolvers.0.name"),
		resource.TestCheckNoResourceAttr(resourceFullName, "resolvers.0.processor"),
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.0.type", "REQUEST"),
		resource.TestCheckNoResourceAttr(resourceFullName, "resolvers.0.value_ref"),
		resource.TestCheckNoResourceAttr(resourceFullName, "resolvers.0.value_type"),
		resource.TestCheckNoResourceAttr(resourceFullName, "resolvers.0.value"),
		resource.TestCheckNoResourceAttr(resourceFullName, "resolvers.0.query"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckFeatureFlag(t, acctest.ENUMFEATUREFLAG_AUTHORIZEPMTF)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             authorize.TrustFrameworkAttribute_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// fullCheck
			{
				Config: testAccTrustFrameworkAttributeConfig_Resolver_Request_Full(resourceName, name),
				Check:  fullCheck,
			},
			{
				Config:  testAccTrustFrameworkAttributeConfig_Resolver_Request_Full(resourceName, name),
				Destroy: true,
			},
			// minimalCheck
			{
				Config: testAccTrustFrameworkAttributeConfig_Resolver_Request_Min(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config:  testAccTrustFrameworkAttributeConfig_Resolver_Request_Min(resourceName, name),
				Destroy: true,
			},
			// Change
			{
				Config: testAccTrustFrameworkAttributeConfig_Resolver_Request_Full(resourceName, name),
				Check:  fullCheck,
			},
			{
				Config: testAccTrustFrameworkAttributeConfig_Resolver_Request_Min(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config: testAccTrustFrameworkAttributeConfig_Resolver_Request_Full(resourceName, name),
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

func TestAccTrustFrameworkAttribute_Resolver_Service(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authorize_trust_framework_attribute.%s", resourceName)

	name := resourceName

	fullCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.#", "1"),
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.0.condition.type", "AND"),
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.0.condition.conditions.#", "3"),
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.0.name", fmt.Sprintf("%s Test Service Resolver", name)),
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.0.processor.name", fmt.Sprintf("%s Test processor 1", name)),
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.0.processor.type", "JSON_PATH"),
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.0.type", "SERVICE"),
		resource.TestMatchResourceAttr(resourceFullName, "resolvers.0.value_ref.id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckNoResourceAttr(resourceFullName, "resolvers.0.value_type"),
		resource.TestCheckNoResourceAttr(resourceFullName, "resolvers.0.value"),
		resource.TestCheckNoResourceAttr(resourceFullName, "resolvers.0.query"),
	)

	minimalCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.#", "1"),
		resource.TestCheckNoResourceAttr(resourceFullName, "resolvers.0.condition"),
		resource.TestCheckNoResourceAttr(resourceFullName, "resolvers.0.name"),
		resource.TestCheckNoResourceAttr(resourceFullName, "resolvers.0.processor"),
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.0.type", "SERVICE"),
		resource.TestMatchResourceAttr(resourceFullName, "resolvers.0.value_ref.id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckNoResourceAttr(resourceFullName, "resolvers.0.value_type"),
		resource.TestCheckNoResourceAttr(resourceFullName, "resolvers.0.value"),
		resource.TestCheckNoResourceAttr(resourceFullName, "resolvers.0.query"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckFeatureFlag(t, acctest.ENUMFEATUREFLAG_AUTHORIZEPMTF)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             authorize.TrustFrameworkAttribute_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// fullCheck
			{
				Config: testAccTrustFrameworkAttributeConfig_Resolver_Service_Full(resourceName, name),
				Check:  fullCheck,
			},
			{
				Config:  testAccTrustFrameworkAttributeConfig_Resolver_Service_Full(resourceName, name),
				Destroy: true,
			},
			// minimalCheck
			{
				Config: testAccTrustFrameworkAttributeConfig_Resolver_Service_Min(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config:  testAccTrustFrameworkAttributeConfig_Resolver_Service_Min(resourceName, name),
				Destroy: true,
			},
			// Change
			{
				Config: testAccTrustFrameworkAttributeConfig_Resolver_Service_Full(resourceName, name),
				Check:  fullCheck,
			},
			{
				Config: testAccTrustFrameworkAttributeConfig_Resolver_Service_Min(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config: testAccTrustFrameworkAttributeConfig_Resolver_Service_Full(resourceName, name),
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

func TestAccTrustFrameworkAttribute_Resolver_System(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authorize_trust_framework_attribute.%s", resourceName)

	name := resourceName

	fullCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.#", "1"),
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.0.condition.type", "AND"),
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.0.condition.conditions.#", "3"),
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.0.name", fmt.Sprintf("%s Test System Resolver", name)),
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.0.processor.name", fmt.Sprintf("%s Test processor 1", name)),
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.0.processor.type", "JSON_PATH"),
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.0.type", "SYSTEM"),
		resource.TestCheckNoResourceAttr(resourceFullName, "resolvers.0.value_ref.id"),
		resource.TestCheckNoResourceAttr(resourceFullName, "resolvers.0.value_type"),
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.0.value", "CURRENT_DATE_TIME"),
		resource.TestCheckNoResourceAttr(resourceFullName, "resolvers.0.query"),
	)

	minimalCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.#", "1"),
		resource.TestCheckNoResourceAttr(resourceFullName, "resolvers.0.condition"),
		resource.TestCheckNoResourceAttr(resourceFullName, "resolvers.0.name"),
		resource.TestCheckNoResourceAttr(resourceFullName, "resolvers.0.processor"),
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.0.type", "SYSTEM"),
		resource.TestCheckNoResourceAttr(resourceFullName, "resolvers.0.value_ref.id"),
		resource.TestCheckNoResourceAttr(resourceFullName, "resolvers.0.value_type"),
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.0.value", "NULL"),
		resource.TestCheckNoResourceAttr(resourceFullName, "resolvers.0.query"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckFeatureFlag(t, acctest.ENUMFEATUREFLAG_AUTHORIZEPMTF)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             authorize.TrustFrameworkAttribute_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// fullCheck
			{
				Config: testAccTrustFrameworkAttributeConfig_Resolver_System_Full(resourceName, name),
				Check:  fullCheck,
			},
			{
				Config:  testAccTrustFrameworkAttributeConfig_Resolver_System_Full(resourceName, name),
				Destroy: true,
			},
			// minimalCheck
			{
				Config: testAccTrustFrameworkAttributeConfig_Resolver_System_Min(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config:  testAccTrustFrameworkAttributeConfig_Resolver_System_Min(resourceName, name),
				Destroy: true,
			},
			// Change
			{
				Config: testAccTrustFrameworkAttributeConfig_Resolver_System_Full(resourceName, name),
				Check:  fullCheck,
			},
			{
				Config: testAccTrustFrameworkAttributeConfig_Resolver_System_Min(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config: testAccTrustFrameworkAttributeConfig_Resolver_System_Full(resourceName, name),
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

func TestAccTrustFrameworkAttribute_Resolver_User(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authorize_trust_framework_attribute.%s", resourceName)

	name := resourceName

	fullCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.#", "1"),
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.0.condition.type", "AND"),
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.0.condition.conditions.#", "3"),
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.0.name", fmt.Sprintf("%s Test User Resolver", name)),
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.0.processor.name", fmt.Sprintf("%s Test processor 1", name)),
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.0.processor.type", "JSON_PATH"),
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.0.type", "USER"),
		resource.TestCheckNoResourceAttr(resourceFullName, "resolvers.0.value_ref.id"),
		resource.TestCheckNoResourceAttr(resourceFullName, "resolvers.0.value_type"),
		resource.TestCheckNoResourceAttr(resourceFullName, "resolvers.0.value"),
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.0.query.type", "USER_ID"),
		resource.TestMatchResourceAttr(resourceFullName, "resolvers.0.query.user_id", verify.P1ResourceIDRegexpFullString),
	)

	minimalCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.#", "1"),
		resource.TestCheckNoResourceAttr(resourceFullName, "resolvers.0.condition"),
		resource.TestCheckNoResourceAttr(resourceFullName, "resolvers.0.name"),
		resource.TestCheckNoResourceAttr(resourceFullName, "resolvers.0.processor"),
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.0.type", "USER"),
		resource.TestCheckNoResourceAttr(resourceFullName, "resolvers.0.value_ref.id"),
		resource.TestCheckNoResourceAttr(resourceFullName, "resolvers.0.value_type"),
		resource.TestCheckNoResourceAttr(resourceFullName, "resolvers.0.value"),
		resource.TestCheckResourceAttr(resourceFullName, "resolvers.0.query.type", "USER_ID"),
		resource.TestMatchResourceAttr(resourceFullName, "resolvers.0.query.user_id", verify.P1ResourceIDRegexpFullString),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckFeatureFlag(t, acctest.ENUMFEATUREFLAG_AUTHORIZEPMTF)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             authorize.TrustFrameworkAttribute_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// fullCheck
			{
				Config: testAccTrustFrameworkAttributeConfig_Resolver_User_Full(resourceName, name),
				Check:  fullCheck,
			},
			{
				Config:  testAccTrustFrameworkAttributeConfig_Resolver_User_Full(resourceName, name),
				Destroy: true,
			},
			// minimalCheck
			{
				Config: testAccTrustFrameworkAttributeConfig_Resolver_User_Min(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config:  testAccTrustFrameworkAttributeConfig_Resolver_User_Min(resourceName, name),
				Destroy: true,
			},
			// Change
			{
				Config: testAccTrustFrameworkAttributeConfig_Resolver_User_Full(resourceName, name),
				Check:  fullCheck,
			},
			{
				Config: testAccTrustFrameworkAttributeConfig_Resolver_User_Min(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config: testAccTrustFrameworkAttributeConfig_Resolver_User_Full(resourceName, name),
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

func TestAccTrustFrameworkAttribute_BadParameters(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authorize_trust_framework_attribute.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckFeatureFlag(t, acctest.ENUMFEATUREFLAG_AUTHORIZEPMTF)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             authorize.TrustFrameworkAttribute_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccTrustFrameworkAttributeConfig_Minimal(resourceName, name),
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

func testAccTrustFrameworkAttributeConfig_NewEnv(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_authorize_trust_framework_attribute" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  name           = "%[3]s"
  description    = "Test attribute"

  value_type = {
    type = "STRING"
  }
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccTrustFrameworkAttributeConfig_Full(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_authorize_trust_framework_attribute" "%[2]s-parent" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s-parent"
  description    = "Test attribute"

  value_type = {
    type = "STRING"
  }
}

resource "pingone_authorize_trust_framework_attribute" "%[2]s-repetition" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s-repetition"
  description    = "Test attribute"

  value_type = {
    type = "COLLECTION"
  }
}

resource "pingone_authorize_trust_framework_attribute" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "Test attribute full"

  parent = {
    id = pingone_authorize_trust_framework_attribute.%[2]s-parent.id
  }

  default_value = "test"

  repetition_source = {
    id = pingone_authorize_trust_framework_attribute.%[2]s-repetition.id
  }

  processor = {
    name = "%[3]s Test processor"
    type = "JSON_PATH"

    expression = "$.data.item.parent"
    value_type = {
      type = "STRING"
    }
  }

  resolvers = [
    {
      type = "CONSTANT"
      value_type = {
        type = "STRING"
      }
      value = "test"
    }
  ]

  value_type = {
    type = "STRING"
  }
}`, acctest.AuthorizePMTFSandboxEnvironment(), resourceName, name)
}

func testAccTrustFrameworkAttributeConfig_Minimal(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_authorize_trust_framework_attribute" "%[2]s-parent" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s-parent"
  description    = "Test attribute"

  value_type = {
    type = "STRING"
  }
}

resource "pingone_authorize_trust_framework_attribute" "%[2]s-repetition" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s-repetition"
  description    = "Test attribute"

  value_type = {
    type = "COLLECTION"
  }
}

resource "pingone_authorize_trust_framework_attribute" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "Test attribute"

  value_type = {
    type = "STRING"
  }
}`, acctest.AuthorizePMTFSandboxEnvironment(), resourceName, name)
}

func testAccTrustFrameworkAttributeConfig_ComplexAttribute1(resourceName, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_authorize_trust_framework_attribute" "%[2]s-resolver" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s-resolver"
  description    = "Test attribute"

  value_type = {
    type = "STRING"
  }
}

resource "pingone_authorize_trust_framework_attribute" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "Test attribute"

  default_value = jsonencode({ "foo" : "bar" })

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
        name = "%[3]s Test chain processor 2"
        type = "JSON_PATH"

        expression = "$.data.item2"
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
    ],
  }

  resolvers = [
    {
      name = "%[3]s Test Attribute Resolver"

      type = "ATTRIBUTE"
      value_ref = {
        id = pingone_authorize_trust_framework_attribute.%[2]s-resolver.id
      }

      condition = {
        type = "AND"

        conditions = [
          {
            type = "EMPTY"
          },
          {
            type       = "COMPARISON"
            comparator = "EQUALS"

            left = {
              type  = "CONSTANT"
              value = "test"
            }

            right = {
              type  = "CONSTANT"
              value = "test1"
            }
          },
          {
            type = "NOT"

            condition = {
              type       = "COMPARISON"
              comparator = "EQUALS"

              left = {
                type  = "CONSTANT"
                value = "test2"
              }

              right = {
                type  = "CONSTANT"
                value = "test3"
              }
            }
          }
        ]
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
            name = "%[3]s Test chain processor 2"
            type = "JSON_PATH"

            expression = "$.data.item2"
            value_type = {
              type = "STRING"
            }
          },
          {
            name = "%[3]s Test chain processor 3"
            type = "JSON_PATH"

            expression = "$.data.item3"
            value_type = {
              type = "BOOLEAN"
            }
          },
        ],
      }
    },
    {
      name = "%[3]s Test Constant Resolver"

      type  = "CONSTANT"
      value = "test"
      value_type = {
        type = "STRING"
      }

      condition = {
        type       = "COMPARISON"
        comparator = "EQUALS"

        left = {
          type  = "CONSTANT"
          value = "test"
        }

        right = {
          type  = "CONSTANT"
          value = "test1"
        }
      }

      processor = {
        name = "%[3]s Test processor 1"
        type = "JSON_PATH"

        expression = "$.data.item1"
        value_type = {
          type = "STRING"
        }
      }
    },
    {
      name = "%[3]s Test Current User ID Resolver"

      type = "CURRENT_USER_ID"

      condition = {
        type       = "COMPARISON"
        comparator = "EQUALS"

        left = {
          type  = "CONSTANT"
          value = "test"
        }

        right = {
          type  = "CONSTANT"
          value = "test1"
        }
      }

      processor = {
        name = "%[3]s Test processor 1"
        type = "JSON_PATH"

        expression = "$.data.item1"
        value_type = {
          type = "STRING"
        }
      }
    },
  ]

  value_type = {
    type = "JSON"
  }
}`, acctest.AuthorizePMTFSandboxEnvironment(), resourceName, name)
}

func testAccTrustFrameworkAttributeConfig_ComplexAttribute2(resourceName, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_authorize_trust_framework_attribute" "%[2]s-resolver" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s-resolver"
  description    = "Test attribute"

  value_type = {
    type = "STRING"
  }
}

resource "pingone_authorize_trust_framework_attribute" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "Test attribute"

  processor = {
    name = "%[3]s Test processor"
    type = "JSON_PATH"

    expression = "$.data.item.item1"
    value_type = {
      type = "STRING"
    }
  }

  resolvers = [
    {
      name = "%[3]s Test Current User ID Resolver"

      type = "CURRENT_USER_ID"

      condition = {
        type       = "COMPARISON"
        comparator = "EQUALS"

        left = {
          type  = "CONSTANT"
          value = "test"
        }

        right = {
          type  = "CONSTANT"
          value = "test1"
        }
      }

      processor = {
        name = "%[3]s Test processor 1"
        type = "JSON_PATH"

        expression = "$.data.item1"
        value_type = {
          type = "STRING"
        }
      }
    },
    {
      name = "%[3]s Test Attribute Resolver"

      type = "ATTRIBUTE"
      value_ref = {
        id = pingone_authorize_trust_framework_attribute.%[2]s-resolver.id
      }

      condition = {
        type = "AND"

        conditions = [
          {
            type = "EMPTY"
          },
          {
            type       = "COMPARISON"
            comparator = "EQUALS"

            left = {
              type  = "CONSTANT"
              value = "test"
            }

            right = {
              type  = "CONSTANT"
              value = "test1"
            }
          },
          {
            type = "NOT"

            condition = {
              type       = "COMPARISON"
              comparator = "EQUALS"

              left = {
                type  = "CONSTANT"
                value = "test2"
              }

              right = {
                type  = "CONSTANT"
                value = "test3"
              }
            }
          }
        ]
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
            name = "%[3]s Test chain processor 2"
            type = "JSON_PATH"

            expression = "$.data.item2"
            value_type = {
              type = "STRING"
            }
          },
          {
            name = "%[3]s Test chain processor 3"
            type = "JSON_PATH"

            expression = "$.data.item3"
            value_type = {
              type = "BOOLEAN"
            }
          },
        ],
      }
    },
    {
      name = "%[3]s Test Constant Resolver"

      type  = "CONSTANT"
      value = "test"
      value_type = {
        type = "STRING"
      }

      condition = {
        type       = "COMPARISON"
        comparator = "EQUALS"

        left = {
          type  = "CONSTANT"
          value = "test"
        }

        right = {
          type  = "CONSTANT"
          value = "test1"
        }
      }

      processor = {
        name = "%[3]s Test processor 1"
        type = "JSON_PATH"

        expression = "$.data.item1"
        value_type = {
          type = "STRING"
        }
      }
    },
  ]

  value_type = {
    type = "DURATION"
  }
}`, acctest.AuthorizePMTFSandboxEnvironment(), resourceName, name)
}

func testAccTrustFrameworkAttributeConfig_Resolver_Attribute_Full(resourceName, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_authorize_trust_framework_attribute" "%[2]s-resolver" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s-resolver"
  description    = "Test attribute"

  value_type = {
    type = "STRING"
  }
}

resource "pingone_authorize_trust_framework_attribute" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "Test attribute"

  resolvers = [
    {
      name = "%[3]s Test Attribute Resolver"

      type = "ATTRIBUTE"
      value_ref = {
        id = pingone_authorize_trust_framework_attribute.%[2]s-resolver.id
      }

      condition = {
        type = "AND"

        conditions = [
          {
            type = "EMPTY"
          },
          {
            type       = "COMPARISON"
            comparator = "EQUALS"

            left = {
              type  = "CONSTANT"
              value = "test"
            }

            right = {
              type  = "CONSTANT"
              value = "test1"
            }
          },
          {
            type = "NOT"

            condition = {
              type       = "COMPARISON"
              comparator = "EQUALS"

              left = {
                type  = "CONSTANT"
                value = "test2"
              }

              right = {
                type  = "CONSTANT"
                value = "test3"
              }
            }
          }
        ]
      }

      processor = {
        name = "%[3]s Test processor 1"
        type = "JSON_PATH"

        expression = "$.data.item1"
        value_type = {
          type = "STRING"
        }
      }
    },
  ]

  value_type = {
    type = "JSON"
  }
}`, acctest.AuthorizePMTFSandboxEnvironment(), resourceName, name)
}

func testAccTrustFrameworkAttributeConfig_Resolver_Attribute_Min(resourceName, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_authorize_trust_framework_attribute" "%[2]s-resolver" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s-resolver"
  description    = "Test attribute"

  value_type = {
    type = "STRING"
  }
}

resource "pingone_authorize_trust_framework_attribute" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "Test attribute"

  resolvers = [
    {
      type = "ATTRIBUTE"
      value_ref = {
        id = pingone_authorize_trust_framework_attribute.%[2]s-resolver.id
      }
    },
  ]

  value_type = {
    type = "JSON"
  }
}`, acctest.AuthorizePMTFSandboxEnvironment(), resourceName, name)
}

func testAccTrustFrameworkAttributeConfig_Resolver_Constant_Full(resourceName, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_authorize_trust_framework_attribute" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "Test attribute"

  resolvers = [
    {
      name = "%[3]s Test Constant Resolver"

      type = "CONSTANT"
      value_type = {
        type = "JSON"
      }
      value = jsonencode({ "foo" : "bar" })

      condition = {
        type = "AND"

        conditions = [
          {
            type = "EMPTY"
          },
          {
            type       = "COMPARISON"
            comparator = "EQUALS"

            left = {
              type  = "CONSTANT"
              value = "test"
            }

            right = {
              type  = "CONSTANT"
              value = "test1"
            }
          },
          {
            type = "NOT"

            condition = {
              type       = "COMPARISON"
              comparator = "EQUALS"

              left = {
                type  = "CONSTANT"
                value = "test2"
              }

              right = {
                type  = "CONSTANT"
                value = "test3"
              }
            }
          }
        ]
      }

      processor = {
        name = "%[3]s Test processor 1"
        type = "JSON_PATH"

        expression = "$.data.item1"
        value_type = {
          type = "STRING"
        }
      }
    },
  ]

  value_type = {
    type = "JSON"
  }
}`, acctest.AuthorizePMTFSandboxEnvironment(), resourceName, name)
}

func testAccTrustFrameworkAttributeConfig_Resolver_Constant_Min(resourceName, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_authorize_trust_framework_attribute" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "Test attribute"

  resolvers = [
    {
      type = "CONSTANT"
      value_type = {
        type = "STRING"
      }
      value = "test"
    }
  ]

  value_type = {
    type = "JSON"
  }
}`, acctest.AuthorizePMTFSandboxEnvironment(), resourceName, name)
}

func testAccTrustFrameworkAttributeConfig_Resolver_CurrentRepetitionValue_Full(resourceName, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_authorize_trust_framework_attribute" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "Test attribute"

  resolvers = [
    {
      name = "%[3]s Test Current Repetition Resolver"

      type = "CURRENT_REPETITION_VALUE"

      condition = {
        type = "AND"

        conditions = [
          {
            type = "EMPTY"
          },
          {
            type       = "COMPARISON"
            comparator = "EQUALS"

            left = {
              type  = "CONSTANT"
              value = "test"
            }

            right = {
              type  = "CONSTANT"
              value = "test1"
            }
          },
          {
            type = "NOT"

            condition = {
              type       = "COMPARISON"
              comparator = "EQUALS"

              left = {
                type  = "CONSTANT"
                value = "test2"
              }

              right = {
                type  = "CONSTANT"
                value = "test3"
              }
            }
          }
        ]
      }

      processor = {
        name = "%[3]s Test processor 1"
        type = "JSON_PATH"

        expression = "$.data.item1"
        value_type = {
          type = "STRING"
        }
      }
    },
  ]

  value_type = {
    type = "JSON"
  }
}`, acctest.AuthorizePMTFSandboxEnvironment(), resourceName, name)
}

func testAccTrustFrameworkAttributeConfig_Resolver_CurrentRepetitionValue_Min(resourceName, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_authorize_trust_framework_attribute" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "Test attribute"

  resolvers = [
    {
      type = "CURRENT_REPETITION_VALUE"
    }
  ]

  value_type = {
    type = "JSON"
  }
}`, acctest.AuthorizePMTFSandboxEnvironment(), resourceName, name)
}

func testAccTrustFrameworkAttributeConfig_Resolver_CurrentUserId_Full(resourceName, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_authorize_trust_framework_attribute" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "Test attribute"

  resolvers = [
    {
      name = "%[3]s Test Current User Id Resolver"

      type = "CURRENT_USER_ID"

      condition = {
        type = "AND"

        conditions = [
          {
            type = "EMPTY"
          },
          {
            type       = "COMPARISON"
            comparator = "EQUALS"

            left = {
              type  = "CONSTANT"
              value = "test"
            }

            right = {
              type  = "CONSTANT"
              value = "test1"
            }
          },
          {
            type = "NOT"

            condition = {
              type       = "COMPARISON"
              comparator = "EQUALS"

              left = {
                type  = "CONSTANT"
                value = "test2"
              }

              right = {
                type  = "CONSTANT"
                value = "test3"
              }
            }
          }
        ]
      }

      processor = {
        name = "%[3]s Test processor 1"
        type = "JSON_PATH"

        expression = "$.data.item1"
        value_type = {
          type = "STRING"
        }
      }
    },
  ]

  value_type = {
    type = "JSON"
  }
}`, acctest.AuthorizePMTFSandboxEnvironment(), resourceName, name)
}

func testAccTrustFrameworkAttributeConfig_Resolver_CurrentUserId_Min(resourceName, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_authorize_trust_framework_attribute" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "Test attribute"

  resolvers = [
    {
      type = "CURRENT_USER_ID"
    }
  ]

  value_type = {
    type = "JSON"
  }
}`, acctest.AuthorizePMTFSandboxEnvironment(), resourceName, name)
}

func testAccTrustFrameworkAttributeConfig_Resolver_Request_Full(resourceName, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_authorize_trust_framework_attribute" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "Test attribute"

  resolvers = [
    {
      name = "%[3]s Test Request Resolver"

      type = "REQUEST"

      condition = {
        type = "AND"

        conditions = [
          {
            type = "EMPTY"
          },
          {
            type       = "COMPARISON"
            comparator = "EQUALS"

            left = {
              type  = "CONSTANT"
              value = "test"
            }

            right = {
              type  = "CONSTANT"
              value = "test1"
            }
          },
          {
            type = "NOT"

            condition = {
              type       = "COMPARISON"
              comparator = "EQUALS"

              left = {
                type  = "CONSTANT"
                value = "test2"
              }

              right = {
                type  = "CONSTANT"
                value = "test3"
              }
            }
          }
        ]
      }

      processor = {
        name = "%[3]s Test processor 1"
        type = "JSON_PATH"

        expression = "$.data.item1"
        value_type = {
          type = "STRING"
        }
      }
    },
  ]

  value_type = {
    type = "JSON"
  }
}`, acctest.AuthorizePMTFSandboxEnvironment(), resourceName, name)
}

func testAccTrustFrameworkAttributeConfig_Resolver_Request_Min(resourceName, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_authorize_trust_framework_attribute" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "Test attribute"

  resolvers = [
    {
      type = "REQUEST"
    }
  ]

  value_type = {
    type = "JSON"
  }
}`, acctest.AuthorizePMTFSandboxEnvironment(), resourceName, name)
}

func testAccTrustFrameworkAttributeConfig_Resolver_Service_Full(resourceName, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_authorize_trust_framework_service" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "Test application service"

  service_type = "NONE"
}

resource "pingone_authorize_trust_framework_attribute" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "Test attribute"

  resolvers = [
    {
      name = "%[3]s Test Service Resolver"

      type = "SERVICE"
      value_ref = {
        id = pingone_authorize_trust_framework_service.%[2]s.id
      }

      condition = {
        type = "AND"

        conditions = [
          {
            type = "EMPTY"
          },
          {
            type       = "COMPARISON"
            comparator = "EQUALS"

            left = {
              type  = "CONSTANT"
              value = "test"
            }

            right = {
              type  = "CONSTANT"
              value = "test1"
            }
          },
          {
            type = "NOT"

            condition = {
              type       = "COMPARISON"
              comparator = "EQUALS"

              left = {
                type  = "CONSTANT"
                value = "test2"
              }

              right = {
                type  = "CONSTANT"
                value = "test3"
              }
            }
          }
        ]
      }

      processor = {
        name = "%[3]s Test processor 1"
        type = "JSON_PATH"

        expression = "$.data.item1"
        value_type = {
          type = "STRING"
        }
      }
    },
  ]

  value_type = {
    type = "JSON"
  }
}`, acctest.AuthorizePMTFSandboxEnvironment(), resourceName, name)
}

func testAccTrustFrameworkAttributeConfig_Resolver_Service_Min(resourceName, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_authorize_trust_framework_service" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "Test application service"

  service_type = "NONE"
}

resource "pingone_authorize_trust_framework_attribute" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "Test attribute"

  resolvers = [
    {
      type = "SERVICE"
      value_ref = {
        id = pingone_authorize_trust_framework_service.%[2]s.id
      }
    }
  ]

  value_type = {
    type = "JSON"
  }
}`, acctest.AuthorizePMTFSandboxEnvironment(), resourceName, name)
}

func testAccTrustFrameworkAttributeConfig_Resolver_System_Full(resourceName, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_authorize_trust_framework_attribute" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "Test attribute"

  resolvers = [
    {
      name = "%[3]s Test System Resolver"

      type  = "SYSTEM"
      value = "CURRENT_DATE_TIME"

      condition = {
        type = "AND"

        conditions = [
          {
            type = "EMPTY"
          },
          {
            type       = "COMPARISON"
            comparator = "EQUALS"

            left = {
              type  = "CONSTANT"
              value = "test"
            }

            right = {
              type  = "CONSTANT"
              value = "test1"
            }
          },
          {
            type = "NOT"

            condition = {
              type       = "COMPARISON"
              comparator = "EQUALS"

              left = {
                type  = "CONSTANT"
                value = "test2"
              }

              right = {
                type  = "CONSTANT"
                value = "test3"
              }
            }
          }
        ]
      }

      processor = {
        name = "%[3]s Test processor 1"
        type = "JSON_PATH"

        expression = "$.data.item1"
        value_type = {
          type = "STRING"
        }
      }
    },
  ]

  value_type = {
    type = "JSON"
  }
}`, acctest.AuthorizePMTFSandboxEnvironment(), resourceName, name)
}

func testAccTrustFrameworkAttributeConfig_Resolver_System_Min(resourceName, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_authorize_trust_framework_attribute" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "Test attribute"

  resolvers = [
    {
      type  = "SYSTEM"
      value = "NULL"
    }
  ]

  value_type = {
    type = "JSON"
  }
}`, acctest.AuthorizePMTFSandboxEnvironment(), resourceName, name)
}

func testAccTrustFrameworkAttributeConfig_Resolver_User_Full(resourceName, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_authorize_trust_framework_attribute" "%[2]s-user" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s-user"
  description    = "Test attribute"

  value_type = {
    type = "STRING"
  }
}

resource "pingone_authorize_trust_framework_attribute" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "Test attribute"

  resolvers = [
    {
      name = "%[3]s Test User Resolver"

      type = "USER"
      query = {
        type    = "USER_ID"
        user_id = pingone_authorize_trust_framework_attribute.%[2]s-user.id
      }

      condition = {
        type = "AND"

        conditions = [
          {
            type = "EMPTY"
          },
          {
            type       = "COMPARISON"
            comparator = "EQUALS"

            left = {
              type  = "CONSTANT"
              value = "test"
            }

            right = {
              type  = "CONSTANT"
              value = "test1"
            }
          },
          {
            type = "NOT"

            condition = {
              type       = "COMPARISON"
              comparator = "EQUALS"

              left = {
                type  = "CONSTANT"
                value = "test2"
              }

              right = {
                type  = "CONSTANT"
                value = "test3"
              }
            }
          }
        ]
      }

      processor = {
        name = "%[3]s Test processor 1"
        type = "JSON_PATH"

        expression = "$.data.item1"
        value_type = {
          type = "STRING"
        }
      }
    },
  ]

  value_type = {
    type = "JSON"
  }
}`, acctest.AuthorizePMTFSandboxEnvironment(), resourceName, name)
}

func testAccTrustFrameworkAttributeConfig_Resolver_User_Min(resourceName, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_authorize_trust_framework_attribute" "%[2]s-user" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s-user"
  description    = "Test attribute"

  value_type = {
    type = "STRING"
  }
}

resource "pingone_authorize_trust_framework_attribute" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "Test attribute"

  resolvers = [
    {
      type = "USER"
      query = {
        type    = "USER_ID"
        user_id = pingone_authorize_trust_framework_attribute.%[2]s-user.id
      }
    }
  ]

  value_type = {
    type = "STRING"
  }
}`, acctest.AuthorizePMTFSandboxEnvironment(), resourceName, name)
}
