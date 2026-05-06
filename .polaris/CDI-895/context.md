# Context: CDI-895

## Source
- Type: jira
- Key: CDI-895
- Issue Type: Defect
- Project: Ping CICD Integrations

## Summary
P1 provider validation for pingone_application CORS origins is out of line with the API

## Description
Customer reported defect: https://pingidentity.slack.com/archives/C03SRQWCNMV/p1777016281540289

Looks like a very easy fix to bump up the validator to 40 for the size of the `oidc_options.cors_settings.origins` attribute.

Opened https://pingidentity.atlassian.net/browse/DOCS-11763 for the doc change.

## Acceptance Criteria
Bump max-size validator for `oidc_options.cors_settings.origins` to 40 (to match the API limit).

## Links
- DOCS-11763: doc change for updated validation limit
