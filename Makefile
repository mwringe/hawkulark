all: build

TAG = v0.0.1-beta.0
GIT_COMMIT = $(shell git rev-parse HEAD)

build: clean deps
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 godep go build -ldflags "-X main.version=$(TAG) -X main.gitcommit=$(GIT_COMMIT)" -o hawkulark github.com/hawkular/hawkulark/agent

clean: 
	echo "Cleaning up project"

deps:
	echo "Configuring dependencies"
	echo "Installing GoDeps"
	go get github.com/tools/godep


