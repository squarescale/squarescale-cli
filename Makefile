GIT_BRANCH    := $(shell git rev-parse --abbrev-ref HEAD || echo 'n/a')
GIT_VERSION   := $(shell ./get-git-ref.sh version)
GIT_REVISION  := $(shell ./get-git-ref.sh revision)
GO_BUILD_DATE ?= $(shell date -u +%FT%T)
GO_ARCH       := $(shell go env GOARCH)
GO_OS         := $(shell go env GOOS)

GO_LD_FLAGS ?= -ldflags "-X main.GoArch=$(GO_ARCH) \
                         -X main.GoOs=$(GO_OS) \
                         -X main.Version=$(GIT_VERSION) \
                         -X main.GitCommit=$(GIT_REVISION) \
                         -X main.GitBranch=$(GIT_BRANCH)   \
                         -X main.BuildDate=$(GO_BUILD_DATE)"

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

docker-alpine-amd64: ## Compile for linux-amd64
	test -d dist || mkdir dist
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO_CMD) -o dist/sqsc-alpine-amd64$(DIST_SUFFIX)

dist%gh-actions: ## Compile for gh-actions
	DIST_SUFFIX='-gh-actions' make dist

build-gh-actions:
	cd ./gh-actions && go build

clean: ## Clean repository
	go clean && rm -rf sqsc sqsc-* dist/

clean-dist: ## Clean repository
	rm -rf dist/

lint: ## Lint Docker
	$(DOCKER_CMD) -v $$PWD:/root/ projectatomic/dockerfile-lint dockerfile_lint
	$(DOCKER_CMD) -i hadolint/hadolint < Dockerfile

tests: ## Run test suites in all packages
	ginkgo -r

coverage: ## Run test suites in all packages with code coverage
	go test ./... -cover -coverprofile=coverage.out
	go tool cover -func=coverage.out

coverage_html: coverage ## Show code coverage html report
	go tool cover -html=coverage.out -o coverage.html

changelog:
	git-chglog v1.1.3.. > CHANGELOG.md
