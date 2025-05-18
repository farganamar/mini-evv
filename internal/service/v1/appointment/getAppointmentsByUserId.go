package service

import (
	"context"
	"time"

	model "github.com/farganamar/evv-service/internal/model/v1/appointment"
	"github.com/farganamar/evv-service/internal/model/v1/appointment/dto"
	"github.com/gofrs/uuid/v5"
	"github.com/guregu/null/v5"
)

func (s *AppointmentServiceImpl) GetAppointmentsByUserId(ctx context.Context, arg dto.GetAppointmentsByUserIdRequest) ([]dto.GetAppointmentsByUserIdResponse, error) {
	res := make([]dto.GetAppointmentsByUserIdResponse, 0)
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	appointments, err := s.BaseService.AppointmentRepository.GetAppointmentsByUserId(ctx, model.Appointment{
		CaregiverID: uuid.FromStringOrNil(arg.UserId),
		Status:      null.StringFrom(arg.Status),
	}, nil)
	if err != nil {
		return res, err
	}

	for _, appointment := range appointments {
		res = append(res, dto.GetAppointmentsByUserIdResponse{
			AppointmentId: appointment.ID.String(),
			CaregiverId:   uuid.NullUUID{UUID: appointment.CaregiverID, Valid: true},
			ClientId:      uuid.NullUUID{UUID: appointment.ClientID, Valid: true},
			StartTime:     appointment.StartTime,
			EndTime:       appointment.EndTime,
			Status:        appointment.Status,
			Notes:         appointment.Notes,
			CreatedAt:     appointment.CreatedAt,
			UpdatedAt:     appointment.UpdatedAt,
		})
	}

	return res, nil
}
