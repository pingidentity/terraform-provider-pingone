# Beta PingOne terraform provider releases
Beta PingOne terraform provider versions allow for releasing of non-GA platform features. Beta releases of the provider will be distinguished by a `-beta` suffix at the end of the version tag, and will follow the same release cadence as non-beta releases.

Functionality specific to the beta version of the provider is not included in the non-beta provider release binary. This is achieved using Go build tags. We use a `beta` build tag on all beta functionality, so that it is only compiled specifically for beta releases.

When adding a new beta resource or data source to the provider, the source file will need to include the `//go:build beta` build tag near the top of the file. The `make betatagscheck` target lints the provider to ensure that this tag is included in beta resources and data sources.

To add a new beta resource or data source to the provider, the corresponding NewXResource or NewXDataSource method should be added in the `service_beta.go` file for the given service. The `service_beta_stub.go` file should not be edited, as indicated in the comments.

To add beta functionality to a non-beta resource or data source, you will need to stub out the corresponding functionality into separate files. The idea is to minimize future maintenance by leaving common functionality in the normal resource file, and stubbing out any custom logic (schema, building of API request bodies, etc.) into conditionally-compiled files. The beta logic will be in a file with the `//go:build beta` tag, and GA logic in a file with the `//go:build !beta` tag. Stubs are necessary so that the provider can compile with or without the beta implementation.

Beta versions are released via git tags, just like non-beta versions, but an extra step will be necessary to ensure the generated docs on the terraform registry include beta-specific functionality. The [registry queries the git tag for the docs](https://developer.hashicorp.com/terraform/registry/providers/docs), so the committed doc files for the given tag are what will show up on the registry. Because of this, the release workflow requires first branching for the beta release, and then generating doc files specific to the beta provider using `make generate BETA=true`, so that beta docs don't get committed to main. The branch name should follow the pattern `vX.X.X-beta-release`, where the release tag will be `vX.X.X-beta` (the branch name can't match the tag name because git refs must be unique). *Only once the beta doc files are in place should the beta tag be created from the new branch*.

Goreleaser will automatically run in a Github action when the beta tag is created. For example, a tag of `v0.0.1` will create a non-beta release, while `v0.0.1-beta` will create a beta release. Goreleaser is configured to set the `beta` build flag based on the tag suffix.

It may be useful when working in VSCode to tell `gopls` to compile with the beta build tag. This can be done by modifying `.vscode/settings.json` to include the build tag:

```
{
    "gopls": {
        "build.buildFlags": [
            // uncomment to build compile beta functionality in vscode
            //"-tags=beta"
        ]
    }
}
```
