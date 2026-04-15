package provider

import (
	"fmt"
	"net/http"
)

type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Created  string `json:"created"`
	Modified string `json:"modified"`
}

func createUser(client *http.Client, host, token, email string) (*User, error) {
	body := map[string]string{"name": email}
	var user User
	err := doRequest(client, host, token, http.MethodPost, "/api/users", body, &user)
	if err != nil {
		return nil, fmt.Errorf("creating user: %w", err)
	}
	return &user, nil
}

func readUser(client *http.Client, host, token, id string) (*User, error) {
	var result struct {
		Items []User `json:"items"`
	}
	err := doRequest(client, host, token, http.MethodGet, fmt.Sprintf("/api/users?id=%s", id), nil, &result)
	if err != nil {
		return nil, fmt.Errorf("reading user: %w", err)
	}
	if len(result.Items) == 0 {
		return nil, nil
	}
	return &result.Items[0], nil
}

func deleteUser(client *http.Client, host, token, id string) error {
	err := doRequest(client, host, token, http.MethodDelete, fmt.Sprintf("/api/users/%s", id), nil, nil)
	if err != nil {
		return fmt.Errorf("deleting user: %w", err)
	}
	return nil
}
