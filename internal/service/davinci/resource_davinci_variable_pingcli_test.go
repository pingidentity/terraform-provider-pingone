// Copyright © 2026 Ping Identity Corporation

package davinci_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	acctestlegacysdk "github.com/pingidentity/terraform-provider-pingone/internal/acctest/legacysdk"
)

func TestAccDavinciVariable_PingCLIAuth_WithEnvironment(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	environmentName := acctest.ResourceNameGenEnvironment()
	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckPingCLIAuthOnly(t)
			acctest.PreCheckPingCLIAuthClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             davinciVariablePingCLI_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: davinciVariablePingCLI_NewEnvHCL(environmentName, licenseID, resourceName),
			},
		},
	})
}

func davinciVariablePingCLI_NewEnvHCL(environmentName, licenseID, resourceName string) string {
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

func davinciVariablePingCLI_CheckDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClientPingCLI(ctx)
	if err != nil {
		return err
	}

	return davinciVariable_CheckDestroyWithClient(ctx, s, p1Client)
}
