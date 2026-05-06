# Task 2 Summary: Update hardcoded origins description in data_source_application.go

## What Was Implemented

Replaced the hardcoded string `"Limited to 20 values."` with `"Limited to 40 values."` in the `originsDescription` string literal at line 776 of `data_source_application.go`. This is a single-character-group change (20 → 40) within the `datasourceApplicationSchemaCorsSettings()` function. No other lines in the file were changed.

## Files Modified

- `internal/service/sso/data_source_application.go` — Changed `Limited to 20 values.` to `Limited to 40 values.` at line 776 within the `originsDescription` string literal

## Verification

- Tests: not run individually (build verified; no test asserts the old limit of 20 per requirements)
- Lint: skipped (no lint command available without a full `make` toolchain invocation)
- Typecheck: passed implicitly via successful build
- Build: passed — `go build ./...` completed with no output and exit code 0
