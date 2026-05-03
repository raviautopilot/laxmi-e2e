#!/bin/bash
set -euo pipefail

# E2E Test Runner
# Usage: ./run.sh [web|api|all|coverage|clean]

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

show_help() {
    cat <<EOF
Usage: ./run.sh [command]

Commands:
  web       Run web tests only (requires build tag "web")
  api       Run API tests only (requires build tag "api")
  all       Run all tests (web + api)
  coverage  Run all tests and generate coverage report
  clean     Remove build artifacts
  help      Show this help message
EOF
}

case "${1:-help}" in
    web)
        echo "Running web tests..."
        go test -tags=web -v .
        ;;
    api)
        echo "Running API tests..."
        go test -tags=api -v .
        ;;
    all)
        echo "Running all tests..."
        go test -tags="web api" -v .
        ;;
    coverage)
        echo "Running all tests with coverage..."
        go test -tags="web api" -coverprofile=coverage.out -covermode=atomic .
        go tool cover -html=coverage.out -o coverage.html
        echo "Coverage report generated: coverage.html"
        ;;
    clean)
        echo "Cleaning build artifacts..."
        rm -f coverage.out coverage.html
        rm -rf /tmp/e2e-tests-*
        echo "Done."
        ;;
    help|*)
        show_help
        ;;
esac
