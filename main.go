package main

import (
	"fmt"
	"gameAppProject/config"
	"gameAppProject/delivery/httpserver"
	"gameAppProject/repository/migrator"
	"gameAppProject/repository/mysql"
	"gameAppProject/repository/mysql/mysqlaccesscontrol"
	"gameAppProject/repository/mysql/mysqluser"
	"gameAppProject/service/authorizationservice"
	"gameAppProject/service/authservice"
	"gameAppProject/service/backofficeuserservice"
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
	// TODO - merge cfg with cfg2
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

	//  TODO - add command for migrations
	mgr := migrator.New(cfg.Mysql)
	mgr.Up()

	authSvc, userSvc, userValidator, backofficeSvc, authorizationSvc := setupServices(cfg)

	server := httpserver.New(cfg, authSvc, userSvc, userValidator, backofficeSvc, authorizationSvc)

	fmt.Println("start echo server")
	server.Serve()
}

func setupServices(cfg config.Config) (
	authservice.Service, userservice.Service, uservalidator.Validator,
	backofficeuserservice.Service, authorizationservice.Service) {
	authSvc := authservice.New(cfg.Auth)

	MysqlRepo := mysql.New(cfg.Mysql)

	userMysql := mysqluser.New(MysqlRepo)
	userSvc := userservice.New(authSvc, userMysql)

	backofficeUserSvc := backofficeuserservice.New()

	aclMysql := mysqlaccesscontrol.New(MysqlRepo)
	authorizationSvc := authorizationservice.New(aclMysql)

	uV := uservalidator.New(userMysql)

	return authSvc, userSvc, uV, backofficeUserSvc, authorizationSvc
}
