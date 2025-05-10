# Define the binary name and paths
BINARY_NAME = istar-api
SWAGGER_GEN_FOLDER = ./docs/swagger
SWAGGER_OUTPUT = swagger.yaml

# Default goal
.DEFAULT_GOAL := build

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	go build -o $(BINARY_NAME) ./cmd/main.go
	@echo "$(BINARY_NAME) built successfully."

# Generate Swagger documentation
swagger:
	@echo "Generating Swagger documentation..."
	#swag init -g main.go -o $(SWAGGER_GEN_FOLDER)
	swag init -g cmd/api/main.go --output docs
	@echo "Swagger documentation generated at $(SWAGGER_GEN_FOLDER)"

# Clean generated files and binary
clean:
	@echo "Cleaning up..."
	rm -f $(BINARY_NAME)
	rm -rf $(SWAGGER_GEN_FOLDER)
	@echo "Cleaned up successfully."

# Run the application
run: build
	@echo "Running $(BINARY_NAME)..."
	./$(BINARY_NAME)

# Full pipeline: clean, build, generate swagger, and run
all: clean swagger build run

.PHONY: build swagger clean run all