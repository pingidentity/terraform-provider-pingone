# Beta PingOne Terraform provider contribution guide

This guide outlines the process for developing, managing, and releasing beta versions of the PingOne Terraform provider.

Beta releases allow us to support non-GA platform features in the terraform provider. They are identified by a `-beta` suffix in the version tag (e.g., `v1.2.3-beta`) and follow the same release cadence as standard releases.

---

## How It Works: Go Build Tags

The core mechanism for separating beta and standard (GA) functionality is **Go build tags**.

We use a `beta` build tag on all files containing beta-specific code. This ensures the beta code is only compiled into the provider binary when explicitly requested.

* **Beta files** must include the build tag comment at the top:
    `//go:build beta`
* **GA-only files** (often used as stubs for beta functionality) must include the inverse tag:
    `//go:build !beta`

The build system, driven by Goreleaser, automatically adds the `-tags=beta` flag when it detects a Git tag ending in `-beta`, ensuring the correct version is built.

---

## Development Workflow

### Adding a New Beta Resource or Data Source

Use this process for resources and data sources that are entirely new and have no existing GA counterpart.

1.  **Create your source file** for the new resource or data source.
2.  **Add the build tag** at the very top of the file.
    ```go
    //go:build beta
    ```
3.  **Register the new resource/data source** by adding its `NewXResource` or `NewXDataSource` method to the `service_beta.go` file for the relevant service.
    * **Do not** edit the `service_beta_stub.go` file.
4.  **Add any required documentation templates** to the `betatemplates` folder. Doc templates for beta resources or data sources must be placed in a separate directory, or tfplugindocs will fail due to being unable to find the resource.

You can validate that your beta tags are correctly placed by running `make betatagscheck`.

### Adding Beta Functionality to an Existing Resource

Use this more complex process when you need to add new beta-only attributes or logic to a resource that already exists in the GA provider. The goal is to isolate only the beta changes, leaving the common code untouched.

This requires splitting the logic into multiple files:

1.  The existing resource file will contain all the common, shared logic (e.g., `Read`, `Delete`, common schema attributes) that applies to both GA and beta versions. This file has **no build tags**.
2.  One or more beta source files will contain the beta-specific logic, such as beta-only schema attributes or custom logic for `Create` and `Update`. These files **must** have the `//go:build beta` tag.
3.  One or more GA source files will contain the GA-specific logic or, more commonly, empty stubs for the functions defined in the `_beta.go` file. This ensures the provider can compile without the beta code. These files **must** have the `//go:build !beta` tag.

---

## Beta Release Process

Publishing a beta release requires a specific sequence of steps to ensure the Terraform Registry displays the correct documentation.

**Important:** The documentation must be generated and committed to a dedicated branch **before** the tag is created.

1.  **Create a Release Branch:** Create a new branch from `main`. The branch name must be unique and should follow the pattern `vX.X.X-beta-release`.
    * *Note: Git refs must be unique, so the branch name cannot be identical to the tag name you will create later. This is why a `-release` suffix is included.*

2.  **Generate Beta Documentation:** 
    On your new branch, copy any beta resource or data source documentation templates from `betatemplates` to the `templates` folder, so that they are handled by `tfplugindocs`.

    Then run the following command to generate documentation that includes the beta-specific resources and attributes:
    ```shell
    make generate BETA=true
    ```
   

3.  **Commit the Documentation:** Commit the newly generated doc files to your release branch.

4.  **Create the Release:** Once the beta documentation is committed and pushed, create a corresponding GitHub release. The tag name **must** match the version number and end with `-beta`, and the release should be marked as pre-release, **not** the latest release. Release notes are not added for beta functionality. You can use the "What's Changed" button to fill out the content of the GitHub release. The main release notes will be included on the non-beta release.

Creating the GitHub release will automatically trigger a GitHub Action that runs Goreleaser, builds the beta provider binaries, and attaches them to the release. The Terraform registry will eventually pick up the new release and make it available to download.

---

## VSCode Configuration Tip

To compile files with a `beta` tag in VSCode, you can instruct the Go language server (`gopls`) to always use the `beta` build tag.

Modify your `.vscode/settings.json` file with the following configuration:

```json
{
    "gopls": {
        "build.buildFlags": [
            // uncomment to build compile beta functionality in vscode
            //"-tags=beta"
        ]
    }
}
```
