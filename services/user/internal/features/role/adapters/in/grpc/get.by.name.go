package rolegrpchandler

import (
	"context"
	"fmt"
	userv1 "go_project/contracts/gen/go/user"
)

func (h RoleGrpcHandler) GetByName(
	ctx context.Context,
	req *userv1.GetByNameRequest,
) (*userv1.GetRoleResponse, error) {
	role, err := h.RoleService.GetRoleByName(ctx, req.Name)
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
