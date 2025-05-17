//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/zorahealth/user-service/configs"
	"github.com/zorahealth/user-service/helpers/auth"
	"github.com/zorahealth/user-service/infras"
	handler "github.com/zorahealth/user-service/internal/handlers"
	"github.com/zorahealth/user-service/internal/repository"
	"github.com/zorahealth/user-service/internal/service"
	"github.com/zorahealth/user-service/transport/http"
	"github.com/zorahealth/user-service/transport/http/middleware"
	"github.com/zorahealth/user-service/transport/http/router"
)

// Wiring for configurations.
var configurationsServiceGen = wire.NewSet(
	configs.Get,
)

// Wiring for persistences.
var persistencesServiceGen = wire.NewSet(
	infras.ProvidePostgresConn,
)

var cache = wire.NewSet(
	infras.ProvideRedis,
)

// Auth
var authService = wire.NewSet(
	auth.NewTokenService,
)

// User services.
var userServiceGen = wire.NewSet(
	service.NewUserService,
	wire.Bind(new(service.UserServiceInterface), new(*service.UserServiceImpl)),
	repository.NewUserRepository,
	wire.Bind(new(repository.UserRepoInterface), new(*repository.UserRepositoryImpl)),
)

var initializeServiceServiceGen = wire.NewSet(
	userServiceGen,
	authService,
)

var authMiddleware = wire.NewSet(
	middleware.NewAuthMiddleware,
	wire.Bind(new(middleware.AuthMiddlewareInterface), new(*middleware.AuthMiddleware)),
)

// Wiring for HTTP routing.
var routingServiceGen = wire.NewSet(
	wire.Struct(new(router.DomainHandlers), "*"),
	handler.NewUserHandler,
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
