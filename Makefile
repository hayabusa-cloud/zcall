GO ?= go

.PHONY: all fmt vet test bench clean

all: fmt vet test

fmt:
	$(GO) fmt ./...

vet:
	@$(GO) vet ./... 2>&1 | grep -v -E "(^#|possible misuse of unsafe.Pointer)" && exit 1 || true

test:
	$(GO) test -timeout 120s -race ./...

bench:
	$(GO) test -timeout 120s -bench=. -benchmem -run=^$$ ./...

clean:
	rm -f coverage.out
	$(GO) clean -testcache
