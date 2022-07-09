# Development Environment Setup

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) 1.0+ (to run acceptance tests)
- [Go](https://golang.org/doc/install) 1.18+ (to build the provider plugin)

## Quick Start

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (please check the [requirements](#requirements) before proceeding).

*Note:* This project uses [Go Modules](https://blog.golang.org/using-go-modules) making it safe to work with it outside of your existing [GOPATH](http://golang.org/doc/code.html#GOPATH). The instructions that follow assume a directory in your home directory outside of the standard GOPATH (i.e `$HOME/development/pingidentity/`).

Clone repository to: `$HOME/development/pingidentity/`

```sh
$ mkdir -p $HOME/development/pingidentity/; cd $HOME/development/pingidentity/
$ git clone git@github.com:pingidentity/terraform-provider-pingone.git
...
```

Enter the provider directory and run `make tools`. This will install the needed tools for the provider.

```sh
$ make tools
```

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make build
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
