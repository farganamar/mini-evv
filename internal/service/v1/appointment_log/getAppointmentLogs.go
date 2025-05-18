package service

import (
	"context"

	model "github.com/farganamar/evv-service/internal/model/v1/appointment_log"
	"github.com/farganamar/evv-service/internal/model/v1/appointment_log/dto"
	"github.com/gofrs/uuid/v5"
)

func (s *AppointmentLogServiceImpl) GetAppointmentLogs(ctx context.Context, arg dto.GetAppointmentLogsRequest) ([]dto.GetAppointmentLogsResponse, error) {
	// Call the repository method to get appointment logs
	appointmentLogs, err := s.BaseService.AppointmentLogRepository.GetAppointmentLogByIdandUserId(ctx, model.AppointmentLog{
		AppointmentID: uuid.FromStringOrNil(arg.AppointmentId),
		CaregiverID:   uuid.FromStringOrNil(arg.UserId),
	}, nil)
	if err != nil {
		return nil, err
	}

	// Map the appointment logs to the response DTO
	response := make([]dto.GetAppointmentLogsResponse, len(appointmentLogs))
	for i, log := range appointmentLogs {
		response[i] = dto.GetAppointmentLogsResponse{
			AppointmentID: log.AppointmentID.String(),
			CaregiverID:   log.CaregiverID.String(),
			LogType:       log.LogType.String,
			LogData:       log.LogDataJSON,
			Latitude:      log.Latitude,
			Longitude:     log.Longitude,
			Timestamp:     log.Timestamp,
			Notes:         log.Notes,
		}
	}

	return response, nil
}
