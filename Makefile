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
	go get ./...

build-docker: build
	CGO_ENABLED=0 $(GO_CMD) -o sqsc-docker
	docker build -t sqsc-cli .

docker-linux-amd64: ## Compile for linux-amd64 in a container
	$(DOCKER_CMD) $(MOUNT_FLAGS) -e GOOS=linux -e GOARCH=amd64 golang:1.7 $(GO_CMD) -o sqsc-linux-amd64

docker-darwin-amd64: ## Compile for darwin-amd64 in a container
	$(DOCKER_CMD) $(MOUNT_FLAGS) -e GOOS=darwin -e GOARCH=amd64 golang:1.7 $(GO_CMD) -o sqsc-darwin-amd64

generate: docker-linux-amd64 docker-darwin-amd64

publish: generate
	python3 publish.py
	aws s3 cp sqsc-linux-amd64 s3://cli-releases/sqsc-linux-amd64-latest --acl public-read
	aws s3 cp sqsc-darwin-amd64 s3://cli-releases/sqsc-darwin-amd64-latest --acl public-read

clean: ## Clean repository
	go clean && rm -f sqsc sqsc-*

lint: ## Lint Docker
	docker run --rm -v $$PWD:/root/ projectatomic/dockerfile-lint dockerfile_lint
	hadolint Dockerfile
