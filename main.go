package main

//go:generate go run github.com/google/wire/cmd/wire

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/zorahealth/user-service/configs"
	"github.com/zorahealth/user-service/helpers/logger"
	"github.com/zorahealth/user-service/helpers/shutdown"
)

var configServiceGen *configs.Config

// @securityDefinitions.apikey EVMOauthToken
// @in header
// @name Authorization
func main() {
	// Initialize logger
	logger.InitLogger()

	// Initialize config
	configServiceGen = configs.Get()

	location, _ := time.LoadLocation(configServiceGen.App.Tz)
	time.Local = location

	// Set desired log level
	logger.SetLogLevel(configServiceGen)

	// Wire everything up
	httpServiceGen := InitializeServiceServiceGen()

	// Run server
	httpServiceGen.SetupAndServe()

	stop, close := gracefulShutdown()
	defer close()
	<-stop

	gracefulShutdown := shutdown.NewGracefulShutdown(
		[]func(){
			httpServiceGen.Shutdown,
		},
		shutdown.SetCleanupPeriodSeconds(configServiceGen.Server.Shutdown.CleanupPeriodSeconds),
		shutdown.SetGracePeriodSeconds(configServiceGen.Server.Shutdown.GracePeriodSeconds),
	)
	gracefulShutdown.Shutdown()
}

func gracefulShutdown() (chan os.Signal, func()) {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	return done, func() {
		close(done)
	}
}
