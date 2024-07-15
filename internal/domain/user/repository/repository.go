package repository

import (
	"github.com/azka-zaydan/synapsis-test/infras"
)

type UserRepository interface {
}

type UserRepositoryMySQL struct {
	DB *infras.MySQLConn
}

func ProvideUserRepositoryMySQL(conn *infras.MySQLConn) *UserRepositoryMySQL {
	return &UserRepositoryMySQL{
		DB: conn,
	}
}
