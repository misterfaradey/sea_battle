package server

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ServerConf struct {
	GinMode        string
	Address        string
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	MaxHeaderBytes uint
}

type Server interface {
	Engine() *gin.Engine
	Run() error
	Shutdown(ctx context.Context) error
}

type Controller interface {
	Actions() []Action
}

type Middleware interface {
	GetMiddleware() gin.HandlerFunc
}

type Action struct {
	HttpMethod   string
	RelativePath string
	ActionExec   func(ctx *gin.Context)
}

type server struct {
	srv    *http.Server
	engine *gin.Engine
}

func NewServer(
	config *ServerConf,
	controllers ...Controller,
) *server {

	s := &server{}

	gin.SetMode(config.GinMode)

	s.engine = gin.New()
	s.engine.Use(gin.Recovery())

	if config.GinMode != gin.ReleaseMode {
		s.engine.Use(gin.Logger())
	}

	for _, controller := range controllers {
		s.setup(controller)
	}

	s.srv = &http.Server{
		Addr:           config.Address,
		Handler:        s.engine,
		ReadTimeout:    config.ReadTimeout,
		WriteTimeout:   config.WriteTimeout,
		MaxHeaderBytes: int(config.MaxHeaderBytes),
	}

	return s
}

func (s *server) setup(controller Controller) {
	for _, action := range controller.Actions() {
		a := action
		s.engine.Handle(a.HttpMethod, a.RelativePath, a.ActionExec)
	}
}

func (s *server) Engine() *gin.Engine {
	return s.engine
}

func (s *server) Run() error {
	return s.srv.ListenAndServe()
}

func (s *server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
