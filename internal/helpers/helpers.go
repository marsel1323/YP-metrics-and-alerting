package helpers

import (
	"log"
	"os"
	"time"
)

func GetEnv(key string, defaultValue string) string {
	env, ok := os.LookupEnv(key)
	if ok {
		return env
	}
	return defaultValue
}

func StringToSeconds(value string) time.Duration {
	sec, err := time.ParseDuration(value)
	if err != nil {
		log.Println(err)
		return 0
	}
	return sec
}
