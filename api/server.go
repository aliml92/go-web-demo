package api

import (
	"github.com/gin-gonic/gin"

	"conduit/config"
	db "conduit/db/sqlc"
)

type Server struct {
	config config.Config
	router *gin.Engine
	store  db.Querier

	// custom logger can be added here
}

func NewServer(config config.Config, store db.Querier) *Server {
	server := &Server{
		config: config,
		router: gin.Default(),
		store:  store,
	}
	return server
}

func (s *Server) MountHandlers() {
	api := s.router.Group("/api")
	api.POST("/users", s.RegisterUser)    // TODO: implement RegisterUser
	api.POST("/users/login", s.LoginUser) // TODO: implement LoginUser
}

func (s *Server) Start(addr string) error {
	return s.router.Run(addr)
}
