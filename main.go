// Copyright © 2025 Ping Identity Corporation

package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6/tf6server"
	"github.com/pingidentity/terraform-provider-pingone/buildflags"
	"github.com/pingidentity/terraform-provider-pingone/internal/provider"
)

var (
	// these will be set by the goreleaser configuration
	// to appropriate values for the compiled binary
	version string = "dev"

	// goreleaser can also pass the specific commit if you want
	// commit  string = ""

	buildFlags []buildflags.BuildFlag
)

func main() {
	debugFlag := flag.Bool("debug", false, "Start provider in debug mode.")
	flag.Parse()

	ctx := context.Background()

	muxServer, err := provider.ProviderServerFactoryV6(ctx, version, buildFlags)
	if err != nil {
		log.Fatal(err)
	}

	var serveOpts []tf6server.ServeOpt

	if *debugFlag {
		serveOpts = append(serveOpts, tf6server.WithManagedDebug())
	}

	err = tf6server.Serve(
		"registry.terraform.io/pingidentity/pingone",
		muxServer,
		serveOpts...,
	)

	if err != nil {
		log.Fatal(err)
	}
}
