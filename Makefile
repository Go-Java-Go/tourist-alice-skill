GOCMD=go
GOBUILD=$(GOCMD) build
BINARY_NAME=./bin/tourist-alice-skill
DOCKER_IMAGE_NAME=tourist-alice-skill

all: build test docker run_server

wire:
	wire gen ./cmd/tourist-alice-skill
	echo "wire build"

build:
	$(GOBUILD) -o $(BINARY_NAME) ./cmd/tourist-alice-skill
	echo "binary build"

test:
	go test -v ./... --tags=test

mock:
	mockery --all --case snake --note "+build test"

docker:
	docker build . -t $(DOCKER_IMAGE_NAME)

run_server:
	 LOG_LEVEL=debug $(BINARY_NAME)

