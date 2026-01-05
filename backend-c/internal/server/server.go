package server

import (
	"fmt"

	"department-eduvault-backend/internal/config"
	"github.com/gin-gonic/gin"
)

type Server struct {
	engine *gin.Engine
	config *config.Config
}

// New creates a server instance from the configured router and settings.
func New(engine *gin.Engine, cfg *config.Config) *Server {
	return &Server{
		engine: engine,
		config: cfg,
	}
}

// Start runs the HTTP server on the configured port.
func (s *Server) Start() error {
	addr := fmt.Sprintf(":%s", s.config.Port)
	return s.engine.Run(addr)
}
