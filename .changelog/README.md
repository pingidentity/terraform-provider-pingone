# Changelog Management

This directory contains changelog entries for changes made to the repository. The changelog entries are used to generate release notes when a new version is released.

## How It Works

1. When you create a PR, you should create a changelog file named `pr-{PR_NUMBER}.txt` in this directory.
2. If you don't create one, GitHub Actions will automatically create one for you.
3. The changelog file should follow the format shown in the examples below.
4. During the release process, all changelog files are combined to generate release notes.
5. After release, all changelog files are moved to the `archive/{version}` directory.

## Changelog File Format

Changelog files use a simple key-value format. Each file should contain the following fields:

- `type:` - The type of change
- `description:` - A description of the change
- `pr:` - The pull request number
- `ticket:` - The Jira ticket (or N/A if not applicable)

The recognized change types are:

- `breaking` - Breaking changes that require user action
- `feature` - New features
- `enhancement` - Enhancements to existing functionality
- `bugfix` - Bug fixes
- `note` - General notes about the release
- `security` - Security-related changes
- `deprecation` - Deprecated features or functionality
- `internal` - Internal changes that are tracked but not included in release notes

## Creating a Changelog Entry

You can create a changelog entry using one of these two methods:

1. **Using the script**:

   ```bash
   ./scripts/create-changelog-entry.sh <PR-NUMBER> <CHANGE-TYPE> "Your change description" [JIRA-TICKET]
   ```

   Example:

   ```bash
   ./scripts/create-changelog-entry.sh 123 feature "Added new calculator function" CDI-456
   ```

2. **Manually creating the file**:
   Create a file named `pr-{PR_NUMBER}.txt` in the `.changelog` directory with the following format:

   ```plaintext
   type: <type>
   description: <sufficient description for portal or other use>
   pr: <pr #>
   ticket: CDI-####, PDI-####, or N/A
   ```

   Where:
   - `type` is one of: breaking, feature, enhancement, bugfix, note, security, deprecation, internal
   - `description` is a clear description of the change
   - `pr` is the pull request number
   - `ticket` is the Jira ticket (CDI-## or PDI-##) or N/A if not applicable

## Examples

```plaintext
type: feature
description: Added new calculator function `Multiply` for multiplication of two numbers
pr: 123
ticket: CDI-456
```

```plaintext
type: bugfix
description: Fixed issue where division by zero would cause application crash
pr: 124
ticket: N/A
```

## Release Process

The release process for managing changelog entries involves the following steps:

1. **Generate Release Notes**: Use the `generate-release-notes.sh` script to create release notes for the new version:

   ```bash
   ./shared-configs/release-notes/scripts/generate-release-notes.sh vX.Y.Z
   ```

   This scriptcreates:
   - `GITHUB_RELEASE_NOTES.md` in the repository root
   - `release-notes/vX.Y.Z/RELEASE_NOTES.adoc` for portal or other documentation
   - `release-notes/vX.Y.Z/GITHUB_RELEASE.md` for GitHub releases

2. **Review the Generated Files**: Verify that the generated files contain the correct information:
   - Confirm that all changes are properly categorized
   - PRs labeled as `internal` are excluded from the release notes
   - Check that the `release-notes/vX.Y.Z/` directory contains all expected files
   - Verify the content of `GITHUB_RELEASE_NOTES.md` (or `.adoc`)

3. **Archive Changelog Entries**: When releasing, archive the changelog entries using:

   ```bash
   ./scripts/archive-changelogs.sh vX.Y.Z
   ```

   This command moves all processed changelog files to the `.changelog/archive/vX.Y.Z/` directory.

This process ensures that each release contains only the changes made since the last release, and that all changelog entries are properly archived for historical reference.
