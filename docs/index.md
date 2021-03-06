---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "buildkite Provider"
subcategory: ""
description: |-

---

# buildkite Provider

This is a Terraform provider for [Buildkite](https://buildkite.com). It can be used to manage a specific organization in Buildkite and accepts an API token and organization slug either via the parameters below or the environment variables `BUILDKITE_ORGANIZATION_SLUG` and `BUILDKITE_TOKEN`. The API token provided must have full GQL access as well as read/write access to the REST API, more documentation [here](https://buildkite.com/docs/apis/managing-api-tokens).

## Example
```hcl
provider "buildkite" {
  version = "0.3.1"
}
```
<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- **api_token** (String)
- **organization_slug** (String)
