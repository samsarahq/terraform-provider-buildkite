package buildkite

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

var testAccUserConfig = fmt.Sprintf(`
data "buildkite_user" "me" {
	email = "%s"
}
`, userEmail)

func TestAccDataSourceUser(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccUserConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.buildkite_user.me", "name"),
					resource.TestCheckResourceAttrSet("data.buildkite_user.me", "uuid"),
					resource.TestCheckResourceAttrSet("data.buildkite_user.me", "id"),
				),
			},
		},
	})
}
