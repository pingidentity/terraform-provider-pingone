// Copyright © 2025 Ping Identity Corporation

package sso_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/sso"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccResourceSecretDataSource_Basic(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_resource_secret.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.%s", resourceFullName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {
				VersionConstraint: "0.11.1",
				Source:            "hashicorp/time",
			},
		},
		CheckDestroy: sso.ResourceSecret_CheckDestroy,
		ErrorCheck:   acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSecretDataSourceConfig_Full(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "resource_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "previous.secret", regexp.MustCompile(`[a-zA-Z0-9-~_]{10,}`)),
					resource.TestMatchResourceAttr(dataSourceFullName, "previous.expires_at", verify.RFC3339Regexp),
					resource.TestCheckNoResourceAttr(dataSourceFullName, "previous.last_used"),
					resource.TestMatchResourceAttr(dataSourceFullName, "secret", regexp.MustCompile(`[a-zA-Z0-9-~_]{10,}`)),
				),
			},
		},
	})
}

func TestAccResourceSecretDataSource_Rotation(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_resource_secret.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.%s", resourceFullName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {
				VersionConstraint: "0.11.1",
				Source:            "hashicorp/time",
			},
		},
		CheckDestroy: sso.ResourceSecret_CheckDestroy,
		ErrorCheck:   acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSecretDataSourceConfig_Rotation1(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr(dataSourceFullName, "previous.secret"),
					resource.TestCheckNoResourceAttr(dataSourceFullName, "previous.expires_at"),
					resource.TestCheckNoResourceAttr(dataSourceFullName, "previous.last_used"),
					resource.TestMatchResourceAttr(dataSourceFullName, "secret", regexp.MustCompile(`[a-zA-Z0-9-~_]{10,}`)),
				),
			},
			{
				Config: testAccResourceSecretDataSourceConfig_Rotation2(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "previous.secret", regexp.MustCompile(`[a-zA-Z0-9-~_]{10,}`)),
					resource.TestMatchResourceAttr(dataSourceFullName, "previous.expires_at", verify.RFC3339Regexp),
					resource.TestCheckNoResourceAttr(dataSourceFullName, "previous.last_used"),
					resource.TestMatchResourceAttr(dataSourceFullName, "secret", regexp.MustCompile(`[a-zA-Z0-9-~_]{10,}`)),
				),
			},
		},
	})
}

func TestAccResourceSecretDataSource_IncorrectResourceType(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.ResourceSecret_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccResourceSecretDataSourceConfig_IncorrectResourceType(resourceName, name),
				ExpectError: regexp.MustCompile("Error when calling `ReadResourceSecret`: The request could not be completed. The requested resource was not found."),
			},
		},
	})
}

func TestAccResourceSecretDataSource_NotFound(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.ResourceSecret_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccResourceSecretDataSourceConfig_NotFoundByID(resourceName),
				ExpectError: regexp.MustCompile("Error when calling `ReadResourceSecret`: The request could not be completed. The requested resource was not found."),
			},
		},
	})
}

func testAccResourceSecretDataSourceConfig_Full(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_resource" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "time_offset" "%[2]s" {
  offset_minutes = 10
}

resource "pingone_resource_secret" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  resource_id    = pingone_resource.%[2]s.id

  previous = {
    expires_at = time_offset.%[2]s.rfc3339
  }
}

data "pingone_resource_secret" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  resource_id    = pingone_resource.%[2]s.id

  depends_on = [
    pingone_resource_secret.%[2]s,
  ]
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccResourceSecretDataSourceConfig_Rotation1(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_resource" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "pingone_resource_secret" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  resource_id    = pingone_resource.%[2]s.id
}

data "pingone_resource_secret" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  resource_id    = pingone_resource.%[2]s.id

  depends_on = [
    pingone_resource_secret.%[2]s,
  ]
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccResourceSecretDataSourceConfig_Rotation2(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_resource" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "time_offset" "%[2]s" {
  offset_minutes = 10
}

resource "pingone_resource_secret" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  resource_id    = pingone_resource.%[2]s.id

  previous = {
    expires_at = time_offset.%[2]s.rfc3339
  }
}

data "pingone_resource_secret" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  resource_id    = pingone_resource.%[2]s.id

  depends_on = [
    pingone_resource_secret.%[2]s,
  ]
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccResourceSecretDataSourceConfig_IncorrectResourceType(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

data "pingone_resource" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "openid"
}

data "pingone_resource_secret" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  resource_id    = data.pingone_resource.%[2]s.id
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccResourceSecretDataSourceConfig_NotFoundByID(resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_resource_secret" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  resource_id    = "9c052a8a-14be-44e4-8f07-2662569994ce" // dummy ID that conforms to UUID v4
}
`, acctest.GenericSandboxEnvironment(), resourceName)
}
