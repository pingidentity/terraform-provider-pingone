// Copyright © 2025 Ping Identity Corporation

// Package provider provides the core provider factory and configuration for the PingOne Terraform provider.
// This package contains functions for creating provider servers that support both SDKv2 and Framework protocols.
package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
	"github.com/pingidentity/terraform-provider-pingone/internal/provider/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/provider/sdkv2"
)

// ProviderServerFactoryV6 creates a multiplexed provider server that supports both SDKv2 and Framework protocols.
// It returns a function that creates a tfprotov6.ProviderServer and any error encountered during setup.
// The version parameter is used to set the provider version string for both underlying providers.
// This factory enables the provider to serve both legacy SDKv2-based resources and new Framework-based resources
// within a single provider instance, allowing for gradual migration from SDKv2 to Framework.
func ProviderServerFactoryV6(ctx context.Context, version string) (func() tfprotov6.ProviderServer, error) {

	p1V5Provider := sdkv2.New(version)()
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
		providerserver.NewProtocol6(p1V6Provider),
	}

	muxServer, err := tf6muxserver.NewMuxServer(ctx, providers...)
	if err != nil {
		return nil, err
	}

	return muxServer.ProviderServer, nil
}
