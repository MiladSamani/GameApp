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
	"gameAppProject/service/presenceservice"
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
	Router                *echo.Echo
}

func New(config config.Config, authSvc authservice.Service, userSvc userservice.Service,
	userValidator uservalidator.Validator,
	backofficeUserSvc backofficeuserservice.Service, authorizationSvc authorizationservice.Service,
	matchingSvc matchingservice.Service,
	matchingValidator matchingvalidator.Validator,
	presenceSvc presenceservice.Service) Server {
	return Server{
		Router:                echo.New(),
		config:                config,
		userHandler:           userhandler.New(config.Auth, authSvc, userSvc, userValidator, presenceSvc),
		backofficeUserHandler: backofficeuserhandler.New(config.Auth, authSvc, backofficeUserSvc, authorizationSvc),
		matchingHandler:       matchinghandler.New(config.Auth, authSvc, matchingSvc, matchingValidator, presenceSvc),
	}
}

func (s Server) Serve() {
	// Middleware
	s.Router.Use(middleware.Logger())
	s.Router.Use(middleware.Recover())

	// Routes
	s.Router.GET("/health-check", s.healthCheck)

	s.userHandler.SetRoutes(s.Router)
	s.backofficeUserHandler.SetRoutes(s.Router)
	s.matchingHandler.SetRoutes(s.Router)

	// Start server
	address := fmt.Sprintf(":%d", s.config.HTTPServer.Port)
	fmt.Printf("start echo server on %s\n", address)
	if err := s.Router.Start(address); err != nil {
		fmt.Println("router start error", err)
	}
}
