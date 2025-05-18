package repository

import (
	"context"
	"database/sql"

	model "github.com/farganamar/evv-service/internal/model/v1/client"
	"github.com/farganamar/evv-service/internal/repository"
)

type ClientRepoInterface interface {
	GetClientDetail(ctx context.Context, arg model.Client, tx *sql.Tx) (model.Client, error)
}

type ClientRepositoryImpl struct {
	*repository.RepositoryImpl
}

func NewClientRepository(db *repository.RepositoryImpl) *ClientRepositoryImpl {
	return &ClientRepositoryImpl{db}
}
