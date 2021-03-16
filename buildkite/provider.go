package buildkite

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/samsara-dev/terraform-provider-buildkite/buildkite/client"
)

// Constants for environment variable names
const (
	OrgEnvVar   = "BUILDKITE_ORGANIZATION_SLUG"
	TokenEnvVar = "BUILDKITE_TOKEN"
)

// Provider returns the sole provider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"organization_slug": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc(OrgEnvVar, nil),
			},
			"api_token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc(TokenEnvVar, nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"buildkite_pipeline":          resourcePipeline(),
			"buildkite_pipeline_schedule": resourcePipelineSchedule(),
			"buildkite_team":              resourceTeam(),
			"buildkite_team_pipeline":     resourceTeamPipeline(),
			"buildkite_team_member":       resourceTeamMember(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"buildkite_user": dataSourceUser(),
		},
		ConfigureFunc: createClient,
	}
}

func createClient(d *schema.ResourceData) (interface{}, error) {
	org := d.Get("organization_slug").(string)
	token := d.Get("api_token").(string)

	cli, err := client.NewClient(org, token)
	if err != nil {
		return nil, err
	}
	return cli, nil
}
