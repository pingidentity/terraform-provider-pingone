package risk_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/patrickcping/pingone-go-sdk-v2/risk"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func testAccCheckRiskPredictorDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.RiskAPIClient
	ctx = context.WithValue(ctx, risk.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})

	apiClientManagement := p1Client.API.ManagementAPIClient
	ctxManagement := context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_risk_predictor" {
			continue
		}

		_, rEnv, err := apiClientManagement.EnvironmentsApi.ReadOneEnvironment(ctxManagement, rs.Primary.Attributes["environment_id"]).Execute()

		if err != nil {

			if rEnv == nil {
				return fmt.Errorf("Response object does not exist and no error detected")
			}

			if rEnv.StatusCode == 404 {
				continue
			}

			return err
		}

		body, r, err := apiClient.RiskAdvancedPredictorsApi.ReadOneRiskPredictor(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.ID).Execute()

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

		return fmt.Errorf("PingOne risk predictor %s still exists", rs.Primary.ID)
	}

	return nil
}

func TestAccRiskPredictor_NewEnv(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_risk_predictor.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckRiskPredictorDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccRiskPredictorConfig_NewEnv(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
				),
			},
		},
	})
}

func TestAccRiskPredictor_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_risk_predictor.%s", resourceName)

	name := resourceName

	fullCheck := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
		resource.TestCheckResourceAttr(resourceFullName, "name", name),
		resource.TestCheckResourceAttr(resourceFullName, "compact_name", fmt.Sprintf("%s1", name)),
		resource.TestCheckResourceAttr(resourceFullName, "description", "When my wife is upset, I let her colour in my black and white tattoos.  She just needs a shoulder to crayon.."),
		resource.TestCheckResourceAttr(resourceFullName, "type", "ANONYMOUS_NETWORK"),
		resource.TestCheckResourceAttr(resourceFullName, "licensed", "true"),
		resource.TestCheckResourceAttr(resourceFullName, "deletable", "true"),
		resource.TestCheckResourceAttr(resourceFullName, "default_decision_value", "MEDIUM"),
	)

	minimalCheck := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
		resource.TestCheckResourceAttr(resourceFullName, "name", name),
		resource.TestCheckResourceAttr(resourceFullName, "compact_name", fmt.Sprintf("%s1", name)),
		resource.TestCheckNoResourceAttr(resourceFullName, "description"),
		resource.TestCheckResourceAttr(resourceFullName, "type", "ANONYMOUS_NETWORK"),
		resource.TestCheckResourceAttr(resourceFullName, "licensed", "true"),
		resource.TestCheckResourceAttr(resourceFullName, "deletable", "true"),
		resource.TestCheckNoResourceAttr(resourceFullName, "default_decision_value"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckRiskPredictorDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			{
				Config: testAccRiskPredictorConfig_Full(resourceName, name),
				Check:  fullCheck,
			},
			{
				Config:  testAccRiskPredictorConfig_Full(resourceName, name),
				Destroy: true,
			},
			// Minimal
			{
				Config: testAccRiskPredictorConfig_Minimal(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config:  testAccRiskPredictorConfig_Minimal(resourceName, name),
				Destroy: true,
			},
			// Change
			{
				Config: testAccRiskPredictorConfig_Full(resourceName, name),
				Check:  fullCheck,
			},
			{
				Config: testAccRiskPredictorConfig_Minimal(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config: testAccRiskPredictorConfig_Full(resourceName, name),
				Check:  fullCheck,
			},
		},
	})
}

func TestAccRiskPredictor_Anonymous_Network(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_risk_predictor.%s", resourceName)

	name := resourceName

	fullCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "type", "ANONYMOUS_NETWORK"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_anonymous_network.#", "1"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_anonymous_network.0.allowed_cidr_list.#", "3"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "predictor_anonymous_network.0.allowed_cidr_list.*", "192.168.0.0/24"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "predictor_anonymous_network.0.allowed_cidr_list.*", "10.0.0.0/8"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "predictor_anonymous_network.0.allowed_cidr_list.*", "172.16.0.0/12"),
	)

	minimalCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "type", "ANONYMOUS_NETWORK"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_anonymous_network.#", "1"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_anonymous_network.0.allowed_cidr_list.#", "0"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckRiskPredictorDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			{
				Config: testAccRiskPredictorConfig_Anonymous_Network_Full(resourceName, name),
				Check:  fullCheck,
			},
			{
				Config:  testAccRiskPredictorConfig_Anonymous_Network_Full(resourceName, name),
				Destroy: true,
			},
			// Minimal
			{
				Config: testAccRiskPredictorConfig_Anonymous_Network_Minimal(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config:  testAccRiskPredictorConfig_Anonymous_Network_Minimal(resourceName, name),
				Destroy: true,
			},
			// Change
			{
				Config: testAccRiskPredictorConfig_Anonymous_Network_Full(resourceName, name),
				Check:  fullCheck,
			},
			{
				Config: testAccRiskPredictorConfig_Anonymous_Network_Minimal(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config: testAccRiskPredictorConfig_Anonymous_Network_Full(resourceName, name),
				Check:  fullCheck,
			},
		},
	})
}

func TestAccRiskPredictor_Anonymous_Network_Override(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_risk_predictor.%s", resourceName)

	name := "Anonymous Network"
	compactName := "anonymousNetwork"

	fullCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "name", "Anonymous Network Detection"),
		resource.TestCheckResourceAttr(resourceFullName, "compact_name", "anonymousNetwork"),
		resource.TestCheckResourceAttr(resourceFullName, "type", "ANONYMOUS_NETWORK"),
		resource.TestCheckResourceAttr(resourceFullName, "deletable", "false"),
		resource.TestCheckResourceAttr(resourceFullName, "default_decision_value", "MEDIUM"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_anonymous_network.#", "1"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_anonymous_network.0.allowed_cidr_list.#", "3"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "predictor_anonymous_network.0.allowed_cidr_list.*", "192.168.0.0/24"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "predictor_anonymous_network.0.allowed_cidr_list.*", "10.0.0.0/8"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "predictor_anonymous_network.0.allowed_cidr_list.*", "172.16.0.0/12"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckRiskPredictorDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			{
				Config: testAccRiskPredictorConfig_Anonymous_Network_Override(resourceName, name, compactName),
				Check:  fullCheck,
			},
		},
	})
}

func TestAccRiskPredictor_Geovelocity(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_risk_predictor.%s", resourceName)

	name := resourceName

	fullCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "type", "GEO_VELOCITY"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_geovelocity.#", "1"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_geovelocity.0.allowed_cidr_list.#", "3"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "predictor_geovelocity.0.allowed_cidr_list.*", "192.168.0.0/24"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "predictor_geovelocity.0.allowed_cidr_list.*", "10.0.0.0/8"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "predictor_geovelocity.0.allowed_cidr_list.*", "172.16.0.0/12"),
	)

	minimalCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "type", "GEO_VELOCITY"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_geovelocity.#", "1"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_geovelocity.0.allowed_cidr_list.#", "0"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckRiskPredictorDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			{
				Config: testAccRiskPredictorConfig_Geovelocity_Full(resourceName, name),
				Check:  fullCheck,
			},
			{
				Config:  testAccRiskPredictorConfig_Geovelocity_Full(resourceName, name),
				Destroy: true,
			},
			// Minimal
			{
				Config: testAccRiskPredictorConfig_Geovelocity_Minimal(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config:  testAccRiskPredictorConfig_Geovelocity_Minimal(resourceName, name),
				Destroy: true,
			},
			// Change
			{
				Config: testAccRiskPredictorConfig_Geovelocity_Full(resourceName, name),
				Check:  fullCheck,
			},
			{
				Config: testAccRiskPredictorConfig_Geovelocity_Minimal(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config: testAccRiskPredictorConfig_Geovelocity_Full(resourceName, name),
				Check:  fullCheck,
			},
		},
	})
}

func TestAccRiskPredictor_Geovelocity_Override(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_risk_predictor.%s", resourceName)

	name := "GeoVelocity"
	compactName := "geoVelocity"

	fullCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "name", "GeoVelocity"),
		resource.TestCheckResourceAttr(resourceFullName, "compact_name", "geoVelocity"),
		resource.TestCheckResourceAttr(resourceFullName, "type", "GEO_VELOCITY"),
		resource.TestCheckResourceAttr(resourceFullName, "deletable", "false"),
		resource.TestCheckResourceAttr(resourceFullName, "default_decision_value", "MEDIUM"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_geovelocity.#", "1"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_geovelocity.0.allowed_cidr_list.#", "3"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "predictor_geovelocity.0.allowed_cidr_list.*", "192.168.0.0/24"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "predictor_geovelocity.0.allowed_cidr_list.*", "10.0.0.0/8"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "predictor_geovelocity.0.allowed_cidr_list.*", "172.16.0.0/12"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckRiskPredictorDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			{
				Config: testAccRiskPredictorConfig_Geovelocity_Override(resourceName, name, compactName),
				Check:  fullCheck,
			},
		},
	})
}

func TestAccRiskPredictor_IP_Reputation(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_risk_predictor.%s", resourceName)

	name := resourceName

	fullCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "type", "IP_REPUTATION"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_ip_reputation.#", "1"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_ip_reputation.0.allowed_cidr_list.#", "3"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "predictor_ip_reputation.0.allowed_cidr_list.*", "192.168.0.0/24"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "predictor_ip_reputation.0.allowed_cidr_list.*", "10.0.0.0/8"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "predictor_ip_reputation.0.allowed_cidr_list.*", "172.16.0.0/12"),
	)

	minimalCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "type", "IP_REPUTATION"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_ip_reputation.#", "1"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_ip_reputation.0.allowed_cidr_list.#", "0"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckRiskPredictorDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			{
				Config: testAccRiskPredictorConfig_IPReputation_Full(resourceName, name),
				Check:  fullCheck,
			},
			{
				Config:  testAccRiskPredictorConfig_IPReputation_Full(resourceName, name),
				Destroy: true,
			},
			// Minimal
			{
				Config: testAccRiskPredictorConfig_IPReputation_Minimal(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config:  testAccRiskPredictorConfig_IPReputation_Minimal(resourceName, name),
				Destroy: true,
			},
			// Change
			{
				Config: testAccRiskPredictorConfig_IPReputation_Full(resourceName, name),
				Check:  fullCheck,
			},
			{
				Config: testAccRiskPredictorConfig_IPReputation_Minimal(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config: testAccRiskPredictorConfig_IPReputation_Full(resourceName, name),
				Check:  fullCheck,
			},
		},
	})
}

func TestAccRiskPredictor_IPReputation_Override(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_risk_predictor.%s", resourceName)

	name := "IP Reputation"
	compactName := "ipRisk"

	fullCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "name", "IP Reputation"),
		resource.TestCheckResourceAttr(resourceFullName, "compact_name", "ipRisk"),
		resource.TestCheckResourceAttr(resourceFullName, "type", "IP_REPUTATION"),
		resource.TestCheckResourceAttr(resourceFullName, "deletable", "false"),
		resource.TestCheckResourceAttr(resourceFullName, "default_decision_value", "MEDIUM"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_ip_reputation.#", "1"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_ip_reputation.0.allowed_cidr_list.#", "3"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "predictor_ip_reputation.0.allowed_cidr_list.*", "192.168.0.0/24"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "predictor_ip_reputation.0.allowed_cidr_list.*", "10.0.0.0/8"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "predictor_ip_reputation.0.allowed_cidr_list.*", "172.16.0.0/12"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckRiskPredictorDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			{
				Config: testAccRiskPredictorConfig_IPReputation_Override(resourceName, name, compactName),
				Check:  fullCheck,
			},
		},
	})
}

func TestAccRiskPredictor_NewDevice(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_risk_predictor.%s", resourceName)

	name := resourceName

	fullCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "type", "DEVICE"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_new_device.#", "1"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_new_device.0.detect", "NEW_DEVICE"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_new_device.0.activation_at", "2023-05-02T00:00:00Z"),
	)

	minimalCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "type", "DEVICE"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_new_device.#", "1"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_new_device.0.detect", "NEW_DEVICE"),
		resource.TestCheckNoResourceAttr(resourceFullName, "predictor_new_device.0.activation_at"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckRiskPredictorDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			{
				Config: testAccRiskPredictorConfig_NewDevice_Full(resourceName, name),
				Check:  fullCheck,
			},
			{
				Config:  testAccRiskPredictorConfig_NewDevice_Full(resourceName, name),
				Destroy: true,
			},
			// Minimal
			{
				Config: testAccRiskPredictorConfig_NewDevice_Minimal(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config:  testAccRiskPredictorConfig_NewDevice_Minimal(resourceName, name),
				Destroy: true,
			},
			// Change
			{
				Config: testAccRiskPredictorConfig_NewDevice_Full(resourceName, name),
				Check:  fullCheck,
			},
			{
				Config: testAccRiskPredictorConfig_NewDevice_Minimal(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config: testAccRiskPredictorConfig_NewDevice_Full(resourceName, name),
				Check:  fullCheck,
			},
		},
	})
}

func TestAccRiskPredictor_NewDevice_Override(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_risk_predictor.%s", resourceName)

	name := "New Device"
	compactName := "newDevice"

	fullCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "name", "New Device"),
		resource.TestCheckResourceAttr(resourceFullName, "compact_name", "newDevice"),
		resource.TestCheckResourceAttr(resourceFullName, "type", "DEVICE"),
		resource.TestCheckResourceAttr(resourceFullName, "deletable", "false"),
		resource.TestCheckResourceAttr(resourceFullName, "default_decision_value", "MEDIUM"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_new_device.#", "1"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_new_device.0.detect", "NEW_DEVICE"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_new_device.0.activation_at", "2023-05-02T00:00:00Z"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckRiskPredictorDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			{
				Config: testAccRiskPredictorConfig_NewDevice_Override(resourceName, name, compactName),
				Check:  fullCheck,
			},
		},
	})
}

func TestAccRiskPredictor_UserLocationAnomaly(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_risk_predictor.%s", resourceName)

	name := resourceName

	fullCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "type", "USER_LOCATION_ANOMALY"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_user_location_anomaly.#", "1"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_user_location_anomaly.0.radius_distance", "100"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_user_location_anomaly.0.radius_distance_unit", "miles"),
	)

	minimalCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "type", "USER_LOCATION_ANOMALY"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_user_location_anomaly.#", "1"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_user_location_anomaly.0.radius_distance", "51"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_user_location_anomaly.0.radius_distance_unit", "kilometers"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckRiskPredictorDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			{
				Config: testAccRiskPredictorConfig_UserLocationAnomaly_Full(resourceName, name),
				Check:  fullCheck,
			},
			{
				Config:  testAccRiskPredictorConfig_UserLocationAnomaly_Full(resourceName, name),
				Destroy: true,
			},
			// Minimal
			{
				Config: testAccRiskPredictorConfig_UserLocationAnomaly_Minimal(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config:  testAccRiskPredictorConfig_UserLocationAnomaly_Minimal(resourceName, name),
				Destroy: true,
			},
			// Change
			{
				Config: testAccRiskPredictorConfig_UserLocationAnomaly_Full(resourceName, name),
				Check:  fullCheck,
			},
			{
				Config: testAccRiskPredictorConfig_UserLocationAnomaly_Minimal(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config: testAccRiskPredictorConfig_UserLocationAnomaly_Full(resourceName, name),
				Check:  fullCheck,
			},
		},
	})
}

func TestAccRiskPredictor_UserLocationAnomaly_Override(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_risk_predictor.%s", resourceName)

	name := "User Location Anomaly"
	compactName := "userLocationAnomaly"

	fullCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "name", "User Location Anomaly"),
		resource.TestCheckResourceAttr(resourceFullName, "compact_name", "userLocationAnomaly"),
		resource.TestCheckResourceAttr(resourceFullName, "type", "USER_LOCATION_ANOMALY"),
		resource.TestCheckResourceAttr(resourceFullName, "deletable", "false"),
		resource.TestCheckResourceAttr(resourceFullName, "default_decision_value", "MEDIUM"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_user_location_anomaly.#", "1"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_user_location_anomaly.0.radius_distance", "100"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_user_location_anomaly.0.radius_distance_unit", "miles"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckRiskPredictorDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			{
				Config: testAccRiskPredictorConfig_UserLocationAnomaly_Override(resourceName, name, compactName),
				Check:  fullCheck,
			},
		},
	})
}

func testAccRiskPredictorConfig_NewEnv(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_risk_predictor" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  name         = "%[4]s"
  compact_name = "%[4]s1"

  predictor_anonymous_network {
    allowed_cidr_list = []
  }
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccRiskPredictorConfig_Full(resourceName, name string) string {
	return testAccRiskPredictorConfig_Anonymous_Network_Full(resourceName, name)
}

func testAccRiskPredictorConfig_Minimal(resourceName, name string) string {
	return testAccRiskPredictorConfig_Anonymous_Network_Minimal(resourceName, name)
}

func testAccRiskPredictorConfig_Anonymous_Network_Full(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_predictor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s"
  compact_name = "%[3]s1"
  description  = "When my wife is upset, I let her colour in my black and white tattoos.  She just needs a shoulder to crayon.."

  default_decision_value = "MEDIUM"

  predictor_anonymous_network {
    allowed_cidr_list = [
      "10.0.0.0/8",
      "172.16.0.0/12",
      "192.168.0.0/24"
    ]
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccRiskPredictorConfig_Anonymous_Network_Minimal(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_predictor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s"
  compact_name = "%[3]s1"

  predictor_anonymous_network {
    allowed_cidr_list = []
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccRiskPredictorConfig_Anonymous_Network_Override(resourceName, name, compactName string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_predictor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s"
  compact_name = "%[4]s"

  default_decision_value = "MEDIUM"

  predictor_anonymous_network {
    allowed_cidr_list = [
      "10.0.0.0/8",
      "172.16.0.0/12",
      "192.168.0.0/24"
    ]
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name, compactName)
}

func testAccRiskPredictorConfig_Geovelocity_Full(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_predictor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s"
  compact_name = "%[3]s1"
  description  = "When my wife is upset, I let her colour in my black and white tattoos.  She just needs a shoulder to crayon.."

  default_decision_value = "MEDIUM"

  predictor_geovelocity {
    allowed_cidr_list = [
      "10.0.0.0/8",
      "172.16.0.0/12",
      "192.168.0.0/24"
    ]
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccRiskPredictorConfig_Geovelocity_Minimal(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_predictor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s"
  compact_name = "%[3]s1"

  predictor_geovelocity {
    allowed_cidr_list = []
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccRiskPredictorConfig_Geovelocity_Override(resourceName, name, compactName string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_predictor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s"
  compact_name = "%[4]s"

  default_decision_value = "MEDIUM"

  predictor_geovelocity {
    allowed_cidr_list = [
      "10.0.0.0/8",
      "172.16.0.0/12",
      "192.168.0.0/24"
    ]
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name, compactName)
}

func testAccRiskPredictorConfig_IPReputation_Full(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_predictor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s"
  compact_name = "%[3]s1"
  description  = "When my wife is upset, I let her colour in my black and white tattoos.  She just needs a shoulder to crayon.."

  default_decision_value = "MEDIUM"

  predictor_ip_reputation {
    allowed_cidr_list = [
      "10.0.0.0/8",
      "172.16.0.0/12",
      "192.168.0.0/24"
    ]
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccRiskPredictorConfig_IPReputation_Minimal(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_predictor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s"
  compact_name = "%[3]s1"

  predictor_ip_reputation {
    allowed_cidr_list = []
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccRiskPredictorConfig_IPReputation_Override(resourceName, name, compactName string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_predictor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s"
  compact_name = "%[4]s"

  default_decision_value = "MEDIUM"

  predictor_ip_reputation {
    allowed_cidr_list = [
      "10.0.0.0/8",
      "172.16.0.0/12",
      "192.168.0.0/24"
    ]
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name, compactName)
}

func testAccRiskPredictorConfig_NewDevice_Full(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_predictor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s"
  compact_name = "%[3]s1"
  description  = "When my wife is upset, I let her colour in my black and white tattoos.  She just needs a shoulder to crayon.."

  default_decision_value = "MEDIUM"

  predictor_new_device {
    detect        = "NEW_DEVICE"
    activation_at = "2023-05-02T00:00:00Z"
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccRiskPredictorConfig_NewDevice_Minimal(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_predictor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s"
  compact_name = "%[3]s1"

  predictor_new_device {
    detect = "NEW_DEVICE"
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccRiskPredictorConfig_NewDevice_Override(resourceName, name, compactName string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_predictor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s"
  compact_name = "%[4]s"

  default_decision_value = "MEDIUM"

  predictor_new_device {
    detect        = "NEW_DEVICE"
    activation_at = "2023-05-02T00:00:00Z"
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name, compactName)
}

func testAccRiskPredictorConfig_UserLocationAnomaly_Full(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_predictor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s"
  compact_name = "%[3]s1"
  description  = "When my wife is upset, I let her colour in my black and white tattoos.  She just needs a shoulder to crayon.."

  default_decision_value = "MEDIUM"

  predictor_user_location_anomaly {
    radius_distance      = 100
    radius_distance_unit = "miles"
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccRiskPredictorConfig_UserLocationAnomaly_Minimal(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_predictor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s"
  compact_name = "%[3]s1"

  predictor_user_location_anomaly {
    radius_distance = 51
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccRiskPredictorConfig_UserLocationAnomaly_Override(resourceName, name, compactName string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_predictor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s"
  compact_name = "%[4]s"

  default_decision_value = "MEDIUM"

  predictor_user_location_anomaly {
    radius_distance      = 100
    radius_distance_unit = "miles"
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name, compactName)
}
