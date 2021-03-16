package client

import (
	"reflect"
	"testing"

	"github.com/shurcooL/graphql"
)

func TestTeamPipelineCRUD(t *testing.T) {
	p, err := setupPipeline()
	if err != nil {
		t.Errorf("Couldn't setup pipeline %s", err)
	}
	defer cli.DeletePipeline(p)

	testPipelineID, err := cli.GetPipelineID(*p.Slug)
	if err != nil {
		t.Errorf("Couldn't get pipeline ID %s", err)
	}
	testPipelineSlug := *p.Slug

	team, err := setupTeam()
	if err != nil {
		t.Errorf("Could not setup team %s", err)
	}
	defer cli.DeleteTeam(team)
	testTeamID := team.ID

	// Test pre-condition.
	teams, err := cli.ReadTeamPipelines(testPipelineID)
	if err != nil {
		t.Errorf("Couldn't check PipelineTeams for %s!", testPipelineSlug)
	}
	if len(teams) != 0 {
		t.Errorf("Test pre-condition failed - %s should contain no teams at test start.", testPipelineSlug)
	}

	tp := &TeamPipeline{
		AccessLevel: "MANAGE_BUILD_AND_READ",
		Team: struct{ ID graphql.String }{
			ID: graphql.String(testTeamID),
		},
		Pipeline: struct{ ID graphql.String }{
			ID: graphql.String(testPipelineID),
		},
	}

	// Test Create.
	err = cli.CreateTeamPipeline(tp)
	if err != nil {
		t.Errorf("Couldn't create PipelineTeam: %s", err)
	}
	if tp.ID == "" {
		t.Errorf("Expected ID to be populated.")
	}

	// Test Update.
	tp.AccessLevel = graphql.String("BUILD_AND_READ")
	err = cli.UpdateTeamPipeline(tp)
	if err != nil {
		t.Errorf("Couldn't update PipelineTeam: %s", err)
	}

	// Test Read.
	updatedTeamPipeline, err := cli.ReadTeamPipeline(string(tp.ID))
	if err != nil {
		t.Errorf("Couldn't read PipelineTeam: %s", err)
	}
	if !reflect.DeepEqual(tp, updatedTeamPipeline) {
		t.Errorf("Updated PipelineTeam did not match expected PipelineTeam")
	}

	// Test Delete.
	err = cli.DeleteTeamPipeline(updatedTeamPipeline)
	if err != nil {
		t.Errorf("Could not delete PipelineTeam: %s", err)
	}

	// Make sure its gone.
	teams, err = cli.ReadTeamPipelines(testPipelineID)
	if err != nil {
		t.Errorf("Could not read PipelineTeams: %s", err)
	}
	if len(teams) != 0 {
		t.Errorf("Expected 0 teams for pipeline, but found %v", len(teams))
	}
}
