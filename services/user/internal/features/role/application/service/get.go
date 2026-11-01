package roleservice

import (
	"context"
	"fmt"
	"user/internal/core/domain"
)

func (s *RoleService) GetRoleById(
	ctx context.Context,
	roleId string,
) (domain.Role, error) {
	role, err := s.RoleRepository.GetById(ctx, roleId)
	if err != nil {
		return domain.Role{}, fmt.Errorf("get role by id=%s; %w", roleId, err)
	}
	return role, nil
}

func (s *RoleService) GetRoleByName(
	ctx context.Context,
	roleName string,
) (domain.Role, error) {
	role, err := s.RoleRepository.GetById(ctx, roleName)
	if err != nil {
		return domain.Role{}, fmt.Errorf("get role by name=%s; %w", roleName, err)
	}
	return role, nil
}
