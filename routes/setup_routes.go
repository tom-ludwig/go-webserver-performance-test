package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ApplyMiddleWare applies the middleware to the gin engine
func ApplyMiddleWare(r *gin.Engine) {
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(gin.ErrorLogger())

	// Use on Groups rather than the engine
	// r.Use(middleware.AuthMiddleware())
}

func InitilizeRoutes(r *gin.Engine, db_pool *pgxpool.Pool) {
	apiGroup := r.Group("/api")

	InitilizeUserRoutes(apiGroup, db_pool)
}
