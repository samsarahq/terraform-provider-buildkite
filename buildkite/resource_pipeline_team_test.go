package buildkite

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/samsara-dev/terraform-provider-buildkite/buildkite/client"
	"github.com/shurcooL/graphql"
)

func testAccTeamPipelineConfig(accessLevel string) string {
	return fmt.Sprintf(`
%s

%s

resource "buildkite_team_pipeline" "tfAccTestTeamPipeline" {
	team_id = "${buildkite_team.devexp.id}"
	pipeline_id = "${buildkite_pipeline.test.id}"
	access_level = "%s"
}`, testAccPipelineConfig("testAccTeamPipeline"), testAccTeamConfig("testAccTeamPipeline"), accessLevel)
}

func TestAccTeamPipeline(t *testing.T) {
	rName := "buildkite_team_pipeline.tfAccTestTeamPipeline"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      testAccTeamPipelineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTeamPipelineConfig("READ_ONLY"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccTeamPipelineExists(rName),
					resource.TestCheckResourceAttr(rName, "access_level", "READ_ONLY"),
				),
			},
			{
				Config: testAccTeamPipelineConfig("BUILD_AND_READ"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccTeamPipelineExists(rName),
					resource.TestCheckResourceAttr(rName, "access_level", "BUILD_AND_READ"),
				),
			},
		},
	})
}

func testAccTeamPipelineExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Resource %s not found", name)
		}
		id := rs.Primary.ID
		if _, err := cli.ReadTeamPipeline(id); err != nil {
			return err
		}
		return nil
	}
}

func testAccTeamPipelineDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "buildkite_team_pipeline" {
			continue
		}
		toDelete := &client.TeamPipeline{ID: graphql.String(rs.Primary.ID)}
		if err := cli.DeleteTeamPipeline(toDelete); err != nil {
			if !strings.Contains(err.Error(), "No team pipeline found") {
				return err
			}
		}
	}
	return nil
}
