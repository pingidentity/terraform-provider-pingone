package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pingidentity/terraform-provider-pingone/internal/provider/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/provider/sdkv2"
)

func ProviderServerFactoryV5(ctx context.Context, version string) (func() tfprotov5.ProviderServer, *schema.Provider, error) {

	p1V5Provider := sdkv2.New(version)()
	p1V6Provider := framework.New(version)()

	providers := []func() tfprotov5.ProviderServer{
		p1V5Provider.GRPCProvider,
		providerserver.NewProtocol5(p1V6Provider),
	}

	muxServer, err := tf5muxserver.NewMuxServer(ctx, providers...)
	if err != nil {
		return nil, nil, err
	}

	return muxServer.ProviderServer, p1V5Provider, nil
}
