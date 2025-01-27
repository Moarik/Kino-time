package api

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"kinotime/internal/configs"
	"kinotime/internal/controller"
	"kinotime/internal/model"

	"github.com/gin-gonic/gin"
)

type Server struct {
	cfg      *configs.Config
	srv      *http.Server
	userRepo *model.UserRepository
}

func NewServer(cfg *configs.Config, db *sql.DB) *Server {
	engine := gin.Default()

	server := &Server{
		cfg:      cfg,
		userRepo: model.NewUserRepository(db),
		srv: &http.Server{
			Addr:    cfg.Port,
			Handler: engine,
		},
	}

	server.setRoutes()

	return server
}

func (s *Server) Start() error {
	go func() {
		if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	<-signalChan
	return nil
}

func (s *Server) setRoutes() {
	s.srv.Handler.(*gin.Engine).POST("/login", func(c *gin.Context) {
		controller.HandleLogin(c, s.userRepo)
	})

	s.srv.Handler.(*gin.Engine).POST("/register", func(c *gin.Context) {
		controller.HandleRegister(c, s.userRepo)
	})
}

func (s *Server) Shutdown() {
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.srv.Shutdown(shutdownCtx); err != nil {
	}
}
