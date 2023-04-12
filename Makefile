# Set the default shell
SHELL := /bin/bash

# Set the default target
.DEFAULT_GOAL := build

# Find all main.go files in the cmd directory
MAIN_GO_FILES := $(shell find cmd -type f -name main.go)

# Define the binary output files based on the main.go files
BINARIES := $(patsubst cmd/%/main.go,bin/%,$(MAIN_GO_FILES))

# Find the latest template file
TEMPLATE := $(shell find infrastructure/cdk.out -type f -name '*.template.json' -print0 | xargs -r -0 ls -1 -t | head -1)

# Build AMD64 binaries
build: $(BINARIES)

# Rule to build each ARM64 binary
bin/%: cmd/%/main.go
	@echo "Building $@"
	@mkdir -p $(dir $@)
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $@ $<

# Clean the bin directory
clean:
	@echo "Cleaning the bin directory"
	@rm -rf bin

.PHONY: test-create-mapping
test-create-mapping: build
	@cd infrastructure && cdk synth --no-staging
	@sam local invoke CreateMapping -t $(TEMPLATE) -e assets/events/create_mapping_event.json

.PHONY: test-control
test-control: build
	@cd infrastructure && cdk synth --no-staging
	@sam local invoke ControlLambda -t $(TEMPLATE) -e assets/events/control_min_event.json
	@sam local invoke ControlLambda -t $(TEMPLATE) -e assets/events/control_max_event.json

# Deploy the infrastructure
.PHONY: deploy
deploy: build
	@cd infrastructure && cdk deploy
