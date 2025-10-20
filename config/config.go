package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var configurations *Config

type DBConfig struct {
	Host          string
	Port          int
	User          string
	Password      string
	Name          string
	EnableSSLMode bool
}

type Config struct {
	Version      string
	ServiceName  string
	HttpPort     int
	JWTSecretKey string
	DB           *DBConfig
}

func loadConfig() {

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Failed to load env fileand err:", err)
		os.Exit(1)
	}
	version := os.Getenv("VERSION")
	if version == "" {
		fmt.Println("Version is required")
		os.Exit(1)
	}

	serviceName := os.Getenv("Service_Name")
	if serviceName == "" {
		fmt.Println("Service name its required")
		os.Exit(1)
	}

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		fmt.Println("Port is required")
		os.Exit(1)
	}

	port, err := strconv.ParseInt(httpPort, 10, 64)

	if err != nil {
		fmt.Println("Port must be number")
		os.Exit(1)
	}

	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")
	if jwtSecretKey == "" {
		fmt.Println("Jwt secret key is needed")
		os.Exit(1)
	}

	dbhost := os.Getenv("DB_HOST")
	if dbhost == "" {
		fmt.Println("DB Host is needed")
		os.Exit(1)
	}

	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		fmt.Println("DB port is needed")
		os.Exit(1)
	}

	dbprt, err := strconv.ParseInt(dbPort, 10, 64)
	if err != nil {
		fmt.Println("DB Port must be number")
		os.Exit(1)
	}

	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		fmt.Println("DB User is needed")
		os.Exit(1)
	}

	dbPass := os.Getenv("DB_PASS")
	if dbPass == "" {
		fmt.Println("DB Password is needed")
		os.Exit(1)
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		fmt.Println("DB Name is needed")
		os.Exit(1)
	}

	enableSSLMode := os.Getenv("ENABLE_SSL_MODE")
	if enableSSLMode == "" {
		fmt.Println("DB ENABLE_SSL_MODE is required")
		os.Exit(1)
	}

	enAblesslmode, err := strconv.ParseBool(enableSSLMode)
	if err != nil {
		fmt.Println("Invalid SSL Mode", err)
		os.Exit(1)
	}

	dbConfig := &DBConfig{
		Host:          dbhost,
		Port:          int(dbprt),
		User:          dbUser,
		Password:      dbPass,
		Name:          dbName,
		EnableSSLMode: enAblesslmode,
	}

	configurations = &Config{
		Version:      version,
		ServiceName:  serviceName,
		HttpPort:     int(port),
		JWTSecretKey: jwtSecretKey,
		DB:           dbConfig,
	}
}

func GetConfig() *Config {
	if configurations == nil {
		loadConfig()
	}
	return configurations
}
