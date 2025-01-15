package restkc

import (
    "context"
    "fmt"
    "github.com/go-resty/resty/v2"
)

type Role struct {
    ID       string `json:"id"`
    Required bool   `json:"required"`
}

type Policy struct {
    Name        string `json:"name"`
    Description string `json:"description"`
    Roles       []Role `json:"roles"`
    FetchRoles  bool   `json:"fetchRoles"`
    Logic       string `json:"logic"`
}

func CreateRolePolicy(ctx context.Context, url, token, tenantID, clientID string, p Policy) (*string, error) {
    baseURL := url + "/admin/realms/%s/clients/%s/authz/resource-server/policy/role"
    endpoint := fmt.Sprintf(baseURL, tenantID, clientID)

    client := resty.New()
    
    var result struct {
        ID string `json:"id"`
    }
    
    resp, err := client.R().
        SetContext(ctx).
        SetHeader("Authorization", "Bearer "+token).
        SetHeader("Content-Type", "application/json").
        SetBody(p).
        SetResult(&result).
        Post(endpoint)

    if err != nil {
        return nil, fmt.Errorf("failed to send request: %w", err)
    }

    if resp.StatusCode() != 201 {
        return nil, fmt.Errorf("failed to create policy. Status: %d, Body: %s", 
            resp.StatusCode(), string(resp.Body()))
    }

    return &result.ID, nil
}