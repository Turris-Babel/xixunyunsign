// web/server.go
package web

import (
	"github.com/gin-gonic/gin"
	"log"
)

type Server struct {
	Router *gin.Engine
}

func NewServer() *Server {
	r := SetupRouter()
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
