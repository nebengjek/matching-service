package main

import (
	"context"
	"fmt"

	"matching-service/bin/pkg/log"
	"matching-service/bin/pkg/redis"
	"matching-service/bin/pkg/utils"
	"net/http"
	"os"
	"os/signal"
	"time"

	"matching-service/bin/config"

	passangerHandler "matching-service/bin/modules/passanger/handlers"
	passangerRepoCommands "matching-service/bin/modules/passanger/repositories/commands"
	passangerRepoQueries "matching-service/bin/modules/passanger/repositories/queries"
	passangerUsecase "matching-service/bin/modules/passanger/usecases"

	driverHandler "matching-service/bin/modules/driver/handlers"
	driverRepoCommands "matching-service/bin/modules/driver/repositories/commands"
	driverRepoQueries "matching-service/bin/modules/driver/repositories/queries"
	driverUsecase "matching-service/bin/modules/driver/usecases"
	kafkaConfluent "matching-service/bin/pkg/kafka/confluent"

	"matching-service/bin/pkg/apm"
	"matching-service/bin/pkg/databases/mongodb"

	"matching-service/bin/pkg/validator"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"go.elastic.co/apm/module/apmechov4"
)

func main() {
	apm.InitConnection()
	redis.LoadConfig()
	redis.InitConnection()
	mongodb.InitConnection()
	kafkaConfluent.InitKafkaConfig()
	log.Init()
	e := echo.New()
	e.Validator = &validator.CustomValidator{Validator: validator.New()}

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Skipper:          middleware.DefaultSkipper,
		Format:           `[ROUTE] ${time_rfc3339} | ${status} | ${latency_human} ${latency} | ${method} | ${uri}` + "\n",
		CustomTimeFormat: "2006-01-02 15:04:05.00000",
	}))
	e.Use(middleware.Recover())
	e.Use(apmechov4.Middleware(apmechov4.WithTracer(apm.GetTracer())))

	e.Use(middleware.CORSWithConfig(middleware.DefaultCORSConfig))
	setConfluentEvents()

	setHttp(e)

	listenerPort := fmt.Sprintf(":%s", config.GetConfig().AppPort)
	e.Logger.Fatal(e.Start(listenerPort))

	server := &http.Server{
		Addr:         listenerPort,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  5 * time.Second,
	}
	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		log.GetLogger().Info("main", "Server message-service is shutting down...", "gracefull", "")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			log.GetLogger().Info("main", fmt.Sprintf("Could not gracefully shutdown the server order-service: %v\n", err), "gracefull", "")
		}
		close(done)
	}()

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.GetLogger().Info("main", fmt.Sprintf("Could not listen on %s: %v\n", config.GetConfig().AppPort, err), "gracefull", "")
	}

	<-done
	log.GetLogger().Info("main", fmt.Sprintf("Server %s stopped", config.GetConfig().AppName), "gracefull", "")
}

func setConfluentEvents() {
	redisClient := redis.GetClient()
	kafkaProducer, err := kafkaConfluent.NewProducer(kafkaConfluent.GetConfig().GetKafkaConfig(), log.GetLogger())
	if err != nil {
		panic(err)
	}
	passangerQueryMongoRepo := passangerRepoQueries.NewQueryMongodbRepository(mongodb.NewMongoDBLogger(mongodb.GetSlaveConn(), mongodb.GetSlaveDBName(), log.GetLogger()))
	passangerCommandRepo := passangerRepoCommands.NewCommandMongodbRepository(mongodb.NewMongoDBLogger(mongodb.GetSlaveConn(), mongodb.GetSlaveDBName(), log.GetLogger()))
	passangerCommandUsecase := passangerUsecase.NewCommandUsecase(passangerQueryMongoRepo, passangerCommandRepo, redisClient, kafkaProducer)
	passangerConsumer, errPassanger := kafkaConfluent.NewConsumer(kafkaConfluent.GetConfig().GetKafkaConfig(), log.GetLogger())

	//
	driverQueryMongoRepo := driverRepoQueries.NewQueryMongodbRepository(mongodb.NewMongoDBLogger(mongodb.GetSlaveConn(), mongodb.GetSlaveDBName(), log.GetLogger()))
	driverCommandRepo := driverRepoCommands.NewCommandMongodbRepository(mongodb.NewMongoDBLogger(mongodb.GetSlaveConn(), mongodb.GetSlaveDBName(), log.GetLogger()))
	driverCommandUsecase := driverUsecase.NewCommandUsecase(driverQueryMongoRepo, driverCommandRepo, redisClient, kafkaProducer)
	driverConsumer, errDriver := kafkaConfluent.NewConsumer(kafkaConfluent.GetConfig().GetKafkaConfig(), log.GetLogger())

	passangerHandler.InitPassangerEventHandler(passangerCommandUsecase, passangerConsumer)
	driverHandler.InitPassangerEventHandler(driverCommandUsecase, driverConsumer)

	if errPassanger != nil {
		log.GetLogger().Error("main", "error registerNewConsumer", "setConfluentEvents", errPassanger.Error())
	}

	if errDriver != nil {
		log.GetLogger().Error("main", "error registerNewConsumer", "setConfluentEvents", errDriver.Error())
	}
}

func setHttp(e *echo.Echo) {
	redisClient := redis.GetClient()
	e.GET("/v1/health-check", func(c echo.Context) error {
		log.GetLogger().Info("main", "This service is running properly", "setConfluentEvents", "")
		return utils.Response(nil, "This service is running properly", 200, c)
	})
	kafkaProducer, err := kafkaConfluent.NewProducer(kafkaConfluent.GetConfig().GetKafkaConfig(), log.GetLogger())
	if err != nil {
		panic(err)
	}

	driverQueryMongodbRepo := driverRepoQueries.NewQueryMongodbRepository(mongodb.NewMongoDBLogger(mongodb.GetSlaveConn(), mongodb.GetSlaveDBName(), log.GetLogger()))
	driverCommandMongodbRepo := driverRepoCommands.NewCommandMongodbRepository(mongodb.NewMongoDBLogger(mongodb.GetMasterConn(), mongodb.GetMasterDBName(), log.GetLogger()))

	driverQueryUsecase := driverUsecase.NewQueryUsecase(driverQueryMongodbRepo, redisClient)
	driverCommandUsecase := driverUsecase.NewCommandUsecase(driverQueryMongodbRepo, driverCommandMongodbRepo, redisClient, kafkaProducer)

	driverHandler.InitDriverHttpHandler(e, driverQueryUsecase, driverCommandUsecase)
}
