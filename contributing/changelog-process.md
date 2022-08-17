# Changelog Process

This project uses [go-changelog](https://github.com/hashicorp/go-changelog) to manage auto-generation of the [CHANGELOG.md](../CHANGELOG.md) file.

Using this library allows the maintainers to have a consistent format of the changelog file, that changes are not accidentally overwritten on PR merge, while ensuring all changes are accounted for.

The following guide shows how to include descriptions of your changes into the changelog generation process.

## Adding change descriptions

In order to add change descriptions to a PR, you add a new file into the `.changelog/` directory.  The name of the file will be `<PR/Issue number>.txt`

For example, if I raise a PR with ID `1234`, the file to create and include in the PR commit will be:

```
├── .changelog
│   ├── 1234.txt
```

The contents of the changelog text file will be one or more change descriptions, described below.

## Change descriptions

The format of the change descriptions follows the following format:
``````markdown
```release-note:<type>
<description text>
```
``````
Where `<type>` is the type of the change, and `<description text>` is a description of the change, depending on the change type.  For example:

``````markdown
```release-note:new-resource
pingone_risk_policy
```
``````

Changelog entry files can have multiple descriptions, for example `.changelog/1234.txt` might contain:
``````markdown
```release-note:new-resource
pingone_risk_policy
```

```release-note:new-data-source
pingone_risk_policy
```

```release-note:new-guide
How to configure a PingOne Risk policy
```
``````

The sections below describe the available change types and example descriptions.

### New Resource

A new resource has the `release-note:new-resource` header and the description will be just the name of the new resource.  Example:

``````markdown
```release-note:new-resource
pingone_risk_policy
```
``````

### New Data Source

A new data source has the `release-note:new-data-source` header and the description will be just the name of the new data source.  Example:

``````markdown
```release-note:new-data-source
pingone_risk_policy
```
``````

### Enhancement

An enhancement to the existing provider code (resources or data sources) has the `release-note:enhancement` header and the description will be the name of the resource and a follow on description.  Example:

``````markdown
```release-note:enhancement
resource/pingone_risk_policy: Add policy definition weights
```
``````

### Bug fix

A new documentation guide has the `release-note:new-guide` header and the description will highlight the fix that has been applied.  Example:

``````markdown
```release-note:bug
resource/pingone_risk_policy: Fix policy definition weights
```
``````

### New Guide

A new documentation guide has the `release-note:new-guide` header and the description will be the title of the guide that has been created.  Example:

``````markdown
```release-note:new-guide
How to configure a PingOne Risk policy
```
``````

### Notes / Deprecations / Dependencies

A general note or deprecation notice should follow the `release-note:note` header and the description relevant to the note.  Example:

``````markdown
```release-note:note
resource/pingone_risk_policy: Deprecate attribute EXAMPLE in favour of EXAMPLE
```
``````

### Breaking Change

A breaking change notice should follow the `release-note:breaking-change` header and the description will provide detail on the breaking change.  Example:

``````markdown
```release-note:breaking-change
resource/pingone_risk_policy: Attribute EXAMPLE has been removed from the resource
```
``````