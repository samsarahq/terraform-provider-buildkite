package buildkite

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/samsara-dev/terraform-provider-buildkite/buildkite/client"
	"github.com/shurcooL/graphql"
)

func testAccTeamConfig(name string) string {
	return fmt.Sprintf(`
resource "buildkite_team" "devexp" {
    name = "%s"
    privacy = "VISIBLE"
    is_default_team = false
    default_member_role = "MAINTAINER"
}`, name)
}

func testAccTeamConfigUpdated(name string) string {
	return fmt.Sprintf(`
resource "buildkite_team" "devexp" {
    name = "%s"
    privacy = "VISIBLE"
    is_default_team = false
    default_member_role = "MEMBER"
}`, name)
}

func TestAccTeam(t *testing.T) {
	rName := acctest.RandString(5)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      testAccTeamDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTeamConfig(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccTeamExists("buildkite_team.devexp"),
					resource.TestCheckResourceAttrSet("buildkite_team.devexp", "id"),
					resource.TestCheckResourceAttr("buildkite_team.devexp", "privacy", "VISIBLE"),
					resource.TestCheckResourceAttr("buildkite_team.devexp", "is_default_team", "false"),
					resource.TestCheckResourceAttr("buildkite_team.devexp", "default_member_role", "MAINTAINER"),
				),
			},
			{
				ResourceName:      "buildkite_team.devexp",
				ImportStateId:     rName,
				ImportState:       true,
				ImportStateVerify: true,
				Config:            testAccTeamConfig(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccTeamExists("buildkite_team.devexp"),
					resource.TestCheckResourceAttrSet("buildkite_team.devexp", "id"),
					resource.TestCheckResourceAttr("buildkite_team.devexp", "privacy", "VISIBLE"),
					resource.TestCheckResourceAttr("buildkite_team.devexp", "is_default_team", "false"),
					resource.TestCheckResourceAttr("buildkite_team.devexp", "default_member_role", "MAINTAINER"),
				),
			},
			{
				Config: testAccTeamConfigUpdated(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccTeamExists("buildkite_team.devexp"),
					resource.TestCheckResourceAttrSet("buildkite_team.devexp", "id"),
					resource.TestCheckResourceAttr("buildkite_team.devexp", "default_member_role", "MEMBER"),
				),
			},
		},
	})
}

func testAccTeamExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Resource %s not found", name)
		}
		id := rs.Primary.ID
		if _, err := cli.ReadTeam(id); err != nil {
			return err
		}
		return nil
	}
}

func testAccTeamDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "buildkite_team" {
			continue
		}
		toDelete := &client.Team{ID: graphql.String(rs.Primary.ID)}
		if err := cli.DeleteTeam(toDelete); err != nil {
			if !strings.Contains(err.Error(), "No team found") {
				return err
			}
		}
	}
	return nil
}
