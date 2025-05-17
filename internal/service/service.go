package service

import (
	"github.com/farganamar/evv-service/configs"
	"github.com/farganamar/evv-service/helpers/auth"
	"github.com/farganamar/evv-service/internal/repository"
	AppointmentRepository "github.com/farganamar/evv-service/internal/repository/v1/appointment"
	UserRepository "github.com/farganamar/evv-service/internal/repository/v1/user"
)

type ServiceInterface interface {
}

type ServiceImpl struct {
	Cfg                   *configs.Config
	BaseRepository        repository.RepoInterface
	AuthService           *auth.TokenService
	UserRepository        UserRepository.UserRepoInterface
	AppointmentRepository AppointmentRepository.AppointmentRepoInterface
}

func NewService(
	cfg *configs.Config,
	baseRepo repository.RepoInterface,
	authService *auth.TokenService,
	userRepo UserRepository.UserRepoInterface,
	appointmentRepo AppointmentRepository.AppointmentRepoInterface,
) *ServiceImpl {
	s := new(ServiceImpl)
	s.Cfg = cfg
	s.BaseRepository = baseRepo
	s.UserRepository = userRepo
	s.AuthService = authService
	s.AppointmentRepository = appointmentRepo
	return s
}
