# Provider Design

## Directory/Module Structure
The following describes how the code is organised in the repository

* `internal/acctest` - The `acctest` directory contains functions that are common to all acceptance tests in the provider
* `internal/client` - The `client` directory contains functions specific to instantiation and modification of the PingOne GO SDK client.  The providers will source the PingOne client from this module
* `internal/framework` - The `sdk` directory contains the code that brokers the interaction between the PingOne SDK functions and the Terraform provider (using the v6 protocol/plugin framework SDK).  For example, restructuring errors to Terraform provider format or processing resource retry conditions.
* `internal/provider` - The `provider` directory contains the core Terraform provider code.  This will initialise the Terraform provider as a whole.  The `provider` package contains a provider factory based on the mux packages described [here](https://developer.hashicorp.com/terraform/plugin/mux).  The provider is being gradually migrated to the v6 protocol/plugin framework SDK and away from the v5 protocol/SDKv2 SDK.  Currently, the factory presents a combined v5 protocol provider for backward compatibility of Terraform CLI earlier than v1.
* `internal/sdk` - The `sdk` directory contains the code that brokers the interaction between the PingOne SDK functions and the Terraform provider (using the v5 protocl/SDKv2 SDK).  For example, restructuring errors to Terraform provider format or processing resource retry conditions.
* `internal/service` - The `service` directory contains the resource and data source code relevant for all of the PingOne services.  Common code to all services and the PingOne platform will be created here.
* `internal/service/base` - The `base` directory, as a subdirectory of `service` contains the resources, data sources and testing code for the PingOne platform components that are shared between multiple services.
* `internal/service/<service name>` - The individual PingOne services are separated into their own directories, under the `service` directory.  This is to ensure full logical separation between the different services and help with ongoing maintainability.  Service names and their support status can be found on the [Services Support](services-support.md) guide

## Migration to Terraform Plugin Framework SDK

The provider is being gradually migrated to the Plugin Framework SDK.  Any new resources and/or data sources should use the Plugin Framework and be registered in the `internal/service/<service name>/service.go` file for inclusion in the provider.  Examples provided in this document are specific to the Plugin Framework SDK.

## PingOne GO SDK and API

The PingOne Terraform provider leverages the [PingOne Platform API](https://apidocs.pingidentity.com/pingone/platform/v1/api/), via an automatically generated [PingOne Go SDK](https://github.com/patrickcping/pingone-go-sdk-v2).  The resources in this provider must use the Go SDK to call PingOne platform endpoints, rather than call API endpoints directly.

For each resource that requires a PingOne SDK call, the client must be retrieved and the PingOne domain suffix (of the PingOne tenant region) applied:

```
func (r *FooResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	resourceConfig, ok := req.ProviderData.(framework.ResourceType)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected the provider client, got: %T. Please report this issue to the provider maintainers.", req.ProviderData),
		)

		return
	}

	preparedClient, err := prepareClient(ctx, resourceConfig)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			err.Error(),
		)

		return
	}

	r.client = preparedClient
}
```

Once initialised in the Configure method (as above), the SDK can be invoked inside a function using a retry and response parsing wrapper.  Example:
```
	var trustedEmailAddress *management.EmailDomainTrustedEmail
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			return r.client.TrustedEmailAddressesApi.CreateTrustedEmailAddress(ctx, plan.EnvironmentId.ValueString(), plan.EmailDomainId.ValueString()).EmailDomainTrustedEmail(*emailDomainTrustedEmail).Execute()
		},
		"CreateTrustedEmailAddress", // This is an ID used for logging and error output
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
		&trustedEmailAddress,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}
```

The purpose of the response wrapper is to standardise and optionally override API error responses and retry conditions.

More information about the SDK methods available can be found at:
* [Go API client for PingOne Management (SSO and Base)](https://pkg.go.dev/github.com/patrickcping/pingone-go-sdk-v2/management)
* [Go API client for PingOne Authorize](https://pkg.go.dev/github.com/patrickcping/pingone-go-sdk-v2/authorize)
* [Go API client for PingOne MFA](https://pkg.go.dev/github.com/patrickcping/pingone-go-sdk-v2/mfa)
* [Go API client for PingOne Risk](https://pkg.go.dev/github.com/patrickcping/pingone-go-sdk-v2/risk)
* [Go API client for PingOne Credentials](https://pkg.go.dev/github.com/patrickcping/pingone-go-sdk-v2/credentials)
* [Go API client for PingOne Verify](https://pkg.go.dev/github.com/patrickcping/pingone-go-sdk-v2/verify)

## Custom Errors

The `CustomError` parameter of the `framework.ParseResponse` function allows the developer to override API errors with a custom message or output.

The following shows an implementation where two overrides are in place; one that evaluates the PingOne API details block for specific validation errors based on the environment region, and the second overriding the error message returned based on a string match:
```
	var trustedEmailAddress *management.EmailDomainTrustedEmail
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			return r.client.TrustedEmailAddressesApi.CreateTrustedEmailAddress(ctx, plan.EnvironmentId.ValueString(), plan.EmailDomainId.ValueString()).EmailDomainTrustedEmail(*emailDomainTrustedEmail).Execute()
		},
		"CreateTrustedEmailAddress", // This is an ID used for logging and error output
		trustedEmailAddressAPIErrors, // This is an overridden error function
		sdk.DefaultCreateReadRetryable,
		&trustedEmailAddress,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}
```

```
func trustedEmailAddressAPIErrors(error model.P1Error) diag.Diagnostics {
	var diags diag.Diagnostics

	// Domain not verified
	if details, ok := error.GetDetailsOk(); ok && details != nil && len(details) > 0 {
		if code, ok := details[0].GetCodeOk(); ok && *code == "INVALID_VALUE" {
			if target, ok := details[0].GetTargetOk(); ok && *target == "trustedEmail" {
				diags.AddError(
					"The domain of the given email address is not verified",
					"Ensure that the domain of the given trusted email address has been verified first.  This can be configured with the `pingone_trusted_email_domain` resource.",
				)

				return diags
			}
		}
	}
	return nil
}
```

The default value for this parameter is the `framework.DefaultCustomError` function.  This can be explicitly set (recommended for readability), or the parameter value can be set to `nil`.

## Custom Retry Conditions

The `Retryable` parameter of the `framework.ParseResponse` function allows the developer to define specific conditions of retry. By default, the provider will retry on network connectivity or timeout responses, but the developer may choose to also retry based on asynchronous latency and eventual consistency.  More information about possible latency conditions can be found at the PingOne API documentation [Accounting for Latency](https://apidocs.pingidentity.com/pingone/platform/v1/api/#accounting-for-latency) section.

The following example shows a custom retry override to account for bootstrapped role assignment for the Terraform client to be able to create populations, based on string match of the error message:
```
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			return apiClient.PopulationsApi.CreatePopulation(ctx, environmentID).Population(population).Execute()
		},
		"CreatePopulation",
		sdk.DefaultCustomError,
		func(ctx context.Context, r *http.Response, p1error *management.P1Error) bool {

			if p1error != nil {
				var err error

				// Permissions may not have propagated by this point
				if m, err := regexp.MatchString("^The actor attempting to perform the request is not authorized.", p1error.GetMessage()); err == nil && m {
					tflog.Warn(ctx, "Insufficient PingOne privileges detected")
					return true
				}
				if err != nil {
					tflog.Warn(ctx, "Cannot match error string for retry")
					return false
				}

			}

			return false
		},
		&population,
	)...)
```

The default value for this parameter is the `sdk.DefaultRetryable` function.  This can be explicitly set (recommended for readability), or the parameter value can be set to `nil`.

*Note:* Retry specific code has not been converted to a framework equivalent and still exists under the `internal/sdk` package.  This is because the plugin framework does not yet include features for retry.