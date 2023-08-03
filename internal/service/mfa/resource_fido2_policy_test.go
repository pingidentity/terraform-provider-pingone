package mfa_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func testAccCheckFIDO2PolicyDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.MFAAPIClient

	apiClientManagement := p1Client.API.ManagementAPIClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_mfa_fido2_policy" {
			continue
		}

		_, rEnv, err := apiClientManagement.EnvironmentsApi.ReadOneEnvironment(ctx, rs.Primary.Attributes["environment_id"]).Execute()

		if err != nil {

			if rEnv == nil {
				return fmt.Errorf("Response object does not exist and no error detected")
			}

			if rEnv.StatusCode == 404 {
				continue
			}

			return err
		}

		body, r, err := apiClient.FIDO2PolicyApi.ReadOneFIDO2Policy(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.ID).Execute()

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

		return fmt.Errorf("PingOne MFA FIDO2 Policy Instance %s still exists", rs.Primary.ID)
	}

	return nil
}

func testAccGetMFAPolicyIDs(resourceName string, environmentID, resourceID *string) resource.TestCheckFunc {
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

func TestAccMFAPolicy_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_policy.%s", resourceName)

	name := resourceName

	var resourceID, environmentID string

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckMFAPolicyDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccMFAPolicyConfig_FullSMS(resourceName, name),
				Check:  testAccGetMFAPolicyIDs(resourceFullName, &environmentID, &resourceID),
			},
			// Replan after removal preconfig
			{
				PreConfig: func() {
					var ctx = context.Background()
					p1Client, err := acctest.TestClient(ctx)

					if err != nil {
						t.Fatalf("Failed to get API client: %v", err)
					}

					apiClient := p1Client.API.MFAAPIClient

					if environmentID == "" || resourceID == "" {
						t.Fatalf("One of environment ID or resource ID cannot be determined. Environment ID: %s, Resource ID: %s", environmentID, resourceID)
					}

					_, err = apiClient.DeviceAuthenticationPolicyApi.DeleteDeviceAuthenticationPolicy(ctx, environmentID, resourceID).Execute()
					if err != nil {
						t.Fatalf("Failed to delete MFA Policy: %v", err)
					}
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccFIDO2Policy_NewEnv(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_fido2_policy.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckFIDO2PolicyDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccFIDO2PolicyConfig_NewEnv(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
				),
			},
		},
	})
}

func TestAccFIDO2Policy_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_fido2_policy.%s", resourceName)

	name := resourceName

	fullCheck := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
		resource.TestCheckResourceAttr(resourceFullName, "name", name),
		resource.TestCheckResourceAttr(resourceFullName, "description", "Test FIDO policy"),
		resource.TestCheckResourceAttr(resourceFullName, "attestation_requirements", "DIRECT"),
		resource.TestCheckResourceAttr(resourceFullName, "authenticator_attachment", "BOTH"),
		resource.TestCheckResourceAttr(resourceFullName, "backup_eligibility.allow", "true"),
		resource.TestCheckResourceAttr(resourceFullName, "backup_eligibility.enforce_during_authentication", "false"),
		resource.TestCheckResourceAttr(resourceFullName, "device_display_name", "Test Device Max"),
		resource.TestCheckResourceAttr(resourceFullName, "discoverable_credentials", "PREFERRED"),
		resource.TestCheckResourceAttr(resourceFullName, "mds_authenticators_requirements.allowed_authenticator_ids.#", "3"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "mds_authenticators_requirements.allowed_authenticator_ids.*", "authenticator_id_1"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "mds_authenticators_requirements.allowed_authenticator_ids.*", "authenticator_id_2"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "mds_authenticators_requirements.allowed_authenticator_ids.*", "authenticator_id_3"),
		resource.TestCheckResourceAttr(resourceFullName, "mds_authenticators_requirements.enforce_during_authentication", "true"),
		resource.TestCheckResourceAttr(resourceFullName, "mds_authenticators_requirements.option", "SPECIFIC"),
		resource.TestCheckResourceAttr(resourceFullName, "relying_party_id", "pingidentity.com"),
		resource.TestCheckResourceAttr(resourceFullName, "user_display_name_attributes.attributes.#", "3"),
		resource.TestCheckResourceAttr(resourceFullName, "user_display_name_attributes.attributes.0.name", "email"),
		resource.TestCheckResourceAttr(resourceFullName, "user_display_name_attributes.attributes.1.name", "name"),
		resource.TestCheckResourceAttr(resourceFullName, "user_display_name_attributes.attributes.1.sub_attributes.#", "2"),
		resource.TestCheckResourceAttr(resourceFullName, "user_display_name_attributes.attributes.1.sub_attributes.0.name", "given"),
		resource.TestCheckResourceAttr(resourceFullName, "user_display_name_attributes.attributes.1.sub_attributes.1.name", "family"),
		resource.TestCheckResourceAttr(resourceFullName, "user_display_name_attributes.attributes.2.name", "username"),
		resource.TestCheckResourceAttr(resourceFullName, "user_verification.enforce_during_authentication", "true"),
		resource.TestCheckResourceAttr(resourceFullName, "user_verification.option", "REQUIRED"),
	)

	minimalCheck := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
		resource.TestCheckResourceAttr(resourceFullName, "name", name),
		resource.TestCheckNoResourceAttr(resourceFullName, "description"),
		resource.TestCheckResourceAttr(resourceFullName, "attestation_requirements", "NONE"),
		resource.TestCheckResourceAttr(resourceFullName, "authenticator_attachment", "PLATFORM"),
		resource.TestCheckResourceAttr(resourceFullName, "backup_eligibility.allow", "false"),
		resource.TestCheckResourceAttr(resourceFullName, "backup_eligibility.enforce_during_authentication", "true"),
		resource.TestCheckResourceAttr(resourceFullName, "device_display_name", "fidoPolicy.deviceDisplayName02"),
		resource.TestCheckResourceAttr(resourceFullName, "discoverable_credentials", "DISCOURAGED"),
		resource.TestCheckResourceAttr(resourceFullName, "mds_authenticators_requirements.allowed_authenticator_ids.#", "0"),
		resource.TestCheckResourceAttr(resourceFullName, "mds_authenticators_requirements.enforce_during_authentication", "false"),
		resource.TestCheckResourceAttr(resourceFullName, "mds_authenticators_requirements.option", "NONE"),
		resource.TestCheckResourceAttr(resourceFullName, "relying_party_id", "ping-devops.com"),
		resource.TestCheckResourceAttr(resourceFullName, "user_display_name_attributes.attributes.#", "1"),
		resource.TestCheckResourceAttr(resourceFullName, "user_display_name_attributes.attributes.0.name", "username"),
		resource.TestCheckResourceAttr(resourceFullName, "user_verification.enforce_during_authentication", "false"),
		resource.TestCheckResourceAttr(resourceFullName, "user_verification.option", "DISCOURAGED"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckFIDO2PolicyDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			{
				Config: testAccFIDO2PolicyConfig_Full(resourceName, name),
				Check:  fullCheck,
			},
			{
				Config:  testAccFIDO2PolicyConfig_Full(resourceName, name),
				Destroy: true,
			},
			// Minimal
			{
				Config: testAccFIDO2PolicyConfig_Minimal(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config:  testAccFIDO2PolicyConfig_Minimal(resourceName, name),
				Destroy: true,
			},
			// Change
			{
				Config: testAccFIDO2PolicyConfig_Full(resourceName, name),
				Check:  fullCheck,
			},
			{
				Config: testAccFIDO2PolicyConfig_Minimal(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config: testAccFIDO2PolicyConfig_Full(resourceName, name),
				Check:  fullCheck,
			},
		},
	})
}

func TestAccFIDO2Policy_Errors(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckFIDO2PolicyDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccFIDO2PolicyConfig_ConflictedOptions_1(resourceName, name),
				ExpectError: regexp.MustCompile(`Invalid argument combination`),
			},
			{
				Config:      testAccFIDO2PolicyConfig_ConflictedOptions_2(resourceName, name),
				ExpectError: regexp.MustCompile(`Invalid argument combination`),
			},
			{
				Config:      testAccFIDO2PolicyConfig_ConflictedOptions_3(resourceName, name),
				ExpectError: regexp.MustCompile(`Invalid Attribute Value`),
			},
		},
	})
}

func testAccFIDO2PolicyConfig_NewEnv(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_mfa_fido2_policy" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  name           = "%[3]s"

  attestation_requirements = "NONE"
  authenticator_attachment = "PLATFORM"

  backup_eligibility = {
    allow                         = false
    enforce_during_authentication = true
  }

  device_display_name = "fidoPolicy.deviceDisplayName02"

  discoverable_credentials = "DISCOURAGED"

  mds_authenticators_requirements = {
    enforce_during_authentication = false
    option                        = "NONE"
  }

  relying_party_id = "ping-devops.com"

  user_display_name_attributes = {
    attributes = [
      {
        name = "username"
      }
    ]
  }

  user_verification = {
    enforce_during_authentication = false
    option                        = "DISCOURAGED"
  }
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccFIDO2PolicyConfig_Full(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_mfa_fido2_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "Test FIDO policy"

  attestation_requirements = "DIRECT"
  authenticator_attachment = "BOTH"

  backup_eligibility = {
    allow                         = true
    enforce_during_authentication = false
  }

  device_display_name = "Test Device Max"

  discoverable_credentials = "PREFERRED"

  mds_authenticators_requirements = {
    allowed_authenticator_ids = [
      "authenticator_id_1",
      "authenticator_id_3",
      "authenticator_id_2",
    ]

    enforce_during_authentication = true
    option                        = "SPECIFIC"
  }

  relying_party_id = "pingidentity.com"

  user_display_name_attributes = {
    attributes = [
      {
        name = "email"
      },
      {
        name = "name",
        sub_attributes = [
          {
            name = "given"
          },
          {
            name = "family"
          }
        ]
      },
      {
        name = "username"
      }
    ]
  }

  user_verification = {
    enforce_during_authentication = true
    option                        = "REQUIRED"
  }

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccFIDO2PolicyConfig_Minimal(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_mfa_fido2_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  attestation_requirements = "NONE"
  authenticator_attachment = "PLATFORM"

  backup_eligibility = {
    allow                         = false
    enforce_during_authentication = true
  }

  device_display_name = "fidoPolicy.deviceDisplayName02"

  discoverable_credentials = "DISCOURAGED"

  mds_authenticators_requirements = {
    enforce_during_authentication = false
    option                        = "NONE"
  }

  relying_party_id = "ping-devops.com"

  user_display_name_attributes = {
    attributes = [
      {
        name = "username"
      }
    ]
  }

  user_verification = {
    enforce_during_authentication = false
    option                        = "DISCOURAGED"
  }

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccFIDO2PolicyConfig_ConflictedOptions_1(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_mfa_fido2_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  attestation_requirements = "NONE"
  authenticator_attachment = "PLATFORM"

  backup_eligibility = {
    allow                         = false
    enforce_during_authentication = true
  }

  device_display_name = "fidoPolicy.deviceDisplayName02"

  discoverable_credentials = "DISCOURAGED"

  mds_authenticators_requirements = {
    enforce_during_authentication = false
    option                        = "CERTIFIED"
  }

  relying_party_id = "ping-devops.com"

  user_display_name_attributes = {
    attributes = [
      {
        name = "username"
      }
    ]
  }

  user_verification = {
    enforce_during_authentication = false
    option                        = "DISCOURAGED"
  }

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccFIDO2PolicyConfig_ConflictedOptions_2(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_mfa_fido2_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  attestation_requirements = "NONE"
  authenticator_attachment = "PLATFORM"

  backup_eligibility = {
    allow                         = false
    enforce_during_authentication = true
  }

  device_display_name = "fidoPolicy.deviceDisplayName02"

  discoverable_credentials = "DISCOURAGED"

  mds_authenticators_requirements = {
    allowed_authenticator_ids = [
      "authenticator_id_1",
      "authenticator_id_3",
      "authenticator_id_2",
    ]

    enforce_during_authentication = false
    option                        = "NONE"
  }

  relying_party_id = "ping-devops.com"

  user_display_name_attributes = {
    attributes = [
      {
        name = "username"
      }
    ]
  }

  user_verification = {
    enforce_during_authentication = false
    option                        = "DISCOURAGED"
  }

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccFIDO2PolicyConfig_ConflictedOptions_3(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_mfa_fido2_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  attestation_requirements = "NONE"
  authenticator_attachment = "PLATFORM"

  backup_eligibility = {
    allow                         = false
    enforce_during_authentication = true
  }

  device_display_name = "fidoPolicy.deviceDisplayName02"

  discoverable_credentials = "DISCOURAGED"

  mds_authenticators_requirements = {
    enforce_during_authentication = false
    option                        = "NONE"
  }

  relying_party_id = "ping-devops.com"

  user_display_name_attributes = {
    attributes = [
      {
        name = "email"
      }
    ]
  }

  user_verification = {
    enforce_during_authentication = false
    option                        = "DISCOURAGED"
  }

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}
