package buildkite

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/samsara-dev/terraform-provider-buildkite/buildkite/client"
	"github.com/shurcooL/graphql"
)

func resourceTeam() *schema.Resource {
	return &schema.Resource{
		Description: "A resource representing a team in Buildkite.",
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"privacy": &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"VISIBLE", "SECRET"}, false),
				Required:     true,
			},
			"is_default_team": &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
			},
			"default_member_role": &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"MAINTAINER", "MEMBER"}, false),
				Required:     true,
			},
		},
		Create: createTeam,
		Read:   readTeam,
		Update: updateTeam,
		Delete: deleteTeam,
		Importer: &schema.ResourceImporter{
			State: importTeam,
		},
	}
}

func createTeam(d *schema.ResourceData, m interface{}) error {
	bk := m.(*client.Client)
	name := d.Get("name").(string)
	privacy := d.Get("privacy").(string)
	isDefaultTeam := d.Get("is_default_team").(bool)
	defaultMemberRole := d.Get("default_member_role").(string)
	team := &client.Team{
		Name:              graphql.String(name),
		Privacy:           graphql.String(privacy),
		IsDefaultTeam:     graphql.Boolean(isDefaultTeam),
		DefaultMemberRole: graphql.String(defaultMemberRole),
	}
	if err := bk.CreateTeam(team); err != nil {
		return err
	}
	d.SetId(string(team.ID))
	return nil
}

func readTeam(d *schema.ResourceData, m interface{}) error {
	bk := m.(*client.Client)
	team, err := bk.ReadTeam(d.Id())
	if err != nil {
		return err
	}
	d.Set("name", team.Name)
	d.Set("privacy", team.Privacy)
	d.Set("is_default_team", team.IsDefaultTeam)
	d.Set("default_member_role", team.DefaultMemberRole)
	return nil
}

func updateTeam(d *schema.ResourceData, m interface{}) error {
	bk := m.(*client.Client)
	id := d.Id()
	name := d.Get("name").(string)
	privacy := d.Get("privacy").(string)
	isDefaultTeam := d.Get("is_default_team").(bool)
	defaultMemberRole := d.Get("default_member_role").(string)
	team := &client.Team{
		ID:                graphql.String(id),
		Name:              graphql.String(name),
		Privacy:           graphql.String(privacy),
		IsDefaultTeam:     graphql.Boolean(isDefaultTeam),
		DefaultMemberRole: graphql.String(defaultMemberRole),
	}
	if err := bk.UpdateTeam(team); err != nil {
		return err
	}
	return nil
}

func deleteTeam(d *schema.ResourceData, m interface{}) error {
	bk := m.(*client.Client)
	id := d.Id()
	if err := bk.DeleteTeam(&client.Team{ID: graphql.String(id)}); err != nil {
		return err
	}
	return nil
}

func importTeam(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	bk := m.(*client.Client)
	team, err := bk.ReadTeamByName(d.Id())
	if err != nil {
		return nil, err
	}
	d.SetId(string(team.ID))
	d.Set("name", team.Name)
	d.Set("privacy", team.Privacy)
	d.Set("is_default_team", team.IsDefaultTeam)
	d.Set("default_member_role", team.DefaultMemberRole)
	return []*schema.ResourceData{d}, nil
}
