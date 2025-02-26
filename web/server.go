package web

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
)

type Server struct {
	Router *gin.Engine
}

func NewServer() *Server {
	r := SetupRouter()

	// 配置CORS中间件
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://sign.utc.edu.rs"},                 // 允许的来源
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},            // 允许的 HTTP 方法
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"}, // 允许的请求头
		ExposeHeaders:    []string{"Content-Length"},                          // 可暴露的头部信息
		AllowCredentials: true,                                                // 是否允许携带凭证
	}))
	// 你可以根据需求自定义CORS配置

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
