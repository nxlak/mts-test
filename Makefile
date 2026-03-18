.PHONY: build test run clean

build:
	go build -o mts-test ./cmd/app

test:
	go test ./...

run: build
	./mts-test $(REPO)

clean:
	go clean
	rm -f mts-test
