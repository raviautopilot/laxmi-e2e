# Phase 1: Immediate Diagnosis (v1)

## Goal
Identify the nil pointer source that causes the segmentation violation in `NewWebDriver`.

## Steps
1. **Examine the stack trace** – The crash occurs in `selenium.newService` at `service.go:230`, called from `NewChromeDriverService`. The nil pointer is likely the `service` variable or a path argument.
2. **Check `browser.go:25`** – Verify that the `ChromeDriverService` binary path is correct and that the service is not nil before use.
3. **Validate config** – Ensure `cfg.Browser` is "chrome" and that the ChromeDriver binary exists at the expected location.
4. **Add nil checks** – In `NewWebDriver`, check if the ChromeDriver service binary path is empty or if the binary is missing. Return an error instead of proceeding.
5. **Improve error wrapping** – Wrap errors from `NewChromeDriverService` with context (e.g., "starting ChromeDriver service: %w").

## Acceptance Criteria
- The segmentation violation no longer occurs when running `./run.sh all`.
- A clear error message is printed if the ChromeDriver binary is missing or the path is invalid.
- The test suite exits gracefully (non-zero exit code) instead of panicking.
