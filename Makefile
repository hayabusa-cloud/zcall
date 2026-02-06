.PHONY: test bench vet clean

test:
	go test -race ./...

bench:
	go test -bench=. -benchmem -run=^$$ ./...

vet:
	go vet ./...

clean:
	rm -f coverage.out
