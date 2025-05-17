package service

import (
	"github.com/farganamar/evv-service/configs"
)

type ServiceInterface interface {
}

type UserServiceImpl struct {
	cfg *configs.Config
}

func NewUserService(cfg *configs.Config) *UserServiceImpl {
	s := new(UserServiceImpl)
	s.cfg = cfg
	return s
}
