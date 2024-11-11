package models

import (
	"os"
	"strconv"
)

type DatabaseConfig struct {
	Host         string
	Port         int
	User         string
	Password     string
	DBName       string
	PoolMaxConns int
}

func DatabaseConfigFromEnvironment() DatabaseConfig {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		port = 5432
	}
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	return DatabaseConfig{
		Host:         host,
		Port:         port,
		User:         user,
		Password:     password,
		DBName:       dbname,
		PoolMaxConns: 10,
	}
}
