# Define variables
GO := go
BINARY_NAME := xtinvasion

# Default target: build the program
all: build

# Build the program
build:
	$(GO) build -o ./build/$(BINARY_NAME)

# Clean the build
clean:
	rm -f $(BINARY_NAME)

# Run the program with default number of aliens (10)
run:
	$(GO) run . --aliens=10 --file=map.txt

# Run the program with a specific number of aliens (e.g., 20)
run-with-aliens:
	$(GO) run . start --aliens=$(num) --file=map.txt

test:
	$(GO) test ./...

# Help target: print available targets
help:
	@echo "Available targets:"
	@echo "  make build        - Build the program"
	@echo "  make clean        - Clean the build"
	@echo "  make run          - Run the program with default number of aliens (10)"
	@echo "  make run-with-aliens num=20  - Run the program with a specific number of aliens (e.g., 20)"
	@echo "  make help         - Print available targets"
