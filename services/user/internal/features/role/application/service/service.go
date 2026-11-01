package roleservice

import (
	roleportout "user/internal/features/role/application/port/out"
)

type RoleService struct {
	RoleRepository  roleportout.RoleRepository
	CacheRepository roleportout.CacheRepository
	TxManager       roleportout.TxManager
}

func NewRoleService(
	roleRepo roleportout.RoleRepository,
	roleCache roleportout.CacheRepository,
	txManager roleportout.TxManager,
) *RoleService {
	return &RoleService{
		RoleRepository:  roleRepo,
		CacheRepository: roleCache,
		TxManager:       txManager,
	}
}
