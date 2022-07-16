# Provider Design

## Directory/Module Structure
The following describes how the code is organised in the repository

* `internal/acctest` - The `acctest` directory contains functions that are common to all acceptance tests in the provider
* `internal/client` - The `client` directory contains functions specific to instantiation and modification of the PingOne GO SDK client.  The providers will source the PingOne client from this module
* `internal/provider` - The `provider` directory contains the core Terraform provider code.  This will initialise the Terraform provider as a whole
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

Once initialised in the function (as above), the SDK can be invoked.  An example of how to retrieve groups is as follows:
```
resp, r, err := apiClient.GroupsApi.ReadOneGroup(ctx, d.Get("environment_id").(string), d.Id()).Execute()
	if err != nil {

		if r.StatusCode == 404 {
			log.Printf("[INFO] PingOne Group %s no longer exists", d.Id())
			d.SetId("")
			return nil
		}
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error when calling `GroupsApi.ReadOneGroup``: %v", err),
			Detail:   fmt.Sprintf("Full HTTP response: %v\n", r.Body),
		})

		return diags
	}
```

More information can be found on the SDK documentation [PingOne Go SDK Readme](https://github.com/patrickcping/pingone-go-sdk-v2/blob/main/README.md)
