---
layout: ""
page_title: "Provider Support"
description: |-
  The guide describes the support provided to customers using the PingOne Terraform provider.
---

# Provider Support

The guide describes the support provided to customers using the PingOne Terraform provider.

Distributions of the provider follows [Semantic Versioning 2.0.0](https://semver.org/) rules.

## "Generally Available" Distributions

Any distribution version that follows the `X.Y.Z` format where the major version number `X` is `1` or above, is considered **Generally Available** / GA and is **fully supported**.  Customers are advised to use these distributions to configure production environments.

With these distributions, customers are able to raise issues with Ping Identity and expect a timely resolution.

- Where issues are raised on the GitHub project repository or raised internally through a Ping contact, these are triaged by the engineering team to assess whether the resolution is with the Terraform project or a dependent API.  Terraform changes are scheduled according to priority and upstream changes to Ping APIs are forwarded internally.
- Where issues are raised to Ping Support, the usual documented SLA applies.

## Pre-release Distributions

Any distribution version that follows the `X.Y.Z-<pre-release>` format where there is a pre-release notation suffix is considered **not Generally Available** and support is limited.  Customers are advised to use these distributions at their own risk as the interfaces may not be stable and breaking changes may be introduced to both the pre-release binary and the underlying API without warning.

With pre-release versions, customers are still able to raise issues with Ping Identity but should not expect a timely resolution.

- Where issues are raised on the GitHub project repository or raised internally through a Ping contact, these are triaged by the engineering team to assess whether the resolution is with the Terraform project or a dependent API.  Terraform changes are scheduled according to priority and upstream changes to Ping APIs are forwarded internally.  Pre-release features carry a lower priority than "Generally Available" features.
- If using Ping Support, customers will be asked to reproduce the error on the latest "Generally Available" provider distribution.
