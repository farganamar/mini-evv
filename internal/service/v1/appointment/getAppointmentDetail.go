package service

import (
	"context"

	model "github.com/farganamar/evv-service/internal/model/v1/appointment"
	"github.com/farganamar/evv-service/internal/model/v1/appointment/dto"
	modelClient "github.com/farganamar/evv-service/internal/model/v1/client"
	"github.com/gofrs/uuid/v5"
)

func (s *AppointmentServiceImpl) GetAppointmentDetail(ctx context.Context, appointmentID string, userID string) (dto.GetAppointmentsByUserIdResponse, error) {
	// Get appointment detail
	appointment, err := s.BaseService.AppointmentRepository.GetAppointmentDetail(ctx, model.Appointment{
		ID:          uuid.FromStringOrNil(appointmentID),
		CaregiverID: uuid.FromStringOrNil(userID),
	}, nil)
	if err != nil {
		return dto.GetAppointmentsByUserIdResponse{}, err
	}

	// Get client detail
	client, err := s.BaseService.ClientRepository.GetClientDetail(ctx, modelClient.Client{
		ID: appointment.ClientID,
	}, nil)
	if err != nil {
		return dto.GetAppointmentsByUserIdResponse{}, err
	}

	return dto.GetAppointmentsByUserIdResponse{
		AppointmentId: appointment.ID.String(),
		CaregiverId:   uuid.NullUUID{UUID: appointment.CaregiverID, Valid: true},
		ClientId:      uuid.NullUUID{UUID: appointment.ClientID, Valid: true},
		StartTime:     appointment.StartTime,
		EndTime:       appointment.EndTime,
		Status:        appointment.Status,
		Notes:         appointment.Notes,
		CreatedAt:     appointment.CreatedAt,
		UpdatedAt:     appointment.UpdatedAt,
		ClientDetail: dto.ClientDetail{
			Name:      client.Name,
			Phone:     client.PhoneNumber,
			Latitude:  client.Latitude,
			Longitude: client.Longitude,
			Address:   client.Address,
			Note:      client.Notes,
		},
	}, nil
}
