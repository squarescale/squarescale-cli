VERSION     = $(shell git describe --always --dirty)
GO_LD_FLAGS = -ldflags '-X main.GitCommit="$(VERSION)"'
GO_CMD      = go build -v $(GO_LD_FLAGS)
DOCKER_CMD  = docker run --rm

.PHONY: all

all: build

help: ## This help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: deps ## Build CLI (./sqsc binary)
	$(GO_CMD) -o sqsc

deps: ## Install dependencies inside $GOPATH
	go get -v -t -d ./...

dist-linux-amd64: ## Compile for linux-amd64 in a container
	test -d dist || mkdir dist
	GOOS=linux GOARCH=amd64 $(GO_CMD) -o dist/sqsc-linux-amd64$(DIST_SUFFIX)

dist-darwin-amd64: ## Compile for darwin-amd64 in a container
	test -d dist || mkdir dist
	GOOS=darwin GOARCH=amd64 $(GO_CMD) -o dist/sqsc-darwin-amd64$(DIST_SUFFIX)

dist-alpine-amd64: ## Compile for linux-amd64 alpine in a container
	test -d dist || mkdir dist
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO_CMD) -o dist/sqsc-alpine-amd64$(DIST_SUFFIX)

dist: dist-linux-amd64 dist-darwin-amd64 dist-alpine-amd64

dist-master:
	DIST_SUFFIX='-staging-latest' make dist

dist-production:
	DIST_SUFFIX='-latest' make dist

clean: ## Clean repository
	go clean && rm -rf sqsc sqsc-* dist/

clean-dist: ## Clean repository
	rm -rf dist/

lint: ## Lint Docker
	$(DOCKER_CMD) -v $$PWD:/root/ projectatomic/dockerfile-lint dockerfile_lint
	$(DOCKER_CMD) -i sjourdan/hadolint < Dockerfile

tests: ## Run test suites in all packages
	ginkgo -r

coverage: ## Run test suites in all packages with code coverage
	go test ./... -cover -coverprofile=coverage.out

coverage_html: ## Show code coverage html report
	go tool cover -html=coverage.out -o coverage.html
