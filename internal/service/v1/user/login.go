package service

import (
	"context"
	"time"

	"github.com/farganamar/evv-service/helpers/auth"
	"github.com/farganamar/evv-service/helpers/failure"
	model "github.com/farganamar/evv-service/internal/model/v1/user"
	"github.com/farganamar/evv-service/internal/model/v1/user/dto"
	"github.com/guregu/null/v5"
)

func (s *UserServiceImpl) Login(ctx context.Context, arg dto.LoginRequest) (dto.LoginResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// tx, err := s.BaseService.BaseRepository.BeginTx(ctx)
	// if err != nil {
	// 	return dto.LoginResponse{}, err
	// }

	// defer func() {
	// 	if tx != nil {
	// 		if rollbackErr := tx.Rollback(); rollbackErr != nil && rollbackErr != sql.ErrTxDone {
	// 			log.Error().Err(rollbackErr).Msg("[login] failed to rollback transaction")
	// 			if err == nil {
	// 				err = rollbackErr
	// 			}
	// 		}
	// 	}
	// }()

	user, err := s.BaseService.UserRepository.FindUser(ctx, model.User{
		Username: null.StringFrom(arg.Username),
	}, nil)

	if err != nil {
		return dto.LoginResponse{}, err
	}

	if user.ID.IsNil() {
		return dto.LoginResponse{}, failure.NotFound("user not found")
	}

	// Generate JWT token
	token, err := s.BaseService.AuthService.GenerateTokenPair(ctx, &auth.ParamsGenerateToken{
		UserID:   user.ID.String(),
		Username: user.Username.String,
		Email:    user.Email.String,
		Roles:    []string{user.Roles},
	})

	if err != nil {
		return dto.LoginResponse{}, err
	}

	// if err := tx.Commit(); err != nil {
	// 	return dto.LoginResponse{}, err
	// }

	return dto.LoginResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiresAt:    null.TimeFrom(token.ExpiresAt),
		IssuedAt:     null.TimeFrom(token.IssuedAt),
	}, nil
}
