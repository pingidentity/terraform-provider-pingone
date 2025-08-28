# Beta PingOne terraform provider releases
Beta PingOne terraform provider versions allow for releasing of non-GA platform features. Beta releases of the provider will be distinguished by a `-beta` suffix at the end of the version tag, and will be released simultaneously with non-beta versions.

Functionality specific to the beta version of the provider is not included in the non-beta provider release binary. This is achieved using Go build tags. We use a `beta` build tag on all beta functionality, so that it is only compiled specifically for beta releases.

When adding a new beta resource or data source to the provider, the source file will need to include the `//go:build beta` build tag near the top of the file. The `make betatagscheck` target lints the provider to ensure that this tag is included in beta resources and data sources.

To add a new beta resource or data source to the provider, the corresponding NewXResource or NewXDataSource method should be added in the service_beta.go file for the given service. The service_beta_stub.go file should not be edited, as indicated in the comments.

Beta versions are released in the same manner as non-beta versions. Goreleaser will automatically run in a Github action when the tag is created. For example, a tag of `v0.0.1` will create a non-beta release, while `v0.0.1-beta` will create a beta release. Goreleaser is configured to set the `beta` build flag based on the tag suffix.

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
