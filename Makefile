VERSION     = $(shell git describe --always --dirty)
GO_LD_FLAGS = -ldflags '-X main.GitCommit="$(VERSION)"'
GO_CMD      = go build -v $(GO_LD_FLAGS)
DOCKER_CMD  = docker run --rm
MOUNT_POINT = /go/src/github.com/squarescale/squarescale-cli
MOUNT_FLAGS = -v "$(PWD)":"$(MOUNT_POINT)" -w "$(MOUNT_POINT)"

.PHONY: all

all: build

help: ## This help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: deps ## Build CLI (./sqsc binary)
	$(GO_CMD) -o sqsc

deps: ## Install dependencies inside $GOPATH
	go get github.com/onsi/ginkgo/ginkgo

docker-linux-amd64: ## Compile for linux-amd64 in a container
	$(DOCKER_CMD) $(MOUNT_FLAGS) -e GOOS=linux -e GOARCH=amd64 golang:1.13 $(GO_CMD) -o sqsc-linux-amd64

docker-darwin-amd64: ## Compile for darwin-amd64 in a container
	$(DOCKER_CMD) $(MOUNT_FLAGS) -e GOOS=darwin -e GOARCH=amd64 golang:1.13 $(GO_CMD) -o sqsc-darwin-amd64

docker-alpine-amd64: ## Compile for linux-amd64 alpine in a container
	$(DOCKER_CMD) $(MOUNT_FLAGS) -e CGO_ENABLED=0 -e GOOS=linux -e GOARCH=amd64 golang:1.13 $(GO_CMD) -o sqsc-alpine-amd64

generate: docker-linux-amd64 docker-darwin-amd64 docker-alpine-amd64

publish-staging: ## Publish existing generated build to github draft and s3 as staging
	aws s3 cp sqsc-linux-amd64 s3://cli-releases/sqsc-linux-amd64-staging-latest --acl public-read
	aws s3 cp sqsc-darwin-amd64 s3://cli-releases/sqsc-darwin-amd64-staging-latest --acl public-read
	aws s3 cp sqsc-alpine-amd64 s3://cli-releases/sqsc-alpine-amd64-staging-latest --acl public-read

publish: ## Publish existing generated build
	aws s3 cp sqsc-linux-amd64 s3://cli-releases/sqsc-linux-amd64-latest --acl public-read
	aws s3 cp sqsc-darwin-amd64 s3://cli-releases/sqsc-darwin-amd64-latest --acl public-read
	aws s3 cp sqsc-alpine-amd64 s3://cli-releases/sqsc-alpine-amd64-latest --acl public-read

clean: ## Clean repository
	go clean && rm -f sqsc sqsc-*

lint: ## Lint Docker
	$(DOCKER_CMD) -v $$PWD:/root/ projectatomic/dockerfile-lint dockerfile_lint
	$(DOCKER_CMD) -i sjourdan/hadolint < Dockerfile

tests: ## Run test suites in all packages
	ginkgo -r
