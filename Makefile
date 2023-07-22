# Define variables
GO := go
BINARY_NAME := xtinvasion
PROJECTNAME := $(shell basename "$(PWD)")

# Default target: build the program
all: build

## build: Build the program
build:
	$(GO) build cmd/cli/cli.go -o ./build/$(BINARY_NAME)

## clean: Clean the build
clean:
	rm -f $(BINARY_NAME)

## start: Run the program with the specified number of aliens, input file, log file, and output file
start:
	$(GO) run cmd/cli/cli.go start --aliens=$(aliens) --input=$(input) --log=$(log) --output=$(output)

## start-tui: Runs the simulation with terminal UI to follow activity
start-tui:
	$(GO) run cmd/tui/tui.go start --aliens=$(aliens) --input=$(input) --log=$(log) --output=$(output) --delay --delay_ms=$(delay_ms) --max_moves=$(max_moves)

## start-help: Print simulation help
make start-help:
	$(GO) run cmd/cli/cli.go start --help

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