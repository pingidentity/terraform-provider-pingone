# Changelog Entry Format Examples

This directory contains changelog entries using markdown format with code blocks.

## Filename Format

The name of the changelog file will comprise of `pr-` followed by the PR number, appended with a `.txt` extension.

For example: `pr-2345.txt`

## Format

Each changelog entry file uses markdown code blocks with the `release-note:<type>` syntax. Multiple entries can be included in a single file, each in its own code block.

## Valid Types

- breaking-change
- feature
- enhancement
- bug
- note
- security
- deprecation
- new-resource (new Terraform resource)
- new-data-source (new Terraform data source)
- new-guide (new Terraform guide)
- internal (tracked but not included in release notes)

## Examples

### Feature example

````plaintext
```release-note:feature
Added new calculator function for advanced mathematical operations
```
````

### Bug fix example

````plaintext
```release-note:bug
Fixed issue where division by zero caused application crash
```
````

### Breaking change example

````plaintext
```release-note:breaking-change
Removed deprecated legacy API endpoints
```
````

### Enhancement with resource reference

````plaintext
```release-note:enhancement
`resource/pingone_risk_predictor`: Added the `predictor_device.should_validate_payload_signature` field to enforce requirement that the Signals SDK payload be provided as a signed JWT.
```
````

### Multiple entries in one file

````plaintext
```release-note:note
bump `github.com/example/sdk` 0.12.1 => 0.12.2
```

```release-note:bug
`resource/example_gateway`: Fixed error when configuring gateways.
```
````

### Internal change (not included in release notes)

````plaintext
```release-note:internal
Refactored internal helper functions for better code organization
```
````
