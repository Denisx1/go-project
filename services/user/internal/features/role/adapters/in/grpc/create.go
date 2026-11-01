package rolegrpchandler

import (
	"context"
	"fmt"
	userv1 "go_project/contracts/gen/go/user"
	"go_project/infrastructure/grpc/server/interceptors"
	"user/internal/core/domain"
)

func (h RoleGrpcHandler) CreateRole(
	ctx context.Context,
	req *userv1.CreateRoleRequest,
) (*userv1.CreateRoleResponse, error) {
	requestId := ctx.Value(interceptors.RequestIdKey{}).(string)
	roleDomain := domain.NewRoleuninitialized(req.Name, req.Description, req.Permissions, requestId)
	roleId, err := h.RoleService.CreateRole(ctx, *roleDomain)
	if err != nil {
		return nil, fmt.Errorf("create role: %w", err)
	}
	return &userv1.CreateRoleResponse{
		Id: roleId,
	}, nil
}
