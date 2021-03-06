---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "buildkite_team_member Resource - terraform-provider-buildkite"
subcategory: ""
description: |-
  An association between a user and a team.
---

# buildkite_team_member (Resource)

An association between a user and a team.

## Example
```hcl
data "buildkite_user" "admin" {
	email = "dev@yourorg.com"
}

resource "buildkite_team" "admins" {
  name = "admins"
  ...
}

resource "buildkite_team_member" "test" {
	user_id = "${data.buildkite_user.admin.id}"
	team_id = "${buildkite_team.admins.id}"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **team_id** (String)
- **user_id** (String)

### Read-Only

- **id** (String) The ID of this resource.


