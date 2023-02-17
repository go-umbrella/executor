package env

import (
	"os"
	"strconv"
)

func Bool(key string, defaultValue bool) bool {
	osValue := os.Getenv(key)
	if osValue == "" {
		return defaultValue
	}

	value, err := strconv.ParseBool(osValue)
	if err != nil {
		return defaultValue
	}

	return value
}

func Uint64(key string, defaultValue uint64) uint64 {
	osValue := os.Getenv(key)
	if osValue == "" {
		return defaultValue
	}

	value, err := strconv.ParseUint(osValue, 10, 64)
	if err != nil {
		return defaultValue
	}

	return value
}
