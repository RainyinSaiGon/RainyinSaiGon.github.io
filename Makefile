.DEFAULT_GOAL := build

.PHONY: fmt vet build build-css generate run serve clean help

fmt:
	go fmt ./...

vet: fmt
	go vet ./...

build: vet
	go build -o portfolio.exe

# Generate site HTML only
run:
	go run .

# Compile Tailwind CSS (requires: npm install)
build-css:
	npx tailwindcss -i assets/input.css -o docs/tailwind.css --minify

# Full build: HTML + CSS
generate: run build-css

# Dev server: HTML + CSS + serve
serve:
	go run .
	npx tailwindcss -i assets/input.css -o docs/tailwind.css --minify
	python -m http.server 8080 --directory docs

clean:
	rm -rf docs/ node_modules/
	rm -f portfolio.exe

help:
	@echo "Portfolio Static Site Generator"
	@echo ""
	@echo "Available targets:"
	@echo "  make build   - Build the site (default)"
	@echo "  make run     - Build and generate HTML"
	@echo "  make clean   - Remove generated files"
	@echo "  make fmt     - Format code"
	@echo "  make vet     - Run go vet"
	@echo "  make help    - Show this help message"