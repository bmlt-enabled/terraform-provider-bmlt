# Terraform Provider for BMLT

default: build

GIT_COMMIT?=$(shell git rev-parse HEAD)
GIT_DIRTY?=$(shell test -n "`git status --porcelain`" && echo "+CHANGES" || true)
GIT_DESCRIBE?=$(shell git describe --tags --always)
GIT_IMPORT=github.com/bmlt-enabled/terraform-provider-bmlt/version
LDFLAGS=-X $(GIT_IMPORT).GitCommit=$(GIT_COMMIT)$(GIT_DIRTY) -X $(GIT_IMPORT).GitDescribe=$(GIT_DESCRIBE)

GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)

build: fmtcheck
	go install -ldflags "$(LDFLAGS)"

test: fmtcheck
	go test -i $(TEST) || exit 1
	echo $(TEST) | \
		xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

testacc: fmtcheck
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m

vet:
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

fmt:
	gofmt -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

lint:
	@echo "==> Checking source code against linters..."
	@golangci-lint run ./...

tools:
	@echo "==> installing required tooling..."
	@sh -c "'$(CURDIR)/scripts/install-tools.sh'"

docs:
	go generate ./...

install: build
	mkdir -p ~/.terraform.d/plugins/bmlt-enabled.org/local/bmlt/1.0.0/darwin_amd64
	cp $${GOPATH}/bin/terraform-provider-bmlt ~/.terraform.d/plugins/bmlt-enabled.org/local/bmlt/1.0.0/darwin_amd64

clean:
	rm -f terraform-provider-bmlt

.PHONY: build test testacc vet fmt fmtcheck lint tools docs install clean
