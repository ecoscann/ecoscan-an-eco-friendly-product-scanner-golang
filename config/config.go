package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Version       string
	ServiceName   string
	HttpPort      int
	JWTSecretKey  string
	CloudinaryURL string
	DatabaseURL   string 
}

var configurations *Config

func loadConfig() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("No .env file found, using environment variables")
	}


	httpPortStr := os.Getenv("PORT")
	if httpPortStr == "" {
		httpPortStr = os.Getenv("HTTP_PORT") 
		if httpPortStr == "" {
			fmt.Println("PORT or HTTP_PORT is required")
			os.Exit(1)
		}
	}
	port, err := strconv.ParseInt(httpPortStr, 10, 64)
	if err != nil {
		fmt.Println("Port must be a number")
		os.Exit(1)
	}

	
	configurations = &Config{
		Version:       os.Getenv("VERSION"),
		ServiceName:   os.Getenv("SERVICE_NAME"),
		HttpPort:      int(port),
		JWTSecretKey:  os.Getenv("JWT_SECRET_KEY"),
		CloudinaryURL: os.Getenv("CLOUDINARY_URL"),
		DatabaseURL:   os.Getenv("DATABASE_URL"), 
	}


	if configurations.DatabaseURL == "" {
		fmt.Println("DATABASE_URL is required")
		os.Exit(1)
	}
	if configurations.JWTSecretKey == "" {
		fmt.Println("JWT_SECRET_KEY is required")
		os.Exit(1)
	}
}

func GetConfig() *Config {
	if configurations == nil {
		loadConfig()
	}
	return configurations
}