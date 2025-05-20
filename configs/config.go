package configs

import (
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/rs/zerolog/log"
)

// Config is a struct that will receive configuration options via environment
// variables.
type Config struct {
	App struct {
		CORS struct {
			AllowCredentials bool
			AllowedHeaders   []string
			AllowedMethods   []string
			AllowedOrigins   []string
			Enable           bool
			MaxAgeSeconds    int
		}
		Name string
		Tz   string
		URL  string
	}
	DB struct {
		Postgres struct {
			Read struct {
				Host     string
				Port     string
				User     string
				Password string
				Name     string
				Timezone string
			}
			Write struct {
				Host     string
				Port     string
				User     string
				Password string
				Name     string
				Timezone string
			}
		}
		SQLite struct {
			Path string
		}
	}
	Server struct {
		Env      string
		LogLevel string
		Port     string
		Shutdown struct {
			CleanupPeriodSeconds int64
			GracePeriodSeconds   int64
		}
	}
	Cache struct {
		Redis struct {
			Primary struct {
				Host       string
				Port       string
				Password   string
				Prefix     string
				RetryCount int
			}
		}
	}
	AccessToken struct {
		ExpiryInHour int
		Secret       string
	}
	RefreshToken struct {
		ExpiryInHour int
		Secret       string
	}
}

var (
	conf Config
	once sync.Once
)

// Get loads configuration from environment variables
func Get() *Config {
	once.Do(func() {
		log.Info().Msg("Loading configuration from environment variables")

		// App CORS config
		conf.App.CORS.AllowCredentials = getBool("APP_CORS_ALLOWCREDENTIALS", true)
		conf.App.CORS.AllowedHeaders = getStringSlice("APP_CORS_ALLOWEDHEADERS", []string{"Accept", "Authorization", "Content-Type"})
		conf.App.CORS.AllowedMethods = getStringSlice("APP_CORS_ALLOWEDMETHODS", []string{"GET", "POST", "PUT", "DELETE"})
		conf.App.CORS.AllowedOrigins = getStringSlice("APP_CORS_ALLOWEDORIGINS", []string{"*"})
		conf.App.CORS.Enable = getBool("APP_CORS_ENABLE", true)
		conf.App.CORS.MaxAgeSeconds = getInt("APP_CORS_MAXAGESECONDS", 300)

		// App config
		conf.App.Name = getEnv("APP_NAME", "mini-evv")
		conf.App.Tz = getEnv("APP_TZ", "UTC")
		conf.App.URL = getEnv("APP_URL", "http://localhost:3200")

		// DB SQLite config
		conf.DB.SQLite.Path = getEnv("DB_SQLITE_PATH", "db/mini-evv.db")

		// Server config
		conf.Server.Env = getEnv("SERVER_ENV", "development")
		conf.Server.LogLevel = getEnv("SERVER_LOGLEVEL", "info")
		conf.Server.Port = getEnv("SERVER_PORT", "3200")
		conf.Server.Shutdown.CleanupPeriodSeconds = getInt64("SERVER_SHUTDOWN_CLEANUP_PERIOD_SECONDS", 15)
		conf.Server.Shutdown.GracePeriodSeconds = getInt64("SERVER_SHUTDOWN_GRACE_PERIOD_SECONDS", 15)

		// Cache Redis config
		conf.Cache.Redis.Primary.Host = getEnv("CACHE_REDIS_PRIMARY_HOST", "")
		conf.Cache.Redis.Primary.Port = getEnv("CACHE_REDIS_PRIMARY_PORT", "6379")
		conf.Cache.Redis.Primary.Password = getEnv("CACHE_REDIS_PRIMARY_PASSWORD", "")
		conf.Cache.Redis.Primary.Prefix = getEnv("CACHE_REDIS_PRIMARY_PREFIX", "mini-evv:")
		conf.Cache.Redis.Primary.RetryCount = getInt("CACHE_REDIS_PRIMARY_RETRY_COUNT", 3)

		// Token configs
		conf.AccessToken.ExpiryInHour = getInt("ACCESSTOKEN_EXPIRYINHOUR", 1)
		conf.AccessToken.Secret = getEnv("ACCESSTOKEN_SECRET", "access-token-secret")
		conf.RefreshToken.ExpiryInHour = getInt("REFRESHTOKEN_EXPIRYINHOUR", 24*7) // 1 week
		conf.RefreshToken.Secret = getEnv("REFRESHTOKEN_SECRET", "refresh-token-secret")

		log.Info().Msg("Configuration loaded successfully")
	})

	return &conf
}

// Helper functions to get environment variables with default values
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getBool(key string, defaultValue bool) bool {
	strValue := os.Getenv(key)
	if strValue == "" {
		return defaultValue
	}

	boolValue, err := strconv.ParseBool(strValue)
	if err != nil {
		log.Warn().Str("key", key).Str("value", strValue).Msg("Invalid boolean value, using default")
		return defaultValue
	}
	return boolValue
}

func getInt(key string, defaultValue int) int {
	strValue := os.Getenv(key)
	if strValue == "" {
		return defaultValue
	}

	intValue, err := strconv.Atoi(strValue)
	if err != nil {
		log.Warn().Str("key", key).Str("value", strValue).Msg("Invalid integer value, using default")
		return defaultValue
	}
	return intValue
}

func getInt64(key string, defaultValue int64) int64 {
	strValue := os.Getenv(key)
	if strValue == "" {
		return defaultValue
	}

	int64Value, err := strconv.ParseInt(strValue, 10, 64)
	if err != nil {
		log.Warn().Str("key", key).Str("value", strValue).Msg("Invalid int64 value, using default")
		return defaultValue
	}
	return int64Value
}

func getStringSlice(key string, defaultValue []string) []string {
	strValue := os.Getenv(key)
	if strValue == "" {
		return defaultValue
	}

	return strings.Split(strValue, ",")
}
