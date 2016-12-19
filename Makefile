COMMIT      = $$(git describe --always)
GO_LD_FLAGS = -ldflags "-X main.GitCommit=\"$(COMMIT)\""
GO_CMD      = go build -v $(GO_LD_FLAGS)
DOCKER_CMD  = docker run --rm
MOUNT_POINT = /go/src/github.com/squarescale/squarescale-cli
MOUNT_FLAGS = -v "$(PWD)":"$(MOUNT_POINT)" -w "$(MOUNT_POINT)"

.PHONY: all

all: build

help: ## This help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: deps ## Build CLI (./sqsc binary)
	go build -o sqsc -ldflags "-X main.GitCommit=\"$(COMMIT)\""

deps: ## Install dependencies inside $GOPATH
	go get ./...

docker-linux-amd64: ## Compile for linux-amd64 in a container
	$(DOCKER_CMD) $(MOUNT_FLAGS) -e GOOS=linux -e GOARCH=amd64 golang:1.7 $(GO_CMD) -o sqsc-linux-amd64

docker-darwin-amd64: ## Compile for darwin-amd64 in a container
	$(DOCKER_CMD) $(MOUNT_FLAGS) -e GOOS=darwin -e GOARCH=amd64 golang:1.7 $(GO_CMD) -o sqsc-darwin-amd64

generate: docker-linux-amd64 docker-darwin-amd64

clean: ## Clean repository
	go clean && rm -f sqsc sqsc-*
