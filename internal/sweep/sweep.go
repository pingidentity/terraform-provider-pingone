package sweep

import (
	"context"
	"os"

	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
)

func SweepClient(ctx context.Context) (*client.Client, error) {

	config := &client.Config{
		ClientID:      os.Getenv("PINGONE_CLIENT_ID"),
		ClientSecret:  os.Getenv("PINGONE_CLIENT_SECRET"),
		EnvironmentID: os.Getenv("PINGONE_ENVIRONMENT_ID"),
		Region:        os.Getenv("PINGONE_REGION"),
		ForceDelete:   false,
	}

	return config.APIClient(ctx)

}
