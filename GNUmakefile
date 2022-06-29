TEST?=$$(go list ./...)
SWEEP_DIR=./internal/provider
NAMESPACE=pingidentity
PKG_NAME=pingone
BINARY=terraform-provider-${NAME}
VERSION=0.1.0
OS_ARCH=linux_amd64

default: build

build: fmtcheck
	go install -ldflags="-X github.com/pingidentity/terraform-provider-pingone/main.version=$(VERSION)"

test: fmtcheck
	go test $(TEST) $(TESTARGS) -timeout=5m

testacc: fmtcheck
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m

sweep:
	@echo "WARNING: This will destroy infrastructure. Use only in development accounts."
	go test $(SWEEP_DIR) -v -sweep="1" $(SWEEPARGS) -timeout 10m

vet:
	@echo "==> Running go vet ."
	@go vet ./... ; if [ $$? -ne 0 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

depscheck:
	@echo "==> Checking source code with go mod tidy..."
	@go mod tidy
	@git diff --exit-code -- go.mod go.sum || \
		(echo; echo "Unexpected difference in go.mod/go.sum files. Run 'go mod tidy' command or revert any go.mod/go.sum changes and commit."; exit 1)

lint: golangci-lint providerlint importlint

golangci-lint:
	@echo "==> Checking source code with golangci-lint..."
	@golangci-lint run ./$(PKG_NAME)/...

importlint:
	@echo "==> Checking source code with importlint..."
	@impi --local . --scheme stdThirdPartyLocal ./$(PKG_NAME)/...

providerlint:
	@echo "==> Checking source code with providerlint..."
	@providerlint \
		-c 1 \
		-AT001.ignored-filename-suffixes=_data_source_test.go \
		-AWSAT006=false \
		-AWSR002=false \
		-AWSV001=false \
		-R001=false \
		-R010=false \
		-R018=false \
		-R019=false \
		-V001=false \
		-V009=false \
		-V011=false \
		-V012=false \
		-V013=false \
		-V014=false \
		-XR001=false \
		-XR002=false \
		-XR003=false \
		-XR004=false \
		-XR005=false \
		-XS001=false \
		-XS002=false \
		./$(PKG_NAME)/provider/...

.PHONY: build test testacc sweep vet fmtcheck depscheck lint golangci-lint importlint providerlint