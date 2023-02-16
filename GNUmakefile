TEST?=$$(go list ./...)
SWEEP_DIR=./internal/sweep
NAMESPACE=pingidentity
PKG_NAME=pingone
BINARY=terraform-provider-${NAME}
VERSION=0.9.0
OS_ARCH=linux_amd64

default: build

tools:
	go generate -tags tools tools/tools.go

build: fmtcheck
	go install -ldflags="-X github.com/pingidentity/terraform-provider-pingone/main.version=$(VERSION)"

generate: fmtcheck
	PINGONE_CLIENT_ID= PINGONE_CLIENT_SECRET= PINGONE_ENVIRONMENT_ID= PINGONE_API_ACCESS_TOKEN= PINGONE_REGION= PINGONE_ORGANIZATION_ID= go generate ./...

test: fmtcheck
	go test $(TEST) $(TESTARGS) -timeout=5m

testacc: fmtcheck
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m

sweep:
	@echo "WARNING: This will destroy infrastructure. Use only in development accounts."
	go test $(SWEEP_DIR) -v -sweep=all $(SWEEPARGS) -timeout 10m

vet:
	@echo "==> Running go vet..."
	@go vet ./... ; if [ $$? -ne 0 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

docscategorycheck:
	@echo "==> Checking for missing category in generated docs..."
	@find ./docs/**/*.md -print | xargs grep "subcategory: \"\""; if [ $$(find ./docs/**/*.md -print | xargs grep "subcategory: \"\"" | wc -l) -ne 0 ]; then \
		echo ""; \
		echo "Documentation check found a blank subcategory for the above files.  Ensure a template is created (./templates) with a subcategory set."; \
		exit 1; \
	fi

depscheck:
	@echo "==> Checking source code with go mod tidy..."
	@go mod tidy
	@git diff --exit-code -- go.mod go.sum || \
		(echo; echo "Unexpected difference in go.mod/go.sum files. Run 'go mod tidy' command or revert any go.mod/go.sum changes and commit."; exit 1)

lint: golangci-lint providerlint importlint tflint

golangci-lint:
	@echo "==> Checking source code with golangci-lint..."
	@golangci-lint run ./...

importlint:
	@echo "==> Checking source code with importlint..."
	@impi --local . --scheme stdThirdPartyLocal ./...

providerlint:
	@echo "==> Checking source code with tfproviderlintx..."
	@tfproviderlintx \
		-c 1 \
		-AT001.ignored-filename-suffixes=_data_source_test.go \
		-XR004=false \
		-XS002=false \
		./internal/provider/... ./internal/service/...

tflint:
	@echo "==> Checking Terraform code with tflint..."
	@tflint --init

terrafmt:
	@echo "==> Formatting embedded Terraform code with terrafmt..."
	@find ./internal/service -type f -name '*_test.go' \
    | sort -u \
    | xargs -I {} terrafmt -f fmt {}

terrafmtcheck:
	@echo "==> Checking embedded Terraform code with terrafmt..."
	@find ./internal/service -type f -name '*_test.go' \
    | sort -u \
    | xargs -I {} terrafmt diff -f --check --fmtcompat {} ; if [ $$? -ne 0 ]; then \
		echo ""; \
		echo "terrafmt found bad formatting of HCL embedded in the test scripts. Please run "; \
		echo "\"make terrafmt\" before submitting the code for review."; \
		exit 1; \
	fi

devcheck: build vet tools generate docscategorycheck lint terrafmtcheck test sweep testacc

.PHONY: tools build generate docscategorycheck test testacc sweep vet fmtcheck depscheck lint golangci-lint importlint providerlint tflint terrafmt terrafmtcheck
