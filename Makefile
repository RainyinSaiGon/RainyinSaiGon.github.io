.DEFAULT_GOAL := build

.PHONY: fmt vet build clean run serve help

fmt:
	go fmt ./...

vet: fmt
	go vet ./...

build: vet
	go build -o portfolio.exe

run: build
	./portfolio.exe

clean:
	rm -rf public/
	rm -f portfolio.exe

serve: run
	@echo "Generated site is in public/ folder"
	@echo "To deploy to GitHub Pages, push public/ as your gh-pages branch"

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