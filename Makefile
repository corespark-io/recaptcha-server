# Variables
APP_NAME := recaptcha-server
DOCKER_REPO := ghcr.io/corespark-io/recaptcha-server
DOCKER_TAG := latest
PORT := 8080

# Default target
.PHONY: all
all: tidy build

# Run go mod tidy
.PHONY: tidy
tidy:
	@echo "	# Variables
	APP_NAME := recaptcha-server
	VERSION := $(shell cat .version)
	DOCKER_TAG := $(VERSION)
	PORT := 8080
	
	# Default target
	.PHONY: all
	all: tidy build
	
	# Run go mod tidy
	.PHONY: tidy
	tidy:
		@echo "Running go mod tidy..."
		@cd app && go mod tidy
	
	# Build Docker image
	.PHONY: build
	build:
		@echo "Building Docker image with tag $(DOCKER_TAG)..."
		docker build -t $(DOCKER_REPO):$(DOCKER_TAG) .
		docker tag $(DOCKER_REPO):$(DOCKER_TAG) $(DOCKER_REPO):latest
	
	# Run Docker container
	.PHONY: run
	run:
		@echo "Starting Docker container..."
		docker run --rm -p $(PORT):$(PORT) \
			-e RECAPTCHA_SECRET_KEY=${RECAPTCHA_SECRET_KEY} \
			-e RECAPTCHA_PORT=$(PORT) \
			-e RECAPTCHA_FRONTEND=${RECAPTCHA_FRONTEND} \
			--name $(APP_NAME) $(DOCKER_REPO):$(DOCKER_TAG)
	
	# Start development instance
	.PHONY: dev
	dev:
		@echo "Starting development server..."
		cd app && go run ./cmd/main.go
	
	# Push Docker image to repository
	.PHONY: push
	push:
		@echo "Pushing image $(DOCKER_REPO):$(DOCKER_TAG) to Docker repository..."
		docker push $(DOCKER_REPO):$(DOCKER_TAG)
		docker push $(DOCKER_REPO):latest
	
	# Clean up
	.PHONY: clean
	clean:
		@echo "Cleaning up..."
		-docker stop $(APP_NAME) 2>/dev/null || true
		-docker rm $(APP_NAME) 2>/dev/null || true
		-docker rmi $(DOCKER_REPO):$(DOCKER_TAG) 2>/dev/null || true
		-docker rmi $(DOCKER_REPO):latest 2>/dev/null || true
	
	# Help
	.PHONY: help
	help:
		@echo "Available targets:"
		@echo "  all    - Run tidy and build (default)"
		@echo "  tidy   - Run go mod tidy"
		@echo "  build  - Build Docker image (tagged with $(VERSION) and latest)"
		@echo "  run    - Run Docker container"
		@echo "  dev    - Run development server"
		@echo "  push   - Push Docker image to repository"
		@echo "  clean  - Clean up Docker resources"
		@echo "Current version: $(VERSION)"go mod tidy..."
	@cd app && go mod tidy

# Build Docker image
.PHONY: build
build:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_REPO):$(DOCKER_TAG) .

# Run Docker container
.PHONY: run
run:
	@echo "Starting Docker container..."
	docker run --rm -p $(PORT):$(PORT) \
		-e RECAPTCHA_SECRET_KEY=${RECAPTCHA_SECRET_KEY} \
		-e RECAPTCHA_PORT=$(PORT) \
		-e RECAPTCHA_FRONTEND=${RECAPTCHA_FRONTEND} \
		--name $(APP_NAME) $(DOCKER_REPO):$(DOCKER_TAG)

# Start development instance
.PHONY: dev
dev:
	@echo "Starting development server..."
	cd app && go run ./cmd/main.go

# Push Docker image to repository
.PHONY: push
push:
	@echo "Pushing image to Docker repository..."
	docker push $(DOCKER_REPO):$(DOCKER_TAG)

# Clean up
.PHONY: clean
clean:
	@echo "Cleaning up..."
	-docker stop $(APP_NAME) 2>/dev/null || true
	-docker rm $(APP_NAME) 2>/dev/null || true
	-docker rmi $(DOCKER_REPO):$(DOCKER_TAG) 2>/dev/null || true

# Help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all    - Run tidy and build (default)"
	@echo "  tidy   - Run go mod tidy"
	@echo "  build  - Build Docker image"
	@echo "  run    - Run Docker container"
	@echo "  dev    - Run development server"
	@echo "  push   - Push Docker image to repository"
	@echo "  clean  - Clean up Docker resources"