package client

import (
	"context"
	"strings"

	"github.com/shurcooL/graphql"
)

// PipelineSchedule represents a schedule of builds to run on a Pipeline.
type PipelineSchedule struct {
	// https://buildkite.com/docs/pipelines/scheduled-builds#schedule-intervals
	Cronline graphql.String
	// Env is a slice of strings of key-value pairs in the form KEY=value.
	Env     []string
	Enabled graphql.Boolean
	Message graphql.String
	Branch  graphql.String
	Commit  graphql.String
	Label   graphql.String
	ID      graphql.String
	// Pipeline is the parent that the pipelineschedule belongs to.
	Pipeline struct {
		ID graphql.String
	}
}

// ReadPipelineSchedules looks up all schedules for a given Pipeline via the Pipeline's Graphql ID.
func (c *Client) ReadPipelineSchedules(pipelineID string) ([]PipelineSchedule, error) {
	type Pipeline struct {
		Schedules struct {
			Edges []struct {
				Node PipelineSchedule
			}
		}
	}
	var query struct {
		Node struct {
			Pipeline Pipeline `graphql:"... on Pipeline"`
		} `graphql:"node(id: $id)"`
	}
	vars := map[string]interface{}{
		"id": pipelineID,
	}
	if err := c.gqlClient.Query(context.TODO(), &query, vars); err != nil {
		return nil, err
	}

	var result []PipelineSchedule
	for _, edge := range query.Node.Pipeline.Schedules.Edges {
		result = append(result, edge.Node)
	}
	return result, nil
}

// ReadPipelineSchedule looks up a PipelineSchedule by given Graphql ID.
func (c *Client) ReadPipelineSchedule(id string) (*PipelineSchedule, error) {
	var query struct {
		Node struct {
			PipelineSchedule PipelineSchedule `graphql:"... on PipelineSchedule"`
		} `graphql:"node(id: $id)"`
	}
	vars := map[string]interface{}{
		"id": id,
	}

	if err := c.gqlClient.Query(context.TODO(), &query, vars); err != nil {
		return nil, err
	}
	return &query.Node.PipelineSchedule, nil
}

// CreatePipelineSchedule creates the provided PipelineSchedule.
func (c *Client) CreatePipelineSchedule(ps *PipelineSchedule) error {
	var mutation struct {
		PipelineScheduleCreate struct {
			PipelineScheduleEdge struct {
				Node PipelineSchedule
			}
		} `graphql:"pipelineScheduleCreate(input: $input)"`
	}

	type PipelineScheduleCreateInput struct {
		PipelineID string `json:"pipelineID"`
		Label      string `json:"label"`
		CronLine   string `json:"cronline"`
		Message    string `json:"message"`
		Commit     string `json:"commit"`
		Branch     string `json:"branch"`
		Enabled    bool   `json:"enabled"`
		Env        string `json:"env"`
	}
	vars := map[string]interface{}{
		"input": PipelineScheduleCreateInput{
			PipelineID: string(ps.Pipeline.ID),
			Label:      string(ps.Label),
			CronLine:   string(ps.Cronline),
			Message:    string(ps.Message),
			Commit:     string(ps.Commit),
			Branch:     string(ps.Branch),
			Enabled:    bool(ps.Enabled),
			Env:        strings.Join(ps.Env, "\n"),
		},
	}

	if err := c.gqlClient.Mutate(context.TODO(), &mutation, vars); err != nil {
		return err
	}

	ps.ID = mutation.PipelineScheduleCreate.PipelineScheduleEdge.Node.ID
	return nil
}

// UpdatePipelineSchedule updates the provided PipelineSchedule.
func (c *Client) UpdatePipelineSchedule(ps *PipelineSchedule) error {
	var mutation struct {
		PipelineScheduleUpdate struct {
			PipelineSchedule PipelineSchedule
		} `graphql:"pipelineScheduleUpdate(input: $input)"`
	}
	type PipelineScheduleUpdateInput struct {
		ID       string `json:"id"`
		Label    string `json:"label"`
		CronLine string `json:"cronline"`
		Message  string `json:"message"`
		Commit   string `json:"commit"`
		Branch   string `json:"branch"`
		Enabled  bool   `json:"enabled"`
		Env      string `json:"env"`
	}
	vars := map[string]interface{}{
		"input": PipelineScheduleUpdateInput{
			ID:       string(ps.ID),
			Label:    string(ps.Label),
			CronLine: string(ps.Cronline),
			Message:  string(ps.Message),
			Commit:   string(ps.Commit),
			Branch:   string(ps.Branch),
			Enabled:  bool(ps.Enabled),
			Env:      strings.Join(ps.Env, "\n"),
		},
	}

	if err := c.gqlClient.Mutate(context.TODO(), &mutation, vars); err != nil {
		return err
	}

	*ps = mutation.PipelineScheduleUpdate.PipelineSchedule
	return nil
}

// DeletePipelineSchedule deletes the provided PipelineSchedule.
func (c *Client) DeletePipelineSchedule(ps *PipelineSchedule) error {
	var mutation struct {
		PipelineScheduleDelete struct {
			DeletedPipelineScheduleID graphql.String `graphql:"deletedPipelineScheduleID"`
		} `graphql:"pipelineScheduleDelete(input: $input)"`
	}
	type PipelineScheduleDeleteInput struct {
		ID string `json:"id"`
	}
	vars := map[string]interface{}{
		"input": PipelineScheduleDeleteInput{
			ID: string(ps.ID),
		},
	}
	return c.gqlClient.Mutate(context.TODO(), &mutation, vars)
}
