TEST?=$$(go list ./... | grep -v '/vendor/' | grep -v '/builtin/bins/')
VETARGS?=-all
GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)

default: test vet

tools:
	go get -u github.com/kardianos/govendor
	go get -u golang.org/x/tools/cmd/stringer
	go get -u golang.org/x/tools/cmd/cover

# bin generates the releaseable binaries for Terraform
bin: fmtcheck generate
	@TF_RELEASE=1 sh -c "'$(CURDIR)/scripts/build.sh'"

# dev creates binaries for testing Terraform locally. These are put
# into ./bin/ as well as $GOPATH/bin
dev: fmtcheck generate
	@TF_DEV=1 sh -c "'$(CURDIR)/scripts/build.sh'"

quickdev: generate
	@TF_DEV=1 sh -c "'$(CURDIR)/scripts/build.sh'"

# Shorthand for quickly building the core of Terraform. Note that some
# changes will require a rebuild of everything, in which case the dev
# target should be used.
core-dev: generate
	go install -tags 'core' github.com/hashicorp/terraform

# Shorthand for quickly testing the core of Terraform (i.e. "not providers")
core-test: generate
	@echo "Testing core packages..." && go test -tags 'core' $(shell go list ./... | grep -v -E 'builtin|vendor')

# Shorthand for building and installing just one plugin for local testing.
# Run as (for example): make plugin-dev PLUGIN=provider-aws
plugin-dev: generate
	go install github.com/hashicorp/terraform/builtin/bins/$(PLUGIN)
	mv $(GOPATH)/bin/$(PLUGIN) $(GOPATH)/bin/terraform-$(PLUGIN)

# test runs the unit tests
test: fmtcheck generate
	TF_ACC= go test $(TEST) $(TESTARGS) -timeout=30s -parallel=4

# testacc runs acceptance tests
testacc: fmtcheck generate
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make testacc TEST=./builtin/providers/aws"; \
		exit 1; \
	fi
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m

# testrace runs the race checker
testrace: fmtcheck generate
	TF_ACC= go test -race $(TEST) $(TESTARGS)

cover:
	@go tool cover 2>/dev/null; if [ $$? -eq 3 ]; then \
		go get -u golang.org/x/tools/cmd/cover; \
	fi
	go test $(TEST) -coverprofile=coverage.out
	go tool cover -html=coverage.out
	rm coverage.out

# vet runs the Go source code static analysis tool `vet` to find
# any common errors.
vet:
	@echo "go tool vet $(VETARGS) ."
	@go tool vet $(VETARGS) $$(ls -d */ | grep -v vendor) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

# generate runs `go generate` to build the dynamically generated
# source files.
generate:
	@which stringer ; if [ $$? -ne 0 ]; then \
	  go get -u golang.org/x/tools/cmd/stringer; \
	fi
	go generate $$(go list ./... | grep -v /vendor/)
	@go fmt command/internal_plugin_list.go > /dev/null

fmt:
	gofmt -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

.PHONY: bin default generate test vet fmt fmtcheck tools
