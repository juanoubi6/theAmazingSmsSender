package config

import (
	"os"
	"bufio"
	"strings"
)

type Config struct {
	ENV  string
	PORT string

	RABBITMQ_USER     string
	RABBITMQ_PASSWORD string
	RABBITMQ_HOST     string
	RABBITMQ_PORT     string

	TWILIO_SID        string
	TWILIO_AUTH_TOKEN string
	TWILIO_ACC_PHONE  string

	WORKER_AMOUNT string
}

var instance *Config

func GetConfig() *Config {
	if instance == nil {
		err := readEnv()
		if err != nil{
			panic(err)
		}
		config := newConfig()
		instance = &config
	}
	return instance
}

func newConfig() Config {
	return Config{
		ENV:  GetEnv("ENV", "develop"),
		PORT: GetEnv("PORT", "5000"),

		TWILIO_SID:        GetEnv("TWILIO_SID", ""),
		TWILIO_AUTH_TOKEN: GetEnv("TWILIO_AUTH_TOKEN", ""),
		TWILIO_ACC_PHONE:  GetEnv("TWILIO_ACC_PHONE", ""),

		RABBITMQ_HOST:     GetEnv("RABBITMQ_HOST", "localhost"),
		RABBITMQ_PORT:     GetEnv("RABBITMQ_PORT", "5672"),
		RABBITMQ_USER:     GetEnv("RABBITMQ_USER", "guest"),
		RABBITMQ_PASSWORD: GetEnv("RABBITMQ_PASSWORD", "guest"),

		WORKER_AMOUNT: GetEnv("WORKER_AMOUNT", "3"),
	}
}

func GetEnv(key, fallback string) string {
	if val, exists := os.LookupEnv(key); exists {
		return val
	}
	return fallback
}

func readEnv() error{
	file, err := os.Open(".env")
	if err != nil {
		return nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		values := strings.Split(scanner.Text(),"=")
		if len(values)==2{
			err = os.Setenv(values[0],values[1])
			if err != nil{
				return err
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}