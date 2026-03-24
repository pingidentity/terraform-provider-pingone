// Copyright © 2026 Ping Identity Corporation

//go:build beta

package davinci_test

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

func TestAccDavinciVariable_AccessTokenAuth_WithEnvironment(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	environmentName := acctest.ResourceNameGenEnvironment()
	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckAccessTokenOnly(t)
			acctest.PreCheckAccessTokenClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             davinciVariableAccessToken_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: davinciVariableAccessToken_NewEnvHCL(environmentName, licenseID, resourceName),
			},
		},
	})
}

func davinciVariableAccessToken_NewEnvHCL(environmentName, licenseID, resourceName string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_davinci_variable" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  context        = "flowInstance"
  data_type      = "string"
  mutable        = true
  name           = "%[2]s"
}
`, acctestlegacysdk.MinimalSandboxEnvironmentNoPopulation(environmentName, licenseID), environmentName, resourceName)
}

func davinciVariableAccessToken_CheckDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClientAccessToken(ctx)
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		switch rs.Type {
		case "pingone_davinci_variable":
			shouldContinue, err := acctest.CheckParentEnvironmentDestroy(ctx, p1Client, rs.Primary.Attributes["environment_id"])
			if err != nil {
				return err
			}

			if shouldContinue {
				continue
			}

			_, r, err := p1Client.DaVinciVariablesApi.GetVariableById(ctx, uuid.MustParse(rs.Primary.Attributes["environment_id"]), uuid.MustParse(rs.Primary.Attributes["id"])).Execute()
			shouldContinue, err = acctest.CheckForResourceDestroy(r, err)
			if err != nil {
				return err
			}

			if shouldContinue {
				continue
			}

			return fmt.Errorf("PingOne davinci_variable Instance %s still exists", rs.Primary.ID)
		case "pingone_environment":
			environmentID, err := uuid.Parse(rs.Primary.ID)
			if err != nil {
				return fmt.Errorf("unable to parse environment id '%s' as uuid: %v", rs.Primary.ID, err)
			}

			_, r, err := p1Client.EnvironmentsApi.GetEnvironmentById(ctx, environmentID).Execute()
			shouldContinue, err := acctest.CheckForResourceDestroy(r, err)
			if err != nil {
				return err
			}

			if shouldContinue {
				continue
			}

			return fmt.Errorf("PingOne environment instance %s still exists", rs.Primary.ID)
		}
	}

	return nil
}
