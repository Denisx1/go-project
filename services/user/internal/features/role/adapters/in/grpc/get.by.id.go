package rolegrpchandler

import (
	"context"
	"fmt"
	userv1 "go_project/contracts/gen/go/user"
)

func (h RoleGrpcHandler) GetById(
	ctx context.Context,
	req *userv1.GetByIdRequest,
) (*userv1.GetRoleResponse, error) {
	role, err := h.RoleService.GetRoleById(ctx, req.Id)
	if err != nil {
		return nil, fmt.Errorf("create role: %w", err)
	}
	return &userv1.GetRoleResponse{
		Id:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		Permissions: role.Permissions,
	}, nil
}
