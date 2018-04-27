APPLICATION_NAME = gosharexserver
VERSION = 0.5.6
BRANCH = $(shell git rev-parse --abbrev-ref HEAD)
COMMIT = $(shell git rev-parse HEAD)

LD_FLAGS = -X "main.applicationName=${APPLICATION_NAME}" -X "main.version=${VERSION}" -X "main.branch=${BRANCH}" -X "main.commit=${COMMIT}"

# add dependencies to vendor folder according to the Gopkg.lock contents
.PHONY: dep
dep:
	@dep ensure -vendor-only

# builds and formats the project with the built-in Golang tool
.PHONY: build
build:
	@go build -ldflags '${LD_FLAGS}' ./cmd/gosharexserver

# build go application for docker usage
.PHONY: build-docker
build-docker:
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags '${LD_FLAGS}' ./cmd/gosharexserver

# installs and formats the project with the built-in Golang tool
install:
	@go install -ldflags '${LD_FLAGS}' ./cmd/gosharexserver

# tests the project by running all test go files
.PHONY: test
test:
	@go test -race $(shell go list ./... | grep -v /vendor/ | grep -v /cmd/)
