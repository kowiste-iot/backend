package command

type PermissionInput struct {
    Token     string
    Resource  string
    Action    string  
    TenantID  string 
    BranchID string
}