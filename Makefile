MAJOR_MINOR = 0.1
BUILD = $(shell  date -u "+%Y%m%d-%H%M%S")
COMMIT_HASH = $(shell git rev-parse --short HEAD 2>/dev/null || echo nocommitinfo)

export CGO_ENABLED ?= 0
export VERSION 	= $(MAJOR_MINOR).$(BUILD).$(COMMIT_HASH)
export FLAGS = $(shell echo "\
        -X gitlab/fenrirunbound/pipeline-queue/internal.buildBranch=$(shell git rev-parse --abbrev-ref HEAD) \
        -X gitlab/fenrirunbound/pipeline-queue/internal.buildCompiler=$(shell go version | cut -f 3 -d' ') \
        -X gitlab/fenrirunbound/pipeline-queue/internal.buildDate=$(shell date -u '+%Y-%m-%dT%H:%M:%SZ') \
        -X gitlab/fenrirunbound/pipeline-queue/internal.buildHash=$(COMMIT_HASH) \
        -X gitlab/fenrirunbound/pipeline-queue/internal.buildUser=$(USER) \
        -X gitlab/fenrirunbound/pipeline-queue/internal.buildVersion=$(VERSION)")

clean:
	rm -rf target vendor

build:
	@./cicd/build.sh

vendor:
	dep ensure

docker:
	docker build --target Builder -t pipeline-queue .

test:
	go test -v ./...
