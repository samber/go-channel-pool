
all: build

deps:
	go get -x github.com/cespare/reflex
	go mod download

build:
	go build -v .

dev:
	go run -v ./example/named/main.go

watch-dev: deps
	reflex -t 50ms -s -- sh -c 'echo \\nBUILDING && make dev && echo Exited \(0\)'

test:
	gotest -v ./...

watch-test: deps
	reflex -t 50ms -s -- sh -c 'make test'

re: clean all
