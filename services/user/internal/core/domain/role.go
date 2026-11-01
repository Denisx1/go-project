package domain

type Role struct {
	ID          string
	Name        string
	Description string
	Permissions string
	TraceId     string
}

func NewRole(
	id string,
	name string,
	description string,
	permissions string,
	traceId string,
) *Role {
	return &Role{
		ID:          id,
		Name:        name,
		Description: description,
		Permissions: permissions,
		TraceId:     traceId,
	}
}
func NewRoleuninitialized(
	name string,
	description string,
	permissions string,
	traceId string,
) *Role {
	return &Role{
		ID:          UninitializedUUID,
		Name:        name,
		Description: description,
		Permissions: permissions,
		TraceId:     traceId,
	}
}
