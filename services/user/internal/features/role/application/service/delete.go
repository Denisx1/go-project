package roleservice

import (
	"context"
	"fmt"
)

func (s *RoleService) DeleteRole(
	ctx context.Context,
	roleId string,
) error {
	if err := s.RoleRepository.DeleteOrder(ctx, roleId); err != nil {
		return fmt.Errorf("delete role by id=%s: %w", roleId, err)
	}
	return nil
}
