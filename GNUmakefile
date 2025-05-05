TEST_PATH?=$$(go list ./...)
SWEEP_DIR=./internal/sweep
NAMESPACE=pingidentity
PKG_NAME=pingone
BINARY=terraform-provider-${NAME}
VERSION=1.7.0
OS_ARCH=linux_amd64

default: install

tools:
	go generate -tags tools tools/main.go

fmtcheck:
	@echo "==> Formatting Terraform documentation examples with terraform fmt..."
	@terraform fmt -recursive ./examples/

build:
	go mod tidy
	go work vendor
	go build -v .

install: build
	go install -ldflags="-X main.version=$(VERSION)"

generate: build fmt
	go tool tfplugindocs generate

test: build
	go test $(TEST_PATH) $(TESTARGS) -timeout=5m

testacc: build
	TF_ACC=1 go test $$(go list ./internal/client/...) -v $(TESTARGS) -timeout 120m
	TF_ACC=1 go test $$(go list ./internal/service/...) -v $(TESTARGS) -timeout 120m

sweep: build
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

lint: golangci-lint providerlint importlint tflint terrafmtcheck

golangci-lint:
	@echo "==> Checking source code with golangci-lint..."
	@go tool golangci-lint run ./...

importlint:
	@echo "==> Checking source code with importlint..."
	@go tool impi --local . --scheme stdThirdPartyLocal ./...

providerlint:
	@echo "==> Checking source code with tfproviderlintx..."
	@go tool tfproviderlintx \
		-c 1 \
		-AT001.ignored-filename-suffixes=_data_source_test.go \
		-XR004=false \
		-XS002=false \
		./internal/provider/... ./internal/service/...

tflint:
	@echo "==> Checking Terraform code with tflint..."
	@go tool tflint --init

terrafmt:
	@echo "==> Formatting embedded Terraform code with terrafmt..."
	@find ./internal/service -type f -name '*_test.go' \
    | sort -u \
    | xargs -I {} go tool terrafmt -f fmt {}

terrafmtcheck:
	@echo "==> Checking embedded Terraform code with terrafmt..."
	@find ./internal/service -type f -name '*_test.go' \
    | sort -u \
    | xargs -I {} go tool terrafmt diff -f --check --fmtcompat {} ; if [ $$? -ne 0 ]; then \
		echo ""; \
		echo "terrafmt found bad formatting of HCL embedded in the test scripts. Please run "; \
		echo "\"make terrafmt\" before submitting the code for review."; \
		exit 1; \
	fi

fmt: terrafmt fmtcheck

devcheck: build vet tools fmt generate docscategorycheck lint test sweep testacc

.PHONY: tools build install generate docscategorycheck test testacc sweep vet fmtcheck depscheck lint golangci-lint importlint providerlint tflint terrafmt terrafmtcheck