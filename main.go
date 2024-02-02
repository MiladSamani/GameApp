package main

import (
	"fmt"
	"gameAppProject/config"
	"gameAppProject/delivery/httpserver"
	"gameAppProject/repository/mysql"
	"gameAppProject/service/authservice"
	"gameAppProject/service/userservice"
	"gameAppProject/validator/uservalidator"
	"time"
)

const (
	JwtSignKey                 = "jwt_secret"
	AccessTokenSubject         = "ac"
	RefreshTokenSubject        = "rt"
	AccessTokenExpireDuration  = time.Hour * 24
	RefreshTokenExpireDuration = time.Hour * 24 * 7
)

func main() {
	// TODO - read config path from command line
	cfg2 := config.Load("config.yml")
	fmt.Printf("cfg2: %+v\n", cfg2)

	cfg := config.Config{
		HTTPServer: config.HTTPServer{Port: 8088},
		Auth: authservice.Config{
			SignKey:               JwtSignKey,
			AccessExpirationTime:  AccessTokenExpireDuration,
			RefreshExpirationTime: RefreshTokenExpireDuration,
			AccessSubject:         AccessTokenSubject,
			RefreshSubject:        RefreshTokenSubject,
		},
		Mysql: mysql.Config{
			Username: "gameapp",
			Password: "gameapp",
			Port:     3306,
			Host:     "localhost",
			DBName:   "gameapp_db",
		},
	}

	//mgr := migrator.New(cfg.Mysql)
	//mgr.Up()

	authSvc, userSvc, userValidator := setupServices(cfg)

	server := httpserver.New(cfg, authSvc, userSvc, userValidator)

	fmt.Println("start echo server")
	server.Serve()
}

func setupServices(cfg config.Config) (authservice.Service, userservice.Service, uservalidator.Validator) {
	authSvc := authservice.New(cfg.Auth)

	MysqlRepo := mysql.New(cfg.Mysql)
	userSvc := userservice.New(authSvc, MysqlRepo)

	uV := uservalidator.New(MysqlRepo)

	return authSvc, userSvc, uV
}
