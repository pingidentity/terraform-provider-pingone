// Copyright © 2025 Ping Identity Corporation

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
	"github.com/pingidentity/terraform-provider-pingone/internal/provider/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/provider/frameworklegacysdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/provider/sdkv2"
)

func ProviderServerFactoryV6(ctx context.Context, version string) (func() tfprotov6.ProviderServer, error) {

	p1V5Provider := sdkv2.New(version)()
	p1V6ProviderLegacySdk := frameworklegacysdk.New(version)()
	p1V6Provider := framework.New(version)()

	upgradedp1V5Provider, err := tf5to6server.UpgradeServer(
		ctx,
		p1V5Provider.GRPCProvider,
	)

	if err != nil {
		return nil, err
	}

	providers := []func() tfprotov6.ProviderServer{
		func() tfprotov6.ProviderServer {
			return upgradedp1V5Provider
		},
		providerserver.NewProtocol6(p1V6ProviderLegacySdk),
		providerserver.NewProtocol6(p1V6Provider),
	}

	muxServer, err := tf6muxserver.NewMuxServer(ctx, providers...)
	if err != nil {
		return nil, err
	}

	return muxServer.ProviderServer, nil
}
