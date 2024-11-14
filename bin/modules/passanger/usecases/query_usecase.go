package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	driver "matching-service/bin/modules/passanger"
	"matching-service/bin/modules/passanger/models"
	httpError "matching-service/bin/pkg/http-error"
	"matching-service/bin/pkg/log"
	"matching-service/bin/pkg/utils"

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

func (q *queryUsecase) DetailTrip(ctx context.Context, psgId string, driverId string) utils.Result {
	var result utils.Result

	key := fmt.Sprintf("USER:ROUTE:%s", psgId)
	var tripPlan models.RouteSummary
	redisData, errRedis := q.redisClient.Get(ctx, key).Result()
	if errRedis != nil || redisData == "" {
		errObj := httpError.NewNotFound()
		errObj.Message = fmt.Sprintf("Error get data from redis: %v", errRedis)
		result.Error = errObj
		log.GetLogger().Error("command_usecase", errObj.Message, "DetailTrip", utils.ConvertString(errRedis))
		return result
	}
	err := json.Unmarshal([]byte(redisData), &tripPlan)
	if err != nil {
		errObj := httpError.NewInternalServerError()
		errObj.Message = fmt.Sprintf("Error unmarshal tripdata: %v", err)
		result.Error = errObj
		log.GetLogger().Error("command_usecase", errObj.Message, "DetailTrip", utils.ConvertString(err))
		return result
	}

	result.Data = tripPlan
	return result
}
