package repository

import (
	"context"
	"database/sql"

	"github.com/azka-zaydan/synapsis-test/internal/domain/user/model"
	"github.com/azka-zaydan/synapsis-test/shared/failure"
	"github.com/azka-zaydan/synapsis-test/shared/logger"
)

type UserRepo interface {
	FindByUsername(ctx context.Context, username string) (res model.User, err error)
	CreateUser(ctx context.Context, user *model.User) (err error)
}

func (repo *UserRepositoryMySQL) FindByUsername(ctx context.Context, username string) (res model.User, err error) {

	exist, err := repo.doesUserExist(ctx, username)
	if err != nil {
		if err == sql.ErrNoRows {
			err = failure.NotFound("user")
			return
		}
		logger.ErrorWithStack(err)
		return
	}
	if !exist {
		err = failure.NotFound("user")
		return
	}

	query := "SELECT id,username,password_hash,created_by,meta_created_at,updated_by,meta_updated_at FROM user WHERE username = ?"

	err = repo.DB.Read.GetContext(ctx, &res, query, username)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	return
}

func (repo *UserRepositoryMySQL) doesUserExist(ctx context.Context, username string) (exist bool, err error) {
	query := `SELECT COUNT(id) FROM user WHERE username = ?`
	var count int
	err = repo.DB.Read.GetContext(ctx, &count, query, username)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return count > 0, nil
}

func (repo *UserRepositoryMySQL) CreateUser(ctx context.Context, user *model.User) (err error) {
	_, err = repo.DB.Write.NamedExecContext(ctx, userInsertQuery, user)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

// queries
var (
	userInsertQuery = `
	INSERT INTO user (id, username, password_hash, created_by, meta_created_at, updated_by)
	VALUES (:id, :username, :password_hash, :created_by, :meta_created_at, :updated_by)`
)
