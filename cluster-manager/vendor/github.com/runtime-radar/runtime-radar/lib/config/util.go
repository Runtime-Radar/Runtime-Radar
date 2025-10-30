package config

import (
	"os"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
)

func LookupEnvString(key string, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}

func LookupEnvBool(key string, defaultVal bool) bool {
	if val, ok := os.LookupEnv(key); ok {
		if val == "" || val == "0" || val == "false" {
			return false
		}
		return true
	}
	return defaultVal
}

func LookupEnvInt(key string, defaultVal int) int {
	if val, ok := os.LookupEnv(key); ok {
		v, err := strconv.Atoi(val)
		if err != nil {
			log.Fatal().Msgf("### Can't convert string '%s' to int: %v", val, err)
		}
		return v
	}
	return defaultVal
}

func LookupEnvFloat64(key string, defaultVal float64) float64 {
	if val, ok := os.LookupEnv(key); ok {
		v, err := strconv.ParseFloat(val, 64)
		if err != nil {
			log.Fatal().Msgf("### Can't convert string '%s' to float64: %v", val, err)
		}
		return v
	}
	return defaultVal
}

func LookupEnvDuration(key string, defaultVal time.Duration) time.Duration {
	if val, ok := os.LookupEnv(key); ok {
		v, err := time.ParseDuration(val)
		if err != nil {
			log.Fatal().Msgf("### Can't convert string '%s' to time.Duration: %v", val, err)
		}
		return v
	}
	return defaultVal
}
