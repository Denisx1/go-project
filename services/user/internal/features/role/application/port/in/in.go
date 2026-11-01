package roleportin

import (
	"context"
	"user/internal/core/domain"
)

type RoleService interface {
	CreateRole(
		ctx context.Context,
		role domain.Role,
	) (string, error)
	GetRoleById(
		ctx context.Context,
		roleId string,
	) (domain.Role, error)
	GetRoleByName(
		ctx context.Context,
		roleName string,
	) (domain.Role, error)
	DeleteRole(
		ctx context.Context,
		roleId string,
	) error
}
