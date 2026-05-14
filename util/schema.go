package util

type collectionSchema struct {
	Workspaces       string
	WorkspaceMembers string
	Users            string
	ApiKeys          string
	Documents        string
	Records          string
	AuditLogs        string
}

type recordFields struct {
	Key, Value, User, Workspace, Label, Type, Format, Created, Updated string
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

type auditLogFields struct {
	User, Action, Collection, RecordId, OldData, NewData, Ip, UserAgent, Workspace string
}

var Coll = collectionSchema{
	Workspaces:       "workspaces",
	WorkspaceMembers: "workspaceMembers",
	Users:            "users",
	ApiKeys:          "apiKeys",
	Documents:        "documents",
	Records:          "records",
	AuditLogs:        "auditLogs",
}

var Fields = struct {
	Workspace       workspaceFields
	User            userFields
	WorkspaceMember memberFields
	ApiKey          apiKeyFields
	Record          recordFields
	Document        documentFields
	AuditLog        auditLogFields
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
	Record: recordFields{
		Key:       "key",
		Value:     "value",
		Type:      "type",
		Format:    "format",
		Label:     "label",
		Workspace: "workspace",
		User:      "user",
		Created:   "created",
		Updated:   "updated",
	},
	AuditLog: auditLogFields{
		User:       "user",
		Action:     "action",
		Collection: "collection",
		RecordId:   "recordId",
		OldData:    "oldData",
		NewData:    "newData",
		Ip:         "ip",
		UserAgent:  "userAgent",
		Workspace:  "workspace",
	},
}
