package models

type WorkspaceModel struct {
	WorkspaceID   int    `json:"workspace_id"`
	DepartmentID  int    `json:"department_id"`
	WorkspaceName string `json:"workspace_name"`
}

type CreateWorkspaceModel struct {
	WorkspaceName string `json:"workspace_name" validate:"required"`
}

type DeleteWorkspaceModel struct {
	WorkspaceName string `json:"workspace_name" validate:"required"`
	WorkspaceID   int    `json:"workspace_id" validate:"required"`
}
