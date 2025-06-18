// Copyright © 2025 Ping Identity Corporation

//go:build beta

package main

import "github.com/pingidentity/terraform-provider-pingone/buildflags"

func init() {
	buildFlags = append(buildFlags, buildflags.Beta)
}
