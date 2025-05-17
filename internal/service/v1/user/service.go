package service

import (
	"context"

	"github.com/farganamar/evv-service/internal/model/v1/user/dto"
	"github.com/farganamar/evv-service/internal/service"
)

type UserService interface {
	Login(ctx context.Context, arg dto.LoginRequest) (dto.LoginResponse, error)
}

type UserServiceImpl struct {
	BaseService *service.ServiceImpl
}

func NewUserService(baseService *service.ServiceImpl) *UserServiceImpl {
	return &UserServiceImpl{
		BaseService: baseService,
	}
}
