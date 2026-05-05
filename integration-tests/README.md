# Integration Tests

Data-driven black-box tests for rubik-store-go. Each test exercises the public API only — no access to internal fields.

## Running

From the `integration-tests` directory:

```bash
cd integration-tests
go test ./...
```

To run with verbose output:

```bash
go test -v ./...
```

To run a specific test by name:

```bash
go test -v -run "TestName/subtest_name"
```

## Adding Tests

Add a new `_test.go` file for each module under test. Use table-driven tests with a named `ops` slice or similar structure to keep cases data-driven and easy to extend.
