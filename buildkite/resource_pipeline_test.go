package buildkite

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/samsara-dev/terraform-provider-buildkite/buildkite/client"
)

func testAccPipelineConfig(name string) string {
	return fmt.Sprintf(`
resource "buildkite_pipeline" "test" {
	name = "%s"
	repository = "%s"
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
`, name, repoName)
}

func testAccPipelineConfigUpdated(name string) string {
	// changes default_branch, provider_settings.build_tags, removes provider_settings.filter_condition
	return fmt.Sprintf(`
resource "buildkite_pipeline" "test" {
	name = "%s"
	repository = "%s"
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
	default_branch = "notmaster"
	description = "The only pipeline you need"
	skip_queued_branch_builds = true
	skip_queued_branch_builds_filter = "!master"

	provider_settings = {
	  build_pull_requests = true
	  pull_request_branch_filter_enabled = true
	  pull_request_branch_filter_configuration = "mobile/*"
	  skip_pull_request_builds_for_existing_commits = true
	  build_tags = true
	  publish_commit_status = true
	  publish_commit_status_per_step = true
	  trigger_mode = "code"
	  filter_enabled = false
	  build_pull_request_forks = false
	  prefix_pull_request_fork_branch_names = true
	  separate_pull_request_statuses  = true
	  publish_blocked_as_pending = true
	}
}
`, name, repoName)
}

func TestAccPipeline(t *testing.T) {
	rName := acctest.RandString(5)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      testAccPipelineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPipelineConfig(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccPipelineExists("buildkite_pipeline.test"),
					resource.TestCheckResourceAttrSet("buildkite_pipeline.test", "id"),
					resource.TestCheckResourceAttrSet("buildkite_pipeline.test", "slug"),
					resource.TestCheckResourceAttr("buildkite_pipeline.test", "name", rName),
					resource.TestCheckResourceAttr("buildkite_pipeline.test", "default_branch", "master"),
					resource.TestCheckResourceAttr("buildkite_pipeline.test", "branch_configuration", "!master"),
					resource.TestCheckResourceAttr("buildkite_pipeline.test", "skip_queued_branch_builds", "true"),
					resource.TestCheckResourceAttr("buildkite_pipeline.test", "provider_settings.pull_request_branch_filter_configuration", "mobile/*"),
					resource.TestCheckResourceAttr("buildkite_pipeline.test", "provider_settings.build_pull_requests", "true"),
					resource.TestCheckResourceAttr("buildkite_pipeline.test", "provider_settings.build_tags", "false"),
					resource.TestCheckResourceAttr("buildkite_pipeline.test", "provider_settings.filter_condition", "build.message != skip"),
				),
			},
			{
				Config: testAccPipelineConfigUpdated(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccPipelineExists("buildkite_pipeline.test"),
					resource.TestCheckResourceAttrSet("buildkite_pipeline.test", "id"),
					resource.TestCheckResourceAttrSet("buildkite_pipeline.test", "slug"),
					resource.TestCheckResourceAttr("buildkite_pipeline.test", "name", rName),
					resource.TestCheckResourceAttr("buildkite_pipeline.test", "default_branch", "notmaster"),
					resource.TestCheckResourceAttr("buildkite_pipeline.test", "branch_configuration", "!master"),
					resource.TestCheckResourceAttr("buildkite_pipeline.test", "skip_queued_branch_builds", "true"),
					resource.TestCheckResourceAttr("buildkite_pipeline.test", "provider_settings.pull_request_branch_filter_configuration", "mobile/*"),
					resource.TestCheckResourceAttr("buildkite_pipeline.test", "provider_settings.build_pull_requests", "true"),
					resource.TestCheckResourceAttr("buildkite_pipeline.test", "provider_settings.build_tags", "true"),
					resource.TestCheckNoResourceAttr("buildkite_pipeline.test", "provider_settings.filter_condition"),
				),
			},
		},
	})
}

func testAccPipelineExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Resource %s not found", name)
		}
		slug := rs.Primary.Attributes["slug"]
		if _, err := cli.ReadPipeline(slug); err != nil {
			return err
		}
		return nil
	}
}

func testAccPipelineDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "buildkite_pipeline" {
			continue
		}
		slug := rs.Primary.Attributes["slug"]
		toDelete := &client.Pipeline{Slug: &slug}
		if err := cli.DeletePipeline(toDelete); err != nil {
			if !strings.Contains(err.Error(), "Not Found") {
				return err
			}
		}
	}
	return nil
}
