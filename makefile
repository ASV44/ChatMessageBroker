# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BROKER_BINARY_NAME=broker_app
CLIENT_BINARY_NAME=client_app

all: build
build:
				$(GOBUILD) -v -o $(BROKER_BINARY_NAME) broker/cmd/main.go
				$(GOBUILD) -v -o $(CLIENT_BINARY_NAME) client/cmd/main.go
test:
				$(GOTEST) -v ./...
clean:
				$(GOCLEAN)
				rm -f $(BROKER_BINARY_NAME)
				rm -f $(CLIENT_BINARY_NAME)
run:
				build
				./$(BROKER_BINARY_NAME)
deps:
