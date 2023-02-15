package framework

import pingone "github.com/pingidentity/terraform-provider-pingone/internal/client"

type ResourceType struct {
	Client *pingone.Client
}
