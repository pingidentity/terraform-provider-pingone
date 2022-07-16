package sweep

import (
	"context"
	"fmt"
	"os"

	"github.com/patrickcping/pingone-go-sdk-v2/management"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
)

func SweepClient(ctx context.Context) (*client.Client, error) {

	config := &client.Config{
		ClientID:      os.Getenv("PINGONE_CLIENT_ID"),
		ClientSecret:  os.Getenv("PINGONE_CLIENT_SECRET"),
		EnvironmentID: os.Getenv("PINGONE_ENVIRONMENT_ID"),
		Region:        os.Getenv("PINGONE_REGION"),
		ForceDelete:   true,
	}

	return config.APIClient(ctx)

}

func FetchTaggedEnvironments(ctx context.Context, apiClient *management.APIClient, region string) ([]management.Environment, error) {

	filter := "name sw \"tf-testacc-\""

	respList, _, err := apiClient.EnvironmentsApi.ReadAllEnvironments(ctx).Filter(filter).Execute()
	if err != nil {
		return nil, fmt.Errorf("Error getting environments: %s", err)
	}

	if environments, ok := respList.Embedded.GetEnvironmentsOk(); ok {

		for _, environment := range environments {
			if environment.GetName() == "Administrators" {
				return nil, fmt.Errorf("Unsafe filter, Administrators environment present: %s", filter)
			}
		}
		return environments, nil
	} else {
		return make([]management.Environment, 0), nil
	}

}
