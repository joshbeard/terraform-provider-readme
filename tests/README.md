# ReadMe Terraform Provider Acceptance Tests

The tests in this directory are acceptance tests that test against a "real"
project with the live ReadMe.com API.

## ReadMe Account for Testing

A _free_ ReadMe account works and is used for these tests, even though a free
account is limited to only API specifications in the web UI. The API
functionality appears to be sufficient for proper testing across all resources,
though this isn't guaranteed to continue to work.

## Running the tests

Ensure a `README_API_TOKEN` variable is set in the environment. These tests are
destructive and should only be ran against a test account.

The conventional method is to use `make test` in the root of the repo:

```shell
make test
```

Run tests and generate a test coverage report:

```shell
make coverage
```

### Run a single test or specific tests

```shell
go run test -v \
  -run TestDocResourceFrontMatter/title/Frontmatter_for_title_attribute \
  ./tests/...
```
