package base_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
)

func TestAccLanguageDataSource_ByNameFull(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_language.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.%s", resourceFullName)

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckLanguageDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLanguageDataSourceConfig_ByNameFull(environmentName, licenseID, resourceName, "fr-FR"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceFullName, "id"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "environment_id"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "name", resourceFullName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "locale", resourceFullName, "locale"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "enabled", resourceFullName, "enabled"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "default", resourceFullName, "default"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "customer_added", resourceFullName, "customer_added"),
				),
			},
		},
	})
}

func TestAccLanguageDataSource_ByIDFull(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_language.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.%s", resourceFullName)

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckLanguageDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLanguageDataSourceConfig_ByIDFull(environmentName, licenseID, resourceName, "fr-FR"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceFullName, "id"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "environment_id"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "name", resourceFullName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "locale", resourceFullName, "locale"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "enabled", resourceFullName, "enabled"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "default", resourceFullName, "default"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "customer_added", resourceFullName, "customer_added"),
				),
			},
		},
	})
}

func TestAccLanguageDataSource_SystemDefined(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_language.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.%s", resourceFullName)

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckLanguageDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLanguageDataSourceConfig_SystemDefined(environmentName, licenseID, resourceName, "fr"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceFullName, "id"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "environment_id"),
					resource.TestCheckResourceAttr(dataSourceFullName, "name", "French"),
					resource.TestCheckResourceAttr(dataSourceFullName, "locale", "fr"),
					resource.TestCheckResourceAttr(dataSourceFullName, "enabled", "false"),
					resource.TestCheckResourceAttr(dataSourceFullName, "default", "false"),
					resource.TestCheckResourceAttr(dataSourceFullName, "customer_added", "false"),
				),
			},
		},
	})
}

func TestAccLanguageDataSource_NotFound(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckLanguageDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccLanguageDataSourceConfig_NotFoundByName(resourceName),
				ExpectError: regexp.MustCompile(`expected locale to be one of \[.*\], got doesnotexist`),
			},
			{
				Config:      testAccLanguageDataSourceConfig_NotFoundByID(resourceName),
				ExpectError: regexp.MustCompile("Error when calling `ReadOneLanguage`: Unable to find language for environmentId=[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12} with id=9c052a8a-14be-44e4-8f07-2662569994ce"),
			},
		},
	})
}

func testAccLanguageDataSourceConfig_ByNameFull(environmentName, licenseID, resourceName, locale string) string {
	return fmt.Sprintf(`


	%[1]s

resource "pingone_language" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  locale = "%[4]s"
}

data "pingone_language" "%[3]s" {

  environment_id = pingone_environment.%[2]s.id

  locale = "%[4]s"

  depends_on = [
    pingone_language.%[3]s
  ]
}
`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, locale)
}

func testAccLanguageDataSourceConfig_ByIDFull(environmentName, licenseID, resourceName, locale string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_language" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  locale = "%[4]s"
}

data "pingone_language" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  language_id = pingone_language.%[3]s.id
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, locale)
}

func testAccLanguageDataSourceConfig_SystemDefined(environmentName, licenseID, resourceName, locale string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_language" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  locale = "%[4]s"
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, locale)
}

func testAccLanguageDataSourceConfig_NotFoundByName(resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_language" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  locale = "doesnotexist"
}
`, acctest.GenericSandboxEnvironment(), resourceName)
}

func testAccLanguageDataSourceConfig_NotFoundByID(resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_language" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  language_id = "9c052a8a-14be-44e4-8f07-2662569994ce" // dummy ID that conforms to UUID v4
}
`, acctest.GenericSandboxEnvironment(), resourceName)
}
