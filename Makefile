.PHONY : test build clean format

build:
	go build github.com/telkomdev/tob/cmd/tob

# build for Apple's OSX 64
build-osx:
	GOOS=darwin GOARCH=amd64 go build -ldflags '-s -w' -o tob github.com/telkomdev/tob/cmd/tob

# build for Linux 64
build-linux:
	GOOS=linux GOARCH=amd64 go build -ldflags '-s -w' -o tob github.com/telkomdev/tob/cmd/tob

# build for Windows 64
build-win:
	GOOS=windows GOARCH=amd64 go build -ldflags '-s -w' -o tob.exe github.com/telkomdev/tob/cmd/tob

# Tob HTTP Agent

# build for Apple's OSX 64
build-http-agent-osx:
	GOOS=darwin GOARCH=amd64 go build -ldflags '-s -w' -o tob-http-agent github.com/telkomdev/tob/cmd/tob-http-agent

# build for Linux 64
build-http-agent-linux:
	GOOS=linux GOARCH=amd64 go build -ldflags '-s -w' -o tob-http-agent github.com/telkomdev/tob/cmd/tob-http-agent

# build for Windows 64
build-http-agent-win:
	GOOS=windows GOARCH=amd64 go build -ldflags '-s -w' -o tob-http-agent.exe github.com/telkomdev/tob/cmd/tob-http-agent

test:
	go test ./...

test-verbose:
	go test -v ./...

format:
	find . -name "*.go" -not -path "*/vendor/*" -not -path "*/.git/*" -not -path "*/volumes/*" | xargs gofmt -s -d -w

clean:
	rm -f tob tob.exe tob-http-agent tob-http-agent.exe *.txt