package auth

type Branch struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Attributes  map[string]string `json:"attributes,omitempty"`
}

const (
	AdminBranch     string = "admin"
	DefaultBranch   string = "default"
	UndefinedBranch string = "undefined"
)

func ForbiddenBranch() []string {
	return []string{
		AdminBranch,
		DefaultBranch,
		UndefinedBranch,
	}

}
