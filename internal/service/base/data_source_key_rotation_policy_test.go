// Copyright © 2026 Ping Identity Corporation

package base_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/base"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccKeyRotationPolicyDataSource_ByNameFull(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	dataSourceFullName := fmt.Sprintf("data.pingone_key_rotation_policy.%s", resourceName)

	// Every environment is provisioned with a default key rotation policy
	// named after the environment, which this test looks up without
	// creating anything.
	environmentName := "tf-testacc-dynamic-general-test"
	policyName := fmt.Sprintf("PingOne Key Rotation Policy for %s environment", environmentName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccKeyRotationPolicyDataSourceConfig_ByNameDefault(resourceName, environmentName, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "key_rotation_policy_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullName, "name", policyName),
					resource.TestCheckResourceAttr(dataSourceFullName, "algorithm", "RSA"),
					resource.TestCheckResourceAttr(dataSourceFullName, "key_length", "2048"),
					resource.TestCheckResourceAttr(dataSourceFullName, "rotation_period", "90"),
					resource.TestCheckResourceAttr(dataSourceFullName, "signature_algorithm", "SHA256withRSA"),
					resource.TestCheckResourceAttr(dataSourceFullName, "usage_type", "SIGNING"),
					resource.TestCheckResourceAttr(dataSourceFullName, "validity_period", "365"),
					resource.TestMatchResourceAttr(dataSourceFullName, "current_key_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "next_key_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "rotated_at", regexp.MustCompile(`^(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d:[0-5]\d\.\d+)|(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d:[0-5]\d)|(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d)$`)),
				),
			},
			// Case insensitivity check
			{
				Config: testAccKeyRotationPolicyDataSourceConfig_ByNameDefault(resourceName, environmentName, true), Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceFullName, "name", policyName),
				),
			},
		},
	})
}

func TestAccKeyRotationPolicyDataSource_ByIDFull(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_key_rotation_policy.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.%s", resourceFullName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.KeyRotationPolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccKeyRotationPolicyDataSourceConfig_ByIDFull(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "key_rotation_policy_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullName, "name", name),
					resource.TestCheckResourceAttr(dataSourceFullName, "algorithm", "RSA"),
					resource.TestCheckResourceAttr(dataSourceFullName, "subject_dn", fmt.Sprintf("CN=%s, OU=Ping Identity, O=Ping Identity, L=, ST=, C=US", name)),
					resource.TestCheckResourceAttr(dataSourceFullName, "key_length", "3072"),
					resource.TestCheckResourceAttr(dataSourceFullName, "rotation_period", "31"),
					resource.TestCheckResourceAttr(dataSourceFullName, "signature_algorithm", "SHA256withRSA"),
					resource.TestCheckResourceAttr(dataSourceFullName, "usage_type", "SIGNING"),
					resource.TestCheckResourceAttr(dataSourceFullName, "validity_period", "340"),
					resource.TestMatchResourceAttr(dataSourceFullName, "current_key_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "next_key_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "rotated_at", regexp.MustCompile(`^(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d:[0-5]\d\.\d+)|(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d:[0-5]\d)|(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d)$`)),
				),
			},
			{
				Config:  testAccKeyRotationPolicyDataSourceConfig_ByIDFull(resourceName, name),
				Destroy: true,
			},
		},
	})
}

func TestAccKeyRotationPolicyDataSource_NotFound(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.KeyRotationPolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccKeyRotationPolicyDataSourceConfig_NotFoundByName(resourceName),
				ExpectError: regexp.MustCompile(`Cannot find key rotation policy from name`),
			},
			{
				Config:      testAccKeyRotationPolicyDataSourceConfig_NotFoundByID(resourceName),
				ExpectError: regexp.MustCompile("Error when calling `GetKeyRotationPolicy`: Could not find Key Rotation Policy"),
			},
		},
	})
}

func testAccKeyRotationPolicyDataSourceConfig_ByNameDefault(resourceName, environmentName string, insensitivityCheck bool) string {
	// If insensitivityCheck is true, alter the case of the name
	nameComparator := fmt.Sprintf("PingOne Key Rotation Policy for %s environment", environmentName)
	if insensitivityCheck {
		nameComparator = acctest.AlterStringCasing(nameComparator)
	}

	return fmt.Sprintf(`
	%[1]s

data "pingone_key_rotation_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}
`, acctest.GenericSandboxEnvironment(), resourceName, nameComparator)
}

func testAccKeyRotationPolicyDataSourceConfig_ByIDFull(resourceName, name string) string {
	return fmt.Sprintf(`
%s

data "pingone_key_rotation_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  key_rotation_policy_id = pingone_key_rotation_policy.%[2]s.id
}`, testAccKeyRotationPolicyConfig_Full(resourceName, name), resourceName)
}

func testAccKeyRotationPolicyDataSourceConfig_NotFoundByName(resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_key_rotation_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "doesnotexist"
}
`, acctest.GenericSandboxEnvironment(), resourceName)
}

func testAccKeyRotationPolicyDataSourceConfig_NotFoundByID(resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_key_rotation_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  key_rotation_policy_id = "9c052a8a-14be-44e4-8f07-2662569994ce" // dummy ID that conforms to UUID v4
}
`, acctest.GenericSandboxEnvironment(), resourceName)
}
