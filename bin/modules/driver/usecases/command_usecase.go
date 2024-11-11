package usecases

import (
	"context"
	"encoding/json"
	"fmt"

	driver "matching-service/bin/modules/driver"
	"matching-service/bin/modules/driver/models"
	"matching-service/bin/pkg/log"
	"matching-service/bin/pkg/utils"

	"github.com/redis/go-redis/v9"
)

type commandUsecase struct {
	driverRepositoryQuery   driver.MongodbRepositoryQuery
	driverRepositoryCommand driver.MongodbRepositoryCommand
	redisClient             redis.UniversalClient
}

func NewCommandUsecase(mq driver.MongodbRepositoryQuery, mc driver.MongodbRepositoryCommand, rc redis.UniversalClient) driver.UsecaseCommand {
	return &commandUsecase{
		driverRepositoryQuery:   mq,
		driverRepositoryCommand: mc,
		redisClient:             rc,
	}
}

func (c *commandUsecase) BroadcastPickupPassanger(payload models.RequestRide, ctx context.Context) error {
	key := fmt.Sprintf("USER:ROUTE:%s", payload.UserId)
	var tripPlan models.RouteSummary
	redisData, errRedis := c.redisClient.Get(ctx, key).Result()
	if errRedis != nil || redisData == "" {
		log.GetLogger().Error("command_usecase", fmt.Sprintf("Error get data from redis: %v", errRedis), "BroadcastPickupPassanger", utils.ConvertString(errRedis))
		return errRedis
	}
	err := json.Unmarshal([]byte(redisData), &tripPlan)
	if err != nil {
		log.GetLogger().Error("command_usecase", fmt.Sprintf("Error unmarshal tripdata: %v", err), "FindDriver", utils.ConvertString(err))
		return err
	}
	radius := 1.0
	drivers, err := c.redisClient.GeoRadius(ctx, "drivers-locations", tripPlan.Route.Origin.Longitude, tripPlan.Route.Origin.Latitude, &redis.GeoRadiusQuery{
		Radius:    radius,
		Unit:      "km",
		WithDist:  true,
		WithCoord: true,
		Sort:      "ASC",
	}).Result()

	if err != nil {
		log.GetLogger().Error("command_usecase", fmt.Sprintf("Error searching drivers: %v", err), "FindDriver", utils.ConvertString(err))
		return err
	}
	if len(drivers) > 0 {
		// get data driver-available
		// data := {
		// 	driverId,socketId
		// }
		// loop then produce event to (send-broadcast,data)
	}
	// send event notif no driver available again
	//
	return nil
}
