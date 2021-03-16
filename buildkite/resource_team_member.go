package buildkite

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/samsara-dev/terraform-provider-buildkite/buildkite/client"
	"github.com/shurcooL/graphql"
)

func resourceTeamMember() *schema.Resource {
	return &schema.Resource{
		Description: "An association between a user and a team.",
		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"team_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
		},
		Create: createTeamMember,
		Read:   readTeamMember,
		Delete: deleteTeamMember,
	}
}

func createTeamMember(d *schema.ResourceData, m interface{}) error {
	bk := m.(*client.Client)
	member := &client.TeamMember{
		UserID: graphql.String(d.Get("user_id").(string)),
		TeamID: graphql.String(d.Get("team_id").(string)),
	}

	if err := bk.CreateTeamMember(member); err != nil {
		return err
	}
	d.SetId(string(member.ID))
	return nil
}

func readTeamMember(d *schema.ResourceData, m interface{}) error {
	bk := m.(*client.Client)
	member, err := bk.ReadTeamMember(d.Id())
	if err != nil {
		return err
	}
	d.Set("user_id", member.UserID)
	d.Set("team_id", member.TeamID)
	return nil
}

func deleteTeamMember(d *schema.ResourceData, m interface{}) error {
	bk := m.(*client.Client)
	member := &client.TeamMember{
		ID: graphql.String(d.Id()),
	}

	if err := bk.DeleteTeamMember(member); err != nil {
		return err
	}
	return nil
}
