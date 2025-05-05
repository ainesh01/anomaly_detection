# Makefile for Anomaly Detection Project

# Variables
BACKEND_BIN = anomaly_detection_server
FRONTEND_DIR = frontend

# Default target
.PHONY: all
all: run-all

# Backend targets
.PHONY: build-backend
build-backend:
	@echo "Building backend..."
	@go build -mod=vendor -o $(BACKEND_BIN) ./cmd/main.go
	@echo "Backend built: $(BACKEND_BIN)"

.PHONY: run-backend
run-backend: build-backend
	@echo "Running backend..."
	@./$(BACKEND_BIN) -file=1743566710-000000000000.jsonl.gz

# Frontend targets
.PHONY: install-frontend
install-frontend:
	@echo "Installing frontend dependencies..."
	@cd $(FRONTEND_DIR) && pnpm install
	@echo "Frontend dependencies installed."

.PHONY: build-frontend
build-frontend: install-frontend
	@echo "Building frontend..."
	@cd $(FRONTEND_DIR) && pnpm build
	@echo "Frontend built."

.PHONY: run-frontend
run-frontend: install-frontend
	@echo "Running frontend dev server..."
	@cd $(FRONTEND_DIR) && pnpm dev

# Combined target
.PHONY: run-all
run-all: 
	@echo "Starting backend and frontend..."
	@$(MAKE) run-backend & $(MAKE) run-frontend
	@echo "Backend and frontend running."

# Clean target
.PHONY: clean
clean:
	@echo "Cleaning up..."
	@rm -f $(BACKEND_BIN)
	@echo "Cleanup complete." 