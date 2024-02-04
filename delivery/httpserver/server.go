package httpserver

import (
	"fmt"
	"gameAppProject/config"
	"gameAppProject/delivery/httpserver/backofficeuserhandler"
	"gameAppProject/delivery/httpserver/matchinghandler"
	"gameAppProject/delivery/httpserver/userhandler"
	"gameAppProject/service/authorizationservice"
	"gameAppProject/service/authservice"
	"gameAppProject/service/backofficeuserservice"
	"gameAppProject/service/matchingservice"
	"gameAppProject/service/userservice"
	"gameAppProject/validator/matchingvalidator"
	"gameAppProject/validator/uservalidator"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	config                config.Config
	userHandler           userhandler.Handler
	backofficeUserHandler backofficeuserhandler.Handler
	matchingHandler       matchinghandler.Handler
}

func New(config config.Config, authSvc authservice.Service, userSvc userservice.Service,
	userValidator uservalidator.Validator,
	backofficeUserSvc backofficeuserservice.Service, authorizationSvc authorizationservice.Service,
	matchingSvc matchingservice.Service,
	matchingValidator matchingvalidator.Validator) Server {
	return Server{
		config:                config,
		userHandler:           userhandler.New(config.Auth, authSvc, userSvc, userValidator),
		backofficeUserHandler: backofficeuserhandler.New(config.Auth, authSvc, backofficeUserSvc, authorizationSvc),
		matchingHandler:       matchinghandler.New(config.Auth, authSvc, matchingSvc, matchingValidator),
	}
}

func (s Server) Serve() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/health-check", s.healthCheck)

	s.userHandler.SetRoutes(e)
	s.backofficeUserHandler.SetRoutes(e)
	s.matchingHandler.SetRoutes(e)

	// Start server
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", s.config.HTTPServer.Port)))
}