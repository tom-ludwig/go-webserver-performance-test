package utils

import (
	"context"
	"fmt"
	"go-webserver-performance-test/models"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Connect_to_database(config models.DatabaseConfig) (*pgxpool.Pool, error) {
	// Format the connection string
	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", config.User, config.Password, config.Host, config.Port, config.DBName)
	log.Print("Conencting to: ")
	log.Println(connectionString)

	poolConfig, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		return nil, err
	}

	poolConfig.MaxConns = int32(config.PoolMaxConns)
	poolConfig.HealthCheckPeriod = 1 * time.Minute

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, err
	}

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	log.Println("Database connection established")
	return pool, nil
}
