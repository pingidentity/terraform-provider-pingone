// Copyright © 2025 Ping Identity Corporation

package authorize_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/authorize"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccTrustFrameworkAttributeDataSource_ByID(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authorize_trust_framework_attribute.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.%s", resourceFullName)

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
			{
				Config: testAccTrustFrameworkAttributeDataSource_ByID(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "environment_id", resourceFullName, "environment_id"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "name", resourceFullName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "description", resourceFullName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "full_name", resourceFullName, "full_name"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "parent", resourceFullName, "parent"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "default_value", resourceFullName, "default_value"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "processor", resourceFullName, "processor"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "resolvers", resourceFullName, "resolvers"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "managed_entity", resourceFullName, "managed_entity"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "repetition_source", resourceFullName, "repetition_source"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "type", resourceFullName, "type"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "value_schema", resourceFullName, "value_schema"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "value_type", resourceFullName, "value_type"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "version", resourceFullName, "version"),
				),
			},
		},
	})
}

func TestAccTrustFrameworkAttributeDataSource_ByName(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authorize_trust_framework_attribute.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.%s", resourceFullName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckFeatureFlag(t, acctest.ENUMFEATUREFLAG_AUTHORIZEPMTF)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             nil,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccTrustFrameworkAttributeDataSource_ByName(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullName, "name", "User"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "description"),
					resource.TestCheckResourceAttr(dataSourceFullName, "full_name", "PingOne.User"),
					resource.TestMatchResourceAttr(dataSourceFullName, "parent.id", verify.P1ResourceIDRegexpFullString),
					// resource.TestCheckResourceAttr(dataSourceFullName, "default_value", "test"),
					// resource.TestCheckResourceAttr(dataSourceFullName, "processor.name", fmt.Sprintf("%s Test processor", name)),
					// resource.TestCheckResourceAttrSet(dataSourceFullName, "resolvers"),
					resource.TestCheckResourceAttr(dataSourceFullName, "managed_entity.restrictions.read_only", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "managed_entity.restrictions.disallow_children", "false"),
					resource.TestCheckNoResourceAttr(dataSourceFullName, "repetition_source"),
					resource.TestCheckResourceAttr(dataSourceFullName, "type", "ATTRIBUTE"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "value_schema"),
					resource.TestCheckResourceAttr(dataSourceFullName, "value_type.type", "JSON"),
					resource.TestMatchResourceAttr(dataSourceFullName, "version", verify.P1ResourceIDRegexpFullString),
				),
			},
		},
	})
}

func TestAccTrustFrameworkAttributeDataSource_FailureChecks(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
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
			{
				Config:      testAccTrustFrameworkAttributeDataSource_FindByIDFail(resourceName, name),
				ExpectError: regexp.MustCompile("Error: Error when calling `GetAttribute`: The request could not be completed. The requested resource was not found"),
			},
			{
				Config:      testAccTrustFrameworkAttributeDataSource_FindByNameFail(resourceName, name),
				ExpectError: regexp.MustCompile("Error: Cannot find the trust framework attribute from the full name"),
			},
		},
	})
}

func testAccTrustFrameworkAttributeDataSource_ByID(resourceName, name string) string {
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
    type = "JSON"
  }

  value_schema = <<EOF
{
	"$schema": "http://json-schema.org/draft-04/schema#",
	"type ": "object"
}
EOF
}

data "pingone_authorize_trust_framework_attribute" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  attribute_id   = pingone_authorize_trust_framework_attribute.%[2]s.id
}`, acctest.AuthorizePMTFSandboxEnvironment(), resourceName, name)
}

func testAccTrustFrameworkAttributeDataSource_ByName(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

data "pingone_authorize_trust_framework_attribute" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  full_name      = "PingOne.User"
}`, acctest.AuthorizePMTFSandboxEnvironment(), resourceName, name)
}

func testAccTrustFrameworkAttributeDataSource_FindByIDFail(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_authorize_trust_framework_attribute" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  attribute_id   = "9c052a8a-14be-44e4-8f07-2662569994ce" // dummy ID that conforms to UUID v4


}`, acctest.AuthorizePMTFSandboxEnvironment(), resourceName, name)
}

func testAccTrustFrameworkAttributeDataSource_FindByNameFail(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_authorize_trust_framework_attribute" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  full_name      = "%[3]s"

}`, acctest.AuthorizePMTFSandboxEnvironment(), resourceName, name)
}
