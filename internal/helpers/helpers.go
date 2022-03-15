package helpers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
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

func Hash(src string, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(src))
	dst := h.Sum(nil)
	return base64.StdEncoding.EncodeToString(dst)
}

func Compare(str string, key string) error {
	log.Println("Compare...")
	data, err := hex.DecodeString(str)
	if err != nil {
		return err
	}
	h := hmac.New(sha256.New, []byte(key))
	h.Write(data)
	sign := h.Sum(nil)
	log.Println("Sign:", sign)
	return nil
}
