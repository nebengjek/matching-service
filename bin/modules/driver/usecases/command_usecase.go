package usecases

import (
	"context"
	"encoding/json"
	"fmt"
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
	//

	marshaledData, _ := json.Marshal(tripOrder)
	log.GetLogger().Info("command_usecase", "marshaled", "kafkaProducer", utils.ConvertString(marshaledData))
	c.kafkaProducer.Publish("trip-created", marshaledData)
	result.Data = tripOrder
	return result
}
