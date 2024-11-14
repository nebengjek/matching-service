package usecases

import (
	"context"
	"fmt"

	driver "matching-service/bin/modules/driver"
	"matching-service/bin/modules/driver/models"
	httpError "matching-service/bin/pkg/http-error"
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

func (c *commandUsecase) DriverAvailable(ctx context.Context, payload models.DriverAvailable) error {
	driver := <-c.driverRepositoryCommand.UpsertDriver(ctx, payload)
	if driver.Error != nil {
		log.GetLogger().Error("command_usecase", fmt.Sprintf("Error Failed update driver-available: %v", driver.Error), "DriverAvailable", utils.ConvertString(driver.Error))
		return httpError.InternalServerError(fmt.Sprintf("Failed update driver-available: %v", driver.Error))
	}
	return nil
}
