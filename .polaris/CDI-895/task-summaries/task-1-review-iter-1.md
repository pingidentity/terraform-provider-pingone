# Task 1 Review — Iter 1

## Verdict
APPROVE

## Summary
Reviewed commit `86684e23` ("Raise CORS origins limit from 20 to 40 in resource_application.go"), which is the sole code-bearing commit for Task 1. The change is a single-line constant edit at line 1981 of `internal/service/sso/resource_application.go`, replacing `const originsMax = 20` with `const originsMax = 40`. All acceptance criteria are fully met. No unintended file changes are present and the build passes.

## Verification
- Tests: NOT_RUN — no test asserts the old limit of 20 per requirements.md; the engineer's note is accurate.
- Lint: NOT_RUN — no lint command available without the full make toolchain. No new code was introduced that would plausibly trigger a lint failure.
- Typecheck: PASS — implicit from build success. The constant remains the same `int` type; no type-system impact.
- Build: PASS — `go build ./...` on branch `polaris/CDI-895-fix-cors-origins-validator-size` exits with code 0 and no output.

## Requirements Check
- [x] AC 1 (`const originsMax` set to `40`): Met. Line 1981 of `resource_application.go` reads `const originsMax = 40`. Verified with `git show` and `sed -n '1981p'`.
- [x] AC 2 (`setvalidator.SizeAtMost` enforces 40): Met. Line 2010 contains `setvalidator.SizeAtMost(originsMax)` and `originsMax` is now 40. No explicit second edit is required; the constant reference propagates automatically.
- [x] AC 3 (description string reads "Limited to 40 values."): Met. Line 1983 uses `fmt.Sprintf("... Limited to %d values.", ..., originsMax)`. With `originsMax = 40` the rendered string will be "Limited to 40 values." Confirmed no other literal "20" remains in the CORS settings function.
- [x] AC 4 (no other lines in `resource_application.go` changed): Met. `git diff main..86684e23 -- internal/service/sso/resource_application.go` shows exactly one line changed (line 1981). No other hunks.
- [x] AC 5 (`go build ./...` succeeds): Met. Verified above.
- [x] Task scope constraint (data source, schema upgrade, utils, tests, docs untouched): Met. `git diff --name-only` for the task commit shows only `internal/service/sso/resource_application.go` and the `.polaris/` summary file were added or modified. No docs, no data source, no schema upgrade file, no test files.

## Findings

### Critical
None.

### Important
None.

### Suggestions
- **internal/service/sso/data_source_application.go:776**: The `data_source_application.go` description still reads "Limited to 20 values." as expected at this stage (that is Task 2's responsibility), but reviewers of the branch as a whole should confirm the tasks are committed and merged together before the branch is squash-merged to main, since the branch is in an intermediate inconsistent state where the resource and data source descriptions disagree. This is not a bug in this task; it is a process observation for the merge sequence.
