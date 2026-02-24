TEST?=./...
GOFMT_FILES?=$$(find . -not -path "./vendor/*" -type f -name '*.go')

.PHONY: default
default: test

.PHONY: build
build: fmt
	go build -o bin/zcl cmd/zcl/main.go

.PHONY: test
test: fmt
	go list -mod=vendor $(TEST) | xargs -t -n4 go test $(TESTARGS) -mod=vendor -timeout=2m -parallel=4

.PHONY: travis
travis:
	go test -race -coverprofile=coverage.txt -covermode=atomic -mod=vendor $(TEST)

.PHONY: cover
cover:
	go test $(TEST) -coverprofile=coverage.out
	go tool cover -html=coverage.out
	rm coverage.out

.PHONY: fmt
fmt:
	gofmt -w $(GOFMT_FILES)

.PHONY: release
release: fmt
	$$(go env GOPATH)/bin/goreleaser release

.PHONY: release-test
release-test: fmt
	$$(go env GOPATH)/bin/goreleaser release --snapshot --skip=publish --clean
