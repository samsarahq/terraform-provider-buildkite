---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "buildkite_pipeline Resource - terraform-provider-buildkite"
subcategory: ""
description: |-
  A resource representing a pipeline in Buildkite.
---

# buildkite_pipeline (Resource)

A resource representing a pipeline in Buildkite.

## Example
```hcl
resource "buildkite_pipeline" "test" {
	name = "build pipeline"
	repository = "git@github.com:your-org/repo.git"
	steps = <<EOF
env:
  IAM_ROLE: "some-role"
steps:
  - label: "test things"
    command: "make test"
  EOF

	branch_configuration = "!master"
	cancel_running_branch_builds = true
	cancel_running_branch_builds_filter = "master"
	default_branch = "master"
	description = "The only pipeline you need"
	skip_queued_branch_builds = true
	skip_queued_branch_builds_filter = "!master"

	provider_settings = {
	  build_pull_requests = true
	  pull_request_branch_filter_enabled = true
	  pull_request_branch_filter_configuration = "mobile/*"
	  skip_pull_request_builds_for_existing_commits = true
	  build_tags = false
	  publish_commit_status = true
	  publish_commit_status_per_step = true
	  trigger_mode = "code"
	  filter_enabled = false
	  filter_condition = "build.message != skip"
	  build_pull_request_forks = false
	  prefix_pull_request_fork_branch_names = true
	  separate_pull_request_statuses  = true
	  publish_blocked_as_pending = true
	}
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **name** (String)
- **repository** (String)
- **steps** (String)

### Optional

- **branch_configuration** (String)
- **cancel_running_branch_builds** (Boolean)
- **cancel_running_branch_builds_filter** (String)
- **default_branch** (String)
- **description** (String)
- **id** (String) The ID of this resource.
- **provider_settings** (Map of String)
- **skip_queued_branch_builds** (Boolean)
- **skip_queued_branch_builds_filter** (String)

### Read-Only

- **slug** (String)


