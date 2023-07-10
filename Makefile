GO_SRCS := $(shell find . -type f -name '*.go' -a ! \( -name 'zz_generated*' -o -name '*_test.go' \))
TAG_NAME = $(shell git describe --tags --abbrev=0 --exact-match 2>/dev/null)
TAG_NAME_DEV = $(shell git describe --tags --abbrev=0 2>/dev/null)
GIT_COMMIT = $(shell git rev-parse --short=7 HEAD)
VERSION = $(or ${TAG_NAME},$(TAG_NAME_DEV)-dev)

bin/benchmark-api: $(GO_SRCS) set-version
	CGO_ENABLED=0 go build -ldflags "-s -w" -o "$@" ./cmd/main.go

bin/benchmark-api-darwin-amd64: $(GO_SRCS) set-version
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w -X main.version=${VERSION}" -o "$@" ./cmd/main.go

bin/benchmark-api-darwin-arm64: $(GO_SRCS) set-version
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags "-s -w -X main.version=${VERSION}" -o "$@" ./cmd/main.go

bin/benchmark-api-freebsd-amd64: $(GO_SRCS) set-version
	CGO_ENABLED=0 GOOS=freebsd GOARCH=amd64 go build -ldflags "-s -w -X main.version=${VERSION}" -o "$@" ./cmd/main.go

bin/benchmark-api-freebsd-arm64: $(GO_SRCS) set-version
	CGO_ENABLED=0 GOOS=freebsd GOARCH=arm64 go build -ldflags "-s -w -X main.version=${VERSION}" -o "$@" ./cmd/main.go

bin/benchmark-api-linux-amd64: $(GO_SRCS) set-version
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w -X main.version=${VERSION}" -o "$@" ./cmd/main.go

bin/benchmark-api-linux-arm64: $(GO_SRCS) set-version
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags "-s -w -X main.version=${VERSION}" -o "$@" ./cmd/main.go

bin/benchmark-api-linux-mips64: $(GO_SRCS) set-version
	CGO_ENABLED=0 GOOS=linux GOARCH=mips64 go build -ldflags "-s -w -X main.version=${VERSION}" -o "$@" ./cmd/main.go

bin/benchmark-api-linux-mips64le: $(GO_SRCS) set-version
	CGO_ENABLED=0 GOOS=linux GOARCH=mips64le go build -ldflags "-s -w -X main.version=${VERSION}" -o "$@" ./cmd/main.go

bin/benchmark-api-linux-ppc64: $(GO_SRCS) set-version
	CGO_ENABLED=0 GOOS=linux GOARCH=ppc64 go build -ldflags "-s -w -X main.version=${VERSION}" -o "$@" ./cmd/main.go

bin/benchmark-api-linux-ppc64le: $(GO_SRCS) set-version
	CGO_ENABLED=0 GOOS=linux GOARCH=ppc64le go build -ldflags "-s -w -X main.version=${VERSION}" -o "$@" ./cmd/main.go

bin/benchmark-api-linux-riscv64: $(GO_SRCS) set-version
	CGO_ENABLED=0 GOOS=linux GOARCH=riscv64 go build -ldflags "-s -w -X main.version=${VERSION}" -o "$@" ./cmd/main.go

bin/benchmark-api-linux-s390x: $(GO_SRCS) set-version
	CGO_ENABLED=0 GOOS=linux GOARCH=s390x go build -ldflags "-s -w -X main.version=${VERSION}" -o "$@" ./cmd/main.go

bins := benchmark-api-darwin-amd64 benchmark-api-darwin-arm64 benchmark-api-freebsd-arm64 benchmark-api-freebsd-arm64 benchmark-api-linux-amd64 benchmark-api-linux-arm64 benchmark-api-linux-mips64 benchmark-api-linux-mips64le benchmark-api-linux-ppc64 benchmark-api-linux-ppc64le benchmark-api-linux-riscv64 benchmark-api-linux-s390x

bin/checksums.txt: $(addprefix bin/,$(bins))
	sha256sum -b $(addprefix bin/,$(bins)) | sed 's/bin\///' > $@

bin/checksums.md: bin/checksums.txt
	@echo "### SHA256 Checksums" > $@
	@echo >> $@
	@echo "\`\`\`" >> $@
	@cat $< >> $@
	@echo "\`\`\`" >> $@

.PHONY:
set-version:
	@sed -Ei 's/Version:(\s+)".*",/Version:\1"$(VERSION)",/g' cmd/main.go

.PHONY: build-all
build-all: $(addprefix bin/,$(bins)) bin/checksums.md

$(golint):
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

.PHONY: lint
lint: $(golint)
	$(golint) run ./...

.PHONY: clean
clean:
	rm -rf bin/

.PHONY: mocks
mocks:
	mockery --all
