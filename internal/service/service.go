package service

import (
	"github.com/farganamar/evv-service/configs"
	"github.com/farganamar/evv-service/helpers/auth"
	"github.com/farganamar/evv-service/internal/repository"
	UserRepository "github.com/farganamar/evv-service/internal/repository/v1/user"
)

type ServiceInterface interface {
}

type ServiceImpl struct {
	Cfg            *configs.Config
	BaseRepository repository.RepoInterface
	UserRepository UserRepository.UserRepoInterface
	AuthService    *auth.TokenService
}

func NewService(
	cfg *configs.Config,
	baseRepo repository.RepoInterface,
	userRepo UserRepository.UserRepoInterface,
	authService *auth.TokenService,
) *ServiceImpl {
	s := new(ServiceImpl)
	s.Cfg = cfg
	s.BaseRepository = baseRepo
	s.UserRepository = userRepo
	s.AuthService = authService
	return s
}
