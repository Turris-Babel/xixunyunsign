// web/router.go
package web

import (
	"github.com/gin-contrib/cors" // Add CORS import
	"github.com/gin-gonic/gin"
	"time" // Add time import for CORS MaxAge
)

// SetupRouter initializes the Gin router with all necessary routes, using the provided handlers.
// It also applies CORS middleware.
func SetupRouter(handlers *Handlers) *gin.Engine { // Accept Handlers struct
	router := gin.Default()

	// Apply CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // 允许所有来源
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 定义API路由
	api := router.Group("/api")
	{
		// Use methods from the handlers instance
		api.POST("/login", handlers.handleLogin)
		api.GET("/query", handlers.handleQuery)
		api.POST("/sign", handlers.handleSign)
		api.GET("/search_school", handlers.handleSearchSchoolID)
		// 根据需要添加更多API端点
	}

	return router
}
