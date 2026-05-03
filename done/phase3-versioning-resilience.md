# Phase 3: Versioning & Resilience (v3)

## Goal
Add version-aware logging and graceful fallback mechanisms to improve resilience.

## Steps
1. **Add version-aware logging** – Log the ChromeDriver version and binary path before starting the service.
2. **Graceful fallback** – If ChromeDriver fails, attempt to use a fallback browser (e.g., Firefox) or skip browser tests with a clear message.
3. **Update `improve-error.md`** – Document the fix and version history (v1, v2, v3) for future reference.

## Acceptance Criteria
- The test suite logs the ChromeDriver version and binary path.
- If ChromeDriver is unavailable, the test suite either falls back to Firefox or skips browser tests with a clear message.
- The test suite exits with a non-zero exit code if no browser is available.
