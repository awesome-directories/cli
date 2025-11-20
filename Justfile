binary_name := "awesome-directories"
version := env_var_or_default("VERSION", "dev")
commit := `git rev-parse --short HEAD 2>/dev/null || echo "unknown"`
date := `date -u +%Y-%m-%dT%H:%M:%SZ`
ldflags := "-ldflags \"-X main.version=" + version + " -X main.commit=" + commit + " -X main.date=" + date + " -X main.builtBy=just\""

default: help

help:
    @just --list

build:
    @echo "Building {{binary_name}}..."
    go build {{ldflags}} -o {{binary_name}} ./cmd/awesome-directories

install:
    @echo "Installing {{binary_name}}..."
    go install {{ldflags}} ./cmd/awesome-directories

test:
    @echo "Running tests..."
    go test -v ./...

clean:
    @echo "Cleaning..."
    rm -f {{binary_name}}
    rm -rf dist/ .goreleaser-dist/

run: build
    ./{{binary_name}}

dev *ARGS:
    go run ./cmd/awesome-directories {{ARGS}}

deps:
    go mod download
    go mod tidy

fmt:
    go fmt ./...

lint:
    golangci-lint run

release-snapshot:
    goreleaser release --snapshot --clean

version: build
    ./{{binary_name}} version
