package provider

import (
	"fmt"
	"net/http"
)

type GroupMember struct {
	ID      string `json:"id"`
	GroupID string `json:"groupId"`
	UserID  string `json:"userId"`
}

func addGroupMember(client *Client, groupID, userID string) (*GroupMember, error) {
	body := map[string]string{
		"userId": userID,
	}

	var member GroupMember
	err := doRequest(client, http.MethodPost, fmt.Sprintf("/api/groups/%s/users", groupID), body, &member)
	if err != nil {
		return nil, fmt.Errorf("error adding group member: %w", err)
	}

	member.GroupID = groupID
	member.UserID = userID
	return &member, nil
}

func readGroupMember(client *Client, groupID, userID string) (*GroupMember, error) {
	var members []GroupMember
	err := doRequest(client, http.MethodGet, fmt.Sprintf("/api/groups/%s/users", groupID), nil, &members)
	if err != nil {
		return nil, fmt.Errorf("error reading group member: %w", err)
	}

	for _, m := range members {
		if m.UserID == userID {
			m.GroupID = groupID
			return &m, nil
		}
	}

	// Return nil, nil when the member is not found so the caller can treat
	// a missing member as a signal to recreate the resource rather than error.
	return nil, nil
}

func removeGroupMember(client *Client, groupID, userID string) error {
	err := doRequest(client, http.MethodDelete, fmt.Sprintf("/api/groups/%s/users/%s", groupID, userID), nil, nil)
	if err != nil {
		return fmt.Errorf("error removing group member: %w", err)
	}
	return nil
}
