package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	driver "matching-service/bin/modules/passanger"
	"matching-service/bin/modules/passanger/models"
	kafkaPkgConfluent "matching-service/bin/pkg/kafka/confluent"
	"matching-service/bin/pkg/log"
	"matching-service/bin/pkg/utils"

	"github.com/redis/go-redis/v9"
	"golang.org/x/exp/rand"
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
	// creat cart order, updated if no driver take pickup
	orderData := <-c.driverRepositoryQuery.FindOrderPassanger(ctx, payload.UserId)
	if orderData.Data == nil {
		// create new
		seed := uint64(time.Now().UnixNano())
		rand.Seed(seed)
		orderID := utils.GenerateOrderID("TRIP")
		trip := models.TripOrder{
			OrderID:     orderID,
			PassengerID: payload.UserId,
			Origin: models.Location{
				Latitude:  payload.RouteSummary.Route.Origin.Latitude,
				Longitude: payload.RouteSummary.Route.Origin.Longitude,
				Address:   payload.RouteSummary.Route.Origin.Address,
			},
			Destination: models.Location{
				Latitude:  payload.RouteSummary.Route.Destination.Latitude,
				Longitude: payload.RouteSummary.Route.Destination.Longitude,
				Address:   payload.RouteSummary.Route.Destination.Address,
			},
			Status:        "request-pickup",
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
			EstimatedFare: payload.RouteSummary.BestRoutePrice,
			DistanceKm:    payload.RouteSummary.BestRouteKm,
		}
		orderCreated := <-c.driverRepositoryCommand.CreateTripOrder(ctx, trip)
		if orderCreated.Error != nil {
			log.GetLogger().Error("command_usecase", fmt.Sprintf("Error create order: %v", orderCreated.Error), "BroadcastPickupPassanger", utils.ConvertString(orderCreated.Error))
		}
	} else {
		tripOrder := orderData.Data.(models.TripOrder)
		tripOrder.Origin = models.Location{
			Latitude:  payload.RouteSummary.Route.Origin.Latitude,
			Longitude: payload.RouteSummary.Route.Origin.Longitude,
			Address:   payload.RouteSummary.Route.Origin.Address,
		}
		tripOrder.Destination = models.Location{
			Latitude:  payload.RouteSummary.Route.Destination.Latitude,
			Longitude: payload.RouteSummary.Route.Destination.Longitude,
			Address:   payload.RouteSummary.Route.Destination.Address,
		}
		tripOrder.EstimatedFare = payload.RouteSummary.BestRoutePrice
		tripOrder.DistanceKm = payload.RouteSummary.BestRouteKm
		tripOrder.CreatedAt = time.Now()
		tripOrder.UpdatedAt = time.Now()

		orderUpdate := <-c.driverRepositoryCommand.UpdateOneTripOrder(ctx, tripOrder.OrderID, tripOrder)
		if orderUpdate.Error != nil {
			log.GetLogger().Error("command_usecase", fmt.Sprintf("Error update order: %v", orderUpdate.Error), "BroadcastPickupPassanger", utils.ConvertString(orderUpdate.Error))
		}
	}
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
				PassangerID:  payload.UserId,
			}
			marshaledData, _ := json.Marshal(kafkaData)
			log.GetLogger().Info("command_usecase", "marshaled", "kafkaProducer", utils.ConvertString(marshaledData))
			c.kafkaProducer.Publish("broadcast-pickup-passanger", marshaledData)
		}
	}
	return nil
}
