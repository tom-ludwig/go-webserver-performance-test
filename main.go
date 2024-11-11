package main

import (
	"context"
	"fmt"
	"go-webserver-performance-test/models"
	"go-webserver-performance-test/routes"
	"go-webserver-performance-test/utils"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	// Load envs
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
		log.Println(err)
	}
	port := 80
	address := fmt.Sprintf(":%d", port)

	// Connect to the database
	// db_config := models.DatabaseConfig{
	// 	Host:         "localhost",
	// 	Port:         5432,
	// 	User:         "dbuser",
	// 	Password:     "password",
	// 	DBName:       "test",
	// 	PoolMaxConns: 10,
	// }

	db_config := models.DatabaseConfigFromEnvironment()

	db_pool, err := utils.Connect_to_database(db_config)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db_pool.Close()

	// Test if the connection is working Optional!!!
	_, err = db_pool.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS test_table (id SERIAL PRIMARY KEY, name TEXT)")

	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	// create a new router instance
	router := gin.Default()

	routes.InitilizeRoutes(router, db_pool)
	log.Printf("Server running on port %d", port)

	if err := router.Run(address); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
