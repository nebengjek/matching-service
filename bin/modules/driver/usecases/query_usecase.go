package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	driver "matching-service/bin/modules/driver"
	"matching-service/bin/modules/driver/models"
	orderModel "matching-service/bin/modules/passanger/models"
	httpError "matching-service/bin/pkg/http-error"
	"matching-service/bin/pkg/log"
	"matching-service/bin/pkg/utils"
	"strconv"

	"github.com/redis/go-redis/v9"
)

type queryUsecase struct {
	driverRepositoryQuery driver.MongodbRepositoryQuery
	redisClient           redis.UniversalClient
}

func NewQueryUsecase(mq driver.MongodbRepositoryQuery, rh redis.UniversalClient) driver.UsecaseQuery {
	return &queryUsecase{
		driverRepositoryQuery: mq,
		redisClient:           rh,
	}
}

func (q *queryUsecase) DetailTrip(ctx context.Context, userId string, orderId string) utils.Result {
	var result utils.Result
	trip := <-q.driverRepositoryQuery.FindActiveOrderPassanger(ctx, orderId)
	if trip.Error != nil {
		errObj := httpError.NewNotFound()
		errObj.Message = fmt.Sprintf("The order has Notfound: %v", trip.Error)
		log.GetLogger().Error("command_usecase", "The order has Notfound", "CompletedTrip", utils.ConvertString(trip.Error))
		result.Error = errObj
		return result
	}
	tripOrder := trip.Data.(orderModel.TripOrder)
	key := fmt.Sprintf("trip:%s", orderId)
	driverTracker, errRedis := q.redisClient.Get(ctx, key).Result()
	if errRedis != nil || driverTracker == "" {
		errObj := httpError.NewNotFound()
		errObj.Message = fmt.Sprintf("Error get data from redis: %v", errRedis)
		result.Error = errObj
		log.GetLogger().Error("command_usecase", errObj.Message, "CompletedTrip", utils.ConvertString(errRedis))
		return result
	}

	var tracker models.TripTracker
	err := json.Unmarshal([]byte(driverTracker), &tracker)
	if err != nil {
		errObj := httpError.NewInternalServerError()
		errObj.Message = fmt.Sprintf("Error unmarshal tripdata: %v", err)
		result.Error = errObj
		log.GetLogger().Error("command_usecase", errObj.Message, "CompletedTrip", utils.ConvertString(err))
		return result
	}
	realDistance, _ := strconv.ParseFloat(tracker.Data.Distance, 64)
	totalFare, _, driverEarnings := CalculateFinalFare(realDistance*3000, 100)

	result.Data = models.TripResponse{
		OrderID:      orderId,
		RealDistance: realDistance,
		PassengerID:  tripOrder.PassengerID,
		Origin: models.Location{
			Latitude:  tripOrder.Origin.Latitude,
			Longitude: tripOrder.Origin.Longitude,
			Address:   tripOrder.Origin.Address,
		},
		Destination: models.Location{
			Latitude:  tripOrder.Destination.Latitude,
			Longitude: tripOrder.Destination.Longitude,
			Address:   tripOrder.Destination.Address,
		},
		Status:        tripOrder.Status,
		Price:         totalFare,
		DriverEarning: driverEarnings,
	}
	return result
}

func CalculateFinalFare(baseFare, discountPercentage float64) (totalFare, adminFee, driverEarnings float64) {
	totalFare = baseFare * (discountPercentage / 100)
	adminFee = totalFare * 0.05
	driverEarnings = totalFare - adminFee
	return
}
