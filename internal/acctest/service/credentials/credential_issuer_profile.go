package credentials

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func CredentialIssuerProfile_CheckDestroy(s *terraform.State) error {
	return nil

	// Note: Issuer Profiles aren't deleted once created. Uncomment and replace Passthrough if this changes.
	/*
	   var ctx = context.Background()

	   p1Client, err := acctest.TestClient(ctx)

	   	if err != nil {
	   		return err
	   	}

	   apiClient := p1Client.API.CredentialsAPIClient

	   mgmtApiClient := p1Client.API.ManagementAPIClient

	   	for _, rs := range s.RootModule().Resources {
	   		if rs.Type != "pingone_credential_issuer_profile" {
	   			continue
	   		}

	   		_, rEnv, err := mgmtApiClient.EnvironmentsApi.ReadOneEnvironment(ctx, rs.Primary.Attributes["environment_id"]).Execute()

	   		if err != nil {

	   			if rEnv == nil {
	   				return fmt.Errorf("Response object does not exist and no error detected")
	   			}

	   			if rEnv.StatusCode == 404 {
	   				continue
	   			}

	   			return err
	   		}

	   		body, r, err := apiClient.CredentialIssuersApi.ReadCredentialIssuerProfile(ctx, rs.Primary.Attributes["environment_id"]).Execute()

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

	   		return fmt.Errorf("PingOne Credential Issuer Profile %s still exists", rs.Primary.ID)
	   	}

	   return nil
	*/
}

func CredentialIssuerProfile_GetIDs(resourceName string, environmentID, resourceID *string) resource.TestCheckFunc {
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
