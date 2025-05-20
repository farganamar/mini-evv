package configs

import (
	"strings"
	"sync"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// Config is a struct that will receive configuration options via environment
// variables.

type Config struct {
	App struct {
		CORS struct {
			AllowCredentials bool     `mapstructure:"ALLOWCREDENTIALS"`
			AllowedHeaders   []string `mapstructure:"ALLOWEDHEADERS"`
			AllowedMethods   []string `mapstructure:"ALLOWEDMETHODS"`
			AllowedOrigins   []string `mapstructure:"ALLOWEDORIGINS"`
			Enable           bool     `mapstructure:"ENABLE"`
			MaxAgeSeconds    int      `mapstructure:"MAXAGESECONDS"`
		} `mapstructure:"cors"`
		Name string `mapstructure:"NAME"`
		Tz   string `mapstructure:"TZ"`
		URL  string `mapstructure:"URL"`
	} `mapstructure:"APP"`
	DB struct {
		Postgres struct {
			Read struct {
				Host     string `mapstructure:"HOST"`
				Port     string `mapstructure:"PORT"`
				User     string `mapstructure:"USER"`
				Password string `mapstructure:"PASSWORD"`
				Name     string `mapstructure:"NAME"`
				Timezone string `mapstructure:"TIMEZONE"`
			} `mapstructure:"read"`
			Write struct {
				Host     string `mapstructure:"HOST"`
				Port     string `mapstructure:"PORT"`
				User     string `mapstructure:"USER"`
				Password string `mapstructure:"PASSWORD"`
				Name     string `mapstructure:"NAME"`
				Timezone string `mapstructure:"TIMEZONE"`
			} `mapstructure:"WRITE"`
		} `mapstructure:"POSTGRES"`
		SQLite struct {
			Path string `mapstructure:"PATH"`
		}
	} `mapstructure:"DB"`
	Server struct {
		Env      string `mapstructure:"ENV"`
		LogLevel string `mapstructure:"LOGLEVEL"`
		Port     string `mapstructure:"PORT"`
		Shutdown struct {
			CleanupPeriodSeconds int64 `mapstructure:"CLEANUP_PERIOD_SECONDS"`
			GracePeriodSeconds   int64 `mapstructure:"GRACE_PERIOD_SECONDS"`
		}
	} `mapstructure:"server"`
	Cache struct {
		Redis struct {
			Primary struct {
				Host       string `mapstructure:"HOST"`
				Port       string `mapstructure:"PORT"`
				Password   string `mapstructure:"PASSWORD"`
				Prefix     string `mapstructure:"PREFIX"`
				RetryCount int    `mapstructure:"RETRY_COUNT"`
			}
		}
	} `mapstructure:"cache"`
	AccessToken struct {
		ExpiryInHour int    `mapstructure:"EXPIRYINHOUR"`
		Secret       string `mapstructure:"SECRET"`
	} `mapstructure:"ACCESSTOKEN"`
	RefreshToken struct {
		ExpiryInHour int    `mapstructure:"EXPIRYINHOUR"`
		Secret       string `mapstructure:"SECRET"`
	} `mapstructure:"REFRESHTOKEN"`
}

var (
	conf Config
	once sync.Once
)

// Get are responsible to load env and get data an return the struct
func Get() *Config {
	// Set up automatic environment variable reading
	viper.AutomaticEnv()

	// Try to read from .env file, but don't fail if it doesn't exist
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		// If .env file doesn't exist or can't be read, just log it and continue
		// using environment variables
		log.Info().Msg("No .env file found or couldn't be read, using environment variables")
	} else {
		log.Info().Msg("Config loaded from .env file")
	}

	// Set up mappings between environment variables and config fields
	// This ensures that environment variables are properly mapped to struct fields
	viper.SetEnvPrefix("") // No prefix for environment variables
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Print out all keys Viper knows about
	for _, key := range viper.AllKeys() {
		val := viper.GetString(key)
		newKey := strings.ReplaceAll(key, "_", ".")
		viper.Set(newKey, val)
	}

	once.Do(func() {
		log.Info().Msg("Service configuration initialized.")
		err := viper.Unmarshal(&conf)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed unmarshall config")
		}
	})

	return &conf
}

// func Get() *Config {
// 	viper.SetConfigFile(".env")
// 	viper.AutomaticEnv()

// 	err := viper.ReadInConfig()
// 	if err != nil {
// 		log.Fatal().Err(err).Msg("Failed reading config file")
// 	}

// 	// Print out all keys Viper knows about
// 	for _, key := range viper.AllKeys() {
// 		val := viper.GetString(key)
// 		newKey := strings.ReplaceAll(key, "_", ".")
// 		viper.Set(newKey, val)
// 	}

// 	once.Do(func() {
// 		log.Info().Msg("Service configuration initialized.")
// 		err = viper.Unmarshal(&conf)
// 		if err != nil {
// 			log.Fatal().Err(err).Msg("Failed unmarshall config")
// 		}
// 	})

// 	return &conf
// }
