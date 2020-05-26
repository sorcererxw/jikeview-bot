GO := GO111MODULE=on go
GO_BUILD := CGO_ENABLED=0 $(GO) build

build:
	for CMD in `ls ./cmd`; do \
	  $(GO_BUILD) -o bin/$$CMD cmd/$$CMD/main.go; \
	done

clean:
	rm -rf build
