package util

type collectionSchema struct {
	Workspaces       string
	WorkspaceMembers string
	Users            string
	ApiKeys          string
	Documents        string
}

type workspaceFields struct {
	Name, Slug, Type, Created, Updated string
}

type userFields struct {
	ActiveWorkspace, ActiveRole, Email, Verified, Avatar, Active string
}

type memberFields struct {
	User, Workspace, Role, Created, Updated string
}

type apiKeyFields struct {
	Token, User, Workspace, Scopes, Created, Updated string
}

type documentFields struct {
	Workspace, User, Title, Content, Created, Updated string
}

var Coll = collectionSchema{
	Workspaces:       "workspaces",
	WorkspaceMembers: "workspaceMembers",
	Users:            "users",
	ApiKeys:          "apiKeys",
	Documents:        "documents",
}

var Fields = struct {
	Workspace       workspaceFields
	User            userFields
	WorkspaceMember memberFields
	ApiKey          apiKeyFields
	Document        documentFields
}{
	Workspace: workspaceFields{
		Name:    "name",
		Slug:    "slug",
		Type:    "type",
		Created: "created",
		Updated: "updated",
	},
	User: userFields{
		ActiveWorkspace: "activeWorkspace",
		ActiveRole:      "activeRole",
		Email:           "email",
		Verified:        "verified",
		Avatar:          "avatar",
		Active:          "active",
	},
	WorkspaceMember: memberFields{
		User:      "user",
		Workspace: "workspace",
		Role:      "role",
		Created:   "created",
		Updated:   "updated",
	},
	ApiKey: apiKeyFields{
		Token:     "token",
		User:      "user",
		Workspace: "workspace",
		Scopes:    "scopes",
		Created:   "created",
		Updated:   "updated",
	},
	Document: documentFields{
		Workspace: "workspace",
		User:      "user",
		Title:     "title",
		Content:   "content",
		Created:   "created",
		Updated:   "updated",
	},
}
