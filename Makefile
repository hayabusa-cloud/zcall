.PHONY: test bench vet clean

test:
	go test -race ./...

bench:
	go test -bench=. -benchmem -run=^$$ ./...

vet:
	go vet ./... 2>&1 | grep -v "possible misuse of unsafe.Pointer"

clean:
	rm -f coverage.out
