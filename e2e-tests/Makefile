.PHONY: test-web test-api test-all coverage lint clean

# Run web tests (requires build tag "web")
test-web:
	go test -tags=web -v .

# Run API tests (requires build tag "api")
test-api:
	go test -tags=api -v .

# Run all tests (both web and API)
test-all:
	go test -tags="web api" -v .

# Generate coverage report
coverage:
	go test -tags="web api" -coverprofile=coverage.out -covermode=atomic .
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run linter (requires golangci-lint)
lint:
	golangci-lint run ./...

# Clean build artifacts
clean:
	rm -f coverage.out coverage.html
	rm -rf /tmp/e2e-tests-*
