.PHONY : test build clean format

build:
	go build github.com/telkomdev/tob/cmd/tob

test:
	go test ./...

test-verbose:
	go test -v ./...

format:
	find . -name "*.go" -not -path "*/vendor/*" -not -path "*/.git/*" -not -path "*/volumes/*" | xargs gofmt -s -d -w

clean:
	rm tob *.txt