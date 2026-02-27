.PHONY: setup dev build clean deps

# Default target
all: deps build

# Install frontend dependencies
deps:
	@echo "Installing frontend dependencies..."
	cd frontend && npm install

# Setup ONNX Runtime library
setup:
	@echo "Setting up ONNX Runtime..."
	@if [ ! -f /usr/lib/x86_64-linux-gnu/libonnxruntime.so.1.18.1 ]; then \
		echo "Downloading ONNX Runtime..."; \
		mkdir -p /tmp/onnx; \
		cd /tmp/onnx && \
		wget -q https://github.com/microsoft/onnxruntime/releases/download/v1.18.1/onnxruntime-linux-x64-1.18.1.tgz && \
		tar -xzf onnxruntime-linux-x64-1.18.1.tgz && \
		sudo cp onnxruntime-linux-x64-1.18.1/lib/* /usr/lib/x86_64-linux-gnu/ && \
		sudo ldconfig; \
		echo "ONNX Runtime installed successfully"; \
	else \
		echo "ONNX Runtime already installed"; \
	fi

# Development mode
dev: deps
	@echo "Starting development server..."
	CGO_ENABLED=1 /home/oscar/.asdf/installs/golang/1.26.0/bin/wails dev

# Build for production
build:
	@echo "Building application..."
	CGO_ENABLED=1 /home/oscar/.asdf/installs/golang/1.26.0/bin/wails build

# Clean build artifacts
clean:
	rm -rf frontend/dist
	rm -rf build/bin
