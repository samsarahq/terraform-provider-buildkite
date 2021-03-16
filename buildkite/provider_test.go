package buildkite

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/samsara-dev/terraform-provider-buildkite/buildkite/client"
)

var testAccProviderFactory = map[string]terraform.ResourceProviderFactory{
	"buildkite": func() (terraform.ResourceProvider, error) {
		return Provider(), nil
	},
}

const repoName = "git@github.com:samsara-dev/terraform-provider-buildkite.git"

var (
	cli       *client.Client
	userEmail string
)

func init() {
	c, err := client.NewClient(os.Getenv(OrgEnvVar), os.Getenv(TokenEnvVar))
	if err != nil {
		panic("Couldn't create client")
	}
	cli = c
	userEmail = os.Getenv("BUILDKITE_USER_EMAIL")
}

func TestProvider(t *testing.T) {
	p := Provider().(*schema.Provider)
	if err := p.InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {
	for _, env := range []string{OrgEnvVar, TokenEnvVar} {
		if err := os.Getenv(env); err == "" {
			t.Fatalf("%s must be set for acceptance tests", env)
		}
	}
}
