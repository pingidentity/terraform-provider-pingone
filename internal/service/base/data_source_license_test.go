package base_test

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
)

func testAccCheckLicenseDestroy(s *terraform.State) error {
	return nil
}

func TestAccLicenseDataSource_ByIDFull(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	dataSourceFullName := fmt.Sprintf("data.pingone_license.%s", resourceName)

	organizationID := os.Getenv("PINGONE_ORGANIZATION_ID")
	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironmentAndOrganisation(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckLicenseDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLicenseDataSourceConfig_ByIDFull(resourceName, organizationID, licenseID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(dataSourceFullName, "license_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(dataSourceFullName, "name", "INTERNAL"),
					resource.TestCheckResourceAttr(dataSourceFullName, "package", "INTERNAL"),
					resource.TestCheckResourceAttr(dataSourceFullName, "status", "ACTIVE"),
					resource.TestMatchResourceAttr(dataSourceFullName, "replaces_license_id", regexp.MustCompile(`^([0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})$|^()$`)),
					resource.TestMatchResourceAttr(dataSourceFullName, "replaced_by_license_id", regexp.MustCompile(`^([0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})$|^()$`)),
					resource.TestMatchResourceAttr(dataSourceFullName, "begins_at", regexp.MustCompile(`^(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d:[0-5]\d\.\d+)|(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d:[0-5]\d)|(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d)$`)),
					resource.TestMatchResourceAttr(dataSourceFullName, "expires_at", regexp.MustCompile(`^(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d:[0-5]\d\.\d+)|(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d:[0-5]\d)|(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d)$`)),
					resource.TestMatchResourceAttr(dataSourceFullName, "terminates_at", regexp.MustCompile(`^((\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d:[0-5]\d\.\d+)|(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d:[0-5]\d)|(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d))$|^()$`)),
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
					resource.TestCheckResourceAttr(dataSourceFullName, "advanced_services.#", "1"),
					resource.TestCheckResourceAttr(dataSourceFullName, "advanced_services.0.pingid.#", "1"),
					resource.TestCheckResourceAttr(dataSourceFullName, "advanced_services.0.pingid.0.included", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "advanced_services.0.pingid.0.type", "FULL"),
					resource.TestCheckResourceAttr(dataSourceFullName, "authorize.#", "1"),
					resource.TestCheckResourceAttr(dataSourceFullName, "authorize.0.allow_api_access_management", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "authorize.0.allow_dynamic_authorization", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "credentials.#", "1"),
					resource.TestCheckResourceAttr(dataSourceFullName, "credentials.0.allow_credentials", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "environments.#", "1"),
					resource.TestCheckResourceAttr(dataSourceFullName, "environments.0.allow_connections", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "environments.0.allow_custom_domain", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "environments.0.allow_custom_schema", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "environments.0.allow_production", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "environments.0.max", "50"),
					resource.TestCheckResourceAttr(dataSourceFullName, "environments.0.regions.#", "4"),
					resource.TestCheckTypeSetElemAttr(dataSourceFullName, "environments.0.regions.*", "EU"),
					resource.TestCheckTypeSetElemAttr(dataSourceFullName, "environments.0.regions.*", "NORTH_AMERICA"),
					resource.TestCheckTypeSetElemAttr(dataSourceFullName, "environments.0.regions.*", "AP"),
					resource.TestCheckTypeSetElemAttr(dataSourceFullName, "environments.0.regions.*", "CA"),
					resource.TestCheckResourceAttr(dataSourceFullName, "fraud.#", "1"),
					resource.TestCheckResourceAttr(dataSourceFullName, "fraud.0.allow_bot_malicious_device_detection", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "fraud.0.allow_account_protection", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "gateways.#", "1"),
					resource.TestCheckResourceAttr(dataSourceFullName, "gateways.0.allow_ldap_gateway", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "gateways.0.allow_kerberos_gateway", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "gateways.0.allow_radius_gateway", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "intelligence.#", "1"),
					resource.TestCheckResourceAttr(dataSourceFullName, "intelligence.0.allow_geo_velocity", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "intelligence.0.allow_anonymous_network_detection", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "intelligence.0.allow_reputation", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "intelligence.0.allow_data_consent", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "intelligence.0.allow_risk", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "intelligence.0.allow_advanced_predictors", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "mfa.#", "1"),
					resource.TestCheckResourceAttr(dataSourceFullName, "mfa.0.allow_push_notification", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "mfa.0.allow_notification_outside_whitelist", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "mfa.0.allow_fido2_devices", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "mfa.0.allow_voice_otp", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "mfa.0.allow_email_otp", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "mfa.0.allow_sms_otp", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "mfa.0.allow_totp", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "orchestrate.#", "1"),
					resource.TestCheckResourceAttr(dataSourceFullName, "orchestrate.0.allow_orchestration", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "users.#", "1"),
					resource.TestCheckResourceAttr(dataSourceFullName, "users.0.allow_password_management_notifications", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "users.0.allow_identity_providers", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "users.0.allow_my_account", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "users.0.allow_password_only_authentication", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "users.0.allow_password_policy", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "users.0.allow_provisioning", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "users.0.allow_inbound_provisioning", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "users.0.allow_role_assignment", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "users.0.allow_verification_flow", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "users.0.allow_update_self", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "users.0.entitled_to_support", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "users.0.max", "10000000"),
					resource.TestCheckResourceAttr(dataSourceFullName, "users.0.max_hard_limit", "11000000"),
					resource.TestCheckResourceAttr(dataSourceFullName, "users.0.annual_active_included", "10000000"),
					resource.TestCheckResourceAttr(dataSourceFullName, "users.0.monthly_active_included", "0"),
					resource.TestCheckResourceAttr(dataSourceFullName, "verify.#", "1"),
					resource.TestCheckResourceAttr(dataSourceFullName, "verify.0.allow_push_notifications", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "verify.0.allow_document_match", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "verify.0.allow_face_match", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "verify.0.allow_manual_id_inspection", "false"),
				),
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
