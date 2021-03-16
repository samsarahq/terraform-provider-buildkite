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

func testAccTeamMemberConfig() string {
	return fmt.Sprintf(`
%s

%s

resource "buildkite_team_member" "test" {
	user_id = "${data.buildkite_user.me.id}"
	team_id = "${buildkite_team.devexp.id}"
}`, testAccTeamConfig("testAccTeamMember"), testAccUserConfig)
}

func TestAccTeamMember(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      testAccTeamMemberDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTeamMemberConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccTeamMemberExists("buildkite_team_member.test"),
					resource.TestCheckResourceAttrSet("buildkite_team_member.test", "id"),
				),
			},
		},
	})
}

func testAccTeamMemberExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Resource %s not found", name)
		}
		id := rs.Primary.ID
		if _, err := cli.ReadTeamMember(id); err != nil {
			return err
		}
		return nil
	}
}

func testAccTeamMemberDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "buildkite_team_member" {
			continue
		}
		toDelete := &client.TeamMember{
			ID: graphql.String(rs.Primary.ID),
		}
		if err := cli.DeleteTeamMember(toDelete); err != nil {
			if !strings.Contains(err.Error(), "No team member found") {
				return err
			}
		}
	}
	return nil
}
