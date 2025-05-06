# Development Environment Setup

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) 1.4+ (to run acceptance tests)
- [Go](https://golang.org/doc/install) 1.23.3+ (to build and test the provider plugin)

## Quick Start

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (please check the [requirements](#requirements) before proceeding).

*Note:* This project uses [Go Modules](https://blog.golang.org/using-go-modules) making it safe to work with it outside of your existing [GOPATH](http://golang.org/doc/code.html#GOPATH). The instructions that follow assume a directory in your home directory outside of the standard GOPATH (i.e `$HOME/development/pingidentity/`).

Clone repository to: `$HOME/development/pingidentity/`

```sh
$ mkdir -p $HOME/development/pingidentity/; cd $HOME/development/pingidentity/
$ git clone git@github.com:pingidentity/terraform-provider-pingone.git
...
```

To compile the provider, run `make build`.

```sh
$ make build
```

To install the provider for local use, run `make install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make install
...
$ $GOPATH/bin/terraform-provider-pingone
...
```

## Testing the Provider

In order to test the provider locally with no connection to PingOne, you can run `make test`.

```sh
$ make test
```

In order to run the full suite of Acceptance tests against a live PingOne tenant, run `make testacc`.

*Note:* Acceptance tests create real configuration in PingOne.  Please ensure you have a trial PingOne account or licensed subscription to proceed with these tests.

```sh
$ make testacc
```

## Using the Provider

With Terraform v0.14 and later, [development overrides for provider developers](https://www.terraform.io/docs/cli/config/config-file.html#development-overrides-for-provider-developers) can be leveraged in order to use the provider built from source.

To do this, populate a Terraform CLI configuration file (`~/.terraformrc` for all platforms other than Windows; `terraform.rc` in the `%APPDATA%` directory when using Windows) with at least the following options:

```hcl
provider_installation {
  dev_overrides {
    "pingidentity/pingone" = "[REPLACE WITH GOPATH]/bin"
  }
  direct {}
}
```

## Local SDK Changes

Occasionally, development may include changes to the [PingOne GO SDK](https://github.com/patrickcping/pingone-go-sdk-v2). If you'd like to develop this provider locally using a local, modified version of the SDK, this can be achieved by adding a `replace` directive in the `go.mod` file.  For example, the start of the `go.mod` file may look like the following example, where the local cloned SDK is in the `../pingone-go` relative path, and we substitute the `management` module:

```
module github.com/pingidentity/terraform-provider-pingone

go 1.23.3

replace github.com/patrickcping/pingone-go-sdk-v2/management => ../pingone-go-sdk-v2/management

require (
	github.com/patrickcping/pingone-go-sdk-v2/management v0.x.x
  
  ...
)

...
```

Once updated, run the following to build the project:

```shell
$ make build
```