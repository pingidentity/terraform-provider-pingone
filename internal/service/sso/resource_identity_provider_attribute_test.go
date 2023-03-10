package sso_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func testAccCheckIdentityProviderAttributeDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_identity_provider_attribute" {
			continue
		}

		_, rEnv, err := apiClient.EnvironmentsApi.ReadOneEnvironment(ctx, rs.Primary.Attributes["environment_id"]).Execute()

		if err != nil {

			if rEnv == nil {
				return fmt.Errorf("Response object does not exist and no error detected")
			}

			if rEnv.StatusCode == 404 {
				continue
			}

			return err
		}

		body, r, err := apiClient.IdentityProviderAttributesApi.ReadOneIdentityProviderAttribute(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["identity_provider_id"], rs.Primary.ID).Execute()

		if err != nil {

			if r == nil {
				return fmt.Errorf("Response object does not exist and no error detected")
			}

			if r.StatusCode == 404 {
				continue
			}

			tflog.Error(ctx, fmt.Sprintf("Error: %v", body))
			return err
		}

		return fmt.Errorf("PingOne Identity Provider attribute %s still exists", rs.Primary.ID)
	}

	return nil
}

func TestAccIdentityProviderAttribute_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_identity_provider_attribute.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckIdentityProviderAttributeDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityProviderAttributeConfig_Full(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "identity_provider_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(resourceFullName, "name", "email"),
					resource.TestCheckResourceAttr(resourceFullName, "update", "EMPTY_ONLY"),
					resource.TestCheckResourceAttr(resourceFullName, "value", "${providerAttributes.emailAddress.value}"),
					resource.TestCheckResourceAttr(resourceFullName, "mapping_type", "CUSTOM"),
				),
			},
		},
	})
}

func TestAccIdentityProviderAttribute_Minimal(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_identity_provider_attribute.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckIdentityProviderAttributeDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityProviderAttributeConfig_Minimal(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "identity_provider_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(resourceFullName, "name", "email"),
					resource.TestCheckResourceAttr(resourceFullName, "update", "ALWAYS"),
					resource.TestCheckResourceAttr(resourceFullName, "value", "${providerAttributes.emailAddress.value}"),
					resource.TestCheckResourceAttr(resourceFullName, "mapping_type", "CUSTOM"),
				),
			},
		},
	})
}

func TestAccIdentityProviderAttribute_Expression(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_identity_provider_attribute.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckIdentityProviderAttributeDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityProviderAttributeConfig_Expression(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "identity_provider_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(resourceFullName, "name", "name.given"),
					resource.TestCheckResourceAttr(resourceFullName, "update", "ALWAYS"),
					resource.TestCheckResourceAttr(resourceFullName, "value", "${providerAttributes.name.givenName}"),
					resource.TestCheckResourceAttr(resourceFullName, "mapping_type", "CUSTOM"),
				),
			},
		},
	})
}

func TestAccIdentityProviderAttribute_Change(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_identity_provider_attribute.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckIdentityProviderAttributeDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityProviderAttributeConfig_Minimal(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "name", "email"),
					resource.TestCheckResourceAttr(resourceFullName, "update", "ALWAYS"),
					resource.TestCheckResourceAttr(resourceFullName, "value", "${providerAttributes.emailAddress.value}"),
				),
			},
			{
				Config: testAccIdentityProviderAttributeConfig_Expression(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "name", "name.given"),
					resource.TestCheckResourceAttr(resourceFullName, "update", "ALWAYS"),
					resource.TestCheckResourceAttr(resourceFullName, "value", "${providerAttributes.name.givenName}"),
				),
			},
			{
				Config: testAccIdentityProviderAttributeConfig_Full(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "name", "email"),
					resource.TestCheckResourceAttr(resourceFullName, "update", "EMPTY_ONLY"),
					resource.TestCheckResourceAttr(resourceFullName, "value", "${providerAttributes.emailAddress.value}"),
				),
			},
			{
				Config:      testAccIdentityProviderAttributeConfig_ReservedAttributeName(resourceName, name),
				ExpectError: regexp.MustCompile(`expected name to not be any of \[[a-zA-Z ]*\], got [a-zA-Z]*`),
			},
			{
				Config: testAccIdentityProviderAttributeConfig_Minimal(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "name", "email"),
					resource.TestCheckResourceAttr(resourceFullName, "update", "ALWAYS"),
					resource.TestCheckResourceAttr(resourceFullName, "value", "${providerAttributes.emailAddress.value}"),
				),
			},
		},
	})
}

func TestAccIdentityProviderAttribute_ReservedAttributeName(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckIdentityProviderAttributeDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccIdentityProviderAttributeConfig_ReservedAttributeName(resourceName, name),
				ExpectError: regexp.MustCompile(`expected name to not be any of \[[a-zA-Z ]*\], got [a-zA-Z]*`),
			},
		},
	})
}

func testAccIdentityProviderAttributeConfig_Full(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_identity_provider" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  google {
    client_id     = "testclientid"
    client_secret = "testclientsecret"
  }
}

resource "pingone_identity_provider_attribute" "%[2]s" {
  environment_id       = data.pingone_environment.general_test.id
  identity_provider_id = pingone_identity_provider.%[2]s.id

  name   = "email"
  update = "EMPTY_ONLY"
  value  = "$${providerAttributes.emailAddress.value}"
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccIdentityProviderAttributeConfig_Minimal(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_identity_provider" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  google {
    client_id     = "testclientid"
    client_secret = "testclientsecret"
  }
}

resource "pingone_identity_provider_attribute" "%[2]s" {
  environment_id       = data.pingone_environment.general_test.id
  identity_provider_id = pingone_identity_provider.%[2]s.id

  name   = "email"
  update = "ALWAYS"
  value  = "$${providerAttributes.emailAddress.value}"
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccIdentityProviderAttributeConfig_Expression(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_identity_provider" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  google {
    client_id     = "testclientid"
    client_secret = "testclientsecret"
  }
}

resource "pingone_identity_provider_attribute" "%[2]s" {
  environment_id       = data.pingone_environment.general_test.id
  identity_provider_id = pingone_identity_provider.%[2]s.id

  name   = "name.given"
  update = "ALWAYS"
  value  = "$${providerAttributes.name.givenName}"
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccIdentityProviderAttributeConfig_ReservedAttributeName(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_identity_provider" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  google {
    client_id     = "testclientid"
    client_secret = "testclientsecret"
  }
}

resource "pingone_identity_provider_attribute" "%[2]s" {
  environment_id       = data.pingone_environment.general_test.id
  identity_provider_id = pingone_identity_provider.%[2]s.id

  name   = "account"
  update = "ALWAYS"
  value  = "$${'test'}"
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}
