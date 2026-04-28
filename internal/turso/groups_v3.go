package turso

import (
	"fmt"
	"net/http"
)

type GroupsV3Client client

func (g *GroupsV3Client) url(orgID, suffix string) string {
	return "/v3/organizations/" + orgID + "/groups" + suffix
}

func (g *GroupsV3Client) List(orgID string) ([]Group, error) {
	r, err := g.client.Get(g.url(orgID, ""), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list groups: %w", err)
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to list groups: %w", parseResponseError(r))
	}

	type response struct {
		Groups []Group `json:"groups"`
	}
	resp, err := unmarshal[response](r)
	if err != nil {
		return nil, err
	}
	return resp.Groups, nil
}

func (g *GroupsV3Client) Get(orgID, groupID string) (Group, error) {
	r, err := g.client.Get(g.url(orgID, "/"+groupID), nil)
	if err != nil {
		return Group{}, fmt.Errorf("failed to get group: %w", err)
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		return Group{}, fmt.Errorf("failed to get group: %w", parseResponseError(r))
	}

	type response struct {
		Group Group `json:"group"`
	}
	resp, err := unmarshal[response](r)
	if err != nil {
		return Group{}, err
	}
	return resp.Group, nil
}
