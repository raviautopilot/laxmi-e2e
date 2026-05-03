# Phase 4: Validation

## Goal
Verify that the segmentation violation is resolved and that the test suite behaves correctly.

## Steps
1. **Run `./run.sh all`** – Confirm the segmentation violation is resolved.
2. **Run individual test suites** – Ensure API tests still pass and web tests either pass or skip gracefully.
3. **Run with missing ChromeDriver** – Simulate the scenario to confirm the error message is clear and the test suite exits gracefully.
4. **Run with missing Firefox** – Simulate the scenario to confirm the fallback or skip behavior.

## Acceptance Criteria
- The segmentation violation no longer occurs.
- API tests pass.
- Web tests either pass or skip gracefully with a clear message.
- The test suite exits with a non-zero exit code if no browser is available.
