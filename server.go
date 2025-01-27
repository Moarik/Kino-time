package api

import (
	"context"
	"database/sql"
	"log/slog"
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
	cfg       *configs.Config
	srv       *http.Server
	userRepo  *model.UserRepository
	movieRepo *model.MovieRepository
}

func NewServer(cfg *configs.Config, db *sql.DB) *Server {
	engine := gin.Default()

	server := &Server{
		cfg:       cfg,
		userRepo:  model.NewUserRepository(db),
		movieRepo: model.NewMovieRepository(db),
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
			slog.Error("Server failed to start " + err.Error())
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	<-signalChan
	return nil
}

func (s *Server) setRoutes() {
	// POST -> /login
	s.srv.Handler.(*gin.Engine).POST("/login", func(c *gin.Context) {
		controller.HandleLogin(c, s.userRepo)
	})

	// POST -> /register
	s.srv.Handler.(*gin.Engine).POST("/register", func(c *gin.Context) {
		controller.HandleRegister(c, s.userRepo)
	})

	// private routes (JWT protected)
	private := s.srv.Handler.(*gin.Engine).Group("/private")
	private.Use(controller.JWTMiddleware())

	{
		// GET -> /private/profile
		private.GET("/profile", func(c *gin.Context) {
			controller.HandleProfile(c, s.userRepo)
		})

		// Movie routes (uses MovieHandler)
		movieHandler := controller.NewMovieHandler(s.movieRepo)

		// POST -> /private/movie
		private.POST("/movie", movieHandler.HandleCreateMovie)

		// GET -> /private/movie/:id
		private.GET("/movie/:id", movieHandler.HandleGetMovieByID)

		// GET -> /private/movies
		private.GET("/movies", movieHandler.HandleGetAllMovies)
	}
}

func (s *Server) Shutdown() {
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("Server failed to shutdown " + err.Error())
	}
}
