package client

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"

	"github.com/shurcooL/graphql"
)

func setupPipeline() (*Pipeline, error) {
	// Use a random suffix to avoid collisions.
	name := fmt.Sprintf("test-pipeline-%d", rand.Int31n(10000))
	p := &Pipeline{
		Name:        strPtr(name),
		Repository:  strPtr(repoName),
		Description: strPtr("integration tests"),
		Configuration: `
env:
  TEST: true
steps:
- command: echo "hello"
`,
		DefaultBranch: strPtr("master"),
	}

	if err := cli.CreatePipeline(p); err != nil {
		return p, err
	}
	return cli.ReadPipeline(name)
}

func TestPipelineScheduleCRUD(t *testing.T) {
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

	// Test pre-condition.
	schedules, err := cli.ReadPipelineSchedules(testPipelineID)
	if err != nil {
		t.Errorf("Couldn't check PipelineSchedules for %s!", testPipelineSlug)
	}
	if len(schedules) != 0 {
		t.Errorf("Test pre-condition failed - %s should contain no schedules at test start.", testPipelineSlug)
	}

	ps := &PipelineSchedule{
		Branch:   graphql.String("master"),
		Commit:   graphql.String("HEAD"),
		Cronline: graphql.String("0 0 1 1 *"),
		Enabled:  graphql.Boolean(false),
		Env:      []string{"KEY1=val1", "KEY2=val2"},
		Label:    graphql.String("Test label"),
		Message:  graphql.String("Test message"),
		Pipeline: struct {
			ID graphql.String
		}{
			ID: graphql.String(testPipelineID),
		},
	}

	// Test Create.
	err = cli.CreatePipelineSchedule(ps)
	if err != nil {
		t.Errorf("Couldn't create PipelineSchedule: %s", err)
	}
	if ps.ID == "" {
		t.Errorf("Expected ID to be populated.")
	}

	// Test Update.
	ps.Enabled = graphql.Boolean(true)
	err = cli.UpdatePipelineSchedule(ps)
	if err != nil {
		t.Errorf("Couldn't update PipelineSchedule: %s", err)
	}

	// Test Read.
	updatedPipelineSchedule, err := cli.ReadPipelineSchedule(string(ps.ID))
	if err != nil {
		t.Errorf("Couldn't read PipelineSchedule: %s", err)
	}
	if !reflect.DeepEqual(ps, updatedPipelineSchedule) {
		t.Errorf("Updated PipelineSchedule did not match expected PipelineSchedule")
	}

	// Test Delete.
	err = cli.DeletePipelineSchedule(updatedPipelineSchedule)
	if err != nil {
		t.Errorf("Could not delete PipelineSchedule: %s", err)
	}

	// Make sure its gone.
	schedules, err = cli.ReadPipelineSchedules(testPipelineID)
	if err != nil {
		t.Errorf("Could not read PipelineSchedules: %s", err)
	}
	if len(schedules) != 0 {
		t.Errorf("Expected 0 schedules for pipeline, but found %v", len(schedules))
	}
}
