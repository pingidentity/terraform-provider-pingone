// Copyright © 2026 Ping Identity Corporation

package service_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	acctestlegacysdk "github.com/pingidentity/terraform-provider-pingone/internal/acctest/legacysdk"
)

// The Ping CLI acceptance tests exercise the feature through a real keychain
// session, one per provider flavour. Ping CLI authentication is enabled by the
// PINGONE_TF_USE_PINGCLI_AUTH environment variable, so the configurations have
// no provider block:
//
//	pingone_language          -> SDKv2 flavour
//	pingone_environment       -> legacy SDK flavour
//	pingone_davinci_variable  -> new SDK framework flavour

func TestAccPingCLIAuth_SDKv2Flavour(t *testing.T) {
	t.Parallel()

	environmentName := acctest.ResourceNameGenEnvironment()
	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckPingCLIAuth(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPingCLIAuthConfig_SDKv2(environmentName, licenseID),
			},
		},
		CheckDestroy: testAccPingCLIAuthCheckDestroyEnvironment,
	})
}

func TestAccPingCLIAuth_LegacyFrameworkFlavour(t *testing.T) {
	t.Parallel()

	environmentName := acctest.ResourceNameGenEnvironment()
	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckPingCLIAuth(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPingCLIAuthConfig_Legacy(environmentName, licenseID),
			},
		},
		CheckDestroy: testAccPingCLIAuthCheckDestroyEnvironment,
	})
}

func TestAccPingCLIAuth_FrameworkFlavour_DaVinciVariable(t *testing.T) {
	t.Parallel()

	environmentName := acctest.ResourceNameGenEnvironment()
	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckPingCLIAuth(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPingCLIAuthConfig_Framework(environmentName, licenseID),
			},
		},
		CheckDestroy: testAccPingCLIAuthCheckDestroyFramework,
	})
}

func testAccPingCLIAuthConfig_SDKv2(environmentName, licenseID string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_language" "test" {
  environment_id = pingone_environment.%[2]s.id
  locale         = "de-DE"
}`, acctestlegacysdk.MinimalSandboxEnvironmentNoPopulation(environmentName, licenseID), environmentName)
}

func testAccPingCLIAuthConfig_Legacy(environmentName, licenseID string) string {
	return acctestlegacysdk.MinimalSandboxEnvironmentNoPopulation(environmentName, licenseID)
}

func testAccPingCLIAuthConfig_Framework(environmentName, licenseID string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_davinci_variable" "test" {
  environment_id = pingone_environment.%[2]s.id
  context        = "flowInstance"
  data_type      = "string"
  mutable        = true
  name           = "%[2]s"
}`, acctestlegacysdk.MinimalSandboxEnvironmentNoPopulation(environmentName, licenseID), environmentName)
}

// testAccPingCLIAuthCheckDestroyEnvironment verifies the environment created by
// the test configuration is destroyed. The environment is the resource the
// Ping CLI session created; the in-environment resource's own lifecycle is
// covered by its service acceptance tests.
func testAccPingCLIAuthCheckDestroyEnvironment(s *terraform.State) error {
	ctx := context.Background()

	p1Client, err := acctest.TestClientPingCLI(ctx)
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_environment" {
			continue
		}

		shouldContinue, err := acctest.CheckParentEnvironmentDestroy(ctx, p1Client, rs.Primary.ID)
		if err != nil {
			return err
		}
		if shouldContinue {
			continue
		}

		return fmt.Errorf("pingone_environment %s still exists", rs.Primary.ID)
	}

	return nil
}

func testAccPingCLIAuthCheckDestroyFramework(s *terraform.State) error {
	ctx := context.Background()

	p1Client, err := acctest.TestClientPingCLI(ctx)
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_davinci_variable" {
			continue
		}

		shouldContinue, err := acctest.CheckParentEnvironmentDestroy(ctx, p1Client, rs.Primary.Attributes["environment_id"])
		if err != nil {
			return err
		}
		if shouldContinue {
			continue
		}

		_, r, err := p1Client.DaVinciVariablesApi.GetVariableById(ctx, uuid.MustParse(rs.Primary.Attributes["environment_id"]), uuid.MustParse(rs.Primary.ID)).Execute()
		shouldContinue, err = acctest.CheckForResourceDestroy(r, err)
		if err != nil {
			return err
		}
		if shouldContinue {
			continue
		}

		return fmt.Errorf("pingone_davinci_variable %s still exists", rs.Primary.ID)
	}

	return nil
}
