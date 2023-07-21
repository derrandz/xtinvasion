# Define variables
GO := go
BINARY_NAME := xtinvasion

# Default target: build the program
all: build

# Build the program
build:
	$(GO) build cmd/main.go -o ./build/$(BINARY_NAME)

# Clean the build
clean:
	rm -f $(BINARY_NAME)

# Run the program with default number of aliens (10)
run:
	$(GO) run cmd/main.go --aliens=10 --file=map.txt

# Run the program with a specific number of aliens (e.g., 20)
run-with-aliens:
	$(GO) run cmd/main.go start --aliens=$(num) --input=./data/map.txt

test:
	$(GO) test ./tests/...

# Help target: print available targets
help:
	@echo "Available targets:"
	@echo "  make build        - Build the program"
	@echo "  make clean        - Clean the build"
	@echo "  make run          - Run the program with default number of aliens (10)"
	@echo "  make run-with-aliens num=20  - Run the program with a specific number of aliens (e.g., 20)"
	@echo "  make help         - Print available targets"
