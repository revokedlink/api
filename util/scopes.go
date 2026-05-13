package util

const (
	ScopeDocumentsRead   = "documents:read"
	ScopeDocumentsCreate = "documents:create"
	ScopeDocumentsUpdate = "documents:update"
	ScopeDocumentsDelete = "documents:delete"

	ScopeWorkspacesRead   = "workspaces:read"
	ScopeWorkspacesCreate = "workspaces:create"
	ScopeWorkspacesUpdate = "workspaces:update"
	ScopeWorkspacesDelete = "workspaces:delete"

	ScopeWorkspaceMembersRead   = "workspaceMembers:read"
	ScopeWorkspaceMembersCreate = "workspaceMembers:create"
	ScopeWorkspaceMembersUpdate = "workspaceMembers:update"
	ScopeWorkspaceMembersDelete = "workspaceMembers:delete"
)

var AllScopes = []string{
	ScopeDocumentsRead,
	ScopeDocumentsCreate,
	ScopeDocumentsUpdate,
	ScopeDocumentsDelete,
	ScopeWorkspacesRead,
	ScopeWorkspacesCreate,
	ScopeWorkspacesUpdate,
	ScopeWorkspacesDelete,
	ScopeWorkspaceMembersRead,
	ScopeWorkspaceMembersCreate,
	ScopeWorkspaceMembersUpdate,
	ScopeWorkspaceMembersDelete,
}
