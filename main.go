package main

import (
	"context"
	"fmt"
	"gameAppProject/adapter/redis"
	"gameAppProject/config"
	"gameAppProject/delivery/httpserver"
	"gameAppProject/repository/migrator"
	"gameAppProject/repository/mysql"
	"gameAppProject/repository/mysql/mysqlaccesscontrol"
	"gameAppProject/repository/mysql/mysqluser"
	"gameAppProject/repository/redis/redismatching"
	"gameAppProject/repository/redis/redispresence"
	"gameAppProject/scheduler"
	"gameAppProject/service/authorizationservice"
	"gameAppProject/service/authservice"
	"gameAppProject/service/backofficeuserservice"
	"gameAppProject/service/matchingservice"
	"gameAppProject/service/presenceservice"
	"gameAppProject/service/userservice"
	"gameAppProject/validator/matchingvalidator"
	"gameAppProject/validator/uservalidator"
	"os"
	"os/signal"
	"sync"
	"time"
)

const (
	JwtSignKey = ""
)

func main() {
	// TODO - read config path from command line
	cfg := config.Load("config.yml")
	fmt.Printf("cfg: %+v\n", cfg)

	// TODO - add command for migrations
	mgr := migrator.New(cfg.Mysql)
	mgr.Up()

	// TODO - add struct and add these returned items as struct field
	authSvc, userSvc, userValidator, backofficeSvc, authorizationSvc, matchingSvc, matchingV, presenceSvc := setupServices(cfg)

	server := httpserver.New(cfg, authSvc, userSvc, userValidator, backofficeSvc, authorizationSvc,
		matchingSvc, matchingV, presenceSvc)
	go func() {
		server.Serve()
	}()

	done := make(chan bool)
	var wg sync.WaitGroup
	go func() {
		sch := scheduler.New(cfg.Scheduler, matchingSvc)

		wg.Add(1)
		sch.Start(done, &wg)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx := context.Background()
	ctxWithTimeout, cancel := context.WithTimeout(ctx, cfg.Application.GracefulShutdownTimeout)
	defer cancel()

	if err := server.Router.Shutdown(ctxWithTimeout); err != nil {
		fmt.Println("http server shutdown error", err)
	}

	fmt.Println("received interrupt signal, shutting down gracefully..")
	done <- true
	time.Sleep(cfg.Application.GracefulShutdownTimeout)

	// TODO - does order of ctx.Done & wg.Wait matter?
	<-ctxWithTimeout.Done()

	wg.Wait()
}

func setupServices(cfg config.Config) (
	authservice.Service, userservice.Service, uservalidator.Validator,
	backofficeuserservice.Service, authorizationservice.Service,
	matchingservice.Service, matchingvalidator.Validator,
	presenceservice.Service,
) {
	authSvc := authservice.New(cfg.Auth)

	MysqlRepo := mysql.New(cfg.Mysql)

	userMysql := mysqluser.New(MysqlRepo)
	userSvc := userservice.New(authSvc, userMysql)

	backofficeUserSvc := backofficeuserservice.New()

	aclMysql := mysqlaccesscontrol.New(MysqlRepo)
	authorizationSvc := authorizationservice.New(aclMysql)

	uV := uservalidator.New(userMysql)

	matchingV := matchingvalidator.New()

	redisAdapter := redis.New(cfg.Redis)

	presenceRepo := redispresence.New(redisAdapter)
	presenceSvc := presenceservice.New(cfg.PresenceService, presenceRepo)

	matchingRepo := redismatching.New(redisAdapter)
	// TODO - panic - replace presenceSvc with presence grpc client
	matchingSvc := matchingservice.New(cfg.MatchingService, matchingRepo, presenceSvc)

	return authSvc, userSvc, uV, backofficeUserSvc, authorizationSvc, matchingSvc, matchingV, presenceSvc
}
