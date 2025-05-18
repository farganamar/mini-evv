package service

import (
	"context"

	"github.com/farganamar/evv-service/internal/model/v1/appointment/dto"
	"github.com/farganamar/evv-service/internal/service"
)

type AppointmentService interface {
	GetAppointmentsByUserId(ctx context.Context, arg dto.GetAppointmentsByUserIdRequest) ([]dto.GetAppointmentsByUserIdResponse, error)
	UpdateAppointmentStatus(ctx context.Context, arg dto.UpdateAppointmentStatusRequest) error
}

type AppointmentServiceImpl struct {
	BaseService *service.ServiceImpl
}

func NewAppointmentService(baseService *service.ServiceImpl) *AppointmentServiceImpl {
	return &AppointmentServiceImpl{
		BaseService: baseService,
	}
}
