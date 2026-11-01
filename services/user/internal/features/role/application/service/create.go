package roleservice

import (
	"context"
	"fmt"
	"user/internal/core/domain"
)

func (s *RoleService) CreateRole(
	ctx context.Context,
	role domain.Role,
) (string, error) {
	roleId, err := s.RoleRepository.CreateRole(ctx, role)
	if err != nil {
		return "", fmt.Errorf("create role: %w", err)
	}
	return roleId, nil
}
