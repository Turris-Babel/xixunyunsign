package web

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

type Server struct {
	Router *gin.Engine
}

func NewServer() *Server {
	r := SetupRouter()

	// 配置CORS中间件，允许所有来源
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // 允许所有来源
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	return &Server{
		Router: r,
	}
}

func (s *Server) Run(addr string) {
	log.Printf("启动Web服务器，监听地址: %s", addr)
	if err := s.Router.Run(addr); err != nil {
		log.Fatalf("Web服务器启动失败: %v", err)
	}
}
