package sso_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func testAccCheckUserDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.ManagementAPIClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_user" {
			continue
		}

		body, r, err := apiClient.UsersApi.ReadUser(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.ID).Execute()

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

		return fmt.Errorf("PingOne User %s still exists", rs.Primary.ID)
	}

	return nil
}

func testAccGetUserIDs(resourceName string, environmentID, resourceID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource Not found: %s", resourceName)
		}

		*resourceID = rs.Primary.ID
		*environmentID = rs.Primary.Attributes["environment_id"]

		return nil
	}
}

func TestAccUser_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_user.%s", resourceName)

	name := resourceName

	var resourceID, environmentID string

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckUserDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccUserConfig_Minimal(resourceName, name),
				Check:  testAccGetUserIDs(resourceFullName, &environmentID, &resourceID),
			},
			// Replan after removal preconfig
			{
				PreConfig: func() {
					var ctx = context.Background()
					p1Client, err := acctest.TestClient(ctx)

					if err != nil {
						t.Fatalf("Failed to get API client: %v", err)
					}

					apiClient := p1Client.API.ManagementAPIClient

					if environmentID == "" || resourceID == "" {
						t.Fatalf("One of environment ID or resource ID cannot be determined. Environment ID: %s, Resource ID: %s", environmentID, resourceID)
					}

					_, err = apiClient.UsersApi.DeleteUser(ctx, environmentID, resourceID).Execute()
					if err != nil {
						t.Fatalf("Failed to delete user: %v", err)
					}
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccUser_NewEnv(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_user.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckUserDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccUserConfig_NewEnv(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "username", name),
				),
			},
		},
	})
}

func TestAccUser_All(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_user.%s", resourceName)

	name := resourceName

	fullTest := resource.TestStep{
		Config: testAccUserConfig_Full(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
			resource.TestCheckResourceAttr(resourceFullName, "username", name),
			resource.TestCheckResourceAttr(resourceFullName, "email", fmt.Sprintf("%s@pingidentity.com", name)),
			//resource.TestCheckNoResourceAttr(resourceFullName, "email_verified"),
			resource.TestMatchResourceAttr(resourceFullName, "population_id", verify.P1ResourceIDRegexp),
			resource.TestCheckResourceAttr(resourceFullName, "enabled", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "status", "DISABLED"),
			resource.TestCheckResourceAttr(resourceFullName, "account.can_authenticate", "false"),
			resource.TestCheckNoResourceAttr(resourceFullName, "account.locked_at"),
			resource.TestCheckResourceAttr(resourceFullName, "account.status", "OK"),
			resource.TestCheckResourceAttr(resourceFullName, "address.country_code", "US"),
			resource.TestCheckResourceAttr(resourceFullName, "address.locality", "Springfield"),
			resource.TestCheckResourceAttr(resourceFullName, "address.postal_code", "BAR7"),
			resource.TestCheckResourceAttr(resourceFullName, "address.region", "Who knows"),
			resource.TestCheckResourceAttr(resourceFullName, "address.street_address", "742 Evergreen Terrace"),
			resource.TestCheckResourceAttr(resourceFullName, "external_id", "12345678-id"),
			resource.TestMatchResourceAttr(resourceFullName, "identity_provider.id", verify.P1ResourceIDRegexp),
			resource.TestCheckResourceAttr(resourceFullName, "identity_provider.type", "LINKEDIN"),
			resource.TestCheckResourceAttr(resourceFullName, "user_lifecycle.status", "VERIFICATION_REQUIRED"),
			resource.TestCheckResourceAttr(resourceFullName, "user_lifecycle.suppress_verification_code", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "locale", "es-419"),
			resource.TestCheckResourceAttr(resourceFullName, "mfa_enabled", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "mobile_phone", "+1 555-4796"),
			resource.TestCheckResourceAttr(resourceFullName, "name.family", "Simpson"),
			resource.TestCheckResourceAttr(resourceFullName, "name.formatted", "Mr. Homer Jay Simpson Jr."),
			resource.TestCheckResourceAttr(resourceFullName, "name.given", "Homer"),
			resource.TestCheckResourceAttr(resourceFullName, "name.middle", "Jay"),
			resource.TestCheckResourceAttr(resourceFullName, "name.honorific_prefix", "Mr."),
			resource.TestCheckResourceAttr(resourceFullName, "name.honorific_suffix", "Jr."),
			resource.TestCheckResourceAttr(resourceFullName, "nickname", "Homie"),
			resource.TestCheckResourceAttr(resourceFullName, "password.force_change", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "password.initial_value", "SuperSecretDummyPassword"),
			resource.TestCheckNoResourceAttr(resourceFullName, "password.external"),
			resource.TestCheckResourceAttr(resourceFullName, "photo.href", "https://www.pingidentity.com/homer-simpson.png"),
			resource.TestCheckResourceAttr(resourceFullName, "preferred_language", "en;q=0.7"),
			resource.TestCheckResourceAttr(resourceFullName, "primary_phone", "555-6832"),
			resource.TestCheckResourceAttr(resourceFullName, "timezone", "America/Los_Angeles"),
			resource.TestCheckResourceAttr(resourceFullName, "title", "President of Mr Plow enterprises, snow plowing entrepreneur"),
			resource.TestCheckResourceAttr(resourceFullName, "type", "Employee"),
			resource.TestCheckResourceAttr(resourceFullName, "verify_status", "ENABLED"),
		),
	}

	minimalTest := resource.TestStep{
		Config: testAccUserConfig_Minimal(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
			resource.TestCheckResourceAttr(resourceFullName, "username", name),
			resource.TestCheckResourceAttr(resourceFullName, "email", fmt.Sprintf("%s@pingidentity.com", name)),
			//resource.TestCheckNoResourceAttr(resourceFullName, "email_verified"),
			resource.TestMatchResourceAttr(resourceFullName, "population_id", verify.P1ResourceIDRegexp),
			resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "status", "ENABLED"),
			resource.TestCheckResourceAttr(resourceFullName, "account.can_authenticate", "true"),
			resource.TestCheckNoResourceAttr(resourceFullName, "account.locked_at"),
			resource.TestCheckResourceAttr(resourceFullName, "account.status", "OK"),
			resource.TestCheckNoResourceAttr(resourceFullName, "address.country_code"),
			resource.TestCheckNoResourceAttr(resourceFullName, "address.locality"),
			resource.TestCheckNoResourceAttr(resourceFullName, "address.postal_code"),
			resource.TestCheckNoResourceAttr(resourceFullName, "address.region"),
			resource.TestCheckNoResourceAttr(resourceFullName, "address.street_address"),
			resource.TestCheckNoResourceAttr(resourceFullName, "external_id"),
			resource.TestCheckNoResourceAttr(resourceFullName, "identity_provider.id"),
			resource.TestCheckResourceAttr(resourceFullName, "identity_provider.type", "PING_ONE"),
			resource.TestCheckResourceAttr(resourceFullName, "user_lifecycle.status", "ACCOUNT_OK"),
			resource.TestCheckNoResourceAttr(resourceFullName, "user_lifecycle.suppress_verification_code"),
			resource.TestCheckNoResourceAttr(resourceFullName, "locale"),
			resource.TestCheckResourceAttr(resourceFullName, "mfa_enabled", "false"),
			resource.TestCheckNoResourceAttr(resourceFullName, "mobile_phone"),
			resource.TestCheckNoResourceAttr(resourceFullName, "name.family"),
			resource.TestCheckNoResourceAttr(resourceFullName, "name.formatted"),
			resource.TestCheckNoResourceAttr(resourceFullName, "name.given"),
			resource.TestCheckNoResourceAttr(resourceFullName, "name.middle"),
			resource.TestCheckNoResourceAttr(resourceFullName, "name.honorific_prefix"),
			resource.TestCheckNoResourceAttr(resourceFullName, "name.honorific_suffix"),
			resource.TestCheckNoResourceAttr(resourceFullName, "nickname"),
			resource.TestCheckNoResourceAttr(resourceFullName, "password"),
			resource.TestCheckNoResourceAttr(resourceFullName, "photo.href"),
			resource.TestCheckNoResourceAttr(resourceFullName, "preferred_language"),
			resource.TestCheckNoResourceAttr(resourceFullName, "primary_phone"),
			resource.TestCheckNoResourceAttr(resourceFullName, "timezone"),
			resource.TestCheckNoResourceAttr(resourceFullName, "title"),
			resource.TestCheckNoResourceAttr(resourceFullName, "type"),
			resource.TestCheckResourceAttr(resourceFullName, "verify_status", "NOT_INITIATED"),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckUserDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full test
			fullTest,
			{
				Config:  testAccUserConfig_Full(resourceName, name),
				Destroy: true,
			},
			// Minimal test
			minimalTest,
			{
				Config:  testAccUserConfig_Minimal(resourceName, name),
				Destroy: true,
			},
			// Change
			fullTest,
			minimalTest,
			fullTest,
		},
	})
}

func TestAccUser_AllWithoutReplacement(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_user.%s", resourceName)

	name := resourceName

	fullTest := resource.TestStep{
		Config: testAccUserConfig_FullWithoutReplaceParams(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
			resource.TestCheckResourceAttr(resourceFullName, "username", name),
			resource.TestCheckResourceAttr(resourceFullName, "email", fmt.Sprintf("%s@pingidentity.com", name)),
			//resource.TestCheckNoResourceAttr(resourceFullName, "email_verified"),
			resource.TestMatchResourceAttr(resourceFullName, "population_id", verify.P1ResourceIDRegexp),
			resource.TestCheckResourceAttr(resourceFullName, "enabled", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "status", "DISABLED"),
			resource.TestCheckResourceAttr(resourceFullName, "account.can_authenticate", "false"),
			resource.TestCheckNoResourceAttr(resourceFullName, "account.locked_at"),
			resource.TestCheckResourceAttr(resourceFullName, "account.status", "OK"),
			resource.TestCheckResourceAttr(resourceFullName, "address.country_code", "US"),
			resource.TestCheckResourceAttr(resourceFullName, "address.locality", "Springfield"),
			resource.TestCheckResourceAttr(resourceFullName, "address.postal_code", "BAR7"),
			resource.TestCheckResourceAttr(resourceFullName, "address.region", "Who knows"),
			resource.TestCheckResourceAttr(resourceFullName, "address.street_address", "742 Evergreen Terrace"),
			resource.TestCheckResourceAttr(resourceFullName, "external_id", "12345678-id"),
			resource.TestMatchResourceAttr(resourceFullName, "identity_provider.id", verify.P1ResourceIDRegexp),
			resource.TestCheckResourceAttr(resourceFullName, "identity_provider.type", "LINKEDIN"),
			resource.TestCheckResourceAttr(resourceFullName, "locale", "es-419"),
			resource.TestCheckResourceAttr(resourceFullName, "mfa_enabled", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "mobile_phone", "+1 555-4796"),
			resource.TestCheckResourceAttr(resourceFullName, "name.family", "Simpson"),
			resource.TestCheckResourceAttr(resourceFullName, "name.formatted", "Mr. Homer Jay Simpson Jr."),
			resource.TestCheckResourceAttr(resourceFullName, "name.given", "Homer"),
			resource.TestCheckResourceAttr(resourceFullName, "name.middle", "Jay"),
			resource.TestCheckResourceAttr(resourceFullName, "name.honorific_prefix", "Mr."),
			resource.TestCheckResourceAttr(resourceFullName, "name.honorific_suffix", "Jr."),
			resource.TestCheckResourceAttr(resourceFullName, "nickname", "Homie"),
			resource.TestCheckResourceAttr(resourceFullName, "password.force_change", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "password.initial_value", "SuperSecretDummyPassword"),
			resource.TestCheckNoResourceAttr(resourceFullName, "password.external"),
			resource.TestCheckResourceAttr(resourceFullName, "photo.href", "https://www.pingidentity.com/homer-simpson.png"),
			resource.TestCheckResourceAttr(resourceFullName, "preferred_language", "en;q=0.7"),
			resource.TestCheckResourceAttr(resourceFullName, "primary_phone", "555-6832"),
			resource.TestCheckResourceAttr(resourceFullName, "timezone", "America/Los_Angeles"),
			resource.TestCheckResourceAttr(resourceFullName, "title", "President of Mr Plow enterprises, snow plowing entrepreneur"),
			resource.TestCheckResourceAttr(resourceFullName, "type", "Employee"),
		),
	}

	minimalTest := resource.TestStep{
		Config: testAccUserConfig_Minimal(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
			resource.TestCheckResourceAttr(resourceFullName, "username", name),
			resource.TestCheckResourceAttr(resourceFullName, "email", fmt.Sprintf("%s@pingidentity.com", name)),
			//resource.TestCheckNoResourceAttr(resourceFullName, "email_verified"),
			resource.TestMatchResourceAttr(resourceFullName, "population_id", verify.P1ResourceIDRegexp),
			resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "status", "ENABLED"),
			resource.TestCheckResourceAttr(resourceFullName, "account.can_authenticate", "true"),
			resource.TestCheckNoResourceAttr(resourceFullName, "account.locked_at"),
			resource.TestCheckResourceAttr(resourceFullName, "account.status", "OK"),
			resource.TestCheckNoResourceAttr(resourceFullName, "address.country_code"),
			resource.TestCheckNoResourceAttr(resourceFullName, "address.locality"),
			resource.TestCheckNoResourceAttr(resourceFullName, "address.postal_code"),
			resource.TestCheckNoResourceAttr(resourceFullName, "address.region"),
			resource.TestCheckNoResourceAttr(resourceFullName, "address.street_address"),
			resource.TestCheckNoResourceAttr(resourceFullName, "external_id"),
			resource.TestCheckNoResourceAttr(resourceFullName, "identity_provider.id"),
			resource.TestCheckResourceAttr(resourceFullName, "identity_provider.type", "PING_ONE"),
			resource.TestCheckNoResourceAttr(resourceFullName, "locale"),
			resource.TestCheckResourceAttr(resourceFullName, "mfa_enabled", "false"),
			resource.TestCheckNoResourceAttr(resourceFullName, "mobile_phone"),
			resource.TestCheckNoResourceAttr(resourceFullName, "name.family"),
			resource.TestCheckNoResourceAttr(resourceFullName, "name.formatted"),
			resource.TestCheckNoResourceAttr(resourceFullName, "name.given"),
			resource.TestCheckNoResourceAttr(resourceFullName, "name.middle"),
			resource.TestCheckNoResourceAttr(resourceFullName, "name.honorific_prefix"),
			resource.TestCheckNoResourceAttr(resourceFullName, "name.honorific_suffix"),
			resource.TestCheckNoResourceAttr(resourceFullName, "nickname"),
			resource.TestCheckNoResourceAttr(resourceFullName, "password"),
			resource.TestCheckNoResourceAttr(resourceFullName, "photo.href"),
			resource.TestCheckNoResourceAttr(resourceFullName, "preferred_language"),
			resource.TestCheckNoResourceAttr(resourceFullName, "primary_phone"),
			resource.TestCheckNoResourceAttr(resourceFullName, "timezone"),
			resource.TestCheckNoResourceAttr(resourceFullName, "title"),
			resource.TestCheckNoResourceAttr(resourceFullName, "type"),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckUserDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full test
			fullTest,
			{
				Config:  testAccUserConfig_FullWithoutReplaceParams(resourceName, name),
				Destroy: true,
			},
			// Minimal test
			minimalTest,
			{
				Config:  testAccUserConfig_Minimal(resourceName, name),
				Destroy: true,
			},
			// Change
			fullTest,
			minimalTest,
			fullTest,
		},
	})
}

func TestAccUser_ChangePopulation(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_user.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckUserDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccUserConfig_Minimal(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(resourceFullName, "username", name),
					resource.TestCheckResourceAttr(resourceFullName, "email", fmt.Sprintf("%s@pingidentity.com", name)),
					resource.TestMatchResourceAttr(resourceFullName, "population_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(resourceFullName, "status", "ENABLED"),
				),
			},
			{
				Config: testAccUserConfig_CustomPopulation(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(resourceFullName, "username", name),
					resource.TestCheckResourceAttr(resourceFullName, "email", fmt.Sprintf("%s@pingidentity.com", name)),
					resource.TestMatchResourceAttr(resourceFullName, "population_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(resourceFullName, "status", "ENABLED"),
				),
			},
		},
	})
}

func TestAccUser_ChangeMFA(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_user.%s", resourceName)

	name := resourceName

	enabledTest := resource.TestStep{
		Config: testAccUserConfig_MinimalWithMFA(resourceName, name, true),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "mfa_enabled", "true"),
		),
	}

	disabledTest := resource.TestStep{
		Config: testAccUserConfig_MinimalWithMFA(resourceName, name, false),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "mfa_enabled", "false"),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckUserDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Change
			enabledTest,
			disabledTest,
			enabledTest,
		},
	})
}

func TestAccUser_ChangeUsernameAndEmail(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_user.%s", resourceName)

	name1 := fmt.Sprintf("%s1", resourceName)
	name2 := fmt.Sprintf("%s2", resourceName)

	fullTest := resource.TestStep{
		Config: testAccUserConfig_Minimal(resourceName, name1),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "username", name1),
			resource.TestCheckResourceAttr(resourceFullName, "email", fmt.Sprintf("%s@pingidentity.com", name1)),
		),
	}

	minimalTest := resource.TestStep{
		Config: testAccUserConfig_Minimal(resourceName, name2),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "username", name2),
			resource.TestCheckResourceAttr(resourceFullName, "email", fmt.Sprintf("%s@pingidentity.com", name2)),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckUserDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Change
			fullTest,
			minimalTest,
			fullTest,
		},
	})
}

func testAccUserConfig_NewEnv(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_population" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  name = "%[4]s"
}

resource "pingone_user" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  username      = "%[4]s"
  email         = "%[4]s@pingidentity.com"
  population_id = pingone_population.%[3]s.id
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccUserConfig_Full(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_identity_provider" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  linkedin {
    client_id     = "dummyclientid1"
    client_secret = "dummyclientsecret1"
  }
}

resource "pingone_population" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "pingone_user" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  username      = "%[3]s"
  email         = "%[3]s@pingidentity.com"
  population_id = pingone_population.%[2]s.id
  enabled       = false

  account = {
    can_authenticate = false
    status           = "OK"
  }

  address = {
    country_code   = "US"
    locality       = "Springfield"
    postal_code    = "BAR7"
    region         = "Who knows"
    street_address = "742 Evergreen Terrace"
  }

  external_id = "12345678-id"

  identity_provider = {
    id = pingone_identity_provider.%[2]s.id
  }

  user_lifecycle = {
    status                     = "VERIFICATION_REQUIRED"
    suppress_verification_code = true
  }

  locale      = "es-419"
  mfa_enabled = true

  mobile_phone  = "+1 555-4796"
  primary_phone = "555-6832"

  name = {
    family           = "Simpson"
    formatted        = "Mr. Homer Jay Simpson Jr."
    given            = "Homer"
    middle           = "Jay"
    honorific_prefix = "Mr."
    honorific_suffix = "Jr."
  }

  nickname = "Homie"

  password = {
    force_change  = true
    initial_value = "SuperSecretDummyPassword"
  }

  photo = {
    href = "https://www.pingidentity.com/homer-simpson.png"
  }

  preferred_language = "en;q=0.7"
  timezone           = "America/Los_Angeles"
  title              = "President of Mr Plow enterprises, snow plowing entrepreneur"
  type               = "Employee"

  verify_status = "ENABLED"
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccUserConfig_FullWithoutReplaceParams(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_identity_provider" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  linkedin {
    client_id     = "dummyclientid1"
    client_secret = "dummyclientsecret1"
  }
}

resource "pingone_population" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "pingone_user" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  username      = "%[3]s"
  email         = "%[3]s@pingidentity.com"
  population_id = pingone_population.%[2]s.id
  enabled       = false

  account = {
    can_authenticate = false
    status           = "OK"
  }

  address = {
    country_code   = "US"
    locality       = "Springfield"
    postal_code    = "BAR7"
    region         = "Who knows"
    street_address = "742 Evergreen Terrace"
  }

  external_id = "12345678-id"

  identity_provider = {
    id = pingone_identity_provider.%[2]s.id
  }

  locale      = "es-419"
  mfa_enabled = false

  mobile_phone  = "+1 555-4796"
  primary_phone = "555-6832"

  name = {
    family           = "Simpson"
    formatted        = "Mr. Homer Jay Simpson Jr."
    given            = "Homer"
    middle           = "Jay"
    honorific_prefix = "Mr."
    honorific_suffix = "Jr."
  }

  nickname = "Homie"

  password = {
    force_change  = true
    initial_value = "SuperSecretDummyPassword"
  }

  photo = {
    href = "https://www.pingidentity.com/homer-simpson.png"
  }

  preferred_language = "en;q=0.7"
  timezone           = "America/Los_Angeles"
  title              = "President of Mr Plow enterprises, snow plowing entrepreneur"
  type               = "Employee"
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccUserConfig_Minimal(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_identity_provider" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  linkedin {
    client_id     = "dummyclientid1"
    client_secret = "dummyclientsecret1"
  }
}

resource "pingone_population" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "pingone_user" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  username      = "%[3]s"
  email         = "%[3]s@pingidentity.com"
  population_id = pingone_population.%[2]s.id
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccUserConfig_MinimalWithMFA(resourceName, name string, mfaEnabled bool) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_identity_provider" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  linkedin {
    client_id     = "dummyclientid1"
    client_secret = "dummyclientsecret1"
  }
}

resource "pingone_population" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "pingone_user" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  username      = "%[3]s"
  email         = "%[3]s@pingidentity.com"
  population_id = pingone_population.%[2]s.id

  mfa_enabled = %[4]t
}`, acctest.GenericSandboxEnvironment(), resourceName, name, mfaEnabled)
}

func testAccUserConfig_CustomPopulation(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_population" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "pingone_user" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  username      = "%[3]s"
  email         = "%[3]s@pingidentity.com"
  population_id = pingone_population.%[2]s.id
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}
