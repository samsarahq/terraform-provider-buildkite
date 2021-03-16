package client

import (
	"context"
	"fmt"

	"github.com/shurcooL/graphql"
)

// Team represents a Buildkite team.
type Team struct {
	ID                graphql.String
	Name              graphql.String
	Privacy           graphql.String
	IsDefaultTeam     graphql.Boolean
	DefaultMemberRole graphql.String
}

// CreateTeam creates a given team and if successful, adds an ID to the given team.
func (c *Client) CreateTeam(team *Team) error {
	var mutation struct {
		TeamCreate struct {
			TeamEdge struct {
				Node Team
			}
		} `graphql:"teamCreate(input: $input)"`
	}
	type TeamCreateInput struct {
		OrganizationID    string `json:"organizationID"`
		Name              string `json:"name"`
		Privacy           string `json:"privacy"`
		IsDefaultTeam     bool   `json:"isDefaultTeam"`
		DefaultMemberRole string `json:"defaultMemberRole"`
	}
	vars := map[string]interface{}{
		"input": TeamCreateInput{
			OrganizationID:    c.orgID,
			Name:              string(team.Name),
			Privacy:           string(team.Privacy),
			IsDefaultTeam:     bool(team.IsDefaultTeam),
			DefaultMemberRole: string(team.DefaultMemberRole),
		},
	}

	if err := c.gqlClient.Mutate(context.TODO(), &mutation, vars); err != nil {
		return err
	}
	team.ID = mutation.TeamCreate.TeamEdge.Node.ID
	return nil
}

// ReadTeamByName uses the human readable name of the team to query for the team struct.
func (c *Client) ReadTeamByName(name string) (*Team, error) {
	var query struct {
		Team Team `graphql:"team(slug: $slug)"`
	}
	vars := map[string]interface{}{
		"slug": fmt.Sprintf("%s/%s", c.orgSlug, name),
	}
	if err := c.gqlClient.Query(context.TODO(), &query, vars); err != nil {
		return nil, err
	}
	return &query.Team, nil
}

// ReadTeam returns a team for a given gql ID.
func (c *Client) ReadTeam(id string) (*Team, error) {
	var query struct {
		Node struct {
			Team Team `graphql:"... on Team"`
		} `graphql:"node(id: $id)"`
	}
	vars := map[string]interface{}{
		"id": id,
	}

	if err := c.gqlClient.Query(context.TODO(), &query, vars); err != nil {
		return nil, err
	}
	return &query.Node.Team, nil
}

// UpdateTeam syncs the local team struct with Buildkite.
func (c *Client) UpdateTeam(team *Team) error {
	var mutation struct {
		TeamUpdate struct {
			Team Team
		} `graphql:"teamUpdate(input: $input)"`
	}
	type TeamUpdateInput struct {
		ID                string `json:"id"`
		Name              string `json:"name"`
		Privacy           string `json:"privacy"`
		IsDefaultTeam     bool   `json:"isDefaultTeam"`
		DefaultMemberRole string `json:"defaultMemberRole"`
	}
	vars := map[string]interface{}{
		"input": TeamUpdateInput{
			ID:                string(team.ID),
			Name:              string(team.Name),
			Privacy:           string(team.Privacy),
			IsDefaultTeam:     bool(team.IsDefaultTeam),
			DefaultMemberRole: string(team.DefaultMemberRole),
		},
	}

	if err := c.gqlClient.Mutate(context.TODO(), &mutation, vars); err != nil {
		return err
	}
	*team = mutation.TeamUpdate.Team
	return nil
}

// DeleteTeam deletes the given team based on the ID field.
func (c *Client) DeleteTeam(team *Team) error {
	var mutation struct {
		TeamDelete struct {
			DeletedTeamID graphql.String `graphql:"deletedTeamID"`
		} `graphql:"teamDelete(input: $input)"`
	}
	type TeamDeleteInput struct {
		ID string `json:"id"`
	}
	vars := map[string]interface{}{
		"input": TeamDeleteInput{
			ID: string(team.ID),
		},
	}
	return c.gqlClient.Mutate(context.TODO(), &mutation, vars)
}
