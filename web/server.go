package web

import (
	// "github.com/gin-contrib/cors" // Moved CORS setup
	"github.com/gin-gonic/gin"
	"log"
	// "time" // No longer needed here
)

type Server struct {
	Router *gin.Engine
}

// NewServer creates a new web server instance using the provided Gin engine.
func NewServer(engine *gin.Engine) *Server {
	// CORS middleware should ideally be applied where the router/engine is created,
	// for example, within SetupRouter or by the injector that provides the engine.
	// If we keep it here, it assumes the passed engine hasn't had CORS applied yet.
	// Let's assume for now it's applied during router setup.

	// Example if CORS needed to be applied here:
	// engine.Use(cors.New(cors.Config{
	// 	AllowOrigins:     []string{"*"},
	// 	AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	// 	AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
	// 	ExposeHeaders:    []string{"Content-Length"},
	// 	AllowCredentials: true,
	// 	MaxAge:           12 * time.Hour,
	// }))

	return &Server{
		Router: engine, // Use the provided engine
	}
}

func (s *Server) Run(addr string) {
	log.Printf("启动Web服务器，监听地址: %s", addr)
	if err := s.Router.Run(addr); err != nil {
		log.Fatalf("Web服务器启动失败: %v", err)
	}
}
