package service

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	modelAppointment "github.com/farganamar/evv-service/internal/model/v1/appointment"
	modelClient "github.com/farganamar/evv-service/internal/model/v1/client"
	"github.com/gofrs/uuid/v5"
	"github.com/guregu/null/v5"
	"github.com/rs/zerolog/log"
)

// CreateSeederAppointmentAndClient creates 10 clients and appointments with nearby locations
func (s *AppointmentServiceImpl) CreateSeederAppointmentAndClient(ctx context.Context, baseLat, baseLong float64) error {
	// Create a new transaction
	tx, err := s.BaseService.BaseRepository.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Error().Err(rbErr).Msg("Failed to rollback transaction")
			}
		}
	}()

	// Number of records to create
	const numRecords = 10

	// Create a local random generator instead of using rand.Seed
	// This is the recommended approach in Go 1.20+
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Generate nearby locations within ~50m radius
	// 0.0005 degrees is approximately 50m
	const latVariation = 0.0005
	const longVariation = 0.0005

	// Use existing caregiver IDs
	caregiverIDs := []uuid.UUID{
		uuid.FromStringOrNil("018e7677-1ee1-7000-8000-000000000001"),
		uuid.FromStringOrNil("018e7677-1ee1-7000-8000-000000000002"),
	}

	// Generate 10 clients and appointments
	for i := 0; i < numRecords; i++ {
		// Generate a nearby location
		clientLat := baseLat + (rng.Float64()*latVariation*2 - latVariation)
		clientLong := baseLong + (rng.Float64()*longVariation*2 - longVariation)

		// Create a client
		clientID, err := uuid.NewV7AtTime(time.Now())
		if err != nil {
			return fmt.Errorf("failed to generate UUID: %w", err)
		}

		client := modelClient.Client{
			ID:          clientID,
			Name:        null.StringFrom(fmt.Sprintf("Client%d", i+1)),
			Address:     null.StringFrom(fmt.Sprintf("%d Main St, City", 100+i*10)),
			PhoneNumber: null.StringFrom(fmt.Sprintf("555-%03d-%04d", i+1, 1000+rng.Intn(9000))),
			Latitude:    null.FloatFrom(clientLat),
			Longitude:   null.FloatFrom(clientLong),
			Notes:       null.StringFrom(fmt.Sprintf("Test client %d", i+1)),
			CreatedAt:   null.TimeFrom(time.Now()),
			UpdatedAt:   null.TimeFrom(time.Now()),
		}

		// Insert client
		createdClient, err := s.BaseService.ClientRepository.CreateClient(ctx, client, tx)
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		// Alternate between the two caregiver IDs
		caregiverID := caregiverIDs[i%len(caregiverIDs)]

		// Generate random start and end times in the next 7 days
		startDaysOffset := rng.Intn(7) + 1 // 1-7 days from now
		startHour := 8 + rng.Intn(8)       // Between 8 AM and 4 PM
		startTime := time.Now().AddDate(0, 0, startDaysOffset).
			Truncate(24 * time.Hour).                 // Start of day
			Add(time.Duration(startHour) * time.Hour) // Add hours

		// End time is 1-3 hours after start time
		durationHours := 1 + rng.Intn(2) // 1-3 hours
		endTime := startTime.Add(time.Duration(durationHours) * time.Hour)

		// Generate a 4-digit verification code
		verificationCode := fmt.Sprintf("%04d", 1000+rng.Intn(9000))

		// Create appointment
		appointmentID, err := uuid.NewV4()
		if err != nil {
			return fmt.Errorf("failed to generate appointment UUID: %w", err)
		}

		apt := modelAppointment.Appointment{
			ID:               appointmentID,
			ClientID:         createdClient.ID,
			CaregiverID:      caregiverID,
			StartTime:        null.TimeFrom(startTime),
			EndTime:          null.TimeFrom(endTime),
			Status:           null.StringFrom("SCHEDULED"),
			Notes:            null.StringFrom(fmt.Sprintf("Sample appointment %d", i+1)),
			VerificationCode: null.StringFrom(verificationCode),
			CreatedAt:        null.TimeFrom(time.Now()),
			UpdatedAt:        null.TimeFrom(time.Now()),
		}

		// Insert appointment
		_, err = s.BaseService.AppointmentRepository.CreateAppointment(ctx, apt, tx)
		if err != nil {
			return fmt.Errorf("failed to create appointment: %w", err)
		}

		log.Info().
			Str("client_id", clientID.String()).
			Str("appointment_id", appointmentID.String()).
			Str("caregiver_id", caregiverID.String()).
			Float64("latitude", clientLat).
			Float64("longitude", clientLong).
			Time("start_time", startTime).
			Time("end_time", endTime).
			Str("verification_code", verificationCode).
			Msgf("Created seeder record %d/%d", i+1, numRecords)
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Info().Msgf("Successfully created %d clients and appointments", numRecords)
	return nil
}
