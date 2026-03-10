## Key Local Development Commands

- `make build` — Build the provider (`go mod tidy` + `go build`)
- `make install` — Build and install provider binary to `$GOPATH/bin`
- `make generate` — Generate documentation
- `make fmt` — Format embedded Terraform code and examples
- `make lint` — Run all linters
- `make vet` — Run `go vet`
- `make devchecknotest` — Pre-PR check: build, vet, fmt, generate, docscategorycheck, lint

### Beta Builds

Pass `BETA=true` to enable beta resources/data sources:

```sh
make install BETA=true
```

Beta resources use Go build tags (`//go:build beta`).

## Development Standards

- Follow typical go coding standards, check merge-readiness with `make devchecknotest`
- Resources use Terraform Plugin Framework (preferred, use for any new endpoints) or SDK v2 (legacy)
- Follow existing patterns in the same service directory for consistency

### Schema Standards

- All attributes must have a description
- Attribute defaults must be defined in schema via plugin-framework methods, and should be included in descriptions
- Use string formatting to inject values into attribute descriptions

### Testing

- All resources/data sources require acceptance tests, which run terraform against a live PingOne environment
- Typical test patterns per resource:
  - **RemovalDrift** — Verify detection when resource is deleted outside Terraform
  - **NewEnv** — Test creation in a fresh environment
  - **Full** — Full lifecycle: minimal → maximal → minimal config, plus import
  - **BadParameters** - Test invalid import parameters
- Example test commands
  - Non-beta: `TF_ACC=1 go test -v -timeout 300s -run ^TestAccPopulation_ github.com/pingidentity/terraform-provider-pingone/internal/service/sso`
  - Beta: `TF_ACC=1 TESTACC_BETA=true go test -tags=beta -v -timeout 300s -run ^TestAccDavinciFlowDeploy github.com/pingidentity/terraform-provider-pingone/internal/service/davinci`
