// Copyright © 2026 Ping Identity Corporation

// Influenced from github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier/*

package setplanmodifier

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

type UnmodifiableDataLossProtectionIfFunc func(context.Context, planmodifier.SetRequest, *UnmodifiableDataLossProtectionIfFuncResponse)

type UnmodifiableDataLossProtectionIfFuncResponse struct {
	Diagnostics diag.Diagnostics
	Error       bool
}
