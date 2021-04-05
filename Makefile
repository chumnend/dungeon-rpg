all: build start

.PHONY: build
build:
	@echo "Building..."
	@mkdir -p bin/
	@cd ./bin && go build ../main.go
	@echo "Build complete."

.PHONY: start
start:
	@echo "Starting App..."
	@./bin/main

.PHONY: clean
clean:
	@echo "Cleaning binaries..."
	@rm -rf bin
	@echo "Clean complete."