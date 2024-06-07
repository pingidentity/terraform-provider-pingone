package authorize

import (
	"context"
	"testing"

	"github.com/patrickcping/pingone-go-sdk-v2/pingone"
)

func ApplicationResource_RemovalDrift_PreConfig(ctx context.Context, apiClient *pingone.Client, t *testing.T, environmentID, resourceID string) {
	if environmentID == "" || resourceID == "" {
		t.Fatalf("One of environment ID, custom resource ID or application resource ID cannot be determined. Environment ID: %s, Application resource ID: %s", environmentID, resourceID)
	}

	applicationResource, r, err := apiClient.AuthorizeAPIClient.ApplicationResourcesApi.ReadOneApplicationResource(ctx, environmentID, resourceID).Execute()
	if err != nil {
		t.Fatalf("Failed to get application resource for delete (authorize): %v", err)
	}
	if r.StatusCode == 404 {
		t.Fatalf("Application Resource %s not found", resourceID)
	}

	_, err = apiClient.ManagementAPIClient.ApplicationResourcesApi.DeleteApplicationResource(ctx, environmentID, applicationResource.Parent.GetId(), resourceID).Execute()
	if err != nil {
		t.Fatalf("Failed to delete application resource (authorize): %v", err)
	}
}

func ApplicationResource_Resource_RemovalDrift_PreConfig(ctx context.Context, apiClient *pingone.Client, t *testing.T, environmentID, resourceID string) {
	if environmentID == "" || resourceID == "" {
		t.Fatalf("One of environment ID, custom resource ID or application resource ID cannot be determined. Environment ID: %s, Application resource ID: %s", environmentID, resourceID)
	}

	applicationResource, r, err := apiClient.AuthorizeAPIClient.ApplicationResourcesApi.ReadOneApplicationResource(ctx, environmentID, resourceID).Execute()
	if err != nil {
		t.Fatalf("Failed to get application resource for delete (authorize): %v", err)
	}
	if r.StatusCode == 404 {
		t.Fatalf("Application Resource %s not found", resourceID)
	}

	_, err = apiClient.ManagementAPIClient.ApplicationResourcesApi.DeleteApplicationResource(ctx, environmentID, applicationResource.Parent.GetId(), resourceID).Execute()
	if err != nil {
		t.Fatalf("Failed to delete resource (authorize): %v", err)
	}
}
