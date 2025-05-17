package http

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/farganamar/evv-service/configs"
	"github.com/farganamar/evv-service/helpers/logger"
	"github.com/farganamar/evv-service/infras"
	"github.com/farganamar/evv-service/transport/http/response"
	"github.com/farganamar/evv-service/transport/http/router"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
)

// ServerState is an
type ServerState int

const (
	// ServerStateReady indicates that the server is ready to serve.
	ServerStateReady ServerState = iota + 1
	// ServerStateInGracePeriod indicates that the server is in its grace
	// period and will shut down after it is done cleaning up.
	ServerStateInGracePeriod
	// ServerStateInCleanupPeriod indicates that the server no longer
	// responds to any requests, is cleaning up its internal state, and
	// will shut down shortly.
	ServerStateInCleanupPeriod
)

// HTTP is the HTTP server.
type HTTP struct {
	Config *configs.Config
	DB     *infras.SQLiteConn
	Router router.Router
	State  ServerState
	mux    *chi.Mux
	Server *http.Server
}

// NewHTTP creates a new HTTP server.
func NewHTTP(
	config *configs.Config,
	db *infras.SQLiteConn,
	router router.Router) *HTTP {
	return &HTTP{
		Config: config,
		DB:     db,
		Router: router,
		Server: &http.Server{
			Addr: ":" + config.Server.Port,
		},
	}
}

// Start starts the HTTP server.
func (h *HTTP) SetupAndServe() {
	h.mux = chi.NewRouter()
	h.setupMiddleware()
	h.setupRoutes()
	h.setupGracefulShutdown()
	h.State = ServerStateReady

	h.logServerInfo()

	log.Info().Str("port", h.Config.Server.Port).Msg("Starting HTTP server.")

	h.Server.Handler = h.mux

	err := h.Server.ListenAndServe()

	if err != nil {
		log.Error().Err(err).Stack().Send()
	}

	// h.mux = chi.NewRouter()
	// h.setupMiddleware()
	// h.setupRoutes()
	// h.setupGracefulShutdown()
	// h.State = ServerStateReady

	// h.logServerInfo()

	// log.Info().Str("port", h.Config.Server.Port).Msg("Starting up HTTP server.")

	// err := http.ListenAndServe(":"+h.Config.Server.Port, h.mux)
	// if err != nil {
	// 	logger.ErrorWithStack(err)
	// }
}

func (h *HTTP) setupRoutes() {
	h.mux.Get("/health", h.HealthCheck)
	decimal.MarshalJSONWithoutQuotes = true
	h.Router.SetupRoutes(h.mux)
}

func (h *HTTP) setupGracefulShutdown() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)
	go h.respondToSigterm(done)
}

func (h *HTTP) respondToSigterm(done chan os.Signal) {
	<-done
	defer os.Exit(0)

	shutdownConfig := h.Config.Server.Shutdown

	log.Info().Msg("Received SIGTERM.")
	log.Info().Int64("seconds", shutdownConfig.GracePeriodSeconds).Msg("Entering grace period.")
	h.State = ServerStateInGracePeriod
	time.Sleep(time.Duration(shutdownConfig.GracePeriodSeconds) * time.Second)

	log.Info().Int64("seconds", shutdownConfig.CleanupPeriodSeconds).Msg("Entering cleanup period.")
	h.State = ServerStateInCleanupPeriod
	time.Sleep(time.Duration(shutdownConfig.CleanupPeriodSeconds) * time.Second)

	log.Info().Msg("Cleaning up completed. Shutting down now.")
}

func (h *HTTP) setupMiddleware() {
	h.mux.Use(middleware.Logger)
	h.mux.Use(middleware.Recoverer)
	h.mux.Use(h.serverStateMiddleware)
	h.setupCORS()
}

func (h *HTTP) logServerInfo() {
	h.logCORSConfigInfo()
}

func (h *HTTP) logCORSConfigInfo() {
	corsConfig := h.Config.App.CORS
	corsHeaderInfo := "CORS Header"
	if corsConfig.Enable {
		log.Info().Msg("CORS Headers and Handlers are enabled.")
		log.Info().Str(corsHeaderInfo, fmt.Sprintf("Access-Control-Allow-Credentials: %t", corsConfig.AllowCredentials)).Msg("")
		log.Info().Str(corsHeaderInfo, fmt.Sprintf("Access-Control-Allow-Headers: %s", strings.Join(corsConfig.AllowedHeaders, ", "))).Msg("")
		log.Info().Str(corsHeaderInfo, fmt.Sprintf("Access-Control-Allow-Methods: %s", strings.Join(corsConfig.AllowedMethods, ", "))).Msg("")
		log.Info().Str(corsHeaderInfo, fmt.Sprintf("Access-Control-Allow-Origin: %s", strings.Join(corsConfig.AllowedOrigins, ", "))).Msg("")
		log.Info().Str(corsHeaderInfo, fmt.Sprintf("Access-Control-Max-Age: %d", corsConfig.MaxAgeSeconds)).Msg("")
	} else {
		log.Info().Msg("CORS Headers are disabled.")
	}
}

func (h *HTTP) serverStateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch h.State {
		case ServerStateReady:
			// Server is ready to serve, don't do anything.
			next.ServeHTTP(w, r)
		case ServerStateInGracePeriod:
			// Server is in grace period. Issue a warning message and continue
			// serving as usual.
			log.Warn().Msg("SERVER IS IN GRACE PERIOD")
			next.ServeHTTP(w, r)
		case ServerStateInCleanupPeriod:
			// Server is in cleanup period. Stop the request from actually
			// invoking any domain services and respond appropriately.
			response.WithPreparingShutdown(w)
		}
	})
}

func (h *HTTP) setupCORS() {
	corsConfig := h.Config.App.CORS
	if corsConfig.Enable {
		h.mux.Use(cors.Handler(cors.Options{
			AllowCredentials: corsConfig.AllowCredentials,
			AllowedHeaders:   corsConfig.AllowedHeaders,
			AllowedMethods:   corsConfig.AllowedMethods,
			AllowedOrigins:   corsConfig.AllowedOrigins,
			MaxAge:           corsConfig.MaxAgeSeconds,
		}))
	}
}

func (h *HTTP) HealthCheck(w http.ResponseWriter, r *http.Request) {
	if err := h.DB.DB.Ping(); err != nil {
		logger.ErrorWithStack(err)
		response.WithUnhealthy(w)
		return
	}
	response.WithMessage(w, http.StatusOK, "OK")
}

func (h *HTTP) Shutdown() {
	gracePeriod := h.Config.Server.Shutdown.GracePeriodSeconds
	cleanupPeriod := h.Config.Server.Shutdown.CleanupPeriodSeconds

	shutdown := (time.Duration(gracePeriod) + time.Duration(cleanupPeriod)) * time.Second
	ctx, cancel := context.WithTimeout(context.TODO(), shutdown)
	defer cancel()

	log.Info().Msg("Rceived SIGTERM.")
	log.Info().Int64("seconds", gracePeriod).Msg("HTTP entering grace period.")
	h.State = ServerStateInGracePeriod

	time.Sleep(time.Duration(gracePeriod) * time.Second)

	log.Info().Int64("seconds", cleanupPeriod).Msg("HTTP entering cleanup period.")
	h.State = ServerStateInCleanupPeriod
	time.Sleep(time.Duration(cleanupPeriod) * time.Second)

	err := h.Server.Shutdown(ctx)
	if err != nil {
		log.Error().Err(err).Stack().Send()
	}
}
