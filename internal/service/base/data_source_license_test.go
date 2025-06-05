// Copyright © 2025 Ping Identity Corporation

package base_test

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/base"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccLicenseDataSource_ByIDFull(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	dataSourceFullName := fmt.Sprintf("data.pingone_license.%s", resourceName)

	organizationID := os.Getenv("PINGONE_ORGANIZATION_ID")
	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckOrganisationID(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.License_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLicenseDataSourceConfig_ByIDFull(resourceName, organizationID, licenseID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "license_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullName, "name", "INTERNAL"),
					resource.TestCheckResourceAttr(dataSourceFullName, "package", "INTERNAL"),
					resource.TestCheckResourceAttr(dataSourceFullName, "status", "ACTIVE"),
					resource.TestMatchResourceAttr(dataSourceFullName, "replaces_license_id", regexp.MustCompile(`^([0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})$|^()$`)),
					resource.TestMatchResourceAttr(dataSourceFullName, "replaced_by_license_id", regexp.MustCompile(`^([0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})$|^()$`)),
					resource.TestMatchResourceAttr(dataSourceFullName, "begins_at", verify.RFC3339Regexp),
					resource.TestMatchResourceAttr(dataSourceFullName, "expires_at", verify.RFC3339Regexp),
					resource.TestCheckNoResourceAttr(dataSourceFullName, "terminates_at"),
					resource.TestCheckResourceAttrWith(dataSourceFullName, "assigned_environments_count", func(value string) error {

						valueInt, err := strconv.Atoi(value)
						if err != nil {
							return err
						}

						if valueInt < 1 {
							return fmt.Errorf("assigned_environments_count should have at least one assigned environment")
						}
						return nil

					}),
					resource.TestCheckResourceAttr(dataSourceFullName, "advanced_services.pingid.included", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "advanced_services.pingid.type", "FULL"),
					resource.TestCheckResourceAttr(dataSourceFullName, "authorize.allow_api_access_management", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "authorize.allow_dynamic_authorization", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "credentials.allow_credentials", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "environments.allow_connections", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "environments.allow_custom_domain", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "environments.allow_custom_schema", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "environments.allow_production", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "environments.max", "500"),
					resource.TestCheckResourceAttr(dataSourceFullName, "environments.regions.#", "4"),
					resource.TestCheckTypeSetElemAttr(dataSourceFullName, "environments.regions.*", "EU"),
					resource.TestCheckTypeSetElemAttr(dataSourceFullName, "environments.regions.*", "NORTH_AMERICA"),
					resource.TestCheckTypeSetElemAttr(dataSourceFullName, "environments.regions.*", "AP"),
					resource.TestCheckTypeSetElemAttr(dataSourceFullName, "environments.regions.*", "CA"),
					resource.TestCheckResourceAttr(dataSourceFullName, "fraud.allow_bot_malicious_device_detection", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "fraud.allow_account_protection", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "gateways.allow_ldap_gateway", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "gateways.allow_kerberos_gateway", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "gateways.allow_radius_gateway", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "intelligence.allow_geo_velocity", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "intelligence.allow_anonymous_network_detection", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "intelligence.allow_reputation", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "intelligence.allow_data_consent", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "intelligence.allow_risk", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "intelligence.allow_advanced_predictors", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "mfa.allow_push_notification", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "mfa.allow_notification_outside_whitelist", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "mfa.allow_fido2_devices", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "mfa.allow_voice_otp", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "mfa.allow_email_otp", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "mfa.allow_sms_otp", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "mfa.allow_totp", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "orchestrate.allow_orchestration", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "users.allow_password_management_notifications", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "users.allow_identity_providers", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "users.allow_my_account", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "users.allow_password_only_authentication", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "users.allow_password_policy", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "users.allow_provisioning", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "users.allow_inbound_provisioning", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "users.allow_role_assignment", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "users.allow_verification_flow", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "users.allow_update_self", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "users.entitled_to_support", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "users.max", "10000000"),
					resource.TestCheckResourceAttr(dataSourceFullName, "users.max_hard_limit", "11000000"),
					resource.TestCheckResourceAttr(dataSourceFullName, "users.annual_active_included", "10000000"),
					resource.TestCheckNoResourceAttr(dataSourceFullName, "users.monthly_active_included"),
					resource.TestCheckResourceAttr(dataSourceFullName, "verify.allow_push_notifications", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "verify.allow_document_match", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "verify.allow_face_match", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "verify.allow_manual_id_inspection", "false"),
				),
			},
		},
	})
}

func TestAccLicenseDataSource_NotFound(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	organizationID := os.Getenv("PINGONE_ORGANIZATION_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// {
			// 	Config:      testAccLicenseDataSourceConfig_NotFoundByName(resourceName),
			// 	ExpectError: regexp.MustCompile("Cannot find license doesnotexist"),
			// },
			{
				Config:      testAccLicenseDataSourceConfig_NotFoundByID(resourceName, organizationID),
				ExpectError: regexp.MustCompile("Error when calling `ReadOneLicense`: The request could not be completed. The requested resource was not found."),
			},
		},
	})
}

func testAccLicenseDataSourceConfig_ByIDFull(resourceName, organizationID, licenseID string) string {
	return fmt.Sprintf(`
data "pingone_license" "%[1]s" {
  organization_id = "%[2]s"
  license_id      = "%[3]s"
}`, resourceName, organizationID, licenseID)
}

func testAccLicenseDataSourceConfig_NotFoundByID(resourceName, organizationID string) string {
	return fmt.Sprintf(`
data "pingone_license" "%[1]s" {
  organization_id = "%[2]s"
  license_id      = "9c052a8a-14be-44e4-8f07-2662569994ce" // dummy ID that conforms to UUID v4
}`, resourceName, organizationID)
}
