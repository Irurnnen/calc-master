package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

type Config struct {
	TimeAdditions       int
	TimeSubtractions    int
	TimeMultiplications int
	TimeDivisions       int
	ServerURL           string
}

func NewConfigExample() *Config {
	return &Config{
		TimeAdditions:       0,
		TimeSubtractions:    0,
		TimeMultiplications: 0,
		TimeDivisions:       0,
		ServerURL:           "0.0.0.0:1234",
	}
}

func NewConfigFromEnv() *Config {
	timeAdditionsMS, err := strconv.Atoi(os.Getenv("TIME_ADDITION_MS"))
	if err != nil {
		log.Fatalf("Fatal error while getting config from env: TIME_ADDITION_MS: %s", os.Getenv("COMPUTING_POWER"))
		return NewConfigExample()
	}
	timeSubtractionMS, err := strconv.Atoi(os.Getenv("TIME_SUBTRACTION_MS"))
	if err != nil {
		log.Fatalf("Fatal error while getting config from env: TIME_SUBTRACTION_MS: %s", os.Getenv("COMPUTING_POWER"))
		return NewConfigExample()
	}
	timeMultiplicationsMS, err := strconv.Atoi(os.Getenv("TIME_MULTIPLICATIONS_MS"))
	if err != nil {
		log.Fatalf("Fatal error while getting config from env: TIME_MULTIPLICATIONS_MS: %s", os.Getenv("COMPUTING_POWER"))
		return NewConfigExample()
	}
	timeDivisionsMS, err := strconv.Atoi(os.Getenv("TIME_DIVISIONS_MS"))
	if err != nil {
		log.Fatalf("Fatal error while getting config from env: TIME_DIVISIONS_MS: %s", os.Getenv("COMPUTING_POWER"))
		return NewConfigExample()
	}
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatalf("Fatal error while getting config from env: PORT: %s", os.Getenv("COMPUTING_POWER"))
		return NewConfigExample()
	}
	return &Config{
		TimeAdditions:       timeAdditionsMS,
		TimeSubtractions:    timeSubtractionMS,
		TimeMultiplications: timeMultiplicationsMS,
		TimeDivisions:       timeDivisionsMS,
		ServerURL:           fmt.Sprintf("0.0.0.0:%d", port),
	}
}
