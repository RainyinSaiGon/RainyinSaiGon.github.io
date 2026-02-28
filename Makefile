.DEFAULT_GOAL := build

.PHONY: fmt vet build clean run serve help

fmt:
	go fmt ./...

vet: fmt
	go vet ./...

build: vet
	go build -o portfolio.exe

run:
	go run .

serve:
	go run .
	python -m http.server 8080 --directory docs

clean:
	rm -rf docs/
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