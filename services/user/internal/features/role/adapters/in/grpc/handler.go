package rolegrpchandler

import (
	userv1 "go_project/contracts/gen/go/user"
	roleportin "user/internal/features/role/application/port/in"
)

type RoleGrpcHandler struct {
	RoleService roleportin.RoleService
	userv1.RoleServiceServer
}

func NewCandidateGrpcHandler(roleService roleportin.RoleService) *RoleGrpcHandler {
	return &RoleGrpcHandler{
		RoleService: roleService,
	}
}
