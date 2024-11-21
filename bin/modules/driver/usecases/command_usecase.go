package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	driver "matching-service/bin/modules/driver"
	"matching-service/bin/modules/driver/models"
	orderModel "matching-service/bin/modules/passanger/models"
	httpError "matching-service/bin/pkg/http-error"
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

func (c *commandUsecase) DriverAvailable(ctx context.Context, payload models.DriverAvailable) error {
	driver := <-c.driverRepositoryCommand.UpsertDriver(ctx, payload)
	if driver.Error != nil {
		log.GetLogger().Error("command_usecase", fmt.Sprintf("Error Failed update driver-available: %v", driver.Error), "DriverAvailable", utils.ConvertString(driver.Error))
		return httpError.InternalServerError(fmt.Sprintf("Failed update driver-available: %v", driver.Error))
	}
	return nil
}

func (c *commandUsecase) PickupPassanger(ctx context.Context, userId string, payload models.PickupPassanger) utils.Result {
	var result utils.Result
	driverInfo := <-c.driverRepositoryQuery.FindDriver(ctx, userId)
	if driverInfo.Error != nil {
		errObj := httpError.BadRequest("Profile Driver not found")
		result.Error = errObj
		return result
	}

	driver, _ := driverInfo.Data.(models.Driver)
	// get data from trip-order if still request-ride then take and update to ontheway
	trip := <-c.driverRepositoryQuery.FindOrderPassanger(ctx, payload.PassangerID)
	if trip.Error != nil {
		errObj := httpError.NewConflict()
		errObj.Message = "The order has already been taken by another driver. Please try again with a different order"
		result.Error = errObj
		log.GetLogger().Error("command_usecase", errObj.Message, "PickupPassanger", utils.ConvertString(trip.Error))
		return result
	}

	tripOrder := trip.Data.(orderModel.TripOrder)
	if tripOrder.DriverID != "" {
		errObj := httpError.NewConflict()
		errObj.Message = "The order has already been taken by another driver. Please try again with a different order"
		result.Error = errObj
		log.GetLogger().Error("command_usecase", errObj.Message, "PickupPassanger", utils.ConvertString(trip.Error))
		return result
	}
	tripOrder.DriverID = driver.Id
	tripOrder.UpdatedAt = time.Now()
	tripOrder.Status = "ontheway"

	// update to mongodb
	orderUpdate := <-c.driverRepositoryCommand.UpdateOneTripOrder(ctx, tripOrder)
	if orderUpdate.Error != nil {
		errObj := httpError.NewInternalServerError()
		errObj.Message = fmt.Sprintf("Internal server error update to db: %v", orderUpdate.Error)
		result.Error = errObj
		log.GetLogger().Error("command_usecase", fmt.Sprintf("Error update order: %v", orderUpdate.Error), "PickupPassanger", utils.ConvertString(orderUpdate.Error))
		return result
	}
	driverAvailable := <-c.driverRepositoryQuery.FindDriverAvailable(ctx, driver.Id)
	statusDriver := driverAvailable.Data.(models.StatusDriver)
	pylDriverAvailable := models.DriverAvailable{
		MetaData: models.MetaData{
			SenderID: statusDriver.SocketID,
			DriverID: driver.Id,
		},
		Available: false,
	}
	UpdateStatusDriver := <-c.driverRepositoryCommand.UpsertDriver(ctx, pylDriverAvailable)
	if UpdateStatusDriver.Error != nil {
		log.GetLogger().Error("command_usecase", fmt.Sprintf("Error update status driver: %v", UpdateStatusDriver.Error), "PickupPassanger", utils.ConvertString(UpdateStatusDriver.Error))
	}
	//

	marshaledData, _ := json.Marshal(tripOrder)
	key := fmt.Sprintf("DRIVER:PICKING-PASSANGER:%s", driver.Id)
	redisErr := c.redisClient.Set(ctx, key, marshaledData, 2*time.Hour).Err()
	if redisErr != nil {
		errObj := httpError.NewInternalServerError()
		errObj.Message = fmt.Sprintf("Internal server error insert to redis: %v", redisErr.Error())
		result.Error = errObj
		log.GetLogger().Error("command_usecase", errObj.Message, "PickupPassanger", utils.ConvertString(redisErr.Error()))
		return result
	}
	log.GetLogger().Info("command_usecase", "marshaled", "kafkaProducer", utils.ConvertString(marshaledData))
	c.kafkaProducer.Publish("trip-created", marshaledData)
	result.Data = tripOrder
	return result
}

func (c *commandUsecase) CompletedTrip(ctx context.Context, userId string, payload models.Trip) utils.Result {
	var result utils.Result
	trip := <-c.driverRepositoryQuery.FindActiveOrderPassanger(ctx, payload.OrderID)
	if trip.Error != nil {
		errObj := httpError.NewNotFound()
		errObj.Message = "The order has Notfound"
		result.Error = errObj
		log.GetLogger().Error("command_usecase", errObj.Message, "CompletedTrip", utils.ConvertString(trip.Error))
		return result
	}

	tripOrder := trip.Data.(orderModel.TripOrder)
	tripOrder.UpdatedAt = time.Now()
	tripOrder.Status = "completed"

	var tracker models.TripTracker

	key := fmt.Sprintf("trip:%s", payload.OrderID)
	driverTracker, errRedis := c.redisClient.Get(ctx, key).Result()
	if errRedis != nil || driverTracker == "" {
		errObj := httpError.NewNotFound()
		errObj.Message = fmt.Sprintf("Error get data from redis: %v", errRedis)
		result.Error = errObj
		log.GetLogger().Error("command_usecase", errObj.Message, "CompletedTrip", utils.ConvertString(errRedis))
		return result
	}
	err := json.Unmarshal([]byte(driverTracker), &tracker)
	if err != nil {
		errObj := httpError.NewInternalServerError()
		errObj.Message = fmt.Sprintf("Error unmarshal tripdata: %v", err)
		result.Error = errObj
		log.GetLogger().Error("command_usecase", errObj.Message, "CompletedTrip", utils.ConvertString(err))
		return result
	}
	// update to mongodb
	orderUpdate := <-c.driverRepositoryCommand.CompletedTripOrder(ctx, tripOrder, tracker)
	if orderUpdate.Error != nil {
		errObj := httpError.NewInternalServerError()
		errObj.Message = fmt.Sprintf("Internal server error update to db: %v", orderUpdate.Error)
		result.Error = errObj
		log.GetLogger().Error("command_usecase", fmt.Sprintf("Error update order: %v", orderUpdate.Error), "CompletedTrip", utils.ConvertString(orderUpdate.Error))
		return result
	}
	//
	realDistance, _ := strconv.ParseFloat(tracker.Data.Distance, 64)
	keyStatusDriver := fmt.Sprintf("DRIVER:PICKING-PASSANGER:%s", tripOrder.DriverID)
	c.redisClient.Del(ctx, keyStatusDriver)
	var dataEvent models.TripOrderCompleted
	dataEvent.FarePercentage = payload.FarePercentage
	dataEvent.OrderID = payload.OrderID
	dataEvent.RealDistance = realDistance
	marshaledBilling, _ := json.Marshal(dataEvent)
	marshaledData, _ := json.Marshal(tripOrder)
	c.kafkaProducer.Publish("trip-created", marshaledData)
	c.kafkaProducer.Publish("create-billing", marshaledBilling)
	result.Data = tripOrder
	return result
}
