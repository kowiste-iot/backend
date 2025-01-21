package resource

type ResourcePermission struct {
	ID          string              `json:"id,omitempty"`
	Name        string              `json:"name"`
	DisplayName string              `json:"displayName,omitempty"`
	Roles       map[string][]string `json:"roles,omitempty"`
}
