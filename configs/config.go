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
			AllowCredentials bool     `mapstructure:"ALLOW_CREDENTIALS"`
			AllowedHeaders   []string `mapstructure:"ALLOWED_HEADERS"`
			AllowedMethods   []string `mapstructure:"ALLOWED_METHODS"`
			AllowedOrigins   []string `mapstructure:"ALLOWED_ORIGINS"`
			Enable           bool     `mapstructure:"ENABLE"`
			MaxAgeSeconds    int      `mapstructure:"MAX_AGE_SECONDS"`
		} `mapstructure:"cors"`
		Name string `mapstructure:"NAME"`
		Tz   string `mapstructure:"TZ"`
		URL  string `mapstructure:"URL"`
	} `mapstructure:"app"`
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
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed reading config file")
	}

	// Print out all keys Viper knows about
	for _, key := range viper.AllKeys() {
		val := viper.GetString(key)
		newKey := strings.ReplaceAll(key, "_", ".")
		viper.Set(newKey, val)
	}

	once.Do(func() {
		log.Info().Msg("Service configuration initialized.")
		err = viper.Unmarshal(&conf)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed unmarshall config")
		}
	})

	return &conf
}
