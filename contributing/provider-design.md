# Provider Design

## Directory/Module Structure
The following describes how the code is organised in the repository

* `internal/acctest` - The `acctest` directory contains functions that are common to all acceptance tests in the provider
* `internal/client` - The `client` directory contains functions specific to instantiation and modification of the PingOne GO SDK client.  The providers will source the PingOne client from this module
* `internal/provider` - The `provider` directory contains the core Terraform provider code.  This will initialise the Terraform provider as a whole
* `internal/sdk` - The `sdk` directory contains the code that brokers the interaction between the PingOne SDK functions and the Terraform provider.  For example, restructuring errors to Terraform provider format or processing resource retry conditions.
* `internal/service` - The `service` directory contains the resource and data source code relevant for all of the PingOne services.  Common code to all services and the PingOne platform will be created here.
* `internal/service/base` - The `base` directory, as a subdirectory of `service` contains the resources, data sources and testing code for the PingOne platform components that are shared between multiple services.
* `internal/service/<service name>` - The individual PingOne services are separated into their own directories, under the `service` directory.  This is to ensure full logical separation between the different services and help with ongoing maintainability.  Service names and their support status can be found on the [Services Support](services-support.md) guide

## PingOne GO SDK and API

The PingOne Terraform provider leverages the [PingOne Platform API](https://apidocs.pingidentity.com/pingone/platform/v1/api/), via an automatically generated [PingOne Go SDK](https://github.com/patrickcping/pingone-go-sdk-v2).  The resources in this provider must use the Go SDK to call PingOne platform endpoints, rather than call API endpoints directly.

For each function that requires an SDK call, the client must be retrieved and the PingOne domain suffix (of the PingOne tenant region) applied:

```
p1Client := meta.(*client.Client)
apiClient := p1Client.API.ManagementAPIClient
ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
	"suffix": p1Client.API.Region.URLSuffix,
})
```

Once initialised in the function (as above), the SDK can be invoked inside a response wrapper.  An example of how to retrieve an environment record is as follows:
```
	resp, diags := sdk.ParseResponse(
		ctx,
		func() (interface{}, *http.Response, error) {
			return apiClient.EnvironmentsApi.ReadOneEnvironment(ctx, environmentID).Execute()
		},
		"ReadOneEnvironment",
		sdk.CustomErrorResourceNotFoundWarning,
		sdk.DefaultRetryable,
	)
	if diags.HasError() {
		return diags
	}

	if resp == nil {
		d.SetId("")
		return nil
	}
	respObject := resp.(*management.Environment)
```

The purpose of the response wrapper is to standardise and optionally override API error responses and retry conditions.

More information about the SDK methods available can be found at:
* [Go API client for PingOne Management (SSO and Base)](https://pkg.go.dev/github.com/patrickcping/pingone-go-sdk-v2/management)
* [Go API client for PingOne MFA](https://pkg.go.dev/github.com/patrickcping/pingone-go-sdk-v2/mfa)
* [Go API client for PingOne Risk](https://pkg.go.dev/github.com/patrickcping/pingone-go-sdk-v2/risk)

## Custom Errors

The `CustomError` parameter of the `sdk.ParseResponse` function allows the developer to override API errors with a custom message or output.

The following shows an implementation where two overrides are in place; one that evaluates the PingOne API details block for specific validation errors based on the environment region, and the second overriding the error message returned based on a string match:
```
	resp, diags := sdk.ParseResponse(
		ctx,
		func() (interface{}, *http.Response, error) {
			return apiClient.EnvironmentsApi.CreateEnvironmentActiveLicense(ctx).Environment(environment).Execute()
		},
		"CreateEnvironmentActiveLicense",
		func(error management.P1Error) diag.Diagnostics {

			// Invalid region
			if details, ok := error.GetDetailsOk(); ok && details != nil && len(details) > 0 {
				if target, ok := details[0].GetTargetOk(); ok && *target == "region" {
					allowedRegions := make([]string, 0)
					for _, allowedRegion := range details[0].GetInnerError().AllowedValues {
						allowedRegions = append(allowedRegions, model.FindRegionByAPICode(management.EnumRegionCode(allowedRegion)).Region)
					}
					diags = diag.FromErr(fmt.Errorf("Incompatible environment region for the organization tenant.  Expecting regions %v, region provided: %s", allowedRegions, model.FindRegionByAPICode(region).Region))

					return diags
				}
			}

			// DV FF
			m, err := regexp.MatchString("^Organization does not have Ping One DaVinci FF enabled", error.GetMessage())
			if err != nil {
				diags = diag.FromErr(fmt.Errorf("Invalid regexp: DV FF error"))
				return diags
			}
			if m {
				diags = diag.FromErr(fmt.Errorf("The PingOne DaVinci service is not enabled in this organization tenant."))

				return diags
			}

			return nil
		},
		sdk.DefaultRetryable,
	)
```

The default value for this parameter is the `sdk.DefaultCustomError` function.  This can be explicitly set (recommended for readability), or the parameter value can be set to `nil`.

## Custom Retry Conditions

The `Retryable` parameter of the `sdk.ParseResponse` function allows the developer to define specific conditions of retry. By default, the provider will retry on network connectivity or timeout responses, but the developer may choose to also retry based on asynchronous latency.  More information about possible latency conditions can be found at the PingOne API documentation [Accounting for Latency](https://apidocs.pingidentity.com/pingone/platform/v1/api/#accounting-for-latency) section.

The following example shows a custom retry override to account for bootstrapped role assignment for the Terraform client to be able to create populations, based on string match of the error message:
```
	resp, diags := sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
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
	)
```

The default value for this parameter is the `sdk.DefaultRetryable` function.  This can be explicitly set (recommended for readability), or the parameter value can be set to `nil`.