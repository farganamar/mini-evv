//go:build wireinject
// +build wireinject

package main

import (
	"github.com/farganamar/evv-service/configs"
	"github.com/farganamar/evv-service/helpers/auth"
	"github.com/farganamar/evv-service/infras"
	handler "github.com/farganamar/evv-service/internal/handlers"
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
// service.NewUserService,
// wire.Bind(new(service.UserServiceInterface), new(*service.UserServiceImpl)),
// repository.NewUserRepository,
// wire.Bind(new(repository.UserRepoInterface), new(*repository.UserRepositoryImpl)),
)

var initializeServiceServiceGen = wire.NewSet(
	ServiceGen,
	authService,
)

var authMiddleware = wire.NewSet(
	middleware.NewAuthMiddleware,
	wire.Bind(new(middleware.AuthMiddlewareInterface), new(*middleware.AuthMiddleware)),
)

// Wiring for HTTP routing.
var routingServiceGen = wire.NewSet(
	wire.Struct(new(router.DomainHandlers), "*"),
	handler.NewHandler,
	router.NewRouter,
)

// Wiring for everything.
func InitializeServiceServiceGen() *http.HTTP {
	wire.Build(
		configurationsServiceGen,
		persistencesServiceGen,
		// cache,
		// initializeServiceServiceGen,
		routingServiceGen,
		http.NewHTTP,
		// authMiddleware,
	)

	return &http.HTTP{}
}
