package repository

import (
	"context"
	"database/sql"

	model "github.com/farganamar/evv-service/internal/model/v1/user"
	"github.com/farganamar/evv-service/internal/repository"
)

type UserRepoInterface interface {
	FindUser(ctx context.Context, arg model.User, tx *sql.Tx) (model.User, error)
}

type UserRepositoryImpl struct {
	*repository.RepositoryImpl
}

func NewUserRepository(db *repository.RepositoryImpl) *UserRepositoryImpl {
	return &UserRepositoryImpl{db}
}
