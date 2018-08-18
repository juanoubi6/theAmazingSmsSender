package config

import "os"

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
		config := newConfig()
		instance = &config
	}
	return instance
}

func newConfig() Config {
	return Config{
		ENV:  GetEnv("ENV", "develop"),
		PORT: GetEnv("PORT", "5000"),

		TWILIO_SID:        GetEnv("TWILIO_SID", "ACb98ded914e1c12b0e276c4f164555f70"),
		TWILIO_AUTH_TOKEN: GetEnv("TWILIO_AUTH_TOKEN", "487d710fa087d1a2f5f757d579a0d367"),
		TWILIO_ACC_PHONE:  GetEnv("TWILIO_ACC_PHONE", "+19852418221"),

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
