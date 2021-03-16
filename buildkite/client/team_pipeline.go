package client

import (
	"context"

	"github.com/shurcooL/graphql"
)

// TeamPipeline represents the association of a team to a pipeline.
type TeamPipeline struct {
	ID          graphql.String
	AccessLevel graphql.String
	// Team is the parent resource that a PipelineTeam belongs to.
	Team struct {
		ID graphql.String
	}
	// Pipeline is the parent resource that a PipelineTeam belongs to.
	Pipeline struct {
		ID graphql.String
	}
}

// ReadTeamPipelines looks up all teams for a given Pipeline via the Pipeline's Graphql ID.
func (c *Client) ReadTeamPipelines(pipelineID string) ([]TeamPipeline, error) {
	var query struct {
		Node struct {
			Fragment struct {
				Teams struct {
					Edges []struct {
						Node TeamPipeline
					}
				} `graphql:"teams(first: 20)"`
			} `graphql:"... on Pipeline"`
		} `graphql:"node(id: $id)"`
	}
	vars := map[string]interface{}{
		"id": pipelineID,
	}
	if err := c.gqlClient.Query(context.TODO(), &query, vars); err != nil {
		return nil, err
	}

	result := []TeamPipeline{}
	for _, edge := range query.Node.Fragment.Teams.Edges {
		result = append(result, edge.Node)
	}
	return result, nil
}

// ReadTeamPipeline returns a PipelineTeam based on its ID.
func (c *Client) ReadTeamPipeline(id string) (*TeamPipeline, error) {
	var query struct {
		Node struct {
			TeamPipeline TeamPipeline `graphql:"... on TeamPipeline"`
		} `graphql:"node(id: $id)"`
	}
	vars := map[string]interface{}{
		"id": id,
	}

	if err := c.gqlClient.Query(context.TODO(), &query, vars); err != nil {
		return nil, err
	}
	return &query.Node.TeamPipeline, nil
}

// CreateTeamPipeline creates the provided PipelineTeam.
func (c *Client) CreateTeamPipeline(tp *TeamPipeline) error {
	var mutation struct {
		TeamPipelineCreate struct {
			TeamPipelineEdge struct {
				Node TeamPipeline
			}
		} `graphql:"teamPipelineCreate(input: $input)"`
	}
	type TeamPipelineCreateInput struct {
		TeamID      string `json:"teamID"`
		PipelineID  string `json:"pipelineID"`
		AccessLevel string `json:"accessLevel"`
	}
	vars := map[string]interface{}{
		"input": TeamPipelineCreateInput{
			TeamID:      string(tp.Team.ID),
			PipelineID:  string(tp.Pipeline.ID),
			AccessLevel: string(tp.AccessLevel),
		},
	}

	if err := c.gqlClient.Mutate(context.TODO(), &mutation, vars); err != nil {
		return err
	}
	tp.ID = mutation.TeamPipelineCreate.TeamPipelineEdge.Node.ID
	return nil
}

// UpdateTeamPipeline updates the provided PipelineTeam.
func (c *Client) UpdateTeamPipeline(tp *TeamPipeline) error {
	var mutation struct {
		TeamPipelineUpdate struct {
			TeamPipeline TeamPipeline
		} `graphql:"teamPipelineUpdate(input: $input)"`
	}
	type TeamPipelineUpdateInput struct {
		ID          string `json:"id"`
		AccessLevel string `json:"accessLevel"`
	}
	vars := map[string]interface{}{
		"input": TeamPipelineUpdateInput{
			ID:          string(tp.ID),
			AccessLevel: string(tp.AccessLevel),
		},
	}

	if err := c.gqlClient.Mutate(context.TODO(), &mutation, vars); err != nil {
		return err
	}

	*tp = mutation.TeamPipelineUpdate.TeamPipeline
	return nil
}

// DeleteTeamPipeline deletes the provided PipelineTeam.
func (c *Client) DeleteTeamPipeline(tp *TeamPipeline) error {
	var mutation struct {
		TeamPipelineDelete struct {
			DeletedTeamPipelineID graphql.String `graphql:"deletedTeamPipelineID"`
		} `graphql:"teamPipelineDelete(input: $input)"`
	}
	type TeamPipelineDeleteInput struct {
		ID    string `json:"id"`
		Force bool   `json:"force"`
	}
	vars := map[string]interface{}{
		"input": TeamPipelineDeleteInput{
			ID:    string(tp.ID),
			Force: false,
		},
	}
	return c.gqlClient.Mutate(context.TODO(), &mutation, vars)
}
