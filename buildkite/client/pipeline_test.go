package client

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/likexian/gokit/assert"
)

const (
	repoName = "git@github.com:samsara-dev/terraform-provider-buildkite.git"
)

func strPtr(s interface{}) *string {
	if s == nil {
		return strPtr("")
	}
	str, ok := s.(string)
	if !ok {
		return strPtr("")
	}
	return &str
}

func TestPipelineCRUD(t *testing.T) {
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

	// Test create.
	if err := cli.CreatePipeline(p); err != nil {
		t.Errorf("Could not create pipeline: %s", err)
	}

	p, err := cli.ReadPipeline(*p.Name)
	if err != nil {
		t.Errorf("Could not read pipeline: %s", err)
	}
	if p.ID == nil || *p.ID == "" {
		t.Errorf("Pipeline ID was blank")
	}

	// Test getting gql ID.
	id, err := cli.GetPipelineID(*p.Slug)
	if err != nil {
		t.Errorf("Could not get pipeline id: %s", err)
	}
	if id == "" {
		t.Errorf("Pipeline ID was blank")
	}

	// Test update.
	p.DefaultBranch = strPtr("notmaster")
	if err := cli.UpdatePipeline(p); err != nil {
		t.Errorf("Could not update pipeline: %s", err)
	}

	updatedP, err := cli.ReadPipeline(*p.Name)
	if err != nil {
		t.Errorf("Could not read pipeline: %s", err)
	}
	assert.Equal(t, p, updatedP)

	// Test delete.
	if err := cli.DeletePipeline(p); err != nil {
		t.Errorf("Could not delete pipeline: %s", err)
	}
}
