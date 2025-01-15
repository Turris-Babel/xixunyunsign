// web/router.go
package web

import (
	"github.com/gin-gonic/gin"
)

// SetupRouter initializes the Gin router with all necessary routes
func SetupRouter() *gin.Engine {
	router := gin.Default()

	// 定义API路由
	api := router.Group("/api")
	{
		api.POST("/login", handleLogin)
		api.GET("/query", handleQuery)
		api.POST("/sign", handleSign)
		api.GET("/search_school", handleSearchSchoolID)
		// 根据需要添加更多API端点
	}

	return router
}
