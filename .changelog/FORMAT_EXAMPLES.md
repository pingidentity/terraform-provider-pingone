# Changelog Entry Format Examples

This directory contains changelog entries using the structured key-value format.

## Filename Format

The name of the changelog file will comprise of `pr-` followed by the PR number, appended with a `.txt` extension.

For example: `pr-2345.txt`

## Format

Each changelog entry file should follow this format:

```plaintext
type: <type>
description: <sufficient description for portal or other use>
pr: <pr #>
ticket: <Jira>
```

## Valid Types

- breaking
- feature
- enhancement
- bugfix
- note
- security
- deprecation
- internal (tracked but not included in release notes)

## Field Descriptions

- **type**: The type of change (required)
- **description**: A clear description of the change (required)
- **pr**: The pull request number (required, or can be inferred from filename)
- **ticket**: The Jira ticket (CDI-### or PDI-###), or N/A if not applicable (required)

## Examples

### Feature with Jira ticket

```plaintext
type: feature
description: Added new calculator function for advanced mathematical operations
pr: 123
ticket: CDI-456
```

### Bug fix without Jira ticket

```plaintext
type: bugfix
description: Fixed issue where division by zero caused application crash
pr: 124
ticket: N/A
```

### Breaking change with Jira ticket

```plaintext
type: breaking
description: Removed deprecated legacy API endpoints
pr: 125
ticket: PDI-789
```

### Internal change (not included in release notes)

```plaintext
type: internal
description: Refactored internal helper functions for better code organization
pr: 126
ticket: N/A
```
