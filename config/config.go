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

	version := os.Getenv("VERSION")
	if version == "" {
		fmt.Println("Version is required")
		os.Exit(1)
	}

	// FIX: Your .env has SERVICE_NAME, not Service_Name
	serviceName := os.Getenv("SERVICE_NAME")
	if serviceName == "" {
		fmt.Println("Service name is required")
		os.Exit(1)
	}

	// FIX: This code now works for BOTH Render (PORT) and Local (HTTP_PORT)
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
	// END FIX

	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")
	if jwtSecretKey == "" {
		fmt.Println("Jwt secret key is needed")
		os.Exit(1)
	}

	cloudinaryURL := os.Getenv("CLOUDINARY_URL")
	if cloudinaryURL == "" {
		fmt.Println("Cloudinary URL is required")
		os.Exit(1)
	}

	// FIX: We only need DATABASE_URL. Remove all DB_HOST, DB_USER, etc.
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		fmt.Println("DATABASE_URL is needed")
		os.Exit(1)
	}
	// END FIX

	configurations = &Config{
		Version:       version,
		ServiceName:   serviceName,
		HttpPort:      int(port),
		JWTSecretKey:  jwtSecretKey,
		CloudinaryURL: cloudinaryURL,
		DatabaseURL:   dbURL,
	}
}

func GetConfig() *Config {
	if configurations == nil {
		loadConfig()
	}
	return configurations
}