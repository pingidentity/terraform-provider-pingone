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

func TestAccTrustFrameworkProcessor_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authorize_trust_framework_processor.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	var fido2PolicyID, environmentID string

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
		CheckDestroy:             authorize.TrustFrameworkProcessor_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccTrustFrameworkProcessorConfig_Minimal(resourceName, name),
				Check:  authorize.TrustFrameworkProcessor_GetIDs(resourceFullName, &environmentID, &fido2PolicyID),
			},
			{
				PreConfig: func() {
					authorize.TrustFrameworkProcessor_RemovalDrift_PreConfig(ctx, p1Client.API.AuthorizeAPIClient, t, environmentID, fido2PolicyID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			// Test removal of the environment
			{
				Config: testAccTrustFrameworkProcessorConfig_NewEnv(environmentName, licenseID, resourceName, name),
				Check:  authorize.TrustFrameworkProcessor_GetIDs(resourceFullName, &environmentID, &fido2PolicyID),
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

func TestAccTrustFrameworkProcessor_NewEnv(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authorize_trust_framework_processor.%s", resourceName)

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
		CheckDestroy:             authorize.TrustFrameworkProcessor_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccTrustFrameworkProcessorConfig_NewEnv(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
				),
			},
		},
	})
}

func TestAccTrustFrameworkProcessor_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authorize_trust_framework_processor.%s", resourceName)

	name := resourceName

	fullCheck := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckResourceAttr(resourceFullName, "name", name),
		resource.TestCheckResourceAttr(resourceFullName, "description", "Test processor"),
		resource.TestCheckResourceAttr(resourceFullName, "full_name", name),
		resource.TestMatchResourceAttr(resourceFullName, "parent.id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckResourceAttr(resourceFullName, "processor.name", fmt.Sprintf("%s Test child processor", name)),
		resource.TestCheckResourceAttr(resourceFullName, "type", "PROCESSOR"),
		resource.TestMatchResourceAttr(resourceFullName, "version", verify.P1ResourceIDRegexpFullString),
	)

	minimalCheck := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckResourceAttr(resourceFullName, "name", name),
		resource.TestCheckNoResourceAttr(resourceFullName, "description"),
		resource.TestCheckNoResourceAttr(resourceFullName, "full_name"),
		resource.TestCheckNoResourceAttr(resourceFullName, "parent"),
		resource.TestCheckResourceAttr(resourceFullName, "processor.name", fmt.Sprintf("%s Test processor", name)),
		resource.TestCheckResourceAttr(resourceFullName, "type", "PROCESSOR"),
		resource.TestMatchResourceAttr(resourceFullName, "version", verify.P1ResourceIDRegexpFullString),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             authorize.TrustFrameworkProcessor_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			{
				Config: testAccTrustFrameworkProcessorConfig_Full(resourceName, name),
				Check:  fullCheck,
			},
			{
				Config:  testAccTrustFrameworkProcessorConfig_Full(resourceName, name),
				Destroy: true,
			},
			// Minimal
			{
				Config: testAccTrustFrameworkProcessorConfig_Minimal(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config:  testAccTrustFrameworkProcessorConfig_Minimal(resourceName, name),
				Destroy: true,
			},
			// Change
			{
				Config: testAccTrustFrameworkProcessorConfig_Full(resourceName, name),
				Check:  fullCheck,
			},
			{
				Config: testAccTrustFrameworkProcessorConfig_Minimal(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config: testAccTrustFrameworkProcessorConfig_Full(resourceName, name),
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

func TestAccTrustFrameworkProcessor_ProcessorType_Chain(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authorize_trust_framework_processor.%s", resourceName)

	name := resourceName

	typeCheck1 := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "processor.name", fmt.Sprintf("%s Test chain processor", name)),
		resource.TestCheckResourceAttr(resourceFullName, "processor.type", "CHAIN"),
		resource.TestCheckNoResourceAttr(resourceFullName, "processor.expression"),
		resource.TestCheckNoResourceAttr(resourceFullName, "processor.predicate"),
		resource.TestCheckNoResourceAttr(resourceFullName, "processor.processor"),
		resource.TestCheckNoResourceAttr(resourceFullName, "processor.processor_ref"),
		resource.TestCheckResourceAttr(resourceFullName, "processor.processors.#", "3"),
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
		resource.TestCheckResourceAttr(resourceFullName, "processor.processors.2.value_type.type", "BOOLEAN"),
		resource.TestCheckNoResourceAttr(resourceFullName, "processor.value_type"),
	)

	typeCheck2 := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "processor.name", fmt.Sprintf("%s Test chain processor", name)),
		resource.TestCheckResourceAttr(resourceFullName, "processor.type", "CHAIN"),
		resource.TestCheckNoResourceAttr(resourceFullName, "processor.expression"),
		resource.TestCheckNoResourceAttr(resourceFullName, "processor.predicate"),
		resource.TestCheckNoResourceAttr(resourceFullName, "processor.processor"),
		resource.TestCheckNoResourceAttr(resourceFullName, "processor.processor_ref"),
		resource.TestCheckResourceAttr(resourceFullName, "processor.processors.#", "3"),
		resource.TestCheckResourceAttr(resourceFullName, "processor.processors.0.name", fmt.Sprintf("%s Test chain processor 1", name)),
		resource.TestCheckResourceAttr(resourceFullName, "processor.processors.0.type", "JSON_PATH"),
		resource.TestCheckResourceAttr(resourceFullName, "processor.processors.0.expression", "$.data.item1"),
		resource.TestCheckResourceAttr(resourceFullName, "processor.processors.0.value_type.type", "STRING"),
		resource.TestCheckResourceAttr(resourceFullName, "processor.processors.1.name", fmt.Sprintf("%s Test chain processor 3", name)),
		resource.TestCheckResourceAttr(resourceFullName, "processor.processors.1.type", "JSON_PATH"),
		resource.TestCheckResourceAttr(resourceFullName, "processor.processors.1.expression", "$.data.item3"),
		resource.TestCheckResourceAttr(resourceFullName, "processor.processors.1.value_type.type", "BOOLEAN"),
		resource.TestCheckResourceAttr(resourceFullName, "processor.processors.2.name", fmt.Sprintf("%s Test chain processor 2", name)),
		resource.TestCheckResourceAttr(resourceFullName, "processor.processors.2.type", "JSON_PATH"),
		resource.TestCheckResourceAttr(resourceFullName, "processor.processors.2.expression", "$.data.item2"),
		resource.TestCheckResourceAttr(resourceFullName, "processor.processors.2.value_type.type", "STRING"),
		resource.TestCheckNoResourceAttr(resourceFullName, "processor.value_type"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             authorize.TrustFrameworkProcessor_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// From scratch
			{
				Config: testAccTrustFrameworkProcessorConfig_Processor_Chain1(resourceName, name),
				Check:  typeCheck1,
			},
			// Change order of processors
			{
				Config: testAccTrustFrameworkProcessorConfig_Processor_Chain2(resourceName, name),
				Check:  typeCheck2,
			},
			{
				Config: testAccTrustFrameworkProcessorConfig_Processor_Chain1(resourceName, name),
				Check:  typeCheck1,
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

func TestAccTrustFrameworkProcessor_ProcessorType_CollectionFilter(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authorize_trust_framework_processor.%s", resourceName)

	name := resourceName

	typeCheck1 := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "processor.name", fmt.Sprintf("%s Test processor", name)),
		resource.TestCheckResourceAttr(resourceFullName, "processor.type", "COLLECTION_FILTER"),
		resource.TestCheckNoResourceAttr(resourceFullName, "processor.expression"),
		resource.TestCheckResourceAttr(resourceFullName, "processor.predicate.name", fmt.Sprintf("%s Test predicate processor 1", name)),
		resource.TestCheckResourceAttr(resourceFullName, "processor.predicate.type", "JSON_PATH"),
		resource.TestCheckResourceAttr(resourceFullName, "processor.predicate.expression", "$.data.item1"),
		resource.TestCheckResourceAttr(resourceFullName, "processor.predicate.value_type", "STRING"),
		resource.TestCheckNoResourceAttr(resourceFullName, "processor.processor"),
		resource.TestCheckNoResourceAttr(resourceFullName, "processor.processor_ref"),
		resource.TestCheckNoResourceAttr(resourceFullName, "processor.processors"),
		resource.TestCheckNoResourceAttr(resourceFullName, "processor.value_type"),
	)

	typeCheck2 := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "processor.name", fmt.Sprintf("%s Test processor", name)),
		resource.TestCheckResourceAttr(resourceFullName, "processor.type", "COLLECTION_FILTER"),
		resource.TestCheckNoResourceAttr(resourceFullName, "processor.expression"),
		resource.TestCheckResourceAttr(resourceFullName, "processor.predicate.name", fmt.Sprintf("%s Test predicate processor 2", name)),
		resource.TestCheckResourceAttr(resourceFullName, "processor.predicate.type", "SPEL"),
		resource.TestCheckResourceAttr(resourceFullName, "processor.predicate.expression", "'Hello SpEL'.concat('!')"),
		resource.TestCheckResourceAttr(resourceFullName, "processor.predicate.value_type", "STRING"),
		resource.TestCheckNoResourceAttr(resourceFullName, "processor.processor"),
		resource.TestCheckNoResourceAttr(resourceFullName, "processor.processor_ref"),
		resource.TestCheckNoResourceAttr(resourceFullName, "processor.processors"),
		resource.TestCheckNoResourceAttr(resourceFullName, "processor.value_type"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             authorize.TrustFrameworkProcessor_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// From scratch
			{
				Config: testAccTrustFrameworkProcessorConfig_Processor_CollectionFilter1(resourceName, name),
				Check:  typeCheck1,
			},
			// Change predicate processor
			{
				Config: testAccTrustFrameworkProcessorConfig_Processor_CollectionFilter2(resourceName, name),
				Check:  typeCheck2,
			},
			{
				Config: testAccTrustFrameworkProcessorConfig_Processor_CollectionFilter1(resourceName, name),
				Check:  typeCheck1,
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

func TestAccTrustFrameworkProcessor_ProcessorType_CollectionTransform(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authorize_trust_framework_processor.%s", resourceName)

	name := resourceName

	typeCheck1 := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "processor.name", fmt.Sprintf("%s Test processor", name)),
		resource.TestCheckResourceAttr(resourceFullName, "processor.type", "COLLECTION_TRANSFORM"),
		resource.TestCheckNoResourceAttr(resourceFullName, "processor.expression"),
		resource.TestCheckNoResourceAttr(resourceFullName, "processor.predicate"),
		resource.TestCheckResourceAttr(resourceFullName, "processor.processor.name", fmt.Sprintf("%s Test collection transform processor 1", name)),
		resource.TestCheckResourceAttr(resourceFullName, "processor.processor.type", "JSON_PATH"),
		resource.TestCheckResourceAttr(resourceFullName, "processor.processor.expression", "$.data.item1"),
		resource.TestCheckResourceAttr(resourceFullName, "processor.processor.value_type", "STRING"),
		resource.TestCheckNoResourceAttr(resourceFullName, "processor.processor_ref"),
		resource.TestCheckNoResourceAttr(resourceFullName, "processor.processors"),
		resource.TestCheckNoResourceAttr(resourceFullName, "processor.value_type"),
	)

	typeCheck2 := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "processor.name", fmt.Sprintf("%s Test processor", name)),
		resource.TestCheckResourceAttr(resourceFullName, "processor.type", "COLLECTION_TRANSFORM"),
		resource.TestCheckNoResourceAttr(resourceFullName, "processor.expression"),
		resource.TestCheckNoResourceAttr(resourceFullName, "processor.predicate"),
		resource.TestCheckResourceAttr(resourceFullName, "processor.processor.name", fmt.Sprintf("%s Test collection transform processor 2", name)),
		resource.TestCheckResourceAttr(resourceFullName, "processor.processor.type", "SPEL"),
		resource.TestCheckResourceAttr(resourceFullName, "processor.processor.expression", "'Hello SpEL'.concat('!')"),
		resource.TestCheckResourceAttr(resourceFullName, "processor.processor.value_type", "STRING"),
		resource.TestCheckNoResourceAttr(resourceFullName, "processor.processor_ref"),
		resource.TestCheckNoResourceAttr(resourceFullName, "processor.processors"),
		resource.TestCheckNoResourceAttr(resourceFullName, "processor.value_type"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             authorize.TrustFrameworkProcessor_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// From scratch
			{
				Config: testAccTrustFrameworkProcessorConfig_Processor_CollectionTransform1(resourceName, name),
				Check:  typeCheck1,
			},
			// Change predicate processor
			{
				Config: testAccTrustFrameworkProcessorConfig_Processor_CollectionTransform2(resourceName, name),
				Check:  typeCheck2,
			},
			{
				Config: testAccTrustFrameworkProcessorConfig_Processor_CollectionTransform1(resourceName, name),
				Check:  typeCheck1,
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

func TestAccTrustFrameworkProcessor_ProcessorType_Json_Path(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authorize_trust_framework_processor.%s", resourceName)

	name := resourceName

	typeCheck1 := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "processor.name", fmt.Sprintf("%s Test processor", name)),
		resource.TestCheckResourceAttr(resourceFullName, "processor.type", "JSON_PATH"),
		resource.TestCheckResourceAttr(resourceFullName, "processor.expression", "$.data.item"),
		resource.TestCheckNoResourceAttr(resourceFullName, "processor.predicate"),
		resource.TestCheckNoResourceAttr(resourceFullName, "processor.processor"),
		resource.TestCheckNoResourceAttr(resourceFullName, "processor.processor_ref"),
		resource.TestCheckNoResourceAttr(resourceFullName, "processor.processors"),
		resource.TestCheckResourceAttr(resourceFullName, "processor.value_type.type", "STRING"),
	)

	typeCheck2 := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "processor.name", fmt.Sprintf("%s Test processor", name)),
		resource.TestCheckResourceAttr(resourceFullName, "processor.type", "JSON_PATH"),
		resource.TestCheckResourceAttr(resourceFullName, "processor.expression", "$.data.item2"),
		resource.TestCheckNoResourceAttr(resourceFullName, "processor.predicate"),
		resource.TestCheckNoResourceAttr(resourceFullName, "processor.processor"),
		resource.TestCheckNoResourceAttr(resourceFullName, "processor.processor_ref"),
		resource.TestCheckNoResourceAttr(resourceFullName, "processor.processors"),
		resource.TestCheckResourceAttr(resourceFullName, "processor.value_type.type", "BOOLEAN"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             authorize.TrustFrameworkProcessor_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// From scratch
			{
				Config: testAccTrustFrameworkProcessorConfig_Processor_Json_Path1(resourceName, name),
				Check:  typeCheck1,
			},
			// Change predicate processor
			{
				Config: testAccTrustFrameworkProcessorConfig_Processor_Json_Path2(resourceName, name),
				Check:  typeCheck2,
			},
			{
				Config: testAccTrustFrameworkProcessorConfig_Processor_Json_Path1(resourceName, name),
				Check:  typeCheck1,
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

func TestAccTrustFrameworkProcessor_ProcessorType_Reference(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authorize_trust_framework_processor.%s", resourceName)

	name := resourceName

	typeCheck1 := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "processor.name", fmt.Sprintf("%s Test processor", name)),
		resource.TestCheckResourceAttr(resourceFullName, "processor.type", "REFERENCE"),
		resource.TestCheckNoResourceAttr(resourceFullName, "processor.expression"),
		resource.TestCheckNoResourceAttr(resourceFullName, "processor.predicate"),
		resource.TestCheckNoResourceAttr(resourceFullName, "processor.processor"),
		resource.TestMatchResourceAttr(resourceFullName, "processor.processor_ref", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckNoResourceAttr(resourceFullName, "processor.processors"),
		resource.TestCheckNoResourceAttr(resourceFullName, "processor.value_type"),
	)

	typeCheck2 := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "processor.name", fmt.Sprintf("%s Test processor", name)),
		resource.TestCheckResourceAttr(resourceFullName, "processor.type", "REFERENCE"),
		resource.TestCheckNoResourceAttr(resourceFullName, "processor.expression"),
		resource.TestCheckNoResourceAttr(resourceFullName, "processor.predicate"),
		resource.TestCheckNoResourceAttr(resourceFullName, "processor.processor"),
		resource.TestMatchResourceAttr(resourceFullName, "processor.processor_ref", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckNoResourceAttr(resourceFullName, "processor.processors"),
		resource.TestCheckNoResourceAttr(resourceFullName, "processor.value_type"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             authorize.TrustFrameworkProcessor_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// From scratch
			{
				Config: testAccTrustFrameworkProcessorConfig_Processor_Reference1(resourceName, name),
				Check:  typeCheck1,
			},
			// Change predicate processor
			{
				Config: testAccTrustFrameworkProcessorConfig_Processor_Reference2(resourceName, name),
				Check:  typeCheck2,
			},
			{
				Config: testAccTrustFrameworkProcessorConfig_Processor_Reference1(resourceName, name),
				Check:  typeCheck1,
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

func TestAccTrustFrameworkProcessor_ProcessorType_SpEL(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authorize_trust_framework_processor.%s", resourceName)

	name := resourceName

	typeCheck1 := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "processor.name", fmt.Sprintf("%s Test processor", name)),
		resource.TestCheckResourceAttr(resourceFullName, "processor.type", "SPEL"),
		resource.TestCheckResourceAttr(resourceFullName, "processor.expression", "'Hello SpEL'.concat('!')"),
		resource.TestCheckNoResourceAttr(resourceFullName, "processor.predicate"),
		resource.TestCheckNoResourceAttr(resourceFullName, "processor.processor"),
		resource.TestCheckNoResourceAttr(resourceFullName, "processor.processor_ref"),
		resource.TestCheckNoResourceAttr(resourceFullName, "processor.processors"),
		resource.TestCheckResourceAttr(resourceFullName, "processor.value_type.type", "STRING"),
	)

	typeCheck2 := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "processor.name", fmt.Sprintf("%s Test processor", name)),
		resource.TestCheckResourceAttr(resourceFullName, "processor.type", "SPEL"),
		resource.TestCheckResourceAttr(resourceFullName, "processor.expression", "'Hello world'.concat('!')"),
		resource.TestCheckNoResourceAttr(resourceFullName, "processor.predicate"),
		resource.TestCheckNoResourceAttr(resourceFullName, "processor.processor"),
		resource.TestCheckNoResourceAttr(resourceFullName, "processor.processor_ref"),
		resource.TestCheckNoResourceAttr(resourceFullName, "processor.processors"),
		resource.TestCheckResourceAttr(resourceFullName, "processor.value_type.type", "STRING"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             authorize.TrustFrameworkProcessor_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// From scratch
			{
				Config: testAccTrustFrameworkProcessorConfig_Processor_SpEL1(resourceName, name),
				Check:  typeCheck1,
			},
			// Change predicate processor
			{
				Config: testAccTrustFrameworkProcessorConfig_Processor_SpEL2(resourceName, name),
				Check:  typeCheck2,
			},
			{
				Config: testAccTrustFrameworkProcessorConfig_Processor_SpEL1(resourceName, name),
				Check:  typeCheck1,
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

func TestAccTrustFrameworkProcessor_ProcessorType_XPath(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authorize_trust_framework_processor.%s", resourceName)

	name := resourceName

	typeCheck1 := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "processor.name", fmt.Sprintf("%s Test processor", name)),
		resource.TestCheckResourceAttr(resourceFullName, "processor.type", "XPATH"),
		resource.TestCheckResourceAttr(resourceFullName, "processor.expression", "/bookstore/book[last()]"),
		resource.TestCheckNoResourceAttr(resourceFullName, "processor.predicate"),
		resource.TestCheckNoResourceAttr(resourceFullName, "processor.processor"),
		resource.TestCheckNoResourceAttr(resourceFullName, "processor.processor_ref"),
		resource.TestCheckNoResourceAttr(resourceFullName, "processor.processors"),
		resource.TestCheckResourceAttr(resourceFullName, "processor.value_type.type", "STRING"),
	)

	typeCheck2 := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "processor.name", fmt.Sprintf("%s Test processor", name)),
		resource.TestCheckResourceAttr(resourceFullName, "processor.type", "XPATH"),
		resource.TestCheckResourceAttr(resourceFullName, "processor.expression", "/bookstore/book[first()]"),
		resource.TestCheckNoResourceAttr(resourceFullName, "processor.predicate"),
		resource.TestCheckNoResourceAttr(resourceFullName, "processor.processor"),
		resource.TestCheckNoResourceAttr(resourceFullName, "processor.processor_ref"),
		resource.TestCheckNoResourceAttr(resourceFullName, "processor.processors"),
		resource.TestCheckResourceAttr(resourceFullName, "processor.value_type.type", "BOOLEAN"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             authorize.TrustFrameworkProcessor_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// From scratch
			{
				Config: testAccTrustFrameworkProcessorConfig_Processor_XPath1(resourceName, name),
				Check:  typeCheck1,
			},
			// Change predicate processor
			{
				Config: testAccTrustFrameworkProcessorConfig_Processor_XPath2(resourceName, name),
				Check:  typeCheck2,
			},
			{
				Config: testAccTrustFrameworkProcessorConfig_Processor_XPath1(resourceName, name),
				Check:  typeCheck1,
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

func TestAccTrustFrameworkProcessor_ProcessorType_Change(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authorize_trust_framework_processor.%s", resourceName)

	name := resourceName

	typeCheck1 := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "processor.name", fmt.Sprintf("%s Test chain processor", name)),
		resource.TestCheckResourceAttr(resourceFullName, "processor.type", "CHAIN"),
		resource.TestCheckNoResourceAttr(resourceFullName, "processor.expression"),
		resource.TestCheckNoResourceAttr(resourceFullName, "processor.predicate"),
		resource.TestCheckNoResourceAttr(resourceFullName, "processor.processor"),
		resource.TestCheckNoResourceAttr(resourceFullName, "processor.processor_ref"),
		resource.TestCheckResourceAttr(resourceFullName, "processor.processors.#", "3"),
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
		resource.TestCheckResourceAttr(resourceFullName, "processor.processors.2.value_type.type", "BOOLEAN"),
		resource.TestCheckNoResourceAttr(resourceFullName, "processor.value_type"),
	)

	typeCheck2 := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "processor.name", fmt.Sprintf("%s Test processor", name)),
		resource.TestCheckResourceAttr(resourceFullName, "processor.type", "XPATH"),
		resource.TestCheckResourceAttr(resourceFullName, "processor.expression", "/bookstore/book[first()]"),
		resource.TestCheckNoResourceAttr(resourceFullName, "processor.predicate"),
		resource.TestCheckNoResourceAttr(resourceFullName, "processor.processor"),
		resource.TestCheckNoResourceAttr(resourceFullName, "processor.processor_ref"),
		resource.TestCheckNoResourceAttr(resourceFullName, "processor.processors"),
		resource.TestCheckResourceAttr(resourceFullName, "processor.value_type.type", "BOOLEAN"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             authorize.TrustFrameworkProcessor_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// From scratch
			{
				Config: testAccTrustFrameworkProcessorConfig_Processor_Chain1(resourceName, name),
				Check:  typeCheck1,
			},
			// Change processor
			{
				Config: testAccTrustFrameworkProcessorConfig_Processor_XPath2(resourceName, name),
				Check:  typeCheck2,
			},
			{
				Config: testAccTrustFrameworkProcessorConfig_Processor_Chain1(resourceName, name),
				Check:  typeCheck1,
			},
		},
	})
}

func TestAccTrustFrameworkProcessor_BadParameters(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authorize_trust_framework_processor.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             authorize.TrustFrameworkProcessor_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccTrustFrameworkProcessorConfig_Minimal(resourceName, name),
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

func testAccTrustFrameworkProcessorConfig_NewEnv(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_authorize_trust_framework_processor" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  name           = "%[3]s"

  processor = {
    name = "%[3]s Test processor"
    type = "JSON_PATH"

    expression = "$.data.item"
    value_type = {
      type = "STRING"
    }
  }
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccTrustFrameworkProcessorConfig_Full(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_authorize_trust_framework_processor" "%[2]s-parent" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  processor = {
    name = "%[3]s Test parent processor"
    type = "JSON_PATH"

    expression = "$.data.item.parent"
    value_type = {
      type = "STRING"
    }
  }
}

resource "pingone_authorize_trust_framework_processor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "Test processor"
  full_name      = "%[3]s"

  parent = {
    id = pingone_authorize_trust_framework_processor.%[2]s-parent.id
  }

  processor = {
    name = "%[3]s Test child processor"
    type = "JSON_PATH"

    expression = "$.data.item2"
    value_type = {
      type = "STRING"
    }
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccTrustFrameworkProcessorConfig_Minimal(resourceName, name string) string {
	return testAccTrustFrameworkProcessorConfig_Processor_Json_Path1(resourceName, name)
}

func testAccTrustFrameworkProcessorConfig_Processor_Chain1(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_authorize_trust_framework_processor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

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
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccTrustFrameworkProcessorConfig_Processor_Chain2(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_authorize_trust_framework_processor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

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
          type = "BOOLEAN"
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
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccTrustFrameworkProcessorConfig_Processor_CollectionFilter1(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_authorize_trust_framework_processor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  processor = {
    name = "%[3]s Test processor"
    type = "COLLECTION_FILTER"

    predicate = {
      name = "%[3]s Test predicate processor 1"
      type = "JSON_PATH"

      expression = "$.data.item1"
      value_type = {
        type = "STRING"
      }
    },
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccTrustFrameworkProcessorConfig_Processor_CollectionFilter2(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_authorize_trust_framework_processor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  processor = {
    name = "%[3]s Test processor"
    type = "COLLECTION_FILTER"

    predicate = {
      name = "%[3]s Test predicate processor 2"
      type = "SPEL"

      expression = "'Hello SpEL'.concat('!')"
      value_type = {
        type = "STRING"
      }
    },
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccTrustFrameworkProcessorConfig_Processor_CollectionTransform1(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_authorize_trust_framework_processor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  processor = {
    name = "%[3]s Test processor"
    type = "COLLECTION_TRANSFORM"

    processor = {
      name = "%[3]s Test collection transform processor 1"
      type = "JSON_PATH"

      expression = "$.data.item1"
      value_type = {
        type = "STRING"
      }
    },
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccTrustFrameworkProcessorConfig_Processor_CollectionTransform2(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_authorize_trust_framework_processor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  processor = {
    name = "%[3]s Test processor"
    type = "COLLECTION_TRANSFORM"

    processor = {
      name = "%[3]s Test collection transform processor 2"
      type = "SPEL"

      expression = "'Hello SpEL'.concat('!')"
      value_type = {
        type = "STRING"
      }
    },
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccTrustFrameworkProcessorConfig_Processor_Json_Path1(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_authorize_trust_framework_processor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  processor = {
    name = "%[3]s Test processor"
    type = "JSON_PATH"

    expression = "$.data.item"
    value_type = {
      type = "STRING"
    }
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccTrustFrameworkProcessorConfig_Processor_Json_Path2(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_authorize_trust_framework_processor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  processor = {
    name = "%[3]s Test processor"
    type = "JSON_PATH"

    expression = "$.data.item2"
    value_type = {
      type = "BOOLEAN"
    }
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccTrustFrameworkProcessorConfig_Processor_Reference1(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_authorize_trust_framework_processor" "%[2]s_ref1" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  processor = {
    name = "%[3]s Test ref processor1"
    type = "JSON_PATH"

    expression = "$.data.item"
    value_type = {
      type = "STRING"
    }
  }
}

resource "pingone_authorize_trust_framework_processor" "%[2]s_ref2" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  processor = {
    name = "%[3]s Test ref processor2"
    type = "JSON_PATH"

    expression = "$.data.item1"
    value_type = {
      type = "BOOLEAN"
    }
  }
}

resource "pingone_authorize_trust_framework_processor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  processor = {
    name = "%[3]s Test processor"
    type = "REFERENCE"

    processor_ref = {
      id = pingone_authorize_trust_framework_processor.%[2]s_ref1.id
    }
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccTrustFrameworkProcessorConfig_Processor_Reference2(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_authorize_trust_framework_processor" "%[2]s_ref1" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  processor = {
    name = "%[3]s Test ref processor1"
    type = "JSON_PATH"

    expression = "$.data.item"
    value_type = {
      type = "STRING"
    }
  }
}

resource "pingone_authorize_trust_framework_processor" "%[2]s_ref2" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  processor = {
    name = "%[3]s Test ref processor2"
    type = "JSON_PATH"

    expression = "$.data.item1"
    value_type = {
      type = "BOOLEAN"
    }
  }
}

resource "pingone_authorize_trust_framework_processor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  processor = {
    name = "%[3]s Test processor"
    type = "REFERENCE"

    processor_ref = {
      id = pingone_authorize_trust_framework_processor.%[2]s_ref2.id
    }
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccTrustFrameworkProcessorConfig_Processor_SpEL1(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_authorize_trust_framework_processor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  processor = {
    name = "%[3]s Test processor"
    type = "SPEL"

    expression = "'Hello SpEL'.concat('!')"
    value_type = {
      type = "STRING"
    }
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccTrustFrameworkProcessorConfig_Processor_SpEL2(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_authorize_trust_framework_processor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  processor = {
    name = "%[3]s Test processor"
    type = "SPEL"

    expression = "'Hello world'.concat('!')"
    value_type = {
      type = "STRING"
    }
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccTrustFrameworkProcessorConfig_Processor_XPath1(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_authorize_trust_framework_processor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  processor = {
    name = "%[3]s Test processor"
    type = "XPATH"

    expression = "/bookstore/book[last()]"
    value_type = {
      type = "STRING"
    }
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccTrustFrameworkProcessorConfig_Processor_XPath2(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_authorize_trust_framework_processor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  processor = {
    name = "%[3]s Test processor"
    type = "XPATH"

    expression = "/bookstore/book[first()]"
    value_type = {
      type = "BOOLEAN"
    }
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}
