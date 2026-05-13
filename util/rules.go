package util

import (
	"fmt"
)

// UserSelfOnly restricts access strictly to the human user who owns the record (field 'user').
func UserSelfOnly() string {
	return "@request.auth.collectionName = 'users' && user = @request.auth.id"
}

// WorkspaceAnyMember allows any verified member of the workspace context to perform an action.
// If scope is provided, API Keys with that granular scope are also allowed.
func WorkspaceAnyMember(scope string) string {
	userPart := "(@request.auth.collectionName = 'users' && workspace = @request.auth.activeWorkspace && @collection.workspaceMembers.workspace ?= workspace && @collection.workspaceMembers.user ?= @request.auth.id)"

	if scope == "" {
		return userPart
	}

	apiKeyPart := fmt.Sprintf("(@request.auth.collectionName = 'apiKeys' && workspace = @request.auth.workspace.id && @request.auth.scopes ~ '%s')", scope)
	return userPart + " || " + apiKeyPart
}

// WorkspaceSelfOnly requires workspace membership AND restricts the action to the record user.
func WorkspaceSelfOnly(scope string) string {
	userPart := "(@request.auth.collectionName = 'users' && workspace = @request.auth.activeWorkspace && @collection.workspaceMembers.workspace ?= workspace && @collection.workspaceMembers.user ?= @request.auth.id && user = @request.auth.id)"

	if scope == "" {
		return userPart
	}

	apiKeyPart := fmt.Sprintf("(@request.auth.collectionName = 'apiKeys' && workspace = @request.auth.workspace.id && @request.auth.scopes ~ '%s')", scope)
	return userPart + " || " + apiKeyPart
}

// WorkspaceAnyAdmin restricts an action to users with an 'admin' role in the workspace.
// targetField is the field name holding the workspace ID (e.g., 'id' for the workspace record).
func WorkspaceAnyAdmin(scope string, targetField string) string {
	userPart := fmt.Sprintf("(@request.auth.collectionName = 'users' && @collection.workspaceMembers.workspace ?= %s && @collection.workspaceMembers.user ?= @request.auth.id && @collection.workspaceMembers.role ?= 'admin')", targetField)

	if scope == "" {
		return userPart
	}

	apiKeyPart := fmt.Sprintf("(@request.auth.collectionName = 'apiKeys' && %s = @request.auth.workspace.id && @request.auth.scopes ~ '%s')", targetField, scope)
	return userPart + " || " + apiKeyPart
}

// WorkspaceAdminSelfOnly requires the user to be a workspace admin AND the user of the record.
func WorkspaceAdminSelfOnly(scope string, targetField string) string {
	userPart := fmt.Sprintf("(@request.auth.collectionName = 'users' && user = @request.auth.id && %s = @request.auth.activeWorkspace && @request.auth.activeRole = 'admin')", targetField)

	if scope == "" {
		return userPart
	}

	apiKeyPart := fmt.Sprintf("(@request.auth.collectionName = 'apiKeys' && %s = @request.auth.workspace.id && @request.auth.scopes ~ '%s')", targetField, scope)
	return userPart + " || " + apiKeyPart
}
