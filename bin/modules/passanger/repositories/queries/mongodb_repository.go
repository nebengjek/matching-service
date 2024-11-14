package queries

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"

	driver "matching-service/bin/modules/passanger"
	"matching-service/bin/modules/passanger/models"
	"matching-service/bin/pkg/databases/mongodb"
	"matching-service/bin/pkg/utils"
)

type queryMongodbRepository struct {
	mongoDb mongodb.MongoDBLogger
}

func NewQueryMongodbRepository(mongodb mongodb.MongoDBLogger) driver.MongodbRepositoryQuery {
	return &queryMongodbRepository{
		mongoDb: mongodb,
	}
}

func (q queryMongodbRepository) FindDriver(ctx context.Context, driverId string) <-chan utils.Result {
	output := make(chan utils.Result)

	go func() {
		defer close(output)
		var driver models.DriverAvailable
		err := q.mongoDb.FindOne(mongodb.FindOne{
			Result:         &driver,
			CollectionName: "driver-available",
			Filter: bson.M{
				"driverId": driverId,
			},
		}, ctx)
		if err != nil {
			output <- utils.Result{
				Error: err,
			}
		}

		output <- utils.Result{
			Data: driver,
		}

	}()

	return output
}