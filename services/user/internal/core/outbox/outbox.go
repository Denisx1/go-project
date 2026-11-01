package useroutbox

import postgres_pool "go_project/infrastructure/repository/postgres/pool"

type UserOutbox struct {
	pool postgres_pool.Pool
}

func NewUserOutbox(pool postgres_pool.Pool) *UserOutbox {
	return &UserOutbox{
		pool: pool,
	}
}


