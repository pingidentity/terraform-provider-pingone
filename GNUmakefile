TEST_PATH?=$$(go list ./...)
SWEEP_DIR=./internal/sweep
NAMESPACE=pingidentity
PKG_NAME=pingone
BINARY=terraform-provider-${NAME}
VERSION=1.11.0
OS_ARCH=linux_amd64
BETA?=false

ifeq ($(BETA),true)
	BUILD_TAGS=-tags=beta
	VERSION_SUFFIX=-beta
else
	BUILD_TAGS=
	VERSION_SUFFIX=
endif

default: install

fmtcheck:
	@echo "==> Formatting Terraform documentation examples with terraform fmt..."
	@terraform fmt -recursive ./examples/

build:
	go mod tidy
	go build $(BUILD_TAGS) -v .

install: build
	go install $(BUILD_TAGS) -ldflags="-X main.version=$(VERSION)$(VERSION_SUFFIX)"

generate: build fmt
	GOFLAGS="$(BUILD_TAGS)" go tool tfplugindocs generate --provider-name terraform-provider-pingone

test: build
	go test $(TEST_PATH) $(TESTARGS) -timeout=5m

testacc: build
	TF_ACC=1 go test $(BUILD_TAGS) $$(go list ./internal/client/...) -v $(TESTARGS) -timeout 120m
	TF_ACC=1 go test $(BUILD_TAGS) $$(go list ./internal/service/...) -v $(TESTARGS) -timeout 120m

sweep: build
	@echo "WARNING: This will destroy infrastructure. Use only in development accounts."
	go test $(SWEEP_DIR) $(BUILD_TAGS) -v -sweep=all $(SWEEPARGS) -timeout 10m

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

lint: golangci-lint providerlint importlint tflint terrafmtcheck betatagscheck

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

betatagscheck:
	@echo "==> Checking beta resources and data sources for correct build tags..."
	@go run scripts/check_beta_build_tags.go

fmt: terrafmt fmtcheck

devcheck: build vet fmt generate docscategorycheck lint test sweep testacc

devchecknotest: build vet fmt generate docscategorycheck lint

.PHONY: build install generate docscategorycheck test testacc sweep vet fmtcheck depscheck lint golangci-lint importlint providerlint tflint terrafmt terrafmtcheck betatagscheck devcheck devchecknotest
