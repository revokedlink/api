package util

import validation "github.com/go-ozzo/ozzo-validation/v4"

type AppError struct {
	ErrorCode string
	ErrorText string
}

var Errors = struct {
	DuplicateWorkspaceMember      AppError
	FailedToCreateWorkspaceMember AppError
	WorkspaceMemberLimitReached   AppError
	PersonalWorkspaceLimitReached AppError
	BusinessWorkspaceLimitReached AppError
	WorkspaceNotFound             AppError
	UserNotFound                  AppError
	ValidationFieldRequired       AppError
	ValidationFieldRestricted     AppError
	NotAuthorized                 AppError
	NotAuthenticated              AppError
	FailedToCreateRecord          AppError
	ForbiddenWorkspaceAccess      AppError
	DuplicateValues               AppError
	InvalidActiveWorkspace        AppError
	ActiveWorkspaceMismatch       AppError
}{
	DuplicateWorkspaceMember: AppError{
		ErrorCode: "duplicate_workspace_member",
		ErrorText: "This user is already a member of this workspace.",
	},
	WorkspaceMemberLimitReached: AppError{
		ErrorCode: "workspace_member_limit_reached",
		ErrorText: "This workspace reached the maximum number of members.",
	},
	FailedToCreateWorkspaceMember: AppError{
		ErrorCode: "failed_to_create_workspace_member",
		ErrorText: "Failed to add user to workspace.",
	},
	ValidationFieldRestricted: AppError{
		ErrorCode: "validation_insufficient_permissions",
		ErrorText: "You do not have permission to modify this field.",
	},
	ValidationFieldRequired: AppError{
		ErrorCode: "validation_required",
		ErrorText: "Missing required value.",
	},
	NotAuthorized: AppError{
		ErrorCode: "not_authorized",
		ErrorText: "You do not have the permission for this request.",
	},
	NotAuthenticated: AppError{
		ErrorCode: "not_authenticated",
		ErrorText: "You are not authenticated.",
	},
	FailedToCreateRecord: AppError{
		ErrorCode: "failed_to_create_record",
		ErrorText: "Failed to create record.",
	},
	DuplicateValues: AppError{
		ErrorCode: "duplicate_values",
		ErrorText: "Duplicate values are not allowed.",
	},
	PersonalWorkspaceLimitReached: AppError{
		ErrorCode: "personal_workspace_limit_reached",
		ErrorText: "This user has reached the maximum number of personal workspaces.",
	},
	BusinessWorkspaceLimitReached: AppError{
		ErrorCode: "business_workspace_limit_reached",
		ErrorText: "This user has reached the maximum number of business workspaces.",
	},
	WorkspaceNotFound: AppError{
		ErrorCode: "workspace_not_found",
		ErrorText: "Workspace not found.",
	},
	UserNotFound: AppError{
		ErrorCode: "user_not_found",
		ErrorText: "User not found.",
	},
	ForbiddenWorkspaceAccess: AppError{
		ErrorCode: "forbidden_workspace_access",
		ErrorText: "This user does not have access to this workspace.",
	},
	InvalidActiveWorkspace: AppError{
		ErrorCode: "invalid_active_workspace",
		ErrorText: "The selected active workspace is invalid.",
	},
	ActiveWorkspaceMismatch: AppError{
		ErrorCode: "active_workspace_mismatch",
		ErrorText: "The selected active workspace is invalid.",
	},
}

func AsValidationError(appErr AppError) error {
	return validation.NewError(appErr.ErrorCode, appErr.ErrorText)
}
