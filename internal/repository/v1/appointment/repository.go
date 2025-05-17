package repository

import (
	"context"
	"database/sql"

	model "github.com/farganamar/evv-service/internal/model/v1/appointment"
	"github.com/farganamar/evv-service/internal/repository"
)

type AppointmentRepoInterface interface {
	GetAppointmentsByUserId(ctx context.Context, arg model.Appointment, tx *sql.Tx) ([]model.Appointment, error)
}

type AppointmentRepositoryImpl struct {
	*repository.RepositoryImpl
}

func NewAppointmentRepository(db *repository.RepositoryImpl) *AppointmentRepositoryImpl {
	return &AppointmentRepositoryImpl{db}
}
