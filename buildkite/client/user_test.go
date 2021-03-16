package client

import "testing"

func TestGetUser(t *testing.T) {
	u, err := cli.GetUser(userEmail)
	if err != nil {
		t.Errorf("Couldn't make query: %s", err)
	}
	if u.Name == "" {
		t.Errorf("Name was not expected, received %s", u.Name)
	}
	if u.ID == "" {
		t.Error("User ID was blank")
	}
}
