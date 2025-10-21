package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)


type DBConfig struct {
	Host        string
	Port        string 
	User        string
	Password    string
	Name        string
	EnableSSLMode bool
}

type Config struct {
	Version       string
	ServiceName   string
	HttpPort      int
	JWTSecretKey  string
	CloudinaryURL string
	DB            *DBConfig 
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

	
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		fmt.Println("DB_HOST is required")
		os.Exit(1)
	}
	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		fmt.Println("DB_PORT is required")
		os.Exit(1)
	}
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		fmt.Println("DB_USER is required")
		os.Exit(1)
	}
	dbPass := os.Getenv("DB_PASS")
	if dbPass == "" {
		fmt.Println("DB_PASS is required")
		os.Exit(1)
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		fmt.Println("DB_NAME is required")
		os.Exit(1)
	}
	enableSSLModeStr := os.Getenv("ENABLE_SSL_MODE")
	if enableSSLModeStr == "" {
		fmt.Println("ENABLE_SSL_MODE is required")
		os.Exit(1)
	}
	enableSSLMode, err := strconv.ParseBool(enableSSLModeStr)
	if err != nil {
		fmt.Println("Invalid SSL Mode")
		os.Exit(1)
	}

	dbconfig := &DBConfig{
		Host:        dbHost,
		Port:        dbPort,
		User:        dbUser,
		Password:    dbPass,
		Name:        dbName,
		EnableSSLMode: enableSSLMode,
	}

	
	configurations = &Config{
		Version:       os.Getenv("VERSION"),
		ServiceName:   os.Getenv("SERVICE_NAME"),
		HttpPort:      int(port),
		JWTSecretKey:  os.Getenv("JWT_SECRET_KEY"),
		CloudinaryURL: os.Getenv("CLOUDINARY_URL"),
		DB:            dbconfig,
	}
}

func GetConfig() *Config {
	if configurations == nil {
		loadConfig()
	}
	return configurations
}