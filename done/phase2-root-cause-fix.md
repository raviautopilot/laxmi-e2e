# Phase 2: Root Cause Fix (v2)

## Goal
Fix the root cause of the nil pointer dereference and ensure robust error handling.

## Steps
1. **Add nil checks** – In `NewWebDriver`, check if the ChromeDriver service binary path is empty or if the binary is missing. Return an error instead of proceeding.
2. **Improve error wrapping** – Wrap errors from `NewChromeDriverService` with context (e.g., "starting ChromeDriver service: %w").
3. **Test with a missing binary** – Simulate the scenario to confirm the fix.
4. **Update `improve-error.md`** – Document the fix and version history (v1, v2, v3) for future reference.

## Acceptance Criteria
- The segmentation violation is resolved.
- The error message clearly indicates the missing binary or invalid path.
- The test suite exits with a non-zero exit code and a descriptive error message.
