//go:build wireinject
// +build wireinject

package main

import (
	"github.com/farganamar/evv-service/configs"
	"github.com/farganamar/evv-service/helpers/auth"
	"github.com/farganamar/evv-service/infras"
	handlerV1 "github.com/farganamar/evv-service/internal/handlers/v1"
	BaseRepository "github.com/farganamar/evv-service/internal/repository"
	AppointmentRepository "github.com/farganamar/evv-service/internal/repository/v1/appointment"
	AppointmentLogRepository "github.com/farganamar/evv-service/internal/repository/v1/appointment_log"
	UserRepository "github.com/farganamar/evv-service/internal/repository/v1/user"
	BaseService "github.com/farganamar/evv-service/internal/service"
	AppointmentServiceV1 "github.com/farganamar/evv-service/internal/service/v1/appointment"
	AppointmentLogServiceV1 "github.com/farganamar/evv-service/internal/service/v1/appointment_log"
	UserServiceV1 "github.com/farganamar/evv-service/internal/service/v1/user"
	"github.com/farganamar/evv-service/transport/http"
	"github.com/farganamar/evv-service/transport/http/middleware"
	"github.com/farganamar/evv-service/transport/http/router"
	"github.com/google/wire"
)

// Wiring for configurations.
var configurationsServiceGen = wire.NewSet(
	configs.Get,
)

// Wiring for persistences.
var persistencesServiceGen = wire.NewSet(
	infras.ProvideSQLiteConn,
)

var cache = wire.NewSet(
	infras.ProvideRedis,
)

// Auth
var authService = wire.NewSet(
	auth.NewTokenService,
)

// Services.
var ServiceGen = wire.NewSet(
	BaseService.NewService,
	wire.Bind(new(BaseService.ServiceInterface), new(*BaseService.ServiceImpl)),

	BaseRepository.NewRepository,
	wire.Bind(new(BaseRepository.RepoInterface), new(*BaseRepository.RepositoryImpl)),
	UserRepository.NewUserRepository,
	wire.Bind(new(UserRepository.UserRepoInterface), new(*UserRepository.UserRepositoryImpl)),
	AppointmentRepository.NewAppointmentRepository,
	wire.Bind(new(AppointmentRepository.AppointmentRepoInterface), new(*AppointmentRepository.AppointmentRepositoryImpl)),
	AppointmentLogRepository.NewAppointmentLogRepository,
	wire.Bind(new(AppointmentLogRepository.AppointmentLogRepoInterface), new(*AppointmentLogRepository.AppointmentLogRepositoryImpl)),
)

// User service.
var userServiceV1Gen = wire.NewSet(
	UserServiceV1.NewUserService,
	wire.Bind(new(UserServiceV1.UserService), new(*UserServiceV1.UserServiceImpl)),
)

// Appointment service.
var appointmentServiceV1Gen = wire.NewSet(
	AppointmentServiceV1.NewAppointmentService,
	wire.Bind(new(AppointmentServiceV1.AppointmentService), new(*AppointmentServiceV1.AppointmentServiceImpl)),
)

// Appointment log service.
var appointmentLogServiceV1Gen = wire.NewSet(
	AppointmentLogServiceV1.NewAppointmentLogService,
	wire.Bind(new(AppointmentLogServiceV1.AppointmentLogService), new(*AppointmentLogServiceV1.AppointmentLogServiceImpl)),
)

var initializeServiceServiceGen = wire.NewSet(
	ServiceGen,
	authService,
	userServiceV1Gen,
	appointmentServiceV1Gen,
	appointmentLogServiceV1Gen,
)

var authMiddleware = wire.NewSet(
	middleware.NewAuthMiddleware,
	wire.Bind(new(middleware.AuthMiddlewareInterface), new(*middleware.AuthMiddleware)),
)

// Wiring for HTTP routing.
var routingServiceGen = wire.NewSet(
	wire.Struct(new(router.DomainHandlers), "*"),
	handlerV1.NewHandler,
	router.NewRouter,
)

// Wiring for everything.
func InitializeServiceServiceGen() *http.HTTP {
	wire.Build(
		configurationsServiceGen,
		persistencesServiceGen,
		cache,
		initializeServiceServiceGen,
		routingServiceGen,
		http.NewHTTP,
		authMiddleware,
	)

	return &http.HTTP{}
}
