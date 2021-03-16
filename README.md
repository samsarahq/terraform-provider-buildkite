[![Build status](https://badge.buildkite.com/ba2febb05f89921c3824ce22bd94d47310dcd481186151de21.svg?branch=main)](https://buildkite.com/samsara/hq-terraform-provider-buildkite)
# Terraform Provider for [Buildkite](https://buildkite.com)

Note: This provider is built for Terraform 0.11 and all docs/examples reflect this.
## Documentation
Documentation for this provider is located in `/docs` with the templates generated according to [Terraform Guidelines](https://www.terraform.io/docs/registry/providers/docs.html#generating-documentation).

## Development
To build binaries run
```
make build
```
which will output the binaries to `./bin` in the format `terraform-provider-buildkite_v0.1.0_${OS}_${ARCH}`.

### Testing
The tests for this provider create and delete *real* resources in a given Buildkite account specified by `BUILDKITE_ORGANIZATION_SLUG` and `BUILDKITE_TOKEN`. A real user already registered in the organization is also required and must be specified via `BUILDKITE_USER_EMAIL`.


Integration tests for just the Buildkite client can be run via:
```
make test
```
while full [Terraform Acceptance Tests](https://www.terraform.io/docs/extend/testing/acceptance-tests/index.html) are run via:
```
make testacc
```
