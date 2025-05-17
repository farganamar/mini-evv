package service

import (
	"github.com/zorahealth/user-service/configs"

	repository "github.com/zorahealth/user-service/internal/repository"
)

type UserServiceInterface interface {
	UserSvc
}

type UserServiceImpl struct {
	UserRepository repository.UserRepoInterface
	cfg            *configs.Config
}

func NewUserService(userRepository repository.UserRepoInterface, cfg *configs.Config) *UserServiceImpl {
	s := new(UserServiceImpl)
	s.UserRepository = userRepository
	s.cfg = cfg
	return s
}
