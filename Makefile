# Define variables
GO := go
BINARY_NAME := xtinvasion
PROJECTNAME := $(shell basename "$(PWD)")

# Default target: build the program
all: build

## build: Build the program
build:
	$(GO) build cmd/main.go -o ./build/$(BINARY_NAME)

## clean: Clean the build
clean:
	rm -f $(BINARY_NAME)

## start: Run the program with the specified number of aliens, input file, log file, and output file
start:
	$(GO) run cmd/main.go start --aliens=$(aliens) --input=$(input) --log=$(log) --output=$(output)

## start-help: Print simulation help
make start-help:
	$(GO) run cmd/main.go start --help

## test: Run the program with default number of aliens (10)
test:
	$(GO) test ./tests/...

## godoc: Run godoc server
godoc:
	godoc -http=:6060

## help: Get more info on make commands.
helpo: Makefile
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
.PHONY: help