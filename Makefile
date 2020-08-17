GO := CGO_ENABLED=0 GO111MODULE=on go

test:
	$(GO) test -p 1 ./...

build: BUILD_OS ?= linux
build: BUILD_FLAG ?= -ldflags '-extldflags "-static"'
build:
	for CMD in `ls ./cmd`; do \
	  GOOS=$(BUILD_OS) $(GO) build $(BUILD_FLAG) -o bin/$$CMD cmd/$$CMD/*.go; \
	done

clean:
	rm -rf build
