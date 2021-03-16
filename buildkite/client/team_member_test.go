package client

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/likexian/gokit/assert"
	"github.com/shurcooL/graphql"
)

func setupTeam() (*Team, error) {
	name := fmt.Sprintf("test-pipeline-%d", rand.Int31n(10000))
	team := &Team{
		Name:              graphql.String(name),
		Privacy:           "VISIBLE",
		IsDefaultTeam:     false,
		DefaultMemberRole: "MAINTAINER",
	}
	return team, cli.CreateTeam(team)
}
func TestTeamMemberCRUD(t *testing.T) {
	team, err := setupTeam()
	if err != nil {
		t.Errorf("Could not setup team %s", err)
	}
	defer cli.DeleteTeam(team)
	teamID := team.ID

	member := &TeamMember{
		TeamID: graphql.String(teamID),
		UserID: graphql.String(userID),
	}

	if err := cli.CreateTeamMember(member); err != nil {
		t.Errorf("Could not create team member: %s", err)
	}
	if member.ID == "" {
		t.Errorf("Member ID was empty.")
	}
	m, err := cli.ReadTeamMember(string(member.ID))
	if err != nil {
		t.Errorf("Could not read team member: %s", err)
	}
	assert.Equal(t, member, m)

	if err := cli.DeleteTeamMember(member); err != nil {
		t.Errorf("Could not delete team member: %s", err)
	}

	m, err = cli.ReadTeamMember(string(member.ID))
	if err != nil {
		t.Errorf("Could not read team member: %s", err)
	}
}
