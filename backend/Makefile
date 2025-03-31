.PHONY: all build run gotool clean help

BINARY="news"

all: gotool build

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o ./bin/${BINARY}

run:
	@go run ./main.go conf/config.yaml

gotool:
	go fmt ./
	go vet ./

clean:
	@if [ -f ${BINARY} ]; then rm ${BINARY} ; fi

help:
	@echo "make - Format Go code and compile it into a binary file"
	@echo "make build - Compile Go code and generate a binary file"
	@echo "make run - Run Go code directly"
	@echo "make clean - Remove binary files and vim swap files"
	@echo "make gotool - Run Go tools 'fmt' and 'vet'"
