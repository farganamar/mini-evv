package repository

import (
	"context"
	"database/sql"

	model "github.com/farganamar/evv-service/internal/model/v1/appointment_log"
	"github.com/farganamar/evv-service/internal/repository"
)

type AppointmentLogRepoInterface interface {
	GetAppointmentLogByIdandUserId(ctx context.Context, arg model.AppointmentLog, tx *sql.Tx) ([]model.AppointmentLog, error)
	CreateLog(ctx context.Context, arg model.AppointmentLog, tx *sql.Tx) error
}

type AppointmentLogRepositoryImpl struct {
	*repository.RepositoryImpl
}

func NewAppointmentLogRepository(db *repository.RepositoryImpl) *AppointmentLogRepositoryImpl {
	return &AppointmentLogRepositoryImpl{db}
}
