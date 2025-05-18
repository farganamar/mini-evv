package service

import (
	"context"
	"sync"
	"time"

	model "github.com/farganamar/evv-service/internal/model/v1/appointment"
	"github.com/farganamar/evv-service/internal/model/v1/appointment/dto"
	modelClient "github.com/farganamar/evv-service/internal/model/v1/client"
	"github.com/gofrs/uuid/v5"
	"github.com/guregu/null/v5"
)

func (s *AppointmentServiceImpl) GetAppointmentsByUserId(ctx context.Context, arg dto.GetAppointmentsByUserIdRequest) ([]dto.GetAppointmentsByUserIdResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	appointments, err := s.BaseService.AppointmentRepository.GetAppointmentsByUserId(ctx, model.Appointment{
		CaregiverID: uuid.FromStringOrNil(arg.UserId),
		Status:      null.StringFrom(arg.Status),
	}, nil)
	if err != nil {
		return nil, err
	}

	// Use a wait group to handle concurrent client data fetching
	var wg sync.WaitGroup
	// Create response slice with exact capacity needed
	res := make([]dto.GetAppointmentsByUserIdResponse, len(appointments))
	// Error channel to collect errors from goroutines
	errChan := make(chan error, len(appointments))
	// Use a mutex to protect concurrent writes to the error channel
	var mu sync.Mutex

	// Process each appointment concurrently
	for i, appointment := range appointments {
		wg.Add(1)
		go func(i int, appt model.Appointment) {
			defer wg.Done()

			// Create the base response object
			response := dto.GetAppointmentsByUserIdResponse{
				AppointmentId: appt.ID.String(),
				CaregiverId:   uuid.NullUUID{UUID: appt.CaregiverID, Valid: true},
				ClientId:      uuid.NullUUID{UUID: appt.ClientID, Valid: true},
				StartTime:     appt.StartTime,
				EndTime:       appt.EndTime,
				Status:        appt.Status,
				Notes:         appt.Notes,
				CreatedAt:     appt.CreatedAt,
				UpdatedAt:     appt.UpdatedAt,
			}

			// Get client details if client ID is valid
			if !appt.ClientID.IsNil() {
				client, err := s.BaseService.ClientRepository.GetClientDetail(ctx, modelClient.Client{
					ID: appt.ClientID,
				}, nil)

				if err != nil {
					mu.Lock()
					errChan <- err
					mu.Unlock()
					return
				}

				// Add client details to response
				response.ClientDetail = dto.ClientDetail{
					Name:      client.Name,
					Phone:     client.PhoneNumber,
					Latitude:  client.Latitude,
					Longitude: client.Longitude,
					Address:   client.Address,
					Note:      client.Notes,
				}
			}

			// Assign to the pre-allocated slot
			res[i] = response
		}(i, appointment)
	}

	// Wait for all goroutines to complete
	wg.Wait()
	close(errChan)

	// Check if any errors occurred during concurrent processing
	if len(errChan) > 0 {
		return nil, <-errChan // Return the first error
	}

	return res, nil
}
