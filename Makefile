BINARY_NAME=terraform-provider-buildkite
VERSION=1.0.0
BIN_PATH ?= ./bin

.PHONY: build
build:
	@GOOS=darwin GOARCH=amd64 go build -o ${BIN_PATH}/${BINARY_NAME}_v${VERSION}_darwin_amd64
	@GOOS=linux GOARCH=amd64 go build -o ${BIN_PATH}/${BINARY_NAME}_v${VERSION}_linux_amd64

.PHONY: test
test:
	@go test -v ./...

.PHONY: testacc
testacc:
	@TF_ACC=1 go test -v ./...

.PHONY: clean
clean:
	@rm -rf ${BIN_PATH}
