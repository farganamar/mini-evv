package service

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
	"github.com/rs/zerolog/log"
	model "github.com/zorahealth/user-service/internal/model/user"
	"github.com/zorahealth/user-service/internal/model/user/dto"
)

type UserSvc interface {
	GetUserByID(ctx context.Context, userID uuid.UUID) (model.User, error)
	CreateUser(ctx context.Context, user dto.CreateNewUser) (model.User, error)
	VerifyUser(ctx context.Context, userID uuid.UUID) error
	Login(ctx context.Context, email string, password string) (dto.LoginResponse, error)
}

func (s *UserServiceImpl) GetUserByID(ctx context.Context, userID uuid.UUID) (model.User, error) {
	return s.UserRepository.GetUser(ctx, userID, nil)
}

func (s *UserServiceImpl) VerifyUser(ctx context.Context, userID uuid.UUID) error {
	isVerified := true
	verifyAt := time.Now()
	err := s.UserRepository.UpdateUser(ctx, model.User{
		ID:         userID,
		IsVerified: &isVerified,
		VerifyAt:   &verifyAt,
	}, nil)

	if err != nil {
		log.Error().Err(err).Msg("[user_service/VerifyUser] failed to update user")
		return err
	}
	return nil
}
