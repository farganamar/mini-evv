package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"math"
	"time"

	"github.com/farganamar/evv-service/helpers/failure"
	model "github.com/farganamar/evv-service/internal/model/v1/appointment"
	"github.com/farganamar/evv-service/internal/model/v1/appointment/dto"
	modelAppointmentLog "github.com/farganamar/evv-service/internal/model/v1/appointment_log"
	modelClient "github.com/farganamar/evv-service/internal/model/v1/client"
	"github.com/gofrs/uuid/v5"
	"github.com/guregu/null/v5"
	"github.com/rs/zerolog/log"
)

func (s *AppointmentServiceImpl) UpdateAppointmentStatus(ctx context.Context, arg dto.UpdateAppointmentStatusRequest) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	tx, err := s.BaseService.BaseRepository.BeginTx(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if tx != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil && rollbackErr != sql.ErrTxDone {
				log.Error().Err(rollbackErr).Msg("[login] failed to rollback transaction")
				if err == nil {
					err = rollbackErr
				}
			}
		}
	}()

	// Chck if appointment exists and is not cancelled and not completed and not expired
	appointment, err := s.BaseService.AppointmentRepository.GetAppointmentDetail(ctx, model.Appointment{
		ID:          uuid.FromStringOrNil(arg.AppointmentId),
		CaregiverID: uuid.FromStringOrNil(arg.UserID),
	}, tx)

	if err != nil {
		return err
	}

	if err := s.validateAppointment(ctx, appointment, arg); err != nil {
		return err
	}

	client, err := s.BaseService.ClientRepository.GetClientDetail(ctx, modelClient.Client{
		ID: uuid.FromStringOrNil(appointment.ClientID.String()),
	}, tx)

	if err != nil {
		return err
	}

	// Update appointment status
	appointment.Status = null.StringFrom(arg.Status)
	appointment.UpdatedAt = null.TimeFrom(time.Now())

	if arg.Status == model.StatusInProgress {
		// Validate latitude and longitude
		if err := s.validateLatLong(ctx, client, arg.Latitude, arg.Longitude); err != nil {
			return err
		}
	}

	// Update appointment in the database
	if err := s.BaseService.AppointmentRepository.UpdateAppointmentStatus(ctx, appointment, tx); err != nil {
		return err
	}

	if err := s.saveAppointmentLog(ctx, appointment, arg, tx); err != nil {
		log.Error().Err(err).Msg("[saveAppointmentLog] failed to save appointment log")
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *AppointmentServiceImpl) validateAppointment(ctx context.Context, appointment model.Appointment, arg dto.UpdateAppointmentStatusRequest) error {
	if appointment.ID.IsNil() {
		return failure.NotFound("appointment not found")
	}

	if appointment.Status.String == model.StatusCancelled {
		return errors.New("appointment cancelled")
	}

	if appointment.Status.String == model.StatusCompleted {
		return errors.New("appointment completed")
	}

	// if appointment.StartTime.Time.Before(time.Now()) {
	// 	return errors.New("appointment expired")
	// }

	if arg.TypeOfNote != "" && appointment.Status.String != model.StatusInProgress {
		return errors.New("type of note is only allowed when status is IN_PROGRESS")
	}

	if arg.Status == model.StatusCompleted && appointment.Status.String != model.StatusInProgress {
		return errors.New("status must be IN_PROGRESS to complete the appointment")
	}

	if arg.VerificationCode != "" && appointment.VerificationCode.String != arg.VerificationCode {
		return errors.New("verification code not match")
	}

	return nil
}

func (s *AppointmentServiceImpl) validateLatLong(ctx context.Context, client modelClient.Client, caregiverLat, caregiverLong float64) error {
	// Constants for validation
	const (
		// Maximum allowed distance in meters
		maxDistance = 500.0
		// Fault tolerance radius in meters (additional allowed distance)
		faultTolerance = 100.0
		// Earth's radius in meters
		earthRadius = 6371000.0
	)

	// Skip validation if client coordinates are not set
	if !client.Latitude.Valid || !client.Longitude.Valid {
		log.Warn().
			Str("clientId", client.ID.String()).
			Msg("Client location not set, skipping location validation")
		return nil
	}

	// Skip validation if caregiver coordinates are not set
	if caregiverLat == 0 && caregiverLong == 0 {
		log.Warn().Msg("Caregiver location not provided, skipping location validation")
		return nil
	}

	// Convert latitude and longitude from degrees to radians
	clientLatRad := client.Latitude.Float64 * (math.Pi / 180.0)
	clientLongRad := client.Longitude.Float64 * (math.Pi / 180.0)
	caregiverLatRad := caregiverLat * (math.Pi / 180.0)
	caregiverLongRad := caregiverLong * (math.Pi / 180.0)

	// Haversine formula to calculate distance between two points on Earth
	dLat := caregiverLatRad - clientLatRad
	dLong := caregiverLongRad - clientLongRad
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(clientLatRad)*math.Cos(caregiverLatRad)*
			math.Sin(dLong/2)*math.Sin(dLong/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	distance := earthRadius * c

	// Check if distance is within acceptable range (with fault tolerance)
	if distance > (maxDistance + faultTolerance) {
		log.Warn().
			Float64("distance", distance).
			Float64("maxAllowed", maxDistance+faultTolerance).
			Float64("clientLat", client.Latitude.Float64).
			Float64("clientLong", client.Longitude.Float64).
			Float64("caregiverLat", caregiverLat).
			Float64("caregiverLong", caregiverLong).
			Msg("Caregiver is too far from client location")

		return errors.New("you are too far from the client's location. Please move closer and try again.")
	}

	log.Info().
		Float64("distance", distance).
		Float64("maxAllowed", maxDistance+faultTolerance).
		Msg("Location validation passed")

	return nil
}

func (s *AppointmentServiceImpl) saveAppointmentLog(ctx context.Context, appoitment model.Appointment, arg dto.UpdateAppointmentStatusRequest, tx *sql.Tx) error {
	var LogType string
	var logDataJSON json.RawMessage
	if appoitment.Status.String == model.StatusInProgress {
		LogType = string(modelAppointmentLog.LogTypeCheckIn)
		data, err := json.Marshal(modelAppointmentLog.CheckInOutLogData{
			Device: arg.MetadataDevice.Device,
			IP:     arg.MetadataDevice.IP,
		})
		if err != nil {
			return err
		}
		logDataJSON = data
	}

	if appoitment.Status.String == model.StatusCompleted {
		LogType = string(modelAppointmentLog.LogTypeCheckOut)
		data, err := json.Marshal(modelAppointmentLog.CheckInOutLogData{
			Device: arg.MetadataDevice.Device,
			IP:     arg.MetadataDevice.IP,
		})
		if err != nil {
			return err
		}
		logDataJSON = data
	}

	if arg.TypeOfNote != "" {
		LogType = string(modelAppointmentLog.LogTypeNote)
		data, err := json.Marshal(modelAppointmentLog.NoteLogData{
			Type: arg.TypeOfNote,
		})
		if err != nil {
			return err
		}

		logDataJSON = data
	}

	// Save appointment log
	appointmentLog := modelAppointmentLog.AppointmentLog{
		ID:            uuid.Must(uuid.NewV7AtTime(time.Now())),
		AppointmentID: appoitment.ID,
		CaregiverID:   uuid.FromStringOrNil(arg.UserID),
		LogType:       null.StringFrom(LogType),
		LogData:       string(logDataJSON),
		CreatedAt:     null.TimeFrom(time.Now()),
		Notes:         null.StringFrom(arg.Note),
		Latitude:      null.FloatFrom(arg.Latitude),
		Longitude:     null.FloatFrom(arg.Longitude),
		Timestamp:     null.TimeFrom(time.Now()),
	}

	if err := s.BaseService.AppointmentLogRepository.CreateLog(ctx, appointmentLog, tx); err != nil {
		return err
	}

	return nil

}
