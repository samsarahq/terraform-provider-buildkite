package client

import (
	"context"

	"github.com/shurcooL/graphql"
)

// TeamMember represents a user's membership with a team.
type TeamMember struct {
	ID     graphql.String `graphql:"id"`
	UserID graphql.String `graphql:"userID"`
	TeamID graphql.String `graphql:"teamID"`
}

func (c *Client) CreateTeamMember(member *TeamMember) error {
	var mutation struct {
		TeamMemberCreate struct {
			TeamMemberEdge struct {
				Node struct {
					ID graphql.String `graphql:"id"`
				} `graphql:"node"`
			} `graphql:"teamMemberEdge`
		} `graphql:"teamMemberCreate(input: $input)"`
	}
	type TeamMemberCreateInput struct {
		UserID string `json:"userID"`
		TeamID string `json:"teamID"`
	}
	vars := map[string]interface{}{
		"input": TeamMemberCreateInput{
			UserID: string(member.UserID),
			TeamID: string(member.TeamID),
		},
	}
	if err := c.gqlClient.Mutate(context.TODO(), &mutation, vars); err != nil {
		return err
	}
	member.ID = mutation.TeamMemberCreate.TeamMemberEdge.Node.ID
	return nil
}

func (c *Client) ReadTeamMember(id string) (*TeamMember, error) {
	var query struct {
		Node struct {
			TeamMember struct {
				User struct {
					ID graphql.String `graphql:"id"`
				} `graphql:"user"`
				Team struct {
					ID graphql.String `graphql:"id"`
				} `graphql:"team"`
			} `graphql:"... on TeamMember"`
		} `graphql:"node(id: $id)"`
	}
	vars := map[string]interface{}{
		"id": id,
	}
	if err := c.gqlClient.Query(context.TODO(), &query, vars); err != nil {
		return nil, err
	}
	member := &TeamMember{
		ID:     graphql.String(id),
		UserID: query.Node.TeamMember.User.ID,
		TeamID: query.Node.TeamMember.Team.ID,
	}
	return member, nil
}

func (c *Client) DeleteTeamMember(member *TeamMember) error {
	var mutation struct {
		TeamMemberDelete struct {
			DeletedTeamMemberID graphql.String `graphql:"deletedTeamMemberID"`
		} `graphql:"teamMemberDelete(input: $input)"`
	}
	type TeamMemberDeleteInput struct {
		ID string `json:"id"`
	}
	vars := map[string]interface{}{
		"input": TeamMemberDeleteInput{
			ID: string(member.ID),
		},
	}
	return c.gqlClient.Mutate(context.TODO(), &mutation, vars)
}
