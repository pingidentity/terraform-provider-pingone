// Copyright © 2025 Ping Identity Corporation

// This file relates to a beta feature described in CDI-492 and should be modified or removed on completion of CDI-631

//go:build beta

package beta

import "fmt"

var appImportFFSandboxEnvironmentName = "tf-testacc-static-app-import-ff-test"

func AppImportFFSandboxEnvironment() string {
	return fmt.Sprintf(`
		data "pingone_environment" "app_import_ff_test" {
			name = "%s"
		}`, appImportFFSandboxEnvironmentName)
}
