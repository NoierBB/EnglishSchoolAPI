package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	HTTPPort   string
}

func Load() Config {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Println("warning: .env not loaded from current dir")
	}

	portStr := os.Getenv("DB_PORT")
	if portStr == "" {
		log.Fatal("DB_PORT is empty (env not loaded)")
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatal("DB_PORT is not a number")
	}

	return Config{
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     port,
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		HTTPPort:   os.Getenv("HTTP_PORT"),
	}
}
