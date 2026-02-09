GOBIN=$(shell go env GOPATH)/bin
# 2. On définit le chemin complet vers Air
AIR=$(GOBIN)/air

BINARY_NAME=mario

all: build

install-deps:
	go install github.com/air-verse/air@latest

build:
	@echo "🔨 Compiling for macOS..."
	@mkdir -p bin
	go build -o bin/$(BINARY_NAME) main.go

dev:
	@if [ ! -f $(AIR) ]; then \
		echo "Air est introuvable dans $(AIR)"; \
		echo "Installation forcée..."; \
		go install github.com/air-verse/air@latest; \
	fi
	@echo "Lancement de Air depuis $(AIR)..."
	$(AIR)

run:
	@go run .

test:
	@go test ./... -v

clean:
	@echo "🧹 Nettoyage..."
	@rm -rf bin
	@go clean

help:
	@echo "Usage:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

doc:
	@echo "Documentation sur http://localhost:8080"
	@pkgsite -open .

.PHONY: all dev build run test clean help install-deps doc