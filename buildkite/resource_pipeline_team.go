package buildkite

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/samsara-dev/terraform-provider-buildkite/buildkite/client"
	"github.com/shurcooL/graphql"
)

func resourceTeamPipeline() *schema.Resource {
	return &schema.Resource{
		Description: "A resource representing a team with permission on a pipeline in Buildkite.",
		Schema: map[string]*schema.Schema{
			"team_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"pipeline_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"access_level": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"MANAGE_BUILD_AND_READ", "BUILD_AND_READ", "READ_ONLY"}, false),
			},
		},
		Create: createPipelineTeam,
		Read:   readPipelineTeam,
		Update: updatePipelineTeam,
		Delete: deletePipelineTeam,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func createPipelineTeam(d *schema.ResourceData, m interface{}) error {
	bk := m.(*client.Client)

	tp := &client.TeamPipeline{
		AccessLevel: graphql.String(d.Get("access_level").(string)),
		Team: struct{ ID graphql.String }{
			ID: graphql.String(d.Get("team_id").(string)),
		},
		Pipeline: struct{ ID graphql.String }{
			ID: graphql.String(d.Get("pipeline_id").(string)),
		},
	}

	if err := bk.CreateTeamPipeline(tp); err != nil {
		return err
	}
	d.SetId(string(tp.ID))
	return nil
}

func readPipelineTeam(d *schema.ResourceData, m interface{}) error {
	bk := m.(*client.Client)
	tp, err := bk.ReadTeamPipeline(d.Id())
	if err != nil {
		return err
	}

	d.Set("team_id", string(tp.Team.ID))
	d.Set("pipeline_id", string(tp.Pipeline.ID))
	d.Set("access_level", string(tp.AccessLevel))
	return nil
}

func updatePipelineTeam(d *schema.ResourceData, m interface{}) error {
	bk := m.(*client.Client)

	tp := &client.TeamPipeline{
		ID:          graphql.String(d.Id()),
		AccessLevel: graphql.String(d.Get("access_level").(string)),
		Team: struct{ ID graphql.String }{
			ID: graphql.String(d.Get("team_id").(string)),
		},
		Pipeline: struct{ ID graphql.String }{
			ID: graphql.String(d.Get("pipeline_id").(string)),
		},
	}
	if err := bk.UpdateTeamPipeline(tp); err != nil {
		return err
	}
	return nil
}

func deletePipelineTeam(d *schema.ResourceData, m interface{}) error {
	bk := m.(*client.Client)
	id := d.Id()
	if err := bk.DeleteTeamPipeline(&client.TeamPipeline{ID: graphql.String(id)}); err != nil {
		return err
	}
	return nil
}
