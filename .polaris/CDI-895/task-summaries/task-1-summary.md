# Task 1 Summary: Raise originsMax constant from 20 to 40 in resource_application.go

## What Was Implemented

Changed the function-local constant `originsMax` from `20` to `40` inside `resourceApplicationSchemaCorsSettings()` in `resource_application.go`. This single-line edit propagates automatically to both:
- The `fmt.Sprintf`-based description string on line 1983, which now produces "Limited to 40 values."
- The `setvalidator.SizeAtMost(originsMax)` validator on line 2010, which now enforces a maximum of 40.

Only line 1981 was changed. No other lines in `resource_application.go` were touched.

## Files Modified

- `internal/service/sso/resource_application.go` — Changed `const originsMax = 20` to `const originsMax = 40` at line 1981 inside `resourceApplicationSchemaCorsSettings()`

## Verification

- Tests: not run individually (build verified; no test asserts the old limit of 20 per requirements)
- Lint: skipped (no lint command available without a full `make` toolchain invocation)
- Typecheck: passed implicitly via successful build
- Build: passed — `go build ./...` completed with no output and exit code 0

## Deviations

None.
