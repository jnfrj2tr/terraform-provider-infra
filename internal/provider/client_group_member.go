package provider

import (
	"fmt"
	"net/http"
)

type GroupMember struct {
	ID      string `json:"id"`
	GroupID string `json:"groupID"`
	UserID  string `json:"userID"`
}

type AddGroupMemberRequest struct {
	UserID string `json:"userID"`
}

func addGroupMember(client *http.Client, host, token, groupID, userID string) (*GroupMember, error) {
	body := AddGroupMemberRequest{UserID: userID}
	var member GroupMember
	err := doRequest(client, http.MethodPost, host, fmt.Sprintf("/api/groups/%s/users", groupID), token, body, &member)
	if err != nil {
		return nil, fmt.Errorf("adding group member: %w", err)
	}
	return &member, nil
}

func readGroupMember(client *http.Client, host, token, groupID, userID string) (*GroupMember, error) {
	var members []GroupMember
	err := doRequest(client, http.MethodGet, host, fmt.Sprintf("/api/groups/%s/users", groupID), token, nil, &members)
	if err != nil {
		return nil, fmt.Errorf("reading group members: %w", err)
	}
	for _, m := range members {
		if m.UserID == userID {
			return &m, nil
		}
	}
	return nil, fmt.Errorf("group member %s not found in group %s", userID, groupID)
}

func removeGroupMember(client *http.Client, host, token, groupID, userID string) error {
	err := doRequest(client, http.MethodDelete, host, fmt.Sprintf("/api/groups/%s/users/%s", groupID, userID), token, nil, nil)
	if err != nil {
		return fmt.Errorf("removing group member: %w", err)
	}
	return nil
}
