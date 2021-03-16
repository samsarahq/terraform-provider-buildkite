package client

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"

	"github.com/shurcooL/graphql"
)

func TestTeamCRUD(t *testing.T) {
	name := fmt.Sprintf("test-pipeline-%d", rand.Int31n(10000))
	team := &Team{
		Name:              graphql.String(name),
		Privacy:           "VISIBLE",
		IsDefaultTeam:     false,
		DefaultMemberRole: "MAINTAINER",
	}
	if err := cli.CreateTeam(team); err != nil {
		t.Errorf("Couldn't create team: %s", err)
	}
	if team.ID == "" {
		t.Errorf("Team ID was blank")
	}

	team.DefaultMemberRole = "MEMBER"
	if err := cli.UpdateTeam(team); err != nil {
		t.Errorf("Couldn't update team: %s", err)
	}

	updatedTeam, err := cli.ReadTeam(string(team.ID))
	if err != nil {
		t.Errorf("Couldn't read team: %s", err)
	}
	if !reflect.DeepEqual(team, updatedTeam) {
		t.Errorf("Actual team not equal to updated team")
	}
	updatedTeam, err = cli.ReadTeamByName(string(team.Name))
	if err != nil {
		t.Errorf("Couldn't read team by name: %s", err)
	}
	if !reflect.DeepEqual(team, updatedTeam) {
		t.Errorf("Actual team not equal to read team")
	}

	if err := cli.DeleteTeam(team); err != nil {
		t.Errorf("Couldn't delete team: %s", err)
	}
	team, err = cli.ReadTeam(string(team.ID))
	if err != nil {
		t.Errorf("Couldn't read team: %s", err)
	}
	if team.ID != "" {
		t.Error("Expected empty team struct after delete")
	}
}
