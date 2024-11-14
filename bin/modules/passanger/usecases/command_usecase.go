package usecases

import (
	"context"
	"encoding/json"
	"fmt"

	driver "matching-service/bin/modules/passanger"
	"matching-service/bin/modules/passanger/models"
	kafkaPkgConfluent "matching-service/bin/pkg/kafka/confluent"
	"matching-service/bin/pkg/log"
	"matching-service/bin/pkg/utils"

	"github.com/redis/go-redis/v9"
)

type commandUsecase struct {
	driverRepositoryQuery   driver.MongodbRepositoryQuery
	driverRepositoryCommand driver.MongodbRepositoryCommand
	redisClient             redis.UniversalClient
	kafkaProducer           kafkaPkgConfluent.Producer
}

func NewCommandUsecase(mq driver.MongodbRepositoryQuery, mc driver.MongodbRepositoryCommand, rc redis.UniversalClient, kp kafkaPkgConfluent.Producer) driver.UsecaseCommand {
	return &commandUsecase{
		driverRepositoryQuery:   mq,
		driverRepositoryCommand: mc,
		redisClient:             rc,
		kafkaProducer:           kp,
	}
}

func (c *commandUsecase) BroadcastPickupPassanger(ctx context.Context, payload models.RequestRide) error {

	radius := 1.0
	drivers, err := c.redisClient.GeoRadius(ctx, "drivers-locations", payload.RouteSummary.Route.Origin.Longitude, payload.RouteSummary.Route.Origin.Latitude, &redis.GeoRadiusQuery{
		Radius:    radius,
		Unit:      "km",
		WithDist:  true,
		WithCoord: true,
		Sort:      "ASC",
	}).Result()

	if err != nil {
		log.GetLogger().Error("command_usecase", fmt.Sprintf("Error searching GeoRadius: %v", err), "BroadcastPickupPassanger", utils.ConvertString(err))
		return err
	}
	if len(drivers) > 0 {
		for _, driver := range drivers {
			geoLocation := driver
			driverMatch := models.DriverMatch{
				DriverID: geoLocation.Name,
			}

			queryRes := <-c.driverRepositoryQuery.FindDriver(ctx, driverMatch.DriverID)
			if queryRes.Error != nil {
				log.GetLogger().Error("command_usecase", fmt.Sprintf("Error searching info driver available: %v", err), "BroadcastPickupPassanger", utils.ConvertString(queryRes.Error))
				continue
			}

			dataDriver := queryRes.Data.(models.DriverAvailable)
			kafkaData := models.BroadcastPickupPassanger{
				DriverID:     driverMatch.DriverID,
				SocketID:     dataDriver.SocketID,
				RouteSummary: payload.RouteSummary,
			}
			marshaledData, _ := json.Marshal(kafkaData)
			log.GetLogger().Info("command_usecase", "marshaled", "kafkaProducer", utils.ConvertString(marshaledData))
			c.kafkaProducer.Publish("broadcast-pickup-passanger", marshaledData)
		}
	}
	return nil
}
