# Discussion Output

## Understanding
The `pingone_application` resource enforces a Terraform-side maximum of 20 entries for `oidc_options.cors_settings.origins`, but the PingOne API supports up to 40. This mismatch prevents customers from configuring as many CORS origins as the API actually allows. The fix is to raise the `originsMax` constant from 20 to 40 in the resource schema.

## Jira Quality Assessment
Jira-based entry. The ticket is clear and well-scoped: a specific attribute, a specific new limit (40), and no ambiguity about intent. No formal ACs section in the ticket, but the description is self-contained. No reproduction steps provided (not strictly needed for a validator change — the symptom is obvious: plan/apply fails for any config with >20 origins).

## Locked Decisions
1. Raise `originsMax` constant from 20 to 40 in `internal/service/sso/resource_application.go`: This is the single source of truth for the validator and the description string — changing the constant updates both automatically.

## Constraints
1. Change is limited to the validator value only — no API behavior, no data model, no other attributes are affected.
2. The description string is templated off the `originsMax` constant, so it will update automatically — no separate string change needed.

## Research Directives
1. Confirm whether the same `originsMax` constant or an equivalent limit appears in the schema upgrade file (`resource_application_schema_upgrade_0_to_1.go`), the data source (`data_source_application.go`), or utils (`utils_application.go`) — and whether any of those also need updating.
2. Confirm whether any tests assert the old limit of 20 and need to be updated to 40.

## Confidence

| Dimension | Level | Rationale |
|-----------|-------|-----------|
| Requirements sufficiency | High | Single-line fix with a clearly stated target value (40). No ambiguity about scope or intent. |
