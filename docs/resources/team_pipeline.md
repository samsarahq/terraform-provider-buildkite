---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "buildkite_team_pipeline Resource - terraform-provider-buildkite"
subcategory: ""
description: |-
  A resource representing a team with permission on a pipeline in Buildkite.
---

# buildkite_team_pipeline (Resource)

A resource representing a team with permission on a pipeline in Buildkite.

## Example
```hcl
resource "buildkite_pipeline" "test" {
  name = "build pipeline"
  ...
}

resource "buildkite_team" "admins" {
  name = "admins"
  ...
}

resource "buildkite_team_pipeline" "tfAccTestTeamPipeline" {
	team_id = "${buildkite_team.admins.id}"
	pipeline_id = "${buildkite_pipeline.test.id}"
	access_level = "%s"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **access_level** (String)
- **pipeline_id** (String)
- **team_id** (String)

### Read-Only

- **id** (String) The ID of this resource.

