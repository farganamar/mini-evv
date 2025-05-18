package service

import (
	"context"

	"github.com/farganamar/evv-service/internal/model/v1/appointment_log/dto"
	"github.com/farganamar/evv-service/internal/service"
)

type AppointmentLogService interface {
	GetAppointmentLogs(ctx context.Context, arg dto.GetAppointmentLogsRequest) ([]dto.GetAppointmentLogsResponse, error)
}

type AppointmentLogServiceImpl struct {
	BaseService *service.ServiceImpl
}

func NewAppointmentLogService(baseService *service.ServiceImpl) *AppointmentLogServiceImpl {
	return &AppointmentLogServiceImpl{
		BaseService: baseService,
	}
}
